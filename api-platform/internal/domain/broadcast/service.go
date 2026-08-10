package broadcast

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	List(ctx context.Context, filter ListFilter, page, limit int) ([]Room, int64, error)
	Get(ctx context.Context, id uint) (*Room, error)
	Create(ctx context.Context, r *Room) error
	Update(ctx context.Context, r *Room) error
	SoftDelete(ctx context.Context, id uint) error
	ReplaceGoods(ctx context.Context, roomID uint, productIDs []uint) error
	ListGoods(ctx context.Context, roomID uint) ([]RoomGoods, error)
	LoadMerName(ctx context.Context, merID uint) (string, error)
	LoadProductMeta(ctx context.Context, productID uint) (storeName, image string, price float64, merID uint, err error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[Room], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, ListFilter{OnlyPublic: true}, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list, false)
	return &PageResult[Room]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, page, limit int) (*PageResult[Room], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.List(ctx, ListFilter{MerID: &merID}, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list, false)
	return &PageResult[Room]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListPlatform(ctx context.Context, filter ListFilter, page, limit int) (*PageResult[Room], error) {
	page, limit = normalize(page, limit)
	// 店铺类别无匹配商户时直接空页，避免 IN () SQL 错误。
	if filter.MerIDs != nil && len(filter.MerIDs) == 0 {
		return &PageResult[Room]{List: []Room{}, Total: 0, Page: page, Limit: limit}, nil
	}
	list, total, err := s.store.List(ctx, filter, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list, false)
	return &PageResult[Room]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint, withGoods bool) (*Room, error) {
	r, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []Room{*r}
	_ = s.enrich(ctx, tmp, withGoods)
	*r = tmp[0]
	return r, nil
}

func (s *Service) GetApp(ctx context.Context, id uint) (*Room, error) {
	r, err := s.Get(ctx, id, true)
	if err != nil {
		return nil, err
	}
	if r.IsShow != 1 || r.Status != AuditApproved {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in SaveInput) (*Room, error) {
	if merID == 0 || strings.TrimSpace(in.Name) == "" {
		return nil, ErrBadParam
	}
	if err := s.assertProducts(ctx, merID, in.ProductIDs); err != nil {
		return nil, err
	}
	start, end := parseRange(in.StartTime, in.EndTime)
	live := LiveStatusPending
	if in.LiveStatus != nil {
		live = *in.LiveStatus
	}
	r := &Room{
		MerID: merID, Name: strings.TrimSpace(in.Name),
		CoverImg: strings.TrimSpace(in.CoverImg), FeedsImg: strings.TrimSpace(in.FeedsImg),
		PlayURL: strings.TrimSpace(in.PlayURL), PushURL: strings.TrimSpace(in.PushURL),
		AnchorName: strings.TrimSpace(in.AnchorName), AnchorWechat: strings.TrimSpace(in.AnchorWechat),
		Phone: strings.TrimSpace(in.Phone),
		StartTime: start, EndTime: end, LiveStatus: live, Status: AuditPending,
		IsShow: 1, Star: 1, Mark: strings.TrimSpace(in.Mark), CreateTime: time.Now(),
	}
	if in.IsShow != nil {
		r.IsShow = int8(*in.IsShow)
	}
	if in.Sort != nil {
		r.Sort = *in.Sort
	}
	if in.Star != nil {
		r.Star = *in.Star
	}
	if err := s.store.Create(ctx, r); err != nil {
		return nil, err
	}
	if err := s.store.ReplaceGoods(ctx, r.BroadcastRoomID, in.ProductIDs); err != nil {
		return nil, err
	}
	return s.Get(ctx, r.BroadcastRoomID, true)
}

func (s *Service) Update(ctx context.Context, merID, id uint, in SaveInput) (*Room, error) {
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if merID > 0 && r.MerID != merID {
		return nil, ErrForbidden
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		r.Name = name
	}
	if in.CoverImg != "" {
		r.CoverImg = strings.TrimSpace(in.CoverImg)
	}
	if in.FeedsImg != "" {
		r.FeedsImg = strings.TrimSpace(in.FeedsImg)
	}
	if in.PlayURL != "" {
		r.PlayURL = strings.TrimSpace(in.PlayURL)
	}
	if in.PushURL != "" {
		r.PushURL = strings.TrimSpace(in.PushURL)
	}
	if in.AnchorName != "" {
		r.AnchorName = strings.TrimSpace(in.AnchorName)
	}
	if in.AnchorWechat != "" {
		r.AnchorWechat = strings.TrimSpace(in.AnchorWechat)
	}
	if in.Phone != "" {
		r.Phone = strings.TrimSpace(in.Phone)
	}
	if in.Mark != "" {
		r.Mark = strings.TrimSpace(in.Mark)
	}
	if in.StartTime != "" || in.EndTime != "" {
		start, end := parseRange(in.StartTime, in.EndTime)
		r.StartTime, r.EndTime = start, end
	}
	if in.LiveStatus != nil {
		r.LiveStatus = *in.LiveStatus
	}
	if in.IsShow != nil {
		r.IsShow = int8(*in.IsShow)
	}
	if in.Sort != nil {
		r.Sort = *in.Sort
	}
	if in.Star != nil {
		r.Star = *in.Star
	}
	if merID > 0 {
		r.Status = AuditPending
		r.Refusal = ""
	}
	if in.ProductIDs != nil {
		owner := merID
		if owner == 0 {
			owner = r.MerID
		}
		if err := s.assertProducts(ctx, owner, in.ProductIDs); err != nil {
			return nil, err
		}
		if err := s.store.ReplaceGoods(ctx, id, in.ProductIDs); err != nil {
			return nil, err
		}
	}
	if err := s.store.Update(ctx, r); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, true)
}

func (s *Service) SetGoods(ctx context.Context, merID, id uint, productIDs []uint) (*Room, error) {
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if merID > 0 && r.MerID != merID {
		return nil, ErrForbidden
	}
	owner := merID
	if owner == 0 {
		owner = r.MerID
	}
	if err := s.assertProducts(ctx, owner, productIDs); err != nil {
		return nil, err
	}
	if err := s.store.ReplaceGoods(ctx, id, productIDs); err != nil {
		return nil, err
	}
	if merID > 0 {
		r.Status = AuditPending
		r.Refusal = ""
		if err := s.store.Update(ctx, r); err != nil {
			return nil, err
		}
	}
	return s.Get(ctx, id, true)
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return err
	}
	if merID > 0 && r.MerID != merID {
		return ErrForbidden
	}
	return s.store.SoftDelete(ctx, id)
}

