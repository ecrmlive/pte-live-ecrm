-- 秒杀活动「已加商品」展开多规格演示：为常用商品 1001 补充双 SKU 投影
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_product_sku_view`
  (`merchant_sku_id`,`product_id`,`sku_key`,`spec_snapshot`,`price`,`stock`,`sale_status`,`version`,`updated_at`)
VALUES
  (61001,1001,'61001',JSON_OBJECT('颜色','米白','规格','单人'),299.00,30,1,1,NOW()),
  (61011,1001,'61011',JSON_OBJECT('颜色','浅蓝','规格','双人'),359.00,20,1,1,NOW()),
  (61002,1002,'61002',JSON_OBJECT('默认','标准'),469.00,50,1,1,NOW()),
  (61012,1002,'61012',JSON_OBJECT('颜色','灰色','规格','加大'),499.00,18,1,1,NOW())
ON DUPLICATE KEY UPDATE
  `sku_key`=VALUES(`sku_key`),
  `spec_snapshot`=VALUES(`spec_snapshot`),
  `price`=VALUES(`price`),
  `stock`=VALUES(`stock`),
  `sale_status`=VALUES(`sale_status`),
  `version`=VALUES(`version`),
  `updated_at`=NOW();
