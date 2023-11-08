package main

import (
	"github.com/ogreks/meeseeks-box/cmd"
)

// @title meeseeks-box-cli
// @description MeeseeksBox CLI
// @name Authorization
// @schemes http https ws wss
// @basePath /
// @license.name MIT
// @license.url https://github.com/ogreks/meeseeks-box/blob/main/LICENSE

func main() {
	cmd.Execute()
}
