-- 阶段 3：购物车 / 地址 / 主单子单 / 两商户可售种子（P0 交易主链）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_user_address` (
  `address_id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户地址id',
  `uid` int(10) unsigned NOT NULL COMMENT '用户id',
  `real_name` varchar(32) NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '收货人电话',
  `province` varchar(64) NOT NULL DEFAULT '' COMMENT '省',
  `province_id` int(10) unsigned DEFAULT '0',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT '市',
  `city_id` int(11) NOT NULL DEFAULT '0',
  `district` varchar(64) NOT NULL DEFAULT '' COMMENT '区',
  `district_id` int(10) unsigned DEFAULT '0',
  `detail` varchar(256) NOT NULL DEFAULT '' COMMENT '详细地址',
  `post_code` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '邮编',
  `is_default` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`address_id`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户地址表';

CREATE TABLE IF NOT EXISTS `qixi_store_cart` (
  `cart_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车ID',
  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户ID',
  `product_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0普通',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `product_attr_unique` varchar(16) NOT NULL DEFAULT '' COMMENT 'SKU unique',
  `cart_num` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '数量',
  `is_pay` tinyint(1) NOT NULL DEFAULT '0',
  `is_del` tinyint(1) NOT NULL DEFAULT '0',
  `is_new` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1立即购买',
  `is_fail` tinyint(3) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`cart_id`),
  KEY `uid_pay` (`uid`,`is_pay`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='购物车表';

CREATE TABLE IF NOT EXISTS `qixi_store_group_order` (
  `group_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_order_sn` varchar(32) NOT NULL COMMENT '主单号',
  `uid` int(10) unsigned NOT NULL,
  `total_postage` decimal(8,2) NOT NULL DEFAULT '0.00',
  `total_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `total_num` int(10) unsigned NOT NULL DEFAULT '0',
  `coupon_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `real_name` varchar(32) NOT NULL DEFAULT '',
  `user_phone` varchar(18) NOT NULL DEFAULT '',
  `user_address` varchar(128) NOT NULL DEFAULT '',
  `pay_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `pay_postage` decimal(8,2) NOT NULL DEFAULT '0.00',
  `cost` decimal(8,2) NOT NULL DEFAULT '0.00',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `pay_time` timestamp NULL DEFAULT NULL,
  `pay_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0余额 7=mock',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_combine` tinyint(3) DEFAULT '0',
  `activity_type` tinyint(3) DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_order_id`),
  UNIQUE KEY `group_order_sn` (`group_order_sn`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户主单';

CREATE TABLE IF NOT EXISTS `qixi_store_order` (
  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_order_id` int(10) unsigned NOT NULL DEFAULT '0',
  `order_sn` varchar(32) NOT NULL,
  `uid` int(10) unsigned NOT NULL,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0',
  `real_name` varchar(32) NOT NULL DEFAULT '',
  `user_phone` varchar(18) NOT NULL DEFAULT '',
  `user_address` varchar(128) NOT NULL DEFAULT '',
  `cart_id` varchar(256) NOT NULL DEFAULT '',
  `total_num` int(10) unsigned NOT NULL DEFAULT '0',
  `total_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `total_postage` decimal(8,2) NOT NULL DEFAULT '0.00',
  `pay_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `pay_postage` decimal(8,2) NOT NULL DEFAULT '0.00',
  `cost` decimal(8,2) NOT NULL DEFAULT '0.00',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `pay_time` timestamp NULL DEFAULT NULL,
  `pay_type` tinyint(1) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待发货1待收货2待评价3完成-1退款',
  `delivery_type` varchar(32) DEFAULT NULL,
  `delivery_name` varchar(50) DEFAULT NULL,
  `delivery_id` varchar(255) DEFAULT NULL,
  `order_type` tinyint(3) DEFAULT '0',
  `is_virtual` tinyint(3) DEFAULT '0',
  `mark` varchar(512) NOT NULL DEFAULT '',
  `verify_code` char(16) DEFAULT NULL,
  `verify_time` timestamp NULL DEFAULT NULL,
  `activity_type` tinyint(3) NOT NULL DEFAULT '0',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_system_del` tinyint(1) DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`order_id`),
  UNIQUE KEY `order_sn` (`order_sn`),
  KEY `idx_mer` (`mer_id`,`paid`,`status`),
  KEY `idx_group` (`group_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户子单';

CREATE TABLE IF NOT EXISTS `qixi_store_order_product` (
  `order_product_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `cart_id` int(10) unsigned NOT NULL DEFAULT '0',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0',
  `product_sku` char(12) NOT NULL DEFAULT '',
  `product_num` int(10) unsigned NOT NULL DEFAULT '0',
  `product_type` int(11) NOT NULL DEFAULT '0',
  `product_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `total_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `cost` decimal(10,2) NOT NULL DEFAULT '0.00',
  `cart_info` text COMMENT '商品快照JSON',
  `is_refund` tinyint(3) NOT NULL DEFAULT '0',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0',
  `is_reply` tinyint(3) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`order_product_id`),
  KEY `order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单商品行';

-- 演示用户余额（余额支付）
UPDATE `qixi_user` SET `now_money` = 1000.00 WHERE `account` = 'demo' AND `now_money` < 100;

-- 第二商户 + 可售商品（两商户一单）
INSERT INTO `qixi_merchant` (`mer_id`, `mer_name`, `real_name`, `mer_phone`, `mer_address`, `mark`, `status`, `mer_state`, `is_del`, `is_audit`)
SELECT 2, '星河数码', '王五', '13900000003', '上海市演示路2号', '阶段3第二商户', 1, 1, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 2);

INSERT INTO `qixi_merchant_admin` (`merchant_admin_id`, `mer_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `level`, `status`, `is_del`)
SELECT 2, 2, 'mer2', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '数码店长', '13900000003', '2', 0, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_admin` WHERE `account` = 'mer2');

INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `cost`, `stock`, `sales`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`, `mer_status`
)
SELECT 3, 2, 'Type-C 快充数据线 1.2m', '阶段3第二商户可售', '数码,线材', 1, 1, 1,
  2, '条', 19.90, 29.90, 5.00, 999, 0, 0, '', '',
  '2', 0, 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 3);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 4, 3, '{}', '默认', 999, 0, 5.00, 29.90, 19.90, 'sku000000004', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 4);

-- 演示用户默认地址
INSERT INTO `qixi_user_address` (
  `address_id`, `uid`, `real_name`, `phone`, `province`, `city`, `district`, `detail`, `post_code`, `is_default`
)
SELECT 1, 1, '演示用户', '13700000001', '上海市', '上海市', '浦东新区', '演示路88号', 200000, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_user_address` WHERE `address_id` = 1);

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase3-trade-cart-order', '阶段3：购物车地址主单子单与两商户可售种子'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase3-trade-cart-order');
