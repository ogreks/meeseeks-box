package migrate

import (
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type Migrate struct {
	command.BaseCommand
}

func NewAutoMigrateCommand() *Migrate {
	a := &Migrate{}

	autoMigrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Meeseeks Box migrate tool",
		Long:  "This is Meeseeks Box migrate tool, it's `init` or `clean` and `reload`",
		RunE:  a.runCommand,
	}

	autoMigrateCmd.AddCommand(NewGenMigrateCommand().GetCommand())
	autoMigrateCmd.AddCommand(NewInitAutoMigrateCommand().GetCommand())

	a.SetCommand(autoMigrateCmd)
	return a
}

func (a *Migrate) runCommand(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
