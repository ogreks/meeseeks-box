package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/pkg/logger"
	"github.com/ogreks/meeseeks-box/pkg/utils"
	"go.uber.org/zap"
	"net"
	"net/http"
	"time"
)

type Api interface {
	Start(ctx context.Context) error
}

type api struct {
	address string
	port    string
}

func NewApi(address, port string) Api {
	return &api{
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

func (a *api) Start(ctx context.Context) error {
	l := initLogger()

	server := InitApiServer()

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
			Addr:    fmt.Sprintf("%s:%s", a.address, a.port),
			Handler: server,
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
