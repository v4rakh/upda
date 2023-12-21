package server

import (
	"git.myservermanager.com/varakh/upda/api"
	"github.com/go-co-op/gocron"
	redislock "github.com/go-co-op/gocron-redis-lock"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type taskService struct {
	updateService     *updateService
	eventService      *eventService
	webhookService    *webhookService
	prometheusService *prometheusService
	appConfig         *appConfig
	taskConfig        *taskConfig
	prometheusConfig  *prometheusConfig
	scheduler         *gocron.Scheduler
}

func newTaskService(u *updateService, e *eventService, w *webhookService, p *prometheusService, ac *appConfig, tc *taskConfig, pc *prometheusConfig) *taskService {
	location, err := time.LoadLocation(ac.timeZone)

	if err != nil {
		zap.L().Sugar().Fatalf("Could not initialize correct timezone for scheduler. Reason: %s", err.Error())
	}

	gocron.SetPanicHandler(func(jobName string, value any) {
		zap.L().Sugar().Errorf("Job '%s' had a panic %v", jobName, value)
	})

	scheduler := gocron.NewScheduler(location)

	if tc.lockRedisEnabled {
		var redisOptions *redis.Options
		redisOptions, err = redis.ParseURL(tc.lockRedisUrl)

		if err != nil {
			zap.L().Sugar().Fatalf("Cannot parse REDIS URL '%s' to set up locking. Reason: %s", tc.lockRedisUrl, err.Error())
		}
		redisClient := redis.NewClient(redisOptions)
		locker, err := redislock.NewRedisLocker(redisClient, redislock.WithTries(1))
		if err != nil {
			zap.L().Sugar().Fatalf("Cannot set up REDIS locker. Reason: %s", err.Error())
		}
		scheduler.WithDistributedLocker(locker)
	}

	return &taskService{
		updateService:     u,
		eventService:      e,
		webhookService:    w,
		prometheusService: p,
		appConfig:         ac,
		taskConfig:        tc,
		prometheusConfig:  pc,
		scheduler:         scheduler,
	}
}

func (s *taskService) init() {
	s.configureCleanupStaleUpdatesTask()
	s.configureCleanupStaleEventsTask()
	s.configurePrometheusRefreshTask()
}

func (s *taskService) stop() {
	zap.L().Sugar().Infof("Stopping %d periodic tasks...", len(s.scheduler.Jobs()))
	s.scheduler.Stop()
	zap.L().Info("Stopped all periodic tasks")
}

func (s *taskService) start() {
	s.scheduler.StartAsync()
	zap.L().Sugar().Infof("Started %d periodic tasks", len(s.scheduler.Jobs()))
}

func (s *taskService) configureCleanupStaleUpdatesTask() {
	if !s.taskConfig.updateCleanStaleEnabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.updateCleanStaleInterval).
		Do(func() {
			t := time.Now()
			t = t.Add(-s.taskConfig.updateCleanStaleMaxAge)

			var err error
			var c int64

			if c, err = s.updateService.cleanStale(t, api.UpdateStateApproved, api.UpdateStateIgnored); err != nil {
				zap.L().Sugar().Errorf("Could not clean up ignored or approved updates older than %s (%s). Reason: %s", s.taskConfig.updateCleanStaleMaxAge, t, err.Error())
				return
			}

			if c > 0 {
				zap.L().Sugar().Infof("Cleaned up '%d' stale updates", c)
			} else {
				zap.L().Info("No stale updates found to clean up")
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale updates. Reason: %s", err.Error())
	}
}

func (s *taskService) configureCleanupStaleEventsTask() {
	if !s.taskConfig.eventCleanStaleEnabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.eventCleanStaleInterval).
		Do(func() {
			t := time.Now()
			t = t.Add(-s.taskConfig.eventCleanStaleMaxAge)

			var err error
			var c int64

			if c, err = s.eventService.cleanStale(t, api.EventStateCreated); err != nil {
				zap.L().Sugar().Errorf("Could not clean up stale events older than %s (%s). Reason: %s", s.taskConfig.eventCleanStaleMaxAge, t, err.Error())
				return
			}

			if c > 0 {
				zap.L().Sugar().Infof("Cleaned up '%d' stale events", c)
			} else {
				zap.L().Info("No stale events found to clean up")
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale events. Reason: %s", err.Error())
	}
}

func (s *taskService) configurePrometheusRefreshTask() {
	if !s.prometheusConfig.enabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.prometheusRefreshInterval).
		Do(func() {
			// updates with labels and collect stats about state
			updates, updatesError := s.updateService.getAll()

			if updatesError = s.prometheusService.setGaugeNoLabels(metricUpdatesTotal, float64(len(updates))); updatesError != nil {
				zap.L().Sugar().Errorf("Could not refresh updates all prometheus metric. Reason: %s", updatesError.Error())
			}

			var pendingTotal int64
			var ignoredTotal int64
			var ackTotal int64

			for _, update := range updates {
				var updateState float64
				if api.UpdateStatePending.Value() == update.State {
					pendingTotal += 1
					updateState = 0
				} else if api.UpdateStateIgnored.Value() == update.State {
					ignoredTotal += 1
					updateState = 2
				} else if api.UpdateStateApproved.Value() == update.State {
					ackTotal += 1
					updateState = 1
				}

				if updatesError = s.prometheusService.setGauge(metricUpdates, []string{update.Application, update.Provider, update.Host}, updateState); updatesError != nil {
					zap.L().Sugar().Errorf("Could not refresh updates prometheus metric. Reason: %s", updatesError.Error())
				}
			}

			if updatesError = s.prometheusService.setGaugeNoLabels(metricUpdatesPending, float64(pendingTotal)); updatesError != nil {
				zap.L().Sugar().Errorf("Could not refresh updates pending prometheus metric. Reason: %s", updatesError.Error())
			}
			if updatesError = s.prometheusService.setGaugeNoLabels(metricUpdatesIgnored, float64(ignoredTotal)); updatesError != nil {
				zap.L().Sugar().Errorf("Could not refresh updates ignored prometheus metric. Reason: %s", updatesError.Error())
			}
			if updatesError = s.prometheusService.setGaugeNoLabels(metricUpdatesApproved, float64(ackTotal)); updatesError != nil {
				zap.L().Sugar().Errorf("Could not refresh updates approved prometheus metric. Reason: %s", updatesError.Error())
			}

			// webhooks
			var webhooksTotal int64
			var webhooksError error
			webhooksTotal, webhooksError = s.webhookService.count()
			if webhooksError = s.prometheusService.setGaugeNoLabels(metricWebhooks, float64(webhooksTotal)); webhooksError != nil {
				zap.L().Sugar().Errorf("Could not refresh webhooks prometheus metric. Reason: %s", webhooksError.Error())
			}

			// events
			var eventsTotal int64
			var eventsError error
			eventsTotal, eventsError = s.eventService.count()
			if eventsError = s.prometheusService.setGaugeNoLabels(metricEvents, float64(eventsTotal)); eventsError != nil {
				zap.L().Sugar().Errorf("Could not refresh events prometheus metric. Reason: %s", eventsError.Error())
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for refreshing prometheus. Reason: %s", err.Error())
	}
}
