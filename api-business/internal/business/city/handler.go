// Package city exposes the published administrative-area projection to C-end clients.
package city

import (
	"net/http"
	"strconv"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler only reads the business-side city projection. It intentionally never
// queries the platform database directly.
type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/system/city/lst", h.List)
	r.GET("/system/city/lst/:pid", h.List)
}

type cityView struct {
	CityID   uint64 `gorm:"column:city_id"`
	ParentID uint64 `gorm:"column:parent_id"`
	Name     string `gorm:"column:name"`
	Level    int8   `gorm:"column:level"`
	IsShow   int8   `gorm:"column:is_show"`
}

func (cityView) TableName() string { return "qixi_crm_b_city_view" }

func (h *Handler) List(c *gin.Context) {
	parentID := uint64(0)
	raw := c.Param("pid")
	if raw == "" {
		raw = c.Query("parent_id")
	}
	if raw != "" {
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "城市父级 ID 错误")
			return
		}
		parentID = parsed
	}
	var rows []cityView
	if err := h.db.WithContext(c.Request.Context()).
		Where("parent_id = ? AND is_show = ?", parentID, 1).
		Order("city_id ASC").Find(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "城市数据服务异常")
		return
	}
	items := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		items = append(items, gin.H{"city_id": row.CityID, "parent_id": row.ParentID, "name": row.Name, "level": row.Level})
	}
	response.OK(c, gin.H{"list": items, "parent_id": parentID})
}
