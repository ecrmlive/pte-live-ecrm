package article

import (
	"context"

	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/article"
	"gorm.io/gorm"
)

type Repo struct{ db *gorm.DB }

func NewRepo(db *gorm.DB) *Repo { return &Repo{db: db} }

func (r *Repo) ListCategory(ctx context.Context) ([]article.Category, error) {
	var rows []article.Category
	err := r.db.WithContext(ctx).Where("is_del = 0").Order("sort DESC, cid ASC").Find(&rows).Error
	return rows, err
}

func (r *Repo) GetCategory(ctx context.Context, id uint) (*article.Category, error) {
	var row article.Category
	err := r.db.WithContext(ctx).Where("cid = ? AND is_del = 0", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateCategory(ctx context.Context, row *article.Category) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateCategory(ctx context.Context, row *article.Category) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteCategory(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&article.Category{}).Where("cid = ?", id).Update("is_del", 1).Error
}

func (r *Repo) ListArticle(ctx context.Context, page, limit int, cid uint, publishedOnly bool) ([]article.Article, int64, error) {
	var total int64
	q := r.db.WithContext(ctx).Model(&article.Article{}).Where("is_del = 0")
	if cid > 0 {
		q = q.Where("cid = ?", cid)
	}
	if publishedOnly {
		q = q.Where("status = 1")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []article.Article
	err := q.Order("sort DESC, article_id DESC").Offset((page - 1) * limit).Limit(limit).Find(&rows).Error
	return rows, total, err
}

func (r *Repo) GetArticle(ctx context.Context, id uint) (*article.Article, error) {
	var row article.Article
	err := r.db.WithContext(ctx).Where("article_id = ? AND is_del = 0", id).First(&row).Error
	return &row, err
}

func (r *Repo) CreateArticle(ctx context.Context, row *article.Article) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *Repo) UpdateArticle(ctx context.Context, row *article.Article) error {
	return r.db.WithContext(ctx).Save(row).Error
}

func (r *Repo) SoftDeleteArticle(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&article.Article{}).Where("article_id = ?", id).Update("is_del", 1).Error
}

func (r *Repo) IncrVisit(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&article.Article{}).Where("article_id = ?", id).
		UpdateColumn("visit", gorm.Expr("visit + 1")).Error
}
