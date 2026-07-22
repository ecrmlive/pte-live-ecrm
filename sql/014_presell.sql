-- 阶段 6：全款预售（presell_type=1；定金+尾款顺延）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_product_presell` (
  `product_presell_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `start_time` datetime NOT NULL COMMENT '预售开始',
  `end_time` datetime NOT NULL COMMENT '预售结束',
  `final_start_time` varchar(30) DEFAULT '' COMMENT '尾款开始(定金模式预留)',
  `final_end_time` varchar(30) DEFAULT '' COMMENT '尾款结束(定金模式预留)',
  `status` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '平台控制 1开启',
  `presell_type` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '1全款 2定金',
  `pay_count` int(10) unsigned DEFAULT 0 COMMENT '限购 0不限',
  `delivery_type` int(10) unsigned NOT NULL DEFAULT 1 COMMENT '1支付后发货 2预售结束后',
  `delivery_day` int(10) unsigned DEFAULT 0 COMMENT '发货天数提示',
  `product_id` int(10) unsigned NOT NULL DEFAULT 0,
  `price` decimal(10,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '预售价',
  `stock` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '活动库存',
  `is_show` tinyint(3) unsigned DEFAULT 1 COMMENT '商户上下架',
  `store_name` varchar(128) NOT NULL DEFAULT '' COMMENT '活动标题',
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `store_info` varchar(255) DEFAULT NULL,
  `is_del` int(10) unsigned NOT NULL DEFAULT 0,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `product_status` int(11) DEFAULT 1 COMMENT '审核 1通过',
  `refusal` varchar(255) DEFAULT NULL,
  `action_status` int(11) DEFAULT 1 COMMENT '1进行中 -1结束',
  `seles` int(10) unsigned DEFAULT 0 COMMENT '销量',
  PRIMARY KEY (`product_presell_id`),
  KEY `idx_time` (`start_time`,`end_time`),
  KEY `idx_product` (`product_id`),
  KEY `idx_mer` (`mer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品预售活动';

-- 演示：商品3 全款预售 ¥15.90，库存 100
INSERT INTO `qixi_store_product_presell` (
  `product_presell_id`, `start_time`, `end_time`, `status`, `presell_type`, `pay_count`,
  `delivery_type`, `delivery_day`, `product_id`, `price`, `stock`, `is_show`, `store_name`,
  `mer_id`, `store_info`, `product_status`, `action_status`, `seles`
)
SELECT 1, DATE_SUB(NOW(), INTERVAL 1 DAY), DATE_ADD(NOW(), INTERVAL 30 DAY), 1, 1, 0,
  1, 3, 3, 15.90, 100, 1, '全款预售 · Type-C 数据线',
  2, '阶段6全款预售竖切演示', 1, 1, 0
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 3 AND `is_del` = 0)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_presell` WHERE `product_presell_id` = 1);

-- 平台菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 26, 17, '/marketing/presell', '', '预售监管', 'MarketingPresell', '', 68, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 26 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',26')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%26%';

-- 商户菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 121, 115, '/marketing/presell', '', '预售活动', 'MerMarketingPresell', '', 70, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 121 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',121')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%121%';

-- DIY 首页入口
UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 5, 'name', '预售', 'icon', '', 'url', '/pages/presell/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '预售', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-presell', '阶段6：全款预售竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-presell');
