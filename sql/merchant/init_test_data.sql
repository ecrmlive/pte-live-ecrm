USE `qixi_crm_merchant`;
SET NAMES utf8mb4;
-- 不初始化店铺账号、密码或任何真实手机号。仅提供无个人信息的店铺读写模型夹具。

INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`brand_name`,`status`,`version`) VALUES
  (1001,1,'轻奢羊绒针织衫',101,'云锦织造','on_sale',1),
  (1002,1,'头层牛皮通勤托特包',101,'栖木皮具','on_sale',1),
  (1003,1,'七禧轻量跑步鞋',106,'逐风运动','on_sale',1),
  (1004,1,'精纺圆领羊毛开衫',101,'云锦织造','on_sale',1),
  (1005,1,'真丝印花方巾礼盒',104,'云锦织造','on_sale',1),
  (1006,1,'七禧香氛礼盒',102,'澄日生活','on_sale',1),
  (1007,1,'柔软亲肤针织披肩',101,'云锦织造','on_sale',1),
  (1008,1,'城市通勤训练跑鞋',106,'逐风运动','on_sale',1),
  (1101,2,'无火藤条香氛礼盒',102,'澄日生活','on_sale',1),
  (1102,2,'晚安助眠香薰蜡烛',102,'澄日生活','on_sale',1),
  (1103,2,'恒温随行保温杯',103,'CRM Live精选','on_sale',1),
  (1104,2,'晨间居家香氛套装',102,'澄日生活','on_sale',1),
  (1107,2,'客厅氛围香薰礼盒',102,'澄日生活','on_sale',1),
  (1108,2,'织物护理香氛喷雾',102,'澄日生活','on_sale',1),
  (1207,3,'桌面恒温杯垫礼盒',103,'CRM Live精选','on_sale',1)
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
  (61004,1004,JSON_OBJECT('默认','标准'),329.00,36,1),
  (61005,1005,JSON_OBJECT('默认','标准'),129.00,50,1),
  (61006,1006,JSON_OBJECT('默认','标准'),239.00,50,1),
  (61008,1007,JSON_OBJECT('默认','标准'),189.00,54,1),
  (61009,1008,JSON_OBJECT('默认','标准'),429.00,31,1),
  (61101,1101,JSON_OBJECT('默认','标准'),239.00,72,1),
  (61102,1102,JSON_OBJECT('默认','标准'),139.00,66,1),
  (61103,1103,JSON_OBJECT('默认','标准'),159.00,80,1),
  (61104,1104,JSON_OBJECT('默认','标准'),268.00,39,1),
  (61107,1107,JSON_OBJECT('默认','标准'),299.00,26,1),
  (61108,1108,JSON_OBJECT('默认','标准'),89.00,103,1),
  (61207,1207,JSON_OBJECT('默认','标准'),219.00,38,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`spec_json`=VALUES(`spec_json`),`price`=VALUES(`price`),`stock`=VALUES(`stock`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_product_detail`
  (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`delivery_way`,`original_price`)
VALUES
  (1001,'轻薄保暖，通勤百搭羊绒针织。','羊绒,针织,通勤','件','/demo/product-knit-v1.png','2',399.00),
  (1002,'头层牛皮托特，容量适中适合通勤。','牛皮,托特包,通勤','件','/demo/product-bag-v1.png','2',599.00),
  (1004,'精纺圆领羊毛开衫，秋冬预售演示商品。','羊毛,开衫,预售','件','/demo/product-knit-v1.png','1,2',429.00),
  (1005,'真丝印花方巾礼盒，好友助力演示商品。','真丝,方巾,礼盒','盒','/demo/product-scarf-v1.png','2',169.00),
  (1007,'柔软亲肤针织披肩，居家与出行两用。','披肩,针织','件','/demo/product-knit-v1.png','2',249.00),
  (1008,'城市通勤训练跑鞋，缓震透气。','跑鞋,通勤,运动','双','/demo/product-shoes-v1.png','2',529.00),
  (1101,'无火藤条香氛礼盒，客厅氛围演示商品。','香氛,藤条,礼盒','盒','/demo/product-fragrance-v1.png','2',299.00),
  (1102,'晚安助眠香薰蜡烛，助眠香调演示。','香薰,蜡烛,助眠','件','/demo/product-fragrance-v1.png','2',189.00),
  (1103,'恒温随行保温杯，出行随手一杯。','保温杯,随行','个','/demo/product-tumbler-v1.png','2',219.00),
  (1104,'晨间居家香氛套装，全款预售演示。','香氛,套装,居家','套','/demo/product-fragrance-v1.png','2',338.00),
  (1107,'客厅氛围香薰礼盒，定金预售演示。','香薰,礼盒,客厅','盒','/demo/product-fragrance-v1.png','2',369.00),
  (1108,'织物护理香氛喷雾，定金待审演示商品。','织物,护理,香氛喷雾','瓶','/demo/product-fragrance-v1.png','1,2',119.00),
  (1207,'桌面恒温杯垫礼盒，好友助力演示商品。','恒温,杯垫,礼盒','盒','/demo/product-tumbler-v1.png','2',279.00)
ON DUPLICATE KEY UPDATE
  `brief`=VALUES(`brief`),`keyword`=VALUES(`keyword`),`unit_name`=VALUES(`unit_name`),
  `cover_url`=VALUES(`cover_url`),`delivery_way`=VALUES(`delivery_way`),`original_price`=VALUES(`original_price`);

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
