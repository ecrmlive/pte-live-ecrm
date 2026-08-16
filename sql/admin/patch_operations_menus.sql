-- 平台「装修」侧栏：叶子与「页面链接」目录同级
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='装修',
    `icon`='ant-design:format-painter-outlined',
    `route_path`='/operations',
    `kind`='directory',
    `sort`=11,
    `parent_id`=0,
    `status`=1
WHERE `id`=110 OR `code`='operations';

-- 二级：叶子 + 页面链接目录
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (111,110,'operations.diy','首页装修','lucide:house','/setting/diy/list','page',1,1),
  (806,110,'operations.product_category','分类装修','lucide:folder-tree','/setting/product_category','page',2,1),
  (811,110,'operations.cart_diy','购物车装修','lucide:shopping-cart','/setting/diy/cart','page',3,1),
  (812,110,'operations.member_diy','我的装修','lucide:circle-user-round','/setting/diy/personal','page',4,1),
  (801,110,'operations.product_detail','详情装修','lucide:package','/setting/diy/product_detail','page',5,1),
  (800,110,'operations.merchant_diy','店铺装修','lucide:store','/setting/merchant/diyList','page',6,1),
  (802,110,'operations.page_config','页面配置','lucide:sliders-horizontal','/setting/system_visualization_data','page',7,1),
  (212,110,'operations.system_form','表单配置','lucide:clipboard-pen','/systemForm/form_list','page',8,1),
  (805,110,'operations.fab','悬浮菜单','lucide:circle-dot','/setting/fab','page',9,1),
  (803,110,'operations.page_links','页面链接','lucide:link','/setting/page','directory',10,1),
  (804,110,'operations.material','素材管理','lucide:images','/config/picture','page',11,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=110,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：页面链接
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (807,803,'operations.page_links.platform_cat','平台页面分类','lucide:folder-tree','/setting/diy/plantform/category/list','page',1,1),
  (808,803,'operations.page_links.platform','平台页面链接','lucide:link','/setting/diy/links/list','page',2,1),
  (809,803,'operations.page_links.merchant_cat','商户页面分类','lucide:folder-tree','/setting/diy/merchant/category/list','page',3,1),
  (810,803,'operations.page_links.merchant','商户页面链接','lucide:link','/setting/diy/merLink/list','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=803,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `parent_id`=111, `route_path`='setting/diy/list' WHERE `id`=20914 OR `code`='operations.diy.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=212, `route_path`='systemForm/form_list' WHERE `id`=21020 OR `code`='operations.system_form.manage';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.id IN (110,111,212,800,801,802,803,804,805,806,807,808,809,810,811,812);
