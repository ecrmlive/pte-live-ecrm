-- IM 远程桥：会话绑定 pte-live-im conversation_id；身份存数值 IM user_id
USE `qixi_mergers`;

SET @col := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_cs_thread' AND COLUMN_NAME = 'im_conversation_id'
);
SET @sql := IF(@col = 0,
  'ALTER TABLE `qixi_cs_thread` ADD COLUMN `im_conversation_id` bigint(20) unsigned NOT NULL DEFAULT 0 COMMENT ''pte-live-im conversation_id'' AFTER `service_id`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col2 := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_im_identity' AND COLUMN_NAME = 'im_user_num'
);
SET @sql2 := IF(@col2 = 0,
  'ALTER TABLE `qixi_im_identity` ADD COLUMN `im_user_num` bigint(20) NOT NULL DEFAULT 0 COMMENT ''数值 IM user_id'' AFTER `im_user_id`, ADD KEY `idx_im_user_num` (`im_user_num`)',
  'SELECT 1');
PREPARE stmt2 FROM @sql2; EXECUTE stmt2; DEALLOCATE PREPARE stmt2;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'im-remote-044', 'IM 远程桥：im_conversation_id / im_user_num'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'im-remote-044');
