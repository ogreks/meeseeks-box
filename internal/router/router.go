package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/pkg/feishu"
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
	Lark              *middleware.Lark
	MessageDispatcher *feishu.UserMessage
}

func InitRouter(rh *RouterHandler) error {

	user.Register(rh.Engine, rh.DB, rh.Log, rh.AuthMiddleware)
	webhook.Register(rh.Engine, rh.Lark)

	return nil
}
