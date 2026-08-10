-- 用户反馈分类：对齐 CRMEB eb_feedback_category 父子树（pid）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_pid := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_user_feedback_category'
    AND COLUMN_NAME = 'pid'
);
SET @sql_pid := IF(
  @col_pid = 0,
  'ALTER TABLE `qixi_crm_b_user_feedback_category` ADD COLUMN `pid` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''上级分类 ID，0=顶级'' AFTER `id`, ADD KEY `idx_feedback_category_pid` (`pid`)',
  'SELECT 1'
);
PREPARE stmt_pid FROM @sql_pid;
EXECUTE stmt_pid;
DEALLOCATE PREPARE stmt_pid;
