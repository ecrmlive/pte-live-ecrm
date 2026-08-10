-- 平台「活动氛围图」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.atmosphere.nav 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (9007,60,'marketing.atmosphere.nav','活动氛围图','lucide:sparkles','/marketing/atmosphere/list','page',9,1),
  (21012,9007,'marketing.atmosphere.read','查看活动氛围','','marketing/atmosphere/list','button',1,1),
  (21013,9007,'marketing.atmosphere.manage','维护活动氛围','','marketing/atmosphere/list','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：导航 + 读/写按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.atmosphere.nav',
    'marketing.atmosphere.read',
    'marketing.atmosphere.manage'
  );
