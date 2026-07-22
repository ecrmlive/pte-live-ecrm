-- 阶段 7：商户订单按钮权限（is_menu=2）+ 演示角色规则
USE `qixi_mergers`;

-- 挂在「订单列表」下，角色树可勾选；path 供前端 hasPerm
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 127, 106, 'order/deliver', '', '发货', 'MerOrderDeliver', '', 10, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 127 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 128, 106, 'order/verify', '', '核销', 'MerOrderVerify', '', 9, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 128 AND `is_mer` = 2);

-- 商户主账号模板：含发货+核销
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',127,128')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%127%';

-- 演示子角色「商户运营」：追加发货按钮（无核销）；不覆盖后续 025 追加的售后规则
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',127')
WHERE `role_id` = 3 AND `role_name` = '商户运营' AND `rules` NOT LIKE '%127%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-button-perms', '阶段7：商户订单发货/核销按钮权限'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-button-perms');
