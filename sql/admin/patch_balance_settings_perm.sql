-- 平台「余额设置」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.balance.settings 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 确保导航页存在且启用
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (667,5126,'marketing.balance.settings','余额设置','lucide:sliders-horizontal','/systemForm/Basics/balance','page',1,1)
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
  (21056,667,'marketing.balance.settings.read','查看余额设置','','systemForm/Basics/balance','button',1,1),
  (21057,667,'marketing.balance.settings.manage','保存余额设置','','systemForm/Basics/balance','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 充值计划/订单按钮挂到「余额充值配置」页，避免与余额设置混淆
UPDATE `qixi_crm_a_menu`
SET `parent_id`=687,
    `route_path`='group/config/69',
    `title`='维护充值计划并查看充值订单'
WHERE `id`=20971 OR `code`='marketing.recharge.manage';

-- 平台 / 运营：导航 + read/manage；充值配置按钮一并绑定
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.balance.dir',
    'marketing.balance.settings',
    'marketing.balance.settings.read',
    'marketing.balance.settings.manage',
    'marketing.balance.config',
    'marketing.recharge.manage'
  );
