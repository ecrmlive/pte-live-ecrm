package finance

import "errors"

var (
	ErrNotFound         = errors.New("提现单不存在")
	ErrForbidden        = errors.New("无权操作")
	ErrBadStatus        = errors.New("提现单状态不允许此操作")
	ErrBadParam         = errors.New("参数错误")
	ErrBalanceNotEnough = errors.New("商户余额不足")
	ErrAlreadyDone      = errors.New("提现已处理")
	ErrRejectNeedMsg    = errors.New("拒绝时必须填写原因")
	ErrMerchantNotFound = errors.New("商户不存在")
)
