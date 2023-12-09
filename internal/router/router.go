package router

import (
	"github.com/gin-gonic/gin"
	feishuCardMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/card"
	feishuUserMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/user"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router/user"
	"github.com/ogreks/meeseeks-box/internal/router/webhook"
	"go.uber.org/zap"
)

type RouterHandler struct {
	Engine *gin.Engine
	DB     orm.Repo
	Log    *zap.Logger

	AuthMiddleware    *middleware.JwtMiddleware
	MessageDispatcher feishuUserMessage.UserMessageInterface
	CardDispatcher    feishuCardMessage.CardMessagerInterface
}

func InitRouter(rh *RouterHandler) error {

	user.Register(rh.Engine, rh.DB, rh.Log, rh.AuthMiddleware)
	webhook.Register(rh.Engine, rh.MessageDispatcher, rh.CardDispatcher)

	return nil
}
