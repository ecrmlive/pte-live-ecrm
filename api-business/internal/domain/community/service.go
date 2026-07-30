package community

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListCategories(ctx context.Context, onlyShow bool) ([]Category, error)
	ListTopics(ctx context.Context, onlyOn bool) ([]Topic, error)
	ListPosts(ctx context.Context, f ListFilter) ([]Post, int64, error)
	GetPost(ctx context.Context, id uint) (*Post, error)
	CreatePost(ctx context.Context, p *Post) error
	UpdatePost(ctx context.Context, p *Post) error
	SoftDeletePost(ctx context.Context, id uint) error
	IncPV(ctx context.Context, id uint) error
	IncReplyCount(ctx context.Context, id uint, delta int) error

	ListReplies(ctx context.Context, communityID uint, page, limit int) ([]Reply, int64, error)
	CreateReply(ctx context.Context, r *Reply) error
	SoftDeleteReply(ctx context.Context, id uint) error
	GetReply(ctx context.Context, id uint) (*Reply, error)

	LoadUserNickname(ctx context.Context, uid uint) (string, error)
	LoadTopicName(ctx context.Context, id uint) (string, error)
	LoadCateName(ctx context.Context, id uint) (string, error)
	LoadProductMeta(ctx context.Context, productID uint) (name string, price float64, merID uint, err error)
}

type ListFilter struct {
	MerID      *uint
	UID        *uint
	Status     *int8
	Keyword    string
	TopicID    uint
	OnlyPublic bool
	Page       int
	Limit      int
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.store.ListCategories(ctx, true)
}

func (s *Service) ListTopics(ctx context.Context) ([]Topic, error) {
	return s.store.ListTopics(ctx, true)
}

// ListHotTopics 热门话题（is_hot=1 优先；已由仓库按 is_hot DESC 排序）。
func (s *Service) ListHotTopics(ctx context.Context, limit int) ([]Topic, error) {
	list, err := s.store.ListTopics(ctx, true)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	hot := make([]Topic, 0, limit)
	for _, t := range list {
		if t.IsHot == 1 {
			hot = append(hot, t)
		}
		if len(hot) >= limit {
			break
		}
	}
	if len(hot) == 0 && len(list) > 0 {
		if len(list) < limit {
			limit = len(list)
		}
		return list[:limit], nil
	}
	return hot, nil
}

