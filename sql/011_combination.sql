-- 阶段 6：拼团活动（秒杀之后的第二种活动竖切）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_product_group` (
  `product_group_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `product_id` int(10) unsigned DEFAULT '0' COMMENT '商品ID',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `time` int(10) unsigned DEFAULT '0' COMMENT '开团时长(小时)',
  `buying_count_num` int(11) DEFAULT '0' COMMENT '成团总人数',
  `buying_num` int(11) DEFAULT '0' COMMENT '最少真实购买人数',
  `pay_count` int(10) unsigned DEFAULT '0' COMMENT '活动购买总人数',
  `once_pay_count` int(10) unsigned DEFAULT '0' COMMENT '单次购买数量',
  `status` int(11) DEFAULT '0' COMMENT '平台控制状态 1通过',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户ID',
  `ficti_status` int(11) DEFAULT '0' COMMENT '虚拟成团状态',
  `ficti_num` int(11) DEFAULT '0' COMMENT '最多虚拟人数',
  `is_show` int(11) DEFAULT '0' COMMENT '上下架',
  `is_del` int(10) unsigned DEFAULT '0',
  `success_num` int(10) unsigned DEFAULT '0' COMMENT '成功团数',
  `product_status` int(11) DEFAULT '0',
  `price` decimal(10,2) DEFAULT '0.00' COMMENT '拼团价',
  `action_status` int(11) DEFAULT '0' COMMENT '活动状态 1进行中',
  `create_time` datetime DEFAULT NULL,
  `refusal` varchar(255) DEFAULT NULL,
  `leader_extension` tinyint(4) DEFAULT '0',
  `leader_rate` decimal(10,2) DEFAULT '0.00',
  PRIMARY KEY (`product_group_id`) USING BTREE,
  KEY `idx_mer_product` (`mer_id`,`product_id`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='拼团商品';

CREATE TABLE IF NOT EXISTS `qixi_store_product_group_buying` (
  `group_buying_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_group_id` int(10) unsigned DEFAULT '0' COMMENT '活动商品ID',
  `status` int(11) DEFAULT '0' COMMENT '0进行中 10已完成 -1失败',
  `ficti_status` int(11) DEFAULT '0',
  `ficti_num` int(10) unsigned DEFAULT '0',
  `buying_count_num` int(10) unsigned DEFAULT '0' COMMENT '成团总人数',
  `buying_num` int(10) unsigned DEFAULT '0',
  `yet_buying_num` int(10) unsigned DEFAULT '0' COMMENT '已参团人数(已支付)',
  `is_del` int(11) DEFAULT '0',
  `mer_id` int(10) unsigned DEFAULT '0',
  `end_time` int(11) DEFAULT NULL COMMENT '结束时间unix',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_hidde` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`group_buying_id`),
  KEY `idx_group` (`product_group_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='拼团团次';

CREATE TABLE IF NOT EXISTS `qixi_store_product_group_user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_buying_id` int(10) unsigned DEFAULT '0' COMMENT '团ID',
  `product_group_id` int(10) unsigned DEFAULT '0' COMMENT '活动商品ID',
  `status` int(11) DEFAULT '0' COMMENT '0待支付 1已支付',
  `is_initiator` int(10) unsigned DEFAULT '0' COMMENT '是否团长',
  `order_id` int(10) unsigned DEFAULT '0' COMMENT '订单ID',
  `uid` int(10) unsigned DEFAULT '0',
  `nickname` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `is_del` int(10) unsigned DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_leader` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_buying_uid` (`group_buying_id`,`uid`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='拼团成员';

-- 演示：商品2 拼团 2人成团 ¥19.90
INSERT INTO `qixi_store_product_group` (
  `product_group_id`, `product_id`, `start_time`, `end_time`, `time`,
  `buying_count_num`, `buying_num`, `once_pay_count`, `status`, `mer_id`,
  `is_show`, `price`, `action_status`, `create_time`, `product_status`
)
SELECT 1, 2, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 24,
  2, 1, 1, 1, 1,
  1, 19.90, 1, NOW(), 1
WHERE EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 2 AND `is_del` = 0)
  AND NOT EXISTS (SELECT 1 FROM `qixi_store_product_group` WHERE `product_group_id` = 1);

-- 平台菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 24, 17, '/marketing/combination', '', '拼团监管', 'MarketingCombination', '', 43, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 24 AND `is_mer` = 1);

-- 商户菜单
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 118, 115, '/marketing/combination', '', '拼团活动', 'MerMarketingCombination', '', 73, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 118 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',24')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%24%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',118')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%118%';

-- DIY 入口菜单补一条拼团（若首页 DIY 已存在则追加 menus 项由后台编辑；此处仅文档种子说明）
-- C 端路由：/pages/combination/list

-- DIY 首页补充拼团入口（若已有 menus 则追加）
UPDATE `qixi_diy`
SET `value` = JSON_SET(
  `value`,
  '$.menus',
  JSON_ARRAY_APPEND(
    COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()),
    '$',
    JSON_OBJECT('id', 4, 'name', '拼团', 'icon', '', 'url', '/pages/combination/list')
  )
)
WHERE `id` = 1
  AND JSON_SEARCH(COALESCE(JSON_EXTRACT(`value`, '$.menus'), JSON_ARRAY()), 'one', '拼团', NULL, '$[*].name') IS NULL;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-combination', '阶段6：拼团活动竖切'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-combination');
