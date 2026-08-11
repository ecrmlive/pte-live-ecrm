-- 公众号「微信菜单」本地演示配置（对齐 CRMEB 截图默认态）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('wechat_menus',0,'[{"name":"商城","type":"view","url":"https://mer.crmeb.net","sub_button":[]},{"name":"一级菜单","type":"click","key":"MENU_LEVEL_1","sub_button":[]}]')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=VALUES(`result`);
