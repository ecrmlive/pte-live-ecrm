// Package customerservice implements the unified-admin service queue.
// Message transport and credentials intentionally remain owned by pte-live-im.
package customerservice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Handler struct {
	adminDB    *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(adminDB, businessDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB, businessDB: businessDB}
}

func (h *Handler) Register(routes gin.IRoutes) {
	routes.GET("/customer-service/agents", h.ListAgents)
	routes.GET("/customer-service/agents/:agent_id/users", h.ListAgentUsers)
	routes.GET("/customer-service/settings", h.GetSettings)
	routes.PUT("/customer-service/settings", h.UpdateSettings)
	routes.GET("/customer-service/threads", h.List)
	routes.GET("/customer-service/threads/:id", h.Detail)
	routes.POST("/customer-service/threads/:id/claim", h.Claim)
	routes.POST("/customer-service/threads/:id/transfer", h.Transfer)
	routes.GET("/customer-service/threads/:id/messages", h.Messages)
	routes.GET("/customer-service/threads/:id/assignment-logs", h.AssignmentLogs)
	routes.GET("/customer-service/threads/:id/order", h.Order)
	routes.GET("/customer-service/threads/:id/delivery", h.Delivery)
	routes.GET("/customer-service/threads/:id/products", h.Products)
	routes.GET("/customer-service/threads/:id/refunds", h.Refunds)
	routes.GET("/customer-service/threads/:id/user", h.User)
	routes.PUT("/customer-service/threads/:id/user-note", h.UpdateUserNote)
	routes.GET("/customer-service/quick-replies", h.ListQuickReplies)
	routes.POST("/customer-service/quick-replies", h.CreateQuickReply)
	routes.PUT("/customer-service/quick-replies/:id", h.UpdateQuickReply)
	routes.DELETE("/customer-service/quick-replies/:id", h.DeleteQuickReply)
}

