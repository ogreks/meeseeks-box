//go:build wireinject
// +build wireinject

package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	iocWeb "github.com/ogreks/meeseeks-box/internal/ioc/api"
)

func InitApiServer() *gin.Engine {
	wire.Build(
		iocWeb.InitApiServer,
	)
	return new(gin.Engine)
}
