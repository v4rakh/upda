package service

import (
	"git.myservermanager.com/varakh/upda/internal/server/auth/sessiongormstore"
	"git.myservermanager.com/varakh/upda/internal/server/config"
	"github.com/go-co-op/gocron/v2"
	"github.com/rs/zerolog/log"
	"github.com/wader/gormstore/v2"
)

type AuthSessionCleanTask struct {
	config      *config.Auth
	taskService *TaskService
}

const (
	jobNameAuthSessionClean = "AUTH_SESSION_CLEAN"
)

func NewAuthSessionCleanTask(c *config.Auth, t *TaskService) *AuthSessionCleanTask {
	return &AuthSessionCleanTask{
		config:      c,
		taskService: t,
	}
}

func (t *AuthSessionCleanTask) GetCleanupFn() sessiongormstore.GormCleanupFunc {
	sessionCleanupUpFunc := func(s *gormstore.Store) {
		if !t.config.SessionCleanupEnabled {
			log.Debug().Msg("Session cleanup is disabled")
			return
		}

		runnable := func() {
			s.Cleanup()
		}

		if _, err := t.taskService.Enqueue(gocron.DurationJob(t.config.SessionCleanupInterval), gocron.NewTask(runnable), jobNameAuthSessionClean); err != nil {
			log.Error().Err(err).Msg("Could not enqueue auth session cleanup")
		}
	}

	return sessionCleanupUpFunc
}
