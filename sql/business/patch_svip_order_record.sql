-- 会员记录：订单补支付方式 / 到期时间（对齐 CRMEB user_order.pay_type / end_time），幂等
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_order' AND COLUMN_NAME='pay_type')=0,
    'ALTER TABLE `qixi_crm_b_svip_order` ADD COLUMN `pay_type` varchar(32) NOT NULL DEFAULT \'\' AFTER `amount`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_order' AND COLUMN_NAME='end_time')=0,
    'ALTER TABLE `qixi_crm_b_svip_order` ADD COLUMN `end_time` datetime DEFAULT NULL AFTER `paid_at`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
