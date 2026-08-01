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

-- 预约服务公开活动夹具。排期余量来自 qixi_crm_b_reservation_booking，测试数据不使用真实用户身份。
INSERT INTO `qixi_crm_b_product_view` (`product_id`,`merchant_id`,`store_id`,`merchant_name`,`store_name`,`category_id`,`title`,`cover_url`,`price`,`original_price`,`product_type`,`sales`,`stock`,`sale_status`,`version`,`updated_at`) VALUES
  (1301,1,1,'七禧服饰商户','七禧服饰旗舰店',10101,'秋日衣橱搭配咨询','/demo/product-knit-v1.png',99.00,129.00,4,46,12,1,1,NOW()),
  (1302,2,2,'七禧居家商户','七禧居家优选店',10201,'居家香氛体验服务','/demo/product-fragrance-v1.png',129.00,169.00,4,38,10,1,1,NOW()),
  (1303,3,3,'七禧数码商户','七禧数码生活店',10301,'跑步鞋试穿与选购服务','/demo/product-shoes-v1.png',79.00,99.00,4,29,8,1,1,NOW())
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

-- PC/H5 演示消费者账号：仅用于本地与测试环境验证，不是生产账号。
-- 密码仅以 bcrypt 哈希保存；客户端、接口和日志不得返回凭据。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`auth_version`) VALUES
  (9101,'七禧体验用户','13500000001',1,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`mobile`=VALUES(`mobile`),`status`=VALUES(`status`),`auth_version`=VALUES(`auth_version`);

INSERT INTO `qixi_crm_b_user_identity` (`user_id`,`channel`,`subject`,`credential_hash`) VALUES
  (9101,'pc','13500000001','$2y$10$qX7mBTgn9Fh5QQ0LGhBW3OU6PFimO83brQEv2St6YsDmFSKGZgXrK'),
  (9101,'h5','13500000001','$2y$10$qX7mBTgn9Fh5QQ0LGhBW3OU6PFimO83brQEv2St6YsDmFSKGZgXrK')
ON DUPLICATE KEY UPDATE `user_id`=VALUES(`user_id`),`credential_hash`=VALUES(`credential_hash`);

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

-- 秒杀展示夹具。规则由后台营销活动投影而来；C 端只读 qixi_crm_b_marketing_activity_view。
INSERT INTO `qixi_crm_b_marketing_activity_view` (`activity_id`,`store_id`,`activity_type`,`name`,`rules`,`status`,`version`,`starts_at`,`ends_at`) VALUES
  (5001,1,'seckill','轻奢羊绒针织衫限时抢购',JSON_OBJECT('product_id',1001,'seckill_price',199.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5002,1,'seckill','头层牛皮托特包限时抢购',JSON_OBJECT('product_id',1002,'seckill_price',329.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5003,2,'seckill','无火藤条香氛礼盒限时抢购',JSON_OBJECT('product_id',1101,'seckill_price',169.00,'time_slots',JSON_ARRAY('00:00','14:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY)),
  (5004,3,'seckill','智能数显保温杯限时抢购',JSON_OBJECT('product_id',1201,'seckill_price',149.00,'time_slots',JSON_ARRAY('07:00','19:00')),1,1,DATE_SUB(NOW(),INTERVAL 1 DAY),DATE_ADD(NOW(),INTERVAL 180 DAY))
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`name`=VALUES(`name`),`rules`=VALUES(`rules`),`status`=VALUES(`status`),`version`=VALUES(`version`),`starts_at`=VALUES(`starts_at`),`ends_at`=VALUES(`ends_at`);

-- 直播公开读模型夹具。直播流地址仅为本地演示标识，客户端不会把它当作生产推流密钥。
INSERT INTO `qixi_crm_b_live_room` (`id`,`merchant_id`,`store_id`,`anchor_user_id`,`title`,`anchor_name`,`cover_url`,`status`,`is_public`,`stream_ref`,`play_url`,`starts_at`,`ended_at`,`sort`) VALUES
  (7001,1,1,NULL,'七禧秋日衣橱直播专场','七禧小夏','/demo/home-hero-accessories-v1.png','living',1,'local-demo-fashion','',DATE_SUB(NOW(), INTERVAL 20 MINUTE),NULL,20),
  (7002,2,2,NULL,'居家香氛与生活好物分享','居家优选主播','/demo/home-hero-fragrance-v1.png','scheduled',1,'local-demo-home','',DATE_ADD(NOW(), INTERVAL 1 HOUR),NULL,10),
  (7003,3,3,NULL,'通勤数码好物直播回放','数码生活主播','/demo/home-tech-wide-v1.png','ended',1,'local-demo-digital','',DATE_SUB(NOW(), INTERVAL 2 DAY),DATE_SUB(NOW(), INTERVAL 1 DAY),5)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`store_id`=VALUES(`store_id`),`title`=VALUES(`title`),`anchor_name`=VALUES(`anchor_name`),`cover_url`=VALUES(`cover_url`),`status`=VALUES(`status`),`is_public`=VALUES(`is_public`),`stream_ref`=VALUES(`stream_ref`),`play_url`=VALUES(`play_url`),`starts_at`=VALUES(`starts_at`),`ended_at`=VALUES(`ended_at`),`sort`=VALUES(`sort`);

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
    'items',JSON_ARRAY(
      JSON_OBJECT('type','banner','name','轮播图','data',JSON_ARRAY(
        JSON_OBJECT('imgName','七禧商城精选','imgUrl','/demo/home-hero-v1.png','linkUrl','/goods?cate_id=101'),
        JSON_OBJECT('imgName','七禧香氛家居','imgUrl','/demo/home-hero-fragrance-v1.png','linkUrl','/goods?cate_id=102'),
        JSON_OBJECT('imgName','七禧箱包配饰','imgUrl','/demo/home-hero-accessories-v1.png','linkUrl','/goods?cate_id=10102')
      )),
      JSON_OBJECT('type','product','name','服饰鞋包展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',101,'showNum',4,'productSort','sales')
      )),
      JSON_OBJECT('type','product','name','家居生活展示类型','params',JSON_OBJECT(
        'source','auto','auto',JSON_OBJECT('category',102,'showNum',4,'productSort','sales')
      ))
    ),
    '_qixi',JSON_OBJECT('title','七禧商城','template_name','home','is_diy',1,'is_show',1,'is_default',1)
  ),'published',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`document`=VALUES(`document`),`status`=VALUES(`status`),`is_active`=VALUES(`is_active`),`updated_at`=NOW();
