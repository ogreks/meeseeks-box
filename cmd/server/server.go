package server

import (
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type Server struct {
	command.BaseCommand
}

func NewServerCommand() *Server {
	s := &Server{}

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Meeseeks Box server tool",
		Long:  "This is Meeseeks Box Server Tool, it's `start` or `stop` and `reload`",
		RunE:  s.runCommand,
	}

	serverCmd.AddCommand(NewStartServer().GetCommand())

	s.SetCommand(serverCmd)
	return s
}

func (s *Server) runCommand(cmd *cobra.Command, args []string) error {
	cmd.Help()
	return nil
}
