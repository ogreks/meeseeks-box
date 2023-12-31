package server

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ogreks/meeseeks-box/configs"
	api "github.com/ogreks/meeseeks-box/internal/bootstrap"
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type StartServer struct {
	command.BaseCommand
	daemon bool
	config string
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
	startServer.initVars()
	return startServer
}

func (s *StartServer) initVars() {
	startServerCmd := s.GetCommand()
	// flags: configs or -c
	startServerCmd.Flags().StringVarP(
		&s.config, "configs", "c", "",
		"runtime configuration files or directory (default: workdir config.yml)",
	)
	// flags: daemon or -d
	startServerCmd.Flags().BoolVarP(
		&s.daemon, "daemon", "d", false,
		"start daemon mode (default: false)",
	)
}

func (s *StartServer) runCommand(cmd *cobra.Command, args []string) error {
	configFile := s.config
	if configFile == "" {
		configFile = app.DefaultConfigFile
	}

	cfg := configs.InitConfig(configFile)

	server := api.NewApi(cfg.GetServer().Addr, cfg.GetServer().Port)
	if s.daemon {
		bin, err := filepath.Abs(os.Args[0])
		if err != nil {
			fmt.Printf("failed to get absolute path for command: %s \n", err.Error())
			return err
		}

		args := []string{"server", "start", "-c", configFile}
		fmt.Printf("execute command: %s %s \n", bin, strings.Join(args, " "))
		execCommand := exec.Command(bin, args...)
		err = execCommand.Start()
		if err != nil {
			fmt.Printf("failed to start daemon thread: %s \n", err.Error())
			return err
		}

		pid := execCommand.Process.Pid
		_ = os.WriteFile(fmt.Sprintf("%s.lock", app.Name), []byte(fmt.Sprintf("%d", pid)), 0666)
		fmt.Printf("service %s daemon thread started with pid %d \n", "meeseeks-box", pid)
		os.Exit(0)
	}

	return server.Start(context.Background())
}
