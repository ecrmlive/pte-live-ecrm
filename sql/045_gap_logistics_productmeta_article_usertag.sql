-- 缺口竖切：物流 / 商品辅资料 / 文章 / 用户标签
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_express` (
  `express_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '快递公司名',
  `code` varchar(32) NOT NULL DEFAULT '' COMMENT '编码',
  `sort` int NOT NULL DEFAULT 0,
  `is_show` tinyint(1) NOT NULL DEFAULT 1,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`express_id`),
  KEY `idx_show` (`is_show`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='快递公司';

CREATE TABLE IF NOT EXISTS `qixi_system_city` (
  `city_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `level` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1省 2市 3区',
  `is_show` tinyint(1) NOT NULL DEFAULT 1,
  PRIMARY KEY (`city_id`),
  KEY `idx_parent` (`parent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='省市区';

CREATE TABLE IF NOT EXISTS `qixi_shipping_template` (
  `template_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1按件数 2按重量 3包邮',
  `appoint` tinyint(1) NOT NULL DEFAULT 0 COMMENT '指定区域',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`template_id`),
  KEY `idx_mer` (`mer_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运费模板';

CREATE TABLE IF NOT EXISTS `qixi_shipping_template_region` (
  `region_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `template_id` int(10) unsigned NOT NULL DEFAULT 0,
  `city_ids` varchar(512) NOT NULL DEFAULT '' COMMENT '城市id逗号分隔，空=默认',
  `first` decimal(10,2) NOT NULL DEFAULT 1.00 COMMENT '首件/首重',
  `first_price` decimal(10,2) NOT NULL DEFAULT 0.00,
  `continue` decimal(10,2) NOT NULL DEFAULT 1.00,
  `continue_price` decimal(10,2) NOT NULL DEFAULT 0.00,
  PRIMARY KEY (`region_id`),
  KEY `idx_tpl` (`template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='运费模板区域';

CREATE TABLE IF NOT EXISTS `qixi_store_product_label` (
  `label_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `info` varchar(255) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`label_id`),
  KEY `idx_mer` (`mer_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品标签';

CREATE TABLE IF NOT EXISTS `qixi_store_product_guarantee` (
  `guarantee_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `content` varchar(512) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`guarantee_id`),
  KEY `idx_mer` (`mer_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='保障服务';

CREATE TABLE IF NOT EXISTS `qixi_store_product_attr_template` (
  `template_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `template_name` varchar(64) NOT NULL DEFAULT '',
  `template_value` text NOT NULL COMMENT '参数JSON',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`template_id`),
  KEY `idx_mer` (`mer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品参数模板';

CREATE TABLE IF NOT EXISTS `qixi_article_category` (
  `cid` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(64) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 1,
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`cid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章分类';

CREATE TABLE IF NOT EXISTS `qixi_article` (
  `article_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `cid` int(10) unsigned NOT NULL DEFAULT 0,
  `title` varchar(128) NOT NULL DEFAULT '',
  `author` varchar(64) NOT NULL DEFAULT '',
  `image` varchar(255) NOT NULL DEFAULT '',
  `synopsis` varchar(255) NOT NULL DEFAULT '',
  `content` longtext NOT NULL,
  `visit` int NOT NULL DEFAULT 0,
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1展示',
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`article_id`),
  KEY `idx_cid` (`cid`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章';

CREATE TABLE IF NOT EXISTS `qixi_user_label` (
  `label_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `label_name` varchar(64) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户标签';

CREATE TABLE IF NOT EXISTS `qixi_user_group` (
  `group_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(64) NOT NULL DEFAULT '',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户分组';

CREATE TABLE IF NOT EXISTS `qixi_user_label_relation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `label_id` int(10) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_uid_label` (`uid`,`label_id`),
  KEY `idx_label` (`label_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户标签关系';

-- 种子
INSERT INTO `qixi_express` (`express_id`, `name`, `code`, `sort`, `is_show`)
SELECT 1, '顺丰速运', 'shunfeng', 100, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_express` WHERE `express_id` = 1);
INSERT INTO `qixi_express` (`express_id`, `name`, `code`, `sort`, `is_show`)
SELECT 2, '中通快递', 'zhongtong', 90, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_express` WHERE `express_id` = 2);

INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 1, 0, '北京市', 1 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 1);
INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 2, 1, '北京市', 2 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 2);
INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 3, 2, '朝阳区', 3 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 3);
INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 10, 0, '上海市', 1 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 10);
INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 11, 10, '上海市', 2 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 11);
INSERT INTO `qixi_system_city` (`city_id`, `parent_id`, `name`, `level`)
SELECT 12, 11, '浦东新区', 3 WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_city` WHERE `city_id` = 12);

INSERT INTO `qixi_shipping_template` (`template_id`, `mer_id`, `name`, `type`, `sort`)
SELECT 1, 1, '演示全国包邮', 3, 0
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_shipping_template` WHERE `template_id` = 1);
INSERT INTO `qixi_shipping_template_region` (`region_id`, `template_id`, `city_ids`, `first`, `first_price`, `continue`, `continue_price`)
SELECT 1, 1, '', 1, 0, 1, 0
WHERE EXISTS (SELECT 1 FROM `qixi_shipping_template` WHERE `template_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_shipping_template_region` WHERE `region_id` = 1);

INSERT INTO `qixi_store_product_label` (`label_id`, `mer_id`, `name`, `info`, `sort`, `status`)
SELECT 1, 1, '热卖', '演示标签', 0, 1
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_label` WHERE `label_id` = 1);

INSERT INTO `qixi_store_product_guarantee` (`guarantee_id`, `mer_id`, `name`, `content`, `sort`, `status`)
SELECT 1, 1, '七天无理由', '演示保障服务', 0, 1
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_guarantee` WHERE `guarantee_id` = 1);

INSERT INTO `qixi_store_product_attr_template` (`template_id`, `mer_id`, `template_name`, `template_value`, `sort`)
SELECT 1, 1, '服装参数', '[{"name":"材质","value":"棉"},{"name":"尺码","value":"均码"}]', 0
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_attr_template` WHERE `template_id` = 1);

INSERT INTO `qixi_article_category` (`cid`, `title`, `status`, `sort`)
SELECT 1, '商城资讯', 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_article_category` WHERE `cid` = 1);

INSERT INTO `qixi_article` (`article_id`, `cid`, `title`, `author`, `synopsis`, `content`, `sort`, `status`)
SELECT 1, 1, '欢迎使用栖息多商户商城', '小栖', '演示文章',
  '<p>本篇为阶段缺口竖切演示文章，可在平台后台编辑。</p>', 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_article` WHERE `article_id` = 1);

INSERT INTO `qixi_user_label` (`label_id`, `label_name`, `sort`)
SELECT 1, '高价值用户', 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_user_label` WHERE `label_id` = 1);
INSERT INTO `qixi_user_group` (`group_id`, `group_name`, `sort`)
SELECT 1, '默认分组', 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_user_group` WHERE `group_id` = 1);

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'gap-045', '物流/商品辅资料/文章/用户标签竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'gap-045');
