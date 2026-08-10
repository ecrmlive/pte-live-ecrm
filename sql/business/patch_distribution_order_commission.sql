-- 分销订单演示佣金：挂到已有演示订单，供平台「分销订单」列表验收
-- 用法：make local-sync-sql
SET NAMES utf8mb4;
USE `qixi_crm_business`;

-- 买家 9101 的上级推广人 9001 获得佣金；订单须已存在
INSERT INTO `qixi_crm_b_commission_ledger`
  (`user_id`,`order_id`,`amount`,`status`,`idempotency_key`,`available_at`,`created_at`)
VALUES
  (9001,9900201,29.90,'available','fixture-commission-order-9900201',NOW(),DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (9001,9900202,19.90,'pending','fixture-commission-order-9900202',NULL,NOW())
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `order_id`=VALUES(`order_id`),
  `amount`=VALUES(`amount`),
  `status`=VALUES(`status`),
  `available_at`=VALUES(`available_at`);
