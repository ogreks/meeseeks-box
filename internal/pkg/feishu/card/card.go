package card

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/ogreks/meeseeks-box/internal/pkg/syncx"
	"go.uber.org/atomic"
)

type CardMessagerInterface interface {
	RegisterCardActionHandler(key string, handler interface{}) CardMessagerInterface
	RegisterRoute(path string, g *gin.RouterGroup)
}

type CardOption func(u *CardMessager)

func WithEncryptKey(encryptKey string) CardOption {
	return func(c *CardMessager) {
		c.EncryptKey.Store(encryptKey)
	}
}

func WithVerificationToken(verificationToken string) CardOption {
	return func(c *CardMessager) {
		c.VerificationToken.Store(verificationToken)
	}
}

type CardMessager struct {
	CardActionHandlers *syncx.Map[string, interface{}]

	cli               *lark.Client
	EncryptKey        *atomic.String
	VerificationToken *atomic.String
}

func NewCardMessager(client *lark.Client, options ...CardOption) *CardMessager {
	card := &CardMessager{
		CardActionHandlers: &syncx.Map[string, interface{}]{},

		cli:               client,
		EncryptKey:        atomic.NewString(""),
		VerificationToken: atomic.NewString(""),
	}

	for _, opt := range options {
		opt(card)
	}

	return card
}

func (c *CardMessager) RegisterCardActionHandler(key string, handler interface{}) CardMessagerInterface {
	c.CardActionHandlers.Store(key, handler)

	return c
}

func (c *CardMessager) RegisterRoute(path string, g *gin.RouterGroup) {
	g.POST(path, sdkginext.NewCardActionHandlerFunc(
		larkcard.NewCardActionHandler(
			c.VerificationToken.Load(),
			c.EncryptKey.Load(),
			c.Handler(),
		),
	))
}

func (c *CardMessager) Handler() func(ctx context.Context, action *larkcard.CardAction) (interface{}, error) {
	return func(ctx context.Context, action *larkcard.CardAction) (interface{}, error) {

		fmt.Println(larkcore.Prettify(action))

		return nil, nil
	}
}
