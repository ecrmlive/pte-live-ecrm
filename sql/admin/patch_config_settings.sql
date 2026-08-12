SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 配置管理跨分类展示配置项，补齐控件与后台类型字段；不导入任何密钥。
-- 字段结构由 patch_config_classification.sql 负责；此处只维护默认配置项。

INSERT IGNORE INTO `qixi_crm_a_config_classification_item`
  (`id`,`classification_id`,`name`,`config_key`,`field_type`,`backend_type`,`content`,`description`,`status`,`sort`,`is_del`) VALUES
  (81501,81407,'网站名称','site_name','input',0,'商城','平台网站名称',1,100,0),
  (81502,81407,'网站开启','site_open','switch',0,'1','关闭后仅允许平台后台访问',1,90,0),
  (81503,81407,'自动解析复制口令','parse_copy_command','switch',0,'1','开启后小程序和 App 自动读取剪贴板口令',1,80,0),
  (81504,81407,'默认赠送复制次数','default_copy_count','number',1,'8','默认给商户赠送的商品采集次数',1,70,0),
  (81505,81402,'启用短信验证码','sms_verification_enabled','switch',0,'1','平台短信验证码开关',1,100,0),
  (81506,81405,'小程序名称','mini_program_name','input',0,'商城小程序','移动端小程序显示名称',1,100,0),
  (81507,81408,'默认客服类型','default_customer_service_type','radio',0,'system','默认由平台系统客服在线接待',1,100,0),
  (81508,81409,'订单支付成功通知','order_payment_notice_enabled','switch',0,'1','订单支付成功后通知会员',1,100,0),
  (81509,81403,'商户入驻开关','merchant_admission_enabled','switch',0,'1','允许商户提交入驻申请',1,100,0),
  (81510,81406,'余额功能开关','balance_enabled','switch',0,'0','平台余额与充值功能开关',1,100,0);
