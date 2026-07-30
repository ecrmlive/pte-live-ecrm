package cart

import "errors"

var (
	ErrNotFound       = errors.New("购物车不存在")
	ErrAddrNotFound   = errors.New("地址不存在")
	ErrAddrInvalid    = errors.New("地址信息不完整")
	ErrInvalidNum     = errors.New("数量无效")
	ErrProductOff     = errors.New("商品不可售")
	ErrStockNotEnough = errors.New("库存不足")
	ErrForbidden      = errors.New("无权操作")
)
