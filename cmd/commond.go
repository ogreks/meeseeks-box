package cmd

import (
	"fmt"

	"github.com/ogreks/meeseeks-box/cmd/server"
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

var (
	FlagPrint = false
)

type commandsBuilder struct {
	commands    []command.Command
	rootCommand *cobra.Command
}

// newCommandsBuilder 是一个用于构建命令的结构体
func newCommandsBuilder() *commandsBuilder {
	// 创建一个 commandsBuilder 实例并返回
	return &commandsBuilder{
		commands: make([]command.Command, 0),
		rootCommand: &cobra.Command{
			Use:               ProjectName,
			Short:             "This is the Meeseeks box CLI.",
			Long:              "This is the Meeseeks box CLI. Serve Mr. Meeseeks! Look at me! \n" + getVersionFmt(),
			Version:           getVersionFmt(),
			TraverseChildren:  true,
			DisableAutoGenTag: true,
			Run: func(cmd *cobra.Command, args []string) {
				if FlagPrint {
					fmt.Println(getVersionFmt())
					return
				}

				fmt.Println(TEMPLATE)
				cmd.Help()
			},
		},
	}
}

func (b *commandsBuilder) GetCommand() *cobra.Command {
	return b.rootCommand
}

func (b *commandsBuilder) addCommands(cmd ...command.Command) *commandsBuilder {
	b.commands = append(b.commands, cmd...)
	return b
}

func (b *commandsBuilder) addAll() *commandsBuilder {
	b.addCommands(
		server.NewServerCommand(),
	)

	return b
}

func (b *commandsBuilder) builder() *commandsBuilder {
	for _, command := range b.commands {
		b.rootCommand.AddCommand(command.GetCommand())
	}
	return b
}
