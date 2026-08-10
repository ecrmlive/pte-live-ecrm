package nativecatalog

import (
	"encoding/json"
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

func (h *Handler) registerEdit(r gin.IRoutes) {
	platform := middleware.RequireAdminRoles("platform")
	audit := middleware.RequireAdminMenu(h.adminDB, "product.audit.submit")
	r.GET("/product-stores", platform, h.listProductStores)
	r.GET("/product-stores/:id/options", platform, h.productStoreOptions)
	r.POST("/products", platform, audit, h.createAdmin)
	r.GET("/products/:id/edit", h.getEdit)
	r.PUT("/products/:id", platform, audit, h.updateAdmin)
	r.GET("/products/:id/operate-logs", h.listOperateLogs)
}

type productStoreOption struct {
	StoreID      uint64 `json:"store_id"`
	StoreName    string `json:"store_name"`
	MerchantID   uint64 `json:"merchant_id"`
	MerchantName string `json:"merchant_name"`
}

func (h *Handler) listProductStores(c *gin.Context) {
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_store AS s").
		Select("s.id AS store_id,s.name AS store_name,m.id AS merchant_id,m.name AS merchant_name").
		Joins("JOIN qixi_crm_m_merchant AS m ON m.id = s.merchant_id").
		Order("s.id ASC")
	if merchantIDs != nil {
		if len(merchantIDs) == 0 {
			response.OK(c, gin.H{"list": []productStoreOption{}})
			return
		}
		q = q.Where("m.id IN ?", merchantIDs)
	}
	var rows []productStoreOption
	if err := q.Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询店铺失败")
		return
	}
	response.OK(c, gin.H{"list": rows})
}

func (h *Handler) productStoreOptions(c *gin.Context) {
	storeID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if storeID == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺 ID 错误")
		return
	}
	var cnt int64
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_store").
		Where("id = ?", storeID).Count(&cnt).Error; err != nil || cnt == 0 {
		response.Fail(c, http.StatusNotFound, "店铺不存在")
		return
	}
	response.OK(c, gin.H{
		"mer_cate_options":  h.loadStoreCategoryOptions(c, storeID),
		"mer_label_options": h.loadStoreTagOptions(c, storeID),
	})
}

type productEditSKU struct {
	SKUId       uint64             `json:"sku_id"`
	Spec        map[string]string  `json:"spec"`
	SpecText    string             `json:"spec_text"`
	Image       string             `json:"image"`
	Price       float64            `json:"price"`
	OtPrice     float64            `json:"ot_price"`
	Stock       int                `json:"stock"`
	BarCode     string             `json:"bar_code"`
	Code        string             `json:"code"`
	Weight      float64            `json:"weight"`
	Volume      float64            `json:"volume"`
	ExtensionOne float64           `json:"extension_one"`
	Status      int8               `json:"status"`
}

type productEditDetail struct {
	ProductID       uint64           `json:"product_id"`
	MerID           uint64           `json:"mer_id"`
	MerName         string           `json:"mer_name"`
	StoreID         uint64           `json:"store_id"`
	StoreName       string           `json:"store_name"`
	Title           string           `json:"title"`
	StoreInfo       string           `json:"store_info"`
	Keyword         string           `json:"keyword"`
	UnitName        string           `json:"unit_name"`
	Image           string           `json:"image"`
	SliderImage     []string         `json:"slider_image"`
	CateID          uint64           `json:"cate_id"`
	CateName        string           `json:"cate_name"`
	CatePath        string           `json:"cate_path"`
	MerCateID       uint64           `json:"mer_cate_id"`
	MerCateName     string           `json:"mer_cate_name"`
	MerCateOptions  []idNameOption   `json:"mer_cate_options"`
	BrandName       string           `json:"brand_name"`
	Price           float64          `json:"price"`
	OtPrice         float64          `json:"ot_price"`
	Stock           int              `json:"stock"`
	Sales           int              `json:"sales"`
	Status          int8             `json:"status"`
	IsShow          int8             `json:"is_show"`
	SpecType        int8             `json:"spec_type"`
	ProductType     int8             `json:"product_type"`
	SVIPPriceType   int8             `json:"svip_price_type"`
	Star            int8             `json:"star"`
	Rank            int              `json:"rank"`
	IsHot           int8             `json:"is_hot"`
	IsBest          int8             `json:"is_best"`
	IsBenefit       int8             `json:"is_benefit"`
	IsNew           int8             `json:"is_new"`
	CateHot         int8             `json:"cate_hot"`
	SysLabels       []string         `json:"sys_labels"`
	MerLabels       []string         `json:"mer_labels"`
	MerLabelIDs     []string         `json:"mer_label_ids"`
	MerLabelOptions []idNameOption   `json:"mer_label_options"`
	ContentHTML     string           `json:"content"`
	RefundSwitch    int8             `json:"refund_switch"`
	OnceMinCount    int              `json:"once_min_count"`
	DeliveryWay     []int            `json:"delivery_way"`
	MerFormID       *uint64          `json:"mer_form_id"`
	MerCategoryName string           `json:"mer_category_name"`
	MerTypeName     string           `json:"mer_type_name"`
	CreateTime      string           `json:"create_time"`
	Refusal         string           `json:"refusal,omitempty"`
	SKUs            []productEditSKU `json:"skus"`
	MerParams       []kvParam        `json:"mer_params"`
	PlatformParams  []kvParam        `json:"platform_params"`
	// 详情页「营销信息」只读字段（缺表/缺列时返回默认或空，前端展示 —）
	MerRecommend      int8     `json:"mer_recommend"`
	CareCount         int      `json:"care_count"`
	Ficti             int      `json:"ficti"`
	SVIPPrice         float64  `json:"svip_price"`
	CommissionText    string   `json:"commission_text"`
	ActivityLabels    []string `json:"activity_labels"`
	SysLabelNames     []string `json:"sys_label_names"`
	IsGiftBag        int8     `json:"is_gift_bag"`
}

type kvParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type idNameOption struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) getEdit(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
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
	detail, err := h.buildEditDetail(c, row)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "加载商品编辑数据失败")
		return
	}
	response.OK(c, detail)
}

func (h *Handler) buildEditDetail(c *gin.Context, row productRow) (*productEditDetail, error) {
	var d detailRow
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").Where("product_id = ?", row.ID).Scan(&d).Error

	var media []struct {
		URL       string `gorm:"column:url"`
		MediaType string `gorm:"column:media_type"`
		Sort      int    `gorm:"column:sort"`
	}
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_media").
		Where("product_id = ? AND media_type = ?", row.ID, "image").Order("sort ASC,id ASC").Find(&media).Error
	slider := make([]string, 0, len(media))
	for _, m := range media {
		if strings.TrimSpace(m.URL) != "" {
			slider = append(slider, m.URL)
		}
	}

	var skus []skuRow
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_sku").
		Where("product_id = ? AND status = 1", row.ID).Order("id ASC").Find(&skus).Error

	ot := 0.0
	if d.OriginalPrice != nil {
		ot = *d.OriginalPrice
	}
	editSKUs := make([]productEditSKU, 0, len(skus))
	stockTotal := 0
	price := 0.0
	for i, s := range skus {
		stockTotal += s.Stock
		if i == 0 {
			price = s.Price
		}
		specMap, specText := parseSpecJSON(s.SpecJSON)
		skuOt := s.OtPrice
		// 旧数据仅商品级划线价时，回填到 SKU 展示
		if skuOt <= 0 && ot > 0 {
			skuOt = ot
		}
		editSKUs = append(editSKUs, productEditSKU{
			SKUId:        s.ID,
			Spec:         specMap,
			SpecText:     specText,
			Image:        s.Image,
			Price:        s.Price,
			OtPrice:      skuOt,
			Stock:        s.Stock,
			BarCode:      s.BarCode,
			Code:         s.Code,
			Weight:       s.Weight,
			Volume:       s.Volume,
			ExtensionOne: s.ExtensionOne,
			Status:       1,
		})
	}

	var view struct {
		Sales       int  `gorm:"column:sales"`
		Stock       int  `gorm:"column:stock"`
		ProductType int8 `gorm:"column:product_type"`
	}
	_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
		Select("sales,stock,product_type").Where("product_id = ?", row.ID).Scan(&view).Error
	if stockTotal == 0 && view.Stock > 0 {
		stockTotal = view.Stock
	}

	cateName, catePath := h.platformCategoryPath(c, row.CategoryID)
	if cateName == "" && row.CategoryID > 0 {
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_category_view").
			Select("name").Where("category_id = ?", row.CategoryID).Scan(&cateName).Error
		catePath = cateName
	}

	var ops productOps
	_ = h.adminDB.WithContext(c.Request.Context()).Where("product_id = ?", row.ID).Take(&ops).Error
	isShow := int8(1)
	if ops.ProductID > 0 {
		isShow = ops.IsUsed
	}

	var review productReview
	_ = h.adminDB.WithContext(c.Request.Context()).Where("product_id = ?", row.ID).Order("id DESC").Take(&review).Error

	merLabelIDs, merLabels := h.loadMerTagBindings(c, row.ID)
	sysLabels := splitCSV(ops.SysLabels)
	merCateName := h.loadStoreCategoryName(c, row.StoreID, row.StoreCategoryID)
	merCateOptions := h.loadStoreCategoryOptions(c, row.StoreID)
	merLabelOptions := h.loadStoreTagOptions(c, row.StoreID)
	deliveryWay := parseDeliveryWayCSV(d.DeliveryWay)
	if len(deliveryWay) == 0 {
		deliveryWay = []int{2}
	}

	merCategoryName, merTypeName := h.loadMerchantMetaNames(c, row.MerchantID)

	var attrs []struct {
		Name  string          `gorm:"column:attribute_name"`
		Value json.RawMessage `gorm:"column:attribute_value"`
	}
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_attribute").
		Where("product_id = ?", row.ID).Find(&attrs).Error
	merParams := make([]kvParam, 0, len(attrs))
	for _, a := range attrs {
		merParams = append(merParams, kvParam{Name: a.Name, Value: stringifyAttrValue(a.Value)})
	}

	specType := int8(0)
	if len(skus) > 1 {
		specType = 1
	}

	refund := int8(1)
	onceMin := 1
	if ops.ProductID > 0 {
		refund = ops.RefundSwitch
		if ops.OnceMinCount > 0 {
			onceMin = ops.OnceMinCount
		}
	}

	sysLabelNames := h.resolveSysLabelNames(c, sysLabels)
	return &productEditDetail{
		ProductID:       row.ID,
		MerID:           row.MerchantID,
		MerName:         row.MerchantName,
		StoreID:         row.StoreID,
		StoreName:       row.StoreName,
		Title:           row.Title,
		StoreInfo:       d.Brief,
		Keyword:         d.Keyword,
		UnitName:        d.UnitName,
		Image:           d.CoverURL,
		SliderImage:     slider,
		CateID:          row.CategoryID,
		CateName:        cateName,
		CatePath:        catePath,
		MerCateID:       row.StoreCategoryID,
		MerCateName:     merCateName,
		MerCateOptions:  merCateOptions,
		BrandName:       row.BrandName,
		Price:           price,
		OtPrice:         ot,
		Stock:           stockTotal,
		Sales:           view.Sales,
		Status:          statusCode(row.Status),
		IsShow:          isShow,
		SpecType:        specType,
		ProductType:     view.ProductType,
		SVIPPriceType:   row.SVIPPriceType,
		Star:            ops.Star,
		Rank:            ops.RankSort,
		IsHot:           ops.IsHot,
		IsBest:          ops.IsBest,
		IsBenefit:       ops.IsBenefit,
		IsNew:           ops.IsNew,
		CateHot:         ops.CateHot,
		SysLabels:       sysLabels,
		MerLabels:       merLabels,
		MerLabelIDs:     merLabelIDs,
		MerLabelOptions: merLabelOptions,
		ContentHTML:     ops.ContentHTML,
		RefundSwitch:    refund,
		OnceMinCount:    onceMin,
		DeliveryWay:     deliveryWay,
		MerFormID:       nil,
		MerCategoryName: merCategoryName,
		MerTypeName:     merTypeName,
		CreateTime:      row.CreatedAt.Format("2006-01-02 15:04:05"),
		Refusal:         review.Reason,
		SKUs:            editSKUs,
		MerParams:       merParams,
		PlatformParams:  []kvParam{},
		MerRecommend:    0,
		CareCount:       0,
		Ficti:           ops.Ficti,
		SVIPPrice:       row.SVIPPrice,
		CommissionText:  "—",
		ActivityLabels:  []string{},
		SysLabelNames:   sysLabelNames,
		IsGiftBag:      row.IsGiftBag,
	}, nil
}

