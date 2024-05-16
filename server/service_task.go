package server

import (
	"git.myservermanager.com/varakh/upda/api"
	redislock "github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
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
	scheduler               gocron.Scheduler
}

const (
	JobUpdatesCleanStale = "UPDATES_CLEAN_STALE"
	JobEventsCleanStale  = "EVENTS_CLEAN_STALE"
	JobActionsEnqueue    = "ACTIONS_ENQUEUE"
	JobActionsInvoke     = "ACTIONS_INVOKE"
	JobActionsCleanStale = "ACTIONS_CLEAN_STALE"
	JobPrometheusRefresh = "PROMETHEUS_REFRESH"
)

var (
	initialTasksStartDelay = time.Now().Add(10 * time.Second)
)

func newTaskService(u *updateService, e *eventService, w *webhookService, a *actionService, ai *actionInvocationService, l lockService, p *prometheusService, ac *appConfig, tc *taskConfig, lc *lockConfig, pc *prometheusConfig) *taskService {
	location, err := time.LoadLocation(ac.timeZone)

	if err != nil {
		zap.L().Sugar().Fatalf("Could not initialize correct timezone for scheduler. Reason: %s", err.Error())
	}

	// global job options
	singletonModeOption := gocron.WithSingletonMode(gocron.LimitModeReschedule)
	errorEventListener := gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
		zap.L().Sugar().Errorf("Job '%s' (%v) had a panic %v", jobName, jobID, err)
	})
	successEventListener := gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
		zap.L().Sugar().Debugf("Job '%s' (%v) finished", jobName, jobID)
	})
	eventListenerOption := gocron.WithEventListeners(successEventListener, errorEventListener)
	startAtOption := gocron.WithStartAt(gocron.WithStartDateTime(initialTasksStartDelay))

	// global scheduler options
	schedulerOptions := []gocron.SchedulerOption{gocron.WithLocation(location), gocron.WithGlobalJobOptions(singletonModeOption, eventListenerOption, startAtOption)}

	if lc.redisEnabled {
		var redisOptions *redis.Options
		redisOptions, err = redis.ParseURL(lc.redisUrl)

		if err != nil {
			zap.L().Sugar().Fatalf("Cannot parse REDIS URL '%s' to set up locking. Reason: %s", lc.redisUrl, err.Error())
		}
		redisClient := redis.NewClient(redisOptions)

		var locker gocron.Locker
		if locker, err = redislock.NewRedisLocker(redisClient, redislock.WithTries(1)); err != nil {
			zap.L().Sugar().Fatalf("Cannot set up REDIS locker. Reason: %s", err.Error())
		}

		schedulerOptions = append(schedulerOptions, gocron.WithDistributedLocker(locker))
	}

	scheduler, _ := gocron.NewScheduler(schedulerOptions...)

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
	if err := s.scheduler.StopJobs(); err != nil {
		zap.L().Sugar().Warnf("Cannot stop periodic tasks. Reason: %v", err)
	}
	if err := s.scheduler.Shutdown(); err != nil {
		zap.L().Sugar().Warnf("Cannot shut down scheduler. Reason: %v", err)
	}
	s.lockService.stop()
	zap.L().Info("Stopped all periodic tasks")
}

func (s *taskService) start() {
	s.scheduler.Start()
	zap.L().Sugar().Infof("Started %d periodic tasks", len(s.scheduler.Jobs()))
}

func (s *taskService) configureCleanupStaleUpdatesTask() {
	if !s.taskConfig.updateCleanStaleEnabled {
		return
	}

	runnable := func() {
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
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.updateCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobUpdatesCleanStale)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale updates. Reason: %s", err.Error())
	}
}

func (s *taskService) configureCleanupStaleEventsTask() {
	if !s.taskConfig.eventCleanStaleEnabled {
		return
	}

	runnable := func() {
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
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.eventCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobEventsCleanStale)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale events. Reason: %s", err.Error())
	}
}

func (s *taskService) configureActionsEnqueueTask() {
	if !s.taskConfig.actionsEnqueueEnabled {
		return
	}

	runnable := func() {
		if err := s.actionInvocationService.enqueue(s.taskConfig.actionsEnqueueBatchSize); err != nil {
			zap.L().Sugar().Errorf("Could enqueue actions. Reason: %s", err.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.actionsEnqueueInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobActionsEnqueue)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for enqueueing actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configureActionsInvokeTask() {
	if !s.taskConfig.actionsInvokeEnabled {
		return
	}

	runnable := func() {
		if err := s.actionInvocationService.invoke(s.taskConfig.actionsInvokeBatchSize, s.taskConfig.actionsInvokeMaxRetries); err != nil {
			zap.L().Sugar().Errorf("Could invoke actions. Reason: %s", err.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.actionsInvokeInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobActionsInvoke)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for invoking actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configureCleanupStaleActionsTask() {
	if !s.taskConfig.actionsCleanStaleEnabled {
		return
	}

	runnable := func() {
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
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.actionsCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobActionsCleanStale)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for cleaning stale actions. Reason: %s", err.Error())
	}
}

func (s *taskService) configurePrometheusRefreshTask() {
	if !s.prometheusConfig.enabled {
		return
	}

	runnable := func() {
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

		var webhooksTotal int64
		var webhooksError error
		webhooksTotal, webhooksError = s.webhookService.count()
		if webhooksError = s.prometheusService.setGaugeNoLabels(metricWebhooks, float64(webhooksTotal)); webhooksError != nil {
			zap.L().Sugar().Errorf("Could not refresh webhooks prometheus metric. Reason: %s", webhooksError.Error())
		}

		var eventsTotal int64
		var eventsError error
		eventsTotal, eventsError = s.eventService.count()
		if eventsError = s.prometheusService.setGaugeNoLabels(metricEvents, float64(eventsTotal)); eventsError != nil {
			zap.L().Sugar().Errorf("Could not refresh events prometheus metric. Reason: %s", eventsError.Error())
		}

		var actionsTotal int64
		var actionsError error
		actionsTotal, actionsError = s.actionService.count()
		if actionsError = s.prometheusService.setGaugeNoLabels(metricActions, float64(actionsTotal)); actionsError != nil {
			zap.L().Sugar().Errorf("Could not refresh actions prometheus metric. Reason: %s", actionsError.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.prometheusRefreshInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(JobPrometheusRefresh)); err != nil {
		zap.L().Sugar().Fatalf("Could not create task for refreshing prometheus. Reason: %s", err.Error())
	}
}
