package assist

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *assist.Service
	adminDB *gorm.DB
}

type visibilityInput struct {
	IsShow *int `json:"is_show"`
}

func NewHandler(svc *assist.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.assist.manage"),
	}
	// 活动商品（配置）
	r.GET("/assist/actives", append(access, h.List)...)
	r.GET("/assist/actives/:id", append(access, h.Get)...)
	r.PUT("/assist/actives/:id", append(access, h.Update)...)
	// 助力活动（用户发起实例）
	r.GET("/assist/sets", append(access, h.ListSets)...)
	r.GET("/assist/sets/:id", append(access, h.GetSet)...)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID, err := optionalMerID(c.Query("mer_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "商户 ID 参数错误")
		return
	}
	res, err := h.svc.ListAdmin(c.Request.Context(), merID, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	h.enrichActiveMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "活动 ID 参数错误")
		return
	}
	row, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []assist.ProductAssist{*row}
	h.enrichActiveMerchants(c, tmp)
	*row = tmp[0]
	response.OK(c, row)
}

func (h *Handler) ListSets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID, err := optionalMerID(c.Query("mer_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "店铺参数错误")
		return
	}
	q := assist.AdminSetQuery{
		Page:     page,
		Limit:    limit,
		MerID:    merID,
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		UserName: strings.TrimSpace(c.Query("user_name")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
	}
	if s := strings.TrimSpace(c.Query("status")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "状态参数错误")
			return
		}
		q.Status = &n
	}
	if s := strings.TrimSpace(c.Query("is_trader")); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil || (n != 0 && n != 1) {
			response.Fail(c, http.StatusBadRequest, "店铺类别参数错误")
			return
		}
		ids, err := h.lookupMerchantIDsByTrader(c, int8(n))
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return
		}
		if len(ids) == 0 {
			response.OK(c, &assist.PageResult[assist.AssistSet]{
				List: []assist.AssistSet{}, Total: 0, Page: page, Limit: limit,
			})
			return
		}
		q.MerIDs = ids
	}
	res, err := h.svc.ListAdminSets(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	h.enrichSetMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) GetSet(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "助力实例 ID 参数错误")
		return
	}
	row, err := h.svc.GetSet(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	list := []assist.AssistSet{*row}
	h.enrichSetMerchants(c, list)
	*row = list[0]
	response.OK(c, row)
}

func optionalMerID(raw string) (*uint, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return nil, errors.New("invalid merchant id")
	}
	value := uint(id)
	return &value, nil
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in visibilityInput
	if err := c.ShouldBindJSON(&in); err != nil || !validVisibility(&in) {
		response.Fail(c, http.StatusBadRequest, "仅允许更新助力活动展示状态")
		return
	}
	row, err := h.svc.SetShow(c.Request.Context(), 0, uint(id), *in.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func validVisibility(in *visibilityInput) bool {
	return in != nil && in.IsShow != nil && (*in.IsShow == 0 || *in.IsShow == 1)
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
		return nil, nil
	}
	var ids []uint
	err := h.adminDB.WithContext(c.Request.Context()).
		Table("qixi_crm_a_merchant_view").
		Where("is_trader = ?", isTrader).
		Pluck("merchant_id", &ids).Error
	return ids, err
}

func (h *Handler) enrichSetMerchants(c *gin.Context, rows []assist.AssistSet) {
	byID := h.loadMerchantMap(c, collectMerIDs(rows, func(r assist.AssistSet) uint { return r.MerID }))
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
		rows[i].TraderName = traderLabel(m.IsTrader)
	}
}

func (h *Handler) enrichActiveMerchants(c *gin.Context, rows []assist.ProductAssist) {
	byID := h.loadMerchantMap(c, collectMerIDs(rows, func(r assist.ProductAssist) uint { return r.MerID }))
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
		rows[i].TraderName = traderLabel(m.IsTrader)
	}
}

type merchantLite struct {
	MerchantID   uint   `gorm:"column:merchant_id"`
	MerchantName string `gorm:"column:merchant_name"`
	IsTrader     int8   `gorm:"column:is_trader"`
}

func collectMerIDs[T any](rows []T, merID func(T) uint) []uint {
	ids := make([]uint, 0, len(rows))
	seen := map[uint]struct{}{}
	for _, row := range rows {
		id := merID(row)
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids
}

func (h *Handler) loadMerchantMap(c *gin.Context, ids []uint) map[uint]merchantLite {
	byID := map[uint]merchantLite{}
	if h.adminDB == nil || len(ids) == 0 {
		return byID
	}
	list := make([]merchantLite, 0, len(ids))
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name, is_trader").
		Where("merchant_id IN ?", ids).Scan(&list)
	for _, m := range list {
		byID[m.MerchantID] = m
	}
	return byID
}

func traderLabel(isTrader int8) string {
	if isTrader == 1 {
		return "自营"
	}
	return "非自营"
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, assist.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, assist.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
