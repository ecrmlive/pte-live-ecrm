-- 平台「商品统计」本地演示数据（可重复导入）。
-- 用法：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4' < sql/business/init_product_stat_demo.sql
-- 全部为明确标识的虚构中文记录，不包含真实用户隐私。
USE `qixi_crm_business`;
SET NAMES utf8mb4;
SET time_zone = '+08:00';

-- 丰富商品类型，供「商品类型」环形图展示。
UPDATE `qixi_crm_b_product_view` SET `product_type` = CASE
  WHEN `product_id` IN (1005, 1105) THEN 1
  WHEN `product_id` IN (1007) THEN 3
  WHEN `product_id` IN (1207) THEN 4
  ELSE `product_type`
END
WHERE `product_id` IN (1005, 1007, 1105, 1207);

-- 演示用户（商品统计专用区间，避免与大屏用户冲突）。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`,`created_at`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 24
)
SELECT 96100 + `n`, CONCAT('商品统计演示用户', LPAD(`n`, 2, '0')), CONCAT('PSTAT-DEMO-', LPAD(`n`, 4, '0')), 1, 0, 1,
  DATE_SUB(NOW(), INTERVAL (`n` % 14) DAY)
FROM `ps_seq`
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`updated_at`=NOW();

-- 近 14 天浏览量：制造环比差异（近 7 天更高）。
DELETE FROM `qixi_crm_b_user_browse_history` WHERE `id` BETWEEN 961001 AND 961400;
INSERT INTO `qixi_crm_b_user_browse_history` (`id`,`user_id`,`product_id`,`store_id`,`viewed_at`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 280
)
SELECT
  961000 + `n`,
  96100 + ((`n` % 24) + 1),
  ELT(((`n` - 1) % 6) + 1, 1001, 1002, 1101, 1103, 1201, 1003),
  ELT(((`n` - 1) % 6) + 1, 1, 1, 2, 2, 3, 1),
  DATE_ADD(
    DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 20) DAY),
    INTERVAL ((`n` * 37) % 20) HOUR
  ) + INTERVAL ((`n` * 11) % 50) MINUTE
FROM `ps_seq`;

-- 收藏：近 14 天少量记录。
DELETE FROM `qixi_crm_b_product_favorite` WHERE `user_id` BETWEEN 96101 AND 96124;
INSERT INTO `qixi_crm_b_product_favorite` (`user_id`,`product_id`,`created_at`) VALUES
  (96101,1001,DATE_SUB(NOW(), INTERVAL 1 DAY)),
  (96102,1002,DATE_SUB(NOW(), INTERVAL 2 DAY)),
  (96103,1101,DATE_SUB(NOW(), INTERVAL 3 DAY)),
  (96104,1201,DATE_SUB(NOW(), INTERVAL 4 DAY)),
  (96105,1003,DATE_SUB(NOW(), INTERVAL 5 DAY)),
  (96106,1103,DATE_SUB(NOW(), INTERVAL 8 DAY)),
  (96107,1001,DATE_SUB(NOW(), INTERVAL 9 DAY)),
  (96108,1201,DATE_SUB(NOW(), INTERVAL 10 DAY)),
  (96109,1002,DATE_SUB(NOW(), INTERVAL 11 DAY))
ON DUPLICATE KEY UPDATE `created_at`=VALUES(`created_at`);

-- 加购：近 14 天。
DELETE FROM `qixi_crm_b_cart` WHERE `id` BETWEEN 961001 AND 961080;
INSERT INTO `qixi_crm_b_cart` (`id`,`user_id`,`store_id`,`product_id`,`sku_key`,`quantity`,`created_at`,`updated_at`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 60
)
SELECT
  961000 + `n`,
  96100 + ((`n` % 24) + 1),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 1001, 1101, 1201),
  CONCAT('pstat-sku-', `n`),
  1 + (`n` % 3),
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 5) DAY), INTERVAL ((`n` * 5) % 18) HOUR),
  NOW()
FROM `ps_seq`;

-- 下单/支付：近 14 天订单（部分未支付制造支付数差异）。
DELETE FROM `qixi_crm_b_order_item` WHERE `id` BETWEEN 961001 AND 961120;
DELETE FROM `qixi_crm_b_order` WHERE `id` BETWEEN 961001 AND 961080;
DELETE FROM `qixi_crm_b_group_order` WHERE `id` BETWEEN 961001 AND 961080;

INSERT INTO `qixi_crm_b_group_order` (`id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,`paid_at`,`idempotency_key`,`remark`,`created_at`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 56
)
SELECT
  961000 + `n`,
  CONCAT('PSTAT-G-', LPAD(`n`, 4, '0')),
  96100 + ((`n` % 24) + 1),
  100 + (`n` % 9) * 20,
  0,
  0,
  100 + (`n` % 9) * 20,
  1 + (`n` % 2),
  0,
  0,
  JSON_OBJECT('name','商品统计演示收件人','mobile','13800000000','address','上海市演示地址'),
  'mock',
  IF(`n` % 5 = 0, 'pending', 'paid'),
  IF(`n` % 5 = 0, NULL, DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 4) DAY), INTERVAL 2 HOUR)),
  CONCAT('pstat-group-', `n`),
  '商品统计演示团单',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 4) DAY), INTERVAL ((`n` * 3) % 16) HOUR)
