-- 阶段 7：商户商品发布/删除 + 平台社区审帖/删帖按钮
USE `qixi_mergers`;

-- 商户：商品列表下
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 134, 103, 'product/create', '', '发布商品', 'MerProductCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 134 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 135, 103, 'product/delete', '', '删除商品', 'MerProductDeleteBtn', '', 3, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 135 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',134,135')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%134%';

-- merprod：可发布，不可删除（上下架仍无）
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',134')
WHERE `role_id` = 5 AND `role_name` = '商品运营' AND `rules` NOT LIKE '%134%';

-- 平台：社区种草下
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 38, 28, 'community/audit', '', '审帖', 'ContentCommunityAuditBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 38 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 39, 28, 'community/delete', '', '删帖', 'ContentCommunityDeleteBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 39 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',38,39')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%38%';

-- auditor：可审帖，不可删帖；补社区菜单入口
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',20,28,38')
WHERE `role_id` = 4 AND `role_name` = '平台运营' AND `rules` NOT LIKE '%38%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-product-community-btns', '阶段7：商品发布/删除 + 社区审帖/删帖按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-product-community-btns');
