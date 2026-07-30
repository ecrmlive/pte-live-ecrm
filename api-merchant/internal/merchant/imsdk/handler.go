// Package imsdk manages a merchant's pte-live-im SDK AppId bindings.
// It only stores non-secret mapping metadata. UserSig remains issued by pte-live-im.
package imsdk

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	statusEnabled  = "enabled"
	statusDisabled = "disabled"
)

type binding struct {
	ID           uint   `gorm:"column:id;primaryKey" json:"id"`
	MerchantID   uint   `gorm:"column:merchant_id" json:"merchant_id"`
	SDKAppID     string `gorm:"column:sdk_app_id" json:"sdk_app_id"`
	Name         string `gorm:"column:name" json:"name"`
	Status       string `gorm:"column:status" json:"status"`
	IsActive     bool   `gorm:"column:is_active" json:"is_active"`
	APIPublicURL string `gorm:"column:api_public_url" json:"api_public_url"`
	WSPublicURL  string `gorm:"column:ws_public_url" json:"ws_public_url"`
	PTEProfileID string `gorm:"column:pte_profile_id" json:"pte_profile_id"`
	CreatedBy    uint   `gorm:"column:created_by" json:"-"`
}

func (binding) TableName() string { return "qixi_crm_m_im_sdk_app" }

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/settings/im-sdk-apps", h.List)
	r.POST("/settings/im-sdk-apps", h.Create)
	r.PUT("/settings/im-sdk-apps/:id", h.Update)
	r.POST("/settings/im-sdk-apps/:id/activate", h.Activate)
	r.POST("/settings/im-sdk-apps/:id/disable", h.Disable)
}

