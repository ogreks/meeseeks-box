package platform

import (
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/service"
	"github.com/ogreks/meeseeks-box/internal/service/platform"
	"net/http"
	"strconv"
	"time"
)

type sessionKeyUriReq struct {
	SessionNo string `uri:"no" form:"session_no" binding:"required"`
}

type sessionKeyReq struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Type     uint32 `json:"type" form:"type"`
	IsEnable uint32 `json:"is_enable" form:"is_enable"`
	Remark   string `json:"remark" form:"remark"`
}

// CreateSessionKey create open api account
// Route /api/open/account [post]
func (h *handle) CreateSessionKey(ctx *gin.Context) {
	var r sessionKeyReq
	if err := ctx.ShouldBind(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    err.Error(),
		})
		return
	}

	err := h.service.SaveSessionKeys(
		ctx.Request.Context(),
		&platform.SessionKeysReq{
			Name:      r.Name,
			Type:      r.Type,
			IsEnabled: r.IsEnable,
			Remark:    r.Remark,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "system error",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

// UpdateSessionKey update open api account
// Route /api/open/account/:session_no [put]
func (h *handle) UpdateSessionKey(ctx *gin.Context) {
	var (
		sessionNo sessionKeyUriReq
		bodyReq   sessionKeyReq
	)
	if err := ctx.ShouldBindUri(&sessionNo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	err := h.service.SaveSessionKeys(
		ctx.Request.Context(),
		&platform.SessionKeysReq{
			SessionNo: sessionNo.SessionNo,
			Name:      bodyReq.Name,
			Type:      bodyReq.Type,
			IsEnabled: bodyReq.IsEnable,
			Remark:    bodyReq.Remark,
		},
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "system error",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    gin.H{},
	})
}

type setStatusReq struct {
	SessionNoSlice []string `json:"session_no"`
}

// SessionKeysSetStatus set status
// Route /api/open/account/set/:status/status
func (h *handle) SessionKeysSetStatus(ctx *gin.Context) {
	statusStr := ctx.Param("status")
	status, err := strconv.ParseUint(statusStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "status in 1/2",
			"data":    err.Error(),
		})
		return
	}
	var r setStatusReq
	if err := ctx.ShouldBindJSON(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	total, err := h.service.SetStatusSessionKeys(
		ctx.Request.Context(),
		uint32(status),
		r.SessionNoSlice...,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "system error",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data": gin.H{
			"total":   len(r.SessionNoSlice),
			"success": total,
			"fail":    len(r.SessionNoSlice) - int(total),
		},
	})
}

// SessionsKeysList get open api account manage
// Route /api/open/account [get]
func (h *handle) SessionsKeysList(ctx *gin.Context) {
	type req struct {
		SessionNo string     `form:"session_no"`
		Name      string     `form:"name"`
		Type      uint32     `form:"type"`
		AppID     string     `form:"app_id"`
		IsEnable  uint32     `form:"is_enable"`
		StartTime *time.Time `form:"start_time" time_format:"2006-01-02 15:04:05"`
		EndTime   *time.Time `form:"end_time" time_format:"2006-01-02 15:04:05"`
		Limit     int        `form:"limit"`
		Page      int        `form:"page"`
	}

	var r req
	if err := ctx.ShouldBindQuery(&r); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bad request",
			"data":    err.Error(),
		})
		return
	}

	resp, err := h.service.SessionsKeysPaginate(ctx.Request.Context(), platform.SessionKeysPaginateReq{
		PaginateReq: service.PaginateReq{
			Limit: r.Limit,
			Page:  r.Page,
		},
		SessionNo: r.SessionNo,
		Name:      r.Name,
		Type:      r.Type,
		AppID:     r.AppID,
		IsEnabled: r.IsEnable,
		StartTime: r.StartTime,
		EndTime:   r.EndTime,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "system error",
			"data":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "success",
		"data":    resp,
	})
}
