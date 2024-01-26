package open

import (
	"github.com/gin-gonic/gin"
	ph "github.com/ogreks/meeseeks-box/internal/api/platform"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	ps "github.com/ogreks/meeseeks-box/internal/service/platform"
	"go.uber.org/zap"
)

type PlatformRouter struct {
	phd    ph.Handle
	db     orm.Repo
	logger *zap.Logger
}

func NewPlatformRouter(
	db orm.Repo,
	logger *zap.Logger,
) *PlatformRouter {
	return &PlatformRouter{
		db:     db,
		logger: logger,
		phd:    ph.New(db, logger),
	}
}

func (p *PlatformRouter) Register(r *gin.Engine) *PlatformRouter {
	rg := r.Group("/api/open", auth.NewOpenMiddleware(ps.New(p.db, p.logger)).Builder())
	{
		rg.GET("/accounts", p.phd.SessionsKeysList)
		rg.POST("/account", p.phd.CreateSessionKey)
		account := rg.Group("/account")
		{
			account.PUT("/set/:status/status", p.phd.SessionKeysSetStatus)
			account.PUT("/:no", p.phd.UpdateSessionKey)
		}
	}

	return p
}
