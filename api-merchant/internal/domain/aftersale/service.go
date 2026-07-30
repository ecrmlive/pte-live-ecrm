package aftersale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	WithTx(fn func(tx Store) error) error

	GetStoreOrder(ctx context.Context, orderID uint) (*StoreOrderBrief, error)
	ListOrderProducts(ctx context.Context, orderID uint) ([]OrderProductLine, error)
	ListOrderProductsByIDs(ctx context.Context, orderID uint, ids []uint) ([]OrderProductLine, error)

	HasActiveRefund(ctx context.Context, orderID uint) (bool, error)
	CreateRefundOrder(ctx context.Context, o *RefundOrder) error
	CreateRefundProducts(ctx context.Context, rows []RefundProduct) error
	CreateRefundStatus(ctx context.Context, log *RefundStatusLog) error

	GetRefund(ctx context.Context, id uint) (*RefundOrder, error)
	ListRefunds(ctx context.Context, filter ListFilter) ([]RefundOrder, int64, error)
	ListPlatformRefundsByRegions(ctx context.Context, regionIDs []uint, status *int8, page, limit int) ([]RefundOrder, int64, error)
	GetPlatformRefundByRegions(ctx context.Context, id uint, regionIDs []uint) (*RefundOrder, error)
	ListRefundProducts(ctx context.Context, refundOrderID uint) ([]RefundProduct, error)

	// CAS: only update when current status matches fromStatus
	UpdateRefundStatus(ctx context.Context, id uint, fromStatus, toStatus int8, failMessage string) (bool, error)

	AddUserBalance(ctx context.Context, uid uint, amount float64) error
	AddProductStock(ctx context.Context, productID uint, num uint) error
	AddSKUStock(ctx context.Context, productID uint, unique string, num uint) error

	MarkOrderProductRefund(ctx context.Context, orderProductID uint, addNum int, isRefund int8) error
	UpdateOrderStatus(ctx context.Context, orderID uint, status int8) error
	CountUnrefundedProducts(ctx context.Context, orderID uint) (int64, error)
}

type ListFilter struct {
	UID    *uint
	MerID  *uint
	Status *int8
	Page   int
	Limit  int
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) Apply(ctx context.Context, uid uint, in ApplyInput) (*RefundOrder, error) {
	return s.applyRefund(ctx, uid, 0, in)
}

// ApplyBehalf 店员代用户发起仅退款（强制 mer_id 隔离，退款单仍挂订单买家 uid）。
func (s *Service) ApplyBehalf(ctx context.Context, merID uint, in ApplyInput) (*RefundOrder, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	return s.applyRefund(ctx, 0, merID, in)
}

