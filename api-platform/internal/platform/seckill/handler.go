package seckill

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/seckill"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *seckill.Service
	adminDB *gorm.DB
}

func NewHandler(svc *seckill.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	configWrite := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.seckill.config"),
	}

	// 秒杀配置（场次）；选择器仍走无 page 的 GET /seckill/times
	r.GET("/seckill/times", h.Times)
	r.POST("/seckill/times", append(configWrite, h.CreateTime)...)
	r.PUT("/seckill/times/:id", append(configWrite, h.UpdateTime)...)
	r.PUT("/seckill/times/:id/status", append(configWrite, h.SetTimeStatus)...)
	r.DELETE("/seckill/times/:id", append(configWrite, h.DeleteTime)...)

	manageWrite := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.seckill.manage"),
	}
	activityWrite := []gin.HandlerFunc{
		middleware.RequireAdminRoles("platform", "operations"),
		middleware.RequireAdminMenu(h.adminDB, "marketing.seckill.activity"),
	}

	// 秒杀活动（活动场）；须在 :id 路由之前注册固定路径
	r.GET("/seckill/activities", h.ListActivities)
	r.POST("/seckill/activities", append(activityWrite, h.CreateActivity)...)
	r.GET("/seckill/activities/:id/stats/people", h.ListActivityStatPeople)
	r.GET("/seckill/activities/:id/stats/orders", h.ListActivityStatOrders)
	r.GET("/seckill/activities/:id/stats/products", h.ListActivityStatProducts)
	r.GET("/seckill/activities/:id/stats", h.GetActivityStats)
	r.GET("/seckill/activities/:id/products", h.ListActivityProducts)
	r.POST("/seckill/activities/:id/products", append(activityWrite, h.SaveActivityProducts)...)
	r.GET("/seckill/activities/:id", h.GetActivity)
	r.PUT("/seckill/activities/:id", append(activityWrite, h.UpdateActivity)...)
	r.PUT("/seckill/activities/:id/status", append(activityWrite, h.SetActivityStatus)...)
	r.POST("/seckill/activities/:id/clone", append(activityWrite, h.CloneActivity)...)
	r.DELETE("/seckill/activities/:id", append(activityWrite, h.DeleteActivity)...)

	// filter 须在 :id 之前
	r.GET("/seckill/actives/filter", h.StatusFilter)
	r.GET("/seckill/actives", h.List)
	r.GET("/seckill/actives/:id", h.Get)
	r.PUT("/seckill/actives/:id", append(manageWrite, h.Update)...)
	r.PUT("/seckill/actives/:id/show", append(manageWrite, h.SetShow)...)
	r.PUT("/seckill/actives/:id/star", append(manageWrite, h.SetStar)...)
	r.PUT("/seckill/actives/:id/labels", append(manageWrite, h.SetLabels)...)
	r.POST("/seckill/actives/force-off", append(manageWrite, h.ForceOff)...)
	r.DELETE("/seckill/actives/:id", append(manageWrite, h.Delete)...)
}

