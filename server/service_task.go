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
	updateService           *updateService
	eventService            *eventService
	actionService           *actionService
	actionInvocationService *actionInvocationService
	webhookService          *webhookService
	lockService             lockService
	prometheusService       *prometheusService
	appConfig               *appConfig
	taskConfig              *taskConfig
	lockConfig              *lockConfig
	prometheusConfig        *prometheusConfig
	scheduler               *gocron.Scheduler
}

const (
	taskLockNameUpdatesCleanStale = "updates_clean_stale"
	taskLockNameEventsCleanStale  = "events_clean_stale"
	taskLockNameActionsEnqueue    = "actions_enqueue"
	taskLockNameActionsInvoke     = "actions_invoke"
	taskLockNameActionsCleanStale = "actions_clean_stale"
	taskLockNamePrometheusUpdate  = "prometheus_update"
)

var (
	initialTasksStartDelay = time.Now().Add(10 * time.Second)
)

func newTaskService(u *updateService, e *eventService, w *webhookService, a *actionService, ai *actionInvocationService, l lockService, p *prometheusService, ac *appConfig, tc *taskConfig, lc *lockConfig, pc *prometheusConfig) *taskService {
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
		updateService:           u,
		eventService:            e,
		actionService:           a,
		actionInvocationService: ai,
		webhookService:          w,
		lockService:             l,
		prometheusService:       p,
		appConfig:               ac,
		taskConfig:              tc,
		lockConfig:              lc,
		prometheusConfig:        pc,
		scheduler:               scheduler,
	}
}

func (s *taskService) init() {
	s.configureCleanupStaleUpdatesTask()
	s.configureCleanupStaleEventsTask()
	s.configureActionsEnqueueTask()
	s.configureActionsInvokeTask()
	s.configureCleanupStaleActionsTask()
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
	_, err := s.scheduler.Every(s.taskConfig.updateCleanStaleInterval).
		StartAt(initialTasksStartDelay).
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
				zap.L().Debug("No stale updates found to clean up")
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
		StartAt(initialTasksStartDelay).
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

			if c, err = s.eventService.cleanStale(t, api.EventStateCreated, api.EventStateEnqueued); err != nil {
				zap.L().Sugar().Errorf("Could not clean up stale events older than %s (%s). Reason: %s", s.taskConfig.eventCleanStaleMaxAge, t, err.Error())
				return
			}

			if c > 0 {
				zap.L().Sugar().Infof("Cleaned up '%d' stale events", c)
			} else {
				zap.L().Debug("No stale events found to clean up")
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale events. Reason: %s", err.Error())
	}
}

func (s *taskService) configureActionsEnqueueTask() {
	if !s.taskConfig.actionsEnqueueEnabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.actionsEnqueueInterval).
		StartAt(initialTasksStartDelay).
		Do(func() {
			resource := taskLockNameActionsEnqueue
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

			if err := s.actionInvocationService.enqueue(s.taskConfig.actionsEnqueueBatchSize); err != nil {
				zap.L().Sugar().Errorf("Could enqueue actions. Reason: %s", err.Error())
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for enqueueing actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configureActionsInvokeTask() {
	if !s.taskConfig.actionsInvokeEnabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.actionsInvokeInterval).
		StartAt(initialTasksStartDelay).
		Do(func() {
			resource := taskLockNameActionsInvoke
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

			if err := s.actionInvocationService.invoke(s.taskConfig.actionsInvokeBatchSize, s.taskConfig.actionsInvokeMaxRetries); err != nil {
				zap.L().Sugar().Errorf("Could invoke actions. Reason: %s", err.Error())
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for invoking actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configureCleanupStaleActionsTask() {
	if !s.taskConfig.actionsCleanStaleEnabled {
		return
	}
	_, err := s.scheduler.Every(s.taskConfig.actionsCleanStaleInterval).
		StartAt(initialTasksStartDelay).
		Do(func() {
			resource := taskLockNameActionsCleanStale
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
			t = t.Add(-s.taskConfig.actionsCleanStaleMaxAge)

			var cError int64
			var err error

			if cError, err = s.actionInvocationService.cleanStale(t, s.taskConfig.actionsInvokeMaxRetries, api.ActionInvocationStateError); err != nil {
				zap.L().Sugar().Errorf("Could not clean up error stale actions older than %s (%s). Reason: %s", s.taskConfig.actionsCleanStaleMaxAge, t, err.Error())
				return
			}

			var cSuccess int64
			if cSuccess, err = s.actionInvocationService.cleanStale(t, 0, api.ActionInvocationStateSuccess); err != nil {
				zap.L().Sugar().Errorf("Could not clean up success stale actions older than %s (%s). Reason: %s", s.taskConfig.actionsCleanStaleMaxAge, t, err.Error())
				return
			}

			c := cError + cSuccess
			if c > 0 {
				zap.L().Sugar().Infof("Cleaned up '%d' stale actions", c)
			} else {
				zap.L().Debug("No stale actions found to clean up")
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configurePrometheusRefreshTask() {
	if !s.prometheusConfig.enabled {
		return
	}

	_, err := s.scheduler.Every(s.taskConfig.prometheusRefreshInterval).
		StartAt(initialTasksStartDelay).
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
				if api.UpdateStatePending.Value() == update.State {
					pendingTotal += 1
				} else if api.UpdateStateIgnored.Value() == update.State {
					ignoredTotal += 1
				} else if api.UpdateStateApproved.Value() == update.State {
					ackTotal += 1
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

			// actions
			var actionsTotal int64
			var actionsError error
			actionsTotal, actionsError = s.actionService.count()
			if actionsError = s.prometheusService.setGaugeNoLabels(metricActions, float64(actionsTotal)); actionsError != nil {
				zap.L().Sugar().Errorf("Could not refresh actions prometheus metric. Reason: %s", actionsError.Error())
			}
		})

	if err != nil {
		zap.L().Sugar().Fatalf("Could not create task for refreshing prometheus. Reason: %s", err.Error())
	}
}
