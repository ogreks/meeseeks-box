package ioc

import (
	"github.com/ogreks/meeseeks-box/pkg/logger"
	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	accessLogger, err := logger.NewJsonLogger(
		// logger.WithDisableConsole(),
		logger.WithTimeLayout("2006-01-02 15:04:05"),
		logger.WithFilePath("./log/learn.log"),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
	}()

	return accessLogger
}
