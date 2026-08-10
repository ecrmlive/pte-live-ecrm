-- 平台「会员类型」API 按钮权限（幂等）：read / manage
-- RequireAdminMenu 仅认 kind=button；page 码 user.svip.plan 不能用于接口鉴权。
-- 导入后需重新登录平台后台以刷新按钮码缓存。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (180,73,'user.svip.plan','会员类型','lucide:badge-check','/user/member/type','page',1,1),
  (21068,180,'user.svip.plan.read','查看会员类型','','user/member/type','button',1,1),
  (20973,180,'user.svip.plan.manage','维护可售会员类型','','user/member/type','button',2,1)
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
    'user.svip.plan',
    'user.svip.plan.read',
    'user.svip.plan.manage'
  );
