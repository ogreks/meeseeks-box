package server

import (
	"fmt"

	"github.com/ogreks/meeseeks-box/pkg/command"
)

type AppServerConfigure struct {
	Name              string
	DefaultConfigFile string
}

var app = AppServerConfigure{
	Name:              fmt.Sprintf("%s-server-api", "meeseeks"),
	DefaultConfigFile: fmt.Sprintf("%s/config.yaml", command.HelpGetWorkDir()),
}
