-- 阶段 7：素材库竖切（分类 + 附件；本地上传）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_system_attachment` (
  `attachment_id` int(11) NOT NULL AUTO_INCREMENT,
  `attachment_category_id` int(10) unsigned NOT NULL DEFAULT 0,
  `attachment_name` varchar(100) NOT NULL DEFAULT '',
  `attachment_src` varchar(255) NOT NULL DEFAULT '',
  `upload_type` tinyint(3) unsigned NOT NULL DEFAULT 1 COMMENT '1本地',
  `user_type` int(11) NOT NULL DEFAULT 0 COMMENT '0平台 >0商户 mer_id -1用户',
  `user_id` int(10) unsigned NOT NULL DEFAULT 0,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `attachment_type` tinyint(1) unsigned NOT NULL DEFAULT 0 COMMENT '0图片',
  PRIMARY KEY (`attachment_id`),
  KEY `idx_cate` (`attachment_category_id`),
  KEY `idx_owner` (`user_type`,`user_id`,`upload_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='附件/素材';

CREATE TABLE IF NOT EXISTS `qixi_system_attachment_category` (
  `attachment_category_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(10) unsigned NOT NULL DEFAULT 0,
  `path` varchar(512) NOT NULL DEFAULT '',
  `attachment_category_name` varchar(32) NOT NULL DEFAULT '',
  `attachment_category_enname` varchar(16) NOT NULL DEFAULT '',
  `sort` smallint(5) unsigned NOT NULL DEFAULT 0,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '0平台',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`attachment_category_id`),
  KEY `idx_mer` (`mer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='素材分类';

INSERT INTO `qixi_system_attachment_category` (
  `attachment_category_id`, `pid`, `path`, `attachment_category_name`, `attachment_category_enname`, `sort`, `mer_id`
)
SELECT 1, 0, '', '默认分类', 'default', 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_attachment_category` WHERE `attachment_category_id` = 1);

-- 商户默认分类（mer_id=1）
INSERT INTO `qixi_system_attachment_category` (
  `attachment_category_id`, `pid`, `path`, `attachment_category_name`, `attachment_category_enname`, `sort`, `mer_id`
)
SELECT 2, 0, '', '店铺素材', 'store', 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_attachment_category` WHERE `attachment_category_id` = 2);

-- 平台：内容 → 素材库
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 33, 20, '/content/attachment', '', '素材库', 'ContentAttachment', '', 52, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 33 AND `is_mer` = 1);

UPDATE `qixi_system_role` SET `rules` = CONCAT(`rules`, ',33')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%33%';

-- 商户：设置 → 素材库（131，避开 125 角色 / 126 子账号 / 127–130 按钮）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 131, 110, '/setting/attachment', '', '素材库', 'MerSettingAttachment', '', 5, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 131 AND `is_mer` = 2);

UPDATE `qixi_system_role` SET `rules` = CONCAT(`rules`, ',131')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%131%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-attachment', '阶段7：素材库分类/附件 + 菜单'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-attachment');
