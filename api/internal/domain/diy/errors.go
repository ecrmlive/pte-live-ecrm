package diy

import "errors"

var (
	ErrNotFound = errors.New("装修页不存在")
	ErrBadParam = errors.New("参数错误")
)
