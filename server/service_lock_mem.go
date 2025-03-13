package server

import (
	"context"
	"errors"
	"git.myservermanager.com/varakh/upda/locker"
	"go.uber.org/zap"
	"time"
)

type lockMemService struct {
	registry *locker.InMemoryLockRegistry
}

var (
	errLockMemNotReleased = newServiceError(conflict, errors.New("lock service: could not release lock"))
)

func newLockMemService() lockService {
	zap.L().Info("Initializing in-memory locking service")
	return &lockMemService{registry: locker.NewInMemoryLockRegistry()}
}

// lock locks a given resource without any options (default expiration)
func (s *lockMemService) lock(ctx context.Context, resource string) (appLock, error) {
	return s.lockWithOptions(ctx, resource, withAppLockOptionExpiry(0))
}

// lockWithOptions locks a given resource, only TTL as option is supported
func (s *lockMemService) lockWithOptions(ctx context.Context, resource string, options ...appLockOption) (appLock, error) {
	if resource == "" {
		return nil, errorValidationNotBlank
	}

	var expiration time.Duration = 0
	if options != nil {
		lockOptions := &appLockOptions{}
		for _, o := range options {
			o.apply(lockOptions)
		}

		if lockOptions.expiry != nil {
			expiration = *lockOptions.expiry
		}
	}

	zap.L().Sugar().Debugf("Trying to lock '%s'", resource)

	s.registry.LockWithTTL(resource, expiration)

	zap.L().Sugar().Debugf("Locked '%s'", resource)

	l := &inMemoryLock{
		registry: s.registry,
		resource: resource,
	}

	return l, nil
}

var _ appLock = (*inMemoryLock)(nil)

type inMemoryLock struct {
	registry *locker.InMemoryLockRegistry
	resource string
}

func (r inMemoryLock) unlock(ctx context.Context) error {
	zap.L().Sugar().Debugf("Unlocking '%s'", r.resource)

	if err := r.registry.Unlock(r.resource); err != nil {
		return errLockMemNotReleased
	}
	return nil
}
