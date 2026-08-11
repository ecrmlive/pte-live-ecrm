-- 平台「发票列表」本地演示数据（幂等，utf8mb4）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_order_invoice`
  (`id`,`order_id`,`invoice_profile_id`,`profile_type`,`invoice_type`,`receipt_sn`,`invoice_amount`,
   `title`,`tax_no`,`email`,`status`,`invoice_no`,`file_url`,`rejection_reason`,`mark`,`requested_at`,`issued_at`)
VALUES
  (9611001,961001,9601,'personal',1,'PT96100120260808001',120.00,
   '奋斗到底','','541251@qq.com','requested','','','','','2026-08-08 10:00:00',NULL),
  (9611002,961002,9601,'enterprise',1,'PT96100220260808002',140.00,
   '111','443','123@qq.com','requested','','','','','2026-08-08 11:00:00',NULL),
  (9611003,961003,9601,'personal',1,'PT96100320260808003',160.00,
   '演示用户04','','demo04@invoice.invalid','issued','INV-961003-001','','','本地已开票演示','2026-08-08 12:00:00','2026-08-08 15:00:00'),
  (9611004,961004,9601,'enterprise',2,'PT96100420260808004',180.00,
   'CRM Live演示科技有限公司','91310000DEMO12345X','finance@invoice.invalid','requested','','','','专用发票演示','2026-08-07 09:30:00',NULL),
  (9611005,961006,9601,'personal',1,'PT96100620260808005',99.00,
   '演示用户07','','demo07@invoice.invalid','issued','INV-961006-001','','','','2026-08-06 14:20:00','2026-08-06 16:00:00'),
  (9611006,961007,9601,'enterprise',1,'PT96100720260808006',210.00,
   'CRM Live演示科技有限公司','91310000DEMO12345X','finance@invoice.invalid','rejected','','','资料不齐','','2026-08-05 10:10:00',NULL),
  (9611007,961008,9601,'personal',1,'PT96100820260808007',88.50,
   '演示用户09','','demo09@invoice.invalid','issued','INV-961008-001','','','','2026-08-04 08:40:00','2026-08-04 11:00:00'),
  (9611008,961009,9601,'enterprise',1,'PT96100920260808008',256.00,
   'CRM Live演示科技有限公司','91310000DEMO12345X','finance@invoice.invalid','requested','','','','','2026-08-03 16:00:00',NULL)
ON DUPLICATE KEY UPDATE
  `order_id`=VALUES(`order_id`),
  `invoice_profile_id`=VALUES(`invoice_profile_id`),
  `profile_type`=VALUES(`profile_type`),
  `invoice_type`=VALUES(`invoice_type`),
  `receipt_sn`=VALUES(`receipt_sn`),
  `invoice_amount`=VALUES(`invoice_amount`),
  `title`=VALUES(`title`),
  `tax_no`=VALUES(`tax_no`),
  `email`=VALUES(`email`),
  `status`=VALUES(`status`),
  `invoice_no`=VALUES(`invoice_no`),
  `file_url`=VALUES(`file_url`),
  `rejection_reason`=VALUES(`rejection_reason`),
  `mark`=VALUES(`mark`),
  `requested_at`=VALUES(`requested_at`),
  `issued_at`=VALUES(`issued_at`);
