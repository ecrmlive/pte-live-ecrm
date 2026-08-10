package nativecatalog

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// productOps stores platform-side display / recommend fields for a merchant product.
type productOps struct {
	ProductID    uint64    `gorm:"column:product_id;primaryKey"`
	IsUsed       int8      `gorm:"column:is_used"`
	Star         int8      `gorm:"column:star"`
	RankSort     int       `gorm:"column:rank_sort"`
	IsHot        int8      `gorm:"column:is_hot"`
	IsBest       int8      `gorm:"column:is_best"`
	IsBenefit    int8      `gorm:"column:is_benefit"`
	IsNew        int8      `gorm:"column:is_new"`
	CateHot      int8      `gorm:"column:cate_hot"`
	SysLabels    string    `gorm:"column:sys_labels"`
	ContentHTML  string    `gorm:"column:content_html"`
	RefundSwitch int8      `gorm:"column:refund_switch"`
	OnceMinCount int       `gorm:"column:once_min_count"`
	Ficti        int       `gorm:"column:ficti"`
	UpdatedBy    uint64    `gorm:"column:updated_by"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (productOps) TableName() string { return "qixi_crm_a_product_ops" }

type statusFilterItem struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func (h *Handler) registerManage(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	audit := middleware.RequireAdminMenu(h.adminDB, "product.audit.submit")
	r.GET("/products/status-filter", h.statusFilter)
	r.POST("/products/:id/show", platform, audit, h.setShow)
	r.POST("/products/:id/force-off", platform, audit, h.forceOff)
	r.POST("/products/batch/force-off", platform, audit, h.batchForceOff)
	r.POST("/products/batch/show", platform, audit, h.batchShow)
	r.POST("/products/batch/labels", platform, audit, h.stubBatchLabels)
	r.POST("/products/batch/recommend", platform, audit, h.stubBatchRecommend)
	r.POST("/products/batch/copy", platform, audit, h.stubBatchCopy)
	r.PUT("/products/:id/ops", platform, audit, h.updateOps)
	r.POST("/products/:id/ficti", platform, audit, h.setFicti)
}

func (h *Handler) statusFilter(c *gin.Context) {
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.OK(c, gin.H{"list": emptyStatusFilters(h.listFilterParams(c).IsGiftBag)})
		return
	}
	baseFilters := h.listFilterParams(c)
	items := statusFilterLabels(baseFilters.IsGiftBag)
	for i := range items {
		q := h.filteredBase(c, merchantIDs, baseFilters)
		q = applyTabType(q, items[i].Type)
		if err := q.Count(&items[i].Count).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "统计商品状态失败")
			return
		}
	}
	response.OK(c, gin.H{"list": items})
}

func statusFilterLabels(isGiftBag *int8) []statusFilterItem {
	if isGiftBag != nil && *isGiftBag == 1 {
		return []statusFilterItem{
			{Type: 1, Name: "出售中礼包"},
			{Type: 2, Name: "仓库中礼包"},
			{Type: 6, Name: "待审核礼包"},
			{Type: 7, Name: "审核未通过礼包"},
			{Type: 5, Name: "回收站礼包"},
		}
	}
	return []statusFilterItem{
		{Type: 1, Name: "出售中商品"},
		{Type: 2, Name: "仓库中商品"},
		{Type: 6, Name: "待审核商品"},
		{Type: 7, Name: "审核未通过商品"},
		{Type: 5, Name: "回收站商品"},
	}
}

func emptyStatusFilters(isGiftBag *int8) []statusFilterItem {
	items := statusFilterLabels(isGiftBag)
	for i := range items {
		items[i].Count = 0
	}
	return items
}

type listFilters struct {
	Keyword       string
	StoreName     string
	BrandName     string
	CateID        uint64
	MerID         uint64
	MerTypeID     uint64
	MerCategoryID uint64
	IsTrader      *int8 // 1 自营店商品（秒杀活动添加用）
	SVIPType      *int8
	ProductType   *int8
	Star          *int8
	IsUsed        *int8
	IsHot         *int8
	CateHot       *int8
	UsStatus      *int8
	IsGiftBag     *int8
	TabType       int
}

func (h *Handler) listFilterParams(c *gin.Context) listFilters {
	f := listFilters{
		Keyword:   strings.TrimSpace(c.Query("keyword")),
		StoreName: strings.TrimSpace(c.Query("store_name")),
		BrandName: strings.TrimSpace(c.Query("brand_name")),
	}
	f.CateID, _ = strconv.ParseUint(c.Query("cate_id"), 10, 64)
	f.MerID, _ = strconv.ParseUint(c.Query("mer_id"), 10, 64)
	f.MerTypeID, _ = strconv.ParseUint(c.Query("mer_type_id"), 10, 64)
	f.MerCategoryID, _ = strconv.ParseUint(c.Query("mer_category_id"), 10, 64)
	f.TabType, _ = strconv.Atoi(c.DefaultQuery("type", "1"))
	if v := strings.TrimSpace(c.Query("is_trader")); v != "" {
		n, err := strconv.ParseInt(v, 10, 8)
		if err == nil && (n == 0 || n == 1) {
			x := int8(n)
			f.IsTrader = &x
		}
	}
	if v := strings.TrimSpace(c.Query("svip_price_type")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.SVIPType = &x
	}
	if v := strings.TrimSpace(c.Query("product_type")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.ProductType = &x
	}
	if v := strings.TrimSpace(c.Query("star")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.Star = &x
	}
	if v := strings.TrimSpace(c.Query("is_used")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.IsUsed = &x
	}
	if v := strings.TrimSpace(c.Query("is_hot")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.IsHot = &x
	}
	if v := strings.TrimSpace(c.Query("cate_hot")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.CateHot = &x
	}
	if v := strings.TrimSpace(c.Query("us_status")); v != "" {
		n, _ := strconv.ParseInt(v, 10, 8)
		x := int8(n)
		f.UsStatus = &x
	}
	if v := strings.TrimSpace(c.Query("is_gift_bag")); v != "" {
		switch v {
		case "0", "1":
			n, _ := strconv.ParseInt(v, 10, 8)
			x := int8(n)
			f.IsGiftBag = &x
		}
	} else {
		// 普通商品列表默认排除分销礼包
		zero := int8(0)
		f.IsGiftBag = &zero
	}
	// legacy status query still accepted when type omitted from old clients
	if strings.TrimSpace(c.Query("type")) == "" {
		if s := strings.TrimSpace(c.Query("status")); s != "" {
			f.UsStatus = nil
			switch statusName(s) {
			case "on_sale":
				f.TabType = 1
			case "off_sale":
				f.TabType = 2
			case "rejected":
				f.TabType = 7
			default:
				f.TabType = 6
			}
		}
	}
	return f
}

func (h *Handler) filteredBase(c *gin.Context, merchantIDs []uint64, f listFilters) *gorm.DB {
	q := h.base(c, merchantIDs)
	q = q.Joins("LEFT JOIN qixi_crm_m_product_recycle_bin AS rb ON rb.product_id = p.id")
	q = q.Joins("LEFT JOIN qixi_crm_m_product_detail AS d ON d.product_id = p.id")
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("(p.title LIKE ? OR d.keyword LIKE ? OR CAST(p.id AS CHAR) LIKE ?)", like, like, like)
	}
	if f.StoreName != "" {
		q = q.Where("s.name LIKE ?", "%"+f.StoreName+"%")
	}
	if f.BrandName != "" {
		q = q.Where("p.brand_name LIKE ?", "%"+f.BrandName+"%")
	}
	if f.CateID > 0 {
		q = q.Where("p.category_id = ?", f.CateID)
	}
	if f.MerID > 0 {
		q = q.Where("s.merchant_id = ?", f.MerID)
	}
	if f.SVIPType != nil {
		q = q.Where("p.svip_price_type = ?", *f.SVIPType)
	}
	if f.MerTypeID > 0 || f.MerCategoryID > 0 {
		ids, err := h.merchantIDsByAdminMeta(c, f.MerTypeID, f.MerCategoryID)
		if err != nil {
			// keep query empty on meta failure
			q = q.Where("1 = 0")
			return q
		}
		if len(ids) == 0 {
			q = q.Where("1 = 0")
			return q
		}
		q = q.Where("s.merchant_id IN ?", ids)
	}
	if f.IsTrader != nil {
		ids, err := h.merchantIDsByTrader(c, *f.IsTrader)
		if err != nil || len(ids) == 0 {
			return q.Where("1 = 0")
		}
		q = q.Where("s.merchant_id IN ?", ids)
	}
	if f.UsStatus != nil {
		q = q.Where("p.status = ?", statusName(strconv.Itoa(int(*f.UsStatus))))
	}
	q = h.applyOpsFilters(c, q, f)
	if f.ProductType != nil {
		ids, err := h.productIDsByType(c, *f.ProductType)
		if err != nil || len(ids) == 0 {
			return q.Where("1 = 0")
		}
		q = q.Where("p.id IN ?", ids)
	}
	if f.IsGiftBag != nil {
		q = q.Where("p.is_gift_bag = ?", *f.IsGiftBag)
	}
	return q
}

func (h *Handler) applyOpsFilters(c *gin.Context, q *gorm.DB, f listFilters) *gorm.DB {
	if f.Star == nil && f.IsUsed == nil && f.IsHot == nil && f.CateHot == nil {
		return q
	}
	// 「显示中」= 无 ops 行或 is_used=1 → 排除 is_used=0
	if f.IsUsed != nil && *f.IsUsed == 1 && f.Star == nil && f.IsHot == nil && f.CateHot == nil {
		hidden, err := h.opsProductIDsBy(c, map[string]any{"is_used": 0})
		if err != nil {
			return q.Where("1 = 0")
		}
		if len(hidden) > 0 {
			return q.Where("p.id NOT IN ?", hidden)
		}
		return q
	}
	conds := map[string]any{}
	if f.Star != nil {
		conds["star"] = *f.Star
	}
	if f.IsUsed != nil {
		conds["is_used"] = *f.IsUsed
	}
	if f.IsHot != nil {
		conds["is_hot"] = *f.IsHot
	}
	if f.CateHot != nil {
		conds["cate_hot"] = *f.CateHot
	}
	ids, err := h.opsProductIDsBy(c, conds)
	if err != nil || len(ids) == 0 {
		return q.Where("1 = 0")
	}
	return q.Where("p.id IN ?", ids)
}

func (h *Handler) opsProductIDsBy(c *gin.Context, conds map[string]any) ([]uint64, error) {
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_product_ops").Select("product_id")
	for k, v := range conds {
		q = q.Where(k+" = ?", v)
	}
	var rows []struct {
		ProductID uint64 `gorm:"column:product_id"`
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ProductID)
	}
	return ids, nil
}

func (h *Handler) productIDsByType(c *gin.Context, productType int8) ([]uint64, error) {
	var rows []struct {
		ProductID uint64 `gorm:"column:product_id"`
	}
	if err := h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Select("product_id").Where("product_type = ?", productType).Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ProductID)
	}
	return ids, nil
}

func applyTabType(q *gorm.DB, tabType int) *gorm.DB {
	switch tabType {
	case 1:
		return q.Where("rb.product_id IS NULL AND p.status = ?", "on_sale")
	case 2:
		return q.Where("rb.product_id IS NULL AND p.status = ?", "off_sale")
	case 6:
		return q.Where("rb.product_id IS NULL AND p.status IN ?", []string{"pending_review", "draft"})
	case 7:
		return q.Where("rb.product_id IS NULL AND p.status = ?", "rejected")
	case 5:
		return q.Where("rb.product_id IS NOT NULL")
	default:
		return q.Where("rb.product_id IS NULL")
	}
}

func (h *Handler) merchantIDsByAdminMeta(c *gin.Context, typeID, categoryID uint64) ([]uint64, error) {
	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").Select("merchant_id")
	if typeID > 0 {
		q = q.Where("type_id = ?", typeID)
	}
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	var rows []struct {
		MerchantID uint64 `gorm:"column:merchant_id"`
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row.MerchantID > 0 {
			ids = append(ids, row.MerchantID)
		}
	}
	return ids, nil
}

func (h *Handler) merchantIDsByTrader(c *gin.Context, isTrader int8) ([]uint64, error) {
	var rows []struct {
		MerchantID uint64 `gorm:"column:merchant_id"`
	}
	err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("merchant_id").Where("is_trader = ?", isTrader).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	for _, row := range rows {
		if row.MerchantID > 0 {
			ids = append(ids, row.MerchantID)
		}
	}
	return ids, nil
}

func (h *Handler) setShow(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Status *int8 `json:"status"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || req.Status == nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	used := int8(0)
	if *req.Status == 1 {
		used = 1
	}
	if err := h.upsertOpsShow(c, []uint64{id}, used); err != nil {
		writeManageErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) batchShow(c *gin.Context) {
	var req struct {
		IDs    []uint64 `json:"ids"`
		Status *int8    `json:"status"`
	}
	if c.ShouldBindJSON(&req) != nil || req.Status == nil || len(req.IDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "请选择商品")
		return
	}
	used := int8(0)
	if *req.Status == 1 {
		used = 1
	}
	if err := h.upsertOpsShow(c, req.IDs, used); err != nil {
		writeManageErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) upsertOpsShow(c *gin.Context, ids []uint64, used int8) error {
	adminID := uint64(middleware.AdminID(c))
	now := time.Now()
	for _, id := range ids {
		if id == 0 {
			continue
		}
		ops := productOps{ProductID: id, IsUsed: used, UpdatedBy: adminID, UpdatedAt: now}
		if err := h.adminDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_used", "updated_by", "updated_at"}),
		}).Create(&ops).Error; err != nil {
			return err
		}
		_ = h.businessDB.WithContext(c.Request.Context()).
			Table("qixi_crm_b_product_view").
			Where("product_id = ?", id).
			Update("sale_status", used).Error
	}
	return nil
}

