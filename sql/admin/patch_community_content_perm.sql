-- 平台「社区内容 / 社区评论」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码不能用于接口鉴权。
-- 导入后需重新登录平台后台以刷新按钮码缓存。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (87,82,'content.community.list','社区内容','lucide:images','/community/list','page',3,1),
  (88,82,'content.community.reply','社区评论','lucide:message-square','/community/reply','page',4,1),
  (21082,87,'content.community_list.read','查看社区内容','','community/list','button',1,1),
  (21083,87,'content.community_list.manage','维护社区内容','','community/list','button',2,1),
  (21084,88,'content.community_reply.read','查看社区评论','','community/reply','button',1,1),
  (21085,88,'content.community_reply.manage','维护社区评论','','community/reply','button',2,1),
  (20920,87,'content.community.audit','审核社区内容','','community/list','button',3,1),
  (20921,87,'content.community.delete','删除社区内容','','community/list','button',4,1)
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
    'content',
    'content.community',
    'content.community.list',
    'content.community.reply',
    'content.community_list.read',
    'content.community_list.manage',
    'content.community_reply.read',
    'content.community_reply.manage',
    'content.community.audit',
    'content.community.delete'
  );
