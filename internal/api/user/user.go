package user

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"

	ujwt "github.com/ogreks/meeseeks-box/internal/pkg/middleware/auth"

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
// Router /api/user/login [post]
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
// @Router /api/user/register [post]
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

// Me user get profile
// @Router /api/user/me [get]
func (h *handler) Me(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	c, ok := claims.(*ujwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	user, err := h.service.GetUserByAccountAid(ctx.Request.Context(), c.Content.(string))
	if err != nil {
		if errors.Is(err, UserSvc.ErrorUserNotFound) || errors.Is(err, UserSvc.ErrorAccountNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "unauthorized",
				"data":    gin.H{},
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "server error",
			"data":    err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    user,
	})
}

// RefersToken user refers token timeout
// Router /api/user/refers/token [put]
func (h *handler) RefersToken(ctx *gin.Context) {
	claims, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	c, ok := claims.(*ujwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "bad request",
			"data":    gin.H{},
		})
		return
	}

	var (
		tk, rest string
		err      error
	)

	if otk := ctx.GetHeader(configs.GetConfig().Jwt.HeaderKey); otk != "" {
		segs := strings.Split(otk, " ")
		if len(segs) == 2 {
			tk, err = h.tokenManager.RefreshToken(ctx.Request.Context(), segs[1], &ujwt.UserClaims{
				StandardClaims: jwt.StandardClaims{
					Subject:   c.Content.(string),
					Audience:  "api",
					Issuer:    configs.GetConfig().Jwt.Issuer,
					ExpiresAt: time.Now().Add(time.Duration(configs.GetConfig().Jwt.Expire) * time.Second).Unix(),
				},
				Content: c.Content,
			}, time.Duration(configs.GetConfig().Jwt.Expire)*time.Second+(20*time.Second))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"code":    http.StatusBadRequest,
					"message": "bad request",
					"data":    err.Error(),
				})
				return
			}
		}
	}

	if orset := ctx.GetHeader(configs.GetConfig().Jwt.HeaderKey); orset != "" {
		rest, err = h.tokenManager.RefreshToken(ctx.Request.Context(), orset, &ujwt.UserClaims{
			StandardClaims: jwt.StandardClaims{
				Subject:   c.Content.(string),
				Audience:  "refresh",
				Issuer:    configs.GetConfig().Jwt.Issuer,
				ExpiresAt: time.Now().Add(time.Duration(configs.GetConfig().Jwt.RefreshTimeout) * time.Second).Unix(),
			},
			Content: c.Content,
		}, time.Duration(configs.GetConfig().Jwt.RefreshTimeout)*time.Second+(20*time.Second))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "bad request refresh",
				"data":    err.Error(),
			})
			return
		}
	}

	ctx.Request.Header.Set(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", tk))
	ctx.Header(configs.GetConfig().Jwt.HeaderKey, fmt.Sprintf("Bearer %s", tk))
	ctx.Request.Header.Set(configs.GetConfig().Jwt.RefersKey, rest)
	ctx.Header(configs.GetConfig().Jwt.RefersKey, rest)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"access_id":            c.Content,
			"token":                tk,
			"expire":               configs.GetConfig().Jwt.Expire - 5,
			"refresh_token":        rest,
			"refresh_token_expire": (time.Duration(configs.GetConfig().Jwt.Expire) * 10) - 5,
		},
	})
}

// Logout logout user token
// Router /api/user/logout [delete]
func (h *handler) Logout(ctx *gin.Context) {
	if tk := ctx.GetHeader(configs.GetConfig().Jwt.HeaderKey); tk != "" {
		_ = h.tokenManager.Store().Delete(tk)
		ctx.Header(configs.GetConfig().Jwt.HeaderKey, "")
		fmt.Printf("tk config delete header %s is exists: %v\n", tk, h.tokenManager.Store().Exists(tk))
	}

	if rtk := ctx.GetHeader(configs.GetConfig().Jwt.RefersKey); rtk != "" {
		_ = h.tokenManager.Store().Delete(rtk)
		ctx.Header(configs.GetConfig().Jwt.RefersKey, "")
		fmt.Printf("tk config refresh delete header %s is exists: %v\n", rtk, h.tokenManager.Store().Exists(rtk))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "You have logged out, looking forward to the next return of the king",
		"data":    gin.H{},
	})
}
