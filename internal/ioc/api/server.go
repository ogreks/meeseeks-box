package api

import (
	"github.com/gin-contrib/cors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/ogreks/meeseeks-box/configs"
	feishuMessage "github.com/ogreks/meeseeks-box/internal/pkg/feishu/message"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	"github.com/ogreks/meeseeks-box/internal/pkg/token"
	"github.com/ogreks/meeseeks-box/internal/repository/orm"
	"github.com/ogreks/meeseeks-box/internal/router"
	"go.uber.org/zap"
)

func InitApiServer(
	db orm.Repo, // db
	logger *zap.Logger, // log
	middlewares []gin.HandlerFunc, // middleware
	client *lark.Client, // fei shu client
	msg feishuMessage.MessageHandleInterface, // fei shu message event
	tokenStore token.Store[string], // store token server
) *gin.Engine {
	g := gin.New()

	if configs.GetConfig().Server.Debug {
		InitTestRoute(g)
		gin.SetMode(gin.DebugMode)
	}

	InitServiceRoute(g)
	InitServerStatus(g)

	g.Use(middlewares...)

	_ = router.InitRouter(&router.Handler{
		Engine:            g,
		DB:                db,
		Log:               logger,
		MessageDispatcher: msg,
		TokenStore:        tokenStore,
	})

	return g
}

// InitTestRoute init debug register hello router
func InitTestRoute(g *gin.Engine) {
	g.GET("/hello", func(ctx *gin.Context) {
		netAddr, err := net.InterfaceAddrs()
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ipAdders []string
		for _, addr := range netAddr {
			ip, ok := addr.(*net.IPNet)
			if ok && !ip.IP.IsLoopback() {
				if ip.IP.To4() != nil {
					ipAdders = append(ipAdders, ip.IP.String())
				}
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "Hello friends from far away",
			"data": gin.H{
				"ips": ipAdders,
			},
		})
	})
}

// InitServerStatus register server run time
func InitServerStatus(g *gin.Engine) {
	registerTime := time.Now()
	g.GET("/status", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "success",
			"data": gin.H{
				"register_time": registerTime.Format(time.RFC3339),
				"run_time":      time.Now().Sub(registerTime).String(),
				"now_time":      time.Now().Format(time.RFC3339),
			},
		})
	})
}

// InitServiceRoute register router common service
// register 404 route not found
// register 405 route method not found
func InitServiceRoute(g *gin.Engine) {

	g.NoMethod(func(ctx *gin.Context) {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": http.StatusMethodNotAllowed,
			"msg":  "route method not allowed",
			"data": gin.H{},
		})
	})

	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "route not found",
			"data": gin.H{},
		})
	})
}

// InitMiddleware init middleware
func InitMiddleware(logger *zap.Logger, cfg configs.Config) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.Recovery(logger),
		middleware.Trace(),
		middleware.CorsWithConfig(cors.Config{
			AllowCredentials: true,
			AllowOrigins:     cfg.Server.CorsAllowOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders: []string{
				"Origin", "Content-Length", "Content-Type", "X-Trace-Id", "Date",
				"Authorization", "Authenticate", "Authorization_At", "Authenticate_At",
			},
			MaxAge: 12 * time.Hour,
		}),
	}
}
