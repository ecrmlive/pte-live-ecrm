package diypersist

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/diy"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repo 只访问 qixi_crm_merchant。用户端通过 outbox/NATS 获取投影，不能把
// 店铺装修表作为 C 端的跨库数据源。
type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo        { return &Repo{db: db} }
func NewStoreAdapter(repo *Repo) *Repo { return repo }

type pageRow struct {
	ID        uint            `gorm:"column:id;primaryKey"`
	StoreID   uint            `gorm:"column:store_id"`
	Name      string          `gorm:"column:name"`
	Document  json.RawMessage `gorm:"column:document"`
	PageType  string          `gorm:"column:page_type"`
	IsActive  bool            `gorm:"column:is_active"`
	Status    string          `gorm:"column:status"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

func (pageRow) TableName() string { return "qixi_crm_m_diy_page" }

type outboxRow struct {
	EventType     string          `gorm:"column:event_type"`
	AggregateType string          `gorm:"column:aggregate_type"`
	AggregateID   string          `gorm:"column:aggregate_id"`
	Payload       json.RawMessage `gorm:"column:payload"`
}

func (outboxRow) TableName() string { return "qixi_crm_m_outbox" }

type diyEventPayload struct {
	PageID   uint            `json:"page_id"`
	StoreID  uint            `json:"store_id"`
	PageType string          `json:"page_type"`
	Name     string          `json:"name"`
	Document json.RawMessage `json:"document"`
	Status   string          `json:"status"`
	IsActive bool            `json:"is_active"`
}

const (
	eventUpsert = "merchant.diy_page.upsert"
	eventDelete = "merchant.diy_page.deleted"
)

func (r *Repo) List(ctx context.Context, f diy.ListFilter) ([]diy.Page, int64, error) {
	page, limit := normalizePage(f.Page, f.Limit)
	q := r.db.WithContext(ctx).Model(&pageRow{}).Where("store_id = ?", f.MerID)
	if f.IsDiy != nil {
		if *f.IsDiy == 1 {
			q = q.Where("page_type = ?", "home")
		} else {
			q = q.Where("page_type = ?", "custom")
		}
	}
	if f.Status != nil {
		q = q.Where("status = ?", statusValue(*f.Status))
	}
	if f.Name != "" {
		q = q.Where("name LIKE ?", "%"+f.Name+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []pageRow
	if err := q.Order("is_active DESC, updated_at DESC, id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	pages := make([]diy.Page, 0, len(rows))
	for _, row := range rows {
		pages = append(pages, toPage(row))
	}
	return pages, total, nil
}

func (r *Repo) Get(ctx context.Context, id uint) (*diy.Page, error) {
	var row pageRow
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error; err != nil {
		return nil, err
	}
	p := toPage(row)
	return &p, nil
}

func (r *Repo) GetActiveHome(ctx context.Context, storeID uint) (*diy.Page, error) {
	var row pageRow
	err := r.db.WithContext(ctx).Where("store_id = ? AND page_type = ? AND status = ? AND is_active = ?", storeID, "home", "published", true).Order("updated_at DESC, id DESC").First(&row).Error
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
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
		p.ID, p.StoreID, p.MerID, p.UpdateTime = row.ID, row.StoreID, row.StoreID, row.UpdatedAt
		return enqueue(tx, eventUpsert, row)
	})
}

func (r *Repo) Update(ctx context.Context, p *diy.Page) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var current pageRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", p.ID).First(&current).Error; err != nil {
			return err
		}
		row, err := fromPage(p)
		if err != nil {
			return err
		}
		row.StoreID = current.StoreID
		if err := tx.Model(&pageRow{}).Where("id = ?", p.ID).Updates(map[string]any{
			"name": row.Name, "document": row.Document, "page_type": row.PageType, "is_active": row.IsActive, "status": row.Status,
		}).Error; err != nil {
			return err
		}
		row.ID = p.ID
		return enqueue(tx, eventUpsert, row)
	})
}

func (r *Repo) ClearActive(ctx context.Context, storeID uint, isDiy int8) error {
	if isDiy != 1 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var rows []pageRow
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("store_id = ? AND page_type = ? AND is_active = ?", storeID, "home", true).Find(&rows).Error; err != nil {
			return err
		}
		if len(rows) == 0 {
			return nil
		}
		if err := tx.Model(&pageRow{}).Where("store_id = ? AND page_type = ? AND is_active = ?", storeID, "home", true).Updates(map[string]any{"is_active": false, "status": "draft"}).Error; err != nil {
			return err
		}
		for _, row := range rows {
			row.IsActive, row.Status = false, "draft"
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
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&row).Error; err != nil {
			return err
		}
		if err := tx.Delete(&pageRow{}, id).Error; err != nil {
			return err
		}
		return enqueue(tx, eventDelete, row)
	})
}

// 店铺端只能消费固定链接；链接定义属于统一后台，不能再落入旧 qixi_m_admin_* 表。
func (r *Repo) ListCategories(context.Context, int8) ([]diy.PageCategory, error) {
	return defaultCategories(), nil
}
func (r *Repo) GetCategory(_ context.Context, id uint) (*diy.PageCategory, error) {
	for _, v := range defaultCategories() {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *Repo) CreateCategory(context.Context, *diy.PageCategory) error    { return gorm.ErrInvalidData }
func (r *Repo) UpdateCategory(context.Context, *diy.PageCategory) error    { return gorm.ErrInvalidData }
func (r *Repo) DeleteCategory(context.Context, uint) error                 { return gorm.ErrInvalidData }
func (r *Repo) CountCategoryChildren(context.Context, uint) (int64, error) { return 0, nil }
func (r *Repo) CountLinksByCategory(context.Context, uint) (int64, error)  { return 0, nil }
func (r *Repo) ListLinks(_ context.Context, f diy.LinkListFilter) ([]diy.PageLink, int64, error) {
	rows := defaultLinks()
	if f.Name != "" {
		filtered := make([]diy.PageLink, 0, len(rows))
		for _, row := range rows {
			if contains(row.Name, f.Name) {
				filtered = append(filtered, row)
			}
		}
		rows = filtered
	}
	return rows, int64(len(rows)), nil
}
func (r *Repo) GetLink(_ context.Context, id uint) (*diy.PageLink, error) {
	for _, v := range defaultLinks() {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *Repo) CreateLink(context.Context, *diy.PageLink) error { return gorm.ErrInvalidData }
func (r *Repo) UpdateLink(context.Context, *diy.PageLink) error { return gorm.ErrInvalidData }
func (r *Repo) DeleteLink(context.Context, uint) error          { return gorm.ErrInvalidData }

func normalizePage(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
func statusValue(status int8) string {
	if status == 1 {
		return "published"
	}
	return "draft"
}
func contains(s, q string) bool { return q == "" || len(s) >= len(q) && (s == q || jsonContains(s, q)) }
func jsonContains(s, q string) bool {
	for i := 0; i+len(q) <= len(s); i++ {
		if s[i:i+len(q)] == q {
			return true
		}
	}
	return false
}

func fromPage(p *diy.Page) (pageRow, error) {
	doc := make(map[string]any)
	if len(p.Value) > 0 {
		if err := json.Unmarshal([]byte(p.Value), &doc); err != nil {
			return pageRow{}, err
		}
	}
	if doc == nil {
		doc = map[string]any{}
	}
	doc["_qixi"] = map[string]any{"title": p.Title, "cover_image": p.CoverImage, "template_name": p.TemplateName, "is_show": p.IsShow, "is_diy": p.IsDiy, "is_bg_color": p.IsBgColor, "is_bg_pic": p.IsBgPic, "color_picker": p.ColorPicker, "bg_pic": p.BgPic, "bg_tab_val": p.BgTabVal, "is_default": p.IsDefault}
	raw, err := json.Marshal(doc)
	if err != nil {
		return pageRow{}, err
	}
	pageType := "custom"
	if p.IsDiy == 1 {
		pageType = "home"
	}
	return pageRow{ID: p.ID, StoreID: p.MerID, Name: p.Name, Document: raw, PageType: pageType, IsActive: p.IsDefault == 1 && p.Status == 1, Status: statusValue(p.Status)}, nil
}

func toPage(row pageRow) diy.Page {
	p := diy.Page{ID: row.ID, StoreID: row.StoreID, MerID: row.StoreID, Name: row.Name, Value: string(row.Document), Status: 0, Type: 0, IsShow: 1, IsDiy: 0, AddTime: row.UpdatedAt, UpdateTime: row.UpdatedAt}
	if row.Status == "published" {
		p.Status = 1
	}
	if row.PageType == "home" {
		p.IsDiy = 1
	}
	if row.IsActive {
		p.IsDefault = 1
	}
	var doc map[string]any
	if json.Unmarshal(row.Document, &doc) == nil {
		if meta, ok := doc["_qixi"].(map[string]any); ok {
			applyMeta(&p, meta)
		}
	}
	return p
}

func applyMeta(p *diy.Page, meta map[string]any) {
	if s, ok := meta["title"].(string); ok {
		p.Title = s
	}
	if s, ok := meta["cover_image"].(string); ok {
		p.CoverImage = s
	}
	if s, ok := meta["template_name"].(string); ok {
		p.TemplateName = s
	}
	if s, ok := meta["color_picker"].(string); ok {
		p.ColorPicker = s
	}
	if s, ok := meta["bg_pic"].(string); ok {
		p.BgPic = s
	}
	if v, ok := meta["is_show"].(float64); ok {
		p.IsShow = int8(v)
	}
	if v, ok := meta["is_bg_color"].(float64); ok {
		p.IsBgColor = int8(v)
	}
	if v, ok := meta["is_bg_pic"].(float64); ok {
		p.IsBgPic = int8(v)
	}
	if v, ok := meta["bg_tab_val"].(float64); ok {
		p.BgTabVal = int8(v)
	}
}

func enqueue(tx *gorm.DB, eventType string, row pageRow) error {
	payload, err := json.Marshal(diyEventPayload{PageID: row.ID, StoreID: row.StoreID, PageType: row.PageType, Name: row.Name, Document: row.Document, Status: row.Status, IsActive: row.IsActive})
	if err != nil {
		return err
	}
	return tx.Create(&outboxRow{EventType: eventType, AggregateType: "diy_page", AggregateID: strconv.FormatUint(uint64(row.ID), 10), Payload: payload}).Error
}

func defaultCategories() []diy.PageCategory {
	return []diy.PageCategory{{ID: 1, Name: "商城页面", Type: "link", Status: 1, Level: 1, IsMer: 1}}
}
func defaultLinks() []diy.PageLink {
	c := &diy.PageCategory{ID: 1, Name: "商城页面", Type: "link", Status: 1, Level: 1, IsMer: 1}
	return []diy.PageLink{{ID: 1, CateID: 1, Name: "店铺首页", URL: "/pages/store/index", Status: 1, IsMer: 1, Category: c}, {ID: 2, CateID: 1, Name: "商品分类", URL: "/pages/goods_cate/index", Status: 1, IsMer: 1, Category: c}, {ID: 3, CateID: 1, Name: "购物车", URL: "/pages/order_addcart/order_addcart", Status: 1, IsMer: 1, Category: c}}
}

var _ diy.Store = (*Repo)(nil)
