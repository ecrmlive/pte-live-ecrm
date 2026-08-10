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
	GetCategory(ctx context.Context, id uint) (*Category, error)
	CreateCategory(ctx context.Context, row *Category) error
	UpdateCategory(ctx context.Context, row *Category) error
	UpdateCategoryShow(ctx context.Context, id uint, isShow int8) error
	DeleteCategory(ctx context.Context, id uint) error
	CountCategoryUsage(ctx context.Context, id uint) (posts int64, topics int64, err error)
	ListTopics(ctx context.Context, onlyOn bool) ([]Topic, error)
	GetTopic(ctx context.Context, id uint) (*Topic, error)
	CreateTopic(ctx context.Context, row *Topic) error
	UpdateTopic(ctx context.Context, row *Topic) error
	UpdateTopicStatus(ctx context.Context, id uint, status int8) error
	UpdateTopicHot(ctx context.Context, id uint, isHot int8) error
	SoftDeleteTopic(ctx context.Context, id uint) error
	TopicNameExists(ctx context.Context, name string, excludeID uint) (bool, error)
	ListPosts(ctx context.Context, f ListFilter) ([]Post, int64, error)
	CountPosts(ctx context.Context, f ListFilter) (int64, error)
	GetPost(ctx context.Context, id uint) (*Post, error)
	CreatePost(ctx context.Context, p *Post) error
	UpdatePost(ctx context.Context, p *Post) error
	SoftDeletePost(ctx context.Context, id uint) error
	IncPV(ctx context.Context, id uint) error
	IncReplyCount(ctx context.Context, id uint, delta int) error

	ListReplies(ctx context.Context, communityID uint, page, limit int) ([]Reply, int64, error)
	ListAllReplies(ctx context.Context, f ReplyListFilter) ([]Reply, int64, error)
	CreateReply(ctx context.Context, r *Reply) error
	UpdateReply(ctx context.Context, r *Reply) error
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
	CategoryID uint
	IsShow     *int8
	IsType     *int8
	AuthorType string
	AuthorKW   string
	TitleOnly  bool
	OnlyPublic bool
	Page       int
	Limit      int
}

type ReplyListFilter struct {
	Keyword  string
	Username string
	DateFrom string
	DateTo   string
	Page     int
	Limit    int
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

// ListCategories 平台后台：返回全部社区分类（含隐藏）。
func (s *Service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.store.ListCategories(ctx, false)
}

func (s *Service) CreateCategory(ctx context.Context, in CategoryInput) (*Category, error) {
	name := strings.TrimSpace(in.CateName)
	if name == "" {
		return nil, ErrBadParam
	}
	row := &Category{
		CateName: name,
		PID:      0,
		IsShow:   1,
		Sort:     in.Sort,
	}
	if in.PID != nil {
		row.PID = *in.PID
	}
	if in.IsShow != nil {
		if *in.IsShow != 0 && *in.IsShow != 1 {
			return nil, ErrBadParam
		}
		row.IsShow = *in.IsShow
	}
	if err := s.store.CreateCategory(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uint, in CategoryInput) (*Category, error) {
	if id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.CateName); name != "" {
		row.CateName = name
	}
	row.Sort = in.Sort
	if in.PID != nil {
		row.PID = *in.PID
	}
	if in.IsShow != nil {
		if *in.IsShow != 0 && *in.IsShow != 1 {
			return nil, ErrBadParam
		}
		row.IsShow = *in.IsShow
	}
	if err := s.store.UpdateCategory(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) SetCategoryShow(ctx context.Context, id uint, isShow int8) error {
	if id == 0 || (isShow != 0 && isShow != 1) {
		return ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.UpdateCategoryShow(ctx, id, isShow)
}

func (s *Service) DeleteCategory(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	posts, topics, err := s.store.CountCategoryUsage(ctx, id)
	if err != nil {
		return err
	}
	if posts > 0 || topics > 0 {
		return ErrForbidden
	}
	return s.store.DeleteCategory(ctx, id)
}

// ListTopics 平台后台：返回全部未删除话题（含隐藏），并填充上级分类名。
func (s *Service) ListTopics(ctx context.Context) ([]Topic, error) {
	list, err := s.store.ListTopics(ctx, false)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.LoadCateName(ctx, list[i].CategoryID); err == nil {
			list[i].CateName = name
		}
	}
	return list, nil
}

func (s *Service) CreateTopic(ctx context.Context, in TopicInput) (*Topic, error) {
	name := strings.TrimSpace(in.TopicName)
	if name == "" || in.CategoryID == 0 {
		return nil, ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, in.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadParam
		}
		return nil, err
	}
	exists, err := s.store.TopicNameExists(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicate
	}
	row := &Topic{
		TopicName:  name,
		Pic:        strings.TrimSpace(in.Pic),
		CategoryID: in.CategoryID,
		Sort:       in.Sort,
		Status:     1,
		IsHot:      0,
		CreateTime: time.Now(),
	}
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return nil, ErrBadParam
		}
		row.Status = *in.Status
	}
	if in.IsHot != nil {
		if *in.IsHot != 0 && *in.IsHot != 1 {
			return nil, ErrBadParam
		}
		row.IsHot = *in.IsHot
	}
	if err := s.store.CreateTopic(ctx, row); err != nil {
		return nil, err
	}
	if name, err := s.store.LoadCateName(ctx, row.CategoryID); err == nil {
		row.CateName = name
	}
	return row, nil
}

