package migrate

import (
	"strings"

	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/ioc"

	// userModel "github.com/ogreks/meeseeks-box/internal/model/user"
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
)

type InitAutoMigrate struct {
	command.BaseCommand

	config   string
	function string
}

func NewInitAutoMigrateCommand() *InitAutoMigrate {
	i := &InitAutoMigrate{}

	initAutoMigrateCmd := &cobra.Command{
		Use:   "init",
		Short: "Meeseeks Box migrate init tool",
		Long:  "This is Meeseeks Box migrate init tool, it's `init` or `clean` and `reload`",
		RunE:  i.runCommand,
	}

	i.SetCommand(initAutoMigrateCmd)
	i.initVars()
	return i
}

func (i *InitAutoMigrate) initVars() {
	initAutoMigrateCmd := i.GetCommand()
	// flags: configs or -c
	initAutoMigrateCmd.Flags().StringVarP(
		&i.config, "configs", "c", "",
		"runtime configuration file or directory (default: workdir config.yml)",
	)
	// flags: function or -f
	initAutoMigrateCmd.Flags().StringVarP(
		&i.function, "function", "f", "",
		"enable features by separating their names with commas; all are enabled by default.",
	)
}

// runCommand is the main function of this command
func (i *InitAutoMigrate) runCommand(cmd *cobra.Command, args []string) error {
	configFile := i.config
	if configFile == "" {
		configFile = defaultConfigFile
	}

	cfg := configs.InitConfig(configFile)
	orm := ioc.InitORM(*cfg)

	db := orm.GetDB().Set("gorm:table_options", "ENGINE=InnoDB")

	var functions []any
	if i.function == "" {
		functions = getAllModels()
	} else {
		functions = getModels(strings.Split(i.function, ",")...)
	}

	return db.AutoMigrate(functions...)
}

// getModels returns models by function
func getModels(functions ...string) []any {
	var models []any
	for _, function := range functions {
		if val, ok := modelConfigs[function]; ok {
			models = append(models, val...)
		}
	}

	return models
}

// getAllModels returns all models
func getAllModels() []any {
	var models []any
	for _, val := range modelConfigs {
		models = append(models, val...)
	}

	return models
}
