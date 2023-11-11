package orm

import (
	"github.com/ogreks/meeseeks-box/config"
	"gorm.io/gorm"
)

type DB_TYPE string

const (
	DB_TYPE_MYSQL  DB_TYPE = "mysql"
	DB_TYPE_SQLITE DB_TYPE = "sqlite"
)

type Repo interface {
	i()
	Connection() error
	GetDB() *gorm.DB
	Close() error
}

type BaseRepo struct {
	db *gorm.DB
	c  config.Database
}

func NewRepo(cfg config.Database) Repo {
	return &BaseRepo{}
}

func (r *BaseRepo) i() {}

func (r *BaseRepo) Connection() error {
	return nil
}

func (r *BaseRepo) SetDB(db *gorm.DB) Repo {
	r.db = db
	return r
}

func (r *BaseRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *BaseRepo) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (r *BaseRepo) SetConfig(c config.Database) Repo {
	r.c = c
	return r
}

func (r *BaseRepo) GetConfig() config.Database {
	return r.c
}
