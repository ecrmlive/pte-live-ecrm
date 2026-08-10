-- 订单分账本地演示数据（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_store_order_profitsharing` (
  `profitsharing_id`,`profitsharing_sn`,`order_id`,`order_sn`,`mer_id`,`mer_name`,
  `transaction_id`,`profitsharing_price`,`profitsharing_refund`,`profitsharing_mer_price`,
  `type`,`status`,`error_msg`,`profitsharing_time`,`create_time`,`is_combine`
) VALUES
  (98001,'pr202608100001',9900201,'CS-DEMO-O-20260803-001',1,'CRM Live服饰旗舰店',
   'mock-tx-ps-98001',299.00,0.00,8.97,'order',1,'',DATE_SUB(NOW(),INTERVAL 2 DAY),DATE_SUB(NOW(),INTERVAL 3 DAY),1),
  (98002,'pr202608100002',9900202,'CS-DEMO-O-20260803-002',1,'CRM Live服饰旗舰店',
   'mock-tx-ps-98002',199.00,0.00,5.97,'order',0,'',NULL,DATE_SUB(NOW(),INTERVAL 1 DAY),1),
  (98003,'pr202608100003',9900201,'CS-DEMO-O-20260803-001',1,'CRM Live服饰旗舰店',
   'mock-tx-ps-98003',50.00,0.00,1.50,'presell',2,'',NULL,DATE_SUB(NOW(),INTERVAL 12 HOUR),1),
  (98004,'pr202608100004',9900202,'CS-DEMO-O-20260803-002',2,'CRM Live居家优选店',
   'mock-tx-ps-98004',139.00,139.00,0.00,'order',-1,'',DATE_SUB(NOW(),INTERVAL 5 DAY),DATE_SUB(NOW(),INTERVAL 6 DAY),1),
  (98005,'pr202608100005',9900201,'CS-DEMO-O-20260803-001',3,'CRM Live数码生活店',
   'mock-tx-ps-98005',219.00,0.00,6.57,'order',-2,'渠道返回：分账接收方未就绪',NULL,DATE_SUB(NOW(),INTERVAL 8 HOUR),1)
ON DUPLICATE KEY UPDATE
  `profitsharing_sn`=VALUES(`profitsharing_sn`),
  `order_id`=VALUES(`order_id`),
  `order_sn`=VALUES(`order_sn`),
  `mer_id`=VALUES(`mer_id`),
  `mer_name`=VALUES(`mer_name`),
  `transaction_id`=VALUES(`transaction_id`),
  `profitsharing_price`=VALUES(`profitsharing_price`),
  `profitsharing_refund`=VALUES(`profitsharing_refund`),
  `profitsharing_mer_price`=VALUES(`profitsharing_mer_price`),
  `type`=VALUES(`type`),
  `status`=VALUES(`status`),
  `error_msg`=VALUES(`error_msg`),
  `profitsharing_time`=VALUES(`profitsharing_time`),
  `is_combine`=VALUES(`is_combine`);
