package user

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/service/user"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Login user login
	// @Tags API.User
	// @Router /api/user/login [post]
	Login(*gin.Context)

	// LoginGITHub login by github account
	// @Tags API.User
	// @Router /api/user/github/login [post]
	LoginGITHub(ctx *gin.Context)

	// Register user register
	// @Tags API.User
	// @Router /api/user/register [post]
	Register(*gin.Context)

	// Me user info
	// @Tags API.User
	// @Router /api/user/me [get]
	Me(*gin.Context)
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
