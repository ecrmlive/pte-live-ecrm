package circle

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	svc        *circle.Service
	adminDB    *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(svc *circle.Service, adminDB, businessDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB, businessDB: businessDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	zoneManage := middleware.RequireAdminMenu(h.adminDB, "region.zone.manage")
	agentManage := middleware.RequireAdminMenu(h.adminDB, "region.agent.manage")
	agentReview := middleware.RequireAdminMenu(h.adminDB, "region.agent.review")
	agentPasswordReset := middleware.RequireAdminMenu(h.adminDB, "region.agent.password.reset")
	r.GET("/business-zones", platform, zoneManage, h.ListCircles)
	r.POST("/business-zones", platform, zoneManage, h.CreateCircle)
	r.GET("/business-zones/options", platform, h.ZoneOptions)
	r.GET("/business-zones/:id/invite", platform, zoneManage, h.InviteCircle)
	r.POST("/business-zones/agents", platform, zoneManage, h.CreateRegionAgent)
	r.GET("/business-zones/:id", platform, zoneManage, h.GetCircle)
	r.PUT("/business-zones/:id", platform, zoneManage, h.UpdateCircle)
	r.PUT("/business-zones/:id/status", platform, zoneManage, h.UpdateCircleStatus)
	r.DELETE("/business-zones/:id", platform, zoneManage, h.DeleteCircle)
	r.GET("/business-zone-agents", platform, agentManage, h.ListAgents)
	// 区域表单选代理人：平台角色即可读选项（与 ZoneOptions 一致）。
	r.GET("/business-zone-agents/options", platform, h.AgentOptions)
	r.GET("/business-zone-agents/settings", platform, agentManage, h.AgentSettings)
	r.GET("/business-zone-agents/:id", platform, agentManage, h.GetAgent)
	r.GET("/business-zone-agents/:id/merchants", platform, agentManage, h.AgentMerchants)
	r.POST("/business-zone-agents", platform, agentManage, h.CreateAgent)
	r.PUT("/business-zone-agents/:id", platform, agentManage, h.UpdateAgent)
	r.POST("/business-zone-agents/:id/audit", platform, agentReview, h.AuditAgent)
	r.POST("/business-zone-agents/:id/password", platform, agentManage, agentPasswordReset, h.ResetAgentPassword)
	r.DELETE("/business-zone-agents/:id", platform, agentManage, h.RevokeAgent)
}

func page(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	l, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	return p, l
}
func optionalStatus(c *gin.Context) (*int8, error) {
	raw := c.Query("status")
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(raw, 10, 8)
	if err != nil {
		return nil, circle.ErrBadParam
	}
	v := int8(n)
	return &v, nil
}
func id(c *gin.Context) (uint, error) {
	n, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || n == 0 {
		return 0, circle.ErrBadParam
	}
	return uint(n), nil
}

func (h *Handler) ListCircles(c *gin.Context) {
	status, err := optionalStatus(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	filter := circle.CircleListFilter{
		Keyword: firstNonEmpty(c.Query("keyword"), c.Query("name")),
		Status:  status,
	}
	if raw := strings.TrimSpace(c.Query("type")); raw != "" {
		n, parseErr := strconv.ParseInt(raw, 10, 8)
		if parseErr != nil {
			writeErr(c, circle.ErrBadParam)
			return
		}
		v := int8(n)
		filter.Type = &v
	}
	if raw := strings.TrimSpace(c.Query("circle_agent_id")); raw != "" {
		n, parseErr := strconv.ParseUint(raw, 10, 64)
		if parseErr != nil {
			writeErr(c, circle.ErrBadParam)
			return
		}
		filter.CircleAgentID = uint(n)
	}
	if raw := strings.TrimSpace(c.Query("pid")); raw != "" {
		n, parseErr := strconv.ParseUint(raw, 10, 64)
		if parseErr != nil {
			writeErr(c, circle.ErrBadParam)
			return
		}
		pid := uint(n)
		filter.PID = &pid
	}
	p, l := page(c)
	out, err := h.svc.ListCircles(c.Request.Context(), filter, p, l)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.enrichCirclePage(c, out))
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

type zoneOptionNode struct {
	Value    uint             `json:"value"`
	Label    string           `json:"label"`
	Children []zoneOptionNode `json:"children,omitempty"`
}

// ZoneOptions 店铺表单「所属商户 / 店铺区域」级联选项；type=1 商户，type=0 区域。
func (h *Handler) ZoneOptions(c *gin.Context) {
	typeRaw := strings.TrimSpace(c.Query("type"))
	type qrow struct {
		CircleID uint   `gorm:"column:circle_id"`
		PID      uint   `gorm:"column:pid"`
		Name     string `gorm:"column:name"`
	}
	q := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_business_zone").
		Select("circle_id, pid, name").
		Where("status = 1").
		Order("sort ASC, circle_id ASC")
	if typeRaw != "" {
		n, err := strconv.ParseInt(typeRaw, 10, 8)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "参数错误")
			return
		}
		q = q.Where("type = ?", int8(n))
	}
	var rows []qrow
	if err := q.Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	byPID := map[uint][]qrow{}
	for _, row := range rows {
		byPID[row.PID] = append(byPID[row.PID], row)
	}
	var build func(pid uint) []zoneOptionNode
	build = func(pid uint) []zoneOptionNode {
		children := byPID[pid]
		out := make([]zoneOptionNode, 0, len(children))
		for _, row := range children {
			node := zoneOptionNode{Value: row.CircleID, Label: row.Name, Children: build(row.CircleID)}
			out = append(out, node)
		}
		return out
	}
	response.OK(c, gin.H{"list": build(0)})
}

