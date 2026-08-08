-- 本地开发：平台价格说明样例（utf8mb4 中文可读假数据）
-- 用法：在 qixi_crm_admin 执行；可重复执行（按名称幂等）

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_price_rule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL COMMENT '名称',
  `cate_ids_json` json NULL COMMENT '关联平台分类 ID 列表；空=全部商品',
  `is_default` tinyint NOT NULL DEFAULT 1 COMMENT '1=未选分类默认全部商品',
  `content` mediumtext NULL COMMENT '价格说明详情 HTML',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示 1显示 0隐藏',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_sort` (`status`,`sort`,`id`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_a_product_price_rule`
  (`name`,`cate_ids_json`,`is_default`,`content`,`sort`,`status`,`created_at`,`updated_at`)
SELECT v.name, v.cate_ids_json, v.is_default, v.content, v.sort, 1, v.created_at, v.updated_at
FROM (
  SELECT
    '集成灶' AS name,
    JSON_ARRAY(7635,7628,7631,7632,7621) AS cate_ids_json,
    0 AS is_default,
    '<p><strong>131321313213131</strong></p>' AS content,
    1 AS sort,
    '2023-08-07 14:33:12' AS created_at,
    '2023-08-07 14:33:12' AS updated_at
  UNION ALL
  SELECT
    'MOCO针织开衫假两件',
    JSON_ARRAY(7635,7628,7631,7632,7621),
    0,
    '<p>演示价格说明：含优惠、规格与运费差异提示，避免售后纠纷。</p>',
    0,
    '2023-08-07 14:31:27',
    '2023-08-07 14:31:27'
  UNION ALL
  SELECT
    '通用价格说明',
    JSON_ARRAY(),
    1,
    '<p>未指定分类时默认适用于全部商品。演示文案：售价以结算页为准。</p>',
    0,
    '2023-08-04 10:00:00',
    '2023-08-04 10:00:00'
  UNION ALL
  SELECT
    '服饰专场说明',
    JSON_ARRAY(7631),
    0,
    '<p>服饰类演示：尺码色差与运费说明。</p>',
    2,
    '2023-08-03 09:20:00',
    '2023-08-03 09:20:00'
) v
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_product_price_rule` t WHERE t.name = v.name
);

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES (48,40,'product.price_description','价格说明','lucide:badge-dollar-sign','/product/priceDescription','page',8,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES (21021,48,'product.price_description.manage','维护平台价格说明','','product/priceDescription','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (48, 21021);
