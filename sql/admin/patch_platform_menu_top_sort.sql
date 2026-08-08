-- 平台侧栏顶层顺序 + 父子结构 + Iconify ant-design 图标
-- 对齐 CRMEB https://mer.crmeb.net/admin/merchant/index （平台菜单）
-- api-platform / Vben：ORDER BY sort ASC,id ASC，故数值越小越靠前。
-- 目标顶层：首页 → 店铺 → 商品 → 订单 → 分销 → 营销 → 用户 → 内容 → 财务 → 应用 → 装修 → 客服 → 设置 → 维护
-- 店铺子级：店铺管理 → 店铺设置 → 商户管理 → 区域代理
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_admin' < sql/admin/patch_platform_menu_top_sort.sql

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 分销一级目录（可重复执行）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (220,0,'promoter','分销','ant-design:send-outlined','/promoter','directory',5,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`code`=VALUES(`code`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- 客服树（CRMEB：客服自动回复 / 客服列表 / 客服设置）
UPDATE `qixi_crm_a_menu`
SET `kind`='directory', `title`='客服', `icon`='ant-design:customer-service-outlined', `route_path`='/service', `code`='service', `parent_id`=0, `sort`=12
WHERE `id`=30;
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (301,30,'service.auto_reply','客服自动回复','lucide:message-square-reply','/systemForm/customer_keyword','page',1,1),
  (302,30,'service.customer.list','客服列表','ant-design:customer-service-outlined','/service/customer/list','page',2,1),
  (303,30,'service.settings','客服设置','lucide:settings-2','/systemForm/Basics/service','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`code`=VALUES(`code`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- ========== 顶层：标题 / 图标 / sort（ASC）==========
UPDATE `qixi_crm_a_menu` SET `title`='首页', `icon`='ant-design:home-outlined', `parent_id`=0, `sort`=1, `kind`='directory' WHERE `id`=1;
UPDATE `qixi_crm_a_menu` SET `title`='店铺', `icon`='ant-design:shop-outlined', `parent_id`=0, `sort`=2, `kind`='directory', `code`='store' WHERE `id`=10;
UPDATE `qixi_crm_a_menu` SET `title`='商品', `icon`='ant-design:shopping-outlined', `parent_id`=0, `sort`=3, `kind`='directory' WHERE `id`=40;
UPDATE `qixi_crm_a_menu` SET `title`='订单', `icon`='ant-design:file-text-outlined', `parent_id`=0, `sort`=4, `kind`='directory' WHERE `id`=50;
UPDATE `qixi_crm_a_menu` SET `title`='分销', `icon`='ant-design:send-outlined', `parent_id`=0, `sort`=5, `kind`='directory' WHERE `id`=220;
UPDATE `qixi_crm_a_menu` SET `title`='营销', `icon`='ant-design:flag-outlined', `parent_id`=0, `sort`=6, `kind`='directory' WHERE `id`=60;
UPDATE `qixi_crm_a_menu` SET `title`='用户', `icon`='ant-design:user-outlined', `parent_id`=0, `sort`=7, `kind`='directory' WHERE `id`=70;
UPDATE `qixi_crm_a_menu` SET `title`='内容', `icon`='ant-design:read-outlined', `parent_id`=0, `sort`=8, `kind`='directory' WHERE `id`=80;
UPDATE `qixi_crm_a_menu` SET `title`='财务', `icon`='ant-design:bar-chart-outlined', `parent_id`=0, `sort`=9, `kind`='directory' WHERE `id`=100;
UPDATE `qixi_crm_a_menu` SET `title`='应用', `icon`='ant-design:appstore-outlined', `parent_id`=0, `sort`=10, `kind`='directory' WHERE `id`=130;
UPDATE `qixi_crm_a_menu` SET `title`='装修', `icon`='ant-design:format-painter-outlined', `parent_id`=0, `sort`=11, `kind`='directory', `code`='operations' WHERE `id`=110;
UPDATE `qixi_crm_a_menu` SET `title`='客服', `icon`='ant-design:customer-service-outlined', `parent_id`=0, `sort`=12 WHERE `id`=30;
UPDATE `qixi_crm_a_menu` SET `title`='设置', `icon`='ant-design:setting-outlined', `parent_id`=0, `sort`=13, `kind`='directory' WHERE `id`=120;
UPDATE `qixi_crm_a_menu` SET `title`='维护', `icon`='ant-design:tool-outlined', `parent_id`=0, `sort`=14, `kind`='directory' WHERE `id`=197;

-- ========== 店铺子级：店铺管理 / 店铺设置 / 商户管理 / 区域代理 ==========
UPDATE `qixi_crm_a_menu` SET `parent_id`=10, `title`='店铺管理', `icon`='ant-design:shop-outlined', `sort`=1, `kind`='directory' WHERE `id`=18;
UPDATE `qixi_crm_a_menu` SET `parent_id`=10, `title`='店铺设置', `icon`='ant-design:setting-outlined', `sort`=2, `kind`='directory' WHERE `id`=19;
UPDATE `qixi_crm_a_menu` SET `parent_id`=10, `title`='商户管理', `icon`='ant-design:team-outlined', `sort`=3, `kind`='directory', `code`='merchant.mgmt' WHERE `id`=25;
UPDATE `qixi_crm_a_menu` SET `parent_id`=10, `title`='区域代理', `icon`='ant-design:cluster-outlined', `sort`=4, `kind`='directory', `code`='region' WHERE `id`=20;

-- 分销：原「营销 → 分销管理」挂到一级「分销」
UPDATE `qixi_crm_a_menu`
SET `parent_id`=220, `title`='分销管理', `icon`='ant-design:send-outlined', `sort`=1, `code`='marketing.spread'
WHERE `id`=65;

-- 物流配送：非 CRMEB 平台顶层，挂到「设置」下，避免打乱顶层顺序（不删业务页）
UPDATE `qixi_crm_a_menu`
SET `parent_id`=120, `title`='物流配送', `icon`='lucide:map-plus', `sort`=20, `kind`='directory'
WHERE `id`=90;

-- 商品 / 订单关键二级命名对齐 CRMEB
UPDATE `qixi_crm_a_menu` SET `title`='商品管理', `sort`=1 WHERE `id`=43;
UPDATE `qixi_crm_a_menu` SET `title`='商品分类', `sort`=2 WHERE `id`=41;
UPDATE `qixi_crm_a_menu` SET `title`='品牌管理', `sort`=3 WHERE `id`=42;
UPDATE `qixi_crm_a_menu` SET `title`='评论管理', `sort`=4 WHERE `id`=47;
UPDATE `qixi_crm_a_menu` SET `title`='保障服务', `sort`=5 WHERE `id`=45;
UPDATE `qixi_crm_a_menu` SET `title`='商品标签', `sort`=6 WHERE `id`=44;
UPDATE `qixi_crm_a_menu` SET `title`='平台商品参数', `sort`=7 WHERE `id`=46;
UPDATE `qixi_crm_a_menu` SET `title`='价格说明', `sort`=8 WHERE `id`=48;
UPDATE `qixi_crm_a_menu` SET `title`='活动标签', `sort`=9 WHERE `id`=49;

UPDATE `qixi_crm_a_menu` SET `title`='订单列表', `sort`=1 WHERE `id`=51;
UPDATE `qixi_crm_a_menu` SET `title`='退款订单', `sort`=2 WHERE `id`=52;
UPDATE `qixi_crm_a_menu` SET `title`='核销记录', `sort`=3 WHERE `id`=53;

-- 平台角色补授权（分销目录 / 客服子页）
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform' AND m.id IN (220,65,301,302,303);
