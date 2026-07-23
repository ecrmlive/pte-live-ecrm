package chat

import "errors"

var (
	ErrBadParam  = errors.New("参数错误")
	ErrNotFound  = errors.New("会话不存在")
	ErrForbidden = errors.New("无权操作")
	ErrClosed    = errors.New("会话已关闭")
)
