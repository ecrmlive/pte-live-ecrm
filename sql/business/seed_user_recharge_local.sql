-- 用户充值记录本地演示数据（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live演示','DEMO-USER-9101',1,0,1),
  (9102,'微信用户','DEMO-USER-9102',1,0,1),
  (9104,'寇小雨','DEMO-USER-9104',1,0,1),
  (9105,'起泡酒','13800009105',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_user_profile` (`user_id`,`avatar_url`,`real_name`,`gender`,`bio`,`source_channel`) VALUES
  (9101,'/demo/avatar-9101.png','CRM Live演示',1,'','wechat'),
  (9102,'/demo/avatar-9102.png','微信用户',0,'','mini_program'),
  (9104,'/demo/avatar-9104.png','寇小雨',2,'','mini_program'),
  (9105,'/demo/avatar-9105.png','起泡酒',2,'','wechat')
ON DUPLICATE KEY UPDATE
  `avatar_url`=VALUES(`avatar_url`),
  `real_name`=VALUES(`real_name`);

INSERT INTO `qixi_crm_b_member_account` (`user_id`,`balance`,`points`,`commission`) VALUES
  (9101,120.50,0,0),
  (9102,36.20,0,0),
  (9104,88.00,0,0),
  (9105,15.00,0,0)
ON DUPLICATE KEY UPDATE
  `balance`=VALUES(`balance`);

INSERT INTO `qixi_crm_b_user_recharge` (
  `recharge_id`,`uid`,`order_id`,`price`,`give_price`,`recharge_type`,`paid`,`pay_time`,`refund_price`,`create_time`
) VALUES
  (135,9105,'wxu20260518100135001',0.01,0.01,'routine',1,DATE_SUB(NOW(),INTERVAL 2 DAY),0.00,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (127,9102,'wxu20260517100127002',1.00,0.50,'weixin',1,DATE_SUB(NOW(),INTERVAL 5 DAY),0.00,DATE_SUB(NOW(),INTERVAL 5 DAY)),
  (126,9104,'wxu20260516100126003',0.10,0.00,'routine',1,DATE_SUB(NOW(),INTERVAL 6 DAY),0.00,DATE_SUB(NOW(),INTERVAL 6 DAY)),
  (125,9101,'wxu20260515100125004',50.00,5.00,'h5',1,DATE_SUB(NOW(),INTERVAL 8 DAY),0.00,DATE_SUB(NOW(),INTERVAL 8 DAY)),
  (124,9105,'wxu20260514100124005',20.00,2.00,'alipay',1,DATE_SUB(NOW(),INTERVAL 10 DAY),0.00,DATE_SUB(NOW(),INTERVAL 10 DAY)),
  (123,9102,'wxu20260513100123006',8.80,0.00,'routine',0,NULL,0.00,DATE_SUB(NOW(),INTERVAL 12 DAY)),
  (122,9104,'wxu20260512100122007',1.14,0.00,'routine',1,DATE_SUB(NOW(),INTERVAL 1 DAY),0.00,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (121,9101,'wxu20260511100121008',30.00,3.00,'weixin',1,DATE_SUB(NOW(),INTERVAL 15 DAY),30.00,DATE_SUB(NOW(),INTERVAL 15 DAY))
ON DUPLICATE KEY UPDATE
  `uid`=VALUES(`uid`),
  `order_id`=VALUES(`order_id`),
  `price`=VALUES(`price`),
  `give_price`=VALUES(`give_price`),
  `recharge_type`=VALUES(`recharge_type`),
  `paid`=VALUES(`paid`),
  `pay_time`=VALUES(`pay_time`),
  `refund_price`=VALUES(`refund_price`);