func (s *Service) UpdateTopic(ctx context.Context, id uint, in TopicInput) (*Topic, error) {
	if id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetTopic(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	name := strings.TrimSpace(in.TopicName)
	if name == "" {
		return nil, ErrBadParam
	}
	if in.CategoryID == 0 {
		return nil, ErrBadParam
	}
	if _, err := s.store.GetCategory(ctx, in.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBadParam
		}
		return nil, err
	}
	exists, err := s.store.TopicNameExists(ctx, name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicate
	}
	row.TopicName = name
	row.Pic = strings.TrimSpace(in.Pic)
	row.CategoryID = in.CategoryID
	row.Sort = in.Sort
	if in.Status != nil {
		if *in.Status != 0 && *in.Status != 1 {
			return nil, ErrBadParam
		}
		row.Status = *in.Status
	}
	if in.IsHot != nil {
		if *in.IsHot != 0 && *in.IsHot != 1 {
			return nil, ErrBadParam
		}
		row.IsHot = *in.IsHot
	}
	if err := s.store.UpdateTopic(ctx, row); err != nil {
		return nil, err
	}
	if name, err := s.store.LoadCateName(ctx, row.CategoryID); err == nil {
		row.CateName = name
	}
	return row, nil
}

func (s *Service) SetTopicStatus(ctx context.Context, id uint, status int8) error {
	if id == 0 || (status != 0 && status != 1) {
		return ErrBadParam
	}
	if _, err := s.store.GetTopic(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.UpdateTopicStatus(ctx, id, status)
}

func (s *Service) SetTopicHot(ctx context.Context, id uint, isHot int8) error {
	if id == 0 || (isHot != 0 && isHot != 1) {
		return ErrBadParam
	}
	if _, err := s.store.GetTopic(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.UpdateTopicHot(ctx, id, isHot)
}

func (s *Service) DeleteTopic(ctx context.Context, id uint) error {
	if id == 0 {
		return ErrBadParam
	}
	if _, err := s.store.GetTopic(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteTopic(ctx, id)
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

func (s *Service) ListPlatform(ctx context.Context, f ListFilter) (*PlatformPostPage, error) {
	f.Page, f.Limit = normalize(f.Page, f.Limit)
	f.Keyword = strings.TrimSpace(f.Keyword)
	f.AuthorKW = strings.TrimSpace(f.AuthorKW)
	f.AuthorType = strings.TrimSpace(f.AuthorType)
	f.TitleOnly = true
	if f.Status != nil &&
		*f.Status != StatusPending && *f.Status != StatusApproved &&
		*f.Status != StatusRejected && *f.Status != StatusForceOff {
		return nil, ErrBadParam
	}
	if f.IsShow != nil && *f.IsShow != 0 && *f.IsShow != 1 {
		return nil, ErrBadParam
	}
	if f.IsType != nil && *f.IsType != TypeImage && *f.IsType != TypeVideo {
		return nil, ErrBadParam
	}
	list, total, err := s.store.ListPosts(ctx, f)
	if err != nil {
		return nil, err
	}
	_ = s.enrichPosts(ctx, list)

	countBase := f
	countBase.Page, countBase.Limit = 1, 1
	countBase.IsType = nil
	imgType := TypeImage
	vidType := TypeVideo
	imgFilter := countBase
	imgFilter.IsType = &imgType
	vidFilter := countBase
	vidFilter.IsType = &vidType
	imageCount, err := s.store.CountPosts(ctx, imgFilter)
	if err != nil {
		return nil, err
	}
	videoCount, err := s.store.CountPosts(ctx, vidFilter)
	if err != nil {
		return nil, err
	}
	return &PlatformPostPage{
		List: list, Total: total, Page: f.Page, Limit: f.Limit,
		ImageCount: imageCount, VideoCount: videoCount,
	}, nil
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
		Status: StatusPending, IsShow: 1, Start: 1, IsType: TypeImage, CreateTime: time.Now(),
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
	if in.Status != StatusPending && in.Status != StatusApproved &&
		in.Status != StatusRejected && in.Status != StatusForceOff {
		return nil, ErrBadParam
	}
	if in.Status == StatusPending {
		if p.Status != StatusApproved || (in.IsShow == nil && in.IsHot == nil) {
			return nil, ErrBadParam
		}
	} else if in.Status == StatusForceOff {
		if p.Status != StatusApproved {
			return nil, ErrBadParam
		}
		if refusal == "" {
			return nil, ErrBadParam
		}
		p.Status = StatusForceOff
		p.Refusal = refusal
		p.IsShow = 0
		p.IsHot = 0
		now := time.Now()
		p.StatusTime = &now
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
	// 驳回/强制下架帖不可在 C 端显示或置顶。
	if in.Status != StatusRejected && in.Status != StatusForceOff {
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

func (s *Service) UpdateStar(ctx context.Context, id uint, start int8) (*Post, error) {
	if start < 1 || start > 5 {
		return nil, ErrBadParam
	}
	p, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	p.Start = start
	if err := s.store.UpdatePost(ctx, p); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, false)
}

func (s *Service) SwitchShow(ctx context.Context, id uint, isShow int8) (*Post, error) {
	if isShow != 0 && isShow != 1 {
		return nil, ErrBadParam
	}
	p, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if p.Status != StatusApproved {
		return nil, ErrBadParam
	}
	p.IsShow = isShow
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

func (s *Service) ListAllReplies(ctx context.Context, f ReplyListFilter) (*PageResult[Reply], error) {
	f.Page, f.Limit = normalize(f.Page, f.Limit)
	f.Keyword = strings.TrimSpace(f.Keyword)
	f.Username = strings.TrimSpace(f.Username)
	f.DateFrom = strings.TrimSpace(f.DateFrom)
	f.DateTo = strings.TrimSpace(f.DateTo)
	list, total, err := s.store.ListAllReplies(ctx, f)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if nick, err := s.store.LoadUserNickname(ctx, list[i].UID); err == nil {
			list[i].Nickname = nick
		}
		if post, err := s.store.GetPost(ctx, list[i].CommunityID); err == nil {
			list[i].PostTitle = post.Title
		}
	}
	return &PageResult[Reply]{List: list, Total: total, Page: f.Page, Limit: f.Limit}, nil
}

func (s *Service) AuditReply(ctx context.Context, id uint, in ReplyAuditInput) (*Reply, error) {
	if in.Status != StatusApproved && in.Status != StatusRejected {
		return nil, ErrBadParam
	}
	refusal := strings.TrimSpace(in.Refusal)
	if in.Status == StatusRejected && refusal == "" {
		return nil, ErrBadParam
	}
	r, err := s.store.GetReply(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if r.Status != StatusPending {
		return nil, ErrBadParam
	}
	r.Status = in.Status
	r.Refusal = refusal
	if err := s.store.UpdateReply(ctx, r); err != nil {
		return nil, err
	}
	if nick, err := s.store.LoadUserNickname(ctx, r.UID); err == nil {
		r.Nickname = nick
	}
	if post, err := s.store.GetPost(ctx, r.CommunityID); err == nil {
		r.PostTitle = post.Title
	}
	return r, nil
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
		if list[i].Start < 1 {
			list[i].Start = 1
		}
		if list[i].IsType == 0 {
			list[i].IsType = TypeImage
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
