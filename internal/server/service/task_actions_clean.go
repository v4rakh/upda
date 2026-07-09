package service

import (
	"time"

	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type ActionsCleanTask struct {
	actionInvocationService *ActionInvocationService
	taskService             *TaskService
	taskConfig              *config.Task
}

const (
	jobActionsCleanStale = "ACTIONS_CLEAN"
)

func NewActionsCleanTask(a *ActionInvocationService, t *TaskService, c *config.Task) *ActionsCleanTask {
	return &ActionsCleanTask{
		actionInvocationService: a,
		taskService:             t,
		taskConfig:              c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewActionsCleanTask
func (s *ActionsCleanTask) Init() error {
	return s.configureActionsCleanTask()
}

func (s *ActionsCleanTask) configureActionsCleanTask() error {
	if !s.taskConfig.ActionsCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.ActionsCleanStaleMaxAge)

		var cError int64
		var err error

		if cError, err = s.actionInvocationService.CleanStale(t, s.taskConfig.ActionsInvokeMaxRetries, constant.ActionInvocationStateError); err != nil {
			log.Error().Err(err).Msgf("Could not clean up error stale actions older than %s (%s)", s.taskConfig.ActionsCleanStaleMaxAge, t)
			return
		}

		var cSuccess int64
		if cSuccess, err = s.actionInvocationService.CleanStale(t, 0, constant.ActionInvocationStateSuccess); err != nil {
			log.Error().Err(err).Msgf("Could not clean up success stale actions older than %s (%s)", s.taskConfig.ActionsCleanStaleMaxAge, t)
			return
		}

		c := cError + cSuccess
		if c > 0 {
			log.Info().Msgf("Cleaned up '%d' stale actions", c)
		} else {
			log.Debug().Msg("No stale actions found to clean up")
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.taskConfig.ActionsCleanStaleInterval), gocron.NewTask(runnable), jobActionsCleanStale)
	return err
}
