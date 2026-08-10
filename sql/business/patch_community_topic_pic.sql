-- 社区话题：补齐图标字段（对齐 CRMEB eb_community_topic.pic，幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_social_topic'
    AND COLUMN_NAME = 'pic'
);
SET @sql := IF(
  @col_exists = 0,
  'ALTER TABLE `qixi_crm_b_social_topic` ADD COLUMN `pic` varchar(255) NOT NULL DEFAULT '''' COMMENT ''话题图标'' AFTER `topic_name`',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
