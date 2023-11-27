package webhook

import (
	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
)

type Lark struct {
	handle *middleware.Lark
}

func NewLark(h *middleware.Lark) *Lark {
	return &Lark{
		handle: h,
	}
}

func (l *Lark) Register(r *gin.Engine) *Lark {
	g := r.Group("webhook/lark/")
	{
		g.POST("/event", sdkginext.NewEventHandlerFunc(l.handle.GetHandle()))
		g.POST("/card", sdkginext.NewCardActionHandlerFunc(l.handle.GetCardHandle()))
	}

	return l
}
