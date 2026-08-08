-- 平台品牌：创建时间（对齐列表「创建时间」列）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_brand' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_brand` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
