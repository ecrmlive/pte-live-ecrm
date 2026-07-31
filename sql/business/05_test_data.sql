USE `qixi_crm_business`;
SET NAMES utf8mb4;
-- 仅提供无个人信息的消费读模型夹具。生产数据必须经商户事件同步写入，不能依赖本文件。
INSERT INTO `qixi_crm_b_store_view` (`store_id`,`merchant_id`,`store_app_id`,`store_name`,`status`) VALUES
  (1,1,'qixi.store.demo.1','七禧服饰旗舰店',1),
  (2,2,'qixi.store.demo.2','七禧居家优选店',1),
  (3,3,'qixi.store.demo.3','七禧数码生活店',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_app_id`=VALUES(`store_app_id`),`store_name`=VALUES(`store_name`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_category_view` (`category_id`,`parent_id`,`name`,`sort`,`status`) VALUES
  (101,0,'服饰鞋包',10,1),(102,0,'家居生活',20,1),(103,0,'数码家电',30,1),
  (104,0,'美妆个护',40,1),(105,0,'食品生鲜',50,1),(106,0,'运动户外',60,1),
  (10101,101,'女装精选',11,1),(10102,101,'箱包配饰',12,1),(10201,102,'香氛家居',21,1),
  (10301,103,'数码配件',31,1),(10401,104,'护肤洗护',41,1),(10601,106,'跑步训练',61,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1001,1,1,'七禧服饰商户','七禧服饰旗舰店',10101,'轻奢羊绒针织衫','/demo/product-knit-v1.png',299.00,399.00,0,158,60,1,1,NOW()),
  (1002,1,1,'七禧服饰商户','七禧服饰旗舰店',10102,'头层牛皮通勤托特包','/demo/product-bag-v1.png',469.00,599.00,0,126,32,1,1,NOW()),
  (1003,1,1,'七禧服饰商户','七禧服饰旗舰店',10601,'轻量缓震跑步鞋','/demo/product-shoes-v1.png',369.00,459.00,0,97,48,1,1,NOW()),
  (1004,1,1,'七禧服饰商户','七禧服饰旗舰店',10101,'精纺圆领羊毛开衫','/demo/product-knit-v1.png',329.00,429.00,0,141,36,1,1,NOW()),
  (1005,1,1,'七禧服饰商户','七禧服饰旗舰店',10102,'真丝印花方巾礼盒','/demo/product-scarf-v1.png',129.00,169.00,0,132,90,1,1,NOW()),
  (1006,1,1,'七禧服饰商户','七禧服饰旗舰店',10102,'都市简约手提斜挎包','/demo/product-bag-v1.png',399.00,529.00,0,88,27,1,1,NOW()),
  (1007,1,1,'七禧服饰商户','七禧服饰旗舰店',10101,'柔软亲肤针织披肩','/demo/product-knit-v1.png',189.00,249.00,0,76,54,1,1,NOW()),
  (1008,1,1,'七禧服饰商户','七禧服饰旗舰店',10601,'城市通勤训练跑鞋','/demo/product-shoes-v1.png',429.00,529.00,0,64,31,1,1,NOW()),
  (1101,2,2,'七禧居家商户','七禧居家优选店',10201,'无火藤条香氛礼盒','/demo/product-fragrance-v1.png',239.00,299.00,0,186,72,1,1,NOW()),
  (1102,2,2,'七禧居家商户','七禧居家优选店',10201,'晚安助眠香薰蜡烛','/demo/product-fragrance-v1.png',139.00,189.00,0,119,66,1,1,NOW()),
  (1103,2,2,'七禧居家商户','七禧居家优选店',10301,'恒温随行保温杯','/demo/product-tumbler-v1.png',159.00,219.00,0,154,80,1,1,NOW()),
  (1104,2,2,'七禧居家商户','七禧居家优选店',10201,'晨间居家香氛套装','/demo/product-fragrance-v1.png',268.00,338.00,0,72,39,1,1,NOW()),
  (1105,2,2,'七禧居家商户','七禧居家优选店',10201,'真丝睡眠眼罩方巾组','/demo/product-scarf-v1.png',99.00,139.00,0,98,88,1,1,NOW()),
  (1106,2,2,'七禧居家商户','七禧居家优选店',10301,'轻量随行运动水杯','/demo/product-tumbler-v1.png',119.00,159.00,0,104,71,1,1,NOW()),
  (1107,2,2,'七禧居家商户','七禧居家优选店',10201,'客厅氛围香薰礼盒','/demo/product-fragrance-v1.png',299.00,369.00,0,57,26,1,1,NOW()),
  (1108,2,2,'七禧居家商户','七禧居家优选店',10201,'织物护理香氛喷雾','/demo/product-fragrance-v1.png',89.00,119.00,0,92,103,1,1,NOW()),
  (1201,3,3,'七禧数码商户','七禧数码生活店',10301,'智能数显保温杯','/demo/product-tumbler-v1.png',199.00,259.00,0,203,110,1,1,NOW()),
  (1202,3,3,'七禧数码商户','七禧数码生活店',10301,'通勤随行杯套组合','/demo/product-tumbler-v1.png',89.00,119.00,0,114,95,1,1,NOW()),
  (1203,3,3,'七禧数码商户','七禧数码生活店',10601,'轻量日常跑步鞋','/demo/product-shoes-v1.png',359.00,449.00,0,83,42,1,1,NOW()),
  (1204,3,3,'七禧数码商户','七禧数码生活店',10301,'便携保温杯清洁套装','/demo/product-tumbler-v1.png',129.00,179.00,0,68,59,1,1,NOW()),
  (1205,3,3,'七禧数码商户','七禧数码生活店',10301,'户外运动补水杯','/demo/product-tumbler-v1.png',149.00,199.00,0,77,64,1,1,NOW()),
  (1206,3,3,'七禧数码商户','七禧数码生活店',10601,'轻缓震训练跑鞋','/demo/product-shoes-v1.png',389.00,489.00,0,70,35,1,1,NOW()),
  (1207,3,3,'七禧数码商户','七禧数码生活店',10301,'桌面恒温杯垫礼盒','/demo/product-tumbler-v1.png',219.00,279.00,0,61,38,1,1,NOW()),
  (1208,3,3,'七禧数码商户','七禧数码生活店',10301,'轻巧旅行随行杯','/demo/product-tumbler-v1.png',109.00,149.00,0,90,85,1,1,NOW())
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`merchant_name`=VALUES(`merchant_name`),`store_name`=VALUES(`store_name`),`category_id`=VALUES(`category_id`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`price`=VALUES(`price`),`original_price`=VALUES(`original_price`),`sales`=VALUES(`sales`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`updated_at`=NOW();

-- 领券中心夹具：store_id=0 为平台券，其余为对应店铺券。领取记录只在用户实际点击领取后生成。
INSERT INTO `qixi_crm_b_coupon_template_view` (`coupon_id`,`store_id`,`name`,`discount_type`,`discount_value`,`min_amount`,`starts_at`,`ends_at`,`status`,`version`) VALUES
  (3001,0,'平台新客满99减10','amount',10.00,99.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3002,0,'平台夏日满299减40','amount',40.00,299.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3003,1,'服饰店满199减30','amount',30.00,199.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3004,1,'服饰店9折券','rate',90.00,399.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3005,2,'居家优选满159减20','amount',20.00,159.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3006,3,'数码生活满299减35','amount',35.00,299.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1)
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`name`=VALUES(`name`),`discount_type`=VALUES(`discount_type`),`discount_value`=VALUES(`discount_value`),`min_amount`=VALUES(`min_amount`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`),`status`=VALUES(`status`),`version`=VALUES(`version`);

-- 秒杀展示夹具。规则由后台营销活动投影而来；C 端只读 qixi_crm_b_marketing_activity_view。
INSERT INTO `qixi_crm_b_marketing_activity_view` (`activity_id`,`store_id`,`activity_type`,`name`,`rules`,`status`,`version`,`starts_at`,`ends_at`) VALUES
  (5001,1,'seckill','轻奢羊绒针织衫限时抢购',JSON_OBJECT('product_id',1001,'seckill_price',199.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5002,1,'seckill','头层牛皮托特包限时抢购',JSON_OBJECT('product_id',1002,'seckill_price',329.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5003,2,'seckill','无火藤条香氛礼盒限时抢购',JSON_OBJECT('product_id',1101,'seckill_price',169.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5004,3,'seckill','智能数显保温杯限时抢购',JSON_OBJECT('product_id',1201,'seckill_price',149.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY))
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`name`=VALUES(`name`),`rules`=VALUES(`rules`),`status`=VALUES(`status`),`version`=VALUES(`version`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`);

INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2001,'notice','七禧商城服务公告','/demo/home-hero-v1.png','七禧商城已上线商品、订单、售后和客服服务。消费者可通过商品详情、购物车和订单中心完成全流程购物。',1,1,NOW(),NOW()),
  (2002,'notice','消费者权益说明','/demo/home-service-wide-v1.png','请在下单前确认商品信息、配送方式和售后规则。如有商品与履约问题，可在订单中心提交售后申请。',1,1,NOW(),NOW()),
  (2003,'notice','夏日居家焕新季：精选家居好物上新','/demo/home-tech-wide-v1.png','居家生活专区已上新香氛、随行杯与织物护理系列，支持按分类、销量和价格快速筛选。',1,1,NOW(),NOW()),
  (2004,'notice','七禧多商户店铺服务规范','/demo/home-service-vertical-v1.png','平台持续完善商户审核、商品审核、订单履约与售后处理规范，为消费者提供清晰可靠的购物体验。',1,1,NOW(),NOW()),
  (2005,'notice','会员积分与优惠券使用说明','/demo/home-beauty-vertical-v1.png','积分、优惠券将按各自规则展示和使用。结算页会提示可用权益与优惠金额。',1,1,NOW(),NOW()),
  (2101,'agreement','sys_user_agree','','欢迎使用七禧商城。使用服务前请阅读并同意本用户协议。',1,1,NOW(),NOW()),
  (2102,'agreement','sys_userr_privacy','','七禧仅在提供服务所必需的范围内处理您的个人信息。',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`body`=VALUES(`body`),`status`=VALUES(`status`),`version`=VALUES(`version`),`published_at`=VALUES(`published_at`),`updated_at`=NOW();

-- 本地投影夹具；生产环境只允许 api-merchant 的 outbox/NATS 事件写入。
INSERT INTO `qixi_crm_b_diy_page_view` (`source`,`page_id`,`store_id`,`page_type`,`name`,`document`,`status`,`is_active`) VALUES
  ('merchant',3001,1,'home','七禧演示店铺首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','七禧演示店铺','title','七禧演示店铺')),
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(JSON_OBJECT('imgName','七禧演示店铺','imgUrl','','linkUrl','/pages/store/index'))),
      JSON_OBJECT('type','navBar','name','导航组','data',JSON_ARRAY(JSON_OBJECT('text','全部商品','imgUrl','','linkUrl','/pages/goods/list'),JSON_OBJECT('text','购物车','imgUrl','','linkUrl','/pages/order_addcart/order_addcart')))
    ),
    '_qixi',JSON_OBJECT('title','七禧演示店铺','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`),`is_active`=VALUES(`is_active`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_diy_page_view` (`source`,`page_id`,`store_id`,`page_type`,`name`,`document`,`status`,`is_active`) VALUES
  ('platform',4001,0,'home','七禧平台首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','七禧商城','title','七禧商城')),
    'items',JSON_ARRAY(JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(
      JSON_OBJECT('imgName','七禧商城精选','imgUrl','/demo/home-hero-v1.png','linkUrl','/goods?cate_id=101'),
      JSON_OBJECT('imgName','七禧香氛家居','imgUrl','/demo/home-hero-fragrance-v1.png','linkUrl','/goods?cate_id=102'),
      JSON_OBJECT('imgName','七禧箱包配饰','imgUrl','/demo/home-hero-accessories-v1.png','linkUrl','/goods?cate_id=10102')
    ))),
    '_qixi',JSON_OBJECT('title','七禧商城','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`),`is_active`=VALUES(`is_active`),`updated_at`=NOW();
