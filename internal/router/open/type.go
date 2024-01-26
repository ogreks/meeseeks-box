package open

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, db orm.Repo, logger *zap.Logger) {
	p := NewPlatformRouter(db, logger)
	p.Register(r)
}
