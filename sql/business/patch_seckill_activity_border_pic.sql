-- 秒杀活动：活动边框图（对齐 CRMEB store_seckill_active.border_pic）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @col_border := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'qixi_crm_b_seckill_activity'
    AND COLUMN_NAME = 'border_pic'
);
SET @sql_border := IF(
  @col_border = 0,
  'ALTER TABLE `qixi_crm_b_seckill_activity` ADD COLUMN `border_pic` varchar(1024) NOT NULL DEFAULT '''' COMMENT ''活动边框图'' AFTER `product_category_ids`',
  'SELECT 1'
);
PREPARE stmt_border FROM @sql_border;
EXECUTE stmt_border;
DEALLOCATE PREPARE stmt_border;

-- 演示：活动 5 挂一张边框示例（沿用商品封面，便于查看态可见）
UPDATE `qixi_crm_b_seckill_activity` a
LEFT JOIN `qixi_crm_b_product_view` p ON p.product_id = 1001
SET a.`border_pic` = IFNULL(NULLIF(TRIM(p.cover_url), ''), a.`border_pic`)
WHERE a.`seckill_activity_id` = 5
  AND a.`delete_time` IS NULL
  AND (a.`border_pic` IS NULL OR a.`border_pic` = '');

-- 演示商品范围（平台分类，不选则为全品类；此处给活动 5 挂上分类便于查看）
UPDATE `qixi_crm_b_seckill_activity`
SET `product_category_ids` = '1,2'
WHERE `seckill_activity_id` = 5
  AND `delete_time` IS NULL
  AND (`product_category_ids` IS NULL OR `product_category_ids` = '');
