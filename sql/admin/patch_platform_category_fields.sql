-- 平台商品分类：图标 / 是否推荐 / 创建时间（对齐 CRMEB store_category.pic / is_hot / create_time）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='pic')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `pic` varchar(1024) NOT NULL DEFAULT '''' AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='is_hot')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `is_hot` tinyint NOT NULL DEFAULT 0 AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `is_hot`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
