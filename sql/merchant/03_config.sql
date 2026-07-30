USE `qixi_crm_merchant`;
INSERT INTO `qixi_crm_m_config` (`store_id`,`config_key`,`config_value`) VALUES
  (1,'store.order.auto_accept',JSON_OBJECT('enabled',false))
ON DUPLICATE KEY UPDATE `config_value`=VALUES(`config_value`);
