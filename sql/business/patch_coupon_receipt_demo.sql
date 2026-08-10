-- 平台「领取记录」演示夹具（脱敏手机号 + 多来源）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

UPDATE `qixi_crm_b_user`
SET `nickname`='微信用户U474322', `mobile`='18600009486'
WHERE `id`=9101;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`) VALUES
  (9102,'微信用户',NULL,1,0),
  (9103,'演示买家',NULL,1,0)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_coupon_user` (`id`,`user_id`,`coupon_id`,`source`,`status`,`obtained_at`,`used_order_id`) VALUES
  (93001,9101,3001,'onboarding','unused',DATE_SUB(NOW(), INTERVAL 2 DAY),NULL),
  (93002,9101,3003,'fixture','used',DATE_SUB(NOW(), INTERVAL 4 DAY),9900201),
  (93003,9101,3002,'platform_manual','expired',DATE_SUB(NOW(), INTERVAL 1 DAY),NULL),
  (93004,9102,3001,'onboarding','unused',DATE_SUB(NOW(), INTERVAL 3 HOUR),NULL),
  (93005,9102,3002,'receive','unused',DATE_SUB(NOW(), INTERVAL 1 HOUR),NULL),
  (93006,9103,3005,'fixture','unused',DATE_SUB(NOW(), INTERVAL 6 HOUR),NULL)
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `coupon_id`=VALUES(`coupon_id`),
  `source`=VALUES(`source`),
  `status`=VALUES(`status`),
  `obtained_at`=VALUES(`obtained_at`),
  `used_order_id`=VALUES(`used_order_id`);
