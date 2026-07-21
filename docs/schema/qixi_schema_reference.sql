-- qixi-live-mergers schema reference
-- Source: CRMEB MER v4.0 install/crmeb_merchant.sql
-- Prefix: eb_ -> qixi_
-- NOTE: Reference DDL for redesign; not a drop-in runtime dump.
-- Charset: utf8mb4

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `qixi_article`;
CREATE TABLE `qixi_article` (

  `article_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章管理ID',
  `cid` int(10) unsigned DEFAULT '0' COMMENT '分类id',
  `title` varchar(64) NOT NULL COMMENT '文章标题',
  `author` varchar(32) DEFAULT NULL COMMENT '文章作者',
  `image_input` varchar(128) NOT NULL COMMENT '文章图片',
  `synopsis` varchar(128) DEFAULT NULL COMMENT '文章简介',
  `visit` varchar(255) DEFAULT NULL COMMENT '浏览次数',
  `sort` int(10) unsigned DEFAULT '0' COMMENT '排序',
  `url` varchar(128) DEFAULT NULL COMMENT '原文链接',
  `admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '管理员id',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户id',
  `is_hot` tinyint(3) unsigned DEFAULT '0' COMMENT '是否热门(小程序)',
  `is_banner` tinyint(3) unsigned DEFAULT '0' COMMENT '是否轮播图(小程序)',
  `status` tinyint(3) unsigned DEFAULT NULL COMMENT '状态',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `wechat_news_id` int(11) DEFAULT '0' COMMENT '微信图文id',
  PRIMARY KEY (`article_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章管理表';

DROP TABLE IF EXISTS `qixi_article_category`;
CREATE TABLE `qixi_article_category` (

  `article_category_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章分类id',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级ID',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `title` varchar(32) NOT NULL COMMENT '文章分类标题',
  `info` varchar(255) DEFAULT NULL COMMENT '文章分类简介',
  `image` varchar(128) NOT NULL COMMENT '文章分类图片',
  `status` tinyint(3) unsigned NOT NULL COMMENT '状态',
  `sort` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`article_category_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章分类表';

DROP TABLE IF EXISTS `qixi_article_content`;
CREATE TABLE `qixi_article_content` (

  `article_content_id` int(10) unsigned NOT NULL COMMENT '文章id',
  `content` text NOT NULL COMMENT '文章内容',
  UNIQUE KEY `article_content_id` (`article_content_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='文章内容表';

DROP TABLE IF EXISTS `qixi_broadcast_assistant`;
CREATE TABLE `qixi_broadcast_assistant` (

  `assistant_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) DEFAULT NULL COMMENT '微信号',
  `nickname` varchar(100) DEFAULT NULL COMMENT '微信昵称',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `is_del` tinyint(1) DEFAULT '0',
  UNIQUE KEY `id` (`assistant_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='直播助手信息';

DROP TABLE IF EXISTS `qixi_broadcast_goods`;
CREATE TABLE `qixi_broadcast_goods` (

  `broadcast_goods_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `goods_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '微信商品ID',
  `audit_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '审核单 id',
  `cover_img` varchar(255) NOT NULL COMMENT '图片',
  `name` varchar(64) NOT NULL COMMENT '商品名称',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '价格',
  `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商品类型',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品 id',
  `error_msg` varchar(255) DEFAULT NULL COMMENT '未通过原因',
  `audit_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0：未审核，1：审核中，2:审核通过，3审核失败',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态0=未审核1=微信审核2=审核通过-1=审核未通过',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否显示',
  `is_mer_show` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商户是否显示',
  `mark` varchar(512) DEFAULT NULL COMMENT '备注',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_mer_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商户是否删除',
  PRIMARY KEY (`broadcast_goods_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `goods_id` (`goods_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='直播商品表';

DROP TABLE IF EXISTS `qixi_broadcast_room`;
CREATE TABLE `qixi_broadcast_room` (

  `broadcast_room_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `room_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '直播间 id',
  `name` varchar(32) NOT NULL COMMENT '直播间名字',
  `cover_img` varchar(255) NOT NULL COMMENT '背景图',
  `share_img` varchar(255) NOT NULL COMMENT '分享图',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '直播计划开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '直播计划结束时间',
  `anchor_name` varchar(32) NOT NULL COMMENT '主播昵称',
  `anchor_wechat` varchar(32) NOT NULL COMMENT '主播微信号',
  `phone` varchar(32) NOT NULL COMMENT '主播手机号',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '直播间类型 【1: 推流，0：手机直播】',
  `screen_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '横屏、竖屏 【1：横屏，0：竖屏】',
  `close_like` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否关闭点赞',
  `close_goods` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否关闭货架',
  `close_comment` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否关闭评论',
  `close_share` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否关闭分享',
  `close_kf` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否关闭客服',
  `error_msg` varchar(255) DEFAULT NULL COMMENT '未通过原因',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '审核状态0=未审核1=微信审核2=审核通过-1=审核未通过',
  `live_status` smallint(5) unsigned NOT NULL DEFAULT '102' COMMENT '直播状态101：直播中，102：未开始，103已结束，104禁播，105：暂停，106：异常，107：已过期',
  `mark` varchar(512) DEFAULT NULL COMMENT '备注',
  `replay_status` tinyint(3) unsigned DEFAULT '0' COMMENT '回放状态',
  `is_mer_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商户是否显示',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否显示',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `star` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT '推荐星级',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `is_mer_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商户是否删除',
  `feeds_img` varchar(255) DEFAULT NULL COMMENT '封面图',
  `push_url` varchar(255) DEFAULT NULL COMMENT '推流地址',
  `assistant_id` varchar(255) DEFAULT NULL COMMENT '小助手ID',
  `is_feeds_public` tinyint(3) unsigned DEFAULT '0' COMMENT '是否开启官方收录，1 开启，0 关闭',
  PRIMARY KEY (`broadcast_room_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='直播间表';

DROP TABLE IF EXISTS `qixi_broadcast_room_goods`;
CREATE TABLE `qixi_broadcast_room_goods` (

  `broadcast_room_id` int(10) unsigned NOT NULL,
  `broadcast_goods_id` int(10) unsigned NOT NULL,
  `on_sale` tinyint(4) DEFAULT '1' COMMENT '商品上下架',
  KEY `broadcast_room_id` (`broadcast_room_id`,`broadcast_goods_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='直播间导入商品表';

DROP TABLE IF EXISTS `qixi_cache`;
CREATE TABLE `qixi_cache` (

  `key` varchar(32) NOT NULL,
  `expire_time` int(11) NOT NULL DEFAULT '0' COMMENT '0=永久',
  `result` longtext NOT NULL COMMENT '缓存数据',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '缓存时间',
  PRIMARY KEY (`key`) USING BTREE,
  KEY `key` (`key`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信缓存表';

DROP TABLE IF EXISTS `qixi_cdkey_library`;
CREATE TABLE `qixi_cdkey_library` (

  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户ID',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '卡密库名称',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `product_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联商品ID',
  `product_attr_value_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联商品skuID',
  `used_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '卡密已使用数量',
  `total_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '卡密总数量',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_circle`;
CREATE TABLE `qixi_circle` (

  `circle_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商圈id',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上级商圈id',
  `path` varchar(50) NOT NULL DEFAULT '' COMMENT '路径',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '商圈名称',
  `circle_agent_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商圈代理id',
  `commission_type` tinyint(3) NOT NULL DEFAULT '0' COMMENT '商圈提成类型(0:按默认设置，1:单独设置)',
  `commission_rate` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '商圈提成比例(0~100%)：commission_type为0时取系统设置等级提成比例，为1时自定义设置)',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '等级:0一级商圈 1二级商圈 2三级商圈',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '说明信息',
  `sort` int(10) NOT NULL DEFAULT '0' COMMENT '排序(数字越大越靠前)',
  `status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态：0禁用 1启用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型：0区域，1商户',
  `role_id` int(10) NOT NULL DEFAULT '0' COMMENT '角色权限id',
  `business_store_category` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户店铺分类',
  `business_store_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户店铺类型',
  PRIMARY KEY (`circle_id`) USING BTREE,
  KEY `idx_pid` (`pid`) USING BTREE,
  KEY `idx_circle_agent_id` (`circle_agent_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_circle_agent`;
CREATE TABLE `qixi_circle_agent` (

  `circle_agent_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商圈代理id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '代理名称',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '联系电话',
  `qualification` text COMMENT '身份资质',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '说明信息',
  `extend` text COMMENT '扩展信息',
  `audit_admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '审核管理员',
  `audit_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '审核拒绝原因',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态：0待审核 1审核通过 -1审核拒绝 -2撤销',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `payment_method` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '结算方式:0:银行卡，1:微信，2:支付宝',
  `payment_name` varchar(16) NOT NULL DEFAULT '' COMMENT '结算名称',
  `payment_account` varchar(30) NOT NULL DEFAULT '' COMMENT '结算账号',
  `payment_bank` varchar(30) NOT NULL DEFAULT '' COMMENT '开户行',
  `payment_qr_img` varchar(200) NOT NULL DEFAULT '' COMMENT '收款二维码图片',
  `balance` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '佣金余额',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型：0区域，1商户',
  `business_name` varchar(64) NOT NULL DEFAULT '' COMMENT '商户名称',
  `business_store_category` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户店铺分类',
  `business_store_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户店铺类型',
  PRIMARY KEY (`circle_agent_id`) USING BTREE,
  KEY `idx_uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商圈代理表';

DROP TABLE IF EXISTS `qixi_circle_brokerage_checkout`;
CREATE TABLE `qixi_circle_brokerage_checkout` (

  `checkout_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `agent_id` int(10) unsigned NOT NULL COMMENT '商圈代理id',
  `agent_name` varchar(64) NOT NULL DEFAULT '' COMMENT '商圈代理名称',
  `agent_phone` varchar(64) NOT NULL DEFAULT '' COMMENT '商圈代理电话',
  `withdrawal_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '提现流水号',
  `withdrawal_amount` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '提现金额',
  `withdrawal_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '提现方式 0:银行卡，1:微信，2:支付宝',
  `transfer_voucher` text COMMENT '转账凭证',
  `remark` varchar(200) NOT NULL DEFAULT '' COMMENT '备注',
  `platform_remark` varchar(200) NOT NULL DEFAULT '' COMMENT '平台备注',
  `transfer_remark` varchar(200) NOT NULL DEFAULT '' COMMENT '转账备注',
  `audit_time` datetime DEFAULT NULL COMMENT '审核时间',
  `audit_status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '审核状态：0待审核 1审核通过 -1审核拒绝 -2撤销',
  `audit_admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '审核管理员',
  `audit_reason` varchar(255) NOT NULL DEFAULT '' COMMENT '审核拒绝原因',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态：0未到帐 1已到帐',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `withdrawal_qr_img` varchar(200) NOT NULL DEFAULT '' COMMENT '收款二维码图片',
  `withdrawal_name` varchar(16) NOT NULL DEFAULT '' COMMENT '结算名称',
  `withdrawal_account` varchar(30) NOT NULL DEFAULT '' COMMENT '结算账号',
  PRIMARY KEY (`checkout_id`) USING BTREE,
  KEY `idx_agent_id` (`agent_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商圈佣金结算表';

DROP TABLE IF EXISTS `qixi_circle_financial_record`;
CREATE TABLE `qixi_circle_financial_record` (

  `record_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `circle_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商圈id',
  `circle_name` varchar(64) NOT NULL DEFAULT '' COMMENT '商圈名称',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户id',
  `mer_name` varchar(64) NOT NULL DEFAULT '' COMMENT '商户名称',
  `agent_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商圈代理id',
  `agent_name` varchar(64) NOT NULL DEFAULT '' COMMENT '商圈代理名称',
  `record_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '流水号',
  `order_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单号',
  `order_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '订单编号',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户 id',
  `user_info` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `amount` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
  `status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '状态：0冻结中 1解冻 -1失效',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`record_id`) USING BTREE,
  KEY `idx_mer_id` (`mer_id`) USING BTREE,
  KEY `idx_circle_id` (`circle_id`) USING BTREE,
  KEY `idx_agent_id` (`agent_id`) USING BTREE,
  KEY `idx_order_sn` (`order_sn`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商圈提成流水';

DROP TABLE IF EXISTS `qixi_city_area`;
CREATE TABLE `qixi_city_area` (

  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `path` varchar(128) NOT NULL DEFAULT '/' COMMENT '省市级别',
  `parent_id` int(11) NOT NULL DEFAULT '0' COMMENT '父级id',
  `type` varchar(32) NOT NULL COMMENT '类型',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '级别',
  `code` varchar(100) NOT NULL DEFAULT '' COMMENT '城市编码',
  `snum` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '子级个数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `parent_id` (`parent_id`) USING BTREE,
  KEY `path` (`path`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=44714 DEFAULT CHARSET=utf8 COMMENT='省市区县数据';

DROP TABLE IF EXISTS `qixi_community`;
CREATE TABLE `qixi_community` (

  `community_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `image` varchar(1000) DEFAULT NULL COMMENT '图片',
  `category_id` int(10) unsigned DEFAULT '0',
  `topic_id` int(10) unsigned DEFAULT '0' COMMENT '话题',
  `uid` int(10) unsigned DEFAULT '0' COMMENT '用户',
  `count_start` int(10) unsigned DEFAULT '0' COMMENT '点赞数',
  `count_reply` int(10) unsigned DEFAULT '0' COMMENT '评论数',
  `count_share` int(10) unsigned DEFAULT '0' COMMENT '分享数',
  `status` tinyint(4) DEFAULT '0' COMMENT '审核状态',
  `is_show` tinyint(4) DEFAULT '0' COMMENT '显示状态',
  `start` tinyint(1) DEFAULT '1' COMMENT '星级排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) DEFAULT '0',
  `content` varchar(1000) DEFAULT NULL,
  `refusal` varchar(255) DEFAULT NULL COMMENT '拒绝理由',
  `is_hot` tinyint(4) DEFAULT '0' COMMENT '是否推荐',
  `order_id` int(10) unsigned DEFAULT '0' COMMENT '关联订单ID',
  `is_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1 图文 2 视频',
  `video_link` varchar(255) DEFAULT NULL COMMENT '视频链接',
  `pv` int(11) DEFAULT '0' COMMENT '浏览量',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `status_time` timestamp NULL DEFAULT NULL COMMENT '审核时间',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户ID',
  PRIMARY KEY (`community_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='社区图文表信息';

DROP TABLE IF EXISTS `qixi_community_category`;
CREATE TABLE `qixi_community_category` (

  `category_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `cate_name` varchar(50) DEFAULT NULL COMMENT '分类名',
  `pid` int(11) DEFAULT NULL COMMENT '父级ID',
  `path` varchar(255) DEFAULT '/' COMMENT '路径 ',
  `is_show` tinyint(4) DEFAULT '1' COMMENT '状态',
  `level` int(11) DEFAULT '0' COMMENT '等级',
  `sort` int(11) DEFAULT NULL,
  PRIMARY KEY (`category_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='社区分类';

DROP TABLE IF EXISTS `qixi_community_reply`;
CREATE TABLE `qixi_community_reply` (

  `reply_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `content` varchar(255) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '评论内容',
  `pid` int(10) unsigned DEFAULT '0' COMMENT '回复id',
  `uid` int(10) unsigned DEFAULT '0' COMMENT '发言人',
  `re_uid` int(10) unsigned DEFAULT '0' COMMENT '回复人',
  `count_start` int(10) unsigned DEFAULT '0' COMMENT '点赞数',
  `count_reply` int(10) unsigned DEFAULT '0' COMMENT '评论数',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 ',
  `community_id` int(10) unsigned DEFAULT '0' COMMENT '文章id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(4) DEFAULT '0',
  `refusal` varchar(255) DEFAULT NULL COMMENT '拒绝原因',
  PRIMARY KEY (`reply_id`) USING BTREE,
  UNIQUE KEY `id` (`reply_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='社区评论';

DROP TABLE IF EXISTS `qixi_community_topic`;
CREATE TABLE `qixi_community_topic` (

  `topic_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `topic_name` varchar(100) DEFAULT NULL COMMENT '话题',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态',
  `is_hot` tinyint(4) DEFAULT '0' COMMENT '推荐',
  `category_id` int(10) unsigned DEFAULT '0' COMMENT '分类id',
  `is_del` tinyint(4) DEFAULT '0',
  `pic` varchar(128) DEFAULT NULL COMMENT '图标',
  `count_use` int(10) unsigned DEFAULT '0' COMMENT '使用次数',
  `count_view` int(10) unsigned DEFAULT '0' COMMENT '浏览量',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `sort` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`topic_id`) USING BTREE,
  UNIQUE KEY `id` (`topic_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='社区话题';

DROP TABLE IF EXISTS `qixi_delivery_config`;
CREATE TABLE `qixi_delivery_config` (

  `delivery_config_id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `mer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商户ID',
  `min_delivery_amount` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '起送价',
  `base_shipping_fee` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '基础运费',
  `free_shipping_amount` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '包邮规则',
  `is_premium_stack_enabled` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '是否开启溢价叠加(0:关 1:开)',
  `distance_premium_config` text NOT NULL COMMENT '距离溢价设置',
  `weight_premium_config` text NOT NULL COMMENT '重量溢价设置',
  `delivery_time_type` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '配送时间类型(1:可选定时送达 2:统一尽快送达)',
  `selectable_days` int(11) unsigned NOT NULL DEFAULT '7' COMMENT '可选天数',
  `delivery_prompt` varchar(200) NOT NULL DEFAULT '' COMMENT '送达文案',
  `status` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '配送状态(0:关 1:开)',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `commission_rate` int(10) NOT NULL DEFAULT '0' COMMENT '配送员提成比例(0~100%)',
  PRIMARY KEY (`delivery_config_id`),
  UNIQUE KEY `mer_id` (`mer_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='配送设置表';

DROP TABLE IF EXISTS `qixi_delivery_order`;
CREATE TABLE `qixi_delivery_order` (

  `delivery_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `station_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '门店ID',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单ID',
  `order_code` varchar(255) NOT NULL DEFAULT '' COMMENT '配送方订单号',
  `city_code` varchar(20) NOT NULL DEFAULT '' COMMENT '所属城市',
  `order_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '订单sn',
  `cargo_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '配送订单价格',
  `finish_code` varchar(255) NOT NULL DEFAULT '' COMMENT '收货码',
  `user_name` varchar(20) NOT NULL DEFAULT '' COMMENT '用户名',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 取消=-1, 待取货＝2,配送中＝3,已完成＝4,物品返回中=9,物品返回完成=10,骑士到店=100',
  `receiver_phone` varchar(11) NOT NULL DEFAULT '' COMMENT '收货人电话',
  `from_address` varchar(255) NOT NULL DEFAULT '' COMMENT '起始位置',
  `to_address` varchar(255) NOT NULL DEFAULT '' COMMENT '结束位置',
  `distance` float NOT NULL DEFAULT '0' COMMENT '配送距离',
  `fee` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '配送费',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户ID',
  `mark` varchar(255) NOT NULL DEFAULT '' COMMENT '订单备注',
  `station_type` int(10) unsigned NOT NULL COMMENT '平台类型',
  `reason` varchar(255) NOT NULL COMMENT '取消原因',
  `from_lat` varchar(255) NOT NULL DEFAULT '',
  `from_lng` varchar(255) NOT NULL DEFAULT '',
  `to_lat` varchar(255) NOT NULL DEFAULT '',
  `to_lng` varchar(255) NOT NULL DEFAULT '',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deduct_fee` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '取消订单违约金',
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `service_id` int(11) NOT NULL DEFAULT '0' COMMENT '服务人员id',
  `service_fee` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '服务费',
  PRIMARY KEY (`delivery_order_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_delivery_service`;
CREATE TABLE `qixi_delivery_service` (

  `service_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '配送员uid',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '类型：0平台1:商户',
  `relation_id` int(11) NOT NULL DEFAULT '0' COMMENT '门店、供应商id',
  `avatar` varchar(250) NOT NULL DEFAULT '' COMMENT '配送员头像',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '配送员名称',
  `phone` varchar(20) NOT NULL DEFAULT '0' COMMENT '手机号码',
  `create_time` int(11) NOT NULL DEFAULT '0' COMMENT '添加时间',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '0隐藏1显示',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户ID',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  PRIMARY KEY (`service_id`) USING BTREE,
  KEY `uid` (`uid`,`is_del`,`status`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_delivery_station`;
CREATE TABLE `qixi_delivery_station` (

  `station_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `station_name` varchar(255) NOT NULL DEFAULT '' COMMENT '门店名称',
  `business` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '支持配送的物品品类',
  `city_name` varchar(100) NOT NULL DEFAULT '' COMMENT '门店所属市',
  `station_address` varchar(255) NOT NULL DEFAULT '' COMMENT '门店地址',
  `lng` char(20) NOT NULL DEFAULT '' COMMENT '门店经度',
  `lat` char(20) NOT NULL DEFAULT '' COMMENT '门店纬度',
  `contact_name` char(10) NOT NULL DEFAULT '' COMMENT '联系人姓名',
  `phone` char(11) NOT NULL DEFAULT '' COMMENT '联系人电话',
  `origin_shop_id` varchar(255) NOT NULL DEFAULT '' COMMENT '门店编码,可自定义,但必须唯一;若不填写,则系统自动生成',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '达达商家app账号(若不需要登陆app,则不用设置)',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '达达商家app密码(若不需要登陆app,则不用设置)\n',
  `status` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '状态 1启用 0关闭',
  `mer_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '商户ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `type` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '类型 0 到店自提 1 达达 2 uu',
  `switch_city` tinyint(2) unsigned NOT NULL DEFAULT '1' COMMENT '同城配送：1 支持 0 不支持',
  `switch_take` tinyint(2) unsigned NOT NULL DEFAULT '0' COMMENT '到店自提：1 支持 0 不支持',
  `business_date` varchar(100) NOT NULL DEFAULT '' COMMENT '营业日期',
  `business_time_start` varchar(20) NOT NULL DEFAULT '' COMMENT '营业开始时间',
  `business_time_end` varchar(20) NOT NULL DEFAULT '' COMMENT '营业结束时间',
  `card_number` varchar(20) NOT NULL DEFAULT '' COMMENT '身份证号',
  `is_del` tinyint(2) unsigned NOT NULL DEFAULT '0',
  `range_type` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '距离设置类型(1:范围 2:行政区 3:电子围栏)',
  `radius` float(8,2) unsigned NOT NULL DEFAULT '1.00' COMMENT '服务半径(km)',
  `region` varchar(200) NOT NULL DEFAULT '' COMMENT '行政区域',
  `fence` text NOT NULL COMMENT '电子围栏配置',
  PRIMARY KEY (`station_id`),
  UNIQUE KEY `id` (`station_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='同城配送门店列表';

DROP TABLE IF EXISTS `qixi_diy`;
CREATE TABLE `qixi_diy` (

  `id` int(11) NOT NULL AUTO_INCREMENT,
  `version` varchar(255) NOT NULL DEFAULT '' COMMENT '版本号',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '页面名称',
  `title` varchar(100) NOT NULL DEFAULT '' COMMENT '网站标题',
  `cover_image` varchar(255) NOT NULL DEFAULT '' COMMENT '封面图',
  `template_name` varchar(255) NOT NULL DEFAULT '' COMMENT '模板名称',
  `default_value` longtext COMMENT '默认数据',
  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否使用',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '页面类型',
  `is_show` tinyint(1) NOT NULL DEFAULT '0' COMMENT '显示首页',
  `is_bg_color` tinyint(1) NOT NULL DEFAULT '0' COMMENT '颜色是否选中',
  `is_bg_pic` tinyint(1) NOT NULL DEFAULT '0' COMMENT '背景图是否选中',
  `is_diy` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是diy数据',
  `color_picker` varchar(50) NOT NULL DEFAULT '' COMMENT '背景颜色',
  `bg_pic` varchar(256) NOT NULL DEFAULT '' COMMENT '背景图',
  `bg_tab_val` tinyint(1) NOT NULL DEFAULT '0' COMMENT '背景图图片样式',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户ID',
  `is_default` tinyint(1) NOT NULL DEFAULT '0' COMMENT '默认模板(1.平台默认 2.商户默认）',
  `scope_type` tinyint(1) DEFAULT '4' COMMENT '适用范围类型：0.全部店铺、1. 指定店铺、2. 指定商户分类、3. 指定店铺类型、4. 指定商户类别',
  `value` longtext COMMENT '页面数据',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `template_name` (`template_name`,`type`) USING BTREE,
  KEY `status_type` (`status`,`type`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_excel`;
CREATE TABLE `qixi_excel` (

  `excel_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) DEFAULT NULL COMMENT '文件名',
  `status` int(11) DEFAULT '0' COMMENT '0.默认，1.完成，2.失败',
  `type` varchar(255) DEFAULT NULL COMMENT '类型',
  `path` varchar(255) DEFAULT NULL COMMENT '文件路径',
  `mer_id` int(11) DEFAULT '0' COMMENT '商户id',
  `admin_id` int(11) DEFAULT NULL COMMENT '操作者id',
  `is_del` int(11) DEFAULT '0' COMMENT '是否删除',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `message` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`excel_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='导出文件记录表';

DROP TABLE IF EXISTS `qixi_express`;
CREATE TABLE `qixi_express` (

  `id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT '快递公司id',
  `code` varchar(50) NOT NULL COMMENT '快递公司简称',
  `name` varchar(50) NOT NULL COMMENT '快递公司全称',
  `mark` varchar(255) DEFAULT '' COMMENT '备注',
  `partner_id` int(11) DEFAULT '0' COMMENT '月结账号',
  `partner_key` int(11) DEFAULT '0' COMMENT '月结密码',
  `net` int(11) DEFAULT '0' COMMENT '取件网点',
  `sort` int(11) NOT NULL COMMENT '排序',
  `is_show` int(11) DEFAULT '1' COMMENT '是否显示',
  `check_man` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否需要承载快递员名称',
  `partner_name` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否需要客户账户名称',
  `is_code` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否需要承载编号',
  `open_mer` varchar(255) NOT NULL DEFAULT '' COMMENT '开启该物流公司的商户',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `code` (`code`) USING BTREE,
  KEY `is_show` (`is_show`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='快递公司表';

DROP TABLE IF EXISTS `qixi_express_partner`;
CREATE TABLE `qixi_express_partner` (

  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `express_id` int(11) NOT NULL COMMENT '快递公司id',
  `account` varchar(20) DEFAULT NULL COMMENT '月结账号',
  `key` varchar(50) DEFAULT NULL COMMENT '月结密码',
  `net_name` varchar(50) DEFAULT NULL COMMENT '取件网点',
  `mer_id` int(11) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '1',
  `check_man` varchar(50) DEFAULT NULL COMMENT '承载快递员名',
  `partner_name` varchar(50) DEFAULT NULL COMMENT '客户账户名称',
  `code` varchar(50) DEFAULT NULL COMMENT '承载编号',
  PRIMARY KEY (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_extend`;
CREATE TABLE `qixi_extend` (

  `extend_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户 id',
  `link_id` int(10) unsigned NOT NULL COMMENT '关联字段',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `extend_type` varchar(32) NOT NULL COMMENT '扩展字段',
  `extend_value` varchar(255) NOT NULL COMMENT '扩展值',
  `update_time` datetime NOT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`extend_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `link_id` (`link_id`,`extend_type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_feedback`;
CREATE TABLE `qixi_feedback` (

  `feedback_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `type` varchar(255) NOT NULL,
  `content` varchar(512) NOT NULL,
  `images` text COMMENT '反馈图片',
  `realname` varchar(24) NOT NULL COMMENT '姓名',
  `contact` varchar(32) NOT NULL COMMENT '联系方式',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint(1) DEFAULT '0' COMMENT '状态',
  `reply` varchar(255) DEFAULT NULL COMMENT '回复，最终给用户的回复内容',
  `remake` varchar(255) DEFAULT NULL COMMENT '备注，后台人员自己查看用',
  `is_del` tinyint(1) NOT NULL DEFAULT '0',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '回复时间',
  PRIMARY KEY (`feedback_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户反馈表';

DROP TABLE IF EXISTS `qixi_feedback_category`;
CREATE TABLE `qixi_feedback_category` (

  `feedback_category_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '商品分类表ID',
  `pid` mediumint(9) NOT NULL COMMENT '父id',
  `cate_name` varchar(100) NOT NULL COMMENT '分类名称',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
  `sort` mediumint(9) NOT NULL COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `level` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`feedback_category_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户反馈分类表';

DROP TABLE IF EXISTS `qixi_financial`;
CREATE TABLE `qixi_financial` (

  `financial_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `financial_sn` varchar(32) NOT NULL COMMENT '单号',
  `mer_money` decimal(12,2) unsigned NOT NULL COMMENT '余额',
  `extract_money` decimal(12,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '提现金额',
  `financial_type` int(10) unsigned DEFAULT '0' COMMENT '收款类型',
  `financial_account` varchar(500) NOT NULL COMMENT '商户账户信息',
  `financial_status` int(10) unsigned DEFAULT '0' COMMENT '转账状态',
  `status` int(11) NOT NULL COMMENT '审核0待审核，1通过 ，-1 未通过',
  `refusal` varchar(32) DEFAULT NULL COMMENT '拒绝理由',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `image` varchar(1000) DEFAULT NULL COMMENT '凭证',
  `admin_id` int(11) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `status_time` timestamp NULL DEFAULT NULL COMMENT '审核时间',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '修改拼凭证时间',
  `is_del` int(10) unsigned DEFAULT '0',
  `mark` varchar(255) DEFAULT NULL COMMENT '商户备注',
  `admin_mark` varchar(255) DEFAULT NULL COMMENT '平台备注',
  `mer_admin_id` int(11) DEFAULT NULL COMMENT '商户管理员',
  `type` int(10) unsigned DEFAULT '0' COMMENT '申请类型 0.余额 1 保证金',
  PRIMARY KEY (`financial_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='商户财务申请提现';

DROP TABLE IF EXISTS `qixi_financial_record`;
CREATE TABLE `qixi_financial_record` (

  `financial_record_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `financial_record_sn` varchar(32) NOT NULL COMMENT '流水号',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单号',
  `order_sn` varchar(32) NOT NULL COMMENT '订单编号',
  `user_info` varchar(32) NOT NULL COMMENT '用户名',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户 id',
  `financial_type` varchar(32) NOT NULL COMMENT '流水类型',
  `financial_pm` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0 = 支出 1 = 获得',
  `number` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '金额',
  `type` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '0:商户  1:公共  2:平台',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `pay_type` int(11) NOT NULL COMMENT '支付类型',
  PRIMARY KEY (`financial_record_id`),
  KEY `mer_id` (`mer_id`),
  KEY `financial_type` (`financial_type`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户财务流水';

DROP TABLE IF EXISTS `qixi_guarantee`;
CREATE TABLE `qixi_guarantee` (

  `guarantee_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `guarantee_name` varchar(255) DEFAULT NULL COMMENT '保障服务名称',
  `guarantee_info` varchar(500) DEFAULT NULL COMMENT '保障服务简介',
  `status` int(11) DEFAULT '1' COMMENT '0.关闭，1开启',
  `image` varchar(255) DEFAULT NULL COMMENT '图标',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `mer_count` int(11) DEFAULT '0' COMMENT '使用的商户数',
  `product_cout` int(11) DEFAULT '0' COMMENT '使用的商品数',
  `is_del` int(11) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`guarantee_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='保障服务选项';

DROP TABLE IF EXISTS `qixi_guarantee_template`;
CREATE TABLE `qixi_guarantee_template` (

  `guarantee_template_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `template_name` varchar(255) DEFAULT NULL,
  `mer_id` int(11) DEFAULT NULL,
  `status` int(10) unsigned DEFAULT '1',
  `sort` int(11) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT NULL,
  `is_del` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`guarantee_template_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='保障服务模板';

DROP TABLE IF EXISTS `qixi_guarantee_value`;
CREATE TABLE `qixi_guarantee_value` (

  `guarantee_value_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `guarantee_id` int(10) unsigned DEFAULT NULL,
  `guarantee_template_id` int(10) unsigned DEFAULT NULL,
  `mer_id` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT '1',
  PRIMARY KEY (`guarantee_value_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='保障服务模板条款';

DROP TABLE IF EXISTS `qixi_label_rule`;
CREATE TABLE `qixi_label_rule` (

  `label_rule_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `label_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '标签 id',
  `type` tinyint(3) unsigned DEFAULT '0' COMMENT '0=订单数 1=订单金额',
  `min` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '最小值',
  `max` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT ' 最大值',
  `user_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户数',
  `update_time` timestamp NULL DEFAULT NULL COMMENT '更新时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`label_rule_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='自定标签规则';

DROP TABLE IF EXISTS `qixi_member_interests`;
CREATE TABLE `qixi_member_interests` (

  `interests_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(20) DEFAULT NULL COMMENT '名称',
  `info` varchar(200) DEFAULT NULL COMMENT '介绍',
  `brokerage_level` tinyint(3) unsigned DEFAULT NULL COMMENT '关联等级',
  `pic` varchar(128) DEFAULT NULL COMMENT '图标',
  `type` tinyint(4) DEFAULT '0' COMMENT '类型1.免费会员 2.付费会员',
  `link` varchar(500) DEFAULT NULL COMMENT '跳转 链接',
  `has_type` int(11) DEFAULT NULL COMMENT '特权类型',
  `value` varchar(500) DEFAULT NULL COMMENT '特权值',
  `on_pic` varchar(128) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL,
  PRIMARY KEY (`interests_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_merchant`;
CREATE TABLE `qixi_merchant` (

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
  `mark` varchar(256) NOT NULL COMMENT '商户备注',
  `reg_admin_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '总后台管理员ID',
  `sort` int(10) unsigned NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '商户是否禁用0锁定,1正常',
  `commission_rate` decimal(6,2) unsigned DEFAULT NULL COMMENT '提成比例',
  `commission_switch` int(11) DEFAULT '0' COMMENT '商户手续费单独设置 0 关闭 1 开启',
  `long` varchar(16) DEFAULT NULL COMMENT '经度',
  `lat` varchar(16) DEFAULT NULL COMMENT '纬度',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0未删除1删除',
  `is_audit` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '添加的产品是否审核0不审核1审核',
  `is_bro_room` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否审核直播间0不审核1审核',
  `is_bro_goods` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否审核直播商品0不审核1审核',
  `is_best` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否推荐',
  `is_trader` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否自营',
  `mer_state` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '商户是否1开启0关闭',
  `mer_info` varchar(256) NOT NULL DEFAULT '' COMMENT '店铺简介',
  `service_phone` varchar(13) NOT NULL DEFAULT '' COMMENT '店铺电话',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `care_count` int(10) unsigned DEFAULT '0' COMMENT '关注总数',
  `copy_product_num` int(10) unsigned DEFAULT '0' COMMENT '剩余复制商品次数',
  `export_dump_num` int(10) unsigned DEFAULT '0' COMMENT '电子面单剩余次数',
  `mer_money` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '商户余额',
  `financial_bank` varchar(255) DEFAULT NULL COMMENT '银行卡转账信息',
  `financial_wechat` varchar(255) DEFAULT NULL COMMENT '微信转账信息',
  `financial_alipay` varchar(255) DEFAULT NULL COMMENT '支付宝转账信息',
  `financial_type` tinyint(3) unsigned DEFAULT '1' COMMENT '默认使用类型',
  `sub_mchid` varchar(16) NOT NULL DEFAULT '' COMMENT '微信支付分配的分账号',
  `delivery_way` varchar(50) DEFAULT '' COMMENT '配送方式',
  `delivery_balance` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '配送余额',
  `margin` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '保证金',
  `margin_remind_time` varchar(255) DEFAULT NULL COMMENT '保证金补缴提醒结束时间，时间点到了就自动关闭店铺',
  `ot_margin` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '保证金额度',
  `is_margin` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否有保证金（0无，1有未支付，10已支付，-1 申请退款, -10 拒绝退款）',
  `offline_switch` tinyint(4) DEFAULT '0' COMMENT '线下支付功能开关',
  `care_ficti` int(11) DEFAULT '0' COMMENT '虚拟关注量',
  `region_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户所属分组',
  `applyment_id` varchar(50) NOT NULL DEFAULT '' COMMENT '特约商户ID',
  `business_id` int(10) NOT NULL DEFAULT '0' COMMENT '店铺所属商户id',
  `applyment_switch` int(11) NOT NULL DEFAULT '1' COMMENT '特约商户设置的分账比例是否合理',
  PRIMARY KEY (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='商户表';

DROP TABLE IF EXISTS `qixi_merchant_admin`;
CREATE TABLE `qixi_merchant_admin` (

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

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='商户管理员表';

DROP TABLE IF EXISTS `qixi_merchant_applyments`;
CREATE TABLE `qixi_merchant_applyments` (

  `mer_applyments_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `out_request_no` varchar(128) DEFAULT NULL COMMENT '业务申请编号',
  `applyment_id` varchar(100) DEFAULT NULL COMMENT '微信支付分配的申请单号',
  `mer_id` int(11) DEFAULT '0' COMMENT '商户ID',
  `sub_mchid` varchar(100) DEFAULT NULL COMMENT '二级商户号/特约商户号',
  `mer_name` varchar(50) DEFAULT NULL COMMENT '商户名',
  `info` text COMMENT '申请资料',
  `status` int(11) DEFAULT '0' COMMENT '申请状态: 0.平台未提交，-1.平台驳回，10.平台提交审核中，11.需用户操作 ，20.已完成，30.已冻结，40.驳回',
  `message` varchar(1000) DEFAULT NULL COMMENT '返回信息',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` int(11) DEFAULT '0' COMMENT '删除',
  `type` tinyint(2) NOT NULL DEFAULT '0' COMMENT '0 平台收付通 1 服务商分账',
  PRIMARY KEY (`mer_applyments_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户申请分账商户号表';

DROP TABLE IF EXISTS `qixi_merchant_category`;
CREATE TABLE `qixi_merchant_category` (

  `merchant_category_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户分类 id',
  `commission_rate` decimal(6,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '手续费',
  `category_name` varchar(32) NOT NULL COMMENT '商户分类名称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`merchant_category_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='商户分类表';

DROP TABLE IF EXISTS `qixi_merchant_intention`;
CREATE TABLE `qixi_merchant_intention` (

  `mer_intention_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `uid` int(10) unsigned DEFAULT '0' COMMENT '用户ID',
  `phone` varchar(11) DEFAULT NULL COMMENT '手机号',
  `mer_name` varchar(30) DEFAULT NULL COMMENT '商户名称',
  `name` varchar(30) DEFAULT NULL COMMENT '客户姓名',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '提交时间',
  `status` tinyint(4) DEFAULT '0' COMMENT '处理状态 1通过 ，2未通过',
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

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户申请表';

DROP TABLE IF EXISTS `qixi_merchant_region`;
CREATE TABLE `qixi_merchant_region` (

  `region_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `info` varchar(255) NOT NULL DEFAULT '' COMMENT '简介',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级ID',
  `path` varchar(255) NOT NULL DEFAULT '/' COMMENT '父级路径',
  `pic` varchar(255) NOT NULL DEFAULT '' COMMENT '图片',
  `lv` int(11) NOT NULL DEFAULT '0' COMMENT '等级',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0 关闭, 1 开启',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '类型:0',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`region_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_merchant_type`;
CREATE TABLE `qixi_merchant_type` (

  `mer_type_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商户类型 id',
  `type_name` varchar(16) NOT NULL COMMENT '类型名称',
  `type_info` varchar(512) DEFAULT NULL COMMENT '类型要求',
  `description` varchar(512) DEFAULT NULL COMMENT '类型说明',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `margin` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '保证金',
  `is_margin` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否有保证金（0无，1有）',
  `mark` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`mer_type_id`)

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='商户类型表';

DROP TABLE IF EXISTS `qixi_open_auth`;
CREATE TABLE `qixi_open_auth` (

  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(50) DEFAULT NULL COMMENT '标题',
  `access_key` varchar(50) DEFAULT NULL,
  `secret_key` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT NULL COMMENT '状态',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `auth` varchar(255) DEFAULT NULL COMMENT '权限',
  `sort` int(11) DEFAULT NULL,
  `is_del` tinyint(4) DEFAULT '0' COMMENT '是否删除',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `last_ip` varchar(50) DEFAULT NULL COMMENT '最后登录的IP',
  `last_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后登录的时间',
  PRIMARY KEY (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_operate_log`;
CREATE TABLE `qixi_operate_log` (

  `operate_log_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `title` varchar(255) DEFAULT NULL COMMENT '标题',
  `relevance_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联id',
  `relevance_title` varchar(255) DEFAULT NULL COMMENT '关联标题',
  `relevance_type` varchar(255) NOT NULL COMMENT '关联类型',
  `type` enum('1','2') NOT NULL DEFAULT '1' COMMENT '1|平台 2|商户',
  `category` varchar(255) NOT NULL COMMENT '类别',
  `action` varchar(255) NOT NULL DEFAULT '' COMMENT '操作类型',
  `operator_role_id` int(11) DEFAULT NULL COMMENT '操作角色id',
  `operator_role_nickname` varchar(255) DEFAULT NULL COMMENT '操作角色昵称',
  `operator_uid` varchar(255) NOT NULL COMMENT '操作用户id',
  `operator_nickname` varchar(255) DEFAULT NULL COMMENT '操作用户昵称',
  `mark` varchar(2000) DEFAULT NULL COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`operate_log_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_page_category`;
CREATE TABLE `qixi_page_category` (

  `id` int(10) NOT NULL AUTO_INCREMENT,
  `pid` int(10) NOT NULL DEFAULT '0' COMMENT '父类id',
  `type` varchar(50) NOT NULL DEFAULT 'link' COMMENT '类型:link、special、product、product_category、custom',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '分类名称',
  `sort` smallint(5) NOT NULL DEFAULT '0' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `level` tinyint(1) NOT NULL DEFAULT '0',
  `is_mer` int(10) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)

) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=utf8 COMMENT='页面链接分类';

DROP TABLE IF EXISTS `qixi_page_link`;
CREATE TABLE `qixi_page_link` (

  `id` int(10) NOT NULL AUTO_INCREMENT,
  `cate_id` int(10) NOT NULL DEFAULT '0' COMMENT '分类id',
  `type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '分组1:基础2:分销3:个人中心',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '页面名称',
  `url` varchar(255) NOT NULL DEFAULT '' COMMENT '页面链接',
  `param` varchar(255) NOT NULL DEFAULT '' COMMENT '参数',
  `example` varchar(255) NOT NULL DEFAULT '' COMMENT '事例',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态',
  `sort` smallint(5) NOT NULL DEFAULT '0' COMMENT '排序',
  `add_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_mer` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1是商户的链接',
  PRIMARY KEY (`id`)

) ENGINE=InnoDB AUTO_INCREMENT=81 DEFAULT CHARSET=utf8 COMMENT='页面链接';

DROP TABLE IF EXISTS `qixi_parameter`;
CREATE TABLE `qixi_parameter` (

  `parameter_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `template_id` int(10) unsigned NOT NULL COMMENT '模板 id',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户 id',
  `name` varchar(32) NOT NULL COMMENT '参数名称',
  `value` varchar(255) DEFAULT NULL COMMENT '参数值',
  `required` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否必填',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `sort` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`parameter_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `template_id` (`template_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_parameter_product`;
CREATE TABLE `qixi_parameter_product` (

  `id` int(11) NOT NULL AUTO_INCREMENT,
  `parameter_value_id` int(11) NOT NULL DEFAULT '0' COMMENT '参数ID',
  `product_id` int(11) NOT NULL DEFAULT '0' COMMENT '商品ID',
  PRIMARY KEY (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_parameter_template`;
CREATE TABLE `qixi_parameter_template` (

  `template_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户id',
  `template_name` varchar(64) NOT NULL COMMENT '模板名称',
  `sort` int(10) unsigned DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`template_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_parameter_value`;
CREATE TABLE `qixi_parameter_value` (

  `parameter_value_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `parameter_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '关联参数 id',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品 id',
  `name` varchar(64) DEFAULT NULL COMMENT '参数名称',
  `value` varchar(64) DEFAULT NULL COMMENT '参数值',
  `sort` int(10) unsigned DEFAULT '0' COMMENT '排序',
  `create_time` datetime DEFAULT NULL,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`parameter_value_id`) USING BTREE,
  KEY `parameter_id` (`parameter_id`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_presell_order`;
CREATE TABLE `qixi_presell_order` (

  `presell_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '预售尾款订单id',
  `presell_order_sn` varchar(32) NOT NULL COMMENT '预售订单号',
  `uid` int(10) unsigned NOT NULL COMMENT '用户 id',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单id',
  `transaction_id` varchar(60) DEFAULT NULL COMMENT '微信支付订单号(分账时有效)',
  `final_start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '支付开始时间',
  `final_end_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '尾款支付结时间',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:未支付 1:已支付',
  `status` tinyint(3) unsigned DEFAULT '1' COMMENT '0:无效 1:有效',
  `pay_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '支付方式 0余额 1微信 2小程序 3,4支付宝',
  `pay_price` decimal(8,2) unsigned NOT NULL COMMENT '尾款',
  `refun_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '退款金额',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_combine` tinyint(3) unsigned DEFAULT '0' COMMENT '是否为合并支付',
  PRIMARY KEY (`presell_order_id`) USING BTREE,
  UNIQUE KEY `order_id` (`order_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_price_rule`;
CREATE TABLE `qixi_price_rule` (

  `rule_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `rule_name` varchar(64) NOT NULL COMMENT '名称',
  `sort` int(10) unsigned DEFAULT '0' COMMENT '排序',
  `is_show` tinyint(3) unsigned DEFAULT '1' COMMENT '是否显示',
  `is_default` tinyint(3) unsigned DEFAULT '0' COMMENT '是否默认',
  `content` longtext COMMENT '内容',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`rule_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_record`;
CREATE TABLE `qixi_record` (

  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(20) NOT NULL DEFAULT '',
  `uid` int(11) NOT NULL DEFAULT '0',
  `link_id` int(11) NOT NULL DEFAULT '0',
  `num` int(11) NOT NULL DEFAULT '0',
  `title` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_relevance`;
CREATE TABLE `qixi_relevance` (

  `relevance_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `left_id` int(10) unsigned NOT NULL,
  `right_id` int(10) unsigned NOT NULL,
  `type` varchar(32) NOT NULL DEFAULT '',
  PRIMARY KEY (`relevance_id`) USING BTREE,
  KEY `type` (`type`,`left_id`,`right_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=1380 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_routine_qrcode`;
CREATE TABLE `qixi_routine_qrcode` (

  `routine_qrcode_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '微信二维码ID',
  `third_type` varchar(32) NOT NULL COMMENT '二维码类型 spread(用户推广) product_spread(商品推广)',
  `third_id` int(10) unsigned NOT NULL COMMENT '用户id',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态 0不可用 1可用',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `page` varchar(255) DEFAULT NULL COMMENT '小程序页面路径带参数',
  `qrcode_url` varchar(255) DEFAULT NULL COMMENT '小程序二维码路径',
  `url_time` timestamp NULL DEFAULT NULL COMMENT '二维码添加时间',
  PRIMARY KEY (`routine_qrcode_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='小程序二维码管理表';

DROP TABLE IF EXISTS `qixi_serve_meal`;
CREATE TABLE `qixi_serve_meal` (

  `meal_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(30) DEFAULT NULL COMMENT '套餐名称',
  `type` int(11) DEFAULT '0' COMMENT '套餐类型,1复制商品，2电子面单',
  `price` decimal(8,2) DEFAULT '0.00' COMMENT '价格',
  `num` int(11) DEFAULT '1' COMMENT '数量',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `status` int(11) DEFAULT '1' COMMENT '状态',
  `is_del` int(11) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`meal_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_serve_order`;
CREATE TABLE `qixi_serve_order` (

  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `meal_id` int(11) DEFAULT NULL COMMENT '套餐ID',
  `pay_type` int(11) DEFAULT NULL COMMENT '支付方式：1微信，2支付宝,3 平台操作',
  `order_sn` varchar(50) DEFAULT NULL COMMENT '订单ID',
  `pay_price` decimal(8,2) DEFAULT NULL COMMENT '价格',
  `order_info` varchar(255) DEFAULT NULL COMMENT '套餐信息',
  `type` int(11) DEFAULT NULL COMMENT '套餐类型 1 采集 2 电子面单 10 保证金 20同城配送 30 会员充值',
  `status` int(11) DEFAULT '0' COMMENT '状态：默认0，支付成功 1，支付失败 -1，20已退款',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID/用户ID',
  `create_time` timestamp NULL DEFAULT NULL,
  `is_del` int(11) DEFAULT '0',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
  PRIMARY KEY (`order_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_shipping_template`;
CREATE TABLE `qixi_shipping_template` (

  `shipping_template_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `name` varchar(255) NOT NULL COMMENT '模板名称',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '计费方式 0=数量 1=重量 2=体积',
  `appoint` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '开启指定包邮',
  `undelivery` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '开启指定区域不配送',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `is_default` tinyint(3) unsigned DEFAULT '0' COMMENT '默认模板',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `info` varchar(1000) DEFAULT NULL COMMENT '运费说明',
  PRIMARY KEY (`shipping_template_id`) USING BTREE,
  KEY `mer_id` (`mer_id`,`sort`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='运费表';

DROP TABLE IF EXISTS `qixi_shipping_template_free`;
CREATE TABLE `qixi_shipping_template_free` (

  `shipping_template_free_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `temp_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '模板ID',
  `city_id` text NOT NULL COMMENT '城市ID /id/id/id/id/',
  `number` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '包邮件数',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '包邮金额',
  PRIMARY KEY (`shipping_template_free_id`) USING BTREE,
  KEY `temp_id` (`temp_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='指定包邮信息表';

DROP TABLE IF EXISTS `qixi_shipping_template_region`;
CREATE TABLE `qixi_shipping_template_region` (

  `shipping_template_region_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `temp_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '模板ID',
  `city_id` text NOT NULL COMMENT '城市ID /id/id/id/',
  `first` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '首件',
  `first_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '首件运费',
  `continue` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '续件',
  `continue_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '续件运费',
  PRIMARY KEY (`shipping_template_region_id`) USING BTREE,
  KEY `temp_id` (`temp_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='配送区域表';

DROP TABLE IF EXISTS `qixi_shipping_template_undelivery`;
CREATE TABLE `qixi_shipping_template_undelivery` (

  `shipping_template_undelivery_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '编号',
  `temp_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '模板ID',
  `city_id` text NOT NULL COMMENT '城市ID /id/id/id/',
  PRIMARY KEY (`shipping_template_undelivery_id`) USING BTREE,
  KEY `temp_id` (`temp_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='指定不配送区域表';

DROP TABLE IF EXISTS `qixi_sms_record`;
CREATE TABLE `qixi_sms_record` (

  `sms_record_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '短信发送记录编号',
  `uid` varchar(255) NOT NULL COMMENT '短信平台账号',
  `phone` char(11) NOT NULL COMMENT '接受短信的手机号',
  `content` text NOT NULL COMMENT '短信内容',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送短信时间',
  `ip` varchar(16) NOT NULL DEFAULT '' COMMENT '添加记录ip',
  `template` varchar(255) NOT NULL COMMENT '短信模板ID',
  `resultcode` int(10) unsigned DEFAULT NULL COMMENT '状态码 100=成功,130=失败,131=空号,132=停机,133=关机,134=无状态',
  `record_id` int(10) unsigned NOT NULL COMMENT '发送记录id',
  PRIMARY KEY (`sms_record_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='短信发送记录表';

DROP TABLE IF EXISTS `qixi_staffs`;
CREATE TABLE `qixi_staffs` (

  `staffs_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '员工id',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `uid` int(11) NOT NULL COMMENT '关联用户uid',
  `photo` varchar(250) NOT NULL COMMENT '证件照',
  `name` varchar(50) NOT NULL COMMENT '姓名',
  `phone` varchar(18) DEFAULT '' COMMENT '电话',
  `remark` varchar(250) DEFAULT '' COMMENT '备注',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '状态：0、关闭；1、开启；',
  `sort` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `delete_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '删除时间',
  PRIMARY KEY (`staffs_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_activity`;
CREATE TABLE `qixi_store_activity` (

  `activity_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `activity_name` varchar(128) NOT NULL DEFAULT '' COMMENT '活动名称',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `pic` varchar(128) DEFAULT '' COMMENT '图片',
  `is_show` tinyint(1) DEFAULT '0' COMMENT '是否显示',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态0未开始，1进行中，2已结束',
  `activity_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1.氛围图 2.边框',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `scope_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '指定类型：0全部商品1指定商品2指定分类3指定商户4秒杀活动',
  `images` varchar(500) DEFAULT '' COMMENT '多图',
  `info` varchar(500) DEFAULT '' COMMENT '简介',
  `color` varchar(128) DEFAULT '' COMMENT '背景色',
  `sort` int(11) DEFAULT '0' COMMENT '排序',
  `mer_id` int(11) DEFAULT '0' COMMENT '商户ID',
  `link_id` int(11) DEFAULT NULL COMMENT '关联ID',
  `update_time` timestamp NULL DEFAULT NULL,
  `count` int(10) unsigned DEFAULT '0' COMMENT '需要的总数',
  `total` int(10) unsigned DEFAULT '0' COMMENT '已有的总数',
  `is_display` tinyint(1) DEFAULT '1' COMMENT '是否在活动列表中显示',
  PRIMARY KEY (`activity_id`) USING BTREE,
  KEY `start_time` (`start_time`,`end_time`) USING BTREE,
  KEY `activity_type` (`activity_type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_activity_cate`;
CREATE TABLE `qixi_store_activity_cate` (

  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `pid` int(11) NOT NULL DEFAULT '0',
  `pic` varchar(255) NOT NULL DEFAULT '',
  `sort` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(2) NOT NULL DEFAULT '0',
  `mer_id` int(11) NOT NULL DEFAULT '0',
  `type` tinyint(2) NOT NULL DEFAULT '0',
  `path` varchar(255) NOT NULL DEFAULT '/',
  `lv` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动标签分类表';

DROP TABLE IF EXISTS `qixi_store_activity_label`;
CREATE TABLE `qixi_store_activity_label` (

  `id` int(11) NOT NULL AUTO_INCREMENT,
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '类型：0平台 2:商户',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `label_cate` int(11) NOT NULL DEFAULT '0' COMMENT '标签分类',
  `label_name` varchar(255) NOT NULL DEFAULT '' COMMENT '标签名称',
  `style_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '样式类型1：自定义2：图片',
  `color` varchar(32) NOT NULL DEFAULT '' COMMENT '颜色',
  `bg_color` varchar(32) NOT NULL DEFAULT '' COMMENT '背景颜色',
  `border_color` varchar(32) NOT NULL DEFAULT '' COMMENT '边框颜色',
  `icon` varchar(255) DEFAULT '' COMMENT '图标',
  `default_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '系统默认标签分类(1:包邮,2:领券,3:上门,4:到店,5:同城,6:拼团,7:秒杀,8:助力,9:预售,10:自营)',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '移动端是否展示',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态是否开启',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `add_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`id`),
  KEY `label_cate` (`label_cate`),
  KEY `type` (`type`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='活动标签表';

DROP TABLE IF EXISTS `qixi_store_activity_related`;
CREATE TABLE `qixi_store_activity_related` (

  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` int(11) NOT NULL DEFAULT '0' COMMENT '活动ID',
  `activity_type` varchar(255) DEFAULT NULL COMMENT '活动类型',
  `keys` varchar(2000) DEFAULT NULL COMMENT '主要信息',
  `value` varchar(2000) DEFAULT NULL COMMENT '活动值',
  `form_value` text COMMENT '表单内容',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `link_id` int(11) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `nickname` varchar(255) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除 1|是 0|否',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `join_id` (`id`) USING BTREE,
  KEY `activity_id` (`activity_id`,`activity_type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_attr_template`;
CREATE TABLE `qixi_store_attr_template` (

  `attr_template_id` int(11) NOT NULL AUTO_INCREMENT,
  `template_name` varchar(32) NOT NULL COMMENT '规格名称',
  `template_value` text NOT NULL COMMENT '规格值',
  `mer_id` int(11) NOT NULL COMMENT '商户 id',
  PRIMARY KEY (`attr_template_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品规则值(规格)表';

DROP TABLE IF EXISTS `qixi_store_brand`;
CREATE TABLE `qixi_store_brand` (

  `brand_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '商品品牌表ID',
  `brand_category_id` mediumint(9) NOT NULL COMMENT '父id',
  `brand_name` varchar(100) NOT NULL COMMENT '品牌名称',
  `sort` mediumint(9) NOT NULL COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`brand_id`) USING BTREE,
  KEY `pid` (`brand_category_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品品牌表';

DROP TABLE IF EXISTS `qixi_store_brand_category`;
CREATE TABLE `qixi_store_brand_category` (

  `store_brand_category_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '品牌分类表ID',
  `pid` mediumint(9) NOT NULL COMMENT '父id',
  `cate_name` varchar(100) NOT NULL COMMENT '分类名称',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
  `sort` mediumint(9) NOT NULL COMMENT '排序',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `level` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`store_brand_category_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='品牌分类表';

DROP TABLE IF EXISTS `qixi_store_cart`;
CREATE TABLE `qixi_store_cart` (

  `cart_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车表ID',
  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `product_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '类型 0=普通产品，2.预售商品',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `product_attr_unique` varchar(16) NOT NULL DEFAULT '' COMMENT '商品属性',
  `cart_num` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '商品数量',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `source` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '来源 1.直播间,2.预售商品,3.助力商品',
  `source_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '来源关联 id',
  `is_pay` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 = 未购买 1 = 已购买',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `is_new` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为立即购买',
  `is_fail` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否失效',
  `spread_id` int(10) unsigned DEFAULT '0' COMMENT '推广人',
  `tourist_unique_key` varchar(20) NOT NULL DEFAULT '' COMMENT '游客唯一标识',
  `reservation_date` varchar(20) DEFAULT '' COMMENT '预约商品的预约日期',
  `reservation_id` int(10) unsigned DEFAULT '0' COMMENT '预约商品的预约时间段ID',
  PRIMARY KEY (`cart_id`) USING BTREE,
  KEY `user_id` (`uid`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='购物车表';

DROP TABLE IF EXISTS `qixi_store_cart_price`;
CREATE TABLE `qixi_store_cart_price` (

  `cart_price_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '购物车价格表ID',
  `cart_id` int(10) unsigned NOT NULL COMMENT '购物车id',
  `old_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '改价前价格',
  `type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '改价类型(0:一口价，1:立减金额，2:折扣率)',
  `reduce_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '立减金额（type为1时）',
  `discount_rate` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '折扣率（type为2时）',
  `new_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '改价后价格(type为0时一口价，或type不为0计算后的金额)',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `is_batch` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为批量改价(0:否，1:是)',
  PRIMARY KEY (`cart_price_id`) USING BTREE,
  KEY `cart_id` (`cart_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_category`;
CREATE TABLE `qixi_store_category` (

  `store_category_id` mediumint(9) NOT NULL AUTO_INCREMENT COMMENT '商品分类表ID',
  `pid` mediumint(9) NOT NULL COMMENT '父id',
  `cate_name` varchar(100) NOT NULL COMMENT '分类名称',
  `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
  `sort` mediumint(9) NOT NULL COMMENT '排序',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示',
  `level` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '等级',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_hot` tinyint(1) DEFAULT '0' COMMENT '是否推荐',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0 商品，1 积分商品',
  PRIMARY KEY (`store_category_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8 COMMENT='商品分类表';

DROP TABLE IF EXISTS `qixi_store_coupon`;
CREATE TABLE `qixi_store_coupon` (

  `coupon_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '优惠券表ID',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `is_timeout` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否限时',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '优惠券领取开启时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '优惠券领取结束时间',
  `is_limited` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否限量',
  `total_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券领取数量',
  `remain_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券剩余领取数量',
  `send_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0=领取 1=消费满赠 2=新人 3=买增 4=首单赠送',
  `full_reduction` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '消费满多少赠送优惠券',
  `title` varchar(64) NOT NULL COMMENT '优惠券名称',
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '优惠券面值',
  `use_min_price` int(11) NOT NULL DEFAULT '0' COMMENT '最低消费多少金额可用优惠券',
  `coupon_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券类型 0=有效天数 1=固定时间段',
  `coupon_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券有效期限（单位：天）',
  `use_start_time` timestamp NULL DEFAULT NULL COMMENT '开始时间',
  `use_end_time` timestamp NULL DEFAULT NULL COMMENT '到期时间',
  `sort` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '排序',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：关闭，1：开启 -1: 失效）',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '优惠券类型 0-店铺 1-商品券 10 平台通用券 11平台品类券 12 平台跨店券',
  PRIMARY KEY (`coupon_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='优惠券表';

DROP TABLE IF EXISTS `qixi_store_coupon_issue_user`;
CREATE TABLE `qixi_store_coupon_issue_user` (

  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '领取优惠券用户ID',
  `coupon_id` int(11) NOT NULL DEFAULT '0' COMMENT '优惠券ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
  KEY `uid` (`uid`,`coupon_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='优惠券前台用户领取记录表';

DROP TABLE IF EXISTS `qixi_store_coupon_product`;
CREATE TABLE `qixi_store_coupon_product` (

  `product_id` int(11) NOT NULL DEFAULT '0' COMMENT '产品id',
  `coupon_id` int(11) NOT NULL DEFAULT '0' COMMENT '优惠卷id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间'

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='优惠卷关联商品辅助表';

DROP TABLE IF EXISTS `qixi_store_coupon_send`;
CREATE TABLE `qixi_store_coupon_send` (

  `coupon_send_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `coupon_id` int(10) unsigned NOT NULL COMMENT '优惠券 id',
  `coupon_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发送数量',
  `mark` varchar(512) NOT NULL COMMENT '发送群体',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:发送中 1:全部发送',
  PRIMARY KEY (`coupon_send_id`),
  KEY `mer_id` (`mer_id`),
  KEY `coupon_id` (`coupon_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='优惠券发送记录';

DROP TABLE IF EXISTS `qixi_store_coupon_user`;
CREATE TABLE `qixi_store_coupon_user` (

  `coupon_user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '优惠券发放记录id',
  `coupon_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '兑换的项目id',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '优惠券所属用户',
  `coupon_title` varchar(32) NOT NULL COMMENT '优惠券名称',
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '优惠券的面值',
  `use_min_price` int(11) NOT NULL DEFAULT '0' COMMENT '最低消费多少金额可用优惠券',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '优惠券创建时间',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '优惠券开启时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '优惠券结束时间',
  `use_time` timestamp NULL DEFAULT NULL COMMENT '使用时间',
  `type` varchar(16) NOT NULL DEFAULT 'send' COMMENT '获取方式(receive:自己领取 send:后台发送  give:满赠  new:新人 buy:买赠送)',
  `send_id` int(10) unsigned DEFAULT '0' COMMENT '批量发送 id',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态（0：未使用，1：已使用, 2:已过期）',
  `is_fail` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否有效',
  PRIMARY KEY (`coupon_user_id`) USING BTREE,
  KEY `coupon_id` (`coupon_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `type` (`type`,`send_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='优惠券发放记录表';

DROP TABLE IF EXISTS `qixi_store_discounts`;
CREATE TABLE `qixi_store_discounts` (

  `discount_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '套餐名称',
  `image` varchar(500) NOT NULL DEFAULT '' COMMENT '组合套餐主图',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '套餐类型0固定1搭配',
  `is_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否限量0不限量1限量',
  `limit_num` int(11) NOT NULL DEFAULT '0' COMMENT '限量个数',
  `link_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '关联标签',
  `product_ids` varchar(255) DEFAULT '' COMMENT '商品IDS',
  `is_time` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否限时0不限时1限时',
  `start_time` int(11) NOT NULL DEFAULT '0' COMMENT '开始时间',
  `stop_time` int(11) NOT NULL DEFAULT '0' COMMENT '结束时间',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `free_shipping` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否包邮0不包邮1包邮',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '平台是否上架0不上架1上架',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除0未删除1已删除',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `is_show` tinyint(1) DEFAULT '0' COMMENT '商户是否上架0不上架1上架',
  `sales` int(10) unsigned DEFAULT NULL COMMENT '销量',
  PRIMARY KEY (`discount_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_discounts_product`;
CREATE TABLE `qixi_store_discounts_product` (

  `discount_product_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增ID',
  `discount_id` int(11) NOT NULL COMMENT '优惠套餐ID',
  `product_id` int(11) NOT NULL COMMENT '商品ID',
  `store_name` varchar(255) NOT NULL COMMENT '商品名称',
  `image` varchar(500) NOT NULL DEFAULT '' COMMENT '商品图',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否搭配0不是1是',
  `temp_id` int(11) NOT NULL DEFAULT '0' COMMENT '运费模版Id',
  `mer_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`discount_product_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_group`;
CREATE TABLE `qixi_store_group` (

  `store_group_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '分组id',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上级商圈id',
  `path` varchar(50) NOT NULL DEFAULT '' COMMENT '路径',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '等级:0一级 1二级 2三级',
  `positioning_status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否开启定位：0否 1是',
  `longitude` varchar(16) NOT NULL DEFAULT '' COMMENT '经度，positioning_status 为1时有效',
  `latitude` varchar(16) NOT NULL DEFAULT '' COMMENT '纬度，positioning_status 为1时有效',
  `address` varchar(100) NOT NULL DEFAULT '' COMMENT '中心点地址，positioning_status 为1时有效',
  `diy_temp_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '首页模板id',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '说明信息',
  `sort` int(10) NOT NULL DEFAULT '0' COMMENT '排序(数字越大越靠前)',
  `status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态：0关闭 1开启',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`store_group_id`) USING BTREE,
  KEY `idx_pid` (`pid`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='分组表';

DROP TABLE IF EXISTS `qixi_store_group_order`;
CREATE TABLE `qixi_store_group_order` (

  `group_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `group_order_sn` varchar(32) NOT NULL COMMENT '订单号',
  `uid` int(10) unsigned NOT NULL COMMENT '用户 ID',
  `total_postage` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '邮费',
  `total_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '订单总额',
  `total_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品数',
  `integral` int(10) unsigned DEFAULT '0' COMMENT '使用积分数量',
  `integral_price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '积分抵扣金额',
  `give_integral` int(10) unsigned DEFAULT '0' COMMENT '赠送积分',
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '优惠金额',
  `real_name` varchar(32) NOT NULL COMMENT '联系人',
  `user_phone` varchar(18) NOT NULL COMMENT '联系电话',
  `user_address` varchar(128) NOT NULL COMMENT '收货地址',
  `pay_price` decimal(8,2) unsigned NOT NULL COMMENT '支付金额',
  `pay_postage` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '支付邮费',
  `cost` decimal(8,2) unsigned NOT NULL COMMENT '成本价',
  `coupon_id` varchar(128) DEFAULT NULL COMMENT ' 平台优惠券',
  `give_coupon_ids` varchar(500) DEFAULT '' COMMENT '赠送优惠券',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否支付',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
  `pay_type` tinyint(1) NOT NULL COMMENT '支付方式 0=余额 1=微信 2=小程序 3=h5',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_remind` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否提醒',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_combine` tinyint(3) unsigned DEFAULT '0' COMMENT '是否为合并支付 ',
  `activity_type` tinyint(3) unsigned DEFAULT '0',
  `is_first` tinyint(1) DEFAULT '0' COMMENT '是否为用户首单',
  `is_behalf` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为代客下单订单',
  PRIMARY KEY (`group_order_id`) USING BTREE,
  UNIQUE KEY `group_order_id` (`group_order_sn`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `paid` (`paid`) USING BTREE,
  KEY `create_time` (`create_time`,`paid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户订单表';

DROP TABLE IF EXISTS `qixi_store_import`;
CREATE TABLE `qixi_store_import` (

  `import_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `import_type` varchar(20) DEFAULT NULL COMMENT 'delivery发货单',
  `type` int(11) DEFAULT '1' COMMENT '类型：1发货，2送货，3虚拟，4电子面单',
  `count` int(11) DEFAULT NULL COMMENT '总数',
  `success` int(11) DEFAULT NULL COMMENT '成功数',
  `status` int(11) DEFAULT '0' COMMENT '0.处理中，1成功，10部分完成，-1失败',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `mer_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`import_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='导入批次记录';

DROP TABLE IF EXISTS `qixi_store_import_delivery`;
CREATE TABLE `qixi_store_import_delivery` (

  `import_delivery_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `import_id` int(10) unsigned NOT NULL,
  `order_sn` varchar(32) DEFAULT NULL COMMENT '订单sn',
  `delivery_type` int(11) DEFAULT '1' COMMENT '类型：1发货，2送货，3虚拟，4电子面单',
  `delivery_name` varchar(64) DEFAULT NULL COMMENT '快递公司',
  `delivery_id` varchar(64) DEFAULT NULL COMMENT '快递单号',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `mer_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`import_delivery_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='导入发货单详细记录';

DROP TABLE IF EXISTS `qixi_store_order`;
CREATE TABLE `qixi_store_order` (

  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单ID',
  `main_id` int(10) unsigned DEFAULT '0' COMMENT '拆单前 id',
  `group_order_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单组 id',
  `order_sn` varchar(32) NOT NULL COMMENT '订单号',
  `uid` int(10) unsigned NOT NULL COMMENT '用户id',
  `spread_uid` int(10) unsigned DEFAULT '0' COMMENT '推荐人id',
  `top_uid` int(10) unsigned DEFAULT '0' COMMENT '二级推荐人 id',
  `real_name` varchar(32) NOT NULL COMMENT '用户姓名',
  `user_phone` varchar(18) NOT NULL COMMENT '用户电话',
  `user_address` varchar(128) NOT NULL COMMENT '详细地址',
  `cart_id` varchar(256) NOT NULL COMMENT '购物车id',
  `total_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '订单商品总数',
  `total_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '订单总价',
  `total_postage` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '邮费',
  `pay_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '实际支付金额',
  `pay_postage` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '支付邮费',
  `is_selfbuy` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否为自购',
  `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '一级佣金',
  `extension_two` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '二级佣金',
  `commission_rate` decimal(7,4) unsigned NOT NULL DEFAULT '0.0000' COMMENT '平台手续费',
  `integral` int(10) unsigned DEFAULT '0' COMMENT '使用积分数量',
  `integral_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '积分抵扣金额',
  `give_integral` int(10) unsigned DEFAULT '0' COMMENT '赠送积分',
  `coupon_id` varchar(128) NOT NULL DEFAULT '' COMMENT '优惠券id',
  `coupon_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '优惠券金额',
  `platform_coupon_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '平台优惠券金额',
  `svip_discount` decimal(8,2) unsigned DEFAULT '0.00' COMMENT 'svip优惠金额',
  `order_type` tinyint(3) unsigned DEFAULT '0' COMMENT '0普通1自提',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '支付状态',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
  `pay_type` tinyint(1) NOT NULL COMMENT '支付方式 0余额 1微信 2小程序 3 h5  4支付宝 5 支付宝扫码 6 微信扫码',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '订单状态（0：待发货；1：待收货；2：待评价；3：已完成； 9: 拼团中 10:  待付尾款 11:尾款超时未付 -1：已退款）',
  `delivery_type` varchar(32) DEFAULT NULL COMMENT '发货类型(1:发货 2: 送货 3: 虚拟,4电子面单，5同城 6 卡密自动发货)',
  `is_virtual` tinyint(3) unsigned DEFAULT '0' COMMENT '0:实物订单 1:虚拟订单',
  `delivery_name` varchar(50) DEFAULT NULL COMMENT '快递名称/送货人姓名',
  `delivery_id` varchar(255) DEFAULT NULL COMMENT '快递单号/手机号',
  `mark` varchar(512) NOT NULL COMMENT '备注',
  `remark` varchar(512) DEFAULT NULL COMMENT '管理员备注',
  `admin_mark` varchar(512) DEFAULT NULL COMMENT '总后台备注',
  `verify_code` char(16) DEFAULT NULL COMMENT '核销码',
  `verify_time` timestamp NULL DEFAULT NULL COMMENT '核销时间/收货时间',
  `verify_service_id` int(10) unsigned DEFAULT NULL COMMENT '核销客服 id',
  `transaction_id` varchar(60) DEFAULT NULL COMMENT '微信支付订单号(分账时有效)',
  `activity_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1:秒杀 2:预售 3:助力 10:套餐',
  `order_extend` varchar(1024) DEFAULT NULL COMMENT '自定义表单数据',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户ID',
  `reconciliation_id` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '对账id',
  `cost` decimal(8,2) unsigned NOT NULL COMMENT '成本价',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `is_system_del` tinyint(1) DEFAULT '0' COMMENT '后台是否删除',
  `verify_status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '核销订单状态0 默认 1 部分核销 2 全部核销',
  `refund_switch` tinyint(3) unsigned DEFAULT '1' COMMENT '是否支持退款',
  `kuaidi_label` varchar(255) DEFAULT NULL COMMENT '快递单号图片',
  `task_id` varchar(255) NOT NULL DEFAULT '' COMMENT '快递单号任务ID',
  `is_behalf` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为代客下单订单',
  `behalf_no_verify` tinyint(1) NOT NULL DEFAULT '0' COMMENT '代客下单订单是否无需核销（0:需要核销，1:无需核销）',
  `enable_assigned` tinyint(1) DEFAULT '0' COMMENT '预约商品订单是否派单 0 领取 1 指派',
  `staffs_id` int(10) unsigned DEFAULT '0' COMMENT '领取/指派的员工ID',
  `is_cancel` int(11) DEFAULT '0' COMMENT '是否可取消预约(0:不可取消,1:可取消)',
  `reservation_service_voucher` text COMMENT '预约订单服务凭证',
  `clock_in_info` text COMMENT '预约订单服务打卡信息',
  `kuaidi_order_id` varchar(255) NOT NULL DEFAULT '' COMMENT '商家寄件的快递单号',
  `is_stock_up` int(11) NOT NULL DEFAULT '0' COMMENT '商家寄件快递是否发出',
  `merchant_take_id` int(11) NOT NULL DEFAULT '0' COMMENT '自提点id',
  `merchant_take_info` text NOT NULL COMMENT '同城配送信息',
  PRIMARY KEY (`order_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `verify_code` (`verify_code`),
  KEY `main_id` (`main_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单表';

DROP TABLE IF EXISTS `qixi_store_order_product`;
CREATE TABLE `qixi_store_order_product` (

  `order_product_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单产品 id',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户 id',
  `cart_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '购物车id',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品ID',
  `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '一级佣金',
  `extension_two` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '二级佣金',
  `integral` int(10) unsigned DEFAULT '0' COMMENT '使用积分(单数)',
  `integral_price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '积分抵扣金额',
  `integral_total` int(10) unsigned DEFAULT '0' COMMENT '使用积分(总数)',
  `coupon_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '优惠金额',
  `platform_coupon_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '平台优惠金额',
  `svip_discount` decimal(8,2) unsigned DEFAULT '0.00' COMMENT 'svip优惠金额',
  `postage_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '运费',
  `product_sku` char(12) NOT NULL COMMENT '商品 sku',
  `is_refund` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否退款 0:未退款 1:退款中 2:部分退款 3=全退',
  `product_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '购买数量',
  `product_type` int(11) NOT NULL DEFAULT '0' COMMENT '0.普通商品 1.秒杀商品,2.预售商品',
  `activity_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动关联 id',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '可申请退货数量',
  `is_reply` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否评价',
  `cost` decimal(10,2) unsigned NOT NULL COMMENT '商品成本价',
  `product_price` decimal(10,2) unsigned NOT NULL COMMENT '商品金额',
  `total_price` decimal(10,2) unsigned NOT NULL COMMENT '商品售价',
  `cart_info` text NOT NULL COMMENT '购买东西的详细信息',
  `refund_switch` tinyint(3) unsigned DEFAULT '1' COMMENT '是否支持退款',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reservation_date` varchar(20) DEFAULT '' COMMENT '预约商品的预约日期',
  `reservation_id` int(10) unsigned DEFAULT '0' COMMENT '预约商品的预约时间段ID',
  `reservation_time_part` varchar(20) NOT NULL DEFAULT '' COMMENT '预约时间段',
  `settlement_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '员工结算价格',
  PRIMARY KEY (`order_product_id`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE,
  KEY `oid` (`order_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单购物详情表';

DROP TABLE IF EXISTS `qixi_store_order_profitsharing`;
CREATE TABLE `qixi_store_order_profitsharing` (

  `profitsharing_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `profitsharing_sn` varchar(32) NOT NULL COMMENT '分账 id',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单 id',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `transaction_id` varchar(60) NOT NULL COMMENT '微信支付订单号',
  `profitsharing_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '分账金额',
  `profitsharing_refund` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '退款金额',
  `profitsharing_mer_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '分账分出去的金额即给平台的手续费',
  `type` varchar(32) NOT NULL COMMENT '分类',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:未分账 1:已分账 -1已退款 -2失败',
  `error_msg` varchar(255) DEFAULT NULL COMMENT '失败原因',
  `profitsharing_time` timestamp NULL DEFAULT NULL COMMENT '分账时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_combine` tinyint(1) DEFAULT '1' COMMENT '分账类型：1 平台收付通 2 服务商',
  PRIMARY KEY (`profitsharing_id`),
  KEY `order_id` (`order_id`),
  KEY `mer_id` (`mer_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='分账表';

DROP TABLE IF EXISTS `qixi_store_order_receipt`;
CREATE TABLE `qixi_store_order_receipt` (

  `order_receipt_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `order_id` varchar(255) NOT NULL DEFAULT '0' COMMENT '订单ID',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `receipt_info` varchar(500) DEFAULT '' COMMENT '发票类型：1.普通发票，2.增值税发票',
  `status` tinyint(4) DEFAULT '0' COMMENT '开票状态：1.已出票,10.已寄出',
  `receipt_sn` varchar(255) DEFAULT '' COMMENT '发票单号',
  `receipt_no` varchar(255) DEFAULT NULL COMMENT '发票编号',
  `delivery_info` varchar(255) DEFAULT NULL COMMENT '收票联系信息',
  `pic` varchar(500) NOT NULL DEFAULT '' COMMENT '发票文件地址',
  `mark` varchar(255) DEFAULT NULL COMMENT '用户备注',
  `receipt_price` decimal(10,2) DEFAULT NULL COMMENT '开票金额',
  `order_price` decimal(10,2) DEFAULT NULL COMMENT '订单金额',
  `status_time` datetime NOT NULL COMMENT '状态变更时间',
  `is_del` tinyint(1) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT NULL,
  `mer_id` int(11) DEFAULT '0',
  `mer_mark` varchar(255) DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`order_receipt_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单发票信息';

DROP TABLE IF EXISTS `qixi_store_order_status`;
CREATE TABLE `qixi_store_order_status` (

  `order_id` int(10) unsigned NOT NULL COMMENT '订单id',
  `order_sn` varchar(256) DEFAULT NULL COMMENT '订单号',
  `type` varchar(20) DEFAULT NULL COMMENT '订单类型',
  `change_type` varchar(32) NOT NULL COMMENT '操作类型',
  `change_message` varchar(256) NOT NULL COMMENT '操作备注',
  `nickname` varchar(20) DEFAULT NULL,
  `uid` int(11) DEFAULT NULL COMMENT '操作者ID',
  `user_type` tinyint(4) DEFAULT NULL COMMENT '操作者类型',
  `change_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  KEY `order_id` (`order_id`) USING BTREE,
  KEY `change_type` (`change_type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单操作记录表';

DROP TABLE IF EXISTS `qixi_store_printer`;
CREATE TABLE `qixi_store_printer` (

  `printer_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `printer_name` varchar(50) NOT NULL DEFAULT '' COMMENT '名称',
  `printer_appkey` varchar(50) NOT NULL DEFAULT '' COMMENT '打印机的应用ID',
  `printer_terminal` varchar(50) NOT NULL DEFAULT '' COMMENT '打印机终端号',
  `printer_appid` varchar(50) NOT NULL DEFAULT '' COMMENT '打印机应用用户ID',
  `printer_secret` varchar(50) NOT NULL DEFAULT '' COMMENT '打印机应用密匙',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `mer_id` int(11) NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `type` tinyint(4) DEFAULT '0' COMMENT '0 易联云 1 飞鹅云',
  `times` int(11) NOT NULL DEFAULT '1' COMMENT '打印联数',
  `print_content` varchar(2000) NOT NULL DEFAULT '' COMMENT '打印内容',
  `print_type` int(11) NOT NULL DEFAULT '1' COMMENT '打印时机1支付后，2下单后',
  PRIMARY KEY (`printer_id`) USING BTREE,
  UNIQUE KEY `id` (`printer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product`;
CREATE TABLE `qixi_store_product` (

  `product_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '商品id',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户Id',
  `store_name` varchar(128) NOT NULL COMMENT '商品名称',
  `store_info` varchar(256) DEFAULT NULL COMMENT '商品简介',
  `keyword` varchar(128) NOT NULL COMMENT '关键字',
  `bar_code` varchar(15) NOT NULL DEFAULT '' COMMENT '产品条码（一维码）',
  `brand_id` int(11) DEFAULT NULL COMMENT '品牌 id',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商户 状态（0:未上架，1:上架，2:定时上架）',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '管理员 状态（0：审核中，1：审核通过 -1: 未通过 -2: 下架）',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `mer_status` tinyint(1) DEFAULT '1' COMMENT '商铺状态是否 1.正常 0. 非正常',
  `cate_id` int(11) NOT NULL COMMENT '分类id',
  `unit_name` varchar(16) NOT NULL COMMENT '单位名',
  `sort` smallint(6) NOT NULL DEFAULT '0' COMMENT '排序',
  `rank` smallint(6) NOT NULL DEFAULT '0' COMMENT '总后台排序',
  `sales` mediumint(8) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '最低价格',
  `cost` decimal(10,2) DEFAULT '0.00' COMMENT '成本价',
  `ot_price` decimal(10,2) DEFAULT '0.00' COMMENT '原价',
  `stock` int(10) unsigned DEFAULT '0' COMMENT '总库存',
  `is_hot` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否热卖',
  `is_benefit` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '促销单品',
  `is_best` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否精品',
  `is_new` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否新品',
  `is_good` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否优品推荐',
  `product_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0.普通商品 1.秒杀商品,2.预售商品，3.助力商品，4.拼团商品',
  `ficti` mediumint(9) DEFAULT '0' COMMENT '虚拟销量',
  `browse` int(11) DEFAULT '0' COMMENT '浏览量',
  `code_path` varchar(64) NOT NULL DEFAULT '' COMMENT '产品二维码地址(用户小程序海报)',
  `video_link` varchar(200) NOT NULL DEFAULT '' COMMENT '主图视频链接',
  `temp_id` int(11) NOT NULL DEFAULT '1' COMMENT '运费模板ID',
  `spec_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '规格 0单 1多',
  `extension_type` tinyint(1) DEFAULT '0' COMMENT '佣金比例 0.系统，1.自定义',
  `refusal` varchar(255) DEFAULT NULL COMMENT '审核拒绝理由',
  `rate` decimal(2,1) DEFAULT '5.0' COMMENT '评价分数',
  `reply_count` int(10) unsigned DEFAULT '0' COMMENT '评论数',
  `give_coupon_ids` varchar(500) DEFAULT NULL COMMENT '赠送优惠券',
  `is_gift_bag` tinyint(1) DEFAULT '0' COMMENT '是否为礼包',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `care_count` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
  `is_used` int(11) DEFAULT '1' COMMENT '显示/隐藏',
  `old_product_id` int(11) DEFAULT '0' COMMENT '原商品ID',
  `image` varchar(256) NOT NULL DEFAULT '' COMMENT '商品图片',
  `slider_image` varchar(2000) NOT NULL COMMENT '轮播图',
  `guarantee_template_id` int(11) DEFAULT '0' COMMENT '保障服务模板',
  `once_max_count` int(11) DEFAULT '0' COMMENT '订单单次购买数量最大限制',
  `once_min_count` int(11) NOT NULL DEFAULT '0' COMMENT '单次购买最低限购',
  `integral_rate` int(11) NOT NULL DEFAULT '-1' COMMENT '积分抵扣比例',
  `integral_total` int(10) unsigned DEFAULT '0' COMMENT '使用积分抵扣总数',
  `integral_price_total` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '使用积分抵扣金额总数',
  `labels` varchar(255) DEFAULT '' COMMENT '标签id',
  `delivery_way` varchar(100) DEFAULT NULL COMMENT '1.仅到店自提2快递计价配送3全国包邮',
  `delivery_free` int(11) DEFAULT '0' COMMENT '全国包邮',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0.实体商品，1.虚拟商品，2 网盘，3 卡密',
  `extend` varchar(1000) NOT NULL DEFAULT '' COMMENT '扩展信息',
  `pay_limit` tinyint(4) NOT NULL DEFAULT '0' COMMENT '购买总数限制 0:不限购，1单次限购 2 长期限购',
  `svip_price_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0不参加，1默认比例，2自定义',
  `svip_price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '会员价',
  `mer_svip_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商户会员状态',
  `param_temp_id` varchar(255) DEFAULT NULL COMMENT '参数模板ID',
  `refund_switch` tinyint(4) DEFAULT '1' COMMENT '是否支持退款',
  `delete` tinyint(4) DEFAULT '0',
  `mer_form_id` int(11) DEFAULT '0' COMMENT '系统表单ID',
  `good_ids` varchar(2000) DEFAULT '' COMMENT '推荐商品',
  `auto_on_time` int(11) DEFAULT NULL COMMENT '自动上架时间',
  `auto_off_time` int(11) DEFAULT NULL COMMENT '自动下架时间',
  `active_id` int(11) DEFAULT NULL COMMENT '秒杀活动ID',
  `cate_hot` tinyint(4) DEFAULT '0' COMMENT '分类大图推荐 1 推荐',
  `bar_code_number` varchar(255) DEFAULT '' COMMENT '商品条码',
  `custom_temp_id` varchar(255) NOT NULL DEFAULT '' COMMENT '自定义参数模版id',
  `activity_label_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '活动标签',
  PRIMARY KEY (`product_id`) USING BTREE,
  KEY `cate_id` (`cate_id`) USING BTREE,
  KEY `sort` (`sort`) USING BTREE,
  KEY `sales` (`sales`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE,
  KEY `bar_code_number` (`bar_code_number`) USING BTREE,
  KEY `status` (`is_show`,`status`,`is_used`,`product_type`,`mer_status`,`is_gift_bag`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='商品表';

DROP TABLE IF EXISTS `qixi_store_product_assist`;
CREATE TABLE `qixi_store_product_assist` (

  `product_assist_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `end_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '结束时间',
  `status` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '平台控制状态：1开启，0.结束',
  `pay_count` int(10) unsigned DEFAULT '0' COMMENT '限购数量，0为不限制',
  `assist_count` int(10) unsigned DEFAULT '0' COMMENT '助力总需人数',
  `assist_user_count` int(10) unsigned DEFAULT '0' COMMENT '单人可助力次数',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品ID',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '商户控制状态 0.下架；1.上架',
  `store_name` varchar(128) NOT NULL COMMENT '商品活动标题',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0',
  `store_info` varchar(255) DEFAULT NULL COMMENT '商品简介',
  `is_del` int(10) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `product_status` int(11) DEFAULT '0' COMMENT '审核状态；0.待审核，1审核通过，-1 审核失败，-2 强制下架',
  `refusal` varchar(255) DEFAULT NULL,
  `action_status` int(11) DEFAULT '1' COMMENT '活动状态1开启，-1 结束',
  PRIMARY KEY (`product_assist_id`) USING BTREE,
  KEY `start_time` (`start_time`,`end_time`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品助力活动表';

DROP TABLE IF EXISTS `qixi_store_product_assist_set`;
CREATE TABLE `qixi_store_product_assist_set` (

  `product_assist_set_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_assist_id` int(10) unsigned NOT NULL,
  `product_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '状态：-1 未完成 ，1 进行中， 10 已完成，20.已支付',
  `assist_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '需助力总人数',
  `assist_user_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单人可助力次数',
  `yet_assist_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '已助力人数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '时间',
  `mer_id` int(10) unsigned DEFAULT '0',
  `share_num` int(10) unsigned DEFAULT '0',
  `view_num` int(10) unsigned DEFAULT '0',
  `is_del` int(10) unsigned DEFAULT '0',
  PRIMARY KEY (`product_assist_set_id`) USING BTREE,
  KEY `product_assist_id` (`product_assist_id`,`product_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='助力发起列表';

DROP TABLE IF EXISTS `qixi_store_product_assist_sku`;
CREATE TABLE `qixi_store_product_assist_sku` (

  `product_assist_id` int(10) unsigned NOT NULL DEFAULT '0',
  `product_id` int(10) unsigned NOT NULL,
  `unique` char(12) NOT NULL,
  `assist_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '助力售价',
  `stock` int(10) unsigned NOT NULL DEFAULT '0',
  `stock_count` int(10) unsigned DEFAULT '0' COMMENT '总限购',
  KEY `product_assist_id` (`product_assist_id`,`product_id`) USING BTREE,
  KEY `unique` (`unique`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_assist_user`;
CREATE TABLE `qixi_store_product_assist_user` (

  `product_assist_user_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_assist_set_id` int(10) unsigned NOT NULL,
  `product_assist_id` int(10) unsigned NOT NULL,
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'id',
  `avatar_img` varchar(256) DEFAULT NULL COMMENT '头像',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_assist_user_id`) USING BTREE,
  KEY `uid` (`uid`,`product_assist_set_id`) USING BTREE,
  KEY `product_assist_id` (`product_assist_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='助力记录表';

DROP TABLE IF EXISTS `qixi_store_product_attr`;
CREATE TABLE `qixi_store_product_attr` (

  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品ID',
  `attr_name` varchar(32) NOT NULL COMMENT '属性名',
  `attr_values` varchar(2000) NOT NULL COMMENT '属性值',
  `type` tinyint(1) DEFAULT '0' COMMENT '活动类型 0=商品',
  KEY `product_id` (`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品属性表';

DROP TABLE IF EXISTS `qixi_store_product_attr_reservation`;
CREATE TABLE `qixi_store_product_attr_reservation` (

  `attr_reservation_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'attr_reservation_id',
  `attr_value_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品属性值表ID',
  `unique` char(12) NOT NULL DEFAULT '' COMMENT '唯一值',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品ID',
  `start_time` varchar(10) NOT NULL DEFAULT '' COMMENT '开始时间段',
  `end_time` varchar(10) NOT NULL DEFAULT '' COMMENT '结束时间段',
  `stock` int(11) NOT NULL DEFAULT '0' COMMENT '可约数量',
  `use_num` int(11) NOT NULL DEFAULT '0' COMMENT '使用数量',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`attr_reservation_id`) USING BTREE,
  KEY `idx_attr_value_id` (`attr_value_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_product_attr_result`;
CREATE TABLE `qixi_store_product_attr_result` (

  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `result` longtext NOT NULL COMMENT '商品属性参数',
  `change_time` int(10) unsigned DEFAULT NULL COMMENT '上次修改时间',
  `type` tinyint(1) DEFAULT NULL COMMENT '活动类型 0=商品，1=秒杀，2=预售，3=助力, 4=拼团',
  PRIMARY KEY (`id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_product_attr_value`;
CREATE TABLE `qixi_store_product_attr_value` (

  `value_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `detail` varchar(1000) NOT NULL DEFAULT '',
  `sku` varchar(128) NOT NULL COMMENT '商品属性索引值 (attr_value|attr_value[|....])',
  `stock` int(10) unsigned NOT NULL COMMENT '属性对应的库存',
  `sales` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '销量',
  `image` varchar(128) DEFAULT NULL COMMENT '图片',
  `bar_code` varchar(50) NOT NULL DEFAULT '' COMMENT '产品条码',
  `cost` decimal(8,2) unsigned NOT NULL COMMENT '成本价',
  `ot_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '原价',
  `price` decimal(8,2) unsigned NOT NULL COMMENT '价格',
  `volume` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '体积',
  `weight` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '重量',
  `type` tinyint(1) DEFAULT '0' COMMENT '活动类型 0=商品; 20 积分商品',
  `extension_one` decimal(8,2) DEFAULT '0.00' COMMENT '一级佣金',
  `extension_two` decimal(8,2) DEFAULT '0.00' COMMENT '二级佣金',
  `unique` char(12) NOT NULL DEFAULT '' COMMENT '唯一值',
  `svip_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '会员价',
  `library_id` int(11) DEFAULT '0' COMMENT '一次性卡密关联卡密库地',
  `bar_code_number` varchar(50) NOT NULL DEFAULT '' COMMENT '规格条码',
  `is_default_select` tinyint(1) DEFAULT '0' COMMENT '默认显示',
  `is_show` tinyint(1) DEFAULT '1',
  `settlement_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '员工结算价格',
  PRIMARY KEY (`value_id`),
  KEY `store_id` (`product_id`,`sku`) USING BTREE,
  KEY `unique` (`unique`,`sku`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 COMMENT='商品属性值表';

DROP TABLE IF EXISTS `qixi_store_product_cate`;
CREATE TABLE `qixi_store_product_cate` (

  `product_id` int(11) DEFAULT NULL,
  `mer_cate_id` int(11) DEFAULT NULL,
  `mer_id` int(11) DEFAULT NULL,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `mer_cate_id` (`mer_cate_id`,`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品商户分类关联表';

DROP TABLE IF EXISTS `qixi_store_product_cdkey`;
CREATE TABLE `qixi_store_product_cdkey` (

  `cdkey_id` int(11) NOT NULL AUTO_INCREMENT,
  `is_type` tinyint(4) DEFAULT '0' COMMENT '卡密类型： 0 固定卡密， 1 一次性卡密',
  `value_id` int(11) DEFAULT NULL COMMENT '商品规格ID',
  `key` varchar(255) DEFAULT NULL COMMENT '卡密内容',
  `pwd` varchar(255) DEFAULT NULL COMMENT '卡密密码',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态： 1 可用  -1 不可用/已使用',
  `product_id` int(11) DEFAULT NULL COMMENT '商品ID',
  `library_id` int(11) DEFAULT '0' COMMENT '卡密库ID',
  `is_use` tinyint(4) DEFAULT '0' COMMENT '是否已使用',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`cdkey_id`) USING BTREE,
  KEY `library_id` (`library_id`) USING BTREE,
  KEY `status` (`status`,`is_use`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_content`;
CREATE TABLE `qixi_store_product_content` (

  `product_id` int(10) unsigned NOT NULL COMMENT '商品id',
  `content` longtext NOT NULL COMMENT '商品详情',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '商品类型 0=普通',
  KEY `product_id` (`product_id`,`type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品详情表';

DROP TABLE IF EXISTS `qixi_store_product_copy`;
CREATE TABLE `qixi_store_product_copy` (

  `store_product_copy_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `type` varchar(255) DEFAULT NULL COMMENT '''taobao'', ''tmall'', ''jd'', ''pinduoduo'', ''suning'', ''yangkeduo''\n',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户id',
  `num` int(11) DEFAULT NULL COMMENT '数量',
  `number` int(11) DEFAULT '1' COMMENT '剩余数量',
  `message` varchar(255) DEFAULT NULL,
  `info` text COMMENT '信息',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`store_product_copy_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_group`;
CREATE TABLE `qixi_store_product_group` (

  `product_group_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `product_id` int(10) unsigned DEFAULT '0' COMMENT '商品ID',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `time` int(10) unsigned DEFAULT '0' COMMENT '开团时长',
  `buying_count_num` int(11) DEFAULT '0' COMMENT '拼团总人数',
  `buying_num` int(11) DEFAULT '0' COMMENT '最少真实购买人数',
  `pay_count` int(10) unsigned DEFAULT '0' COMMENT '活动购买总人数',
  `once_pay_count` int(10) unsigned DEFAULT '0' COMMENT '单次购买数量',
  `status` int(11) DEFAULT '0' COMMENT '平台控制状态',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户ID',
  `ficti_status` int(11) DEFAULT '0' COMMENT '虚拟成团状态',
  `ficti_num` int(11) DEFAULT '0' COMMENT '最多虚拟人数',
  `is_show` int(11) DEFAULT '0' COMMENT '上下架',
  `is_del` int(10) unsigned DEFAULT '0',
  `success_num` int(10) unsigned DEFAULT '0' COMMENT '成功团数',
  `product_status` int(11) DEFAULT '0',
  `price` decimal(10,2) DEFAULT '0.00',
  `action_status` int(11) DEFAULT '0' COMMENT '活动状态',
  `create_time` datetime DEFAULT NULL,
  `refusal` varchar(255) DEFAULT NULL,
  `leader_extension` tinyint(4) DEFAULT '0' COMMENT '团长分销',
  `leader_rate` decimal(10,2) DEFAULT '0.00' COMMENT '分销比例',
  PRIMARY KEY (`product_group_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='拼团商品信息表';

DROP TABLE IF EXISTS `qixi_store_product_group_buying`;
CREATE TABLE `qixi_store_product_group_buying` (

  `group_buying_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_group_id` int(10) unsigned DEFAULT '0' COMMENT '活动商品ID',
  `status` int(11) DEFAULT '0' COMMENT '状态：0。默认，进行中，10.已完成，-1 时间到未完成',
  `ficti_status` int(11) DEFAULT '0' COMMENT '虚拟成团状态0.未开启，1开启',
  `ficti_num` int(10) unsigned DEFAULT '0' COMMENT '虚拟成团人数',
  `buying_count_num` int(10) unsigned DEFAULT '0' COMMENT '成团总人数',
  `buying_num` int(10) unsigned DEFAULT '0' COMMENT '真实人数',
  `yet_buying_num` int(10) unsigned DEFAULT '0' COMMENT '已参团人数',
  `is_del` int(11) DEFAULT '0',
  `mer_id` int(10) unsigned DEFAULT '0',
  `end_time` int(11) DEFAULT NULL COMMENT '结束时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_hidde` tinyint(1) DEFAULT '0' COMMENT '是否隐藏团信息 0 否 1 是',
  PRIMARY KEY (`group_buying_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='拼团活动表';

DROP TABLE IF EXISTS `qixi_store_product_group_sku`;
CREATE TABLE `qixi_store_product_group_sku` (

  `product_group_id` int(10) unsigned NOT NULL DEFAULT '0',
  `product_id` int(10) unsigned NOT NULL,
  `unique` char(12) NOT NULL,
  `active_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '活动价',
  `stock` int(10) unsigned NOT NULL DEFAULT '0',
  `stock_count` int(10) unsigned DEFAULT '0',
  KEY `product_group_id` (`product_group_id`,`product_id`),
  KEY `unique` (`unique`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_group_user`;
CREATE TABLE `qixi_store_product_group_user` (

  `group_buying_id` int(10) unsigned DEFAULT '0' COMMENT '团ID',
  `product_group_id` int(10) unsigned DEFAULT '0' COMMENT '活动商品ID',
  `status` int(11) DEFAULT '0' COMMENT '状态',
  `is_initiator` int(10) unsigned DEFAULT '0' COMMENT '是否为 团长',
  `order_id` int(10) unsigned DEFAULT '0' COMMENT '订单ID',
  `uid` int(10) unsigned DEFAULT '0' COMMENT '用户ID ',
  `nickname` varchar(255) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(255) DEFAULT NULL COMMENT '头像',
  `is_del` int(10) unsigned DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_leader` tinyint(1) DEFAULT '0' COMMENT '是否为创建者'

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='拼团成员表';

DROP TABLE IF EXISTS `qixi_store_product_label`;
CREATE TABLE `qixi_store_product_label` (

  `product_label_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `label_name` varchar(50) DEFAULT NULL COMMENT '标签名',
  `status` tinyint(1) DEFAULT NULL COMMENT '状态',
  `info` varchar(255) DEFAULT NULL COMMENT '说明',
  `sort` int(11) DEFAULT NULL COMMENT '排序',
  `type` int(11) DEFAULT '0' COMMENT '类型 ',
  `mer_id` int(11) DEFAULT '0' COMMENT '商户ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`product_label_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_presell`;
CREATE TABLE `qixi_store_product_presell` (

  `product_presell_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '预售开始时间',
  `end_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '预售结束时间',
  `final_start_time` varchar(30) DEFAULT '' COMMENT '尾款支付开始时间',
  `final_end_time` varchar(30) DEFAULT '' COMMENT '尾款支付结时间',
  `status` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '平台控制状态：1开启，0.结束',
  `presell_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '预售类型：1.全款预售，2.定金预售',
  `pay_count` int(10) unsigned DEFAULT '0' COMMENT '限购数量，0为不限制',
  `delivery_type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '发货类型：1.支付成功后 ； 2. 预售结束后',
  `delivery_day` int(10) unsigned DEFAULT '0' COMMENT '发货时间',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品ID',
  `price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '预售最低价',
  `is_show` tinyint(3) unsigned DEFAULT NULL COMMENT '商户控制状态 0.下架；1.上架',
  `store_name` varchar(128) NOT NULL COMMENT '商品活动标题',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0',
  `store_info` varchar(255) DEFAULT NULL COMMENT '商品简介',
  `is_del` int(10) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `product_status` int(11) DEFAULT '0' COMMENT '审核状态；0.待审核，1审核通过，-1 审核失败，-2 强制下架',
  `refusal` varchar(255) DEFAULT NULL,
  `action_status` int(11) DEFAULT '1' COMMENT '活动状态1开启，-1 结束',
  PRIMARY KEY (`product_presell_id`) USING BTREE,
  KEY `start_time` (`start_time`,`end_time`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品预售活动表';

DROP TABLE IF EXISTS `qixi_store_product_presell_sku`;
CREATE TABLE `qixi_store_product_presell_sku` (

  `product_presell_id` int(10) unsigned NOT NULL DEFAULT '0',
  `product_id` int(10) unsigned NOT NULL,
  `unique` char(12) NOT NULL,
  `presell_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '预售价',
  `stock` int(10) unsigned NOT NULL DEFAULT '0',
  `stock_count` int(11) NOT NULL DEFAULT '0' COMMENT '总限购',
  `down_price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '订金',
  `final_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '尾款金额',
  `one_take` int(10) unsigned DEFAULT '0' COMMENT '第一阶段参与人数',
  `one_pay` int(10) unsigned DEFAULT '0' COMMENT '第一阶段支付人数',
  `two_pay` int(10) unsigned DEFAULT '0' COMMENT '第二阶段支付人数',
  `seles` int(10) unsigned DEFAULT '0' COMMENT '销量',
  KEY `product_presell_id` (`product_presell_id`,`product_id`) USING BTREE,
  KEY `unique` (`unique`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_reply`;
CREATE TABLE `qixi_store_product_reply` (

  `reply_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '评论ID',
  `uid` int(11) NOT NULL COMMENT '用户ID',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `order_product_id` int(11) NOT NULL COMMENT '订单商品ID',
  `unique` char(12) DEFAULT NULL COMMENT '商品 sku',
  `product_id` int(11) NOT NULL COMMENT '商品id',
  `product_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0=普通商品',
  `product_score` tinyint(1) NOT NULL COMMENT '商品分数',
  `service_score` tinyint(1) NOT NULL COMMENT '服务分数',
  `postage_score` tinyint(1) NOT NULL COMMENT '物流分数',
  `rate` float(2,1) DEFAULT '5.0' COMMENT '平均值',
  `comment` varchar(512) NOT NULL COMMENT '评论内容',
  `pics` text COMMENT '评论图片',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评论时间',
  `merchant_reply_content` varchar(300) DEFAULT NULL COMMENT '管理员回复内容',
  `merchant_reply_time` timestamp NULL DEFAULT NULL COMMENT '管理员回复时间',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '1' COMMENT '商家排序',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0未删除1已删除',
  `is_reply` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0未回复1已回复',
  `is_virtual` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0不是虚拟评价1是虚拟评价',
  `nickname` varchar(64) NOT NULL COMMENT '用户名称',
  `avatar` varchar(255) NOT NULL COMMENT '用户头像',
  PRIMARY KEY (`reply_id`) USING BTREE,
  UNIQUE KEY `order_id` (`order_product_id`,`unique`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `product_id` (`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品评论表';

DROP TABLE IF EXISTS `qixi_store_product_reservation`;
CREATE TABLE `qixi_store_product_reservation` (

  `product_reservation_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'product_reservation_id',
  `product_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商品id',
  `reservation_time_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '预约时段划分类型(1:自动划分,2:自定义划分)',
  `reservation_start_time` varchar(20) NOT NULL DEFAULT '' COMMENT '预约开始时间',
  `reservation_end_time` varchar(20) NOT NULL DEFAULT '' COMMENT '预约结束时间',
  `reservation_time_interval` int(11) NOT NULL DEFAULT '10' COMMENT '时间跨度,以分钟为单位',
  `time_period` text NOT NULL COMMENT '时间段',
  `reservation_type` tinyint(1) NOT NULL DEFAULT '2' COMMENT '预约类型(1:到店服务,2:上门服务,3:上门+到店服务)',
  `show_num_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否展示可约数量(0:不展示,1:展示)',
  `sale_time_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '可售日期(1:每天,2:自定义时间)',
  `sale_time_start_day` varchar(20) NOT NULL DEFAULT '' COMMENT '可售日期自定义开始时间',
  `sale_time_end_day` varchar(20) NOT NULL DEFAULT '' COMMENT '可售日期自定义结束时间',
  `sale_time_week` varchar(20) NOT NULL DEFAULT '' COMMENT '可售日期周数据',
  `show_reservation_days` int(11) NOT NULL DEFAULT '1' COMMENT '显示日期范围',
  `is_advance` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否提前预约(0:不提前,1:提前)',
  `advance_time` int(11) NOT NULL DEFAULT '1' COMMENT '提前预约时间,以小时为单位',
  `is_cancel_reservation` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否可取消预约(0:不可取消,1:可取消)',
  `cancel_reservation_time` int(11) NOT NULL DEFAULT '1' COMMENT '取消预约提前时间,以小时为单位',
  `reservation_form_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '预约表单类型(1:每个预约提交,2:每单提交)',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
  `is_del` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`product_reservation_id`) USING BTREE,
  KEY `idx_product_id` (`product_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_store_product_sku`;
CREATE TABLE `qixi_store_product_sku` (

  `product_sku_id` int(11) NOT NULL AUTO_INCREMENT,
  `active_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动ID',
  `active_product_id` int(10) unsigned DEFAULT NULL COMMENT '活动商品的ID',
  `active_type` int(10) unsigned DEFAULT '0' COMMENT '活动类型',
  `product_id` int(10) unsigned NOT NULL,
  `unique` char(12) NOT NULL,
  `price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '原价',
  `active_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '活动售价',
  `stock` int(10) unsigned NOT NULL DEFAULT '0',
  `stock_count` int(10) unsigned DEFAULT '0' COMMENT '总限购',
  PRIMARY KEY (`product_sku_id`) USING BTREE,
  KEY `active_id` (`active_id`,`product_id`) USING BTREE,
  KEY `active_id_2` (`active_id`,`unique`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_product_take`;
CREATE TABLE `qixi_store_product_take` (

  `product_take_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) DEFAULT NULL,
  `unique` char(12) DEFAULT NULL,
  `status` int(10) unsigned DEFAULT '0' COMMENT '默认0，发送 1',
  `uid` int(11) DEFAULT NULL COMMENT '用户',
  `type` varchar(255) DEFAULT NULL COMMENT '1.PC,2.公众号,3.小程序',
  `is_del` int(11) DEFAULT '0',
  PRIMARY KEY (`product_take_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户到货通知记录';

DROP TABLE IF EXISTS `qixi_store_product_unit`;
CREATE TABLE `qixi_store_product_unit` (

  `product_unit_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `value` varchar(255) NOT NULL DEFAULT '' COMMENT '值',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态 1|开启 0|关闭',
  `is_del` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除 0|正常 1|删除',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`product_unit_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_refund_order`;
CREATE TABLE `qixi_store_refund_order` (

  `refund_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '退款单id',
  `refund_order_sn` varchar(32) NOT NULL COMMENT '退款单号',
  `order_id` int(10) unsigned NOT NULL COMMENT '订单id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户 id',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '退还一级佣金',
  `extension_two` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '退还二级佣金',
  `integral` int(10) unsigned DEFAULT '0' COMMENT '退还积分',
  `delivery_type` varchar(32) DEFAULT NULL COMMENT '快递公司',
  `delivery_id` varchar(32) DEFAULT NULL COMMENT '快递单号',
  `delivery_mark` varchar(200) DEFAULT NULL COMMENT '快递备注',
  `delivery_pics` varchar(255) DEFAULT NULL COMMENT '快递凭证',
  `delivery_phone` varchar(18) DEFAULT NULL COMMENT '联系电话',
  `mer_delivery_user` varchar(32) DEFAULT NULL COMMENT '收货人',
  `mer_delivery_address` varchar(32) DEFAULT NULL COMMENT '收货地址',
  `phone` varchar(18) DEFAULT NULL COMMENT '联系电话',
  `mark` varchar(200) DEFAULT '' COMMENT '备注',
  `mer_mark` varchar(255) DEFAULT '' COMMENT '商户备注',
  `admin_mark` varchar(255) DEFAULT '' COMMENT '平台备注',
  `pics` text COMMENT '图片',
  `refund_type` tinyint(1) NOT NULL COMMENT '退款类型 1:退款 2:退款退货',
  `refund_message` varchar(128) NOT NULL COMMENT '退款原因',
  `refund_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '退款金额',
  `platform_refund_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '退款平台优惠券金额',
  `refund_postage` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '退的运费',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '退款数',
  `fail_message` varchar(200) DEFAULT NULL COMMENT '未通过原因',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态 0:待审核 -1:审核未通过 1:待退货 2:待收货 3:已退款',
  `status_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '状态改变时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `reconciliation_id` int(10) unsigned DEFAULT '0' COMMENT '对账id',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_system_del` tinyint(1) DEFAULT '0' COMMENT '商户删除',
  `admin_id` int(10) unsigned DEFAULT '0' COMMENT '管理/客服ID',
  `user_type` tinyint(1) DEFAULT '1' COMMENT '用户类型 1 用户 2 平台 3 商户 4 客服 ',
  PRIMARY KEY (`refund_order_id`) USING BTREE,
  UNIQUE KEY `refund_order_sn` (`refund_order_sn`) USING BTREE,
  KEY `oid` (`order_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单退款表';

DROP TABLE IF EXISTS `qixi_store_refund_product`;
CREATE TABLE `qixi_store_refund_product` (

  `refund_product_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单产品 id',
  `refund_order_id` int(10) unsigned NOT NULL COMMENT '退款单',
  `order_product_id` int(10) unsigned NOT NULL COMMENT '订单产品id',
  `refund_price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '退款金额',
  `platform_refund_price` decimal(8,2) DEFAULT NULL COMMENT '平台券退款金额',
  `refund_postage` decimal(8,2) DEFAULT NULL COMMENT '退邮费金额',
  `refund_integral` int(10) unsigned DEFAULT '0' COMMENT '退金额',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '退货数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`refund_product_id`) USING BTREE,
  KEY `refund_order_id` (`refund_order_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='退款单产品表';

DROP TABLE IF EXISTS `qixi_store_refund_status`;
CREATE TABLE `qixi_store_refund_status` (

  `refund_order_id` int(10) unsigned NOT NULL COMMENT '退款单订单id',
  `change_type` varchar(32) NOT NULL COMMENT '操作类型',
  `change_message` varchar(256) NOT NULL COMMENT '操作备注',
  `change_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
  KEY `refund_order_id` (`refund_order_id`) USING BTREE,
  KEY `change_type` (`change_type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='订单操作记录表';

DROP TABLE IF EXISTS `qixi_store_seckill_active`;
CREATE TABLE `qixi_store_seckill_active` (

  `seckill_active_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` varchar(64) NOT NULL COMMENT '活动名称',
  `seckill_time_ids` varchar(255) DEFAULT '' COMMENT '活动场次',
  `start_day` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '开始日期',
  `end_day` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '结束日期',
  `start_time` int(10) unsigned NOT NULL COMMENT '开始时间',
  `end_time` int(10) unsigned NOT NULL COMMENT '结束时间',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户ID',
  `all_pay_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '活动有效期内每个用户可购买该商品总数限制',
  `once_pay_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '单次购买最大数量限制',
  `product_id` int(10) unsigned NOT NULL COMMENT '商品ID',
  `product_category_ids` varchar(255) DEFAULT NULL COMMENT '平台一级商品分类/为空均可参与',
  `merchant_count` int(11) NOT NULL DEFAULT '0' COMMENT '商户数量',
  `product_count` int(11) NOT NULL DEFAULT '0' COMMENT '商品数量',
  `active_status` enum('0','1','-1') DEFAULT '0' COMMENT '活动状态: 0未开始 1 进行中 -1 已结束',
  `sign` int(11) NOT NULL COMMENT '标识 1=秒杀活动',
  `status` tinyint(3) unsigned DEFAULT '0' COMMENT '0=未开启,1=已开启',
  `create_time` int(11) NOT NULL COMMENT '创建时间',
  `update_time` int(11) NOT NULL COMMENT '修改时间',
  `delete_time` int(11) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`seckill_active_id`) USING BTREE,
  KEY `start_day` (`start_day`,`end_day`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `status` (`status`,`active_status`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='商户设置秒杀商品关联表';

DROP TABLE IF EXISTS `qixi_store_seckill_time`;
CREATE TABLE `qixi_store_seckill_time` (

  `seckill_time_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) DEFAULT '',
  `start_time` int(10) unsigned NOT NULL COMMENT '开始时间',
  `end_time` int(10) unsigned NOT NULL COMMENT '结束时间',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '1，0状态',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `pic` varchar(255) DEFAULT NULL COMMENT '图片',
  PRIMARY KEY (`seckill_time_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='秒杀时间段配置';

DROP TABLE IF EXISTS `qixi_store_service`;
CREATE TABLE `qixi_store_service` (

  `service_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '客服id',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `uid` int(11) NOT NULL COMMENT '客服uid',
  `avatar` varchar(250) NOT NULL COMMENT '客服头像',
  `nickname` varchar(50) NOT NULL COMMENT '客服名称',
  `account` varchar(32) DEFAULT NULL COMMENT '客服账号',
  `pwd` varchar(64) DEFAULT NULL COMMENT '客服密码',
  `is_open` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '开启 pc 登录',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0隐藏1显示',
  `notify` int(11) DEFAULT '0' COMMENT '订单通知1开启0关闭',
  `phone` varchar(18) DEFAULT '' COMMENT '电话',
  `customer` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否展示统计管理',
  `is_verify` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否有核销权限',
  `is_goods` tinyint(3) unsigned DEFAULT '0' COMMENT '是否有商品管理权限',
  `sort` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`service_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `account` (`account`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客服表';

DROP TABLE IF EXISTS `qixi_store_service_log`;
CREATE TABLE `qixi_store_service_log` (

  `service_log_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '客服用户对话记录表ID',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户id',
  `msn` varchar(200) CHARACTER SET utf8mb4 NOT NULL COMMENT '消息内容',
  `uid` int(11) NOT NULL COMMENT '发送人uid',
  `service_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '客服 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '发送时间',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读（0：否；1：是；）',
  `service_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '客服是否已读（0：否；1：是；）',
  `remind` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否提醒过（0：否；1：是；）',
  `send_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:用户发送 1:客服回复',
  `msn_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '消息类型 1=文字 2=表情 3=图片 4=商品 5=订单 6=退款单',
  PRIMARY KEY (`service_log_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `service_id` (`service_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='客服用户对话记录表';

DROP TABLE IF EXISTS `qixi_store_service_reply`;
CREATE TABLE `qixi_store_service_reply` (

  `service_reply_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1:文字 2:图片',
  `keyword` varchar(64) NOT NULL COMMENT '回复的关键字',
  `content` varchar(512) NOT NULL COMMENT '回复内容',
  `status` tinyint(3) unsigned DEFAULT '1' COMMENT '是否开启',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`service_reply_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_service_user`;
CREATE TABLE `qixi_store_service_user` (

  `service_user_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '聊天用户 id',
  `service_id` int(10) unsigned NOT NULL COMMENT '客服 id',
  `uid` int(11) NOT NULL COMMENT '用户 id',
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `is_online` tinyint(3) unsigned DEFAULT '0' COMMENT '是否在线',
  `service_unread` smallint(5) unsigned DEFAULT '0' COMMENT '客服未读数',
  `user_unread` smallint(5) unsigned DEFAULT '0' COMMENT '用户未读数',
  `last_log_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后一条记录 id',
  `last_time` datetime NOT NULL COMMENT '最后发送时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`service_user_id`) USING BTREE,
  KEY `service_id` (`service_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE,
  KEY `user_id` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_store_spu`;
CREATE TABLE `qixi_store_spu` (

  `spu_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户ID',
  `product_id` int(10) unsigned DEFAULT '0' COMMENT '商品ID',
  `product_type` int(10) unsigned DEFAULT '0' COMMENT '活动类型0普通，1秒杀，2预售，3助力',
  `activity_id` int(10) unsigned DEFAULT '0' COMMENT '活动ID',
  `status` int(11) DEFAULT '0' COMMENT '0.下架，1.上架',
  `store_name` varchar(128) DEFAULT NULL COMMENT '商品名称',
  `ot_price` decimal(10,2) unsigned DEFAULT '0.00',
  `keyword` varchar(255) DEFAULT NULL COMMENT '关键词',
  `price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '最低价格',
  `rank` int(11) DEFAULT NULL COMMENT '排序',
  `create_time` datetime DEFAULT NULL,
  `temp_id` int(10) unsigned DEFAULT '0' COMMENT '运费模板',
  `sort` int(10) unsigned DEFAULT '0' COMMENT '商户排序',
  `star` int(11) DEFAULT '1' COMMENT '星级',
  `image` varchar(255) DEFAULT NULL COMMENT '主图',
  `is_del` int(10) unsigned DEFAULT '0',
  `mer_labels` varchar(255) DEFAULT '' COMMENT '标签id',
  `sys_labels` varchar(255) DEFAULT '' COMMENT '标签id',
  PRIMARY KEY (`spu_id`),
  KEY `mer_id` (`mer_id`,`product_id`),
  KEY `activity_id` (`activity_id`,`product_type`),
  KEY `product_id` (`product_id`,`product_type`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8 COMMENT='商品搜索信息表';

DROP TABLE IF EXISTS `qixi_system_admin`;
CREATE TABLE `qixi_system_admin` (

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

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='后台管理员表';

DROP TABLE IF EXISTS `qixi_system_attachment`;
CREATE TABLE `qixi_system_attachment` (

  `attachment_id` int(11) NOT NULL AUTO_INCREMENT,
  `attachment_category_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '分类ID 0编辑器,1产品图片,2拼团图片,3砍价图片,4秒杀图片,5文章图片,6组合数据图',
  `attachment_name` varchar(100) NOT NULL COMMENT '附件名称',
  `attachment_src` varchar(200) NOT NULL COMMENT '附件路径',
  `upload_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '图片上传类型 1本地 2七牛云 3OSS 4COS ',
  `user_type` int(11) NOT NULL DEFAULT '0' COMMENT '图片上传模块类型 0总后台后台  >0商户后台  -1用户生成',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '上传用户的 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `attachment_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '素材类型',
  PRIMARY KEY (`attachment_id`) USING BTREE,
  KEY `attachment_category_id` (`attachment_category_id`) USING BTREE,
  KEY `user_type` (`user_type`,`user_id`,`upload_type`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=94 DEFAULT CHARSET=utf8 COMMENT='附件管理表';

DROP TABLE IF EXISTS `qixi_system_attachment_category`;
CREATE TABLE `qixi_system_attachment_category` (

  `attachment_category_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  `path` varchar(512) NOT NULL DEFAULT '' COMMENT '路径',
  `attachment_category_name` varchar(32) NOT NULL COMMENT '分类名称',
  `attachment_category_enname` varchar(16) NOT NULL COMMENT '分类目录',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`attachment_category_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='附件分类表';

DROP TABLE IF EXISTS `qixi_system_config`;
CREATE TABLE `qixi_system_config` (

  `config_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置id',
  `config_classify_id` int(10) unsigned NOT NULL COMMENT '配置分类id',
  `config_name` varchar(64) NOT NULL COMMENT '字段名称',
  `config_key` varchar(64) NOT NULL COMMENT '字段 key',
  `config_type` varchar(20) NOT NULL DEFAULT 'input' COMMENT '配置类型',
  `config_rule` varchar(255) DEFAULT NULL COMMENT '规则',
  `config_props` varchar(255) DEFAULT '' COMMENT '配置',
  `required` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '必填',
  `info` varchar(255) DEFAULT '' COMMENT '配置说明',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `user_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0=总后台配置 1=商户后台配置',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否显示',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `linked_status` tinyint(3) unsigned DEFAULT '0' COMMENT '是否开启联动显示 0/1',
  `linked_id` int(11) DEFAULT '0' COMMENT '联动显示的id 信息',
  `linked_value` int(11) DEFAULT '0' COMMENT '联动显示的值',
  PRIMARY KEY (`config_id`) USING BTREE,
  UNIQUE KEY `config_name` (`config_key`) USING BTREE,
  KEY `config_classify_id` (`config_classify_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=574 DEFAULT CHARSET=utf8 COMMENT='配置表';

DROP TABLE IF EXISTS `qixi_system_config_classify`;
CREATE TABLE `qixi_system_config_classify` (

  `config_classify_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置分类id',
  `pid` int(11) DEFAULT '0' COMMENT '父级ID',
  `classify_name` varchar(255) NOT NULL COMMENT '配置分类名称',
  `classify_key` varchar(255) NOT NULL COMMENT '配置分类英文名称',
  `info` varchar(30) DEFAULT NULL COMMENT '配置分类说明',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `icon` varchar(30) DEFAULT NULL COMMENT '图标',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '配置分类状态',
  PRIMARY KEY (`config_classify_id`) USING BTREE,
  UNIQUE KEY `classify_key` (`classify_key`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=114 DEFAULT CHARSET=utf8 COMMENT='配置分类表';

DROP TABLE IF EXISTS `qixi_system_config_value`;
CREATE TABLE `qixi_system_config_value` (

  `config_value_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置id',
  `config_key` varchar(32) NOT NULL COMMENT '配置分类key',
  `value` text NOT NULL COMMENT '值',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`config_value_id`) USING BTREE,
  UNIQUE KEY `config_name` (`config_key`,`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=1511 DEFAULT CHARSET=utf8 COMMENT='配置表';

DROP TABLE IF EXISTS `qixi_system_form`;
CREATE TABLE `qixi_system_form` (

  `form_id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '表单标题',
  `form_keys` text COMMENT '表单所有的key',
  `value` text COMMENT '表单内容',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态',
  `is_del` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `mer_id` int(11) NOT NULL DEFAULT '0' COMMENT '商户ID',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`form_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_system_group`;
CREATE TABLE `qixi_system_group` (

  `group_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '组合数据ID',
  `group_name` varchar(50) NOT NULL COMMENT '数据组名称',
  `group_info` varchar(256) NOT NULL COMMENT '数据提示',
  `group_key` varchar(50) NOT NULL COMMENT '数据字段',
  `fields` text COMMENT '数据组字段以及类型（json数据）',
  `user_type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0=总后台配置 1=商户后台配置',
  `sort` smallint(5) unsigned NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_id`) USING BTREE,
  UNIQUE KEY `group_key` (`group_key`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=104 DEFAULT CHARSET=utf8 COMMENT='组合数据表';

DROP TABLE IF EXISTS `qixi_system_group_data`;
CREATE TABLE `qixi_system_group_data` (

  `group_data_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '组合数据详情ID',
  `group_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '对应的数据组id',
  `value` text NOT NULL COMMENT '数据组对应的数据值（json数据）',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '数据排序',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态（1：开启；0：关闭；）',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加数据时间',
  PRIMARY KEY (`group_data_id`) USING BTREE,
  KEY `group_id` (`group_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=935 DEFAULT CHARSET=utf8 COMMENT='组合数据详情表';

DROP TABLE IF EXISTS `qixi_system_log`;
CREATE TABLE `qixi_system_log` (

  `log_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '管理员操作记录ID',
  `admin_id` int(10) unsigned NOT NULL COMMENT '管理员id',
  `admin_name` varchar(64) NOT NULL COMMENT '管理员姓名',
  `route` varchar(128) NOT NULL COMMENT '路由',
  `method` varchar(12) NOT NULL COMMENT '方式',
  `url` varchar(256) NOT NULL COMMENT '链接',
  `ip` varchar(16) NOT NULL COMMENT '登录IP',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`log_id`) USING BTREE,
  KEY `admin_id` (`admin_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=248 DEFAULT CHARSET=utf8 COMMENT='管理员操作记录表';

DROP TABLE IF EXISTS `qixi_system_menu`;
CREATE TABLE `qixi_system_menu` (

  `menu_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级id',
  `path` varchar(512) NOT NULL COMMENT '路径',
  `icon` varchar(32) DEFAULT '' COMMENT '图标',
  `menu_name` varchar(128) NOT NULL DEFAULT '' COMMENT '按钮名',
  `route` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL COMMENT '路由名称',
  `params` varchar(128) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '参数',
  `sort` tinyint(4) NOT NULL DEFAULT '1' COMMENT '排序',
  `is_show` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否显示',
  `is_mer` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '模块，1 平台， 2商户',
  `is_menu` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '类型，1菜单 2 权限',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_agent` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0:平台,1:区域,2:商户',
  PRIMARY KEY (`menu_id`) USING BTREE,
  KEY `pid` (`pid`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=10186 DEFAULT CHARSET=utf8 COMMENT='菜单表';

DROP TABLE IF EXISTS `qixi_system_notice`;
CREATE TABLE `qixi_system_notice` (

  `notice_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `admin_id` int(10) unsigned NOT NULL COMMENT '管理员 id',
  `notice_title` varchar(128) NOT NULL COMMENT '通知标题',
  `notice_content` text NOT NULL COMMENT '通知内容',
  `type` tinyint(3) unsigned NOT NULL COMMENT '通知类型',
  `type_str` varchar(512) NOT NULL COMMENT '通知说明',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0:正常 1:删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` int(10) NOT NULL DEFAULT '1' COMMENT '状态（0:关;1:开）',
  PRIMARY KEY (`notice_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户公告';

DROP TABLE IF EXISTS `qixi_system_notice_config`;
CREATE TABLE `qixi_system_notice_config` (

  `notice_config_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `notice_title` varchar(20) DEFAULT NULL COMMENT '消息名称',
  `const_key` varchar(100) DEFAULT NULL COMMENT '通知标识',
  `notice_info` varchar(50) DEFAULT NULL COMMENT '消息说明',
  `notice_sys` tinyint(4) DEFAULT '-1' COMMENT '站内消息',
  `notice_sms` tinyint(4) DEFAULT '-1' COMMENT '短信消息',
  `notice_wechat` tinyint(4) DEFAULT '-1' COMMENT '公众号模板消息',
  `wechat_tempkey` varchar(100) DEFAULT NULL COMMENT '微信模板关联ID',
  `wechat_content` varchar(255) NOT NULL DEFAULT '' COMMENT '微信模板内容',
  `wechat_tempid` varchar(255) DEFAULT NULL COMMENT '微信模板ID',
  `notice_routine` tinyint(4) DEFAULT '-1' COMMENT '小程序订阅消息',
  `routine_tempkey` varchar(100) DEFAULT NULL COMMENT '订阅消息关联ID',
  `routine_content` varchar(255) NOT NULL DEFAULT '' COMMENT '小程序订阅消息内容',
  `routine_tempid` varchar(255) DEFAULT NULL COMMENT '小程序消息ID',
  `type` int(10) unsigned DEFAULT '0' COMMENT '1商户通知， 0用户通知',
  `sms_tempid` varchar(50) DEFAULT NULL COMMENT '一号通短信模板ID',
  `sms_ali_tempid` varchar(50) DEFAULT NULL COMMENT '阿里云短信模板ID',
  `sms_content` varchar(100) DEFAULT NULL COMMENT '阿里云短信模板内容',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间 ',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `kid` char(100) DEFAULT '0',
  PRIMARY KEY (`notice_config_id`) USING BTREE,
  UNIQUE KEY `notic_config_id` (`notice_config_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_system_notice_log`;
CREATE TABLE `qixi_system_notice_log` (

  `notice_log_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL COMMENT '商户 id',
  `notice_id` int(10) unsigned NOT NULL COMMENT '公告 id',
  `is_read` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否已读',
  `read_time` timestamp NULL DEFAULT NULL COMMENT '读取时间',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  PRIMARY KEY (`notice_log_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_system_role`;
CREATE TABLE `qixi_system_role` (

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

) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8 COMMENT='身份管理表';

DROP TABLE IF EXISTS `qixi_system_storage`;
CREATE TABLE `qixi_system_storage` (

  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `type` tinyint(4) NOT NULL DEFAULT '1',
  `access_key` varchar(100) NOT NULL DEFAULT '',
  `name` varchar(100) NOT NULL DEFAULT '',
  `region` varchar(100) NOT NULL DEFAULT '',
  `acl` varchar(100) NOT NULL DEFAULT '',
  `domain` varchar(255) NOT NULL DEFAULT '',
  `cdn` varchar(255) DEFAULT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `is_del` tinyint(4) NOT NULL DEFAULT '0',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `qixi_template_message`;
CREATE TABLE `qixi_template_message` (

  `template_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '模板id',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0=订阅消息,1=微信模板消息',
  `tempkey` char(50) NOT NULL COMMENT '模板编号',
  `name` char(100) NOT NULL COMMENT '模板名',
  `content` varchar(1000) NOT NULL COMMENT '回复内容',
  `tempid` char(100) DEFAULT NULL COMMENT '模板ID',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '添加时间',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '状态',
  `kid` char(100) DEFAULT '0',
  PRIMARY KEY (`template_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信模板';

DROP TABLE IF EXISTS `qixi_user`;
CREATE TABLE `qixi_user` (

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
  `avatar` varchar(256) NOT NULL COMMENT '用户头像',
  `phone` char(15) DEFAULT NULL COMMENT '手机号码',
  `addres` varchar(128) DEFAULT NULL COMMENT '地址',
  `cancel_time` timestamp NULL DEFAULT NULL COMMENT '注销时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `last_time` timestamp NULL DEFAULT NULL COMMENT '最后一次登录时间',
  `last_ip` varchar(16) NOT NULL COMMENT '最后一次登录ip',
  `now_money` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '用户余额',
  `brokerage_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '佣金金额',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1为正常，0为禁止',
  `spread_uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '推广员id',
  `spread_time` timestamp NULL DEFAULT NULL COMMENT '推广员关联时间',
  `spread_limit` timestamp NULL DEFAULT NULL COMMENT '推广员到期时间',
  `brokerage_level` int(10) unsigned DEFAULT '0' COMMENT '推广员等级',
  `user_type` varchar(32) NOT NULL COMMENT '用户类型',
  `promoter_time` timestamp NULL DEFAULT NULL COMMENT '成功推广时间',
  `is_promoter` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否为推广员',
  `main_uid` int(10) unsigned DEFAULT '0' COMMENT '主账号',
  `pay_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户购买次数',
  `pay_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '用户消费金额',
  `spread_count` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '下级人数',
  `spread_pay_count` int(10) unsigned DEFAULT '0' COMMENT '下级订单数',
  `spread_pay_price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '下级订单金额',
  `integral` int(10) unsigned DEFAULT '0' COMMENT '积分',
  `member_level` int(10) unsigned DEFAULT '0' COMMENT '免费会员等级',
  `member_value` int(10) unsigned DEFAULT '0' COMMENT '免费会员成长值',
  `count_start` int(10) unsigned DEFAULT '0' COMMENT '用户获赞数',
  `count_fans` int(10) unsigned DEFAULT '0' COMMENT '用户粉丝数',
  `count_content` int(10) unsigned DEFAULT '0' COMMENT '用户内容数量',
  `is_svip` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否为付费会员 -1未开通过 0到期 1体验卡 2 有效期 3 永久',
  `svip_endtime` timestamp NULL DEFAULT NULL COMMENT '会员结束时间',
  `svip_save_money` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '会员节省金额',
  `promoter_switch` tinyint(1) NOT NULL DEFAULT '1' COMMENT '分销资格 0无 1有',
  PRIMARY KEY (`uid`) USING BTREE,
  KEY `account` (`account`) USING BTREE,
  KEY `spreaduid` (`spread_uid`) USING BTREE,
  KEY `wechat_user_id` (`wechat_user_id`) USING BTREE,
  KEY `main_uid` (`main_uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';

DROP TABLE IF EXISTS `qixi_user_address`;
CREATE TABLE `qixi_user_address` (

  `address_id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户地址id',
  `uid` int(10) unsigned NOT NULL COMMENT '用户id',
  `real_name` varchar(32) NOT NULL DEFAULT '' COMMENT '收货人姓名',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '收货人电话',
  `province` varchar(64) NOT NULL DEFAULT '' COMMENT '收货人所在省',
  `province_id` int(10) unsigned DEFAULT '0' COMMENT '省 id',
  `city` varchar(64) NOT NULL DEFAULT '' COMMENT '收货人所在市',
  `city_id` int(11) NOT NULL DEFAULT '0' COMMENT '城市id',
  `district` varchar(64) NOT NULL DEFAULT '' COMMENT '收货人所在区',
  `district_id` int(10) unsigned DEFAULT '0' COMMENT '区域 id',
  `street` varchar(64) DEFAULT NULL COMMENT '街/镇',
  `street_id` int(10) unsigned DEFAULT '0' COMMENT '街镇 id',
  `detail` varchar(256) NOT NULL DEFAULT '' COMMENT '收货人详细地址',
  `post_code` int(10) unsigned NOT NULL COMMENT '邮编',
  `longitude` varchar(16) NOT NULL DEFAULT '0' COMMENT '经度',
  `latitude` varchar(16) NOT NULL DEFAULT '0' COMMENT '纬度',
  `is_default` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否默认',
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否删除',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `tourist_unique_key` varchar(20) NOT NULL DEFAULT '' COMMENT '游客唯一标识',
  PRIMARY KEY (`address_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户地址表';

DROP TABLE IF EXISTS `qixi_user_bill`;
CREATE TABLE `qixi_user_bill` (

  `bill_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户账单id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户uid',
  `link_id` varchar(32) NOT NULL DEFAULT '0' COMMENT '关联id',
  `pm` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0 = 支出 1 = 获得',
  `title` varchar(64) NOT NULL COMMENT '账单标题',
  `category` varchar(64) NOT NULL COMMENT '明细种类',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '明细类型',
  `number` decimal(11,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '明细数字',
  `balance` decimal(11,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '剩余',
  `mark` varchar(512) NOT NULL COMMENT '备注',
  `mer_id` int(10) unsigned DEFAULT '0' COMMENT '商户 id',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0 = 待确定 1 = 有效 -1 = 无效',
  PRIMARY KEY (`bill_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE,
  KEY `type` (`category`,`type`,`link_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE

) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='用户账单表';

DROP TABLE IF EXISTS `qixi_user_brokerage`;
CREATE TABLE `qixi_user_brokerage` (

  `user_brokerage_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `brokerage_level` tinyint(3) unsigned NOT NULL COMMENT '等级',
  `brokerage_name` varchar(32) NOT NULL COMMENT 'vip 名称',
  `brokerage_icon` varchar(128) NOT NULL COMMENT 'vip 图标',
  `brokerage_rule` varchar(1500) NOT NULL COMMENT '升级规则',
  `user_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'vip 人数',
  `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '一级佣金',
  `extension_two` decimal(8,2) unsigned NOT NULL COMMENT '二级佣金',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '默认0分销会员等级，1 免费会员等级',
  PRIMARY KEY (`user_brokerage_id`) USING BTREE,
  UNIQUE KEY `vip_level` (`brokerage_level`,`type`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_user_extract`;
CREATE TABLE `qixi_user_extract` (

  `extract_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '用户 id',
  `extract_sn` varchar(255) DEFAULT NULL,
  `real_name` varchar(64) DEFAULT NULL COMMENT '姓名',
  `extract_type` tinyint(1) DEFAULT '0' COMMENT '0 银行卡 1 支付宝 2微信 3 零钱',
  `bank_code` varchar(32) DEFAULT '0' COMMENT '银行卡',
  `bank_address` varchar(256) DEFAULT '' COMMENT '开户地址',
  `alipay_code` varchar(64) DEFAULT '' COMMENT '支付宝账号',
  `wechat` varchar(15) DEFAULT NULL COMMENT '微信号',
  `extract_pic` varchar(128) DEFAULT NULL COMMENT '收款码',
  `extract_price` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '提现金额',
  `balance` decimal(8,2) unsigned DEFAULT '0.00' COMMENT '余额',
  `mark` varchar(512) DEFAULT NULL COMMENT '管理员备注',
  `admin_id` int(11) DEFAULT '0' COMMENT '审核管理员',
  `fail_msg` varchar(128) DEFAULT NULL COMMENT '无效原因',
  `status_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '无效时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `status` tinyint(4) DEFAULT '0' COMMENT '-1 未通过 0 审核中 1 已提现',
  `bank_name` varchar(255) DEFAULT NULL COMMENT '银行名称',
  `wechat_status` varchar(50) NOT NULL DEFAULT '' COMMENT '微信转账状态：ACCEPTED: 转账已受理;PROCESSING: 转账锁定资金中;WAIT_USER_CONFIRM: 待收款用户确认;TRANSFERING: 转账中;SUCCESS: 转账成功;FAIL: 转账失败;CANCELING: 撤销中;CANCELLED: 转账撤销完成',
  `package_info` varchar(255) NOT NULL DEFAULT '' COMMENT '跳转微信支付收款页的package信息',
  `transfer_bill_no` varchar(100) NOT NULL DEFAULT '' COMMENT '微信转账单号',
  `wechat_app_id` varchar(100) NOT NULL DEFAULT '' COMMENT '微信appid',
  `wechat_mch_id` varchar(100) NOT NULL DEFAULT '' COMMENT '微信mchid',
  PRIMARY KEY (`extract_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户提现表';

DROP TABLE IF EXISTS `qixi_user_fields`;
CREATE TABLE `qixi_user_fields` (

  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL DEFAULT '0' COMMENT '用户id',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_user_group`;
CREATE TABLE `qixi_user_group` (

  `group_id` smallint(5) unsigned NOT NULL AUTO_INCREMENT,
  `group_name` varchar(64) NOT NULL COMMENT '用户分组名称',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`group_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户分组表';

DROP TABLE IF EXISTS `qixi_user_history`;
CREATE TABLE `qixi_user_history` (

  `user_history_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `res_id` int(10) unsigned DEFAULT NULL COMMENT '历史记录对象的ID',
  `res_type` int(11) DEFAULT NULL COMMENT '历史记录类型',
  `uid` int(11) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` int(11) DEFAULT NULL,
  PRIMARY KEY (`user_history_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='浏览记录表';

DROP TABLE IF EXISTS `qixi_user_info`;
CREATE TABLE `qixi_user_info` (

  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `field` varchar(255) NOT NULL DEFAULT '' COMMENT '字段',
  `title` varchar(255) NOT NULL DEFAULT '字段名' COMMENT '字段名',
  `is_used` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否使用',
  `is_require` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否必填',
  `is_show` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否在用户端展示',
  `type` varchar(255) NOT NULL DEFAULT '' COMMENT '信息格式',
  `msg` varchar(255) NOT NULL DEFAULT '' COMMENT '提示信息',
  `content` varchar(255) DEFAULT NULL COMMENT '配置内容',
  `is_default` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否系统默认字段',
  `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `id` (`id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_user_label`;
CREATE TABLE `qixi_user_label` (

  `label_id` int(11) NOT NULL AUTO_INCREMENT,
  `label_name` varchar(255) NOT NULL DEFAULT '' COMMENT '标签名称',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0=手动标签 1=自动标签',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`label_id`) USING BTREE,
  KEY `mer_id` (`mer_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户标签表';

DROP TABLE IF EXISTS `qixi_user_merchant`;
CREATE TABLE `qixi_user_merchant` (

  `user_merchant_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL COMMENT '用户 id',
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `first_pay_time` timestamp NULL DEFAULT NULL COMMENT '首次消费时间',
  `last_pay_time` timestamp NULL DEFAULT NULL COMMENT '最后一次消费时间',
  `pay_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '消费次数',
  `pay_price` decimal(10,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '消费金额',
  `label_id` varchar(256) DEFAULT NULL COMMENT '用户标签',
  `last_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后一次访问时间',
  `status` tinyint(3) unsigned DEFAULT '1' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_merchant_id`),
  UNIQUE KEY `uid` (`uid`,`mer_id`) USING BTREE,
  KEY `mer_id` (`mer_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商户用户表';

DROP TABLE IF EXISTS `qixi_user_order`;
CREATE TABLE `qixi_user_order` (

  `order_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `link_id` int(11) DEFAULT NULL COMMENT '关联ID',
  `pay_type` varchar(10) DEFAULT NULL COMMENT '支付方式：1微信，2支付宝',
  `title` varchar(50) NOT NULL,
  `order_sn` varchar(50) DEFAULT NULL COMMENT '订单ID',
  `pay_price` decimal(8,2) unsigned DEFAULT NULL COMMENT '价格',
  `order_info` varchar(255) DEFAULT NULL COMMENT '订单信息信息',
  `order_type` varchar(255) DEFAULT NULL COMMENT '订单类型 S 付费会员 ',
  `paid` int(11) NOT NULL DEFAULT '0' COMMENT '支付状态',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '支付时间',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态：默认1',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `create_time` timestamp NULL DEFAULT NULL,
  `is_del` int(11) NOT NULL DEFAULT '0',
  `admin_id` int(10) unsigned DEFAULT NULL COMMENT '管理员ID',
  `other` varchar(50) DEFAULT NULL COMMENT '其他参数',
  `end_time` timestamp NULL DEFAULT NULL,
  `transaction_id` varchar(255) DEFAULT '',
  PRIMARY KEY (`order_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='支付订单信息';

DROP TABLE IF EXISTS `qixi_user_receipt`;
CREATE TABLE `qixi_user_receipt` (

  `user_receipt_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `receipt_type` tinyint(1) DEFAULT '0' COMMENT '发票类型：1.普通发票，2.增值税发票',
  `receipt_title` varchar(128) DEFAULT '' COMMENT '发票抬头',
  `receipt_title_type` varchar(255) DEFAULT '0' COMMENT '发票抬头类型：1.个人，2.企业',
  `duty_paragraph` varchar(255) DEFAULT '' COMMENT '税号',
  `email` varchar(255) DEFAULT '' COMMENT '邮箱',
  `bank_name` varchar(255) DEFAULT '' COMMENT '开户行',
  `bank_code` varchar(255) DEFAULT '0' COMMENT '银行账号',
  `address` varchar(255) DEFAULT '' COMMENT '企业地址',
  `tel` varchar(255) DEFAULT '' COMMENT '企业电话',
  `is_default` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否默认',
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `is_del` tinyint(1) DEFAULT '0',
  `create_time` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`user_receipt_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户发票信息';

DROP TABLE IF EXISTS `qixi_user_recharge`;
CREATE TABLE `qixi_user_recharge` (

  `recharge_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL COMMENT '充值用户UID',
  `order_id` varchar(32) NOT NULL COMMENT '订单号',
  `price` decimal(8,2) unsigned NOT NULL DEFAULT '0.00' COMMENT '充值金额',
  `give_price` decimal(8,2) NOT NULL DEFAULT '0.00' COMMENT '购买赠送金额',
  `recharge_type` varchar(32) NOT NULL COMMENT '充值类型',
  `paid` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否充值',
  `pay_time` timestamp NULL DEFAULT NULL COMMENT '充值支付时间',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '充值时间',
  `refund_price` decimal(10,2) unsigned DEFAULT '0.00' COMMENT '退款金额',
  PRIMARY KEY (`recharge_id`) USING BTREE,
  UNIQUE KEY `order_id` (`order_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户充值表';

DROP TABLE IF EXISTS `qixi_user_relation`;
CREATE TABLE `qixi_user_relation` (

  `uid` int(10) unsigned NOT NULL COMMENT '用户ID',
  `type_id` int(10) unsigned NOT NULL COMMENT '类型的 id',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '关联类型(0 普通商品、1秒杀2、预售3、助力4、拼团、10 = 店铺、12=购买过)',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  UNIQUE KEY `uid` (`type`,`type_id`,`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户记录表';

DROP TABLE IF EXISTS `qixi_user_sign`;
CREATE TABLE `qixi_user_sign` (

  `sign_id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL DEFAULT '0' COMMENT '用户uid',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '签到说明',
  `number` int(11) NOT NULL DEFAULT '0' COMMENT '获得积分',
  `integral` int(11) NOT NULL DEFAULT '0' COMMENT '剩余积分',
  `sign_num` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '连续签到天数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`sign_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='签到记录表';

DROP TABLE IF EXISTS `qixi_user_spread_log`;
CREATE TABLE `qixi_user_spread_log` (

  `user_spread_log_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL COMMENT 'uid',
  `old_spread_uid` int(10) unsigned NOT NULL COMMENT '原来的推荐人uid',
  `spread_uid` int(10) unsigned NOT NULL COMMENT '新的推荐人 uid',
  `admin_id` int(10) unsigned NOT NULL COMMENT '修改的管理员',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_spread_log_id`),
  KEY `uid` (`uid`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `qixi_user_visit`;
CREATE TABLE `qixi_user_visit` (

  `user_visit_id` int(11) NOT NULL AUTO_INCREMENT,
  `uid` int(11) DEFAULT NULL COMMENT '用户ID',
  `type` varchar(32) NOT NULL COMMENT '记录类型',
  `type_id` int(11) NOT NULL DEFAULT '0' COMMENT '商品ID',
  `content` varchar(255) DEFAULT NULL COMMENT '备注描述',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`user_visit_id`) USING BTREE,
  KEY `uid` (`uid`) USING BTREE,
  KEY `type` (`type`,`type_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='商品浏览分析表';

DROP TABLE IF EXISTS `qixi_wechat_news`;
CREATE TABLE `qixi_wechat_news` (

  `wechat_news_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '图文消息管理ID',
  `mer_id` int(11) DEFAULT '0' COMMENT '商户id',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`wechat_news_id`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='图文消息管理表';

DROP TABLE IF EXISTS `qixi_wechat_qrcode`;
CREATE TABLE `qixi_wechat_qrcode` (

  `wechat_qrcode_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '微信二维码ID',
  `third_type` varchar(32) NOT NULL COMMENT '二维码类型',
  `third_id` int(10) unsigned NOT NULL COMMENT '类型id',
  `ticket` varchar(255) NOT NULL COMMENT '二维码参数',
  `expire_seconds` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '二维码有效时间',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态',
  `url` varchar(255) NOT NULL COMMENT '微信访问url',
  `qrcode_url` varchar(255) NOT NULL COMMENT '微信二维码url',
  `scan` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '被扫的次数',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`wechat_qrcode_id`) USING BTREE,
  UNIQUE KEY `third_type` (`third_type`,`third_id`) USING BTREE,
  KEY `ticket` (`ticket`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信二维码管理表';

DROP TABLE IF EXISTS `qixi_wechat_reply`;
CREATE TABLE `qixi_wechat_reply` (

  `wechat_reply_id` mediumint(8) unsigned NOT NULL AUTO_INCREMENT COMMENT '微信关键字回复id',
  `key` varchar(64) NOT NULL COMMENT '关键字',
  `type` varchar(32) NOT NULL COMMENT '回复类型',
  `data` text NOT NULL COMMENT '回复数据',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '0=不可用  1 =可用',
  `hidden` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '是否显示',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`wechat_reply_id`) USING BTREE,
  UNIQUE KEY `key` (`key`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信关键字回复表';

DROP TABLE IF EXISTS `qixi_wechat_user`;
CREATE TABLE `qixi_wechat_user` (

  `wechat_user_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '微信用户id',
  `unionid` varchar(60) DEFAULT NULL COMMENT '只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段',
  `openid` varchar(30) DEFAULT NULL COMMENT '用户的标识，对当前公众号唯一',
  `routine_openid` varchar(32) DEFAULT NULL COMMENT '小程序唯一身份ID',
  `nickname` varchar(64) NOT NULL COMMENT '用户的昵称',
  `headimgurl` varchar(256) NOT NULL COMMENT '用户头像',
  `sex` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '用户的性别，值为1时是男性，值为2时是女性，值为0时是未知',
  `city` varchar(32) NOT NULL COMMENT '用户所在城市',
  `language` varchar(32) NOT NULL COMMENT '用户的语言，简体中文为zh_CN',
  `province` varchar(32) NOT NULL COMMENT '用户所在省份',
  `country` varchar(32) NOT NULL COMMENT '用户所在国家',
  `remark` varchar(256) DEFAULT NULL COMMENT '公众号运营者对粉丝的备注，公众号运营者可在微信公众平台用户管理界面对粉丝添加备注',
  `groupid` smallint(5) unsigned DEFAULT '0' COMMENT '用户所在的分组ID（兼容旧的用户分组接口）',
  `tagid_list` varchar(256) DEFAULT NULL COMMENT '用户被打上的标签ID列表',
  `subscribe` tinyint(3) unsigned DEFAULT '0' COMMENT '用户是否订阅该公众号标识',
  `subscribe_time` int(10) unsigned DEFAULT NULL COMMENT '关注公众号时间',
  `session_key` varchar(32) DEFAULT NULL COMMENT '小程序用户会话密匙',
  `user_type` varchar(32) DEFAULT 'wechat' COMMENT '用户类型',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`wechat_user_id`) USING BTREE,
  KEY `groupid` (`groupid`) USING BTREE,
  KEY `subscribe_time` (`subscribe_time`) USING BTREE,
  KEY `unionid` (`unionid`) USING BTREE,
  KEY `routine_openid` (`routine_openid`) USING BTREE,
  KEY `openid` (`openid`) USING BTREE

) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='微信用户表';

SET FOREIGN_KEY_CHECKS = 1;