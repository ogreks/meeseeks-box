package feishu

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu"
)

// OnP2MessageReceiveV1 接收消息
func (u *UserMessage) OnP2MessageReceiveV1() feishu.OnP2MessageReceiveV1 {
	return func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
		fmt.Println("OnP2MessageReceiveV1")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())

		if !u.EnableOnP2MessageReceiveV1.Load() {
			return u.SendFeatureShutdownReminder(ctx, *event.Event.Message.MessageId)
		}

		eventMessage := event.EventV2Base.Header.EventType
		if condition, ok := u.EventDispatcher.Load(eventMessage); ok {
			return condition.(feishu.OnP2MessageReceiveV1)(ctx, event)
		}

		return u.SendFeatureUndevelopedReminder(ctx, *event.Event.Message.MessageId)
	}
}

// OnP2MessageReadV1 消息已读
func (u *UserMessage) OnP2MessageReadV1() feishu.OnP2MessageReadV1 {
	return func(ctx context.Context, event *larkim.P2MessageReadV1) error {
		fmt.Println("OnP2MessageReadV1")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())
		// 关闭服务提示
		// if !u.EnableOnP2MessageReadV1.Load() {
		// 	// return u.SendFeatureShutdownReminder(ctx, *event.Event.Reader.ReaderId.UserId)
		// }

		eventMessage := event.EventV2Base.Header.EventType
		if condition, ok := u.EventDispatcher.Load(eventMessage); ok {
			return condition.(feishu.OnP2MessageReadV1)(ctx, event)
		}

		// TODO 功能未开放提示
		return nil
	}
}

// OnP2UserCreatedV3 员工入职 事件通知
func (u *UserMessage) OnP2UserCreatedV3() feishu.OnP2UserCreatedV3 {
	return func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
		fmt.Println("OnP2UserCreatedV3")
		fmt.Println(larkcore.Prettify(event))
		fmt.Println(event.RequestId())

		// Todo 功能未开放提示
		return nil
	}
}
