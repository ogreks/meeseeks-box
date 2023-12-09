package webhook

import (
	"github.com/gin-gonic/gin"
	feishuCardMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/card"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
)

func Register(g *gin.Engine, messageEvent feishuUserMessage.UserMessageInterface, cardEvent feishuCardMessage.CardMessagerInterface) {
	l := NewLark()

	r := g.Group("/webhook")

	l.Register(r, messageEvent, cardEvent)
}
