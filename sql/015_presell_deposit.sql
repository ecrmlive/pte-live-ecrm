-- 阶段 6：定金预售（presell_type=2）+ 尾款单
USE `qixi_mergers`;

SET @db := DATABASE();

SET @has_down := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_store_product_presell' AND COLUMN_NAME = 'down_price'
);
SET @sql := IF(@has_down = 0,
  'ALTER TABLE `qixi_store_product_presell`
     ADD COLUMN `down_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT ''定金'' AFTER `price`,
     ADD COLUMN `final_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT ''尾款'' AFTER `down_price`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `qixi_presell_order` (
  `presell_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '预售尾款订单id',
  `presell_order_sn` varchar(32) NOT NULL COMMENT '预售订单号',
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `order_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '主商品订单id',
  `product_presell_id` int(10) unsigned NOT NULL DEFAULT 0,
  `final_start_time` datetime NOT NULL,
  `final_end_time` datetime NOT NULL,
  `paid` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT '0未支付 1已支付',
  `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '0无效 1有效',
  `pay_type` tinyint(3) NOT NULL DEFAULT 0,
  `pay_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '尾款',
  `refun_price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00,
  `pay_time` datetime DEFAULT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`presell_order_id`),
  KEY `idx_order` (`order_id`),
  KEY `idx_uid_paid` (`uid`,`paid`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='预售尾款单';

-- 演示商品 12：定金预售 ¥99（定金 20 + 尾款 79）
INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `cost`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`, `mer_status`
)
SELECT 12, 1, '定金预售 · 演示套装', '定金预售竖切：先付定金再付尾款', '预售,定金', 1, 1, 1,
  2, '套', 99.00, 129.00, 30.00, 999, 0, '', '',
  '2', 0, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 12);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 12, 12, '默认', 'deposit-demo', 999, 0, 30.00, 129.00, 99.00, 'sku000000012', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 12);

INSERT INTO `qixi_store_product_presell` (
  `product_presell_id`, `start_time`, `end_time`, `final_start_time`, `final_end_time`,
  `status`, `presell_type`, `pay_count`, `delivery_type`, `delivery_day`,
  `product_id`, `price`, `down_price`, `final_price`, `stock`, `is_show`, `store_name`,
  `mer_id`, `store_info`, `product_status`, `action_status`
)
SELECT 2, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 30 DAY),
  DATE_FORMAT(NOW(), '%Y-%m-%d %H:%i:%s'), DATE_FORMAT(DATE_ADD(NOW(), INTERVAL 14 DAY), '%Y-%m-%d %H:%i:%s'),
  1, 2, 0, 1, 7,
  12, 99.00, 20.00, 79.00, 50, 1, '定金预售 · 演示套装',
  1, '阶段6定金预售竖切', 1, 1
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 12 AND `is_del` = 0)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_presell` WHERE `product_presell_id` = 2);

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-presell-deposit', '阶段6：定金预售+尾款单'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-presell-deposit');
