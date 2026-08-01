package diypersist

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/diy"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func NewStoreAdapter(repo *Repo) *Repo { return repo }

type pageRow struct {
	ID        uint            `gorm:"column:id;primaryKey"`
	PageType  string          `gorm:"column:page_type"`
	Name      string          `gorm:"column:name"`
	Document  json.RawMessage `gorm:"column:document"`
	Status    string          `gorm:"column:status"`
	UpdatedBy uint            `gorm:"column:updated_by"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

func (pageRow) TableName() string { return "qixi_crm_a_diy_page" }

type outboxRow struct {
	EventType     string          `gorm:"column:event_type"`
	AggregateType string          `gorm:"column:aggregate_type"`
	AggregateID   string          `gorm:"column:aggregate_id"`
	Payload       json.RawMessage `gorm:"column:payload"`
}

func (outboxRow) TableName() string { return "qixi_crm_a_outbox" }

const (
	eventUpsert = "platform.diy_page.upsert"
	eventDelete = "platform.diy_page.deleted"
)

func (r *Repo) List(ctx context.Context, f diy.ListFilter) ([]diy.Page, int64, error) {
	page, limit := f.Page, f.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := r.db.WithContext(ctx).Model(&pageRow{})
	if f.IsDiy != nil {
		if *f.IsDiy == 1 {
			q = q.Where("page_type = ?", "home")
		} else {
			q = q.Where("page_type = ?", "custom")
		}
	}
	if f.Status != nil {
		if *f.Status == 1 {
			q = q.Where("status = ?", "published")
		} else {
			q = q.Where("status = ?", "draft")
		}
	}
	if name := f.Name; name != "" {
		q = q.Where("name LIKE ?", "%"+name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []pageRow
	err := q.Order("updated_at DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	pages := make([]diy.Page, 0, len(rows))
	for _, row := range rows {
		pages = append(pages, toPage(row))
	}
	return pages, total, err
}

func (r *Repo) Get(ctx context.Context, id uint) (*diy.Page, error) {
	var row pageRow
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
	if err != nil {
		return nil, err
	}
	p := toPage(row)
	return &p, nil
}

func (r *Repo) GetActiveHome(ctx context.Context, merID uint) (*diy.Page, error) {
	var row pageRow
	err := r.db.WithContext(ctx).
		Where("page_type = ? AND status = ?", "home", "published").
		Order("id DESC").
		First(&row).Error
	if err != nil {
		return nil, err
	}
	p := toPage(row)
	return &p, nil
}

func (r *Repo) Create(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row, err := fromPage(p)
		if err != nil {
			return err
		}
		if err = tx.Create(&row).Error; err != nil {
			return err
		}
		p.ID = row.ID
		return enqueue(tx, eventUpsert, row)
	})
}

func (r *Repo) Update(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		row, err := fromPage(p)
		if err != nil {
			return err
		}
		if err = tx.Model(&pageRow{}).Where("id = ?", p.ID).Updates(map[string]interface{}{"name": row.Name, "document": row.Document, "page_type": row.PageType, "status": row.Status}).Error; err != nil {
			return err
		}
		return enqueue(tx, eventUpsert, row)
	})
}

func (r *Repo) ClearActive(ctx context.Context, merID uint, isDiy int8) error {
	if isDiy != 1 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rows []pageRow
		if err := tx.Where("page_type = ? AND status = ?", "home", "published").Find(&rows).Error; err != nil {
			return err
		}
		if err := tx.Model(&pageRow{}).Where("page_type = ? AND status = ?", "home", "published").Update("status", "draft").Error; err != nil {
			return err
		}
		for _, row := range rows {
			row.Status = "draft"
			if err := enqueue(tx, eventUpsert, row); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repo) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var row pageRow
		if err := tx.Where("id = ?", id).First(&row).Error; err != nil {
			return err
		}
		if err := tx.Delete(&pageRow{}, id).Error; err != nil {
			return err
		}
		return enqueue(tx, eventDelete, row)
	})
}

func fromPage(p *diy.Page) (pageRow, error) {
	var doc map[string]any
	if err := json.Unmarshal([]byte(p.Value), &doc); err != nil {
		return pageRow{}, err
	}
	doc["_qixi"] = map[string]any{"title": p.Title, "cover_image": p.CoverImage, "template_name": p.TemplateName, "is_show": p.IsShow, "is_diy": p.IsDiy, "is_default": p.IsDefault}
	raw, err := json.Marshal(doc)
	if err != nil {
		return pageRow{}, err
	}
	pageType := "custom"
	if p.IsDiy == 1 {
		pageType = "home"
	}
	status := "draft"
	if p.Status == 1 {
		status = "published"
	}
	return pageRow{ID: p.ID, PageType: pageType, Name: p.Name, Document: raw, Status: status}, nil
}
func toPage(row pageRow) diy.Page {
	p := diy.Page{ID: row.ID, Name: row.Name, Value: string(row.Document), Status: 0, IsShow: 1, Type: 0, IsDiy: 0, UpdateTime: row.UpdatedAt, AddTime: row.UpdatedAt}
	if row.Status == "published" {
		p.Status = 1
	}
	if row.PageType == "home" {
		p.IsDiy = 1
	}
	var doc map[string]any
	if json.Unmarshal(row.Document, &doc) == nil {
		if m, ok := doc["_qixi"].(map[string]any); ok {
			if s, ok := m["title"].(string); ok {
				p.Title = s
			}
			if s, ok := m["cover_image"].(string); ok {
				p.CoverImage = s
			}
			if s, ok := m["template_name"].(string); ok {
				p.TemplateName = s
			}
			if v, ok := m["is_show"].(float64); ok {
				p.IsShow = int8(v)
			}
			if v, ok := m["is_default"].(float64); ok {
				p.IsDefault = int8(v)
			}
		}
	}
	return p
}

func enqueue(tx *gorm.DB, eventType string, row pageRow) error {
	payload, err := json.Marshal(map[string]any{"page_id": row.ID, "page_type": row.PageType, "name": row.Name, "document": json.RawMessage(row.Document), "status": row.Status, "is_active": row.PageType == "home" && row.Status == "published"})
	if err != nil {
		return err
	}
	return tx.Create(&outboxRow{EventType: eventType, AggregateType: "diy_page", AggregateID: strconv.FormatUint(uint64(row.ID), 10), Payload: payload}).Error
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
