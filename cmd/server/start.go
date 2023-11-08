package server

import (
	"fmt"

	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type StartServer struct {
	command.BaseCommand
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

	startServer.SetCommand(startServerCmd)
	return startServer
}

func (s *StartServer) runCommand(cmd *cobra.Command, args []string) error {
	fmt.Println("Starting server...")
	return nil
}
