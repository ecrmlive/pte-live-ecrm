// Package payment manages one store's own payment merchant accounts. It never
// reads or falls back to the platform payment account.
package payment

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/pkg/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	channelWechat = "wechat"
	channelAlipay = "alipay"
)

type paymentChannel struct {
	StoreID   uint   `gorm:"column:store_id;primaryKey"`
	Channel   string `gorm:"column:channel;primaryKey" json:"channel"`
	Enabled   bool   `gorm:"column:enabled" json:"enabled"`
	UpdatedBy uint   `gorm:"column:updated_by"`
}

func (paymentChannel) TableName() string { return "qixi_crm_m_payment_channel" }

type Handler struct {
	db         *gorm.DB
	businessDB *gorm.DB
	configs    *paymentconfig.Store
}

// businessDB is a one-way encrypted projection target. It lets api-business
// use this store's own account without querying the merchant database.
func NewHandler(db, businessDB *gorm.DB, masterSecret string) *Handler {
	configs, _ := paymentconfig.NewStore(businessDB, masterSecret)
	return &Handler{db: db, businessDB: businessDB, configs: configs}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/payment/channels", h.List)
	r.PUT("/payment/channels/:channel", h.Save)
}

// List always returns both supported channels so the Vben UI has a stable
// contract before a store has saved any choice.
func (h *Handler) List(c *gin.Context) {
	storeID := middleware.StoreID(c)
	var rows []paymentChannel
	if err := h.db.WithContext(c.Request.Context()).Where("store_id = ?", storeID).Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询支付方式失败")
		return
	}
	byChannel := make(map[string]paymentChannel, len(rows))
	for _, row := range rows {
		byChannel[row.Channel] = row
	}
	var configs []paymentConfigRow
	if err := h.db.WithContext(c.Request.Context()).Where("store_id = ?", storeID).Find(&configs).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺支付配置失败")
		return
	}
	configured := make(map[string]bool, len(configs))
	for _, row := range configs {
		configured[row.Channel] = row.Ciphertext != ""
	}
	response.OK(c, gin.H{"list": []gin.H{
		{"channel": channelWechat, "enabled": byChannel[channelWechat].Enabled, "configured": configured[channelWechat]},
		{"channel": channelAlipay, "enabled": byChannel[channelAlipay].Enabled, "configured": configured[channelAlipay]},
	}})
}

type saveRequest struct {
	Values map[string]string `json:"values" binding:"required"`
}

func (h *Handler) Save(c *gin.Context) {
	channel, ok := normalizeChannel(c.Param("channel"))
	if !ok {
		response.Fail(c, http.StatusBadRequest, "不支持的支付方式")
		return
	}
	var req saveRequest
	if err := c.ShouldBindJSON(&req); err != nil || len(req.Values) == 0 || h.configs == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if !validValues(channel, req.Values) {
		response.Fail(c, http.StatusBadRequest, "支付配置字段不完整或不合法")
		return
	}
	ciphertext, err := h.configs.Seal(req.Values)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "支付配置加密失败")
		return
	}
	enabled := req.Values["enabled"] == "true"
	row := paymentChannel{StoreID: middleware.StoreID(c), Channel: channel, Enabled: enabled, UpdatedBy: middleware.AdminID(c)}
	err = h.db.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "channel"}},
		DoUpdates: clause.AssignmentColumns([]string{"enabled", "updated_by", "updated_at"}),
	}).Create(&row).Error
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存支付方式失败")
		return
	}
	configRow := paymentConfigRow{StoreID: row.StoreID, Channel: row.Channel, Ciphertext: ciphertext, UpdatedBy: row.UpdatedBy}
	if err := h.db.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "store_id"}, {Name: "channel"}},
		DoUpdates: clause.AssignmentColumns([]string{"ciphertext", "updated_by", "updated_at"}),
	}).Create(&configRow).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存店铺支付配置失败")
		return
	}
	if h.businessDB != nil {
		if err := h.configs.PublishStore(c.Request.Context(), row.StoreID, req.Values); err != nil {
			response.Fail(c, http.StatusInternalServerError, "支付配置同步失败，请重试保存")
			return
		}
		if err := h.businessDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "store_id"}, {Name: "channel"}},
			DoUpdates: clause.AssignmentColumns([]string{"enabled", "updated_at"}),
		}).Create(&storePaymentChannel{StoreID: row.StoreID, Channel: row.Channel, Enabled: row.Enabled}).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "支付方式同步失败，请重试保存")
			return
		}
	}
	response.OK(c, gin.H{"channel": channel, "enabled": enabled, "configured": true})
}

type paymentConfigRow struct {
	StoreID    uint   `gorm:"column:store_id;primaryKey"`
	Channel    string `gorm:"column:channel;primaryKey"`
	Ciphertext string `gorm:"column:ciphertext"`
	UpdatedBy  uint   `gorm:"column:updated_by"`
}

func (paymentConfigRow) TableName() string { return "qixi_crm_m_payment_config" }

type storePaymentChannel struct {
	StoreID uint   `gorm:"column:store_id;primaryKey"`
	Channel string `gorm:"column:channel;primaryKey"`
	Enabled bool   `gorm:"column:enabled"`
}

func (storePaymentChannel) TableName() string { return "qixi_crm_b_store_payment_channel" }

func normalizeChannel(value string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case channelWechat:
		return channelWechat, true
	case channelAlipay:
		return channelAlipay, true
	default:
		return "", false
	}
}

func validValues(channel string, values map[string]string) bool {
	if values["enabled"] != "true" && values["enabled"] != "false" {
		return false
	}
	required := []string{"notify_url"}
	switch channel {
	case channelWechat:
		required = append(required, "app_id", "mch_id", "api_v3_key", "serial_no", "private_key")
	case channelAlipay:
		required = append(required, "app_id", "private_key", "public_key")
	}
	if values["enabled"] == "false" {
		return true
	}
	for _, key := range required {
		if strings.TrimSpace(values[key]) == "" || len(values[key]) > 16*1024 {
			return false
		}
	}
	return true
}
