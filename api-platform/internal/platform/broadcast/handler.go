package broadcast

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/broadcast"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *broadcast.Service
	adminDB *gorm.DB
}

func NewHandler(svc *broadcast.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/broadcast/rooms", h.List)
	r.GET("/broadcast/rooms/:id", h.Get)
	r.PUT("/broadcast/rooms/:id/status", middleware.RequireAdminRoles("platform", "operations"), middleware.RequireAdminMenu(h.adminDB, "marketing.broadcast.audit"), h.UpdateStatus)
	r.PUT("/broadcast/rooms/:id/show", middleware.RequireAdminRoles("platform", "operations"), middleware.RequireAdminMenu(h.adminDB, "marketing.broadcast.audit"), h.UpdateShow)
	r.PUT("/broadcast/rooms/:id/recommend", middleware.RequireAdminRoles("platform", "operations"), middleware.RequireAdminMenu(h.adminDB, "marketing.broadcast.audit"), h.UpdateRecommend)
	r.DELETE("/broadcast/rooms/:id", middleware.RequireAdminRoles("platform", "operations"), middleware.RequireAdminMenu(h.adminDB, "marketing.broadcast.audit"), h.Delete)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := broadcast.ListFilter{
		Keyword: strings.TrimSpace(c.Query("keyword")),
	}
	if v := strings.TrimSpace(c.Query("mer_id")); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			mid := uint(id)
			filter.MerID = &mid
		}
	}
	if v := strings.TrimSpace(c.Query("status_tag")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			tag := int8(n)
			filter.StatusTag = &tag
		}
	} else if v := strings.TrimSpace(c.Query("status")); v != "" {
		// 兼容旧参数 status
		if n, err := strconv.Atoi(v); err == nil {
			tag := int8(n)
			filter.StatusTag = &tag
		}
	}
	if v := strings.TrimSpace(c.Query("show_type")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && (n == 0 || n == 1) {
			st := int8(n)
			filter.ShowType = &st
		}
	}
	if v := strings.TrimSpace(c.Query("live_status")); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ls := int16(n)
			filter.LiveStatus = &ls
		}
	}
	if v := strings.TrimSpace(c.Query("star")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 && n <= 5 {
			star := n
			filter.Star = &star
		}
	}
	if v := strings.TrimSpace(c.Query("is_trader")); v != "" {
		if n, err := strconv.Atoi(v); err == nil && (n == 0 || n == 1) {
			ids, err := h.lookupMerchantIDsByTrader(c, int8(n))
			if err != nil {
				response.Fail(c, http.StatusInternalServerError, "查询失败")
				return
			}
			filter.MerIDs = ids
		}
	}
	res, err := h.svc.ListPlatform(c.Request.Context(), filter, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	h.enrichMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	room, err := h.svc.Get(c.Request.Context(), uint(id), true)
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []broadcast.Room{*room}
	h.enrichMerchants(c, tmp)
	*room = tmp[0]
	response.OK(c, room)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in broadcast.AuditInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.Audit(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []broadcast.Room{*room}
	h.enrichMerchants(c, tmp)
	*room = tmp[0]
	response.OK(c, room)
}

func (h *Handler) UpdateShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsShow *int `json:"is_show"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.SetShow(c.Request.Context(), uint(id), int8(*body.IsShow))
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []broadcast.Room{*room}
	h.enrichMerchants(c, tmp)
	*room = tmp[0]
	response.OK(c, room)
}

func (h *Handler) UpdateRecommend(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in broadcast.RecommendInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	room, err := h.svc.SetRecommend(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []broadcast.Room{*room}
	h.enrichMerchants(c, tmp)
	*room = tmp[0]
	response.OK(c, room)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) lookupMerchantIDsByTrader(c *gin.Context, isTrader int8) ([]uint, error) {
	if h.adminDB == nil {
		return []uint{}, nil
	}
	var ids []uint
	err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_view").
		Where("is_trader = ?", isTrader).
		Pluck("merchant_id", &ids).Error
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []uint{}
	}
	return ids, nil
}

func (h *Handler) enrichMerchants(c *gin.Context, rows []broadcast.Room) {
	if h.adminDB == nil || len(rows) == 0 {
		return
	}
	ids := make([]uint, 0, len(rows))
	seen := map[uint]struct{}{}
	for _, row := range rows {
		if row.MerID == 0 {
			continue
		}
		if _, ok := seen[row.MerID]; ok {
			continue
		}
		seen[row.MerID] = struct{}{}
		ids = append(ids, row.MerID)
	}
	if len(ids) == 0 {
		return
	}
	type merRow struct {
		MerchantID   uint   `gorm:"column:merchant_id"`
		MerchantName string `gorm:"column:merchant_name"`
		IsTrader     int8   `gorm:"column:is_trader"`
	}
	list := make([]merRow, 0, len(ids))
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name, is_trader").
		Where("merchant_id IN ?", ids).Scan(&list)
	byID := map[uint]merRow{}
	for _, m := range list {
		byID[m.MerchantID] = m
	}
	for i := range rows {
		m, ok := byID[rows[i].MerID]
		if !ok {
			if rows[i].TraderName == "" {
				rows[i].TraderName = "非自营"
			}
			continue
		}
		if rows[i].MerName == "" {
			rows[i].MerName = m.MerchantName
		}
		rows[i].IsTrader = m.IsTrader
		if m.IsTrader == 1 {
			rows[i].TraderName = "自营"
		} else {
			rows[i].TraderName = "非自营"
		}
	}
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, broadcast.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, broadcast.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
