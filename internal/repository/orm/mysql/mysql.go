package mysql

import (
	"fmt"
	"time"

	"github.com/ogreks/meeseeks-box/config"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type MysqlRepo struct {
	orm.BaseRepo
}

func NewMysqlRepo(cfg config.Database) orm.Repo {
	repo := &MysqlRepo{}

	repo.SetConfig(cfg)

	return repo
}

func (m *MysqlRepo) Connection() error {
	cfg := m.GetConfig()
	driver, err := gorm.Open(mysql.Open(cfg.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   cfg.Prefix,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(cfg.Mode)),
	})

	if err != nil {
		return err
	}

	driver.Set("gorm:table_options", fmt.Sprintf("CHARSET=%s", cfg.Charset))

	db, err := driver.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConn)

	db.SetMaxIdleConns(cfg.MaxIdleConn)

	db.SetConnMaxLifetime(time.Minute * time.Duration(cfg.MaxLifetime))

	m.SetDB(driver)

	return nil
}
