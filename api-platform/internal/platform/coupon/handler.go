package coupon

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/promotion"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *promotion.Service
	adminDB *gorm.DB
}

func NewHandler(svc *promotion.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/coupons", h.List)
	sendRead := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.coupon.send.read"),
	}
	// /coupons/sends、/coupons/store 须注册在 /coupons/:id 之前
	r.GET("/coupons/sends", append(sendRead, h.ListSends)...)
	r.GET("/coupons/sends/:id", append(sendRead, h.SendDetail)...)
	r.GET("/coupons/sends/:id/users", append(sendRead, h.ListSendUsers)...)
	r.GET("/coupons/store", h.ListStore)
	r.GET("/coupons/store/:id", h.StoreDetail)
	r.GET("/coupons/:id", h.Detail)
	write := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.coupon.manage")
	r.POST("/coupons", write, manage, h.Create)
	r.PUT("/coupons/:id", write, manage, h.Update)
	r.DELETE("/coupons/:id", write, manage, h.Delete)
	r.POST("/coupons/:id/status", write, manage, h.Status)
	r.POST("/coupons/:id/clone", write, manage, h.Clone)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := promotion.CouponListFilter{
		Keyword: strings.TrimSpace(c.Query("keyword")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		st := int8(v)
		filter.Status = &st
	}
	if raw := strings.TrimSpace(c.Query("send_type")); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "获取方式参数错误")
			return
		}
		st := int8(v)
		filter.SendType = &st
	}
	res, err := h.svc.ListAdmin(c.Request.Context(), 0, promotion.CouponTypePlatform, page, limit, filter)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.GetAdminDetail(c.Request.Context(), 0, promotion.CouponTypePlatform, uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListStore(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := promotion.CouponListFilter{
		Keyword: strings.TrimSpace(c.Query("keyword")),
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		st := int8(v)
		filter.Status = &st
	}
	if raw := strings.TrimSpace(c.Query("is_trader")); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "店铺类别参数错误")
			return
		}
		ids, err := h.merchantIDsByTrader(c, int8(v))
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "店铺查询失败")
			return
		}
		if len(ids) == 0 {
			response.OK(c, promotion.PageResult[promotion.StoreCouponListItem]{
				List: []promotion.StoreCouponListItem{}, Total: 0, Page: page, Limit: limit,
			})
			return
		}
		filter.MerIDs = ids
	}
	res, err := h.svc.ListStoreCouponsAdmin(c.Request.Context(), page, limit, filter)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	items, err := h.enrichStoreCoupons(c, res.List)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, promotion.PageResult[promotion.StoreCouponListItem]{
		List: items, Total: res.Total, Page: res.Page, Limit: res.Limit,
	})
}

func (h *Handler) StoreDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.GetStoreCouponDetail(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	info, _ := h.merchantInfoMap(c, []uint{row.MerID})
	if m, ok := info[row.MerID]; ok {
		row.MerName = m.Name
		row.IsTrader = m.IsTrader
		row.TraderName = traderLabel(m.IsTrader)
	}
	response.OK(c, row)
}

type merchantBrief struct {
	Name     string
	IsTrader int8
}

func (h *Handler) merchantIDsByTrader(c *gin.Context, isTrader int8) ([]uint, error) {
	var ids []uint
	err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_view").
		Where("is_trader = ?", isTrader).
		Pluck("merchant_id", &ids).Error
	return ids, err
}

