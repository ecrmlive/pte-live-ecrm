-- 阶段 7：菜单接线修复（可重复执行）
-- 1) 平台退款/财务/提现（原 005 与品牌 menu_id=14 冲突导致未插入）
-- 2) 社区种草挂到「内容」父级
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 30, 8, '/order/refund', '', '退款监管', 'OrderRefund', '', 68, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/order/refund' AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 31, 0, '/accounts', 'AccountBookOutlined', '财务', 'Accounts', '', 60, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/accounts' AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 32, 31, '/accounts/withdraw', '', '提现审核', 'AccountsWithdraw', '', 59, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/accounts/withdraw' AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 113, 105, '/order/refund', '', '售后处理', 'MerOrderRefund', '', 77, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/order/refund' AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 114, 108, '/finance/withdraw', '', '提现申请', 'MerFinanceWithdraw', '', 68, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/finance/withdraw' AND `is_mer` = 2);

UPDATE `qixi_system_menu`
SET `pid` = 20
WHERE `menu_id` = 28 AND `is_mer` = 1 AND `path` = '/content/community' AND `pid` <> 20;

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',30,31,32')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%30%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',113,114')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%113%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-menu-fix', '阶段7：补退款/提现菜单；社区挂内容父级'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-menu-fix');
