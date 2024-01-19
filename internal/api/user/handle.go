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