func (s *Service) applyRefund(ctx context.Context, uid, staffMerID uint, in ApplyInput) (*RefundOrder, error) {
	if in.OrderID == 0 {
		return nil, ErrBadParam
	}
	if in.RefundType != RefundTypeMoneyOnly {
		return nil, ErrBadParam
	}
	msg := strings.TrimSpace(in.RefundMessage)
	if msg == "" {
		return nil, ErrBadParam
	}

	order, err := s.store.GetStoreOrder(ctx, in.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	if order.IsDel == 1 {
		return nil, ErrForbidden
	}
	if staffMerID > 0 {
		if order.MerID != staffMerID {
			return nil, ErrForbidden
		}
		uid = order.UID
	} else if order.UID != uid {
		return nil, ErrForbidden
	}
	if order.Paid != 1 {
		return nil, ErrOrderNotPaid
	}
	if order.Status == OrderStatusRefunded {
		return nil, ErrOrderRefunded
	}

	active, err := s.store.HasActiveRefund(ctx, in.OrderID)
	if err != nil {
		return nil, err
	}
	if active {
		return nil, ErrRefundInProgress
	}

	var products []OrderProductLine
	if len(in.OrderProductIDs) == 0 {
		products, err = s.store.ListOrderProducts(ctx, in.OrderID)
	} else {
		products, err = s.store.ListOrderProductsByIDs(ctx, in.OrderID, in.OrderProductIDs)
	}
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, ErrProductInvalid
	}
	if len(in.OrderProductIDs) > 0 && len(products) != len(uniqueUint(in.OrderProductIDs)) {
		return nil, ErrProductInvalid
	}

	var refundPrice float64
	var refundNum int
	rps := make([]RefundProduct, 0, len(products))
	for _, p := range products {
		remain := p.ProductNum - p.RefundNum
		if remain <= 0 || p.IsRefund == OrderProductRefundFull {
			return nil, ErrProductInvalid
		}
		refundPrice += p.TotalPrice
		refundNum += remain
		rps = append(rps, RefundProduct{
			OrderProductID: p.OrderProductID,
			RefundPrice:    p.TotalPrice,
			RefundNum:      remain,
		})
	}

	var created *RefundOrder
	err = s.store.WithTx(func(tx Store) error {
		again, err := tx.HasActiveRefund(ctx, in.OrderID)
		if err != nil {
			return err
		}
		if again {
			return ErrRefundInProgress
		}
		now := time.Now()
		ro := &RefundOrder{
			RefundOrderSN: genSN("R"),
			OrderID:       in.OrderID,
			UID:           uid,
			MerID:         order.MerID,
			RefundType:    RefundTypeMoneyOnly,
			RefundMessage: msg,
			RefundPrice:   refundPrice,
			RefundNum:     refundNum,
			Status:        StatusWait,
			StatusTime:    now,
			CreateTime:    now,
		}
		if err := tx.CreateRefundOrder(ctx, ro); err != nil {
			return err
		}
		for i := range rps {
			rps[i].RefundOrderID = ro.RefundOrderID
			rps[i].CreateTime = now
		}
		if err := tx.CreateRefundProducts(ctx, rps); err != nil {
			return err
		}
		changeMsg := "用户申请仅退款"
		if staffMerID > 0 {
			changeMsg = "店员代申请仅退款"
		}
		if err := tx.CreateRefundStatus(ctx, &RefundStatusLog{
			RefundOrderID: ro.RefundOrderID,
			ChangeType:    "apply",
			ChangeMessage: changeMsg,
			ChangeTime:    now,
		}); err != nil {
			return err
		}
		for _, p := range products {
			if err := tx.MarkOrderProductRefund(ctx, p.OrderProductID, 0, OrderProductRefunding); err != nil {
				return err
			}
		}
		ro.Products = rps
		created = ro
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *Service) ListUser(ctx context.Context, uid uint, page, limit int) (*PageResult[RefundOrder], error) {
	list, total, err := s.store.ListRefunds(ctx, ListFilter{UID: &uid, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[RefundOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetUser(ctx context.Context, uid, id uint) (*RefundOrder, error) {
	ro, err := s.getAttach(ctx, id)
	if err != nil {
		return nil, err
	}
	if ro.UID != uid {
		return nil, ErrForbidden
	}
	return ro, nil
}

func (s *Service) Cancel(ctx context.Context, uid, id uint) error {
	ro, err := s.store.GetRefund(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if ro.UID != uid || ro.IsDel == 1 {
		return ErrForbidden
	}
	if ro.Status != StatusWait {
		return ErrBadStatus
	}
	return s.store.WithTx(func(tx Store) error {
		ok, err := tx.UpdateRefundStatus(ctx, id, StatusWait, StatusCancel, "")
		if err != nil {
			return err
		}
		if !ok {
			return ErrBadStatus
		}
		now := time.Now()
		if err := tx.CreateRefundStatus(ctx, &RefundStatusLog{
			RefundOrderID: id,
			ChangeType:    "cancel",
			ChangeMessage: "用户取消退款",
			ChangeTime:    now,
		}); err != nil {
			return err
		}
		products, err := tx.ListRefundProducts(ctx, id)
		if err != nil {
			return err
		}
		for _, rp := range products {
			lines, err := tx.ListOrderProductsByIDs(ctx, ro.OrderID, []uint{rp.OrderProductID})
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				continue
			}
			line := lines[0]
			isRefund := OrderProductRefundNone
			if line.RefundNum > 0 {
				if line.RefundNum >= line.ProductNum {
					isRefund = OrderProductRefundFull
				} else {
					isRefund = OrderProductRefundPartial
				}
			}
			if err := tx.MarkOrderProductRefund(ctx, rp.OrderProductID, 0, isRefund); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) ListMerchant(ctx context.Context, merID uint, status *int8, page, limit int) (*PageResult[RefundOrder], error) {
	list, total, err := s.store.ListRefunds(ctx, ListFilter{MerID: &merID, Status: status, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[RefundOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetMerchant(ctx context.Context, merID, id uint) (*RefundOrder, error) {
	ro, err := s.getAttach(ctx, id)
	if err != nil {
		return nil, err
	}
	if ro.MerID != merID {
		return nil, ErrForbidden
	}
	return ro, nil
}

func (s *Service) Approve(ctx context.Context, merID, id uint) error {
	ro, err := s.store.GetRefund(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	// merID=0：平台监管代审；>0：商户本店隔离
	if (merID > 0 && ro.MerID != merID) || ro.IsDel == 1 {
		return ErrForbidden
	}
	if ro.Status == StatusRefunded {
		return ErrAlreadyDone
	}
	fromStatus, err := refundAuditFromStatus(merID, ro.Status)
	if err != nil {
		return err
	}
	if ro.RefundType != RefundTypeMoneyOnly {
		return ErrBadParam
	}
	approveMsg := "商户同意仅退款，已退款"
	if merID == 0 {
		approveMsg = "平台同意仅退款，已退款"
	}

	return s.store.WithTx(func(tx Store) error {
		ok, err := tx.UpdateRefundStatus(ctx, id, fromStatus, StatusRefunded, "")
		if err != nil {
			return err
		}
		if !ok {
			return ErrAlreadyDone
		}
		now := time.Now()
		if err := tx.CreateRefundStatus(ctx, &RefundStatusLog{
			RefundOrderID: id,
			ChangeType:    "approve",
			ChangeMessage: approveMsg,
			ChangeTime:    now,
		}); err != nil {
			return err
		}

		products, err := tx.ListRefundProducts(ctx, id)
		if err != nil {
			return err
		}
		for _, rp := range products {
			lines, err := tx.ListOrderProductsByIDs(ctx, ro.OrderID, []uint{rp.OrderProductID})
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				return ErrProductInvalid
			}
			line := lines[0]
			n := uint(rp.RefundNum)
			if n == 0 {
				continue
			}
			if err := tx.AddProductStock(ctx, line.ProductID, n); err != nil {
				return err
			}
			if err := tx.AddSKUStock(ctx, line.ProductID, skuUnique(line), n); err != nil {
				return err
			}
			isRefund := OrderProductRefundPartial
			if line.RefundNum+rp.RefundNum >= line.ProductNum {
				isRefund = OrderProductRefundFull
			}
			if err := tx.MarkOrderProductRefund(ctx, line.OrderProductID, rp.RefundNum, isRefund); err != nil {
				return err
			}
		}

		if ro.RefundPrice > 0 {
			if err := tx.AddUserBalance(ctx, ro.UID, ro.RefundPrice); err != nil {
				return err
			}
		}

		left, err := tx.CountUnrefundedProducts(ctx, ro.OrderID)
		if err != nil {
			return err
		}
		if left == 0 {
			if err := tx.UpdateOrderStatus(ctx, ro.OrderID, OrderStatusRefunded); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) Reject(ctx context.Context, merID, id uint, in RejectInput) error {
	msg := strings.TrimSpace(in.FailMessage)
	if msg == "" {
		return ErrRejectNeedMsg
	}
	ro, err := s.store.GetRefund(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if (merID > 0 && ro.MerID != merID) || ro.IsDel == 1 {
		return ErrForbidden
	}
	fromStatus, err := refundAuditFromStatus(merID, ro.Status)
	if err != nil {
		return err
	}
	// 当前售后闭环仅实现仅退款；退货退款的寄回、收货确认等链路尚未实现，
	// 不得通过“拒绝”接口把它提前落入终态。
	if ro.RefundType != RefundTypeMoneyOnly {
		return ErrBadParam
	}
	rejectPrefix := "商户拒绝："
	if merID == 0 {
		rejectPrefix = "平台拒绝："
	}
	return s.store.WithTx(func(tx Store) error {
		ok, err := tx.UpdateRefundStatus(ctx, id, fromStatus, StatusReject, msg)
		if err != nil {
			return err
		}
		if !ok {
			return ErrBadStatus
		}
		now := time.Now()
		if err := tx.CreateRefundStatus(ctx, &RefundStatusLog{
			RefundOrderID: id,
			ChangeType:    "reject",
			ChangeMessage: rejectPrefix + msg,
			ChangeTime:    now,
		}); err != nil {
			return err
		}
		products, err := tx.ListRefundProducts(ctx, id)
		if err != nil {
			return err
		}
		for _, rp := range products {
			lines, err := tx.ListOrderProductsByIDs(ctx, ro.OrderID, []uint{rp.OrderProductID})
			if err != nil {
				return err
			}
			if len(lines) == 0 {
				continue
			}
			line := lines[0]
			isRefund := OrderProductRefundNone
			if line.RefundNum > 0 {
				if line.RefundNum >= line.ProductNum {
					isRefund = OrderProductRefundFull
				} else {
					isRefund = OrderProductRefundPartial
				}
			}
			if err := tx.MarkOrderProductRefund(ctx, rp.OrderProductID, 0, isRefund); err != nil {
				return err
			}
		}
		return nil
	})
}

// RequestPlatform 用户申请平台介入：待审(0) → 平台介入(4)。
func (s *Service) RequestPlatform(ctx context.Context, uid, id uint) error {
	ro, err := s.store.GetRefund(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if ro.UID != uid || ro.IsDel == 1 {
		return ErrForbidden
	}
	if ro.Status != StatusWait {
		return ErrBadStatus
	}
	return s.store.WithTx(func(tx Store) error {
		ok, err := tx.UpdateRefundStatus(ctx, id, StatusWait, StatusPlatform, "")
		if err != nil {
			return err
		}
		if !ok {
			return ErrBadStatus
		}
		return tx.CreateRefundStatus(ctx, &RefundStatusLog{
			RefundOrderID: id,
			ChangeType:    "platform",
			ChangeMessage: "用户申请平台介入",
			ChangeTime:    time.Now(),
		})
	})
}

// refundAuditFromStatus 商户仅审待审；平台可审待审或平台介入。
func refundAuditFromStatus(merID uint, status int8) (int8, error) {
	if merID == 0 {
		if status == StatusWait || status == StatusPlatform {
			return status, nil
		}
		return 0, ErrBadStatus
	}
	if status != StatusWait {
		return 0, ErrBadStatus
	}
	return StatusWait, nil
}

func (s *Service) ListPlatform(ctx context.Context, status *int8, page, limit int) (*PageResult[RefundOrder], error) {
	list, total, err := s.store.ListRefunds(ctx, ListFilter{Status: status, Page: page, Limit: limit})
	if err != nil {
		return nil, err
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[RefundOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetPlatform(ctx context.Context, id uint) (*RefundOrder, error) {
	return s.getAttach(ctx, id)
}

// ListPlatformByRegions 返回区域管理员可监管的退款单；nil 表示平台全量，空范围明确返回空列表。
func (s *Service) ListPlatformByRegions(ctx context.Context, regionIDs []uint, status *int8, page, limit int) (*PageResult[RefundOrder], error) {
	if regionIDs == nil {
		return s.ListPlatform(ctx, status, page, limit)
	}
	page, limit = normalizePage(page, limit)
	if len(regionIDs) == 0 {
		return &PageResult[RefundOrder]{List: []RefundOrder{}, Total: 0, Page: page, Limit: limit}, nil
	}
	list, total, err := s.store.ListPlatformRefundsByRegions(ctx, regionIDs, status, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[RefundOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

// GetPlatformByRegions 返回区域管理员可监管的退款详情。查询必须在 SQL 层同时关联商户区域，避免先读后判造成数据泄露。
func (s *Service) GetPlatformByRegions(ctx context.Context, id uint, regionIDs []uint) (*RefundOrder, error) {
	if regionIDs == nil {
		return s.GetPlatform(ctx, id)
	}
	if len(regionIDs) == 0 {
		return nil, ErrNotFound
	}
	ro, err := s.store.GetPlatformRefundByRegions(ctx, id, regionIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	products, err := s.store.ListRefundProducts(ctx, id)
	if err != nil {
		return nil, err
	}
	ro.Products = products
	return ro, nil
}

func (s *Service) getAttach(ctx context.Context, id uint) (*RefundOrder, error) {
	ro, err := s.store.GetRefund(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if ro.IsDel == 1 {
		return nil, ErrNotFound
	}
	ps, err := s.store.ListRefundProducts(ctx, id)
	if err != nil {
		return nil, err
	}
	ro.Products = ps
	return ro, nil
}

func skuUnique(p OrderProductLine) string {
	var snap struct {
		Unique string `json:"unique"`
	}
	_ = json.Unmarshal([]byte(p.CartInfo), &snap)
	if snap.Unique != "" {
		return snap.Unique
	}
	return p.ProductSKU
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

func uniqueUint(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	out := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
