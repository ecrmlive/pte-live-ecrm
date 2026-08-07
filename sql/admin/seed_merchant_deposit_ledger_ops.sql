-- 店铺保证金「操作记录」演示流水：同一商户多条 fund/deduct，对齐 CRMEB 弹窗验收
SET NAMES utf8mb4;
USE qixi_crm_admin;

DELETE FROM qixi_crm_a_merchant_deposit_ledger WHERE id BETWEEN 961001 AND 961100;

-- 商户 1：补缴 → 扣除 → 再扣除（结余 0），对应截图风格
INSERT INTO qixi_crm_a_merchant_deposit_ledger
 (`id`,`merchant_id`,`entry_type`,`amount`,`balance_after`,`reason`,`idempotency_key`,`operator_admin_id`,`created_at`)
VALUES
(961001,1,'fund',1000.00,1000.00,'线下补缴保证金','demo-ops-1-fund',1,'2026-03-20 09:48:03'),
(961002,1,'deduct',1000.00,0.00,'111','demo-ops-1-deduct-1',1,'2026-04-29 16:45:24'),
(961003,1,'deduct',0.00,0.00,'11111','demo-ops-1-deduct-2',1,'2026-04-29 16:45:55'),

-- 商户 1004：补缴 → 扣除部分 → 再补缴
(961011,1004,'fund',1000.00,1000.00,'线下补缴保证金','demo-ops-1004-fund-1',1,'2026-03-18 10:00:00'),
(961012,1004,'deduct',200.00,800.00,'虚构违规扣减','demo-ops-1004-deduct-1',1,'2026-04-01 14:20:00'),
(961013,1004,'fund',200.00,1000.00,'线下补缴保证金','demo-ops-1004-fund-2',1,'2026-04-15 11:30:00'),
(961014,1004,'deduct',700.00,300.00,'本地验收扣减','demo-ops-1004-deduct-2',1,'2026-07-20 16:00:00'),

-- 商户 1100：完整路径 补缴→扣除→退还
(961021,1100,'fund',500.00,500.00,'线下补缴保证金','demo-ops-1100-fund',1,'2026-02-10 09:00:00'),
(961022,1100,'deduct',100.00,400.00,'虚构违规扣减','demo-ops-1100-deduct',1,'2026-05-01 15:00:00'),
(961023,1100,'refund_approved',0.00,400.00,'退款审核通过','demo-ops-1100-refund-ok',1,'2026-08-03 09:30:00'),
(961024,1100,'refund_paid',400.00,0.00,'退还保证金已打款','demo-ops-1100-refund-paid',1,'2026-08-03 10:03:00'),

-- 商户 2：补缴 + 拒绝退款记录
(961031,2,'fund',800.00,800.00,'线下补缴保证金','demo-ops-2-fund',1,'2026-01-12 08:30:00'),
(961032,2,'deduct',50.00,750.00,'虚构违规扣减','demo-ops-2-deduct',1,'2026-06-08 13:10:00'),
(961033,2,'refund_rejected',0.00,750.00,'退款申请驳回','demo-ops-2-refund-rej',1,'2026-08-01 09:05:00');

SELECT merchant_id, COUNT(*) c
FROM qixi_crm_a_merchant_deposit_ledger
WHERE id BETWEEN 961001 AND 961100
GROUP BY merchant_id
ORDER BY merchant_id;
