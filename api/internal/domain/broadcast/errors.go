package broadcast

import "errors"

var (
	ErrNotFound  = errors.New("直播间不存在")
	ErrBadParam  = errors.New("参数错误")
	ErrForbidden = errors.New("无权操作该直播间")
)
