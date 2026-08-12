package notification

import "errors"

var (
	ErrBadParam        = errors.New("参数错误")
	ErrNotFound        = errors.New("通知配置不存在")
	ErrSyncUnavailable = errors.New("未配置外部平台凭据，未执行同步")
)
