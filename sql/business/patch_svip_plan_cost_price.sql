-- 付费会员类型：补原价字段（对齐 CRMEB cost_price），幂等
USE `qixi_crm_business`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_plan' AND COLUMN_NAME='cost_price')=0,
    'ALTER TABLE `qixi_crm_b_svip_plan` ADD COLUMN `cost_price` decimal(12,2) NOT NULL DEFAULT 0.00 AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 历史行：若原价仍为 0，用售价回填，避免列表空白
UPDATE `qixi_crm_b_svip_plan`
SET `cost_price` = `price`
WHERE `cost_price` = 0 AND `price` > 0;
