-- 本地演示：积分订单（待发货 / 待收货 / 已完成 / 用户已删除）
-- utf8mb4；幂等 UPSERT。不含真实手机号。
SET NAMES utf8mb4;

-- 清理本夹具旧行（先子后父）
DELETE FROM `qixi_crm_b_order_delivery` WHERE `order_id` BETWEEN 988860101 AND 988860104;
DELETE FROM `qixi_crm_b_order_item` WHERE `order_id` BETWEEN 988860101 AND 988860104;
DELETE FROM `qixi_crm_b_order` WHERE `id` BETWEEN 988860101 AND 988860104;
DELETE FROM `qixi_crm_b_group_order` WHERE `id` BETWEEN 988860101 AND 988860104;

INSERT INTO `qixi_crm_b_group_order` (
  `id`,`order_no`,`user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,
  `total_quantity`,`activity_type`,`points_amount`,`recipient_snapshot`,`pay_channel`,`pay_status`,
  `paid_at`,`idempotency_key`,`remark`,`user_archived_at`,`created_at`
) VALUES
  (988860101,'PTS-DEMO-G-20260810-001',9101,0.00,0.00,0.00,0.00,1,20,120,
   JSON_OBJECT('recipient','演示收件人小满','mobile','13800001001','province','上海市','city','上海市','district','浦东新区','detail','积分兑换演示地址一号'),
   'balance','paid',DATE_SUB(NOW(),INTERVAL 2 DAY),'fixture-pts-order-988860101','请工作日发货',NULL,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (988860102,'PTS-DEMO-G-20260810-002',9101,9.90,0.00,0.00,9.90,1,20,180,
   JSON_OBJECT('recipient','演示收件人小雪','mobile','13800001002','province','浙江省','city','杭州市','district','西湖区','detail','积分兑换演示地址二号'),
   'wechat','paid',DATE_SUB(NOW(),INTERVAL 1 DAY),'fixture-pts-order-988860102','',NULL,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (988860103,'PTS-DEMO-G-20260810-003',9101,0.00,0.00,0.00,0.00,2,20,500,
   JSON_OBJECT('recipient','演示收件人阿强','mobile','13800001003','province','江苏省','city','苏州市','district','姑苏区','detail','积分兑换演示地址三号'),
   'balance','paid',DATE_SUB(NOW(),INTERVAL 5 DAY),'fixture-pts-order-988860103','包装精美一点',NULL,DATE_SUB(NOW(),INTERVAL 5 DAY)),
  (988860104,'PTS-DEMO-G-20260810-004',9101,0.00,0.00,0.00,0.00,1,20,220,
   JSON_OBJECT('recipient','演示收件人已删','mobile','13800001004','province','广东省','city','深圳市','district','南山区','detail','积分兑换演示地址四号'),
   'balance','paid',DATE_SUB(NOW(),INTERVAL 8 DAY),'fixture-pts-order-988860104','用户侧已删除订单',DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_SUB(NOW(),INTERVAL 8 DAY))
ON DUPLICATE KEY UPDATE
  `pay_status`=VALUES(`pay_status`),`points_amount`=VALUES(`points_amount`),`pay_amount`=VALUES(`pay_amount`),
  `remark`=VALUES(`remark`),`user_archived_at`=VALUES(`user_archived_at`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_order` (
  `id`,`group_order_id`,`order_no`,`merchant_id`,`merchant_name_snapshot`,`store_id`,`store_name_snapshot`,
  `user_id`,`total_amount`,`discount_amount`,`freight_amount`,`pay_amount`,`total_quantity`,`activity_type`,
  `points_amount`,`recipient_snapshot`,`remark`,`merchant_remark`,`is_system_del`,`status`,`paid_at`,`created_at`
) VALUES
  (988860101,988860101,'PTS-DEMO-O-20260810-001',1,'CRM Live服饰商户',1,'CRM Live服饰旗舰店',9101,
   0.00,0.00,0.00,0.00,1,20,120,
   JSON_OBJECT('recipient','演示收件人小满','mobile','13800001001','province','上海市','city','上海市','district','浦东新区','detail','积分兑换演示地址一号'),
   '请工作日发货','优先发出真丝方巾礼盒',0,'paid',DATE_SUB(NOW(),INTERVAL 2 DAY),DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (988860102,988860102,'PTS-DEMO-O-20260810-002',2,'CRM Live居家商户',2,'CRM Live居家优选店',9101,
   9.90,0.00,0.00,9.90,1,20,180,
   JSON_OBJECT('recipient','演示收件人小雪','mobile','13800001002','province','浙江省','city','杭州市','district','西湖区','detail','积分兑换演示地址二号'),
   '','已联系快递揽收',0,'shipped',DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (988860103,988860103,'PTS-DEMO-O-20260810-003',1,'CRM Live服饰商户',1,'CRM Live服饰旗舰店',9101,
   0.00,0.00,0.00,0.00,2,20,500,
   JSON_OBJECT('recipient','演示收件人阿强','mobile','13800001003','province','江苏省','city','苏州市','district','姑苏区','detail','积分兑换演示地址三号'),
   '包装精美一点','签收无误',0,'completed',DATE_SUB(NOW(),INTERVAL 5 DAY),DATE_SUB(NOW(),INTERVAL 5 DAY)),
  (988860104,988860104,'PTS-DEMO-O-20260810-004',3,'CRM Live数码商户',3,'CRM Live数码生活店',9101,
   0.00,0.00,0.00,0.00,1,20,220,
   JSON_OBJECT('recipient','演示收件人已删','mobile','13800001004','province','广东省','city','深圳市','district','南山区','detail','积分兑换演示地址四号'),
   '用户侧已删除订单','',0,'fulfilling',DATE_SUB(NOW(),INTERVAL 8 DAY),DATE_SUB(NOW(),INTERVAL 8 DAY))
ON DUPLICATE KEY UPDATE
  `status`=VALUES(`status`),`points_amount`=VALUES(`points_amount`),`pay_amount`=VALUES(`pay_amount`),
  `remark`=VALUES(`remark`),`merchant_remark`=VALUES(`merchant_remark`),`is_system_del`=VALUES(`is_system_del`),
  `paid_at`=VALUES(`paid_at`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_order_item` (
  `id`,`order_id`,`product_id`,`merchant_sku_id`,`sku_key`,`title_snapshot`,`cover_url_snapshot`,
  `spec_snapshot`,`unit_price`,`quantity`,`refund_quantity`
) VALUES
  (988860101,988860101,1005,0,'pts-1005','真丝印花方巾礼盒','/demo/product-scarf-v1.png',
   JSON_OBJECT('规格','礼盒装'),0.00,1,0),
  (988860102,988860102,1105,0,'pts-1105','真丝睡眠眼罩方巾组','/demo/product-scarf-v1.png',
   JSON_OBJECT('规格','标准款'),9.90,1,0),
  (988860103,988860103,1401,0,'pts-1401','轻奢羊绒围巾兑换券','/demo/product-knit-v1.png',
   JSON_OBJECT('规格','兑换券'),0.00,2,0),
  (988860104,988860104,1204,0,'pts-1204','便携保温杯清洁套装','/demo/product-tumbler-v1.png',
   JSON_OBJECT('规格','清洁套装'),0.00,1,0)
ON DUPLICATE KEY UPDATE
  `title_snapshot`=VALUES(`title_snapshot`),`cover_url_snapshot`=VALUES(`cover_url_snapshot`),
  `spec_snapshot`=VALUES(`spec_snapshot`),`unit_price`=VALUES(`unit_price`),`quantity`=VALUES(`quantity`);

INSERT INTO `qixi_crm_b_order_delivery` (
  `id`,`order_id`,`delivery_type`,`carrier_code`,`tracking_no`,`status`,`delivered_at`
) VALUES
  (988860102,988860102,'express','顺丰速运','SFPTS20260810002','shipped',NULL),
  (988860103,988860103,'express','中通快递','ZTOPTS20260810003','shipped',DATE_SUB(NOW(),INTERVAL 3 DAY))
ON DUPLICATE KEY UPDATE
  `carrier_code`=VALUES(`carrier_code`),`tracking_no`=VALUES(`tracking_no`),
  `status`=VALUES(`status`),`delivered_at`=VALUES(`delivered_at`);
