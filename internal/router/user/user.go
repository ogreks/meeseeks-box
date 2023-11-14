package user

import (
	"github.com/gin-gonic/gin"
	userHandler "github.com/ogreks/meeseeks-box/internal/api/user"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type UserRouter struct {
	userHandle    userHandler.Handler
	jwtMiddleware *middleware.JwtMiddleware
}

func NewUserRouter(db orm.Repo, logger *zap.Logger, jwtMiddleware *middleware.JwtMiddleware) *UserRouter {
	return &UserRouter{
		jwtMiddleware: jwtMiddleware,
		userHandle:    userHandler.New(db, logger),
	}
}

func (uRoute *UserRouter) Register(r *gin.Engine) *UserRouter {
	userRouter := r.Group("/api/user")
	{
		userRouter.POST("/login", uRoute.userHandle.Login)
		userRouter.POST("/register", uRoute.userHandle.Register)

		mUserRouter := userRouter.Group("/", uRoute.jwtMiddleware.Builder())
		{
			mUserRouter.GET("/me", uRoute.userHandle.Me)
		}
	}
	return uRoute
}
