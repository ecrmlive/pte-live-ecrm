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
  -- 区域后台账号可选地绑定一个已审核通过的区域代理；NULL 表示非区域账号。
  `circle_agent_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_username` (`username`), UNIQUE KEY `uk_circle_agent_id` (`circle_agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 已初始化数据库的向前兼容迁移。唯一约束保证一个代理最多绑定一个统一后台账号；
-- 若历史数据存在重复绑定，迁移应失败并要求先人工厘清账号归属，不能静默选择其中一个。
SET @qixi_admin_circle_agent_column_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_admin_user' AND COLUMN_NAME = 'circle_agent_id'
);
SET @qixi_admin_circle_agent_column_ddl := IF(
  @qixi_admin_circle_agent_column_exists = 0,
  'ALTER TABLE `qixi_crm_a_admin_user` ADD COLUMN `circle_agent_id` bigint unsigned DEFAULT NULL AFTER `data_scope_version`',
  'SELECT 1'
);
PREPARE qixi_admin_circle_agent_column_stmt FROM @qixi_admin_circle_agent_column_ddl;
EXECUTE qixi_admin_circle_agent_column_stmt;
DEALLOCATE PREPARE qixi_admin_circle_agent_column_stmt;

SET @qixi_admin_circle_agent_index_exists := (
  SELECT COUNT(*) FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_admin_user' AND INDEX_NAME = 'uk_circle_agent_id'
);
SET @qixi_admin_circle_agent_index_ddl := IF(
  @qixi_admin_circle_agent_index_exists = 0,
  'ALTER TABLE `qixi_crm_a_admin_user` ADD UNIQUE KEY `uk_circle_agent_id` (`circle_agent_id`)',
  'SELECT 1'
);
PREPARE qixi_admin_circle_agent_index_stmt FROM @qixi_admin_circle_agent_index_ddl;
EXECUTE qixi_admin_circle_agent_index_stmt;
DEALLOCATE PREPARE qixi_admin_circle_agent_index_stmt;

-- 逻辑删除用于保留后台操作、客服转接和资金审批的历史归属；已删除账号
-- 不能重新登录，也不会再出现在可管理账号或客服转接坐席中。
SET @qixi_admin_deleted_at_column_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_admin_user' AND COLUMN_NAME = 'deleted_at'
);
SET @qixi_admin_deleted_at_column_ddl := IF(
  @qixi_admin_deleted_at_column_exists = 0,
  'ALTER TABLE `qixi_crm_a_admin_user` ADD COLUMN `deleted_at` datetime DEFAULT NULL AFTER `updated_at`',
  'SELECT 1'
);
PREPARE qixi_admin_deleted_at_column_stmt FROM @qixi_admin_deleted_at_column_ddl;
EXECUTE qixi_admin_deleted_at_column_stmt;
DEALLOCATE PREPARE qixi_admin_deleted_at_column_stmt;

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
  `source_application_id` bigint unsigned DEFAULT NULL, `category_name` varchar(128) NOT NULL DEFAULT '', `merchant_type` varchar(64) NOT NULL DEFAULT '', `license_key` varchar(1024) NOT NULL DEFAULT '', `license_url` varchar(1024) NOT NULL DEFAULT '',
  `status` enum('draft','pending','approved','rejected') NOT NULL DEFAULT 'draft', `reviewed_by` bigint unsigned DEFAULT NULL,
  `review_note` varchar(500) NOT NULL DEFAULT '', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `reviewed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_source_application` (`source_application_id`), KEY `idx_region_status` (`region_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 兼容已初始化的本地/测试库。MySQL 8.4 不支持 ALTER TABLE ... IF NOT EXISTS，
-- 因此与上方菜单迁移一样使用 information_schema + 动态 SQL 保持可重复执行。
SET @qixi_application_source_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND COLUMN_NAME = 'source_application_id');
SET @qixi_application_source_ddl := IF(@qixi_application_source_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD COLUMN `source_application_id` bigint unsigned DEFAULT NULL', 'SELECT 1');
PREPARE qixi_application_source_stmt FROM @qixi_application_source_ddl;
EXECUTE qixi_application_source_stmt;
DEALLOCATE PREPARE qixi_application_source_stmt;
SET @qixi_application_category_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND COLUMN_NAME = 'category_name');
SET @qixi_application_category_ddl := IF(@qixi_application_category_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD COLUMN `category_name` varchar(128) NOT NULL DEFAULT ''''', 'SELECT 1');
PREPARE qixi_application_category_stmt FROM @qixi_application_category_ddl;
EXECUTE qixi_application_category_stmt;
DEALLOCATE PREPARE qixi_application_category_stmt;
SET @qixi_application_type_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND COLUMN_NAME = 'merchant_type');
SET @qixi_application_type_ddl := IF(@qixi_application_type_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD COLUMN `merchant_type` varchar(64) NOT NULL DEFAULT ''''', 'SELECT 1');
PREPARE qixi_application_type_stmt FROM @qixi_application_type_ddl;
EXECUTE qixi_application_type_stmt;
DEALLOCATE PREPARE qixi_application_type_stmt;
SET @qixi_application_license_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND COLUMN_NAME = 'license_url');
SET @qixi_application_license_ddl := IF(@qixi_application_license_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD COLUMN `license_url` varchar(1024) NOT NULL DEFAULT ''''', 'SELECT 1');
PREPARE qixi_application_license_stmt FROM @qixi_application_license_ddl;
EXECUTE qixi_application_license_stmt;
DEALLOCATE PREPARE qixi_application_license_stmt;
SET @qixi_application_license_key_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND COLUMN_NAME = 'license_key');
SET @qixi_application_license_key_ddl := IF(@qixi_application_license_key_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD COLUMN `license_key` varchar(1024) NOT NULL DEFAULT ''''', 'SELECT 1');
PREPARE qixi_application_license_key_stmt FROM @qixi_application_license_key_ddl;
EXECUTE qixi_application_license_key_stmt;
DEALLOCATE PREPARE qixi_application_license_key_stmt;
SET @qixi_application_source_index_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_merchant_application' AND INDEX_NAME = 'uk_source_application');
SET @qixi_application_source_index_ddl := IF(@qixi_application_source_index_exists = 0, 'ALTER TABLE `qixi_crm_a_merchant_application` ADD UNIQUE INDEX `uk_source_application` (`source_application_id`)', 'SELECT 1');
PREPARE qixi_application_source_index_stmt FROM @qixi_application_source_index_ddl;
EXECUTE qixi_application_source_index_stmt;
DEALLOCATE PREPARE qixi_application_source_index_stmt;
-- 平台监管只读取该投影，不直连 qixi_crm_merchant。由 api-merchant 的受控
-- 命令和 NATS 事件同步，保证服务与数据库边界独立。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_view` (
  `merchant_id` bigint unsigned NOT NULL, `merchant_name` varchar(128) NOT NULL,
  `contact_name` varchar(64) NOT NULL DEFAULT '', `contact_mobile` varchar(32) NOT NULL DEFAULT '',
  `region_id` bigint unsigned DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`merchant_id`), KEY `idx_region_status` (`region_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台侧的店铺结算监管投影。结算事实由店铺库产生并通过事件同步，平台不能直连
