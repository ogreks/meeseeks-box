package platform

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ogreks/meeseeks-box/internal/service"
	nofservice "github.com/ogreks/meeseeks-box/internal/service/notify"
	"net/http"
	"time"
)

type AgreementReq struct {
	AgreementNo string     `json:"agreement_no" form:"agreement_no"`
	Title       string     `json:"title" form:"title" binding:"required"`
	Type        int32      `json:"type" form:"type" binding:"required"`
	Content     string     `json:"content" form:"content" binding:"required"`
	Version     string     `json:"version" form:"version"`
	Status      uint32     `json:"status" form:"status" binding:"required"`
	PublishAt   *time.Time `json:"publish_at" form:"publish_at" time_format:"2006-01-02 15:04:05"`
	Limit       int        `json:"limit" form:"limit"`
	Page        int        `json:"page" form:"page"`
}

// CreateAgreement create agreement content
// Route /api/open/agreement [post]
func (h *handle) CreateAgreement(ctx *gin.Context) {
	var req AgreementReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	err := h.notify.SaveAgreement(ctx.Request.Context(), nofservice.AgreementReq{
		Title:     req.Title,
		Type:      req.Type,
		Content:   req.Content,
		Version:   req.Version,
		Status:    req.Status,
		PublishAt: req.PublishAt,
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
		"data":    gin.H{},
	})
}

// UpdateAgreement update agreement detail
// Route /api/open/agreement/:agreement_no
func (h *handle) UpdateAgreement(ctx *gin.Context) {
	agreementNo := ctx.Param("agreement_no")
	var req AgreementReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	err := h.notify.SaveAgreement(ctx.Request.Context(), nofservice.AgreementReq{
		AgreementNo: agreementNo,
		Title:       req.Title,
		Type:        req.Type,
		Content:     req.Content,
		Version:     req.Version,
		Status:      req.Status,
		PublishAt:   req.PublishAt,
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
		"data":    gin.H{},
	})
}

// DeleteAgreement delete agreement
// Route /api/open/agreement [delete]
func (h *handle) DeleteAgreement(ctx *gin.Context) {
	agreementNo := ctx.Param("agreement_no")
	err := h.notify.AgreementDelete(ctx.Request.Context(), agreementNo)
	if err != nil {
		if errors.Is(err, nofservice.ErrDeleteAgreement) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"meesage": err.Error(),
				"data":    gin.H{},
			})
			return
		}

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

// AgreementPaginate get agreement paginate list
// Route /api/open/agreements [get]
func (h *handle) AgreementPaginate(ctx *gin.Context) {
	var req AgreementReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	resp, err := h.notify.AgreementPaginate(ctx.Request.Context(), nofservice.AgreementReq{
		AgreementNo: req.AgreementNo,
		Title:       req.Title,
		Type:        req.Type,
		Content:     req.Content,
		Version:     req.Version,
		Status:      req.Status,
		PublishAt:   req.PublishAt,
		PaginateReq: service.PaginateReq{
			Limit: req.Limit,
			Page:  req.Page,
		},
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
