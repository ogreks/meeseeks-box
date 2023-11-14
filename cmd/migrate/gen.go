package migrate

import (
	"gorm.io/gorm"
	"strings"

	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/ioc"
	"github.com/ogreks/meeseeks-box/pkg/command"
	"github.com/spf13/cobra"
	"gorm.io/gen"
)

type GenMigrate struct {
	command.BaseCommand
	config string
	tablse string
}

func NewGenMigrateCommand() *GenMigrate {
	g := &GenMigrate{}

	genMigrateCmd := &cobra.Command{
		Use:   "gen",
		Short: "Meeseeks Box migrate gen model/dao",
		Long:  "This is Meeseeks Box migrate gen model/dao, it's -c configs.yaml",
		RunE:  g.runCommand,
	}

	g.SetCommand(genMigrateCmd)
	g.initVars()
	return g
}

func (g *GenMigrate) initVars() {
	genMigrateCmd := g.GetCommand()
	// flags: configs or -c
	genMigrateCmd.Flags().StringVarP(
		&g.config, "configs", "c", "",
		"runtime configuration files or directory (default: workdir configs.yaml)",
	)
	// flags: table or -t
	genMigrateCmd.Flags().StringVarP(
		&g.tablse, "table", "t", "",
		"table name `,` split (default: all tables)",
	)
}

func (g *GenMigrate) runCommand(cmd *cobra.Command, args []string) error {
	configFile := g.config
	if configFile == "" {
		configFile = defaultConfigFile
	}

	cfg := configs.InitConfig(configFile)
	orm := ioc.InitORM(*cfg)

	ggorm := gen.NewGenerator(gen.Config{
		OutPath:           daoPath,
		FieldNullable:     false,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		ModelPkgPath:      "model",
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	ggorm.UseDB(orm.GetDB())
	ggorm.WithDataTypeMap(map[string]func(columnType gorm.ColumnType) (dataType string){
		"timestamp": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "*time.Time"
			}
			return "time.Time"
		},
	})

	setModels := make([]interface{}, 0)
	if g.tablse != "" {
		var (
			tables = strings.Split(g.tablse, ",")
		)

		for index, table := range tables {
			setModels[index] = ggorm.GenerateModel(table)
		}

	} else {
		setModels = append(setModels, ggorm.GenerateAllTable()...)
	}

	ggorm.ApplyBasic(
		ggorm.GenerateAllTable()...,
	)

	ggorm.ApplyInterface(func() {})

	ggorm.Execute()

	return nil
}
