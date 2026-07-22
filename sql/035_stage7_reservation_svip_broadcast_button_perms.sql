-- 阶段 7：预约配置 / SVIP（商户叠加 + 平台设会员）/ 直播建房·删除
USE `qixi_mergers`;

-- 商户：直播创建/删除（pid=122 直播间）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 151, 122, 'broadcast/create', '', '新建直播间', 'MerBroadcastCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 151 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 152, 122, 'broadcast/delete', '', '删除直播间', 'MerBroadcastDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 152 AND `is_mer` = 2);

-- 商户：预约时段配置（pid=119）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 153, 119, 'reservation/config', '', '预约时段配置', 'MerReservationConfigBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 153 AND `is_mer` = 2);

-- 商户：SVIP 叠加设置（pid=120）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 154, 120, 'svip/update', '', '会员价叠加设置', 'MerSvipUpdateBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 154 AND `is_mer` = 2);

-- 平台：设置用户付费会员（pid=25 付费会员）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 44, 25, 'svip/update', '', '设置付费会员', 'UserSvipUpdateBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 44 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',151,152,153,154')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%151%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',44')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%44%';

-- meract：无直播/预约/SVIP 写按钮（不对称）
-- auditor：可看付费会员列表菜单则另赋；默认不赋设置会员（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-reservation-svip-broadcast-btns', '阶段7：预约配置/SVIP设置/直播创建删除按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-reservation-svip-broadcast-btns');
