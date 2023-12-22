package server

import (
	"git.myservermanager.com/varakh/upda/util"
	"go.uber.org/zap"
)

type lockMemService struct {
	registry *util.InMemoryLockRegistry
}

func newLockMemService() lockService {
	return &lockMemService{registry: util.NewInMemoryLockRegistry()}
}

func (s *lockMemService) init() error {
	zap.L().Info("Initialized in-memory locking service")
	return nil
}

func (s *lockMemService) tryLock(resource string) error {
	if resource == "" {
		return errorValidationNotBlank
	}

	zap.L().Sugar().Debugf("Trying to lock '%s'", resource)
	s.registry.Lock(resource)
	zap.L().Sugar().Debugf("Locked '%s'", resource)

	return nil
}

func (s *lockMemService) release(resource string) error {
	if resource == "" {
		return errorValidationNotBlank
	}

	zap.L().Sugar().Debugf("Releasing lock '%s'", resource)
	err := s.registry.Unlock(resource)
	zap.L().Sugar().Debugf("Released lock '%s'", resource)

	return err
}

func (s *lockMemService) exists(resource string) bool {
	return s.registry.Exists(resource)
}

func (s *lockMemService) stop() {
	zap.L().Info("Clearing in-memory locking service")
	s.registry.Clear()
}
