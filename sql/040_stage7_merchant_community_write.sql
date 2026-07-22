-- 阶段 7：商户社区种草写（create/update/delete）；对齐 features CRUD
USE `qixi_mergers`;

-- 商户：逛逛社区（pid=123）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 165, 123, 'community/create', '', '发帖', 'MerCommunityCreateBtn', '', 3, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 165 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 166, 123, 'community/update', '', '编辑帖', 'MerCommunityUpdateBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 166 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 167, 123, 'community/delete', '', '删帖', 'MerCommunityDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 167 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',165,166,167')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%165%';

-- meract / mersub：不赋社区写（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-merchant-community-write', '阶段7：商户社区发帖/编辑/删除按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-merchant-community-write');
