package trade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/cart"
	"github.com/crmlive/qixi-live-ecrm/api-platform/internal/domain/promotion"
	"gorm.io/gorm"
)

type Store interface {
	WithTx(fn func(tx Store) error) error

	CreateGroupOrder(ctx context.Context, g *GroupOrder) error
	CreateStoreOrder(ctx context.Context, o *StoreOrder) error
	CreateOrderProduct(ctx context.Context, p *OrderProduct) error

	GetGroupOrder(ctx context.Context, id uint) (*GroupOrder, error)
	ListGroupOrders(ctx context.Context, uid uint, page, limit int) ([]GroupOrder, int64, error)
	ListStoreOrdersByGroup(ctx context.Context, groupID uint) ([]StoreOrder, error)
	ListOrderProductsByOrder(ctx context.Context, orderID uint) ([]OrderProduct, error)
	ListOrderProductsByOrders(ctx context.Context, orderIDs []uint) ([]OrderProduct, error)

	GetStoreOrder(ctx context.Context, id uint) (*StoreOrder, error)
	ListStoreOrders(ctx context.Context, merID *uint, page, limit int) ([]StoreOrder, int64, error)
	ListStoreOrdersFiltered(ctx context.Context, merID *uint, paid, status *int8, page, limit int) ([]StoreOrder, int64, error)
	ListPlatformOrdersByRegions(ctx context.Context, regionIDs []uint, paid *int8, page, limit int) ([]StoreOrder, int64, error)
	GetPlatformOrderByRegions(ctx context.Context, orderID uint, regionIDs []uint) (*StoreOrder, error)
	MerchantName(ctx context.Context, merID uint) (string, error)

	MarkGroupPaid(ctx context.Context, id uint, payType int8, payTime time.Time) error
	MarkChildrenPaid(ctx context.Context, groupID uint, payType int8, payTime time.Time) error
	MarkCartsPaid(ctx context.Context, uid uint, cartIDs []uint64) error
	DeductProductStock(ctx context.Context, productID uint, num uint) error
	DeductSKUStock(ctx context.Context, productID uint, unique string, num uint) error
	GetUserBalance(ctx context.Context, uid uint) (float64, error)
	DeductUserBalance(ctx context.Context, uid uint, amount float64) error
	GetUserIntegral(ctx context.Context, uid uint) (int, error)
	DeductUserIntegral(ctx context.Context, uid uint, amount int) error
	AddUserIntegral(ctx context.Context, uid uint, amount int) (int, error)
	MerchantsIntegralEnabled(ctx context.Context, merIDs []uint) (bool, error)
	GetUserSVIP(ctx context.Context, uid uint) (isSvip int8, end *time.Time, err error)
	MerchantsSVIPCouponMerge(ctx context.Context, merIDs []uint) (map[uint]int8, error)
	HasBill(ctx context.Context, uid uint, category, typ, linkID string) (bool, error)
	CreateBill(ctx context.Context, b *UserBill) error
	LoadPointsProduct(ctx context.Context, productID uint, unique string) (*PointsProductView, error)
	MarkCouponUsersUsed(ctx context.Context, uid uint, ids []uint, at time.Time) (int64, error)

	DeliverOrder(ctx context.Context, orderID, merID uint, name, deliveryID, deliveryType string) error
	UpdateOrderStatus(ctx context.Context, orderID uint, status int8) error
	MarkOrderVerified(ctx context.Context, orderID, merID, serviceID uint, at time.Time) error
	GetStoreOrderByVerifyCode(ctx context.Context, merID uint, code string) (*StoreOrder, error)
	ListStoreOrdersInStatuses(ctx context.Context, merID uint, paid *int8, statuses []int8, page, limit int) ([]StoreOrder, int64, error)
	SoftDeleteGroup(ctx context.Context, id, uid uint) error
	SoftDeleteOrdersByGroup(ctx context.Context, groupID uint) error
	ListExpiredUnpaidGroups(ctx context.Context, before time.Time, limit int) ([]GroupOrder, error)

	CreatePresellOrder(ctx context.Context, o *PresellOrder) error
	GetPresellOrder(ctx context.Context, id uint) (*PresellOrder, error)
	GetPresellOrderByOrderID(ctx context.Context, orderID uint) (*PresellOrder, error)
	ListPresellOrdersByUID(ctx context.Context, uid uint, unpaidOnly bool, page, limit int) ([]PresellOrder, int64, error)
	MarkPresellOrderPaid(ctx context.Context, id uint, payType int8, at time.Time) error
	InvalidatePresellOrder(ctx context.Context, id uint) error
}

// PriceQuoter 可选秒杀报价（场次内覆盖行价）。
type PriceQuoter interface {
	QuotePrice(ctx context.Context, productID uint) (price float64, activeID uint, oncePay int, ok bool, err error)
}

type Service struct {
	store   Store
	cartSvc *cart.Service
	promo   *promotion.Service
	seckill PriceQuoter
	combo   CombinationHook
	reserve ReservationHook
	presell PresellHook
	assist  AssistHook
	payment PaymentSettings
}

func NewService(store Store, cartSvc *cart.Service, promo *promotion.Service) *Service {
	return &Service{store: store, cartSvc: cartSvc, promo: promo}
}

func (s *Service) SetSeckill(q PriceQuoter) { s.seckill = q }

func (s *Service) applySeckillPrices(ctx context.Context, rows []cart.Cart) error {
	if s.seckill == nil {
		return nil
	}
	for i := range rows {
		price, activeID, oncePay, ok, err := s.seckill.QuotePrice(ctx, rows[i].ProductID)
		if err != nil {
			return err
		}
		if !ok || price <= 0 || activeID == 0 {
			continue
		}
		if oncePay > 0 && int(rows[i].CartNum) > oncePay {
			return ErrSeckillLimit
		}
		rows[i].Price = price
		rows[i].SeckillActiveID = activeID
		rows[i].OncePayCount = oncePay
		rows[i].ProductType = ActivityTypeSeckill
	}
	return nil
}

