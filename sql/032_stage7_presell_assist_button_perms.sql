-- 阶段 7：预售/助力上下架按钮（店铺券启停见 031）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 139, 121, 'presell/toggle', '', '预售上下架', 'MerPresellToggleBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 139 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 140, 124, 'assist/toggle', '', '助力上下架', 'MerAssistToggleBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 140 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',139,140')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%139%';

-- meract 仅秒杀/券启停，不赋预售/助力上下架（不对称演示）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-presell-assist-btns', '阶段7：预售/助力上下架按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-presell-assist-btns');
