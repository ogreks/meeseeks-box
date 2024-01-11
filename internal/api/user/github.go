package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	UserSvc "github.com/ogreks/meeseeks-box/internal/service/user"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// GITHub user login github
// @Summary user login
// @Description github account api login
// @Tags API.user
// @Accept application/json
// @Produce json
// @Param connect_account_id Body json true "account id"
// @Param token Body json true "github token"
// @Param refresh_token Body json true "github refresh token"
// @Param expire_at Body json true "github token expire time"
// @Param id Body json true "github id"
// @Param nickname Body json true "github nickname"
// @Param name Body json true "github name"
// @Param email Body json true "github bind email"
// @Param options Body json true "github account more json"
// @Router /api/user/github/login [post]
// @Security Login
func (h *handler) LoginGITHub(ctx *gin.Context) {
	type githubReq struct {
		Token        string `json:"token" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
		ExpireAt     int    `json:"expire_at" binding:"required"`

		ID       int    `json:"id" binding:"required"`
		NickName string `json:"nickname" binding:"required"`
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Options  string `json:"options" binding:"required"`
	}

	var r githubReq
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    err.Error(),
		})
		return
	}

	refreshTokenTimeAt := time.Now().Add(time.Duration(r.ExpireAt-10) * time.Second)
	ac, err := h.service.LoginUserByGITHub(ctx.Request.Context(), UserSvc.AccountPlatform{
		Aid:                  xid.New().String(),
		PlatformID:           1,
		AccountID:            fmt.Sprintf("%d", r.ID),
		Token:                r.Token,
		RefreshToken:         r.RefreshToken,
		RefreshTokenExpireAt: &refreshTokenTimeAt,
		UserName:             r.Name,
		NickName:             r.NickName,
		Email:                r.Email,
		MoreJson:             r.Options,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})

		h.logger.Error("create user error", zap.Error(err))

		return
	}

	claims := middleware.NewGlobalJWT(ac.Aid, time.Duration(configs.GetConfig().Jwt.Expire)*time.Second)
	token, err := claims.CreateToken()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})
		return
	}

	ctx.Request.Header.Set(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", token))
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"access_id": ac.Aid,
			"token":     token,
			"expire":    claims.ExpiresAt,
		},
	})
}
