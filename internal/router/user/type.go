package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, db orm.Repo, logger *zap.Logger, jwtMiddleware *middleware.JwtMiddleware) {
	ur := NewUserRouter(db, logger, jwtMiddleware)
	ur.Register(r)
}
