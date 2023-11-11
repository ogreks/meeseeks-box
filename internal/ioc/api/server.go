package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router"
	"go.uber.org/zap"
)

func InitApiServer(db orm.Repo, logger *zap.Logger) *gin.Engine {
	g := gin.Default()

	_ = router.InitRouter(g, logger)

	return g
}
