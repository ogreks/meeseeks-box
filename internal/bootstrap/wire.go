//go:build wireinject
// +build wireinject

package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/ogreks/meeseeks-box/internal/ioc"
	iocApi "github.com/ogreks/meeseeks-box/internal/ioc/api"
)

//go:generate wire
func InitApiServer() *gin.Engine {
	wire.Build(
		// init configs
		ioc.InitConfig,
		// init logger
		ioc.InitLogDriver,
		ioc.InitLogger,
		// init orm
		ioc.InitORM,
		// web hook
		ioc.InitLarkClient,
		ioc.InitLarkMessageDispatcher,
		// init middleware
		iocApi.InitMiddleware,
		iocApi.InitJwtMiddleware,
		// init api server
		iocApi.InitApiServer,
	)
	return new(gin.Engine)
}
