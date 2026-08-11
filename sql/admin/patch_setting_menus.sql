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
  (1501,1500,'setting.serve','一号通','lucide:cloud','/serve','directory',1,1),
  (1502,1500,'setting.shop.dir','商城设置','lucide:store','/shop','directory',2,1),
  (1503,1500,'setting.delivery.dir','配送配置','lucide:truck','/delivery_config','directory',3,1),
  (1504,1500,'setting.rbac.dir','权限管理','lucide:shield-check','/setting/rbac','directory',4,1),
  (1505,1500,'setting.notice.dir','消息管理','lucide:bell','/notice','directory',5,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1500,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='directory',
  `sort`=VALUES(`sort`),
  `status`=1;

-- L4：一号通（保留入口页即可；短信等可再展开）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1510,1501,'setting.serve.login','登陆入口','lucide:log-in','/setting/sms/sms_config/index','page',1,1),
  (1511,1501,'setting.serve.config','服务配置','lucide:settings-2','/service/settings','page',2,1),
  (1512,1501,'setting.serve.sms','短信设置','lucide:mail','/sms','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1501,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- L4：商城设置
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1520,1502,'setting.shop.form','商城设置','lucide:store','/systemForm/Basics/shop_tabs','page',1,1),
  (1521,1502,'setting.shop.hot','热门搜索','lucide:search','/group/config/67','page',2,1),
  (1522,1502,'setting.shop.agreements','协议规则','lucide:file-text','/setting/agreements','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1502,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

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
    1500,1501,1502,1503,1504,1505,
    1510,1511,1512,
    1520,1521,1522,
    1530,
    1540,1541,1542,
    1550,1551
  );
