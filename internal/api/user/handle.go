package user

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	Login(*gin.Context)
}

type handler struct {
	logger *zap.Logger
}

func New(logger *zap.Logger) Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) i() {}
