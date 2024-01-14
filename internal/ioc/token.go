package ioc

import (
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"go.uber.org/zap"
)

// NewStore redis token store
func NewStore(cfg configs.Config, logger *zap.Logger) token.Store[string] {
	tk, err := token.NewDefaultStore[string](cfg.RCache.Driver, logger)
	if err != nil {
		panic(err)
	}

	return tk
	//return token.NewRStore(r, logger)
}
