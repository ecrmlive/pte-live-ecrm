-- 平台「积分订单」API 按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 marketing.integral.orders 不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21054,9120,'marketing.integral.orders.read','查看积分订单','','marketing/integral/orderList','button',1,1),
  (21055,9120,'marketing.integral.orders.manage','管理积分订单','','marketing/integral/orderList','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：页面导航 + 读取/管理按钮
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.integral.dir',
    'marketing.integral.orders',
    'marketing.integral.orders.read',
    'marketing.integral.orders.manage'
  );
