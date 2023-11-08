package utils

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, logger *zap.Logger, handle func(ctx context.Context) (func(), error)) error {
	state := 1
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	f, err := handle(ctx)
	if err != nil {
		return err
	}

	defer func() {
		f()
		logger.Info("Server exit, bye...")
		time.Sleep(time.Millisecond * 100)
		os.Exit(state)
	}()

	for {
		notify := <-c
		logger.Info("signal notify", zap.String("signal", notify.String()))

		switch notify {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			return nil
		case syscall.SIGHUP:
		default:
			return nil
		}
	}
}
