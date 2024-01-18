package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"net/http"
	"strings"
)

type UserClaims struct {
	jwt.StandardClaims

	Content any `json:"content"`
}

type UserJwtMiddleware struct {
	HeaderKey string
	Token     token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)]
}

func NewUserJwtMiddleware(headerKey string, tk token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)]) *UserJwtMiddleware {
	return &UserJwtMiddleware{
		HeaderKey: headerKey,
		Token:     tk,
	}
}

func (uj *UserJwtMiddleware) Builder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get token from header
		tokenHeader := ctx.GetHeader(uj.HeaderKey)
		if tokenHeader == "" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) < 2 {
			_ = ctx.AbortWithError(http.StatusUnauthorized, errors.New("token is invalid"))
			return
		}

		tk := segs[1]
		// parse token
		c, err := uj.Token.Validate(tk)
		if err != nil {
			if !errors.Is(err, token.ErrTokenNotFound) {
				_ = ctx.AbortWithError(http.StatusUnauthorized, err)
				return
			}
		}

		// set claims to context
		ctx.Set("claims", c)
		ctx.Next()
	}
}
