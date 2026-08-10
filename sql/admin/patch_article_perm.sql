-- 平台「文章管理」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 content.article 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (84,534,'content.article','文章管理','lucide:pen-line','/cms/article','page',1,1),
  (21073,84,'content.article.read','查看文章','','cms/article','button',1,1),
  (20910,84,'content.article.manage','维护文章','','cms/article','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：页面 + 读/写按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'content.article',
    'content.article.read',
    'content.article.manage'
  );