func (s *Service) Audit(ctx context.Context, id uint, in AuditInput) (*Room, error) {
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if in.Status != 0 && in.Status != AuditApproved && in.Status != AuditRejected {
		return nil, ErrBadParam
	}
	refusal := strings.TrimSpace(in.Refusal)
	if in.Status == AuditRejected && refusal == "" {
		return nil, ErrBadParam
	}
	// 审核结论只能从待审状态产生；通过后的可见性调整仍可单独执行。
	if in.Status != 0 && r.Status != AuditPending {
		return nil, ErrBadParam
	}
	if in.Status != 0 {
		r.Status = in.Status
	}
	if in.Refusal != "" || in.Status == AuditRejected {
		r.Refusal = refusal
	}
	if in.IsShow != nil {
		r.IsShow = int8(*in.IsShow)
	}
	if in.Status == 0 && in.IsShow == nil {
		return nil, ErrBadParam
	}
	if err := s.store.Update(ctx, r); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, true)
}

func (s *Service) SetShow(ctx context.Context, id uint, isShow int8) (*Room, error) {
	if isShow != 0 && isShow != 1 {
		return nil, ErrBadParam
	}
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	r.IsShow = isShow
	if err := s.store.Update(ctx, r); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, true)
}

func (s *Service) SetRecommend(ctx context.Context, id uint, in RecommendInput) (*Room, error) {
	r, err := s.Get(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if in.Star != nil {
		if *in.Star < 0 || *in.Star > 5 {
			return nil, ErrBadParam
		}
		r.Star = *in.Star
	}
	if in.Sort != nil {
		r.Sort = *in.Sort
	}
	if in.Star == nil && in.Sort == nil {
		return nil, ErrBadParam
	}
	if err := s.store.Update(ctx, r); err != nil {
		return nil, err
	}
	return s.Get(ctx, id, true)
}

func (s *Service) assertProducts(ctx context.Context, merID uint, ids []uint) error {
	for _, pid := range ids {
		if pid == 0 {
			continue
		}
		_, _, _, pMer, err := s.store.LoadProductMeta(ctx, pid)
		if err != nil || pMer != merID {
			return ErrBadParam
		}
	}
	return nil
}

func (s *Service) enrich(ctx context.Context, list []Room, withGoods bool) error {
	for i := range list {
		if name, err := s.store.LoadMerName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
		if !withGoods {
			continue
		}
		goods, err := s.store.ListGoods(ctx, list[i].BroadcastRoomID)
		if err != nil {
			continue
		}
		for j := range goods {
			sn, img, price, _, err := s.store.LoadProductMeta(ctx, goods[j].ProductID)
			if err != nil {
				continue
			}
			goods[j].StoreName = sn
			goods[j].Image = img
			goods[j].Price = price
		}
		list[i].Goods = goods
	}
	return nil
}

func parseRange(startS, endS string) (*time.Time, *time.Time) {
	now := time.Now()
	start := now
	end := now.Add(2 * time.Hour)
	if t, ok := parseTime(startS); ok {
		start = t
	}
	if t, ok := parseTime(endS); ok {
		end = t
	}
	if !end.After(start) {
		end = start.Add(2 * time.Hour)
	}
	return &start, &end
}

func parseTime(s string) (time.Time, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
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
