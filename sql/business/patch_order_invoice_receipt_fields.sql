-- 平台发票列表对齐 CRMEB store_order_receipt 字段（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_order_invoice' AND COLUMN_NAME='invoice_type')=0,
    'ALTER TABLE `qixi_crm_b_order_invoice`
      ADD COLUMN `invoice_type` tinyint NOT NULL DEFAULT 1 COMMENT ''1普通发票 2专用发票'' AFTER `profile_type`,
      ADD COLUMN `receipt_sn` varchar(64) NOT NULL DEFAULT '''' COMMENT ''发票申请单号'' AFTER `invoice_type`,
      ADD COLUMN `invoice_amount` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT ''发票金额'' AFTER `receipt_sn`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
     WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_order_invoice' AND COLUMN_NAME='mark')=0,
    'ALTER TABLE `qixi_crm_b_order_invoice`
      ADD COLUMN `mark` varchar(500) NOT NULL DEFAULT '''' COMMENT ''发票备注'' AFTER `rejection_reason`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 兼容：旧行把 invoice_no 回填到 receipt_sn
UPDATE `qixi_crm_b_order_invoice`
SET `receipt_sn` = `invoice_no`
WHERE `receipt_sn` = '' AND `invoice_no` <> '';
