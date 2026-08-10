-- 平台「报名活动」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.application.nav 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (9217,60,'marketing.application.nav','报名活动','lucide:clipboard-list','/marketing/application/list','page',14,1),
  (21018,9217,'marketing.application.read','查看活动报名','','marketing/application/list','button',1,1),
  (21019,9217,'marketing.application.manage','维护活动报名','','marketing/application/list','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 旧菜单隐藏，避免与 nav 重复
UPDATE `qixi_crm_a_menu`
SET `status`=0, `sort`=1004
WHERE `id`=211 OR `code`='marketing.application';

-- 平台 / 运营：导航 + 读/写按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.application.nav',
    'marketing.application.read',
    'marketing.application.manage'
  );
