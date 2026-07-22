-- 阶段 6：DIY 首页装修 + 秒杀活动（阶段5顺延的「一种活动」）
USE `qixi_mergers`;

-- 订单行活动 id（秒杀 active_id 等）
SET @db := DATABASE();
SET @has_aid := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_store_order_product' AND COLUMN_NAME = 'activity_id'
);
SET @sql := IF(@has_aid = 0,
  'ALTER TABLE `qixi_store_order_product` ADD COLUMN `activity_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''活动关联 id'' AFTER `product_type`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `qixi_diy` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `version` varchar(255) NOT NULL DEFAULT '' COMMENT '版本号',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '页面名称',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '网站标题',
  `cover_image` varchar(255) NOT NULL DEFAULT '',
  `template_name` varchar(255) NOT NULL DEFAULT '' COMMENT '模板名称',
  `value` longtext COMMENT '页面 JSON',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '页面类型 0首页',
  `is_show` tinyint(1) NOT NULL DEFAULT '0' COMMENT '显示首页',
  `is_diy` tinyint(1) NOT NULL DEFAULT '1',
  `is_del` tinyint(1) NOT NULL DEFAULT '0',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '0=平台',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1平台默认',
  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_type` (`status`,`type`,`mer_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='DIY装修页';

CREATE TABLE IF NOT EXISTS `qixi_store_seckill_time` (
  `seckill_time_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT '',
  `start_time` int(10) unsigned NOT NULL COMMENT '开始小时 0-23',
  `end_time` int(10) unsigned NOT NULL COMMENT '结束小时 1-24',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `pic` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`seckill_time_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='秒杀场次';

CREATE TABLE IF NOT EXISTS `qixi_store_seckill_active` (
  `seckill_active_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `seckill_time_ids` varchar(255) DEFAULT '' COMMENT '场次id逗号分隔',
  `start_day` date NOT NULL,
  `end_day` date NOT NULL,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `product_id` int(10) unsigned NOT NULL,
  `seckill_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '秒杀价（本仓库扩展）',
  `once_pay_count` int(10) unsigned NOT NULL DEFAULT 1,
  `all_pay_count` int(10) unsigned NOT NULL DEFAULT 0,
  `active_status` tinyint(4) NOT NULL DEFAULT 1 COMMENT '0未开始 1进行中 -1结束',
  `status` tinyint(3) unsigned NOT NULL DEFAULT 1 COMMENT '0关 1开',
  `create_time` int(11) NOT NULL DEFAULT 0,
  `update_time` int(11) NOT NULL DEFAULT 0,
  `delete_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`seckill_active_id`),
  KEY `idx_mer_product` (`mer_id`,`product_id`,`status`),
  KEY `idx_day` (`start_day`,`end_day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='秒杀活动商品';

-- 平台默认 DIY 首页
INSERT INTO `qixi_diy` (`id`, `version`, `name`, `title`, `template_name`, `value`, `status`, `type`, `is_show`, `is_diy`, `mer_id`, `is_default`)
SELECT 1, '1.0', '平台首页', '栖息商城', 'home',
  '{"banners":[{"id":1,"title":"夏日秒杀","image":"","url":"/pages/seckill/list"},{"id":2,"title":"积分好物","image":"","url":"/pages/points/list"}],"menus":[{"id":1,"name":"秒杀","icon":"","url":"/pages/seckill/list"},{"id":2,"name":"积分","icon":"","url":"/pages/points/list"},{"id":3,"name":"领券","icon":"","url":"/pages/coupon/center"}]}',
  1, 0, 1, 1, 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_diy` WHERE `id` = 1);

-- 全天场次（演示）
INSERT INTO `qixi_store_seckill_time` (`seckill_time_id`, `title`, `start_time`, `end_time`, `status`)
SELECT 1, '全天场', 0, 24, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_seckill_time` WHERE `seckill_time_id` = 1);

-- 商品1 秒杀（若存在）
INSERT INTO `qixi_store_seckill_active` (
  `seckill_active_id`, `name`, `seckill_time_ids`, `start_day`, `end_day`,
  `mer_id`, `product_id`, `seckill_price`, `once_pay_count`, `active_status`, `status`,
  `create_time`, `update_time`
)
SELECT 1, '演示秒杀 · 商品1', '1', CURDATE(), DATE_ADD(CURDATE(), INTERVAL 30 DAY),
  1, 1, 9.90, 2, 1, 1, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 1 AND `is_del` = 0)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_seckill_active` WHERE `seckill_active_id` = 1);

-- 平台菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 22 AS menu_id, 20 AS pid, '/content/diy' AS path, '' AS icon, '页面装修' AS menu_name, 'ContentDiy' AS route, '' AS params, 53 AS sort, 1 AS is_show, 1 AS is_mer, 1 AS is_menu, 0 AS is_agent
  UNION ALL SELECT 23, 17, '/marketing/seckill', '', '秒杀监管', 'MarketingSeckill', '', 44, 1, 1, 1, 0
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 22 AND `is_mer` = 1);

-- 商户菜单：秒杀活动
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 117, 115, '/marketing/seckill', '', '秒杀活动', 'MerMarketingSeckill', '', 74, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 117 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',22,23')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%22%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',117')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%117%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-diy-seckill', '阶段6：DIY首页 + 秒杀活动'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-diy-seckill');
