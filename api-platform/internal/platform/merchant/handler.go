package merchant

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/adminscope"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantonboarding"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	svc        *merchant.Service
	id         *identity.Service
	adminDB    *gorm.DB
	onboarding *merchantonboarding.Client
}

func NewHandler(svc *merchant.Service, id *identity.Service, adminDB *gorm.DB, onboarding *merchantonboarding.Client) *Handler {
	return &Handler{svc: svc, id: id, adminDB: adminDB, onboarding: onboarding}
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	categoryManage := middleware.RequireAdminMenu(h.adminDB, "merchant.category.manage")
	statusManage := middleware.RequireAdminMenu(h.adminDB, "merchant.status.manage")
	r.GET("/merchants", h.ListMerchants)
	r.GET("/merchants/:id", h.GetMerchant)
	r.GET("/merchants/:id/operate-logs", h.ListMerchantOperateLogs)
	r.POST("/merchants", platformOnly, statusManage, h.CreateMerchant)
	r.PUT("/merchants/:id", platformOnly, statusManage, h.UpdateMerchant)
	r.PUT("/merchants/:id/status", platformOnly, statusManage, h.SetMerchantStatus)
	r.PUT("/merchants/:id/recommend", platformOnly, statusManage, h.SetMerchantRecommend)

	r.GET("/merchant-intentions", h.ListIntentions)
	r.GET("/merchant-intentions/:id", h.GetIntention)
	r.POST("/merchant-intentions/:id/assign-region", middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.assign_region"), h.AssignIntentionRegion)
	r.POST("/merchant-intentions/:id/audit", middleware.RequireAdminRoles("platform", "region"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.audit"), h.AuditIntention)

	r.GET("/merchant-categories", platformOnly, h.ListCategories)
	r.POST("/merchant-categories", platformOnly, categoryManage, h.CreateCategory)
	r.PUT("/merchant-categories/:id", platformOnly, categoryManage, h.UpdateCategory)
	r.DELETE("/merchant-categories/:id", platformOnly, categoryManage, h.DeleteCategory)
}

func (h *Handler) ListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := merchant.ListFilter{
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
		Page:     page,
		Limit:    limit,
	}
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		filter.Status = &st
	}
	if s := c.Query("category_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		id := uint(v)
		filter.CategoryID = &id
	}
	if s := c.Query("type_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		id := uint(v)
		filter.TypeID = &id
	}
	if s := c.Query("region_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		id := uint(v)
		filter.RegionID = &id
	}
	if s := c.Query("is_best"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		filter.IsBest = &st
	}
	if s := c.Query("offline_pay"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		filter.OfflinePay = &st
	}
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	res, err := h.svc.ListMerchants(c.Request.Context(), filter, scope)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetMerchant(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), scope); err != nil {
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

type merchantOperateLogRow struct {
	ID           uint64    `json:"id"`
	Action       string    `json:"action"`
	ActionLabel  string    `json:"action_label"`
	Terminal     string    `json:"terminal"`
	RoleCode     string    `json:"role_code"`
	RoleLabel    string    `json:"role_label"`
	OperatorID   uint64    `json:"operator_id"`
	OperatorName string    `json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
}

// ListMerchantOperateLogs returns platform audit rows for one shop drawer「操作记录」tab.
func (h *Handler) ListMerchantOperateLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺 ID 参数错误")
		return
	}
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), scope); err != nil {
		writeErr(c, err)
		return
	}
	if h.adminDB == nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志不可用")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	terminal := strings.TrimSpace(c.Query("terminal"))
	start, end, ok := parseOperateDateRange(c.Query("start_date"), c.Query("end_date"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "操作日期参数错误")
		return
	}

	resourceID := strconv.FormatUint(id, 10)
	db := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_operation_log AS l").
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = l.admin_user_id").
		Where("l.resource_type = ? AND l.resource_id = ?", "merchants", resourceID)
	switch terminal {
	case "platform":
		db = db.Where("l.role_code LIKE ?", "%platform%")
	case "merchant":
		db = db.Where("l.role_code LIKE ?", "%merchant%")
	case "":
	default:
		response.Fail(c, http.StatusBadRequest, "操作端参数错误")
		return
	}
	if !start.IsZero() {
		db = db.Where("l.created_at >= ?", start)
	}
	if !end.IsZero() {
		db = db.Where("l.created_at < ?", end)
	}

	var total int64
	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}

	type scanRow struct {
		ID          uint64    `gorm:"column:id"`
		Action      string    `gorm:"column:action"`
		RoleCode    string    `gorm:"column:role_code"`
		AdminUserID uint64    `gorm:"column:admin_user_id"`
		Username    string    `gorm:"column:username"`
		DisplayName string    `gorm:"column:display_name"`
		CreatedAt   time.Time `gorm:"column:created_at"`
	}
	rows := make([]scanRow, 0)
	if err := db.Select("l.id,l.action,l.role_code,l.admin_user_id,COALESCE(u.username,'') AS username,COALESCE(u.display_name,'') AS display_name,l.created_at").
		Order("l.id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "操作日志查询失败")
		return
	}

	list := make([]merchantOperateLogRow, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.DisplayName)
		if name == "" {
			name = strings.TrimSpace(row.Username)
		}
		if name == "" {
			name = "未知"
		}
		list = append(list, merchantOperateLogRow{
			ID:           row.ID,
			Action:       row.Action,
			ActionLabel:  merchantOperateActionLabel(row.Action),
			Terminal:     merchantOperateTerminal(row.RoleCode),
			RoleCode:     row.RoleCode,
			RoleLabel:    merchantOperateRoleLabel(row.RoleCode),
			OperatorID:   row.AdminUserID,
			OperatorName: name,
			CreatedAt:    row.CreatedAt,
		})
	}
	response.OK(c, gin.H{"list": list, "total": total})
}

func parseOperateDateRange(startRaw, endRaw string) (time.Time, time.Time, bool) {
	startRaw, endRaw = strings.TrimSpace(startRaw), strings.TrimSpace(endRaw)
	var start, end time.Time
	var err error
	if startRaw != "" {
		start, err = time.ParseInLocation("2006-01-02", startRaw, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, false
		}
	}
	if endRaw != "" {
		end, err = time.ParseInLocation("2006-01-02", endRaw, time.Local)
		if err != nil {
			return time.Time{}, time.Time{}, false
		}
		end = end.AddDate(0, 0, 1)
	}
	if !start.IsZero() && !end.IsZero() && !start.Before(end) {
		return time.Time{}, time.Time{}, false
	}
	return start, end, true
}

func merchantOperateActionLabel(action string) string {
	a := strings.ToUpper(strings.TrimSpace(action))
	switch {
	case strings.Contains(a, "STATUS"):
		return "开启/关闭"
	case strings.Contains(a, "RECOMMEND"):
		return "推荐变更"
	case strings.HasPrefix(a, "POST "):
		return "新增"
	case strings.HasPrefix(a, "PUT "), strings.HasPrefix(a, "PATCH "):
		return "编辑"
	case strings.HasPrefix(a, "DELETE "):
		return "删除"
	default:
		if action == "" {
			return "操作"
		}
		return action
	}
}

func merchantOperateTerminal(roleCode string) string {
	if strings.Contains(roleCode, "merchant") {
		return "商户操作"
	}
	return "平台操作"
}

func merchantOperateRoleLabel(roleCode string) string {
	parts := strings.Split(roleCode, ",")
	labels := make([]string, 0, len(parts))
	for _, part := range parts {
		switch strings.TrimSpace(part) {
		case "platform":
			labels = append(labels, "平台管理员")
		case "region":
			labels = append(labels, "区域管理员")
		case "merchant":
			labels = append(labels, "商户管理员")
		case "operations":
			labels = append(labels, "运营")
		case "":
		default:
			labels = append(labels, part)
		}
	}
	if len(labels) == 0 {
		return "—"
	}
	return strings.Join(labels, "/")
}

type statusReq struct {
	Enabled bool `json:"enabled"`
}

func (h *Handler) SetMerchantStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), scope); err != nil {
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

func (h *Handler) SetMerchantRecommend(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), scope); err != nil {
		writeErr(c, err)
		return
	}
	var req statusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetMerchantRecommend(c.Request.Context(), uint(id), req.Enabled); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) CreateMerchant(c *gin.Context) {
	var req merchant.ShopProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateMerchant(c.Request.Context(), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateMerchant(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	scope, ok := h.merchantScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置商户监管数据范围")
		return
	}
	if err := h.svc.RequireMerchantScope(c.Request.Context(), uint(id), scope); err != nil {
		writeErr(c, err)
		return
	}
	var req merchant.ShopProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateShopProfile(c.Request.Context(), uint(id), req)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

// merchantScope resolves either the direct merchant IDs or region IDs assigned
// to the unified-admin account. Missing configuration denies access rather than
// falling back to a full-platform view.
func (h *Handler) merchantScope(c *gin.Context) (merchant.MerchantScope, bool) {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return merchant.MerchantScope{}, false
	}
	resolved, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, claims)
	if err != nil {
		return merchant.MerchantScope{}, false
	}
	if resolved.Full {
		return merchant.MerchantScope{}, true
	}
	return merchant.MerchantScope{MerchantIDs: uints(resolved.MerchantIDs), RegionIDs: uints(resolved.RegionIDs)}, true
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
	regionIDs, ok := h.intentionScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "当前角色无商户入驻审核范围")
		return
	}
	res, err := h.svc.ListIntentions(c.Request.Context(), c.Query("keyword"), statusPtr, regionIDs, page, limit, strings.TrimSpace(c.Query("date_from")), strings.TrimSpace(c.Query("date_to")))
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetIntention(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	regionIDs, ok := h.intentionScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "当前角色无商户入驻审核范围")
		return
	}
	row, err := h.svc.GetIntention(c.Request.Context(), uint(id), regionIDs)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

type assignIntentionRegionRequest struct {
	RegionID uint `json:"region_id"`
}

func (h *Handler) AssignIntentionRegion(c *gin.Context) {
	if !isPlatform(c) {
		response.Fail(c, http.StatusForbidden, "只有平台角色可以分配入驻审核区域")
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req assignIntentionRegionRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.RegionID == 0 {
		response.Fail(c, http.StatusBadRequest, "请选择有效区域")
		return
	}
	row, err := h.svc.AssignIntentionRegion(c.Request.Context(), uint(id), req.RegionID)
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
	regionIDs, ok := h.intentionScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "当前角色无商户入驻审核范围")
		return
	}
	if req.Status == merchant.IntentionApproved {
		if req.RegionID == 0 || strings.TrimSpace(req.Account) == "" || len(req.Password) < 8 {
			response.Fail(c, http.StatusBadRequest, "通过入驻必须填写区域、账号和至少 8 位初始密码")
			return
		}
		row, err := h.svc.GetIntention(c.Request.Context(), uint(id), regionIDs)
		if err != nil {
			writeErr(c, err)
			return
		}
		if regionIDs != nil && req.RegionID != row.CircleID {
			response.Fail(c, http.StatusForbidden, "区域审核只能开通分配给本区域的申请")
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "初始密码处理失败")
			return
		}
		provisioned, err := h.onboarding.Provision(c.Request.Context(), merchantonboarding.Request{ApplicationID: uint(id), RegionID: req.RegionID, MerchantName: row.MerName, ContactName: row.Name, ContactMobile: row.Phone, Account: strings.TrimSpace(req.Account), PasswordHash: string(hash)})
		if err != nil {
			response.Fail(c, http.StatusServiceUnavailable, err.Error())
			return
		}
		res, err := h.svc.FinalizeIntentionApproval(c.Request.Context(), uint(id), req, provisioned.MerchantID, regionIDs)
		if err != nil {
			writeErr(c, err)
			return
		}
		response.OK(c, res)
		return
	}
	res, err := h.svc.AuditIntention(c.Request.Context(), uint(id), req, regionIDs)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) intentionScope(c *gin.Context) ([]uint, bool) {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return nil, false
	}
	for _, role := range claims.Roles {
		if role == "platform" {
			return nil, true
		}
	}
	if !hasRole(claims.Roles, "region") {
		return nil, false
	}
	resolved, err := adminscope.ResolveMerchantScope(c.Request.Context(), h.adminDB, claims)
	if err != nil || len(resolved.RegionIDs) == 0 {
		return nil, false
	}
	return uints(resolved.RegionIDs), true
}

func isPlatform(c *gin.Context) bool {
	claims := middleware.ClaimsFrom(c)
	return claims != nil && hasRole(claims.Roles, "platform")
}

func uints(values []uint64) []uint {
	out := make([]uint, 0, len(values))
	for _, value := range values {
		if value > 0 {
			out = append(out, uint(value))
		}
	}
	return out
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
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
	id, parseErr := strconv.ParseUint(c.Param("id"), 10, 64)
	var req categoryReq
	if parseErr != nil || id == 0 || c.ShouldBindJSON(&req) != nil {
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
	id, parseErr := strconv.ParseUint(c.Param("id"), 10, 64)
	if parseErr != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
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
	case errors.Is(err, merchant.ErrConflict):
		response.Fail(c, http.StatusConflict, err.Error())
	case errors.Is(err, merchant.ErrAlreadyAudited),
		errors.Is(err, merchant.ErrBadStatus),
		errors.Is(err, merchant.ErrRejectNeedMsg),
		errors.Is(err, merchant.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
