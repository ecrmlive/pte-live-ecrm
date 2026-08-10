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

-- 逻辑删除用于保留后台操作、客服转接和资金审批的历史归属；已删除账号
-- 不能重新登录，也不会再出现在可管理账号或客服转接坐席中。

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
-- 平台监管只读取该投影，不直连 qixi_crm_merchant。由 api-merchant 的受控
-- 命令和 NATS 事件同步，保证服务与数据库边界独立。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_view` (
  `merchant_id` bigint unsigned NOT NULL, `merchant_name` varchar(128) NOT NULL,
  `owner_name` varchar(128) NOT NULL DEFAULT '', `contact_name` varchar(64) NOT NULL DEFAULT '',
  `contact_mobile` varchar(32) NOT NULL DEFAULT '', `address` varchar(255) NOT NULL DEFAULT '',
  `category_id` bigint unsigned NOT NULL DEFAULT 0, `type_id` bigint unsigned NOT NULL DEFAULT 0,
  `business_id` bigint unsigned NOT NULL DEFAULT 0,
  `region_id` bigint unsigned DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `is_best` tinyint NOT NULL DEFAULT 0, `offline_pay` tinyint NOT NULL DEFAULT 0,
  `is_trader` tinyint NOT NULL DEFAULT 0, `is_audit` tinyint NOT NULL DEFAULT 1,
  `is_bro_room` tinyint NOT NULL DEFAULT 0, `is_bro_goods` tinyint NOT NULL DEFAULT 0,
  `commission_switch` tinyint NOT NULL DEFAULT 0, `commission_rate` decimal(5,2) NOT NULL DEFAULT 0,
  `mer_keyword` varchar(255) NOT NULL DEFAULT '', `mer_info` varchar(1000) NOT NULL DEFAULT '',
  `mer_account` varchar(64) NOT NULL DEFAULT '', `sub_mchid` varchar(64) NOT NULL DEFAULT '',
  `applyment_id` varchar(64) NOT NULL DEFAULT '',
  `care_count` int NOT NULL DEFAULT 0, `care_ficti` int NOT NULL DEFAULT 0,
  `goods_type` varchar(64) NOT NULL DEFAULT '', `platform_category_ids` varchar(500) NOT NULL DEFAULT '',
  `mer_star` tinyint NOT NULL DEFAULT 5,
  `mer_avatar` varchar(1024) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `mark` varchar(500) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`merchant_id`), KEY `idx_region_status` (`region_id`,`status`),
  KEY `idx_category_type` (`category_id`,`type_id`), KEY `idx_best_sort` (`is_best`,`sort`,`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- CREATE IF NOT EXISTS 不会给旧表补列；幂等对齐 CRMEB 店铺列表/编辑抽屉监管字段。
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='owner_name')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `owner_name` varchar(128) NOT NULL DEFAULT '''' AFTER `merchant_name`, ADD COLUMN `address` varchar(255) NOT NULL DEFAULT '''' AFTER `contact_mobile`, ADD COLUMN `category_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `address`, ADD COLUMN `type_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `category_id`, ADD COLUMN `is_best` tinyint NOT NULL DEFAULT 0 AFTER `status`, ADD COLUMN `offline_pay` tinyint NOT NULL DEFAULT 0 AFTER `is_best`, ADD COLUMN `sort` int NOT NULL DEFAULT 0 AFTER `offline_pay`, ADD COLUMN `mark` varchar(500) NOT NULL DEFAULT '''' AFTER `sort`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='business_id')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `business_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `type_id`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='mer_keyword')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `is_trader` tinyint NOT NULL DEFAULT 0 AFTER `offline_pay`, ADD COLUMN `is_audit` tinyint NOT NULL DEFAULT 1 AFTER `is_trader`, ADD COLUMN `is_bro_room` tinyint NOT NULL DEFAULT 0 AFTER `is_audit`, ADD COLUMN `is_bro_goods` tinyint NOT NULL DEFAULT 0 AFTER `is_bro_room`, ADD COLUMN `commission_switch` tinyint NOT NULL DEFAULT 0 AFTER `is_bro_goods`, ADD COLUMN `commission_rate` decimal(5,2) NOT NULL DEFAULT 0 AFTER `commission_switch`, ADD COLUMN `mer_keyword` varchar(255) NOT NULL DEFAULT '''' AFTER `commission_rate`, ADD COLUMN `mer_info` varchar(1000) NOT NULL DEFAULT '''' AFTER `mer_keyword`, ADD COLUMN `mer_account` varchar(64) NOT NULL DEFAULT '''' AFTER `mer_info`, ADD COLUMN `sub_mchid` varchar(64) NOT NULL DEFAULT '''' AFTER `mer_account`, ADD COLUMN `applyment_id` varchar(64) NOT NULL DEFAULT '''' AFTER `sub_mchid`, ADD COLUMN `care_count` int NOT NULL DEFAULT 0 AFTER `applyment_id`, ADD COLUMN `care_ficti` int NOT NULL DEFAULT 0 AFTER `care_count`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='goods_type')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `goods_type` varchar(64) NOT NULL DEFAULT '''' AFTER `care_ficti`, ADD COLUMN `platform_category_ids` varchar(500) NOT NULL DEFAULT '''' AFTER `goods_type`, ADD COLUMN `mer_star` tinyint NOT NULL DEFAULT 5 AFTER `platform_category_ids`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='mer_avatar')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `mer_avatar` varchar(1024) NOT NULL DEFAULT '''' AFTER `mer_star`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
-- 平台侧的店铺结算监管投影。结算事实由店铺库产生并通过事件同步，平台不能直连
-- qixi_crm_merchant；本表只承担监管查询，不承担扣款、审批或打款。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_settlement_view` (
  `source_settlement_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '', `region_id` bigint unsigned DEFAULT NULL,
  `period_start` datetime NOT NULL, `period_end` datetime NOT NULL, `amount` decimal(16,2) NOT NULL,
  `status` enum('bill_pending','bill_frozen','withdraw_applied','approved','paid','rejected','cancelled') NOT NULL, `updated_at` datetime NOT NULL,
  PRIMARY KEY (`source_settlement_id`), KEY `idx_region_status_time` (`region_id`,`status`,`updated_at`), KEY `idx_merchant_period` (`merchant_id`,`period_start`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
  `mer_id` bigint unsigned NOT NULL DEFAULT 0, `is_system` tinyint NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`attachment_category_id`), KEY `idx_owner_sort` (`mer_id`,`sort`),
  KEY `idx_owner_enname` (`mer_id`,`attachment_category_enname`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_attachment_asset` (
  `attachment_id` bigint unsigned NOT NULL AUTO_INCREMENT, `attachment_category_id` bigint unsigned NOT NULL DEFAULT 0,
  `attachment_name` varchar(255) NOT NULL, `attachment_src` varchar(1024) NOT NULL,
  `upload_type` tinyint NOT NULL DEFAULT 1, `user_type` int NOT NULL DEFAULT 0, `user_id` bigint unsigned NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `attachment_type` tinyint NOT NULL DEFAULT 0,
  `is_system` tinyint NOT NULL DEFAULT 0 COMMENT '1=系统预置素材（侧栏系统素材）',
  PRIMARY KEY (`attachment_id`), KEY `idx_owner_category` (`user_type`,`attachment_category_id`,`attachment_type`),
  KEY `idx_owner_system` (`user_type`,`is_system`,`attachment_type`)
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
  `goods_type` varchar(64) NOT NULL DEFAULT '', `platform_category_ids` varchar(500) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`circle_id`), KEY `idx_parent_type` (`pid`,`type`), KEY `idx_listing` (`status`,`sort`,`circle_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_business_zone' AND COLUMN_NAME='goods_type')=0,
    'ALTER TABLE `qixi_crm_a_business_zone` ADD COLUMN `goods_type` varchar(64) NOT NULL DEFAULT '''' AFTER `business_store_type`, ADD COLUMN `platform_category_ids` varchar(500) NOT NULL DEFAULT '''' AFTER `goods_type`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_business_zone_agent` (
  `circle_agent_id` bigint unsigned NOT NULL AUTO_INCREMENT, `uid` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL, `phone` varchar(32) NOT NULL, `qualification` varchar(2000) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '', `extend` text,
  `audit_admin_id` bigint unsigned NOT NULL DEFAULT 0,
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
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_business_zone_agent' AND COLUMN_NAME='extend')=0,
    'ALTER TABLE `qixi_crm_a_business_zone_agent` ADD COLUMN `extend` text NULL AFTER `remark`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
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
  `pic` varchar(1024) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `is_hot` tinyint NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 既有库补齐商品分类图标 / 推荐 / 创建时间（幂等）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='pic')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `pic` varchar(1024) NOT NULL DEFAULT '''' AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='is_hot')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `is_hot` tinyint NOT NULL DEFAULT 0 AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_category' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_category` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `is_hot`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
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
-- 商户分账申请为平台审核投影；对齐 CRMEB eb_merchant_applyments 列表字段，不存渠道密钥或收款账户。
-- status: applied 待审核 / platform_rejected 平台驳回 / auditing 审核中 / shop_verify 店铺验证 /
--         completed 已完成 / frozen 已冻结 / wechat_rejected 微信驳回（保留 approved/rejected 兼容旧夹具）
CREATE TABLE IF NOT EXISTS `qixi_crm_a_merchant_profitsharing_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '',
  `application_no` varchar(64) NOT NULL COMMENT '业务申请编号',
  `applyment_id` varchar(100) NOT NULL DEFAULT '' COMMENT '微信支付申请单号',
  `status` enum('applied','approved','rejected','platform_rejected','auditing','shop_verify','completed','frozen','wechat_rejected') NOT NULL DEFAULT 'applied',
  `description` varchar(500) NOT NULL DEFAULT '',
  `message` varchar(1000) NOT NULL DEFAULT '' COMMENT '审核结果',
  `review_note` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `reviewed_by` bigint unsigned NOT NULL DEFAULT 0, `reviewed_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_application_no` (`application_no`), KEY `idx_merchant_status` (`merchant_id`,`status`,`id`),
  KEY `idx_merchant_name` (`merchant_name`), KEY `idx_applyment_id` (`applyment_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 幂等补齐旧库列与状态枚举（CREATE IF NOT EXISTS 不会改已有表）。
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_profitsharing_application' AND COLUMN_NAME='applyment_id')=0,
    'ALTER TABLE `qixi_crm_a_merchant_profitsharing_application` ADD COLUMN `merchant_name` varchar(128) NOT NULL DEFAULT '''' AFTER `merchant_id`, ADD COLUMN `applyment_id` varchar(100) NOT NULL DEFAULT '''' AFTER `application_no`, ADD COLUMN `message` varchar(1000) NOT NULL DEFAULT '''' AFTER `description`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
ALTER TABLE `qixi_crm_a_merchant_profitsharing_application`
  MODIFY COLUMN `status` enum('applied','approved','rejected','platform_rejected','auditing','shop_verify','completed','frozen','wechat_rejected') NOT NULL DEFAULT 'applied';
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
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `category_id` bigint unsigned NOT NULL DEFAULT 0, `logo_url` varchar(1024) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_category` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 既有库补齐品牌创建时间（幂等）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_brand' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_brand` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_platform_brand_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_parent_name` (`parent_id`,`name`),
  KEY `idx_parent_sort` (`parent_id`,`sort`,`id`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 既有库补齐品牌分类创建时间（幂等）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_brand_category' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_brand_category` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_review` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `source_event_id` bigint unsigned DEFAULT NULL,
  `status` enum('pending','approved','rejected','off_sale') NOT NULL DEFAULT 'pending', `reason` varchar(1000) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `reviewed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_source_event` (`source_event_id`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
-- 平台对店铺商品的运营字段（显示/星级/排序/推荐），不改写商户库商品主状态。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_ops` (
  `product_id` bigint unsigned NOT NULL,
  `is_used` tinyint NOT NULL DEFAULT 1 COMMENT '平台是否显示',
  `star` tinyint NOT NULL DEFAULT 0 COMMENT '推荐级别 0-5',
  `rank_sort` int NOT NULL DEFAULT 0 COMMENT '平台排序',
  `is_hot` tinyint NOT NULL DEFAULT 0,
  `is_best` tinyint NOT NULL DEFAULT 0,
  `is_benefit` tinyint NOT NULL DEFAULT 0,
  `is_new` tinyint NOT NULL DEFAULT 0,
  `cate_hot` tinyint NOT NULL DEFAULT 0,
  `sys_labels` varchar(500) NOT NULL DEFAULT '',
  `content_html` mediumtext COMMENT '平台侧商品详情/营销详情 HTML',
  `refund_switch` tinyint NOT NULL DEFAULT 1 COMMENT '支持退款',
  `once_min_count` int NOT NULL DEFAULT 1 COMMENT '最少购买件数',
  `ficti` int NOT NULL DEFAULT 0 COMMENT '虚拟已售数量',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`), KEY `idx_used_star` (`is_used`,`star`,`rank_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_ops' AND COLUMN_NAME='content_html')=0,
    'ALTER TABLE `qixi_crm_a_product_ops` ADD COLUMN `content_html` mediumtext NULL AFTER `sys_labels`, ADD COLUMN `refund_switch` tinyint NOT NULL DEFAULT 1 AFTER `content_html`, ADD COLUMN `once_min_count` int NOT NULL DEFAULT 1 AFTER `refund_switch`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_ops' AND COLUMN_NAME='ficti')=0,
    'ALTER TABLE `qixi_crm_a_product_ops` ADD COLUMN `ficti` int NOT NULL DEFAULT 0 COMMENT ''虚拟已售数量'' AFTER `once_min_count`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
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
  `mer_count` int NOT NULL DEFAULT 0 COMMENT '使用的店铺数（计数缓存，对齐 CRMEB mer_count）',
  `product_count` int NOT NULL DEFAULT 0 COMMENT '使用商品数（计数缓存，对齐 CRMEB product_cout）',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_parameter_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '模板名称',
  `cate_ids_json` json NOT NULL COMMENT '关联平台分类 ID 列表',
  `params_json` json NOT NULL COMMENT '参数项 [{name,values,required,sort}]',
  `values_json` json NOT NULL COMMENT '兼容旧候选值数组（首参数 values 快照）',
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_name` (`name`), KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_price_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '名称',
  `cate_ids_json` json NULL COMMENT '关联平台分类 ID 列表；空=全部商品',
  `is_default` tinyint NOT NULL DEFAULT 1 COMMENT '1=未选分类默认全部商品',
  `content` mediumtext NULL COMMENT '价格说明详情 HTML',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示 1显示 0隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_status_sort` (`status`,`sort`,`id`), KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 既有库补齐平台参数模板分类/参数项（幂等；完整迁移见 patch_product_parameter_template_crmeb.sql）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_parameter_template' AND COLUMN_NAME='cate_ids_json')=0,
    'ALTER TABLE `qixi_crm_a_product_parameter_template` ADD COLUMN `cate_ids_json` json NULL COMMENT ''关联平台分类 ID 列表'' AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_parameter_template' AND COLUMN_NAME='params_json')=0,
    'ALTER TABLE `qixi_crm_a_product_parameter_template` ADD COLUMN `params_json` json NULL COMMENT ''参数项 [{name,values,required,sort}]'' AFTER `cate_ids_json`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_a_marketing_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `rule_type` enum('coupon','seckill','combination','presell','assist','svip','distribution') NOT NULL,
  `name` varchar(128) NOT NULL, `rule` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, `updated_by` bigint unsigned NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_type_status_time` (`rule_type`,`status`,`starts_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台营销装饰（氛围图/边框图/专题；报名已迁至 qixi_crm_a_signup_activity），取代 setting_cache list stub；不含密钥。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_marketing_decor` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `decor_type` enum('atmosphere','border','topic','application') NOT NULL,
  `name` varchar(128) NOT NULL, `code` varchar(64) NOT NULL DEFAULT '',
  `cover_url` varchar(1024) NOT NULL DEFAULT '', `remark` varchar(500) NOT NULL DEFAULT '',
  `payload` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_type_visible` (`decor_type`,`is_del`,`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台报名活动（对齐 CRMEB store_activity form）
CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_activity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL,
  `info` varchar(500) NOT NULL DEFAULT '',
  `cover_url` varchar(1024) NOT NULL DEFAULT '',
  `poster_url` varchar(1024) NOT NULL DEFAULT '',
  `color` varchar(32) NOT NULL DEFAULT '',
  `form_id` bigint unsigned NOT NULL DEFAULT 0,
  `quota` int unsigned NOT NULL DEFAULT 0,
  `total` int unsigned NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `sort` int NOT NULL DEFAULT 0,
  `starts_at` datetime DEFAULT NULL,
  `ends_at` datetime DEFAULT NULL,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_form` (`form_id`),
  KEY `idx_visible` (`is_del`,`status`,`sort`),
  KEY `idx_time` (`starts_at`,`ends_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台报名活动用户记录
CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL DEFAULT 0,
  `nickname` varchar(128) NOT NULL DEFAULT '',
  `mobile` varchar(32) NOT NULL DEFAULT '',
  `avatar` varchar(1024) NOT NULL DEFAULT '',
  `form_value` json NOT NULL,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_activity` (`activity_id`,`is_del`,`created_at`),
  KEY `idx_user` (`user_id`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台配置条目（热门搜索/组合数据/系统表单/备份登记），取代 setting_cache list stub。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_config_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `item_type` enum('hot_search','group_data','system_form','backup') NOT NULL,
  `name` varchar(128) NOT NULL, `code` varchar(64) NOT NULL DEFAULT '',
  `remark` varchar(500) NOT NULL DEFAULT '', `payload` json NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_type_visible` (`item_type`,`is_del`,`status`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台报名活动（对齐 CRMEB store_activity form）+ 用户报名记录
CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_activity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL COMMENT '活动名称',
  `info` varchar(500) NOT NULL DEFAULT '' COMMENT '活动简介',
  `cover_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '封面图 750*350',
  `poster_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '分享海报 750*1250',
  `color` varchar(32) NOT NULL DEFAULT '' COMMENT '活动背景色',
  `form_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '关联系统表单',
  `quota` int unsigned NOT NULL DEFAULT 0 COMMENT '人数上限，0=不限制',
  `total` int unsigned NOT NULL DEFAULT 0 COMMENT '已报名人数',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示/开启',
  `sort` int NOT NULL DEFAULT 0,
  `starts_at` datetime DEFAULT NULL,
  `ends_at` datetime DEFAULT NULL,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_form` (`form_id`),
  KEY `idx_visible` (`is_del`,`status`,`sort`),
  KEY `idx_time` (`starts_at`,`ends_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台报名活动';
CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` bigint unsigned NOT NULL,
  `user_id` bigint unsigned NOT NULL DEFAULT 0,
  `nickname` varchar(128) NOT NULL DEFAULT '',
  `mobile` varchar(32) NOT NULL DEFAULT '',
  `avatar` varchar(1024) NOT NULL DEFAULT '',
  `form_value` json NOT NULL,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_activity` (`activity_id`,`is_del`,`created_at`),
  KEY `idx_user` (`user_id`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台报名活动用户记录';
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
