package content

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

type Handler struct {
	svc *content.Service
	id  *identity.Service
}

func NewHandler(svc *content.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/notices", h.List)
	r.POST("/notices", h.Create)
	r.PUT("/notices/:id", h.Update)
	r.DELETE("/notices/:id", h.Delete)
	r.GET("/agreements", h.ListAgreements)
	r.GET("/agreements/:key", h.GetAgreement)
	r.PUT("/agreements/:key", middleware.RequirePlatformMenu(h.id, identity.PlatPermAgreementUpdate), h.SaveAgreement)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListAdmin(c.Request.Context(), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	var in content.NoticeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	n, err := h.svc.Create(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, n)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in content.NoticeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	n, err := h.svc.Update(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, n)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListAgreements(c *gin.Context) {
	list, err := h.svc.ListAgreements(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) GetAgreement(c *gin.Context) {
	row, err := h.svc.GetAgreement(c.Request.Context(), c.Param("key"))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SaveAgreement(c *gin.Context) {
	var in content.AgreeSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SaveAgreement(c.Request.Context(), c.Param("key"), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, content.ErrNotFound), errors.Is(err, content.ErrAgreeNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, content.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
