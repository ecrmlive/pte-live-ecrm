-- 阶段 5：平台券 + 用户券 + 分销绑定日志（营销竖切）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_coupon` (
  `coupon_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '优惠券表ID',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '0=平台',
  `is_timeout` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `start_time` timestamp NULL DEFAULT NULL,
  `end_time` timestamp NULL DEFAULT NULL,
  `is_limited` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `total_count` int(10) unsigned NOT NULL DEFAULT '0',
  `remain_count` int(10) unsigned NOT NULL DEFAULT '0',
  `send_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0领取',
  `full_reduction` decimal(8,2) unsigned NOT NULL DEFAULT '0.00',
  `title` varchar(64) NOT NULL,
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00',
  `use_min_price` int(11) NOT NULL DEFAULT '0' COMMENT '最低消费（元取整）',
  `coupon_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0有效天数 1固定时段',
  `coupon_time` int(10) unsigned NOT NULL DEFAULT '0',
  `use_start_time` timestamp NULL DEFAULT NULL,
  `use_end_time` timestamp NULL DEFAULT NULL,
  `sort` int(10) unsigned NOT NULL DEFAULT '1',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0关 1开 -1失效',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `type` tinyint(4) NOT NULL DEFAULT '10' COMMENT '0店铺 10平台通用',
  PRIMARY KEY (`coupon_id`),
  KEY `mer_id` (`mer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券';

CREATE TABLE IF NOT EXISTS `qixi_store_coupon_issue_user` (
  `uid` int(11) NOT NULL DEFAULT '0',
  `coupon_id` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY `uid_coupon` (`uid`,`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='领券记录';

CREATE TABLE IF NOT EXISTS `qixi_store_coupon_user` (
  `coupon_user_id` int(11) NOT NULL AUTO_INCREMENT,
  `coupon_id` int(10) unsigned NOT NULL DEFAULT '0',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0',
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `coupon_title` varchar(32) NOT NULL DEFAULT '',
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00',
  `use_min_price` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `start_time` timestamp NULL DEFAULT NULL,
  `end_time` timestamp NULL DEFAULT NULL,
  `use_time` timestamp NULL DEFAULT NULL,
  `type` varchar(16) NOT NULL DEFAULT 'receive',
  `send_id` int(10) unsigned DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未用 1已用 2过期',
  `is_fail` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`coupon_user_id`),
  KEY `uid` (`uid`),
  KEY `coupon_id` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券';

CREATE TABLE IF NOT EXISTS `qixi_user_spread_log` (
  `user_spread_log_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `old_spread_uid` int(10) unsigned NOT NULL DEFAULT '0',
  `spread_uid` int(10) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_spread_log_id`),
  KEY `uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='分销绑定日志';

-- 主单记录选用的用户券
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_group_order' AND COLUMN_NAME = 'coupon_id'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_group_order` ADD COLUMN `coupon_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''用户券coupon_user_id'' AFTER `coupon_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 主单/子单券字段（幂等）
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'coupon_price'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_order`
     ADD COLUMN `coupon_id` varchar(128) NOT NULL DEFAULT '''' COMMENT ''店铺券 coupon_user_id'' AFTER `cost`,
     ADD COLUMN `coupon_price` decimal(8,2) NOT NULL DEFAULT 0.00 COMMENT ''店铺券优惠'' AFTER `coupon_id`,
     ADD COLUMN `platform_coupon_price` decimal(8,2) NOT NULL DEFAULT 0.00 COMMENT ''分摊平台券优惠'' AFTER `coupon_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `qixi_user_bill` (
  `bill_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `link_id` varchar(32) NOT NULL DEFAULT '0',
  `pm` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0支出 1获得',
  `title` varchar(64) NOT NULL DEFAULT '',
  `category` varchar(64) NOT NULL DEFAULT '',
  `type` varchar(64) NOT NULL DEFAULT '',
  `number` decimal(11,2) unsigned NOT NULL DEFAULT '0.00',
  `balance` decimal(11,2) unsigned NOT NULL DEFAULT '0.00',
  `mark` varchar(512) NOT NULL DEFAULT '',
  `mer_id` int(10) unsigned DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) NOT NULL DEFAULT '1',
  PRIMARY KEY (`bill_id`),
  KEY `uid` (`uid`),
  KEY `type` (`category`,`type`,`link_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户账单/佣金流水';

-- 演示平台券：满 30 减 5，有效 30 天，可领
INSERT INTO `qixi_store_coupon` (
  `coupon_id`, `mer_id`, `is_timeout`, `is_limited`, `total_count`, `remain_count`,
  `send_type`, `title`, `coupon_price`, `use_min_price`, `coupon_type`, `coupon_time`,
  `sort`, `status`, `is_del`, `type`
)
SELECT 1, 0, 0, 1, 1000, 1000, 0, '平台新人满减券', 5.00, 30, 0, 30, 1, 1, 0, 10
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_coupon` WHERE `coupon_id` = 1);

-- 演示店铺券（商户1）：满 20 减 3
INSERT INTO `qixi_store_coupon` (
  `coupon_id`, `mer_id`, `is_timeout`, `is_limited`, `total_count`, `remain_count`,
  `send_type`, `title`, `coupon_price`, `use_min_price`, `coupon_type`, `coupon_time`,
  `sort`, `status`, `is_del`, `type`
)
SELECT 2, 1, 0, 1, 500, 500, 0, '店铺满减券', 3.00, 20, 0, 30, 1, 1, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_coupon` WHERE `coupon_id` = 2);

-- 演示推广员：给 demo 用户标记可推广（uid=1）
UPDATE `qixi_user` SET `is_promoter` = 1 WHERE `account` = 'demo' AND `is_promoter` = 0;

-- 平台菜单：营销
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 17 AS menu_id, 0 AS pid, '/marketing' AS path, 'GiftOutlined' AS icon, '营销' AS menu_name, 'Marketing' AS route, '' AS params, 65 AS sort, 1 AS is_show, 1 AS is_mer, 1 AS is_menu, 0 AS is_agent
  UNION ALL SELECT 18, 17, '/marketing/coupon', '', '平台优惠券', 'MarketingCoupon', '', 64, 1, 1, 1, 0
  UNION ALL SELECT 19, 17, '/marketing/spread', '', '分销监管', 'MarketingSpread', '', 63, 1, 1, 1, 0
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 17 AND `is_mer` = 1);

-- 商户菜单：店铺券（P0 页可后接，菜单先预埋）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 115 AS menu_id, 0 AS pid, '/marketing' AS path, 'GiftOutlined' AS icon, '营销' AS menu_name, 'MerMarketing' AS route, '' AS params, 75 AS sort, 1 AS is_show, 2 AS is_mer, 1 AS is_menu, 2 AS is_agent
  UNION ALL SELECT 116, 115, '/marketing/coupon', '', '店铺优惠券', 'MerMarketingCoupon', '', 74, 1, 2, 1, 2
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 115 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',17,18,19')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%17%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',115,116')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%115%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase5-promotion-coupon-spread', '阶段5：平台券领取选券 + 分销绑定日志'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase5-promotion-coupon-spread');
