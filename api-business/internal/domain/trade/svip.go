package trade

import (
	"context"
	"time"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cart"
)

// DefaultSvipRatio svip_price_type=1 时默认会员折扣（相对售价）。
const DefaultSvipRatio = 0.9

// applySvipPrices 会员价覆盖行价（秒杀行跳过）；返回是否用了 SVIP、是否应禁用店铺券。
func (s *Service) applySvipPrices(ctx context.Context, uid uint, rows []cart.Cart) (usedSvip bool, skipStoreCoupon bool, err error) {
	if len(rows) == 0 {
		return false, false, nil
	}
	isSvip, end, err := s.store.GetUserSVIP(ctx, uid)
	if err != nil {
		return false, false, err
	}
	if !UserSvipActiveAt(isSvip, end, time.Now()) {
		return false, false, nil
	}
	merUsed := map[uint]struct{}{}
	for i := range rows {
		if rows[i].SeckillActiveID > 0 || rows[i].PresellActiveID > 0 {
			continue
		}
		if rows[i].MerSvipStatus == 0 || rows[i].SvipPriceType == 0 {
			continue
		}
		vip := resolveSvipLinePrice(rows[i])
		if vip <= 0 || vip >= rows[i].Price {
			continue
		}
		diff := round2((rows[i].Price - vip) * float64(rows[i].CartNum))
		rows[i].Price = vip
		rows[i].UsedSvip = true
		rows[i].SvipDiscount = diff
		usedSvip = true
		merUsed[rows[i].MerID] = struct{}{}
	}
	if !usedSvip {
		return false, false, nil
	}
	merIDs := make([]uint, 0, len(merUsed))
	for m := range merUsed {
		merIDs = append(merIDs, m)
	}
	mergeMap, err := s.store.MerchantsSVIPCouponMerge(ctx, merIDs)
	if err != nil {
		return false, false, err
	}
	// §8.2：任一用了 SVIP 的商户 merge≠1 → 清空店铺券
	for _, m := range merIDs {
		if mergeMap[m] != 1 {
			skipStoreCoupon = true
			break
		}
	}
	return usedSvip, skipStoreCoupon, nil
}

func resolveSvipLinePrice(r cart.Cart) float64 {
	switch r.SvipPriceType {
	case 1:
		if r.Price <= 0 {
			return 0
		}
		return round2(r.Price * DefaultSvipRatio)
	case 2:
		return r.SvipPrice
	default:
		return 0
	}
}

func sumSvipDiscount(rows []cart.Cart) float64 {
	var sum float64
	for _, r := range rows {
		sum += r.SvipDiscount
	}
	return round2(sum)
}

func merSvipDiscount(rows []cart.Cart, merID uint) float64 {
	var sum float64
	for _, r := range rows {
		if r.MerID == merID {
			sum += r.SvipDiscount
		}
	}
	return round2(sum)
}

func cartHasActivityGoods(rows []cart.Cart) bool {
	for _, r := range rows {
		if r.SeckillActiveID > 0 || r.PresellActiveID > 0 || r.ProductType != 0 {
			return true
		}
	}
	return false
}

// UserSvipActiveAt is_svip: 1体验 2有效期 3永久。
func UserSvipActiveAt(isSvip int8, end *time.Time, now time.Time) bool {
	switch isSvip {
	case 1, 3:
		return true
	case 2:
		return end != nil && end.After(now)
	default:
		return false
	}
}
