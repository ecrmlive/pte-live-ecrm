-- 阶段 6e：社区种草竖切（图文+评论+平台审核；可挂商品）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_community_category` (
  `category_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cate_name` varchar(50) NOT NULL DEFAULT '',
  `pid` int(11) NOT NULL DEFAULT 0,
  `is_show` tinyint(4) NOT NULL DEFAULT 1,
  `sort` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区分类';

CREATE TABLE IF NOT EXISTS `qixi_community_topic` (
  `topic_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `topic_name` varchar(100) NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT 1,
  `is_hot` tinyint(4) NOT NULL DEFAULT 0,
  `category_id` int(10) unsigned NOT NULL DEFAULT 0,
  `is_del` tinyint(4) NOT NULL DEFAULT 0,
  `count_use` int(10) unsigned NOT NULL DEFAULT 0,
  `sort` int(10) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`topic_id`),
  KEY `idx_cate` (`category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区话题';

CREATE TABLE IF NOT EXISTS `qixi_community` (
  `community_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL DEFAULT '',
  `image` varchar(1000) NOT NULL DEFAULT '',
  `category_id` int(10) unsigned NOT NULL DEFAULT 0,
  `topic_id` int(10) unsigned NOT NULL DEFAULT 0,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `product_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '种草商品',
  `count_start` int(10) unsigned NOT NULL DEFAULT 0,
  `count_reply` int(10) unsigned NOT NULL DEFAULT 0,
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '0待审 1通过 -1驳回',
  `is_show` tinyint(4) NOT NULL DEFAULT 1,
  `is_hot` tinyint(4) NOT NULL DEFAULT 0,
  `is_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1图文',
  `content` varchar(1000) NOT NULL DEFAULT '',
  `refusal` varchar(255) DEFAULT NULL,
  `pv` int(11) NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status_time` datetime DEFAULT NULL,
  PRIMARY KEY (`community_id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_show` (`status`,`is_show`,`is_del`),
  KEY `idx_topic` (`topic_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区图文';

CREATE TABLE IF NOT EXISTS `qixi_community_reply` (
  `reply_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) NOT NULL DEFAULT '',
  `pid` int(10) unsigned NOT NULL DEFAULT 0,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `community_id` int(10) unsigned NOT NULL DEFAULT 0,
  `status` tinyint(4) NOT NULL DEFAULT 1,
  `is_del` tinyint(4) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`reply_id`),
  KEY `idx_community` (`community_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社区评论';

INSERT INTO `qixi_community_category` (`category_id`, `cate_name`, `pid`, `is_show`, `sort`)
SELECT 1, '种草分享', 0, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_community_category` WHERE `category_id` = 1);

INSERT INTO `qixi_community_topic` (`topic_id`, `topic_name`, `status`, `is_hot`, `category_id`, `sort`)
SELECT 1, '好物推荐', 1, 1, 1, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_community_topic` WHERE `topic_id` = 1);

INSERT INTO `qixi_community` (
  `community_id`, `title`, `image`, `category_id`, `topic_id`, `uid`, `mer_id`, `product_id`,
  `count_start`, `count_reply`, `status`, `is_show`, `is_hot`, `content`
)
SELECT 1, '春季上新 · 亲测好用', '', 1, 1, 1, 1, 1,
  3, 1, 1, 1, 1, '演示种草帖：这款商品手感不错，回购中。'
WHERE EXISTS (SELECT 1 FROM `qixi_user` WHERE `uid` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_community` WHERE `community_id` = 1);

INSERT INTO `qixi_community_reply` (`reply_id`, `content`, `uid`, `community_id`, `status`)
SELECT 1, '同款已下单，等收货！', 1, 1, 1
WHERE EXISTS (SELECT 1 FROM `qixi_community` WHERE `community_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_community_reply` WHERE `reply_id` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 28, 20, '/content/community', '', '社区种草', 'ContentCommunity', '', 70, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 28 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',28')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%28%';

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 123, 115, '/marketing/community', '', '逛逛社区', 'MerCommunityList', '', 70, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 123 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',123')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%123%';

UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 7, 'name', '社区', 'icon', '', 'url', '/pages/community/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '社区', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-community', '阶段6e：社区种草竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-community');
