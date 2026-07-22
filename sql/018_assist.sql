-- 阶段 6：助力活动竖切（邀请助力 → 满员可下单）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_product_assist` (
  `product_assist_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `start_time` datetime NOT NULL,
  `end_time` datetime NOT NULL,
  `status` int(10) NOT NULL DEFAULT 1 COMMENT '平台 1开启 0结束',
  `pay_count` int(10) unsigned NOT NULL DEFAULT 0,
  `assist_count` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '需助力人数',
  `assist_user_count` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '单人可助力次数',
  `product_id` int(10) unsigned NOT NULL DEFAULT 0,
  `assist_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '助力价',
  `stock` int(10) unsigned NOT NULL DEFAULT 0,
  `is_show` tinyint(3) NOT NULL DEFAULT 1,
  `store_name` varchar(128) NOT NULL DEFAULT '',
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `store_info` varchar(255) DEFAULT NULL,
  `is_del` int(10) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `product_status` int(11) NOT NULL DEFAULT 1 COMMENT '1通过',
  `refusal` varchar(255) DEFAULT NULL,
  `action_status` int(11) NOT NULL DEFAULT 1,
  PRIMARY KEY (`product_assist_id`),
  KEY `idx_mer` (`mer_id`),
  KEY `idx_product` (`product_id`),
  KEY `idx_window` (`start_time`,`end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='助力活动';

CREATE TABLE IF NOT EXISTS `qixi_store_product_assist_set` (
  `product_assist_set_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_assist_id` int(10) unsigned NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `status` int(11) NOT NULL DEFAULT 1 COMMENT '-1失败 1进行中 10已完成 20已支付',
  `assist_count` int(10) unsigned NOT NULL DEFAULT 0,
  `assist_user_count` int(10) unsigned NOT NULL DEFAULT 0,
  `yet_assist_count` int(10) unsigned NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `is_del` int(10) NOT NULL DEFAULT 0,
  PRIMARY KEY (`product_assist_set_id`),
  KEY `idx_assist` (`product_assist_id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='助力发起';

CREATE TABLE IF NOT EXISTS `qixi_store_product_assist_user` (
  `product_assist_user_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_assist_set_id` int(10) unsigned NOT NULL,
  `product_assist_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `nickname` varchar(50) DEFAULT NULL,
  `avatar_img` varchar(256) DEFAULT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_assist_user_id`),
  UNIQUE KEY `uk_set_uid` (`product_assist_set_id`,`uid`),
  KEY `idx_assist` (`product_assist_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='助力记录';

-- 演示商品 13：助力价 ¥9.90（原价 29.90）
INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `cost`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`, `mer_status`
)
SELECT 13, 1, '助力好物 · 演示礼盒', '邀请好友助力后低价购', '助力', 1, 1, 1,
  2, '盒', 29.90, 39.90, 8.00, 999, 0, '', '',
  '2', 0, 3, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 13);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 13, 13, '默认', 'assist-demo', 999, 0, 8.00, 39.90, 29.90, 'sku000000013', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 13);

INSERT INTO `qixi_store_product_assist` (
  `product_assist_id`, `start_time`, `end_time`, `status`, `assist_count`, `assist_user_count`,
  `product_id`, `assist_price`, `stock`, `is_show`, `store_name`, `mer_id`, `store_info`,
  `product_status`, `action_status`
)
SELECT 1, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 30 DAY), 1, 1, 1,
  13, 9.90, 50, 1, '助力好物 · 演示礼盒', 1, '阶段6助力竖切',
  1, 1
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 13)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_assist` WHERE `product_assist_id` = 1);

-- 演示：已满员待下单的助力单（uid=1）
INSERT INTO `qixi_store_product_assist_set` (
  `product_assist_set_id`, `product_assist_id`, `product_id`, `uid`, `status`,
  `assist_count`, `assist_user_count`, `yet_assist_count`, `mer_id`
)
SELECT 1, 1, 13, 1, 10, 1, 1, 1, 1
WHERE EXISTS (SELECT 1 FROM `qixi_user` WHERE `uid` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_assist_set` WHERE `product_assist_set_id` = 1);

INSERT INTO `qixi_store_product_assist_user` (
  `product_assist_user_id`, `product_assist_set_id`, `product_assist_id`, `uid`, `nickname`
)
SELECT 1, 1, 1, 2, '助力好友'
WHERE EXISTS (SELECT 1 FROM `qixi_store_product_assist_set` WHERE `product_assist_set_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_assist_user` WHERE `product_assist_user_id` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 29, 17, '/marketing/assist', '', '助力监管', 'MarketingAssist', '', 72, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 29 AND `is_mer` = 1);

UPDATE `qixi_system_role` SET `rules` = CONCAT(`rules`, ',29')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%29%';

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 124, 115, '/marketing/assist', '', '助力活动', 'MerMarketingAssist', '', 72, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 124 AND `is_mer` = 2);

UPDATE `qixi_system_role` SET `rules` = CONCAT(`rules`, ',124')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%124%';

UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 8, 'name', '助力', 'icon', '', 'url', '/pages/assist/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '助力', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-assist', '阶段6：助力活动竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-assist');
