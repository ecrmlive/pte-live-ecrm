CREATE DATABASE IF NOT EXISTS `qixi_crm_admin` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_admin_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(64) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(64) NOT NULL,
  `phone` varchar(32) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1,
  `auth_version` bigint unsigned NOT NULL DEFAULT 1,
  `data_scope_version` bigint unsigned NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `code` varchar(32) NOT NULL,
  `name` varchar(64) NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_admin_user_role` (
  `admin_user_id` bigint unsigned NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`admin_user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `code` varchar(128) NOT NULL,
  `title` varchar(64) NOT NULL,
  `icon` varchar(96) NOT NULL DEFAULT '',
  `route_path` varchar(255) NOT NULL DEFAULT '',
  `kind` enum('directory','page','button') NOT NULL,
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 兼容已初始化的 local/test 库：MySQL 8.4 不支持 ADD COLUMN IF NOT EXISTS，
-- 因此通过 information_schema 判断后再执行幂等迁移。
SET @qixi_menu_icon_exists := (
  SELECT COUNT(*)
  FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = 'qixi_crm_admin'
    AND TABLE_NAME = 'qixi_crm_a_menu'
    AND COLUMN_NAME = 'icon'
);
SET @qixi_menu_icon_ddl := IF(
  @qixi_menu_icon_exists = 0,
  'ALTER TABLE `qixi_crm_a_menu` ADD COLUMN `icon` varchar(96) NOT NULL DEFAULT '''' AFTER `title`',
  'SELECT 1'
);
PREPARE qixi_menu_icon_stmt FROM @qixi_menu_icon_ddl;
EXECUTE qixi_menu_icon_stmt;
DEALLOCATE PREPARE qixi_menu_icon_stmt;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_role_menu` (
  `role_id` bigint unsigned NOT NULL,
  `menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_data_scope` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `admin_user_id` bigint unsigned NOT NULL,
  `scope_type` enum('all','merchant','region','store','service_queue','operations') NOT NULL,
  `scope_value` json NOT NULL,
  `version` bigint unsigned NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user_type` (`admin_user_id`,`scope_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_config` (
  `config_key` varchar(128) NOT NULL,
  `config_value` json NOT NULL,
  `updated_by` bigint unsigned DEFAULT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_cloud_config` (
  `provider` varchar(64) NOT NULL,
  `config_key` varchar(128) NOT NULL,
  `ciphertext` text NOT NULL,
  `key_version` varchar(32) NOT NULL,
  `updated_by` bigint unsigned DEFAULT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`provider`,`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_operation_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `admin_user_id` bigint unsigned NOT NULL,
  `role_code` varchar(32) NOT NULL,
  `action` varchar(128) NOT NULL,
  `resource_type` varchar(64) NOT NULL,
  `resource_id` varchar(64) NOT NULL DEFAULT '',
  `request_id` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_request` (`request_id`), KEY `idx_user_time` (`admin_user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_region` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `code` varchar(64) NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `applicant_user_id` bigint unsigned DEFAULT NULL, `merchant_name` varchar(128) NOT NULL,
  `contact_name` varchar(64) NOT NULL, `contact_mobile` varchar(32) NOT NULL, `region_id` bigint unsigned DEFAULT NULL,
  `status` enum('draft','pending','approved','rejected') NOT NULL DEFAULT 'draft', `reviewed_by` bigint unsigned DEFAULT NULL,
  `review_note` varchar(500) NOT NULL DEFAULT '', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `reviewed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_region_status` (`region_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台监管只读取该投影，不直连 qixi_crm_merchant。由 api-merchant 的受控
-- 命令和 NATS 事件同步，保证服务与数据库边界独立。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_view` (
  `merchant_id` bigint unsigned NOT NULL, `merchant_name` varchar(128) NOT NULL,
  `contact_name` varchar(64) NOT NULL DEFAULT '', `contact_mobile` varchar(32) NOT NULL DEFAULT '',
  `region_id` bigint unsigned DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`merchant_id`), KEY `idx_region_status` (`region_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `owner_scope` enum('platform','operations','merchant') NOT NULL,
  `owner_id` bigint unsigned NOT NULL DEFAULT 0, `media_type` enum('image','video','file') NOT NULL, `url` varchar(1024) NOT NULL,
  `name` varchar(255) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_owner` (`owner_scope`,`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_content` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `content_type` enum('article','notice','agreement','banner') NOT NULL,
  `title` varchar(255) NOT NULL, `body` longtext NOT NULL, `status` enum('draft','published','hidden') NOT NULL DEFAULT 'draft',
  `published_at` datetime DEFAULT NULL, `created_by` bigint unsigned NOT NULL, PRIMARY KEY (`id`), KEY `idx_type_status` (`content_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_diy_page` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `page_type` enum('home','store_street','member','custom') NOT NULL,
  `name` varchar(128) NOT NULL, `document` json NOT NULL, `status` enum('draft','published') NOT NULL DEFAULT 'draft',
  `updated_by` bigint unsigned NOT NULL, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_type_status` (`page_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_diy_link_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `pid` bigint unsigned NOT NULL DEFAULT 0,
  `type` varchar(32) NOT NULL DEFAULT 'link', `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1, `level` tinyint NOT NULL DEFAULT 1, `is_mer` tinyint NOT NULL DEFAULT 0,
  `add_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_scope_parent` (`is_mer`,`pid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_diy_link` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `cate_id` bigint unsigned NOT NULL, `type` tinyint NOT NULL DEFAULT 1,
  `name` varchar(128) NOT NULL, `url` varchar(1024) NOT NULL, `param` varchar(1000) NOT NULL DEFAULT '',
  `example` varchar(1000) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0,
  `is_mer` tinyint NOT NULL DEFAULT 0, `add_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_scope_category` (`is_mer`,`cate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `event_type` varchar(128) NOT NULL, `aggregate_type` varchar(64) NOT NULL,
  `aggregate_id` varchar(64) NOT NULL, `payload` json NOT NULL, `status` enum('pending','published','failed') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `published_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_status_time` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_export_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `admin_user_id` bigint unsigned NOT NULL, `resource_type` varchar(64) NOT NULL,
  `filters` json NOT NULL, `file_url` varchar(1024) DEFAULT NULL, `status` enum('pending','completed','failed') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_user_time` (`admin_user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 统一后台鉴权、数据范围、通知与审计域。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_login_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `admin_user_id` bigint unsigned DEFAULT NULL, `username` varchar(64) NOT NULL,
  `role_code` varchar(32) NOT NULL DEFAULT '', `success` tinyint NOT NULL, `ip` varchar(64) NOT NULL DEFAULT '',
  `user_agent` varchar(512) NOT NULL DEFAULT '', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user_time` (`admin_user_id`,`created_at`), KEY `idx_username_time` (`username`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_token_revocation` (
  `jti` varchar(191) NOT NULL, `admin_user_id` bigint unsigned NOT NULL, `expires_at` datetime NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`jti`), KEY `idx_expire` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_data_scope_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `role_id` bigint unsigned NOT NULL, `scope_type` enum('all','merchant','region','store','service_queue','operations') NOT NULL,
  `rule` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), KEY `idx_role_type` (`role_id`,`scope_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_notification_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `code` varchar(128) NOT NULL, `channel` enum('in_app','sms','wechat','email') NOT NULL,
  `title_template` varchar(255) NOT NULL, `body_template` text NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `updated_by` bigint unsigned DEFAULT NULL, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code_channel` (`code`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_admin_notification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `admin_user_id` bigint unsigned NOT NULL, `category` varchar(64) NOT NULL,
  `title` varchar(255) NOT NULL, `body` varchar(2000) NOT NULL, `payload` json NOT NULL, `read_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_user_read_time` (`admin_user_id`,`read_at`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_config_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `config_key` varchar(128) NOT NULL, `old_value` json DEFAULT NULL,
  `new_value` json NOT NULL, `updated_by` bigint unsigned NOT NULL, `request_id` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_key_time` (`config_key`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 平台、区域与运营规则；业务事实保持在 business/merchant 库。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_district` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `region_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `boundary` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), KEY `idx_region` (`region_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0, `name` varchar(128) NOT NULL,
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), KEY `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_brand` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `logo_url` varchar(1024) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_review` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `status` enum('pending','approved','rejected','off_sale') NOT NULL DEFAULT 'pending', `reason` varchar(1000) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `reviewed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_marketing_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `rule_type` enum('coupon','seckill','combination','presell','assist','svip','distribution') NOT NULL,
  `name` varchar(128) NOT NULL, `rule` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, `updated_by` bigint unsigned NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_type_status_time` (`rule_type`,`status`,`starts_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_approval_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `task_type` varchar(64) NOT NULL, `subject_type` varchar(64) NOT NULL,
  `subject_id` varchar(64) NOT NULL, `status` enum('pending','approved','rejected','cancelled') NOT NULL DEFAULT 'pending',
  `payload` json NOT NULL, `handled_by` bigint unsigned DEFAULT NULL, `handle_note` varchar(1000) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `handled_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_type_status` (`task_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_job_definition` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `code` varchar(128) NOT NULL, `name` varchar(128) NOT NULL,
  `cron_expr` varchar(128) NOT NULL, `payload` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `updated_by` bigint unsigned DEFAULT NULL, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
