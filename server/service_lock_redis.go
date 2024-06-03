package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	redsyncgoredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type lockRedisService struct {
	rs *redsync.Redsync
}

var (
	errLockRedisNotObtained = newServiceError(Conflict, errors.New("lock service: could not obtain lock"))
	errLockRedisNotReleased = newServiceError(Conflict, errors.New("lock service: could not release lock"))
)

func newLockRedisService(lc *lockConfig) (lockService, error) {
	zap.L().Info("Initializing REDIS locking service")

	var err error
	var redisOptions *redis.Options
	redisOptions, err = redis.ParseURL(lc.redisUrl)

	if err != nil {
		return nil, fmt.Errorf("lock service: cannot parse REDIS URL '%s' to set up locking: %s", lc.redisUrl, err)
	}

	c := redis.NewClient(redisOptions)
	if err = c.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("lock service: failed to connect to REDIS: %w", err)
	}

	pool := redsyncgoredis.NewPool(c)
	rs := redsync.New(pool)

	return &lockRedisService{rs: rs}, nil
}

// lock locks a given resource without any options
func (s *lockRedisService) lock(ctx context.Context, resource string) (appLock, error) {
	return s.lockWithOptions(ctx, resource, nil)
}

// lockWithOptions locks a given resource considering all options
func (s *lockRedisService) lockWithOptions(ctx context.Context, resource string, options ...appLockOption) (appLock, error) {
	if resource == "" {
		return nil, errorValidationNotBlank
	}

	var rsOptions []redsync.Option

	if options != nil {
		lockOptions := &appLockOptions{}
		for _, o := range options {
			o.apply(lockOptions)
		}

		if lockOptions.expiry != nil {
			rsOptions = append(rsOptions, redsync.WithExpiry(*lockOptions.expiry))
		}
		if lockOptions.maxRetries != nil {
			rsOptions = append(rsOptions, redsync.WithTries(*lockOptions.maxRetries))
		}
		if lockOptions.retryDelay != nil {
			rsOptions = append(rsOptions, redsync.WithRetryDelay(*lockOptions.retryDelay))
		}
	}

	mu := s.rs.NewMutex(resource, rsOptions...)

	zap.L().Sugar().Debugf("Trying to lock '%s'", resource)

	if err := mu.LockContext(ctx); err != nil {
		return nil, errLockRedisNotObtained
	}

	zap.L().Sugar().Debugf("Locked '%s'", resource)

	l := &redisLock{
		mu: mu,
	}

	return l, nil
}

var _ appLock = (*redisLock)(nil)

type redisLock struct {
	mu *redsync.Mutex
}

func (r redisLock) unlock(ctx context.Context) error {
	zap.L().Sugar().Debugf("Unlocking '%s'", r.mu.Name())

	unlocked, err := r.mu.UnlockContext(ctx)
	if err != nil {
		return errLockRedisNotReleased
	}
	if !unlocked {
		return errLockRedisNotReleased
	}

	return nil
}
