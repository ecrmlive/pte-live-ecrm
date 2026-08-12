-- 平台「设置」侧栏：系统设置下 L2→L3→L4（截图结构）
-- 旧扁平权限/商城等叶子侧栏隐藏或挂入新树；执行后需重新登录。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='设置',
    `icon`='ant-design:setting-outlined',
    `route_path`='/setting',
    `kind`='directory',
    `sort`=13,
    `parent_id`=0,
    `status`=1
WHERE `id`=120 OR `code`='setting';

-- 旧扁平：侧栏隐藏（功能由系统设置树承接）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=121 OR `code`='setting.admin';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=991 WHERE `id`=122 OR `code`='setting.role';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=992 WHERE `id`=123 OR `code`='setting.menu';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=993 WHERE `id`=124 OR `code`='setting.agreements';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=994 WHERE `id`=125 OR `code`='setting.cloud_config';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=995 WHERE `id`=126 OR `code`='setting.sms';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=996 WHERE `id`=187 OR `code`='setting.operation_log';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=997 WHERE `id`=189 OR `code`='setting.login_log';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=998 WHERE `id`=190 OR `code`='setting.shop';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=999 WHERE `id`=191 OR `code`='setting.pay';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1000 WHERE `id`=196 OR `code`='setting.storage';
-- 历史“快递公司”页面会在父目录隐藏后成为孤儿一级菜单；必须一并隐藏。
-- 正式入口固定为「设置 → 系统设置 → 配送配置 → 物流公司」。
UPDATE `qixi_crm_a_menu`
SET `status`=0, `sort`=1001
WHERE `id` IN (90,91) OR `code` IN ('freight','freight.express');

-- L2：系统设置
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1500,120,'setting.system','系统设置','lucide:settings','/sys','directory',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=120,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='directory',
  `sort`=1,
  `status`=1;

-- L3 目录（截图）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1502,1500,'setting.shop.dir','商城设置','lucide:store','/shop','directory',3,1),
  (1503,1500,'setting.delivery.dir','配送配置','lucide:truck','/delivery_config','directory',4,1),
  (1504,1500,'setting.rbac.dir','权限管理','lucide:shield-check','/setting/rbac','directory',5,1),
  (1505,1500,'setting.notice.dir','消息管理','lucide:bell','/notice','directory',6,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1500,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='directory',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 系统设置仅保留服务配置与平台短信验证码配置。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1511,1500,'setting.system.service','服务配置','lucide:settings-2','/service/settings','page',1,1),
  (1512,1500,'setting.system.sms','短信配置','lucide:mail','/setting/sms','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1500,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 原一号通目录及其呼叫系统登录入口不属于本项目，移除菜单与授权。
DELETE FROM `qixi_crm_a_role_menu` WHERE `menu_id` IN (1501,1510);
DELETE FROM `qixi_crm_a_menu` WHERE `id` IN (1501,1510) OR `code` IN ('setting.serve','setting.serve.login');
UPDATE `qixi_crm_a_menu`
SET `parent_id`=1512, `route_path`='setting/sms', `title`='维护平台短信验证码配置'
WHERE `id`=20936 OR `code`='setting.sms.manage';

-- L4：商城设置
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1520,1502,'setting.shop.form','商城设置','lucide:store','/systemForm/Basics/shop_tabs','page',1,1),
  (1521,1502,'setting.shop.hot','热门搜索','lucide:search','/group/config/67','page',2,1),
  (1522,1502,'setting.shop.agreements','协议设置','lucide:file-text','/setting/agreements','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1502,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 热门搜索按商城设置新菜单授权，旧维护菜单不再作为接口权限来源。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21570,1521,'setting.shop.hot.read','查看热门搜索','','group/config/67','button',1,1),
  (21571,1521,'setting.shop.hot.manage','维护热门搜索','','group/config/67','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1521,
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform' AND m.code IN ('setting.shop.hot','setting.shop.hot.read','setting.shop.hot.manage');

-- L4：配送配置
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1530,1503,'setting.delivery.express','物流公司','lucide:truck','/freight/express','page',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1503,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 物流公司维护权限挂在正式配送配置入口，避免历史按钮成为孤儿或一级菜单。
UPDATE `qixi_crm_a_menu`
SET `parent_id`=1530, `route_path`='freight/express', `title`='维护物流公司', `status`=1
WHERE `id`=20932 OR `code`='freight.express.manage';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform' AND m.code='freight.express.manage';

-- L4：权限管理
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1540,1504,'setting.rbac.role','角色权限','lucide:shield-check','/setting/role','page',1,1),
  (1541,1504,'setting.rbac.admin','管理员管理','lucide:user-round-cog','/setting/admin','page',2,1),
  (1542,1504,'setting.rbac.menu','菜单管理','lucide:menu','/setting/menu','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1504,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- L4：消息管理
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1550,1505,'setting.notice.station','公告管理','lucide:megaphone','/station/notice','page',1,1),
  (1551,1505,'setting.notice.list','消息管理','lucide:bell','/setting/notification/index','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1505,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 按钮挂到新页面节点（保留旧 code）
UPDATE `qixi_crm_a_menu` SET `parent_id`=1522, `route_path`='setting/agreements' WHERE `id`=20913 OR `code`='setting.agreement.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1520, `route_path`='systemForm/Basics/shop_tabs' WHERE `id`=20984 OR `code`='setting.shop.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1540 WHERE `code` LIKE 'setting.role%' AND `kind`='button';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1541 WHERE `code` LIKE 'setting.admin%' AND `kind`='button';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1542 WHERE `code` LIKE 'setting.menu%' AND `kind`='button';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (
    1500,1502,1503,1504,1505,
    1511,1512,
    1520,1521,1522,
    1530,
    1540,1541,1542,
    1550,1551
  );
