package seckill

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListTimes(ctx context.Context) ([]TimeSlot, error)
	ListActives(ctx context.Context, merID *uint, onlyOn bool, page, limit int) ([]Active, int64, error)
	GetActive(ctx context.Context, id uint) (*Active, error)
	GetActiveByProduct(ctx context.Context, productID uint) (*Active, error)
	CreateActive(ctx context.Context, a *Active) error
	UpdateActive(ctx context.Context, a *Active) error
	SoftDeleteActive(ctx context.Context, id uint) error
	LoadProductMeta(ctx context.Context, productID uint) (storeName, image, merName string, price float64, merID uint, err error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListTimes(ctx context.Context) ([]TimeSlot, error) {
	return s.store.ListTimes(ctx)
}

func (s *Service) ListAdmin(ctx context.Context, merID *uint, page, limit int) (*PageResult[Active], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListActives(ctx, merID, false, page, limit)
	if err != nil {
		return nil, err
	}
	_ = s.enrich(ctx, list)
	return &PageResult[Active]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[Active], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListActives(ctx, nil, true, page, limit)
	if err != nil {
		return nil, err
	}
	slots, _ := s.store.ListTimes(ctx)
	now := time.Now()
	for i := range list {
		list[i].InWindow = inWindow(&list[i], slots, now)
	}
	_ = s.enrich(ctx, list)
	return &PageResult[Active]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Get(ctx context.Context, id uint) (*Active, error) {
	a, err := s.store.GetActive(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	tmp := []Active{*a}
	_ = s.enrich(ctx, tmp)
	*a = tmp[0]
	slots, _ := s.store.ListTimes(ctx)
	a.InWindow = inWindow(a, slots, time.Now())
	return a, nil
}

func (s *Service) Create(ctx context.Context, merID uint, in ActiveInput) (*Active, error) {
	if merID == 0 || in.ProductID == 0 || strings.TrimSpace(in.Name) == "" || in.SeckillPrice <= 0 {
		return nil, ErrBadParam
	}
	_, _, _, _, pMer, err := s.store.LoadProductMeta(ctx, in.ProductID)
	if err != nil {
		return nil, ErrBadParam
	}
	if pMer != merID {
		return nil, ErrBadParam
	}
	now := time.Now().Unix()
	a := &Active{
		Name: strings.TrimSpace(in.Name), SeckillTimeIDs: defaultTimes(in.SeckillTimeIDs),
		StartDay: in.StartDay, EndDay: in.EndDay, MerID: merID, ProductID: in.ProductID,
		SeckillPrice: in.SeckillPrice, OncePayCount: in.OncePayCount, ActiveStatus: 1, Status: 1,
		CreateTime: now, UpdateTime: now,
	}
	if a.OncePayCount <= 0 {
		a.OncePayCount = 1
	}
	if a.StartDay == "" {
		a.StartDay = time.Now().Format("2006-01-02")
	}
	if a.EndDay == "" {
		a.EndDay = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}
	if in.Status != nil {
		a.Status = *in.Status
	}
	if !validActivityDates(a.StartDay, a.EndDay) {
		return nil, ErrBadParam
	}
	if err := s.store.CreateActive(ctx, a); err != nil {
		return nil, err
	}
	return s.Get(ctx, a.SeckillActiveID)
}

func (s *Service) Update(ctx context.Context, merID, id uint, in ActiveInput) (*Active, error) {
	if in.Status != nil && *in.Status != 0 && *in.Status != 1 {
		return nil, ErrBadParam
	}
	a, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if merID > 0 && a.MerID != merID {
		return nil, ErrNotFound
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		a.Name = name
	}
	if in.SeckillTimeIDs != "" {
		a.SeckillTimeIDs = in.SeckillTimeIDs
	}
	if in.StartDay != "" {
		a.StartDay = in.StartDay
	}
	if in.EndDay != "" {
		a.EndDay = in.EndDay
	}
	if in.SeckillPrice > 0 {
		a.SeckillPrice = in.SeckillPrice
	}
	if in.OncePayCount > 0 {
		a.OncePayCount = in.OncePayCount
	}
	if in.Status != nil {
		a.Status = *in.Status
	}
	if !validActivityDates(a.StartDay, a.EndDay) {
		return nil, ErrBadParam
	}
	a.UpdateTime = time.Now().Unix()
	if err := s.store.UpdateActive(ctx, a); err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, merID, id uint) error {
	a, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if merID > 0 && a.MerID != merID {
		return ErrNotFound
	}
	return s.store.SoftDeleteActive(ctx, id)
}

// SetStatus 启停秒杀活动（status 1开启 / 0关闭）。
func (s *Service) SetStatus(ctx context.Context, merID, id uint, status int8) (*Active, error) {
	if status != 0 && status != 1 {
		return nil, ErrBadParam
	}
	st := status
	return s.Update(ctx, merID, id, ActiveInput{Status: &st})
}

// QuotePrice 若商品当前在秒杀场次内，返回秒杀价与限购。
func (s *Service) QuotePrice(ctx context.Context, productID uint) (price float64, activeID uint, oncePay int, ok bool, err error) {
	a, err := s.store.GetActiveByProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, 0, 0, false, nil
		}
		return 0, 0, 0, false, err
	}
	slots, err := s.store.ListTimes(ctx)
	if err != nil {
		return 0, 0, 0, false, err
	}
	if !inWindow(a, slots, time.Now()) {
		return 0, 0, 0, false, nil
	}
	once := a.OncePayCount
	if once <= 0 {
		once = 1
	}
	return a.SeckillPrice, a.SeckillActiveID, once, true, nil
}

func (s *Service) enrich(ctx context.Context, list []Active) error {
	for i := range list {
		name, img, merName, price, _, err := s.store.LoadProductMeta(ctx, list[i].ProductID)
		if err != nil {
			continue
		}
		list[i].StoreName = name
		list[i].Image = img
		list[i].MerName = merName
		list[i].Price = price
	}
	return nil
}

func inWindow(a *Active, slots []TimeSlot, now time.Time) bool {
	if a == nil || a.Status != 1 || a.ActiveStatus != 1 || a.DeleteTime != nil {
		return false
	}
	day := now.Format("2006-01-02")
	if day < a.StartDay || day > a.EndDay {
		return false
	}
	hour := now.Hour()
	ids := map[uint]struct{}{}
	for _, p := range strings.Split(a.SeckillTimeIDs, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, _ := strconv.ParseUint(p, 10, 64)
		if n > 0 {
			ids[uint(n)] = struct{}{}
		}
	}
	if len(ids) == 0 {
		return true
	}
	for _, sl := range slots {
		if sl.Status != 1 {
			continue
		}
		if _, ok := ids[sl.SeckillTimeID]; !ok {
			continue
		}
		if hour >= sl.StartTime && hour < sl.EndTime {
			return true
		}
	}
	return false
}

func defaultTimes(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "1"
	}
	return s
}

func validActivityDates(startDay, endDay string) bool {
	start, err := time.Parse("2006-01-02", startDay)
	if err != nil {
		return false
	}
	end, err := time.Parse("2006-01-02", endDay)
	return err == nil && !end.Before(start)
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
