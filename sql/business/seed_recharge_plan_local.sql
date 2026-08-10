-- 本地演示：余额充值档位（幂等，utf8mb4）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_recharge_plan` (`id`,`name`,`amount`,`bonus_amount`,`status`,`sort`,`version`,`created_at`) VALUES
  (980001,'1元充值',1.00,0.00,1,1,1,'2026-08-01 10:00:00'),
  (980002,'5元充值',5.00,0.50,1,2,1,'2026-08-01 10:05:00'),
  (980003,'10元充值',10.00,1.00,1,3,1,'2026-08-01 10:10:00'),
  (980004,'50元充值',50.00,5.00,1,4,1,'2026-08-02 09:00:00'),
  (980005,'100元充值',100.00,10.00,1,5,1,'2026-08-02 09:10:00'),
  (980006,'400元充值',400.00,40.00,1,6,1,'2026-08-03 11:00:00')
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `amount`=VALUES(`amount`),
  `bonus_amount`=VALUES(`bonus_amount`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `updated_at`=NOW();