func (h *Handler) resolveSysLabelNames(c *gin.Context, ids []string) []string {
	if len(ids) == 0 {
		return []string{}
	}
	numIDs := make([]uint64, 0, len(ids))
	for _, id := range ids {
		n, err := strconv.ParseUint(strings.TrimSpace(id), 10, 64)
		if err != nil || n == 0 {
			continue
		}
		numIDs = append(numIDs, n)
	}
	if len(numIDs) == 0 {
		return ids
	}
	var rows []struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_product_label").
		Select("id,name").Where("id IN ?", numIDs).Find(&rows).Error
	nameByID := map[uint64]string{}
	for _, r := range rows {
		nameByID[r.ID] = r.Name
	}
	out := make([]string, 0, len(numIDs))
	for _, id := range numIDs {
		if name := nameByID[id]; name != "" {
			out = append(out, name)
		}
	}
	return out
}

func parseSpecJSON(raw []byte) (map[string]string, string) {
	out := map[string]string{}
	if len(raw) == 0 {
		return out, ""
	}
	var anyMap map[string]any
	if err := json.Unmarshal(raw, &anyMap); err != nil {
		return out, string(raw)
	}
	parts := make([]string, 0, len(anyMap))
	for k, v := range anyMap {
		s := strings.TrimSpace(toString(v))
		out[k] = s
		if s != "" {
			parts = append(parts, k+":"+s)
		}
	}
	return out, strings.Join(parts, " / ")
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case json.Number:
		return t.String()
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func stringifyAttrValue(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		return strings.Join(arr, "、")
	}
	return string(raw)
}

func splitCSV(v string) []string {
	v = strings.TrimSpace(v)
	if v == "" {
		return []string{}
	}
	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func (h *Handler) loadMerTagBindings(c *gin.Context, productID uint64) (ids []string, names []string) {
	var rows []struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	_ = h.merchantDB.WithContext(c.Request.Context()).
		Table("qixi_crm_m_product_tag_binding AS b").
		Select("t.id,t.name").
		Joins("JOIN qixi_crm_m_product_tag AS t ON t.id = b.tag_id").
		Where("b.product_id = ?", productID).
		Order("t.sort DESC,t.id ASC").
		Find(&rows).Error
	ids = make([]string, 0, len(rows))
	names = make([]string, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, strconv.FormatUint(r.ID, 10))
		if strings.TrimSpace(r.Name) != "" {
			names = append(names, r.Name)
		}
	}
	return ids, names
}

func (h *Handler) loadStoreCategoryName(c *gin.Context, storeID, cateID uint64) string {
	if cateID == 0 {
		return ""
	}
	var name string
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_category").
		Select("name").Where("id = ? AND store_id = ?", cateID, storeID).Scan(&name).Error
	return name
}

