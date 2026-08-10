-- 平台「秒杀活动」活动场（对齐 CRMEB store_seckill_active；与秒杀商品 seckill_active 分离）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_seckill_activity` (
  `seckill_activity_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL DEFAULT '',
  `seckill_time_ids` varchar(255) NOT NULL DEFAULT '',
  `start_day` date NOT NULL,
  `end_day` date NOT NULL,
  `once_pay_count` int NOT NULL DEFAULT 1,
  `all_pay_count` int NOT NULL DEFAULT 0,
  `product_category_ids` varchar(255) NOT NULL DEFAULT '',
  `border_pic` varchar(1024) NOT NULL DEFAULT '' COMMENT '活动边框图',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '活动开关 1开 0关',
  `active_status` tinyint NOT NULL DEFAULT 0 COMMENT '0未开始 1进行中 -1已结束',
  `product_count` int NOT NULL DEFAULT 0,
  `merchant_count` int NOT NULL DEFAULT 0,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `delete_time` datetime DEFAULT NULL,
  PRIMARY KEY (`seckill_activity_id`),
  KEY `idx_status_day` (`status`,`active_status`,`start_day`,`end_day`),
  KEY `idx_delete` (`delete_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 秒杀商品关联活动
SET @col_aid := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_seckill_active'
    AND COLUMN_NAME = 'activity_id'
);
SET @sql_aid := IF(
  @col_aid = 0,
  'ALTER TABLE `qixi_crm_b_seckill_active` ADD COLUMN `activity_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''所属秒杀活动'' AFTER `seckill_active_id`, ADD KEY `idx_activity` (`activity_id`)',
  'SELECT 1'
);
PREPARE stmt_aid FROM @sql_aid;
EXECUTE stmt_aid;
DEALLOCATE PREPARE stmt_aid;

-- 对齐 CRMEB「秒杀活动」列表截图的活动场（开关默认关闭；日期使今天能区分进行中/已结束）
INSERT INTO `qixi_crm_b_seckill_activity`
  (`seckill_activity_id`,`name`,`seckill_time_ids`,`start_day`,`end_day`,`once_pay_count`,`all_pay_count`,
   `status`,`active_status`,`product_count`,`merchant_count`,`create_time`,`delete_time`)
VALUES
  (5,'秒杀活动','1,2,3,4','2026-07-15','2027-01-31',1,0,0,1,18,2,'2025-12-12 22:45:39',NULL),
  (4,'晚上','4','2026-07-15','2026-08-31',1,0,0,1,7,2,'2025-07-23 14:39:07',NULL),
  (3,'下午','3','2026-07-15','2026-08-31',1,0,0,1,2,1,'2025-07-23 14:38:46',NULL),
  (2,'早上','2','2025-07-15','2026-01-31',1,0,0,-1,1,1,'2025-07-23 14:38:29',NULL),
  (1,'全天','1,2','2025-07-15','2025-10-25',1,0,0,-1,3,2,'2025-07-23 14:38:03',NULL)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `seckill_time_ids`=VALUES(`seckill_time_ids`),
  `start_day`=VALUES(`start_day`),
  `end_day`=VALUES(`end_day`),
  `once_pay_count`=VALUES(`once_pay_count`),
  `all_pay_count`=VALUES(`all_pay_count`),
  `status`=VALUES(`status`),
  `active_status`=VALUES(`active_status`),
  `product_count`=VALUES(`product_count`),
  `merchant_count`=VALUES(`merchant_count`),
  `create_time`=VALUES(`create_time`),
  `delete_time`=NULL;

-- 刷新活动态（以本地 CURDATE 为准）
UPDATE `qixi_crm_b_seckill_activity`
SET `active_status` = CASE
  WHEN CURDATE() < `start_day` THEN 0
  WHEN CURDATE() > `end_day` THEN -1
  ELSE 1
END
WHERE `delete_time` IS NULL
  AND `seckill_activity_id` IN (1,2,3,4,5);

