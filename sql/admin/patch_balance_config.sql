-- 平台「余额设置」夹具（对齐 CRMEB systemForm/Basics/balance）
-- 存储：qixi_crm_a_setting_cache.key = balance_config
-- 字段：balance_func_status / recharge_switch / store_user_min_recharge / recharge_attention
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('balance_config',0,'{"balance_func_status":1,"recharge_switch":1,"store_user_min_recharge":1,"recharge_attention":"1、账户充值仅限用于购买商城内商品，不可提现\\n2、账户充值成功后，一般1～5分钟到账\\n3、如有疑问，请联系客服"}')
ON DUPLICATE KEY UPDATE `expire_time`=VALUES(`expire_time`),`result`=VALUES(`result`);
