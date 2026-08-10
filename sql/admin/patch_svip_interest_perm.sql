-- 平台「会员权益」API 按钮权限（幂等）：read / manage
-- RequireAdminMenu 仅认 kind=button；page 码 user.svip.interest 不能用于接口鉴权。
-- 导入后需重新登录平台后台以刷新按钮码缓存。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (184,73,'user.svip.interest','会员权益','lucide:heart-handshake','/user/member/equity','page',2,1),
  (21069,184,'user.svip.interest.read','查看会员权益','','user/member/equity','button',1,1),
  (20979,184,'user.svip.interest.manage','维护付费会员权益','','user/member/equity','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'user.svip.interest',
    'user.svip.interest.read',
    'user.svip.interest.manage'
  );
