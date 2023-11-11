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
	mysql := &MysqlRepo{}

	mysql.SetConfig(cfg)

	return mysql
}

func (m *MysqlRepo) Connection() error {
	config := m.GetConfig()
	driver, err := gorm.Open(mysql.Open(config.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			TablePrefix:   config.Prefix,
		},
		Logger: logger.Default.LogMode(logger.LogLevel(config.Mode)),
	})

	if err != nil {
		return err
	}

	driver.Set("gorm:table_options", fmt.Sprintf("CHARSET=%s", config.Charset))

	db, err := driver.DB()
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(config.MaxOpenConn)

	db.SetMaxIdleConns(config.MaxIdleConn)

	db.SetConnMaxLifetime(time.Minute * time.Duration(config.MaxLifetime))

	m.SetDB(driver)

	return nil
}
