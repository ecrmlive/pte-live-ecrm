-- 菜单管理：平台 / 商户 / 区域三类菜单及创建时间。
-- 既有菜单默认归属平台；新建菜单由后台接口校验父子归属和菜单标识唯一性。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @has_menu_scope := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_a_menu' AND COLUMN_NAME = 'menu_scope'
);
SET @sql := IF(@has_menu_scope = 0,
  'ALTER TABLE `qixi_crm_a_menu` ADD COLUMN `menu_scope` enum(''platform'',''merchant'',''region'') NOT NULL DEFAULT ''platform'' AFTER `status`',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;

SET @has_created_at := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_a_menu' AND COLUMN_NAME = 'created_at'
);
SET @sql := IF(@has_created_at = 0,
  'ALTER TABLE `qixi_crm_a_menu` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `menu_scope`',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;

SET @has_scope_parent_index := (
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_a_menu' AND INDEX_NAME = 'idx_menu_scope_parent'
);
SET @sql := IF(@has_scope_parent_index = 0,
  'ALTER TABLE `qixi_crm_a_menu` ADD KEY `idx_menu_scope_parent` (`menu_scope`, `parent_id`)',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;
