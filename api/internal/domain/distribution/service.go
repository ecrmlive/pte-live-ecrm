package distribution

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserSpread struct {
	UID        uint  `gorm:"column:uid"`
	SpreadUID  uint  `gorm:"column:spread_uid"`
	IsPromoter uint8 `gorm:"column:is_promoter"`
	Status     int8  `gorm:"column:status"`
}

type Store interface {
	WithTx(ctx context.Context, fn func(ctx context.Context, tx Store) error) error
	GetUser(ctx context.Context, uid uint) (*UserSpread, error)
	SetSpreadUID(ctx context.Context, uid, spreadUID uint) (bool, error)
	CreateLog(ctx context.Context, log *SpreadLog) error
	CountChildren(ctx context.Context, spreadUID uint) (int64, error)
	ListLogs(ctx context.Context, page, limit int) ([]SpreadLog, int64, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Bind(ctx context.Context, uid uint, in BindInput) (*MeInfo, error) {
	if uid == 0 || in.SpreadUID == 0 {
		return nil, ErrBadParam
	}
	if uid == in.SpreadUID {
		return nil, ErrSelfBind
	}
	err := s.store.WithTx(ctx, func(ctx context.Context, tx Store) error {
		me, err := tx.GetUser(ctx, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrUserNotFound
			}
			return err
		}
		if me.SpreadUID == in.SpreadUID {
			// 幂等：已绑同一人
			return nil
		}
		if me.SpreadUID > 0 {
			return ErrAlreadyBound
		}
		sp, err := tx.GetUser(ctx, in.SpreadUID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrSpreadInvalid
			}
			return err
		}
		if sp.Status != 1 {
			return ErrSpreadInvalid
		}
		ok, err := tx.SetSpreadUID(ctx, uid, in.SpreadUID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAlreadyBound
		}
		return tx.CreateLog(ctx, &SpreadLog{
			UID:          uid,
			OldSpreadUID: 0,
			SpreadUID:    in.SpreadUID,
		})
	})
	if err != nil {
		return nil, err
	}
	return s.Me(ctx, uid)
}

func (s *Service) Me(ctx context.Context, uid uint) (*MeInfo, error) {
	u, err := s.store.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	n, err := s.store.CountChildren(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &MeInfo{
		UID:         u.UID,
		SpreadUID:   u.SpreadUID,
		IsPromoter:  u.IsPromoter,
		SpreadCount: n,
	}, nil
}

func (s *Service) ListLogs(ctx context.Context, page, limit int) (*PageResult[SpreadLog], error) {
	list, total, err := s.store.ListLogs(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return &PageResult[SpreadLog]{List: list, Total: total, Page: page, Limit: limit}, nil
}
