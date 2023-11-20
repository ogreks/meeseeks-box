package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/configs"
	"github.com/ogreks/meeseeks-box/internal/pkg/middleware"
	UserSvc "github.com/ogreks/meeseeks-box/internal/service/user"
	"github.com/rs/xid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const (
	Account = iota + 1
	Email
	Phone
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
	type login struct {
		Type     int    `json:"type" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"base64"`
	}

	var l login
	if err := ctx.ShouldBindJSON(&l); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	if l.Type < Account || l.Type > Phone {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	password, err := base64.StdEncoding.DecodeString(l.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	var (
		u *UserSvc.UserAccount
	)

	switch l.Type {
	case Account:
		u, err = h.service.GetUserByUserName(ctx.Request.Context(), l.Username, string(password))
	case Email:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "login type not support",
			"data":    gin.H{},
		})
		return
	case Phone:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "login type not support",
			"data":    gin.H{},
		})
		return
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "login type not support",
			"data":    gin.H{},
		})
		return
	}

	if errors.Is(err, UserSvc.ErrorAccountOrPassword) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	if errors.Is(err, UserSvc.ErrorAccountNotEnable) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})
		return
	}

	claims := middleware.NewGlobalJWT(u.Aid, time.Duration(configs.GetConfig().Jwt.Expire)*time.Second)
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
			"access_id": u.Aid,
			"token":     token,
			"expire":    claims.ExpiresAt,
		},
	})
}

// Register user register
// @Summary user register
// @Description user api register
// @Tags API.user
// @Accept application/json
// @Produce json
// @Param username Body json true "username/email/phone ..."
// @Param password Body json true "md5 for password"
// @Param register_type Body json true "register type 1: username 2: email 3: phone"
// @Router /api/user/register [post]
// @Security Register
func (h *handler) Register(ctx *gin.Context) {
	type register struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required,base64"`
		VerifyCode   string `json:"verify_code"`
		NickName     string `json:"nick_name"`
		RegisterType int    `json:"type" binding:"required"`
	}

	var r register
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    err.Error(),
		})
		return
	}

	if r.RegisterType < Account || r.RegisterType > Phone {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	password, err := base64.StdEncoding.DecodeString(r.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	switch r.RegisterType {
	case Account:
		err = h.service.CreateUserByUserName(ctx.Request.Context(), xid.New().String(), r.Username, string(hash))
	case Email:
		err = h.service.CreateUserByEmail(ctx.Request.Context(), xid.New().String(), r.Username, string(hash))
	case Phone:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "register type not support",
			"data":    gin.H{},
		})

		return
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})

		h.logger.Error("create user error", zap.Error(err))

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "register success",
		"data":    gin.H{},
	})
}

func (h *handler) Me(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    gin.H{},
	})
}
