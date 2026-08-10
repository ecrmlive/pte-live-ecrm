-- 平台「等级管理」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 user.member.level 不能用于接口鉴权。
-- 导入后需重新登录平台后台以刷新按钮码缓存。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (530,70,'user.level.dir','用户等级','lucide:medal','/user/member','directory',6,1),
  (183,530,'user.member.level','等级管理','lucide:medal','/user/member/list','page',1,1),
  (21067,183,'user.member.level.read','查看会员等级','','user/member/list','button',1,1),
  (20978,183,'user.member.level.manage','维护会员等级','','user/member/list','button',2,1)
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
    'user.level.dir',
    'user.member.level',
    'user.member.level.read',
    'user.member.level.manage'
  );
