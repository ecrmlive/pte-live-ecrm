-- 搜索记录按钮文案对齐 CRMEB：导出搜索记录 / 一键清空
-- 权限码不变；仅更新标题。若按钮文案未刷新可重新登录。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (20975,182,'user.search_record.read','查看用户搜索记录','','user/search-record','button',1,1),
  (20976,182,'user.search_record.clear','一键清空搜索记录','','user/search-record','button',2,1),
  (20977,182,'user.search_record.export','导出搜索记录','','user/search-record','button',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=182,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台超管角色补齐按钮权限（幂等）
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.`id`, m.`id`
FROM `qixi_crm_a_role` r
CROSS JOIN `qixi_crm_a_menu` m
WHERE r.`code`='platform'
  AND m.`code` IN (
    'user.search_record',
    'user.search_record.read',
    'user.search_record.clear',
    'user.search_record.export'
  );