func (h *Handler) List(c *gin.Context) {
	var rows []binding
	if err := h.db.WithContext(c.Request.Context()).Where("merchant_id = ?", middleware.MerID(c)).Order("is_active DESC, id DESC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询 IM 应用失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

type saveRequest struct {
	SDKAppID     string `json:"sdk_app_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Status       string `json:"status"`
	APIPublicURL string `json:"api_public_url"`
	WSPublicURL  string `json:"ws_public_url"`
	PTEProfileID string `json:"pte_profile_id"`
}

func (h *Handler) Create(c *gin.Context) {
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "IM SDK AppId 配置不合法")
		return
	}
	row := binding{
		MerchantID: middleware.MerID(c), SDKAppID: strings.TrimSpace(req.SDKAppID), Name: strings.TrimSpace(req.Name),
		Status: normalizeStatus(req.Status), APIPublicURL: strings.TrimSpace(req.APIPublicURL),
		WSPublicURL: strings.TrimSpace(req.WSPublicURL), PTEProfileID: strings.TrimSpace(req.PTEProfileID), CreatedBy: middleware.AdminID(c),
	}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusConflict, "IM SDK AppId 已存在")
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "应用标识错误")
		return
	}
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || !valid(req) {
		response.Fail(c, http.StatusBadRequest, "IM SDK AppId 配置不合法")
		return
	}
	updates := map[string]any{
		"name": strings.TrimSpace(req.Name), "status": normalizeStatus(req.Status),
		"api_public_url": strings.TrimSpace(req.APIPublicURL), "ws_public_url": strings.TrimSpace(req.WSPublicURL),
		"pte_profile_id": strings.TrimSpace(req.PTEProfileID),
	}
	if updates["status"] == statusDisabled {
		updates["is_active"] = false
	}
	merchantID := middleware.MerID(c)
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var current binding
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND merchant_id = ?", id, merchantID).First(&current).Error; err != nil {
			return err
		}
		if err := tx.Model(&binding{}).Where("id = ? AND merchant_id = ?", id, merchantID).Updates(updates).Error; err != nil {
			return err
		}
		if !current.IsActive {
			return nil
		}
		var next binding
		if err := tx.Where("id = ? AND merchant_id = ?", id, merchantID).First(&next).Error; err != nil {
			return err
		}
		if next.Status == statusEnabled {
			return enqueueIMEvent(tx, "merchant.im_sdk_app.activated", next)
		}
		return enqueueIMEvent(tx, "merchant.im_sdk_app.deactivated", next)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "IM 应用不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新 IM 应用失败")
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *Handler) Activate(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "应用标识错误")
		return
	}
	merchantID := middleware.MerID(c)
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var target binding
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND merchant_id = ? AND status = ?", id, merchantID, statusEnabled).First(&target).Error; err != nil {
			return err
		}
		if err := tx.Model(&binding{}).Where("merchant_id = ?", merchantID).Update("is_active", false).Error; err != nil {
			return err
		}
		if err := tx.Model(&binding{}).Where("id = ? AND merchant_id = ?", id, merchantID).Update("is_active", true).Error; err != nil {
			return err
		}
		return enqueueIMEvent(tx, "merchant.im_sdk_app.activated", target)
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "可启用的 IM 应用不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "启用 IM 应用失败")
		return
	}
	response.OK(c, gin.H{"id": id, "is_active": true})
}

func (h *Handler) Disable(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Fail(c, http.StatusBadRequest, "应用标识错误")
		return
	}
	merchantID := middleware.MerID(c)
	err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		var current binding
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ? AND merchant_id = ?", id, merchantID).First(&current).Error; err != nil {
			return err
		}
		if err := tx.Model(&binding{}).Where("id = ? AND merchant_id = ?", id, merchantID).Updates(map[string]any{"status": statusDisabled, "is_active": false}).Error; err != nil {
			return err
		}
		if current.IsActive {
			current.Status = statusDisabled
			current.IsActive = false
			return enqueueIMEvent(tx, "merchant.im_sdk_app.deactivated", current)
		}
		return nil
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "IM 应用不存在")
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "停用 IM 应用失败")
		return
	}
	response.OK(c, gin.H{"id": id, "is_active": false})
}

// outboxEvent carries no PTE token or UserSig. A NATS dispatcher can replay pending
// records safely and api-business treats the event as the only merchant IM projection source.
type outboxEvent struct {
	EventType     string `gorm:"column:event_type"`
	AggregateType string `gorm:"column:aggregate_type"`
	AggregateID   string `gorm:"column:aggregate_id"`
	Payload       []byte `gorm:"column:payload"`
}

func (outboxEvent) TableName() string { return "qixi_crm_m_outbox" }

func enqueueIMEvent(tx *gorm.DB, eventType string, row binding) error {
	payload, err := imEventPayload(row)
	if err != nil {
		return err
	}
	return tx.Create(&outboxEvent{
		EventType: eventType, AggregateType: "merchant_im_sdk_app",
		AggregateID: strconv.FormatUint(uint64(row.MerchantID), 10), Payload: payload,
	}).Error
}

func imEventPayload(row binding) ([]byte, error) {
	return json.Marshal(gin.H{
		"merchant_id":    row.MerchantID,
		"sdk_app_id":     row.SDKAppID,
		"status":         row.Status,
		"is_active":      row.IsActive,
		"api_public_url": row.APIPublicURL,
		"ws_public_url":  row.WSPublicURL,
		"pte_profile_id": row.PTEProfileID,
	})
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	return uint(id), err == nil && id > 0
}

func normalizeStatus(value string) string {
	if strings.EqualFold(strings.TrimSpace(value), statusEnabled) {
		return statusEnabled
	}
	return statusDisabled
}

func valid(req saveRequest) bool {
	return len(strings.TrimSpace(req.SDKAppID)) <= 64 && strings.TrimSpace(req.SDKAppID) != "" &&
		len(strings.TrimSpace(req.Name)) <= 128 && strings.TrimSpace(req.Name) != "" &&
		len(strings.TrimSpace(req.APIPublicURL)) <= 1024 && len(strings.TrimSpace(req.WSPublicURL)) <= 1024 &&
		len(strings.TrimSpace(req.PTEProfileID)) <= 128
}
