-- 平台「提现管理」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 accounts.withdraw 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21098,101,'accounts.withdraw.read','查看提现管理','','accounts/extract','button',1,1),
  (20923,101,'accounts.withdraw.review','审核用户提现','','accounts/extract','button',2,1),
  (21099,101,'accounts.withdraw.export','导出提现列表','','accounts/extract','button',3,1)
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
    'accounts.user.dir',
    'accounts.withdraw',
    'accounts.withdraw.read',
    'accounts.withdraw.review',
    'accounts.withdraw.export'
  );
