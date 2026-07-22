-- 阶段 6e：直播间竖切（不接微信推流；本店商品挂货）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_broadcast_room` (
  `broadcast_room_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '直播间名',
  `cover_img` varchar(255) NOT NULL DEFAULT '' COMMENT '背景图',
  `feeds_img` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `anchor_name` varchar(64) NOT NULL DEFAULT '' COMMENT '主播昵称',
  `phone` varchar(32) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '审核 0待审 2通过 -1驳回',
  `live_status` smallint(5) NOT NULL DEFAULT 102 COMMENT '101直播中 102未开始 103已结束',
  `is_show` tinyint(3) NOT NULL DEFAULT 1,
  `is_del` tinyint(3) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sort` smallint(5) NOT NULL DEFAULT 0,
  `star` smallint(5) NOT NULL DEFAULT 1 COMMENT '推荐星级',
  `mark` varchar(512) DEFAULT NULL,
  `refusal` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`broadcast_room_id`),
  KEY `idx_mer` (`mer_id`),
  KEY `idx_live` (`live_status`,`status`,`is_show`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='直播间';

CREATE TABLE IF NOT EXISTS `qixi_broadcast_room_goods` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `broadcast_room_id` int(10) unsigned NOT NULL DEFAULT 0,
  `product_id` int(10) unsigned NOT NULL DEFAULT 0,
  `on_sale` tinyint(4) NOT NULL DEFAULT 1,
  `sort` smallint(5) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_room_product` (`broadcast_room_id`,`product_id`),
  KEY `idx_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='直播间挂货';

INSERT INTO `qixi_broadcast_room` (
  `broadcast_room_id`, `mer_id`, `name`, `cover_img`, `feeds_img`, `start_time`, `end_time`,
  `anchor_name`, `phone`, `status`, `live_status`, `is_show`, `sort`, `star`, `mark`
)
SELECT 1, 1, '演示直播间 · 春季上新', '', '',
  DATE_SUB(NOW(), INTERVAL 1 HOUR), DATE_ADD(NOW(), INTERVAL 2 HOUR),
  '小柒主播', '13800000001', 2, 101, 1, 0, 1, '阶段6e直播竖切演示'
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_broadcast_room` WHERE `broadcast_room_id` = 1);

INSERT INTO `qixi_broadcast_room_goods` (`broadcast_room_id`, `product_id`, `on_sale`, `sort`)
SELECT 1, 1, 1, 1
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 1 AND `is_del` = 0)
  AND NOT EXISTS (SELECT 1 FROM `qixi_broadcast_room_goods` WHERE `broadcast_room_id` = 1 AND `product_id` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 27, 17, '/marketing/broadcast', '', '直播监管', 'MarketingBroadcast', '', 66, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 27 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',27')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%27%';

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 122, 115, '/marketing/broadcast', '', '直播间', 'MerMarketingBroadcast', '', 68, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 122 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',122')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%122%';

UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 6, 'name', '直播', 'icon', '', 'url', '/pages/live/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '直播', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-broadcast', '阶段6e：直播间竖切（无微信推流）'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-broadcast');
