package message

import (
	"context"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"go.uber.org/atomic"
)

type MessageHandle struct {
	cli      *lark.Client
	dispatch *dispatcher.EventDispatcher

	EncryptKey        *atomic.String
	VerificationToken *atomic.String

	actions []Action
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

	// register action
	m.RegisterActions(
		&EmptyAction{},
		&VersionAction{},
		&HelpAction{},
	)

	return m
}

func (m *MessageHandle) init() {
	if m.VerificationToken.Load() == "" || m.EncryptKey.Load() == "" {
		return
	}

	m.dispatch = dispatcher.NewEventDispatcher(
		m.VerificationToken.Load(),
		m.EncryptKey.Load(),
	).
		OnP2MessageReceiveV1(m.msgReceivedHandler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			return nil
		})
}

func (m *MessageHandle) RegisterRoute(path string, g gin.IRouter) {
	g.POST(path, sdkginext.NewEventHandlerFunc(m.dispatch))
}
