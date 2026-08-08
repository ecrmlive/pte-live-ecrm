-- 本地：确保「店铺商品参数」菜单与按钮权限存在（幂等，与 patch_product_menus 对齐）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (46,40,'product.parameter','商品参数','lucide:list-tree','/product/specsMain','directory',7,1),
  (56,46,'product.parameter.store','店铺商品参数','lucide:store','/product/merSpecs','page',1,1),
  (57,46,'product.parameter.platform','平台商品参数','lucide:list-tree','/product/specs','page',2,1),
  (20907,56,'product.parameter.store.manage','维护店铺商品参数模板','','product/merSpecs','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (46,56,57,20907);
