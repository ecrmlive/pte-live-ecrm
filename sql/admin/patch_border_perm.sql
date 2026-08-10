-- 平台「活动边框图」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.border.nav 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (9008,60,'marketing.border.nav','活动边框图','lucide:frame','/marketing/border/list','page',10,1),
  (21014,9008,'marketing.border.read','查看活动边框','','marketing/border/list','button',1,1),
  (21015,9008,'marketing.border.manage','维护活动边框','','marketing/border/list','button',2,1)
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
    'marketing.border.nav',
    'marketing.border.read',
    'marketing.border.manage'
  );
