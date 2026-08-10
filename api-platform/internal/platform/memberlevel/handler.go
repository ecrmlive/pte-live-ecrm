// Package memberlevel maintains normal membership level configuration. It does
// not alter a user's current level; those high-impact changes stay in the
// separately audited user-list member-level command.
package memberlevel

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	menuMemberLevelRead   = "user.member.level.read"
	menuMemberLevelManage = "user.member.level.manage"
)

type Handler struct{ business, admin *gorm.DB }

func New(business, admin *gorm.DB) *Handler {
	return &Handler{business: business, admin: admin}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := middleware.RequireAdminRoles("platform", "operations")
	read := middleware.RequireAdminMenu(h.admin, menuMemberLevelRead)
	manage := middleware.RequireAdminMenu(h.admin, menuMemberLevelManage)
	r.GET("/member-levels", access, read, h.List)
	r.POST("/member-levels", access, manage, h.Create)
	r.PUT("/member-levels/:id", access, manage, h.Update)
	r.DELETE("/member-levels/:id", access, manage, h.Delete)
}

type levelRow struct {
	ID            uint64     `gorm:"column:id" json:"id"`
	Name          string     `gorm:"column:name" json:"name"`
	Rank          int        `gorm:"column:rank" json:"rank"`
	IconURL       string     `gorm:"column:icon_url" json:"icon_url"`
	Rules         string     `gorm:"column:rules" json:"-"`
	Benefits      string     `gorm:"column:benefits" json:"-"`
	Status        int        `gorm:"column:status" json:"status"`
	Version       uint64     `gorm:"column:version" json:"version"`
	AssignedCount int64      `gorm:"column:assigned_count" json:"user_count"`
	CreatedAt     *time.Time `gorm:"column:created_at" json:"-"`
}

type levelView struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
	IconURL     string `json:"icon_url"`
	GrowthValue int64  `json:"growth_value"`
	BgImage     string `json:"bg_image"`
	Status      int    `json:"status"`
	Version     uint64 `json:"version"`
	UserCount   int64  `json:"user_count"`
	CreatedAt   string `json:"created_at"`
}

type input struct {
	Name        string `json:"name"`
	Rank        int    `json:"rank"`
	IconURL     string `json:"icon_url"`
	GrowthValue int64  `json:"growth_value"`
	BgImage     string `json:"bg_image"`
	Status      int    `json:"status"`
	Version     uint64 `json:"version"`
}

type rulesPayload struct {
	Value       int64  `json:"value"`
	Image       string `json:"image"`
	Description string `json:"description"`
}

func (h *Handler) List(c *gin.Context) {
	rows := make([]levelRow, 0)
	err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level AS l").
		Select("l.id,l.name,l.rank,l.icon_url,l.rules,l.benefits,l.status,l.version,l.created_at,COUNT(a.user_id) AS assigned_count").
		Joins("LEFT JOIN qixi_crm_b_member_account AS a ON a.level_id=l.id").
		Where("l.deleted_at IS NULL").
		Group("l.id,l.name,l.rank,l.icon_url,l.rules,l.benefits,l.status,l.version,l.created_at").
		Order("l.rank ASC,l.id ASC").Scan(&rows).Error
	if err != nil {
		fail(c, "会员等级查询失败")
		return
	}
	list := make([]levelView, 0, len(rows))
	for _, row := range rows {
		list = append(list, toView(row))
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) Create(c *gin.Context) {
	var in input
	if c.ShouldBindJSON(&in) != nil || !valid(&in, false) {
		response.Fail(c, http.StatusBadRequest, "会员等级参数错误")
		return
	}
	payload := persistValues(&in)
	err := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").Create(map[string]any{
		"name":       payload.Name,
		"rank":       payload.Rank,
		"icon_url":   payload.IconURL,
		"rules":      payload.Rules,
		"benefits":   payload.Benefits,
		"status":     payload.Status,
		"version":    uint64(1),
		"created_at": time.Now(),
	}).Error
	if err != nil {
		conflictOrFail(c, err, "会员等级序号已存在", "会员等级创建失败")
		return
	}
	response.OK(c, gin.H{"name": payload.Name, "rank": payload.Rank})
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in input
	if id == 0 || c.ShouldBindJSON(&in) != nil || !valid(&in, true) {
		response.Fail(c, http.StatusBadRequest, "会员等级参数错误")
		return
	}
	payload := persistValues(&in)
	changes := map[string]any{
		"name":     payload.Name,
		"rank":     payload.Rank,
		"icon_url": payload.IconURL,
		"rules":    payload.Rules,
		"benefits": payload.Benefits,
		"status":   payload.Status,
		"version":  gorm.Expr("version + 1"),
	}
	res := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").
		Where("id=? AND deleted_at IS NULL AND version=?", id, in.Version).Updates(changes)
	if res.Error != nil {
		conflictOrFail(c, res.Error, "会员等级序号已存在", "会员等级更新失败")
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
	res := h.business.WithContext(c.Request.Context()).Table("qixi_crm_b_member_level").
		Where("id=? AND deleted_at IS NULL", id).Update("deleted_at", gorm.Expr("NOW()"))
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

func toView(row levelRow) levelView {
	rules := parseRules(row.Rules)
	created := ""
	if row.CreatedAt != nil && !row.CreatedAt.IsZero() {
		created = row.CreatedAt.In(time.Local).Format("2006-01-02 15:04:05")
	}
	return levelView{
		ID:          row.ID,
		Name:        row.Name,
		Rank:        row.Rank,
		IconURL:     strings.TrimSpace(row.IconURL),
		GrowthValue: rules.Value,
		BgImage:     strings.TrimSpace(rules.Image),
		Status:      row.Status,
		Version:     row.Version,
		UserCount:   row.AssignedCount,
		CreatedAt:   created,
	}
}

func persistValues(in *input) levelRow {
	rules := rulesPayload{
		Value:       in.GrowthValue,
		Image:       strings.TrimSpace(in.BgImage),
		Description: "累计成长值达到要求后自动升级",
	}
	raw, _ := json.Marshal(rules)
	benefits, _ := json.Marshal([]string{"会员专享活动"})
	return levelRow{
		Name:     strings.TrimSpace(in.Name),
		Rank:     in.Rank,
		IconURL:  strings.TrimSpace(in.IconURL),
		Rules:    string(raw),
		Benefits: string(benefits),
		Status:   in.Status,
	}
}

func valid(in *input, update bool) bool {
	if in == nil {
		return false
	}
	name := strings.TrimSpace(in.Name)
	if name == "" || len([]rune(name)) > 64 {
		return false
	}
	if in.Rank < 1 || in.Rank > 10000 {
		return false
	}
	if in.GrowthValue < 0 || in.GrowthValue > 100000000 {
		return false
	}
	if len([]rune(strings.TrimSpace(in.IconURL))) > 1024 {
		return false
	}
	if len([]rune(strings.TrimSpace(in.BgImage))) > 1024 {
		return false
	}
	if in.Status != 0 && in.Status != 1 {
		return false
	}
	if update && in.Version == 0 {
		return false
	}
	return true
}

func parseRules(raw string) rulesPayload {
	var rules rulesPayload
	_ = json.Unmarshal([]byte(raw), &rules)
	return rules
}

func conflictOrFail(c *gin.Context, err error, conflict, internal string) {
	if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
		response.Fail(c, http.StatusConflict, conflict)
		return
	}
	fail(c, internal)
}

func fail(c *gin.Context, message string) {
	response.Fail(c, http.StatusInternalServerError, message)
}
