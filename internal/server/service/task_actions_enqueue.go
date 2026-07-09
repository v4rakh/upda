package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type ActionsEnqueueTask struct {
	actionInvocationService *ActionInvocationService
	taskService             *TaskService
	taskConfig              *config.Task
}

const (
	jobActionsEnqueue = "ACTIONS_ENQUEUE"
)

func NewActionsEnqueueTask(a *ActionInvocationService, t *TaskService, c *config.Task) *ActionsEnqueueTask {
	return &ActionsEnqueueTask{
		actionInvocationService: a,
		taskService:             t,
		taskConfig:              c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewActionsEnqueueTask
func (s *ActionsEnqueueTask) Init() error {
	return s.configureActionsEnqueueTask()
}

func (s *ActionsEnqueueTask) configureActionsEnqueueTask() error {
	if !s.taskConfig.ActionsEnqueueEnabled {
		return nil
	}

	runnable := func() {
		if err := s.actionInvocationService.Enqueue(s.taskConfig.ActionsEnqueueBatchSize); err != nil {
			log.Error().Err(err).Msg("Could not enqueue actions")
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.taskConfig.ActionsEnqueueInterval), gocron.NewTask(runnable), jobActionsEnqueue)
	return err
}
