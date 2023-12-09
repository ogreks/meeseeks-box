package local

import (
	"os"
	"testing"
	"time"
)

//go:test -timeout 300s -run ^Test_LocalDriver$ github.com/ogreks/meeseeks-box/pkg/logger/driver/local
func Test_LocalDriver(t *testing.T) {
	d := NewLocalDriver("test.log", BMinitue)

	defer func() {
		os.RemoveAll(d.Logger.Filename)
		d.Close()
	}()

	tc := time.NewTimer(time.Minute * 3)

OutTop:
	for {
		select {
		case <-tc.C:
			tc.Stop()
			break OutTop
		default:
			d.Write([]byte(time.Now().Format(time.RFC3339) + "\n"))
			t.Log(time.Now().Format(time.RFC3339))
		}
	}
}
