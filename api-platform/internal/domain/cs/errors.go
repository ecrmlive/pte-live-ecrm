package cs

import "errors"

var (
	ErrBadParam  = errors.New("参数错误")
	ErrNotFound  = errors.New("记录不存在")
	ErrForbidden = errors.New("无权操作")
)
