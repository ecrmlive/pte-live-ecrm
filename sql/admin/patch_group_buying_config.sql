-- 平台「拼团设置」夹具（对齐 CRMEB ConfigOthers::getGroupBuying / setGroupBuying）
-- 存储：qixi_crm_a_setting_cache.key = group_buying_config
-- 字段：ficti_status（虚拟成团启用 0/1）、group_buying_rate（真实成团最小比例 0～100）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('group_buying_config',0,'{"ficti_status":1,"group_buying_rate":30}')
ON DUPLICATE KEY UPDATE `expire_time`=VALUES(`expire_time`),`result`=VALUES(`result`);
