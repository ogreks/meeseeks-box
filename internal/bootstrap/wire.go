//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ogreks/meeseeks-box/internal/ioc"
	iocApi "github.com/ogreks/meeseeks-box/internal/ioc/api"
)

func InitApiServer() *gin.Engine {
	wire.Build(
		// init configs
		ioc.InitConfig,
		// init logger
		ioc.InitLogger,
		// init orm
		ioc.InitORM,
		// web hook
		ioc.NewLarkClient,
		// init middleware
		iocApi.InitMiddleware,
		iocApi.InitJwtMiddleware,
		iocApi.InitWebHook,
		// init api server
		iocApi.InitApiServer,
	)
	return new(gin.Engine)
}
