package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/config"
)

func (h *handler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": config.GetConfig(),
	})
}
