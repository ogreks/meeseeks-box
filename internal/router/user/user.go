package user

import (
	"github.com/gin-gonic/gin"
	userHandler "github.com/ogreks/meeseeks-box/internal/api/user"
	"github.com/ogreks/meeseeks-box/internal/router"
	"go.uber.org/zap"
)

type UserRouter struct {
	userHandle userHandler.Handler
}

func NewUserRouter(logger *zap.Logger) *UserRouter {
	return &UserRouter{
		userHandle: userHandler.New(logger),
	}
}

func (uRoute *UserRouter) Register(r *gin.Engine) router.Router {
	userRouter := r.Group("/api/user")
	{
		userRouter.POST("/login", uRoute.userHandle.Login)
	}
	return uRoute
}
