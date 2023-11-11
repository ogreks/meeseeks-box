package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Register(r *gin.Engine, logger *zap.Logger) {
	ur := NewUserRouter(logger)
	ur.Register(r)
}
