package immessagereceive

import (
	"context"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
)

type Handler interface {
	Version() feishu.OnP2MessageReceiveV1
}

type handler struct {
	msgDispatcher feishuUserMessage.UserMessageInterface
}

func NewHandler(msg feishuUserMessage.UserMessageInterface) Handler {
	return &handler{
		msgDispatcher: msg,
	}
}

func (h *handler) Version() feishu.OnP2MessageReceiveV1 {
	return func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {

		return nil
	}
}
