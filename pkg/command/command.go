package command

import (
	"github.com/spf13/cobra"
)

type Command interface {
	GetCommand() *cobra.Command
}

type BaseCommand struct {
	cmd *cobra.Command
}

func (b *BaseCommand) GetCommand() *cobra.Command {
	return b.cmd
}

func (b *BaseCommand) SetCommand(cmd *cobra.Command) {
	b.cmd = cmd
}
