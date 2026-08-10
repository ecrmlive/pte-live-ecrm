-- 平台「转账记录」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 accounts.transfer 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21092,539,'accounts.transfer.read','查看转账记录','','accounts/transferRecord','button',1,1),
  (21093,539,'accounts.transfer.manage','审核备注并登记转账','','accounts/transferRecord','button',2,1),
  (21094,539,'accounts.transfer.export','导出转账记录','','accounts/transferRecord','button',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：页面导航 + 读取 / 管理 / 导出按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'accounts',
    'accounts.merchant.dir',
    'accounts.transfer',
    'accounts.transfer.read',
    'accounts.transfer.manage',
    'accounts.transfer.export'
  );