-- 确保秒杀管理扩展列已存在（本文件字母序可能早于 patch_seckill_manage.sql）
SET @col_is_show := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_seckill_active' AND COLUMN_NAME = 'is_show'
);
SET @sql_is_show := IF(
  @col_is_show = 0,
  'ALTER TABLE `qixi_crm_b_seckill_active`
     ADD COLUMN `is_show` tinyint NOT NULL DEFAULT 1 COMMENT ''是否显示'' AFTER `status`,
     ADD COLUMN `product_status` tinyint NOT NULL DEFAULT 1 COMMENT ''1通过 0待审 -1未通过 -2强制下架'' AFTER `is_show`,
     ADD COLUMN `star` tinyint NOT NULL DEFAULT 0 COMMENT ''推荐星级'' AFTER `product_status`,
     ADD COLUMN `sort` int NOT NULL DEFAULT 0 COMMENT ''排序'' AFTER `star`,
     ADD COLUMN `stock` int NOT NULL DEFAULT 0 COMMENT ''限量剩余'' AFTER `sort`,
     ADD COLUMN `sales` int NOT NULL DEFAULT 0 COMMENT ''已售数量'' AFTER `stock`,
     ADD COLUMN `sys_labels` varchar(255) NOT NULL DEFAULT '''' COMMENT ''平台标签ID逗号分隔'' AFTER `sales`,
     ADD COLUMN `refusal` varchar(500) NOT NULL DEFAULT '''' COMMENT ''拒绝/下架原因'' AFTER `sys_labels`',
  'SELECT 1'
);
PREPARE stmt_is_show FROM @sql_is_show;
EXECUTE stmt_is_show;
DEALLOCATE PREPARE stmt_is_show;

-- 用户操作向演示秒杀商品（销量/库存/店铺参与），挂到对应活动
-- product_id 仅用 1001/1002（两端商品视图均有）；mer_id 用 1/2 拉出店铺数
INSERT INTO `qixi_crm_b_seckill_active`
  (`seckill_active_id`,`activity_id`,`name`,`seckill_time_ids`,`start_day`,`end_day`,`mer_id`,`product_id`,`seckill_price`,
   `once_pay_count`,`all_pay_count`,`active_status`,`status`,`is_show`,`product_status`,`star`,`sort`,`stock`,`sales`,
   `sys_labels`,`refusal`,`create_time`,`update_time`,`delete_time`)
