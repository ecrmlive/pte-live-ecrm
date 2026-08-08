-- 商品：店铺分类 + 配送方式（平台后台编辑可改）
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product' AND COLUMN_NAME='store_category_id')=0,
    'ALTER TABLE `qixi_crm_m_product` ADD COLUMN `store_category_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `category_id`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_detail' AND COLUMN_NAME='delivery_way')=0,
    'ALTER TABLE `qixi_crm_m_product_detail` ADD COLUMN `delivery_way` varchar(64) NOT NULL DEFAULT ''2'' AFTER `cover_url`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
