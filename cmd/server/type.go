package server

import (
	"fmt"
	"os"
)

type AppServerConfigure struct {
	Name              string
	DefaultConfigFile string
}

var app = AppServerConfigure{
	Name:              fmt.Sprintf("%s-server-api", "meeseeks"),
	DefaultConfigFile: fmt.Sprintf("%s/config/dev.yaml", workdir()),
}

func workdir() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workdir
}
