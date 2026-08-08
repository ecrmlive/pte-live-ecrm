-- 平台「核销记录」本地演示：为已存在的到店自提已完成订单补齐 used 核销事实。
-- 幂等：固定 id 区间；可重复执行。utf8mb4 中文。
SET NAMES utf8mb4;

-- 修正演示店员展示名（避免乱码占位）
UPDATE `qixi_crm_merchant`.`qixi_crm_m_account`
SET `display_name` = CASE `id`
  WHEN 1 THEN '演示店主甲'
  WHEN 2 THEN '演示店主乙'
  WHEN 3 THEN '演示店主丙'
  ELSE `display_name`
END
WHERE `id` IN (1, 2, 3);

-- 为核销演示订单补齐可读收货人 / 商品名，并按序号打散支付渠道（自提单 id 对 4 同余，不能用 id%4）
UPDATE `qixi_crm_b_order` o
JOIN `qixi_crm_b_group_order` g ON g.id = o.group_order_id
JOIN (
  SELECT d.order_id,
         ROW_NUMBER() OVER (ORDER BY d.order_id) AS rn
  FROM `qixi_crm_b_order_delivery` d
  JOIN `qixi_crm_b_order` ox ON ox.id = d.order_id
  WHERE d.delivery_type = 'pickup'
    AND ox.status IN ('paid', 'fulfilling', 'completed', 'shipped')
  ORDER BY d.order_id
  LIMIT 24
) t ON t.order_id = o.id
SET
  o.recipient_snapshot = JSON_OBJECT(
    'recipient', ELT(((t.rn - 1) % 5) + 1, '小杨', '阿强', '小林', '陈晨', '王芳'),
    'mobile', CONCAT('139', LPAD((o.id % 100000000), 8, '0')),
    'province', '广东省', 'city', '广州市', 'district', '天河区',
    'detail', CONCAT('演示路', (o.id % 90) + 1, '号')
  ),
  o.status = 'completed',
  g.pay_channel = ELT(((t.rn - 1) % 3) + 1, 'balance', 'wechat', 'alipay');

UPDATE `qixi_crm_b_order_item` oi
JOIN (
  SELECT d.order_id
  FROM `qixi_crm_b_order_delivery` d
  JOIN `qixi_crm_b_order` ox ON ox.id = d.order_id
  WHERE d.delivery_type = 'pickup'
    AND ox.status = 'completed'
  ORDER BY d.order_id
  LIMIT 24
) t ON t.order_id = oi.order_id
SET
  oi.title_snapshot = ELT(((oi.order_id - 1) % 4) + 1,
    'square houlest运动五分裤 紧身短裤',
    '演示自提咖啡豆 454g',
    '门店核销体验卡 次卡',
    '本地生活到店套餐券'),
  oi.cover_url_snapshot = IF(oi.cover_url_snapshot = '' OR oi.cover_url_snapshot IS NULL
      OR oi.cover_url_snapshot LIKE '/demo/%'
      OR oi.cover_url_snapshot LIKE 'demo/%',
    'https://cos.qxkejiwl.top/pte-live-ecrm/platform/20260807/7dc13a394086786f9aba4a9606ad1eb2.png',
    oi.cover_url_snapshot),
  oi.spec_snapshot = JSON_OBJECT('规格', ELT(((oi.order_id - 1) % 3) + 1, '紧身短裤', '标准装', '体验版'));

DELETE FROM `qixi_crm_b_order_verification` WHERE `id` BETWEEN 962401 AND 962424;

INSERT INTO `qixi_crm_b_order_verification` (
  `id`, `order_id`, `verify_code`, `verify_code_hash`, `status`,
  `verified_by_account_id`, `verified_at`, `created_at`
)
SELECT
  962400 + ROW_NUMBER() OVER (ORDER BY d.order_id),
  d.order_id,
  '',
  SHA2(CONCAT('verify-demo-', d.order_id), 256),
  'used',
  CASE
    WHEN (d.order_id % 5) = 0 THEN NULL
    WHEN (d.order_id % 5) = 1 THEN 9527
    ELSE ((d.order_id - 1) % 3) + 1
  END,
  DATE_ADD(COALESCE(ox.paid_at, ox.created_at), INTERVAL 2 HOUR),
  COALESCE(ox.paid_at, ox.created_at)
FROM `qixi_crm_b_order_delivery` d
JOIN `qixi_crm_b_order` ox ON ox.id = d.order_id
WHERE d.delivery_type = 'pickup'
  AND ox.status = 'completed'
ORDER BY d.order_id
LIMIT 24;
