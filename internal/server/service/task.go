package service

import (
	"fmt"
	"git.myservermanager.com/varakh/upda/api"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/go-co-op/gocron-redis-lock/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

type TaskService struct {
	updateService           *UpdateService
	eventService            *EventService
	actionService           *ActionService
	actionInvocationService *ActionInvocationService
	webhookService          *WebhookService
	lockService             LockService
	prometheusService       *PrometheusService
	appConfig               *config.App
	taskConfig              *config.Task
	lockConfig              *config.Lock
	prometheusConfig        *config.Prometheus
	scheduler               gocron.Scheduler
}

const (
	jobUpdatesCleanStale = "UPDATES_CLEAN_STALE"
	jobEventsCleanStale  = "EVENTS_CLEAN_STALE"
	jobActionsEnqueue    = "ACTIONS_ENQUEUE"
	jobActionsInvoke     = "ACTIONS_INVOKE"
	jobActionsCleanStale = "ACTIONS_CLEAN_STALE"
	jobPrometheusRefresh = "PROMETHEUS_REFRESH"
)

var (
	initialTasksStartDelay = time.Now().Add(10 * time.Second)
)

func NewTaskService(u *UpdateService, e *EventService, w *WebhookService, a *ActionService, ai *ActionInvocationService, l LockService, p *PrometheusService, ac *config.App, tc *config.Task, lc *config.Lock, pc *config.Prometheus) (*TaskService, error) {
	var err error
	var location *time.Location
	if location, err = time.LoadLocation(ac.TimeZone); err != nil {
		return nil, fmt.Errorf("could not initialize correct timezone for scheduler: %s", err)
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

	if lc.RedisEnabled {
		var redisOptions *redis.Options
		redisOptions, err = redis.ParseURL(lc.RedisUrl)

		if err != nil {
			return nil, fmt.Errorf("cannot parse REDIS URL '%s' to set up locking for scheduler: %s", lc.RedisUrl, err)
		}
		redisClient := redis.NewClient(redisOptions)

		var locker gocron.Locker
		if locker, err = redislock.NewRedisLocker(redisClient, redislock.WithTries(1), redislock.WithExpiry(30*time.Second), redislock.WithRetryDelay(5*time.Second)); err != nil {
			return nil, fmt.Errorf("cannot set up REDIS locker for scheduler: %s", err)
		}

		schedulerOptions = append(schedulerOptions, gocron.WithDistributedLocker(locker))
	}

	scheduler, _ := gocron.NewScheduler(schedulerOptions...)

	return &TaskService{
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
	}, nil
}

func (s *TaskService) Init() error {
	if err := s.configureCleanupStaleUpdatesTask(); err != nil {
		return err
	}
	if err := s.configureCleanupStaleEventsTask(); err != nil {
		return err
	}
	if err := s.configureActionsEnqueueTask(); err != nil {
		return err
	}
	if err := s.configureActionsInvokeTask(); err != nil {
		return err
	}
	if err := s.configureCleanupStaleActionsTask(); err != nil {
		return err
	}
	if err := s.configurePrometheusRefreshTask(); err != nil {
		return err
	}

	return nil
}

func (s *TaskService) Stop() {
	zap.L().Sugar().Infof("Stopping %d periodic tasks...", len(s.scheduler.Jobs()))
	if err := s.scheduler.StopJobs(); err != nil {
		zap.L().Sugar().Warnf("Cannot stop periodic tasks. Reason: %v", err)
	}
	if err := s.scheduler.Shutdown(); err != nil {
		zap.L().Sugar().Warnf("Cannot shut down scheduler. Reason: %v", err)
	}
	zap.L().Info("Stopped all periodic tasks")
}

func (s *TaskService) Start() {
	s.scheduler.Start()
	zap.L().Sugar().Infof("Started %d periodic tasks", len(s.scheduler.Jobs()))
}

func (s *TaskService) configureCleanupStaleUpdatesTask() error {
	if !s.taskConfig.UpdateCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.UpdateCleanStaleMaxAge)

		var err error
		var c int64

		if c, err = s.updateService.CleanStale(t, api.UpdateStateApproved, api.UpdateStateIgnored); err != nil {
			zap.L().Sugar().Errorf("Could not clean up ignored or approved updates older than %s (%s). Reason: %s", s.taskConfig.UpdateCleanStaleMaxAge, t, err.Error())
			return
		}

		if c > 0 {
			zap.L().Sugar().Infof("Cleaned up '%d' stale updates", c)
		} else {
			zap.L().Debug("No stale updates found to clean up")
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.UpdateCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobUpdatesCleanStale)); err != nil {
		return fmt.Errorf("could not create task for cleaning stale updates: %w", err)
	}

	return nil
}

func (s *TaskService) configureCleanupStaleEventsTask() error {
	if !s.taskConfig.EventCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.EventCleanStaleMaxAge)

		var err error
		var c int64

		if c, err = s.eventService.CleanStale(t, api.EventStateCreated, api.EventStateEnqueued); err != nil {
			zap.L().Sugar().Errorf("Could not clean up stale events older than %s (%s). Reason: %s", s.taskConfig.EventCleanStaleMaxAge, t, err.Error())
			return
		}

		if c > 0 {
			zap.L().Sugar().Infof("Cleaned up '%d' stale events", c)
		} else {
			zap.L().Debug("No stale events found to clean up")
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.EventCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobEventsCleanStale)); err != nil {
		return fmt.Errorf("could not create task for cleaning stale events: %w", err)
	}

	return nil
}

