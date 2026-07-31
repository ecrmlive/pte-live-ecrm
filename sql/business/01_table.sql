CREATE DATABASE IF NOT EXISTS `qixi_crm_business` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `nickname` varchar(64) NOT NULL DEFAULT '',
  `mobile` varchar(32) DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1, `auth_version` bigint unsigned NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_identity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `channel` enum('wechat','mini_program','h5','pc','ios','android','harmony') NOT NULL, `subject` varchar(191) NOT NULL,
  `credential_hash` varchar(255) DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_channel_subject` (`channel`,`subject`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- pte-tools-captcha 校验令牌仅保存 SHA-256 摘要；令牌只可按用途消费一次。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_auth_captcha_token` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `token_hash` char(64) NOT NULL,
  `action` enum('login_password','login_sms','register') NOT NULL,
  `expires_at` datetime NOT NULL,
  `consumed_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_token_hash` (`token_hash`),
  KEY `idx_action_available` (`action`,`consumed_at`,`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `recipient` varchar(64) NOT NULL,
  `mobile` varchar(32) NOT NULL, `province` varchar(64) NOT NULL DEFAULT '', `city` varchar(64) NOT NULL DEFAULT '',
  `district` varchar(64) NOT NULL DEFAULT '', `region_code` varchar(32) NOT NULL DEFAULT '', `detail` varchar(255) NOT NULL,
  `post_code` int unsigned NOT NULL DEFAULT 0, `is_default` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_view` (
  `product_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '', `store_name` varchar(128) NOT NULL DEFAULT '',
  `category_id` bigint unsigned NOT NULL DEFAULT 0, `title` varchar(255) NOT NULL, `cover_url` varchar(1024) NOT NULL DEFAULT '',
  `price` decimal(12,2) NOT NULL, `original_price` decimal(12,2) DEFAULT NULL, `product_type` tinyint NOT NULL DEFAULT 0,
  `sales` int NOT NULL DEFAULT 0, `stock` int NOT NULL, `sale_status` tinyint NOT NULL, `version` bigint unsigned NOT NULL,
  `updated_at` datetime NOT NULL, PRIMARY KEY (`product_id`), KEY `idx_merchant_sale` (`merchant_id`,`sale_status`), KEY `idx_store_sale` (`store_id`,`sale_status`), KEY `idx_category_sale` (`category_id`,`sale_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_category_view` (
  `category_id` bigint unsigned NOT NULL, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `sort` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`category_id`), KEY `idx_parent_sort` (`parent_id`,`sort`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_cart` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL, `sku_key` varchar(128) NOT NULL, `quantity` int NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_product_sku` (`user_id`,`product_id`,`sku_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_group_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `total_amount` decimal(12,2) NOT NULL, `discount_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `freight_amount` decimal(12,2) NOT NULL DEFAULT 0, `pay_amount` decimal(12,2) NOT NULL,
  `total_quantity` int unsigned NOT NULL, `recipient_snapshot` json NOT NULL,
  `pay_channel` enum('balance','wechat','alipay','mock') DEFAULT NULL,
  `pay_status` enum('pending','paid','closed','refunded') NOT NULL DEFAULT 'pending', `paid_at` datetime DEFAULT NULL,
  `idempotency_key` varchar(128) NOT NULL, `remark` varchar(500) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_no` (`order_no`), UNIQUE KEY `uk_user_idempotency` (`user_id`,`idempotency_key`), KEY `idx_user_status_time` (`user_id`,`pay_status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_order_id` bigint unsigned NOT NULL, `order_no` varchar(64) NOT NULL,
  `merchant_id` bigint unsigned NOT NULL, `merchant_name_snapshot` varchar(128) NOT NULL, `store_id` bigint unsigned NOT NULL,
  `store_name_snapshot` varchar(128) NOT NULL, `user_id` bigint unsigned NOT NULL, `total_amount` decimal(12,2) NOT NULL,
  `discount_amount` decimal(12,2) NOT NULL DEFAULT 0, `freight_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `pay_amount` decimal(12,2) NOT NULL, `total_quantity` int unsigned NOT NULL,
  `recipient_snapshot` json NOT NULL, `remark` varchar(500) NOT NULL DEFAULT '',
  `status` enum('pending_pay','paid','fulfilling','shipped','completed','cancelled','aftersale') NOT NULL DEFAULT 'pending_pay',
  `paid_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_no` (`order_no`), KEY `idx_merchant_status` (`merchant_id`,`status`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `refund_no` varchar(64) NOT NULL,
  `reason` varchar(500) NOT NULL, `amount` decimal(12,2) NOT NULL,
  `order_status_before` enum('pending_pay','paid','fulfilling','shipped','completed','cancelled','aftersale') NOT NULL,
  `status` enum('applied','merchant_handling','platform_intervene','refunding','refunded','rejected','cancelled') NOT NULL DEFAULT 'applied',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_no` (`refund_no`), UNIQUE KEY `uk_order_idempotency` (`order_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_config` (
  `config_key` varchar(128) NOT NULL, `config_value` json NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`config_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_payment_channel` (
  `store_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay') NOT NULL,
  `enabled` tinyint NOT NULL DEFAULT 0, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`,`channel`), KEY `idx_channel_enabled` (`channel`,`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_favorite` (
  `user_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`,`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_search_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `keyword` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_user_time` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `sku_key` varchar(128) NOT NULL, `title_snapshot` varchar(255) NOT NULL, `cover_url_snapshot` varchar(1024) NOT NULL DEFAULT '',
  `spec_snapshot` json NOT NULL, `unit_price` decimal(12,2) NOT NULL, `quantity` int unsigned NOT NULL, `refund_quantity` int unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_payment_transaction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_order_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay','balance','mock') NOT NULL,
  `transaction_no` varchar(128) NOT NULL, `amount` decimal(12,2) NOT NULL,
  `status` enum('created','processing','succeeded','failed','closed','refunded') NOT NULL DEFAULT 'created',
  `provider_transaction_no` varchar(128) DEFAULT NULL, `callback_idempotency_key` varchar(128) DEFAULT NULL, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_transaction` (`transaction_no`), UNIQUE KEY `uk_callback_key` (`callback_idempotency_key`), KEY `idx_group_status` (`group_order_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_delivery` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `delivery_type` enum('express','pickup','city','service') NOT NULL,
  `carrier_code` varchar(64) DEFAULT NULL, `tracking_no` varchar(128) DEFAULT NULL, `status` varchar(32) NOT NULL,
  `delivered_at` datetime DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_coupon_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `coupon_id` bigint unsigned NOT NULL,
  `source` varchar(32) NOT NULL, `status` enum('unused','locked','used','expired') NOT NULL DEFAULT 'unused',
  `obtained_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `used_order_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_coupon` (`user_id`,`coupon_id`), KEY `idx_user_status` (`user_id`,`status`), KEY `idx_coupon` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_coupon_user_index_exists := (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = 'qixi_crm_business' AND TABLE_NAME = 'qixi_crm_b_coupon_user' AND INDEX_NAME = 'uk_user_coupon');
SET @qixi_coupon_user_index_ddl := IF(@qixi_coupon_user_index_exists = 0, 'ALTER TABLE `qixi_crm_b_coupon_user` ADD UNIQUE INDEX `uk_user_coupon` (`user_id`,`coupon_id`)', 'SELECT 1');
PREPARE qixi_coupon_user_index_stmt FROM @qixi_coupon_user_index_ddl;
EXECUTE qixi_coupon_user_index_stmt;
DEALLOCATE PREPARE qixi_coupon_user_index_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_member_account` (
  `user_id` bigint unsigned NOT NULL, `level_id` bigint unsigned DEFAULT NULL, `points` bigint NOT NULL DEFAULT 0,
  `balance` decimal(12,2) NOT NULL DEFAULT 0, `commission` decimal(12,2) NOT NULL DEFAULT 0,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_asset_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `asset_type` enum('balance','points','commission') NOT NULL, `amount` decimal(16,2) NOT NULL,
  `reference_type` varchar(64) NOT NULL, `reference_id` varchar(64) NOT NULL, `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE KEY `uk_asset_idempotency` (`asset_type`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_community_post` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `store_id` bigint unsigned DEFAULT NULL,
  `content` text NOT NULL, `media` json NOT NULL, `status` enum('pending','approved','rejected','hidden') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_status_time` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_community_reply` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `post_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL,
  `content` varchar(1000) NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_post_time` (`post_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_binding` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `order_id` bigint unsigned DEFAULT NULL, `im_conversation_id` varchar(128) NOT NULL, `status` enum('open','closed') NOT NULL DEFAULT 'open',
  `assigned_admin_id` bigint unsigned DEFAULT NULL, `assigned_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_conversation` (`im_conversation_id`),
  KEY `idx_store_status_assignee` (`store_id`,`status`,`assigned_admin_id`), KEY `idx_assignee_status` (`assigned_admin_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_view` (
  `store_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_app_id` varchar(64) NOT NULL,
  `store_name` varchar(128) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 1,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`), UNIQUE KEY `uk_store_app_id` (`store_app_id`), UNIQUE KEY `uk_merchant` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_merchant_im_sdk_app_view` (
  `merchant_id` bigint unsigned NOT NULL, `sdk_app_id` varchar(64) NOT NULL,
  `api_public_url` varchar(1024) NOT NULL DEFAULT '', `ws_public_url` varchar(1024) NOT NULL DEFAULT '',
  `pte_profile_id` varchar(128) NOT NULL DEFAULT '', `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`merchant_id`), UNIQUE KEY `uk_sdk_app_id` (`sdk_app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 用户端只读取本库装修投影。装修原文归属平台/店铺库，严禁 C 端跨库直连。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_diy_page_view` (
  `source` enum('platform','merchant') NOT NULL,
  `page_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL DEFAULT 0,
  `page_type` enum('home','store_street','member','custom') NOT NULL DEFAULT 'home',
  `name` varchar(128) NOT NULL,
  `document` json NOT NULL,
  `status` enum('draft','published') NOT NULL DEFAULT 'draft',
  `is_active` tinyint NOT NULL DEFAULT 0,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`source`,`page_id`),
  KEY `idx_store_home_active` (`source`,`store_id`,`page_type`,`status`,`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `event_type` varchar(128) NOT NULL, `aggregate_type` varchar(64) NOT NULL,
  `aggregate_id` varchar(64) NOT NULL, `payload` json NOT NULL, `status` enum('pending','published','failed') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `published_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_status_time` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 用户、内容与商品消费域：所有用户端（PC、uni-app、原生 App）共用。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_profile` (
  `user_id` bigint unsigned NOT NULL, `avatar_url` varchar(1024) NOT NULL DEFAULT '', `gender` tinyint NOT NULL DEFAULT 0,
  `birthday` date DEFAULT NULL, `bio` varchar(500) NOT NULL DEFAULT '', `source_channel` enum('wechat','mini_program','h5','pc','ios','android','harmony') NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- C 端商户入驻申请。后台库只保存该记录的监管投影，C 端不得直写 qixi_crm_admin。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_merchant_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `applicant_user_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL,
  `contact_name` varchar(64) NOT NULL,
  `contact_mobile` varchar(32) NOT NULL,
  `category_name` varchar(128) NOT NULL DEFAULT '',
  `merchant_type` varchar(64) NOT NULL DEFAULT '',
  `license_url` varchar(1024) NOT NULL DEFAULT '',
  `status` enum('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user_created` (`applicant_user_id`,`created_at`), KEY `idx_status_created` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_browse_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL, `viewed_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user_viewed` (`user_id`,`viewed_at`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_follow_store` (
  `user_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`,`store_id`), KEY `idx_store` (`store_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_invoice_profile` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `type` enum('personal','enterprise') NOT NULL,
  `title` varchar(255) NOT NULL, `tax_no` varchar(64) NOT NULL DEFAULT '', `email` varchar(255) NOT NULL DEFAULT '',
  `is_default` tinyint NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_item_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `score` tinyint NOT NULL, `content` varchar(2000) NOT NULL DEFAULT '',
  `media` json NOT NULL, `reply_content` varchar(2000) NOT NULL DEFAULT '', `status` enum('pending','published','hidden') NOT NULL DEFAULT 'pending',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `replied_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_item` (`order_item_id`), KEY `idx_product_status` (`product_id`,`status`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_question` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL, `question` varchar(1000) NOT NULL, `answer` varchar(2000) NOT NULL DEFAULT '',
  `status` enum('pending','published','hidden') NOT NULL DEFAULT 'pending', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `answered_at` datetime DEFAULT NULL, PRIMARY KEY (`id`), KEY `idx_product_status` (`product_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_content_view` (
  `content_id` bigint unsigned NOT NULL, `content_type` enum('article','notice','agreement','banner') NOT NULL,
  `title` varchar(255) NOT NULL, `cover_url` varchar(1024) NOT NULL DEFAULT '', `body` longtext NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `version` bigint unsigned NOT NULL, `published_at` datetime DEFAULT NULL, `updated_at` datetime NOT NULL,
  PRIMARY KEY (`content_id`), KEY `idx_type_status_time` (`content_type`,`status`,`published_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 交易、支付、履约与售后域：每一次外部回调和状态迁移必须由幂等键保护。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_address_snapshot` (
  `order_id` bigint unsigned NOT NULL, `recipient` varchar(64) NOT NULL, `mobile` varchar(32) NOT NULL,
  `region_code` varchar(32) NOT NULL, `region_name` varchar(255) NOT NULL, `detail` varchar(255) NOT NULL,
  PRIMARY KEY (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_payment_callback` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `channel` enum('wechat','alipay','balance') NOT NULL, `provider_event_id` varchar(191) NOT NULL,
  `transaction_no` varchar(128) NOT NULL, `payload` json NOT NULL, `verified` tinyint NOT NULL DEFAULT 0,
  `processed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_channel_event` (`channel`,`provider_event_id`), KEY `idx_transaction` (`transaction_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_transaction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay','balance') NOT NULL,
  `provider_refund_no` varchar(128) NOT NULL, `amount` decimal(12,2) NOT NULL, `status` enum('created','processing','succeeded','failed') NOT NULL DEFAULT 'created',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `completed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_provider_refund` (`provider_refund_no`), UNIQUE KEY `uk_refund_idempotency` (`refund_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_invoice` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `invoice_profile_id` bigint unsigned NOT NULL,
  `status` enum('requested','issued','rejected','voided') NOT NULL DEFAULT 'requested', `invoice_no` varchar(128) NOT NULL DEFAULT '',
  `file_url` varchar(1024) NOT NULL DEFAULT '', `requested_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `issued_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_verification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `verify_code_hash` char(64) NOT NULL,
  `status` enum('unused','used','expired','cancelled') NOT NULL DEFAULT 'unused', `verified_by_account_id` bigint unsigned DEFAULT NULL,
  `expires_at` datetime DEFAULT NULL, `verified_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_verify_code_hash` (`verify_code_hash`), KEY `idx_order_status` (`order_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_aftersale_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL, `order_item_id` bigint unsigned NOT NULL,
  `quantity` int NOT NULL, `amount` decimal(12,2) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_item` (`refund_id`,`order_item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_aftersale_evidence` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL, `actor_type` enum('user','merchant','platform') NOT NULL,
  `actor_id` bigint unsigned NOT NULL, `content` varchar(2000) NOT NULL DEFAULT '', `attachments` json NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_refund_time` (`refund_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_event` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL,
  `from_status` varchar(32) NOT NULL DEFAULT '', `to_status` varchar(32) NOT NULL,
  `actor_type` enum('user','merchant','platform','system') NOT NULL, `actor_id` bigint unsigned NOT NULL,
  `reason` varchar(500) NOT NULL DEFAULT '', `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_transition_key` (`refund_id`,`idempotency_key`), KEY `idx_refund_created` (`refund_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 营销、会员、资产和分销域。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_coupon_template_view` (
  `coupon_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `discount_type` enum('amount','rate') NOT NULL, `discount_value` decimal(12,2) NOT NULL, `min_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1, `version` bigint unsigned NOT NULL,
  PRIMARY KEY (`coupon_id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_marketing_activity_view` (
  `activity_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `activity_type` enum('seckill','combination','presell','assist','discount','svip') NOT NULL, `name` varchar(128) NOT NULL,
  `rules` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `version` bigint unsigned NOT NULL,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, PRIMARY KEY (`activity_id`), KEY `idx_type_time` (`activity_type`,`status`,`starts_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_member_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL, `rank` int NOT NULL, `rules` json NOT NULL,
  `benefits` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, PRIMARY KEY (`id`), UNIQUE KEY `uk_rank` (`rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_recharge_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `recharge_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `amount` decimal(12,2) NOT NULL, `bonus_amount` decimal(12,2) NOT NULL DEFAULT 0, `status` enum('pending','paid','closed') NOT NULL DEFAULT 'pending',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_recharge_no` (`recharge_no`), UNIQUE KEY `uk_user_idempotency` (`user_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_withdrawal_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `withdrawal_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `amount` decimal(12,2) NOT NULL, `channel` enum('wechat','bank') NOT NULL, `account_snapshot` json NOT NULL,
  `status` enum('applied','reviewing','approved','paying','paid','rejected') NOT NULL DEFAULT 'applied', `review_note` varchar(500) NOT NULL DEFAULT '',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_withdrawal_no` (`withdrawal_no`), UNIQUE KEY `uk_user_withdrawal_key` (`user_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_relation` (
  `user_id` bigint unsigned NOT NULL, `parent_user_id` bigint unsigned DEFAULT NULL, `bound_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`), KEY `idx_parent` (`parent_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_commission_ledger` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `order_id` bigint unsigned NOT NULL,
  `amount` decimal(12,2) NOT NULL, `status` enum('pending','available','settled','voided') NOT NULL DEFAULT 'pending',
  `idempotency_key` varchar(128) NOT NULL, `available_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_commission_key` (`idempotency_key`), KEY `idx_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 直播、预约、客服和通知：消息与会话实体仍由 pte-live-im 管理。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_live_room` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `merchant_id` bigint unsigned NOT NULL DEFAULT 0, `store_id` bigint unsigned NOT NULL, `anchor_user_id` bigint unsigned DEFAULT NULL,
  `title` varchar(255) NOT NULL, `anchor_name` varchar(128) NOT NULL DEFAULT '', `cover_url` varchar(1024) NOT NULL DEFAULT '',
  `status` enum('draft','scheduled','living','ended','closed') NOT NULL DEFAULT 'draft', `is_public` tinyint NOT NULL DEFAULT 0,
  `stream_ref` varchar(191) NOT NULL DEFAULT '', `play_url` varchar(1024) NOT NULL DEFAULT '',
  `starts_at` datetime DEFAULT NULL, `ended_at` datetime DEFAULT NULL, `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_public_status_sort` (`is_public`,`status`,`sort`,`id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_live_room_product` (
  `live_room_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL, `sort` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`live_room_id`,`product_id`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_live_reservation` (
  `live_room_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`live_room_id`,`user_id`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_appointment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL,
  `service_product_id` bigint unsigned NOT NULL, `schedule_at` datetime NOT NULL, `status` enum('pending','confirmed','completed','cancelled') NOT NULL DEFAULT 'pending',
  `remark` varchar(1000) NOT NULL DEFAULT '', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_store_schedule` (`store_id`,`schedule_at`), KEY `idx_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_notification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `category` varchar(64) NOT NULL,
  `title` varchar(255) NOT NULL, `body` varchar(2000) NOT NULL, `payload` json NOT NULL, `read_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_user_read_time` (`user_id`,`read_at`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
