package middleware

import (
	"context"
	"fmt"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type LarkConfig struct {
	AppId             string `json:"app_id"`
	Secret            string `json:"-"`
	EncryptKey        string `json:"-"`
	VerificationToken string `json:"-"`
}

type Lark struct {
	cfg    LarkConfig
	handle *dispatcher.EventDispatcher
}

func NewLarkMiddleware(cfg LarkConfig) *Lark {
	return &Lark{
		cfg: cfg,
		handle: dispatcher.NewEventDispatcher(cfg.VerificationToken, cfg.EncryptKey).
			OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}).
			OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}).
			OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}),
	}
}

func (l *Lark) GetHandle() *dispatcher.EventDispatcher {
	return l.handle
}
