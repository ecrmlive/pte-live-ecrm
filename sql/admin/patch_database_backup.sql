-- 本地数据库备份及表维护记录；备份文件写入 Docker 挂载的 release/backups，不入 Git。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_database_backup` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `file_name` varchar(255) NOT NULL,
  `database_scope` enum('admin','business') NOT NULL,
  `table_names_json` longtext NOT NULL,
  `table_count` int unsigned NOT NULL DEFAULT 0,
  `size_bytes` bigint unsigned NOT NULL DEFAULT 0,
  `status` enum('ready','deleted','failed') NOT NULL DEFAULT 'ready',
  `created_by` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_backup_file_name` (`file_name`),
  KEY `idx_backup_status_time` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='本地数据库备份记录';

CREATE TABLE IF NOT EXISTS `qixi_crm_a_database_maintenance_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `action` enum('backup','optimize','repair','delete') NOT NULL,
  `database_scope` enum('admin','business') NOT NULL,
  `table_names_json` longtext NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL,
  `detail` varchar(1024) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_maintenance_action_time` (`action`,`created_at`),
  KEY `idx_maintenance_operator_time` (`operator_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库维护审计记录';
