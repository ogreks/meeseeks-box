package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router interface {
	Register(r *gin.Engine, logger *zap.Logger) Router
}

func InitRouter(g *gin.Engine, log *zap.Logger) {
}
