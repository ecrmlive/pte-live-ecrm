package fulfillment

import "errors"

var (
	ErrBadParam = errors.New("参数错误")
	ErrNotFound = errors.New("配送员不存在")
)