func (h *Handler) merchantInfoMap(c *gin.Context, merIDs []uint) (map[uint]merchantBrief, error) {
	out := make(map[uint]merchantBrief, len(merIDs))
	if len(merIDs) == 0 || h.adminDB == nil {
		return out, nil
	}
	type row struct {
		MerchantID   uint   `gorm:"column:merchant_id"`
		MerchantName string `gorm:"column:merchant_name"`
		IsTrader     int8   `gorm:"column:is_trader"`
	}
	var rows []row
	err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name, is_trader").
		Where("merchant_id IN ?", merIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, item := range rows {
		out[item.MerchantID] = merchantBrief{Name: item.MerchantName, IsTrader: item.IsTrader}
	}
	return out, nil
}

func (h *Handler) enrichStoreCoupons(c *gin.Context, list []promotion.Coupon) ([]promotion.StoreCouponListItem, error) {
	out := make([]promotion.StoreCouponListItem, 0, len(list))
	if len(list) == 0 {
		return out, nil
	}
	ids := make([]uint, 0, len(list))
	merIDs := make([]uint, 0, len(list))
	seenMer := map[uint]struct{}{}
	for _, row := range list {
		ids = append(ids, row.CouponID)
		if _, ok := seenMer[row.MerID]; !ok {
			seenMer[row.MerID] = struct{}{}
			merIDs = append(merIDs, row.MerID)
		}
	}
	usage, err := h.svc.CountCouponUsageBatch(c.Request.Context(), ids)
	if err != nil {
		return nil, err
	}
	merchants, err := h.merchantInfoMap(c, merIDs)
	if err != nil {
		return nil, err
	}
	for _, row := range list {
		item := promotion.StoreCouponListItem{Coupon: row}
		item.CouponTypeName = storeCouponTypeLabel(row.Type)
		item.ClaimText = claimPeriodText(row)
		item.ValidityText = validityPeriodText(row)
		if m, ok := merchants[row.MerID]; ok {
			item.MerName = m.Name
			item.IsTrader = m.IsTrader
			item.TraderName = traderLabel(m.IsTrader)
		} else {
			item.MerName = "—"
			item.TraderName = "非自营"
		}
		if u, ok := usage[row.CouponID]; ok {
			item.ReceivedTotal = u[0]
			item.UsedTotal = u[1]
		}
		out = append(out, item)
	}
	return out, nil
}

func traderLabel(isTrader int8) string {
	if isTrader == 1 {
		return "自营"
	}
	return "非自营"
}

func storeCouponTypeLabel(typ int) string {
	switch typ {
	case promotion.CouponTypeProduct:
		return "商品券"
	case promotion.CouponTypePlatform:
		return "平台通用券"
	default:
		return "店铺券"
	}
}

func claimPeriodText(c promotion.Coupon) string {
	if c.IsTimeout == 1 && c.StartTime != nil && c.EndTime != nil {
		return c.StartTime.Format("2006-01-02 15:04:05") + " ~ " + c.EndTime.Format("2006-01-02 15:04:05")
	}
	return "不限时"
}

func validityPeriodText(c promotion.Coupon) string {
	if c.CouponType == int8(promotion.TemplateFixed) {
		return "时间段"
	}
	return strconv.FormatUint(uint64(c.CouponTime), 10) + "天"
}

func (h *Handler) Create(c *gin.Context) {
	var in promotion.CouponSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateAdmin(c.Request.Context(), 0, promotion.CouponTypePlatform, in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in promotion.CouponSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateAdmin(c.Request.Context(), 0, promotion.CouponTypePlatform, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Status(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in promotion.StatusInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SetStatus(c.Request.Context(), 0, promotion.CouponTypePlatform, uint(id), in.Status); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteAdmin(c.Request.Context(), 0, promotion.CouponTypePlatform, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Clone(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CloneAdmin(c.Request.Context(), 0, promotion.CouponTypePlatform, uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListSends(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	filter := promotion.CouponSendListFilter{
		DateFrom:   strings.TrimSpace(c.Query("date_from")),
		DateTo:     strings.TrimSpace(c.Query("date_to")),
		CouponName: strings.TrimSpace(c.Query("coupon_name")),
	}
	if raw := strings.TrimSpace(c.Query("coupon_type")); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "优惠券类型参数错误")
			return
		}
		filter.CouponType = &v
	}
	if raw := strings.TrimSpace(c.Query("status")); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		st := int8(v)
		filter.SendStatus = &st
	}
	res, err := h.svc.ListPlatformCouponSends(c.Request.Context(), page, limit, filter)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "发送记录查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) SendDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.GetPlatformCouponSend(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListSendUsers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	res, err := h.svc.ListPlatformCouponSendUsers(c.Request.Context(), uint(id), page, limit)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, promotion.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, promotion.ErrForbidden):
		response.Fail(c, http.StatusForbidden, err.Error())
	case errors.Is(err, promotion.ErrBadParam),
		errors.Is(err, promotion.ErrBadStatus):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "服务异常")
	}
}
