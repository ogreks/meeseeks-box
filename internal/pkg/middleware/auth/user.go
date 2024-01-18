package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"net/http"
	"strings"
	"time"
)

type UserClaims struct {
	jwt.StandardClaims

	Content any `json:"content"`
}

type UserJwtMiddleware struct {
	HeaderKey      string
	RefreshKey     string
	RefreshTimeout time.Duration
	Token          token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)]
}

func NewUserJwtMiddleware(
	headerKey, RefreshKey string,
	refreshTimeout time.Duration,
	tk token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)]
) *UserJwtMiddleware {
	return &UserJwtMiddleware{
		HeaderKey: headerKey,
		RefreshKey: RefreshKey,
		RefreshTimeout: refreshTimeout,
		Token:     tk,
	}
}

// RefreshToken inner refresh token
func (uj *UserJwtMiddleware) RefreshToken(ctx *gin.Context, claim *UserClaims) {
	if time.Now().Sub(time.Unix(claim.ExpiresAt, 0)) > uj.RefreshTimeout {
		return
	}

	// refresh token
	rft := ctx.Request.Header.Get(uj.RefreshKey)
	if rft == "" {
		return
	}

	// TODO refresh token set header
}

// Builder middleware
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
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		// true
		uj.RefreshToken(ctx, c.(*UserClaims))

		// set claims to context
		ctx.Set("claims", c)
		ctx.Next()
	}
}
