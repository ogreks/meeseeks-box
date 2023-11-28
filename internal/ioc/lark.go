package ioc

import (
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/ogreks/meeseeks-box/configs"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
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
func InitLarkMessageDispatcher(cfg configs.Config, log *zap.Logger, db orm.Repo, client *lark.Client) feishuUserMessage.UserMessageInterface {
	return feishuUserMessage.NewUserMessage(
		log,
		db,
		client,
		feishuUserMessage.WithOnP2MessageReceiveV1(),
		feishuUserMessage.WithOnP2MessageReadV1(),
		feishuUserMessage.WithOnP2UserCreatedV3(),
		feishuUserMessage.WithEncryptKey(cfg.WebHook.Feishu.EncryptKey),
		feishuUserMessage.WithVerificationToken(cfg.WebHook.Feishu.VerificationToken),
	)
}
