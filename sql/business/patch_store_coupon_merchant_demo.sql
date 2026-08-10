-- 平台「商户优惠券」列表演示数据（mer_id > 0）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 使用时间（领取记录弹窗「使用记录」）
SET @col_used_at := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_coupon_user'
    AND COLUMN_NAME = 'used_at'
);
SET @sql_used_at := IF(
  @col_used_at = 0,
  'ALTER TABLE `qixi_crm_b_coupon_user` ADD COLUMN `used_at` datetime DEFAULT NULL COMMENT ''核销/使用时间'' AFTER `obtained_at`',
  'SELECT 1'
);
PREPARE stmt_used_at FROM @sql_used_at;
EXECUTE stmt_used_at;
DEALLOCATE PREPARE stmt_used_at;

INSERT INTO `qixi_crm_b_store_coupon`
  (`coupon_id`,`mer_id`,`is_timeout`,`is_limited`,`total_count`,`remain_count`,`send_type`,`title`,`coupon_price`,`use_min_price`,`coupon_type`,`coupon_time`,`sort`,`status`,`type`,`is_del`)
VALUES
  (9501,1,0,0,0,0,0,'1000无门槛优惠券',1000.00,0,0,365,100,1,0,0),
  (9502,1,0,0,0,0,0,'店铺满减券50',50.00,200,0,30,90,1,0,0),
  (9503,2,0,1,4,4,0,'aaaaa',10.00,0,0,4,80,1,1,0),
  (9504,2,0,0,0,0,0,'11122',20.00,100,0,30,70,0,0,0),
  (9505,1004,0,0,0,0,0,'自营店新人券',15.00,0,0,15,60,1,0,0),
  (9506,1001,0,1,100,88,0,'商品专享券',30.00,99,0,7,50,1,1,0)
ON DUPLICATE KEY UPDATE
  `mer_id`=VALUES(`mer_id`),
  `title`=VALUES(`title`),
  `coupon_price`=VALUES(`coupon_price`),
  `use_min_price`=VALUES(`use_min_price`),
  `coupon_time`=VALUES(`coupon_time`),
  `is_limited`=VALUES(`is_limited`),
  `total_count`=VALUES(`total_count`),
  `remain_count`=VALUES(`remain_count`),
  `status`=VALUES(`status`),
  `type`=VALUES(`type`),
  `sort`=VALUES(`sort`),
  `is_del`=0;

-- 领取/使用演示（含已过期，便于领取记录验收）
INSERT INTO `qixi_crm_b_coupon_user`
  (`id`,`user_id`,`coupon_id`,`source`,`send_id`,`status`,`obtained_at`,`used_at`,`used_order_id`)
VALUES
  (95071,9101,9501,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 2 DAY),NULL,NULL),
  (95072,9102,9501,'receive',0,'used',DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_SUB(NOW(), INTERVAL 20 HOUR),9900301),
  (95073,9103,9502,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 5 HOUR),NULL,NULL),
  (95074,9101,9503,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 3 HOUR),NULL,NULL),
  (95075,9102,9505,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 1 HOUR),NULL,NULL),
  (95076,9103,9504,'receive',0,'expired',DATE_SUB(NOW(), INTERVAL 40 DAY),NULL,NULL),
  (95077,9210,9501,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 6 HOUR),NULL,NULL),
  (95078,9211,9502,'receive',0,'expired',DATE_SUB(NOW(), INTERVAL 35 DAY),NULL,NULL),
  (95079,9212,9506,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 2 HOUR),NULL,NULL),
  (95080,9213,9505,'receive',0,'unused',DATE_SUB(NOW(), INTERVAL 90 MINUTE),NULL,NULL)
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `coupon_id`=VALUES(`coupon_id`),
  `source`=VALUES(`source`),
  `status`=VALUES(`status`),
  `obtained_at`=VALUES(`obtained_at`),
  `used_at`=VALUES(`used_at`),
  `used_order_id`=VALUES(`used_order_id`);

-- 补几条已使用，便于「使用记录」弹窗验收（避开 uk_user_coupon 冲突）
UPDATE `qixi_crm_b_coupon_user`
SET
  `status`='used',
  `used_at`=DATE_SUB(NOW(), INTERVAL 18 HOUR),
  `used_order_id`=COALESCE(`used_order_id`, 9900301)
WHERE `id`=95072;

UPDATE `qixi_crm_b_coupon_user`
SET
  `status`='used',
  `used_at`=DATE_SUB(NOW(), INTERVAL 2 DAY),
  `used_order_id`=COALESCE(`used_order_id`, 9900302)
WHERE `id`=95077;

UPDATE `qixi_crm_b_coupon_user`
SET `used_at` = COALESCE(`used_at`, `obtained_at`)
WHERE `status` = 'used' AND `used_at` IS NULL;

INSERT INTO `qixi_crm_b_coupon_template_view`
  (`coupon_id`,`store_id`,`name`,`discount_type`,`discount_value`,`min_amount`,`status`,`version`)
VALUES
  (9501,1,'1000无门槛优惠券','amount',1000.00,0,1,1),
  (9502,1,'店铺满减券50','amount',50.00,200,1,1),
  (9503,2,'aaaaa','amount',10.00,0,1,1),
  (9504,2,'11122','amount',20.00,100,0,1),
  (9505,1004,'自营店新人券','amount',15.00,0,1,1),
  (9506,1001,'商品专享券','amount',30.00,99,1,1)
ON DUPLICATE KEY UPDATE
  `store_id`=VALUES(`store_id`),
  `name`=VALUES(`name`),
  `discount_value`=VALUES(`discount_value`),
  `min_amount`=VALUES(`min_amount`),
  `status`=VALUES(`status`);
