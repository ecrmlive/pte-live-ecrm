package logistics

import "errors"

var (
	ErrBadParam        = errors.New("参数错误")
	ErrExpressExists   = errors.New("物流公司名称或编码已存在")
	ErrNotFound        = errors.New("记录不存在")
	ErrDefaultTemplate = errors.New("默认运费模板不能删除，请先设置其他模板为默认")
)
