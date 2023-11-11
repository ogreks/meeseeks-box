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
		// init config
		ioc.InitConfig,
		// init orm
		ioc.InitORM,
		// init api server
		iocApi.InitApiServer,
	)
	return new(gin.Engine)
}
