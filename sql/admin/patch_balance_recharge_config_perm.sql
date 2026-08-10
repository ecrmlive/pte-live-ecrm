-- 平台「余额充值设置」档位列表 API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.balance.config 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 导航页：标题对齐 CRMEB「余额充值设置」
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (687,5126,'marketing.balance.config','余额充值设置','lucide:badge-dollar-sign','/group/config/69','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=5126,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21058,687,'marketing.balance.config.read','查看余额充值设置','','group/config/69','button',1,1),
  (21059,687,'marketing.balance.config.manage','管理余额充值设置','','group/config/69','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 旧按钮保留并挂到本页，避免历史角色绑定悬空
UPDATE `qixi_crm_a_menu`
SET `parent_id`=687,
    `route_path`='group/config/69',
    `title`='维护充值计划并查看充值订单',
    `status`=1
WHERE `id`=20971 OR `code`='marketing.recharge.manage';

-- 平台 / 运营：导航 + read/manage；旧 manage 一并绑定
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.balance.dir',
    'marketing.balance.config',
    'marketing.balance.config.read',
    'marketing.balance.config.manage',
    'marketing.recharge.manage'
  );
