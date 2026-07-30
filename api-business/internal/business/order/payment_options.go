package order

import (
	"context"
	"errors"

	"github.com/qixi-live/qixi-live-mergers/api-business/internal/paymentconfig"
	"gorm.io/gorm"
)

var ErrStoreChannelDisabled = errors.New("订单包含的店铺未启用该支付方式")

type PaymentChannelView struct {
	Channel string `json:"channel"`
	Enabled bool   `json:"enabled"`
}

// AvailablePaymentChannels selects the payment subject from the order itself:
// platform goods use only the platform bundle; merchant goods use only that
// store's bundle. There is deliberately no fallback between the two.
func AvailablePaymentChannels(ctx context.Context, db *gorm.DB, configs *paymentconfig.Store, userID, groupOrderID uint64) ([]PaymentChannelView, error) {
	if err := assertGroupOwner(ctx, db, userID, groupOrderID); err != nil {
		return nil, err
	}
	var order orderRow
	if err := db.WithContext(ctx).Where("group_order_id = ? AND user_id = ?", groupOrderID, userID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderOwnership
		}
		return nil, err
	}
	platformOwned := order.MerchantID == 0 && order.StoreID == 0
	var values paymentconfig.Values
	var err error
	if platformOwned {
		values, err = configs.Load(ctx)
	} else {
		values, err = configs.LoadStore(ctx, uint(order.StoreID))
	}
	if err != nil && !errors.Is(err, paymentconfig.ErrNotConfigured) {
		return nil, err
	}
	channels := []PaymentChannelView{{Channel: "wechat"}, {Channel: "alipay"}}
	for i := range channels {
		if err != nil {
			continue
		}
		if platformOwned {
			channels[i].Enabled = paymentconfig.ChannelReady(values, channels[i].Channel)
		} else {
			channels[i].Enabled = paymentconfig.StoreChannelReady(values, channels[i].Channel)
		}
	}
	return channels, nil
}

func AssertPaymentChannelAvailable(ctx context.Context, db *gorm.DB, configs *paymentconfig.Store, userID, groupOrderID uint64, channel string) error {
	channels, err := AvailablePaymentChannels(ctx, db, configs, userID, groupOrderID)
	if err != nil {
		return err
	}
	for _, item := range channels {
		if item.Channel == channel {
			if !item.Enabled {
				return ErrStoreChannelDisabled
			}
			return nil
		}
	}
	return ErrPayChannel
}

func assertGroupOwner(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) error {
	var count int64
	if err := db.WithContext(ctx).Model(&groupRow{}).Where("id = ? AND user_id = ?", groupOrderID, userID).Count(&count).Error; err != nil {
		return err
	}
	if count != 1 {
		return ErrOrderOwnership
	}
	return nil
}
