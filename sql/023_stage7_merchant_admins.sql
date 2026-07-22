-- 阶段 7：商户后台「子账号」菜单（merchant_admin + 角色绑定）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 126, 110, '/setting/admins', '', '子账号', 'MerSettingAdmins', '', 7, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 126 AND `is_mer` = 2);

-- 若曾以旧 path 插入，对齐路由
UPDATE `qixi_system_menu`
SET `path` = '/setting/admins', `route` = 'MerSettingAdmins', `menu_name` = '子账号'
WHERE `menu_id` = 126 AND `is_mer` = 2 AND `path` <> '/setting/admins';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',126')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%126%';

-- 演示：商户1 本店角色 + 子账号 mersub / admin123
INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 3, '商户运营', '101,105,106,107,108,109', 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 3);

INSERT INTO `qixi_merchant_admin` (`merchant_admin_id`, `mer_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `level`, `status`, `is_del`)
SELECT 3, 1, 'mersub', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '运营子账号', '13900000011', '3', 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_admin` WHERE `account` = 'mersub');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-merchant-admins', '阶段7：商户子账号菜单 + 演示 mersub'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-merchant-admins');
