-- 阶段 7：平台按钮权限（退款监管 / 提现审核）+ 演示子账号 auditor
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 34, 30, 'order/refund/approve', '', '同意退款', 'OrderRefundApproveBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 34 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 35, 30, 'order/refund/reject', '', '拒绝退款', 'OrderRefundRejectBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 35 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 36, 32, 'accounts/withdraw/approve', '', '通过提现', 'AccountsWithdrawApproveBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 36 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 37, 32, 'accounts/withdraw/reject', '', '拒绝提现', 'AccountsWithdrawRejectBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 37 AND `is_mer` = 1);

-- 超管角色追加按钮
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',34,35,36,37')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%34%';

-- 演示：平台运营（可退款同意/提现通过，不可拒绝）
INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 4, '平台运营', '1,8,9,30,31,32,34,36', 1, 0, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 4);

INSERT INTO `qixi_system_admin` (`admin_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `status`, `level`, `is_del`)
SELECT 2, 'auditor', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '平台运营', '13800000002', '4', 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_admin` WHERE `account` = 'auditor');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-platform-button-perms', '阶段7：平台退款/提现按钮权限 + auditor'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-platform-button-perms');
