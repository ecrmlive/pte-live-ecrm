-- 区域代理 extend（对齐 CRMEB）：JSON，含 avatar「区域代理」标识图。
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_business_zone_agent' AND COLUMN_NAME='extend')=0,
    'ALTER TABLE `qixi_crm_a_business_zone_agent` ADD COLUMN `extend` text NULL AFTER `remark`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
