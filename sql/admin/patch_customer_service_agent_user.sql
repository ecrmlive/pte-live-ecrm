-- 客服账号关联 C 端用户（仅保存引用与客服头像覆盖值）。
-- 执行前请确认已备份；本补丁仅修改 qixi_crm_admin.qixi_crm_a_admin_user。
USE `qixi_crm_admin`;

SET @has_linked_user_id := (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'qixi_crm_a_admin_user'
    AND column_name = 'linked_user_id'
);
SET @sql := IF(
  @has_linked_user_id = 0,
  'ALTER TABLE `qixi_crm_a_admin_user` ADD COLUMN `linked_user_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `display_name`',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;

SET @has_avatar_url := (
  SELECT COUNT(*)
  FROM information_schema.columns
  WHERE table_schema = DATABASE()
    AND table_name = 'qixi_crm_a_admin_user'
    AND column_name = 'avatar_url'
);
SET @sql := IF(
  @has_avatar_url = 0,
  'ALTER TABLE `qixi_crm_a_admin_user` ADD COLUMN `avatar_url` varchar(1024) NOT NULL DEFAULT '''' AFTER `linked_user_id`',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;
