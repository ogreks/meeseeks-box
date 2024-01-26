package router

import (
	"github.com/gin-gonic/gin"
	feishuMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/message"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router/open"
	"github.com/ogreks/meeseeks-box/internal/router/user"
	"github.com/ogreks/meeseeks-box/internal/router/webhook"
	"go.uber.org/zap"
)

type Handler struct {
	Engine *gin.Engine
	DB     orm.Repo
	Log    *zap.Logger

	MessageDispatcher feishuMessage.MessageHandleInterface
	TokenStore        token.Store[string]
}

func InitRouter(rh *Handler) error {

	user.Register(rh.Engine, rh.DB, rh.Log, rh.TokenStore)
	open.Register(rh.Engine, rh.DB, rh.Log)
	webhook.Register(rh.Engine, rh.MessageDispatcher)

	return nil
}
