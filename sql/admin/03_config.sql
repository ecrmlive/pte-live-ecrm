USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_config` (`config_key`,`config_value`) VALUES
  ('system.locale', JSON_OBJECT('default','zh-CN')),
  ('system.timezone', JSON_OBJECT('default','Asia/Shanghai'))
ON DUPLICATE KEY UPDATE `config_value`=VALUES(`config_value`);
