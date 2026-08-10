-- 平台「专场列表」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.topic.nav 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1470,60,'marketing.topic.nav','专场列表','lucide:layout-template','/group/topic/94','page',11,1),
  (21016,1470,'marketing.topic.read','查看专场','','group/topic/94','button',1,1),
  (21017,1470,'marketing.topic.manage','维护专场','','group/topic/94','button',2,1)
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
    'marketing.topic.nav',
    'marketing.topic.read',
    'marketing.topic.manage'
  );
