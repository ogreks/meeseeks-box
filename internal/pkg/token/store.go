package token

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/signal"
	"regexp"
	"sync"
	"time"

	"go.uber.org/zap"
)

//go:generate mockgen -source=./store.go -package=tkmocks -destination=mocks/store.mock.go Store
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

type DefaultStore[T Type] struct {
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

	var (
		m      = make(map[T]time.Time)
		delVal []T
	)

	ds.container.Range(func(key, value any) bool {
		tv := value.(time.Time)
		if tv.Before(time.Now()) || tv.Equal(time.Now()) {
			delVal = append(delVal, key.(T))
			return true
		}
		m[key.(T)] = value.(time.Time)
		return true
	})

	// delete timeout key
	for _, v := range delVal {
		ds.container.Delete(v)
	}

	_ = ds.file.Truncate(0)
	if len(m) == 0 {
		return nil
	}

	data, err := json.MarshalIndent(&m, "", "    ")
	if err != nil {
		return err
	}

	bf := bufio.NewWriter(ds.file)
	_, err = bf.Write(bytes.TrimSpace(data))
	if err != nil {
		return err
	}

	return bf.Flush()
}

func (ds *DefaultStore[T]) sync() error {
	t := time.NewTimer(60 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err := ds.write()
			if err != nil {
				ds.logger.Error("write file error", zap.Error(err))
				return err
			}
		case <-ds.interrupter:
			signal.Stop(ds.interrupter)
			return errors.New("interrupt closed")
		}
	}
}

func (ds *DefaultStore[T]) loadSyncMap() {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	var buffer = make([]byte, 1024)
	var ret []byte

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

		tmb := bytes.TrimSpace(buffer[:readLength])
		if len(tmb) == 0 {
			continue
		}

		ret = append(ret, buffer[:readLength]...)
	}

	if len(ret) == 0 {
		return
	}

	invalidCharRegex := regexp.MustCompile("[\x00-\x1F]")
	ret = invalidCharRegex.ReplaceAll(ret, []byte{})

	var m = make(map[T]time.Time)
	err := json.Unmarshal(ret, &m)
	if err != nil {
		panic(err)
		ds.logger.Panic("unmarshal error", zap.Error(err))
		return
	}

	for k, v := range m {
		if v.Before(time.Now()) {
			continue
		}
		ds.container.Store(k, v)
	}

	ds.container.Range(func(key, value any) bool {
		return true
	})
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
	if err != nil && !errors.Is(err, errors.New("interrupt closed")) {
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

func NewDefaultStore[T Type](file string, logger *zap.Logger) (Store[T], error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
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

func (ds *DefaultStore[T]) Set(token T, expiry time.Duration) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ds.container.Store(token, time.Now().Add(expiry))

	return nil
}

func (ds *DefaultStore[T]) Exists(token T) bool {
	value, ok := ds.container.Load(token)

	if !ok {
		return false
	}

	if value.(time.Time).Before(time.Now()) {
		ds.container.Delete(token)
		return false
	}

	return true
}

func (ds *DefaultStore[T]) Delete(token T) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	ds.container.Delete(token)

	return nil
}