func (h *Handler) loadStoreCategoryOptions(c *gin.Context, storeID uint64) []idNameOption {
	var rows []struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_category").
		Select("id,name").Where("store_id = ? AND status = 1", storeID).
		Order("sort ASC,id ASC").Find(&rows).Error
	out := make([]idNameOption, 0, len(rows))
	for _, r := range rows {
		out = append(out, idNameOption{ID: r.ID, Name: r.Name})
	}
	return out
}

func (h *Handler) loadStoreTagOptions(c *gin.Context, storeID uint64) []idNameOption {
	var rows []struct {
		ID   uint64 `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_tag").
		Select("id,name").Where("store_id = ? AND status = 1", storeID).
		Order("sort DESC,id ASC").Find(&rows).Error
	out := make([]idNameOption, 0, len(rows))
	for _, r := range rows {
		out = append(out, idNameOption{ID: r.ID, Name: r.Name})
	}
	return out
}

func (h *Handler) platformCategoryPath(c *gin.Context, id uint64) (name, path string) {
	if id == 0 {
		return "", ""
	}
	parts := make([]string, 0, 4)
	cur := id
	for i := 0; i < 8 && cur > 0; i++ {
		var row struct {
			ID       uint64 `gorm:"column:id"`
			ParentID uint64 `gorm:"column:parent_id"`
			Name     string `gorm:"column:name"`
		}
		if err := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_platform_category").
			Select("id,parent_id,name").Where("id = ?", cur).Take(&row).Error; err != nil || row.ID == 0 {
			break
		}
		parts = append([]string{row.Name}, parts...)
		if name == "" {
			name = row.Name
		}
		cur = row.ParentID
	}
	if len(parts) == 0 {
		return "", ""
	}
	return name, strings.Join(parts, " / ")
}

func parseDeliveryWayCSV(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	seen := map[int]struct{}{}
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || (n != 1 && n != 2) {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
}

func encodeDeliveryWayCSV(ways []int) string {
	seen := map[int]struct{}{}
	parts := make([]string, 0, 2)
	for _, n := range ways {
		if n != 1 && n != 2 {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		parts = append(parts, strconv.Itoa(n))
	}
	if len(parts) == 0 {
		return "2"
	}
	return strings.Join(parts, ",")
}

func (h *Handler) loadMerchantMetaNames(c *gin.Context, merchantID uint64) (categoryName, typeName string) {
	var meta struct {
		CategoryID uint64 `gorm:"column:category_id"`
		TypeID     uint64 `gorm:"column:type_id"`
	}
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_view").
		Select("category_id,type_id").Where("merchant_id = ?", merchantID).Scan(&meta).Error
	if meta.CategoryID > 0 {
		_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_category").
			Select("name").Where("id = ?", meta.CategoryID).Scan(&categoryName).Error
	}
	if meta.TypeID > 0 {
		_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_merchant_type").
			Select("name").Where("id = ?", meta.TypeID).Scan(&typeName).Error
	}
	if categoryName == "" {
		categoryName = "—"
	}
	if typeName == "" {
		typeName = "—"
	}
	return categoryName, typeName
}

type adminSKUInput struct {
	SKUId        uint64            `json:"sku_id"`
	Spec         map[string]string `json:"spec"`
	Image        string            `json:"image"`
	Price        float64           `json:"price"`
	OtPrice      float64           `json:"ot_price"`
	Stock        int               `json:"stock"`
	Code         string            `json:"code"`
	BarCode      string            `json:"bar_code"`
	Weight       float64           `json:"weight"`
	Volume       float64           `json:"volume"`
	ExtensionOne float64           `json:"extension_one"`
	Status       *int8             `json:"status"`
}

type adminUpdateReq struct {
	Title        *string         `json:"title"`
	StoreInfo    *string         `json:"store_info"`
	Keyword      *string         `json:"keyword"`
	UnitName     *string         `json:"unit_name"`
	BrandName    *string         `json:"brand_name"`
	CateID       *uint64         `json:"cate_id"`
	MerCateID    *uint64         `json:"mer_cate_id"`
	MerLabelIDs  []string        `json:"mer_label_ids"`
	DeliveryWay  []int           `json:"delivery_way"`
	Image        *string         `json:"image"`
	SliderImage  []string        `json:"slider_image"`
	OtPrice      *float64        `json:"ot_price"`
	Skus         []adminSKUInput `json:"skus"`
	Star         *int8           `json:"star"`
	Rank         *int            `json:"rank"`
	IsHot        *int8           `json:"is_hot"`
	IsBest       *int8           `json:"is_best"`
	IsBenefit    *int8           `json:"is_benefit"`
	IsNew        *int8           `json:"is_new"`
	CateHot      *int8           `json:"cate_hot"`
	SysLabels    []string        `json:"sys_labels"`
	Content      *string         `json:"content"`
	RefundSwitch *int8           `json:"refund_switch"`
	OnceMinCount *int            `json:"once_min_count"`
}

type adminCreateReq struct {
	StoreID uint64 `json:"store_id"`
	adminUpdateReq
}

func validateAdminBasic(req *adminUpdateReq) string {
	if req.Title == nil || strings.TrimSpace(*req.Title) == "" {
		return "请填写商品名称"
	}
	if req.Image == nil || strings.TrimSpace(*req.Image) == "" {
		return "请选择封面图"
	}
	if req.SliderImage == nil || len(req.SliderImage) == 0 {
		return "请至少选择一张轮播图"
	}
	if req.StoreInfo == nil || strings.TrimSpace(*req.StoreInfo) == "" {
		return "请填写商品简介"
	}
	if req.CateID == nil || *req.CateID == 0 {
		return "请选择平台分类"
	}
	if req.MerCateID == nil || *req.MerCateID == 0 {
		return "请选择店铺分类"
	}
	if req.MerLabelIDs == nil || len(req.MerLabelIDs) == 0 {
		return "请选择商品标签"
	}
	if req.SysLabels == nil || len(req.SysLabels) == 0 {
		return "请选择平台标签"
	}
	if req.BrandName == nil || strings.TrimSpace(*req.BrandName) == "" {
		return "请选择或填写品牌"
	}
	if req.UnitName == nil || strings.TrimSpace(*req.UnitName) == "" {
		return "请填写单位"
	}
	if req.Keyword == nil || strings.TrimSpace(*req.Keyword) == "" {
		return "请填写关键字"
	}
	if req.DeliveryWay == nil || len(req.DeliveryWay) == 0 {
		return "请至少选择一种配送方式"
	}
	if len(req.Skus) == 0 {
		return "请至少配置一个规格/SKU"
	}
	for _, sku := range req.Skus {
		if sku.Price < 0 {
			return "规格售价不能为负数"
		}
		if sku.OtPrice < 0 {
			return "规格划线价不能为负数"
		}
		if sku.Stock < 0 {
			return "规格库存不能为负数"
		}
		if sku.Weight < 0 {
			return "规格重量不能为负数"
		}
		if sku.Volume < 0 {
			return "规格体积不能为负数"
		}
		if sku.ExtensionOne < 0 {
			return "一级返佣不能为负数"
		}
	}
	return ""
}

func (h *Handler) createAdmin(c *gin.Context) {
	var req adminCreateReq
	if c.ShouldBindJSON(&req) != nil || req.StoreID == 0 {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	if msg := validateAdminBasic(&req.adminUpdateReq); msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	var store struct {
		ID         uint64 `gorm:"column:id"`
		MerchantID uint64 `gorm:"column:merchant_id"`
	}
	q := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_store AS s").
		Select("s.id,s.merchant_id").Where("s.id = ?", req.StoreID)
	if merchantIDs != nil {
		q = q.Where("s.merchant_id IN ?", merchantIDs)
	}
	if err := q.Scan(&store).Error; err != nil || store.ID == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺不存在或无权限")
		return
	}
	var merCateCnt int64
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_category").
		Where("id = ? AND store_id = ?", *req.MerCateID, store.ID).Count(&merCateCnt).Error; err != nil || merCateCnt == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺分类不存在")
		return
	}

	type productInsert struct {
		ID              uint64 `gorm:"column:id;primaryKey"`
		StoreID         uint64 `gorm:"column:store_id"`
		Title           string `gorm:"column:title"`
		CategoryID      uint64 `gorm:"column:category_id"`
		StoreCategoryID uint64 `gorm:"column:store_category_id"`
		BrandName       string `gorm:"column:brand_name"`
		Status          string `gorm:"column:status"`
		Version         uint64 `gorm:"column:version"`
	}
	created := productInsert{
		StoreID: store.ID, Title: strings.TrimSpace(*req.Title), CategoryID: *req.CateID,
		StoreCategoryID: *req.MerCateID, BrandName: strings.TrimSpace(*req.BrandName),
		Status: "on_sale", Version: 1,
	}
	err = h.merchantDB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("qixi_crm_m_product").Create(&created).Error; err != nil {
			return err
		}
		ot := 0.0
		if req.OtPrice != nil {
			ot = *req.OtPrice
		} else if len(req.Skus) > 0 {
			ot = req.Skus[0].OtPrice
		}
		detail := map[string]any{
			"product_id":     created.ID,
			"brief":          strings.TrimSpace(*req.StoreInfo),
			"keyword":        strings.TrimSpace(*req.Keyword),
			"unit_name":      strings.TrimSpace(*req.UnitName),
			"cover_url":      strings.TrimSpace(*req.Image),
			"delivery_way":   encodeDeliveryWayCSV(req.DeliveryWay),
			"original_price": ot,
		}
		return tx.Table("qixi_crm_m_product_detail").Create(detail).Error
	})
	if err != nil || created.ID == 0 {
		response.Fail(c, http.StatusInternalServerError, "创建商品失败")
		return
	}
	productID := created.ID

	if err := h.replaceProductSlider(c, productID, req.SliderImage); err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存轮播图失败")
		return
	}
	if err := h.replaceMerTagBindings(c, productID, store.ID, req.MerLabelIDs); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}
	if _, firstOt, _, err := h.replaceProductSKUs(c, productID, req.Skus); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	} else if req.OtPrice == nil {
		_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").
			Where("product_id = ?", productID).Update("original_price", firstOt).Error
	}

	ops := productOps{
		ProductID: productID, IsUsed: 1, RefundSwitch: 1, OnceMinCount: 1,
		UpdatedBy: uint64(middleware.AdminID(c)), UpdatedAt: time.Now(),
	}
	if req.Star != nil {
		ops.Star = *req.Star
	}
	if req.Rank != nil {
		ops.RankSort = *req.Rank
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
	ops.SysLabels = strings.Join(req.SysLabels, ",")
	if req.Content != nil {
		ops.ContentHTML = *req.Content
	}
	if req.RefundSwitch != nil {
		ops.RefundSwitch = *req.RefundSwitch
	}
	if req.OnceMinCount != nil {
		ops.OnceMinCount = *req.OnceMinCount
	}
	_ = h.adminDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"star", "rank_sort", "is_hot", "is_best", "is_benefit", "is_new", "cate_hot",
			"sys_labels", "content_html", "refund_switch", "once_min_count", "updated_by", "updated_at",
		}),
	}).Create(&ops).Error

	_ = h.dispatchProjection(c.Request.Context(), productID, "upsert")
	h.writeProductOperateLog(c, productID, "product.create", "新增商品")
	response.OK(c, gin.H{"product_id": productID, "ok": true})
}

func (h *Handler) replaceProductSKUs(c *gin.Context, productID uint64, skus []adminSKUInput) (firstPrice float64, firstOtPrice float64, totalStock int, err error) {
	if len(skus) == 0 {
		return 0, 0, 0, errors.New("请至少配置一个规格/SKU")
	}
	type row struct {
		SpecJSON     string
		Image        string
		Price        float64
		OtPrice      float64
		Stock        int
		Code         string
		BarCode      string
		Weight       float64
		Volume       float64
		ExtensionOne float64
		Status       int8
	}
	rows := make([]row, 0, len(skus))
	for _, sku := range skus {
		spec := sku.Spec
		if len(spec) == 0 {
			spec = map[string]string{"默认": "标准"}
		}
		raw, mErr := json.Marshal(spec)
		if mErr != nil {
			return 0, 0, 0, errors.New("规格数据不合法")
		}
		status := int8(1)
		if sku.Status != nil {
			status = *sku.Status
		}
		if status != 0 && status != 1 {
			status = 1
		}
		rows = append(rows, row{
			SpecJSON:     string(raw),
			Image:        strings.TrimSpace(sku.Image),
			Price:        sku.Price,
			OtPrice:      sku.OtPrice,
			Stock:        sku.Stock,
			Code:         strings.TrimSpace(sku.Code),
			BarCode:      strings.TrimSpace(sku.BarCode),
			Weight:       sku.Weight,
			Volume:       sku.Volume,
			ExtensionOne: sku.ExtensionOne,
			Status:       status,
		})
		if status == 1 {
			totalStock += sku.Stock
		}
	}
	firstPrice = rows[0].Price
	firstOtPrice = rows[0].OtPrice

	tx := h.merchantDB.WithContext(c.Request.Context()).Begin()
	if tx.Error != nil {
		return 0, 0, 0, tx.Error
	}
	if err := tx.Table("qixi_crm_m_product_sku").Where("product_id = ?", productID).Delete(nil).Error; err != nil {
		_ = tx.Rollback()
		return 0, 0, 0, err
	}
	for _, r := range rows {
		if err := tx.Table("qixi_crm_m_product_sku").Create(map[string]any{
			"product_id":    productID,
			"spec_json":     r.SpecJSON,
			"image":         r.Image,
			"price":         r.Price,
			"ot_price":      r.OtPrice,
			"stock":         r.Stock,
			"code":          r.Code,
			"bar_code":      r.BarCode,
			"weight":        r.Weight,
			"volume":        r.Volume,
			"extension_one": r.ExtensionOne,
			"status":        r.Status,
		}).Error; err != nil {
			_ = tx.Rollback()
			return 0, 0, 0, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return 0, 0, 0, err
	}
	return firstPrice, firstOtPrice, totalStock, nil
}

func (h *Handler) updateAdmin(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req adminUpdateReq
	if id == 0 || c.ShouldBindJSON(&req) != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误")
		return
	}
	merchantIDs, err := h.merchantScope(c)
	if err != nil {
		response.Fail(c, http.StatusForbidden, "未配置商品监管数据范围")
		return
	}
	var row productRow
	q := h.base(c, merchantIDs).Where("p.id = ?", id)
	if err := q.Scan(&row).Error; err != nil || row.ID == 0 {
		response.Fail(c, http.StatusNotFound, "商品不存在")
		return
	}

	if msg := validateAdminBasic(&req); msg != "" {
		response.Fail(c, http.StatusBadRequest, msg)
		return
	}

	productUpdates := map[string]any{}
	viewUpdates := map[string]any{}
	title := strings.TrimSpace(*req.Title)
	productUpdates["title"] = title
	viewUpdates["title"] = title
	brand := strings.TrimSpace(*req.BrandName)
	productUpdates["brand_name"] = brand
	viewUpdates["brand_name"] = brand
	productUpdates["category_id"] = *req.CateID
	viewUpdates["category_id"] = *req.CateID
	merCateID := *req.MerCateID
	var merCateCnt int64
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_category").
		Where("id = ? AND store_id = ?", merCateID, row.StoreID).Count(&merCateCnt).Error; err != nil || merCateCnt == 0 {
		response.Fail(c, http.StatusBadRequest, "店铺分类不存在")
		return
	}
	productUpdates["store_category_id"] = merCateID
	productUpdates["version"] = gorm.Expr("version + 1")
	if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product").
		Where("id = ?", id).Updates(productUpdates).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新商品基本信息失败")
		return
	}

	cover := strings.TrimSpace(*req.Image)
	detailUpdates := map[string]any{
		"brief":        strings.TrimSpace(*req.StoreInfo),
		"keyword":      strings.TrimSpace(*req.Keyword),
		"unit_name":    strings.TrimSpace(*req.UnitName),
		"cover_url":    cover,
		"delivery_way": encodeDeliveryWayCSV(req.DeliveryWay),
	}
	if req.OtPrice != nil {
		if *req.OtPrice < 0 {
			response.Fail(c, http.StatusBadRequest, "划线价不能为负数")
			return
		}
		detailUpdates["original_price"] = *req.OtPrice
		viewUpdates["original_price"] = *req.OtPrice
	} else if len(req.Skus) > 0 {
		detailUpdates["original_price"] = req.Skus[0].OtPrice
		viewUpdates["original_price"] = req.Skus[0].OtPrice
	}
	viewUpdates["cover_url"] = cover
	var detailCount int64
	_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").
		Where("product_id = ?", id).Count(&detailCount).Error
	if detailCount > 0 {
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").
			Where("product_id = ?", id).Updates(detailUpdates).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "更新商品详情字段失败")
			return
		}
	} else {
		detailUpdates["product_id"] = id
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").
			Create(detailUpdates).Error; err != nil {
			response.Fail(c, http.StatusInternalServerError, "更新商品详情字段失败")
			return
		}
	}

	if err := h.replaceProductSlider(c, id, req.SliderImage); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新轮播图失败")
		return
	}

	if err := h.replaceMerTagBindings(c, id, row.StoreID, req.MerLabelIDs); err != nil {
		response.Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Skus != nil {
		firstPrice, firstOt, totalStock, err := h.replaceProductSKUs(c, id, req.Skus)
		if err != nil {
			response.Fail(c, http.StatusBadRequest, err.Error())
			return
		}
		viewUpdates["price"] = firstPrice
		viewUpdates["stock"] = totalStock
		if req.OtPrice == nil {
			viewUpdates["original_price"] = firstOt
			_ = h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_detail").
				Where("product_id = ?", id).Update("original_price", firstOt).Error
		}
	}

	if len(viewUpdates) > 0 {
		_ = h.businessDB.WithContext(c.Request.Context()).Table("qixi_crm_b_product_view").
			Where("product_id = ?", id).Updates(viewUpdates).Error
	}
	if row.Status == "on_sale" {
		_ = h.dispatchProjection(c.Request.Context(), id, "upsert")
	}

	ops := productOps{ProductID: id, IsUsed: 1, RefundSwitch: 1, OnceMinCount: 1, UpdatedBy: uint64(middleware.AdminID(c)), UpdatedAt: time.Now()}
	var existing productOps
	if err := h.adminDB.WithContext(c.Request.Context()).Where("product_id = ?", id).Take(&existing).Error; err == nil {
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
	if req.Rank != nil {
		ops.RankSort = *req.Rank
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
	if req.SysLabels != nil {
		ops.SysLabels = strings.Join(req.SysLabels, ",")
	}
	if req.Content != nil {
		ops.ContentHTML = *req.Content
	}
	if req.RefundSwitch != nil {
		ops.RefundSwitch = *req.RefundSwitch
	}
	if req.OnceMinCount != nil {
		if *req.OnceMinCount < 0 {
			response.Fail(c, http.StatusBadRequest, "最少购买件数不能为负数")
			return
		}
		if *req.OnceMinCount == 0 {
			ops.OnceMinCount = 0
		} else {
			ops.OnceMinCount = *req.OnceMinCount
		}
	}
	if err := h.adminDB.WithContext(c.Request.Context()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "product_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"star", "rank_sort", "is_hot", "is_best", "is_benefit", "is_new", "cate_hot",
			"sys_labels", "content_html", "refund_switch", "once_min_count", "updated_by", "updated_at",
		}),
	}).Create(&ops).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "保存运营字段失败")
		return
	}
	h.writeProductOperateLog(c, id, "product.update", "编辑商品")
	response.OK(c, gin.H{"ok": true})
}

func (h *Handler) replaceMerTagBindings(c *gin.Context, productID, storeID uint64, labelIDs []string) error {
	ids := make([]uint64, 0, len(labelIDs))
	seen := map[uint64]struct{}{}
	for _, raw := range labelIDs {
		id, err := strconv.ParseUint(strings.TrimSpace(raw), 10, 64)
		if err != nil || id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		var cnt int64
		if err := h.merchantDB.WithContext(c.Request.Context()).Table("qixi_crm_m_product_tag").
			Where("store_id = ? AND id IN ?", storeID, ids).Count(&cnt).Error; err != nil {
			return err
		}
		if cnt != int64(len(ids)) {
			return errors.New("商品标签不存在或不属于该店铺")
		}
	}
	tx := h.merchantDB.WithContext(c.Request.Context()).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("qixi_crm_m_product_tag_binding").Where("product_id = ?", productID).Delete(nil).Error; err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, id := range ids {
		if err := tx.Table("qixi_crm_m_product_tag_binding").Create(map[string]any{
			"product_id": productID,
			"tag_id":     id,
		}).Error; err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (h *Handler) replaceProductSlider(c *gin.Context, productID uint64, urls []string) error {
	ctx := c.Request.Context()
	tx := h.merchantDB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	if err := tx.Table("qixi_crm_m_product_media").
		Where("product_id = ? AND media_type = ?", productID, "image").
		Delete(nil).Error; err != nil {
		_ = tx.Rollback()
		return err
	}
	for i, raw := range urls {
		url := strings.TrimSpace(raw)
		if url == "" {
			continue
		}
		if err := tx.Table("qixi_crm_m_product_media").Create(map[string]any{
			"product_id": productID,
			"media_type": "image",
			"url":        url,
			"sort":       i,
		}).Error; err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (h *Handler) writeProductOperateLog(c *gin.Context, productID uint64, action, label string) {
	claims := middleware.ClaimsFrom(c)
	role := ""
	if claims != nil && len(claims.Roles) > 0 {
		role = claims.Roles[0]
	}
	_ = h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_operation_log").Create(map[string]any{
		"admin_user_id": middleware.AdminID(c),
		"role_code":     role,
		"action":        action,
		"resource_type": "product",
		"resource_id":   strconv.FormatUint(productID, 10),
		"request_id":    c.GetHeader("X-Request-Id"),
		"created_at":    time.Now(),
	}).Error
	_ = label
}

type operateLogRow struct {
	ID           uint64 `json:"id"`
	Index        int    `json:"index"`
	ActionLabel  string `json:"action_label"`
	Terminal     string `json:"terminal"`
	RoleName     string `json:"role_name"`
	OperatorName string `json:"operator_name"`
	CreatedAt    string `json:"created_at"`
}

func (h *Handler) listOperateLogs(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.Fail(c, http.StatusBadRequest, "商品 ID 错误")
		return
	}
	pageN, limitN := page(c)
	terminal := strings.TrimSpace(c.Query("terminal"))
	dateFrom := strings.TrimSpace(c.Query("date_from"))
	dateTo := strings.TrimSpace(c.Query("date_to"))

	q := h.adminDB.WithContext(c.Request.Context()).Table("qixi_crm_a_operation_log AS l").
		Select("l.id,l.action,l.role_code,l.admin_user_id,l.created_at,COALESCE(u.username,'') AS username").
		Joins("LEFT JOIN qixi_crm_a_admin_user AS u ON u.id = l.admin_user_id").
		Where("l.resource_type = ? AND l.resource_id = ?", "product", strconv.FormatUint(id, 10))
	if terminal != "" {
		q = q.Where("l.role_code = ?", terminal)
	}
	if dateFrom != "" {
		q = q.Where("l.created_at >= ?", dateFrom+" 00:00:00")
	}
	if dateTo != "" {
		q = q.Where("l.created_at <= ?", dateTo+" 23:59:59")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询操作记录失败")
		return
	}
	var rows []struct {
		ID          uint64    `gorm:"column:id"`
		Action      string    `gorm:"column:action"`
		RoleCode    string    `gorm:"column:role_code"`
		AdminUserID uint64    `gorm:"column:admin_user_id"`
		CreatedAt   time.Time `gorm:"column:created_at"`
		Username    string    `gorm:"column:username"`
	}
	if err := q.Order("l.id DESC").Offset((pageN - 1) * limitN).Limit(limitN).Scan(&rows).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, "查询操作记录失败")
		return
	}
	list := make([]operateLogRow, 0, len(rows))
	for i, r := range rows {
		list = append(list, operateLogRow{
			ID:           r.ID,
			Index:        (pageN-1)*limitN + i + 1,
			ActionLabel:  operateActionLabel(r.Action),
			Terminal:     terminalLabel(r.RoleCode),
			RoleName:     roleLabel(r.RoleCode),
			OperatorName: dashIfEmpty(r.Username),
			CreatedAt:    r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response.OK(c, gin.H{"list": list, "total": total, "page": pageN, "limit": limitN})
}

func operateActionLabel(action string) string {
	switch action {
	case "product.update":
		return "编辑商品"
	case "product.force_off":
		return "强制下架"
	case "product.show":
		return "修改显示状态"
	case "product.audit":
		return "审核商品"
	default:
		if action == "" {
			return "—"
		}
		return action
	}
}

func terminalLabel(role string) string {
	switch role {
	case "platform", "operations", "customer_service", "region":
		return "平台"
	case "merchant":
		return "店铺"
	default:
		if role == "" {
			return "—"
		}
		return role
	}
}

func roleLabel(role string) string {
	switch role {
	case "platform":
		return "平台管理员"
	case "operations":
		return "运营"
	case "customer_service":
		return "客服"
	case "region":
		return "区域"
	case "merchant":
		return "店铺账号"
	default:
		return dashIfEmpty(role)
	}
}

func dashIfEmpty(v string) string {
	if strings.TrimSpace(v) == "" {
		return "—"
	}
	return v
}
