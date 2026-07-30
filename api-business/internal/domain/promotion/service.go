package promotion

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	WithTx(fn func(tx Store) error) error

	CreateCoupon(ctx context.Context, c *Coupon) error
	GetCoupon(ctx context.Context, id uint) (*Coupon, error)
	UpdateCoupon(ctx context.Context, c *Coupon) error
	UpdateCouponStatus(ctx context.Context, id uint, merID *uint, status int8) (bool, error)
	SoftDeleteCoupon(ctx context.Context, id uint, merID *uint) (bool, error)
	ListCoupons(ctx context.Context, merID *uint, typ *int, page, limit int) ([]Coupon, int64, error)
	ListCenter(ctx context.Context, page, limit int) ([]Coupon, int64, error)
	DecRemain(ctx context.Context, couponID uint) (bool, error)

	HasReceived(ctx context.Context, uid, couponID uint) (bool, error)
	CreateIssueUser(ctx context.Context, row *IssueUser) error
	CreateCouponUser(ctx context.Context, u *CouponUser) error
	CreateCouponSend(ctx context.Context, row *CouponSend) error
	MarkCouponSendDone(ctx context.Context, id uint) error
	GetCouponUser(ctx context.Context, id uint) (*CouponUser, error)
	ListCouponUsers(ctx context.Context, uid uint, status *int, page, limit int) ([]CouponUser, int64, error)
	HasMerchantCouponUser(ctx context.Context, merID, uid, couponID uint) (bool, error)
	ListMerchantCouponUsers(ctx context.Context, merID uint, couponID *uint, page, limit int) ([]CouponUser, int64, error)
	ListMerchantCouponSends(ctx context.Context, merID uint, page, limit int) ([]CouponSend, int64, error)
	ListCouponUsersByIDs(ctx context.Context, uid uint, ids []uint) ([]CouponUser, error)
	ListUsablePlatform(ctx context.Context, uid uint, orderAmount float64) ([]CouponUser, error)
	MarkCouponUsersUsed(ctx context.Context, uid uint, ids []uint, at time.Time) (int64, error)
	HasMerchantPaidOrderUser(ctx context.Context, merID, uid uint) (bool, error)

	GetUserSpread(ctx context.Context, uid uint) (spreadUID uint, isPromoter int8, err error)
	IsPromoter(ctx context.Context, uid uint) (bool, error)
	SetUserSpread(ctx context.Context, uid, spreadUID uint) error
	CreateSpreadLog(ctx context.Context, log *SpreadLog) error
	ListSpreadLogs(ctx context.Context, page, limit int) ([]SpreadLog, int64, error)
	CountSpreadChildren(ctx context.Context, spreadUID uint) (int64, error)

	AddBrokerage(ctx context.Context, uid uint, amount float64) (balance float64, err error)
	CreateBill(ctx context.Context, b *UserBill) error
	HasBill(ctx context.Context, uid uint, category, typ, linkID string) (bool, error)
	ListBills(ctx context.Context, uid *uint, category string, page, limit int) ([]UserBill, int64, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) CreatePlatformCoupon(ctx context.Context, in CreateCouponInput) (*Coupon, error) {
	return s.createCoupon(ctx, 0, CouponTypePlatform, in)
}

func (s *Service) CreateMerchantCoupon(ctx context.Context, merID uint, in CreateCouponInput) (*Coupon, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	return s.createCoupon(ctx, merID, CouponTypeStore, in)
}

