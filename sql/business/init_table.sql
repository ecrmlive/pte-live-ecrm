CREATE DATABASE IF NOT EXISTS `qixi_crm_business` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `nickname` varchar(64) NOT NULL DEFAULT '',
  `mobile` varchar(32) DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1, `group_id` bigint unsigned NOT NULL DEFAULT 0, `auth_version` bigint unsigned NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_mobile` (`mobile`), KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_user' AND COLUMN_NAME='group_id')=0,
    'ALTER TABLE `qixi_crm_b_user` ADD COLUMN `group_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `status`, ADD KEY `idx_group_id` (`group_id`)',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
-- 用户标签与分组属于消费者运营事实，后台通过 RBAC 管理，数据统一存放业务库。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_label` (
  `label_id` bigint unsigned NOT NULL AUTO_INCREMENT, `label_name` varchar(64) NOT NULL,
  `sort` int NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`label_id`), KEY `idx_listing` (`is_del`,`sort`,`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_group` (
  `group_id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_name` varchar(64) NOT NULL,
  `sort` int NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_id`), KEY `idx_listing` (`is_del`,`sort`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_label_relation` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `uid` bigint unsigned NOT NULL,
  `label_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_label` (`uid`,`label_id`), KEY `idx_label` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台批量调整用户分组的命令审计。一次命令只调整运营归属，不触碰会员、订单、资产或佣金事实。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_group_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_ids_json` json NOT NULL, `group_id` bigint unsigned NOT NULL DEFAULT 0,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_idempotency_key` (`idempotency_key`), KEY `idx_group_created` (`group_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台批量替换用户运营标签的命令审计；空标签集合表示清除运营标签。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_label_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_ids_json` json NOT NULL, `label_ids_json` json NOT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台启停用户的不可变审计。状态切换递增 auth_version，使既有 C 端令牌立即失效。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_status_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `from_status` tinyint NOT NULL, `to_status` tinyint NOT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_idempotency_key` (`idempotency_key`), KEY `idx_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台创建用户、维护资料、重置密码的命令审计。只保存不可逆请求摘要，绝不保存密码或身份凭据。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_admin_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `action` enum('create','profile_update','password_reset') NOT NULL,
  `user_id` bigint unsigned NOT NULL, `request_fingerprint` char(64) NOT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_action_idempotency` (`action`,`idempotency_key`), KEY `idx_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 用户信息导出审计仅保存筛选摘要、行数、原因与操作人，不保存导出的个人资料副本。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_export_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `query_fingerprint` char(64) NOT NULL,
  `row_count` int unsigned NOT NULL, `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_operator_created` (`operator_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_identity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `channel` enum('wechat','mini_program','h5','pc','ios','android','harmony') NOT NULL, `subject` varchar(191) NOT NULL,
  `credential_hash` varchar(255) DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_channel_subject` (`channel`,`subject`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- pte-tools-captcha 校验令牌仅保存 SHA-256 摘要；令牌只可按用途消费一次。
-- 短信验证码只保存摘要；同手机号和用途的明文验证码绝不写入数据库、日志或响应。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_auth_sms_code` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `mobile` varchar(32) NOT NULL, `purpose` enum('login','binding','reset_password','change_mobile') NOT NULL,
  `code_hash` char(64) NOT NULL, `expires_at` datetime NOT NULL, `consumed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_sms_available` (`mobile`,`purpose`,`consumed_at`,`expires_at`,`id`), KEY `idx_sms_created` (`mobile`,`purpose`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
-- 注销确认令牌只保留 SHA-256 摘要且十分钟过期；账户注销不会删除交易、售后和资金事实。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_cancellation_confirmation` (
  `user_id` bigint unsigned NOT NULL, `token_hash` char(64) NOT NULL,
  `blockers` json NOT NULL, `expires_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`), KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_cancellation_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `confirmation_required` tinyint NOT NULL, `blockers` json NOT NULL,
  `cancelled_at` datetime NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_cancellation` (`user_id`), KEY `idx_cancelled_at` (`cancelled_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_address` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `recipient` varchar(64) NOT NULL,
  `mobile` varchar(32) NOT NULL, `province` varchar(64) NOT NULL DEFAULT '', `city` varchar(64) NOT NULL DEFAULT '',
  `district` varchar(64) NOT NULL DEFAULT '', `region_code` varchar(32) NOT NULL DEFAULT '', `detail` varchar(255) NOT NULL,
  `post_code` int unsigned NOT NULL DEFAULT 0, `is_default` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 行政区划是由平台配送配置发布到业务库的只读投影。C 端只读取可见节点，
-- 不直接跨库读取 qixi_crm_a_city，也不允许通过用户端写入。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_city_view` (
  `city_id` bigint unsigned NOT NULL, `parent_id` bigint unsigned NOT NULL DEFAULT 0,
  `name` varchar(128) NOT NULL, `level` tinyint NOT NULL DEFAULT 0, `is_show` tinyint NOT NULL DEFAULT 1,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`city_id`), KEY `idx_parent_visible` (`parent_id`,`is_show`,`city_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_view` (
  `product_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '', `store_name` varchar(128) NOT NULL DEFAULT '',
  `category_id` bigint unsigned NOT NULL DEFAULT 0, `brand_name` varchar(64) NOT NULL DEFAULT '', `title` varchar(255) NOT NULL, `cover_url` varchar(1024) NOT NULL DEFAULT '',
  `price` decimal(12,2) NOT NULL, `original_price` decimal(12,2) DEFAULT NULL, `svip_price_type` tinyint NOT NULL DEFAULT 0, `svip_price` decimal(12,2) NOT NULL DEFAULT 0, `product_type` tinyint NOT NULL DEFAULT 0,
  `sales` int NOT NULL DEFAULT 0, `stock` int NOT NULL, `sale_status` tinyint NOT NULL, `version` bigint unsigned NOT NULL,
  `updated_at` datetime NOT NULL, PRIMARY KEY (`product_id`), KEY `idx_merchant_sale` (`merchant_id`,`sale_status`), KEY `idx_store_sale` (`store_id`,`sale_status`), KEY `idx_category_sale` (`category_id`,`sale_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- CREATE IF NOT EXISTS 不会给旧表补列；幂等补齐品牌与 SVIP 价字段。
SET @qixi_ddl := (
  SELECT IF(
    (
      SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'qixi_crm_b_product_view'
        AND COLUMN_NAME = 'brand_name'
    ) = 0,
    'ALTER TABLE `qixi_crm_b_product_view` ADD COLUMN `brand_name` varchar(64) NOT NULL DEFAULT '''' AFTER `category_id`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (
      SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'qixi_crm_b_product_view'
        AND COLUMN_NAME = 'svip_price_type'
    ) = 0,
    'ALTER TABLE `qixi_crm_b_product_view` ADD COLUMN `svip_price_type` tinyint NOT NULL DEFAULT 0 AFTER `original_price`, ADD COLUMN `svip_price` decimal(12,2) NOT NULL DEFAULT 0 AFTER `svip_price_type`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 商品消费视图只保存由商户域发布的可售 SKU 映射。C 端购物车和订单明细
-- 以 merchant_sku_id 追溯库存事实，不能由商品 ID 或展示规格文字猜测库存行。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_sku_view` (
  `merchant_sku_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `sku_key` varchar(128) NOT NULL, `spec_snapshot` json NOT NULL,
  `price` decimal(12,2) NOT NULL, `stock` int NOT NULL, `sale_status` tinyint NOT NULL,
  `version` bigint unsigned NOT NULL DEFAULT 1, `updated_at` datetime NOT NULL,
  PRIMARY KEY (`merchant_sku_id`), UNIQUE KEY `uk_product_sku_key` (`product_id`,`sku_key`),
  KEY `idx_product_sale` (`product_id`,`sale_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 直播间与挂货属于业务事实。统一后台仅监管审核、显示状态和业务投影，
-- 不返回推流地址或主播手机号等敏感字段。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_broadcast_room` (
  `broadcast_room_id` bigint unsigned NOT NULL AUTO_INCREMENT, `mer_id` bigint unsigned NOT NULL,
  `name` varchar(255) NOT NULL, `cover_img` varchar(1024) NOT NULL DEFAULT '', `feeds_img` varchar(1024) NOT NULL DEFAULT '',
  `play_url` varchar(2048) NOT NULL DEFAULT '', `push_url` varchar(2048) NOT NULL DEFAULT '',
  `start_time` datetime DEFAULT NULL, `end_time` datetime DEFAULT NULL, `anchor_name` varchar(128) NOT NULL DEFAULT '',
  `phone` varchar(32) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 0, `live_status` smallint NOT NULL DEFAULT 102,
  `is_show` tinyint NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sort` int NOT NULL DEFAULT 0, `star` int NOT NULL DEFAULT 0, `mark` varchar(500) NOT NULL DEFAULT '', `refusal` varchar(500) NOT NULL DEFAULT '',
  PRIMARY KEY (`broadcast_room_id`), KEY `idx_merchant` (`mer_id`,`is_del`,`create_time`), KEY `idx_review_visibility` (`status`,`is_show`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_broadcast_room_goods` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `broadcast_room_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `on_sale` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_room_product` (`broadcast_room_id`,`product_id`), KEY `idx_room_sort` (`broadcast_room_id`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 预约服务：活动配置、可售时段和日占位账本分离。日账本是余量的唯一来源，
-- 不再从历史 qixi_m_app_store_order 反查，避免跨表前缀和并发超卖。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_points_product_view` (
  `product_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '', `store_name` varchar(128) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL, `cover_url` varchar(1024) NOT NULL DEFAULT '', `original_price` decimal(12,2) NOT NULL DEFAULT 0,
  `points_required` bigint NOT NULL, `stock` int NOT NULL, `sale_status` tinyint NOT NULL DEFAULT 1,
  `version` bigint unsigned NOT NULL DEFAULT 1, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`), KEY `idx_store_sale` (`store_id`,`sale_status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_reservation_activity` (
  `product_reservation_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL DEFAULT 0,
  `reservation_type` tinyint NOT NULL DEFAULT 1, `show_reservation_days` int NOT NULL DEFAULT 7,
  `is_cancel_reservation` tinyint NOT NULL DEFAULT 1, `time_period` text NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_reservation_id`), UNIQUE KEY `uk_product` (`product_id`), KEY `idx_merchant_status` (`merchant_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_reservation_slot` (
  `attr_reservation_id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL,
  `slot_key` varchar(64) NOT NULL, `start_time` char(5) NOT NULL, `end_time` char(5) NOT NULL,
  `stock` int NOT NULL, `use_num` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`attr_reservation_id`), UNIQUE KEY `uk_product_slot_key` (`product_id`,`slot_key`),
  KEY `idx_product_time` (`product_id`,`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_reservation_booking` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_id` bigint unsigned NOT NULL, `slot_id` bigint unsigned NOT NULL,
  `booking_date` date NOT NULL, `order_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `verify_code` varchar(32) NOT NULL DEFAULT '', `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order` (`order_id`), KEY `idx_slot_day_status` (`slot_id`,`booking_date`,`status`),
  KEY `idx_user_day` (`user_id`,`booking_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 秒杀活动与时段独立保存于业务库。平台监管、商户配置和 C 端报价均只读写这一份状态，
-- 禁止回退到历史 qixi_m_admin_store_seckill_* 表。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_seckill_time` (
  `seckill_time_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(64) NOT NULL, `start_time` int NOT NULL, `end_time` int NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `pic` varchar(1024) NOT NULL DEFAULT '',
  PRIMARY KEY (`seckill_time_id`), UNIQUE KEY `uk_start_end` (`start_time`,`end_time`), KEY `idx_status_start` (`status`,`start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_seckill_active` (
  `seckill_active_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL, `seckill_time_ids` varchar(255) NOT NULL DEFAULT '',
  `start_day` date NOT NULL, `end_day` date NOT NULL, `mer_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL, `seckill_price` decimal(12,2) NOT NULL,
  `once_pay_count` int NOT NULL DEFAULT 1, `all_pay_count` int NOT NULL DEFAULT 0,
  `active_status` tinyint NOT NULL DEFAULT 1, `status` tinyint NOT NULL DEFAULT 1,
  `create_time` bigint NOT NULL, `update_time` bigint NOT NULL, `delete_time` bigint DEFAULT NULL,
  PRIMARY KEY (`seckill_active_id`), KEY `idx_listing` (`delete_time`,`status`,`active_status`,`start_day`,`end_day`),
  KEY `idx_merchant` (`mer_id`,`delete_time`), KEY `idx_product` (`product_id`,`delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 拼团状态机：活动定义、开团记录、团员记录三表。订单支付前仅占位，支付成功后由业务服务推进成员和成团状态。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_combination_group` (
  `product_group_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` bigint unsigned NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL,
  `time` int NOT NULL DEFAULT 24, `buying_count_num` int NOT NULL, `buying_num` int NOT NULL DEFAULT 1,
  `pay_count` int NOT NULL DEFAULT 0, `once_pay_count` int NOT NULL DEFAULT 1,
  `status` tinyint NOT NULL DEFAULT 1, `mer_id` bigint unsigned NOT NULL,
  `is_show` tinyint NOT NULL DEFAULT 1, `is_del` tinyint NOT NULL DEFAULT 0,
  `success_num` int NOT NULL DEFAULT 0, `product_status` tinyint NOT NULL DEFAULT 1,
  `price` decimal(12,2) NOT NULL, `action_status` tinyint NOT NULL DEFAULT 1,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_group_id`), KEY `idx_app_listing` (`is_del`,`is_show`,`status`,`action_status`,`start_time`,`end_time`),
  KEY `idx_merchant` (`mer_id`,`is_del`,`create_time`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_combination_buying` (
  `group_buying_id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_group_id` bigint unsigned NOT NULL,
  `status` tinyint NOT NULL DEFAULT 0, `buying_count_num` int NOT NULL, `buying_num` int NOT NULL DEFAULT 1,
  `yet_buying_num` int NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0, `mer_id` bigint unsigned NOT NULL,
  `end_time` bigint NOT NULL, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_buying_id`), KEY `idx_open` (`product_group_id`,`status`,`is_del`,`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_combination_member` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_buying_id` bigint unsigned NOT NULL,
  `product_group_id` bigint unsigned NOT NULL, `status` tinyint NOT NULL DEFAULT 0,
  `is_initiator` tinyint NOT NULL DEFAULT 0, `order_id` bigint unsigned NOT NULL DEFAULT 0,
  `uid` bigint unsigned NOT NULL, `nickname` varchar(64) NOT NULL DEFAULT '', `avatar` varchar(1024) NOT NULL DEFAULT '',
  `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_leader` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), KEY `idx_buying` (`group_buying_id`,`is_del`,`id`), KEY `idx_order` (`order_id`,`is_del`),
  KEY `idx_user` (`uid`,`product_group_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 预售活动及定金尾款记录。活动库存由服务层在事务中扣减，尾款记录独立保存。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_presell` (
  `product_presell_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `start_time` datetime NOT NULL, `end_time` datetime NOT NULL,
  `final_start_time` varchar(32) NOT NULL DEFAULT '', `final_end_time` varchar(32) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1, `presell_type` tinyint NOT NULL DEFAULT 1,
  `pay_count` int NOT NULL DEFAULT 0, `delivery_type` tinyint NOT NULL DEFAULT 1, `delivery_day` int NOT NULL DEFAULT 0,
  `product_id` bigint unsigned NOT NULL, `price` decimal(12,2) NOT NULL,
  `down_price` decimal(12,2) NOT NULL DEFAULT 0, `final_price` decimal(12,2) NOT NULL DEFAULT 0,
  `stock` int NOT NULL DEFAULT 0, `is_show` tinyint NOT NULL DEFAULT 1,
  `store_name` varchar(255) NOT NULL DEFAULT '', `mer_id` bigint unsigned NOT NULL, `store_info` text NOT NULL,
  `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `product_status` tinyint NOT NULL DEFAULT 1, `refusal` varchar(500) NOT NULL DEFAULT '',
  `action_status` tinyint NOT NULL DEFAULT 1, `seles` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`product_presell_id`), KEY `idx_app_listing` (`is_del`,`status`,`is_show`,`product_status`,`action_status`,`start_time`,`end_time`),
  KEY `idx_merchant` (`mer_id`,`is_del`,`create_time`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_presell_order` (
  `presell_order_id` bigint unsigned NOT NULL AUTO_INCREMENT, `presell_order_sn` varchar(64) NOT NULL,
  `uid` bigint unsigned NOT NULL, `mer_id` bigint unsigned NOT NULL, `order_id` bigint unsigned NOT NULL,
  `product_presell_id` bigint unsigned NOT NULL, `final_start_time` datetime NOT NULL, `final_end_time` datetime NOT NULL,
  `paid` tinyint NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 1, `pay_type` tinyint NOT NULL DEFAULT 0,
  `pay_price` decimal(12,2) NOT NULL, `pay_time` datetime DEFAULT NULL, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`presell_order_id`), UNIQUE KEY `uk_order` (`order_id`), KEY `idx_user_status` (`uid`,`paid`,`status`,`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 助力活动、助力发起人和助力人记录。状态 1/10/11/20/-1 的推进由服务层定义。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_assist` (
  `product_assist_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `pay_count` int NOT NULL DEFAULT 0, `assist_count` int NOT NULL, `assist_user_count` int NOT NULL,
  `product_id` bigint unsigned NOT NULL, `assist_price` decimal(12,2) NOT NULL, `stock` int NOT NULL DEFAULT 0,
  `is_show` tinyint NOT NULL DEFAULT 1, `store_name` varchar(255) NOT NULL DEFAULT '', `mer_id` bigint unsigned NOT NULL,
  `store_info` text NOT NULL, `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `product_status` tinyint NOT NULL DEFAULT 1, `refusal` varchar(500) NOT NULL DEFAULT '', `action_status` tinyint NOT NULL DEFAULT 1,
  PRIMARY KEY (`product_assist_id`), KEY `idx_app_listing` (`is_del`,`status`,`is_show`,`product_status`,`action_status`,`start_time`,`end_time`),
  KEY `idx_merchant` (`mer_id`,`is_del`,`create_time`), KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_assist_set` (
  `product_assist_set_id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_assist_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL, `uid` bigint unsigned NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `assist_count` int NOT NULL, `assist_user_count` int NOT NULL, `yet_assist_count` int NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `mer_id` bigint unsigned NOT NULL, `is_del` tinyint NOT NULL DEFAULT 0,
  PRIMARY KEY (`product_assist_set_id`), KEY `idx_open` (`product_assist_id`,`status`,`is_del`,`product_assist_set_id`), KEY `idx_user` (`uid`,`product_assist_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_assist_user` (
  `product_assist_user_id` bigint unsigned NOT NULL AUTO_INCREMENT, `product_assist_set_id` bigint unsigned NOT NULL,
  `product_assist_id` bigint unsigned NOT NULL, `uid` bigint unsigned NOT NULL, `nickname` varchar(64) NOT NULL DEFAULT '',
  `avatar_img` varchar(1024) NOT NULL DEFAULT '', `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_assist_user_id`), UNIQUE KEY `uk_set_user` (`product_assist_set_id`,`uid`), KEY `idx_assist_user` (`product_assist_id`,`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 社区旧模型的兼容业务表。与新的精简社区投影分开，避免字段语义和审核状态混用。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_category` (
  `category_id` bigint unsigned NOT NULL AUTO_INCREMENT, `cate_name` varchar(128) NOT NULL, `pid` int NOT NULL DEFAULT 0,
  `is_show` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, PRIMARY KEY (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_topic` (
  `topic_id` bigint unsigned NOT NULL AUTO_INCREMENT, `topic_name` varchar(128) NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `is_hot` tinyint NOT NULL DEFAULT 0, `category_id` bigint unsigned NOT NULL DEFAULT 0, `is_del` tinyint NOT NULL DEFAULT 0,
  `count_use` int NOT NULL DEFAULT 0, `sort` int NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`topic_id`), KEY `idx_public` (`is_del`,`status`,`is_hot`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_post` (
  `community_id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(255) NOT NULL, `image` varchar(1024) NOT NULL DEFAULT '',
  `category_id` bigint unsigned NOT NULL DEFAULT 0, `topic_id` bigint unsigned NOT NULL DEFAULT 0, `uid` bigint unsigned NOT NULL,
  `mer_id` bigint unsigned NOT NULL DEFAULT 0, `product_id` bigint unsigned NOT NULL DEFAULT 0,
  `count_start` int NOT NULL DEFAULT 0, `count_reply` int NOT NULL DEFAULT 0, `status` tinyint NOT NULL DEFAULT 0,
  `is_show` tinyint NOT NULL DEFAULT 1, `is_hot` tinyint NOT NULL DEFAULT 0, `is_type` tinyint NOT NULL DEFAULT 0,
  `content` text NOT NULL, `refusal` varchar(500) NOT NULL DEFAULT '', `pv` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `status_time` datetime DEFAULT NULL,
  PRIMARY KEY (`community_id`), KEY `idx_public` (`is_del`,`status`,`is_show`,`is_hot`,`community_id`), KEY `idx_topic` (`topic_id`,`is_del`,`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_social_reply` (
  `reply_id` bigint unsigned NOT NULL AUTO_INCREMENT, `content` varchar(1000) NOT NULL, `pid` bigint unsigned NOT NULL DEFAULT 0,
  `uid` bigint unsigned NOT NULL, `community_id` bigint unsigned NOT NULL, `status` tinyint NOT NULL DEFAULT 1,
  `is_del` tinyint NOT NULL DEFAULT 0, `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`reply_id`), KEY `idx_post` (`community_id`,`is_del`,`status`,`reply_id`)
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
  `total_quantity` int unsigned NOT NULL, `activity_type` tinyint NOT NULL DEFAULT 0, `points_amount` bigint NOT NULL DEFAULT 0, `recipient_snapshot` json NOT NULL,
  `pay_channel` enum('balance','wechat','alipay','mock') DEFAULT NULL,
  `pay_status` enum('pending','paid','closed','refunded') NOT NULL DEFAULT 'pending', `paid_at` datetime DEFAULT NULL,
  `idempotency_key` varchar(128) NOT NULL, `remark` varchar(500) NOT NULL DEFAULT '', `user_archived_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_no` (`order_no`), UNIQUE KEY `uk_user_idempotency` (`user_id`,`idempotency_key`), KEY `idx_user_status_time` (`user_id`,`pay_status`,`created_at`), KEY `idx_user_archive_time` (`user_id`,`user_archived_at`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_order_id` bigint unsigned NOT NULL, `order_no` varchar(64) NOT NULL,
  `merchant_id` bigint unsigned NOT NULL, `merchant_name_snapshot` varchar(128) NOT NULL, `store_id` bigint unsigned NOT NULL,
  `store_name_snapshot` varchar(128) NOT NULL, `user_id` bigint unsigned NOT NULL, `total_amount` decimal(12,2) NOT NULL,
  `discount_amount` decimal(12,2) NOT NULL DEFAULT 0, `freight_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `pay_amount` decimal(12,2) NOT NULL, `total_quantity` int unsigned NOT NULL, `activity_type` tinyint NOT NULL DEFAULT 0, `points_amount` bigint NOT NULL DEFAULT 0,
  `recipient_snapshot` json NOT NULL, `remark` varchar(500) NOT NULL DEFAULT '',
  `status` enum('pending_pay','paid','awaiting_final','final_timeout','fulfilling','shipped','completed','cancelled','aftersale') NOT NULL DEFAULT 'pending_pay',
  `paid_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_no` (`order_no`), KEY `idx_merchant_status` (`merchant_id`,`status`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `refund_no` varchar(64) NOT NULL,
  `reason` varchar(500) NOT NULL, `amount` decimal(12,2) NOT NULL,
  `refund_type` enum('money_only','return_and_refund') NOT NULL DEFAULT 'money_only',
  `order_status_before` enum('pending_pay','paid','fulfilling','shipped','completed','cancelled','aftersale') NOT NULL,
  `status` enum('applied','merchant_handling','awaiting_return','awaiting_receipt','platform_intervene','refunding','refunded','rejected','cancelled') NOT NULL DEFAULT 'applied',
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
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), UNIQUE KEY `uk_user_keyword` (`user_id`,`keyword`), KEY `idx_user_time` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `product_id` bigint unsigned NOT NULL,
  `merchant_sku_id` bigint unsigned NOT NULL DEFAULT 0,
  `sku_key` varchar(128) NOT NULL, `title_snapshot` varchar(255) NOT NULL, `cover_url_snapshot` varchar(1024) NOT NULL DEFAULT '',
  `spec_snapshot` json NOT NULL, `unit_price` decimal(12,2) NOT NULL, `quantity` int unsigned NOT NULL, `refund_quantity` int unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`), KEY `idx_order` (`order_id`), KEY `idx_merchant_sku` (`merchant_sku_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 历史订单以 0 明确标记为“未映射 SKU”；新建订单由订单服务拒绝该状态，
-- 业务库库存命令 outbox 与订单创建在同一事务提交。状态 accepted 仅表示商户库
-- 已完成幂等预留；网络失败保持 pending，由投递器重试，禁止在订单事务中越库扣减。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_stock_command_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `action` enum('reserve','confirm','release','restock') NOT NULL,
  `order_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `merchant_sku_id` bigint unsigned NOT NULL,
  `quantity` int unsigned NOT NULL, `expires_at` datetime DEFAULT NULL,
  `idempotency_key` varchar(128) NOT NULL, `status` enum('pending','accepted','failed') NOT NULL DEFAULT 'pending',
  `attempts` int unsigned NOT NULL DEFAULT 0, `last_error` varchar(500) NOT NULL DEFAULT '', `processed_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_action_order_sku` (`action`,`order_id`,`merchant_sku_id`),
  UNIQUE KEY `uk_idempotency_key` (`idempotency_key`), KEY `idx_status_created` (`status`,`created_at`), KEY `idx_order_action` (`order_id`,`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 商户结算命令 outbox。确认收货和可信退款在业务订单事务内写入，禁止跨库改结算。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_settlement_command_outbox` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `action` enum('accrue','reverse') NOT NULL,
  `order_id` bigint unsigned NOT NULL, `refund_id` bigint unsigned NOT NULL DEFAULT 0,
  `store_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `amount` decimal(16,2) NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `status` enum('pending','accepted','failed') NOT NULL DEFAULT 'pending',
  `attempts` int unsigned NOT NULL DEFAULT 0, `last_error` varchar(500) NOT NULL DEFAULT '', `processed_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_action_order_refund` (`action`,`order_id`,`refund_id`),
  UNIQUE KEY `uk_idempotency_key` (`idempotency_key`), KEY `idx_status_created` (`status`,`created_at`), KEY `idx_order_action` (`order_id`,`action`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 特殊营销订单与活动实例的稳定关联。支付、取消和超时关闭都必须通过该表
-- 推进/回滚活动状态，禁止从备注字段或旧 qixi_m_* 订单反查。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_activity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_order_id` bigint unsigned NOT NULL,
  `activity_type` tinyint NOT NULL, `activity_id` bigint unsigned NOT NULL,
  `related_activity_id` bigint unsigned NOT NULL DEFAULT 0, `quantity` int unsigned NOT NULL DEFAULT 1,
  `status` enum('reserved','paid','released') NOT NULL DEFAULT 'reserved',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_group_activity` (`group_order_id`,`activity_type`),
  KEY `idx_activity_instance` (`activity_type`,`activity_id`), KEY `idx_activity_status` (`activity_type`,`related_activity_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_payment_transaction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `group_order_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay','balance','mock') NOT NULL,
  `transaction_no` varchar(128) NOT NULL, `amount` decimal(12,2) NOT NULL,
  `status` enum('created','processing','succeeded','failed','closed','refunded') NOT NULL DEFAULT 'created',
  `provider_transaction_no` varchar(128) DEFAULT NULL, `provider_payload` json DEFAULT NULL, `callback_idempotency_key` varchar(128) DEFAULT NULL, `paid_at` datetime DEFAULT NULL,
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
-- 平台人工发券/撤销审计。撤销只改变未锁定券为 expired，绝不删除已被订单引用的用户券。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_coupon_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `coupon_id` bigint unsigned NOT NULL, `coupon_user_id` bigint unsigned NOT NULL,
  `action` enum('issue','revoke') NOT NULL, `from_status` varchar(16) NOT NULL DEFAULT '',
  `to_status` varchar(16) NOT NULL, `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL, `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_coupon_command_idempotency` (`action`,`idempotency_key`),
  KEY `idx_user_coupon_command` (`user_id`,`coupon_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
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
-- 平台人工调账的不可变审计事实。账户余额与台账必须在同一业务库事务内更新，
-- 幂等键由调用方生成，禁止由浏览器以外的支付或结算回调复用该入口。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_asset_adjustment_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `asset_type` enum('balance','points') NOT NULL, `amount` decimal(16,2) NOT NULL,
  `balance_before` decimal(16,2) NOT NULL, `balance_after` decimal(16,2) NOT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_asset_adjustment_idempotency` (`asset_type`,`idempotency_key`),
  KEY `idx_user_asset_adjustment` (`user_id`,`created_at`)
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
  `last_msg` varchar(500) NOT NULL DEFAULT '', `last_time` datetime DEFAULT NULL,
  `user_unread` int unsigned NOT NULL DEFAULT 0, `service_unread` int unsigned NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_conversation` (`im_conversation_id`), UNIQUE KEY `uk_user_store` (`user_id`,`store_id`),
  KEY `idx_store_status_assignee` (`store_id`,`status`,`assigned_admin_id`), KEY `idx_assignee_status` (`assigned_admin_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 会话归属转接是业务审计事实，不能覆写到绑定表的备注或 IM 消息正文中。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_assignment_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `binding_id` bigint unsigned NOT NULL,
  `from_admin_id` bigint unsigned DEFAULT NULL, `target_admin_id` bigint unsigned NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL, `reason` varchar(500) NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_binding_transfer_key` (`binding_id`,`idempotency_key`),
  KEY `idx_binding_created` (`binding_id`,`created_at`), KEY `idx_target_created` (`target_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `binding_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL,
  `sender_role` enum('user','service','system') NOT NULL, `sender_id` bigint unsigned NOT NULL DEFAULT 0,
  `msg_type` enum('order','system') NOT NULL, `content` text NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_binding_created` (`binding_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 客服快捷回复属于业务数据；统一后台按客服授权店铺维护，删除采用软删除以保留审计线索。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_quick_reply` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `store_id` bigint unsigned NOT NULL,
  `title` varchar(64) NOT NULL, `content` varchar(2000) NOT NULL,
  `status` enum('enabled','disabled') NOT NULL DEFAULT 'enabled',
  `created_by` bigint unsigned NOT NULL, `updated_by` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_store_status_updated` (`store_id`,`status`,`updated_at`), KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 客服备注以用户和店铺联合隔离，必须通过已授权会话写入，不能以任意用户 ID 直写。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_user_note` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `content` varchar(500) NOT NULL, `updated_by` bigint unsigned NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_store` (`user_id`,`store_id`), KEY `idx_store_updated` (`store_id`,`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_im_identity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `portal` enum('app','service') NOT NULL, `local_id` bigint unsigned NOT NULL,
  `im_user_id` varchar(64) NOT NULL, `im_user_num` bigint NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_portal_local` (`portal`,`local_id`), UNIQUE KEY `uk_im_user_id` (`im_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 坐席由商户后台投影维护；C 端只读取可用坐席，不跨库读取后台账号表。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_customer_service_agent_view` (
  `admin_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `display_name` varchar(64) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1, `available_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`admin_id`,`store_id`), KEY `idx_store_available` (`store_id`,`status`,`available_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 店铺消费投影；积分抵扣规则由商户配置事件同步，默认关闭，C 端不得跨库读取商户库。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_view` (
  `store_id` bigint unsigned NOT NULL, `merchant_id` bigint unsigned NOT NULL, `store_app_id` varchar(64) NOT NULL,
  `store_name` varchar(128) NOT NULL DEFAULT '', `status` tinyint NOT NULL DEFAULT 1,
  `integral_enabled` tinyint NOT NULL DEFAULT 0, `integral_points_per_yuan` bigint NOT NULL DEFAULT 100, `integral_max_deduction_bps` int NOT NULL DEFAULT 2000,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_id`), UNIQUE KEY `uk_store_app_id` (`store_app_id`), UNIQUE KEY `uk_merchant` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- CREATE IF NOT EXISTS 不会给旧表补列；幂等补齐积分抵扣消费投影字段。
SET @qixi_ddl := (
  SELECT IF(
    (
      SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE()
        AND TABLE_NAME = 'qixi_crm_b_store_view'
        AND COLUMN_NAME = 'integral_enabled'
    ) = 0,
    'ALTER TABLE `qixi_crm_b_store_view` ADD COLUMN `integral_enabled` tinyint NOT NULL DEFAULT 0 AFTER `status`, ADD COLUMN `integral_points_per_yuan` bigint NOT NULL DEFAULT 100 AFTER `integral_enabled`, ADD COLUMN `integral_max_deduction_bps` int NOT NULL DEFAULT 2000 AFTER `integral_points_per_yuan`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl;
EXECUTE qixi_stmt;
DEALLOCATE PREPARE qixi_stmt;
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
  `license_key` varchar(1024) NOT NULL DEFAULT '',
  `license_url` varchar(1024) NOT NULL DEFAULT '',
  `status` enum('pending','approved','rejected') NOT NULL DEFAULT 'pending',
  `review_note` varchar(500) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user_created` (`applicant_user_id`,`created_at`), KEY `idx_status_created` (`status`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_upload_object` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `owner_user_id` bigint unsigned NOT NULL,
  `purpose` varchar(64) NOT NULL,
  `object_key` varchar(512) NOT NULL,
  `original_name` varchar(255) NOT NULL,
  `content_type` varchar(128) NOT NULL,
  `size` bigint unsigned NOT NULL,
  `status` enum('issued','completed','expired') NOT NULL DEFAULT 'issued',
  `expires_at` datetime NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `completed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_object_key` (`object_key`), KEY `idx_owner_purpose_status` (`owner_user_id`,`purpose`,`status`)
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
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_comment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_item_id` bigint unsigned DEFAULT NULL, `user_id` bigint unsigned NOT NULL,
  `product_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `score` tinyint NOT NULL, `content` varchar(2000) NOT NULL DEFAULT '',
  `media` json NOT NULL, `reply_content` varchar(2000) NOT NULL DEFAULT '', `status` enum('pending','published','hidden') NOT NULL DEFAULT 'pending',
  `source` enum('user','virtual') NOT NULL DEFAULT 'user', `virtual_author_name` varchar(64) NOT NULL DEFAULT '', `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `replied_at` datetime DEFAULT NULL, `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order_item` (`order_item_id`), KEY `idx_product_status` (`product_id`,`status`), KEY `idx_store_status` (`store_id`,`status`), KEY `idx_product_visible_sort` (`product_id`,`deleted_at`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_product_comment' AND COLUMN_NAME='source')=0,
    'ALTER TABLE `qixi_crm_b_product_comment` ADD COLUMN `source` enum(''user'',''virtual'') NOT NULL DEFAULT ''user'' AFTER `status`, ADD COLUMN `virtual_author_name` varchar(64) NOT NULL DEFAULT '''' AFTER `source`, ADD COLUMN `sort` int NOT NULL DEFAULT 0 AFTER `virtual_author_name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_product_comment' AND COLUMN_NAME='deleted_at')=0,
    'ALTER TABLE `qixi_crm_b_product_comment` ADD COLUMN `deleted_at` datetime DEFAULT NULL AFTER `replied_at`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
-- 旧表曾把 order_item_id 建成 NOT NULL；虚拟评论需要可空。
SET @qixi_ddl := (
  SELECT IF(
    (
      SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_product_comment'
        AND COLUMN_NAME='order_item_id' AND IS_NULLABLE='NO'
    )>0,
    'ALTER TABLE `qixi_crm_b_product_comment` MODIFY COLUMN `order_item_id` bigint unsigned DEFAULT NULL',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_product_comment_moderation_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `comment_id` bigint unsigned NOT NULL,
  `from_status` varchar(16) NOT NULL, `to_status` varchar(16) NOT NULL, `action` varchar(16) NOT NULL,
  `note` varchar(500) NOT NULL DEFAULT '', `operator_admin_id` bigint unsigned NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_idempotency` (`idempotency_key`), KEY `idx_comment_time` (`comment_id`,`created_at`)
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
-- 退款回调与支付回调分表：避免事件类型、外部退款单号和原交易号混用。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_callback` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `channel` enum('wechat','alipay','balance','mock') NOT NULL,
  `provider_event_id` varchar(191) NOT NULL, `provider_refund_no` varchar(128) NOT NULL,
  `out_trade_no` varchar(128) NOT NULL, `payload` json NOT NULL, `verified` tinyint NOT NULL DEFAULT 0,
  `processed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_callback_event` (`channel`,`provider_event_id`), KEY `idx_refund_provider_no` (`provider_refund_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_transaction` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL, `channel` enum('wechat','alipay','balance','mock') NOT NULL,
  `provider_refund_no` varchar(128) NOT NULL, `amount` decimal(12,2) NOT NULL, `status` enum('created','processing','succeeded','failed') NOT NULL DEFAULT 'created',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, `completed_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_provider_refund` (`provider_refund_no`), UNIQUE KEY `uk_refund_idempotency` (`refund_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- mock 仅用于显式 sandbox 的本地闭环；生产环境退款执行器拒绝该渠道。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_invoice` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL, `invoice_profile_id` bigint unsigned NOT NULL,
  `profile_type` enum('personal', 'enterprise') NOT NULL, `title` varchar(255) NOT NULL,
  `tax_no` varchar(64) NOT NULL DEFAULT '', `email` varchar(255) NOT NULL DEFAULT '',
  `status` enum('requested','issued','rejected','voided') NOT NULL DEFAULT 'requested', `invoice_no` varchar(128) NOT NULL DEFAULT '',
  `file_url` varchar(1024) NOT NULL DEFAULT '', `rejection_reason` varchar(500) NOT NULL DEFAULT '', `requested_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `issued_at` datetime DEFAULT NULL,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_order_verification` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_id` bigint unsigned NOT NULL,
  `verify_code` varchar(32) NOT NULL DEFAULT '', `verify_code_hash` char(64) NOT NULL,
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
-- 退货物流使用独立、不可覆盖的登记记录；单个售后单只能有一条当前寄回单，避免把快递单号混入备注或证据文本。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_return_shipment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL,
  `carrier_name` varchar(128) NOT NULL, `tracking_no` varchar(128) NOT NULL,
  `remark` varchar(500) NOT NULL DEFAULT '', `submitted_by` bigint unsigned NOT NULL,
  `submitted_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_return_shipment` (`refund_id`), KEY `idx_tracking_no` (`tracking_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_event` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `refund_id` bigint unsigned NOT NULL,
  `from_status` varchar(32) NOT NULL DEFAULT '', `to_status` varchar(32) NOT NULL,
  `actor_type` enum('user','merchant','platform','system') NOT NULL, `actor_id` bigint unsigned NOT NULL,
  `reason` varchar(500) NOT NULL DEFAULT '', `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_refund_transition_key` (`refund_id`,`idempotency_key`), KEY `idx_refund_created` (`refund_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 退款监管导出只记录执行人、筛选摘要与导出行数；不保存退款原因、物流或用户资料副本。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_refund_export_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `query_fingerprint` char(64) NOT NULL,
  `row_count` int unsigned NOT NULL, `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_operator_created` (`operator_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 营销、会员、资产和分销域。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_coupon_template_view` (
  `coupon_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL, `name` varchar(128) NOT NULL,
  `discount_type` enum('amount','rate') NOT NULL, `discount_value` decimal(12,2) NOT NULL, `min_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, `status` tinyint NOT NULL DEFAULT 1, `version` bigint unsigned NOT NULL,
  PRIMARY KEY (`coupon_id`), KEY `idx_store_status` (`store_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 新人权益配置与券映射由平台运营维护；注册发券和展示均只读取本业务库事实。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_onboarding_policy` (
  `id` tinyint unsigned NOT NULL, `enabled` tinyint NOT NULL DEFAULT 0, `coupon_enabled` tinyint NOT NULL DEFAULT 0,
  `title` varchar(64) NOT NULL DEFAULT '新人礼', `description` varchar(255) NOT NULL DEFAULT '注册即享专属优惠券',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_onboarding_coupon` (
  `coupon_id` bigint unsigned NOT NULL, `enabled` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0,
  PRIMARY KEY (`coupon_id`), KEY `idx_enabled_sort` (`enabled`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_marketing_activity_view` (
  `activity_id` bigint unsigned NOT NULL, `store_id` bigint unsigned NOT NULL,
  `activity_type` enum('seckill','combination','presell','assist','discount','svip') NOT NULL, `name` varchar(128) NOT NULL,
  `rules` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `version` bigint unsigned NOT NULL,
  `starts_at` datetime DEFAULT NULL, `ends_at` datetime DEFAULT NULL, PRIMARY KEY (`activity_id`), KEY `idx_type_time` (`activity_type`,`status`,`starts_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_member_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL, `rank` int NOT NULL, `rules` json NOT NULL,
  `benefits` json NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `version` bigint unsigned NOT NULL DEFAULT 1,
  `deleted_at` datetime DEFAULT NULL, PRIMARY KEY (`id`), UNIQUE KEY `uk_rank` (`rank`), KEY `idx_member_level_visible` (`deleted_at`,`status`,`rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='version')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `version` bigint unsigned NOT NULL DEFAULT 1 AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='deleted_at')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `deleted_at` datetime DEFAULT NULL AFTER `version`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_member_level_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `level_id` bigint unsigned DEFAULT NULL, `previous_level_id` bigint unsigned DEFAULT NULL,
  `change_type` enum('initial', 'upgrade', 'downgrade', 'manual') NOT NULL, `note` varchar(500) NOT NULL DEFAULT '',
  `idempotency_key` varchar(128) DEFAULT NULL, `operator_admin_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (`id`), KEY `idx_user_created` (`user_id`,`created_at`),
  UNIQUE KEY `uk_member_level_idempotency` (`user_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_sign` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL, `sign_date` date NOT NULL,
  `points` bigint NOT NULL, `continuous_days` int NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_sign_date` (`user_id`,`sign_date`), KEY `idx_user_sign_date` (`user_id`,`sign_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 搜索历史只保存用户 ID、关键词、来源和时间；不复制账号、手机、地址或设备标识。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_search_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `keyword` varchar(128) NOT NULL, `source` enum('pc','h5','mini') NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), KEY `idx_user_search_created` (`user_id`,`created_at`), KEY `idx_search_visible` (`deleted_at`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 清理必须留有可重放的审计，但审计不复制被清理的关键词内容。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_search_record_clear_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `reason` varchar(500) NOT NULL, `idempotency_key` varchar(128) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL, `cleared_count` int unsigned NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_search_clear_operator_key` (`operator_admin_id`,`idempotency_key`), KEY `idx_search_clear_user_created` (`user_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_search_record_export_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `query_fingerprint` char(64) NOT NULL,
  `row_count` int unsigned NOT NULL, `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_search_export_operator_created` (`operator_admin_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_recharge_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `recharge_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `amount` decimal(12,2) NOT NULL, `bonus_amount` decimal(12,2) NOT NULL DEFAULT 0, `status` enum('pending','paid','closed') NOT NULL DEFAULT 'pending',
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_recharge_no` (`recharge_no`), UNIQUE KEY `uk_user_idempotency` (`user_id`,`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 用户充值与 SVIP 购买独立于商品交易：商品支付表带 group_order_id，不能复用。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_recharge_plan` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `amount` decimal(12,2) NOT NULL, `bonus_amount` decimal(12,2) NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, `version` bigint unsigned NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_recharge_plan_visible` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_svip_plan` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `price` decimal(12,2) NOT NULL, `plan_type` enum('trial','period','lifetime') NOT NULL,
  `duration_days` int unsigned DEFAULT NULL, `benefits` json NOT NULL,
  `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_svip_plan_visible` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 付费会员权益是可复用的 C 端展示主数据；套餐保存权益名称快照，已生效订单不受后续配置改动影响。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_svip_interest` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `name` varchar(64) NOT NULL,
  `description` varchar(500) NOT NULL DEFAULT '', `icon_url` varchar(1024) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1, `sort` int NOT NULL DEFAULT 0, `version` bigint unsigned NOT NULL DEFAULT 1,
  `deleted_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_svip_interest_name` (`name`), KEY `idx_svip_interest_visible` (`deleted_at`,`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_svip_order` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `plan_id` bigint unsigned NOT NULL, `plan_name` varchar(64) NOT NULL, `plan_type` enum('trial','period','lifetime') NOT NULL,
  `duration_days` int unsigned DEFAULT NULL, `amount` decimal(12,2) NOT NULL,
  `status` enum('pending','paid','closed') NOT NULL DEFAULT 'pending', `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_svip_order_no` (`order_no`), UNIQUE KEY `uk_svip_user_key` (`user_id`,`idempotency_key`), KEY `idx_svip_user_status` (`user_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_svip` (
  `user_id` bigint unsigned NOT NULL, `status` enum('trial','period','lifetime') NOT NULL,
  `expires_at` datetime DEFAULT NULL, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`), KEY `idx_svip_expires` (`status`,`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_funding_payment` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `order_type` enum('recharge','svip') NOT NULL,
  `funding_order_id` bigint unsigned NOT NULL, `user_id` bigint unsigned NOT NULL,
  `channel` enum('wechat') NOT NULL, `out_trade_no` varchar(64) NOT NULL, `amount` decimal(12,2) NOT NULL,
  `status` enum('created','processing','succeeded','failed','closed') NOT NULL DEFAULT 'created',
  `provider_transaction_no` varchar(128) DEFAULT NULL, `provider_payload` json DEFAULT NULL,
  `callback_idempotency_key` varchar(191) DEFAULT NULL, `paid_at` datetime DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_funding_out_trade_no` (`out_trade_no`), UNIQUE KEY `uk_funding_order_channel` (`order_type`,`funding_order_id`,`channel`), UNIQUE KEY `uk_funding_callback_key` (`callback_idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_funding_payment_callback` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `channel` enum('wechat') NOT NULL,
  `provider_event_id` varchar(191) NOT NULL, `out_trade_no` varchar(64) NOT NULL, `payload` json NOT NULL,
  `verified` tinyint NOT NULL DEFAULT 0, `processed_at` datetime DEFAULT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_funding_channel_event` (`channel`,`provider_event_id`), KEY `idx_funding_callback_order` (`out_trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_withdrawal_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `withdrawal_no` varchar(64) NOT NULL, `user_id` bigint unsigned NOT NULL,
  `amount` decimal(12,2) NOT NULL, `channel` enum('wechat','bank') NOT NULL, `account_snapshot` json NOT NULL,
  `status` enum('applied','reviewing','approved','paying','paid','rejected') NOT NULL DEFAULT 'applied', `review_note` varchar(500) NOT NULL DEFAULT '', `reviewed_by` bigint unsigned DEFAULT NULL, `reviewed_at` datetime DEFAULT NULL,
  `payout_idempotency_key` varchar(128) DEFAULT NULL, `payout_reference` varchar(128) DEFAULT NULL, `paid_by` bigint unsigned DEFAULT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, `paid_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_withdrawal_no` (`withdrawal_no`), UNIQUE KEY `uk_user_withdrawal_key` (`user_id`,`idempotency_key`), UNIQUE KEY `uk_user_payout_key` (`user_id`,`payout_idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_withdrawal_application' AND COLUMN_NAME='reviewed_by')=0,
    'ALTER TABLE `qixi_crm_b_withdrawal_application` ADD COLUMN `reviewed_by` bigint unsigned DEFAULT NULL AFTER `review_note`, ADD COLUMN `reviewed_at` datetime DEFAULT NULL AFTER `reviewed_by`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_withdrawal_application' AND COLUMN_NAME='payout_idempotency_key')=0,
    'ALTER TABLE `qixi_crm_b_withdrawal_application` ADD COLUMN `payout_idempotency_key` varchar(128) DEFAULT NULL AFTER `reviewed_at`, ADD COLUMN `payout_reference` varchar(128) DEFAULT NULL AFTER `payout_idempotency_key`, ADD COLUMN `paid_by` bigint unsigned DEFAULT NULL AFTER `payout_reference`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_promoter` (
  `user_id` bigint unsigned NOT NULL, `status` tinyint NOT NULL DEFAULT 1, `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`), KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
-- 平台批量设置推广员资格的审计；资格仅影响后续分销准入，不重算历史佣金或关系事实。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_promoter_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_ids_json` json NOT NULL,
  `status` tinyint NOT NULL, `reason` varchar(500) NOT NULL,
  `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0, `idempotency_key` varchar(128) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_relation` (
  `user_id` bigint unsigned NOT NULL, `parent_user_id` bigint unsigned DEFAULT NULL, `bound_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`), KEY `idx_parent` (`parent_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_relation_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `previous_parent_user_id` bigint unsigned DEFAULT NULL, `parent_user_id` bigint unsigned DEFAULT NULL,
  `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_distribution_relation_idempotency` (`user_id`,`idempotency_key`),
  KEY `idx_distribution_relation_audit` (`user_id`,`created_at`)
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
-- 平台向用户投递站内图文通知的审计。公众号/短信等外部渠道必须另走具备真实凭据的受控投递服务。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_notification_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `user_id` bigint unsigned NOT NULL,
  `notification_id` bigint unsigned NOT NULL, `title` varchar(255) NOT NULL,
  `request_fingerprint` char(64) NOT NULL, `reason` varchar(500) NOT NULL, `operator_admin_id` bigint unsigned NOT NULL DEFAULT 0,
  `idempotency_key` varchar(128) NOT NULL, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), UNIQUE KEY `uk_user_idempotency` (`user_id`,`idempotency_key`), KEY `idx_notification` (`notification_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE IF NOT EXISTS qixi_crm_b_user_feedback (
  id bigint unsigned NOT NULL AUTO_INCREMENT, user_id bigint unsigned NOT NULL,
  category_id bigint unsigned DEFAULT NULL, type varchar(32) NOT NULL, content varchar(1000) NOT NULL, status varchar(16) NOT NULL, reply varchar(1000) NOT NULL,
  deleted_at datetime DEFAULT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), KEY idx_user_time (user_id,created_at), KEY idx_feedback_status_deleted (status,deleted_at), KEY idx_feedback_category (category_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS qixi_crm_b_user_feedback_category (
  id bigint unsigned NOT NULL AUTO_INCREMENT, name varchar(32) NOT NULL, sort int unsigned NOT NULL DEFAULT 0,
  status tinyint NOT NULL DEFAULT 1, deleted_at datetime DEFAULT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id), UNIQUE KEY uk_feedback_category_name (name), KEY idx_feedback_category_active (status,deleted_at,sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS qixi_crm_b_user_feedback_category_audit (
  id bigint unsigned NOT NULL AUTO_INCREMENT, category_id bigint unsigned NOT NULL, action varchar(32) NOT NULL,
  name varchar(32) NOT NULL DEFAULT '', sort int unsigned NOT NULL DEFAULT 0, status tinyint NOT NULL DEFAULT 1,
  operator_admin_id bigint unsigned NOT NULL, idempotency_key varchar(128) NOT NULL, created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id), UNIQUE KEY uk_feedback_category_idempotency (idempotency_key), KEY idx_feedback_category_audit (category_id,created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
CREATE TABLE IF NOT EXISTS qixi_crm_b_user_feedback_audit (
  id bigint unsigned NOT NULL AUTO_INCREMENT, feedback_id bigint unsigned NOT NULL,
  from_status varchar(16) NOT NULL, to_status varchar(16) NOT NULL, action varchar(16) NOT NULL, reply varchar(1000) NOT NULL DEFAULT '',
  operator_admin_id bigint unsigned NOT NULL, idempotency_key varchar(128) NOT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id), UNIQUE KEY uk_idempotency (idempotency_key), KEY idx_feedback_time (feedback_id,created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 订单仅对用户软归档；财务、售后和履约记录保留。

-- CRMEB open_screen 的 C 端开屏配置；关闭时接口只返回 enabled=false。
CREATE TABLE IF NOT EXISTS `qixi_crm_b_open_screen_campaign` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT, `title` varchar(128) NOT NULL DEFAULT '',
  `image_url` varchar(1024) NOT NULL, `link_url` varchar(1024) NOT NULL DEFAULT '',
  `duration_sec` tinyint unsigned NOT NULL DEFAULT 3, `space_hours` int unsigned NOT NULL DEFAULT 24,
  `enabled` tinyint NOT NULL DEFAULT 0, `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`), KEY `idx_enabled` (`enabled`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
