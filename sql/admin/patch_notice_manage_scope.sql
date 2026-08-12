-- 公告管理：公告投放范围必须使用店铺、店铺类别或店铺分类的真实关联，禁止保存输入式 ID。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @has_scope_type := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_a_notice' AND COLUMN_NAME = 'scope_type'
);
SET @sql := IF(@has_scope_type = 0,
  'ALTER TABLE `qixi_crm_a_notice` ADD COLUMN `scope_type` enum(''all'',''store_name'',''store_type'',''store_category'') NOT NULL DEFAULT ''all'' AFTER `sort`, ADD KEY `idx_notice_scope_type` (`scope_type`)',
  'SELECT 1'
);
PREPARE statement FROM @sql;
EXECUTE statement;
DEALLOCATE PREPARE statement;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_notice_scope` (
  `notice_id` bigint unsigned NOT NULL,
  `scope_kind` enum('store_name','store_type','store_category') NOT NULL,
  `scope_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`notice_id`,`scope_kind`,`scope_id`),
  KEY `idx_scope_lookup` (`scope_kind`,`scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公告店铺范围关联';
