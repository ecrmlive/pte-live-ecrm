-- 阶段 7：商户「角色权限」菜单（RBAC 树勾选）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 125, 110, '/setting/role', '', '角色权限', 'MerSettingRole', '', 6, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 125 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',125')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%125%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-rbac-menu', '阶段7：商户角色权限菜单'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-rbac-menu');
