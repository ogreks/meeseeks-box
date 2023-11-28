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

func InitApiServer() *gin.Engine {
	config := ioc.InitConfig()
	repo := ioc.InitORM(config)
	logger := ioc.InitLogger()
	v := api.InitMiddleware(logger)
	jwtMiddleware := api.InitJwtMiddleware(config)
	client := ioc.InitLarkClient(config, logger)
	userMessageInterface := ioc.InitLarkMessageDispatcher(config, logger, repo, client)
	engine := api.InitApiServer(repo, logger, v, jwtMiddleware, client, userMessageInterface)
	return engine
}
