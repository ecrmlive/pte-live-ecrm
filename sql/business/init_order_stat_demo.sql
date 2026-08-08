-- 平台「订单统计」本地演示数据（可重复导入）。
-- 用法：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" exec mysql --protocol=socket -uroot --default-character-set=utf8mb4' < sql/business/init_order_stat_demo.sql
-- 全部为明确标识的虚构中文记录，不包含真实用户隐私。
USE `qixi_crm_business`;
SET NAMES utf8mb4;
SET time_zone = '+08:00';

-- 演示用户（订单统计专用区间，避免与商品统计/大屏冲突）。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`,`created_at`) WITH RECURSIVE `os_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `os_seq` WHERE `n` < 32
)
SELECT 96200 + `n`, CONCAT('订单统计演示用户', LPAD(`n`, 2, '0')), CONCAT('OSTAT-DEMO-', LPAD(`n`, 4, '0')), 1, 0, 1,
  DATE_SUB(NOW(), INTERVAL (`n` % 20) DAY)
FROM `os_seq`
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`updated_at`=NOW();

-- 清理本批演示订单相关行。
DELETE FROM `qixi_crm_b_order_delivery` WHERE `id` BETWEEN 962001 AND 962200;
DELETE FROM `qixi_crm_b_refund` WHERE `id` BETWEEN 962001 AND 962040;
DELETE FROM `qixi_crm_b_order_item` WHERE `id` BETWEEN 962001 AND 962120;
DELETE FROM `qixi_crm_b_order` WHERE `id` BETWEEN 962001 AND 962120;
DELETE FROM `qixi_crm_b_group_order` WHERE `id` BETWEEN 962001 AND 962120;

-- 近 14 天已支付订单：制造环比、活动类型与金额差异。
INSERT INTO `qixi_crm_b_group_order` (
  `id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,
  `total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,
  `paid_at`,`idempotency_key`,`remark`,`created_at`
) WITH RECURSIVE `os_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `os_seq` WHERE `n` < 84
)
SELECT
  962000 + `n`,
  CONCAT('OSTAT-G-', LPAD(`n`, 4, '0')),
  96200 + ((`n` % 32) + 1),
  80 + (`n` % 11) * 25,
  IF(`n` % 4 = 0, 10 + (`n` % 5) * 5, 0),
  IF(`n` % 6 = 0, 8, 0),
  GREATEST(50, 80 + (`n` % 11) * 25 - IF(`n` % 4 = 0, 10 + (`n` % 5) * 5, 0)),
  1 + (`n` % 2),
  ELT(((`n` - 1) % 8) + 1, 0, 1, 2, 3, 4, 20, 10, 0),
  IF(((`n` - 1) % 8) = 5, 100 + (`n` % 50), 0),
  JSON_OBJECT('name','订单统计演示收件人','mobile','13900000000','address','上海市订单统计演示地址'),
  'mock',
  'paid',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 6) DAY), INTERVAL 1 HOUR),
  CONCAT('ostat-group-', `n`),
  '订单统计演示团单',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 6) DAY), INTERVAL ((`n` * 5) % 18) HOUR)
FROM `os_seq`;

INSERT INTO `qixi_crm_b_order` (
  `id`,`group_order_id`,`order_no`,`merchant_id`,`merchant_name_snapshot`,`store_id`,`store_name_snapshot`,
  `user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,
  `points_amount`,`recipient_snapshot`,`remark`,`status`,`paid_at`,`created_at`
) WITH RECURSIVE `os_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `os_seq` WHERE `n` < 84
)
SELECT
  962000 + `n`,
  962000 + `n`,
  CONCAT('OSTAT-O-', LPAD(`n`, 4, '0')),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰商户','CRM Live居家商户','CRM Live数码商户'),
  ELT(((`n` - 1) % 3) + 1, 1, 2, 3),
  ELT(((`n` - 1) % 3) + 1, 'CRM Live服饰旗舰店','CRM Live居家优选店','CRM Live数码生活店'),
  96200 + ((`n` % 32) + 1),
  80 + (`n` % 11) * 25,
  IF(`n` % 4 = 0, 10 + (`n` % 5) * 5, 0),
  IF(`n` % 6 = 0, 8, 0),
  GREATEST(50, 80 + (`n` % 11) * 25 - IF(`n` % 4 = 0, 10 + (`n` % 5) * 5, 0)),
  1 + (`n` % 2),
  ELT(((`n` - 1) % 8) + 1, 0, 1, 2, 3, 4, 20, 10, 0),
  IF(((`n` - 1) % 8) = 5, 100 + (`n` % 50), 0),
  JSON_OBJECT('name','订单统计演示收件人','mobile','13900000000','address','上海市订单统计演示地址'),
  '订单统计演示订单',
  ELT(((`n` - 1) % 4) + 1, 'paid', 'fulfilling', 'shipped', 'completed'),
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 6) DAY), INTERVAL 1 HOUR),
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 6) DAY), INTERVAL ((`n` * 5) % 18) HOUR)
FROM `os_seq`;

INSERT INTO `qixi_crm_b_order_item` (
  `id`,`order_id`,`product_id`,`merchant_sku_id`,`sku_key`,`title_snapshot`,`cover_url_snapshot`,
  `spec_snapshot`,`unit_price`,`quantity`,`refund_quantity`
) WITH RECURSIVE `os_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `os_seq` WHERE `n` < 84
)
SELECT
  962000 + `n`,
  962000 + `n`,
  ELT(((`n` - 1) % 3) + 1, 1001, 1101, 1201),
  0,
  CONCAT('ostat-item-', `n`),
  ELT(((`n` - 1) % 3) + 1, '轻奢羊绒针织衫','无火藤条香氛礼盒','智能数显保温杯'),
  'https://cos.qxkejiwl.top/pte-live-ecrm/platform/20260807/7dc13a394086786f9aba4a9606ad1eb2.png',
  JSON_OBJECT('默认','标准'),
  80 + (`n` % 11) * 25,
  1 + (`n` % 2),
  0
