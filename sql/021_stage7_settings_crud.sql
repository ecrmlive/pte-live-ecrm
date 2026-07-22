-- 阶段 7：恢复设置类菜单（最小 CRUD 已接线）
USE `qixi_mergers`;

UPDATE `qixi_system_menu`
SET `is_show` = 1
WHERE `is_mer` = 1
  AND `menu_id` IN (10, 11, 12, 13)
  AND `path` IN ('/setting', '/setting/admin', '/setting/role', '/setting/menu');

UPDATE `qixi_system_menu`
SET `is_show` = 1
WHERE `is_mer` = 2
  AND `menu_id` IN (111, 112)
  AND `path` IN ('/setting/shop', '/setting/staff');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-settings-crud', '阶段7：设置最小 CRUD，恢复菜单可见'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-settings-crud');
