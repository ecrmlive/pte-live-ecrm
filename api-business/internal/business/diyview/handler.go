// Package diyview exposes published DIY pages from the business projection.
// It deliberately has no merchant/admin database dependency.
package diyview

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/diy/home", h.Home)
	r.GET("/diy/pages/:id", h.Get)
}

type pageView struct {
	Source   string          `gorm:"column:source"`
	PageID   uint64          `gorm:"column:page_id"`
	StoreID  uint64          `gorm:"column:store_id"`
	PageType string          `gorm:"column:page_type"`
	Name     string          `gorm:"column:name"`
	Document json.RawMessage `gorm:"column:document"`
	Status   string          `gorm:"column:status"`
	IsActive bool            `gorm:"column:is_active"`
}

func (pageView) TableName() string { return "qixi_crm_b_diy_page_view" }

func (h *Handler) Home(c *gin.Context) {
	storeID, ok := h.resolveStore(c)
	if !ok {
		return
	}
	var row pageView
	query := h.db.WithContext(c.Request.Context()).Where("page_type = ? AND status = ? AND is_active = ?", "home", "published", true)
	if storeID == 0 {
		query = query.Where("source = ? AND store_id = ?", "platform", 0)
	} else {
		query = query.Where("source = ? AND store_id = ?", "merchant", storeID)
	}
	if err := query.Order("page_id DESC").First(&row).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.OK(c, emptyPayload())
			return
		}
		response.Fail(c, http.StatusInternalServerError, "查询装修页失败")
		return
	}
	response.OK(c, payload(row))
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "装修页标识错误")
		return
	}
	storeID, ok := h.resolveStore(c)
	if !ok {
		return
	}
	var row pageView
	query := h.db.WithContext(c.Request.Context()).Where("page_id = ? AND status = ?", id, "published")
	if storeID == 0 {
		query = query.Where("source = ? AND store_id = ?", "platform", 0)
	} else {
		query = query.Where("source = ? AND store_id = ?", "merchant", storeID)
	}
	if err := query.First(&row).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "装修页不存在")
		return
	}
	response.OK(c, payload(row))
}

func (h *Handler) resolveStore(c *gin.Context) (uint64, bool) {
	appID := c.GetHeader("X-AppId")
	if appID == "" {
		return 0, true
	}
	var row struct {
		StoreID uint64 `gorm:"column:store_id"`
	}
	if err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_store_view").Where("store_app_id = ? AND status = ?", appID, 1).First(&row).Error; err != nil {
		response.Fail(c, http.StatusNotFound, "店铺不存在或已关闭")
		return 0, false
	}
	return row.StoreID, true
}

func emptyPayload() gin.H {
	return gin.H{"id": 0, "title": "", "name": "", "page": map[string]any{}, "items": []any{}, "banners": []any{}, "menus": []any{}}
}

func payload(row pageView) gin.H {
	var doc struct {
		Page  map[string]any   `json:"page"`
		Items []map[string]any `json:"items"`
		Meta  map[string]any   `json:"_qixi"`
	}
	if json.Unmarshal(row.Document, &doc) != nil {
		doc.Page, doc.Items = map[string]any{}, []map[string]any{}
	}
	if doc.Page == nil {
		doc.Page = map[string]any{}
	}
	if doc.Items == nil {
		doc.Items = []map[string]any{}
	}
	title, _ := doc.Meta["title"].(string)
	if title == "" {
		title = row.Name
	}
	banners, menus := legacyCompat(doc.Items)
	return gin.H{"id": row.PageID, "store_id": row.StoreID, "name": row.Name, "title": title, "page_type": row.PageType, "page": doc.Page, "items": doc.Items, "banners": banners, "menus": menus, "value": json.RawMessage(row.Document)}
}

func legacyCompat(items []map[string]any) ([]map[string]any, []map[string]any) {
	banners, menus := []map[string]any{}, []map[string]any{}
	for _, item := range items {
		typ, _ := item["type"].(string)
		data, _ := item["data"].([]any)
		if data == nil {
			if list, ok := item["data"].([]map[string]any); ok {
				for _, entry := range list {
					data = append(data, entry)
				}
			}
		}
		for i, entry := range data {
			value, _ := entry.(map[string]any)
			if value == nil {
				continue
			}
			if typ == "banner" {
				banners = append(banners, map[string]any{"id": i + 1, "title": value["imgName"], "image": value["imgUrl"], "url": value["linkUrl"]})
			}
			if typ == "navBar" || typ == "option" {
				text := value["text"]
				if text == nil {
					text = value["title"]
				}
				menus = append(menus, map[string]any{"id": i + 1, "name": text, "icon": value["imgUrl"], "url": value["linkUrl"]})
			}
		}
	}
	return banners, menus
}
