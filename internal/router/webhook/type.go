package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
)

func Register(g *gin.Engine, lark *middleware.Lark) {
	l := NewLark(lark)
	l.Register(g)
}