func (h *Handler) GetCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	row, err := h.svc.GetCircle(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.circleDetail(c, row))
}
func (h *Handler) CreateCircle(c *gin.Context) {
	var in circle.CircleInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.CreateCircle(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	if err = h.syncCircleMerchants(c, row.CircleID, in.MerchantIDs); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.circleDetail(c, row))
}
func (h *Handler) UpdateCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.CircleInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.UpdateCircle(c.Request.Context(), key, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	if err = h.syncCircleMerchants(c, row.CircleID, in.MerchantIDs); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.circleDetail(c, row))
}

func (h *Handler) UpdateCircleStatus(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var body struct {
		Status int8 `json:"status"`
	}
	if c.ShouldBindJSON(&body) != nil || (body.Status != 0 && body.Status != 1) {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.GetCircle(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	in := circle.CircleInput{
		PID: row.PID, Name: row.Name, CircleAgentID: row.CircleAgentID,
		CommissionType: row.CommissionType, CommissionRate: row.CommissionRate,
		Remark: row.Remark, Sort: row.Sort, Status: body.Status, Type: row.Type,
		RoleID: row.RoleID, BusinessStoreCategory: row.BusinessStoreCategory,
		BusinessStoreType: row.BusinessStoreType,
		GoodsTypes:          circle.ParseIntCSV(row.GoodsType),
		PlatformCategoryIDs: circle.ParseUintCSV(row.PlatformCategoryIDs),
	}
	updated, err := h.svc.UpdateCircle(c.Request.Context(), key, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, updated)
}
func (h *Handler) DeleteCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	if err = h.svc.DeleteCircle(c.Request.Context(), key); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

// InviteCircle 返回 H5 邀请入驻 URL（扫码进入商户入驻页并带上区域 ID）。
func (h *Handler) InviteCircle(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	row, err := h.svc.GetCircle(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	siteURL := h.loadMallSiteURL(c)
	h5URL := buildRegionInviteH5URL(siteURL, row.CircleID)
	response.OK(c, gin.H{
		"circle_id": row.CircleID,
		"name":      row.Name,
		"site_url":  siteURL,
		"h5_url":    h5URL,
		"label":     "H5邀请入驻",
	})
}

func (h *Handler) loadMallSiteURL(c *gin.Context) string {
	var raw string
	_ = h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_setting_cache").
		Select("result").
		Where("`key` = ?", "mall_shop_config").
		Limit(1).
		Scan(&raw).Error
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	var cfg struct {
		SiteURL string `json:"site_url"`
	}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return ""
	}
	return strings.TrimRight(strings.TrimSpace(cfg.SiteURL), "/")
}

func buildRegionInviteH5URL(siteURL string, circleID uint) string {
	path := "/pages/merchant/apply/index?region_id=" + strconv.FormatUint(uint64(circleID), 10)
	if siteURL == "" {
		return path
	}
	return siteURL + path
}
func (h *Handler) ListAgents(c *gin.Context) {
	status, err := optionalStatus(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var agentType *int8
	if raw := strings.TrimSpace(c.Query("type")); raw != "" {
		n, parseErr := strconv.ParseInt(raw, 10, 8)
		if parseErr != nil {
			writeErr(c, circle.ErrBadParam)
			return
		}
		v := int8(n)
		agentType = &v
	}
	p, l := page(c)
	filter := circle.AgentListFilter{
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		Name:     strings.TrimSpace(c.Query("name")),
		Phone:    strings.TrimSpace(c.Query("phone")),
		Status:   status,
		Type:     agentType,
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
	}
	// 用户搜索三选一：uid | user_phone | nickname（仅生效当前类型，互斥）。
	uidRaw := strings.TrimSpace(c.Query("uid"))
	userPhone := strings.TrimSpace(c.Query("user_phone"))
	nickname := strings.TrimSpace(c.Query("nickname"))
	switch {
	case uidRaw != "":
		n, parseErr := strconv.ParseUint(uidRaw, 10, 64)
		if parseErr != nil || n == 0 {
			writeErr(c, circle.ErrBadParam)
			return
		}
		v := uint(n)
		filter.UID = &v
	case userPhone != "", nickname != "":
		if err = h.applyAgentUserSearch(c.Request.Context(), userPhone, nickname, &filter); err != nil {
			writeErr(c, err)
			return
		}
	}
	// 兼容旧前端：仅传 name 时按姓名模糊搜；keyword 仍覆盖姓名/手机/商户名。
	out, err := h.svc.ListAgents(c.Request.Context(), filter, p, l)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.enrichAgentPage(c, out))
}

// applyAgentUserSearch 按 C 端用户手机号或昵称预解析 uid 列表（二者互斥，手机号优先）。
func (h *Handler) applyAgentUserSearch(ctx context.Context, userPhone, nickname string, filter *circle.AgentListFilter) error {
	userPhone = strings.TrimSpace(userPhone)
	nickname = strings.TrimSpace(nickname)
	if userPhone == "" && nickname == "" {
		return nil
	}
	if h.businessDB == nil {
		empty := []uint{}
		filter.UIDs = &empty
		return nil
	}
	q := h.businessDB.WithContext(ctx).Table("qixi_crm_b_user").Select("id")
	if userPhone != "" {
		q = q.Where("mobile LIKE ?", "%"+userPhone+"%")
	} else {
		q = q.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	var ids []uint
	if err := q.Limit(500).Pluck("id", &ids).Error; err != nil {
		return err
	}
	if ids == nil {
		ids = []uint{}
	}
	filter.UIDs = &ids
	return nil
}
func (h *Handler) GetAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	row, err := h.svc.GetAgent(c.Request.Context(), key)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.enrichAgentItem(c, row))
}

func (h *Handler) AgentOptions(c *gin.Context) {
	approved := int8(circle.AgentApproved)
	var agentType *int8
	if raw := strings.TrimSpace(c.Query("type")); raw != "" {
		n, parseErr := strconv.ParseInt(raw, 10, 8)
		if parseErr != nil {
			writeErr(c, circle.ErrBadParam)
			return
		}
		v := int8(n)
		agentType = &v
	}
	page, err := h.svc.ListAgents(c.Request.Context(), circle.AgentListFilter{
		Status: &approved, Type: agentType,
	}, 1, 100)
	if err != nil {
		writeErr(c, err)
		return
	}
	options := make([]gin.H, 0, len(page.List))
	for _, item := range page.List {
		options = append(options, gin.H{"circle_agent_id": item.CircleAgentID, "name": item.Name, "phone": item.Phone, "type": item.Type})
	}
	response.OK(c, gin.H{"list": options})
}

// AgentSettings 对齐 CRMEB 仅只读的“代理设置”入口：规则来自服务端实际约束，
// 不把支付资料、结算账号或外部通道参数暴露给页面。
func (h *Handler) AgentSettings(c *gin.Context) {
	var rows []struct {
		Status int8  `gorm:"column:status"`
		Total  int64 `gorm:"column:total"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_business_zone_agent").Select("status, COUNT(*) AS total").Group("status").Scan(&rows).Error; err != nil {
		writeErr(c, err)
		return
	}
	counts := gin.H{"pending": 0, "approved": 0, "rejected": 0, "revoked": 0}
	for _, row := range rows {
		switch row.Status {
		case circle.AgentPending:
			counts["pending"] = row.Total
		case circle.AgentApproved:
			counts["approved"] = row.Total
		case circle.AgentRejected:
			counts["rejected"] = row.Total
		case circle.AgentRevoked:
			counts["revoked"] = row.Total
		}
	}
	response.OK(c, gin.H{
		"status_counts": counts,
		"review": gin.H{
			"platform_review_required":  true,
			"rejection_reason_required": true,
		},
		"security": gin.H{
			"payment_credentials_write_only": true,
			"admin_binding_required":         true,
			"password_min_length":            12,
			"password_max_length":            72,
		},
		"revocation": gin.H{
			"hard_delete":  false,
			"blocked_when": []string{"未通过审核", "已关联区域", "佣金余额非零", "仍关联统一后台账号"},
		},
	})
}

func (h *Handler) AgentMerchants(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	if _, err = h.svc.GetAgent(c.Request.Context(), key); err != nil {
		writeErr(c, err)
		return
	}
	rows := make([]struct {
		MerchantID   uint64 `gorm:"column:merchant_id" json:"merchant_id"`
		MerchantName string `gorm:"column:merchant_name" json:"merchant_name"`
		RegionID     uint64 `gorm:"column:region_id" json:"region_id"`
		Status       int    `gorm:"column:status" json:"status"`
	}, 0)
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_business_zone AS z").Select("m.merchant_id,m.merchant_name,m.region_id,m.status").Joins("JOIN qixi_crm_a_merchant_view AS m ON m.region_id=z.circle_id").Where("z.circle_agent_id=?", key).Order("m.merchant_id ASC").Scan(&rows).Error; err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"list": rows})
}
func (h *Handler) CreateAgent(c *gin.Context) {
	var in circle.AgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	h.createAgentWithInput(c, in)
}

// CreateRegionAgent 区域列表表单内「添加代理人」：强制 type=0 且立即通过。
func (h *Handler) CreateRegionAgent(c *gin.Context) {
	var in circle.AgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	in.Type = 0
	in.AutoApprove = true
	h.createAgentWithInput(c, in)
}

func (h *Handler) createAgentWithInput(c *gin.Context, in circle.AgentInput) {
	account, password := strings.TrimSpace(in.Account), strings.TrimSpace(in.Password)
	// 先校验登录账号与昵称，避免代理行已写入后因冲突留下脏数据。
	if account != "" && password != "" {
		if err := h.ensureAdminAccountAvailable(c, account); err != nil {
			writeErr(c, err)
			return
		}
		display := strings.TrimSpace(in.Name)
		if display == "" {
			display = account
		}
		if err := h.ensureAdminDisplayNameAvailable(c, display, 0); err != nil {
			writeErr(c, err)
			return
		}
	}
	row, err := h.svc.CreateAgent(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	if account != "" && password != "" {
		if err = h.bindAgentAdminAccount(c, row.CircleAgentID, account, password, row.Name, row.Phone); err != nil {
			writeErr(c, err)
			return
		}
	}
	if err = h.syncAgentCircles(c, row.CircleAgentID, in.CircleIDs); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.enrichAgentItem(c, row))
}
func (h *Handler) UpdateAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.AgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	row, err := h.svc.UpdateAgent(c.Request.Context(), key, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	if err = h.syncBoundAdminProfile(c, key, row.Name, row.Phone, strings.TrimSpace(in.Password)); err != nil {
		writeErr(c, err)
		return
	}
	if err = h.syncAgentCircles(c, key, in.CircleIDs); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, h.enrichAgentItem(c, row))
}
func (h *Handler) AuditAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in circle.AuditInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	if err = h.svc.AuditAgent(c.Request.Context(), key, in, middleware.AdminID(c)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

type revokeAgentInput struct {
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *Handler) RevokeAgent(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in revokeAgentInput
	if c.ShouldBindJSON(&in) != nil {
		writeErr(c, circle.ErrBadParam)
		return
	}
	replayed, err := h.svc.RevokeAgent(c.Request.Context(), key, in.Reason, in.IdempotencyKey, middleware.AdminID(c))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"circle_agent_id": key, "status": circle.AgentRevoked, "replayed": replayed})
}

type resetAgentPasswordInput struct {
	Password       string `json:"password"`
	Reason         string `json:"reason"`
	IdempotencyKey string `json:"idempotency_key"`
}

func validAgentPasswordReset(in *resetAgentPasswordInput) bool {
	if in == nil {
		return false
	}
	in.Reason, in.IdempotencyKey = strings.TrimSpace(in.Reason), strings.TrimSpace(in.IdempotencyKey)
	return len(in.Password) >= 12 && len(in.Password) <= 72 && len([]rune(in.Reason)) >= 2 && len([]rune(in.Reason)) <= 500 && len([]rune(in.IdempotencyKey)) >= 8 && len([]rune(in.IdempotencyKey)) <= 128
}

// ResetAgentPassword 仅重置已关联且启用的统一后台区域账号。密码只在本次请求中
// 参与 bcrypt 计算，审计表不保存密码、散列或可用于离线猜测的派生值。
func (h *Handler) ResetAgentPassword(c *gin.Context) {
	key, err := id(c)
	if err != nil {
		writeErr(c, err)
		return
	}
	var in resetAgentPasswordInput
	if c.ShouldBindJSON(&in) != nil || !validAgentPasswordReset(&in) {
		writeErr(c, circle.ErrBadParam)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "密码处理失败")
		return
	}
	err = h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var agent struct {
			Status int8 `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_a_business_zone_agent").Clauses(clause.Locking{Strength: "UPDATE"}).Select("status").Where("circle_agent_id=?", key).Take(&agent).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return circle.ErrNotFound
			}
			return err
		}
		if err := tx.Table("qixi_crm_a_business_zone_agent_password_reset_audit").Where("circle_agent_id=? AND idempotency_key=?", key, in.IdempotencyKey).Take(&struct{}{}).Error; err == nil {
			// 审计刻意不保存口令或其派生摘要，无法安全识别“同一口令”的重试；
			// 因而该高风险命令采用单次提交语义，避免把不同新口令伪装成重放成功。
			return circle.ErrCommandConflict
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if agent.Status != circle.AgentApproved {
			return circle.ErrAgentNotApproved
		}
		var admin struct {
			ID     uint64 `gorm:"column:id"`
			Status int8   `gorm:"column:status"`
		}
		if err := tx.Table("qixi_crm_a_admin_user").Clauses(clause.Locking{Strength: "UPDATE"}).Select("id,status").Where("circle_agent_id=?", key).Take(&admin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return circle.ErrAgentAdminUnbound
			}
			return err
		}
		if admin.Status != 1 {
			return circle.ErrAgentAdminUnbound
		}
		if result := tx.Table("qixi_crm_a_admin_user").Where("id=? AND status=1", admin.ID).Updates(map[string]any{"password_hash": string(hash), "auth_version": gorm.Expr("auth_version + 1")}); result.Error != nil {
			return result.Error
		} else if result.RowsAffected != 1 {
			return circle.ErrAgentAdminUnbound
		}
		return tx.Table("qixi_crm_a_business_zone_agent_password_reset_audit").Create(map[string]any{
			"circle_agent_id": key, "admin_user_id": admin.ID, "reason": in.Reason,
			"operator_admin_id": middleware.AdminID(c), "idempotency_key": in.IdempotencyKey,
		}).Error
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"circle_agent_id": key, "replayed": false})
}

