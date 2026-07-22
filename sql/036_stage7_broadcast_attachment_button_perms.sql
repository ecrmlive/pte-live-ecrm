-- 阶段 7：直播开播/挂货 + 平台/商户素材上传·删除
USE `qixi_mergers`;

-- 商户：直播（pid=122）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 155, 122, 'broadcast/live', '', '开播结束', 'MerBroadcastLiveBtn', '', 3, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 155 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 156, 122, 'broadcast/goods', '', '直播挂货', 'MerBroadcastGoodsBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 156 AND `is_mer` = 2);

-- 商户：素材库（pid=131）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 157, 131, 'attachment/upload', '', '上传素材', 'MerAttachmentUploadBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 157 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 158, 131, 'attachment/delete', '', '删除素材', 'MerAttachmentDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 158 AND `is_mer` = 2);

-- 平台：素材库（pid=33）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 45, 33, 'attachment/upload', '', '上传素材', 'ContentAttachmentUploadBtn', '', 4, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 45 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 46, 33, 'attachment/delete', '', '删除素材', 'ContentAttachmentDeleteBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 46 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',155,156,157,158')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%155%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',45,46')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%45%';

-- meract / auditor：不赋开播挂货与素材写（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-broadcast-attachment-btns', '阶段7：直播开播挂货 + 素材上传删除'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-broadcast-attachment-btns');
