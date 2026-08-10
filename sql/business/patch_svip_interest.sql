-- 付费会员权益：对齐 CRMEB 展示名称 / 未开通与已开通双图标 / 跳转链接
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_interest' AND COLUMN_NAME='display_name')=0,
    'ALTER TABLE `qixi_crm_b_svip_interest` ADD COLUMN `display_name` varchar(64) NOT NULL DEFAULT \'\' AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_interest' AND COLUMN_NAME='on_icon_url')=0,
    'ALTER TABLE `qixi_crm_b_svip_interest` ADD COLUMN `on_icon_url` varchar(1024) NOT NULL DEFAULT \'\' AFTER `icon_url`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_interest' AND COLUMN_NAME='link')=0,
    'ALTER TABLE `qixi_crm_b_svip_interest` ADD COLUMN `link` varchar(500) NOT NULL DEFAULT \'\' AFTER `on_icon_url`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 兼容旧行：展示名默认回填权益名称，已开通图标默认回填未开通图标
UPDATE `qixi_crm_b_svip_interest`
SET `display_name` = IF(`display_name` = '' OR `display_name` IS NULL, `name`, `display_name`),
    `on_icon_url` = IF(`on_icon_url` = '' OR `on_icon_url` IS NULL, `icon_url`, `on_icon_url`)
WHERE `deleted_at` IS NULL;
