package webhook

import (
	"github.com/gin-gonic/gin"
	feishuMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/message"
)

func Register(g *gin.Engine, messageEvent feishuMessage.MessageHandleInterface) {
	l := NewLark()

	r := g.Group("/webhook")

	l.Register(r, messageEvent)
}
