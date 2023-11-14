package mysql

import (
	"fmt"

	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Mysql struct {
	orm.BaseRepo
}

func NewMysql(cfg configs.Database) orm.Repo {
	repo := &Mysql{}

	repo.SetConfig(cfg)

	return repo
}

func (m *Mysql) Connection() error {
	cfg := m.GetConfig()
	driver, err := gorm.Open(mysql.Open(cfg.Source), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.Mode)),
	})

	if err != nil {
		return err
	}

	driver.Set("gorm:table_options", fmt.Sprintf("CHARSET=%s", cfg.Charset))

	m.SetDB(driver)

	return nil
}
