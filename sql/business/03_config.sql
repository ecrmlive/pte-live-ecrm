USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_config` (`config_key`,`config_value`) VALUES
  ('trade.currency', JSON_OBJECT('code','CNY','scale',2)),
  ('trade.order_timeout_minutes', JSON_OBJECT('value',30))
ON DUPLICATE KEY UPDATE `config_value`=VALUES(`config_value`);
