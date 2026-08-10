-- 积分商品分类（对齐 CRMEB eb_store_category.type=1）
-- 表：qixi_crm_b_points_category；积分商品投影可挂 cate_id 便于 has_product 统计
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_points_category` (
  `store_category_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '积分商品分类 ID',
  `pid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '父分类，积分分类默认 0',
  `cate_name` varchar(100) NOT NULL DEFAULT '' COMMENT '分类名称',
  `path` varchar(255) NOT NULL DEFAULT '/' COMMENT '路径',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序，越大越靠前',
  `pic` varchar(128) NOT NULL DEFAULT '' COMMENT '图标',
  `is_show` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示：1显示 0隐藏',
  `level` int unsigned NOT NULL DEFAULT 0 COMMENT '层级',
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '商户 ID，平台为 0',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_hot` tinyint NOT NULL DEFAULT 0 COMMENT '是否推荐',
  `type` tinyint NOT NULL DEFAULT 1 COMMENT '1=积分商品分类',
  `is_del` tinyint NOT NULL DEFAULT 0 COMMENT '软删：1已删',
  PRIMARY KEY (`store_category_id`),
  KEY `idx_mer_type_del_sort` (`mer_id`,`type`,`is_del`,`sort`),
  KEY `idx_show` (`is_show`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='积分商品分类';

-- 积分商品投影挂分类（幂等）
SET @col_cate := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_points_product_view'
    AND COLUMN_NAME = 'cate_id'
);
SET @sql_cate := IF(
  @col_cate = 0,
  'ALTER TABLE `qixi_crm_b_points_product_view` ADD COLUMN `cate_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''积分商品分类 ID'' AFTER `title`, ADD KEY `idx_cate` (`cate_id`)',
  'SELECT 1'
);
PREPARE stmt_cate FROM @sql_cate;
EXECUTE stmt_cate;
DEALLOCATE PREPARE stmt_cate;

-- 演示分类（中文）
INSERT INTO `qixi_crm_b_points_category`
  (`store_category_id`,`pid`,`cate_name`,`path`,`sort`,`pic`,`is_show`,`level`,`mer_id`,`is_hot`,`type`,`is_del`)
SELECT * FROM (
  SELECT 71 AS store_category_id, 0 AS pid, '数码家电' AS cate_name, '/' AS path, 100 AS sort, '' AS pic, 1 AS is_show, 0 AS level, 0 AS mer_id, 0 AS is_hot, 1 AS type, 0 AS is_del
  UNION ALL SELECT 72, 0, '美妆护肤', '/', 90, '', 1, 0, 0, 0, 1, 0
  UNION ALL SELECT 73, 0, '居家日用', '/', 80, '', 1, 0, 0, 0, 1, 0
  UNION ALL SELECT 74, 0, '服饰箱包', '/', 70, '', 1, 0, 0, 0, 1, 0
  UNION ALL SELECT 75, 0, '食品饮料', '/', 60, '', 0, 0, 0, 0, 1, 0
) AS seed
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_b_points_category` c WHERE c.store_category_id = seed.store_category_id
);

-- 演示商品挂到分类，便于 has_product / 删除提示
UPDATE `qixi_crm_b_points_product_view` SET `cate_id` = 71 WHERE `product_id` = 1005 AND (`cate_id` IS NULL OR `cate_id` = 0);
UPDATE `qixi_crm_b_points_product_view` SET `cate_id` = 72 WHERE `product_id` = 1105 AND (`cate_id` IS NULL OR `cate_id` = 0);
UPDATE `qixi_crm_b_points_product_view` SET `cate_id` = 73 WHERE `product_id` = 1204 AND (`cate_id` IS NULL OR `cate_id` = 0);
