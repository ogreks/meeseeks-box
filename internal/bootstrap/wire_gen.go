// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/ioc"
	"github.com/ogreks/meeseeks-box/internal/ioc/api"
)

// Injectors from wire.go:

//go:generate wire
func InitApiServer() *gin.Engine {
	config := ioc.InitConfig()
	repo := ioc.InitORM(config)
	driver := ioc.InitLogDriver(config)
	logger := ioc.InitLogger(config, driver)
	v := api.InitMiddleware(logger, config)
	client := ioc.InitLarkClient(config, logger)
	messageHandleInterface := ioc.InitLarkMessageDispatcher(config, logger, repo, client)
	store := ioc.NewStore(config, logger)
	engine := api.InitApiServer(repo, logger, v, client, messageHandleInterface, store)
	return engine
}
