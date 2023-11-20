package server

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type StopServer struct {
	command.BaseCommand
}

func NewStopServer() *StopServer {
	stopServer := &StopServer{}

	stopServerCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop Server",
		Long:  `Stop the server`,
		RunE:  stopServer.runCommand,
	}

	stopServer.SetCommand(stopServerCmd)
	return stopServer
}

func (s *StopServer) runCommand(cmd *cobra.Command, args []string) error {
	lockFile := fmt.Sprintf("%s.lock", app.Name)
	pid, err := os.ReadFile(lockFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("The corresponding executable file was not found. Please check whether the service is running.\ncode: %d", -4004)
		}
		return err
	}

	execCommand := exec.Command("kill", string(pid))
	err = execCommand.Start()
	if err != nil {
		return err
	}

	err = os.Remove(lockFile)
	if err != nil {
		return fmt.Errorf("can't remove %s.lock. %s", "meeseeks-box", err.Error())
	}

	fmt.Printf("service %s stopped \n", "meeseeks-box")
	return nil
}
