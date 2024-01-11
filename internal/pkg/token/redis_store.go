package token

import (
	"context"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type RStore[T string] struct {
	c redis.Cmdable

	logger *zap.Logger
}

func (rs *RStore[T]) Set(token T, expiry time.Duration) error {
	_, err := rs.c.Set(string(token), time.Now().Format(time.RFC3339), expiry).Result()
	return err
}

func (rs *RStore[T]) Exists(token T) bool {
	i, e := rs.c.Exists(string(token)).Result()
	if e != nil {
		return false
	}

	return !(i == 0)
}

func (rs *RStore[T]) Delete(token T) error {
	_, e := rs.c.Del(string(token)).Result()
	if e != nil && e != redis.Nil {
		return e
	}

	return nil
}

func (rs *RStore[T]) Shutdown(ctx context.Context) <-chan error { return nil }

func NewRStore[T string](c redis.Cmdable, logger *zap.Logger) Store[T] {
	return &RStore[T]{
		c:      c,
		logger: logger,
	}
}
