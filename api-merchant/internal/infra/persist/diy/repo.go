package diypersist

import (
	"context"

	"github.com/qixi-live/qixi-live-mergers/api-merchant/internal/domain/diy"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func NewStoreAdapter(repo *Repo) *Repo { return repo }

func (r *Repo) List(ctx context.Context, f diy.ListFilter) ([]diy.Page, int64, error) {
	page, limit := f.Page, f.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&diy.Page{}).Where("is_del = 0 AND mer_id = ?", f.MerID)
	if f.IsDiy != nil {
		q = q.Where("is_diy = ?", *f.IsDiy)
	}
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if name := f.Name; name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []diy.Page
	err := q.Order("status DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*diy.Page, error) {
	var row diy.Page
	err := r.db.WithContext(ctx).Where("id = ? AND is_del = 0", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetActiveHome(ctx context.Context, merID uint) (*diy.Page, error) {
	var row diy.Page
	err := r.db.WithContext(ctx).
		Where("mer_id = ? AND status = 1 AND type = 0 AND is_del = 0 AND is_diy = 1", merID).
		Order("id DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) Create(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *Repo) Update(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"name":          p.Name,
		"title":         p.Title,
		"cover_image":   p.CoverImage,
		"template_name": p.TemplateName,
		"value":         p.Value,
		"version":       p.Version,
		"status":        p.Status,
		"type":          p.Type,
		"is_show":       p.IsShow,
		"is_diy":        p.IsDiy,
		"is_bg_color":   p.IsBgColor,
		"is_bg_pic":     p.IsBgPic,
		"color_picker":  p.ColorPicker,
		"bg_pic":        p.BgPic,
		"bg_tab_val":    p.BgTabVal,
		"is_default":    p.IsDefault,
	}).Error
}

func (r *Repo) ClearActive(ctx context.Context, merID uint, isDiy int8) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).
		Where("mer_id = ? AND status = 1 AND is_del = 0 AND is_diy = ?", merID, isDiy).
		Updates(map[string]interface{}{"status": 0, "is_default": 0}).Error
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&diy.Page{}).
		Where("id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) ListCategories(ctx context.Context, isMer int8) ([]diy.PageCategory, error) {
	var rows []diy.PageCategory
	err := r.db.WithContext(ctx).Where("is_mer = ? AND type = ?", isMer, "link").
		Order("sort DESC, add_time DESC, id DESC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetCategory(ctx context.Context, id uint) (*diy.PageCategory, error) {
	var row diy.PageCategory
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateCategory(ctx context.Context, category *diy.PageCategory) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, category *diy.PageCategory) error {
	return r.db.WithContext(ctx).Model(&diy.PageCategory{}).Where("id = ?", category.ID).Updates(map[string]interface{}{
		"pid": category.PID, "name": category.Name, "sort": category.Sort,
		"status": category.Status, "level": category.Level,
	}).Error
}

func (r *Repo) DeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&diy.PageCategory{}, id).Error
}

func (r *Repo) CountCategoryChildren(ctx context.Context, id uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&diy.PageCategory{}).Where("pid = ?", id).Count(&total).Error
	return total, err
}

func (r *Repo) CountLinksByCategory(ctx context.Context, id uint) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).Model(&diy.PageLink{}).Where("cate_id = ?", id).Count(&total).Error
	return total, err
}

func (r *Repo) ListLinks(ctx context.Context, f diy.LinkListFilter) ([]diy.PageLink, int64, error) {
	page, limit := f.Page, f.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&diy.PageLink{}).Where("is_mer = ?", f.IsMer)
	if f.Status != nil {
		q = q.Where("status = ?", *f.Status)
	}
	if f.Name != "" {
		q = q.Where("name LIKE ?", "%"+f.Name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []diy.PageLink
	err := q.Order("sort DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, err
	}
	categoryIDs := make([]uint, 0, len(rows))
	for _, row := range rows {
		categoryIDs = append(categoryIDs, row.CateID)
	}
	if len(categoryIDs) == 0 {
		return rows, total, nil
	}
	var categories []diy.PageCategory
	if err := r.db.WithContext(ctx).Where("id IN ?", categoryIDs).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	byID := make(map[uint]*diy.PageCategory, len(categories))
	for i := range categories {
		byID[categories[i].ID] = &categories[i]
	}
	for i := range rows {
		rows[i].Category = byID[rows[i].CateID]
	}
	return rows, total, nil
}

func (r *Repo) GetLink(ctx context.Context, id uint) (*diy.PageLink, error) {
	var row diy.PageLink
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) CreateLink(ctx context.Context, link *diy.PageLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *Repo) UpdateLink(ctx context.Context, link *diy.PageLink) error {
	return r.db.WithContext(ctx).Model(&diy.PageLink{}).Where("id = ?", link.ID).Updates(map[string]interface{}{
		"cate_id": link.CateID, "name": link.Name, "url": link.URL, "param": link.Param,
		"example": link.Example, "status": link.Status, "sort": link.Sort,
	}).Error
}

func (r *Repo) DeleteLink(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&diy.PageLink{}, id).Error
}

var _ diy.Store = (*Repo)(nil)
