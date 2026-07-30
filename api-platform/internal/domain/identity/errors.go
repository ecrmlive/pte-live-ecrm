package identity

import "errors"

var (
	ErrInvalidCredentials = errors.New("账号或密码错误")
	ErrAccountDisabled    = errors.New("账号已禁用")
	ErrMerchantDisabled   = errors.New("商户已关闭或禁用")
	ErrNotFound           = errors.New("记录不存在")
	ErrBadPassword        = errors.New("原密码不正确")
	ErrWeakPassword       = errors.New("新密码至少 6 位")
	ErrAccountExists      = errors.New("账号已存在")
	ErrNoVerifyPerm       = errors.New("无核销权限")
	ErrNoDeliverPerm      = errors.New("无发货权限")
	ErrNoCustomerPerm     = errors.New("无客服权限")
	ErrNoPerm             = errors.New("无操作权限")
	ErrBadParam           = errors.New("参数错误")
)