-- qixi_crm_merchant；本表只承担监管查询，不承担扣款、审批或打款。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_settlement_view` (
  `source_settlement_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '', `region_id` bigint unsigned DEFAULT NULL,
  `period_start` datetime NOT NULL, `period_end` datetime NOT NULL, `amount` decimal(16,2) NOT NULL,
  `status` enum('bill_pending','bill_frozen','withdraw_applied','approved','paid','rejected','cancelled') NOT NULL, `updated_at` datetime NOT NULL,
  PRIMARY KEY (`source_settlement_id`), KEY `idx_region_status_time` (`region_id`,`status`,`updated_at`), KEY `idx_merchant_period` (`merchant_id`,`period_start`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 兼容已初始化的历史投影；保留 cancelled 只用于历史展示，新增结算不会产生该状态。
ALTER TABLE `qixi_crm_a_merchant_settlement_view`
  MODIFY COLUMN `status` enum('bill_pending','bill_frozen','withdraw_applied','approved','paid','rejected','cancelled') NOT NULL;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `owner_scope` enum('platform','operations','merchant') NOT NULL,
  `owner_id` bigint unsigned NOT NULL DEFAULT 0, `media_type` enum('image','video','file') NOT NULL, `url` varchar(1024) NOT NULL,
  `name` varchar(255) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_owner` (`owner_scope`,`owner_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 统一后台素材库沿用 CRMEB 兼容字段，但使用独立 qixi_crm_a_ 表，
-- 使平台/运营素材与商户端素材边界可审计。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment_category` (
  `attachment_category_id` bigint unsigned NOT NULL AUTO_INCREMENT, `pid` bigint unsigned NOT NULL DEFAULT 0,
  `path` varchar(1000) NOT NULL DEFAULT '', `attachment_category_name` varchar(128) NOT NULL,
  `attachment_category_enname` varchar(128) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 0,
  `mer_id` bigint unsigned NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`attachment_category_id`), KEY `idx_owner_sort` (`mer_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment_asset` (
  `attachment_id` bigint unsigned NOT NULL AUTO_INCREMENT, `attachment_category_id` bigint unsigned NOT NULL DEFAULT 0,
  `attachment_name` varchar(255) NOT NULL, `attachment_src` varchar(1024) NOT NULL,
  `upload_type` tinyint NOT NULL DEFAULT 1, `user_type` int NOT NULL DEFAULT 0, `user_id` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `attachment_type` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`attachment_id`), KEY `idx_owner_category` (`user_type`,`attachment_category_id`,`attachment_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_content` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `content_type` enum('article','notice','agreement','banner') NOT NULL,
  `title` varchar(255) NOT NULL, `body` longtext NOT NULL, `status` enum('draft','published','hidden') NOT NULL DEFAULT 'draft',
  `published_at` datetime DEFAULT NULL, `created_by` bigint unsigned NOT NULL, PRIMARY KEY (`id`), KEY `idx_type_status` (`content_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 公告、协议和通道 stub 都是统一后台配置，不复用 CRMEB 旧缓存表。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_notice` (
  `notice_id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL, `content` longtext NOT NULL,
  `is_show` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`notice_id`), KEY `idx_visible_sort` (`is_del`,`is_show`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_setting_cache` (
  `key` varchar(128) NOT NULL, `expire_time` int NOT NULL DEFAULT 0, `result` longtext NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 文章与分类由统一后台直接维护；不再依赖 CRMEB 旧后台表。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_article_category` (
  `cid` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(128) NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`cid`), KEY `idx_visible_sort` (`is_del`,`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_article` (
  `article_id` bigint unsigned NOT NULL AUTO_INCREMENT, `cid` bigint unsigned NOT NULL DEFAULT 0,
  `title` varchar(255) NOT NULL, `author` varchar(128) NOT NULL DEFAULT '', `image` varchar(1024) NOT NULL DEFAULT '',
  `synopsis` varchar(1000) NOT NULL DEFAULT '', `content` longtext NOT NULL, `visit` int NOT NULL DEFAULT 0,
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, `is_del` tinyint NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`article_id`), KEY `idx_category_visible` (`cid`,`is_del`,`status`,`sort`)
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
-- 区域商圈和代理审核属于统一后台监管数据；代理结算资料不在接口响应中回传。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone` (
  `circle_id` bigint unsigned NOT NULL AUTO_INCREMENT, `pid` bigint unsigned NOT NULL DEFAULT 0,
  `path` varchar(255) NOT NULL DEFAULT '', `name` varchar(64) NOT NULL,
  `circle_agent_id` bigint unsigned NOT NULL DEFAULT 0, `commission_type` tinyint NOT NULL DEFAULT 0,
  `commission_rate` decimal(6,2) NOT NULL DEFAULT 0, `level` tinyint unsigned NOT NULL DEFAULT 0,
  `remark` varchar(500) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `type` tinyint NOT NULL DEFAULT 0, `role_id` bigint unsigned NOT NULL DEFAULT 0,
  `business_store_category` bigint unsigned NOT NULL DEFAULT 0, `business_store_type` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`circle_id`), KEY `idx_parent_type` (`pid`,`type`), KEY `idx_listing` (`status`,`sort`,`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone_agent` (
  `circle_agent_id` bigint unsigned NOT NULL AUTO_INCREMENT, `uid` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL, `phone` varchar(32) NOT NULL, `qualification` varchar(2000) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '', `audit_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `audit_reason` varchar(500) NOT NULL DEFAULT '', `audit_time` datetime DEFAULT NULL,
  `status` tinyint NOT NULL DEFAULT 0, `payment_method` tinyint unsigned NOT NULL DEFAULT 0,
  `payment_name` varchar(128) NOT NULL DEFAULT '', `payment_account` varchar(255) NOT NULL DEFAULT '',
  `payment_bank` varchar(255) NOT NULL DEFAULT '', `payment_qr_img` varchar(1024) NOT NULL DEFAULT '',
  `balance` decimal(12,2) NOT NULL DEFAULT 0, `type` tinyint NOT NULL DEFAULT 0,
  `business_name` varchar(128) NOT NULL DEFAULT '', `business_store_category` bigint unsigned NOT NULL DEFAULT 0,
  `business_store_type` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`circle_agent_id`), KEY `idx_status_listing` (`status`,`circle_agent_id`), KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 代理撤销不可硬删：保留资格状态与幂等审计，避免破坏已结算或关联区域的事实链路。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone_agent_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `circle_agent_id` bigint unsigned NOT NULL,
  `action` enum('revoke') NOT NULL, `from_status` tinyint NOT NULL, `to_status` tinyint NOT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_agent_action_idempotency` (`circle_agent_id`,`action`,`idempotency_key`), KEY `idx_agent_created` (`circle_agent_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 密码本身及其派生值均不落库；审计只保存代理、受影响的后台账号、原因、操作人和幂等键。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone_agent_password_reset_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `circle_agent_id` bigint unsigned NOT NULL,
  `admin_user_id` bigint unsigned NOT NULL,
  `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_agent_password_reset_idempotency` (`circle_agent_id`,`idempotency_key`),
  KEY `idx_agent_password_reset_created` (`circle_agent_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 快递公司、行政区划和平台运费模板均为全局监管配置，不向商户/区域角色开放直写。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_express` (
  `express_id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `code` varchar(64) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `is_show` tinyint NOT NULL DEFAULT 1, `is_del` tinyint NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`express_id`), UNIQUE KEY `uk_code` (`code`), KEY `idx_listing` (`is_del`,`is_show`,`sort`,`express_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_city` (
  `city_id` bigint unsigned NOT NULL, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `level` tinyint NOT NULL DEFAULT 0, `is_show` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`city_id`), KEY `idx_parent_visible` (`parent_id`,`is_show`,`city_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0, `name` varchar(128) NOT NULL,
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), KEY `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 商户入驻分类与商品分类不是同一事实。佣金比例由平台维护，不能复用商品类目表。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL,
  `commission_rate` decimal(5,2) NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_merchant_category_name` (`name`), KEY `idx_status_id` (`status`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 店铺类型定义入驻保证金规则及可见说明；不复用商品类目，授权菜单以统一后台 menu code 保存。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_type` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL,
  `type_info` varchar(500) NOT NULL DEFAULT '', `is_margin` tinyint NOT NULL DEFAULT 0,
  `margin` decimal(12,2) NOT NULL DEFAULT 0, `description` text NOT NULL,
  `remark` varchar(500) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_merchant_type_name` (`name`), KEY `idx_status_id` (`status`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_type_menu` (
  `merchant_type_id` bigint unsigned NOT NULL, `menu_code` varchar(128) NOT NULL,
  PRIMARY KEY (`merchant_type_id`,`menu_code`), KEY `idx_menu_code` (`menu_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 保证金余额、不可变流水和退款申请分表保存。所有余额改动必须由同一 admin 库事务持锁完成。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_deposit_account` (
  `merchant_id` bigint unsigned NOT NULL, `required_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `available_amount` decimal(12,2) NOT NULL DEFAULT 0, `state` enum('not_required','pending','funded','shortfall','refund_pending','refunded') NOT NULL DEFAULT 'not_required',
  `version` bigint unsigned NOT NULL DEFAULT 1, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`merchant_id`), KEY `idx_state` (`state`,`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_deposit_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL,
  `entry_type` enum('fund','deduct','refund_approved','refund_rejected','refund_paid') NOT NULL,
  `amount` decimal(12,2) NOT NULL, `balance_after` decimal(12,2) NOT NULL,
  `reason` varchar(500) NOT NULL DEFAULT '', `idempotency_key` varchar(128) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_merchant_idempotency` (`merchant_id`,`idempotency_key`), KEY `idx_merchant_time` (`merchant_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_deposit_refund` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL, `amount` decimal(12,2) NOT NULL,
  `status` enum('applied','approved','rejected','paid') NOT NULL DEFAULT 'applied',
  `reason` varchar(500) NOT NULL DEFAULT '', `review_note` varchar(500) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned NOT NULL DEFAULT 0, `reviewed_at` datetime DEFAULT NULL,
  `payout_idempotency_key` varchar(128) DEFAULT NULL, `payout_reference` varchar(128) DEFAULT NULL,
  `paid_by` bigint unsigned NOT NULL DEFAULT 0, `paid_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_payout_idempotency` (`merchant_id`,`payout_idempotency_key`), KEY `idx_merchant_status` (`merchant_id`,`status`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 商户分账申请为平台审核投影；仅记录申请标识、状态、说明和审核轨迹，不存渠道密钥或收款账户。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_profitsharing_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL,
  `application_no` varchar(64) NOT NULL, `status` enum('applied','approved','rejected') NOT NULL DEFAULT 'applied',
  `description` varchar(500) NOT NULL DEFAULT '', `review_note` varchar(500) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned NOT NULL DEFAULT 0, `reviewed_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_application_no` (`application_no`), KEY `idx_merchant_status` (`merchant_id`,`status`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_profitsharing_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `application_id` bigint unsigned NOT NULL,
  `from_status` varchar(16) NOT NULL, `to_status` varchar(16) NOT NULL, `note` varchar(500) NOT NULL DEFAULT '',
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_application_time` (`application_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 店铺分组是平台运营视图，不复用区域商圈或商户入驻分类。关联表不设跨库外键，
-- 商户合法性由统一后台 qixi_crm_a_merchant_view 投影在写入事务中校验。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_store_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `path` varchar(500) NOT NULL DEFAULT '/', `level` tinyint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `diy_page_id` bigint unsigned NOT NULL DEFAULT 0, `positioning_status` tinyint NOT NULL DEFAULT 0,
  `longitude` decimal(10,7) DEFAULT NULL, `latitude` decimal(10,7) DEFAULT NULL,
  `address` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_parent_name` (`parent_id`,`name`),
  KEY `idx_parent_sort` (`parent_id`,`sort`,`id`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_store_group_merchant` (
  `store_group_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_group_id`,`merchant_id`), KEY `idx_merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_brand` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `logo_url` varchar(1024) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_brand_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_parent_name` (`parent_id`,`name`),
  KEY `idx_parent_sort` (`parent_id`,`sort`,`id`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_crm_a_brand_category_id_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_platform_brand' AND COLUMN_NAME = 'category_id');
SET @qixi_crm_a_brand_category_id_ddl := IF(@qixi_crm_a_brand_category_id_exists = 0, 'ALTER TABLE `qixi_crm_a_platform_brand` ADD COLUMN `category_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `name`, ADD KEY `idx_category` (`category_id`)', 'SELECT 1');
PREPARE qixi_crm_a_brand_category_id_stmt FROM @qixi_crm_a_brand_category_id_ddl; EXECUTE qixi_crm_a_brand_category_id_stmt; DEALLOCATE PREPARE qixi_crm_a_brand_category_id_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_review` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `source_event_id` bigint unsigned DEFAULT NULL,
  `status` enum('pending','approved','rejected','off_sale') NOT NULL DEFAULT 'pending', `reason` varchar(1000) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `reviewed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_source_event` (`source_event_id`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_crm_a_product_review_source_event_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_product_review' AND COLUMN_NAME = 'source_event_id');
SET @qixi_crm_a_product_review_source_event_ddl := IF(@qixi_crm_a_product_review_source_event_exists = 0, 'ALTER TABLE `qixi_crm_a_product_review` ADD COLUMN `source_event_id` bigint unsigned DEFAULT NULL AFTER `store_id`', 'SELECT 1');
PREPARE qixi_crm_a_product_review_source_event_stmt FROM @qixi_crm_a_product_review_source_event_ddl; EXECUTE qixi_crm_a_product_review_source_event_stmt; DEALLOCATE PREPARE qixi_crm_a_product_review_source_event_stmt;
SET @qixi_crm_a_product_review_source_event_index_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_product_review' AND INDEX_NAME = 'uk_source_event');
SET @qixi_crm_a_product_review_source_event_index_ddl := IF(@qixi_crm_a_product_review_source_event_index_exists = 0, 'ALTER TABLE `qixi_crm_a_product_review` ADD UNIQUE KEY `uk_source_event` (`source_event_id`)', 'SELECT 1');
PREPARE qixi_crm_a_product_review_source_event_index_stmt FROM @qixi_crm_a_product_review_source_event_index_ddl; EXECUTE qixi_crm_a_product_review_source_event_index_stmt; DEALLOCATE PREPARE qixi_crm_a_product_review_source_event_index_stmt;
-- 商品审核跨库投影 outbox：审核事实与待投递命令必须同库提交，业务消费视图
-- 写入失败时保留 pending 以重试，禁止出现无法恢复的“已审核未投影”。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_projection_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint unsigned NOT NULL, `source_event_id` bigint unsigned DEFAULT NULL, `action` enum('upsert','delete') NOT NULL,
  `payload` json NOT NULL, `status` enum('pending','processing','published','failed') NOT NULL DEFAULT 'pending',
  `attempts` int unsigned NOT NULL DEFAULT 0, `last_error` varchar(500) NOT NULL DEFAULT '',
  `processed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_source_event` (`source_event_id`), UNIQUE KEY `uk_product_action` (`product_id`,`action`),
  KEY `idx_status_created` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 兼容已执行过早期建表脚本的环境：投影器以 processing 租约原子领取命令。
ALTER TABLE `qixi_crm_a_product_projection_outbox`
  MODIFY COLUMN `status` enum('pending','processing','published','failed') NOT NULL DEFAULT 'pending';
SET @qixi_crm_a_product_projection_source_event_exists := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_product_projection_outbox' AND COLUMN_NAME = 'source_event_id');
SET @qixi_crm_a_product_projection_source_event_ddl := IF(@qixi_crm_a_product_projection_source_event_exists = 0, 'ALTER TABLE `qixi_crm_a_product_projection_outbox` ADD COLUMN `source_event_id` bigint unsigned DEFAULT NULL AFTER `product_id`', 'SELECT 1');
PREPARE qixi_crm_a_product_projection_source_event_stmt FROM @qixi_crm_a_product_projection_source_event_ddl; EXECUTE qixi_crm_a_product_projection_source_event_stmt; DEALLOCATE PREPARE qixi_crm_a_product_projection_source_event_stmt;
SET @qixi_crm_a_product_projection_source_event_index_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'qixi_crm_admin' AND TABLE_NAME = 'qixi_crm_a_product_projection_outbox' AND INDEX_NAME = 'uk_source_event');
SET @qixi_crm_a_product_projection_source_event_index_ddl := IF(@qixi_crm_a_product_projection_source_event_index_exists = 0, 'ALTER TABLE `qixi_crm_a_product_projection_outbox` ADD UNIQUE KEY `uk_source_event` (`source_event_id`)', 'SELECT 1');
PREPARE qixi_crm_a_product_projection_source_event_index_stmt FROM @qixi_crm_a_product_projection_source_event_index_ddl; EXECUTE qixi_crm_a_product_projection_source_event_index_stmt; DEALLOCATE PREPARE qixi_crm_a_product_projection_source_event_index_stmt;
-- 平台商品元数据独立于店铺商品配置；只保存运营展示与筛选规则，不存商户私有规格。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_label` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT '', `color` varchar(32) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_guarantee` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `content` varchar(1000) NOT NULL DEFAULT '', `icon_url` varchar(1024) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_parameter_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `values_json` json NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_status_sort` (`status`,`sort`,`id`)
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
