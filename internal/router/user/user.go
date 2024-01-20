package user

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	userHandler "github.com/ogreks/meeseeks-box/internal/api/user"
	userJwt "github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

type UserRouter struct {
	userHandle    userHandler.Handler
	jwtMiddleware *userJwt.UserJwtMiddleware
	TokenManager  token.Token[string, func() (jwt.SigningMethod, []byte)]
}

func NewUserRouter(
	db orm.Repo,
	logger *zap.Logger,
	jwtMiddleware *userJwt.UserJwtMiddleware,
	token token.Token[string, func() (jwt.SigningMethod, []byte)],
) *UserRouter {
	return &UserRouter{
		TokenManager:  token,
		jwtMiddleware: jwtMiddleware,
		userHandle:    userHandler.New(db, logger, token),
	}
}

func (uRoute *UserRouter) Register(r *gin.Engine) *UserRouter {
	userRouter := r.Group("/api/user")
	{
		userRouter.POST("/login", uRoute.userHandle.Login)
		userRouter.POST("/register", uRoute.userHandle.Register)
		userRouter.POST("/github/login", uRoute.userHandle.LoginGITHub)

		mUserRouter := userRouter.Group("/", uRoute.jwtMiddleware.Builder())
		{
			mUserRouter.GET("/me", uRoute.userHandle.Me)
			mUserRouter.PUT("/refresh/token", uRoute.userHandle.RefersToken)
			mUserRouter.DELETE("/logout", uRoute.userHandle.Logout)
		}
	}
	return uRoute
}
