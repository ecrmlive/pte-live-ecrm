package catalog

import "errors"

var (
	ErrNotFound      = errors.New("记录不存在")
	ErrBadStatus     = errors.New("审核状态无效")
	ErrRejectNeedMsg = errors.New("拒绝时必须填写原因")
	ErrNameRequired  = errors.New("名称不能为空")
	ErrForbidden     = errors.New("无权操作该商品")
	ErrInvalidPrice  = errors.New("价格或库存无效")
	ErrNotOnSale     = errors.New("商品未上架或未过审")
	ErrCateRequired  = errors.New("请选择分类")
)
