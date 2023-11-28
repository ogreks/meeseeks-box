package feishu

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/syncx"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type UserMessageInterface interface {
	RegisterEvent(key string, handler interface{}) UserMessageInterface
	RegisterRoute(path string, g *gin.RouterGroup)
}

type UserMessageOption func(u *UserMessage)

type UserMessage struct {
	EventDispatcher *syncx.Map[string, interface{}]

	cli *lark.Client

	EnableOnP2MessageReceiveV1 *atomic.Bool
	EnableOnP2MessageReadV1    *atomic.Bool
	EnableOnP2UserCreatedV3    *atomic.Bool

	EncryptKey        *atomic.String
	VerificationToken *atomic.String

	handle *dispatcher.EventDispatcher
}

func NewUserMessage(log *zap.Logger, db orm.Repo, client *lark.Client, opts ...UserMessageOption) *UserMessage {
	u := &UserMessage{
		cli:             client,
		EventDispatcher: &syncx.Map[string, interface{}]{},

		EnableOnP2MessageReceiveV1: atomic.NewBool(false),
		EnableOnP2MessageReadV1:    atomic.NewBool(false),
		EnableOnP2UserCreatedV3:    atomic.NewBool(false),

		EncryptKey:        atomic.NewString(""),
		VerificationToken: atomic.NewString(""),
	}

	for _, opt := range opts {
		opt(u)
	}

	u.init()

	return u
}

func (u *UserMessage) init() {
	if u.VerificationToken.Load() == "" || u.EncryptKey.Load() == "" {
		return
	}

	u.handle = dispatcher.NewEventDispatcher(u.VerificationToken.Load(), u.EncryptKey.Load())

	// register event callback
	u.handle.OnP2MessageReadV1(u.OnP2MessageReadV1())
	u.handle.OnP2MessageReceiveV1(u.OnP2MessageReceiveV1())
	u.handle.OnP2UserCreatedV3(u.OnP2UserCreatedV3())
}

func (u *UserMessage) RegisterEvent(key string, handler interface{}) UserMessageInterface {
	u.EventDispatcher.Store(key, handler)
	return u
}

func (u *UserMessage) RegisterRoute(path string, g *gin.RouterGroup) {
	g.POST(path, sdkginext.NewEventHandlerFunc(u.handle))
}

// SendCloseMessage send close message
func (u *UserMessage) SendReplyMessage(ctx context.Context, messageId string, message string, msgType string) error {
	replyBody := larkim.NewReplyMessageReqBodyBuilder().
		MsgType(msgType).
		Content(message).
		Build()

	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(replyBody).
		Build()

	resp, err := u.cli.Im.Message.Reply(ctx, replyReq)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

// SendFeatureShutdownReminder Send feature shutdown reminder
func (u *UserMessage) SendFeatureShutdownReminder(ctx context.Context, messageId string) error {
	msg := "{\"text\": \"管理员关闭了这个功能，如果你需要开启，请联系管理员\"}"
	return u.SendReplyMessage(ctx, messageId, msg, larkim.MsgTypePost)
}

// SendFeatureUndevelopedReminder Send feature undeveloped reminder
func (u *UserMessage) SendFeatureUndevelopedReminder(ctx context.Context, messageId string) error {
	msg := TemplateUnopenedAbilityCard()
	return u.SendReplyMessage(ctx, messageId, msg, larkim.MsgTypeInteractive)
}
