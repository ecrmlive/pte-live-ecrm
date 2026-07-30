package invoice

import "errors"

var (
	ErrBadParam  = errors.New("参数错误")
	ErrNotFound  = errors.New("发票申请不存在")
	ErrForbidden = errors.New("无权操作")
	ErrExists    = errors.New("该订单已申请发票")
	ErrOrder     = errors.New("订单不可开票")
)