func (h *Handler) forceOff(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Reason string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&req)
	reason := strings.TrimSpace(req.Reason)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	if reason == "" {
		response.Fail(c, http.StatusBadRequest, "请输入强制下架原因")
		return
	}
	if err := h.forceOffProducts(c, []uint64{id}, reason); err != nil {
		writeManageErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) batchForceOff(c *gin.Context) {
	var req struct {
		IDs    []uint64 `json:"ids"`
		Reason string   `json:"reason"`
	}
	if c.ShouldBindJSON(&req) != nil || len(req.IDs) == 0 {
		response.Fail(c, http.StatusBadRequest, "请选择商品")
		return
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		response.Fail(c, http.StatusBadRequest, "请输入强制下架原因")
		return
	}
	if err := h.forceOffProducts(c, req.IDs, reason); err != nil {
		writeManageErr(c, err)
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) forceOffProducts(c *gin.Context, ids []uint64, reason string) error {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return errForceOffReason
	}
	for _, id := range ids {
		if id == 0 {
			continue
		}
		var row productRow
		var command productAuditOutbox
		changed := false
		err := h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Table("qixi_crm_m_product AS p").
				Select("p.id,p.store_id,s.merchant_id,m.name AS merchant_name,s.name AS store_name,p.title,p.category_id,p.brand_name,p.svip_price_type,p.svip_price,p.status,p.version,p.created_at").
				Joins("JOIN qixi_crm_m_store AS s ON s.id = p.store_id").
				Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id").
				Where("p.id = ?", id).Scan(&row).Error; err != nil {
				return err
			}
			if row.ID == 0 {
				return gorm.ErrRecordNotFound
			}
			if row.Status == "off_sale" {
				return nil
			}
			if row.Status != "on_sale" && row.Status != "pending_review" && row.Status != "draft" {
				return errForceOffState
			}
			if err := tx.Table("qixi_crm_m_product").Where("id = ?", id).
				Updates(map[string]any{"status": "off_sale", "version": gorm.Expr("version + 1")}).Error; err != nil {
				return err
			}
			changed = true
			command = productAuditOutbox{
				ProductID:    row.ID,
				StoreID:      row.StoreID,
				Action:       "delete",
				ReviewStatus: "approved",
				Reason:       reason,
				ReviewedBy:   uint64(middleware.AdminID(c)),
				Status:       "pending",
			}
			return tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "product_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"action", "review_status", "reason", "reviewed_by", "status", "attempts", "last_error", "updated_at"}),
			}).Create(&command).Error
		})
		if err != nil {
			return err
		}
		if command.ProductID > 0 {
			_, _ = h.processAuditOutbox(c.Request.Context(), command)
		}
		if changed {
			_ = h.adminDB.WithContext(c.Request.Context()).
				Table("qixi_crm_a_product_review").
				Where("product_id = ?", id).
				Updates(map[string]any{"status": "off_sale", "reason": reason, "reviewed_by": middleware.AdminID(c), "reviewed_at": time.Now()}).Error
			h.writeProductOperateLog(c, id, "product.force_off", reason)
		}
	}
	return nil
}

