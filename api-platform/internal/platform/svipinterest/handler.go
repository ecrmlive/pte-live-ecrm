// Package svipinterest manages the reusable C-end SVIP benefit catalogue.
// Plans store benefit-name snapshots, while new and edited plans are checked
// against this active catalogue to prevent dangling customer-facing benefits.
package svipinterest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler { return &Handler{business: business, admin: admin} }
func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.admin, "user.svip.interest.manage")
	r.GET("/svip/interests", access, manage, h.List)
	r.POST("/svip/interests", access, manage, h.Create)
	r.PUT("/svip/interests/:id", access, manage, h.Update)
	r.DELETE("/svip/interests/:id", access, manage, h.Delete)
}

type interest struct {
	ID          uint64 `gorm:"column:id" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
	IconURL     string `gorm:"column:icon_url" json:"icon_url"`
	Status      int    `gorm:"column:status" json:"status"`
	Sort        int    `gorm:"column:sort" json:"sort"`
	Version     uint64 `gorm:"column:version" json:"version"`
	PlanCount   int64  `gorm:"column:plan_count" json:"plan_count"`
}
type input struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"icon_url"`
	Status      int    `json:"status"`
	Sort        int    `json:"sort"`
	Version     uint64 `json:"version"`
}

func (h *Handler) List(c *gin.Context) {
	rows := make([]interest, 0)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest AS i").Select("i.id,i.name,i.description,i.icon_url,i.status,i.sort,i.version,0 AS plan_count").Where("i.deleted_at IS NULL").Order("i.sort ASC,i.id ASC").Scan(&rows).Error; err != nil {
		fail(c, "会员权益查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) Create(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || !valid(&in, false) {
		response.Fail(c, http.StatusBadRequest, "会员权益参数错误")
		return
	}
	row := fromInput(&in)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest").Create(map[string]any{"name": row.Name, "description": row.Description, "icon_url": row.IconURL, "status": row.Status, "sort": row.Sort, "version": 1}).Error; err != nil {
		conflictOrFail(c, err, "会员权益名称已存在", "会员权益创建失败")
		return
	}
	response.OK(c, gin.H{"name": row.Name, "version": 1})
}
func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in input
	if id == 0 || c.ShouldBindJSON(&in) != nil || !valid(&in, true) {
		response.Fail(c, http.StatusBadRequest, "会员权益参数错误")
		return
	}
	row := fromInput(&in)
	res := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest").Where("id=? AND deleted_at IS NULL AND version=?", id, in.Version).Updates(map[string]any{"name": row.Name, "description": row.Description, "icon_url": row.IconURL, "status": row.Status, "sort": row.Sort, "version": gorm.Expr("version + 1")})
	if res.Error != nil {
		conflictOrFail(c, res.Error, "会员权益名称已存在", "会员权益更新失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusConflict, "会员权益已被修改或不存在，请刷新后重试")
		return
	}
	response.OK(c, gin.H{"id": id, "version": in.Version + 1})
}
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "会员权益 ID 错误")
		return
	}
	var item interest
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest").Select("id,name").Where("id=? AND deleted_at IS NULL", id).Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "会员权益不存在")
		} else {
			fail(c, "会员权益删除校验失败")
		}
		return
	}
	var plans int64
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_plan").Where("JSON_CONTAINS(benefits, JSON_QUOTE(?)) AND status=1", item.Name).Count(&plans).Error; err != nil {
		fail(c, "会员权益删除校验失败")
		return
	}
	if plans > 0 {
		response.Fail(c, http.StatusConflict, "该权益仍被启用会员类型使用，请先调整或停用相关类型")
		return
	}
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_svip_interest").Where("id=? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("NOW()")).Error; err != nil {
		fail(c, "会员权益删除失败")
		return
	}
	response.OK(c, gin.H{"id": id})
}
func fromInput(in *input) interest {
	return interest{Name: strings.TrimSpace(in.Name), Description: strings.TrimSpace(in.Description), IconURL: strings.TrimSpace(in.IconURL), Status: in.Status, Sort: in.Sort}
}
func valid(in *input, update bool) bool {
	if in == nil || strings.TrimSpace(in.Name) == "" || len([]rune(strings.TrimSpace(in.Name))) > 64 || len([]rune(strings.TrimSpace(in.Description))) > 500 || (in.IconURL != "" && !strings.HasPrefix(in.IconURL, "/demo/") && !strings.HasPrefix(in.IconURL, "https://")) || (in.Status != 0 && in.Status != 1) || in.Sort < 0 || in.Sort > 999999 || (update && in.Version == 0) {
		return false
	}
	return true
}
func conflictOrFail(c *gin.Context, err error, conflict, internal string) {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		response.Fail(c, http.StatusConflict, conflict)
		return
	}
	fail(c, internal)
}
func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
