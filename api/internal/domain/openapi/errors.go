package openapi

import "errors"

var (
	ErrBadParam     = errors.New("参数错误")
	ErrUnauthorized = errors.New("凭证无效")
	ErrForbidden    = errors.New("无权访问")
	ErrDisabled     = errors.New("凭证已禁用")
)
