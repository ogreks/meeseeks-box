package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/ogreks/meeseeks-box/configs"
	feishuMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/message"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router"
	"go.uber.org/zap"
)

func InitApiServer(
	db orm.Repo, // db
	logger *zap.Logger, // log
	middlewares []gin.HandlerFunc, // middleware
	jwtMiddleware *middleware.JwtMiddleware, // jwt middleware
	client *lark.Client, // feishu client
	msg feishuMessage.MessageHandleInterface, // feishu message event
) *gin.Engine {
	g := gin.New()

	g.Use(middlewares...)

	_ = router.InitRouter(&router.RouterHandler{
		Engine:            g,
		DB:                db,
		Log:               logger,
		AuthMiddleware:    jwtMiddleware,
		MessageDispatcher: msg,
	})

	return g
}

func InitMiddleware(logger *zap.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Recovery(logger),
		middleware.Trace(),
	}
}

func InitJwtMiddleware(cfg configs.Config) *middleware.JwtMiddleware {
	return middleware.NewJWTMiddleware(
		middleware.WithKeyFunc(func() (any, error) {
			return []byte(cfg.Jwt.Secret), nil
		}),
		middleware.WithClaims(&middleware.GlobalJWT{}),
		middleware.WithSigningMethod(jwt.SigningMethodHS512),
		middleware.WithJWTHeaderKey(cfg.Jwt.HeaderKey),
	)
}
