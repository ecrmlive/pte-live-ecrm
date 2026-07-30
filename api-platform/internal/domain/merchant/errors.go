package merchant

import "errors"

var (
	ErrNotFound       = errors.New("记录不存在")
	ErrAlreadyAudited = errors.New("申请已审核")
	ErrBadStatus      = errors.New("审核状态无效")
	ErrRejectNeedMsg  = errors.New("拒绝时必须填写原因")
	ErrBadParam       = errors.New("参数错误")
)
