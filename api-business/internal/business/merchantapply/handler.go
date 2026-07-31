// Package merchantapply owns the C-end merchant application write model.
// The platform console receives a NATS projection; this package never connects
// to qixi_crm_admin directly.
package merchantapply

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

const submittedEvent = "business.merchant_application.submitted"

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/merchant-applications", h.Create)
	r.GET("/merchant-applications/mine", h.ListMine)
}

type createRequest struct {
	MerchantName  string `json:"merchant_name" binding:"required"`
	ContactName   string `json:"contact_name" binding:"required"`
	ContactMobile string `json:"contact_mobile" binding:"required"`
	CategoryName  string `json:"category_name"`
	MerchantType  string `json:"merchant_type"`
	LicenseURL    string `json:"license_url"`
}

type applicationRow struct {
	ID              uint64    `gorm:"column:id" json:"id"`
	ApplicantUserID uint64    `gorm:"column:applicant_user_id" json:"applicant_user_id"`
	MerchantName    string    `gorm:"column:merchant_name" json:"merchant_name"`
	ContactName     string    `gorm:"column:contact_name" json:"contact_name"`
	ContactMobile   string    `gorm:"column:contact_mobile" json:"contact_mobile"`
	CategoryName    string    `gorm:"column:category_name" json:"category_name"`
	MerchantType    string    `gorm:"column:merchant_type" json:"merchant_type"`
	LicenseURL      string    `gorm:"column:license_url" json:"license_url"`
	Status          string    `gorm:"column:status" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

func (applicationRow) TableName() string { return "qixi_crm_b_merchant_application" }

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	uid := uint64(middleware.UID(c))
	row := applicationRow{ApplicantUserID: uid, MerchantName: strings.TrimSpace(req.MerchantName), ContactName: strings.TrimSpace(req.ContactName), ContactMobile: strings.TrimSpace(req.ContactMobile), CategoryName: strings.TrimSpace(req.CategoryName), MerchantType: strings.TrimSpace(req.MerchantType), LicenseURL: strings.TrimSpace(req.LicenseURL), Status: "pending"}
	if uid == 0 || len(row.MerchantName) < 2 || row.ContactName == "" || len(row.ContactMobile) < 6 || row.CategoryName == "" || row.MerchantType == "" || row.LicenseURL == "" {
		response.Fail(c, http.StatusBadRequest, "请完整填写入驻资料")
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		payload, err := json.Marshal(row)
		if err != nil {
			return err
		}
		return tx.Table("qixi_crm_b_outbox").Create(map[string]any{"event_type": submittedEvent, "aggregate_type": "merchant_application", "aggregate_id": row.ID, "payload": payload, "status": "pending"}).Error
	}); err != nil {
		response.Fail(c, http.StatusInternalServerError, "提交申请失败")
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListMine(c *gin.Context) {
	var rows []applicationRow
	if err := h.db.WithContext(c.Request.Context()).Where("applicant_user_id = ?", middleware.UID(c)).Order("id DESC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询申请记录失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}
