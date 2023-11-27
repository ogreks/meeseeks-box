package middleware

import (
	"context"
	"fmt"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
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
	cfg         LarkConfig
	eventHandle *dispatcher.EventDispatcher
	cardHandle  *larkcard.CardActionHandler
}

func NewLarkMiddleware(cfg LarkConfig) *Lark {
	return &Lark{
		cfg: cfg,
		eventHandle: dispatcher.NewEventDispatcher(cfg.VerificationToken, cfg.EncryptKey).
			OnP2MessageReceiveV1(func(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
				fmt.Println("OnP2MessageReceiveV1")
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}).
			OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
				fmt.Println("OnP2MessageReadV1")
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}).
			OnP2UserCreatedV3(func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error {
				fmt.Println("OnP2UserCreatedV3")
				fmt.Println(larkcore.Prettify(event))
				fmt.Println(event.RequestId())
				return nil
			}),
		cardHandle: larkcard.NewCardActionHandler(cfg.VerificationToken, cfg.Secret, func(ctx context.Context, ca *larkcard.CardAction) (interface{}, error) {
			fmt.Println("card action NewCardActionHandler")
			fmt.Println(larkcore.Prettify(ca))
			fmt.Println(ca.RequestId())
			return nil, nil
		}),
	}
}

func (l *Lark) GetHandle() *dispatcher.EventDispatcher {
	return l.eventHandle
}

func (l *Lark) GetCardHandle() *larkcard.CardActionHandler {
	return l.cardHandle
}
