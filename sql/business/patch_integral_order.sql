-- 积分订单：店铺备注 / 后台软删（幂等）
-- 对齐 CRMEB store_order.remark / is_system_del；用户删除用 group_order.user_archived_at。
SET NAMES utf8mb4;

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
   WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_order' AND COLUMN_NAME = 'merchant_remark') = 0,
  'ALTER TABLE `qixi_crm_b_order` ADD COLUMN `merchant_remark` varchar(500) NOT NULL DEFAULT '''' COMMENT ''店铺备注'' AFTER `remark`',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS
   WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_order' AND COLUMN_NAME = 'is_system_del') = 0,
  'ALTER TABLE `qixi_crm_b_order` ADD COLUMN `is_system_del` tinyint NOT NULL DEFAULT 0 COMMENT ''后台软删：1已删'' AFTER `merchant_remark`, ADD KEY `idx_activity_system_del` (`activity_type`,`is_system_del`,`id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
