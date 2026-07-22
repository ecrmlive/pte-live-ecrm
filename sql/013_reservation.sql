-- 阶段 6e：预约服务（到店核销）最小竖切
USE `qixi_mergers`;

-- 订单预约字段
SET @db := DATABASE();
SET @has_rd := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'reservation_date'
);
SET @sql := IF(@has_rd = 0,
  'ALTER TABLE `qixi_store_order`
     ADD COLUMN `reservation_date` varchar(20) DEFAULT '''' COMMENT ''预约日期'' AFTER `verify_code`,
     ADD COLUMN `reservation_id` int(10) DEFAULT 0 COMMENT ''预约时段槽位 id'' AFTER `reservation_date`,
     ADD COLUMN `reservation_time_part` varchar(32) NOT NULL DEFAULT '''' COMMENT ''预约时段文案'' AFTER `reservation_id`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `qixi_store_product_reservation` (
  `product_reservation_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) NOT NULL DEFAULT 0,
  `reservation_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1到店 2上门',
  `show_reservation_days` int(11) NOT NULL DEFAULT 7 COMMENT '可约天数',
  `is_cancel_reservation` tinyint(1) NOT NULL DEFAULT 1,
  `time_period` text COMMENT '时段 JSON [{start,end,stock}]',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`product_reservation_id`),
  UNIQUE KEY `uk_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约商品配置';

CREATE TABLE IF NOT EXISTS `qixi_store_product_attr_reservation` (
  `attr_reservation_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(10) NOT NULL DEFAULT 0,
  `unique` char(12) NOT NULL DEFAULT '',
  `start_time` varchar(10) NOT NULL DEFAULT '',
  `end_time` varchar(10) NOT NULL DEFAULT '',
  `stock` int(11) NOT NULL DEFAULT 0 COMMENT '每日可约数',
  `use_num` int(11) NOT NULL DEFAULT 0 COMMENT '累计预约（演示统计）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`attr_reservation_id`),
  KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预约时段槽位';

-- 演示预约商品 type=4（CRMEB DEFINE_TYPE_RESERVATION）
INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `cost`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`, `mer_status`
)
SELECT 11, 1, '到店护理 · 体验预约', '预约到店服务，核销码核销', '预约', 1, 1, 1,
  2, '次', 99.00, 129.00, 20.00, 999, 0, '', '',
  '2', 4, 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 11);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 11, 11, '默认', 'reserve-care', 999, 0, 20.00, 129.00, 99.00, 'sku000000011', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 11);

INSERT INTO `qixi_store_product_reservation` (
  `product_reservation_id`, `product_id`, `reservation_type`, `show_reservation_days`, `is_cancel_reservation`, `time_period`
)
SELECT 1, 11, 1, 7, 1, '[{"start":"09:00","end":"11:00","stock":5},{"start":"14:00","end":"16:00","stock":5},{"start":"16:00","end":"18:00","stock":3}]'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_reservation` WHERE `product_id` = 11);

INSERT INTO `qixi_store_product_attr_reservation` (`attr_reservation_id`, `product_id`, `unique`, `start_time`, `end_time`, `stock`)
SELECT 1, 11, 'rsv000000001', '09:00', '11:00', 5
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_reservation` WHERE `attr_reservation_id` = 1);
INSERT INTO `qixi_store_product_attr_reservation` (`attr_reservation_id`, `product_id`, `unique`, `start_time`, `end_time`, `stock`)
SELECT 2, 11, 'rsv000000002', '14:00', '16:00', 5
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_reservation` WHERE `attr_reservation_id` = 2);
INSERT INTO `qixi_store_product_attr_reservation` (`attr_reservation_id`, `product_id`, `unique`, `start_time`, `end_time`, `stock`)
SELECT 3, 11, 'rsv000000003', '16:00', '18:00', 3
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_reservation` WHERE `attr_reservation_id` = 3);

-- 菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 119, 115, '/marketing/reservation', '', '预约服务', 'MerMarketingReservation', '', 72, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 119 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',119')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%119%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6e-reservation', '阶段6e：预约服务最小竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6e-reservation');
