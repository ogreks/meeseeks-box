package ioc

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"go.uber.org/zap"
)

func InitLarkClient(cfg configs.Config, logger *zap.Logger) *lark.Client {
	client := lark.NewClient(
		cfg.WebHook.Feishu.AppId,
		cfg.WebHook.Feishu.AppSecret,
		lark.WithLogLevel(larkcore.LogLevelDebug),
		lark.WithEnableTokenCache(true),
		lark.WithLogReqAtDebug(true),
	)

	return client
}

// 处理飞书消息
func InitLarkMessageDispatcher(log *zap.Logger, db orm.Repo, client *lark.Client) *feishu.UserMessage {
	return feishu.NewUserMessage(
		log,
		db,
		client,
		feishu.WithOnP2MessageReceiveV1(),
		feishu.WithOnP2MessageReadV1(),
		feishu.WithOnP2UserCreatedV3(),
	)
}
