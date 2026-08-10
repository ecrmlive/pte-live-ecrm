package seckill

import "errors"

var (
	ErrNotFound         = errors.New("秒杀商品不存在")
	ErrActivityNotFound = errors.New("秒杀活动不存在")
	ErrTimeNotFound     = errors.New("秒杀场次不存在")
	ErrBadParam         = errors.New("参数错误")
	ErrTimeOverlap      = errors.New("所选时间段与已有场次重叠")
	ErrNotInWindow      = errors.New("不在秒杀场次内")
	ErrInactive         = errors.New("秒杀未开启")
)
