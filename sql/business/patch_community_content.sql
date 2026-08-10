-- 社区内容/评论：补齐星级、视频链接、评论点赞/回复数与拒绝原因（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_social_post' AND COLUMN_NAME='start')=0,
    'ALTER TABLE `qixi_crm_b_social_post` ADD COLUMN `start` tinyint NOT NULL DEFAULT 1 COMMENT ''推荐星级'' AFTER `is_hot`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_social_post' AND COLUMN_NAME='video_link')=0,
    'ALTER TABLE `qixi_crm_b_social_post` ADD COLUMN `video_link` varchar(512) NOT NULL DEFAULT '''' COMMENT ''短视频链接'' AFTER `is_type`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_social_reply' AND COLUMN_NAME='count_start')=0,
    'ALTER TABLE `qixi_crm_b_social_reply` ADD COLUMN `count_start` int NOT NULL DEFAULT 0 COMMENT ''评论点赞数'' AFTER `uid`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_social_reply' AND COLUMN_NAME='count_reply')=0,
    'ALTER TABLE `qixi_crm_b_social_reply` ADD COLUMN `count_reply` int NOT NULL DEFAULT 0 COMMENT ''评论条数'' AFTER `count_start`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_social_reply' AND COLUMN_NAME='refusal')=0,
    'ALTER TABLE `qixi_crm_b_social_reply` ADD COLUMN `refusal` varchar(500) NOT NULL DEFAULT '''' COMMENT ''拒绝原因'' AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 历史图文帖 is_type=0 视为图文
UPDATE `qixi_crm_b_social_post` SET `is_type`=1 WHERE `is_type`=0;
UPDATE `qixi_crm_b_social_post` SET `start`=1 WHERE `start` IS NULL OR `start`<1;
