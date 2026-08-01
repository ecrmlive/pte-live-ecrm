// Package upload implements the C-end direct-to-COS upload handshake.
package upload

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	storage "github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/upload"
	"gorm.io/gorm"
)

const merchantApplicationLicense = "merchant_application_license"

type Storage interface {
	PresignPut(context.Context, storage.PresignInput) (*storage.PresignResult, error)
	Exists(context.Context, string) error
}

type Handler struct {
	db      *gorm.DB
	storage Storage
}

func NewHandler(db *gorm.DB, storage Storage) *Handler { return &Handler{db: db, storage: storage} }

func (h *Handler) Register(r gin.IRoutes) {
	r.POST("/uploads/presign", h.Presign)
	r.POST("/uploads/complete", h.Complete)
}

type presignRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size" binding:"required"`
	Purpose     string `json:"purpose" binding:"required"`
}

type completeRequest struct {
	ObjectKey string `json:"object_key" binding:"required"`
	Purpose   string `json:"purpose" binding:"required"`
}

type objectRow struct {
	ID          uint64     `gorm:"column:id;primaryKey"`
	OwnerUserID uint       `gorm:"column:owner_user_id"`
	Purpose     string     `gorm:"column:purpose"`
	ObjectKey   string     `gorm:"column:object_key"`
	Original    string     `gorm:"column:original_name"`
	ContentType string     `gorm:"column:content_type"`
	Size        int64      `gorm:"column:size"`
	Status      string     `gorm:"column:status"`
	ExpiresAt   time.Time  `gorm:"column:expires_at"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	CompletedAt *time.Time `gorm:"column:completed_at"`
}

func (objectRow) TableName() string { return "qixi_crm_b_upload_object" }

func (h *Handler) Presign(c *gin.Context) {
	var req presignRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Purpose != merchantApplicationLicense {
		response.Fail(c, http.StatusBadRequest, "上传参数错误")
		return
	}
	uid := middleware.UID(c)
	scope := "app/merchant-applications/" + uintString(uid)
	intent, err := h.storage.PresignPut(c.Request.Context(), storage.PresignInput{Scope: scope, Filename: req.Filename, ContentType: req.ContentType, Size: req.Size})
	if err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	row := objectRow{OwnerUserID: uid, Purpose: req.Purpose, ObjectKey: intent.ObjectKey, Original: strings.TrimSpace(req.Filename), ContentType: intent.Headers["Content-Type"], Size: req.Size, Status: "issued", ExpiresAt: intent.ExpiresAt}
	if err := h.db.WithContext(c.Request.Context()).Create(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存上传任务失败")
		return
	}
	response.OK(c, intent)
}

func (h *Handler) Complete(c *gin.Context) {
	var req completeRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Purpose != merchantApplicationLicense {
		response.Fail(c, http.StatusBadRequest, "上传参数错误")
		return
	}
	var row objectRow
	err := h.db.WithContext(c.Request.Context()).Where("owner_user_id = ? AND purpose = ? AND object_key = ?", middleware.UID(c), req.Purpose, strings.TrimSpace(req.ObjectKey)).First(&row).Error
	if err != nil {
		response.Fail(c, http.StatusNotFound, "上传任务不存在")
		return
	}
	if row.Status == "completed" {
		response.OK(c, gin.H{"object_key": row.ObjectKey})
		return
	}
	if row.Status != "issued" || time.Now().After(row.ExpiresAt) {
		response.Fail(c, http.StatusBadRequest, "上传地址已过期，请重新选择文件")
		return
	}
	if err := h.storage.Exists(c.Request.Context(), row.ObjectKey); err != nil {
		response.Fail(c, http.StatusBadRequest, "文件尚未成功上传")
		return
	}
	now := time.Now()
	if err := h.db.WithContext(c.Request.Context()).Model(&objectRow{}).Where("id = ? AND status = 'issued'", row.ID).Updates(map[string]any{"status": "completed", "completed_at": now}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存上传对象失败")
		return
	}
	response.OK(c, gin.H{"object_key": row.ObjectKey})
}

func uintString(value uint) string {
	const digits = "0123456789"
	if value == 0 {
		return "0"
	}
	buf := [20]byte{}
	i := len(buf)
	for value > 0 {
		i--
		buf[i] = digits[value%10]
		value /= 10
	}
	return string(buf[i:])
}
