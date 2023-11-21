package token

import (
	"go.uber.org/zap"
	"io/fs"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Store[T Type] interface {
	// Set a token to the store with the specified expiry time.
	Set(token T, expiry time.Duration) error
	// Exists checks if a token exists in the store.
	Exists(token T) bool
	// Delete a token from the store.
	Delete(token T) error
	// Shutdown the store.
	Shutdown() error
}

type DefaultStore[T string] struct {
	file fs.File

	lock *sync.RWMutex

	container   *sync.Map
	complete    chan error
	interrupter chan os.Signal // 关闭信号通知

	logger *zap.Logger
}

func (ds *DefaultStore[T]) write() error {
	return nil
}

func (ds *DefaultStore[T]) sync() error {
	return nil
}

func (ds *DefaultStore[T]) loadSyncMap() {

}

func (ds *DefaultStore[T]) init() {
	signal.Notify(ds.interrupter, os.Interrupt)

	go ds.loadSyncMap()

	go func() {
		ds.complete <- ds.sync()
	}()

	select {
	case err := <-ds.complete:
		close(ds.complete)
		_ = ds.Shutdown()
		if err != nil {
			ds.logger.Error("sync error", zap.Error(err))
			return
		}
	}

	ds.logger.Info("sync success")
}

func (ds *DefaultStore[T]) Shutdown() error {
	return nil
}

func NewDefaultStore[T string](file string, logger *zap.Logger) (*DefaultStore[T], error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	d := &DefaultStore[T]{
		file:   f,
		lock:   &sync.RWMutex{},
		logger: logger,

		container:   &sync.Map{},
		complete:    make(chan error),
		interrupter: make(chan os.Signal),
	}

	go func() {
		d.init()
	}()

	return d, nil
}
