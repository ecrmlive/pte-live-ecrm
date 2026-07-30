package merchant

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/merchant"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
)

type Handler struct {
	svc *merchant.Service
	id  *identity.Service
}

func NewHandler(svc *merchant.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/merchants", h.ListMerchants)
	r.GET("/merchants/:id", h.GetMerchant)
	r.PUT("/merchants/:id/status", h.SetMerchantStatus)

	r.GET("/merchant-intentions", h.ListIntentions)
	r.GET("/merchant-intentions/:id", h.GetIntention)
	r.POST("/merchant-intentions/:id/audit", h.AuditIntention)

	r.GET("/merchant-categories", h.ListCategories)
	r.POST("/merchant-categories", h.CreateCategory)
	r.PUT("/merchant-categories/:id", h.UpdateCategory)
	r.DELETE("/merchant-categories/:id", h.DeleteCategory)
}

func (h *Handler) ListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	regionIDs := merchantScope(c)
	res, err := h.svc.ListMerchants(c.Request.Context(), c.Query("keyword"), statusPtr, regionIDs, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetMerchant(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), merchantScope(c)); err != nil {
		writeErr(c, err)
		return
	}
	m, err := h.svc.GetMerchant(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, m)
}

type statusReq struct {
	Enabled bool `json:"enabled"`
}

func (h *Handler) SetMerchantStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), merchantScope(c)); err != nil {
		writeErr(c, err)
		return
	}
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetMerchantEnabled(c.Request.Context(), uint(id), req.Enabled); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// merchantScope 不再读取旧 qixi_m_admin_system_admin。平台角色拥有全量
// 监管视图；其余角色在正式 qixi_crm_a_data_scope 投影接入前默认无商户范围，
// 避免以“全量”作为不安全回退。区域角色另由 RestrictRegionConsole 拦截。
func merchantScope(c *gin.Context) []uint {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return []uint{}
	}
	for _, role := range claims.Roles {
		if role == "platform" {
			return nil
		}
	}
	return []uint{}
}

func (h *Handler) ListIntentions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	res, err := h.svc.ListIntentions(c.Request.Context(), c.Query("keyword"), statusPtr, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetIntention(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetIntention(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) AuditIntention(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req merchant.AuditIntentionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.svc.AuditIntention(c.Request.Context(), uint(id), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

type categoryReq struct {
	CategoryName   string  `json:"category_name"`
	CommissionRate float64 `json:"commission_rate"`
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateCategory(c.Request.Context(), req.CategoryName, req.CommissionRate)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req categoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.UpdateCategory(c.Request.Context(), uint(id), req.CategoryName, req.CommissionRate); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteCategory(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, merchant.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, merchant.ErrAlreadyAudited),
		errors.Is(err, merchant.ErrBadStatus),
		errors.Is(err, merchant.ErrRejectNeedMsg):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		if err != nil && err.Error() == "分类名称不能为空" {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
