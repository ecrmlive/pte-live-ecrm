-- 用户提现本地演示数据（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9102,'微信用户','DEMO-USER-9102',1,0,1),
  (9103,'已注销用户','DEMO-USER-9103',0,0,1),
  (9104,'寇小雨','DEMO-USER-9104',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_user_extract` (
  `extract_id`,`uid`,`extract_sn`,`real_name`,`extract_type`,
  `bank_code`,`bank_address`,`bank_name`,`alipay_code`,`wechat`,`extract_pic`,
  `extract_price`,`balance`,`mark`,`admin_id`,`fail_msg`,`status`,`status_time`,`create_time`
) VALUES
  (97001,9104,'UE202608110001','寇小雨',0,
   '6217****1234','上海市浦东新区','中国银行','','','',
   1.00,88.50,'',9001,'',1,DATE_SUB(NOW(),INTERVAL 2 DAY),DATE_SUB(NOW(),INTERVAL 3 DAY)),
  (97002,9102,'UE202608110002','',4,
   '','','','','','',
   0.01,12.30,'',0,'',0,NULL,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (97003,9101,'UE202608110003','CRM Live演示',3,
   '','','','','wx_demo_9101','/demo/extract-qr-9101.png',
   12.00,54.20,'',9001,'',1,DATE_SUB(NOW(),INTERVAL 5 DAY),DATE_SUB(NOW(),INTERVAL 6 DAY)),
  (97004,9103,'UE202608110004','',0,
   '6222****5678','上海市徐汇区','工商银行','','','',
   20.00,0.00,'',9001,'信息有误,请完善',-1,DATE_SUB(NOW(),INTERVAL 4 DAY),DATE_SUB(NOW(),INTERVAL 5 DAY)),
  (97005,9102,'UE202608110005','微信用户',2,
   '','','','ali****@demo.invalid','','/demo/extract-qr-9102.png',
   8.80,3.20,'',0,'',0,NULL,DATE_SUB(NOW(),INTERVAL 8 HOUR)),
  (97006,9101,'UE202608110006','CRM Live演示',1,
   '','','','','wx_demo_9101','',
   5.50,48.70,'',9001,'',1,DATE_SUB(NOW(),INTERVAL 10 DAY),DATE_SUB(NOW(),INTERVAL 11 DAY)),
  (97007,9104,'UE202608110007','寇小雨',0,
   '6217****8899','上海市静安区','建设银行','','','',
   15.00,73.50,'',0,'',0,NULL,DATE_SUB(NOW(),INTERVAL 3 HOUR)),
  (97008,9103,'UE202608110008','',4,
   '','','','','','',
   2.00,0.00,'',9001,'账户异常',-1,DATE_SUB(NOW(),INTERVAL 7 DAY),DATE_SUB(NOW(),INTERVAL 8 DAY))
ON DUPLICATE KEY UPDATE
  `uid`=VALUES(`uid`),
  `extract_sn`=VALUES(`extract_sn`),
  `real_name`=VALUES(`real_name`),
  `extract_type`=VALUES(`extract_type`),
  `bank_code`=VALUES(`bank_code`),
  `bank_address`=VALUES(`bank_address`),
  `bank_name`=VALUES(`bank_name`),
  `alipay_code`=VALUES(`alipay_code`),
  `wechat`=VALUES(`wechat`),
  `extract_pic`=VALUES(`extract_pic`),
  `extract_price`=VALUES(`extract_price`),
  `balance`=VALUES(`balance`),
  `mark`=VALUES(`mark`),
  `admin_id`=VALUES(`admin_id`),
  `fail_msg`=VALUES(`fail_msg`),
  `status`=VALUES(`status`),
  `status_time`=VALUES(`status_time`);
