-- 阶段 7：隐藏未实现的设置类占位菜单；DIY 补预约；商户社区 path（可重复执行）
USE `qixi_mergers`;

UPDATE `qixi_system_menu`
SET `is_show` = 0
WHERE `is_mer` = 1
  AND `menu_id` IN (10, 11, 12, 13)
  AND `path` IN ('/setting', '/setting/admin', '/setting/role', '/setting/menu');

UPDATE `qixi_system_menu`
SET `is_show` = 0
WHERE `is_mer` = 2
  AND `menu_id` IN (111, 112)
  AND `path` IN ('/setting/shop', '/setting/staff');

UPDATE `qixi_system_menu`
SET `path` = '/marketing/community'
WHERE `menu_id` = 123 AND `is_mer` = 2 AND `path` = '/community/list';

UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 9, 'name', '预约', 'icon', '', 'url', '/pages/reservation/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '预约', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-hide-settings', '阶段7：隐藏设置占位；商户社区 path；DIY 补预约'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-hide-settings');
