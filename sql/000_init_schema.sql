-- Phase 0：库初始化占位。完整表结构见 docs/schema/qixi_schema_reference.sql
-- 业务迁移按阶段追加 001_*.sql …

CREATE DATABASE IF NOT EXISTS `qixi_mergers`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `qixi_mergers`;

-- 健康探测用占位表（非业务）
CREATE TABLE IF NOT EXISTS `qixi_schema_meta` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `version` VARCHAR(64) NOT NULL DEFAULT '',
  `note` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase0', '骨架已就绪，业务表按阶段迁入'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase0');
