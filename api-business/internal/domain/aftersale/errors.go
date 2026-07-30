package aftersale

import "errors"

var (
	ErrNotFound         = errors.New("退款单不存在")
	ErrForbidden        = errors.New("无权操作")
	ErrBadStatus        = errors.New("退款单状态不允许此操作")
	ErrBadParam         = errors.New("参数错误")
	ErrOrderNotFound    = errors.New("订单不存在")
	ErrOrderNotPaid     = errors.New("仅已支付订单可申请退款")
	ErrOrderRefunded    = errors.New("订单已退款")
	ErrRefundInProgress = errors.New("该订单已有进行中的退款")
	ErrProductInvalid   = errors.New("退款商品无效")
	ErrAlreadyDone      = errors.New("退款已处理")
	ErrRejectNeedMsg    = errors.New("拒绝时必须填写原因")
)
