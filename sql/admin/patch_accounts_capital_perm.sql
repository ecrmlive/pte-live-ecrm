-- 平台「资金记录 / 资金流水」API 按钮权限（幂等）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (20924,102,'accounts.user_assets.read','查看资金记录','','accounts/capital','button',1,1),
  (21102,102,'accounts.user_assets.export','导出资金记录','','accounts/capital','button',2,1),
  (21103,542,'accounts.capital_flow.read','查看资金流水','','accounts/capitalFlow','button',1,1),
  (21104,542,'accounts.capital_flow.export','导出资金流水','','accounts/capitalFlow','button',2,1)
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
    'accounts.user_assets',
    'accounts.user_assets.read',
    'accounts.user_assets.export',
    'accounts.capital_flow',
    'accounts.capital_flow.read',
    'accounts.capital_flow.export'
  );
