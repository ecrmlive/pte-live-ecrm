-- 平台/店铺优惠券模板（对齐 CRMEB store_coupon 字段；替换已废弃的 qixi_m_admin_store_coupon）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_coupon` (
  `coupon_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '0=平台券',
  `is_timeout` tinyint NOT NULL DEFAULT 0,
  `start_time` datetime DEFAULT NULL,
  `end_time` datetime DEFAULT NULL,
  `is_limited` tinyint NOT NULL DEFAULT 0,
  `total_count` int unsigned NOT NULL DEFAULT 0,
  `remain_count` int unsigned NOT NULL DEFAULT 0,
  `send_type` tinyint NOT NULL DEFAULT 0,
  `full_reduction` decimal(12,2) NOT NULL DEFAULT 0.00,
  `title` varchar(64) NOT NULL DEFAULT '',
  `coupon_price` decimal(12,2) NOT NULL DEFAULT 0.00,
  `use_min_price` int NOT NULL DEFAULT 0,
  `coupon_type` tinyint NOT NULL DEFAULT 0 COMMENT '0领取后N天 1固定时段',
  `coupon_time` int unsigned NOT NULL DEFAULT 0,
  `use_start_time` datetime DEFAULT NULL,
  `use_end_time` datetime DEFAULT NULL,
  `sort` int unsigned NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `type` int NOT NULL DEFAULT 0 COMMENT '0店铺券 10平台券',
  PRIMARY KEY (`coupon_id`),
  KEY `idx_mer_type_status` (`mer_id`,`type`,`status`,`is_del`),
  KEY `idx_sort` (`sort`,`coupon_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- 本地演示：平台券（type=10, mer_id=0）
INSERT INTO `qixi_crm_b_store_coupon`
  (`coupon_id`,`mer_id`,`is_limited`,`total_count`,`remain_count`,`title`,`coupon_price`,`use_min_price`,`coupon_type`,`coupon_time`,`sort`,`status`,`type`,`is_del`)
VALUES
  (9401,0,0,0,0,'平台新客满99减10',10.00,99,0,30,20,1,10,0),
  (9402,0,1,100,88,'平台夏日满299减40',40.00,299,0,60,10,1,10,0),
  (9403,0,0,0,0,'测试隐藏券',5.00,50,0,7,1,0,10,0)
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `coupon_price`=VALUES(`coupon_price`),
  `use_min_price`=VALUES(`use_min_price`),
  `coupon_time`=VALUES(`coupon_time`),
  `is_limited`=VALUES(`is_limited`),
  `total_count`=VALUES(`total_count`),
  `remain_count`=VALUES(`remain_count`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `type`=VALUES(`type`),
  `is_del`=0;
