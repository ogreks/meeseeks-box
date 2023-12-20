package webhook

import (
	"github.com/gin-gonic/gin"
	feishuCardMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/card"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
)

type Lark struct {
}

func NewLark() *Lark {
	return &Lark{}
}

func (l *Lark) Register(g *gin.RouterGroup, MessageDispatcher feishuUserMessage.UserMessageInterface, CardDispatcher feishuCardMessage.CardMessagerInterface) *Lark {
	MessageDispatcher.RegisterRoute("/lark/event", g)
	{
		// message event
	}

	CardDispatcher.RegisterRoute("/lark/card", g)
	{
		// card event
	}

	return l
}
