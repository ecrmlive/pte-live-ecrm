// Package live exposes the C-end read model for live rooms.  It deliberately
// reads qixi_crm_business only: stream publishing and merchant administration
// remain merchant-side concerns, while clients receive at most a public play URL.
package live

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
	"gorm.io/gorm"
)

type Handler struct{ db *gorm.DB }

func NewHandler(db *gorm.DB) *Handler { return &Handler{db: db} }

func (h *Handler) RegisterPublic(r gin.IRoutes) {
	r.GET("/live/rooms", h.List)
	r.GET("/live/rooms/:id", h.Get)
}

func (h *Handler) RegisterAuthed(r gin.IRoutes) {
	r.POST("/live/rooms/:id/reservation", h.Reserve)
	r.DELETE("/live/rooms/:id/reservation", h.CancelReservation)
}

type roomRow struct {
	ID           uint64     `gorm:"column:id"`
	MerchantID   uint64     `gorm:"column:merchant_id"`
	StoreID      uint64     `gorm:"column:store_id"`
	Title        string     `gorm:"column:title"`
	AnchorName   string     `gorm:"column:anchor_name"`
	CoverURL     string     `gorm:"column:cover_url"`
	Status       string     `gorm:"column:status"`
	PlayURL      string     `gorm:"column:play_url"`
	StartsAt     *time.Time `gorm:"column:starts_at"`
	EndedAt      *time.Time `gorm:"column:ended_at"`
	MerchantName string     `gorm:"column:merchant_name"`
	StoreName    string     `gorm:"column:store_name"`
}

func (roomRow) TableName() string { return "qixi_crm_b_live_room" }

type liveProductRow struct {
	ProductID uint64  `gorm:"column:product_id"`
	Title     string  `gorm:"column:title"`
	CoverURL  string  `gorm:"column:cover_url"`
	Price     float64 `gorm:"column:price"`
}

type reservationRow struct {
	LiveRoomID uint64 `gorm:"column:live_room_id"`
	UserID     uint64 `gorm:"column:user_id"`
}

func (reservationRow) TableName() string { return "qixi_crm_b_live_reservation" }

func (h *Handler) List(c *gin.Context) {
	page, limit := paging(c)
	var total int64
	base := h.rooms(c).Where("qixi_crm_b_live_room.is_public = 1 AND qixi_crm_b_live_room.status IN ?", []string{"scheduled", "living", "ended"})
	if err := base.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "直播间加载失败")
		return
	}
	var rows []roomRow
	if err := base.Order("qixi_crm_b_live_room.status = 'living' DESC, qixi_crm_b_live_room.sort DESC, qixi_crm_b_live_room.id DESC").
		Offset((page - 1) * limit).Limit(limit).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "直播间加载失败")
		return
	}
	list := make([]gin.H, 0, len(rows))
	for _, row := range rows {
		list = append(list, roomResponse(row, nil))
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": page, "limit": limit})
}

func (h *Handler) Get(c *gin.Context) {
	id, ok := roomID(c)
	if !ok {
		return
	}
	row, err := h.findPublicRoom(c, id)
	if err != nil {
		writeRoomError(c, err)
		return
	}
	goods, err := h.goods(c, row.ID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "直播商品加载失败")
		return
	}
	response.OK(c, roomResponse(*row, goods))
}

func (h *Handler) Reserve(c *gin.Context) {
	id, ok := roomID(c)
	if !ok {
		return
	}
	if _, err := h.findPublicRoom(c, id); err != nil {
		writeRoomError(c, err)
		return
	}
	row := reservationRow{LiveRoomID: id, UserID: uint64(middleware.UID(c))}
	if err := h.db.WithContext(c.Request.Context()).Where("live_room_id = ? AND user_id = ?", row.LiveRoomID, row.UserID).FirstOrCreate(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "预约失败")
		return
	}
	response.OK(c, gin.H{"live_room_id": id, "reserved": true})
}

func (h *Handler) CancelReservation(c *gin.Context) {
	id, ok := roomID(c)
	if !ok {
		return
	}
	if err := h.db.WithContext(c.Request.Context()).Where("live_room_id = ? AND user_id = ?", id, uint64(middleware.UID(c))).Delete(&reservationRow{}).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "取消预约失败")
		return
	}
	response.OK(c, gin.H{"live_room_id": id, "reserved": false})
}

func (h *Handler) rooms(c *gin.Context) *gorm.DB {
	return h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_live_room").
		Select("qixi_crm_b_live_room.id, qixi_crm_b_live_room.merchant_id, qixi_crm_b_live_room.store_id, qixi_crm_b_live_room.title, qixi_crm_b_live_room.anchor_name, qixi_crm_b_live_room.cover_url, qixi_crm_b_live_room.status, qixi_crm_b_live_room.play_url, qixi_crm_b_live_room.starts_at, qixi_crm_b_live_room.ended_at, qixi_crm_b_store_view.merchant_name, qixi_crm_b_store_view.store_name").
		Joins("LEFT JOIN qixi_crm_b_store_view ON qixi_crm_b_store_view.store_id = qixi_crm_b_live_room.store_id")
}

func (h *Handler) findPublicRoom(c *gin.Context, id uint64) (*roomRow, error) {
	var row roomRow
	err := h.rooms(c).Where("qixi_crm_b_live_room.id = ? AND qixi_crm_b_live_room.is_public = 1 AND qixi_crm_b_live_room.status IN ?", id, []string{"scheduled", "living", "ended"}).Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (h *Handler) goods(c *gin.Context, roomID uint64) ([]liveProductRow, error) {
	var rows []liveProductRow
	err := h.db.WithContext(c.Request.Context()).Table("qixi_crm_b_live_room_product AS rel").
		Select("rel.product_id, product.title, product.cover_url, product.price").
		Joins("JOIN qixi_crm_b_product_view AS product ON product.product_id = rel.product_id AND product.sale_status = 1").
		Where("rel.live_room_id = ?", roomID).Order("rel.sort ASC, rel.product_id ASC").Scan(&rows).Error
	return rows, err
}

func roomResponse(row roomRow, goods []liveProductRow) gin.H {
	items := make([]gin.H, 0, len(goods))
	for _, item := range goods {
		items = append(items, gin.H{"product_id": item.ProductID, "store_name": item.Title, "image": item.CoverURL, "price": item.Price, "on_sale": 1})
	}
	playURL := ""
	if row.Status == "living" {
		playURL = row.PlayURL
	}
	return gin.H{
		"id": row.ID, "broadcast_room_id": row.ID, "merchant_id": row.MerchantID, "store_id": row.StoreID,
		"name": row.Title, "title": row.Title, "anchor_name": row.AnchorName, "cover_img": row.CoverURL,
		"mer_name": row.MerchantName, "store_name": row.StoreName, "live_status": legacyStatus(row.Status),
		"status": row.Status, "play_url": playURL, "starts_at": row.StartsAt, "ended_at": row.EndedAt, "goods": items,
	}
}

func legacyStatus(status string) int {
	switch status {
	case "living":
		return 101
	case "ended":
		return 103
	default:
		return 102
	}
}

func paging(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit
}

func roomID(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Fail(c, http.StatusBadRequest, "直播间 ID 错误")
		return 0, false
	}
	return id, true
}

func writeRoomError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusNotFound, "直播间不存在或未公开")
		return
	}
	response.Fail(c, http.StatusInternalServerError, "直播间加载失败")
}
