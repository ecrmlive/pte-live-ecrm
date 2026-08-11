-- 平台「充值记录」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 accounts.recharge_record 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21100,541,'accounts.recharge_record.read','查看充值记录','','accounts/bill','button',1,1),
  (21101,541,'accounts.recharge_record.refund','充值退款','','accounts/bill','button',2,1)
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
    'accounts.recharge_record',
    'accounts.recharge_record.read',
    'accounts.recharge_record.refund'
  );
