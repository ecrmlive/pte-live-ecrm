-- 权限管理 / 角色权限：三类身份、创建时间和更新时间。
-- make local-sync-sql 会在本机 qixi_crm_admin 自动执行本补丁。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @role_type_exists := (
  SELECT COUNT(*) FROM information_schema.columns
  WHERE table_schema = DATABASE() AND table_name = 'qixi_crm_a_role' AND column_name = 'role_type'
);
SET @role_type_sql := IF(@role_type_exists = 0,
  'ALTER TABLE `qixi_crm_a_role` ADD COLUMN `role_type` enum(''platform'',''merchant'',''region'') NOT NULL DEFAULT ''platform'' AFTER `status`',
  'SELECT 1');
PREPARE role_type_stmt FROM @role_type_sql;
EXECUTE role_type_stmt;
DEALLOCATE PREPARE role_type_stmt;

SET @role_created_at_exists := (
  SELECT COUNT(*) FROM information_schema.columns
  WHERE table_schema = DATABASE() AND table_name = 'qixi_crm_a_role' AND column_name = 'created_at'
);
SET @role_created_at_sql := IF(@role_created_at_exists = 0,
  'ALTER TABLE `qixi_crm_a_role` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `role_type`',
  'SELECT 1');
PREPARE role_created_at_stmt FROM @role_created_at_sql;
EXECUTE role_created_at_stmt;
DEALLOCATE PREPARE role_created_at_stmt;

SET @role_updated_at_exists := (
  SELECT COUNT(*) FROM information_schema.columns
  WHERE table_schema = DATABASE() AND table_name = 'qixi_crm_a_role' AND column_name = 'updated_at'
);
SET @role_updated_at_sql := IF(@role_updated_at_exists = 0,
  'ALTER TABLE `qixi_crm_a_role` ADD COLUMN `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP AFTER `created_at`',
  'SELECT 1');
PREPARE role_updated_at_stmt FROM @role_updated_at_sql;
EXECUTE role_updated_at_stmt;
DEALLOCATE PREPARE role_updated_at_stmt;

UPDATE `qixi_crm_a_role`
SET `role_type` = CASE
  WHEN `code` = 'merchant' OR `code` LIKE 'merchant_%' THEN 'merchant'
  WHEN `code` = 'region' OR `code` LIKE 'region_%' THEN 'region'
  ELSE `role_type`
END
WHERE (`code` = 'merchant' OR `code` LIKE 'merchant_%' OR `code` = 'region' OR `code` LIKE 'region_%')
  AND `role_type` = 'platform';
