-- 平台「分销配置」夹具（对齐 CRMEB ConfigOthers / distribution_tabs）
-- 存储：qixi_crm_a_setting_cache.key = distribution_config
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('distribution_config',0,'{"extension_status":true,"extension_self":true,"extension_limit":false,"extension_limit_day":15,"promoter_type":0,"promoter_low_money":0,"extension_pop":0,"extension_one_rate":0.15,"extension_two_rate":0.05,"user_extract_min":10,"lock_brokerage_timer":7,"sys_extension_type":0,"withdraw_type":["0","1","2"],"extract_switch":1,"transfer_scene_id":0,"max_bag_number":10}')
ON DUPLICATE KEY UPDATE `expire_time`=VALUES(`expire_time`),`result`=VALUES(`result`);