func (h *Handler) enrichCirclePage(c *gin.Context, page *circle.PageResult[circle.Circle]) gin.H {
	list := make([]gin.H, 0, len(page.List))
	for i := range page.List {
		list = append(list, h.circleListItem(c, &page.List[i]))
	}
	return gin.H{"list": list, "total": page.Total, "page": page.Page, "limit": page.Limit}
}

func (h *Handler) circleListItem(c *gin.Context, row *circle.Circle) gin.H {
	hasChild := false
	if row.Level < 2 {
		if n, err := h.svc.CountCircleChildren(c.Request.Context(), row.CircleID); err == nil {
			hasChild = n > 0
		}
	}
	item := gin.H{
		"circle_id": row.CircleID, "pid": row.PID, "path": row.Path, "name": row.Name,
		"circle_agent_id": row.CircleAgentID, "commission_type": row.CommissionType,
		"commission_rate": row.CommissionRate, "level": row.Level, "remark": row.Remark,
		"sort": row.Sort, "status": row.Status, "type": row.Type, "role_id": row.RoleID,
		"business_store_category": row.BusinessStoreCategory, "business_store_type": row.BusinessStoreType,
		"goods_type": circle.ParseIntCSV(row.GoodsType),
		"platform_category_ids": circle.ParseUintCSV(row.PlatformCategoryIDs),
		"create_time": row.CreateTime, "merchant_count": h.countCircleMerchants(c, row.CircleID),
		"has_child": hasChild,
	}
	if row.CircleAgentID > 0 {
		if agent, err := h.svc.GetAgent(c.Request.Context(), row.CircleAgentID); err == nil {
			item["circle_agent"] = gin.H{"circle_agent_id": agent.CircleAgentID, "name": agent.Name, "phone": agent.Phone}
		}
	}
	return item
}

