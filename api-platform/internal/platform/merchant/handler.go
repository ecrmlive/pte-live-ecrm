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
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
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
	merchantDB *gorm.DB
	storeJWT   *authjwt.Manager
	onboarding *merchantonboarding.Client
}

func NewHandler(svc *merchant.Service, id *identity.Service, adminDB *gorm.DB, onboarding *merchantonboarding.Client) *Handler {
	return &Handler{svc: svc, id: id, adminDB: adminDB, onboarding: onboarding}
}

// WithStoreLogin enables platform→store console handoff (read merchant DB, issue store_console JWT).
func (h *Handler) WithStoreLogin(merchantDB *gorm.DB, storeJWT *authjwt.Manager) *Handler {
	if h == nil {
		return nil
	}
	h.merchantDB = merchantDB
	h.storeJWT = storeJWT
	return h
}

func (h *Handler) Register(r gin.IRoutes) {
	platformOnly := middleware.RequireAdminRoles("platform")
	categoryManage := middleware.RequireAdminMenu(h.adminDB, "merchant.category.manage")
	statusManage := middleware.RequireAdminMenu(h.adminDB, "merchant.status.manage")
	r.GET("/merchants", h.ListMerchants)
	r.GET("/merchants/:id", h.GetMerchant)
	r.GET("/merchants/:id/operate-logs", h.ListMerchantOperateLogs)
	r.POST("/merchants/:id/login", h.LoginMerchant)
	r.POST("/merchants", platformOnly, statusManage, h.CreateMerchant)
	r.PUT("/merchants/:id", platformOnly, statusManage, h.UpdateMerchant)
	r.PUT("/merchants/:id/status", platformOnly, statusManage, h.SetMerchantStatus)
	r.PUT("/merchants/:id/recommend", platformOnly, statusManage, h.SetMerchantRecommend)

	r.GET("/merchant-intentions", h.ListIntentions)
	r.GET("/merchant-intentions/:id", h.GetIntention)
	r.POST("/merchant-intentions", middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.create"), h.CreateIntention)
	r.POST("/merchant-intentions/:id/assign-region", middleware.RequireAdminRoles("platform"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.assign_region"), h.AssignIntentionRegion)
	r.POST("/merchant-intentions/:id/audit", middleware.RequireAdminRoles("platform", "region"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.audit"), h.AuditIntention)
	r.DELETE("/merchant-intentions/:id", middleware.RequireAdminRoles("platform", "region"), middleware.RequireAdminMenu(h.adminDB, "merchant.intention.delete"), h.DeleteIntention)

	r.GET("/merchant-categories", platformOnly, h.ListCategories)
	r.POST("/merchant-categories", platformOnly, categoryManage, h.CreateCategory)
	r.PUT("/merchant-categories/:id", platformOnly, categoryManage, h.UpdateCategory)
	r.DELETE("/merchant-categories/:id", platformOnly, categoryManage, h.DeleteCategory)
}

func (h *Handler) ListMerchants(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
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
		c.Error(err)
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

type storeOwnerAccount struct {
	AccountID    uint64 `gorm:"column:account_id"`
	StoreID      uint64 `gorm:"column:store_id"`
	MerchantID   uint64 `gorm:"column:merchant_id"`
	StoreAppID   string `gorm:"column:store_app_id"`
	IMSDKAppID   string `gorm:"column:im_sdk_app_id"`
	Username     string `gorm:"column:username"`
	RoleCode     string `gorm:"column:role_code"`
	AuthVersion  uint64 `gorm:"column:auth_version"`
	StoreName    string `gorm:"column:store_name"`
	MerchantName string `gorm:"column:merchant_name"`
}

// LoginMerchant issues a store_console session for the shop owner (CRMEB system/merchant/login/:id).
// Platform never writes qixi_crm_m_*; it only reads the owner account and signs with the store JWT secret.
func (h *Handler) LoginMerchant(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if h.merchantDB == nil || h.storeJWT == nil {
		response.Fail(c, http.StatusServiceUnavailable, "店铺一键登录未配置")
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
	m, err := h.svc.GetMerchant(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	owner, err := h.findStoreOwner(c, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Fail(c, http.StatusNotFound, "店铺账号不存在或尚未开通")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询店铺账号失败")
		return
	}
	if owner.StoreAppID == "" {
		response.Fail(c, http.StatusConflict, "店铺缺少 AppId，无法签发后台会话")
		return
	}
	pair, err := h.storeJWT.IssueStoreConsoleWithIdentityVersion(
		uint(owner.AccountID),
		uint(owner.MerchantID),
		uint(owner.StoreID),
		owner.StoreAppID,
		owner.IMSDKAppID,
		owner.Username,
		owner.RoleCode,
		owner.AuthVersion,
	)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "签发令牌失败")
		return
	}
	response.OK(c, gin.H{
		"token":        pair,
		"mer_id":       owner.MerchantID,
		"store_id":     owner.StoreID,
		"account":      owner.Username,
		"mer_name":     firstNonEmpty(owner.MerchantName, m.MerName),
		"store_name":   firstNonEmpty(owner.StoreName, m.MerName),
		"store_app_id": owner.StoreAppID,
		"path":         "/auth/login",
	})
}

func (h *Handler) findStoreOwner(c *gin.Context, merchantID uint) (*storeOwnerAccount, error) {
	var row storeOwnerAccount
	// Prefer owner; fall back to any enabled account on the merchant's unique store.
	err := h.merchantDB.WithContext(c.Request.Context()).
		Table("qixi_crm_m_account AS a").
		Select(`a.id AS account_id, a.store_id, a.username, a.role_code, a.auth_version,
			s.merchant_id, s.app_id AS store_app_id, s.name AS store_name, m.name AS merchant_name,
			COALESCE(im.sdk_app_id, '') AS im_sdk_app_id`).
		Joins("INNER JOIN qixi_crm_m_store AS s ON s.id = a.store_id AND s.status = 1").
		Joins("INNER JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id AND m.status = 1").
		Joins("LEFT JOIN qixi_crm_m_im_sdk_app AS im ON im.merchant_id = s.merchant_id AND im.status = 'enabled' AND im.is_active = 1").
		Where("s.merchant_id = ? AND a.status = 1", merchantID).
		Order("CASE WHEN a.role_code = 'owner' THEN 0 ELSE 1 END, a.id ASC").
		Take(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
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
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	var statusPtr *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		statusPtr = &st
	}
	var categoryID, typeID *uint
	if raw := strings.TrimSpace(c.Query("category_id")); raw != "" {
		v, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && v > 0 {
			id := uint(v)
			categoryID = &id
		}
	}
	if raw := strings.TrimSpace(c.Query("type_id")); raw != "" {
		v, err := strconv.ParseUint(raw, 10, 64)
		if err == nil && v > 0 {
			id := uint(v)
			typeID = &id
		}
	}
	regionIDs, ok := h.intentionScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "当前角色无店铺入驻申请范围")
		return
	}
	res, err := h.svc.ListIntentions(c.Request.Context(), c.Query("keyword"), statusPtr, regionIDs, page, limit, strings.TrimSpace(c.Query("date_from")), strings.TrimSpace(c.Query("date_to")), categoryID, typeID)
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
		response.Fail(c, http.StatusForbidden, "当前角色无店铺入驻申请范围")
		return
	}
	row, err := h.svc.GetIntention(c.Request.Context(), uint(id), regionIDs)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) CreateIntention(c *gin.Context) {
	var req merchant.CreateIntentionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateIntention(c.Request.Context(), req)
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
		row, err := h.svc.GetIntention(c.Request.Context(), uint(id), regionIDs)
		if err != nil {
			writeErr(c, err)
			return
		}
		// 对齐 CRMEB：账号=申请人手机，初始密码=手机号后六位（不足 8 位时左侧补 Qx），区域=申请已分配 region。
		account := strings.TrimSpace(row.Phone)
		if account == "" {
			response.Fail(c, http.StatusBadRequest, "申请缺少联系手机，无法开通商户管理账号")
			return
		}
		password := account
		if len(account) >= 6 {
			password = account[len(account)-6:]
		}
		if len(password) < 8 {
			password = "Qx" + password
			if len(password) < 8 {
				password = (password + "00000000")[:8]
			}
		}
		regionID := row.CircleID
		if regionID == 0 {
			response.Fail(c, http.StatusBadRequest, "请先为该入驻申请分配所属区域后再审核通过")
			return
		}
		if regionIDs != nil && regionID != row.CircleID {
			response.Fail(c, http.StatusForbidden, "区域审核只能开通分配给本区域的申请")
			return
		}
		req.Account = account
		req.Password = password
		req.RegionID = regionID
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "初始密码处理失败")
			return
		}
		provisioned, err := h.onboarding.Provision(c.Request.Context(), merchantonboarding.Request{
			ApplicationID: uint(id),
			RegionID:      regionID,
			MerchantName:  row.MerName,
			ContactName:   row.Name,
			ContactMobile: row.Phone,
			Account:       account,
			PasswordHash:  string(hash),
		})
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

func (h *Handler) DeleteIntention(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	regionIDs, ok := h.intentionScope(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "当前角色无商户入驻审核范围")
		return
	}
	if err := h.svc.DeleteIntention(c.Request.Context(), uint(id), regionIDs); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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
