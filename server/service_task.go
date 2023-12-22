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
	lockService       lockService
	prometheusService *prometheusService
	appConfig         *appConfig
	taskConfig        *taskConfig
	lockConfig        *lockConfig
	prometheusConfig  *prometheusConfig
	scheduler         *gocron.Scheduler
}

const (
	taskLockNameUpdatesCleanStale = "updates_clean_stale"
	taskLockNameEventsCleanStale  = "events_clean_stale"
	taskLockNamePrometheusUpdate  = "prometheus_update"
)

func newTaskService(u *updateService, e *eventService, w *webhookService, l lockService, p *prometheusService, ac *appConfig, tc *taskConfig, lc *lockConfig, pc *prometheusConfig) *taskService {
	location, err := time.LoadLocation(ac.timeZone)

	if err != nil {
		zap.L().Sugar().Fatalf("Could not initialize correct timezone for scheduler. Reason: %s", err.Error())
	}

	gocron.SetPanicHandler(func(jobName string, value any) {
		zap.L().Sugar().Errorf("Job '%s' had a panic %v", jobName, value)
	})

	scheduler := gocron.NewScheduler(location)

	if lc.redisEnabled {
		var redisOptions *redis.Options
		redisOptions, err = redis.ParseURL(lc.redisUrl)

		if err != nil {
			zap.L().Sugar().Fatalf("Cannot parse REDIS URL '%s' to set up locking. Reason: %s", lc.redisUrl, err.Error())
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
		lockService:       l,
		prometheusService: p,
		appConfig:         ac,
		taskConfig:        tc,
		lockConfig:        lc,
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
	s.lockService.stop()
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
	initialDelay := time.Now().Add(10 * time.Second)
	_, err := s.scheduler.Every(s.taskConfig.updateCleanStaleInterval).
		StartAt(initialDelay).
		Do(func() {
			resource := taskLockNameUpdatesCleanStale
			// distributed lock handled via gocron-redis-lock for tasks
			if !s.lockConfig.redisEnabled {
				// skip execution if lock already exists, wait otherwise
				if lockExists := s.lockService.exists(resource); lockExists {
					zap.L().Sugar().Debugf("Skipping task execution because task lock '%s' exists", resource)
					return
				}
				_ = s.lockService.tryLock(resource)
				defer func(lockService lockService, resource string) {
					err := lockService.release(resource)
					if err != nil {
						zap.L().Sugar().Warnf("Could not release task lock '%s'", resource)
					}
				}(s.lockService, resource)
			}

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

	initialDelay := time.Now().Add(5 * time.Second)
	_, err := s.scheduler.Every(s.taskConfig.eventCleanStaleInterval).
		StartAt(initialDelay).
		Do(func() {
			resource := taskLockNameEventsCleanStale
			// distributed lock handled via gocron-redis-lock for tasks
			if !s.lockConfig.redisEnabled {
				// skip execution if lock already exists, wait otherwise
				if lockExists := s.lockService.exists(resource); lockExists {
					zap.L().Sugar().Debugf("Skipping task execution because task lock '%s' exists", resource)
					return
				}
				_ = s.lockService.tryLock(resource)
				defer func(lockService lockService, resource string) {
					err := lockService.release(resource)
					if err != nil {
						zap.L().Sugar().Warnf("Could not release task lock '%s'", resource)
					}
				}(s.lockService, resource)
			}

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

	initialDelay := time.Now().Add(10 * time.Second)
	_, err := s.scheduler.Every(s.taskConfig.prometheusRefreshInterval).
		StartAt(initialDelay).
		Do(func() {
			resource := taskLockNamePrometheusUpdate
			// distributed lock handled via gocron-redis-lock for tasks
			if !s.lockConfig.redisEnabled {
				// skip execution if lock already exists, wait otherwise
				if lockExists := s.lockService.exists(resource); lockExists {
					zap.L().Sugar().Debugf("Skipping task execution because task lock '%s' exists", resource)
					return
				}
				_ = s.lockService.tryLock(resource)
				defer func(lockService lockService, resource string) {
					err := lockService.release(resource)
					if err != nil {
						zap.L().Sugar().Warnf("Could not release task lock '%s'", resource)
					}
				}(s.lockService, resource)
			}

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