FROM `ps_seq`;

INSERT INTO `qixi_crm_b_order` (`id`,`group_order_id`,`order_no`,`merchant_id`,`merchant_name_snapshot`,`store_id`,`store_name_snapshot`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`remark`,`status`,`paid_at`,`created_at`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 56
)
SELECT
  961000 + `n`,
  961000 + `n`,
  CONCAT('PSTAT-O-', LPAD(`n`, 4, '0')),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰商户','CRM Live居家商户','CRM Live数码商户'),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰旗舰店','CRM Live居家优选店','CRM Live数码生活店'),
  96100 + ((`n` % 24) + 1),
  100 + (`n` % 9) * 20,
  0,
  0,
  100 + (`n` % 9) * 20,
  1 + (`n` % 2),
  0,
  0,
  JSON_OBJECT('name','商品统计演示收件人','mobile','13800000000','address','上海市演示地址'),
  '商品统计演示订单',
  IF(`n` % 5 = 0, 'pending_pay', 'paid'),
  IF(`n` % 5 = 0, NULL, DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 4) DAY), INTERVAL 2 HOUR)),
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 4) DAY), INTERVAL ((`n` * 3) % 16) HOUR)
FROM `ps_seq`;

INSERT INTO `qixi_crm_b_order_item` (`id`,`order_id`,`product_id`,`merchant_sku_id`,`sku_key`,`title_snapshot`,`cover_url_snapshot`,`spec_snapshot`,`unit_price`,`quantity`,`refund_quantity`) WITH RECURSIVE `ps_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `ps_seq` WHERE `n` < 56
)
SELECT
  961000 + `n`,
  961000 + `n`,
  ELT(((`n` - 1) % 3) + 1, 1001, 1101, 1201),
  0,
  CONCAT('pstat-item-', `n`),
  ELT(((`n` - 1) % 3) + 1, '轻奢羊绒针织衫','无火藤条香氛礼盒','智能数显保温杯'),
  '/demo/product-knit-v1.png',
  JSON_OBJECT('默认','标准'),
  100 + (`n` % 9) * 20,
  1 + (`n` % 2),
  0
FROM `ps_seq`;

-- 少量退款（偏上期，便于近 7 天环比为 0 或负）。
DELETE FROM `qixi_crm_b_refund` WHERE `id` BETWEEN 961001 AND 961010;
INSERT INTO `qixi_crm_b_refund` (`id`,`order_id`,`refund_no`,`reason`,`amount`,`refund_type`,`order_status_before`,`status`,`idempotency_key`,`created_at`) VALUES
  (961001,961008,'PSTAT-R-0001','商品统计演示退款',120.00,'money_only','paid','refunded','pstat-refund-1',DATE_SUB(NOW(), INTERVAL 9 DAY)),
  (961002,961012,'PSTAT-R-0002','商品统计演示退款',140.00,'money_only','paid','applied','pstat-refund-2',DATE_SUB(NOW(), INTERVAL 10 DAY)),
  (961003,961016,'PSTAT-R-0003','商品统计演示退款',160.00,'money_only','paid','refunded','pstat-refund-3',DATE_SUB(NOW(), INTERVAL 2 DAY))
ON DUPLICATE KEY UPDATE `reason`=VALUES(`reason`),`amount`=VALUES(`amount`),`status`=VALUES(`status`),`created_at`=VALUES(`created_at`);
