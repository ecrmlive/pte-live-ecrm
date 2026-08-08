package circle

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	svc     *circle.Service
	adminDB *gorm.DB
}

func NewHandler(svc *circle.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
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
	r.GET("/business-zones/:id", platform, zoneManage, h.GetCircle)
	r.PUT("/business-zones/:id", platform, zoneManage, h.UpdateCircle)
	r.PUT("/business-zones/:id/status", platform, zoneManage, h.UpdateCircleStatus)
	r.DELETE("/business-zones/:id", platform, zoneManage, h.DeleteCircle)
	r.GET("/business-zone-agents", platform, agentManage, h.ListAgents)
	r.GET("/business-zone-agents/options", platform, agentManage, h.AgentOptions)
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
	keyword := firstNonEmpty(c.Query("keyword"), c.Query("name"))
	out, err := h.svc.ListAgents(c.Request.Context(), keyword, status, agentType, p, l)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, out)
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
	response.OK(c, row)
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
	page, err := h.svc.ListAgents(c.Request.Context(), "", &approved, agentType, 1, 100)
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
	row, err := h.svc.CreateAgent(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	account, password := strings.TrimSpace(in.Account), strings.TrimSpace(in.Password)
	if account != "" && password != "" {
		if err = h.bindAgentAdminAccount(c, row.CircleAgentID, account, password, row.Name, row.Phone); err != nil {
			writeErr(c, err)
			return
		}
	}
	response.OK(c, row)
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
	response.OK(c, row)
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
	item := gin.H{
		"circle_id": row.CircleID, "pid": row.PID, "path": row.Path, "name": row.Name,
		"circle_agent_id": row.CircleAgentID, "commission_type": row.CommissionType,
		"commission_rate": row.CommissionRate, "level": row.Level, "remark": row.Remark,
		"sort": row.Sort, "status": row.Status, "type": row.Type, "role_id": row.RoleID,
		"business_store_category": row.BusinessStoreCategory, "business_store_type": row.BusinessStoreType,
		"goods_type": circle.ParseIntCSV(row.GoodsType),
		"platform_category_ids": circle.ParseUintCSV(row.PlatformCategoryIDs),
		"create_time": row.CreateTime, "merchant_count": h.countCircleMerchants(c, row.CircleID),
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

func (h *Handler) bindAgentAdminAccount(c *gin.Context, agentID uint, account, password, name, phone string) error {
	var exists int64
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").
		Where("username = ? AND deleted_at IS NULL", account).Count(&exists).Error; err != nil {
		return err
	}
	if exists > 0 {
		return errors.New("登录账号已存在")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	display := strings.TrimSpace(name)
	if display == "" {
		display = account
	}
	return h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user").Create(map[string]any{
		"username": account, "password_hash": string(hash), "display_name": display,
		"phone": phone, "status": 1, "auth_version": 1, "data_scope_version": 1,
		"circle_agent_id": agentID,
	}).Error
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, circle.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, circle.ErrBadParam), errors.Is(err, circle.ErrHasChildren), errors.Is(err, circle.ErrAlreadyAudited), errors.Is(err, circle.ErrAgentBound), errors.Is(err, circle.ErrAgentBalance), errors.Is(err, circle.ErrAgentAdminBound), errors.Is(err, circle.ErrAgentRevoked), errors.Is(err, circle.ErrAgentNotApproved), errors.Is(err, circle.ErrAgentAdminUnbound), errors.Is(err, circle.ErrCommandConflict):
		response.Fail(c, http.StatusBadRequest, err.Error())
	case err != nil && (strings.Contains(err.Error(), "登录账号已存在") || strings.Contains(err.Error(), "区域创建后不支持")):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "区域代理服务异常")
	}
}
