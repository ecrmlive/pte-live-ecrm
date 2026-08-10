-- 平台「优惠套餐」按钮权限（幂等）
-- RequireAdminMenu 仅认 kind=button；导航 page 码不能用于接口鉴权。
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 确保导航页存在且启用
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1629,60,'marketing.discounts.nav','优惠套餐','lucide:package','/marketing/discounts/list','page',12,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21010,1629,'marketing.discounts.read','查看优惠套餐','','marketing/discounts/list','button',1,1),
  (21011,1629,'marketing.discounts.manage','上下架优惠套餐','','marketing/discounts/list','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 平台 / 运营：导航 + read/manage
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'marketing.discounts.nav',
    'marketing.discounts.read',
    'marketing.discounts.manage'
  );