func (s *Service) createCoupon(ctx context.Context, merID uint, typ int, in CreateCouponInput) (*Coupon, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" || in.CouponPrice <= 0 {
		return nil, ErrBadParam
	}
	if in.CouponTime == 0 {
		in.CouponTime = 30
	}
	c := &Coupon{
		MerID:       merID,
		Title:       title,
		CouponPrice: round2(in.CouponPrice),
		UseMinPrice: in.UseMinPrice,
		CouponType:  int8(TemplateDays),
		CouponTime:  in.CouponTime,
		IsLimited:   in.IsLimited,
		TotalCount:  in.TotalCount,
		RemainCount: in.TotalCount,
		SendType:    0,
		Sort:        1,
		Status:      1,
		Type:        typ,
	}
	if c.IsLimited == 0 {
		c.TotalCount = 0
		c.RemainCount = 0
	}
	if err := s.store.CreateCoupon(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListPlatformCoupons(ctx context.Context, page, limit int) (*PageResult[Coupon], error) {
	typ := CouponTypePlatform
	mer := uint(0)
	list, total, err := s.store.ListCoupons(ctx, &mer, &typ, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[Coupon]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchantCoupons(ctx context.Context, merID uint, page, limit int) (*PageResult[Coupon], error) {
	typ := CouponTypeStore
	list, total, err := s.store.ListCoupons(ctx, &merID, &typ, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[Coupon]{List: list, Total: total, Page: page, Limit: limit}, nil
}

// ListAdmin / CreateAdmin / UpdateAdmin / SetStatus / DeleteAdmin：管理端统一入口。
func (s *Service) ListAdmin(ctx context.Context, merID uint, typ int, page, limit int) (*PageResult[Coupon], error) {
	if typ == CouponTypePlatform {
		return s.ListPlatformCoupons(ctx, page, limit)
	}
	return s.ListMerchantCoupons(ctx, merID, page, limit)
}

// SendMerchantCoupon 将一张已启用的店铺券定向发给本商户已有支付订单的用户。
// 所有收件人、库存扣减、批次记录和用户券在同一事务中处理，避免限量券超发。
func (s *Service) SendMerchantCoupon(ctx context.Context, merID, couponID uint, in CouponSendInput) (*CouponSend, error) {
	if merID == 0 || couponID == 0 {
		return nil, ErrBadParam
	}
	uids := uniqueUint(in.UIDs)
	if len(uids) == 0 || len(uids) > 100 {
		return nil, ErrBadParam
	}
	for _, uid := range uids {
		ok, err := s.store.HasMerchantPaidOrderUser(ctx, merID, uid)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrForbidden
		}
	}

	mark := strings.TrimSpace(in.Mark)
	if mark == "" {
		mark = "商户后台定向发券"
	}
	var created *CouponSend
	err := s.store.WithTx(func(tx Store) error {
		c, err := tx.GetCoupon(ctx, couponID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if c.IsDel == 1 || c.MerID != merID || c.Type != CouponTypeStore {
			return ErrForbidden
		}
		if c.Status != 1 {
			return ErrClosed
		}
		for _, uid := range uids {
			got, err := tx.HasMerchantCouponUser(ctx, merID, uid, couponID)
			if err != nil {
				return err
			}
			if got {
				return ErrAlreadyReceived
			}
		}

		now := time.Now()
		send := &CouponSend{MerID: merID, CouponID: couponID, CouponNum: uint(len(uids)), Mark: mark, CreateTime: now, Status: 0}
		if err := tx.CreateCouponSend(ctx, send); err != nil {
			return err
		}
		validDays := c.CouponTime
		if validDays == 0 {
			validDays = 30
		}
		for _, uid := range uids {
			if c.IsLimited == 1 {
				ok, err := tx.DecRemain(ctx, couponID)
				if err != nil {
					return err
				}
				if !ok {
					return ErrSoldOut
				}
			}
			start, end := now, now.Add(time.Duration(validDays)*24*time.Hour)
			row := &CouponUser{
				CouponID: couponID, MerID: merID, UID: uid, CouponTitle: truncate(c.Title, 32),
				CouponPrice: c.CouponPrice, UseMinPrice: c.UseMinPrice, StartTime: &start, EndTime: &end,
				Type: "send", SendID: send.CouponSendID, Status: UserUnused,
			}
			if err := tx.CreateCouponUser(ctx, row); err != nil {
				return err
			}
		}
		if err := tx.MarkCouponSendDone(ctx, send.CouponSendID); err != nil {
			return err
		}
		send.Status = 1
		created = send
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) ListMerchantCouponUsers(ctx context.Context, merID uint, couponID *uint, page, limit int) (*PageResult[CouponUser], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	list, total, err := s.store.ListMerchantCouponUsers(ctx, merID, couponID, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[CouponUser]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) ListMerchantCouponSends(ctx context.Context, merID uint, page, limit int) (*PageResult[CouponSend], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	list, total, err := s.store.ListMerchantCouponSends(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[CouponSend]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreateAdmin(ctx context.Context, merID uint, typ int, in CouponSaveInput) (*Coupon, error) {
	if typ == CouponTypePlatform {
		return s.CreatePlatformCoupon(ctx, in)
	}
	return s.CreateMerchantCoupon(ctx, merID, in)
}

func (s *Service) UpdateAdmin(ctx context.Context, merID uint, typ int, id uint, in CouponSaveInput) (*Coupon, error) {
	c, err := s.store.GetCoupon(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if c.IsDel == 1 || c.MerID != merID || c.Type != typ {
		return nil, ErrForbidden
	}
	title := strings.TrimSpace(in.Title)
	if title == "" || in.CouponPrice <= 0 {
		return nil, ErrBadParam
	}
	c.Title = title
	c.CouponPrice = round2(in.CouponPrice)
	c.UseMinPrice = in.UseMinPrice
	c.IsLimited = in.IsLimited
	if in.IsLimited == 1 {
		if in.TotalCount > 0 {
			c.TotalCount = in.TotalCount
			if c.RemainCount > in.TotalCount {
				c.RemainCount = in.TotalCount
			}
		}
	} else {
		c.TotalCount = 0
		c.RemainCount = 0
	}
	if in.CouponTime > 0 {
		c.CouponTime = in.CouponTime
	}
	if in.Sort > 0 {
		c.Sort = in.Sort
	}
	if in.Status != nil {
		c.Status = *in.Status
	}
	if err := s.store.UpdateCoupon(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) SetStatus(ctx context.Context, merID uint, typ int, id uint, status int8) error {
	if status != 0 && status != 1 {
		return ErrBadStatus
	}
	c, err := s.store.GetCoupon(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if c.IsDel == 1 || c.MerID != merID || c.Type != typ {
		return ErrForbidden
	}
	ok, err := s.store.UpdateCouponStatus(ctx, id, &merID, status)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

// DeleteAdmin 软删优惠券（is_del=1）。
func (s *Service) DeleteAdmin(ctx context.Context, merID uint, typ int, id uint) error {
	c, err := s.store.GetCoupon(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if c.IsDel == 1 || c.MerID != merID || c.Type != typ {
		return ErrForbidden
	}
	ok, err := s.store.SoftDeleteCoupon(ctx, id, &merID)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s *Service) Center(ctx context.Context, page, limit int) (*PageResult[Coupon], error) {
	list, total, err := s.store.ListCenter(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[Coupon]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Mine(ctx context.Context, uid uint, status *int, page, limit int) (*PageResult[CouponUser], error) {
	return s.MyCoupons(ctx, uid, status, page, limit)
}

func (s *Service) Usable(ctx context.Context, uid uint, orderAmount float64) ([]CouponUser, error) {
	return s.store.ListUsablePlatform(ctx, uid, orderAmount)
}

// ValidateAndQuote P0：单张平台券报价。
func (s *Service) ValidateAndQuote(ctx context.Context, uid, couponUserID uint, orderAmount float64) (float64, *CouponUser, error) {
	if couponUserID == 0 {
		return 0, nil, nil
	}
	q, err := s.Quote(ctx, uid, QuoteInput{
		MerTotals:     []MerTotal{{MerID: 0, TotalPrice: orderAmount}},
		CouponUserIDs: []uint{couponUserID},
	})
	if err != nil {
		return 0, nil, err
	}
	if q.PlatformCouponUserID == 0 {
		return 0, nil, ErrPlatformOnly
	}
	cu, err := s.store.GetCouponUser(ctx, couponUserID)
	if err != nil {
		return 0, nil, ErrCouponInvalid
	}
	cu.CouponKind = CouponTypePlatform
	return q.PlatformDiscount, cu, nil
}

// MarkUsed CAS 核销单张用户券。
func (s *Service) MarkUsed(ctx context.Context, couponUserID uint) error {
	if couponUserID == 0 {
		return nil
	}
	cu, err := s.store.GetCouponUser(ctx, couponUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCouponInvalid
		}
		return err
	}
	return s.ConsumeMarks(ctx, cu.UID, []uint{couponUserID})
}

func (s *Service) SpreadMe(ctx context.Context, uid uint) (spreadUID uint, isPromoter int8, spreadCount int64, err error) {
	spreadUID, isPromoter, err = s.store.GetUserSpread(ctx, uid)
	if err != nil {
		return 0, 0, 0, err
	}
	spreadCount, err = s.store.CountSpreadChildren(ctx, uid)
	return
}

func (s *Service) ListReceivable(ctx context.Context, uid uint, merID *uint) ([]Coupon, error) {
	var list []Coupon
	var err error
	if merID != nil && *merID > 0 {
		typ := CouponTypeStore
		list, _, err = s.store.ListCoupons(ctx, merID, &typ, 1, 50)
	} else {
		// 领券中心：平台券 + 各店可领店铺券
		list, _, err = s.store.ListCenter(ctx, 1, 50)
	}
	if err != nil {
		return nil, err
	}
	out := make([]Coupon, 0, len(list))
	for _, c := range list {
		if c.Status != 1 || c.IsDel == 1 {
			continue
		}
		if c.IsLimited == 1 && c.RemainCount == 0 {
			continue
		}
		got, err := s.store.HasReceived(ctx, uid, c.CouponID)
		if err != nil {
			return nil, err
		}
		c.Received = got
		out = append(out, c)
	}
	return out, nil
}

func (s *Service) Receive(ctx context.Context, uid, couponID uint) (*CouponUser, error) {
	var created *CouponUser
	err := s.store.WithTx(func(tx Store) error {
		c, err := tx.GetCoupon(ctx, couponID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if c.Status != 1 || c.IsDel == 1 {
			return ErrClosed
		}
		// P0 C 端领取仅平台券
		if c.Type != CouponTypePlatform || c.MerID != 0 {
			return ErrNotAvailable
		}
		got, err := tx.HasReceived(ctx, uid, couponID)
		if err != nil {
			return err
		}
		if got {
			return ErrAlreadyReceived
		}
		if c.IsLimited == 1 {
			ok, err := tx.DecRemain(ctx, couponID)
			if err != nil {
				return err
			}
			if !ok {
				return ErrSoldOut
			}
		}
		now := time.Now()
		start := now
		end := now.Add(time.Duration(c.CouponTime) * 24 * time.Hour)
		if c.CouponTime == 0 {
			end = now.Add(30 * 24 * time.Hour)
		}
		u := &CouponUser{
			CouponID:    c.CouponID,
			MerID:       c.MerID,
			UID:         uid,
			CouponTitle: truncate(c.Title, 32),
			CouponPrice: c.CouponPrice,
			UseMinPrice: c.UseMinPrice,
			StartTime:   &start,
			EndTime:     &end,
			Type:        "receive",
			Status:      UserUnused,
		}
		if err := tx.CreateCouponUser(ctx, u); err != nil {
			return err
		}
		if err := tx.CreateIssueUser(ctx, &IssueUser{UID: uid, CouponID: couponID, CreateTime: now}); err != nil {
			return err
		}
		u.CouponKind = c.Type
		created = u
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) MyCoupons(ctx context.Context, uid uint, status *int, page, limit int) (*PageResult[CouponUser], error) {
	list, total, err := s.store.ListCouponUsers(ctx, uid, status, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if c, err := s.store.GetCoupon(ctx, list[i].CouponID); err == nil {
			list[i].CouponKind = c.Type
		}
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[CouponUser]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) Quote(ctx context.Context, uid uint, in QuoteInput) (*QuoteResult, error) {
	if len(in.CouponUserIDs) == 0 {
		q := CalcQuote(in.MerTotals, nil, nil)
		return &q, nil
	}
	rows, err := s.store.ListCouponUsersByIDs(ctx, uid, in.CouponUserIDs)
	if err != nil {
		return nil, err
	}
	if len(rows) != len(uniqueUint(in.CouponUserIDs)) {
		return nil, ErrCouponInvalid
	}
	now := time.Now()
	storeByMer := map[uint]CouponUser{}
	var platform *CouponUser
	for _, u := range rows {
		if u.Status != UserUnused || u.IsFail == 1 {
			return nil, ErrCouponInvalid
		}
		if u.EndTime != nil && u.EndTime.Before(now) {
			return nil, ErrCouponInvalid
		}
		if u.StartTime != nil && u.StartTime.After(now) {
			return nil, ErrCouponInvalid
		}
		c, err := s.store.GetCoupon(ctx, u.CouponID)
		if err != nil {
			return nil, ErrCouponInvalid
		}
		u.CouponKind = c.Type
		switch c.Type {
		case CouponTypePlatform:
			if in.SkipPlatformCoupon {
				continue
			}
			if platform != nil {
				return nil, ErrCouponConflict
			}
			cp := u
			platform = &cp
		case CouponTypeStore:
			if in.SkipStoreCoupon {
				continue // 秒杀等：忽略店铺券
			}
			if u.MerID == 0 {
				return nil, ErrCouponInvalid
			}
			if _, ok := storeByMer[u.MerID]; ok {
				return nil, ErrCouponConflict
			}
			found := false
			for _, m := range in.MerTotals {
				if m.MerID == u.MerID {
					found = true
					break
				}
			}
			if !found {
				return nil, ErrCouponInvalid
			}
			storeByMer[u.MerID] = u
		default:
			return nil, ErrCouponInvalid
		}
	}
	// 门槛未满足时视为不可用（创建时报错；check 可先不算）
	q := CalcQuote(in.MerTotals, storeByMer, platform)
	if platform != nil && q.PlatformDiscount == 0 {
		return nil, ErrCouponMinNotMet
	}
	for merID, cu := range storeByMer {
		if q.MerStoreDiscount[merID] == 0 && cu.CouponUserID != 0 {
			return nil, ErrCouponMinNotMet
		}
	}
	return &q, nil
}

// ConsumeMarks 在订单事务内核销用户券（由 trade persist 调用同库更新亦可）。
func (s *Service) ConsumeMarks(ctx context.Context, uid uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	n, err := s.store.MarkCouponUsersUsed(ctx, uid, uniqueUint(ids), time.Now())
	if err != nil {
		return err
	}
	if int(n) != len(uniqueUint(ids)) {
		return ErrCouponInvalid
	}
	return nil
}

func (s *Service) BindSpread(ctx context.Context, uid, spreadUID uint) error {
	if uid == 0 || spreadUID == 0 {
		return ErrBadParam
	}
	if uid == spreadUID {
		return ErrSpreadSelf
	}
	return s.store.WithTx(func(tx Store) error {
		cur, _, err := tx.GetUserSpread(ctx, uid)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNotFound
			}
			return err
		}
		if cur != 0 {
			return ErrSpreadBound
		}
		ok, err := tx.IsPromoter(ctx, spreadUID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrSpreadInvalid
		}
		if err := tx.SetUserSpread(ctx, uid, spreadUID); err != nil {
			return err
		}
		return tx.CreateSpreadLog(ctx, &SpreadLog{
			UID:          uid,
			OldSpreadUID: 0,
			SpreadUID:    spreadUID,
			CreateTime:   time.Now(),
		})
	})
}

func (s *Service) ListSpreadLogs(ctx context.Context, page, limit int) (*PageResult[SpreadLog], error) {
	list, total, err := s.store.ListSpreadLogs(ctx, page, limit)
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[SpreadLog]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) CreditSpreadOnPay(ctx context.Context, buyerUID, groupOrderID uint, payPrice float64) error {
	if payPrice <= 0 || buyerUID == 0 {
		return nil
	}
	spreadUID, _, err := s.store.GetUserSpread(ctx, buyerUID)
	if err != nil || spreadUID == 0 {
		return nil
	}
	link := fmt.Sprintf("%d", groupOrderID)
	dup, err := s.store.HasBill(ctx, spreadUID, "now_money", "brokerage", link)
	if err != nil {
		return err
	}
	if dup {
		return nil
	}
	amount := round2(payPrice * SpreadRate)
	if amount <= 0 {
		return nil
	}
	return s.store.WithTx(func(tx Store) error {
		bal, err := tx.AddBrokerage(ctx, spreadUID, amount)
		if err != nil {
			return err
		}
		return tx.CreateBill(ctx, &UserBill{
			UID:      spreadUID,
			LinkID:   link,
			PM:       BillPMIn,
			Title:    "推广佣金",
			Category: "now_money",
			Type:     "brokerage",
			Number:   amount,
			Balance:  bal,
			Mark:     fmt.Sprintf("下级订单%d佣金", groupOrderID),
			Status:   1,
		})
	})
}

func (s *Service) ListBrokerageBills(ctx context.Context, uid *uint, page, limit int) (*PageResult[UserBill], error) {
	list, total, err := s.store.ListBills(ctx, uid, "now_money", page, limit)
	if err != nil {
		return nil, err
	}
	// 仅佣金类型
	filtered := make([]UserBill, 0, len(list))
	for _, b := range list {
		if b.Type == "brokerage" {
			filtered = append(filtered, b)
		}
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[UserBill]{List: filtered, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) MyBrokerage(ctx context.Context, uid uint, page, limit int) (*PageResult[UserBill], error) {
	return s.ListBrokerageBills(ctx, &uid, page, limit)
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

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func uniqueUint(ids []uint) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
