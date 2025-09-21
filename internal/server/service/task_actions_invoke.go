package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type ActionsInvokeTask struct {
	actionInvocationService *ActionInvocationService
	taskService             *TaskService
	taskConfig              *config.Task
}

const (
	jobActionsInvoke = "ACTIONS_INVOKE"
)

func NewActionsInvokeTask(a *ActionInvocationService, t *TaskService, c *config.Task) *ActionsInvokeTask {
	return &ActionsInvokeTask{
		actionInvocationService: a,
		taskService:             t,
		taskConfig:              c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewActionsInvokeTask
func (s *ActionsInvokeTask) Init() error {
	return s.configureActionsInvokeTask()
}

func (s *ActionsInvokeTask) configureActionsInvokeTask() error {
	if !s.taskConfig.ActionsInvokeEnabled {
		return nil
	}

	runnable := func() {
		if err := s.actionInvocationService.Invoke(s.taskConfig.ActionsInvokeBatchSize, s.taskConfig.ActionsInvokeMaxRetries); err != nil {
			log.Error().Msgf("Could invoke actions. Reason: %s", err.Error())
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.taskConfig.ActionsInvokeInterval), gocron.NewTask(runnable), jobActionsInvoke)
	return err
}
