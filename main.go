package main

import (
	"github.com/ogreks/meeseeks-box/cmd"
	"os"
	"time"
)

// title meeseeks box cli
func main() {

	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			panic("timezone error: " + tz)
		}
	}

	cmd.Execute()
}
