// Package customerservice implements the unified-admin service queue.
// Message transport and credentials intentionally remain owned by pte-live-im.
package customerservice

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct {
	adminDB    *gorm.DB
	businessDB *gorm.DB
}

func NewHandler(adminDB, businessDB *gorm.DB) *Handler {
	return &Handler{adminDB: adminDB, businessDB: businessDB}
}

func (h *Handler) Register(routes gin.IRoutes) {
	routes.GET("/customer-service/threads", h.List)
	routes.GET("/customer-service/threads/:id", h.Detail)
	routes.POST("/customer-service/threads/:id/claim", h.Claim)
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

type dataScope struct {
	StoreIDs []uint64 `json:"store_ids"`
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
	var rows []struct {
		ScopeValue json.RawMessage `gorm:"column:scope_value"`
	}
	if err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_data_scope").
		Select("scope_value").
		Where("admin_user_id = ? AND scope_type = ?", claims.AdminID, "service_queue").
		Find(&rows).Error; err != nil {
		return false, nil, false
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
		return false, nil, false
	}
	storeIDs = make([]uint64, 0, len(set))
	for id := range set {
		storeIDs = append(storeIDs, id)
	}
	return false, storeIDs, true
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
			Where("b.id = ?", id)
		if !all {
			q = q.Where("b.store_id IN ?", storeIDs)
		}
		if err := q.Take(&row).Error; err != nil {
			return err
		}
		if row.Status != "open" {
			return errThreadClosed
		}
		if row.AssignedAdminID != nil && *row.AssignedAdminID != uint64(middleware.AdminID(c)) {
			return errThreadTaken
		}
		return tx.Table("qixi_crm_b_customer_service_binding").Where("id = ?", id).
			Updates(map[string]any{"assigned_admin_id": middleware.AdminID(c), "assigned_at": time.Now()}).Error
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
	errThreadClosed = &serviceError{"thread closed"}
	errThreadTaken  = &serviceError{"thread already claimed"}
)

type serviceError struct{ text string }

func (e *serviceError) Error() string { return e.text }
