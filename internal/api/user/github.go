package user

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/configs"
	userJwt "github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"
	UserSvc "github.com/ogreks/meeseeks-box/internal/service/user"
	"github.com/rs/xid"
	"go.uber.org/zap"
)

// LoginGITHub user login GitHub
// @Summary user login
// @Description GitHub account api login
// @Tags API.user
// @Accept application/json
// @Produce json
// @Param connect_account_id Body json true "account id"
// @Param token Body json true "GitHub token"
// @Param refresh_token Body json true "GitHub refresh token"
// @Param expire_at Body json true "GitHub token expire time"
// @Param id Body json true "GitHub id"
// @Param nickname Body json true "GitHub nickname"
// @Param name Body json true "GitHub name"
// @Param email Body json true "GitHub bind email"
// @Param options Body json true "GitHub account more json"
// @Router /api/user/github/login [post]
// @Security Login
func (h *handler) LoginGITHub(ctx *gin.Context) {
	type githubReq struct {
		Token        string `json:"token" binding:"required"`
		RefreshToken string `json:"refresh_token" binding:"required"`
		ExpireAt     int    `json:"expire_at" binding:"required"`

		ID       int    `json:"id" binding:"required"`
		NickName string `json:"nickname"`
		Name     string `json:"name"`
		Email    string `json:"email"`
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

	tk, err := h.tokenManager.CreateToken(ctx.Request.Context(), &userJwt.UserClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    configs.GetConfig().Jwt.Issuer,
			ExpiresAt: time.Now().Add(time.Duration(configs.GetConfig().Jwt.Expire) * time.Second).Unix(),
		},
		Content: ac.Aid,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})
		return
	}

	ctx.Request.Header.Set(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", tk))
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", tk))

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"access_id": ac.Aid,
			"token":     tk,
			"expire":    configs.GetConfig().Jwt.Expire - 5,
		},
	})
}
