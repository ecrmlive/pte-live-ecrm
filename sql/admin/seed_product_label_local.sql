-- 本地开发：平台商品标签样例（utf8mb4 中文可读假数据）
-- 用法：在 qixi_crm_admin 执行；可重复执行（按名称幂等）

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_label` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT '',
  `color` varchar(32) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_a_product_label` (`name`,`description`,`color`,`sort`,`status`,`created_at`)
SELECT v.name, v.description, '', v.sort, 1, v.created_at
FROM (
  SELECT '全网最低' AS name, '' AS description, 2 AS sort, '2026-01-23 16:06:05' AS created_at UNION ALL
  SELECT '性价比', '', 1, '2026-01-23 16:06:11' UNION ALL
  SELECT '新品', '', 0, '2026-01-12 17:06:52' UNION ALL
  SELECT '包邮', '', 0, '2025-09-22 09:16:13' UNION ALL
  SELECT '三方快递', '', 0, '2025-08-15 09:20:03'
) v
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_product_label` t WHERE t.name = v.name
);

-- 确保菜单与按钮权限存在（与 patch_product_menus / init_data 对齐幂等）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES (44,40,'product.label','商品标签','lucide:tags','/product/label','page',6,1)
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
VALUES (20942,44,'product.label.manage','维护平台商品标签','','product/label','button',1,1)
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
  AND m.id IN (44, 20942);
