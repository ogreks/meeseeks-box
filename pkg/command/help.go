package command

import "os"

// HelpGetWorkDir get current workdir
func HelpGetWorkDir() string {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return workdir
}
