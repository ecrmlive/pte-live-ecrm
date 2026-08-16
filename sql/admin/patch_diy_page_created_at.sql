-- 页面装修模板的添加时间：兼容已有本地库，并允许本补丁重复执行。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @diy_page_created_at_exists := (
  SELECT COUNT(*) FROM information_schema.columns
  WHERE table_schema = DATABASE() AND table_name = 'qixi_crm_a_diy_page' AND column_name = 'created_at'
);
SET @diy_page_created_at_sql := IF(@diy_page_created_at_exists = 0,
  'ALTER TABLE `qixi_crm_a_diy_page` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `updated_by`',
  'SELECT 1');
PREPARE diy_page_created_at_stmt FROM @diy_page_created_at_sql;
EXECUTE diy_page_created_at_stmt;
DEALLOCATE PREPARE diy_page_created_at_stmt;
