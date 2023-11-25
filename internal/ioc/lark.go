package ioc

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/ogreks/meeseeks-box/configs"
	"go.uber.org/zap"
)

func NewLarkClient(cfg configs.Config, logger *zap.Logger) *lark.Client {
	client := lark.NewClient(
		cfg.WebHook.Feishu.AppId,
		cfg.WebHook.Feishu.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithEnableTokenCache(true),
		lark.WithLogReqAtDebug(true),
	)

	return client
}
