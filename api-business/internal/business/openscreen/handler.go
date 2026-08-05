package openscreen

import (
	"errors"
	"net/http"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) { r.GET("/open-screen", h.Get) }

type campaign struct {
	ID          uint64 `gorm:"column:id"`
	Title       string `gorm:"column:title"`
	ImageURL    string `gorm:"column:image_url"`
	LinkURL     string `gorm:"column:link_url"`
	DurationSec int    `gorm:"column:duration_sec"`
	SpaceHours  int    `gorm:"column:space_hours"`
	Enabled     int8   `gorm:"column:enabled"`
}

func (campaign) TableName() string { return "qixi_crm_b_open_screen_campaign" }

func (h *Handler) Get(c *gin.Context) {
	var row campaign
	err := h.db.WithContext(c.Request.Context()).Where("enabled = ?", 1).Order("id DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.OK(c, gin.H{"enabled": false})
		return
	}
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "开屏广告服务异常")
		return
	}
	if strings.TrimSpace(row.ImageURL) == "" || row.DurationSec < 1 || row.DurationSec > 10 || row.SpaceHours < 0 || row.SpaceHours > 24*365 {
		response.Fail(c, http.StatusInternalServerError, "开屏广告配置无效")
		return
	}
	response.OK(c, gin.H{"enabled": true, "id": row.ID, "title": row.Title, "image": row.ImageURL, "url": row.LinkURL, "duration": row.DurationSec, "space_hours": row.SpaceHours})
}
