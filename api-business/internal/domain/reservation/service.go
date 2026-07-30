package reservation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListProducts(ctx context.Context, merID *uint, page, limit int) ([]ProductView, int64, error)
	GetProduct(ctx context.Context, productID uint) (*ProductView, error)
	GetConfig(ctx context.Context, productID uint) (*Config, error)
	UpsertConfig(ctx context.Context, c *Config) error
	ListSlots(ctx context.Context, productID uint) ([]Slot, error)
	ReplaceSlots(ctx context.Context, productID uint, slots []Slot) error
	GetSlot(ctx context.Context, slotID uint) (*Slot, error)
	CountBooked(ctx context.Context, productID, slotID uint, date string) (int64, error)
	BumpSlotUse(ctx context.Context, slotID uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListApp(ctx context.Context, page, limit int) (*PageResult[ProductView], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListProducts(ctx, nil, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[ProductView]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, page, limit int) (*PageResult[ProductView], error) {
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListProducts(ctx, &merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[ProductView]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetProduct(ctx context.Context, productID uint) (*ProductView, error) {
	p, err := s.store.GetProduct(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) GetConfig(ctx context.Context, merID, productID uint) (*Config, []Slot, error) {
	p, err := s.GetProduct(ctx, productID)
	if err != nil {
		return nil, nil, err
	}
	if merID > 0 && p.MerID != merID {
		return nil, nil, ErrForbidden
	}
	cfg, err := s.store.GetConfig(ctx, productID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}
	slots, err := s.store.ListSlots(ctx, productID)
	if err != nil {
		return nil, nil, err
	}
	return cfg, slots, nil
}

func (s *Service) SaveConfig(ctx context.Context, merID uint, in ConfigSaveInput) error {
	if in.ProductID == 0 || len(in.Slots) == 0 {
		return ErrBadParam
	}
	p, err := s.GetProduct(ctx, in.ProductID)
	if err != nil {
		return err
	}
	if p.MerID != merID {
		return ErrForbidden
	}
	days := in.ShowReservationDays
	if days <= 0 {
		days = 7
	}
	rt := in.ReservationType
	if rt == 0 {
		rt = 1
	}
	period, _ := json.Marshal(in.Slots)
	cfg := &Config{
		ProductID:           in.ProductID,
		ReservationType:     rt,
		ShowReservationDays: days,
		IsCancelReservation: 1,
		TimePeriod:          string(period),
	}
	if err := s.store.UpsertConfig(ctx, cfg); err != nil {
		return err
	}
	clean := make([]Slot, 0, len(in.Slots))
	for i, sl := range in.Slots {
		st := strings.TrimSpace(sl.StartTime)
		et := strings.TrimSpace(sl.EndTime)
		if st == "" || et == "" || sl.Stock < 0 {
			return ErrBadParam
		}
		u := strings.TrimSpace(sl.Unique)
		if u == "" {
			u = fmt.Sprintf("rsv%08d", i+1)
		}
		clean = append(clean, Slot{
			ProductID: in.ProductID,
			Unique:    u,
			StartTime: st,
			EndTime:   et,
			Stock:     sl.Stock,
		})
	}
	return s.store.ReplaceSlots(ctx, in.ProductID, clean)
}

func (s *Service) DaySlots(ctx context.Context, productID uint, date string) ([]SlotDayView, error) {
	if err := validateDate(date); err != nil {
		return nil, err
	}
	p, err := s.GetProduct(ctx, productID)
	if err != nil {
		return nil, err
	}
	cfg, err := s.store.GetConfig(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	days := cfg.ShowReservationDays
	if days <= 0 {
		days = 7
	}
	d, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	today := time.Now().Truncate(24 * time.Hour)
	if d.Before(today) || d.After(today.AddDate(0, 0, days-1)) {
		return nil, ErrBadDate
	}
	slots, err := s.store.ListSlots(ctx, productID)
	if err != nil {
		return nil, err
	}
	out := make([]SlotDayView, 0, len(slots))
	for _, sl := range slots {
		booked, err := s.store.CountBooked(ctx, productID, sl.AttrReservationID, date)
		if err != nil {
			return nil, err
		}
		remain := sl.Stock - int(booked)
		if remain < 0 {
			remain = 0
		}
		out = append(out, SlotDayView{
			AttrReservationID: sl.AttrReservationID,
			StartTime:         sl.StartTime,
			EndTime:           sl.EndTime,
			Stock:             sl.Stock,
			Booked:            booked,
			Remain:            remain,
			Label:             sl.StartTime + "-" + sl.EndTime,
		})
	}
	_ = p
	return out, nil
}

// ValidateBook 下单前校验时段余量。
func (s *Service) ValidateBook(ctx context.Context, productID, slotID uint, date string) (*ProductView, *Slot, error) {
	if err := validateDate(date); err != nil {
		return nil, nil, err
	}
	p, err := s.GetProduct(ctx, productID)
	if err != nil {
		return nil, nil, err
	}
	slot, err := s.store.GetSlot(ctx, slotID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrNoSlot
		}
		return nil, nil, err
	}
	if slot.ProductID != productID {
		return nil, nil, ErrNoSlot
	}
	views, err := s.DaySlots(ctx, productID, date)
	if err != nil {
		return nil, nil, err
	}
	for _, v := range views {
		if v.AttrReservationID == slotID {
			if v.Remain <= 0 {
				return nil, nil, ErrFull
			}
			return p, slot, nil
		}
	}
	return nil, nil, ErrNoSlot
}

func (s *Service) AfterBooked(ctx context.Context, slotID uint) error {
	return s.store.BumpSlotUse(ctx, slotID)
}

func validateDate(date string) error {
	date = strings.TrimSpace(date)
	if date == "" {
		return ErrBadDate
	}
	if _, err := time.ParseInLocation("2006-01-02", date, time.Local); err != nil {
		return ErrBadDate
	}
	return nil
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
