package combination

import "errors"

var (
	ErrNotFound     = errors.New("拼团活动不存在")
	ErrBadParam     = errors.New("参数错误")
	ErrInactive     = errors.New("拼团未开启或不在活动期")
	ErrBuyingFull   = errors.New("该团已满员")
	ErrBuyingClosed = errors.New("该团已结束")
	ErrAlreadyJoined = errors.New("已参加该团")
	ErrForbidden    = errors.New("无权操作")
)
