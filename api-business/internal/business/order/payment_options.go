package order

import (
	"context"
	"errors"

	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
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
	return availablePaymentChannels(ctx, db, configs, nil, userID, groupOrderID)
}

func availablePaymentChannels(ctx context.Context, db *gorm.DB, configs *paymentconfig.Store, platformConfigs *cloudconfig.Service, userID, groupOrderID uint64) ([]PaymentChannelView, error) {
	if err := assertGroupOwner(ctx, db, userID, groupOrderID); err != nil {
		return nil, err
	}
	order, err := paymentSubject(ctx, db, userID, groupOrderID)
	if err != nil {
		return nil, err
	}
	platformOwned := order.MerchantID == 0 && order.StoreID == 0
	values, err := loadPaymentValues(ctx, configs, platformConfigs, order, platformOwned)
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

func paymentSubject(ctx context.Context, db *gorm.DB, userID, groupOrderID uint64) (orderRow, error) {
	var orders []orderRow
	if err := db.WithContext(ctx).Where("group_order_id = ? AND user_id = ?", groupOrderID, userID).Find(&orders).Error; err != nil {
		return orderRow{}, err
	}
	if len(orders) == 0 {
		return orderRow{}, ErrOrderOwnership
	}
	first := orders[0]
	for _, order := range orders[1:] {
		if order.MerchantID != first.MerchantID || order.StoreID != first.StoreID {
			return orderRow{}, ErrMixedPaySubject
		}
	}
	return first, nil
}

func loadPaymentValues(ctx context.Context, configs *paymentconfig.Store, platformConfigs *cloudconfig.Service, order orderRow, platformOwned bool) (paymentconfig.Values, error) {
	if configs == nil {
		return nil, paymentconfig.ErrNotConfigured
	}
	if !platformOwned {
		return configs.LoadStore(ctx, uint(order.StoreID))
	}
	if platformConfigs != nil {
		values, err := platformConfigs.Values(ctx, "payment")
		if err != nil {
			return nil, err
		}
		return paymentconfig.Values(values), nil
	}
	return configs.Load(ctx)
}

func AssertPaymentChannelAvailable(ctx context.Context, db *gorm.DB, configs *paymentconfig.Store, userID, groupOrderID uint64, channel string) error {
	return assertPaymentChannelAvailable(ctx, db, configs, nil, userID, groupOrderID, channel)
}

func assertPaymentChannelAvailable(ctx context.Context, db *gorm.DB, configs *paymentconfig.Store, platformConfigs *cloudconfig.Service, userID, groupOrderID uint64, channel string) error {
	channels, err := availablePaymentChannels(ctx, db, configs, platformConfigs, userID, groupOrderID)
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
