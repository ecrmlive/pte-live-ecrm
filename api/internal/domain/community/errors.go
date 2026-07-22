package community

import "errors"

var (
	ErrNotFound  = errors.New("内容不存在")
	ErrBadParam  = errors.New("参数错误")
	ErrForbidden = errors.New("无权操作")
)
