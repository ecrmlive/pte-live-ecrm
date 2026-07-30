package presell

import "errors"

var (
	ErrNotFound       = errors.New("预售活动不存在")
	ErrBadParam       = errors.New("参数错误")
	ErrInactive       = errors.New("预售活动未开始或已结束")
	ErrSoldOut        = errors.New("预售库存不足")
	ErrNotFullPay     = errors.New("暂仅支持全款预售")
	ErrFinalNotFound  = errors.New("尾款单不存在")
	ErrFinalNotOpen   = errors.New("未到尾款支付时间或已截止")
	ErrFinalPaid      = errors.New("尾款已支付")
	ErrFinalInvalid   = errors.New("尾款单无效")
)
