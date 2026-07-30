package catalogpersist

import (
	"context"
	"strconv"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/catalog"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListPlatformCategories(ctx context.Context) ([]catalog.Category, error) {
	var rows []catalog.Category
	err := r.db.WithContext(ctx).
		Where("mer_id = 0 AND type = 0").
		Order("sort DESC, store_category_id ASC").
		Find(&rows).Error
	return rows, err
}

func (r *Repo) GetCategory(ctx context.Context, id uint) (*catalog.Category, error) {
	var row catalog.Category
	err := r.db.WithContext(ctx).Where("store_category_id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateCategory(ctx context.Context, c *catalog.Category) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, c *catalog.Category) error {
	return r.db.WithContext(ctx).Model(&catalog.Category{}).
		Where("store_category_id = ?", c.StoreCategoryID).
		Updates(map[string]interface{}{
			"pid":       c.PID,
			"cate_name": c.CateName,
			"path":      c.Path,
			"sort":      c.Sort,
			"pic":       c.Pic,
			"is_show":   c.IsShow,
			"level":     c.Level,
		}).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalog.Category{}, id).Error
}

func (r *Repo) ListBrands(ctx context.Context) ([]catalog.Brand, error) {
	var rows []catalog.Brand
	err := r.db.WithContext(ctx).Order("sort DESC, brand_id ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) CreateBrand(ctx context.Context, b *catalog.Brand) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *Repo) UpdateBrand(ctx context.Context, b *catalog.Brand) error {
	return r.db.WithContext(ctx).Model(&catalog.Brand{}).
		Where("brand_id = ?", b.BrandID).
		Updates(map[string]interface{}{
			"brand_category_id": b.BrandCategoryID,
			"brand_name":        b.BrandName,
			"sort":              b.Sort,
			"pic":               b.Pic,
			"is_show":           b.IsShow,
		}).Error
}

func (r *Repo) DeleteBrand(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&catalog.Brand{}, id).Error
}

type ListProductFilter struct {
	Status  *int8
	Keyword string
	MerID   *uint
	Page    int
	Limit   int
}

func (r *Repo) ListProducts(ctx context.Context, f ListProductFilter) ([]catalog.Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&catalog.Product{}).Where("is_del = 0")
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.MerID != nil {
		q = q.Where("mer_id = ?", *f.MerID)
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		q = q.Where("store_name LIKE ? OR keyword LIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := f.Page, f.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []catalog.Product
	err := q.Order("status ASC, product_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetProduct(ctx context.Context, id uint) (*catalog.Product, error) {
	var row catalog.Product
	err := r.db.WithContext(ctx).Where("product_id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) UpdateProductAudit(ctx context.Context, id uint, status int8, refusal string) error {
	return r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("product_id = ? AND is_del = 0", id).
		Updates(map[string]interface{}{
			"status":  status,
			"refusal": refusal,
		}).Error
}

func (r *Repo) MerchantName(ctx context.Context, merID uint) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Select("mer_name").Where("mer_id = ?", merID).Scan(&name).Error
	return name, err
}

func (r *Repo) CategoryName(ctx context.Context, cateID int) (string, error) {
	var name string
	err := r.db.WithContext(ctx).Table("qixi_m_admin_store_category").
		Select("cate_name").Where("store_category_id = ?", cateID).Scan(&name).Error
	return name, err
}

func (r *Repo) MerchantNeedAudit(ctx context.Context, merID uint) (bool, error) {
	var isAudit int8
	err := r.db.WithContext(ctx).Table("qixi_m_admin_merchant").
		Select("is_audit").Where("mer_id = ? AND is_del = 0", merID).Scan(&isAudit).Error
	return isAudit == 1, err
}

func (r *Repo) CreateProduct(ctx context.Context, p *catalog.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) UpdateProduct(ctx context.Context, p *catalog.Product) error {
	return r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("product_id = ? AND mer_id = ? AND is_del = 0", p.ProductID, p.MerID).
		Updates(map[string]interface{}{
			"store_name":   p.StoreName,
			"store_info":   p.StoreInfo,
			"keyword":      p.Keyword,
			"brand_id":     p.BrandID,
			"is_show":      p.IsShow,
			"status":       p.Status,
			"refusal":      p.Refusal,
			"cate_id":      p.CateID,
			"unit_name":    p.UnitName,
			"price":        p.Price,
			"ot_price":     p.OtPrice,
			"stock":        p.Stock,
			"spec_type":    p.SpecType,
			"image":        p.Image,
			"slider_image": p.SliderImage,
			"delivery_way": p.DeliveryWay,
			"type":         p.Type,
		}).Error
}

func (r *Repo) SoftDeleteProduct(ctx context.Context, id, merID uint) error {
	return r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("product_id = ? AND mer_id = ?", id, merID).
		Update("is_del", 1).Error
}

func (r *Repo) UpsertDefaultSKU(ctx context.Context, v *catalog.AttrValue) error {
	var existing catalog.AttrValue
	err := r.db.WithContext(ctx).
		Where("product_id = ? AND sku = ?", v.ProductID, v.SKU).
		First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(v).Error
	}
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&catalog.AttrValue{}).
		Where("value_id = ?", existing.ValueID).
		Updates(map[string]interface{}{
			"stock":    v.Stock,
			"price":    v.Price,
			"ot_price": v.OtPrice,
			"cost":     v.Cost,
			"unique":   v.Unique,
			"is_show":  v.IsShow,
		}).Error
}

