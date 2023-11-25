package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router/user"
	"go.uber.org/zap"
)

type Router interface {
	Register(r *gin.Engine) Router
}

func InitRouter(g *gin.Engine, db orm.Repo, log *zap.Logger, authMiddleware *middleware.JwtMiddleware, client *middleware.Lark) error {

	user.Register(g, db, log, authMiddleware)

	return nil
}
