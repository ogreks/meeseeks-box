package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/configs"
	userJwt "github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, db orm.Repo, logger *zap.Logger, tokenStore token.Store[string]) {
	// token manager
	tk := token.NewDefaultToken[string, func() (jwt.SigningMethod, []byte)](
		token.WithStore(tokenStore),
		token.WithExpire(time.Duration(configs.GetConfig().Jwt.Expire)*time.Second),
		token.WithFun(func() (jwt.SigningMethod, []byte) {
			return jwt.SigningMethodHS512, []byte(configs.GetConfig().Jwt.Secret)
		}),
		token.WithClaims(&userJwt.UserClaims{}),
	)

	// init jwt middleware
	j := userJwt.NewUserJwtMiddleware(
		configs.GetConfig().Jwt.HeaderKey,
		configs.GetConfig().Jwt.RefersKey,
		time.Duration(configs.GetConfig().Jwt.RefreshTimeout)*time.Second,
		tk,
	)

	ur := NewUserRouter(db, logger, j, tk)
	ur.Register(r)
}