func (r *Repo) ListDefaultSKUs(ctx context.Context, productIDs []uint) (map[uint]catalog.AttrValue, error) {
	out := map[uint]catalog.AttrValue{}
	if len(productIDs) == 0 {
		return out, nil
	}
	var rows []catalog.AttrValue
	err := r.db.WithContext(ctx).Where("product_id IN ?", productIDs).Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if _, ok := out[row.ProductID]; !ok {
			out[row.ProductID] = row
		}
	}
	return out, nil
}

func (r *Repo) ListCategoryIDsUnder(ctx context.Context, cateID int) ([]int, error) {
	if cateID <= 0 {
		return nil, nil
	}
	var ids []int
	err := r.db.WithContext(ctx).Model(&catalog.Category{}).
		Where("store_category_id = ? OR pid = ? OR path LIKE ?", cateID, cateID, "%/"+itoa(cateID)+"/%").
		Pluck("store_category_id", &ids).Error
	return ids, err
}

func (r *Repo) ListOnSaleProducts(ctx context.Context, cateID int, keyword string, merID *uint, page, limit int) ([]catalog.Product, int64, error) {
	// 默认仅普通商品；积分商品（product_type=20）走 ListPointsProducts
	q := r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("is_del = 0 AND status = 1 AND is_show = 1 AND mer_status = 1 AND product_type <> ?", catalog.ProductTypePoints)
	if merID != nil {
		q = q.Where("mer_id = ?", *merID)
	}
	if cateID > 0 {
		ids, err := r.ListCategoryIDsUnder(ctx, cateID)
		if err != nil {
			return nil, 0, err
		}
		if len(ids) == 0 {
			ids = []int{cateID}
		}
		q = q.Where("cate_id IN ?", ids)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("store_name LIKE ? OR keyword LIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []catalog.Product
	err := q.Order("product_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) ListPointsProducts(ctx context.Context, page, limit int) ([]catalog.Product, int64, error) {
	q := r.db.WithContext(ctx).Model(&catalog.Product{}).
		Where("is_del = 0 AND status = 1 AND is_show = 1 AND mer_status = 1").
		Where("(type = 1 OR product_type = ?)", catalog.ProductTypePoints)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	var rows []catalog.Product
	err := q.Order("product_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
