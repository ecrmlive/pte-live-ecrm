package combination

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/combination"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	svc     *combination.Service
	adminDB *gorm.DB
}

func NewHandler(svc *combination.Service, adminDB *gorm.DB) *Handler {
	return &Handler{svc: svc, adminDB: adminDB}
}

func (h *Handler) Register(r gin.IRoutes) {
	r.GET("/combination/groups", h.List)
	r.GET("/combination/groups/:id", h.Get)
	r.GET("/combination/buyings", h.ListBuyings)
	r.GET("/combination/buyings/:id", h.GetBuying)
	write := middleware.RequireAdminRoles("platform", "operations")
	manage := middleware.RequireAdminMenu(h.adminDB, "marketing.combination.manage")
	access := []gin.HandlerFunc{write, manage}
	r.PUT("/combination/groups/:id", write, manage, h.Update)
	r.DELETE("/combination/groups/:id", write, manage, h.Delete)
	r.POST("/combination/groups/audit", append(access, h.Audit)...)
	r.POST("/combination/groups/force-off", append(access, h.ForceOff)...)
}

func (h *Handler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	var merID *uint
	if s := c.Query("mer_id"); s != "" {
		v, _ := strconv.ParseUint(s, 10, 64)
		u := uint(v)
		merID = &u
	}
	res, err := h.svc.ListAdmin(c.Request.Context(), merID, page, limit)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	g, err := h.svc.Get(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var in combination.SaveInput
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.Update(c.Request.Context(), 0, uint(id), in)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := h.svc.Delete(c.Request.Context(), 0, uint(id)); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) Audit(c *gin.Context) {
	var body struct {
		ID      uint   `json:"id"`
		Status  int    `json:"status"`
		Refusal string `json:"refusal"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	g, err := h.svc.Audit(c.Request.Context(), body.ID, body.Status, body.Refusal)
	if err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, g)
}

func (h *Handler) ForceOff(c *gin.Context) {
	var body struct {
		IDs    []uint `json:"ids"`
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || len(body.IDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if err := h.svc.ForceOff(c.Request.Context(), body.IDs, body.Reason); err != nil {
		writeErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) ListBuyings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	merID, err := optionalMerID(c.Query("mer_id"))
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "店铺参数错误")
		return
	}
	q := combination.AdminBuyingQuery{
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
			response.OK(c, &combination.PageResult[combination.Buying]{
				List: []combination.Buying{}, Total: 0, Page: page, Limit: limit,
			})
			return
		}
		q.MerIDs = ids
	}
	res, err := h.svc.ListAdminBuyings(c.Request.Context(), q)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询失败")
		return
	}
	h.enrichBuyingMerchants(c, res.List)
	response.OK(c, res)
}

func (h *Handler) GetBuying(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "开团 ID 参数错误")
		return
	}
	row, err := h.svc.GetBuying(c.Request.Context(), uint(id))
	if err != nil {
		writeErr(c, err)
		return
	}
	list := []combination.Buying{*row}
	h.enrichBuyingMerchants(c, list)
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

type merchantLite struct {
	MerchantID   uint   `gorm:"column:merchant_id"`
	MerchantName string `gorm:"column:merchant_name"`
	IsTrader     int8   `gorm:"column:is_trader"`
}

func (h *Handler) enrichBuyingMerchants(c *gin.Context, rows []combination.Buying) {
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
	byID := map[uint]merchantLite{}
	if h.adminDB != nil && len(ids) > 0 {
		list := make([]merchantLite, 0, len(ids))
		_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
			Select("merchant_id, merchant_name, is_trader").
			Where("merchant_id IN ?", ids).Scan(&list)
		for _, m := range list {
			byID[m.MerchantID] = m
		}
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
	case errors.Is(err, combination.ErrNotFound):
		response.Fail(c, http.StatusNotFound, err.Error())
	case errors.Is(err, combination.ErrBadParam):
		response.Fail(c, http.StatusBadRequest, err.Error())
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
