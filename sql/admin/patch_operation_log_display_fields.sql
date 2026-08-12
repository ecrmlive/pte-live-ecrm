-- 操作日志展示字段：保留原始审计字段，并增加请求方式、链接、IP 与权限名称。
-- MySQL 8.4 兼容：不依赖 ADD COLUMN IF NOT EXISTS，补丁可以安全重复执行。
SET @request_method_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_operation_log' AND COLUMN_NAME='request_method')=0,
  'ALTER TABLE `qixi_crm_a_operation_log` ADD COLUMN `request_method` varchar(16) NOT NULL DEFAULT '''' AFTER `request_id`',
  'SELECT 1');
PREPARE request_method_stmt FROM @request_method_sql;
EXECUTE request_method_stmt;
DEALLOCATE PREPARE request_method_stmt;

SET @request_path_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_operation_log' AND COLUMN_NAME='request_path')=0,
  'ALTER TABLE `qixi_crm_a_operation_log` ADD COLUMN `request_path` varchar(512) NOT NULL DEFAULT '''' AFTER `request_method`',
  'SELECT 1');
PREPARE request_path_stmt FROM @request_path_sql;
EXECUTE request_path_stmt;
DEALLOCATE PREPARE request_path_stmt;

SET @request_ip_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_operation_log' AND COLUMN_NAME='request_ip')=0,
  'ALTER TABLE `qixi_crm_a_operation_log` ADD COLUMN `request_ip` varchar(64) NOT NULL DEFAULT '''' AFTER `request_path`',
  'SELECT 1');
PREPARE request_ip_stmt FROM @request_ip_sql;
EXECUTE request_ip_stmt;
DEALLOCATE PREPARE request_ip_stmt;

SET @permission_name_sql := IF(
  (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_operation_log' AND COLUMN_NAME='permission_name')=0,
  'ALTER TABLE `qixi_crm_a_operation_log` ADD COLUMN `permission_name` varchar(128) NOT NULL DEFAULT '''' AFTER `request_ip`',
  'SELECT 1');
PREPARE permission_name_stmt FROM @permission_name_sql;
EXECUTE permission_name_stmt;
DEALLOCATE PREPARE permission_name_stmt;

SET @idx_method_time_sql := IF(
  (SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_operation_log' AND INDEX_NAME='idx_method_time')=0,
  'CREATE INDEX `idx_method_time` ON `qixi_crm_a_operation_log` (`request_method`,`created_at`)',
  'SELECT 1');
PREPARE idx_method_time_stmt FROM @idx_method_time_sql;
EXECUTE idx_method_time_stmt;
DEALLOCATE PREPARE idx_method_time_stmt;

-- 历史审计记录没有独立展示字段时，从原 action/resource_type 回填可公开展示的信息。
UPDATE `qixi_crm_a_operation_log`
SET
  `request_method` = CASE
    WHEN `request_method`='' AND SUBSTRING_INDEX(`action`,' ',1) IN ('GET','POST','PUT','PATCH','DELETE')
      THEN SUBSTRING_INDEX(`action`,' ',1)
    ELSE `request_method`
  END,
  `request_path` = CASE
    WHEN `request_path`='' AND SUBSTRING_INDEX(`action`,' ',1) IN ('GET','POST','PUT','PATCH','DELETE')
      THEN TRIM(SUBSTRING(`action`, LOCATE(' ',`action`)+1))
    ELSE `request_path`
  END,
  `permission_name` = CASE
    WHEN `permission_name`<>'' THEN `permission_name`
    WHEN `resource_type` IN ('setting','maintain') THEN '系统设置'
    WHEN `resource_type`='customer-service' THEN '客服管理'
    WHEN `resource_type` IN ('product','products') THEN '商品管理'
    WHEN `resource_type`='merchants' THEN '商户管理'
    WHEN `resource_type`='stores' THEN '店铺管理'
    WHEN `resource_type` IN ('users','user-list') THEN '用户管理'
    ELSE '平台管理'
  END;
