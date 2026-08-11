-- 平台「资金记录」本地演示数据（幂等，utf8mb4）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live演示','DEMO-USER-9101',1,0,1),
  (9102,'微信用户','DEMO-USER-9102',1,0,1),
  (9104,'寇小雨','DEMO-USER-9104',1,0,1),
  (9106,'行元','13500004334',1,0,1),
  (9107,'微信用户','18100008838',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_user_bill`
  (`bill_id`,`uid`,`link_id`,`pm`,`title`,`category`,`type`,`number`,`balance`,`mark`,`mer_id`,`status`,`create_time`)
VALUES
  (882001,9106,'GIFT-20260811-01',1,'新人赠送余额','now_money','sys_inc_money',5000.00,5000.00,'系统增加了5000余额',0,1,'2026-08-11 10:23:39'),
  (882002,9107,'GIFT-20260811-02',1,'新人赠送余额','now_money','sys_inc_money',5000.00,5000.00,'系统增加了5000余额',0,1,'2026-08-11 10:22:18'),
  (882003,9102,'GIFT-20260811-03',1,'新人赠送余额','now_money','sys_inc_money',5000.00,5036.20,'系统增加了5000余额',0,1,'2026-08-11 10:20:05'),
  (882004,9101,'RCH-20260810-01',1,'余额充值','now_money','recharge',100.00,220.50,'本地演示充值',0,1,'2026-08-10 15:00:00'),
  (882005,9104,'PAY-20260809-01',0,'购买商品','now_money','pay_product',28.00,60.00,'余额支付订单（演示）',0,1,'2026-08-09 18:30:00'),
  (882006,9106,'DEC-20260808-01',0,'系统减少余额','now_money','sys_dec_money',20.00,4980.00,'本地演示系统扣余额',0,1,'2026-08-08 11:10:00'),
  (882007,9101,'SVIP-20260807-01',0,'付费会员支付','svip_pay','svip_pay',99.00,121.50,'开通付费会员（演示）',0,1,'2026-08-07 09:40:00'),
  (882008,9107,'BRK-20260806-01',1,'佣金转入余额','now_money','brokerage',12.50,5012.50,'佣金转余额（演示）',0,1,'2026-08-06 14:20:00')
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
  `status`=VALUES(`status`),
  `create_time`=VALUES(`create_time`);
