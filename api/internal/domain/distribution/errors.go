package distribution

import "errors"

var (
	ErrBadParam      = errors.New("参数错误")
	ErrSelfBind      = errors.New("不能绑定自己")
	ErrAlreadyBound  = errors.New("已绑定推广关系")
	ErrUserNotFound  = errors.New("用户不存在")
	ErrSpreadInvalid = errors.New("推广员无效")
)
