-- 平台优惠券「发送记录」批次表 + 用户券关联 send_id（对齐 CRMEB store_coupon_send）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_coupon_send` (
  `coupon_send_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '0=平台',
  `coupon_id` bigint unsigned NOT NULL,
  `coupon_num` int unsigned NOT NULL DEFAULT 0 COMMENT '发放数量',
  `mark` json DEFAULT NULL COMMENT '筛选条件 JSON：{type,search}',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '0发送中 1已完成',
  PRIMARY KEY (`coupon_send_id`),
  KEY `idx_mer_time` (`mer_id`,`create_time`),
  KEY `idx_coupon` (`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- coupon_user 增加 send_id（幂等）
SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_coupon_user'
    AND COLUMN_NAME = 'send_id'
);
SET @sql := IF(
  @col_exists = 0,
  'ALTER TABLE `qixi_crm_b_coupon_user` ADD COLUMN `send_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''发送批次'' AFTER `source`, ADD KEY `idx_send_id` (`send_id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 纠正演示券中文标题（避免本地乱码残留）
UPDATE `qixi_crm_b_store_coupon`
SET `title`='平台新客满99减10', `send_type`=6
WHERE `coupon_id`=9401;
UPDATE `qixi_crm_b_store_coupon`
SET `title`='平台夏日满299减40', `send_type`=6
WHERE `coupon_id`=9402;
UPDATE `qixi_crm_b_store_coupon`
SET `title`='测试隐藏券'
WHERE `coupon_id`=9403;

INSERT INTO `qixi_crm_b_store_coupon`
  (`coupon_id`,`mer_id`,`is_limited`,`total_count`,`remain_count`,`send_type`,`title`,`coupon_price`,`use_min_price`,`coupon_type`,`coupon_time`,`sort`,`status`,`type`,`is_del`)
VALUES
  (9404,0,0,0,0,6,'明月优惠券',50.00,0,0,365,5,1,10,0),
  (9405,0,0,0,0,6,'电影',20.00,100,0,10,4,1,10,0)
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `coupon_price`=VALUES(`coupon_price`),
  `use_min_price`=VALUES(`use_min_price`),
  `coupon_time`=VALUES(`coupon_time`),
  `send_type`=VALUES(`send_type`),
  `status`=VALUES(`status`),
  `type`=VALUES(`type`),
  `is_del`=0;

-- 发送批次演示（45 号批次 25 人，便于使用记录分页验收）
INSERT INTO `qixi_crm_b_store_coupon_send`
  (`coupon_send_id`,`mer_id`,`coupon_id`,`coupon_num`,`mark`,`create_time`,`status`)
VALUES
  (45,0,9401,25,JSON_OBJECT('type',1,'search',JSON_OBJECT('user_type',1)),'2026-07-16 10:58:39',1),
  (44,0,9404,2,NULL,'2026-07-15 14:20:11',1),
  (43,0,9405,1,JSON_OBJECT('type',1,'search',JSON_OBJECT('user_type',2)),'2026-07-14 09:12:03',1),
  (42,0,9402,2,NULL,'2026-07-10 16:40:00',1)
ON DUPLICATE KEY UPDATE
  `coupon_id`=VALUES(`coupon_id`),
  `coupon_num`=VALUES(`coupon_num`),
  `mark`=VALUES(`mark`),
  `create_time`=VALUES(`create_time`),
  `status`=VALUES(`status`);

-- 使用记录分页演示用户（send_id=45）
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`) VALUES
  (9101,'181****8838',NULL,1,0),
  (9102,'小程序用户',NULL,1,0),
  (9103,'微信用户',NULL,1,0),
  (9210,'181****1001',NULL,1,0),
  (9211,'181****1002',NULL,1,0),
  (9212,'小程序用户',NULL,1,0),
  (9213,'微信用户',NULL,1,0),
  (9214,'aaa222',NULL,1,0),
  (9215,'董斌',NULL,1,0),
  (9216,'181****1003',NULL,1,0),
  (9217,'演示用户A',NULL,1,0),
  (9218,'演示用户B',NULL,1,0),
  (9219,'181****1004',NULL,1,0),
  (9220,'小程序用户',NULL,1,0),
  (9221,'微信用户',NULL,1,0),
  (9222,'181****1005',NULL,1,0),
  (9223,'演示用户C',NULL,1,0),
  (9224,'演示用户D',NULL,1,0),
  (9225,'181****1006',NULL,1,0),
  (9226,'小程序用户',NULL,1,0),
  (9227,'微信用户',NULL,1,0),
  (9228,'181****1007',NULL,1,0),
  (9229,'演示用户E',NULL,1,0),
  (9230,'演示用户F',NULL,1,0),
  (9231,'181****1008',NULL,1,0)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `status`=VALUES(`status`);

-- 关联用户券到发送批次（覆盖/补充演示行）
INSERT INTO `qixi_crm_b_coupon_user`
  (`id`,`user_id`,`coupon_id`,`source`,`send_id`,`status`,`obtained_at`,`used_order_id`)
VALUES
  (94051,9101,9401,'platform_manual',45,'unused','2026-07-16 10:58:40',NULL),
  (94052,9102,9401,'platform_manual',45,'unused','2026-07-16 10:58:41',NULL),
  (94053,9103,9401,'platform_manual',45,'used','2026-07-16 11:20:00',9900202),
  (94060,9210,9401,'platform_manual',45,'unused','2026-07-16 10:59:01',NULL),
  (94061,9211,9401,'platform_manual',45,'unused','2026-07-16 10:59:02',NULL),
  (94062,9212,9401,'platform_manual',45,'unused','2026-07-16 10:59:03',NULL),
  (94063,9213,9401,'platform_manual',45,'unused','2026-07-16 10:59:04',NULL),
  (94064,9214,9401,'platform_manual',45,'unused','2026-07-16 10:59:05',NULL),
  (94065,9215,9401,'platform_manual',45,'unused','2026-07-16 10:59:06',NULL),
  (94066,9216,9401,'platform_manual',45,'unused','2026-07-16 10:59:07',NULL),
  (94067,9217,9401,'platform_manual',45,'unused','2026-07-16 10:59:08',NULL),
  (94068,9218,9401,'platform_manual',45,'unused','2026-07-16 10:59:09',NULL),
  (94069,9219,9401,'platform_manual',45,'unused','2026-07-16 10:59:10',NULL),
  (94070,9220,9401,'platform_manual',45,'unused','2026-07-16 10:59:11',NULL),
  (94071,9221,9401,'platform_manual',45,'unused','2026-07-16 10:59:12',NULL),
  (94072,9222,9401,'platform_manual',45,'unused','2026-07-16 10:59:13',NULL),
  (94073,9223,9401,'platform_manual',45,'unused','2026-07-16 10:59:14',NULL),
  (94074,9224,9401,'platform_manual',45,'unused','2026-07-16 10:59:15',NULL),
  (94075,9225,9401,'platform_manual',45,'unused','2026-07-16 10:59:16',NULL),
  (94076,9226,9401,'platform_manual',45,'unused','2026-07-16 10:59:17',NULL),
  (94077,9227,9401,'platform_manual',45,'unused','2026-07-16 10:59:18',NULL),
  (94078,9228,9401,'platform_manual',45,'unused','2026-07-16 10:59:19',NULL),
  (94079,9229,9401,'platform_manual',45,'unused','2026-07-16 10:59:20',NULL),
  (94080,9230,9401,'platform_manual',45,'unused','2026-07-16 10:59:21',NULL),
  (94081,9231,9401,'platform_manual',45,'unused','2026-07-16 10:59:22',NULL),
  (94054,9101,9404,'platform_manual',44,'unused','2026-07-15 14:20:12',NULL),
  (94055,9102,9404,'platform_manual',44,'unused','2026-07-15 14:20:13',NULL),
  (94056,9101,9405,'platform_manual',43,'unused','2026-07-14 09:12:05',NULL),
  (94057,9102,9402,'platform_manual',42,'unused','2026-07-10 16:40:05',NULL),
  (94058,9103,9402,'platform_manual',42,'unused','2026-07-10 16:40:06',NULL)
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `coupon_id`=VALUES(`coupon_id`),
  `source`=VALUES(`source`),
  `send_id`=VALUES(`send_id`),
  `status`=VALUES(`status`),
  `obtained_at`=VALUES(`obtained_at`),
  `used_order_id`=VALUES(`used_order_id`);

-- 头像演示
INSERT INTO `qixi_crm_b_user_profile` (`user_id`,`avatar_url`,`gender`,`bio`,`source_channel`)
VALUES
  (9101,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9102,'https://via.placeholder.com/64',0,'', 'mini_program'),
  (9103,'https://via.placeholder.com/64',0,'', 'h5'),
  (9210,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9211,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9212,'https://via.placeholder.com/64',0,'', 'mini_program'),
  (9213,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9214,'https://via.placeholder.com/64',0,'', 'h5'),
  (9215,'https://via.placeholder.com/64',0,'', 'h5'),
  (9216,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9217,'https://via.placeholder.com/64',0,'', 'pc'),
  (9218,'https://via.placeholder.com/64',0,'', 'pc'),
  (9219,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9220,'https://via.placeholder.com/64',0,'', 'mini_program'),
  (9221,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9222,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9223,'https://via.placeholder.com/64',0,'', 'pc'),
  (9224,'https://via.placeholder.com/64',0,'', 'pc'),
  (9225,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9226,'https://via.placeholder.com/64',0,'', 'mini_program'),
  (9227,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9228,'https://via.placeholder.com/64',0,'', 'wechat'),
  (9229,'https://via.placeholder.com/64',0,'', 'pc'),
  (9230,'https://via.placeholder.com/64',0,'', 'pc'),
  (9231,'https://via.placeholder.com/64',0,'', 'wechat')
ON DUPLICATE KEY UPDATE
  `avatar_url`=IF(`avatar_url`='', VALUES(`avatar_url`), `avatar_url`);
