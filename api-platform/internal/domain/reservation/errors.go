package reservation

import "errors"

var (
	ErrNotFound   = errors.New("预约商品不存在")
	ErrBadParam   = errors.New("参数错误")
	ErrForbidden  = errors.New("无权操作")
	ErrNoSlot     = errors.New("时段不可约")
	ErrFull       = errors.New("该时段已约满")
	ErrBadDate    = errors.New("预约日期无效")
)
