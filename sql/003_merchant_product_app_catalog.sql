-- 阶段 2（商户商品 + C 端可售）：SKU 表、可售种子、schema 版本
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_product_attr_value` (
  `value_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `detail` varchar(1000) NOT NULL DEFAULT '' COMMENT '规格详情 JSON/文案',
  `sku` varchar(128) NOT NULL DEFAULT '' COMMENT '规格索引',
  `stock` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '库存',
  `sales` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `image` varchar(128) DEFAULT NULL COMMENT '图片',
  `bar_code` varchar(50) NOT NULL DEFAULT '' COMMENT '条码',
  `cost` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '成本价',
  `ot_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '原价',
  `price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '价格',
  `volume` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '体积',
  `weight` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '重量',
  `type` tinyint(1) DEFAULT '0' COMMENT '0商品',
  `extension_one` decimal(8,2) DEFAULT '0.00',
  `extension_two` decimal(8,2) DEFAULT '0.00',
  `unique` char(12) NOT NULL DEFAULT '' COMMENT '唯一值',
  `svip_price` decimal(10,2) NOT NULL DEFAULT '0.00',
  `library_id` int(11) DEFAULT '0',
  `bar_code_number` varchar(50) NOT NULL DEFAULT '',
  `is_default_select` tinyint(1) DEFAULT '0',
  `is_show` tinyint(1) DEFAULT '1',
  `settlement_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  PRIMARY KEY (`value_id`) USING BTREE,
  KEY `idx_product` (`product_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品 SKU';

-- 演示商户：一单规格已审可售 + 一多规格待审
UPDATE `qixi_store_product`
SET `status` = 1, `is_show` = 1, `store_name` = '高山冻干草莓脆 100g', `store_info` = '阶段2可售演示',
    `price` = 29.90, `ot_price` = 39.90, `stock` = 500, `spec_type` = 0, `cate_id` = 2
WHERE `product_id` = 1 AND `mer_id` = 1;

INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `cost`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`
)
SELECT 2, 1, '纯棉圆领短袖 T 恤', '多规格演示（待审）', '服饰', 1, 1, 0,
  2, '件', 69.00, 99.00, 30.00, 200, 1, '', '',
  '2', 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 2);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 1, 1, '默认', 'default', 500, 0, 12.00, 39.90, 29.90, 'sku000000001', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 1);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 2, 2, '白色|M', '白|M', 80, 0, 30.00, 99.00, 69.00, 'sku000000002', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 2);

INSERT INTO `qixi_store_product_attr_value` (
  `value_id`, `product_id`, `detail`, `sku`, `stock`, `sales`, `cost`, `ot_price`, `price`, `unique`, `is_show`
)
SELECT 3, 2, '黑色|L', '黑|L', 120, 0, 30.00, 99.00, 72.00, 'sku000000003', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_value` WHERE `value_id` = 3);

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase2-merchant-product-app', '阶段2：商户商品SKU + C端可售目录'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase2-merchant-product-app');
