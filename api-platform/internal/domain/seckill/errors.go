package seckill

import "errors"

var (
	ErrNotFound   = errors.New("秒杀活动不存在")
	ErrBadParam   = errors.New("参数错误")
	ErrNotInWindow = errors.New("不在秒杀场次内")
	ErrInactive   = errors.New("秒杀未开启")
)
