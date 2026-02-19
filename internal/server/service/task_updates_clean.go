package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"git.myservermanager.com/varakh/upda/internal/server/constant"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"time"
)

type UpdatesCleanTask struct {
	updateService *UpdateService
	taskService   *TaskService
	taskConfig    *config.Task
}

const (
	jobUpdatesCleanStale = "UPDATES_CLEAN"
)

func NewUpdatesCleanTask(u *UpdateService, t *TaskService, c *config.Task) *UpdatesCleanTask {
	return &UpdatesCleanTask{
		updateService: u,
		taskService:   t,
		taskConfig:    c,
	}
}

// Init initializes background tasks for the service, should be called directly after NewUpdatesCleanTask
func (s *UpdatesCleanTask) Init() error {
	return s.configureUpdatesCleanTask()
}

func (s *UpdatesCleanTask) configureUpdatesCleanTask() error {
	if !s.taskConfig.UpdateCleanStaleEnabled {
		return nil
	}

	runnable := func() {
		t := time.Now()
		t = t.Add(-s.taskConfig.UpdateCleanStaleMaxAge)

		var err error
		var c int64

		cleanupStateTargets := make([]constant.UpdateState, 0)

		for _, state := range s.taskConfig.UpdateCleanStaleStates {
			var updateState constant.UpdateState

			if updateState, err = constant.ParseUpdateState(state); err != nil {
				log.Warn().Msgf("Skipping stale update cleanup for unknown state '%s'. Reason: %s", state, err.Error())
				continue
			}

			cleanupStateTargets = append(cleanupStateTargets, updateState)
		}

		if c, err = s.updateService.CleanStale(t, cleanupStateTargets...); err != nil {
			log.Error().Msgf("Could not clean up updates in state '%s' older than %s (%s). Reason: %s", cleanupStateTargets, s.taskConfig.UpdateCleanStaleMaxAge, t, err.Error())
			return
		}

		if c > 0 {
			log.Info().Msgf("Cleaned up '%d' stale updates", c)
		} else {
			log.Debug().Msg("No stale updates found to clean up")
		}
	}

	_, err := s.taskService.Enqueue(gocron.DurationJob(s.taskConfig.UpdateCleanStaleInterval), gocron.NewTask(runnable), jobUpdatesCleanStale)
	return err
}
