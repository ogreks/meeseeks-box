package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/router/user"
	"go.uber.org/zap"
)

type Router interface {
	Register(r *gin.Engine) Router
}

func InitRouter(g *gin.Engine, log *zap.Logger) error {

	user.Register(g, log)

	return nil
}
