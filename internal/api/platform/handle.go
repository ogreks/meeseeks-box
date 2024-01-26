package platform

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	pservice "github.com/ogreks/meeseeks-box/internal/service/platform"
	"go.uber.org/zap"
)

type Handle interface {
	i()

	// CreateSessionKey create open api account
	// Route /api/open/account [post]
	CreateSessionKey(ctx *gin.Context)

	// UpdateSessionKey update open api account
	// Route /api/open/account/:session_no [put]
	UpdateSessionKey(ctx *gin.Context)

	// SessionKeysSetStatus set status
	// Route /api/open/account/set/:status/status
	SessionKeysSetStatus(ctx *gin.Context)

	// SessionsKeysList get open api account manage
	// Route /api/open/account [get]
	SessionsKeysList(ctx *gin.Context)
}

type handle struct {
	logger  *zap.Logger
	service pservice.Service
}

func New(db orm.Repo, logger *zap.Logger) Handle {
	return &handle{
		logger:  logger,
		service: pservice.New(db, logger),
	}
}

func (h *handle) i() {}
