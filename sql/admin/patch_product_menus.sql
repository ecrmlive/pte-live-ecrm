-- 平台「商品」侧栏对齐图片1 / CRMEB 嵌套结构
-- 品牌管理、商品参数改为二级目录；活动标签软隐藏（页面与路由保留）
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_admin' < sql/admin/patch_product_menus.sql
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 一级「商品」
UPDATE `qixi_crm_a_menu`
SET `title`='商品',
    `icon`='ant-design:shopping-outlined',
    `route_path`='/product',
    `kind`='directory',
    `sort`=3,
    `parent_id`=0,
    `status`=1
WHERE `id`=40 OR `code`='product';

-- 直属叶子
UPDATE `qixi_crm_a_menu`
SET `title`='商品管理', `icon`='lucide:shield-check', `route_path`='/product/audit', `kind`='page', `sort`=1, `parent_id`=40, `status`=1
WHERE `id`=43 OR `code`='product.audit';

UPDATE `qixi_crm_a_menu`
SET `title`='商品分类', `icon`='lucide:folder-tree', `route_path`='/product/category', `kind`='page', `sort`=2, `parent_id`=40, `status`=1
WHERE `id`=41 OR `code`='product.category';

-- 品牌管理（目录）+ 子页
UPDATE `qixi_crm_a_menu`
SET `title`='品牌管理', `icon`='lucide:award', `route_path`='/product/brand', `kind`='directory', `sort`=3, `parent_id`=40, `code`='product.brand', `status`=1
WHERE `id`=42 OR `code`='product.brand';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (54,42,'product.brand.category','品牌分类','lucide:folder-tree','/product/band/brandClassify','page',1,1),
  (55,42,'product.brand.list','品牌列表','lucide:award','/product/band/brandList','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu`
SET `title`='评论管理', `icon`='lucide:message-square', `route_path`='/product/comment', `kind`='page', `sort`=4, `parent_id`=40, `status`=1
WHERE `id`=47 OR `code`='product.comment';

UPDATE `qixi_crm_a_menu`
SET `title`='保障服务', `icon`='lucide:shield', `route_path`='/product/guarantee', `kind`='page', `sort`=5, `parent_id`=40, `status`=1
WHERE `id`=45 OR `code`='product.guarantee';

UPDATE `qixi_crm_a_menu`
SET `title`='商品标签', `icon`='lucide:tags', `route_path`='/product/label', `kind`='page', `sort`=6, `parent_id`=40, `status`=1
WHERE `id`=44 OR `code`='product.label';

-- 商品参数（目录）+ 子页
UPDATE `qixi_crm_a_menu`
SET `title`='商品参数', `icon`='lucide:list-tree', `route_path`='/product/specsMain', `kind`='directory', `sort`=7, `parent_id`=40, `code`='product.parameter', `status`=1
WHERE `id`=46 OR `code`='product.parameter';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (56,46,'product.parameter.store','店铺商品参数','lucide:store','/product/merSpecs','page',1,1),
  (57,46,'product.parameter.platform','平台商品参数','lucide:list-tree','/product/specs','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu`
SET `title`='价格说明', `icon`='lucide:badge-dollar-sign', `route_path`='/product/priceDescription', `kind`='page', `sort`=8, `parent_id`=40, `status`=1
WHERE `id`=48 OR `code`='product.price_description';

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

-- 活动标签：功能页保留，侧栏不展示（对齐图片1）
UPDATE `qixi_crm_a_menu`
SET `title`='活动标签', `icon`='lucide:badge', `route_path`='/product/activityLabel', `kind`='page', `sort`=99, `parent_id`=40, `status`=0
WHERE `id`=49 OR `code`='product.activity_label';

-- 按钮权限挂到实际页面节点
UPDATE `qixi_crm_a_menu`
SET `parent_id`=55, `title`='维护商品品牌', `route_path`='product/band/brandList', `kind`='button', `sort`=1, `status`=1
WHERE `id`=20905 OR `code`='product.brand.manage';

UPDATE `qixi_crm_a_menu`
SET `parent_id`=57, `title`='维护平台商品参数模板', `route_path`='product/specs', `kind`='button', `sort`=1, `status`=1
WHERE `id`=20944 OR `code`='product.parameter.manage';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (20906,54,'product.brand.category.manage','维护品牌分类','','product/band/brandClassify','button',1,1),
  (20907,56,'product.parameter.store.manage','维护店铺商品参数模板','','product/merSpecs','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台角色补授权新菜单/按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (48,54,55,56,57,20906,20907,21021);