// Times：带 page 时返回后台分页列表；否则返回启用中的场次（活动选择器用）。
func (h *Handler) Times(c *gin.Context) {
	if c.Query("page") != "" || c.Query("limit") != "" || c.Query("status") != "" {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		q := seckill.TimeSlotQuery{Page: page, Limit: limit}
		if s := c.Query("status"); s != "" {
			v, err := strconv.ParseInt(s, 10, 8)
			if err != nil || (v != 0 && v != 1) {
				response.Fail(c, http.StatusBadRequest, "是否显示参数错误")
				return
			}
			st := int8(v)
			q.Status = &st
		}
		res, err := h.svc.ListTimesAdmin(c.Request.Context(), q)
		if err != nil {
			response.Fail(c, http.StatusInternalServerError, "查询失败")
			return
		}
		response.OK(c, res)
		return
	}
	rows, err := h.svc.ListTimes(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) CreateTime(c *gin.Context) {
	var in seckill.TimeSlotInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateTime(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateTime(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in seckill.TimeSlotInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateTime(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetTimeStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetTimeStatus(c.Request.Context(), uint(id), *body.Status)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteTime(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteTime(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) parseActiveQuery(c *gin.Context) (seckill.ActiveAdminQuery, bool) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	tabType, _ := strconv.Atoi(c.DefaultQuery("type", "1"))
	q := seckill.ActiveAdminQuery{
		Type:       tabType,
		Page:       page,
		Limit:      limit,
		ActiveName: strings.TrimSpace(c.Query("active_name")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		SysLabels:  strings.TrimSpace(c.Query("sys_labels")),
	}
	if s := c.Query("mer_id"); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil || v == 0 {
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
	if s := c.Query("us_status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "活动状态参数错误")
			return q, false
		}
		us := int8(v)
		q.UsStatus = &us
	}
	return q, true
}

func (h *Handler) List(c *gin.Context) {
	q, ok := h.parseActiveQuery(c)
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

func (h *Handler) StatusFilter(c *gin.Context) {
	q, ok := h.parseActiveQuery(c)
	if !ok {
		return
	}
	list, err := h.svc.StatusFilter(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "统计失败")
		return
	}
	response.OK(c, gin.H{"list": list})
}

func (h *Handler) SetShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		IsShow *int8 `json:"is_show"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.IsShow == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetShow(c.Request.Context(), uint(id), *body.IsShow)
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
	if err := c.ShouldBindJSON(&body); err != nil || body.Star == nil {
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
	if err := c.ShouldBindJSON(&body); err != nil {
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

func (h *Handler) ForceOff(c *gin.Context) {
	var body struct {
		IDs    []uint `json:"ids"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.ForceOff(c.Request.Context(), body.IDs, body.Reason); err != nil {
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

func (h *Handler) enrichMerchants(c *gin.Context, rows []seckill.Active) {
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

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	a, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in seckill.ActiveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	a, err := h.svc.Update(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, a)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListActivities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	q := seckill.ActivityQuery{
		Name:     strings.TrimSpace(c.Query("name")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
		Page:     page,
		Limit:    limit,
	}
	if s := c.Query("status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || (v != 0 && v != 1) {
			response.Fail(c, http.StatusBadRequest, "是否开启参数错误")
			return
		}
		st := int8(v)
		q.Status = &st
	}
	if s := c.Query("active_status"); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || (v != 0 && v != 1 && v != -1) {
			response.Fail(c, http.StatusBadRequest, "活动状态参数错误")
			return
		}
		st := int8(v)
		q.ActiveStatus = &st
	}
	res, err := h.svc.ListActivities(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.GetActivity(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) ListActivityProducts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	q := seckill.ActivityProductQuery{
		ActivityID: uint(id),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
		Page:       page,
		Limit:      limit,
	}
	if s := strings.TrimSpace(c.Query("product_status")); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil || (v != 0 && v != 1 && v != -1) {
			response.Fail(c, http.StatusBadRequest, "商品审核状态参数错误")
			return
		}
		st := int8(v)
		q.ProductStatus = &st
	}
	res, err := h.svc.ListActivityProductsAdmin(c.Request.Context(), q)
	if err != nil {
		writeErr(c, err)
		return
	}
	h.enrichActivityProductMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) SaveActivityProducts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in seckill.ActivityProductsSaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.SaveActivityProducts(c.Request.Context(), uint(id), in); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) enrichActivityProductMerchants(c *gin.Context, rows []seckill.ActivityProductItem) {
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
	}
	list := make([]merRow, 0, len(ids))
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id, merchant_name").
		Where("merchant_id IN ?", ids).Scan(&list)
	byID := map[uint]string{}
	for _, m := range list {
		byID[m.MerchantID] = m.MerchantName
	}
	for i := range rows {
		if rows[i].MerName == "" {
			rows[i].MerName = byID[rows[i].MerID]
		}
	}
}

func (h *Handler) GetActivityStats(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var merID *uint
	if s := strings.TrimSpace(c.Query("mer_id")); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, "店铺参数错误")
			return
		}
		u := uint(v)
		merID = &u
	}
	row, err := h.svc.GetActivityStats(c.Request.Context(), uint(id), merID)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func parseActivityStatQuery(c *gin.Context) (seckill.ActivityStatQuery, error) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	q := seckill.ActivityStatQuery{
		Keyword:  strings.TrimSpace(c.Query("keyword")),
		DateFrom: strings.TrimSpace(c.Query("date_from")),
		DateTo:   strings.TrimSpace(c.Query("date_to")),
		Page:     page,
		Limit:    limit,
	}
	if s := strings.TrimSpace(c.Query("mer_id")); s != "" {
		v, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return q, err
		}
		u := uint(v)
		q.MerID = &u
	}
	if s := strings.TrimSpace(c.Query("status")); s != "" {
		v, err := strconv.ParseInt(s, 10, 8)
		if err != nil {
			return q, err
		}
		st := int8(v)
		q.Status = &st
	}
	return q, nil
}

func (h *Handler) ListActivityStatPeople(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	q, err := parseActivityStatQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.svc.ListActivityStatPeople(c.Request.Context(), uint(id), q)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListActivityStatOrders(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	q, err := parseActivityStatQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.svc.ListActivityStatOrders(c.Request.Context(), uint(id), q)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) ListActivityStatProducts(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	q, err := parseActivityStatQuery(c)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	res, err := h.svc.ListActivityStatProducts(c.Request.Context(), uint(id), q)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) CreateActivity(c *gin.Context) {
	var in seckill.ActivityInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.CreateActivity(c.Request.Context(), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) UpdateActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in seckill.ActivityInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.UpdateActivity(c.Request.Context(), uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) SetActivityStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var body struct {
		Status *int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	row, err := h.svc.SetActivityStatus(c.Request.Context(), uint(id), *body.Status)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) CloneActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	row, err := h.svc.CloneActivity(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, row)
}

func (h *Handler) DeleteActivity(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.DeleteActivity(c.Request.Context(), uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func writeErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, seckill.ErrNotFound),
		errors.Is(err, seckill.ErrTimeNotFound),
		errors.Is(err, seckill.ErrActivityNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, seckill.ErrBadParam), errors.Is(err, seckill.ErrTimeOverlap):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
