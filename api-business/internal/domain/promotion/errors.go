package promotion

import "errors"

var (
	ErrNotFound        = errors.New("优惠券不存在")
	ErrForbidden       = errors.New("无权操作")
	ErrBadParam        = errors.New("参数错误")
	ErrBadStatus       = errors.New("状态无效")
	ErrClosed          = errors.New("优惠券已关闭")
	ErrSoldOut         = errors.New("优惠券已领完")
	ErrAlreadyReceived = errors.New("已领取过该券")
	ErrNotAvailable    = errors.New("优惠券不可领取")
	ErrRemainEmpty     = errors.New("优惠券已领完")
	ErrCouponInvalid   = errors.New("优惠券不可用")
	ErrCouponMinNotMet = errors.New("未达到用券门槛")
	ErrCouponConflict  = errors.New("同类优惠券只能选一张")
	ErrPlatformOnly    = errors.New("结算仅支持平台券")
	ErrSpreadSelf      = errors.New("不能绑定自己为推广人")
	ErrSpreadInvalid   = errors.New("推广人不存在或非推广员")
	ErrSpreadBound      = errors.New("已绑定推广人")
	ErrBillDup         = errors.New("佣金已入账")
)