func cartHasSeckill(rows []cart.Cart) bool {
	for _, r := range rows {
		if r.SeckillActiveID > 0 {
			return true
		}
	}
	return false
}

func (s *Service) V2Check(ctx context.Context, uid uint, in CheckInput) (*CheckResult, error) {
	rows, err := s.cartSvc.LoadForCheckout(ctx, uid, in.CartIDs)
	if err != nil {
		return nil, err
	}
	if err := rejectPointsInNormal(rows); err != nil {
		return nil, err
	}
	if err := s.applySeckillPrices(ctx, rows); err != nil {
		return nil, err
	}
	usedSvip, skipStoreSvip, err := s.applySvipPrices(ctx, uid, rows)
	if err != nil {
		return nil, err
	}
	hasSeckill := cartHasSeckill(rows)
	check := buildCheck(rows)
	check.UsedSvip = usedSvip
	check.SvipDiscount = sumSvipDiscount(rows)
	// 店铺券：秒杀禁用。平台券：仅普通。预售走专用 /order/presell/*，不经购物车。
	if err := s.applyCoupons(ctx, uid, check, in.NormalizedCouponUserIDs(), hasSeckill || skipStoreSvip, hasSeckill); err != nil {
		return nil, err
	}
	if err := s.applyIntegral(ctx, uid, check, in.UseIntegral, cartHasActivityGoods(rows)); err != nil {
		return nil, err
	}
	return check, nil
}

