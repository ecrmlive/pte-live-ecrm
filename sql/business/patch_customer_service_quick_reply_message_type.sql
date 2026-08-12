-- 客服自动回复支持文字与图片消息；已有记录默认视为文字消息。
SET @message_type_column_exists := (
  SELECT COUNT(*)
  FROM `information_schema`.`columns`
  WHERE `table_schema` = DATABASE()
    AND `table_name` = 'qixi_crm_b_customer_service_quick_reply'
    AND `column_name` = 'message_type'
);
SET @add_message_type_column := IF(
  @message_type_column_exists = 0,
  'ALTER TABLE `qixi_crm_b_customer_service_quick_reply` ADD COLUMN `message_type` enum(\'text\',\'image\') NOT NULL DEFAULT \'text\' AFTER `content`',
  'SELECT 1'
);
PREPARE `add_message_type_column_stmt` FROM @add_message_type_column;
EXECUTE `add_message_type_column_stmt`;
DEALLOCATE PREPARE `add_message_type_column_stmt`;
