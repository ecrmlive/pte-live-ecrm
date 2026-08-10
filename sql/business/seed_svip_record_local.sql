-- 本地演示：付费会员「会员记录」列表 + 汇总（幂等，utf8mb4）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 确保记录页所需列存在
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_order' AND COLUMN_NAME='pay_type')=0,
    'ALTER TABLE `qixi_crm_b_svip_order` ADD COLUMN `pay_type` varchar(32) NOT NULL DEFAULT \'\' AFTER `amount`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_svip_order' AND COLUMN_NAME='end_time')=0,
    'ALTER TABLE `qixi_crm_b_svip_order` ADD COLUMN `end_time` datetime DEFAULT NULL AFTER `paid_at`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 演示用户（含已过期会员）
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live体验用户','13800009101',1,0,1),
  (9121,'演示会员阿强','13900009121',1,0,1),
  (9122,'演示会员小满','13700009122',1,0,1),
  (9123,'演示会员过期用户','13600009123',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_user_profile` (`user_id`,`avatar_url`,`gender`,`bio`,`source_channel`) VALUES
  (9101,'',0,'会员记录演示','pc'),
  (9121,'',1,'会员记录演示','h5'),
  (9122,'',2,'会员记录演示','mini_program'),
  (9123,'',0,'会员记录演示-已过期','pc')
ON DUPLICATE KEY UPDATE
  `bio`=VALUES(`bio`),
  `source_channel`=VALUES(`source_channel`);

-- 当前会员状态（含已过期，供「累计已过期人数」）
INSERT INTO `qixi_crm_b_user_svip` (`user_id`,`status`,`expires_at`,`updated_at`) VALUES
  (9101,'period','2026-12-31 23:59:59',NOW()),
  (9121,'period','2026-11-10 18:00:00',NOW()),
  (9122,'lifetime',NULL,NOW()),
  (9123,'period','2026-01-15 12:00:00',NOW())
ON DUPLICATE KEY UPDATE
  `status`=VALUES(`status`),
  `expires_at`=VALUES(`expires_at`),
  `updated_at`=NOW();

-- 会员购买记录演示（含已支付 / 待支付 / 免费 / 平台赠送）
INSERT INTO `qixi_crm_b_svip_order`
  (`id`,`order_no`,`user_id`,`plan_id`,`plan_name`,`plan_type`,`duration_days`,`amount`,`pay_type`,`status`,`idempotency_key`,`created_at`,`paid_at`,`end_time`)
VALUES
  (981001,'SVIP-REC-20260810-001',9101,980001,'青铜会员','period',30,29.00,'weixin','paid','fixture-svip-rec-981001','2026-08-01 10:20:00','2026-08-01 10:21:30','2026-08-31 10:21:30'),
  (981002,'SVIP-REC-20260810-002',9121,980002,'白银会员','period',90,79.00,'alipay','paid','fixture-svip-rec-981002','2026-08-03 14:05:00','2026-08-03 14:06:12','2026-11-01 14:06:12'),
  (981003,'SVIP-REC-20260810-003',9122,980004,'钻石会员','lifetime',NULL,599.00,'routine','paid','fixture-svip-rec-981003','2026-08-05 09:00:00','2026-08-05 09:01:08',NULL),
  (981004,'SVIP-REC-20260810-004',9123,980001,'青铜会员','period',30,29.00,'weixin','paid','fixture-svip-rec-981004','2025-12-16 11:30:00','2025-12-16 11:31:00','2026-01-15 11:31:00'),
  (981005,'SVIP-REC-20260810-005',9101,980003,'黄金会员','period',365,299.00,'weixin','pending','fixture-svip-rec-981005','2026-08-09 16:40:00',NULL,NULL),
  (981006,'SVIP-REC-20260810-006',9121,980001,'青铜会员','period',30,0.00,'free','paid','fixture-svip-rec-981006','2026-07-20 08:15:00','2026-07-20 08:15:00','2026-08-19 08:15:00'),
  (981007,'SVIP-REC-20260810-007',9122,980002,'白银会员','period',90,0.00,'sys','paid','fixture-svip-rec-981007','2026-07-25 19:22:00','2026-07-25 19:22:00','2026-10-23 19:22:00')
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `plan_id`=VALUES(`plan_id`),
  `plan_name`=VALUES(`plan_name`),
  `plan_type`=VALUES(`plan_type`),
  `duration_days`=VALUES(`duration_days`),
  `amount`=VALUES(`amount`),
  `pay_type`=VALUES(`pay_type`),
  `status`=VALUES(`status`),
  `created_at`=VALUES(`created_at`),
  `paid_at`=VALUES(`paid_at`),
  `end_time`=VALUES(`end_time`);
