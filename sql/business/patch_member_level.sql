-- 会员等级：对齐 CRMEB 等级管理（图标 / 成长值 / 创建时间）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='icon_url')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `icon_url` varchar(1024) NOT NULL DEFAULT \'\' AFTER `rank`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `version`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
