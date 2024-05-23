package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	testLockName = "test_lock"
)

func TestLockExpiresAndCannotBeReleased(t *testing.T) {
	a := assert.New(t)

	s := newLockMemService()
	ctx := context.Background()

	lock, lockErr := s.lockWithOptions(ctx, testLockName, withAppLockOptionExpiry(250*time.Millisecond))
	a.Nil(lockErr)
	a.NotNil(lock)

	time.Sleep(251 * time.Millisecond)

	unlockErr := lock.unlock(ctx)
	a.NotNil(unlockErr)
	a.ErrorContains(unlockErr, "could not release lock")
}
