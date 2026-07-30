package content

import "errors"

var (
	ErrNotFound     = errors.New("公告不存在")
	ErrAgreeNotFound = errors.New("协议不存在")
	ErrBadParam     = errors.New("参数错误")
)