VALUES
  -- 活动5「秒杀活动」：18 件 / 2 店
  (6201,5,'夏凉被秒杀','1,2,3,4','2026-07-15','2027-01-31',1,1001,69.90,1,0,1,1,1,1,4,100,120,36,'','',UNIX_TIMESTAMP('2025-12-13 10:20:00'),UNIX_TIMESTAMP('2026-08-01 09:00:00'),NULL),
  (6202,5,'纯棉四件套','1,2,3,4','2026-07-15','2027-01-31',1,1002,129.00,1,0,1,1,1,1,3,90,88,41,'','',UNIX_TIMESTAMP('2025-12-13 10:21:00'),UNIX_TIMESTAMP('2026-08-01 09:00:00'),NULL),
  (6203,5,'运动水杯','1,2,3,4','2026-07-15','2027-01-31',2,1001,19.90,2,5,1,1,1,1,2,80,200,95,'','',UNIX_TIMESTAMP('2025-12-14 11:05:00'),UNIX_TIMESTAMP('2026-08-02 11:00:00'),NULL),
  (6204,5,'蓝牙耳机','1,2,3,4','2026-07-15','2027-01-31',2,1002,99.00,1,0,1,1,1,1,5,70,60,52,'','',UNIX_TIMESTAMP('2025-12-14 11:06:00'),UNIX_TIMESTAMP('2026-08-02 11:00:00'),NULL),
  (6205,5,'懒人沙发','1,2,3,4','2026-07-15','2027-01-31',1,1001,159.00,1,0,1,1,1,1,2,60,25,18,'','',UNIX_TIMESTAMP('2025-12-15 15:30:00'),UNIX_TIMESTAMP('2026-08-03 12:00:00'),NULL),
  (6206,5,'厨房置物架','1,2,3,4','2026-07-15','2027-01-31',1,1002,49.90,1,3,1,1,1,1,1,50,70,33,'','',UNIX_TIMESTAMP('2025-12-15 15:31:00'),UNIX_TIMESTAMP('2026-08-03 12:00:00'),NULL),
  (6207,5,'儿童绘本套装','1,2,3,4','2026-07-15','2027-01-31',2,1001,39.90,1,0,1,1,1,1,3,40,45,27,'','',UNIX_TIMESTAMP('2025-12-16 09:12:00'),UNIX_TIMESTAMP('2026-08-04 10:00:00'),NULL),
  (6208,5,'护肤礼盒','1,2,3,4','2026-07-15','2027-01-31',2,1002,89.00,1,0,1,1,1,1,4,30,35,22,'','',UNIX_TIMESTAMP('2025-12-16 09:13:00'),UNIX_TIMESTAMP('2026-08-04 10:00:00'),NULL),
  (6209,5,'露营折叠椅','1,2,3,4','2026-07-15','2027-01-31',1,1001,79.00,1,0,1,1,1,1,2,20,40,19,'','',UNIX_TIMESTAMP('2025-12-17 16:40:00'),UNIX_TIMESTAMP('2026-08-05 14:00:00'),NULL),
  (6210,5,'便携榨汁杯','1,2,3,4','2026-07-15','2027-01-31',1,1002,59.90,1,0,1,1,1,1,1,10,55,31,'','',UNIX_TIMESTAMP('2025-12-17 16:41:00'),UNIX_TIMESTAMP('2026-08-05 14:00:00'),NULL),
  (6211,5,'速干运动套装','1,2,3,4','2026-07-15','2027-01-31',2,1001,109.00,1,0,1,1,1,1,3,8,28,14,'','',UNIX_TIMESTAMP('2026-01-08 13:20:00'),UNIX_TIMESTAMP('2026-08-06 15:00:00'),NULL),
  (6212,5,'香薰机','1,2,3,4','2026-07-15','2027-01-31',2,1002,69.00,1,0,1,1,1,1,2,7,32,16,'','',UNIX_TIMESTAMP('2026-01-08 13:21:00'),UNIX_TIMESTAMP('2026-08-06 15:00:00'),NULL),
  (6213,5,'数据线三件套','1,2,3,4','2026-07-15','2027-01-31',1,1001,15.90,3,10,1,1,1,1,1,6,180,120,'','',UNIX_TIMESTAMP('2026-02-11 10:05:00'),UNIX_TIMESTAMP('2026-08-07 09:30:00'),NULL),
  (6214,5,'保温杯','1,2,3,4','2026-07-15','2027-01-31',1,1002,45.00,1,0,1,1,1,1,2,5,90,48,'','',UNIX_TIMESTAMP('2026-02-11 10:06:00'),UNIX_TIMESTAMP('2026-08-07 09:30:00'),NULL),
  (6215,5,'瑜伽垫','1,2,3,4','2026-07-15','2027-01-31',2,1001,35.00,1,0,1,1,1,1,1,4,66,29,'','',UNIX_TIMESTAMP('2026-03-03 18:18:00'),UNIX_TIMESTAMP('2026-08-08 11:00:00'),NULL),
  (6216,5,'办公桌收纳','1,2,3,4','2026-07-15','2027-01-31',2,1002,29.90,1,0,1,1,1,1,1,3,74,37,'','',UNIX_TIMESTAMP('2026-03-03 18:19:00'),UNIX_TIMESTAMP('2026-08-08 11:00:00'),NULL),
  (6217,5,'猫砂盆','1,2,3,4','2026-07-15','2027-01-31',1,1001,55.00,1,0,1,1,1,1,0,2,22,9,'','',UNIX_TIMESTAMP('2026-04-20 12:00:00'),UNIX_TIMESTAMP('2026-08-08 16:00:00'),NULL),
  (6218,5,'护眼台灯','1,2,3,4','2026-07-15','2027-01-31',2,1002,119.00,1,0,1,1,1,1,4,1,18,11,'','',UNIX_TIMESTAMP('2026-04-20 12:01:00'),UNIX_TIMESTAMP('2026-08-08 16:00:00'),NULL),

  -- 活动4「晚上」：7 件 / 2 店
  (6219,4,'夜宵螺蛳粉','4','2026-07-15','2026-08-31',1,1001,12.90,2,6,1,1,1,1,2,40,150,88,'','',UNIX_TIMESTAMP('2025-07-24 20:10:00'),UNIX_TIMESTAMP('2026-08-01 20:00:00'),NULL),
  (6220,4,'夜用眼罩','4','2026-07-15','2026-08-31',1,1002,9.90,1,0,1,1,1,1,1,30,90,45,'','',UNIX_TIMESTAMP('2025-07-24 20:11:00'),UNIX_TIMESTAMP('2026-08-01 20:00:00'),NULL),
  (6221,4,'夜宵烧烤套餐','4','2026-07-15','2026-08-31',2,1001,39.90,1,0,1,1,1,1,3,20,40,26,'','',UNIX_TIMESTAMP('2025-07-25 19:30:00'),UNIX_TIMESTAMP('2026-08-02 21:00:00'),NULL),
  (6222,4,'安眠香薰','4','2026-07-15','2026-08-31',2,1002,49.00,1,0,1,1,1,1,2,15,28,13,'','',UNIX_TIMESTAMP('2025-07-25 19:31:00'),UNIX_TIMESTAMP('2026-08-02 21:00:00'),NULL),
  (6223,4,'夜跑腰包','4','2026-07-15','2026-08-31',1,1001,25.00,1,0,1,1,1,1,1,10,35,17,'','',UNIX_TIMESTAMP('2025-07-26 21:05:00'),UNIX_TIMESTAMP('2026-08-03 22:00:00'),NULL),
  (6224,4,'夜宵奶茶券','4','2026-07-15','2026-08-31',2,1002,6.90,5,20,1,1,1,1,0,8,300,210,'','',UNIX_TIMESTAMP('2025-07-26 21:06:00'),UNIX_TIMESTAMP('2026-08-03 22:00:00'),NULL),
  (6225,4,'夜读台灯','4','2026-07-15','2026-08-31',1,1001,89.00,1,0,1,1,1,1,4,5,16,8,'','',UNIX_TIMESTAMP('2025-07-27 22:15:00'),UNIX_TIMESTAMP('2026-08-04 23:00:00'),NULL),

  -- 活动3「下午」：2 件 / 1 店
  (6226,3,'下午茶糕点','3','2026-07-15','2026-08-31',1,1001,18.80,1,0,1,1,1,1,2,25,60,34,'','',UNIX_TIMESTAMP('2025-07-24 14:20:00'),UNIX_TIMESTAMP('2026-08-05 15:00:00'),NULL),
  (6227,3,'冰美式咖啡','3','2026-07-15','2026-08-31',1,1002,9.90,2,8,1,1,1,1,1,20,120,76,'','',UNIX_TIMESTAMP('2025-07-24 14:21:00'),UNIX_TIMESTAMP('2026-08-05 15:00:00'),NULL),

  -- 活动2「早上」（已结束）：1 件 / 1 店
  (6228,2,'早餐豆浆油条','2','2025-07-15','2026-01-31',1,1001,8.80,1,0,-1,0,0,1,0,1,0,64,'','活动已结束',UNIX_TIMESTAMP('2025-07-24 07:30:00'),UNIX_TIMESTAMP('2026-02-01 08:00:00'),NULL),

  -- 活动1「全天」（已结束）：3 件 / 2 店
  (6229,1,'全天水果盒','1,2','2025-07-15','2025-10-25',1,1001,22.00,1,0,-1,0,0,1,1,3,0,41,'','活动已结束',UNIX_TIMESTAMP('2025-07-24 09:00:00'),UNIX_TIMESTAMP('2025-10-26 10:00:00'),NULL),
  (6230,1,'全天坚果礼包','1,2','2025-07-15','2025-10-25',2,1002,45.00,1,0,-1,0,0,1,2,2,0,29,'','活动已结束',UNIX_TIMESTAMP('2025-07-24 09:01:00'),UNIX_TIMESTAMP('2025-10-26 10:00:00'),NULL),
  (6231,1,'全天面包券','1,2','2025-07-15','2025-10-25',1,1001,5.90,3,15,-1,0,0,1,0,1,0,102,'','活动已结束',UNIX_TIMESTAMP('2025-07-24 09:02:00'),UNIX_TIMESTAMP('2025-10-26 10:00:00'),NULL)
