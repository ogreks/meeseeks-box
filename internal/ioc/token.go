package ioc

import (
	"github.com/go-redis/redis"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"go.uber.org/zap"
)

// NewStore redis token store
func NewStore(r redis.Cmdable, logger *zap.Logger) token.Store[string] {
	return token.NewRStore(r, logger)
}
