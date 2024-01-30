package open

import (
	"github.com/gin-gonic/gin"
	ph "github.com/ogreks/meeseeks-box/internal/api/platform"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	ps "github.com/ogreks/meeseeks-box/internal/service/platform"
	"go.uber.org/zap"
)

type NotifyRouter struct {
	phd    ph.Handle
	db     orm.Repo
	logger *zap.Logger
}

func NewNotifyRouter(
	db orm.Repo,
	logger *zap.Logger,
) *NotifyRouter {
	return &NotifyRouter{
		db:     db,
		logger: logger,
		phd:    ph.New(db, logger),
	}
}

func (p *NotifyRouter) Register(r *gin.Engine) *NotifyRouter {
	rg := r.Group("/api/open", auth.NewOpenMiddleware(ps.New(p.db, p.logger)).Builder())
	{
		rg.GET("/agreements", p.phd.AgreementPaginate)
		rg.POST("/agreement", p.phd.CreateAgreement)
		rg.PUT("/agreement/:agreement_no", p.phd.UpdateAgreement)
		rg.DELETE("/agreement/:agreement_no", p.phd.DeleteAgreement)
	}

	return p
}
