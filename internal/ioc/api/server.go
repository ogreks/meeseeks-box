package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router"
	"go.uber.org/zap"
)

func InitApiServer(db orm.Repo, logger *zap.Logger, middlewares []gin.HandlerFunc, jwtMiddleware *middleware.JwtMiddleware) *gin.Engine {
	g := gin.New()

	g.Use(middlewares...)

	_ = router.InitRouter(g, db, logger, jwtMiddleware)

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
