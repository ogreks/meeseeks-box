package message

import (
	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	"go.uber.org/atomic"
)

type MessageHandle struct {
	cli      *lark.Client
	dispatch *dispatcher.EventDispatcher

	EncryptKey        *atomic.String
	VerificationToken *atomic.String
}

func NewMessageHandler(cli *lark.Client, opts ...MessageHandleOptions) MessageHandleInterface {
	m := &MessageHandle{
		cli:               cli,
		EncryptKey:        atomic.NewString(""),
		VerificationToken: atomic.NewString(""),
	}

	for _, opt := range opts {
		opt(m)
	}

	m.init()

	return m
}

func (m *MessageHandle) init() {
	if m.VerificationToken.Load() == "" || m.EncryptKey.Load() == "" {
		return
	}

	m.dispatch = dispatcher.NewEventDispatcher(
		m.VerificationToken.Load(),
		m.EncryptKey.Load(),
	).OnP1MessageReceiveV1(m.msgReceivedHandler)
}

func (m *MessageHandle) RegisterRoute(path string, g gin.IRouter) {
	g.POST(path, sdkginext.NewEventHandlerFunc(m.dispatch))
}
