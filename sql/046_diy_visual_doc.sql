-- DIY 可视化文档协议：页面外观字段 + 默认首页升级为 {page,items[]}
USE `qixi_mergers`;

SET @db := DATABASE();

SET @has_bg := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_diy' AND COLUMN_NAME = 'is_bg_color'
);
SET @sql := IF(@has_bg = 0,
  'ALTER TABLE `qixi_diy`
    ADD COLUMN `is_bg_color` tinyint(1) NOT NULL DEFAULT 0 COMMENT ''是否背景色'' AFTER `is_diy`,
    ADD COLUMN `is_bg_pic` tinyint(1) NOT NULL DEFAULT 0 COMMENT ''是否背景图'' AFTER `is_bg_color`,
    ADD COLUMN `color_picker` varchar(50) NOT NULL DEFAULT '''' COMMENT ''背景色'' AFTER `is_bg_pic`,
    ADD COLUMN `bg_pic` varchar(255) NOT NULL DEFAULT '''' COMMENT ''背景图'' AFTER `color_picker`,
    ADD COLUMN `bg_tab_val` tinyint(1) NOT NULL DEFAULT 0 COMMENT ''背景图模式'' AFTER `bg_pic`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 将旧 banners/menus 种子升级为可视化文档（仅当 value 仍是旧协议时）
UPDATE `qixi_diy`
SET `value` = '{"page":{"type":"page","name":"页面设置","params":{"name":"平台首页","title":"栖息商城","share_title":"栖息商城"},"style":{"titleTextColor":"black","titleBackgroundColor":"#ffffff"}},"items":[{"type":"banner","name":"轮播图","params":{},"style":{"btnColor":"#ffffff","btnShape":"round","indicator":"1"},"data":[{"imgUrl":"","linkUrl":"/pages/seckill/list","imgName":"夏日秒杀"},{"imgUrl":"","linkUrl":"/pages/points/list","imgName":"积分好物"}]},{"type":"navBar","name":"导航组","params":{},"style":{"background":"#ffffff","rowsNum":"4","show_title":"1"},"data":[{"text":"秒杀","imgUrl":"","linkUrl":"/pages/seckill/list"},{"text":"积分","imgUrl":"","linkUrl":"/pages/points/list"},{"text":"领券","imgUrl":"","linkUrl":"/pages/coupon/center"}]}]}',
    `version` = '2.0'
WHERE `id` = 1
  AND (`value` LIKE '%"banners"%' OR `value` NOT LIKE '%"items"%');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'diy-046-visual-doc', 'DIY 可视化文档协议与外观字段'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'diy-046-visual-doc');
