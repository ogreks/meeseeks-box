package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/configs"
	userJwt "github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/service/user"
	"go.uber.org/zap"
	"time"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Login user login
	// Router /api/user/login [post]
	Login(*gin.Context)

	// LoginGITHub login by github account
	// Router /api/user/github/login [post]
	LoginGITHub(ctx *gin.Context)

	// Register user register
	// Router /api/user/register [post]
	Register(*gin.Context)

	// Me user info
	// Router /api/user/me [get]
	Me(*gin.Context)

	// RefersToken refers token
	// Router /api/user/refersh/token [put]
	RefersToken(ctx *gin.Context)

	// Logout logout user token
	// Router /api/user/logout [delete]
	Logout(ctx *gin.Context)
}

type handler struct {
	logger       *zap.Logger
	service      user.Service
	tokenManager token.Token[string, func() (jwt.SigningMethod, []byte)]
}

func New(db orm.Repo, logger *zap.Logger, tkm token.Token[string, func() (jwt.SigningMethod, []byte)]) Handler {
	return &handler{
		logger:       logger,
		service:      user.New(db, logger),
		tokenManager: tkm,
	}
}

func (h *handler) i() {}

// createToken create user token
func (h *handler) createToken(ctx *gin.Context, aid string) (string, string, error) {
	tk, err := h.tokenManager.CreateToken(ctx.Request.Context(), &userJwt.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   aid,
			Audience:  "api",
			Issuer:    configs.GetConfig().Jwt.Issuer,
			ExpiresAt: time.Now().Add(time.Duration(configs.GetConfig().Jwt.Expire) * time.Second).Unix(),
		},
		Content: aid,
	}, time.Duration(configs.GetConfig().Jwt.Expire)*time.Second+(20*time.Second))
	if err != nil {
		return "", "", err
	}

	ctx.Request.Header.Set(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", tk))
	ctx.Header(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", tk))

	rest, err := h.tokenManager.CreateToken(ctx.Request.Context(), &userJwt.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   aid,
			Audience:  "refresh",
			Issuer:    configs.GetConfig().Jwt.Issuer,
			ExpiresAt: time.Now().Add(time.Duration(configs.GetConfig().Jwt.RefreshTimeout) * time.Second).Unix(),
		},
		Content: aid,
	}, time.Duration(configs.GetConfig().Jwt.RefreshTimeout)*time.Second+(20*time.Second))
	if err != nil {
		return "", "", err
	}

	ctx.Request.Header.Set(configs.GetConfig().Jwt.RefersKey, rest)
	ctx.Header(configs.GetConfig().Jwt.RefersKey, rest)

	return tk, rest, nil
}