func (h *Handler) circleDetail(c *gin.Context, row *circle.Circle) gin.H {
	item := h.circleListItem(c, row)
	merchants := h.listCircleMerchants(c, row.CircleID)
	ids := make([]uint, 0, len(merchants))
	for _, m := range merchants {
		switch v := m["mer_id"].(type) {
		case uint64:
			ids = append(ids, uint(v))
		case int64:
			ids = append(ids, uint(v))
		case uint:
			ids = append(ids, v)
		}
	}
	item["merchant"] = merchants
	item["merchant_ids"] = ids
	return item
}

func (h *Handler) countCircleMerchants(c *gin.Context, circleID uint) int64 {
	var total int64
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Where("business_id = ?", circleID).Count(&total).Error
	return total
}

func (h *Handler) listCircleMerchants(c *gin.Context, circleID uint) []gin.H {
	rows := make([]struct {
		MerID    uint64 `gorm:"column:merchant_id"`
		MerName  string `gorm:"column:merchant_name"`
		RealName string `gorm:"column:contact_name"`
		MerPhone string `gorm:"column:contact_mobile"`
	}, 0)
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name, contact_name, contact_mobile").
		Where("business_id = ?", circleID).Order("merchant_id ASC").Scan(&rows).Error
	out := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		out = append(out, gin.H{
			"mer_id": row.MerID, "mer_name": row.MerName,
			"real_name": row.RealName, "mer_phone": row.MerPhone,
		})
	}
	return out
}

