package feishu

import (
	"context"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type OnP2MessageReceiveV1 func(ctx context.Context, event *larkim.P2MessageReceiveV1) error
type OnP2MessageReadV1 func(ctx context.Context, event *larkim.P2MessageReadV1) error
type OnP2UserCreatedV3 func(ctx context.Context, event *larkcontact.P2UserCreatedV3) error
type CardActionHandler func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
