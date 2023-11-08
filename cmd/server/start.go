package server

import (
	"context"
	"fmt"
	"github.com/ogreks/meeseeks-box/internal/api"
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type StartServer struct {
	command.BaseCommand
	daemon bool
}

func NewStartServer() *StartServer {
	startServer := &StartServer{}

	startServerCmd := &cobra.Command{
		Use:     "start",
		Short:   "Start server",
		Long:    `Start the server`,
		RunE:    startServer.runCommand,
		Aliases: []string{"s"},
	}

	startServerCmd.Flags().BoolVarP(&startServer.daemon, "daemon", "d", false, "daemon")

	startServer.SetCommand(startServerCmd)
	return startServer
}

func (s *StartServer) runCommand(cmd *cobra.Command, args []string) error {
	server := api.NewApi("localhost", "8000")

	if s.daemon {
		bin, err := filepath.Abs(os.Args[0])
		if err != nil {
			fmt.Printf("failed to get absolute path for command: %s \n", err.Error())
			return err
		}

		args := []string{"server", "start"}
		fmt.Printf("execute command: %s %s \n", bin, strings.Join(args, " "))
		execCommand := exec.Command(bin, args...)
		err = execCommand.Start()
		if err != nil {
			fmt.Printf("failed to start daemon thread: %s \n", err.Error())
			return err
		}

		pid := execCommand.Process.Pid
		_ = os.WriteFile(fmt.Sprintf("%s.lock", "meeseeks-box"), []byte(fmt.Sprintf("%d", pid)), 0666)
		fmt.Printf("service %s daemon thread started with pid %d \n", "meeseeks-box", pid)
		os.Exit(0)
	}

	return server.Start(context.Background())
}
