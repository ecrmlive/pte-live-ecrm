-- 客服自动回复支持文字与图片消息；已有记录默认视为文字消息。
ALTER TABLE `qixi_crm_b_customer_service_quick_reply`
  ADD COLUMN IF NOT EXISTS `message_type` enum('text','image') NOT NULL DEFAULT 'text' AFTER `content`;
