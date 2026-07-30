USE `qixi_crm_merchant`;
-- 不初始化店铺账号、密码或任何真实手机号。仅提供无个人信息的店铺读写模型夹具。

INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`status`,`version`) VALUES
  (1001,1,'七禧臻选羊绒针织衫',101,'on_sale',1),
  (1002,1,'七禧头层牛皮通勤包',101,'on_sale',1),
  (1003,1,'七禧轻量跑步鞋',106,'on_sale',1),
  (1004,1,'七禧智能保温杯',103,'on_sale',1),
  (1005,1,'七禧真丝方巾',104,'on_sale',1),
  (1006,1,'七禧香氛礼盒',102,'on_sale',1)
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`category_id`=VALUES(`category_id`),`status`=VALUES(`status`),`version`=VALUES(`version`);

INSERT INTO `qixi_crm_m_product_sku` (`product_id`,`spec_json`,`price`,`stock`,`status`)
SELECT p.id, JSON_OBJECT('默认','标准'), CASE p.id
  WHEN 1001 THEN 299.00 WHEN 1002 THEN 469.00 WHEN 1003 THEN 369.00
  WHEN 1004 THEN 159.00 WHEN 1005 THEN 129.00 ELSE 239.00 END,
  50, 1
FROM `qixi_crm_m_product` AS p
WHERE p.id BETWEEN 1001 AND 1006
  AND NOT EXISTS (SELECT 1 FROM `qixi_crm_m_product_sku` AS s WHERE s.product_id = p.id);

INSERT INTO `qixi_crm_m_finance_ledger` (`store_id`,`entry_type`,`amount`,`reference_type`,`reference_id`,`idempotency_key`) VALUES
  (1,'test_seed',1000.00,'fixture','opening-balance','fixture-store-1-opening-balance')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`);
