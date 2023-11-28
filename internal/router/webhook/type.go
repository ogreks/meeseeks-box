package webhook

import (
	"github.com/gin-gonic/gin"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
)

func Register(g *gin.Engine, messageEvent feishuUserMessage.UserMessageInterface) {
	l := NewLark()

	r := g.Group("/webhook")

	l.Register(r, messageEvent)
}
