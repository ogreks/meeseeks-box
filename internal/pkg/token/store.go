package token

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Store[T Type] interface {
	// Set a token to the store with the specified expiry time.
	Set(token T, expiry time.Duration) error
	// Exists checks if a token exists in the store.
	Exists(token T) bool
	// Delete a token from the store.
	Delete(token T) error
	// Shutdown the store.
	Shutdown(ctx context.Context) <-chan error
}

type DefaultStore[T string] struct {
	file *os.File

	lock *sync.RWMutex

	container   *sync.Map
	complete    chan error
	interrupter chan os.Signal // 关闭信号通知

	logger *zap.Logger
}

func (ds *DefaultStore[T]) write() error {
	ds.lock.Lock()
	defer func() {
		_ = ds.file.Sync()
		ds.lock.Unlock()
	}()

	var m = make(map[T]time.Time)
	ds.container.Range(func(key, value any) bool {
		m[key.(T)] = value.(time.Time)
		return true
	})

	data, err := json.Marshal(&m)
	if err != nil {
		return err
	}

	var b = bufio.NewWriter(ds.file)
	defer b.Flush()
	_, err = b.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (ds *DefaultStore[T]) sync() error {
	for {
		select {
		case <-time.After(5 * time.Second):
			err := ds.write()
			if err != nil {
				ds.logger.Error("write file error", zap.Error(err))
				return err
			}
		case <-ds.interrupter:
			signal.Stop(ds.interrupter)
			ds.complete <- errors.New("interrupt closed")
			return nil
		}
	}
}

func (ds *DefaultStore[T]) loadSyncMap() {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	var buffer = make([]byte, 1024)
	ret := []byte{}

	for {
		readLength, err := ds.file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				ret = append(ret, buffer[:readLength]...)
				break
			}
			ds.logger.Panic("read file error", zap.Error(err))
			return
		}

		ret = append(ret, buffer[:readLength]...)
	}

	if len(ret) == 0 {
		return
	}

	var m = make(map[T]time.Time)
	err := json.Unmarshal(ret, &m)
	if err != nil {
		ds.logger.Panic("unmarshal error", zap.Error(err))
		return
	}

	for k, v := range m {
		ds.container.Store(k, v)
	}
}

func (ds *DefaultStore[T]) init() {
	signal.Notify(ds.interrupter, os.Interrupt)

	go ds.loadSyncMap()

	go func() {
		ds.complete <- ds.sync()
	}()

	err := <-ds.complete
	close(ds.complete)
	_ = ds.shutdown(context.Background())
	if err != nil {
		if !errors.Is(err, errors.New("interrupt closed")) {
			ds.logger.Error("sync error", zap.Error(err))
		}

		return
	}

	ds.logger.Info("sync success")
}

func (ds *DefaultStore[T]) Shutdown(ctx context.Context) <-chan error {
	ch := make(chan error)

	go func(ch chan<- error) {
		ch <- ds.shutdown(ctx)
	}(ch)

	return ch
}

func (ds *DefaultStore[T]) shutdown(ctx context.Context) (err error) {
	defer func() {
		ctx.Done()
		err = ds.file.Close()
	}()

	err = ds.write()
	return
}

func NewDefaultStore[T string](file string, logger *zap.Logger) (*DefaultStore[T], error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_SYNC|os.O_TRUNC, 0666)
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

func (d *DefaultStore[T]) Set(token T, expiry time.Duration) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.container.Store(token, time.Now().Add(expiry))

	return nil
}

func (d *DefaultStore[T]) Exists(token T) bool {
	value, ok := d.container.Load(token)

	if !ok {
		return false
	}

	if value.(time.Time).Before(time.Now()) {
		d.container.Delete(token)
		return false
	}

	return true
}

func (d *DefaultStore[T]) Delete(token T) error {
	d.lock.Lock()
	defer d.lock.Unlock()

	d.container.Delete(token)

	return nil
}
