-- 平台「积分配置」夹具（对齐 CRMEB UserIntegral::getConfig / saveConfig）
-- 存储：qixi_crm_a_setting_cache.key = integral_config
-- 字段：integral_status / integral_money / integral_order_rate / integral_freeze /
--       integral_clear_time / integral_user_give / integral_community_give /
--       integral_community_give_limit / rule（积分说明，对齐 sys_integral_rule）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('integral_config',0,'{"integral_status":1,"integral_money":0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10,"integral_community_give_limit":10,"rule":""}')
ON DUPLICATE KEY UPDATE `expire_time`=VALUES(`expire_time`),`result`=VALUES(`result`);
