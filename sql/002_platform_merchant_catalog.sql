-- 阶段 2（平台侧）：商户入驻审核 + 平台类目/品牌 + 商品审核最小集
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_merchant_category` (
  `merchant_category_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户分类 id',
  `commission_rate` decimal(6,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '手续费',
  `category_name` varchar(32) NOT NULL COMMENT '商户分类名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`merchant_category_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户分类表';

CREATE TABLE IF NOT EXISTS `qixi_merchant_intention` (
  `mer_intention_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int(10) unsigned DEFAULT '0' COMMENT '用户ID',
  `phone` varchar(11) DEFAULT NULL COMMENT '手机号',
  `mer_name` varchar(30) DEFAULT NULL COMMENT '商户名称',
  `name` varchar(30) DEFAULT NULL COMMENT '客户姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `status` tinyint(4) DEFAULT '0' COMMENT '处理状态 0待审 1通过 2未通过',
  `fail_msg` varchar(255) DEFAULT NULL COMMENT '未通过原因',
  `is_del` tinyint(4) DEFAULT '0' COMMENT '删除状态 1删除 ，0未删除',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '关联商户',
  `images` varchar(2000) DEFAULT NULL COMMENT '多图',
  `merchant_category_id` int(10) unsigned DEFAULT '0' COMMENT '商户分类',
  `mer_type_id` int(10) unsigned DEFAULT '0' COMMENT '店铺类型',
  `circle_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商圈id',
  `business_id` int(10) NOT NULL DEFAULT '0' COMMENT '店铺所属商户id',
  PRIMARY KEY (`mer_intention_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户申请表';

CREATE TABLE IF NOT EXISTS `qixi_store_category` (
  `store_category_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '商品分类表ID',
  `pid` mediumint(9) NOT NULL DEFAULT '0' COMMENT '父id',
  `cate_name` varchar(100) NOT NULL COMMENT '分类名称',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
  `sort` mediumint(9) NOT NULL DEFAULT '0' COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `level` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户id，0=平台',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_hot` tinyint(1) DEFAULT '0' COMMENT '是否推荐',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 商品，1 积分商品',
  PRIMARY KEY (`store_category_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品分类表';

CREATE TABLE IF NOT EXISTS `qixi_store_brand` (
  `brand_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '商品品牌表ID',
  `brand_category_id` mediumint(9) NOT NULL DEFAULT '0' COMMENT '品牌分类id',
  `brand_name` varchar(100) NOT NULL COMMENT '品牌名称',
  `sort` mediumint(9) NOT NULL DEFAULT '0' COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`brand_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品品牌表';

-- 商品表（阶段2审核所需字段；后续阶段可 ALTER 补齐）
CREATE TABLE IF NOT EXISTS `qixi_store_product` (
  `product_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品id',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户Id',
  `store_name` varchar(128) NOT NULL COMMENT '商品名称',
  `store_info` varchar(256) DEFAULT NULL COMMENT '商品简介',
  `keyword` varchar(128) NOT NULL DEFAULT '' COMMENT '关键字',
  `bar_code` varchar(15) NOT NULL DEFAULT '' COMMENT '产品条码',
  `brand_id` int(11) DEFAULT '0' COMMENT '品牌 id',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商户状态 0未上架 1上架',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '平台审核 0审核中 1通过 -1未通过 -2下架',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `mer_status` tinyint(1) DEFAULT '1' COMMENT '商铺状态 1正常 0非正常',
  `cate_id` int(11) NOT NULL DEFAULT '0' COMMENT '分类id',
  `unit_name` varchar(16) NOT NULL DEFAULT '件' COMMENT '单位名',
  `sort` smallint(6) NOT NULL DEFAULT '0' COMMENT '排序',
  `rank` smallint(6) NOT NULL DEFAULT '0' COMMENT '总后台排序',
  `sales` mediumint(8) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '最低价格',
  `cost` decimal(10,2) DEFAULT '0.00' COMMENT '成本价',
  `ot_price` decimal(10,2) DEFAULT '0.00' COMMENT '原价',
  `stock` int(10) unsigned DEFAULT '0' COMMENT '总库存',
  `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0普通',
  `spec_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规格 0单 1多',
  `refusal` varchar(255) DEFAULT NULL COMMENT '审核拒绝理由',
  `image` varchar(256) NOT NULL DEFAULT '' COMMENT '商品图片',
  `slider_image` varchar(2000) NOT NULL DEFAULT '' COMMENT '轮播图',
  `delivery_way` varchar(100) DEFAULT '2' COMMENT '配送方式',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0实体 1虚拟',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`product_id`) USING BTREE,
  KEY `idx_mer_status` (`mer_id`,`status`,`is_del`),
  KEY `idx_cate` (`cate_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- 平台菜单补品牌
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 14, 5, '/product/brand', '', '品牌管理', 'ProductBrand', '', 77, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 14);

UPDATE `qixi_system_role`
SET `rules` = '1,2,3,4,5,6,7,8,9,10,11,12,13,14'
WHERE `role_id` = 1;

INSERT INTO `qixi_merchant_category` (`merchant_category_id`, `commission_rate`, `category_name`)
SELECT 1, 0.0600, '综合零售'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_category` WHERE `merchant_category_id` = 1);

INSERT INTO `qixi_merchant_intention` (`mer_intention_id`, `uid`, `phone`, `mer_name`, `name`, `status`, `merchant_category_id`, `mark`)
SELECT 1, 0, '13900000002', '待审旗舰店', '李四', 0, 1, '阶段2演示入驻申请'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_intention` WHERE `mer_intention_id` = 1);

INSERT INTO `qixi_store_category` (`store_category_id`, `pid`, `cate_name`, `path`, `sort`, `is_show`, `level`, `mer_id`, `type`)
SELECT 1, 0, '默认分类', '/1/', 100, 1, 0, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_category` WHERE `store_category_id` = 1);

INSERT INTO `qixi_store_category` (`store_category_id`, `pid`, `cate_name`, `path`, `sort`, `is_show`, `level`, `mer_id`, `type`)
SELECT 2, 1, '日用百货', '/1/2/', 90, 1, 1, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_category` WHERE `store_category_id` = 2);

INSERT INTO `qixi_store_brand` (`brand_id`, `brand_category_id`, `brand_name`, `sort`, `is_show`)
SELECT 1, 0, '栖息自有', 100, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_brand` WHERE `brand_id` = 1);

INSERT INTO `qixi_store_product` (
  `product_id`, `mer_id`, `store_name`, `store_info`, `keyword`, `brand_id`, `is_show`, `status`,
  `cate_id`, `unit_name`, `price`, `ot_price`, `stock`, `spec_type`, `image`, `slider_image`,
  `delivery_way`, `type`, `product_type`
)
SELECT 1, 1, '演示待审商品', '平台审核演示', '演示', 1, 1, 0,
  2, '件', 29.90, 39.90, 100, 0, '', '',
  '2', 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_product` WHERE `product_id` = 1);

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase2-platform-merchant-catalog', '阶段2平台：商户审核/类目品牌/商品审核'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase2-platform-merchant-catalog');
