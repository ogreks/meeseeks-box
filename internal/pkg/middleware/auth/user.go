package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"net/http"
	"strconv"
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
	HeaderTimeout  time.Duration
	RefreshTimeout time.Duration
	Token          token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)]
}

func NewUserJwtMiddleware(
	headerKey, RefreshKey string,
	headerTimeout, refreshTimeout time.Duration,
	tk token.Token[string, func() (jwt.SigningMethod, []byte, jwt.Claims)],
) *UserJwtMiddleware {
	return &UserJwtMiddleware{
		HeaderKey:      headerKey,
		RefreshKey:     RefreshKey,
		HeaderTimeout:  headerTimeout,
		RefreshTimeout: refreshTimeout,
		Token:          tk,
	}
}

// RefreshToken inner refresh token
func (uj *UserJwtMiddleware) RefreshToken(ctx *gin.Context, userTk string, claim *UserClaims) {
	if time.Now().Sub(time.Unix(claim.ExpiresAt, 0)).Abs() > (uj.RefreshTimeout * time.Second) {
		return
	}

	// refresh token
	rft := ctx.Request.Header.Get(uj.RefreshKey)
	if rft == "" {
		return
	}

	rftc, err := uj.Token.Validate(rft)
	if err != nil {
		return
	}

	claim.ExpiresAt = time.Now().Add(uj.HeaderTimeout * time.Second).Unix()
	utk, err := uj.Token.RefreshToken(ctx.Request.Context(), userTk, claim, (uj.HeaderTimeout+20)*time.Second)
	if err != nil {
		// TODO 记录日志
	}

	// refresh user api token
	if utk != "" {
		ctx.Writer.Header().Set(uj.HeaderKey, fmt.Sprintf("Bearer %s", utk))
		//ctx.Request.Response.Header.Set(uj.HeaderKey, fmt.Sprintf("Bearer %s", utk))
		ctx.Writer.Header().Set(fmt.Sprintf("%s_At", uj.HeaderKey), strconv.FormatInt(int64(uj.HeaderTimeout-5), 10))
		ctx.Set("claims", claim)
	}

	uc, ok := rftc.(*UserClaims)
	if !ok {
		return
	}

	// determine whether the token needs to be refreshed
	if time.Now().Sub(time.Unix(uc.ExpiresAt, 0)).Abs() > (uj.HeaderTimeout+10)*time.Second {
		return
	}

	uc.ExpiresAt = time.Now().Add(uj.RefreshTimeout * time.Second).Unix()
	refreshToken, err := uj.Token.RefreshToken(ctx.Request.Context(), rft, uc, (uj.RefreshTimeout+20)*time.Second)
	if err != nil {
		return
	}

	ctx.Writer.Header().Set(uj.RefreshKey, refreshToken)
	//ctx.Request.Response.Header.Set(uj.RefreshKey, refreshToken)
	ctx.Writer.Header().Set(fmt.Sprintf("%s_At", uj.RefreshKey), strconv.FormatInt(int64(uj.RefreshTimeout-5), 10))
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
		if err != nil && !errors.Is(err, token.ErrTokenTimeout) {
			_ = ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		if errors.Is(err, token.ErrTokenTimeout) {
			if ctx.Request.Header.Get(uj.RefreshKey) == "" {
				_ = ctx.AbortWithError(http.StatusUnauthorized, err)
				return
			}
		}

		if c.(*UserClaims).Audience != "api" {
			_ = ctx.AbortWithError(http.StatusUnauthorized, token.ErrInvalidKey)
			return
		}

		ctx.Set("claims", c)

		// true
		uj.RefreshToken(ctx, tk, c.(*UserClaims))

		// set claims to context
		ctx.Next()
	}
}