ON DUPLICATE KEY UPDATE
  `activity_id`=VALUES(`activity_id`),
  `name`=VALUES(`name`),
  `seckill_time_ids`=VALUES(`seckill_time_ids`),
  `start_day`=VALUES(`start_day`),
  `end_day`=VALUES(`end_day`),
  `mer_id`=VALUES(`mer_id`),
  `product_id`=VALUES(`product_id`),
  `seckill_price`=VALUES(`seckill_price`),
  `once_pay_count`=VALUES(`once_pay_count`),
  `all_pay_count`=VALUES(`all_pay_count`),
  `active_status`=VALUES(`active_status`),
  `status`=VALUES(`status`),
  `is_show`=VALUES(`is_show`),
  `product_status`=VALUES(`product_status`),
  `star`=VALUES(`star`),
  `sort`=VALUES(`sort`),
  `stock`=VALUES(`stock`),
  `sales`=VALUES(`sales`),
  `refusal`=VALUES(`refusal`),
  `delete_time`=VALUES(`delete_time`),
  `update_time`=VALUES(`update_time`);

-- 秒杀管理旧夹具(6001~6007)仍服务「秒杀管理」Tab；不计入活动场商品数，避免冲掉截图数量
UPDATE `qixi_crm_b_seckill_active`
SET `activity_id`=0
WHERE `seckill_active_id` IN (6001,6002,6003,6004,6005,6006,6007);

-- 按真实参与商品回填数量（与截图一致：18/7/2/1/3）
UPDATE `qixi_crm_b_seckill_activity` a
SET
  `product_count` = (
    SELECT COUNT(*) FROM `qixi_crm_b_seckill_active` p
    WHERE p.`activity_id` = a.`seckill_activity_id` AND p.`delete_time` IS NULL
  ),
  `merchant_count` = (
    SELECT COUNT(DISTINCT p.`mer_id`) FROM `qixi_crm_b_seckill_active` p
    WHERE p.`activity_id` = a.`seckill_activity_id` AND p.`delete_time` IS NULL
  )
WHERE a.`seckill_activity_id` IN (1,2,3,4,5) AND a.`delete_time` IS NULL;
