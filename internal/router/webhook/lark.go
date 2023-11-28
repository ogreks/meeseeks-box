package webhook

import (
	"github.com/gin-gonic/gin"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
)

type Lark struct {
}

func NewLark() *Lark {
	return &Lark{}
}

func (l *Lark) Register(g *gin.RouterGroup, MessageDispatcher feishuUserMessage.UserMessageInterface) *Lark {
	MessageDispatcher.RegisterRoute("/lark/event", g)

	return l
}
