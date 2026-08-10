-- 平台「平台账单」本地演示数据（幂等，utf8mb4 中文可读）
SET NAMES utf8mb4;

-- 充值 / 余额消费（统计卡：充值金额、充值消费金额）
INSERT INTO `qixi_crm_b_user_bill`
  (`bill_id`,`uid`,`link_id`,`pm`,`title`,`category`,`type`,`number`,`balance`,`mark`,`mer_id`,`status`,`create_time`)
VALUES
  (881001,9101,'RCH-20260801',1,'余额充值','now_money','recharge',200.00,236.50,'本地演示充值',0,1,'2026-08-01 10:00:00'),
  (881002,9102,'RCH-20260803',1,'系统增加余额','now_money','sys_inc_money',50.00,50.00,'本地演示系统加余额',0,1,'2026-08-03 11:00:00'),
  (881003,9101,'PAY-20260802',0,'购买商品','now_money','pay_product',80.00,156.50,'余额支付订单（演示）',0,1,'2026-08-02 15:30:00'),
  (881004,9103,'RCH-20260715',1,'余额充值','now_money','recharge',100.00,110.00,'七月演示充值',0,1,'2026-07-15 09:20:00'),
  (881005,9103,'PAY-20260716',0,'购买商品','now_money','pay_product',40.00,70.00,'七月余额消费（演示）',0,1,'2026-07-16 14:10:00')
ON DUPLICATE KEY UPDATE
  `uid`=VALUES(`uid`),
  `link_id`=VALUES(`link_id`),
  `pm`=VALUES(`pm`),
  `title`=VALUES(`title`),
  `category`=VALUES(`category`),
  `type`=VALUES(`type`),
  `number`=VALUES(`number`),
  `balance`=VALUES(`balance`),
  `mark`=VALUES(`mark`),
  `mer_id`=VALUES(`mer_id`),
  `status`=VALUES(`status`),
  `create_time`=VALUES(`create_time`);

-- 财务流水（列表日/月账单 + 汇总统计）
INSERT INTO `qixi_crm_b_financial_record`
  (`financial_record_id`,`financial_record_sn`,`order_id`,`order_sn`,`user_info`,`user_id`,`financial_type`,`financial_pm`,`number`,`type`,`mer_id`,`pay_type`,`create_time`)
