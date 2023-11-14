package orm

import (
	"time"

	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
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
	c  configs.Database
}

func NewRepo(cfg configs.Database) Repo {
	return &BaseRepo{}
}

func (r *BaseRepo) i() {}

func (r *BaseRepo) Connection() error {
	return nil
}

// setDBConfig set db configs
func (r *BaseRepo) setDBConfig(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(r.c.MaxIdleConn)
		sqlDB.SetMaxOpenConns(r.c.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Minute * time.Duration(r.c.MaxLifetime))
	}
}

// setNowFunc set now func
func (r *BaseRepo) setNowFunc(db *gorm.DB) {
	db.NowFunc = func() time.Time {
		return time.Now().Local()
	}
}

// setLogger set logger
func (r *BaseRepo) setLogger(db *gorm.DB) {
	loggerWiths := []logger.Option{
		logger.WithTimeLayout("2006-01-02 15:04:05"),
		logger.WithFilePath(r.c.LogPath),
	}

	if glogger.LogLevel(r.c.Mode) == glogger.Info {
		loggerWiths = append(loggerWiths, logger.WithLevel(zap.DebugLevel))
	}

	databaseLogger, err := logger.NewJsonLogger(loggerWiths...)
	if err != nil {
		panic(err)
	}

	defer func() {
		databaseLogger.Sync()
	}()

	log := NewLogger(databaseLogger)
	log.SetAsDefault() // // optional: configure gorm to use this zapgorm.Logger for callbacks
	db.Logger = log
}

// SetDB set db
func (r *BaseRepo) SetDB(db *gorm.DB) Repo {
	r.setDBConfig(db)
	r.setNowFunc(db)
	r.setLogger(db)

	r.db = db

	return r
}

// GetDB get db
func (r *BaseRepo) GetDB() *gorm.DB {
	return r.db
}

// Close db
func (r *BaseRepo) Close() error {
	db, err := r.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

// SetConfig set configs
func (r *BaseRepo) SetConfig(c configs.Database) Repo {
	r.c = c
	return r
}

// GetConfig get configs
func (r *BaseRepo) GetConfig() configs.Database {
	return r.c
}
