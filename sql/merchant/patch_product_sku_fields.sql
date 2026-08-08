-- 商品 SKU：对齐 CRMEB attr_value 常用字段（图片/划线价/编码/条码/重量体积/一级返佣）
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='image')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `image` varchar(1024) NOT NULL DEFAULT '''' AFTER `spec_json`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='ot_price')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `ot_price` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `price`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='code')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `code` varchar(64) NOT NULL DEFAULT '''' AFTER `stock`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='bar_code')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `bar_code` varchar(64) NOT NULL DEFAULT '''' AFTER `code`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='weight')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `weight` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `bar_code`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='volume')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `volume` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `weight`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product_sku' AND COLUMN_NAME='extension_one')=0,
    'ALTER TABLE `qixi_crm_m_product_sku` ADD COLUMN `extension_one` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `volume`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