FROM `os_seq`;

-- 发货方式：express/city/service/pickup；末批约 14 单故意不写发货 → 自动发货。
INSERT INTO `qixi_crm_b_order_delivery` (
  `id`,`order_id`,`delivery_type`,`carrier_code`,`tracking_no`,`status`,`delivered_at`
) WITH RECURSIVE `os_seq` AS (
  SELECT 1 AS `n` UNION ALL SELECT `n` + 1 FROM `os_seq` WHERE `n` < 70
)
SELECT
  962000 + `n`,
  962000 + `n`,
  ELT(((`n` - 1) % 4) + 1, 'express', 'city', 'service', 'pickup'),
  IF(((`n` - 1) % 4) = 0, 'SF', NULL),
  IF(((`n` - 1) % 4) = 0, CONCAT('SF-OSTAT-', LPAD(`n`, 4, '0')), NULL),
  'shipped',
  DATE_ADD(DATE_SUB(CURDATE(), INTERVAL ((`n` - 1) DIV 6) DAY), INTERVAL 6 HOUR)
FROM `os_seq`;

-- 为近 14 天其它无发货记录的已支付订单补齐发货方式，避免饼图被「自动发货」淹没。
DELETE FROM `qixi_crm_b_order_delivery` WHERE `id` BETWEEN 963001 AND 963400;
INSERT INTO `qixi_crm_b_order_delivery` (
  `id`,`order_id`,`delivery_type`,`carrier_code`,`tracking_no`,`status`,`delivered_at`
)
SELECT
  963000 + ROW_NUMBER() OVER (ORDER BY o.id),
  o.id,
  ELT(((o.id - 1) % 4) + 1, 'express', 'city', 'service', 'pickup'),
  IF(((o.id - 1) % 4) = 0, 'YTO', NULL),
  IF(((o.id - 1) % 4) = 0, CONCAT('YTO-OSTAT-', LPAD(o.id % 10000, 4, '0')), NULL),
  'shipped',
  DATE_ADD(o.created_at, INTERVAL 4 HOUR)
FROM `qixi_crm_b_order` AS o
LEFT JOIN `qixi_crm_b_order_delivery` AS d ON d.order_id = o.id
WHERE d.id IS NULL
  AND o.status IN ('paid','fulfilling','shipped','completed')
  AND o.created_at >= DATE_SUB(CURDATE(), INTERVAL 14 DAY)
  AND o.id NOT BETWEEN 962071 AND 962084
LIMIT 360;

-- 退款：近 14 天混合状态（退款金额仅计 refunded；退款订单数排除 cancelled）。
INSERT INTO `qixi_crm_b_refund` (
  `id`,`order_id`,`refund_no`,`reason`,`amount`,`refund_type`,`order_status_before`,
  `status`,`idempotency_key`,`created_at`
) VALUES
  (962001,962008,'OSTAT-R-0001','订单统计演示退款',68.00,'money_only','paid','refunded','ostat-refund-1',DATE_SUB(NOW(), INTERVAL 1 DAY)),
  (962002,962012,'OSTAT-R-0002','订单统计演示退款',95.00,'money_only','shipped','refunded','ostat-refund-2',DATE_SUB(NOW(), INTERVAL 2 DAY)),
  (962003,962016,'OSTAT-R-0003','订单统计演示退款',120.00,'return_and_refund','completed','applied','ostat-refund-3',DATE_SUB(NOW(), INTERVAL 3 DAY)),
  (962004,962020,'OSTAT-R-0004','订单统计演示退款',55.00,'money_only','paid','refunded','ostat-refund-4',DATE_SUB(NOW(), INTERVAL 8 DAY)),
  (962005,962024,'OSTAT-R-0005','订单统计演示退款',88.00,'money_only','paid','cancelled','ostat-refund-5',DATE_SUB(NOW(), INTERVAL 9 DAY)),
  (962006,962028,'OSTAT-R-0006','订单统计演示退款',140.00,'money_only','fulfilling','refunded','ostat-refund-6',DATE_SUB(NOW(), INTERVAL 10 DAY)),
  (962007,962032,'OSTAT-R-0007','订单统计演示退款',72.00,'money_only','paid','merchant_handling','ostat-refund-7',DATE_SUB(NOW(), INTERVAL 4 DAY)),
  (962008,962036,'OSTAT-R-0008','订单统计演示退款',110.00,'money_only','shipped','refunded','ostat-refund-8',DATE_SUB(NOW(), INTERVAL 5 DAY))
ON DUPLICATE KEY UPDATE
  `reason`=VALUES(`reason`),`amount`=VALUES(`amount`),`status`=VALUES(`status`),`created_at`=VALUES(`created_at`);