func (s *Service) ListApp(ctx context.Context, topicID uint, page, limit int) (*PageResult[Post], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListPosts(ctx, ListFilter{TopicID: topicID, OnlyPublic: true, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	_ = s.enrichPosts(ctx, list)
	return &PageResult[Post]{List: list, Total: total, Page: page, Limit: limit}, nil
}

// ListUser 返回当前用户自己的所有未删除帖子，包含待审和驳回状态，供 C 端“我的发布”查看。
func (s *Service) ListUser(ctx context.Context, uid uint, page, limit int) (*PageResult[Post], error) {
	if uid == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListPosts(ctx, ListFilter{UID: &uid, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	_ = s.enrichPosts(ctx, list)
	return &PageResult[Post]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListPlatform(ctx context.Context, status *int8, keyword string, page, limit int) (*PageResult[Post], error) {
	page, limit = normalize(page, limit)
	if status != nil && *status != StatusPending && *status != StatusApproved && *status != StatusRejected {
		return nil, ErrBadParam
	}
	list, total, err := s.store.ListPosts(ctx, ListFilter{
		Status: status, Keyword: strings.TrimSpace(keyword), Page: page, Limit: limit,
	})
	if err != nil {
		return nil, err
	}
	_ = s.enrichPosts(ctx, list)
	return &PageResult[Post]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, page, limit int) (*PageResult[Post], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListPosts(ctx, ListFilter{MerID: &merID, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	_ = s.enrichPosts(ctx, list)
	return &PageResult[Post]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint, bumpPV bool) (*Post, error) {
	p, err := s.store.GetPost(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []Post{*p}
	_ = s.enrichPosts(ctx, tmp)
	*p = tmp[0]
	if bumpPV {
		_ = s.store.IncPV(ctx, id)
		p.PV++
	}
	return p, nil
}

func (s *Service) CreatePost(ctx context.Context, uid uint, in CreatePostInput) (*Post, error) {
	if uid == 0 {
		return nil, ErrBadParam
	}
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || content == "" {
		return nil, ErrBadParam
	}
	cateID := in.CategoryID
	if cateID == 0 {
		cateID = 1
	}
	topicID := in.TopicID
	if topicID == 0 {
		topicID = 1
	}
	merID := in.MerID
	if in.ProductID > 0 {
		_, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
		if err != nil {
			return nil, ErrBadParam
		}
		if merID == 0 {
			merID = pMer
		}
	}
	p := &Post{
		Title: title, Content: content, Image: strings.TrimSpace(in.Image),
		CategoryID: cateID, TopicID: topicID, UID: uid, MerID: merID, ProductID: in.ProductID,
		Status: StatusPending, IsShow: 1, IsType: 1, CreateTime: time.Now(),
	}
	if err := s.store.CreatePost(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, p.CommunityID, false)
}

func (s *Service) Audit(ctx context.Context, id uint, in AuditInput) (*Post, error) {
	p, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	refusal := strings.TrimSpace(in.Refusal)
	if in.Status != StatusPending && in.Status != StatusApproved && in.Status != StatusRejected {
		return nil, ErrBadParam
	}
	if in.Status == StatusPending {
		if p.Status != StatusApproved || (in.IsShow == nil && in.IsHot == nil) {
			return nil, ErrBadParam
		}
	} else {
		if p.Status != StatusPending || (in.Status == StatusRejected && refusal == "") {
			return nil, ErrBadParam
		}
		p.Status = in.Status
		p.Refusal = refusal
		now := time.Now()
		p.StatusTime = &now
		if in.Status == StatusRejected {
			p.IsShow = 0
			p.IsHot = 0
		}
	}
	// 驳回帖无论请求体携带何种展示参数，均不可在 C 端显示或置顶。
	if in.Status != StatusRejected {
		if in.IsShow != nil {
			p.IsShow = int8(*in.IsShow)
		}
		if in.IsHot != nil {
			p.IsHot = int8(*in.IsHot)
		}
	}
	if err := s.store.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, false)
}

func (s *Service) DeletePost(ctx context.Context, id uint) error {
	if _, err := s.Get(ctx, id, false); err != nil {
		return err
	}
	return s.store.SoftDeletePost(ctx, id)
}

// CreateMerchantPost 商户发帖：强制 mer_id；商品须属本店；待平台审。
func (s *Service) CreateMerchantPost(ctx context.Context, merID, authorUID uint, in CreatePostInput) (*Post, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	if authorUID == 0 {
		authorUID = 1 // 演示默认挂 C 端 demo
	}
	in.MerID = merID
	if in.ProductID > 0 {
		_, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
		if err != nil {
			return nil, ErrBadParam
		}
		if pMer != merID {
			return nil, ErrForbidden
		}
	}
	return s.CreatePost(ctx, authorUID, in)
}

// UpdateMerchantPost 商户改帖：仅本店；改后重回待审。
func (s *Service) UpdateMerchantPost(ctx context.Context, merID, id uint, in CreatePostInput) (*Post, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	p, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if p.MerID != merID {
		return nil, ErrForbidden
	}
	title := strings.TrimSpace(in.Title)
	content := strings.TrimSpace(in.Content)
	if title == "" || content == "" {
		return nil, ErrBadParam
	}
	if in.ProductID > 0 {
		_, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
		if err != nil {
			return nil, ErrBadParam
		}
		if pMer != merID {
			return nil, ErrForbidden
		}
		p.ProductID = in.ProductID
	}
	p.Title = title
	p.Content = content
	if img := strings.TrimSpace(in.Image); img != "" {
		p.Image = img
	}
	if in.CategoryID > 0 {
		p.CategoryID = in.CategoryID
	}
	if in.TopicID > 0 {
		p.TopicID = in.TopicID
	}
	p.Status = StatusPending
	p.Refusal = ""
	p.IsShow = 1
	if err := s.store.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, false)
}

// DeleteMerchantPost 商户删帖：仅本店。
func (s *Service) DeleteMerchantPost(ctx context.Context, merID, id uint) error {
	if merID == 0 || id == 0 {
		return ErrBadParam
	}
	p, err := s.Get(ctx, id, false)
	if err != nil {
		return err
	}
	if p.MerID != merID {
		return ErrForbidden
	}
	return s.store.SoftDeletePost(ctx, id)
}

func (s *Service) ListReplies(ctx context.Context, communityID uint, page, limit int) (*PageResult[Reply], error) {
	if communityID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListReplies(ctx, communityID, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if nick, err := s.store.LoadUserNickname(ctx, list[i].UID); err == nil {
			list[i].Nickname = nick
		}
	}
	return &PageResult[Reply]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateReply(ctx context.Context, uid, communityID uint, in CreateReplyInput) (*Reply, error) {
	if uid == 0 || communityID == 0 {
		return nil, ErrBadParam
	}
	content := strings.TrimSpace(in.Content)
	if content == "" {
		return nil, ErrBadParam
	}
	p, err := s.Get(ctx, communityID, false)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusApproved || p.IsShow != 1 {
		return nil, ErrForbidden
	}
	r := &Reply{
		Content: content, PID: in.PID, UID: uid, CommunityID: communityID,
		Status: 1, CreateTime: time.Now(),
	}
	if err := s.store.CreateReply(ctx, r); err != nil {
		return nil, err
	}
	_ = s.store.IncReplyCount(ctx, communityID, 1)
	if nick, err := s.store.LoadUserNickname(ctx, uid); err == nil {
		r.Nickname = nick
	}
	return r, nil
}

func (s *Service) DeleteReply(ctx context.Context, id uint) error {
	r, err := s.store.GetReply(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := s.store.SoftDeleteReply(ctx, id); err != nil {
		return err
	}
	_ = s.store.IncReplyCount(ctx, r.CommunityID, -1)
	return nil
}

func (s *Service) enrichPosts(ctx context.Context, list []Post) error {
	for i := range list {
		if nick, err := s.store.LoadUserNickname(ctx, list[i].UID); err == nil {
			list[i].Nickname = nick
		}
		if name, err := s.store.LoadTopicName(ctx, list[i].TopicID); err == nil {
			list[i].TopicName = name
		}
		if name, err := s.store.LoadCateName(ctx, list[i].CategoryID); err == nil {
			list[i].CateName = name
		}
		if list[i].ProductID > 0 {
			if name, price, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID); err == nil {
				list[i].ProductName = name
				list[i].ProductPrice = price
			}
		}
	}
	return nil
}

func normalize(page, limit int) (int, int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
