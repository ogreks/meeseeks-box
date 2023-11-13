package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/config"
)

// Login user login
// @Summary user login
// @Description user api login
// @Tags API.user
// @Accept application/json
// @Produce json
// @Param username Body json true "username/email/phone ..."
// @Param password Body json true "md5 for password"
// @Router /api/user/login [post]
// @Security Login
func (h *handler) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": config.GetConfig(),
	})
}
