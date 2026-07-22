-- 阶段 7：商户售后按钮权限（同意/拒绝退款）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 129, 113, 'order/refund/approve', '', '同意退款', 'MerOrderRefundApproveBtn', '', 2, 1, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 129 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 130, 113, 'order/refund/reject', '', '拒绝退款', 'MerOrderRefundRejectBtn', '', 1, 1, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 130 AND `is_mer` = 2);

-- 商户主账号：售后页 + 全部按钮
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',129,130')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%129%';

-- 演示子账号 mersub（角色3）：售后页 + 仅同意，无拒绝
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',113')
WHERE `role_id` = 3 AND `rules` NOT LIKE '%113%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',129')
WHERE `role_id` = 3 AND `rules` NOT LIKE '%129%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-refund-button-perms', '阶段7：商户同意/拒绝退款按钮权限'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-refund-button-perms');