type serviceBinding struct {
	ID               uint64     `gorm:"column:id" json:"id"`
	UserID           uint64     `gorm:"column:user_id" json:"user_id"`
	StoreID          uint64     `gorm:"column:store_id" json:"store_id"`
	StoreName        string     `gorm:"column:store_name" json:"store_name"`
	MerchantID       uint64     `gorm:"column:merchant_id" json:"merchant_id"`
	IMSDKAppID       string     `gorm:"column:im_sdk_app_id" json:"im_sdk_app_id"`
	OrderID          *uint64    `gorm:"column:order_id" json:"order_id"`
	IMConversationID string     `gorm:"column:im_conversation_id" json:"im_conversation_id"`
	Status           string     `gorm:"column:status" json:"status"`
	AssignedAdminID  *uint64    `gorm:"column:assigned_admin_id" json:"assigned_admin_id"`
	AssignedAt       *time.Time `gorm:"column:assigned_at" json:"assigned_at"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// messageRow is a business projection for order/system events only. IM body
// transport remains in pte-live-im and is never fabricated from this table.
type messageRow struct {
	ID         uint64    `gorm:"column:id" json:"id"`
	BindingID  uint64    `gorm:"column:binding_id" json:"binding_id"`
	SenderRole string    `gorm:"column:sender_role" json:"sender_role"`
	SenderID   uint64    `gorm:"column:sender_id" json:"sender_id"`
	MsgType    string    `gorm:"column:msg_type" json:"msg_type"`
	Content    string    `gorm:"column:content" json:"content"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

// assignmentLogRow is a write-once transfer audit record. It deliberately
// excludes the idempotency key from API responses because that key is a caller
// retry secret, not a business-facing audit field.
type assignmentLogRow struct {
	ID              uint64    `gorm:"column:id" json:"id"`
	BindingID       uint64    `gorm:"column:binding_id" json:"binding_id"`
	FromAdminID     *uint64   `gorm:"column:from_admin_id" json:"from_admin_id"`
	TargetAdminID   uint64    `gorm:"column:target_admin_id" json:"target_admin_id"`
	OperatorAdminID uint64    `gorm:"column:operator_admin_id" json:"operator_admin_id"`
	Reason          string    `gorm:"column:reason" json:"reason"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

type dataScope struct {
	StoreIDs []uint64 `json:"store_ids"`
}

type quickReply struct {
	ID          uint64     `gorm:"column:id" json:"id"`
	StoreID     uint64     `gorm:"column:store_id" json:"store_id"`
	Title       string     `gorm:"column:title" json:"title"`
	Content     string     `gorm:"column:content" json:"content"`
	MessageType string     `gorm:"column:message_type" json:"message_type"`
	Status      string     `gorm:"column:status" json:"status"`
	CreatedBy   uint64     `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   uint64     `gorm:"column:updated_by" json:"updated_by"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"-"`
}

type quickReplyInput struct {
	StoreID     uint64 `json:"store_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	MessageType string `json:"message_type"`
	Status      string `json:"status"`
}

type userNoteInput struct {
	Content string `json:"content"`
}

type transferInput struct {
	TargetAdminID uint64 `json:"target_admin_id"`
	Reason        string `json:"reason"`
}

// serviceAgent is a deliberately small projection for service dispatch. It
// never exposes an account password or IM credential. Phone is only returned
// to the platform supervisor so that it can safely edit an existing account
// without accidentally clearing the persisted value.
type serviceAgent struct {
	ID              uint64    `json:"id"`
	Account         string    `json:"account"`
	AvatarURL       string    `json:"avatar_url"`
	CreatedAt       time.Time `json:"created_at"`
	DisplayName     string    `json:"display_name"`
	LinkedUserID    uint64    `json:"linked_user_id"`
	Phone           string    `json:"phone"`
	Roles           string    `json:"roles"`
	Status          int8      `json:"status"`
	ServiceStoreIDs []uint64  `json:"service_store_ids"`
	WechatUsername  string    `json:"wechat_username"`
}

// serviceAgentUser is a constrained roster projection. Mobile numbers stay
// masked in the staff overview; a permitted agent can use the existing
// thread-level support panel when order fulfilment requires more context.
type serviceAgentUser struct {
	BindingID uint64    `gorm:"column:binding_id" json:"binding_id"`
	UserID    uint64    `gorm:"column:user_id" json:"user_id"`
	Nickname  string    `gorm:"column:nickname" json:"nickname"`
	Mobile    string    `gorm:"column:mobile" json:"mobile"`
	StoreID   uint64    `gorm:"column:store_id" json:"store_id"`
	StoreName string    `gorm:"column:store_name" json:"store_name"`
	Status    string    `gorm:"column:status" json:"status"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// serviceSettings covers customer entry and workflow settings owned by this
// console. It stores only public entry addresses and identifiers; IM transport,
// UserSig and third-party credentials remain owned by pte-live-im/cloud config.
type serviceSettings struct {
	AutoReplyEnabled       bool   `json:"auto_reply_enabled"`
	AutoReplyText          string `json:"auto_reply_text"`
	EnterpriseWechatCorpID string `json:"enterprise_wechat_corp_id"`
	EnterpriseWechatURL    string `json:"enterprise_wechat_url"`
	QueueMode              string `json:"queue_mode"`
	MaxSessionsPerAgent    int    `json:"max_sessions_per_agent"`
	RedirectURL            string `json:"redirect_url"`
	ServicePhone           string `json:"service_phone"`
	ServiceType            string `json:"service_type"`
}

const serviceSettingsConfigKey = "customer_service.settings"

type serviceConfigRow struct {
	ConfigKey   string          `gorm:"column:config_key;primaryKey"`
	ConfigValue json.RawMessage `gorm:"column:config_value"`
	UpdatedBy   *uint64         `gorm:"column:updated_by"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (h *Handler) access(c *gin.Context) (all bool, storeIDs []uint64, ok bool) {
	claims := middleware.ClaimsFrom(c)
	if claims == nil {
		return false, nil, false
	}
	if hasRole(claims.Roles, "platform") {
		return true, nil, true
	}
	if !hasRole(claims.Roles, "customer_service") {
		return false, nil, false
	}
	storeIDs, err := h.serviceQueueStores(c.Request.Context(), uint64(claims.AdminID))
	if err != nil || len(storeIDs) == 0 {
		return false, nil, false
	}
	return false, storeIDs, true
}

// ListAgents provides the customer-service roster and the safe target list
// used when a service conversation is transferred. A customer-service user
// sees only peers sharing at least one queue store; platform users supervise
// the full roster. It deliberately projects no password or IM credential.
func (h *Handler) ListAgents(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	var rows []struct {
		ID           uint64    `gorm:"column:id"`
		Account      string    `gorm:"column:account"`
		AvatarURL    string    `gorm:"column:avatar_url"`
		CreatedAt    time.Time `gorm:"column:created_at"`
		DisplayName  string    `gorm:"column:display_name"`
		LinkedUserID uint64    `gorm:"column:linked_user_id"`
		Phone        string    `gorm:"column:phone"`
		Status       int8      `gorm:"column:status"`
	}
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_admin_user AS u").
		Select("DISTINCT u.id,u.username AS account,u.avatar_url,u.display_name,u.linked_user_id,u.phone,u.status,u.created_at").
		Joins("JOIN qixi_crm_a_admin_user_role AS ur ON ur.admin_user_id = u.id").
		Joins("JOIN qixi_crm_a_role AS r ON r.id = ur.role_id").
		Where("u.deleted_at IS NULL AND r.code = ? AND r.status = ?", "customer_service", 1).
		Order("u.created_at DESC,u.id DESC")
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("(u.username LIKE ? OR u.display_name LIKE ? OR u.phone LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if rawStatus := strings.TrimSpace(c.Query("status")); rawStatus == "0" || rawStatus == "1" {
		q = q.Where("u.status = ?", rawStatus)
	}
	if err := q.Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服人员失败")
		return
	}
	linkedUsers := make(map[uint64]struct {
		AvatarURL string `gorm:"column:avatar_url"`
		Nickname  string `gorm:"column:nickname"`
	})
	linkedIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row.LinkedUserID != 0 {
			linkedIDs = append(linkedIDs, row.LinkedUserID)
		}
	}
	if len(linkedIDs) > 0 {
		var profiles []struct {
			ID        uint64 `gorm:"column:id"`
			AvatarURL string `gorm:"column:avatar_url"`
			Nickname  string `gorm:"column:nickname"`
		}
		if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user AS u").
			Select("u.id,COALESCE(u.nickname, '') AS nickname,COALESCE(p.avatar_url, '') AS avatar_url").
			Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id = u.id").
			Where("u.id IN ?", linkedIDs).Scan(&profiles).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询客服关联用户失败")
			return
		}
		for _, profile := range profiles {
			linkedUsers[profile.ID] = struct {
				AvatarURL string `gorm:"column:avatar_url"`
				Nickname  string `gorm:"column:nickname"`
			}{AvatarURL: profile.AvatarURL, Nickname: profile.Nickname}
		}
	}

	agents := make([]serviceAgent, 0, len(rows))
	for _, row := range rows {
		assignedStores, err := h.serviceQueueStores(c.Request.Context(), row.ID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询客服服务范围失败")
			return
		}
		// Keep the HTTP contract stable for historical customer-service accounts
		// that have not yet received a queue scope: the frontend expects an array,
		// never a JSON null value.
		assignedStores = serviceStoreIDs(assignedStores)
		if !all && !sharesStore(storeIDs, assignedStores) {
			continue
		}
		roles, err := h.agentRoleCodes(c.Request.Context(), row.ID)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询客服角色失败")
			return
		}
		phone := ""
		if all {
			phone = row.Phone
		}
		linkedUser := linkedUsers[row.LinkedUserID]
		avatarURL := row.AvatarURL
		if avatarURL == "" {
			avatarURL = linkedUser.AvatarURL
		}
		wechatUsername := linkedUser.Nickname
		if wechatUsername == "" {
			wechatUsername = row.Account
		}
		agents = append(agents, serviceAgent{ID: row.ID, Account: row.Account, AvatarURL: avatarURL, CreatedAt: row.CreatedAt, DisplayName: row.DisplayName, LinkedUserID: row.LinkedUserID, Phone: phone, Roles: strings.Join(roles, ","), Status: row.Status, ServiceStoreIDs: assignedStores, WechatUsername: wechatUsername})
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 100
	}
	total := len(agents)
	start := (page - 1) * limit
	if start >= total {
		agents = []serviceAgent{}
	} else {
		end := start + limit
		if end > total {
			end = total
		}
		agents = agents[start:end]
	}
	response.OK(c, gin.H{"list": agents, "total": total, "page": page, "limit": limit})
}

func (h *Handler) agentRoleCodes(ctx context.Context, adminID uint64) ([]string, error) {
	var roles []string
	err := h.adminDB.WithContext(ctx).Table("qixi_crm_a_role AS r").
		Select("r.code").
		Joins("JOIN qixi_crm_a_admin_user_role AS ur ON ur.role_id = r.id").
		Where("ur.admin_user_id = ? AND r.status = ?", adminID, 1).
		Order("r.code ASC").Scan(&roles).Error
	return roles, err
}

// ListAgentUsers covers CRMEB's "客服的全部用户" and "客服的聊天用户列表"
// through the authoritative assignment bindings. It cannot be used to list a
// customer outside the caller's service-queue stores.
func (h *Handler) ListAgentUsers(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	agentID, err := strconv.ParseUint(c.Param("agent_id"), 10, 64)
	if err != nil || agentID == 0 {
		response.Fail(c, http.StatusBadRequest, "客服人员编号错误")
		return
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 100
	}
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_customer_service_binding AS b").
		Select(`b.id AS binding_id,b.user_id,COALESCE(u.nickname, '') AS nickname,COALESCE(u.mobile, '') AS mobile,
			b.store_id,COALESCE(s.store_name, '') AS store_name,b.status,b.updated_at`).
		Joins("LEFT JOIN qixi_crm_b_user AS u ON u.id = b.user_id").
		Joins("LEFT JOIN qixi_crm_b_store_view AS s ON s.store_id = b.store_id").
		Where("b.assigned_admin_id = ?", agentID)
	if !all {
		q = q.Where("b.store_id IN ?", storeIDs)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服服务用户失败")
		return
	}
	var rows []serviceAgentUser
	if err := q.Order("b.updated_at DESC,b.id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服服务用户失败")
		return
	}
	for i := range rows {
		rows[i].Mobile = maskMobile(rows[i].Mobile)
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func defaultServiceSettings() serviceSettings {
	return serviceSettings{ServiceType: "system", QueueMode: "manual", MaxSessionsPerAgent: 20}
}

func (h *Handler) GetSettings(c *gin.Context) {
	if _, _, ok := h.access(c); !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	settings, updatedAt, err := h.loadSettings(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取客服设置失败")
		return
	}
	response.OK(c, gin.H{"settings": settings, "updated_at": updatedAt})
}

// UpdateSettings is platform-only. The settings contain business text and
// queue policy only; no endpoint accepts an IM token, password, or cloud key.
func (h *Handler) UpdateSettings(c *gin.Context) {
	all, _, ok := h.access(c)
	if !ok || !all {
		response.Fail(c, http.StatusForbidden, "仅平台管理员可以修改客服设置")
		return
	}
	var settings serviceSettings
	if c.ShouldBindJSON(&settings) != nil || !validServiceSettings(&settings) {
		response.Fail(c, http.StatusBadRequest, "客服设置参数错误")
		return
	}
	payload, err := json.Marshal(settings)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "序列化客服设置失败")
		return
	}
	operatorID := uint64(middleware.AdminID(c))
	row := serviceConfigRow{ConfigKey: serviceSettingsConfigKey, ConfigValue: payload, UpdatedBy: &operatorID}
	if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_config").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"config_value", "updated_by", "updated_at"}),
	}).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存客服设置失败")
		return
	}
	settings, updatedAt, err := h.loadSettings(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取客服设置失败")
		return
	}
	response.OK(c, gin.H{"settings": settings, "updated_at": updatedAt})
}

func (h *Handler) loadSettings(ctx context.Context) (serviceSettings, *time.Time, error) {
	settings := defaultServiceSettings()
	var row serviceConfigRow
	err := h.adminDB.WithContext(ctx).Table("qixi_crm_a_config").Where("config_key = ?", serviceSettingsConfigKey).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return settings, nil, nil
	}
	if err != nil {
		return serviceSettings{}, nil, err
	}
	if err := json.Unmarshal(row.ConfigValue, &settings); err != nil || !validServiceSettings(&settings) {
		return serviceSettings{}, nil, errInvalidServiceSettings
	}
	return settings, &row.UpdatedAt, nil
}

func (h *Handler) serviceQueueStores(ctx context.Context, adminID uint64) ([]uint64, error) {
	var rows []struct {
		ScopeValue json.RawMessage `gorm:"column:scope_value"`
	}
	if err := h.adminDB.WithContext(ctx).
		Table("qixi_crm_a_data_scope").
		Select("scope_value").
		Where("admin_user_id = ? AND scope_type = ?", adminID, "service_queue").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	set := map[uint64]struct{}{}
	for _, row := range rows {
		var value dataScope
		if json.Unmarshal(row.ScopeValue, &value) != nil {
			continue
		}
		for _, id := range value.StoreIDs {
			if id > 0 {
				set[id] = struct{}{}
			}
		}
	}
	if len(set) == 0 {
		return nil, nil
	}
	storeIDs := make([]uint64, 0, len(set))
	for id := range set {
		storeIDs = append(storeIDs, id)
	}
	return storeIDs, nil
}

func (h *Handler) query(c *gin.Context, all bool, storeIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_customer_service_binding AS b").
		Select(`b.id, b.user_id, b.store_id, COALESCE(s.store_name, '') AS store_name,
				COALESCE(s.merchant_id, 0) AS merchant_id, COALESCE(im.sdk_app_id, '') AS im_sdk_app_id, b.order_id, b.im_conversation_id,
				b.status, b.assigned_admin_id, b.assigned_at, b.created_at, b.updated_at`).
		Joins("LEFT JOIN qixi_crm_b_store_view AS s ON s.store_id = b.store_id").
		Joins("LEFT JOIN qixi_crm_b_merchant_im_sdk_app_view AS im ON im.merchant_id = s.merchant_id")
	if !all {
		q = q.Where("b.store_id IN ?", storeIDs)
	}
	return q
}

func (h *Handler) List(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 100
	}
	q := h.query(c, all, storeIDs)
	status := strings.TrimSpace(c.Query("status"))
	if status == "open" || status == "closed" {
		q = q.Where("b.status = ?", status)
	}
	if c.Query("mine") == "1" {
		q = q.Where("b.assigned_admin_id = ?", middleware.AdminID(c))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服会话失败")
		return
	}
	var rows []serviceBinding
	if err := q.Order("CASE WHEN b.assigned_admin_id IS NULL THEN 0 ELSE 1 END, b.updated_at DESC").
		Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服会话失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Detail(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "会话编号错误")
		return
	}
	var row serviceBinding
	if err := h.query(c, all, storeIDs).Where("b.id = ?", id).Scan(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服会话失败")
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "客服会话不存在")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Claim(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "会话编号错误")
		return
	}
	var row serviceBinding
	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		q := tx.Table("qixi_crm_b_customer_service_binding AS b").
			Select("b.id, b.store_id, b.status, b.assigned_admin_id").
			Where("b.id = ?", id).Clauses(clause.Locking{Strength: "UPDATE"})
		if !all {
			q = q.Where("b.store_id IN ?", storeIDs)
		}
		if err := q.Take(&row).Error; err != nil {
			return err
		}
		if err := validateClaim(row.Status, row.AssignedAdminID, uint64(middleware.AdminID(c))); err != nil {
			return err
		}
		result := tx.Table("qixi_crm_b_customer_service_binding").
			Where("id = ? AND status = ? AND (assigned_admin_id IS NULL OR assigned_admin_id = ?)", id, "open", middleware.AdminID(c)).
			Updates(map[string]any{"assigned_admin_id": middleware.AdminID(c), "assigned_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errThreadTaken
		}
		return nil
	})
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response.Fail(c, http.StatusNotFound, "客服会话不存在")
		case errThreadClosed:
			response.Fail(c, http.StatusConflict, "客服会话已关闭")
		case errThreadTaken:
			response.Fail(c, http.StatusConflict, "客服会话已被其他客服领取")
		default:
			response.Fail(c, http.StatusInternalServerError, "领取客服会话失败")
		}
		return
	}
	// Reload through the same scope-filtered query so no write response leaks a foreign store.
	if err := h.query(c, all, storeIDs).Where("b.id = ?", id).Scan(&row).Error; err != nil || row.ID == 0 {
		response.Fail(c, http.StatusInternalServerError, "读取领取结果失败")
		return
	}
	response.OK(c, row)
}

// Transfer keeps queue ownership inside the same authorized store. The current
// owner may transfer; a platform operator may dispatch an unclaimed thread.
// Every transfer requires a caller-provided idempotency key and leaves an
// immutable business-side audit row.
func (h *Handler) Transfer(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "会话编号错误")
		return
	}
	key := strings.TrimSpace(c.GetHeader("X-Idempotency-Key"))
	var in transferInput
	if c.ShouldBindJSON(&in) != nil || !validTransfer(&in, key) {
		response.Fail(c, http.StatusBadRequest, "转接客服、原因或幂等键错误")
		return
	}
	if !h.transferTargetEligible(c.Request.Context(), in.TargetAdminID, 0) {
		response.Fail(c, http.StatusBadRequest, "目标客服未启用或未配置服务队列")
		return
	}
	claims := middleware.ClaimsFrom(c)
	operatorID := uint64(middleware.AdminID(c))
	var row serviceBinding
	err := h.businessDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		q := tx.Table("qixi_crm_b_customer_service_binding AS b").
			Select("b.id,b.store_id,b.status,b.assigned_admin_id").
			Where("b.id = ?", id).Clauses(clause.Locking{Strength: "UPDATE"})
		if !all {
			q = q.Where("b.store_id IN ?", storeIDs)
		}
		if err := q.Take(&row).Error; err != nil {
			return err
		}
		if row.Status != "open" {
			return errThreadClosed
		}
		if !h.transferTargetEligible(c.Request.Context(), in.TargetAdminID, row.StoreID) {
			return errTargetNotEligible
		}
		var logRow struct {
			TargetAdminID   uint64 `gorm:"column:target_admin_id"`
			OperatorAdminID uint64 `gorm:"column:operator_admin_id"`
			Reason          string `gorm:"column:reason"`
		}
		lookup := tx.Table("qixi_crm_b_customer_service_assignment_log").Where("binding_id = ? AND idempotency_key = ?", id, key).Limit(1).Find(&logRow)
		if lookup.Error != nil {
			return lookup.Error
		}
		if lookup.RowsAffected == 1 {
			if logRow.OperatorAdminID != operatorID {
				return errTransferForbidden
			}
			if logRow.TargetAdminID != in.TargetAdminID || logRow.Reason != in.Reason {
				return errTransferKeyConflict
			}
			return nil
		}
		if claims == nil || (!hasRole(claims.Roles, "platform") && (row.AssignedAdminID == nil || *row.AssignedAdminID != operatorID)) {
			return errTransferForbidden
		}
		if row.AssignedAdminID != nil && *row.AssignedAdminID == in.TargetAdminID {
			return errTransferSameOwner
		}
		update := tx.Table("qixi_crm_b_customer_service_binding").Where("id = ? AND status = ?", id, "open")
		if row.AssignedAdminID == nil {
			update = update.Where("assigned_admin_id IS NULL")
		} else {
			update = update.Where("assigned_admin_id = ?", *row.AssignedAdminID)
		}
		result := update.Updates(map[string]any{"assigned_admin_id": in.TargetAdminID, "assigned_at": time.Now()})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return errThreadTaken
		}
		return tx.Table("qixi_crm_b_customer_service_assignment_log").Create(map[string]any{
			"binding_id": id, "from_admin_id": row.AssignedAdminID, "target_admin_id": in.TargetAdminID,
			"operator_admin_id": operatorID, "reason": in.Reason, "idempotency_key": key,
		}).Error
	})
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			response.Fail(c, http.StatusNotFound, "客服会话不存在")
		case errThreadClosed, errThreadTaken, errTransferSameOwner, errTransferKeyConflict:
			response.Fail(c, http.StatusConflict, "客服会话转接冲突")
		case errTransferForbidden:
			response.Fail(c, http.StatusForbidden, "仅当前领取客服或平台可转接会话")
		case errTargetNotEligible:
			response.Fail(c, http.StatusBadRequest, "目标客服无该店铺服务范围")
		default:
			response.Fail(c, http.StatusInternalServerError, "转接客服会话失败")
		}
		return
	}
	if err := h.query(c, all, storeIDs).Where("b.id = ?", id).Scan(&row).Error; err != nil || row.ID == 0 {
		response.Fail(c, http.StatusInternalServerError, "读取转接结果失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) transferTargetEligible(ctx context.Context, adminID, storeID uint64) bool {
	if adminID == 0 {
		return false
	}
	var total int64
	if err := h.adminDB.WithContext(ctx).Table("qixi_crm_a_admin_user AS u").
		Joins("JOIN qixi_crm_a_admin_user_role AS ur ON ur.admin_user_id = u.id").
		Joins("JOIN qixi_crm_a_role AS r ON r.id = ur.role_id").
		Where("u.id = ? AND u.status = 1 AND u.deleted_at IS NULL AND r.code = ? AND r.status = 1", adminID, "customer_service").Count(&total).Error; err != nil || total == 0 {
		return false
	}
	stores, err := h.serviceQueueStores(ctx, adminID)
	return err == nil && (storeID == 0 || includesStore(stores, storeID))
}

func (h *Handler) Messages(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "会话编号错误")
		return
	}
	var binding serviceBinding
	if err := h.query(c, all, storeIDs).Where("b.id = ?", id).Scan(&binding).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服会话失败")
		return
	}
	if binding.ID == 0 {
		response.Fail(c, http.StatusNotFound, "客服会话不存在")
		return
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "50"), 50)
	if limit > 100 {
		limit = 100
	}
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_customer_service_message").Where("binding_id = ?", id)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服事件失败")
		return
	}
	var rows []messageRow
	if err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服事件失败")
		return
	}
	response.OK(c, gin.H{"conversation_id": binding.IMConversationID, "list": rows, "total": total, "page": page, "limit": limit, "note": "仅订单和系统事件投影；聊天正文由 pte-live-im 提供"})
}

// AssignmentLogs exposes transfer history only after the same scoped binding
// authorization as the thread itself. This prevents a caller from using a
// binding ID to enumerate the staffing history of an out-of-scope store.
func (h *Handler) AssignmentLogs(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 100
	}
	q := h.businessDB.WithContext(c.Request.Context()).
		Table("qixi_crm_b_customer_service_assignment_log").
		Where("binding_id = ?", binding.ID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会话转接审计失败")
		return
	}
	var rows []assignmentLogRow
	if err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询会话转接审计失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) quickReplyQuery(c *gin.Context, all bool, storeIDs []uint64) *gorm.DB {
	q := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_customer_service_quick_reply").Where("deleted_at IS NULL")
	if !all {
		q = q.Where("store_id IN ?", storeIDs)
	}
	return q
}

func (h *Handler) ListQuickReplies(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	page := positiveInt(c.DefaultQuery("page", "1"), 1)
	limit := positiveInt(c.DefaultQuery("limit", "20"), 20)
	if limit > 100 {
		limit = 100
	}
	q := h.quickReplyQuery(c, all, storeIDs)
	if keyword := strings.TrimSpace(c.Query("keyword")); keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	if raw := strings.TrimSpace(c.Query("store_id")); raw != "" {
		storeID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil || storeID == 0 || (!all && !includesStore(storeIDs, storeID)) {
			response.Fail(c, http.StatusForbidden, "无权查看该店铺的快捷回复")
			return
		}
		q = q.Where("store_id = ?", storeID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询快捷回复失败")
		return
	}
	var rows []quickReply
	if err := q.Order("updated_at DESC,id DESC").Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询快捷回复失败")
		return
	}
	response.OK(c, gin.H{"list": rows, "total": total, "page": page, "limit": limit})
}

func (h *Handler) CreateQuickReply(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	var in quickReplyInput
	if c.ShouldBindJSON(&in) != nil || !validQuickReply(&in) {
		response.Fail(c, http.StatusBadRequest, "快捷回复标题、内容或状态错误")
		return
	}
	if !all && !includesStore(storeIDs, in.StoreID) {
		response.Fail(c, http.StatusForbidden, "无权维护该店铺的快捷回复")
		return
	}
	row := quickReply{StoreID: in.StoreID, Title: in.Title, Content: in.Content, MessageType: in.MessageType, Status: in.Status, CreatedBy: uint64(middleware.AdminID(c)), UpdatedBy: uint64(middleware.AdminID(c))}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_customer_service_quick_reply").Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建快捷回复失败")
		return
	}
	if err := h.quickReplyQuery(c, all, storeIDs).Where("id = ?", row.ID).Take(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取快捷回复失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateQuickReply(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "快捷回复编号错误")
		return
	}
	var in quickReplyInput
	if c.ShouldBindJSON(&in) != nil || !validQuickReplyUpdate(&in) {
		response.Fail(c, http.StatusBadRequest, "快捷回复标题、内容或状态错误")
		return
	}
	if in.StoreID != 0 && !all && !includesStore(storeIDs, in.StoreID) {
		response.Fail(c, http.StatusForbidden, "无权关联该店铺的快捷回复")
		return
	}
	q := h.quickReplyQuery(c, all, storeIDs).Where("id = ?", id)
	updates := map[string]any{"title": in.Title, "content": in.Content, "message_type": in.MessageType, "status": in.Status, "updated_by": middleware.AdminID(c), "updated_at": time.Now()}
	if in.StoreID != 0 {
		updates["store_id"] = in.StoreID
	}
	result := q.Updates(updates)
	if result.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "更新快捷回复失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "快捷回复不存在或无权操作")
		return
	}
	var row quickReply
	if err := h.quickReplyQuery(c, all, storeIDs).Where("id = ?", id).Take(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "读取快捷回复失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteQuickReply(c *gin.Context) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "快捷回复编号错误")
		return
	}
	result := h.quickReplyQuery(c, all, storeIDs).Where("id = ?", id).Updates(map[string]any{"deleted_at": time.Now(), "updated_by": middleware.AdminID(c), "updated_at": time.Now()})
	if result.Error != nil {
		response.Fail(c, http.StatusInternalServerError, "删除快捷回复失败")
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, "快捷回复不存在或无权操作")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func validQuickReply(in *quickReplyInput) bool {
	if in == nil || in.StoreID == 0 {
		return false
	}
	return validQuickReplyUpdate(in)
}

func validQuickReplyUpdate(in *quickReplyInput) bool {
	if in == nil {
		return false
	}
	in.Title = strings.TrimSpace(in.Title)
	in.Content = strings.TrimSpace(in.Content)
	in.MessageType = strings.TrimSpace(in.MessageType)
	in.Status = strings.TrimSpace(in.Status)
	if in.MessageType == "" {
		in.MessageType = "text"
	}
	if in.Status == "" {
		in.Status = "enabled"
	}
	return in.Title != "" && len([]rune(in.Title)) <= 64 && in.Content != "" && len([]rune(in.Content)) <= 2000 && (in.MessageType == "text" || in.MessageType == "image") && (in.Status == "enabled" || in.Status == "disabled")
}

func validServiceSettings(in *serviceSettings) bool {
	if in == nil {
		return false
	}
	in.AutoReplyText = strings.TrimSpace(in.AutoReplyText)
	in.EnterpriseWechatCorpID = strings.TrimSpace(in.EnterpriseWechatCorpID)
	in.EnterpriseWechatURL = strings.TrimSpace(in.EnterpriseWechatURL)
	in.QueueMode = strings.TrimSpace(in.QueueMode)
	in.RedirectURL = strings.TrimSpace(in.RedirectURL)
	in.ServicePhone = strings.TrimSpace(in.ServicePhone)
	in.ServiceType = strings.TrimSpace(in.ServiceType)
	if in.ServiceType == "" {
		// Compatibility with the existing queue-only configuration.
		in.ServiceType = "system"
	}
	if in.QueueMode == "" {
		in.QueueMode = "manual"
	}
	if in.MaxSessionsPerAgent == 0 {
		in.MaxSessionsPerAgent = defaultServiceSettings().MaxSessionsPerAgent
	}
	if in.MaxSessionsPerAgent < 1 || in.MaxSessionsPerAgent > 200 || (in.QueueMode != "manual" && in.QueueMode != "round_robin") || len([]rune(in.AutoReplyText)) > 500 {
		return false
	}
	switch in.ServiceType {
	case "disabled", "mini_program":
		in.AutoReplyEnabled = false
	case "system":
		return !in.AutoReplyEnabled || in.AutoReplyText != ""
	case "phone":
		in.AutoReplyEnabled = false
		return validServicePhone(in.ServicePhone)
	case "enterprise_wechat":
		in.AutoReplyEnabled = false
		return validServiceURL(in.EnterpriseWechatURL) && validEnterpriseWechatCorpID(in.EnterpriseWechatCorpID)
	case "link":
		in.AutoReplyEnabled = false
		return validServiceURL(in.RedirectURL)
	default:
		return false
	}
	return true
}

func validServiceURL(raw string) bool {
	parsed, err := url.ParseRequestURI(raw)
	return err == nil && (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Host != ""
}

func validServicePhone(raw string) bool {
	if len([]rune(raw)) < 3 || len([]rune(raw)) > 32 {
		return false
	}
	digits := 0
	for _, char := range raw {
		if char >= '0' && char <= '9' {
			digits++
			continue
		}
		switch char {
		case '+', '-', ' ', '(', ')':
		default:
			return false
		}
	}
	return digits >= 3
}

func validEnterpriseWechatCorpID(raw string) bool {
	if len(raw) < 2 || len(raw) > 128 {
		return false
	}
	for _, char := range raw {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '_' || char == '-' {
			continue
		}
		return false
	}
	return true
}

func validUserNote(in *userNoteInput) bool {
	if in == nil {
		return false
	}
	in.Content = strings.TrimSpace(in.Content)
	return in.Content != "" && len([]rune(in.Content)) <= 500
}

func validTransfer(in *transferInput, idempotencyKey string) bool {
	if in == nil || in.TargetAdminID == 0 || idempotencyKey == "" || len([]rune(idempotencyKey)) > 128 {
		return false
	}
	in.Reason = strings.TrimSpace(in.Reason)
	return in.Reason != "" && len([]rune(in.Reason)) <= 500
}

func includesStore(storeIDs []uint64, storeID uint64) bool {
	for _, id := range storeIDs {
		if id == storeID {
			return true
		}
	}
	return false
}

func serviceStoreIDs(storeIDs []uint64) []uint64 {
	if storeIDs == nil {
		return []uint64{}
	}
	return storeIDs
}

func sharesStore(left, right []uint64) bool {
	for _, storeID := range left {
		if includesStore(right, storeID) {
			return true
		}
	}
	return false
}

func maskMobile(raw string) string {
	value := strings.TrimSpace(raw)
	if len([]rune(value)) < 7 {
		return ""
	}
	runes := []rune(value)
	return string(runes[:3]) + "****" + string(runes[len(runes)-4:])
}

func (h *Handler) scopedBinding(c *gin.Context) (*serviceBinding, bool) {
	all, storeIDs, ok := h.access(c)
	if !ok {
		response.Fail(c, http.StatusForbidden, "未配置客服队列数据范围")
		return nil, false
	}
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "会话编号错误")
		return nil, false
	}
	var binding serviceBinding
	if err := h.query(c, all, storeIDs).Where("b.id = ?", id).Scan(&binding).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询客服会话失败")
		return nil, false
	}
	if binding.ID == 0 {
		response.Fail(c, http.StatusNotFound, "客服会话不存在")
		return nil, false
	}
	return &binding, true
}

func (h *Handler) Order(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	if binding.OrderID == nil || *binding.OrderID == 0 {
		response.Fail(c, http.StatusNotFound, "该会话未关联订单")
		return
	}
	var row map[string]any
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order").
		Select("id,order_no,store_id,store_name_snapshot,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,status,recipient_snapshot,remark,created_at,updated_at").
		Where("id = ? AND store_id = ?", *binding.OrderID, binding.StoreID).Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Delivery(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	if binding.OrderID == nil || *binding.OrderID == 0 {
		response.Fail(c, http.StatusNotFound, "该会话未关联订单")
		return
	}
	var row map[string]any
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_delivery AS d").
		Select("d.id,d.order_id,d.delivery_type,d.carrier_code,d.tracking_no,d.status,d.delivered_at").
		Joins("INNER JOIN qixi_crm_b_order AS o ON o.id = d.order_id AND o.store_id = ?", binding.StoreID).
		Where("d.order_id = ?", *binding.OrderID).Order("d.id DESC").Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "订单暂无配送信息")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询订单配送信息失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Products(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	if binding.OrderID == nil || *binding.OrderID == 0 {
		response.OK(c, gin.H{"list": []any{}})
		return
	}
	var rows []map[string]any
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_order_item AS i").
		Select("i.id AS order_item_id,i.product_id,i.sku_key,i.title_snapshot,i.cover_url_snapshot,i.spec_snapshot,i.unit_price,i.quantity,i.refund_quantity,COALESCE(p.title, '') AS current_title,COALESCE(p.cover_url, '') AS current_cover_url,COALESCE(p.price, 0) AS current_price,COALESCE(p.stock, 0) AS current_stock,COALESCE(p.sale_status, 0) AS current_sale_status").
		Joins("LEFT JOIN qixi_crm_b_product_view AS p ON p.product_id = i.product_id AND p.store_id = ?", binding.StoreID).
		Where("i.order_id = ?", *binding.OrderID).Order("i.id ASC").Find(&rows).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询关联商品失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) Refunds(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	if binding.OrderID == nil || *binding.OrderID == 0 {
		response.OK(c, gin.H{"list": []any{}})
		return
	}
	var rows []map[string]any
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_refund AS r").
		Select("r.id,r.refund_no,r.order_id,r.reason,r.amount,r.order_status_before,r.status,r.created_at,r.updated_at").
		Joins("INNER JOIN qixi_crm_b_order AS o ON o.id = r.order_id AND o.store_id = ?", binding.StoreID).
		Where("r.order_id = ?", *binding.OrderID).Order("r.id DESC").Find(&rows).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询退款单失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) User(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	var row map[string]any
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_user AS u").
		Select("u.id,u.nickname,u.mobile,u.status,u.created_at,u.updated_at,COALESCE(p.avatar_url, '') AS avatar_url,COALESCE(p.bio, '') AS bio,COALESCE(p.source_channel, '') AS source_channel,COALESCE(n.content, '') AS service_note").
		Joins("LEFT JOIN qixi_crm_b_user_profile AS p ON p.user_id = u.id").
		Joins("LEFT JOIN qixi_crm_b_customer_service_user_note AS n ON n.user_id = u.id AND n.store_id = ?", binding.StoreID).
		Where("u.id = ?", binding.UserID).Take(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Fail(c, http.StatusNotFound, "用户不存在")
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询用户资料失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateUserNote(c *gin.Context) {
	binding, ok := h.scopedBinding(c)
	if !ok {
		return
	}
	var in userNoteInput
	if c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "用户备注参数错误")
		return
	}
	if !validUserNote(&in) {
		response.Fail(c, http.StatusBadRequest, "用户备注不能为空且不能超过 500 字")
		return
	}
	row := map[string]any{"user_id": binding.UserID, "store_id": binding.StoreID, "content": in.Content, "updated_by": middleware.AdminID(c)}
	err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_customer_service_user_note").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "store_id"}},
		DoUpdates: clause.Assignments(map[string]any{"content": in.Content, "updated_by": middleware.AdminID(c), "updated_at": time.Now()}),
	}).Create(row).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存用户备注失败")
		return
	}
	response.OK(c, gin.H{"user_id": binding.UserID, "store_id": binding.StoreID, "content": in.Content})
}

func hasRole(roles []string, expected string) bool {
	for _, role := range roles {
		if role == expected {
			return true
		}
	}
	return false
}

func positiveInt(raw string, fallback int) int {
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}

func parseID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return id, err == nil && id > 0
}

var (
	errThreadClosed           = &serviceError{"thread closed"}
	errThreadTaken            = &serviceError{"thread already claimed"}
	errTargetNotEligible      = &serviceError{"transfer target is not eligible"}
	errTransferForbidden      = &serviceError{"transfer forbidden"}
	errTransferKeyConflict    = &serviceError{"transfer idempotency key conflict"}
	errTransferSameOwner      = &serviceError{"transfer target is current owner"}
	errInvalidServiceSettings = &serviceError{"invalid customer service settings"}
)

func validateClaim(status string, assignedAdminID *uint64, adminID uint64) error {
	if status != "open" {
		return errThreadClosed
	}
	if adminID == 0 || (assignedAdminID != nil && *assignedAdminID != adminID) {
		return errThreadTaken
	}
	return nil
}

type serviceError struct{ text string }

func (e *serviceError) Error() string { return e.text }
