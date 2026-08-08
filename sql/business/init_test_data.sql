USE `qixi_crm_business`;
SET NAMES utf8mb4;
-- 仅提供无个人信息的消费读模型夹具。生产数据必须经商户事件同步写入，不能依赖本文件。
INSERT INTO `qixi_crm_b_store_view` (`store_id`,`merchant_id`,`store_app_id`,`store_name`,`status`,`integral_enabled`,`integral_points_per_yuan`,`integral_max_deduction_bps`) VALUES
  (1,1,'qixi.store.demo.1','CRM Live服饰旗舰店',1,1,100,2000),
  (2,2,'qixi.store.demo.2','CRM Live居家优选店',1,0,100,2000),
  (3,3,'qixi.store.demo.3','CRM Live数码生活店',1,0,100,2000)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_app_id`=VALUES(`store_app_id`),`store_name`=VALUES(`store_name`),`status`=VALUES(`status`),`integral_enabled`=VALUES(`integral_enabled`),`integral_points_per_yuan`=VALUES(`integral_points_per_yuan`),`integral_max_deduction_bps`=VALUES(`integral_max_deduction_bps`);

INSERT INTO `qixi_crm_b_category_view` (`category_id`,`parent_id`,`name`,`sort`,`status`) VALUES
  (101,0,'服饰鞋包',10,1),(102,0,'家居生活',20,1),(103,0,'数码家电',30,1),
  (104,0,'美妆个护',40,1),(105,0,'食品生鲜',50,1),(106,0,'运动户外',60,1),
  (10101,101,'女装精选',11,1),(10102,101,'箱包配饰',12,1),(10201,102,'香氛家居',21,1),
  (10301,103,'数码配件',31,1),(10401,104,'护肤洗护',41,1),(10601,106,'跑步训练',61,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

-- 行政区划消费投影：全部为本地中文模拟数据，由平台配送配置发布；不包含定位坐标或个人地址。
INSERT INTO `qixi_crm_b_city_view` (`city_id`,`parent_id`,`name`,`level`,`is_show`) VALUES
  (1,0,'中国',1,1),
  (310000,1,'上海市',2,1),(310100,310000,'上海市区',3,1),
  (310101,310100,'黄浦区',4,1),(310104,310100,'徐汇区',4,1),(310105,310100,'长宁区',4,1),
  (310106,310100,'静安区',4,1),(310107,310100,'普陀区',4,1),(310115,310100,'浦东新区',4,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`level`=VALUES(`level`),`is_show`=VALUES(`is_show`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1001,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10101,'轻奢羊绒针织衫','/demo/product-knit-v1.png',299.00,399.00,0,158,60,1,1,NOW()),
  (1002,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10102,'头层牛皮通勤托特包','/demo/product-bag-v1.png',469.00,599.00,0,126,32,1,1,NOW()),
  (1003,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10601,'轻量缓震跑步鞋','/demo/product-shoes-v1.png',369.00,459.00,0,97,48,1,1,NOW()),
  (1004,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10101,'精纺圆领羊毛开衫','/demo/product-knit-v1.png',329.00,429.00,0,141,36,1,1,NOW()),
  (1005,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10102,'真丝印花方巾礼盒','/demo/product-scarf-v1.png',129.00,169.00,0,132,90,1,1,NOW()),
  (1006,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10102,'都市简约手提斜挎包','/demo/product-bag-v1.png',399.00,529.00,0,88,27,1,1,NOW()),
  (1007,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10101,'柔软亲肤针织披肩','/demo/product-knit-v1.png',189.00,249.00,0,76,54,1,1,NOW()),
  (1008,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10601,'城市通勤训练跑鞋','/demo/product-shoes-v1.png',429.00,529.00,0,64,31,1,1,NOW()),
  (1101,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'无火藤条香氛礼盒','/demo/product-fragrance-v1.png',239.00,299.00,0,186,72,1,1,NOW()),
  (1102,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'晚安助眠香薰蜡烛','/demo/product-fragrance-v1.png',139.00,189.00,0,119,66,1,1,NOW()),
  (1103,2,2,'CRM Live居家商户','CRM Live居家优选店',10301,'恒温随行保温杯','/demo/product-tumbler-v1.png',159.00,219.00,0,154,80,1,1,NOW()),
  (1104,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'晨间居家香氛套装','/demo/product-fragrance-v1.png',268.00,338.00,0,72,39,1,1,NOW()),
  (1105,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'真丝睡眠眼罩方巾组','/demo/product-scarf-v1.png',99.00,139.00,0,98,88,1,1,NOW()),
  (1106,2,2,'CRM Live居家商户','CRM Live居家优选店',10301,'轻量随行运动水杯','/demo/product-tumbler-v1.png',119.00,159.00,0,104,71,1,1,NOW()),
  (1107,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'客厅氛围香薰礼盒','/demo/product-fragrance-v1.png',299.00,369.00,0,57,26,1,1,NOW()),
  (1108,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'织物护理香氛喷雾','/demo/product-fragrance-v1.png',89.00,119.00,0,92,103,1,1,NOW()),
  (1201,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'智能数显保温杯','/demo/product-tumbler-v1.png',199.00,259.00,0,203,110,1,1,NOW()),
  (1202,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'通勤随行杯套组合','/demo/product-tumbler-v1.png',89.00,119.00,0,114,95,1,1,NOW()),
  (1203,3,3,'CRM Live数码商户','CRM Live数码生活店',10601,'轻量日常跑步鞋','/demo/product-shoes-v1.png',359.00,449.00,0,83,42,1,1,NOW()),
  (1204,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'便携保温杯清洁套装','/demo/product-tumbler-v1.png',129.00,179.00,0,68,59,1,1,NOW()),
  (1205,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'户外运动补水杯','/demo/product-tumbler-v1.png',149.00,199.00,0,77,64,1,1,NOW()),
  (1206,3,3,'CRM Live数码商户','CRM Live数码生活店',10601,'轻缓震训练跑鞋','/demo/product-shoes-v1.png',389.00,489.00,0,70,35,1,1,NOW()),
  (1207,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'桌面恒温杯垫礼盒','/demo/product-tumbler-v1.png',219.00,279.00,0,61,38,1,1,NOW()),
  (1208,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'轻巧旅行随行杯','/demo/product-tumbler-v1.png',109.00,149.00,0,90,85,1,1,NOW())
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`merchant_name`=VALUES(`merchant_name`),`store_name`=VALUES(`store_name`),`category_id`=VALUES(`category_id`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`price`=VALUES(`price`),`original_price`=VALUES(`original_price`),`sales`=VALUES(`sales`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`updated_at`=NOW();

-- SVIP 价格夹具：默认九折与固定专享价均由商品消费视图提供，未登录用户不会被降价。
UPDATE `qixi_crm_b_product_view` SET `svip_price_type`=1, `svip_price`=0 WHERE `product_id`=1001;
UPDATE `qixi_crm_b_product_view` SET `svip_price_type`=2, `svip_price`=429.00 WHERE `product_id`=1002;

-- 中文品牌夹具来自已审核商品投影；C 端只按品牌名精确筛选，不接受客户端写入。
UPDATE `qixi_crm_b_product_view` SET `brand_name` = CASE
  WHEN `product_id` IN (1001,1004,1007) THEN '云锦织造'
  WHEN `product_id` IN (1002,1006,1201) THEN '栖木皮具'
  WHEN `product_id` IN (1003,1008,1203,1206) THEN '逐风运动'
  WHEN `product_id` IN (1101,1102,1104,1107) THEN '澄日生活'
  ELSE 'CRM Live精选' END;

-- 订单/库存闭环的业务侧 SKU 消费投影。主键来自商户库的虚构固定 SKU，不包含
-- 商户连接信息、真实用户身份或凭据。
INSERT INTO `qixi_crm_b_product_sku_view` (`merchant_sku_id`,`product_id`,`sku_key`,`spec_snapshot`,`price`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (61001,1001,'61001',JSON_OBJECT('默认','标准'),299.00,50,1,1,NOW()),
  (61002,1002,'61002',JSON_OBJECT('默认','标准'),469.00,50,1,1,NOW()),
  (61003,1003,'61003',JSON_OBJECT('颜色','晨雾灰','尺码','40'),369.00,24,1,1,NOW()),
  (61007,1003,'61007',JSON_OBJECT('颜色','星曜蓝','尺码','41'),389.00,26,1,1,NOW()),
  (61004,1004,'61004',JSON_OBJECT('默认','标准'),159.00,50,1,1,NOW()),
  (61005,1005,'61005',JSON_OBJECT('默认','标准'),129.00,50,1,1,NOW()),
  (61006,1006,'61006',JSON_OBJECT('默认','标准'),239.00,50,1,1,NOW())
ON DUPLICATE KEY UPDATE `sku_key`=VALUES(`sku_key`),`spec_snapshot`=VALUES(`spec_snapshot`),`price`=VALUES(`price`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`version`=VALUES(`version`),`updated_at`=NOW();

-- 直播监管夹具仅包含虚构房间、商品和审核状态；推流地址、主播手机号始终为空，不写入任何凭据。
INSERT INTO `qixi_crm_b_broadcast_room` (`broadcast_room_id`,`mer_id`,`name`,`cover_img`,`feeds_img`,`play_url`,`push_url`,`start_time`,`end_time`,`anchor_name`,`phone`,`status`,`live_status`,`is_show`,`is_del`,`sort`,`star`,`mark`,`refusal`) VALUES
  (7101,1,'CRM Live服饰秋日穿搭直播间','/demo/live-fashion-cover.png','/demo/live-fashion-feed.png','','',DATE_ADD(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 26 HOUR),'虚构主播小七','',0,102,0,0,10,5,'中文模拟直播审核夹具',''),
  (7102,2,'CRM Live居家香氛新品直播间','/demo/live-home-cover.png','/demo/live-home-feed.png','','',DATE_SUB(NOW(),INTERVAL 1 HOUR),DATE_ADD(NOW(),INTERVAL 1 HOUR),'虚构主播小居','',2,101,1,0,8,4,'中文模拟已审核直播夹具','')
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`cover_img`=VALUES(`cover_img`),`feeds_img`=VALUES(`feeds_img`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`anchor_name`=VALUES(`anchor_name`),`status`=VALUES(`status`),`live_status`=VALUES(`live_status`),`is_show`=VALUES(`is_show`),`mark`=VALUES(`mark`),`refusal`=VALUES(`refusal`);
INSERT INTO `qixi_crm_b_broadcast_room_goods` (`broadcast_room_id`,`product_id`,`on_sale`,`sort`) VALUES
  (7101,1001,1,1),(7101,1002,1,2),(7102,1101,1,1),(7102,1104,1,2)
ON DUPLICATE KEY UPDATE `on_sale`=VALUES(`on_sale`),`sort`=VALUES(`sort`);

-- 积分商城消费投影：库存与所需积分均由业务事件同步，以下为本地中文验收夹具。
INSERT INTO `qixi_crm_b_points_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`title`,`cover_url`,`original_price`,`points_required`,`stock`,`sale_status`,`version`) VALUES
  (1005,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店','真丝印花方巾礼盒','/demo/product-scarf-v1.png',169.00,120,20,1,1),
  (1105,2,2,'CRM Live居家商户','CRM Live居家优选店','真丝睡眠眼罩方巾组','/demo/product-scarf-v1.png',139.00,180,16,1,1),
  (1204,3,3,'CRM Live数码商户','CRM Live数码生活店','便携保温杯清洁套装','/demo/product-tumbler-v1.png',179.00,220,12,1,1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`merchant_name`=VALUES(`merchant_name`),`store_name`=VALUES(`store_name`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`original_price`=VALUES(`original_price`),`points_required`=VALUES(`points_required`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`version`=VALUES(`version`),`updated_at`=NOW();

-- 预约服务公开活动夹具。排期余量来自 qixi_crm_b_reservation_booking，测试数据不使用真实用户身份。
INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1301,1,1,'CRM Live服饰商户','CRM Live服饰旗舰店',10101,'秋日衣橱搭配咨询','/demo/product-knit-v1.png',99.00,129.00,4,46,12,1,1,NOW()),
  (1302,2,2,'CRM Live居家商户','CRM Live居家优选店',10201,'居家香氛体验服务','/demo/product-fragrance-v1.png',129.00,169.00,4,38,10,1,1,NOW()),
  (1303,3,3,'CRM Live数码商户','CRM Live数码生活店',10301,'跑步鞋试穿与选购服务','/demo/product-shoes-v1.png',79.00,99.00,4,29,8,1,1,NOW())
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`merchant_name`=VALUES(`merchant_name`),`store_name`=VALUES(`store_name`),`category_id`=VALUES(`category_id`),`title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`price`=VALUES(`price`),`original_price`=VALUES(`original_price`),`product_type`=VALUES(`product_type`),`sales`=VALUES(`sales`),`stock`=VALUES(`stock`),`sale_status`=VALUES(`sale_status`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_reservation_activity` (`product_reservation_id`,`product_id`,`merchant_id`,`store_id`,`reservation_type`,`show_reservation_days`,`is_cancel_reservation`,`time_period`,`status`) VALUES
  (6701,1301,1,1,1,7,1,'[]',1),(6702,1302,2,2,1,7,1,'[]',1),(6703,1303,3,3,1,7,1,'[]',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`reservation_type`=VALUES(`reservation_type`),`show_reservation_days`=VALUES(`show_reservation_days`),`is_cancel_reservation`=VALUES(`is_cancel_reservation`),`time_period`=VALUES(`time_period`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_b_reservation_slot` (`attr_reservation_id`,`product_id`,`slot_key`,`start_time`,`end_time`,`stock`,`use_num`) VALUES
  (6711,1301,'fashion-am','10:00','11:00',6,0),(6712,1301,'fashion-pm','14:00','15:00',6,0),
  (6721,1302,'home-am','10:30','11:30',5,0),(6722,1302,'home-pm','15:00','16:00',5,0),
  (6731,1303,'shoe-am','09:30','10:30',4,0),(6732,1303,'shoe-pm','16:00','17:00',4,0)
ON DUPLICATE KEY UPDATE `slot_key`=VALUES(`slot_key`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`stock`=VALUES(`stock`);

-- 拼团公开展示夹具。活动规则、价格和库存仍由服务端状态机校验，前端不伪造团购成功。
-- 秒杀夹具只包含虚构中文商品和商户，不含个人信息；时段和活动均在业务库内验收。
INSERT INTO `qixi_crm_b_seckill_time` (`seckill_time_id`,`title`,`start_time`,`end_time`,`status`,`pic`) VALUES
  (1,'上午场',9,12,1,''),(2,'晚间场',19,22,1,'')
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`status`=VALUES(`status`),`pic`=VALUES(`pic`);
INSERT INTO `qixi_crm_b_seckill_active` (`seckill_active_id`,`name`,`seckill_time_ids`,`start_day`,`end_day`,`mer_id`,`product_id`,`seckill_price`,`once_pay_count`,`all_pay_count`,`active_status`,`status`,`create_time`,`update_time`,`delete_time`) VALUES
  (6001,'轻奢羊绒针织衫限时秒杀','1,2',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),1,1001,199.00,1,0,1,1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL),
  (6002,'头层牛皮托特包限时秒杀','2',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),1,1002,329.00,1,0,1,1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL),
  (6003,'无火藤条香氛礼盒限时秒杀','1',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),2,1101,169.00,1,0,1,1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL),
  (6004,'智能数显保温杯限时秒杀','2',DATE_SUB(CURDATE(),INTERVAL 1 DAY),DATE_ADD(CURDATE(),INTERVAL 90 DAY),3,1201,149.00,1,0,1,1,UNIX_TIMESTAMP(),UNIX_TIMESTAMP(),NULL)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`seckill_time_ids`=VALUES(`seckill_time_ids`),`start_day`=VALUES(`start_day`),`end_day`=VALUES(`end_day`),`mer_id`=VALUES(`mer_id`),`product_id`=VALUES(`product_id`),`seckill_price`=VALUES(`seckill_price`),`once_pay_count`=VALUES(`once_pay_count`),`active_status`=VALUES(`active_status`),`status`=VALUES(`status`),`update_time`=VALUES(`update_time`),`delete_time`=VALUES(`delete_time`);

INSERT INTO `qixi_crm_b_combination_group` (`product_group_id`,`product_id`,`start_time`,`end_time`,`time`,`buying_count_num`,`buying_num`,`pay_count`,`once_pay_count`,`status`,`mer_id`,`is_show`,`is_del`,`success_num`,`product_status`,`price`,`action_status`) VALUES
  (6101,1001,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 90 DAY),24,2,1,0,1,1,1,1,0,18,1,239.00,1),
  (6102,1002,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 90 DAY),24,3,1,0,1,1,1,1,0,12,1,369.00,1),
  (6103,1101,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 90 DAY),24,2,1,0,1,1,2,1,0,26,1,189.00,1),
  (6104,1201,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 90 DAY),24,4,1,0,1,1,3,1,0,9,1,159.00,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`time`=VALUES(`time`),`buying_count_num`=VALUES(`buying_count_num`),`mer_id`=VALUES(`mer_id`),`is_show`=VALUES(`is_show`),`is_del`=VALUES(`is_del`),`success_num`=VALUES(`success_num`),`price`=VALUES(`price`),`action_status`=VALUES(`action_status`);

-- 预售与好友助力公开展示夹具。没有登录时只允许浏览；创建订单仍需要地址、JWT 和后端状态机校验。
INSERT INTO `qixi_crm_b_presell` (`product_presell_id`,`start_time`,`end_time`,`final_start_time`,`final_end_time`,`status`,`presell_type`,`pay_count`,`delivery_type`,`delivery_day`,`product_id`,`price`,`down_price`,`final_price`,`stock`,`is_show`,`store_name`,`mer_id`,`store_info`,`is_del`,`product_status`,`action_status`,`seles`) VALUES
  (6201,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),DATE_ADD(NOW(),INTERVAL 61 DAY),DATE_ADD(NOW(),INTERVAL 75 DAY),1,2,0,1,15,1004,299.00,59.00,240.00,36,1,'精纺圆领羊毛开衫',1,'秋冬预售，尾款支付后按订单地址发货。',0,1,1,48),
  (6202,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),'','',1,1,0,1,10,1104,238.00,0.00,0.00,39,1,'晨间居家香氛套装',2,'全款预售，预计十日内发货。',0,1,1,37),
  (6203,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),DATE_ADD(NOW(),INTERVAL 61 DAY),DATE_ADD(NOW(),INTERVAL 75 DAY),1,2,0,1,20,1206,329.00,69.00,260.00,35,1,'轻缓震训练跑鞋',3,'新品预售，支持定金与尾款支付。',0,1,1,29)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`presell_type`=VALUES(`presell_type`),`price`=VALUES(`price`),`down_price`=VALUES(`down_price`),`final_price`=VALUES(`final_price`),`stock`=VALUES(`stock`),`is_show`=VALUES(`is_show`),`store_name`=VALUES(`store_name`),`store_info`=VALUES(`store_info`),`is_del`=VALUES(`is_del`),`product_status`=VALUES(`product_status`),`action_status`=VALUES(`action_status`);

INSERT INTO `qixi_crm_b_assist` (`product_assist_id`,`start_time`,`end_time`,`status`,`pay_count`,`assist_count`,`assist_user_count`,`product_id`,`assist_price`,`stock`,`is_show`,`store_name`,`mer_id`,`store_info`,`is_del`,`product_status`,`action_status`) VALUES
  (6301,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),1,0,2,1,1005,99.00,90,1,'真丝印花方巾礼盒',1,'好友助力满员后可按助力价下单。',0,1,1),
  (6302,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),1,0,3,2,1107,239.00,26,1,'客厅氛围香薰礼盒',2,'邀请好友共同助力，名额用完即止。',0,1,1),
  (6303,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),1,0,4,3,1207,179.00,38,1,'桌面恒温杯垫礼盒',3,'助力价与普通售价以活动页实时校验为准。',0,1,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`start_time`=VALUES(`start_time`),`end_time`=VALUES(`end_time`),`assist_count`=VALUES(`assist_count`),`assist_user_count`=VALUES(`assist_user_count`),`assist_price`=VALUES(`assist_price`),`stock`=VALUES(`stock`),`is_show`=VALUES(`is_show`),`store_name`=VALUES(`store_name`),`store_info`=VALUES(`store_info`),`is_del`=VALUES(`is_del`),`product_status`=VALUES(`product_status`),`action_status`=VALUES(`action_status`);

-- PC/H5 演示消费者主体：仅用于本地与测试环境监管验收，不是生产账号。
-- 本夹具不提供可登录密码或身份凭据；客户端、接口和日志不得返回凭据。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live体验用户','DEMO-USER-9101',1,9501,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`group_id`=VALUES(`group_id`),`auth_version`=VALUES(`auth_version`);

-- 用户搜索记录夹具不含联系方式或设备指纹，仅供平台检索、清理和 CSV 导出验收。
INSERT INTO `qixi_crm_b_user_search_record` (`id`,`user_id`,`keyword`,`source`,`created_at`,`deleted_at`) VALUES
  (970001,9101,'夏季亚麻衬衫','pc','2026-08-03 10:20:00',NULL),
  (970002,9101,'居家香薰礼盒','h5','2026-08-03 11:35:00',NULL),
  (970003,9001,'儿童雨靴','mini','2026-08-04 09:10:00',NULL)
ON DUPLICATE KEY UPDATE `keyword`=VALUES(`keyword`),`source`=VALUES(`source`),`created_at`=VALUES(`created_at`),`deleted_at`=NULL;

-- 商品评论夹具仅使用虚构中文内容。平台可审核展示状态，但不能篡改用户原始评论。
INSERT INTO `qixi_crm_b_product_comment` (`id`,`order_item_id`,`user_id`,`product_id`,`store_id`,`score`,`content`,`media`,`reply_content`,`status`) VALUES
  (8801,980001,9101,1001,1,5,'虚构中文评论：羊绒柔软，版型合身。',JSON_ARRAY(),'感谢您的认可。','pending'),
  (8802,980002,9101,1101,2,4,'虚构中文评论：香氛淡雅，包装完整。',JSON_ARRAY(),'','published')
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`store_id`=VALUES(`store_id`),`score`=VALUES(`score`),`content`=VALUES(`content`),`media`=VALUES(`media`),`reply_content`=VALUES(`reply_content`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_b_product_comment` (`id`,`order_item_id`,`user_id`,`product_id`,`store_id`,`score`,`content`,`media`,`source`,`virtual_author_name`,`virtual_author_avatar`,`sort`,`status`) VALUES
  (8803,NULL,0,1001,1,5,'虚构中文虚拟评论：上身显气质，换季也很百搭。',JSON_ARRAY(),'virtual','演示用户小满','',80,'published')
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`store_id`=VALUES(`store_id`),`score`=VALUES(`score`),`content`=VALUES(`content`),`media`=VALUES(`media`),`source`=VALUES(`source`),`virtual_author_name`=VALUES(`virtual_author_name`),`virtual_author_avatar`=VALUES(`virtual_author_avatar`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`deleted_at`=NULL;

-- 用户运营夹具均为虚构中文名称，仅用于后台标签、分组与打标闭环验收。
INSERT INTO `qixi_crm_b_user_label` (`label_id`,`label_name`,`sort`,`is_del`) VALUES
  (9401,'高频复购用户',30,0),(9402,'香氛兴趣用户',20,0),(9403,'售后关怀用户',10,0)
ON DUPLICATE KEY UPDATE `label_name`=VALUES(`label_name`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_b_user_group` (`group_id`,`group_name`,`sort`,`is_del`) VALUES
  (9501,'CRM Live精选会员',20,0),(9502,'新品体验用户',10,0)
ON DUPLICATE KEY UPDATE `group_name`=VALUES(`group_name`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_b_user_label_relation` (`uid`,`label_id`) VALUES
  (9101,9401),(9101,9402)
ON DUPLICATE KEY UPDATE `label_id`=VALUES(`label_id`);

INSERT INTO `qixi_crm_b_user_identity` (`user_id`,`channel`,`subject`,`credential_hash`) VALUES
  (9101,'pc','DEMO-USER-9101',NULL),
  (9101,'h5','DEMO-USER-9101',NULL)
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`subject`=VALUES(`subject`),`credential_hash`=NULL;

-- 助力单夹具仅提供“已满员可下单”起点；订单、库存冻结与支付状态必须由 API 状态机生成。
INSERT INTO `qixi_crm_b_assist_set` (`product_assist_set_id`,`product_assist_id`,`product_id`,`uid`,`status`,`assist_count`,`assist_user_count`,`yet_assist_count`,`mer_id`,`is_del`) VALUES
  (63101,6301,1005,9101,10,2,1,2,1,0)
ON DUPLICATE KEY UPDATE `status`=VALUES(`status`),`assist_count`=VALUES(`assist_count`),`assist_user_count`=VALUES(`assist_user_count`),`yet_assist_count`=VALUES(`yet_assist_count`),`is_del`=VALUES(`is_del`);
INSERT INTO `qixi_crm_b_assist_user` (`product_assist_set_id`,`product_assist_id`,`uid`,`nickname`,`avatar_img`) VALUES
  (63101,6301,9201,'小林',''),(63101,6301,9202,'阿澈','')
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`avatar_img`=VALUES(`avatar_img`);

-- 会员等级与成长记录夹具仅用于本地验收，不代表真实消费权益或付费会员开通结果。
INSERT INTO `qixi_crm_b_member_level` (`id`,`name`,`rank`,`rules`,`benefits`,`status`,`version`,`deleted_at`) VALUES
  (8101,'普通会员',1,JSON_OBJECT('description','注册即享基础会员服务'),JSON_ARRAY('积分兑换','优惠券领取'),1,1,NULL),
  (8102,'悦享会员',2,JSON_OBJECT('description','满足成长规则后自动升级'),JSON_ARRAY('专属优惠提醒','会员活动优先参与'),1,1,NULL),
  (8103,'尊享会员',3,JSON_OBJECT('description','以平台实时会员规则为准'),JSON_ARRAY('优先客服','会员专享活动'),1,1,NULL)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`rank`=VALUES(`rank`),`rules`=VALUES(`rules`),`benefits`=VALUES(`benefits`),`status`=VALUES(`status`),`version`=VALUES(`version`),`deleted_at`=NULL;

INSERT INTO `qixi_crm_b_member_account` (`user_id`,`level_id`,`points`,`balance`,`commission`) VALUES
  (9101,8102,268,36.50,0.00)
ON DUPLICATE KEY UPDATE `level_id`=VALUES(`level_id`),`points`=VALUES(`points`),`balance`=VALUES(`balance`),`commission`=VALUES(`commission`);

-- 用户提现待登记打款夹具：账户快照为虚构脱敏值，平台只显示内部凭证编号，不读取该字段。
INSERT INTO `qixi_crm_b_withdrawal_application` (`id`,`withdrawal_no`,`user_id`,`amount`,`channel`,`account_snapshot`,`status`,`review_note`,`reviewed_by`,`reviewed_at`,`idempotency_key`) VALUES
  (9201001,'WD-DEMO-20260803-001',9101,8.50,'wechat',JSON_OBJECT('account_name','虚构用户','account_no','已脱敏演示账户'),'approved','中文演示审批通过，等待登记内部打款凭证。',9001,NOW(),'fixture-withdraw-9101-01')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`status`=VALUES(`status`),`review_note`=VALUES(`review_note`),`reviewed_by`=VALUES(`reviewed_by`),`reviewed_at`=VALUES(`reviewed_at`),`payout_idempotency_key`=NULL,`payout_reference`=NULL,`paid_by`=NULL,`paid_at`=NULL;

-- 分销夹具：推广资格由业务后台授权，用户绑定关系一经建立不可更换。
INSERT INTO `qixi_crm_b_distribution_promoter` (`user_id`,`status`) VALUES
  (9001,1),(9101,1)
ON DUPLICATE KEY UPDATE `status`=VALUES(`status`);
-- 推荐关系只使用虚构演示账号；已绑定关系不可由平台用户详情页修改。
INSERT INTO `qixi_crm_b_distribution_relation` (`user_id`,`parent_user_id`,`bound_at`) VALUES
  (9101,9001,DATE_SUB(NOW(),INTERVAL 7 DAY))
ON DUPLICATE KEY UPDATE `parent_user_id`=VALUES(`parent_user_id`),`bound_at`=VALUES(`bound_at`);
INSERT INTO `qixi_crm_b_commission_ledger` (`user_id`,`order_id`,`amount`,`status`,`idempotency_key`,`available_at`,`created_at`) VALUES
  (9101,0,12.80,'available','fixture-commission-9101-01',NOW(),DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (9101,0,6.50,'pending','fixture-commission-9101-02',NULL,DATE_SUB(NOW(),INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`status`=VALUES(`status`),`available_at`=VALUES(`available_at`),`created_at`=VALUES(`created_at`);

-- 余额与站内通知夹具仅用于本地验收；收支与通知均为虚构中文示例。
INSERT INTO `qixi_crm_b_asset_ledger` (`user_id`,`asset_type`,`amount`,`reference_type`,`reference_id`,`idempotency_key`,`created_at`) VALUES
  (9101,'balance',66.50,'demo_recharge','balance-demo-01','fixture-balance-income-9101',DATE_SUB(NOW(),INTERVAL 3 DAY)),
  (9101,'balance',-30.00,'demo_order','balance-demo-02','fixture-balance-expense-9101',DATE_SUB(NOW(),INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`),`reference_type`=VALUES(`reference_type`),`reference_id`=VALUES(`reference_id`),`created_at`=VALUES(`created_at`);
INSERT INTO `qixi_crm_b_notification` (`id`,`user_id`,`category`,`title`,`body`,`payload`,`read_at`,`created_at`) VALUES
  (9901,9101,'order','订单服务提醒','你的演示订单状态已更新，请在订单中心查看详情。',JSON_OBJECT('source','fixture'),NULL,DATE_SUB(NOW(),INTERVAL 2 HOUR)),
  (9902,9101,'member','会员权益提示','当前会员等级可享积分兑换与优惠券领取服务。',JSON_OBJECT('source','fixture'),NULL,DATE_SUB(NOW(),INTERVAL 1 HOUR)),
  (9903,9101,'system','平台服务通知','本地演示环境已准备好中文模拟数据，可继续体验用户端功能。',JSON_OBJECT('source','fixture'),DATE_SUB(NOW(),INTERVAL 30 MINUTE),DATE_SUB(NOW(),INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`body`=VALUES(`body`),`payload`=VALUES(`payload`),`read_at`=VALUES(`read_at`),`created_at`=VALUES(`created_at`);

INSERT INTO `qixi_crm_b_member_level_log` (`id`,`user_id`,`level_id`,`previous_level_id`,`change_type`,`note`,`created_at`) VALUES
  (81001,9101,8102,8101,'upgrade','完成演示会员成长任务，已升级为悦享会员。',DATE_SUB(NOW(),INTERVAL 3 DAY))
ON DUPLICATE KEY UPDATE `level_id`=VALUES(`level_id`),`previous_level_id`=VALUES(`previous_level_id`),`change_type`=VALUES(`change_type`),`note`=VALUES(`note`),`created_at`=VALUES(`created_at`);

INSERT INTO `qixi_crm_b_user_sign` (`id`,`user_id`,`sign_date`,`points`,`continuous_days`,`created_at`) VALUES
  (81101,9101,DATE_SUB(CURDATE(),INTERVAL 2 DAY),5,1,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (81102,9101,DATE_SUB(CURDATE(),INTERVAL 1 DAY),5,2,DATE_SUB(NOW(),INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE `points`=VALUES(`points`),`continuous_days`=VALUES(`continuous_days`),`created_at`=VALUES(`created_at`);

-- 发票抬头夹具仅用于本地验收，企业名称、税号和邮箱均为虚构示例，不对应真实主体。
INSERT INTO `qixi_crm_b_user_invoice_profile` (`id`,`user_id`,`type`,`title`,`tax_no`,`email`,`is_default`) VALUES
  (9601,9101,'enterprise','CRM Live演示科技有限公司','91310000DEMO12345X','finance.invalid',1)
ON DUPLICATE KEY UPDATE `type`=VALUES(`type`),`title`=VALUES(`title`),`tax_no`=VALUES(`tax_no`),`email`=VALUES(`email`),`is_default`=VALUES(`is_default`);

-- 意见反馈夹具仅用于本地验收，内容为虚构中文建议，不包含联系方式。
INSERT INTO `qixi_crm_b_user_feedback_category` (`id`,`name`,`sort`,`status`) VALUES
  (9711,'功能建议',10,1),(9712,'订单问题',20,1),(9713,'使用体验',30,1),(9714,'历史分类',90,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`),`deleted_at`=NULL;
INSERT INTO `qixi_crm_b_user_feedback` (`id`,`user_id`,`category_id`,`type`,`content`,`status`,`reply`,`created_at`,`updated_at`) VALUES
  (9701,9101,9711,'功能建议','希望商品列表支持销量和价格排序。','replied','已在商品列表提供销量、价格升降序筛选。',DATE_SUB(NOW(),INTERVAL 2 DAY),NOW()),
  (9702,9101,9712,'订单问题','希望待付款订单可以直接取消。','pending','',DATE_SUB(NOW(),INTERVAL 1 DAY),NOW()),
  (9703,9101,9713,'使用体验','希望反馈处理完成后能看到关闭说明。','closed','已完成本次模拟验收并归档。',DATE_SUB(NOW(),INTERVAL 3 HOUR),NOW())
ON DUPLICATE KEY UPDATE `category_id`=VALUES(`category_id`),`type`=VALUES(`type`),`content`=VALUES(`content`),`status`=VALUES(`status`),`reply`=VALUES(`reply`),`deleted_at`=NULL,`updated_at`=NOW();

-- 浏览记录夹具只关联本地演示用户与公开商品，不包含真实访问行为。
INSERT INTO `qixi_crm_b_user_browse_history` (`id`,`user_id`,`product_id`,`store_id`,`viewed_at`) VALUES
  (9801,9101,1001,1,DATE_SUB(NOW(),INTERVAL 2 HOUR)),
  (9802,9101,1101,2,DATE_SUB(NOW(),INTERVAL 1 HOUR))
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`product_id`=VALUES(`product_id`),`store_id`=VALUES(`store_id`),`viewed_at`=VALUES(`viewed_at`);

-- 社区公开内容仅为演示数据，不包含真实个人信息。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`) VALUES
  (9001,'晴空漫游者',NULL,1),(9002,'居家研究员',NULL,1),(9003,'通勤玩家',NULL,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_b_social_category` (`category_id`,`cate_name`,`pid`,`is_show`,`sort`) VALUES
  (6401,'穿搭分享',0,1,10),(6402,'居家生活',0,1,20),(6403,'数码体验',0,1,30)
ON DUPLICATE KEY UPDATE `cate_name`=VALUES(`cate_name`),`is_show`=VALUES(`is_show`),`sort`=VALUES(`sort`);
INSERT INTO `qixi_crm_b_social_topic` (`topic_id`,`topic_name`,`status`,`is_hot`,`category_id`,`is_del`,`count_use`,`sort`) VALUES
  (6501,'通勤穿搭',1,1,6401,0,18,10),(6502,'居家氛围',1,1,6402,0,12,20),(6503,'好物测评',1,0,6403,0,9,30)
ON DUPLICATE KEY UPDATE `topic_name`=VALUES(`topic_name`),`status`=VALUES(`status`),`is_hot`=VALUES(`is_hot`),`category_id`=VALUES(`category_id`),`is_del`=VALUES(`is_del`),`count_use`=VALUES(`count_use`),`sort`=VALUES(`sort`);
INSERT INTO `qixi_crm_b_social_post` (`community_id`,`title`,`image`,`category_id`,`topic_id`,`uid`,`mer_id`,`product_id`,`count_start`,`count_reply`,`status`,`is_show`,`is_hot`,`is_type`,`content`,`refusal`,`pv`,`is_del`,`status_time`) VALUES
  (6601,'通勤针织衫的三种叠穿思路','/demo/product-knit-v1.png',6401,6501,9001,1,1001,26,3,1,1,1,0,'柔软的针织衫适合和衬衫、半裙搭配，通勤和周末都能穿。','',128,0,NOW()),
  (6602,'让客厅更放松的香氛组合','/demo/product-fragrance-v1.png',6402,6502,9002,2,1101,18,2,1,1,0,0,'无火香薰放在玄关，蜡烛安排在晚间阅读角，层次更自然。','',96,0,NOW()),
  (6603,'桌面恒温杯垫值得入手吗','/demo/product-tumbler-v1.png',6403,6503,9003,3,1207,11,1,1,1,0,0,'适合常坐在电脑前的人，配保温杯使用更方便。','',65,0,NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`image`=VALUES(`image`),`category_id`=VALUES(`category_id`),`topic_id`=VALUES(`topic_id`),`uid`=VALUES(`uid`),`mer_id`=VALUES(`mer_id`),`product_id`=VALUES(`product_id`),`count_start`=VALUES(`count_start`),`count_reply`=VALUES(`count_reply`),`status`=VALUES(`status`),`is_show`=VALUES(`is_show`),`is_hot`=VALUES(`is_hot`),`content`=VALUES(`content`),`pv`=VALUES(`pv`),`is_del`=VALUES(`is_del`),`status_time`=NOW();

-- 领券中心夹具：store_id=0 为平台券，其余为对应店铺券。领取记录只在用户实际点击领取后生成。
INSERT INTO `qixi_crm_b_coupon_template_view` (`coupon_id`,`store_id`,`name`,`discount_type`,`discount_value`,`min_amount`,`starts_at`,`ends_at`,`status`,`version`) VALUES
  (3001,0,'平台新客满99减10','amount',10.00,99.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3002,0,'平台夏日满299减40','amount',40.00,299.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3003,1,'服饰店满199减30','amount',30.00,199.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3004,1,'服饰店9折券','rate',90.00,399.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3005,2,'居家优选满159减20','amount',20.00,159.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1),
  (3006,3,'数码生活满299减35','amount',35.00,299.00,DATE_SUB(NOW(), INTERVAL 1 DAY),DATE_ADD(NOW(), INTERVAL 180 DAY),1,1)
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`name`=VALUES(`name`),`discount_type`=VALUES(`discount_type`),`discount_value`=VALUES(`discount_value`),`min_amount`=VALUES(`min_amount`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`),`status`=VALUES(`status`),`version`=VALUES(`version`);
-- 新人礼中文夹具：新注册用户将在创建账户的同一事务内自动获得下列有效券。
INSERT INTO `qixi_crm_b_onboarding_policy` (`id`,`enabled`,`coupon_enabled`,`title`,`description`) VALUES
  (1,1,1,'新人专享礼','注册成功后优惠券已自动放入账户')
ON DUPLICATE KEY UPDATE `enabled`=VALUES(`enabled`),`coupon_enabled`=VALUES(`coupon_enabled`),`title`=VALUES(`title`),`description`=VALUES(`description`);
INSERT INTO `qixi_crm_b_onboarding_coupon` (`coupon_id`,`enabled`,`sort`) VALUES
  (3001,1,10),(3002,1,20)
ON DUPLICATE KEY UPDATE `enabled`=VALUES(`enabled`),`sort`=VALUES(`sort`);
INSERT INTO `qixi_crm_b_coupon_user` (`id`,`user_id`,`coupon_id`,`source`,`status`,`obtained_at`) VALUES
  (93001,9101,3001,'fixture','unused',DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (93002,9101,3003,'fixture','used',DATE_SUB(NOW(),INTERVAL 4 DAY)),
  (93003,9101,3002,'platform_manual','expired',DATE_SUB(NOW(),INTERVAL 1 DAY))
ON DUPLICATE KEY UPDATE `source`=VALUES(`source`),`status`=VALUES(`status`),`obtained_at`=VALUES(`obtained_at`);
-- 平台人工发券审计夹具。仅用于后台发送记录演示，未伪造支付、核销或真实用户资料。
INSERT INTO `qixi_crm_b_user_coupon_command_audit` (`id`,`user_id`,`coupon_id`,`coupon_user_id`,`action`,`from_status`,`to_status`,`reason`,`operator_admin_id`,`idempotency_key`,`created_at`) VALUES
  (94001,9101,3002,93003,'issue','','unused','虚构中文工单：补发夏日活动体验券。',9201,'fixture-coupon-issue-94001',DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (94002,9101,3002,93003,'revoke','unused','expired','虚构中文工单：活动资格调整，撤销未使用体验券。',9201,'fixture-coupon-revoke-94002',DATE_SUB(NOW(),INTERVAL 23 HOUR))
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`coupon_id`=VALUES(`coupon_id`),`coupon_user_id`=VALUES(`coupon_user_id`),`from_status`=VALUES(`from_status`),`to_status`=VALUES(`to_status`),`reason`=VALUES(`reason`),`operator_admin_id`=VALUES(`operator_admin_id`),`created_at`=VALUES(`created_at`);

-- 秒杀展示夹具。规则由后台营销活动投影而来；C 端只读 qixi_crm_b_marketing_activity_view。
INSERT INTO `qixi_crm_b_marketing_activity_view` (`activity_id`,`store_id`,`activity_type`,`name`,`rules`,`status`,`version`,`starts_at`,`ends_at`) VALUES
  (5001,1,'seckill','轻奢羊绒针织衫限时抢购',JSON_OBJECT('product_id',1001,'seckill_price',199.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5002,1,'seckill','头层牛皮托特包限时抢购',JSON_OBJECT('product_id',1002,'seckill_price',329.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5003,2,'seckill','无火藤条香氛礼盒限时抢购',JSON_OBJECT('product_id',1101,'seckill_price',169.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5004,3,'seckill','智能数显保温杯限时抢购',JSON_OBJECT('product_id',1201,'seckill_price',149.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5101,1,'discount','夏日香氛随行套餐',JSON_OBJECT('package_price',199.00,'product_ids',JSON_ARRAY(1004,1006),'free_shipping',true,'remark','中文演示套餐'),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY))
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`name`=VALUES(`name`),`rules`=VALUES(`rules`),`status`=VALUES(`status`),`version`=VALUES(`version`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`);

-- 直播公开读模型夹具。直播流地址仅为本地演示标识，客户端不会把它当作生产推流密钥。
INSERT INTO `qixi_crm_b_live_room` (`id`,`merchant_id`,`store_id`,`anchor_user_id`,`title`,`anchor_name`,`cover_url`,`status`,`is_public`,`stream_ref`,`play_url`,`starts_at`,`ended_at`,`sort`) VALUES
  (7001,1,1,NULL,'CRM Live秋日衣橱直播专场','CRM Live小夏','/demo/home-hero-accessories-v1.png','living',1,'local-demo-fashion','',DATE_SUB(NOW(), INTERVAL 20 MINUTE),NULL,20),
  (7002,2,2,NULL,'居家香氛与生活好物分享','居家优选主播','/demo/home-hero-fragrance-v1.png','scheduled',1,'local-demo-home','',DATE_ADD(NOW(), INTERVAL 1 HOUR),NULL,10),
  (7003,3,3,NULL,'通勤数码好物直播回放','数码生活主播','/demo/home-tech-wide-v1.png','ended',1,'local-demo-digital','',DATE_SUB(NOW(), INTERVAL 2 DAY),DATE_SUB(NOW(), INTERVAL 1 DAY),5)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`title`=VALUES(`title`),`anchor_name`=VALUES(`anchor_name`),`cover_url`=VALUES(`cover_url`),`status`=VALUES(`status`),`is_public`=VALUES(`is_public`),`stream_ref`=VALUES(`stream_ref`),`play_url`=VALUES(`play_url`),`starts_at`=VALUES(`starts_at`),`ended_at`=VALUES(`ended_at`),`sort`=VALUES(`sort`);

INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2001,'notice','CRM Live商城服务公告','/demo/home-hero-v1.png','CRM Live商城已上线商品、订单、售后和客服服务。消费者可通过商品详情、购物车和订单中心完成全流程购物。',1,1,NOW(),NOW()),
  (2002,'notice','消费者权益说明','/demo/home-service-wide-v1.png','请在下单前确认商品信息、配送方式和售后规则。如有商品与履约问题，可在订单中心提交售后申请。',1,1,NOW(),NOW()),
  (2003,'notice','夏日居家焕新季：精选家居好物上新','/demo/home-tech-wide-v1.png','居家生活专区已上新香氛、随行杯与织物护理系列，支持按分类、销量和价格快速筛选。',1,1,NOW(),NOW()),
  (2004,'notice','CRM Live多商户店铺服务规范','/demo/home-service-vertical-v1.png','平台持续完善商户审核、商品审核、订单履约与售后处理规范，为消费者提供清晰可靠的购物体验。',1,1,NOW(),NOW()),
  (2005,'notice','会员积分与优惠券使用说明','/demo/home-beauty-vertical-v1.png','积分、优惠券将按各自规则展示和使用。结算页会提示可用权益与优惠金额。',1,1,NOW(),NOW()),
  (2101,'agreement','sys_user_agree','','欢迎使用CRM Live商城。使用服务前请阅读并同意本用户协议。',1,1,NOW(),NOW()),
  (2102,'agreement','sys_userr_privacy','','CRM Live仅在提供服务所必需的范围内处理您的个人信息。',1,1,NOW(),NOW()),
  (2103,'agreement','sys_integral_agree','','积分可用于积分商城兑换。兑换订单独立结算，积分余额与订单状态以服务端实时校验结果为准。',1,1,NOW(),NOW()),
  (2104,'agreement','sys_about_us','','CRM Live 致力于为消费者提供清晰、可靠的多商户购物服务。平台持续完善商品、订单、售后和权益保障能力。',1,1,NOW(),NOW()),
  (2201,'article','夏日焕新购物指南','/demo/home-hero-v1.png','精选商品已按类目、销量与价格开放筛选。下单前请确认商品规格、配送方式与售后说明。',1,1,NOW(),NOW()),
  (2202,'article','积分商城兑换说明','/demo/home-service-wide-v1.png','积分兑换订单独立结算，积分余额与库存以提交订单时服务端实时校验结果为准。',1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`cover_url`=VALUES(`cover_url`),`body`=VALUES(`body`),`status`=VALUES(`status`),`version`=VALUES(`version`),`published_at`=VALUES(`published_at`),`updated_at`=NOW();

-- 本地投影夹具；生产环境只允许 api-merchant 的 outbox/NATS 事件写入。
INSERT INTO `qixi_crm_b_diy_page_view` (`source`,`page_id`,`store_id`,`page_type`,`name`,`document`,`status`,`is_active`) VALUES
  ('merchant',3001,1,'home','CRM Live演示店铺首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','CRM Live演示店铺','title','CRM Live演示店铺')),
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(JSON_OBJECT('imgName','CRM Live演示店铺','imgUrl','','linkUrl','/pages/store/index'))),
      JSON_OBJECT('type','navBar','name','导航组','data',JSON_ARRAY(JSON_OBJECT('text','全部商品','imgUrl','','linkUrl','/pages/goods/list'),JSON_OBJECT('text','购物车','imgUrl','','linkUrl','/pages/order_addcart/order_addcart')))
    ),
    '_qixi',JSON_OBJECT('title','CRM Live演示店铺','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`),`is_active`=VALUES(`is_active`),`updated_at`=NOW();

INSERT INTO `qixi_crm_b_diy_page_view` (`source`,`page_id`,`store_id`,`page_type`,`name`,`document`,`status`,`is_active`) VALUES
  ('platform',4001,0,'home','CRM Live平台首页',JSON_OBJECT(
    'page',JSON_OBJECT('type','page','name','页面设置','params',JSON_OBJECT('name','CRM Live商城','title','CRM Live商城')),
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(
        JSON_OBJECT('imgName','CRM Live商城精选','imgUrl','/demo/home-hero-v1.png','linkUrl','/goods?cate_id=101'),
        JSON_OBJECT('imgName','CRM Live香氛家居','imgUrl','/demo/home-hero-fragrance-v1.png','linkUrl','/goods?cate_id=102'),
        JSON_OBJECT('imgName','CRM Live箱包配饰','imgUrl','/demo/home-hero-accessories-v1.png','linkUrl','/goods?cate_id=10102')
      )),
      JSON_OBJECT('type','product','name','服饰鞋包展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',101,'showNum',4,'productSort','sales')
      )),
      JSON_OBJECT('type','product','name','家居生活展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',102,'showNum',4,'productSort','sales')
      ))
    ),
    '_qixi',JSON_OBJECT('title','CRM Live商城','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`),`is_active`=VALUES(`is_active`),`updated_at`=NOW();

-- 客服坐席中文夹具：仅验证会话分配，不包含 IM 凭证、真实会话或客户消息。
INSERT INTO qixi_crm_b_customer_service_agent_view (admin_id,store_id,display_name,status,available_at) VALUES
  (9301,1,'客服小林',1,NOW()),(9302,2,'客服阿澈',1,NOW()),(9303,3,'客服小夏',1,NOW())
ON DUPLICATE KEY UPDATE display_name=VALUES(display_name),status=VALUES(status),available_at=VALUES(available_at);

-- 客服快捷回复中文夹具：均为明显虚构文本，仅用于授权店铺范围和软删除回归。
INSERT INTO qixi_crm_b_customer_service_quick_reply (id,store_id,title,content,status,created_by,updated_by,deleted_at) VALUES
  (9900101,1,'发货时效','您好，CRM Live演示茶铺将在 48 小时内为您安排发货。','enabled',9301,9301,NULL),
  (9900102,2,'售后指引','您好，虚构订单可在订单详情中提交售后申请，我们会尽快协助。','disabled',9302,9302,NULL)
ON DUPLICATE KEY UPDATE store_id=VALUES(store_id),title=VALUES(title),content=VALUES(content),status=VALUES(status),updated_by=VALUES(updated_by),deleted_at=NULL,updated_at=NOW();

INSERT INTO qixi_crm_b_customer_service_user_note (user_id,store_id,content,updated_by) VALUES
  (9101,1,'虚构用户备注：已告知CRM Live演示茶铺的售后处理时效。',9301)
ON DUPLICATE KEY UPDATE content=VALUES(content),updated_by=VALUES(updated_by),updated_at=NOW();

-- 客服订单辅助夹具：全部为本地虚构编号，用于授权会话→订单→商品/物流/售后的闭环验收。
INSERT INTO qixi_crm_b_group_order (id,order_no,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,pay_channel,pay_status,paid_at,idempotency_key,remark) VALUES
  (9900201,'CS-DEMO-G-20260803-001',9101,299.00,0.00,0.00,299.00,1,0,0,JSON_OBJECT('recipient','虚构收件人','mobile','演示号已脱敏'),'mock','paid',NOW(),'fixture-cs-order-9900201','客服订单辅助夹具')
ON DUPLICATE KEY UPDATE pay_status=VALUES(pay_status),paid_at=VALUES(paid_at),pay_amount=VALUES(pay_amount),updated_at=NOW();
-- 模拟支付流水仅用于本地 sandbox 验收；不是渠道订单号，也不含真实凭据。
INSERT INTO qixi_crm_b_payment_transaction (id,group_order_id,channel,transaction_no,amount,status,provider_transaction_no,callback_idempotency_key,paid_at) VALUES
  (9900201,9900201,'mock','CS-DEMO-G-20260803-001',299.00,'succeeded','mock-payment-CS-20260803-001','fixture-payment-callback-9900201',NOW())
ON DUPLICATE KEY UPDATE channel=VALUES(channel),amount=VALUES(amount),status=VALUES(status),provider_transaction_no=VALUES(provider_transaction_no),callback_idempotency_key=VALUES(callback_idempotency_key),paid_at=VALUES(paid_at);
INSERT INTO qixi_crm_b_order (id,group_order_id,order_no,merchant_id,merchant_name_snapshot,store_id,store_name_snapshot,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,remark,status,paid_at) VALUES
  (9900201,9900201,'CS-DEMO-O-20260803-001',1,'CRM Live服饰商户',1,'CRM Live服饰旗舰店',9101,299.00,0.00,0.00,299.00,1,0,0,JSON_OBJECT('recipient','虚构收件人','mobile','演示号已脱敏'),'客服辅助演示订单','shipped',NOW())
ON DUPLICATE KEY UPDATE status=VALUES(status),pay_amount=VALUES(pay_amount),paid_at=VALUES(paid_at),updated_at=NOW();
INSERT INTO qixi_crm_b_order_item (id,order_id,product_id,merchant_sku_id,sku_key,title_snapshot,cover_url_snapshot,spec_snapshot,unit_price,quantity,refund_quantity) VALUES
  (9900201,9900201,1001,61001,'61001','轻奢羊绒针织衫','/demo/product-knit-v1.png',JSON_OBJECT('默认','标准'),299.00,1,1)
ON DUPLICATE KEY UPDATE merchant_sku_id=VALUES(merchant_sku_id),sku_key=VALUES(sku_key),title_snapshot=VALUES(title_snapshot),spec_snapshot=VALUES(spec_snapshot),unit_price=VALUES(unit_price),quantity=VALUES(quantity),refund_quantity=VALUES(refund_quantity);
INSERT INTO qixi_crm_b_order_delivery (id,order_id,delivery_type,carrier_code,tracking_no,status,delivered_at) VALUES
  (9900201,9900201,'express','CRM Live演示快递','CSDEMO20260803001','shipped',NOW())
ON DUPLICATE KEY UPDATE carrier_code=VALUES(carrier_code),tracking_no=VALUES(tracking_no),status=VALUES(status),delivered_at=VALUES(delivered_at);
-- 订单发票夹具仅用于平台只读监管；税号、邮箱均为虚构示例，后台接口会脱敏返回。
INSERT INTO qixi_crm_b_order_invoice (id,order_id,invoice_profile_id,profile_type,title,tax_no,email,status,invoice_no,file_url,rejection_reason,requested_at,issued_at) VALUES
  (9900201,9900201,9601,'enterprise','CRM Live演示科技有限公司','91310000DEMO12345X','finance@invoice.invalid','issued','DEMO-INV-20260803-001','/demo/invoice-9900201.pdf','',DATE_SUB(NOW(),INTERVAL 1 DAY),NOW())
ON DUPLICATE KEY UPDATE title=VALUES(title),tax_no=VALUES(tax_no),email=VALUES(email),status=VALUES(status),invoice_no=VALUES(invoice_no),file_url=VALUES(file_url),rejection_reason=VALUES(rejection_reason),issued_at=VALUES(issued_at);
-- 待审发票夹具：供商户后台 invoice.audit 开票/驳回联调（独立订单，避免与上条 uk_order 冲突）。
INSERT INTO qixi_crm_b_group_order (id,order_no,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,pay_channel,pay_status,paid_at,idempotency_key,remark) VALUES
  (9900202,'CS-DEMO-G-20260803-002',9101,199.00,0.00,0.00,199.00,1,0,0,JSON_OBJECT('recipient','虚构收件人','mobile','演示号已脱敏'),'mock','paid',NOW(),'fixture-cs-order-9900202','商户发票待审夹具')
ON DUPLICATE KEY UPDATE pay_status=VALUES(pay_status),pay_amount=VALUES(pay_amount),paid_at=VALUES(paid_at),remark=VALUES(remark),updated_at=NOW();
INSERT INTO qixi_crm_b_payment_transaction (id,group_order_id,channel,transaction_no,amount,status,provider_transaction_no,callback_idempotency_key,paid_at) VALUES
  (9900202,9900202,'mock','CS-DEMO-G-20260803-002',199.00,'succeeded','mock-payment-CS-20260803-002','fixture-payment-callback-9900202',NOW())
ON DUPLICATE KEY UPDATE channel=VALUES(channel),amount=VALUES(amount),status=VALUES(status),provider_transaction_no=VALUES(provider_transaction_no),callback_idempotency_key=VALUES(callback_idempotency_key),paid_at=VALUES(paid_at);
INSERT INTO qixi_crm_b_order (id,group_order_id,order_no,merchant_id,merchant_name_snapshot,store_id,store_name_snapshot,user_id,total_amount,discount_amount,freight_amount,pay_amount,total_quantity,activity_type,points_amount,recipient_snapshot,remark,status,paid_at) VALUES
  (9900202,9900202,'CS-DEMO-O-20260803-002',1,'CRM Live服饰商户',1,'CRM Live服饰旗舰店',9101,199.00,0.00,0.00,199.00,1,0,0,JSON_OBJECT('recipient','虚构收件人','mobile','演示号已脱敏'),'商户发票待审演示订单','paid',NOW())
ON DUPLICATE KEY UPDATE status=VALUES(status),pay_amount=VALUES(pay_amount),paid_at=VALUES(paid_at),remark=VALUES(remark),updated_at=NOW();
INSERT INTO qixi_crm_b_order_item (id,order_id,product_id,merchant_sku_id,sku_key,title_snapshot,cover_url_snapshot,spec_snapshot,unit_price,quantity,refund_quantity) VALUES
  (9900202,9900202,1001,61001,'61001','轻奢羊绒针织衫','/demo/product-knit-v1.png',JSON_OBJECT('默认','标准'),199.00,1,0)
ON DUPLICATE KEY UPDATE merchant_sku_id=VALUES(merchant_sku_id),sku_key=VALUES(sku_key),title_snapshot=VALUES(title_snapshot),unit_price=VALUES(unit_price),quantity=VALUES(quantity),refund_quantity=VALUES(refund_quantity);
INSERT INTO qixi_crm_b_order_invoice (id,order_id,invoice_profile_id,profile_type,title,tax_no,email,status,invoice_no,file_url,rejection_reason,requested_at,issued_at) VALUES
  (9900202,9900202,9601,'enterprise','CRM Live演示科技有限公司','91310000DEMO12345X','finance@invoice.invalid','requested','','','',NOW(),NULL)
ON DUPLICATE KEY UPDATE title=VALUES(title),tax_no=VALUES(tax_no),email=VALUES(email),status=VALUES(status),invoice_no=VALUES(invoice_no),file_url=VALUES(file_url),rejection_reason=VALUES(rejection_reason),issued_at=VALUES(issued_at);
INSERT INTO qixi_crm_b_refund (id,order_id,refund_no,reason,amount,refund_type,order_status_before,status,idempotency_key) VALUES
  (9900201,9900201,'CS-DEMO-R-20260803-001','尺寸不合适，申请虚构退货退款演示。',299.00,'return_and_refund','shipped','awaiting_receipt','fixture-cs-refund-9900201')
ON DUPLICATE KEY UPDATE reason=VALUES(reason),amount=VALUES(amount),refund_type=VALUES(refund_type),status=VALUES(status),updated_at=NOW();
INSERT INTO qixi_crm_b_refund_return_shipment (id,refund_id,carrier_name,tracking_no,remark,submitted_by,submitted_at) VALUES
  (9900201,9900201,'CRM Live演示快递','RETURNDEMO20260803001','虚构退货物流夹具，仅用于本地闭环验收。',9101,NOW())
ON DUPLICATE KEY UPDATE carrier_name=VALUES(carrier_name),tracking_no=VALUES(tracking_no),remark=VALUES(remark),submitted_by=VALUES(submitted_by),submitted_at=NOW();
INSERT INTO qixi_crm_b_refund_event (id,refund_id,from_status,to_status,actor_type,actor_id,reason,idempotency_key,created_at) VALUES
  (9900201,9900201,'awaiting_return','awaiting_receipt','user',9101,'虚构演示：用户已提交退货物流，等待商户确认收货。','fixture-refund-event-9900201',NOW())
ON DUPLICATE KEY UPDATE reason=VALUES(reason),to_status=VALUES(to_status),created_at=VALUES(created_at);
INSERT INTO qixi_crm_b_customer_service_binding (id,user_id,store_id,order_id,im_conversation_id,status,assigned_admin_id,assigned_at,last_msg,last_time) VALUES
  (9900201,9101,1,9900201,'cs-demo-CRM Live-20260803-001','open',9301,NOW(),'客服订单辅助演示会话已创建。',NOW())
ON DUPLICATE KEY UPDATE order_id=VALUES(order_id),status=VALUES(status),assigned_admin_id=VALUES(assigned_admin_id),assigned_at=VALUES(assigned_at),last_msg=VALUES(last_msg),last_time=VALUES(last_time),updated_at=NOW();
-- 客服转接审计夹具：仅使用虚构坐席 ID 和中文原因，不含会话正文或 IM 凭证。
INSERT INTO qixi_crm_b_customer_service_assignment_log (binding_id,from_admin_id,target_admin_id,operator_admin_id,reason,idempotency_key) VALUES
  (9900201,9301,9302,9301,'虚构演示：用户咨询居家商品，转交对应队列客服。','fixture-cs-transfer-9900201')
ON DUPLICATE KEY UPDATE reason=VALUES(reason),target_admin_id=VALUES(target_admin_id),operator_admin_id=VALUES(operator_admin_id);


-- 资金域中文模拟数据：仅提供可展示方案，不创建“已支付”订单，不伪造渠道到账。
INSERT INTO `qixi_crm_b_recharge_plan` (`id`,`name`,`amount`,`bonus_amount`,`status`,`sort`) VALUES
  (980001,'100 元充值',100.00,8.00,1,10),
  (980002,'300 元充值',300.00,35.00,1,20),
  (980003,'500 元充值',500.00,80.00,1,30)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`amount`=VALUES(`amount`),`bonus_amount`=VALUES(`bonus_amount`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`updated_at`=NOW();
INSERT INTO `qixi_crm_b_svip_interest` (`id`,`name`,`description`,`icon_url`,`status`,`sort`,`version`,`deleted_at`) VALUES
  (981001,'会员专享价','符合条件的商品可按会员专享价结算。','/demo/svip-price.png',1,10,1,NULL),
  (981002,'专属客服优先响应','会员咨询进入专属客服优先队列。','/demo/svip-service.png',1,20,1,NULL),
  (981003,'每月权益提醒','每月通过站内消息提醒可用会员权益。','/demo/svip-reminder.png',1,30,1,NULL),
  (981004,'年度权益礼包','年度会员可查看年度权益礼包说明。','/demo/svip-gift.png',1,40,1,NULL)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`description`=VALUES(`description`),`icon_url`=VALUES(`icon_url`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`version`=VALUES(`version`),`deleted_at`=NULL,`updated_at`=NOW();
INSERT INTO `qixi_crm_b_svip_plan` (`id`,`name`,`price`,`plan_type`,`duration_days`,`benefits`,`status`,`sort`) VALUES
  (980001,'SVIP 月度会员',29.00,'period',30,JSON_ARRAY('会员专享价','专属客服优先响应'),1,10),
  (980002,'SVIP 季度会员',79.00,'period',90,JSON_ARRAY('会员专享价','专属客服优先响应','每月权益提醒'),1,20),
  (980003,'SVIP 年度会员',299.00,'period',365,JSON_ARRAY('会员专享价','专属客服优先响应','年度权益礼包'),1,30)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`price`=VALUES(`price`),`plan_type`=VALUES(`plan_type`),`duration_days`=VALUES(`duration_days`),`benefits`=VALUES(`benefits`),`status`=VALUES(`status`),`sort`=VALUES(`sort`),`updated_at`=NOW();
-- SVIP 订单夹具只覆盖待支付与已关闭分支；不伪造任何支付成功、渠道流水或会员到账事实。
INSERT INTO `qixi_crm_b_svip_order` (`id`,`order_no`,`user_id`,`plan_id`,`plan_name`,`plan_type`,`duration_days`,`amount`,`status`,`idempotency_key`,`created_at`,`paid_at`) VALUES
  (980010,'SVIP-DEMO-20260804-001',9101,980001,'SVIP 月度会员','period',30,29.00,'pending','fixture-svip-pending-980010','2026-08-04 09:30:00',NULL),
  (980011,'SVIP-DEMO-20260804-002',9101,980002,'SVIP 季度会员','period',90,79.00,'closed','fixture-svip-closed-980011','2026-08-04 08:10:00',NULL)
ON DUPLICATE KEY UPDATE `plan_name`=VALUES(`plan_name`),`plan_type`=VALUES(`plan_type`),`duration_days`=VALUES(`duration_days`),`amount`=VALUES(`amount`),`status`=VALUES(`status`),`paid_at`=NULL;

-- 开屏广告中文夹具；图片地址可由运营替换，H5 会按 24 小时间隔展示。
INSERT INTO `qixi_crm_b_open_screen_campaign` (`id`,`title`,`image_url`,`link_url`,`duration_sec`,`space_hours`,`enabled`) VALUES
  (8801,'夏日焕新季','/demo/open-screen-summer.png','/pages/goods/list',3,24,1)
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`image_url`=VALUES(`image_url`),`link_url`=VALUES(`link_url`),`duration_sec`=VALUES(`duration_sec`),`space_hours`=VALUES(`space_hours`),`enabled`=VALUES(`enabled`);

-- 数据大屏演示数据见独立可重复导入脚本（360 单 + thumbs 封面）：
--   sql/business/init_data_screen_demo.sql
-- 本地初始化后请额外执行该文件，或使用 scripts 中的导入流程一并加载。
