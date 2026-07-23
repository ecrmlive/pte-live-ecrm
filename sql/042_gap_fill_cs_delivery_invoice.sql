-- 缺口补齐：客服会话桥 / 配送员 / 发票 / 短信配置 / 协议键 / 直播 URL
USE `qixi_mergers`;

-- 客服会话（本仓业务真相；IM 远程可选）
CREATE TABLE IF NOT EXISTS `qixi_cs_thread` (
  `thread_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `service_id` int(10) unsigned NOT NULL DEFAULT 0 COMMENT '承接坐席，0=未分配',
  `last_msg` varchar(255) NOT NULL DEFAULT '',
  `last_time` datetime DEFAULT NULL,
  `user_unread` int(10) unsigned NOT NULL DEFAULT 0,
  `service_unread` int(10) unsigned NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1进行中 2已关闭',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`thread_id`),
  UNIQUE KEY `uk_mer_uid` (`mer_id`,`uid`),
  KEY `idx_mer_service` (`mer_id`,`service_id`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客服会话';

CREATE TABLE IF NOT EXISTS `qixi_cs_message` (
  `msg_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `thread_id` int(10) unsigned NOT NULL DEFAULT 0,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `sender_role` varchar(16) NOT NULL DEFAULT '' COMMENT 'user|service',
  `sender_id` int(10) unsigned NOT NULL DEFAULT 0,
  `msg_type` varchar(16) NOT NULL DEFAULT 'text' COMMENT 'text|order|system',
  `content` text NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`msg_id`),
  KEY `idx_thread` (`thread_id`,`msg_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客服消息（本地桥）';

CREATE TABLE IF NOT EXISTS `qixi_im_identity` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `portal` varchar(16) NOT NULL DEFAULT '' COMMENT 'app|service',
  `local_id` int(10) unsigned NOT NULL DEFAULT 0,
  `im_user_id` varchar(64) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_portal_local` (`portal`,`local_id`),
  UNIQUE KEY `uk_im_user` (`im_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商城↔IM 身份映射';

-- 配送员
CREATE TABLE IF NOT EXISTS `qixi_delivery_staff` (
  `staff_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `phone` varchar(32) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1启用 0停用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`staff_id`),
  KEY `idx_mer` (`mer_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户配送员';

INSERT INTO `qixi_delivery_staff` (`staff_id`, `mer_id`, `name`, `phone`, `status`)
SELECT 1, 1, '演示骑手小栖', '13900000001', 1
WHERE EXISTS (SELECT 1 FROM `qixi_merchant` WHERE `mer_id` = 1)
  AND NOT EXISTS (SELECT 1 FROM `qixi_delivery_staff` WHERE `staff_id` = 1);

-- 发票申请
CREATE TABLE IF NOT EXISTS `qixi_store_order_invoice` (
  `invoice_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL DEFAULT 0,
  `order_id` int(10) unsigned NOT NULL DEFAULT 0,
  `mer_id` int(10) unsigned NOT NULL DEFAULT 0,
  `invoice_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1普通 2专用',
  `header_type` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1个人 2单位',
  `header` varchar(128) NOT NULL DEFAULT '',
  `tax_no` varchar(64) NOT NULL DEFAULT '',
  `email` varchar(128) NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0待开 1已开 -1驳回',
  `mark` varchar(255) NOT NULL DEFAULT '',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`invoice_id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_mer` (`mer_id`,`status`),
  KEY `idx_order` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单发票申请';

-- 直播播放/推流 URL stub
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_broadcast_room' AND COLUMN_NAME = 'play_url'
);
SET @sql := IF(@col_exists = 0,
  'ALTER TABLE `qixi_broadcast_room` ADD COLUMN `play_url` varchar(512) NOT NULL DEFAULT '''' COMMENT ''播放地址 stub'' AFTER `feeds_img`, ADD COLUMN `push_url` varchar(512) NOT NULL DEFAULT '''' COMMENT ''推流地址 stub'' AFTER `play_url`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

UPDATE `qixi_broadcast_room`
SET `play_url` = 'https://example.local/live/demo/play.m3u8',
    `push_url` = 'rtmp://example.local/live/demo'
WHERE `broadcast_room_id` = 1 AND (`play_url` = '' OR `play_url` IS NULL);

-- 协议键补齐（key ≤ 32）
INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_refund_agree', 0, '【退款协议】售后退款需符合商品状态与时效；审核通过后原路退回。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_refund_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_cancel_agree', 0, '【取消订单说明】未支付订单可取消；已支付订单按售后流程处理。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_cancel_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_recharge_agree', 0, '【充值协议】余额充值仅用于本平台消费，不支持提现到银行卡。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_recharge_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_integral_agree', 0, '【积分规则】积分可用于指定商品抵扣；清零与获取规则以平台配置为准。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_integral_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'mer_settle_agree', 0, '【商户结算说明】订单完成后按结算周期打款；手续费与抽成以入驻协议为准。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'mer_settle_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_lottery_agree', 0, '【抽奖活动说明】奖品以活动页展示为准；作弊账号取消获奖资格。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_lottery_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_deposit_agree', 0, '【保证金说明】商户保证金用于保障交易；违规可按规定扣减。演示文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_deposit_agree');

-- 短信配置 stub（JSON）
INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sms_config', 0, '{"enabled":false,"provider":"stub","sign":"栖息商城","template_login":"您的验证码是{code}","remark":"仅配置存储，未接真实短信通道"}'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sms_config');

-- 商户菜单：配送员 / 发票
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 170, 105, '/order/delivery-staff', '', '配送员', 'OrderDeliveryStaff', '', 45, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 170 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 171, 105, '/order/invoice', '', '发票管理', 'OrderInvoice', '', 44, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 171 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',170,171')
WHERE `mer_id` > 0 AND `rules` NOT LIKE '%170%';

-- 平台：短信配置页
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 172, 10, '/setting/sms', '', '短信配置', 'SettingSms', '', 5, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 172 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',172')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%172%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'gap-fill-042', '缺口：客服会话桥/配送员/发票/短信/协议键/直播URL'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'gap-fill-042');