func (h *Handler) syncCircleMerchants(c *gin.Context, circleID uint, merchantIDs []uint) error {
	if merchantIDs == nil {
		return nil
	}
	ctx := c.Request.Context()
	return h.adminDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_merchant_view").
			Where("business_id = ? AND merchant_id NOT IN ?", circleID, nonzeroOrPlaceholder(merchantIDs)).
			Update("business_id", 0).Error; err != nil {
			return err
		}
		if len(merchantIDs) == 0 {
			return nil
		}
		return tx.Table("qixi_crm_a_merchant_view").
			Where("merchant_id IN ?", merchantIDs).
			Update("business_id", circleID).Error
	})
}

func nonzeroOrPlaceholder(ids []uint) []uint {
	if len(ids) == 0 {
		return []uint{0}
	}
	return ids
}

func (h *Handler) ensureAdminAccountAvailable(c *gin.Context, account string) error {
	var exists int64
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Where("username = ? AND deleted_at IS NULL", account).Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return circle.ErrAccountExists
	}
	return nil
}

func (h *Handler) ensureAdminDisplayNameAvailable(c *gin.Context, displayName string, excludeAdminID uint) error {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return circle.ErrBadParam
	}
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Where("display_name = ? AND deleted_at IS NULL", displayName)
	if excludeAdminID > 0 {
		q = q.Where("id <> ?", excludeAdminID)
	}
	var exists int64
	if err := q.Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return circle.ErrDisplayNameExists
	}
	return nil
}

