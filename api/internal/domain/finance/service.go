package finance

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	WithTx(fn func(tx Store) error) error

	GetMerMoney(ctx context.Context, merID uint) (float64, error)
	DeductMerMoney(ctx context.Context, merID uint, amount float64) error
	AddMerMoney(ctx context.Context, merID uint, amount float64) error

	CreateFinancial(ctx context.Context, f *Financial) error
	CreateRecord(ctx context.Context, rec *FinancialRecord) error
	GetFinancial(ctx context.Context, id uint) (*Financial, error)
	ListFinancials(ctx context.Context, filter ListFilter) ([]Financial, int64, error)
	UpdateFinancialStatus(ctx context.Context, id uint, fromStatus, toStatus int, refusal string, adminID uint) (bool, error)
}

type ListFilter struct {
	MerID  *uint
	Status *int
	Page   int
	Limit  int
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Balance(ctx context.Context, merID uint) (*Balance, error) {
	money, err := s.store.GetMerMoney(ctx, merID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMerchantNotFound
		}
		return nil, err
	}
	return &Balance{MerID: merID, MerMoney: money}, nil
}

func (s *Service) ApplyWithdraw(ctx context.Context, merID uint, merAdminID uint, in WithdrawInput) (*Financial, error) {
	if in.ExtractMoney <= 0 {
		return nil, ErrBadParam
	}
	if in.FinancialType < AccountBank || in.FinancialType > AccountAlipay {
		return nil, ErrBadParam
	}
	account := strings.TrimSpace(in.FinancialAccount)
	if account == "" {
		return nil, ErrBadParam
	}
	mark := strings.TrimSpace(in.Mark)

	var created *Financial
	err := s.store.WithTx(func(tx Store) error {
		bal, err := tx.GetMerMoney(ctx, merID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrMerchantNotFound
			}
			return err
		}
		if bal < in.ExtractMoney {
			return ErrBalanceNotEnough
		}
		if err := tx.DeductMerMoney(ctx, merID, in.ExtractMoney); err != nil {
			return err
		}
		now := time.Now()
		adminID := int(merAdminID)
		f := &Financial{
			FinancialSN:      genSN("F"),
			MerMoney:         bal,
			ExtractMoney:     in.ExtractMoney,
			FinancialType:    in.FinancialType,
			FinancialAccount: account,
			FinancialStatus:  FinancialStatusUnpaid,
			Status:           StatusWait,
			MerID:            merID,
			CreateTime:       &now,
			Mark:             mark,
			MerAdminID:       &adminID,
			Type:             TypeExtract,
		}
		if err := tx.CreateFinancial(ctx, f); err != nil {
			return err
		}
		rec := &FinancialRecord{
			FinancialRecordSN: genSN("FR"),
			OrderSN:           f.FinancialSN,
			UserInfo:          fmt.Sprintf("mer:%d", merID),
			FinancialType:     "extract",
			FinancialPM:       RecordPMOut,
			Number:            in.ExtractMoney,
			Type:              0,
			MerID:             merID,
			CreateTime:        now,
		}
		if err := tx.CreateRecord(ctx, rec); err != nil {
			return err
		}
		created = f
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, page, limit int) (*PageResult[Financial], error) {
	list, total, err := s.store.ListFinancials(ctx, ListFilter{MerID: &merID, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[Financial]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetMerchant(ctx context.Context, merID, id uint) (*Financial, error) {
	f, err := s.store.GetFinancial(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if f.IsDel != 0 {
		return nil, ErrNotFound
	}
	if f.MerID != merID {
		return nil, ErrForbidden
	}
	return f, nil
}

func (s *Service) ListPlatform(ctx context.Context, status *int, page, limit int) (*PageResult[Financial], error) {
	list, total, err := s.store.ListFinancials(ctx, ListFilter{Status: status, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[Financial]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetPlatform(ctx context.Context, id uint) (*Financial, error) {
	f, err := s.store.GetFinancial(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if f.IsDel != 0 {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *Service) Approve(ctx context.Context, adminID, id uint) error {
	f, err := s.store.GetFinancial(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if f.IsDel != 0 {
		return ErrNotFound
	}
	if f.Status == StatusPass {
		return ErrAlreadyDone
	}
	if f.Status != StatusWait {
		return ErrBadStatus
	}
	ok, err := s.store.UpdateFinancialStatus(ctx, id, StatusWait, StatusPass, "", adminID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAlreadyDone
	}
	return nil
}

func (s *Service) Reject(ctx context.Context, adminID, id uint, in RejectInput) error {
	refusal := strings.TrimSpace(in.Refusal)
	if refusal == "" {
		return ErrRejectNeedMsg
	}
	f, err := s.store.GetFinancial(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if f.IsDel != 0 {
		return ErrNotFound
	}
	if f.Status == StatusReject {
		return ErrAlreadyDone
	}
	if f.Status != StatusWait {
		return ErrBadStatus
	}

	return s.store.WithTx(func(tx Store) error {
		ok, err := tx.UpdateFinancialStatus(ctx, id, StatusWait, StatusReject, refusal, adminID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrAlreadyDone
		}
		if err := tx.AddMerMoney(ctx, f.MerID, f.ExtractMoney); err != nil {
			return err
		}
		now := time.Now()
		rec := &FinancialRecord{
			FinancialRecordSN: genSN("FR"),
			OrderSN:           f.FinancialSN,
			UserInfo:          fmt.Sprintf("mer:%d", f.MerID),
			FinancialType:     "extract_reject",
			FinancialPM:       RecordPMIn,
			Number:            f.ExtractMoney,
			Type:              0,
			MerID:             f.MerID,
			CreateTime:        now,
		}
		return tx.CreateRecord(ctx, rec)
	})
}

func genSN(prefix string) string {
	return fmt.Sprintf("%s%s%04d", prefix, time.Now().Format("20060102150405"), rand.Intn(10000))
}

func normalizePage(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
