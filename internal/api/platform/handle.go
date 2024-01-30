package platform

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	nofyservice "github.com/ogreks/meeseeks-box/internal/service/notify"
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

	// CreateAgreement create agreement content
	// Route /api/open/agreement [post]
	CreateAgreement(ctx *gin.Context)

	// UpdateAgreement update agreement detail
	// Route /api/open/agreement/:agreement_no
	UpdateAgreement(ctx *gin.Context)

	// DeleteAgreement delete agreement
	// Route /api/open/agreement [delete]
	DeleteAgreement(ctx *gin.Context)

	// AgreementPaginate get agreement paginate list
	// Route /api/open/agreements [get]
	AgreementPaginate(ctx *gin.Context)
}

type handle struct {
	logger   *zap.Logger
	platform pservice.Service
	notify   nofyservice.Service
}

func New(db orm.Repo, logger *zap.Logger) Handle {
	return &handle{
		logger:   logger,
		platform: pservice.New(db, logger),
		notify:   nofyservice.New(db, logger),
	}
}

func (h *handle) i() {}
