package locker

import (
	"context"

	"github.com/bsm/redislock"
)

type Lock interface {
	Release(ctx context.Context) error
}

type lock struct {
	redisLock *redislock.Lock
}

func (l *lock) Release(ctx context.Context) error {
	return l.redisLock.Release(ctx)
}
