package chat

import "errors"

var (
	ErrBadParam  = errors.New("参数错误")
	ErrIMRemoteRequired = errors.New("客服密聊需 im.mode=remote 且配置 integration_token")
	ErrTextViaCS        = errors.New("文本消息请经 PTE IM SDK 发送，本仓 /cs 仅保留订单卡片等元数据")
	ErrNotFound  = errors.New("会话不存在")
	ErrForbidden = errors.New("无权操作")
	ErrClosed    = errors.New("会话已关闭")
)
