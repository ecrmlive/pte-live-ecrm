SET NAMES utf8mb4;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_config_classification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL,
  `classify_key` varchar(64) NOT NULL,
  `description` varchar(500) NOT NULL DEFAULT '',
  `icon` varchar(96) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1,
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_classify_key` (`classify_key`),
  KEY `idx_visible_sort` (`is_del`,`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台配置分类';

CREATE TABLE IF NOT EXISTS `qixi_crm_a_config_classification_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `classification_id` bigint unsigned NOT NULL,
  `name` varchar(128) NOT NULL,
  `config_key` varchar(128) NOT NULL,
  `field_type` varchar(32) NOT NULL DEFAULT 'input' COMMENT '配置控件类型',
  `backend_type` tinyint NOT NULL DEFAULT 0 COMMENT '0=总后台 1=商户后台',
  `content` text NOT NULL,
  `description` varchar(500) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1,
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_classification_config_key` (`classification_id`,`config_key`),
  KEY `idx_classification_visible_sort` (`classification_id`,`is_del`,`status`,`sort`,`id`),
  KEY `idx_backend_visible_sort` (`backend_type`,`is_del`,`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台配置分类项';

-- MySQL 8.4 兼容：不依赖 ADD COLUMN IF NOT EXISTS，补丁可以安全重复执行。
SET @field_type_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_config_classification_item' AND COLUMN_NAME='field_type')=0,
  'ALTER TABLE `qixi_crm_a_config_classification_item` ADD COLUMN `field_type` varchar(32) NOT NULL DEFAULT ''input'' COMMENT ''配置控件类型'' AFTER `config_key`',
  'SELECT 1');
PREPARE field_type_stmt FROM @field_type_sql;
EXECUTE field_type_stmt;
DEALLOCATE PREPARE field_type_stmt;

SET @backend_type_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_config_classification_item' AND COLUMN_NAME='backend_type')=0,
  'ALTER TABLE `qixi_crm_a_config_classification_item` ADD COLUMN `backend_type` tinyint NOT NULL DEFAULT 0 COMMENT ''0=总后台 1=商户后台'' AFTER `field_type`',
  'SELECT 1');
PREPARE backend_type_stmt FROM @backend_type_sql;
EXECUTE backend_type_stmt;
DEALLOCATE PREPARE backend_type_stmt;

SET @idx_backend_visible_sort_sql := IF(
  (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_config_classification_item' AND INDEX_NAME='idx_backend_visible_sort')=0,
  'CREATE INDEX `idx_backend_visible_sort` ON `qixi_crm_a_config_classification_item` (`backend_type`,`is_del`,`status`,`sort`,`id`)',
  'SELECT 1');
PREPARE idx_backend_visible_sort_stmt FROM @idx_backend_visible_sort_sql;
EXECUTE idx_backend_visible_sort_stmt;
DEALLOCATE PREPARE idx_backend_visible_sort_stmt;
