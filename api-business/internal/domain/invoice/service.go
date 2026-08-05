package invoice

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OrderMeta struct {
	OrderID uint64
	UserID  uint64
	Status  string
}

type Store interface {
	ListProfiles(ctx context.Context, userID uint64) ([]InvoiceProfile, error)
	GetProfile(ctx context.Context, userID, id uint64) (*InvoiceProfile, error)
	CreateProfile(ctx context.Context, row *InvoiceProfile) error
	UpdateProfile(ctx context.Context, row *InvoiceProfile) error
	DeleteProfile(ctx context.Context, userID, id uint64) error
	SetDefaultProfile(ctx context.Context, userID, id uint64) error
	ListByUID(ctx context.Context, userID uint64, page, limit int) ([]Invoice, int64, error)
	GetByUID(ctx context.Context, userID, id uint64) (*Invoice, error)
	FindByOrder(ctx context.Context, orderID uint64) (*Invoice, error)
	Create(ctx context.Context, row *Invoice) error
	LoadOrder(ctx context.Context, orderID uint64) (*OrderMeta, error)
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListProfiles(ctx context.Context, userID uint64) ([]InvoiceProfile, error) {
	if userID == 0 {
		return nil, ErrBadParam
	}
	return s.store.ListProfiles(ctx, userID)
}

func (s *Service) CreateProfile(ctx context.Context, userID uint64, in ProfileInput) (*InvoiceProfile, error) {
	if userID == 0 {
		return nil, ErrBadParam
	}
	row, err := profileFromInput(in)
	if err != nil {
		return nil, err
	}
	row.UserID = userID
	if err := s.store.CreateProfile(ctx, &row); err != nil {
		return nil, err
	}
	if row.IsDefault {
		if err := s.store.SetDefaultProfile(ctx, userID, row.ID); err != nil {
			return nil, err
		}
	}
	return &row, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID, id uint64, in ProfileInput) (*InvoiceProfile, error) {
	if userID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := profileFromInput(in)
	if err != nil {
		return nil, err
	}
	current, err := s.store.GetProfile(ctx, userID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	row.ID, row.UserID, row.CreatedAt = current.ID, current.UserID, current.CreatedAt
	if err := s.store.UpdateProfile(ctx, &row); err != nil {
		return nil, err
	}
	if row.IsDefault {
		if err := s.store.SetDefaultProfile(ctx, userID, id); err != nil {
			return nil, err
		}
	}
	return &row, nil
}

func (s *Service) DeleteProfile(ctx context.Context, userID, id uint64) error {
	if userID == 0 || id == 0 {
		return ErrBadParam
	}
	if err := s.store.DeleteProfile(ctx, userID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProfileNotFound
		}
		return err
	}
	return nil
}

func (s *Service) SetDefaultProfile(ctx context.Context, userID, id uint64) error {
	if userID == 0 || id == 0 {
		return ErrBadParam
	}
	if err := s.store.SetDefaultProfile(ctx, userID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProfileNotFound
		}
		return err
	}
	return nil
}

func (s *Service) Apply(ctx context.Context, userID uint64, in ApplyInput) (*Invoice, error) {
	if userID == 0 || in.OrderID == 0 || in.InvoiceProfileID == 0 {
		return nil, ErrBadParam
	}
	profile, err := s.store.GetProfile(ctx, userID, in.InvoiceProfileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProfileNotFound
		}
		return nil, err
	}
	order, err := s.store.LoadOrder(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrder
		}
		return nil, err
	}
	if order.UserID != userID || !invoiceableOrderStatus(order.Status) {
		return nil, ErrOrder
	}
	if _, err := s.store.FindByOrder(ctx, in.OrderID); err == nil {
		return nil, ErrExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	now := time.Now()
	row := &Invoice{
		OrderID: in.OrderID, InvoiceProfileID: profile.ID, ProfileType: profile.Type,
		Title: profile.Title, TaxNo: profile.TaxNo, Email: profile.Email,
		Status: StatusRequested, RequestedAt: now, UpdatedAt: now,
	}
	if err := s.store.Create(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) ListMine(ctx context.Context, userID uint64, page, limit int) (*PageResult[Invoice], error) {
	if userID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListByUID(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[Invoice]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetMine(ctx context.Context, userID, id uint64) (*Invoice, error) {
	if userID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetByUID(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return row, err
}

func profileFromInput(in ProfileInput) (InvoiceProfile, error) {
	row := InvoiceProfile{Type: strings.TrimSpace(in.Type), Title: strings.TrimSpace(in.Title), TaxNo: strings.TrimSpace(in.TaxNo), Email: strings.TrimSpace(in.Email), IsDefault: in.IsDefault}
	if (row.Type != ProfilePersonal && row.Type != ProfileEnterprise) || row.Title == "" || len([]rune(row.Title)) > 255 || len([]rune(row.TaxNo)) > 64 || len([]rune(row.Email)) > 255 {
		return InvoiceProfile{}, ErrBadParam
	}
	if row.Type == ProfileEnterprise && row.TaxNo == "" {
		return InvoiceProfile{}, ErrBadParam
	}
	if row.Email != "" {
		parsed, err := mail.ParseAddress(row.Email)
		if err != nil || parsed.Address != row.Email {
			return InvoiceProfile{}, ErrBadParam
		}
	}
	return row, nil
}

func invoiceableOrderStatus(status string) bool {
	return status == "paid" || status == "fulfilling" || status == "shipped" || status == "completed"
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
