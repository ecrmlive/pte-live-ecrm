-- 平台「用户统计」本地演示数据（可重复导入）。
-- 用法：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4' < sql/business/init_user_stat_demo.sql
-- 全部为明确标识的虚构中文记录，不包含真实用户隐私。
USE `qixi_crm_business`;
SET NAMES utf8mb4;
SET time_zone = '+08:00';

-- 演示用户（用户统计专用区间，避免与商品统计/大屏冲突）。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`,`created_at`) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 80
)
SELECT
  97100 + `n`,
  CONCAT('用户统计演示', LPAD(`n`, 2, '0')),
  CONCAT('USTAT-DEMO-', LPAD(`n`, 4, '0')),
  1,
  0,
  1,
  DATE_ADD(
    DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 14) DAY),
    INTERVAL ((`n` * 17) % 20) HOUR
  ) + INTERVAL ((`n` * 7) % 50) MINUTE
FROM `us_seq`
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`created_at`=VALUES(`created_at`),`updated_at`=NOW();

-- 近 14 天浏览：制造活跃用户与环比差异。
DELETE FROM `qixi_crm_b_user_browse_history` WHERE `id` BETWEEN 971001 AND 971500;
INSERT INTO `qixi_crm_b_user_browse_history` (`id`,`user_id`,`product_id`,`store_id`,`viewed_at`) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 360
)
SELECT
  971000 + `n`,
  97100 + ((`n` % 80) + 1),
  ELT(((`n` - 1) % 6) + 1, 1001, 1002, 1101, 1103, 1201, 1003),
  ELT(((`n` - 1) % 6) + 1, 1, 1, 2, 2, 3, 1),
  DATE_ADD(
    DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 26) DAY),
    INTERVAL ((`n` * 37) % 20) HOUR
  ) + INTERVAL ((`n` * 11) % 50) MINUTE
FROM `us_seq`;

-- 付费会员现状。
DELETE FROM `qixi_crm_b_user_svip` WHERE `user_id` BETWEEN 97101 AND 97180;
INSERT INTO `qixi_crm_b_user_svip` (`user_id`,`status`,`expires_at`,`updated_at`) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 28
)
SELECT
  97100 + `n`,
  IF(`n` % 7 = 0, 'lifetime', 'period'),
  IF(`n` % 7 = 0, NULL, DATE_ADD(NOW(), INTERVAL (30 + `n`) DAY)),
  NOW()
FROM `us_seq`;

-- 新增付费会员订单（近 14 天部分已支付）。
DELETE FROM `qixi_crm_b_svip_order` WHERE `id` BETWEEN 971001 AND 971040;
INSERT INTO `qixi_crm_b_svip_order` (
  `id`,`order_no`,`user_id`,`plan_id`,`plan_name`,`plan_type`,`duration_days`,`amount`,`status`,`idempotency_key`,`created_at`,`paid_at`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 24
)
SELECT
  971000 + `n`,
  CONCAT('SVIP-USTAT-', LPAD(`n`, 4, '0')),
  97100 + `n`,
  980001,
  'SVIP 月度会员',
  'period',
  30,
  29.00,
  IF(`n` % 5 = 0, 'pending', 'paid'),
  CONCAT('fixture-ustat-svip-', `n`),
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 12) DAY), INTERVAL 9 HOUR),
  IF(`n` % 5 = 0, NULL, DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 12) DAY), INTERVAL 10 HOUR))
FROM `us_seq`;

-- 成交用户：老用户（历史首次支付）+ 新用户（近窗首次支付）。
DELETE FROM `qixi_crm_b_order_item` WHERE `id` BETWEEN 971001 AND 971120;
DELETE FROM `qixi_crm_b_order` WHERE `id` BETWEEN 971001 AND 971080;
DELETE FROM `qixi_crm_b_group_order` WHERE `id` BETWEEN 971001 AND 971080;

-- 老用户：14 天前已有首单。
INSERT INTO `qixi_crm_b_group_order` (
  `id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,
  `activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,`paid_at`,`idempotency_key`,`remark`,`created_at`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 20
)
SELECT
  971000 + `n`,
  CONCAT('GO-USTAT-OLD-', LPAD(`n`, 4, '0')),
  97100 + `n`,
  100.00, 0, 0, 100.00, 1, 0, 0,
  JSON_OBJECT('name', CONCAT('老用户收件人', `n`), 'mobile', CONCAT('1380000', LPAD(`n`, 4, '0')), 'address', '上海市演示地址'),
  'mock', 'paid',
  DATE_SUB(CURDATE(), INTERVAL (20 + (`n` % 5)) DAY) + INTERVAL 11 HOUR,
  CONCAT('fixture-ustat-old-first-', `n`),
  '用户统计-老用户首单',
  DATE_SUB(CURDATE(), INTERVAL (20 + (`n` % 5)) DAY) + INTERVAL 10 HOUR
FROM `us_seq`;

-- 老用户近 7 天复购。
INSERT INTO `qixi_crm_b_group_order` (
  `id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,
  `activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,`paid_at`,`idempotency_key`,`remark`,`created_at`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 20
)
SELECT
  971020 + `n`,
  CONCAT('GO-USTAT-OLD2-', LPAD(`n`, 4, '0')),
  97100 + `n`,
  88.00, 0, 0, 88.00, 1, 0, 0,
  JSON_OBJECT('name', CONCAT('老用户收件人', `n`), 'mobile', CONCAT('1380000', LPAD(`n`, 4, '0')), 'address', '上海市演示地址'),
  'mock', 'paid',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 7) DAY), INTERVAL 14 HOUR),
  CONCAT('fixture-ustat-old-repeat-', `n`),
  '用户统计-老用户复购',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 7) DAY), INTERVAL 13 HOUR)
FROM `us_seq`;

-- 新用户：近 7 天首次支付。
INSERT INTO `qixi_crm_b_group_order` (
  `id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,
  `activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,`paid_at`,`idempotency_key`,`remark`,`created_at`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 18
)
SELECT
  971040 + `n`,
  CONCAT('GO-USTAT-NEW-', LPAD(`n`, 4, '0')),
  97120 + `n`,
  66.00, 0, 0, 66.00, 1, 0, 0,
  JSON_OBJECT('name', CONCAT('新用户收件人', `n`), 'mobile', CONCAT('1390000', LPAD(`n`, 4, '0')), 'address', '上海市演示地址'),
  'mock', 'paid',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 7) DAY), INTERVAL 15 HOUR),
  CONCAT('fixture-ustat-new-first-', `n`),
  '用户统计-新用户首单',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) % 7) DAY), INTERVAL 14 HOUR)
FROM `us_seq`;

INSERT INTO `qixi_crm_b_order` (
  `id`,`group_order_id`,`order_no`,`merchant_id`,`merchant_name_snapshot`,`store_id`,`store_name_snapshot`,
  `user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,
  `points_amount`,`recipient_snapshot`,`remark`,`status`,`paid_at`,`created_at`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 58
)
SELECT
  971000 + `n`,
  971000 + `n`,
  CONCAT('ORD-USTAT-', LPAD(`n`, 4, '0')),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰商户','CRM Live居家商户','CRM Live数码商户'),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰旗舰店','CRM Live居家优选店','CRM Live数码生活店'),
  CASE
    WHEN `n` <= 20 THEN 97100 + `n`
    WHEN `n` <= 40 THEN 97100 + (`n` - 20)
    ELSE 97120 + (`n` - 40)
  END,
  CASE WHEN `n` <= 20 THEN 100.00 WHEN `n` <= 40 THEN 88.00 ELSE 66.00 END,
  0, 0,
  CASE WHEN `n` <= 20 THEN 100.00 WHEN `n` <= 40 THEN 88.00 ELSE 66.00 END,
  1, 0, 0,
  JSON_OBJECT('name','用户统计演示收件人','mobile','13800000000','address','上海市演示地址'),
  '用户统计演示订单',
  'paid',
  (SELECT `paid_at` FROM `qixi_crm_b_group_order` WHERE `id` = 971000 + `n`),
  (SELECT `created_at` FROM `qixi_crm_b_group_order` WHERE `id` = 971000 + `n`)
FROM `us_seq`;

INSERT INTO `qixi_crm_b_order_item` (
  `id`,`order_id`,`product_id`,`merchant_sku_id`,`sku_key`,`title_snapshot`,`cover_url_snapshot`,`spec_snapshot`,`unit_price`,`quantity`,`refund_quantity`
) WITH RECURSIVE `us_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `us_seq` WHERE `n` < 58
)
SELECT
  971000 + `n`,
  971000 + `n`,
  ELT(((`n` - 1) % 3) + 1, 1001, 1101, 1201),
  0,
  CONCAT('ustat-item-', `n`),
  CONCAT('用户统计演示商品', `n`),
  '/demo/product-knit-v1.png',
  JSON_OBJECT('默认','标准'),
  CASE WHEN `n` <= 20 THEN 100.00 WHEN `n` <= 40 THEN 88.00 ELSE 66.00 END,
  1,
  0
FROM `us_seq`;
