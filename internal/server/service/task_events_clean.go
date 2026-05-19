package service

import (
	"time"

	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
)

type EventsCleanTask struct {
	eventService *EventService
	taskService  *TaskService
	taskConfig   *config.Task
}

const (
	jobEventsCleanStale = "EVENTS_CLEAN"
)

func NewEventsCleanTask(e *EventService, t *TaskService, c *config.Task) *EventsCleanTask {
	return &EventsCleanTask{
		eventService: e,
		taskService:  t,
		taskConfig:   c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewEventsCleanTask
func (s *EventsCleanTask) Init() error {
	return s.configureEventsCleanTask()
}

func (s *EventsCleanTask) configureEventsCleanTask() error {
	if !s.taskConfig.EventCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.EventCleanStaleMaxAge)

		var err error
		var c int64

		if c, err = s.eventService.CleanStale(t, constant.EventStateCreated, constant.EventStateEnqueued); err != nil {
			log.Error().Msgf("Could not clean up stale events older than %s (%s). Reason: %s", s.taskConfig.EventCleanStaleMaxAge, t, err.Error())
			return
		}

		if c > 0 {
			log.Info().Msgf("Cleaned up '%d' stale events", c)
		} else {
			log.Debug().Msg("No stale events found to clean up")
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.taskConfig.EventCleanStaleInterval), gocron.NewTask(runnable), jobEventsCleanStale)
	return err
}
