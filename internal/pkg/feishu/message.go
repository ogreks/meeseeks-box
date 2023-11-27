package feishu

import (
	"context"
	"fmt"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/syncx"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type UserMessageOption func(u *UserMessage)

type UserMessage struct {
	EventDispatcher *syncx.Map[string, interface{}]

	cli *lark.Client

	EnableOnP2MessageReceiveV1 atomic.Bool
	EnableOnP2MessageReadV1    atomic.Bool
	EnableOnP2UserCreatedV3    atomic.Bool
}

func WithOnP2MessageReceiveV1(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}

		u.EnableOnP2MessageReceiveV1.Store(isEnabled)
	}
}

func WithOnP2MessageReadV1(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}

		u.EnableOnP2MessageReadV1.Store(isEnabled)
	}
}

func WithOnP2UserCreatedV3(enable ...bool) UserMessageOption {
	return func(u *UserMessage) {
		isEnabled := true
		if len(enable) >= 1 {
			isEnabled = enable[0]
		}
		u.EnableOnP2UserCreatedV3.Store(isEnabled)
	}
}

func NewUserMessage(log *zap.Logger, db orm.Repo, client *lark.Client, opts ...UserMessageOption) *UserMessage {
	u := &UserMessage{
		cli:             client,
		EventDispatcher: &syncx.Map[string, interface{}]{},
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

func (u *UserMessage) RegisterEvent(key string, handler interface{}) *UserMessage {
	u.EventDispatcher.Store(key, handler)
	return u
}

// SendCloseMessage send close message
func (u *UserMessage) SendOnP2MessageReceiveV1CloseMessage(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	replyBody := larkim.NewReplyMessageReqBodyBuilder().
		MsgType(larkim.MsgTypePost).
		Content("{\"text\": \"管理员关闭了这个功能，如果你需要开启，请联系管理员\"}").
		Build()

	replyReq := larkim.NewReplyMessageReqBuilder().
		MessageId(*event.Event.Message.MessageId).
		Body(replyBody).
		Build()

	resp, err := u.cli.Im.Message.Reply(ctx, replyReq)
	if err != nil {
		return err
	}

	fmt.Println(resp)
	return nil
}

func (u *UserMessage) OnP2MessageReceiveV1() OnP2MessageReceiveV1 {
	return func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println("OnP2MessageReceiveV1")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())

		if !u.EnableOnP2MessageReceiveV1.Load() {
			return u.SendOnP2MessageReceiveV1CloseMessage(ctx, event)
		}

		eventMessage := event.EventV2Base.Header.EventType
		if condition, ok := u.EventDispatcher.Load(eventMessage); ok {
			return condition.(OnP2MessageReceiveV1)(ctx, event)
		}

		return nil
	}
}

// OnP2MessageReadV1 接受机器人读取单聊消息
func (u *UserMessage) OnP2MessageReadV1() OnP2MessageReadV1 {
	return func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println("OnP2MessageReadV1")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		// 关闭服务提示
		// if u.EnableOnP2MessageReadV1.Load() == false {
		// 	// todo 关闭服务提示
		// }

		eventMessage := event.EventV2Base.Header.EventType
		if condition, ok := u.EventDispatcher.Load(eventMessage); ok {
			return condition.(OnP2MessageReadV1)(ctx, event)
		}

		// TODO 功能未开放提示
		return nil
	}
}

// OnP2UserCreatedV3 员工入职 事件通知
func (u *UserMessage) OnP2UserCreatedV3() OnP2UserCreatedV3 {
	return func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println("OnP2UserCreatedV3")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())

		// Todo 功能未开放提示
		return nil
	}
}
