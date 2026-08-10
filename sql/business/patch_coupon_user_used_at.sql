-- 用户券使用时间（领取/使用记录弹窗）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_used_at := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_coupon_user'
    AND COLUMN_NAME = 'used_at'
);
SET @sql_used_at := IF(
  @col_used_at = 0,
  'ALTER TABLE `qixi_crm_b_coupon_user` ADD COLUMN `used_at` datetime DEFAULT NULL COMMENT ''核销/使用时间'' AFTER `obtained_at`',
  'SELECT 1'
);
PREPARE stmt_used_at FROM @sql_used_at;
EXECUTE stmt_used_at;
DEALLOCATE PREPARE stmt_used_at;

UPDATE `qixi_crm_b_coupon_user`
SET `used_at` = COALESCE(`used_at`, `obtained_at`)
WHERE `status` = 'used' AND `used_at` IS NULL;
