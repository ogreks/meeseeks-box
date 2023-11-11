package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/config"
	"github.com/ogreks/meeseeks-box/pkg/logger"
	"github.com/ogreks/meeseeks-box/pkg/utils"
	"go.uber.org/zap"
)

type Api interface {
	Start(ctx context.Context) error
}

type server struct {
	address string
	port    int
}

func NewApi(address string, port int) Api {
	return &server{
		address: address,
		port:    port,
	}
}

func initLogger() *zap.Logger {
	l, err := logger.NewJsonLogger()
	if err != nil {
		panic(err)
	}

	defer zap.ReplaceGlobals(l)

	return l
}

func (a *server) Start(ctx context.Context) error {
	l := initLogger()

	server := InitApiServer()

	if config.GetConfig().GetServer().Debug {
		gin.SetMode(gin.DebugMode)
	}

	server.GET("/hello", func(ctx *gin.Context) {
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
			"msg":  "hello",
			"data": gin.H{
				"ips": ipAdders,
			},
		})
	})

	return utils.Run(ctx, l, func(ctx context.Context) (func(), error) {
		httpSvc := &http.Server{
			Addr:        fmt.Sprintf("%s:%d", a.address, a.port),
			Handler:     server,
			ReadTimeout: time.Second * time.Duration(config.GetConfig().GetServer().ReadTimeout),
		}

		go func() {
			if err := httpSvc.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				l.Error("Failed to listen http server", zap.Error(err))
			}
		}()

		return func() {
			ctx, canal := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
			defer canal()

			httpSvc.SetKeepAlivesEnabled(false)
			if err := httpSvc.Shutdown(ctx); err != nil {
				l.Error("Failed to shutdown http server", zap.Error(err))
			}
		}, nil
	})
}