func (s *Service) V2Create(ctx context.Context, uid uint, in CreateInput) (*GroupOrder, error) {
	if len(in.CartIDs) == 0 {
		return nil, ErrEmptyCart
	}
	if in.AddressID == 0 {
		return nil, ErrAddressRequired
	}
	addr, err := s.cartSvc.GetAddress(ctx, uid, in.AddressID)
	if err != nil {
		return nil, err
	}
	rows, err := s.cartSvc.LoadForCheckout(ctx, uid, in.CartIDs)
	if err != nil {
		return nil, err
	}
	if err := rejectPointsInNormal(rows); err != nil {
		return nil, err
	}
	if err := s.applySeckillPrices(ctx, rows); err != nil {
		return nil, err
	}
	usedSvip, skipStoreSvip, err := s.applySvipPrices(ctx, uid, rows)
	if err != nil {
		return nil, err
	}
	check := buildCheck(rows)
	hasSeckill := cartHasSeckill(rows)
	check.UsedSvip = usedSvip
	check.SvipDiscount = sumSvipDiscount(rows)
	if err := s.applyCoupons(ctx, uid, check, in.NormalizedCouponUserIDs(), hasSeckill || skipStoreSvip, hasSeckill); err != nil {
		return nil, err
	}
	if err := s.applyIntegral(ctx, uid, check, in.UseIntegral, cartHasActivityGoods(rows)); err != nil {
		return nil, err
	}
	fullAddr := addr.FullAddress()
	mark := strings.TrimSpace(in.Mark)
	// 仅核销实际生效的券（秒杀会跳过店铺券）
	useIDs := make([]uint, 0)
	if check.PlatformCouponUserID > 0 {
		useIDs = append(useIDs, check.PlatformCouponUserID)
	}
	for _, id := range check.merCouponUserIDs {
		if id > 0 {
			useIDs = append(useIDs, id)
		}
	}
	useIDs = uniqueUintIDs(useIDs)
	var actType int8
	if hasSeckill {
		actType = ActivityTypeSeckill
	}

	var created *GroupOrder
	err = s.store.WithTx(func(tx Store) error {
		g := &GroupOrder{
			GroupOrderSN:  genSN("G"),
			UID:           uid,
			TotalPostage:  0,
			TotalPrice:    check.TotalPrice,
			TotalNum:      check.TotalNum,
			CouponPrice:   check.CouponPrice,
			CouponID:      check.PlatformCouponUserID,
			Integral:      check.Integral,
			IntegralPrice: check.IntegralPrice,
			GiveIntegral:  check.GiveIntegral,
			RealName:      addr.RealName,
			UserPhone:     addr.Phone,
			UserAddress:   fullAddr,
			PayPrice:      check.PayPrice,
			PayPostage:    0,
			Cost:          sumCost(rows),
			Paid:          0,
			ActivityType:  actType,
			CreateTime:    time.Now(),
		}
		if err := tx.CreateGroupOrder(ctx, g); err != nil {
			return err
		}
		byMer := map[uint][]cart.Cart{}
		merOrder := make([]uint, 0)
		for _, r := range rows {
			if _, ok := byMer[r.MerID]; !ok {
				merOrder = append(merOrder, r.MerID)
			}
			byMer[r.MerID] = append(byMer[r.MerID], r)
		}
		orders := make([]StoreOrder, 0, len(merOrder))
		for _, merID := range merOrder {
			items := byMer[merID]
			var totalPrice float64
			var totalNum int
			var cost float64
			cartParts := make([]string, 0, len(items))
			for _, it := range items {
				totalPrice += it.Price * float64(it.CartNum)
				totalNum += int(it.CartNum)
				cost += it.Cost * float64(it.CartNum)
				cartParts = append(cartParts, strconv.FormatUint(it.CartID, 10))
			}
			var merCheck *CheckMerchant
			for i := range check.Merchants {
				if check.Merchants[i].MerID == merID {
					merCheck = &check.Merchants[i]
					break
				}
			}
			storeDisc, platShare := 0.0, 0.0
			var couponID uint
			payPrice := totalPrice
			merIntegral, merIntegralPrice, merGive := 0, 0.0, 0
			if merCheck != nil {
				storeDisc = merCheck.CouponPrice
				platShare = merCheck.PlatformCouponPrice
				payPrice = merCheck.PayPrice
				couponID = merCouponUserID(check, merID)
				merIntegral = merCheck.Integral
				merIntegralPrice = merCheck.IntegralPrice
				merGive = merCheck.GiveIntegral
			}
			merAct := int8(0)
			for _, it := range items {
				if it.SeckillActiveID > 0 {
					merAct = ActivityTypeSeckill
					break
				}
			}
			o := &StoreOrder{
				GroupOrderID:        g.GroupOrderID,
				OrderSN:             genSN("O"),
				UID:                 uid,
				RealName:            addr.RealName,
				UserPhone:           addr.Phone,
				UserAddress:         fullAddr,
				CartID:              strings.Join(cartParts, ","),
				TotalNum:            totalNum,
				TotalPrice:          totalPrice,
				PayPrice:            payPrice,
				Paid:                0,
				Status:              OrderStatusAwaitShip,
				Mark:                mark,
				MerID:               merID,
				Cost:                cost,
				CouponID:            strconv.FormatUint(uint64(couponID), 10),
				CouponPrice:         storeDisc,
				PlatformCouponPrice: platShare,
				SvipDiscount:        merSvipDiscount(items, merID),
				Integral:            merIntegral,
				IntegralPrice:       merIntegralPrice,
				GiveIntegral:        merGive,
				VerifyCode:          ensureVerifyCode(""),
				ActivityType:        merAct,
				CreateTime:          time.Now(),
			}
			if err := tx.CreateStoreOrder(ctx, o); err != nil {
				return err
			}
			for _, it := range items {
				info, _ := json.Marshal(map[string]interface{}{
					"product_id":        it.ProductID,
					"store_name":        it.StoreName,
					"image":             it.Image,
					"price":             it.Price,
					"unique":            it.ProductAttrUnique,
					"seckill_active_id": it.SeckillActiveID,
				})
				op := &OrderProduct{
					OrderID:      o.OrderID,
					UID:          uid,
					CartID:       int(it.CartID),
					ProductID:    it.ProductID,
					ProductSKU:   trimSKU(it.ProductAttrUnique),
					ProductNum:   int(it.CartNum),
					ProductType:  it.ProductType,
					ActivityID:   it.SeckillActiveID,
					Cost:         it.Cost,
					ProductPrice: it.Price,
					TotalPrice:   it.Price * float64(it.CartNum),
					ProductInfo:  string(info),
				}
				if err := tx.CreateOrderProduct(ctx, op); err != nil {
					return err
				}
				o.Products = append(o.Products, *op)
			}
			if name, err := tx.MerchantName(ctx, merID); err == nil {
				o.MerName = name
			}
			orders = append(orders, *o)
		}
		if len(useIDs) > 0 {
			n, err := tx.MarkCouponUsersUsed(ctx, uid, useIDs, time.Now())
			if err != nil {
				return err
			}
			if int(n) != len(useIDs) {
				return ErrCoupon
			}
		}
		if check.Integral > 0 {
			if err := tx.DeductUserIntegral(ctx, uid, check.Integral); err != nil {
				return err
			}
			bal, err := tx.GetUserIntegral(ctx, uid)
			if err != nil {
				return err
			}
			if err := tx.CreateBill(ctx, &UserBill{
				UID:      uid,
				LinkID:   fmt.Sprintf("%d", g.GroupOrderID),
				PM:       BillPMOut,
				Title:    "积分抵扣",
				Category: "integral",
				Type:     "deduction",
				Number:   float64(check.Integral),
				Balance:  float64(bal),
				Mark:     fmt.Sprintf("订单%s积分抵扣", g.GroupOrderSN),
				Status:   1,
			}); err != nil {
				return err
			}
		}
		g.Orders = orders
		created = g
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 创建成功后购物车仍保持 is_pay=0，直至 PaySuccess
	return created, nil
}

func (s *Service) applyCoupons(ctx context.Context, uid uint, check *CheckResult, couponUserIDs []uint, skipStoreCoupon, skipPlatformCoupon bool) error {
	if s.promo == nil || len(couponUserIDs) == 0 {
		return nil
	}
	mers := make([]promotion.MerTotal, 0, len(check.Merchants))
	for _, m := range check.Merchants {
		mers = append(mers, promotion.MerTotal{MerID: m.MerID, TotalPrice: m.TotalPrice})
	}
	q, err := s.promo.Quote(ctx, uid, promotion.QuoteInput{
		MerTotals: mers, CouponUserIDs: couponUserIDs,
		SkipStoreCoupon: skipStoreCoupon, SkipPlatformCoupon: skipPlatformCoupon,
	})
	if err != nil {
		return err
	}
	check.CouponPrice = q.CouponPrice
	check.PayPrice = q.PayPrice
	check.PlatformCouponUserID = q.PlatformCouponUserID
	for i := range check.Merchants {
		merID := check.Merchants[i].MerID
		store := q.MerStoreDiscount[merID]
		plat := q.MerPlatformShare[merID]
		check.Merchants[i].CouponPrice = store
		check.Merchants[i].PlatformCouponPrice = plat
		pay := check.Merchants[i].TotalPrice - store - plat
		if pay < 0 {
			pay = 0
		}
		check.Merchants[i].PayPrice = pay
	}
	check.merCouponUserIDs = q.MerCouponUserID
	return nil
}

func (s *Service) applyIntegral(ctx context.Context, uid uint, check *CheckResult, useIntegral int, hasActivity bool) error {
	bal, err := s.store.GetUserIntegral(ctx, uid)
	if err != nil {
		return err
	}
	check.UserIntegral = bal
	if useIntegral <= 0 {
		check.GiveIntegral = int(check.PayPrice * GiveIntegralRatio)
		return nil
	}
	if hasActivity {
		return ErrIntegralOnActivity
	}
	merIDs := make([]uint, 0, len(check.Merchants))
	for _, m := range check.Merchants {
		merIDs = append(merIDs, m.MerID)
	}
	ok, err := s.store.MerchantsIntegralEnabled(ctx, merIDs)
	if err != nil {
		return err
	}
	if !ok {
		return ErrMerIntegralOff
	}
	// 券后应付作为 TotalPrice，避免 ApplyPricing 重算优惠券
	before := make(map[uint]float64, len(check.Merchants))
	mers := make([]MerAmount, len(check.Merchants))
	for i, m := range check.Merchants {
		before[m.MerID] = m.PayPrice
		mers[i] = MerAmount{MerID: m.MerID, TotalPrice: m.PayPrice, PayPrice: m.PayPrice}
	}
	priced := ApplyPricing(PricingInput{
		MerAmounts:   mers,
		UseIntegral:  useIntegral,
		UserIntegral: bal,
	})
	check.Integral = priced.Integral
	check.IntegralPrice = priced.IntegralPrice
	check.PayPrice = priced.PayPrice
	check.GiveIntegral = priced.GiveIntegral
	leftIntegral := priced.Integral
	leftGive := priced.GiveIntegral
	for i := range check.Merchants {
		merID := check.Merchants[i].MerID
		var pm *MerAmount
		for j := range priced.Merchants {
			if priced.Merchants[j].MerID == merID {
				pm = &priced.Merchants[j]
				break
			}
		}
		if pm == nil {
			continue
		}
		ip := round2(before[merID] - pm.PayPrice)
		check.Merchants[i].IntegralPrice = ip
		check.Merchants[i].PayPrice = pm.PayPrice
		if i == len(check.Merchants)-1 {
			check.Merchants[i].Integral = leftIntegral
			check.Merchants[i].GiveIntegral = leftGive
			continue
		}
		share := 0
		if priced.IntegralPrice > 0 && ip > 0 {
			share = int(float64(priced.Integral) * (ip / priced.IntegralPrice))
		}
		if share > leftIntegral {
			share = leftIntegral
		}
		check.Merchants[i].Integral = share
		leftIntegral -= share
		gShare := 0
		if check.PayPrice > 0 {
			gShare = int(float64(priced.GiveIntegral) * (pm.PayPrice / check.PayPrice))
		}
		if gShare > leftGive {
			gShare = leftGive
		}
		check.Merchants[i].GiveIntegral = gShare
		leftGive -= gShare
	}
	return nil
}

func rejectPointsInNormal(rows []cart.Cart) error {
	for _, r := range rows {
		if r.GoodsType == GoodsTypePoints || int8(r.ProductType) == ActivityTypePoints {
			return ErrPointsProductMix
		}
	}
	return nil
}

func merCouponUserID(check *CheckResult, merID uint) uint {
	if check == nil || check.merCouponUserIDs == nil {
		return 0
	}
	return check.merCouponUserIDs[merID]
}

func uniqueUintIDs(ids []uint) []uint {
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

// PaySuccess mock/balance 支付成功；已支付则幂等返回。
func (s *Service) PaySuccess(ctx context.Context, uid uint, groupOrderID uint, payTypeStr string) (*GroupOrder, error) {
	payType, err := parsePayType(payTypeStr)
	if err != nil {
		return nil, err
	}
	// mock 只能用于明确开启的本地/测试闭环，绝不能成为线上用户可调用的成功路径。
	// 微信/支付宝由 CreatePayIntent + 受验签的 NotifyChannelPay 完成，不能直接调用此方法。
	if (payType == PayTypeMock && !s.payment.Sandbox) || payType == PayTypeWechat || payType == PayTypeAlipay {
		return nil, ErrPaymentConfig
	}
	g, err := s.store.GetGroupOrder(ctx, groupOrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if g.UID != uid || g.IsDel == 1 {
		return nil, ErrForbidden
	}
	if g.Paid == 1 {
		return s.attachGroup(ctx, g)
	}

	children, err := s.store.ListStoreOrdersByGroup(ctx, groupOrderID)
	if err != nil {
		return nil, err
	}
	orderIDs := make([]uint, 0, len(children))
	for _, c := range children {
		orderIDs = append(orderIDs, c.OrderID)
	}
	products, err := s.store.ListOrderProductsByOrders(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	cartIDs := make([]uint64, 0, len(products))
	for _, p := range products {
		if p.CartID > 0 {
			cartIDs = append(cartIDs, uint64(p.CartID))
		}
	}

	err = s.store.WithTx(func(tx Store) error {
		cur, err := tx.GetGroupOrder(ctx, groupOrderID)
		if err != nil {
			return err
		}
		if cur.Paid == 1 {
			return nil
		}
		if cur.ActivityType == ActivityTypePoints {
			return ErrNotPointsProduct // 积分商城单走 /order/points/pay/:id
		}
		if payType == PayTypeBalance {
			bal, err := tx.GetUserBalance(ctx, uid)
			if err != nil {
				return err
			}
			if bal < cur.PayPrice {
				return ErrBalanceNotEnough
			}
			if err := tx.DeductUserBalance(ctx, uid, cur.PayPrice); err != nil {
				return err
			}
		}
		if payType == PayTypeIntegral {
			return ErrInvalidPayType
		}
		// 普通单积分抵扣已在 V2Create 落库并扣减，此处不再扣
		now := time.Now()
		if err := tx.MarkGroupPaid(ctx, groupOrderID, payType, now); err != nil {
			return err
		}
		if err := tx.MarkChildrenPaid(ctx, groupOrderID, payType, now); err != nil {
			return err
		}
		for _, p := range products {
			n := uint(p.ProductNum)
			if n == 0 {
				continue
			}
			if err := tx.DeductProductStock(ctx, p.ProductID, n); err != nil {
				return err
			}
			if err := tx.DeductSKUStock(ctx, p.ProductID, skuUnique(p), n); err != nil {
				return err
			}
		}
		if len(cartIDs) > 0 {
			if err := tx.MarkCartsPaid(ctx, uid, cartIDs); err != nil {
				return err
			}
		}
		if cur.GiveIntegral > 0 {
			if err := creditGiveIntegral(ctx, tx, uid, groupOrderID, cur.GiveIntegral); err != nil {
				return err
			}
		}
		if cur.ActivityType == ActivityTypeCombination && s.combo != nil {
			for _, c := range children {
				ok, orderIDs, err := s.combo.OnOrderPaid(ctx, c.OrderID)
				if err != nil {
					return err
				}
				if ok {
					for _, oid := range orderIDs {
						if err := tx.UpdateOrderStatus(ctx, oid, OrderStatusAwaitShip); err != nil {
							return err
						}
					}
				}
			}
		}
		if cur.ActivityType == ActivityTypePresell && s.presell != nil {
			for _, p := range products {
				if p.ProductType != ActivityTypePresell || p.ActivityID == 0 {
					continue
				}
				if err := s.presell.OnOrderPaid(ctx, p.ActivityID, p.ProductNum); err != nil {
					return err
				}
			}
			for _, c := range children {
				po, err := tx.GetPresellOrderByOrderID(ctx, c.OrderID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						continue
					}
					return err
				}
				if po != nil && po.Paid == 0 && po.Status == 1 {
					if err := tx.UpdateOrderStatus(ctx, c.OrderID, OrderStatusAwaitFinal); err != nil {
						return err
					}
				}
			}
		}
		if cur.ActivityType == ActivityTypeAssist && s.assist != nil {
			for _, p := range products {
				if p.ProductType != ActivityTypeAssist || p.ActivityID == 0 {
					continue
				}
				if err := s.assist.MarkSetPaid(ctx, p.ActivityID); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	g, err = s.store.GetGroupOrder(ctx, groupOrderID)
	if err != nil {
		return nil, err
	}
	if s.promo != nil && g != nil && g.ActivityType != ActivityTypePoints {
		_ = s.promo.CreditSpreadOnPay(ctx, uid, g.GroupOrderID, g.PayPrice)
	}
	return s.attachGroup(ctx, g)
}

func creditGiveIntegral(ctx context.Context, tx Store, uid, groupOrderID uint, amount int) error {
	link := fmt.Sprintf("give:%d", groupOrderID)
	dup, err := tx.HasBill(ctx, uid, "integral", "gain", link)
	if err != nil {
		return err
	}
	if dup {
		return nil
	}
	bal, err := tx.AddUserIntegral(ctx, uid, amount)
	if err != nil {
		return err
	}
	return tx.CreateBill(ctx, &UserBill{
		UID:      uid,
		LinkID:   link,
		PM:       BillPMIn,
		Title:    "下单赠送积分",
		Category: "integral",
		Type:     "gain",
		Number:   float64(amount),
		Balance:  float64(bal),
		Mark:     fmt.Sprintf("订单%d支付赠送", groupOrderID),
		Status:   1,
	})
}

// UserIntegral 查询用户积分余额。
func (s *Service) UserIntegral(ctx context.Context, uid uint) (int, error) {
	return s.store.GetUserIntegral(ctx, uid)
}

func (s *Service) ListGroupOrders(ctx context.Context, uid uint, page, limit int) (*PageResult[GroupOrder], error) {
	list, total, err := s.store.ListGroupOrders(ctx, uid, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		orders, err := s.store.ListStoreOrdersByGroup(ctx, list[i].GroupOrderID)
		if err != nil {
			return nil, err
		}
		list[i].Orders = orders
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[GroupOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetGroupOrder(ctx context.Context, uid, id uint) (*GroupOrder, error) {
	g, err := s.store.GetGroupOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if g.UID != uid || g.IsDel == 1 {
		return nil, ErrForbidden
	}
	return s.attachGroup(ctx, g)
}

func (s *Service) GetStoreOrderForUser(ctx context.Context, uid, id uint) (*StoreOrder, error) {
	o, err := s.store.GetStoreOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if o.UID != uid || o.IsDel == 1 {
		return nil, ErrForbidden
	}
	return s.attachOrder(ctx, o)
}

func (s *Service) ListMerchantOrders(ctx context.Context, merID uint, page, limit int) (*PageResult[StoreOrder], error) {
	return s.MerchantList(ctx, merID, nil, nil, page, limit)
}

func (s *Service) MerchantList(ctx context.Context, merID uint, paid, status *int8, page, limit int) (*PageResult[StoreOrder], error) {
	list, total, err := s.store.ListStoreOrdersFiltered(ctx, &merID, paid, status, page, limit)
	if err != nil {
		return nil, err
	}
	return s.pageOrders(ctx, list, total, page, limit)
}

// MerchantListAwaitVerify 待核销：paid=1 且 status IN (0,1,2)
func (s *Service) MerchantListAwaitVerify(ctx context.Context, merID uint, page, limit int) (*PageResult[StoreOrder], error) {
	paid := int8(1)
	list, total, err := s.store.ListStoreOrdersInStatuses(ctx, merID, &paid, []int8{0, 1, 2}, page, limit)
	if err != nil {
		return nil, err
	}
	return s.pageOrders(ctx, list, total, page, limit)
}

func (s *Service) pageOrders(ctx context.Context, list []StoreOrder, total int64, page, limit int) (*PageResult[StoreOrder], error) {
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
		prods, _ := s.store.ListOrderProductsByOrder(ctx, list[i].OrderID)
		list[i].Products = prods
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[StoreOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) MerchantGet(ctx context.Context, merID, id uint) (*StoreOrder, error) {
	return s.GetMerchantOrder(ctx, merID, id)
}

func (s *Service) GetMerchantOrder(ctx context.Context, merID, id uint) (*StoreOrder, error) {
	o, err := s.store.GetStoreOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if o.MerID != merID || o.IsDel == 1 {
		return nil, ErrForbidden
	}
	return s.attachOrder(ctx, o)
}

func (s *Service) Deliver(ctx context.Context, merID, orderID uint, in DeliveryInput) error {
	name := strings.TrimSpace(in.DeliveryName)
	did := strings.TrimSpace(in.DeliveryID)
	if name == "" || did == "" {
		return ErrDeliveryParam
	}
	o, err := s.GetMerchantOrder(ctx, merID, orderID)
	if err != nil {
		return err
	}
	if o.Paid != 1 || o.Status != OrderStatusAwaitShip {
		return ErrBadStatus
	}
	typ := strings.TrimSpace(in.DeliveryType)
	if typ == "" {
		typ = "express"
	}
	return s.store.DeliverOrder(ctx, orderID, merID, name, did, typ)
}

// Verify 核销订单。若订单有 verify_code 则必须匹配；空码则仅校验 mer_id+paid。
// 已 status=3 幂等成功；写入 verify_time 与 status=3。
func (s *Service) Verify(ctx context.Context, merID, orderID uint, verifyCode string) error {
	return s.VerifyByStaff(ctx, merID, orderID, 0, verifyCode)
}

// VerifyByStaff 同 Verify，额外记录核销店员 service_id。
func (s *Service) VerifyByStaff(ctx context.Context, merID, orderID, serviceID uint, verifyCode string) error {
	o, err := s.GetMerchantOrder(ctx, merID, orderID)
	if err != nil {
		return err
	}
	if o.Paid != 1 {
		return ErrNotPaid
	}
	if o.Status == OrderStatusDone {
		return nil
	}
	want := strings.TrimSpace(o.VerifyCode)
	if want != "" && strings.TrimSpace(verifyCode) != want {
		return ErrVerifyCodeMismatch
	}
	return s.store.MarkOrderVerified(ctx, orderID, merID, serviceID, time.Now())
}

// GetStoreOrderByVerifyCode 按核销码查本店订单。
func (s *Service) GetStoreOrderByVerifyCode(ctx context.Context, merID uint, code string) (*StoreOrder, error) {
	code = strings.TrimSpace(code)
	if merID == 0 || code == "" {
		return nil, ErrBadParam
	}
	o, err := s.store.GetStoreOrderByVerifyCode(ctx, merID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.attachOrder(ctx, o)
}

// VerifyByCode 店员核销码核销（mer_id 隔离）。
func (s *Service) VerifyByCode(ctx context.Context, merID uint, code string) (*StoreOrder, error) {
	code = strings.TrimSpace(code)
	if merID == 0 || code == "" {
		return nil, ErrBadParam
	}
	o, err := s.store.GetStoreOrderByVerifyCode(ctx, merID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := s.VerifyByStaff(ctx, merID, o.OrderID, 0, code); err != nil {
		return nil, err
	}
	return s.GetMerchantOrder(ctx, merID, o.OrderID)
}

// VerifyByCodeStaff 店员核销并记录 service_id。
func (s *Service) VerifyByCodeStaff(ctx context.Context, merID, serviceID uint, code string) (*StoreOrder, error) {
	code = strings.TrimSpace(code)
	if merID == 0 || code == "" {
		return nil, ErrBadParam
	}
	o, err := s.store.GetStoreOrderByVerifyCode(ctx, merID, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := s.VerifyByStaff(ctx, merID, o.OrderID, serviceID, code); err != nil {
		return nil, err
	}
	return s.GetMerchantOrder(ctx, merID, o.OrderID)
}

func (s *Service) ConfirmReceive(ctx context.Context, uid, orderID uint) error {
	o, err := s.GetStoreOrderForUser(ctx, uid, orderID)
	if err != nil {
		return err
	}
	if o.Paid != 1 || o.Status != OrderStatusShipped {
		return ErrBadStatus
	}
	return s.store.UpdateOrderStatus(ctx, orderID, OrderStatusDone)
}

func (s *Service) CancelGroup(ctx context.Context, uid, id uint) error {
	g, err := s.GetGroupOrder(ctx, uid, id)
	if err != nil {
		return err
	}
	if g.Paid == 1 {
		return ErrAlreadyPaid
	}
	return s.cancelUnpaidGroup(ctx, g)
}

// CloseExpiredUnpaid 关闭超时未支付主单（job 调用）；复用取消路径退回抵扣积分。
func (s *Service) CloseExpiredUnpaid(ctx context.Context, olderThan time.Duration, limit int) (int, error) {
	if olderThan <= 0 {
		olderThan = 30 * time.Minute
	}
	if limit <= 0 {
		limit = 50
	}
	list, err := s.store.ListExpiredUnpaidGroups(ctx, time.Now().Add(-olderThan), limit)
	if err != nil {
		return 0, err
	}
	n := 0
	for i := range list {
		g := list[i]
		if err := s.cancelUnpaidGroup(ctx, &g); err != nil {
			continue
		}
		n++
	}
	return n, nil
}

func (s *Service) cancelUnpaidGroup(ctx context.Context, g *GroupOrder) error {
	if g == nil || g.Paid == 1 {
		return ErrAlreadyPaid
	}
	id, uid := g.GroupOrderID, g.UID
	var restore []OrderProduct
	var comboOrderIDs []uint
	err := s.store.WithTx(func(tx Store) error {
		orders, err := tx.ListStoreOrdersByGroup(ctx, id)
		if err != nil {
			return err
		}
		for i := range orders {
			products, err := tx.ListOrderProductsByOrder(ctx, orders[i].OrderID)
			if err != nil {
				return err
			}
			restore = append(restore, products...)
			// 定金预售未付定金时已建尾款单：随主单取消一并作废
			if g.ActivityType == ActivityTypePresell {
				po, err := tx.GetPresellOrderByOrderID(ctx, orders[i].OrderID)
				if err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return err
					}
				} else if po.Paid == 0 && po.Status == 1 {
					if err := tx.InvalidatePresellOrder(ctx, po.PresellOrderID); err != nil {
						return err
					}
				}
			}
			// 拼团席位在事务成功后再释放（组合域独立会话，避免关单失败却丢席位）
			if orders[i].ActivityType == ActivityTypeCombination {
				comboOrderIDs = append(comboOrderIDs, orders[i].OrderID)
			}
		}
		// 先软删主/子单，避免关单失败重试时活动库存被重复 IncStock
		if err := tx.SoftDeleteGroup(ctx, id, uid); err != nil {
			return err
		}
		if err := tx.SoftDeleteOrdersByGroup(ctx, id); err != nil {
			return err
		}
		// 普通单创建时已扣抵扣积分：取消未支付单退回
		if g.ActivityType != ActivityTypePoints && g.Integral > 0 {
			link := fmt.Sprintf("refund:%d", id)
			dup, err := tx.HasBill(ctx, uid, "integral", "gain", link)
			if err != nil {
				return err
			}
			if !dup {
				bal, err := tx.AddUserIntegral(ctx, uid, g.Integral)
				if err != nil {
					return err
				}
				if err := tx.CreateBill(ctx, &UserBill{
					UID: uid, LinkID: link, PM: BillPMIn,
					Title: "取消订单退回积分", Category: "integral", Type: "gain",
					Number: float64(g.Integral), Balance: float64(bal),
					Mark: fmt.Sprintf("取消订单%d退回抵扣积分", id), Status: 1,
				}); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	for i := range restore {
		if err := s.restoreActivityStock(ctx, &restore[i]); err != nil {
			return fmt.Errorf("restore activity stock order_product=%d: %w", restore[i].OrderProductID, err)
		}
	}
	if s.combo != nil {
		for _, oid := range comboOrderIDs {
			if err := s.combo.CancelUnpaid(ctx, oid); err != nil {
				return fmt.Errorf("release combination seat order=%d: %w", oid, err)
			}
		}
	}
	return nil
}

// restoreActivityStock 归还创建时预扣的活动库存（预售/助力）；普通 v2 在支付时扣库存，此处跳过。
func (s *Service) restoreActivityStock(ctx context.Context, p *OrderProduct) error {
	if p == nil {
		return nil
	}
	n := p.ProductNum
	if n <= 0 {
		n = 1
	}
	switch p.ProductType {
	case ActivityTypePresell:
		if s.presell == nil || p.ActivityID == 0 {
			return nil
		}
		return s.presell.RestoreStock(ctx, p.ActivityID, n)
	case ActivityTypeAssist:
		if s.assist == nil {
			return nil
		}
		assistID := assistIDFromProductInfo(p.ProductInfo)
		if assistID == 0 {
			return nil
		}
		return s.assist.RestoreStock(ctx, assistID, n)
	default:
		return nil
	}
}

func assistIDFromProductInfo(raw string) uint {
	if raw == "" {
		return 0
	}
	var info struct {
		ProductAssistID uint `json:"product_assist_id"`
	}
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		return 0
	}
	return info.ProductAssistID
}

// PayGroup / ListUserGroups / GetGroupDetail：handler 友好别名
func (s *Service) PayGroup(ctx context.Context, uid, id uint, payType string) (*GroupOrder, error) {
	return s.PaySuccess(ctx, uid, id, payType)
}

func (s *Service) ListUserGroups(ctx context.Context, uid uint, paid *int8, page, limit int) (*PageResult[GroupOrder], error) {
	_ = paid // P0 列表不过滤 paid，前端可本地筛
	return s.ListGroupOrders(ctx, uid, page, limit)
}

func (s *Service) GetGroupDetail(ctx context.Context, uid, id uint) (*GroupOrder, error) {
	return s.GetGroupOrder(ctx, uid, id)
}

func (s *Service) ListPlatformOrders(ctx context.Context, page, limit int) (*PageResult[StoreOrder], error) {
	return s.PlatformList(ctx, nil, page, limit)
}

func (s *Service) PlatformList(ctx context.Context, paid *int8, page, limit int) (*PageResult[StoreOrder], error) {
	list, total, err := s.store.ListStoreOrdersFiltered(ctx, nil, paid, nil, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[StoreOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

// PlatformListByRegions 返回区域管理员可监管的订单；nil 表示平台全量，空范围明确返回空列表。
func (s *Service) PlatformListByRegions(ctx context.Context, paid *int8, regionIDs []uint, page, limit int) (*PageResult[StoreOrder], error) {
	if regionIDs == nil {
		return s.PlatformList(ctx, paid, page, limit)
	}
	if len(regionIDs) == 0 {
		page, limit = normalizePage(page, limit)
		return &PageResult[StoreOrder]{List: []StoreOrder{}, Total: 0, Page: page, Limit: limit}, nil
	}
	list, total, err := s.store.ListPlatformOrdersByRegions(ctx, regionIDs, paid, page, limit)
	if err != nil {
		return nil, err
	}
	for i := range list {
		if name, err := s.store.MerchantName(ctx, list[i].MerID); err == nil {
			list[i].MerName = name
		}
	}
	page, limit = normalizePage(page, limit)
	return &PageResult[StoreOrder]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) PlatformGet(ctx context.Context, id uint) (*StoreOrder, error) {
	return s.GetPlatformOrder(ctx, id)
}

func (s *Service) PlatformGetByRegions(ctx context.Context, id uint, regionIDs []uint) (*StoreOrder, error) {
	if regionIDs == nil {
		return s.PlatformGet(ctx, id)
	}
	if len(regionIDs) == 0 {
		return nil, ErrNotFound
	}
	o, err := s.store.GetPlatformOrderByRegions(ctx, id, regionIDs)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.attachOrder(ctx, o)
}

func (s *Service) GetPlatformOrder(ctx context.Context, id uint) (*StoreOrder, error) {
	o, err := s.store.GetStoreOrder(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if o.IsDel == 1 {
		return nil, ErrNotFound
	}
	return s.attachOrder(ctx, o)
}

func (s *Service) attachGroup(ctx context.Context, g *GroupOrder) (*GroupOrder, error) {
	orders, err := s.store.ListStoreOrdersByGroup(ctx, g.GroupOrderID)
	if err != nil {
		return nil, err
	}
	for i := range orders {
		full, err := s.attachOrder(ctx, &orders[i])
		if err != nil {
			return nil, err
		}
		orders[i] = *full
	}
	g.Orders = orders
	return g, nil
}

func (s *Service) attachOrder(ctx context.Context, o *StoreOrder) (*StoreOrder, error) {
	if name, err := s.store.MerchantName(ctx, o.MerID); err == nil {
		o.MerName = name
	}
	ps, err := s.store.ListOrderProductsByOrder(ctx, o.OrderID)
	if err != nil {
		return nil, err
	}
	o.Products = ps
	return o, nil
}

func buildCheck(rows []cart.Cart) *CheckResult {
	byMer := map[uint]*CheckMerchant{}
	order := make([]uint, 0)
	var totalPrice float64
	var totalNum int
	for _, r := range rows {
		m := byMer[r.MerID]
		if m == nil {
			m = &CheckMerchant{MerID: r.MerID, MerName: r.MerName, Postage: 0, Items: []CheckItem{}}
			byMer[r.MerID] = m
			order = append(order, r.MerID)
		}
		sub := r.Price * float64(r.CartNum)
		m.Items = append(m.Items, CheckItem{
			CartID: r.CartID, ProductID: r.ProductID, ProductAttrUnique: r.ProductAttrUnique,
			StoreName: r.StoreName, Image: r.Image, Price: r.Price, CartNum: r.CartNum, Subtotal: sub,
		})
		m.TotalPrice += sub
		m.TotalNum += int(r.CartNum)
		totalPrice += sub
		totalNum += int(r.CartNum)
	}
	merchants := make([]CheckMerchant, 0, len(order))
	for _, id := range order {
		merchants = append(merchants, *byMer[id])
	}
	for i := range merchants {
		merchants[i].PayPrice = merchants[i].TotalPrice
	}
	return &CheckResult{
		Merchants: merchants, TotalPrice: totalPrice, TotalPostage: 0,
		CouponPrice: 0, PayPrice: totalPrice, TotalNum: totalNum,
	}
}

func parsePayType(s string) (int8, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "balance":
		return PayTypeBalance, nil
	case "wechat":
		return PayTypeWechat, nil
	case "alipay":
		return PayTypeAlipay, nil
	case "mock":
		return PayTypeMock, nil
	case "integral":
		return PayTypeIntegral, nil
	default:
		return 0, ErrInvalidPayType
	}
}

func genSN(prefix string) string {
	return fmt.Sprintf("%s%s%04d", prefix, time.Now().Format("20060102150405"), rand.Intn(10000))
}

// genVerifyCode 生成 8~12 位数字核销码（唯一性尽力，无强约束）。
func genVerifyCode() string {
	n := 8 + rand.Intn(5) // 8..12
	var b strings.Builder
	b.Grow(n)
	b.WriteByte(byte('1' + rand.Intn(9))) // 首位非 0
	for i := 1; i < n; i++ {
		b.WriteByte(byte('0' + rand.Intn(10)))
	}
	return b.String()
}

// ensureVerifyCode 创建订单时若码为空则生成。
func ensureVerifyCode(code string) string {
	if strings.TrimSpace(code) == "" {
		return genVerifyCode()
	}
	return strings.TrimSpace(code)
}

func sumCost(rows []cart.Cart) float64 {
	var c float64
	for _, r := range rows {
		c += r.Cost * float64(r.CartNum)
	}
	return c
}

func trimSKU(unique string) string {
	if len(unique) > 12 {
		return unique[:12]
	}
	return unique
}

func skuUnique(p OrderProduct) string {
	var snap struct {
		Unique string `json:"unique"`
	}
	_ = json.Unmarshal([]byte(p.ProductInfo), &snap)
	if snap.Unique != "" {
		return snap.Unique
	}
	return p.ProductSKU
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
