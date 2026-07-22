-- 阶段 6：积分抵扣/积分商城字段 + 积分商品种子；附带 6a 公告种子
USE `qixi_mergers`;

-- 主单积分字段（幂等）
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_group_order' AND COLUMN_NAME = 'integral'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_group_order`
     ADD COLUMN `integral` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''使用积分'' AFTER `total_num`,
     ADD COLUMN `integral_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT ''积分抵扣金额'' AFTER `integral`,
     ADD COLUMN `give_integral` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''赠送积分'' AFTER `integral_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'integral'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_order`
     ADD COLUMN `integral` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''使用积分'' AFTER `cost`,
     ADD COLUMN `integral_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT ''积分抵扣金额'' AFTER `integral`,
     ADD COLUMN `give_integral` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''赠送积分'' AFTER `integral_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 若仅缺 give_integral
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'give_integral'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_order` ADD COLUMN `give_integral` int(10) unsigned NOT NULL DEFAULT 0 COMMENT ''赠送积分'' AFTER `integral_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 平台公告
CREATE TABLE IF NOT EXISTS `qixi_system_notice` (
  `notice_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(128) NOT NULL DEFAULT '',
  `content` text NOT NULL,
  `is_show` tinyint(1) NOT NULL DEFAULT '1',
  `sort` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`notice_id`),
  KEY `idx_show` (`is_show`,`is_del`,`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='平台公告';

INSERT INTO `qixi_system_notice` (`notice_id`, `title`, `content`, `is_show`, `sort`)
SELECT 1, '积分商城上线', '可用积分兑换积分商品；普通单可勾选积分抵扣（最多应付20%，100积分=1元）。', 1, 100
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_notice` WHERE `notice_id` = 1);

-- 商品兑换积分列 + 商户抵扣开关
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_product' AND COLUMN_NAME = 'integral'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_store_product` ADD COLUMN `integral` int(10) NOT NULL DEFAULT 0 COMMENT ''积分商城兑换积分'' AFTER `price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_merchant' AND COLUMN_NAME = 'mer_integral_status'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_merchant`
     ADD COLUMN `mer_integral_status` tinyint(1) NOT NULL DEFAULT 1 COMMENT ''是否允许积分抵扣'' AFTER `mer_money`,
     ADD COLUMN `mer_integral_rate` int(10) NOT NULL DEFAULT 100 COMMENT ''多少积分=1元'' AFTER `mer_integral_status`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE `qixi_merchant` SET `mer_integral_status` = 1, `mer_integral_rate` = 100 WHERE `mer_id` IN (1, 2);

-- 积分商品：type=1 + product_type=20；integral=所需积分
INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `integral`, `ot_price`, `cost`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`, `mer_status`
)
SELECT 10, 1, '积分兑换 · 帆布袋', '积分商城演示商品', '积分', 1, 1, 1,
  2, '个', 0.00, 500, 500.00, 1.00, 200, 0, '', '',
  '2', 1, 20, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 10);

-- 幂等修正：已存在的演示商品对齐积分标记
UPDATE `qixi_store_product`
SET `type` = 1, `product_type` = 20, `integral` = 500, `ot_price` = 500.00, `price` = 0.00
WHERE `product_id` = 10;

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 10, 10, '默认', 'points-bag', 200, 0, 1.00, 500.00, 0.01, 'sku000000010', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 10);

-- 演示用户积分
UPDATE `qixi_user` SET `integral` = 2000 WHERE `account` = 'demo' AND `integral` < 2000;

-- 平台菜单：内容公告
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 20 AS menu_id, 0 AS pid, '/content' AS path, 'NotificationOutlined' AS icon, '内容' AS menu_name, 'Content' AS route, '' AS params, 55 AS sort, 1 AS is_show, 1 AS is_mer, 1 AS is_menu, 0 AS is_agent
  UNION ALL SELECT 21, 20, '/content/notice', '', '平台公告', 'ContentNotice', '', 54, 1, 1, 1, 0
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 20 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',20,21')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%20%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-loyalty-points', '阶段6：积分商城字段/种子 + 公告'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` IN ('phase6-loyalty-points','phase5-loyalty-points'));