func (h *Handler) updateOps(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Star     *int8 `json:"star"`
		RankSort *int  `json:"rank"`
		IsHot    *int8 `json:"is_hot"`
		IsBest   *int8 `json:"is_best"`
		IsBenefit *int8 `json:"is_benefit"`
		IsNew    *int8 `json:"is_new"`
		CateHot  *int8 `json:"cate_hot"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	ops := productOps{ProductID: id, IsUsed: 1, UpdatedBy: uint64(middleware.AdminID(c)), UpdatedAt: time.Now()}
	var existing productOps
	_ = h.adminDB.WithContext(c.Request.Context()).Where("product_id = ?", id).Take(&existing).Error
	if existing.ProductID > 0 {
		ops = existing
		ops.UpdatedBy = uint64(middleware.AdminID(c))
		ops.UpdatedAt = time.Now()
	}
	if req.Star != nil {
		if *req.Star < 0 || *req.Star > 5 {
			response.Fail(c, http.StatusBadRequest, "推荐级别须为 0-5")
			return
		}
		ops.Star = *req.Star
	}
	if req.RankSort != nil {
		ops.RankSort = *req.RankSort
	}
	if req.IsHot != nil {
		ops.IsHot = *req.IsHot
	}
	if req.IsBest != nil {
		ops.IsBest = *req.IsBest
	}
	if req.IsBenefit != nil {
		ops.IsBenefit = *req.IsBenefit
	}
	if req.IsNew != nil {
		ops.IsNew = *req.IsNew
	}
	if req.CateHot != nil {
		ops.CateHot = *req.CateHot
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"star", "rank_sort", "is_hot", "is_best", "is_benefit", "is_new", "cate_hot", "updated_by", "updated_at"}),
	}).Create(&ops).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存运营字段失败")
		return
	}
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) stubBatchLabels(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "TODO: 批量设置标签尚未接入平台标签绑定（qixi_crm_a_product_ops.sys_labels）")
}

func (h *Handler) stubBatchRecommend(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "TODO: 批量设置推荐尚未接入；请使用单行编辑推荐级别或 PUT /products/:id/ops")
}

func (h *Handler) stubBatchCopy(c *gin.Context) {
	response.Fail(c, http.StatusNotImplemented, "TODO: 批量复制商品到店铺尚未实现（需商户域复制流水）")
}

// setFicti adjusts virtual sold count (ficti) and mirrors the delta into product_view.sales
// (aligned with CRMEB StoreProduct::addFicti). type=1 increase, type=2 decrease.
func (h *Handler) setFicti(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req struct {
		Type  int `json:"type"`
		Ficti int `json:"ficti"`
	}
	if id == 0 || c.ShouldBindJSON(&req) != nil || (req.Type != 1 && req.Type != 2) {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if req.Ficti <= 0 {
		response.Fail(c, http.StatusBadRequest, "已售数量必须大于0")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	if merchantIDs != nil && len(merchantIDs) == 0 {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}
	var row productRow
	if err := h.base(c, merchantIDs).Where("p.id = ?", id).Scan(&row).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询商品失败")
		return
	}
	if row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}

	var ops productOps
	err = h.adminDB.WithContext(c.Request.Context()).Where("product_id = ?", id).Take(&ops).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.Fail(c, http.StatusInternalServerError, "查询已售数量失败")
		return
	}
	current := 0
	if ops.ProductID > 0 {
		current = ops.Ficti
	}
	delta := req.Ficti
	if req.Type == 2 {
		if current < req.Ficti {
			response.Fail(c, http.StatusBadRequest, "已售数量不足")
			return
		}
		delta = -req.Ficti
	}
	newFicti := current + delta
	if newFicti < 0 {
		response.Fail(c, http.StatusBadRequest, "已售数量不能为负数")
		return
	}

	now := time.Now()
	adminID := uint64(middleware.AdminID(c))
	next := productOps{
		ProductID:    id,
		IsUsed:       1,
		RefundSwitch: 1,
		OnceMinCount: 1,
		Ficti:        newFicti,
		UpdatedBy:    adminID,
		UpdatedAt:    now,
	}
	if ops.ProductID > 0 {
		next = ops
		next.Ficti = newFicti
		next.UpdatedBy = adminID
		next.UpdatedAt = now
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"ficti", "updated_by", "updated_at"}),
	}).Create(&next).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存已售数量失败")
		return
	}

	var viewSales int
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Select("sales").Where("product_id = ?", id).Scan(&viewSales).Error
	newSales := viewSales + delta
	if newSales < 0 {
		newSales = 0
	}
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Where("product_id = ?", id).Update("sales", newSales).Error

	response.OK(c, gin.H{"ok": true, "ficti": newFicti, "sales": newSales})
}

var (
	errForceOffState  = errors.New("product cannot force off")
	errForceOffReason = errors.New("force off reason required")
)

func writeManageErr(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		response.Fail(c, http.StatusNotFound, "商品不存在")
	case errors.Is(err, errForceOffReason):
		response.Fail(c, http.StatusBadRequest, "请输入强制下架原因")
	case errors.Is(err, errForceOffState):
		response.Fail(c, http.StatusConflict, "当前状态不可强制下架")
	default:
		response.Fail(c, http.StatusInternalServerError, "操作失败")
	}
}