func (h *Handler) bindAgentAdminAccount(c *gin.Context, agentID uint, account, password, name, phone string) error {
	if err := h.ensureAdminAccountAvailable(c, account); err != nil {
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	display := strings.TrimSpace(name)
	if display == "" {
		display = account
	}
	if err := h.ensureAdminDisplayNameAvailable(c, display, 0); err != nil {
		return err
	}
	err = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").Create(map[string]any{
		"username": account, "password_hash": string(hash), "display_name": display,
		"phone": phone, "status": 1, "auth_version": 1, "data_scope_version": 1,
		"circle_agent_id": agentID,
	}).Error
	if err != nil && isMySQLDuplicate(err) {
		return circle.ErrAccountExists
	}
	return err
}

func (h *Handler) enrichAgentPage(c *gin.Context, page *circle.PageResult[circle.Agent]) gin.H {
	agentIDs := make([]uint, 0, len(page.List))
	for i := range page.List {
		if page.List[i].CircleAgentID > 0 {
			agentIDs = append(agentIDs, page.List[i].CircleAgentID)
		}
	}
	circlesByAgent := h.listAgentsCircles(c, agentIDs)
	list := make([]gin.H, 0, len(page.List))
	for i := range page.List {
		list = append(list, h.enrichAgentItem(c, &page.List[i], circlesByAgent[page.List[i].CircleAgentID]))
	}
	return gin.H{"list": list, "total": page.Total, "page": page.Page, "limit": page.Limit}
}

