-- 阶段 1：身份 / RBAC 最小集（平台管理员、商户管理员、角色、菜单）
-- 表前缀 qixi_；字段对齐 docs/schema/
-- 种子账号密码均为 admin123（bcrypt），仅本地开发用

USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_system_admin` (
  `admin_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '后台管理员表ID',
  `account` varchar(32) NOT NULL COMMENT '后台管理员账号',
  `pwd` varchar(64) NOT NULL COMMENT '后台管理员密码',
  `real_name` varchar(16) NOT NULL COMMENT '后台管理员姓名',
  `phone` varchar(12) DEFAULT NULL COMMENT '联系电话',
  `roles` varchar(128) NOT NULL COMMENT '后台管理员权限(role_id), 多个逗号分隔',
  `last_ip` varchar(16) DEFAULT NULL COMMENT '后台管理员最后一次登录ip',
  `last_time` timestamp NULL DEFAULT NULL COMMENT '后台管理员最后一次登录时间',
  `login_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '登录次数',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '后台管理员状态 1有效0无效',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '1',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '后台管理员添加时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '后台管理员编辑时间',
  `region_ids` varchar(200) NOT NULL DEFAULT '' COMMENT '分组ID',
  `is_agent` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为区域代理：0否1是',
  `circle_agent_id` int(10) NOT NULL DEFAULT '0' COMMENT '商圈代理id',
  PRIMARY KEY (`admin_id`) USING BTREE,
  KEY `account` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='后台管理员表';

CREATE TABLE IF NOT EXISTS `qixi_system_role` (
  `role_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '身份管理id',
  `role_name` varchar(32) NOT NULL COMMENT '身份管理名称',
  `rules` text NOT NULL COMMENT '身份管理权限(menus_id)',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `is_agent` tinyint(1) NOT NULL DEFAULT '0' COMMENT '角色类型：0:平台,1:区域,2:商户',
  `circle_id` int(10) NOT NULL DEFAULT '0' COMMENT '商圈ID',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否默认角色',
  PRIMARY KEY (`role_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='身份管理表';

CREATE TABLE IF NOT EXISTS `qixi_system_menu` (
  `menu_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `path` varchar(512) NOT NULL COMMENT '路径',
  `icon` varchar(32) DEFAULT '' COMMENT '图标',
  `menu_name` varchar(128) NOT NULL DEFAULT '' COMMENT '按钮名',
  `route` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL COMMENT '路由名称',
  `params` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL DEFAULT '' COMMENT '参数',
  `sort` tinyint(4) NOT NULL DEFAULT '1' COMMENT '排序',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否显示',
  `is_mer` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '模块，1 平台， 2商户',
  `is_menu` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '类型，1菜单 2 权限',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_agent` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:平台,1:区域,2:商户',
  PRIMARY KEY (`menu_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

CREATE TABLE IF NOT EXISTS `qixi_merchant` (
  `mer_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户id',
  `category_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户分类 id',
  `type_id` int(10) unsigned DEFAULT '0' COMMENT '店铺类型 id',
  `mer_name` varchar(32) NOT NULL DEFAULT '' COMMENT '商户名称',
  `real_name` varchar(32) NOT NULL DEFAULT '' COMMENT '商户姓名',
  `mer_phone` varchar(13) NOT NULL DEFAULT '' COMMENT '商户手机号',
  `mer_address` varchar(64) NOT NULL DEFAULT '' COMMENT '商户地址',
  `mer_keyword` varchar(64) NOT NULL DEFAULT '' COMMENT '商户关键字',
  `mer_avatar` varchar(128) DEFAULT NULL COMMENT '商户头像',
  `mer_banner` varchar(128) DEFAULT NULL COMMENT '商户banner图片',
  `mini_banner` varchar(128) DEFAULT NULL COMMENT '商户店店铺街图片',
  `sales` int(10) unsigned DEFAULT '0' COMMENT '销量',
  `product_score` decimal(11,1) DEFAULT '5.0' COMMENT '商品描述评分',
  `service_score` decimal(11,1) DEFAULT '5.0' COMMENT '服务评分',
  `postage_score` decimal(11,1) DEFAULT '5.0' COMMENT '物流评分',
  `mark` varchar(256) NOT NULL DEFAULT '' COMMENT '商户备注',
  `reg_admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '总后台管理员ID',
  `sort` int(10) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '商户是否禁用0锁定,1正常',
  `commission_rate` decimal(6,2) DEFAULT NULL COMMENT '提成比例',
  `commission_switch` int(11) DEFAULT '0' COMMENT '商户手续费单独设置 0 关闭 1 开启',
  `long` varchar(16) DEFAULT NULL COMMENT '经度',
  `lat` varchar(16) DEFAULT NULL COMMENT '纬度',
  `is_del` tinyint(3) NOT NULL DEFAULT '0' COMMENT '0未删除1删除',
  `is_audit` tinyint(3) NOT NULL DEFAULT '0' COMMENT '添加的产品是否审核0不审核1审核',
  `is_bro_room` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否审核直播间0不审核1审核',
  `is_bro_goods` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否审核直播商品0不审核1审核',
  `is_best` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否推荐',
  `is_trader` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否自营',
  `mer_state` tinyint(3) NOT NULL DEFAULT '0' COMMENT '商户是否1开启0关闭',
  `mer_info` varchar(256) NOT NULL DEFAULT '' COMMENT '店铺简介',
  `service_phone` varchar(13) NOT NULL DEFAULT '' COMMENT '店铺电话',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `care_count` int(10) DEFAULT '0' COMMENT '关注总数',
  `copy_product_num` int(10) DEFAULT '0' COMMENT '剩余复制商品次数',
  `export_dump_num` int(10) DEFAULT '0' COMMENT '电子面单剩余次数',
  `mer_money` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '商户余额',
  `financial_bank` varchar(255) DEFAULT NULL COMMENT '银行卡转账信息',
  `financial_wechat` varchar(255) DEFAULT NULL COMMENT '微信转账信息',
  `financial_alipay` varchar(255) DEFAULT NULL COMMENT '支付宝转账信息',
  `financial_type` tinyint(3) DEFAULT '1' COMMENT '默认使用类型',
  `sub_mchid` varchar(16) NOT NULL DEFAULT '' COMMENT '微信支付分配的分账号',
  `delivery_way` varchar(50) DEFAULT '' COMMENT '配送方式',
  `delivery_balance` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '配送余额',
  `margin` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '保证金',
  `margin_remind_time` varchar(255) DEFAULT NULL COMMENT '保证金补缴提醒结束时间',
  `ot_margin` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '保证金额度',
  `is_margin` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否有保证金',
  `offline_switch` tinyint(4) DEFAULT '0' COMMENT '线下支付功能开关',
  `care_ficti` int(11) DEFAULT '0' COMMENT '虚拟关注量',
  `region_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户所属分组',
  `applyment_id` varchar(50) NOT NULL DEFAULT '' COMMENT '特约商户ID',
  `business_id` int(10) NOT NULL DEFAULT '0' COMMENT '店铺所属商户id',
  `applyment_switch` int(11) NOT NULL DEFAULT '1' COMMENT '特约商户分账比例是否合理',
  PRIMARY KEY (`mer_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户表';

CREATE TABLE IF NOT EXISTS `qixi_merchant_admin` (
  `merchant_admin_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户管理员表ID',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户ID(属于哪一个商户)',
  `account` varchar(32) NOT NULL COMMENT '商户管理员账号',
  `pwd` char(64) NOT NULL COMMENT '商户管理员密码',
  `real_name` varchar(16) NOT NULL COMMENT '商户管理员姓名',
  `phone` varchar(13) DEFAULT NULL COMMENT '商户管理员手机号',
  `last_ip` varchar(16) DEFAULT NULL COMMENT '商户管理员最后一次登录IP地址',
  `last_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '商户管理员最后一次登录时间',
  `roles` varchar(128) DEFAULT '',
  `login_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户管理员登录次数',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商户管理员等级(管理员添加的为0, 商户添加的为1)',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否有效 1有效 0无效 ',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '商户管理员添加时间',
  PRIMARY KEY (`merchant_admin_id`) USING BTREE,
  KEY `account` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户管理员表';

-- 平台菜单子集（is_mer=1）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 1 AS menu_id, 0 AS pid, '/dashboard' AS path, 'DashboardOutlined' AS icon, '工作台' AS menu_name, 'Dashboard' AS route, '' AS params, 100 AS sort, 1 AS is_show, 1 AS is_mer, 1 AS is_menu, 0 AS is_agent
  UNION ALL SELECT 2, 0, '/merchant', 'ShopOutlined', '商户', 'Merchant', '', 90, 1, 1, 1, 0
  UNION ALL SELECT 3, 2, '/merchant/list', '', '商户列表', 'MerchantList', '', 89, 1, 1, 1, 0
  UNION ALL SELECT 4, 2, '/merchant/audit', '', '入驻审核', 'MerchantAudit', '', 88, 1, 1, 1, 0
  UNION ALL SELECT 5, 0, '/product', 'ShoppingOutlined', '商品', 'Product', '', 80, 1, 1, 1, 0
  UNION ALL SELECT 6, 5, '/product/audit', '', '商品审核', 'ProductAudit', '', 79, 1, 1, 1, 0
  UNION ALL SELECT 7, 5, '/product/category', '', '平台分类', 'ProductCategory', '', 78, 1, 1, 1, 0
  UNION ALL SELECT 8, 0, '/order', 'ProfileOutlined', '订单', 'Order', '', 70, 1, 1, 1, 0
  UNION ALL SELECT 9, 8, '/order/list', '', '订单监管', 'OrderList', '', 69, 1, 1, 1, 0
  UNION ALL SELECT 10, 0, '/setting', 'SettingOutlined', '设置', 'Setting', '', 10, 1, 1, 1, 0
  UNION ALL SELECT 11, 10, '/setting/admin', '', '管理员', 'SettingAdmin', '', 9, 1, 1, 1, 0
  UNION ALL SELECT 12, 10, '/setting/role', '', '角色管理', 'SettingRole', '', 8, 1, 1, 1, 0
  UNION ALL SELECT 13, 10, '/setting/menu', '', '菜单管理', 'SettingMenu', '', 7, 1, 1, 1, 0
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 1 AND `is_mer` = 1);

-- 商户菜单子集（is_mer=2）
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT * FROM (
  SELECT 101 AS menu_id, 0 AS pid, '/dashboard' AS path, 'DashboardOutlined' AS icon, '工作台' AS menu_name, 'MerDashboard' AS route, '' AS params, 100 AS sort, 1 AS is_show, 2 AS is_mer, 1 AS is_menu, 2 AS is_agent
  UNION ALL SELECT 102, 0, '/product', 'ShoppingOutlined', '商品', 'MerProduct', '', 90, 1, 2, 1, 2
  UNION ALL SELECT 103, 102, '/product/list', '', '商品列表', 'MerProductList', '', 89, 1, 2, 1, 2
  UNION ALL SELECT 104, 102, '/product/edit', '', '发布商品', 'MerProductEdit', '', 88, 1, 2, 1, 2
  UNION ALL SELECT 105, 0, '/order', 'ProfileOutlined', '订单', 'MerOrder', '', 80, 1, 2, 1, 2
  UNION ALL SELECT 106, 105, '/order/list', '', '订单列表', 'MerOrderList', '', 79, 1, 2, 1, 2
  UNION ALL SELECT 107, 105, '/order/delivery', '', '待发货', 'MerOrderDelivery', '', 78, 1, 2, 1, 2
  UNION ALL SELECT 108, 0, '/finance', 'AccountBookOutlined', '财务', 'MerFinance', '', 70, 1, 2, 1, 2
  UNION ALL SELECT 109, 108, '/finance/balance', '', '店铺余额', 'MerFinanceBalance', '', 69, 1, 2, 1, 2
  UNION ALL SELECT 110, 0, '/setting', 'SettingOutlined', '设置', 'MerSetting', '', 10, 1, 2, 1, 2
  UNION ALL SELECT 111, 110, '/setting/shop', '', '店铺资料', 'MerSettingShop', '', 9, 1, 2, 1, 2
  UNION ALL SELECT 112, 110, '/setting/staff', '', '员工账号', 'MerSettingStaff', '', 8, 1, 2, 1, 2
) AS seed
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 101 AND `is_mer` = 2);

INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 1, '平台超级管理员', '1,2,3,4,5,6,7,8,9,10,11,12,13', 1, 0, 0, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 1);

INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 2, '商户主账号', '101,102,103,104,105,106,107,108,109,110,111,112', 1, 0, 2, 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 2);

-- 密码 admin123 → bcrypt
INSERT INTO `qixi_system_admin` (`admin_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `status`, `level`, `is_del`)
SELECT 1, 'admin', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '平台管理员', '13800000000', '1', 1, 0, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_admin` WHERE `account` = 'admin');

INSERT INTO `qixi_merchant` (`mer_id`, `mer_name`, `real_name`, `mer_phone`, `mer_address`, `mark`, `status`, `mer_state`, `is_del`)
SELECT 1, '演示商户', '张三', '13900000001', '上海市演示路1号', '阶段1种子商户', 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1);

INSERT INTO `qixi_merchant_admin` (`merchant_admin_id`, `mer_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `level`, `status`, `is_del`)
SELECT 1, 1, 'meradmin', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '商户管理员', '13900000001', '2', 0, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_admin` WHERE `account` = 'meradmin');

-- C 端用户（阶段 1 登录/注册最小集；字段对齐 docs/schema/domain-user.md）
CREATE TABLE IF NOT EXISTS `qixi_user` (
  `uid` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `wechat_user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '微信用户 id',
  `account` varchar(32) NOT NULL COMMENT '用户账号',
  `pwd` varchar(128) NOT NULL COMMENT '用户密码',
  `real_name` varchar(25) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `sex` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '性别',
  `birthday` date DEFAULT NULL COMMENT '生日',
  `card_id` varchar(20) NOT NULL DEFAULT '' COMMENT '身份证号码',
  `mark` varchar(255) NOT NULL DEFAULT '' COMMENT '用户备注',
  `label_id` varchar(64) DEFAULT NULL COMMENT '用户标签 id',
  `group_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户分组id',
  `nickname` varchar(16) NOT NULL COMMENT '用户昵称',
  `avatar` varchar(256) NOT NULL DEFAULT '' COMMENT '用户头像',
  `phone` char(15) DEFAULT NULL COMMENT '手机号码',
  `addres` varchar(128) DEFAULT NULL COMMENT '地址',
  `cancel_time` timestamp NULL DEFAULT NULL COMMENT '注销时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `last_time` timestamp NULL DEFAULT NULL COMMENT '最后一次登录时间',
  `last_ip` varchar(16) NOT NULL DEFAULT '' COMMENT '最后一次登录ip',
  `now_money` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '用户余额',
  `brokerage_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '佣金金额',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1为正常，0为禁止',
  `spread_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '推广员id',
  `spread_time` timestamp NULL DEFAULT NULL COMMENT '推广员关联时间',
  `spread_limit` timestamp NULL DEFAULT NULL COMMENT '推广员到期时间',
  `brokerage_level` int(10) DEFAULT '0' COMMENT '推广员等级',
  `user_type` varchar(32) NOT NULL DEFAULT 'h5' COMMENT '用户类型',
  `promoter_time` timestamp NULL DEFAULT NULL COMMENT '成功推广时间',
  `is_promoter` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否为推广员',
  `main_uid` int(10) unsigned DEFAULT '0' COMMENT '主账号',
  `pay_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户购买次数',
  `pay_price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '用户消费金额',
  `spread_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '下级人数',
  `spread_pay_count` int(10) unsigned DEFAULT '0' COMMENT '下级订单数',
  `spread_pay_price` decimal(10,2) DEFAULT '0.00' COMMENT '下级订单金额',
  `integral` int(10) DEFAULT '0' COMMENT '积分',
  `member_level` int(10) DEFAULT '0' COMMENT '免费会员等级',
  `member_value` int(10) DEFAULT '0' COMMENT '免费会员成长值',
  `count_start` int(10) DEFAULT '0' COMMENT '用户获赞数',
  `count_fans` int(10) DEFAULT '0' COMMENT '用户粉丝数',
  `count_content` int(10) DEFAULT '0' COMMENT '用户内容数量',
  `is_svip` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否为付费会员',
  `svip_endtime` timestamp NULL DEFAULT NULL COMMENT '会员结束时间',
  `svip_save_money` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '会员节省金额',
  `promoter_switch` tinyint(1) NOT NULL DEFAULT '1' COMMENT '分销资格 0无 1有',
  PRIMARY KEY (`uid`) USING BTREE,
  UNIQUE KEY `account` (`account`) USING BTREE,
  KEY `phone` (`phone`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

INSERT INTO `qixi_user` (`uid`, `account`, `pwd`, `nickname`, `avatar`, `phone`, `status`, `user_type`, `last_ip`)
SELECT 1, 'demo', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '演示用户', '', '13700000001', 1, 'h5', ''
WHERE NOT EXISTS (SELECT 1 FROM `qixi_user` WHERE `account` = 'demo');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase1-identity', '阶段1身份表：平台/商户管理员、C端用户、菜单子集'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase1-identity');
