package presell

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/presell"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *presell.Service
	adminDB *gorm.DB
}

func NewHandler(svc *presell.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	access := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.presell.manage"),
	}
	// filter 须在 :id 之前
	r.GET("/presell/actives/filter", append(access, h.TypeFilter)...)
	r.GET("/presell/actives", append(access, h.List)...)
	r.GET("/presell/actives/:id", append(access, h.Get)...)
	r.PUT("/presell/actives/:id", append(access, h.Update)...)
	r.PUT("/presell/actives/:id/show", append(access, h.SetShow)...)
	r.PUT("/presell/actives/:id/star", append(access, h.SetStar)...)
	r.PUT("/presell/actives/:id/labels", append(access, h.SetLabels)...)
	r.POST("/presell/actives/audit", append(access, h.Audit)...)
	r.POST("/presell/actives/force-off", append(access, h.ForceOff)...)
	r.DELETE("/presell/actives/:id", append(access, h.Delete)...)
}

func (h *Handler) parseAdminQuery(c *gin.Context) (presell.AdminQuery, bool) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	q := presell.AdminQuery{
		Page:      page,
		Limit:     limit,
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		SysLabels: strings.TrimSpace(c.Query("sys_labels")),
	}
	if s := c.Query("presell_type"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || (v != 1 && v != 2) {
			response.Fail(c, http.StatusBadRequest, "预售类型参数错误")
			return q, false
		}
		q.PresellType = v
	}
	if s := c.Query("mer_id"); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "店铺参数错误")
			return q, false
		}
		u := uint(v)
		q.MerID = &u
	}
	if s := c.Query("is_trader"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "店铺类别参数错误")
			return q, false
		}
		tr := int8(v)
		q.IsTrader = &tr
		ids, err := h.lookupMerchantIDsByTrader(c, tr)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return q, false
		}
		if len(ids) == 0 {
			q.MerIDs = []uint{0}
		} else {
			q.MerIDs = ids
		}
	}
	if s := c.Query("star"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || v < 0 || v > 5 {
			response.Fail(c, http.StatusBadRequest, "推荐级别参数错误")
			return q, false
		}
		st := int8(v)
		q.Star = &st
	}
	if s := c.Query("product_status"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "审核状态参数错误")
			return q, false
		}
		q.ProductStatus = &v
	}
	if s := c.Query("type"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || (v != 0 && v != 1 && v != 2) {
			response.Fail(c, http.StatusBadRequest, "活动状态参数错误")
			return q, false
		}
		at := int8(v)
		q.ActivityType = &at
	}
	if s := c.Query("us_status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "商品状态参数错误")
			return q, false
		}
		us := int8(v)
		q.UsStatus = &us
	}
	return q, true
}

func (h *Handler) List(c *gin.Context) {
	q, ok := h.parseAdminQuery(c)
	if !ok {
		return
	}
	res, err := h.svc.ListAdminFiltered(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	h.enrichMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) TypeFilter(c *gin.Context) {
	q, ok := h.parseAdminQuery(c)
	if !ok {
		return
	}
	list, err := h.svc.TypeFilter(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "统计失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	tmp := []presell.ProductPresell{*row}
	h.enrichMerchants(c, tmp)
	*row = tmp[0]
	response.OK(c, row)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in presell.SaveInput
	if id == 0 || c.ShouldBindJSON(&in) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if in.Status != nil && *in.Status != 0 && *in.Status != 1 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Update(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsShow *int `json:"is_show"`
	}
	if id == 0 || c.ShouldBindJSON(&body) != nil || body.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetShow(c.Request.Context(), 0, uint(id), *body.IsShow)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetStar(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Star *int8 `json:"star"`
	}
	if id == 0 || c.ShouldBindJSON(&body) != nil || body.Star == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetStar(c.Request.Context(), uint(id), *body.Star)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetLabels(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		SysLabels string `json:"sys_labels"`
	}
	if id == 0 || c.ShouldBindJSON(&body) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetLabels(c.Request.Context(), uint(id), body.SysLabels)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) Audit(c *gin.Context) {
	var body struct {
		ID      uint   `json:"id"`
		Status  int    `json:"status"`
		Refusal string `json:"refusal"`
	}
	if c.ShouldBindJSON(&body) != nil || body.ID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.Audit(c.Request.Context(), body.ID, body.Status, body.Refusal)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ForceOff(c *gin.Context) {
	var body struct {
		IDs    []uint `json:"ids"`
		Reason string `json:"reason"`
	}
	if c.ShouldBindJSON(&body) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.ForceOff(c.Request.Context(), body.IDs, body.Reason); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
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

func (h *Handler) enrichMerchants(c *gin.Context, rows []presell.ProductPresell) {
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
	case errors.Is(err, presell.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, presell.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
