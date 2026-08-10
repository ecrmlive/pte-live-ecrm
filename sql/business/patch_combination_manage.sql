-- 拼团商品平台监管：拒绝/下架原因 + 待审演示行（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @db := DATABASE();
SET @has_refusal := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_crm_b_combination_group' AND COLUMN_NAME = 'refusal'
);
SET @sql := IF(
  @has_refusal = 0,
  'ALTER TABLE `qixi_crm_b_combination_group` ADD COLUMN `refusal` varchar(500) NOT NULL DEFAULT '''' COMMENT ''拒绝/下架原因'' AFTER `product_status`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

UPDATE `qixi_crm_b_combination_group`
SET `product_status` = 0, `is_show` = 0, `refusal` = ''
WHERE `product_group_id` = 6104 AND `product_status` = 1;
