package server

import (
	"context"
	"math"
	"time"
)

// lockService provides methods for locking resources, behavior depends on underlying implementation
type lockService interface {
	// lock locks a resource applying default options (varies for implementations)
	lock(ctx context.Context, resource string) (appLock, error)

	// lockWithOptions locks a resource with given options, not all options are applied (varies for implementations)
	lockWithOptions(ctx context.Context, resource string, options ...appLockOption) (appLock, error)
}

type appLock interface {
	// unlock unlocks a lock
	unlock(ctx context.Context) error
}

type appLockOption interface {
	apply(l *appLockOptions)
}

type appLockOptionFunc func(o *appLockOptions)

func (f appLockOptionFunc) apply(o *appLockOptions) {
	f(o)
}

type appLockOptions struct {
	expiry     *time.Duration
	retryDelay *time.Duration
	maxRetries *int
}

func withAppLockOptionExpiry(expiry time.Duration) appLockOption {
	return appLockOptionFunc(func(o *appLockOptions) {
		o.expiry = &expiry
	})
}

func withAppLockOptionRetries(retries int) appLockOption {
	return appLockOptionFunc(func(o *appLockOptions) {
		o.maxRetries = &retries
	})
}

var (
	appLockOptionMaxRetries = math.MaxInt32
)

func withAppLockOptionInfiniteRetries() appLockOption {
	return appLockOptionFunc(func(o *appLockOptions) {
		o.maxRetries = &appLockOptionMaxRetries
	})
}

func withAppLockOptionRetryDelay(retryDelay time.Duration) appLockOption {
	return appLockOptionFunc(func(o *appLockOptions) {
		o.retryDelay = &retryDelay
	})
}
