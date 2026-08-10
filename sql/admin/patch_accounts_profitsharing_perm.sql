-- 平台「分账管理」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 accounts.profitsharing 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21095,540,'accounts.profitsharing.read','查看分账管理','','merchant/applyList','button',1,1),
  (21096,540,'accounts.profitsharing.manage','重新发起分账','','merchant/applyList','button',2,1),
  (21097,540,'accounts.profitsharing.export','导出分账列表','','merchant/applyList','button',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
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
    'accounts',
    'accounts.merchant.dir',
    'accounts.profitsharing',
    'accounts.profitsharing.read',
    'accounts.profitsharing.manage',
    'accounts.profitsharing.export'
  );