func (s *TaskService) configureActionsEnqueueTask() error {
	if !s.taskConfig.ActionsEnqueueEnabled {
		return nil
	}

	runnable := func() {
		if err := s.actionInvocationService.Enqueue(s.taskConfig.ActionsEnqueueBatchSize); err != nil {
			zap.L().Sugar().Errorf("Could enqueue actions. Reason: %s", err.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.ActionsEnqueueInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobActionsEnqueue)); err != nil {
		return fmt.Errorf("could not create task for enqueueing actions: %w", err)
	}

	return nil
}

func (s *TaskService) configureActionsInvokeTask() error {
	if !s.taskConfig.ActionsInvokeEnabled {
		return nil
	}

	runnable := func() {
		if err := s.actionInvocationService.Invoke(s.taskConfig.ActionsInvokeBatchSize, s.taskConfig.ActionsInvokeMaxRetries); err != nil {
			zap.L().Sugar().Errorf("Could invoke actions. Reason: %s", err.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.ActionsInvokeInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobActionsInvoke)); err != nil {
		return fmt.Errorf("could not create task for invoking actions: %w", err)
	}

	return nil
}

func (s *TaskService) configureCleanupStaleActionsTask() error {
	if !s.taskConfig.ActionsCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.ActionsCleanStaleMaxAge)

		var cError int64
		var err error

		if cError, err = s.actionInvocationService.CleanStale(t, s.taskConfig.ActionsInvokeMaxRetries, api.ActionInvocationStateError); err != nil {
			zap.L().Sugar().Errorf("Could not clean up error stale actions older than %s (%s). Reason: %s", s.taskConfig.ActionsCleanStaleMaxAge, t, err.Error())
			return
		}

		var cSuccess int64
		if cSuccess, err = s.actionInvocationService.CleanStale(t, 0, api.ActionInvocationStateSuccess); err != nil {
			zap.L().Sugar().Errorf("Could not clean up success stale actions older than %s (%s). Reason: %s", s.taskConfig.ActionsCleanStaleMaxAge, t, err.Error())
			return
		}

		c := cError + cSuccess
		if c > 0 {
			zap.L().Sugar().Infof("Cleaned up '%d' stale actions", c)
		} else {
			zap.L().Debug("No stale actions found to clean up")
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.ActionsCleanStaleInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobActionsCleanStale)); err != nil {
		return fmt.Errorf("could not create task for cleaning stale actions: %w", err)
	}

	return nil
}

func (s *TaskService) configurePrometheusRefreshTask() error {
	if !s.prometheusConfig.Enabled {
		return nil
	}

	runnable := func() {
		updates, updatesError := s.updateService.GetAll()

		if updatesError = s.prometheusService.SetGaugeNoLabels(metricUpdatesTotal, float64(len(updates))); updatesError != nil {
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

		if updatesError = s.prometheusService.SetGaugeNoLabels(metricUpdatesPending, float64(pendingTotal)); updatesError != nil {
			zap.L().Sugar().Errorf("Could not refresh updates pending prometheus metric. Reason: %s", updatesError.Error())
		}
		if updatesError = s.prometheusService.SetGaugeNoLabels(metricUpdatesIgnored, float64(ignoredTotal)); updatesError != nil {
			zap.L().Sugar().Errorf("Could not refresh updates ignored prometheus metric. Reason: %s", updatesError.Error())
		}
		if updatesError = s.prometheusService.SetGaugeNoLabels(metricUpdatesApproved, float64(ackTotal)); updatesError != nil {
			zap.L().Sugar().Errorf("Could not refresh updates approved prometheus metric. Reason: %s", updatesError.Error())
		}

		var webhooksTotal int64
		var webhooksError error
		webhooksTotal, webhooksError = s.webhookService.Count()
		if webhooksError = s.prometheusService.SetGaugeNoLabels(metricWebhooks, float64(webhooksTotal)); webhooksError != nil {
			zap.L().Sugar().Errorf("Could not refresh webhooks prometheus metric. Reason: %s", webhooksError.Error())
		}

		var eventsTotal int64
		var eventsError error
		eventsTotal, eventsError = s.eventService.Count()
		if eventsError = s.prometheusService.SetGaugeNoLabels(metricEvents, float64(eventsTotal)); eventsError != nil {
			zap.L().Sugar().Errorf("Could not refresh events prometheus metric. Reason: %s", eventsError.Error())
		}

		var actionsTotal int64
		var actionsError error
		actionsTotal, actionsError = s.actionService.Count()
		if actionsError = s.prometheusService.SetGaugeNoLabels(metricActions, float64(actionsTotal)); actionsError != nil {
			zap.L().Sugar().Errorf("Could not refresh actions prometheus metric. Reason: %s", actionsError.Error())
		}
	}

	scheduledJob := gocron.DurationJob(s.taskConfig.PrometheusRefreshInterval)
	if _, err := s.scheduler.NewJob(scheduledJob, gocron.NewTask(runnable), gocron.WithName(jobPrometheusRefresh)); err != nil {
		return fmt.Errorf("could not create task for refreshing prometheus: %w", err)
	}

	return nil
}
