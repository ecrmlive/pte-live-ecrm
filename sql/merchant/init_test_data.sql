USE `qixi_crm_merchant`;
SET NAMES utf8mb4;
-- 不初始化店铺账号、密码或任何真实手机号。仅提供无个人信息的店铺读写模型夹具。

INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`brand_name`,`status`,`version`) VALUES
  (1001,1,'七禧臻选羊绒针织衫',101,'云锦织造','on_sale',1),
  (1002,1,'七禧头层牛皮通勤包',101,'栖木皮具','on_sale',1),
  (1003,1,'七禧轻量跑步鞋',106,'逐风运动','on_sale',1),
  (1004,1,'七禧智能保温杯',103,'澄日生活','on_sale',1),
  (1005,1,'七禧真丝方巾',104,'云锦织造','on_sale',1),
  (1006,1,'七禧香氛礼盒',102,'澄日生活','on_sale',1)
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`category_id`=VALUES(`category_id`),`brand_name`=VALUES(`brand_name`),`status`=VALUES(`status`),`version`=VALUES(`version`);

-- SVIP 价格夹具：商品审核后同步到业务消费视图；九折和固定专享价均不含真实会员资料。
UPDATE `qixi_crm_m_product` SET `svip_price_type`=1, `svip_price`=0 WHERE `id`=1001;
UPDATE `qixi_crm_m_product` SET `svip_price_type`=2, `svip_price`=429.00 WHERE `id`=1002;

-- SKU 主键固定，方便三库隔离的中文闭环夹具可靠验证；均为虚构商品。
INSERT INTO `qixi_crm_m_product_sku` (`id`,`product_id`,`spec_json`,`price`,`stock`,`status`) VALUES
  (61001,1001,JSON_OBJECT('默认','标准'),299.00,50,1),
  (61002,1002,JSON_OBJECT('默认','标准'),469.00,50,1),
  (61003,1003,JSON_OBJECT('颜色','晨雾灰','尺码','40'),369.00,24,1),
  (61007,1003,JSON_OBJECT('颜色','星曜蓝','尺码','41'),389.00,26,1),
  (61004,1004,JSON_OBJECT('默认','标准'),159.00,50,1),
  (61005,1005,JSON_OBJECT('默认','标准'),129.00,50,1),
  (61006,1006,JSON_OBJECT('默认','标准'),239.00,50,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`spec_json`=VALUES(`spec_json`),`price`=VALUES(`price`),`stock`=VALUES(`stock`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_finance_ledger` (`store_id`,`entry_type`,`amount`,`reference_type`,`reference_id`,`idempotency_key`) VALUES
  (1,'test_seed',1000.00,'fixture','opening-balance','fixture-store-1-opening-balance')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`);

-- 结算状态机夹具：仅供本地验收，均为虚构账期与金额，不含收款资料或支付凭据。
INSERT INTO `qixi_crm_m_settlement_bill`
  (`id`,`store_id`,`merchant_id`,`period_start`,`period_end`,`amount`,`status`,`version`)
VALUES
  (9001,1,1,'2026-07-01 00:00:00','2026-07-31 23:59:59',1280.50,'bill_frozen',1),
  (9002,2,2,'2026-06-01 00:00:00','2026-06-30 23:59:59',960.00,'paid',3)
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`status`=VALUES(`status`),`version`=VALUES(`version`);

-- 订单级结算账本中文夹具：均为虚构订单和退款编号，仅验证应计与负向调整可追溯，
-- 不含收款账户、真实手机号、支付流水或任何外部凭据。
INSERT INTO `qixi_crm_m_settlement_entry`
  (`store_id`,`merchant_id`,`order_id`,`refund_id`,`entry_type`,`amount`,`idempotency_key`,`occurred_at`)
VALUES
  (1,1,99001,0,'order_accrual',299.00,'settlement:accrue:99001','2026-07-18 10:00:00'),
  (1,1,99001,88001,'refund_reversal',-299.00,'settlement:reverse:88001','2026-08-02 09:00:00')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`occurred_at`=VALUES(`occurred_at`);

-- 优惠套餐夹具：店铺营销活动事实；不含收款或密钥。
INSERT INTO `qixi_crm_m_marketing_activity`
  (`id`,`store_id`,`activity_type`,`name`,`rules`,`status`,`starts_at`,`ends_at`)
VALUES
  (5101,1,'discount','夏日香氛随行套餐',JSON_OBJECT('package_price',199.00,'product_ids',JSON_ARRAY(1004,1006),'free_shipping',true,'remark','中文演示套餐'),'active',DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY))
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`rules`=VALUES(`rules`),`status`=VALUES(`status`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`);

-- 售后备注夹具：关联业务库的虚构退款 9900201，仅验证店铺操作审计；不含买家资料、支付凭据或真实客服账号。
INSERT INTO `qixi_crm_m_aftersale_action`
  (`refund_id`,`store_id`,`account_id`,`action`,`note`,`attachments`,`idempotency_key`)
VALUES
  (9900201,1,9901,'remark','虚构演示：仓库已核验退货包裹外观，等待确认收货。',JSON_ARRAY(),'fixture-refund-remark-9900201')
ON DUPLICATE KEY UPDATE `note`=VALUES(`note`),`attachments`=VALUES(`attachments`);

-- 仅用于本地联调：店铺发布页的真实消费来源仍是 qixi_crm_m_outbox 事件。
INSERT INTO `qixi_crm_m_diy_page` (`id`,`store_id`,`name`,`document`,`page_type`,`is_active`,`status`) VALUES
  (3001,1,'七禧演示店铺首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','七禧演示店铺','title','七禧演示店铺')),
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(JSON_OBJECT('imgName','七禧演示店铺','imgUrl','','linkUrl','/pages/store/index'))),
      JSON_OBJECT('type','navBar','name','导航组','data',JSON_ARRAY(JSON_OBJECT('text','全部商品','imgUrl','','linkUrl','/pages/goods/list'),JSON_OBJECT('text','购物车','imgUrl','','linkUrl','/pages/order_addcart/order_addcart')))
    ),
    '_qixi',JSON_OBJECT('title','七禧演示店铺','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'home',1,'published')
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`page_type`=VALUES(`page_type`),`is_active`=VALUES(`is_active`),`status`=VALUES(`status`);
