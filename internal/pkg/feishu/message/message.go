package message

import (
	"context"

	"github.com/gin-gonic/gin"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type (
	CardKind     string
	CardChatType string
)

var (
	GroupChatType = CardChatType("group")
	UserChatType  = CardChatType("personal")
)

const (
	GroupHandler = "group"
	UserHandler  = "personal"
)

type MessageHandleInterface interface {
	RegisterRoute(path string, g gin.IRouter)
}

// chain
func (m *MessageHandle) chain() bool {

	return true
}

func (m *MessageHandle) msgReceivedHandler(ctx context.Context, event *larkim.P1MessageReceiveV1) error {

	return nil
}
