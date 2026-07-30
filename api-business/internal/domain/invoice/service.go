package invoice

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderMeta struct {
	OrderID  uint
	UID      uint
	MerID    uint
	Paid     bool
	Status   int8
}

type Store interface {
	ListByUID(ctx context.Context, uid uint, page, limit int) ([]Invoice, int64, error)
	ListByMer(ctx context.Context, merID uint, page, limit int) ([]Invoice, int64, error)
	Get(ctx context.Context, id uint) (*Invoice, error)
	FindByOrder(ctx context.Context, orderID uint) (*Invoice, error)
	Create(ctx context.Context, row *Invoice) error
	Update(ctx context.Context, row *Invoice) error
	LoadOrder(ctx context.Context, orderID uint) (*OrderMeta, error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) Apply(ctx context.Context, uid uint, in ApplyInput) (*Invoice, error) {
	if uid == 0 || in.OrderID == 0 {
		return nil, ErrBadParam
	}
	header := strings.TrimSpace(in.Header)
	if header == "" {
		return nil, ErrBadParam
	}
	if _, err := s.store.FindByOrder(ctx, in.OrderID); err == nil {
		return nil, ErrExists
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	ord, err := s.store.LoadOrder(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrder
		}
		return nil, err
	}
	if ord.UID != uid || !ord.Paid {
		return nil, ErrOrder
	}
	invType := in.InvoiceType
	if invType == 0 {
		invType = 1
	}
	headerType := in.HeaderType
	if headerType == 0 {
		headerType = 1
	}
	row := &Invoice{
		UID: uid, OrderID: in.OrderID, MerID: ord.MerID,
		InvoiceType: invType, HeaderType: headerType,
		Header: header, TaxNo: strings.TrimSpace(in.TaxNo), Email: strings.TrimSpace(in.Email),
		Status: StatusPending, CreateTime: time.Now(),
	}
	if err := s.store.Create(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListMine(ctx context.Context, uid uint, page, limit int) (*PageResult[Invoice], error) {
	if uid == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListByUID(ctx, uid, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Invoice]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, page, limit int) (*PageResult[Invoice], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListByMer(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Invoice]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Audit(ctx context.Context, merID, id uint, in AuditInput) (*Invoice, error) {
	row, err := s.store.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if row.MerID != merID || row.IsDel == 1 {
		return nil, ErrForbidden
	}
	if in.Status != StatusIssued && in.Status != StatusRejected {
		return nil, ErrBadParam
	}
	row.Status = in.Status
	row.Mark = strings.TrimSpace(in.Mark)
	if err := s.store.Update(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
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
