CREATE DATABASE IF NOT EXISTS `qixi_crm_merchant` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `qixi_crm_merchant`;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_merchant` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(128) NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `region_id` bigint unsigned DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_store` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL, `app_id` varchar(64) NOT NULL, `name` varchar(128) NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_merchant` (`merchant_id`), UNIQUE KEY `uk_app_id` (`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_account` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `username` varchar(64) NOT NULL,
  `password_hash` varchar(255) NOT NULL, `role_code` enum('owner','manager','clerk','delivery','service') NOT NULL,
  `display_name` varchar(64) NOT NULL DEFAULT '', `phone` varchar(32) NOT NULL DEFAULT '',
  `can_accept_orders` tinyint NOT NULL DEFAULT 1, `can_verify_orders` tinyint NOT NULL DEFAULT 1, `can_ship_orders` tinyint NOT NULL DEFAULT 1,
  `status` tinyint NOT NULL DEFAULT 1, `auth_version` bigint unsigned NOT NULL DEFAULT 1, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_username` (`username`), KEY `idx_store_role` (`store_id`,`role_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_menu` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `code` varchar(128) NOT NULL, `name` varchar(64) NOT NULL, `path` varchar(255) NOT NULL,
  `component` varchar(255) NOT NULL DEFAULT '', `icon` varchar(128) NOT NULL DEFAULT '',
  `is_menu` tinyint NOT NULL DEFAULT 1, `is_route` tinyint NOT NULL DEFAULT 1,
  `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 已有库补列（local docker mysql 可重复执行）：
-- ALTER TABLE `qixi_crm_m_menu` ADD COLUMN IF NOT EXISTS `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_role_menu` (
  `role_code` varchar(64) NOT NULL, `menu_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_code`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `title` varchar(255) NOT NULL,
  `category_id` bigint unsigned DEFAULT NULL, `brand_name` varchar(64) NOT NULL DEFAULT '', `svip_price_type` tinyint NOT NULL DEFAULT 0, `svip_price` decimal(12,2) NOT NULL DEFAULT 0, `status` enum('draft','pending_review','on_sale','off_sale','rejected') NOT NULL DEFAULT 'draft',
  `version` bigint unsigned NOT NULL DEFAULT 1, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_sku` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `spec_json` json NOT NULL,
  `price` decimal(12,2) NOT NULL, `stock` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 商品状态与审核出站命令必须在商户库同一事务提交；平台库和业务库均为可重试投影。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_audit_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL, `action` enum('upsert','delete') NOT NULL,
  `review_status` enum('approved','rejected') NOT NULL, `reason` varchar(1000) NOT NULL DEFAULT '',
  `reviewed_by` bigint unsigned NOT NULL, `status` enum('pending','processing','published','failed') NOT NULL DEFAULT 'pending',
  `attempts` int unsigned NOT NULL DEFAULT 0, `last_error` varchar(500) NOT NULL DEFAULT '',
  `processed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_product_audit_event` (`product_id`),
  KEY `idx_status_created` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 商品主表保持审核与库存边界；面向编辑器的文案、单位和图片单独存放，避免把
-- CRMEB 的旧 qixi_m_admin_store_product 字段继续带入新模型。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_detail` (
  `product_id` bigint unsigned NOT NULL, `brief` varchar(2000) NOT NULL DEFAULT '',
  `keyword` varchar(255) NOT NULL DEFAULT '', `unit_name` varchar(32) NOT NULL DEFAULT '件',
  `cover_url` varchar(1024) NOT NULL DEFAULT '', `original_price` decimal(12,2) DEFAULT NULL,
  PRIMARY KEY (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_stock_reservation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `sku_id` bigint unsigned NOT NULL, `order_id` bigint unsigned NOT NULL,
  `quantity` int NOT NULL, `status` enum('reserved','confirmed','released') NOT NULL DEFAULT 'reserved',
  `expires_at` datetime NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_sku` (`order_id`,`sku_id`), KEY `idx_expire` (`status`,`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_fulfillment_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `assignee_account_id` bigint unsigned DEFAULT NULL, `task_type` enum('ship','verify','delivery','service') NOT NULL,
  `status` enum('pending','accepted','completed','cancelled') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_diy_page` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `document` json NOT NULL, `page_type` enum('home','custom') NOT NULL DEFAULT 'home',
  `is_active` tinyint NOT NULL DEFAULT 0, `status` enum('draft','published') NOT NULL DEFAULT 'draft',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_store_home_active` (`store_id`,`page_type`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS `qixi_crm_m_config` (
  `store_id` bigint unsigned NOT NULL, `config_key` varchar(128) NOT NULL, `config_value` json NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`,`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_payment_channel` (
  `store_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay') NOT NULL,
  `enabled` tinyint NOT NULL DEFAULT 0, `updated_by` bigint unsigned DEFAULT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`,`channel`), KEY `idx_channel_enabled` (`channel`,`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_payment_config` (
  `store_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay') NOT NULL,
  `ciphertext` longtext NOT NULL, `updated_by` bigint unsigned DEFAULT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`,`channel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_im_sdk_app` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL,
  `sdk_app_id` varchar(64) NOT NULL, `name` varchar(128) NOT NULL,
  `status` enum('enabled','disabled') NOT NULL DEFAULT 'disabled', `is_active` tinyint NOT NULL DEFAULT 0,
  `api_public_url` varchar(1024) NOT NULL DEFAULT '', `ws_public_url` varchar(1024) NOT NULL DEFAULT '',
  `pte_profile_id` varchar(128) NOT NULL DEFAULT '', `created_by` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_merchant_sdk_app` (`merchant_id`,`sdk_app_id`),
  KEY `idx_merchant_active` (`merchant_id`,`is_active`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_category` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL DEFAULT 0,
  `parent_id` bigint unsigned NOT NULL DEFAULT 0, `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), KEY `idx_store_parent` (`store_id`,`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_coupon` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `discount_type` enum('amount','rate') NOT NULL, `discount_value` decimal(12,2) NOT NULL, `min_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `total_quantity` int NOT NULL, `issued_quantity` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_marketing_activity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL,
  `activity_type` enum('seckill','combination','presell','assist','discount','svip') NOT NULL, `name` varchar(128) NOT NULL,
  `rules` json NOT NULL, `status` enum('draft','pending','active','closed','rejected') NOT NULL DEFAULT 'draft',
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_store_type_status` (`store_id`,`activity_type`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_shipping_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `rules` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), KEY `idx_store` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_finance_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `entry_type` varchar(64) NOT NULL,
  `amount` decimal(16,2) NOT NULL, `reference_type` varchar(64) NOT NULL, `reference_id` varchar(64) NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_entry_idempotency` (`entry_type`,`idempotency_key`), KEY `idx_store_time` (`store_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 订单级结算账本：每一笔确认收货或可信退款均独立留痕。账期单只汇总未冻结账期，
-- 已冻结/已打款账期绝不可因退款被回写，退款会进入当期的负向调整。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_settlement_entry` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned NOT NULL, `refund_id` bigint unsigned NOT NULL DEFAULT 0,
  `entry_type` enum('order_accrual','refund_reversal') NOT NULL, `amount` decimal(16,2) NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `occurred_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_entry_idempotency` (`idempotency_key`),
  KEY `idx_store_order` (`store_id`,`order_id`), KEY `idx_merchant_time` (`merchant_id`,`occurred_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 店铺结算事实。平台只消费 outbox 投影，不得跨库直接读写本表。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_settlement_bill` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL,
  `period_start` datetime NOT NULL, `period_end` datetime NOT NULL, `amount` decimal(16,2) NOT NULL,
  `status` enum('bill_pending','bill_frozen','withdraw_applied','approved','paid','rejected') NOT NULL DEFAULT 'bill_pending',
  `idempotency_key` varchar(128) DEFAULT NULL, `application_no` varchar(64) DEFAULT NULL,
  `applied_by_account_id` bigint unsigned DEFAULT NULL, `applied_at` datetime DEFAULT NULL,
  `reviewed_by_admin_id` bigint unsigned DEFAULT NULL, `review_idempotency_key` varchar(128) DEFAULT NULL,
  `review_note` varchar(500) NOT NULL DEFAULT '', `reviewed_at` datetime DEFAULT NULL,
  `payout_idempotency_key` varchar(128) DEFAULT NULL, `payout_reference` varchar(128) DEFAULT NULL, `paid_at` datetime DEFAULT NULL,
  `version` bigint unsigned NOT NULL DEFAULT 1, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_store_apply_key` (`store_id`,`idempotency_key`),
  UNIQUE KEY `uk_store_review_key` (`store_id`,`review_idempotency_key`), UNIQUE KEY `uk_store_payout_key` (`store_id`,`payout_idempotency_key`),
  KEY `idx_store_status_time` (`store_id`,`status`,`updated_at`), KEY `idx_merchant_period` (`merchant_id`,`period_start`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `event_type` varchar(128) NOT NULL, `aggregate_type` varchar(64) NOT NULL,
  `aggregate_id` varchar(64) NOT NULL, `payload` json NOT NULL, `status` enum('pending','published','failed') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `published_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_status_time` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_store_integral_policy` (
  `store_id` bigint unsigned NOT NULL, `enabled` tinyint NOT NULL DEFAULT 0,
  `points_per_yuan` bigint NOT NULL DEFAULT 100, `max_deduction_bps` int NOT NULL DEFAULT 2000,
  `updated_by_account_id` bigint unsigned NOT NULL DEFAULT 0, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 店铺组织、资质与员工权限。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_account_role` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `code` varchar(64) NOT NULL,
  `name` varchar(64) NOT NULL, `permissions` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_store_code` (`store_id`,`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_account_role_binding` (
  `account_id` bigint unsigned NOT NULL, `role_id` bigint unsigned NOT NULL, PRIMARY KEY (`account_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_merchant_qualification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL, `qualification_type` varchar(64) NOT NULL,
  `document_no` varchar(128) NOT NULL DEFAULT '', `document_url` varchar(1024) NOT NULL, `expires_at` datetime DEFAULT NULL,
  `status` enum('pending','approved','rejected','expired') NOT NULL DEFAULT 'pending', `review_note` varchar(500) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_merchant_status` (`merchant_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_store_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `address_type` enum('ship_from','return','pickup') NOT NULL,
  `contact_name` varchar(64) NOT NULL, `mobile` varchar(32) NOT NULL, `region_code` varchar(32) NOT NULL, `detail` varchar(255) NOT NULL,
  `is_default` tinyint NOT NULL DEFAULT 0, PRIMARY KEY (`id`), KEY `idx_store_type` (`store_id`,`address_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_operation_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `account_id` bigint unsigned NOT NULL,
  `action` varchar(128) NOT NULL, `resource_type` varchar(64) NOT NULL, `resource_id` varchar(64) NOT NULL DEFAULT '',
  `request_id` varchar(64) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_store_time` (`store_id`,`created_at`), KEY `idx_request` (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 商品、库存与素材域。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_brand` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `logo_url` varchar(1024) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), UNIQUE KEY `uk_store_name` (`store_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_media` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `media_type` enum('image','video') NOT NULL,
  `url` varchar(1024) NOT NULL, `sort` int NOT NULL DEFAULT 0, PRIMARY KEY (`id`), KEY `idx_product_sort` (`product_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_attribute` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `attribute_name` varchar(128) NOT NULL,
  `attribute_value` json NOT NULL, PRIMARY KEY (`id`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_tag` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `name` varchar(64) NOT NULL,
  `info` varchar(255) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE KEY `uk_store_tag` (`store_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_tag_binding` (
  `product_id` bigint unsigned NOT NULL, `tag_id` bigint unsigned NOT NULL, PRIMARY KEY (`product_id`,`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_stock_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `sku_id` bigint unsigned NOT NULL, `change_quantity` int NOT NULL,
  `balance_quantity` int NOT NULL, `reason_type` varchar(64) NOT NULL, `reference_type` varchar(64) NOT NULL,
  `reference_id` varchar(64) NOT NULL, `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_stock_idempotency` (`sku_id`,`idempotency_key`), KEY `idx_sku_time` (`sku_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_recycle_bin` (
  `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `deleted_by_account_id` bigint unsigned NOT NULL,
  `deleted_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `restore_until` datetime NOT NULL, PRIMARY KEY (`product_id`), KEY `idx_store_deleted` (`store_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_attachment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `media_type` enum('image','video','file') NOT NULL,
  `name` varchar(255) NOT NULL, `url` varchar(1024) NOT NULL, `size_bytes` bigint unsigned NOT NULL DEFAULT 0,
  `created_by_account_id` bigint unsigned NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_store_time` (`store_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 物流、履约、售后与店铺财务域。业务订单只保存 ID，绝不跨库建外键。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_delivery_provider` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `provider_code` varchar(64) NOT NULL,
  `name` varchar(128) NOT NULL, `config` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_store_provider` (`store_id`,`provider_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_shipment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `order_id` bigint unsigned NOT NULL,
  `shipping_template_id` bigint unsigned DEFAULT NULL, `provider_code` varchar(64) NOT NULL DEFAULT '', `tracking_no` varchar(128) NOT NULL DEFAULT '',
  `status` enum('pending','shipped','delivered','returned','cancelled') NOT NULL DEFAULT 'pending', `shipped_by_account_id` bigint unsigned DEFAULT NULL,
  `shipped_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order` (`order_id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_order_action` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `account_id` bigint unsigned NOT NULL, `action` varchar(64) NOT NULL, `payload` json NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_order_time` (`order_id`,`created_at`), KEY `idx_store_time` (`store_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_aftersale_action` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `account_id` bigint unsigned NOT NULL, `action` varchar(64) NOT NULL, `note` varchar(2000) NOT NULL DEFAULT '',
  `attachments` json NOT NULL, `idempotency_key` varchar(128) DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_refund_time` (`refund_id`,`created_at`), UNIQUE KEY `uk_refund_action_idempotency` (`refund_id`,`store_id`,`account_id`,`action`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- CRMEB 的“删除退款单”在本系统是店铺视图隐藏：不能删除业务退款、退款交易或平台监管事实。
CREATE TABLE IF NOT EXISTS `qixi_crm_m_refund_hidden` (
  `refund_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `deleted_by_account_id` bigint unsigned NOT NULL,
  `reason` varchar(500) NOT NULL, `idempotency_key` varchar(128) NOT NULL, `deleted_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`refund_id`), UNIQUE KEY `uk_store_idempotency` (`store_id`,`idempotency_key`), KEY `idx_store_deleted` (`store_id`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_withdrawal_view` (
  `withdrawal_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `amount` decimal(12,2) NOT NULL,
  `status` varchar(32) NOT NULL, `updated_at` datetime NOT NULL, PRIMARY KEY (`withdrawal_id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_m_settlement_view` (
  `settlement_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `period_start` datetime NOT NULL, `period_end` datetime NOT NULL,
  `amount` decimal(12,2) NOT NULL, `status` varchar(32) NOT NULL, `updated_at` datetime NOT NULL,
  PRIMARY KEY (`settlement_id`), KEY `idx_store_period` (`store_id`,`period_start`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