VALUES
  -- 2026-08-01
  (890001,'FR-LOCAL-890001',70001,'ORD-LOCAL-70001','体验用户',9101,'order',1,328.00,2,1001,1,'2026-08-01 12:10:00'),
  (890002,'FR-LOCAL-890002',70001,'ORD-LOCAL-70001','体验用户',9101,'order_true',0,280.00,2,1001,1,'2026-08-01 12:10:01'),
  (890003,'FR-LOCAL-890003',70001,'ORD-LOCAL-70001','体验用户',9101,'order_charge',1,28.00,2,1001,1,'2026-08-01 12:10:02'),
  (890004,'FR-LOCAL-890004',70001,'ORD-LOCAL-70001','体验用户',9101,'brokerage_one',0,12.00,2,1001,1,'2026-08-01 12:10:03'),
  (890005,'FR-LOCAL-890005',70001,'ORD-LOCAL-70001','体验用户',9101,'order_platform_coupon',0,8.00,2,1001,1,'2026-08-01 12:10:04'),
  -- 2026-08-02 含线下收款
  (890006,'FR-LOCAL-890006',70002,'ORD-LOCAL-70002','阿强',9102,'order',1,199.00,2,1002,7,'2026-08-02 16:20:00'),
  (890007,'FR-LOCAL-890007',70002,'ORD-LOCAL-70002','阿强',9102,'order_true',0,170.00,2,1002,7,'2026-08-02 16:20:01'),
  (890008,'FR-LOCAL-890008',70002,'ORD-LOCAL-70002','阿强',9102,'order_charge',1,15.00,2,1002,7,'2026-08-02 16:20:02'),
  (890009,'FR-LOCAL-890009',70002,'ORD-LOCAL-70002','阿强',9102,'brokerage_two',0,6.00,2,1002,7,'2026-08-02 16:20:03'),
  -- 2026-08-03 退款
  (890010,'FR-LOCAL-890010',70003,'ORD-LOCAL-70003','小林',9103,'order',1,88.00,2,1001,1,'2026-08-03 09:40:00'),
  (890011,'FR-LOCAL-890011',70003,'ORD-LOCAL-70003','小林',9103,'order_true',0,75.00,2,1001,1,'2026-08-03 09:40:01'),
  (890012,'FR-LOCAL-890012',70003,'ORD-LOCAL-70003','小林',9103,'order_charge',1,7.00,2,1001,1,'2026-08-03 09:40:02'),
  (890013,'FR-LOCAL-890013',70003,'ORD-LOCAL-70003','小林',9103,'refund_order',0,88.00,2,1001,1,'2026-08-03 18:05:00'),
  (890014,'FR-LOCAL-890014',70003,'ORD-LOCAL-70003','小林',9103,'refund_charge',0,7.00,2,1001,1,'2026-08-03 18:05:01'),
  (890015,'FR-LOCAL-890015',70003,'ORD-LOCAL-70003','小林',9103,'refund_brokerage_one',1,3.00,2,1001,1,'2026-08-03 18:05:02'),
  (890016,'FR-LOCAL-890016',70003,'ORD-LOCAL-70003','小林',9103,'refund_platform_coupon',1,5.00,2,1001,1,'2026-08-03 18:05:03'),
  -- 2026-08-05
  (890017,'FR-LOCAL-890017',70004,'ORD-LOCAL-70004','体验用户',9101,'order',1,520.00,2,1003,1,'2026-08-05 11:15:00'),
  (890018,'FR-LOCAL-890018',70004,'ORD-LOCAL-70004','体验用户',9101,'order_true',0,450.00,2,1003,1,'2026-08-05 11:15:01'),
  (890019,'FR-LOCAL-890019',70004,'ORD-LOCAL-70004','体验用户',9101,'order_charge',1,40.00,2,1003,1,'2026-08-05 11:15:02'),
  (890020,'FR-LOCAL-890020',70004,'ORD-LOCAL-70004','体验用户',9101,'brokerage_one',0,18.00,2,1003,1,'2026-08-05 11:15:03'),
  (890021,'FR-LOCAL-890021',70004,'ORD-LOCAL-70004','体验用户',9101,'order_svip_coupon',0,12.00,2,1003,1,'2026-08-05 11:15:04'),
  -- 2026-08-08 线下 + 线上
  (890022,'FR-LOCAL-890022',70005,'ORD-LOCAL-70005','阿强',9102,'order',1,260.00,2,1002,1,'2026-08-08 13:30:00'),
  (890023,'FR-LOCAL-890023',70005,'ORD-LOCAL-70005','阿强',9102,'order_true',0,220.00,2,1002,1,'2026-08-08 13:30:01'),
  (890024,'FR-LOCAL-890024',70005,'ORD-LOCAL-70005','阿强',9102,'order_charge',1,22.00,2,1002,1,'2026-08-08 13:30:02'),
  (890025,'FR-LOCAL-890025',70006,'ORD-LOCAL-70006','小林',9103,'order',1,150.00,2,1001,7,'2026-08-08 19:10:00'),
  (890026,'FR-LOCAL-890026',70006,'ORD-LOCAL-70006','小林',9103,'order_true',0,130.00,2,1001,7,'2026-08-08 19:10:01'),
  (890027,'FR-LOCAL-890027',70006,'ORD-LOCAL-70006','小林',9103,'order_charge',1,12.00,2,1001,7,'2026-08-08 19:10:02'),
  -- 2026-07 月账单
  (890028,'FR-LOCAL-890028',70007,'ORD-LOCAL-70007','体验用户',9101,'order',1,680.00,2,1001,1,'2026-07-10 10:00:00'),
  (890029,'FR-LOCAL-890029',70007,'ORD-LOCAL-70007','体验用户',9101,'order_true',0,580.00,2,1001,1,'2026-07-10 10:00:01'),
  (890030,'FR-LOCAL-890030',70007,'ORD-LOCAL-70007','体验用户',9101,'order_charge',1,55.00,2,1001,1,'2026-07-10 10:00:02'),
  (890031,'FR-LOCAL-890031',70007,'ORD-LOCAL-70007','体验用户',9101,'brokerage_one',0,25.00,2,1001,1,'2026-07-10 10:00:03'),
  (890032,'FR-LOCAL-890032',70007,'ORD-LOCAL-70007','体验用户',9101,'order_platform_coupon',0,20.00,2,1001,1,'2026-07-10 10:00:04'),
  (890033,'FR-LOCAL-890033',70008,'ORD-LOCAL-70008','阿强',9102,'order',1,210.00,2,1002,7,'2026-07-22 15:40:00'),
  (890034,'FR-LOCAL-890034',70008,'ORD-LOCAL-70008','阿强',9102,'order_true',0,180.00,2,1002,7,'2026-07-22 15:40:01'),
  (890035,'FR-LOCAL-890035',70008,'ORD-LOCAL-70008','阿强',9102,'order_charge',1,18.00,2,1002,7,'2026-07-22 15:40:02'),
  (890036,'FR-LOCAL-890036',70008,'ORD-LOCAL-70008','阿强',9102,'refund_order',0,50.00,2,1002,1,'2026-07-25 11:00:00')
ON DUPLICATE KEY UPDATE
  `financial_record_sn`=VALUES(`financial_record_sn`),
  `order_id`=VALUES(`order_id`),
  `order_sn`=VALUES(`order_sn`),
  `user_info`=VALUES(`user_info`),
  `user_id`=VALUES(`user_id`),
  `financial_type`=VALUES(`financial_type`),
  `financial_pm`=VALUES(`financial_pm`),
  `number`=VALUES(`number`),
  `type`=VALUES(`type`),
  `mer_id`=VALUES(`mer_id`),
  `pay_type`=VALUES(`pay_type`),
  `create_time`=VALUES(`create_time`);
