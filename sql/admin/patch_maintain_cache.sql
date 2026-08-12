-- 平台缓存清理审计。缓存键只允许 ecrm:platform:* 命名空间，避免影响 IM 等共享服务。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_cache_clear_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `scope` varchar(32) NOT NULL,
  `deleted_keys` bigint unsigned NOT NULL DEFAULT 0,
  `updated_assets` bigint unsigned NOT NULL DEFAULT 0,
  `operator_admin_id` bigint unsigned NOT NULL,
  `detail` varchar(512) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_scope_time` (`scope`,`created_at`), KEY `idx_operator_time` (`operator_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
