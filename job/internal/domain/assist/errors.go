package assist

import "errors"

var (
	ErrNotFound   = errors.New("助力活动不存在")
	ErrBadParam   = errors.New("参数错误")
	ErrForbidden  = errors.New("无权操作")
	ErrInactive   = errors.New("助力活动未开始或已结束")
	ErrSoldOut    = errors.New("助力库存不足")
	ErrSetNotOpen = errors.New("助力未完成或不可下单")
	ErrAlreadyHelped = errors.New("已助力过")
	ErrSelfHelp   = errors.New("不能为自己助力")
	ErrSetClosed  = errors.New("该助力已结束")
)
