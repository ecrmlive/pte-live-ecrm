-- 积分商品投影字段补齐（对齐 CRMEB 积分商品列表：创建时间/排序/已兑换/软删/来源商品）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_points_product_view` (
  `product_id` bigint unsigned NOT NULL,
  `merchant_id` bigint unsigned NOT NULL,
  `store_id` bigint unsigned NOT NULL,
  `merchant_name` varchar(128) NOT NULL DEFAULT '',
  `store_name` varchar(128) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL,
  `cate_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '积分商品分类 ID',
  `cover_url` varchar(1024) NOT NULL DEFAULT '',
  `original_price` decimal(12,2) NOT NULL DEFAULT 0 COMMENT '兑换金额',
  `points_required` bigint NOT NULL COMMENT '兑换积分',
  `stock` int NOT NULL,
  `sales` int NOT NULL DEFAULT 0 COMMENT '已兑换数量',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序，越大越靠前',
  `sale_status` tinyint NOT NULL DEFAULT 1 COMMENT '上架状态：1上架 0下架',
  `source_product_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '快速添加来源普通商品 ID',
  `is_del` tinyint NOT NULL DEFAULT 0 COMMENT '软删：1已删',
  `version` bigint unsigned NOT NULL DEFAULT 1,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`),
  KEY `idx_store_sale` (`store_id`,`sale_status`),
  KEY `idx_cate` (`cate_id`),
  KEY `idx_del_sort` (`is_del`,`sort`,`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商品消费投影';

-- 幂等补列
SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND COLUMN_NAME = 'sales');
SET @sql := IF(@col = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `sales` int NOT NULL DEFAULT 0 COMMENT ''已兑换数量'' AFTER `stock`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND COLUMN_NAME = 'sort');
SET @sql := IF(@col = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `sort` int NOT NULL DEFAULT 0 COMMENT ''排序，越大越靠前'' AFTER `sales`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND COLUMN_NAME = 'source_product_id');
SET @sql := IF(@col = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `source_product_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''快速添加来源普通商品 ID'' AFTER `sale_status`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND COLUMN_NAME = 'is_del');
SET @sql := IF(@col = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `is_del` tinyint NOT NULL DEFAULT 0 COMMENT ''软删：1已删'' AFTER `source_product_id`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @col := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND COLUMN_NAME = 'create_time');
SET @sql := IF(@col = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'' AFTER `version`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @idx := (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_crm_b_points_product_view' AND INDEX_NAME = 'idx_del_sort');
SET @sql := IF(@idx = 0, 'ALTER TABLE `qixi_crm_b_points_product_view` ADD KEY `idx_del_sort` (`is_del`,`sort`,`product_id`)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 旧行 create_time 异常则回填 updated_at（避免严格模式零日期比较报错）
UPDATE `qixi_crm_b_points_product_view`
SET `create_time` = COALESCE(NULLIF(`updated_at`, '1970-01-01 00:00:00'), NOW())
WHERE `create_time` IS NULL
   OR YEAR(`create_time`) < 1971;

-- 演示积分商品（中文可读，幂等）
INSERT INTO `qixi_crm_b_points_product_view`
  (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`title`,`cate_id`,`cover_url`,`original_price`,`points_required`,`stock`,`sales`,`sort`,`sale_status`,`source_product_id`,`is_del`,`version`,`create_time`)
VALUES
  (1005,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店','真丝印花方巾礼盒',71,'/demo/product-scarf-v1.png',16.90,120,20,8,100,1,0,0,1,DATE_SUB(NOW(), INTERVAL 9 DAY)),
  (1105,2,2,'CRM Live居家商户','CRM Live居家优选店','真丝睡眠眼罩方巾组',72,'/demo/product-scarf-v1.png',13.90,180,16,5,90,1,0,0,1,DATE_SUB(NOW(), INTERVAL 7 DAY)),
  (1204,3,3,'CRM Live数码商户','CRM Live数码生活店','便携保温杯清洁套装',73,'/demo/product-tumbler-v1.png',17.90,220,12,3,80,1,0,0,1,DATE_SUB(NOW(), INTERVAL 5 DAY)),
  (1401,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店','轻奢羊绒围巾兑换券',71,'/demo/product-knit-v1.png',39.90,500,30,12,70,1,1001,0,1,DATE_SUB(NOW(), INTERVAL 3 DAY)),
  (1402,2,2,'CRM Live居家商户','CRM Live居家优选店','居家香氛体验兑换',72,'/demo/product-fragrance-v1.png',29.90,360,18,0,60,0,1302,0,1,DATE_SUB(NOW(), INTERVAL 2 DAY)),
  (1403,3,3,'CRM Live数码商户','CRM Live数码生活店','跑步鞋护理礼包',73,'/demo/product-shoes-v1.png',49.90,680,10,2,50,1,0,0,1,DATE_SUB(NOW(), INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE
  `merchant_id`=VALUES(`merchant_id`),
  `store_id`=VALUES(`store_id`),
  `merchant_name`=VALUES(`merchant_name`),
  `store_name`=VALUES(`store_name`),
  `title`=VALUES(`title`),
  `cate_id`=VALUES(`cate_id`),
  `cover_url`=VALUES(`cover_url`),
  `original_price`=VALUES(`original_price`),
  `points_required`=VALUES(`points_required`),
  `stock`=VALUES(`stock`),
  `sales`=VALUES(`sales`),
  `sort`=VALUES(`sort`),
  `sale_status`=VALUES(`sale_status`),
  `source_product_id`=VALUES(`source_product_id`),
  `is_del`=0,
  `updated_at`=NOW();
