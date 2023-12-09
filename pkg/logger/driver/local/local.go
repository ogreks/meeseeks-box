package local

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/ogreks/meeseeks-box/pkg/logger/driver"
	"gopkg.in/natefinch/lumberjack.v2"
)

type BackType int

const (
	BDay BackType = iota + 1
	BHour
	BMinitue
)

var _ driver.Driver = (*LocalDriver)(nil)

type LocalDriver struct {
	lumberjack.Logger

	backType BackType
}

func WithCompress(comparable bool) driver.Option {
	return func(d driver.Driver) {
		d.(*LocalDriver).Compress = comparable
	}
}

func NewLocalDriver(file string, backType ...BackType) *LocalDriver {
	if len(backType) == 0 {
		backType = []BackType{BDay}
	}
	dir := filepath.Dir(file)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	d := &LocalDriver{
		backType: backType[0],
		Logger: lumberjack.Logger{
			Filename:   file,
			MaxSize:    128,
			MaxBackups: 300,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  true,
		},
	}

	go d.back()

	return d
}

func (d *LocalDriver) caleTickerTime() (t time.Duration) {
	var td time.Time
	ts := time.Now()
	switch d.backType {
	case BDay:
		td = time.Date(ts.Year(), ts.Month(), ts.Day()+1, 0, 0, 0, 0, ts.Location())
	case BHour:
		td = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour()+1, 0, 0, 0, ts.Location())
	case BMinitue:
		td = time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute()+1, 0, 0, ts.Location())
	default:
		td = time.Date(ts.Year(), ts.Month(), ts.Day()+1, 0, 0, 0, 0, ts.Location())
	}

	return td.Sub(ts)
}

func (d *LocalDriver) back() {
	t := time.NewTicker(d.caleTickerTime())
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	for {
		select {
		case <-t.C:
			d.Rotate()
		case <-c:
			close(c)
			t.Stop()
			return
		}
	}
}
