package ioc

import (
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/pkg/logger"
	"github.com/ogreks/meeseeks-box/pkg/logger/driver"
	"github.com/ogreks/meeseeks-box/pkg/logger/driver/local"
	"go.uber.org/zap"
)

func InitLogger(cfg configs.Config, driver driver.Driver) *zap.Logger {
	accessLogger, err := logger.NewJsonLogger(
		logger.WithDisableConsole(),
		logger.WithTimeLayout("2006-01-02 15:04:05"),
		logger.WithDriver(driver),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()

	return accessLogger
}

func InitLogDriver(cfg configs.Config) driver.Driver {
	return local.NewLocalDriver(cfg.Server.LogPath)
}
