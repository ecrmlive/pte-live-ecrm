-- 平台「反馈分类」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 user.feedback.category 不能用于接口鉴权。
-- 导入后需重新登录平台后台以刷新按钮码缓存。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (76,74,'user.feedback.category','反馈分类','lucide:folder-tree','/user/feedback/categories','page',2,1),
  (21064,76,'user.feedback.category.read','查看反馈分类','','user/feedback/categories','button',1,1),
  (20951,76,'user.feedback.category.manage','维护反馈分类','','user/feedback/categories','button',2,1)
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
    'user.feedback',
    'user.feedback.category',
    'user.feedback.category.read',
    'user.feedback.category.manage'
  );
