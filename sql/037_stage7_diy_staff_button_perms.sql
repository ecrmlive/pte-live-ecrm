-- 阶段 7：平台 DIY 按钮 + 店员发货细权（is_goods）
USE `qixi_mergers`;

-- 平台：DIY（pid=22 页面装修）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 47, 22, 'diy/create', '', '新建装修页', 'ContentDiyCreateBtn', '', 4, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 47 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 48, 22, 'diy/update', '', '编辑装修页', 'ContentDiyUpdateBtn', '', 3, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 48 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 49, 22, 'diy/delete', '', '删除装修页', 'ContentDiyDeleteBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 49 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 50, 22, 'diy/active', '', '启用装修页', 'ContentDiyActiveBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 50 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 51, 22, 'diy/pick', '', '素材选图', 'ContentDiyPickBtn', '', 5, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 51 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',47,48,49,50,51')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%47%';

-- auditor：不赋 DIY 写/选图（不对称）

-- 店员发货：复用 is_goods=1；演示 staff1 开启发货
UPDATE `qixi_store_service`
SET `is_goods` = 1
WHERE `account` = 'staff1' AND `is_del` = 0;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-diy-staff-btns', '阶段7：DIY按钮 + 店员发货细权(is_goods)'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-diy-staff-btns');
