-- 平台「平台账单」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 accounts.statement 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21090,538,'accounts.statement.read','查看平台账单','','accounts/statement','button',1,1),
  (21091,538,'accounts.statement.download','下载平台账单','','accounts/statement','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：页面导航 + 读取 / 下载按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'accounts',
    'accounts.merchant.dir',
    'accounts.statement',
    'accounts.statement.read',
    'accounts.statement.download'
  );