func (h *Handler) enrichAgentItem(c *gin.Context, row *circle.Agent, preloaded ...[]agentCircleBrief) gin.H {
	if row == nil {
		return gin.H{
			"nickname": "", "account": "", "avatar": "", "extend": "",
			"circles": []gin.H{}, "circle_ids": []uint{},
		}
	}
	item := gin.H{
		"circle_agent_id": row.CircleAgentID, "uid": row.UID, "name": row.Name, "phone": row.Phone,
		"qualification": row.Qualification, "remark": row.Remark, "extend": row.Extend,
		"avatar": circle.AgentAvatar(row.Extend), "audit_admin_id": row.AuditAdminID,
		"audit_reason": row.AuditReason, "audit_time": row.AuditTime, "status": row.Status,
		"payment_method": row.PaymentMethod, "payment_name": row.PaymentName,
		"payment_configured": row.PaymentConfigured, "balance": row.Balance, "type": row.Type,
		"business_name": row.BusinessName, "business_store_category": row.BusinessStoreCategory,
		"business_store_type": row.BusinessStoreType, "create_time": row.CreateTime,
		"nickname": "", "account": "", "circles": []gin.H{}, "circle_ids": []uint{},
	}
	var circles []agentCircleBrief
	if len(preloaded) > 0 {
		circles = preloaded[0]
	} else {
		circles = h.listAgentCircles(c, row.CircleAgentID)
	}
	if circles == nil {
		circles = []agentCircleBrief{}
	}
	circlePayload := make([]gin.H, 0, len(circles))
	ids := make([]uint, 0, len(circles))
	for _, circleRow := range circles {
		ids = append(ids, circleRow.CircleID)
		circlePayload = append(circlePayload, gin.H{
			"circle_id":        circleRow.CircleID,
			"name":             circleRow.Name,
			"type":             circleRow.Type,
			"status":           circleRow.Status,
			"commission_type":  circleRow.CommissionType,
			"commission_rate":  circleRow.CommissionRate,
		})
	}
	item["circles"] = circlePayload
	item["circle_ids"] = ids
	if row.UID > 0 && h.businessDB != nil {
		var user struct {
			Nickname  string `gorm:"column:nickname"`
			AvatarURL string `gorm:"column:avatar_url"`
		}
		_ = h.businessDB.WithContext(c.Request.Context()).
			Table("qixi_crm_b_user AS u").
			Select("u.nickname, COALESCE(p.avatar_url,'') AS avatar_url").
			Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id = u.id").
			Where("u.id = ?", row.UID).Limit(1).Scan(&user).Error
		item["nickname"] = user.Nickname
		if user.AvatarURL != "" {
			item["avatar_url"] = user.AvatarURL
		}
	}
	var admin struct {
		ID       uint   `gorm:"column:id"`
		Username string `gorm:"column:username"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Select("id, username").
		Where("circle_agent_id = ? AND deleted_at IS NULL", row.CircleAgentID).
		Limit(1).Scan(&admin).Error; err == nil && admin.ID > 0 {
		item["account"] = admin.Username
		item["admin"] = gin.H{"admin_id": admin.ID, "account": admin.Username}
	}
	return item
}

type agentCircleBrief struct {
	CircleID       uint    `gorm:"column:circle_id" json:"circle_id"`
	Name           string  `gorm:"column:name" json:"name"`
	Type           int8    `gorm:"column:type" json:"type"`
	Status         int8    `gorm:"column:status" json:"status"`
	CommissionType int8    `gorm:"column:commission_type" json:"commission_type"`
	CommissionRate float64 `gorm:"column:commission_rate" json:"commission_rate"`
}

// listAgentCircles 返回代理绑定的区域/商户（含提成比例，对齐 CRMEB「负责区域」展示）。
func (h *Handler) listAgentCircles(c *gin.Context, agentID uint) []agentCircleBrief {
	if agentID == 0 {
		return []agentCircleBrief{}
	}
	return h.listAgentsCircles(c, []uint{agentID})[agentID]
}

func (h *Handler) listAgentsCircles(c *gin.Context, agentIDs []uint) map[uint][]agentCircleBrief {
	out := make(map[uint][]agentCircleBrief, len(agentIDs))
	for _, id := range agentIDs {
		out[id] = []agentCircleBrief{}
	}
	if len(agentIDs) == 0 || h.adminDB == nil {
		return out
	}
	type agentCircleRow struct {
		CircleID       uint    `gorm:"column:circle_id"`
		Name           string  `gorm:"column:name"`
		Type           int8    `gorm:"column:type"`
		Status         int8    `gorm:"column:status"`
		CommissionType int8    `gorm:"column:commission_type"`
		CommissionRate float64 `gorm:"column:commission_rate"`
		CircleAgentID  uint    `gorm:"column:circle_agent_id"`
	}
	bound := make([]agentCircleRow, 0)
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_business_zone").
		Select("circle_id, name, type, status, commission_type, commission_rate, circle_agent_id").
		Where("circle_agent_id IN ?", agentIDs).
		Order("type ASC, circle_id ASC").Scan(&bound).Error
	for _, row := range bound {
		out[row.CircleAgentID] = append(out[row.CircleAgentID], agentCircleBrief{
			CircleID:       row.CircleID,
			Name:           row.Name,
			Type:           row.Type,
			Status:         row.Status,
			CommissionType: row.CommissionType,
			CommissionRate: row.CommissionRate,
		})
	}
	return out
}

func (h *Handler) syncAgentCircles(c *gin.Context, agentID uint, circleIDs *[]uint) error {
	if circleIDs == nil || agentID == 0 {
		return nil
	}
	ids := uniqueUints(*circleIDs)
	now := time.Now()
	return h.adminDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_a_business_zone").
			Where("circle_agent_id = ? AND type = 1 AND circle_id NOT IN ?", agentID, nonzeroOrPlaceholder(ids)).
			Updates(map[string]any{"circle_agent_id": 0, "update_time": now}).Error; err != nil {
			return err
		}
		if len(ids) == 0 {
			return nil
		}
		return tx.Table("qixi_crm_a_business_zone").
			Where("circle_id IN ? AND type = 1", ids).
			Updates(map[string]any{"circle_agent_id": agentID, "update_time": now}).Error
	})
}

func (h *Handler) syncBoundAdminProfile(c *gin.Context, agentID uint, name, phone, password string) error {
	var admin struct {
		ID uint `gorm:"column:id"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Select("id").Where("circle_agent_id = ? AND deleted_at IS NULL", agentID).
		Limit(1).Scan(&admin).Error; err != nil {
		return err
	}
	if admin.ID == 0 {
		return nil
	}
	display := strings.TrimSpace(name)
	if display == "" {
		return circle.ErrBadParam
	}
	if err := h.ensureAdminDisplayNameAvailable(c, display, admin.ID); err != nil {
		return err
	}
	updates := map[string]any{
		"display_name": display,
		"phone":        strings.TrimSpace(phone),
	}
	if password != "" {
		if len([]rune(password)) < 6 {
			return circle.ErrBadParam
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updates["password_hash"] = string(hash)
		updates["auth_version"] = gorm.Expr("auth_version + 1")
	}
	return h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Where("id = ?", admin.ID).Updates(updates).Error
}

func uniqueUints(values []uint) []uint {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[uint]struct{}, len(values))
	out := make([]uint, 0, len(values))
	for _, v := range values {
		if v == 0 {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func isMySQLDuplicate(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "Error 1062") || strings.Contains(msg, "Duplicate entry")
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, circle.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, circle.ErrBadParam), errors.Is(err, circle.ErrHasChildren), errors.Is(err, circle.ErrAlreadyAudited), errors.Is(err, circle.ErrAgentBound), errors.Is(err, circle.ErrAgentBalance), errors.Is(err, circle.ErrAgentAdminBound), errors.Is(err, circle.ErrAgentRevoked), errors.Is(err, circle.ErrAgentNotApproved), errors.Is(err, circle.ErrAgentAdminUnbound), errors.Is(err, circle.ErrCommandConflict), errors.Is(err, circle.ErrAccountExists), errors.Is(err, circle.ErrDisplayNameExists):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case err != nil && strings.Contains(err.Error(), "区域创建后不支持"):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "区域代理服务异常")
	}
}
