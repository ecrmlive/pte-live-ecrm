-- 用户资料补充：真实姓名、身份证（平台后台「创建用户」表单对齐 CRMEB）
-- 幂等：列已存在则跳过
USE `qixi_crm_business`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_user_profile' AND COLUMN_NAME='real_name')=0,
    'ALTER TABLE `qixi_crm_b_user_profile` ADD COLUMN `real_name` varchar(64) NOT NULL DEFAULT \'\' AFTER `avatar_url`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_user_profile' AND COLUMN_NAME='id_card')=0,
    'ALTER TABLE `qixi_crm_b_user_profile` ADD COLUMN `id_card` varchar(32) NOT NULL DEFAULT \'\' AFTER `real_name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
