-- 平台价格说明（对齐 CRMEB eb_price_rule；分类关联用 cate_ids_json）
-- 可重复执行

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

-- 菜单页（幂等）
UPDATE `qixi_crm_a_menu`
SET `title`='价格说明', `icon`='lucide:badge-dollar-sign', `route_path`='/product/priceDescription',
    `kind`='page', `sort`=8, `parent_id`=40, `status`=1
WHERE `id`=48 OR `code`='product.price_description';

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

-- 维护按钮
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
