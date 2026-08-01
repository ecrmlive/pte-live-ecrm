package setting

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/domain/merchant"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-merchant/internal/pkg/response"
)

type Handler struct {
	id  *identity.Service
	mer *merchant.Service
}

func NewHandler(id *identity.Service, mer *merchant.Service) *Handler {
	return &Handler{id: id, mer: mer}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/setting/shop", h.GetShop)
	r.PUT("/setting/shop", middleware.RequireMerchantMenu(h.id, identity.MerPermShopUpdate), h.UpdateShop)
	r.GET("/setting/staff", h.ListStaff)
	r.POST("/setting/staff", middleware.RequireMerchantMenu(h.id, identity.MerPermStaffWrite), h.CreateStaff)
	r.PUT("/setting/staff/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermStaffWrite), h.UpdateStaff)
	r.GET("/setting/admins", h.ListAdmins)
	r.POST("/setting/admins", middleware.RequireMerchantMenu(h.id, identity.MerPermAdminsWrite), h.CreateAdmin)
	r.PUT("/setting/admins/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermAdminsWrite), h.UpdateAdmin)
	r.GET("/setting/roles", h.ListRoles)
	r.POST("/setting/roles", middleware.RequireMerchantMenu(h.id, identity.MerPermRolesWrite), h.CreateRole)
	r.PUT("/setting/roles/:id", middleware.RequireMerchantMenu(h.id, identity.MerPermRolesWrite), h.UpdateRole)
	r.GET("/setting/menus/tree", h.MenuTree)
}

func (h *Handler) GetShop(c *gin.Context) {
	row, err := h.mer.GetMerchant(c.Request.Context(), middleware.MerID(c))
	if err != nil {
		writeMerErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateShop(c *gin.Context) {
	var in merchant.ShopProfileInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.mer.UpdateShopProfile(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeMerErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListStaff(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.id.ListStaff(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateStaff(c *gin.Context) {
	var in identity.StaffSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.id.CreateStaff(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateStaff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in identity.StaffSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.id.UpdateStaff(c.Request.Context(), middleware.MerID(c), uint(id), in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListAdmins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.id.ListMerchantAdmins(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateAdmin(c *gin.Context) {
	var in identity.AdminSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.id.CreateMerchantAdmin(c.Request.Context(), middleware.MerID(c), in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateAdmin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in identity.AdminSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.id.UpdateMerchantAdmin(c.Request.Context(), middleware.MerID(c), uint(id), middleware.AdminID(c), in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	res, err := h.id.ListMerchantRoles(c.Request.Context(), middleware.MerID(c), page, limit)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateRole(c *gin.Context) {
	var in identity.RoleSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.id.CreateRole(c.Request.Context(), middleware.MerID(c), 2, in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in identity.RoleSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	merID := middleware.MerID(c)
	row, err := h.id.UpdateRole(c.Request.Context(), uint(id), &merID, in)
	if err != nil {
		writeIDErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) MenuTree(c *gin.Context) {
	tree, err := h.id.MenuTree(c.Request.Context(), 2)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": tree})
}

func writeIDErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identity.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, identity.ErrBadParam), errors.Is(err, identity.ErrWeakPassword),
		errors.Is(err, identity.ErrAccountExists):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}

func writeMerErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, merchant.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, merchant.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
