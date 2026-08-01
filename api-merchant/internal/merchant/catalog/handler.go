package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/catalog"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	svc *catalog.Service
	id  *identity.Service
}

func NewHandler(svc *catalog.Service, id *identity.Service) *Handler {
	return &Handler{svc: svc, id: id}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/product-categories", h.CategoryTree)
	r.GET("/brands", h.ListBrands)
	r.GET("/products", h.ListProducts)
	r.GET("/products/:id", h.GetProduct)
	r.POST("/products", middleware.RequireMerchantMenu(h.id, identity.MerPermProductCreate), h.CreateProduct)
	r.PUT("/products/:id", h.UpdateProduct)
	r.DELETE("/products/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermProductDelete), h.DeleteProduct)
	r.PUT("/products/:id/show", middleware.RequireMerchantMenu(h.id, identity.MerPermProductShow), h.SetShow)
	r.PUT("/products/:id/stock", middleware.RequireMerchantMenu(h.id, identity.MerPermProductStock), h.SetStock)
}

func (h *Handler) CategoryTree(c *gin.Context) {
	tree, err := h.svc.CategoryTree(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": tree})
}

func (h *Handler) ListBrands(c *gin.Context) {
	list, err := h.svc.ListBrands(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	res, err := h.svc.ListMerchantProducts(c.Request.Context(), middleware.MerID(c), statusPtr, c.Query("keyword"), page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	p, err := h.svc.GetMerchantProduct(c.Request.Context(), middleware.MerID(c), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var in catalog.ProductSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.CreateMerchantProduct(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) UpdateProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in catalog.ProductSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	p, err := h.svc.UpdateMerchantProduct(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) DeleteProduct(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteMerchantProduct(c.Request.Context(), middleware.MerID(c), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type showReq struct {
	IsShow *uint8 `json:"is_show" binding:"required"`
}

func (h *Handler) SetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req showReq
	if err := c.ShouldBindJSON(&req); err != nil || req.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetMerchantProductShow(c.Request.Context(), middleware.MerID(c), uint(id), *req.IsShow == 1); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type stockReq struct {
	Stock *uint `json:"stock" binding:"required"`
}

func (h *Handler) SetStock(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req stockReq
	if err := c.ShouldBindJSON(&req); err != nil || req.Stock == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetMerchantProductStock(c.Request.Context(), middleware.MerID(c), uint(id), *req.Stock); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, catalog.ErrNotFound), errors.Is(err, catalog.ErrNotOnSale):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, catalog.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, catalog.ErrNameRequired),
		errors.Is(err, catalog.ErrInvalidPrice),
		errors.Is(err, catalog.ErrCateRequired),
		errors.Is(err, catalog.ErrBadStatus):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
