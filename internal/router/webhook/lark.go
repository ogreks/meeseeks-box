package webhook

import (
	"github.com/gin-gonic/gin"
	feishuMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/message"
)

type Lark struct {
}

func NewLark() *Lark {
	return &Lark{}
}

func (l *Lark) Register(g *gin.RouterGroup, MessageDispatcher feishuMessage.MessageHandleInterface) *Lark {
	MessageDispatcher.RegisterRoute("/lark/event", g)
	{
		// message event
	}

	return l
}
