// Package recharge governs recharge plans without changing recharge orders,
// payment callbacks, or member balances. Those facts belong to api-business.
package recharge

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ businessDB, adminDB *gorm.DB }

func NewHandler(businessDB, adminDB *gorm.DB) *Handler { return &Handler{businessDB, adminDB} }

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.adminDB, "marketing.balance.config.read")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.balance.config.manage")
	r.GET("/recharge/plans", access, read, h.ListPlans)
	r.POST("/recharge/plans", access, manage, h.CreatePlan)
	r.PUT("/recharge/plans/:id", access, manage, h.UpdatePlan)
	r.PUT("/recharge/plans/:id/status", access, manage, h.UpdatePlanStatus)
	r.DELETE("/recharge/plans/:id", access, manage, h.DeletePlan)
	r.GET("/recharge/orders", access, manage, h.ListOrders)
}

type plan struct {
	ID          uint64  `gorm:"column:id" json:"id"`
	Name        string  `gorm:"column:name" json:"name"`
	Amount      float64 `gorm:"column:amount" json:"amount"`
	BonusAmount float64 `gorm:"column:bonus_amount" json:"bonus_amount"`
	Status      int     `gorm:"column:status" json:"status"`
	Sort        int     `gorm:"column:sort" json:"sort"`
	Version     uint64  `gorm:"column:version" json:"version"`
	CreatedAt   string  `gorm:"column:created_at" json:"created_at"`
}

type input struct {
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	BonusAmount float64 `json:"bonus_amount"`
	Status      int     `json:"status"`
	Sort        int     `json:"sort"`
	Version     uint64  `json:"version"`
}

type order struct {
	ID          uint64  `gorm:"column:id" json:"id"`
	RechargeNo  string  `gorm:"column:recharge_no" json:"recharge_no"`
	UserID      uint64  `gorm:"column:user_id" json:"user_id"`
	Amount      float64 `gorm:"column:amount" json:"amount"`
	BonusAmount float64 `gorm:"column:bonus_amount" json:"bonus_amount"`
	Status      string  `gorm:"column:status" json:"status"`
	CreatedAt   string  `gorm:"column:created_at" json:"created_at"`
	PaidAt      *string `gorm:"column:paid_at" json:"paid_at,omitempty"`
}

func (h *Handler) ListPlans(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan")
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		status, err := strconv.Atoi(raw)
		if err != nil || (status != 0 && status != 1) {
			response.Fail(c, http.StatusBadRequest, "显示状态参数错误")
			return
		}
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "充值档位查询失败")
		return
	}
	var rows []plan
	if err := q.Select("id,name,amount,bonus_amount,status,sort,version,created_at").
		Order("sort ASC,id ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error; err != nil {
		fail(c, "充值档位查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || !valid(in, false) {
		response.Fail(c, http.StatusBadRequest, "充值档位参数错误")
		return
	}
	row := plan{
		Name:        planName(in),
		Amount:      in.Amount,
		BonusAmount: in.BonusAmount,
		Status:      in.Status,
		Sort:        in.Sort,
		Version:     1,
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan").Create(&row).Error; err != nil {
		fail(c, "充值档位创建失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in input
	if id == 0 || c.ShouldBindJSON(&in) != nil || in.Version == 0 || !valid(in, false) {
		response.Fail(c, http.StatusBadRequest, "充值档位参数错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan").
		Where("id=? AND version=?", id, in.Version).
		Updates(map[string]any{
			"name":         planName(in),
			"amount":       in.Amount,
			"bonus_amount": in.BonusAmount,
			"status":       in.Status,
			"sort":         in.Sort,
			"version":      gorm.Expr("version + 1"),
		})
	if res.Error != nil {
		fail(c, "充值档位更新失败")
		return
	}
	if res.RowsAffected == 0 {
		var total int64
		if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan").Where("id = ?", id).Count(&total).Error; err != nil {
			fail(c, "充值档位查询失败")
			return
		}
		if total == 0 {
			response.Fail(c, http.StatusNotFound, "充值档位不存在")
			return
		}
		response.Fail(c, http.StatusConflict, "充值档位已被其他操作更新，请刷新后重试")
		return
	}
	response.OK(c, gin.H{"id": id, "version": in.Version + 1})
}

func (h *Handler) UpdatePlanStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in struct {
		Status int `json:"status"`
	}
	if id == 0 || c.ShouldBindJSON(&in) != nil || (in.Status != 0 && in.Status != 1) {
		response.Fail(c, http.StatusBadRequest, "显示状态参数错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan").
		Where("id = ?", id).
		Updates(map[string]any{"status": in.Status, "version": gorm.Expr("version + 1")})
	if res.Error != nil {
		fail(c, "充值档位状态更新失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "充值档位不存在")
		return
	}
	response.OK(c, gin.H{"id": id, "status": in.Status})
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "充值档位参数错误")
		return
	}
	res := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_plan").Where("id = ?", id).Delete(nil)
	if res.Error != nil {
		fail(c, "充值档位删除失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "充值档位不存在")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) ListOrders(c *gin.Context) {
	page, limit := pagination(c)
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_recharge_order")
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		if raw != "pending" && raw != "paid" && raw != "closed" {
			response.Fail(c, http.StatusBadRequest, "充值订单状态参数错误")
			return
		}
		q = q.Where("status=?", raw)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		fail(c, "充值订单查询失败")
		return
	}
	var rows []order
	if err := q.Select("id,recharge_no,user_id,amount,bonus_amount,status,created_at,paid_at").Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		fail(c, "充值订单查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func planName(in input) string {
	name := strings.TrimSpace(in.Name)
	if name != "" {
		return name
	}
	if in.Amount == math.Trunc(in.Amount) {
		return fmt.Sprintf("%.0f元充值", in.Amount)
	}
	return fmt.Sprintf("%.2f元充值", in.Amount)
}

func valid(in input, requireName bool) bool {
	name := strings.TrimSpace(in.Name)
	if requireName && name == "" {
		return false
	}
	if len([]rune(name)) > 64 {
		return false
	}
	return validMoney(in.Amount, false) && validMoney(in.BonusAmount, true) &&
		(in.Status == 0 || in.Status == 1) && in.Sort >= 0 && in.Sort <= 999999
}

// validMoney ensures the API accepts only values that can be represented by
// the DECIMAL(12,2) columns. It deliberately validates before GORM writes so
// MySQL never rounds or truncates an operator's requested plan price.
func validMoney(value float64, allowZero bool) bool {
	if math.IsNaN(value) || math.IsInf(value, 0) || value > 1_000_000 {
		return false
	}
	if value < 0 || (!allowZero && value == 0) {
		return false
	}
	return math.Abs(value*100-math.Round(value*100)) < 0.000001
}

func pagination(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if p < 1 {
		p = 1
	}
	if l < 1 || l > 100 {
		l = 10
	}
	return p, l
}

func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
