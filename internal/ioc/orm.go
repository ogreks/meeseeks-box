package ioc

import (
	"errors"
	"fmt"

	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/dao"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/repository/orm/mysql"
)

func InitORM(cfg configs.Config) (r orm.Repo) {
	database := cfg.Database
	driver := orm.DB_TYPE(database.Driver)
	switch driver {
	case orm.DB_TYPE_MYSQL:
		r = mysql.NewMysql(database)
	default:
		panic(errors.New(fmt.Sprintf("orm configure error: input [%s] unknown, use mysql、sqlite", driver)))
	}

	err := r.Connection()
	if err != nil {
		panic(err)
	}

	dao.SetDefault(r.GetDB())

	return
}
