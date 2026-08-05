// Package memberlevel maintains normal membership level configuration. It does
// not alter a user's current level; those high-impact changes stay in the
// separately audited user-list member-level command.
package memberlevel

import (
	"encoding/json"
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
	manage := []gin.HandlerFunc{middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.admin, "user.member.level.manage")}
	r.GET("/member-levels", manage[0], manage[1], h.List)
	r.POST("/member-levels", manage[0], manage[1], h.Create)
	r.PUT("/member-levels/:id", manage[0], manage[1], h.Update)
	r.DELETE("/member-levels/:id", manage[0], manage[1], h.Delete)
}

type level struct {
	ID       uint64 `gorm:"column:id" json:"id"`
	Name     string `gorm:"column:name" json:"name"`
	Rank     int    `gorm:"column:rank" json:"rank"`
	Rules    string `gorm:"column:rules" json:"rules"`
	Benefits string `gorm:"column:benefits" json:"benefits"`
	Status   int    `gorm:"column:status" json:"status"`
	Version  uint64 `gorm:"column:version" json:"version"`
	Assigned int64  `gorm:"column:assigned_count" json:"assigned_count"`
}
type input struct {
	Name     string `json:"name"`
	Rank     int    `json:"rank"`
	Rules    string `json:"rules"`
	Benefits string `json:"benefits"`
	Status   int    `json:"status"`
	Version  uint64 `json:"version"`
}

func (h *Handler) List(c *gin.Context) {
	rows := make([]level, 0)
	err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level AS l").
		Select("l.id,l.name,l.rank,l.rules,l.benefits,l.status,l.version,COUNT(a.user_id) AS assigned_count").
		Joins("LEFT JOIN qixi_crm_b_member_account AS a ON a.level_id=l.id").
		Where("l.deleted_at IS NULL").Group("l.id,l.name,l.rank,l.rules,l.benefits,l.status,l.version").Order("l.rank ASC,l.id ASC").Scan(&rows).Error
	if err != nil {
		fail(c, "会员等级查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) Create(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || !valid(&in, false) {
		response.Fail(c, http.StatusBadRequest, "会员等级参数错误")
		return
	}
	row := values(&in)
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").Create(map[string]any{"name": row.Name, "rank": row.Rank, "rules": row.Rules, "benefits": row.Benefits, "status": row.Status, "version": row.Version}).Error; err != nil {
		conflictOrFail(c, err, "会员等级排序等级已存在", "会员等级创建失败")
		return
	}
	response.OK(c, gin.H{"name": row.Name, "rank": row.Rank, "version": row.Version})
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in input
	if id == 0 || c.ShouldBindJSON(&in) != nil || !valid(&in, true) {
		response.Fail(c, http.StatusBadRequest, "会员等级参数错误")
		return
	}
	row := values(&in)
	changes := map[string]any{"name": row.Name, "rank": row.Rank, "rules": row.Rules, "benefits": row.Benefits, "status": row.Status, "version": gorm.Expr("version + 1")}
	res := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").Where("id=? AND deleted_at IS NULL AND version=?", id, in.Version).Updates(changes)
	if res.Error != nil {
		conflictOrFail(c, res.Error, "会员等级排序等级已存在", "会员等级更新失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusConflict, "会员等级已被修改或不存在，请刷新后重试")
		return
	}
	response.OK(c, gin.H{"id": id, "version": in.Version + 1})
}

// Delete is a logical hide. Existing assignments and immutable change logs
// retain their historical level snapshot, and cannot be orphaned by this API.
func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "会员等级 ID 错误")
		return
	}
	var references int64
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_account").Where("level_id=?", id).Count(&references).Error; err != nil {
		fail(c, "会员等级删除校验失败")
		return
	}
	if references > 0 {
		response.Fail(c, http.StatusConflict, "该等级仍有用户使用，不能删除")
		return
	}
	if err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level_log").Where("level_id=? OR previous_level_id=?", id, id).Count(&references).Error; err != nil {
		fail(c, "会员等级删除校验失败")
		return
	}
	if references > 0 {
		response.Fail(c, http.StatusConflict, "该等级已有历史变更记录，不能删除")
		return
	}
	res := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").Where("id=? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("NOW()"))
	if res.Error != nil {
		fail(c, "会员等级删除失败")
		return
	}
	if res.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "会员等级不存在")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func values(in *input) level {
	return level{Name: strings.TrimSpace(in.Name), Rank: in.Rank, Rules: canonicalJSON(in.Rules), Benefits: canonicalJSON(in.Benefits), Status: in.Status, Version: 1}
}
func valid(in *input, update bool) bool {
	if in == nil || strings.TrimSpace(in.Name) == "" || len([]rune(strings.TrimSpace(in.Name))) > 64 || in.Rank < 1 || in.Rank > 10000 || (in.Status != 0 && in.Status != 1) || (update && in.Version == 0) {
		return false
	}
	return validRules(in.Rules) && validBenefits(in.Benefits)
}
func validRules(raw string) bool {
	var value map[string]any
	return json.Unmarshal([]byte(raw), &value) == nil && len(value) > 0
}
func validBenefits(raw string) bool {
	var values []string
	if json.Unmarshal([]byte(raw), &values) != nil || len(values) == 0 || len(values) > 20 {
		return false
	}
	seen := map[string]struct{}{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || len([]rune(value)) > 64 {
			return false
		}
		if _, exists := seen[value]; exists {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}
func canonicalJSON(raw string) string {
	var value any
	_ = json.Unmarshal([]byte(raw), &value)
	result, _ := json.Marshal(value)
	return string(result)
}
func conflictOrFail(c *gin.Context, err error, conflict, internal string) {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		response.Fail(c, http.StatusConflict, conflict)
		return
	}
	fail(c, internal)
}
func fail(c *gin.Context, message string) { response.Fail(c, http.StatusInternalServerError, message) }
