-- 预售审核 Drawer 依赖 GET /products/:id/edit（读 qixi_crm_m_product）。
-- 业务库 qixi_crm_b_product_view 已有 1007/1008/110x，但商户源表缺失会导致「商品不存在」。
-- 幂等补齐预售演示关联商品（含详情/封面/轮播/SKU）。
USE `qixi_crm_merchant`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_m_product`
  (`id`,`store_id`,`title`,`category_id`,`store_category_id`,`brand_name`,`status`,`version`)
VALUES
  (1001,1,'轻奢羊绒针织衫',101,0,'云锦织造','on_sale',1),
  (1002,1,'头层牛皮通勤托特包',101,0,'栖木皮具','on_sale',1),
  (1004,1,'精纺圆领羊毛开衫',101,0,'云锦织造','on_sale',1),
  (1007,1,'柔软亲肤针织披肩',101,0,'云锦织造','on_sale',1),
  (1008,1,'城市通勤训练跑鞋',106,0,'逐风运动','on_sale',1),
  (1101,2,'无火藤条香氛礼盒',102,0,'澄日生活','on_sale',1),
  (1102,2,'晚安助眠香薰蜡烛',102,0,'澄日生活','on_sale',1),
  (1103,2,'恒温随行保温杯',103,0,'CRM Live精选','on_sale',1),
  (1104,2,'晨间居家香氛套装',102,0,'澄日生活','on_sale',1),
  (1107,2,'客厅氛围香薰礼盒',102,0,'澄日生活','on_sale',1),
  (1108,2,'织物护理香氛喷雾',102,0,'澄日生活','on_sale',1)
ON DUPLICATE KEY UPDATE
  `store_id`=VALUES(`store_id`),
  `title`=VALUES(`title`),
  `category_id`=VALUES(`category_id`),
  `brand_name`=VALUES(`brand_name`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_product_detail`
  (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`delivery_way`,`original_price`)
VALUES
  (1001,'轻薄保暖，通勤百搭羊绒针织。','羊绒,针织,通勤','件','/demo/product-knit-v1.png','2',399.00),
  (1002,'头层牛皮托特，容量适中适合通勤。','牛皮,托特包,通勤','件','/demo/product-bag-v1.png','2',599.00),
  (1004,'精纺圆领羊毛开衫，秋冬预售演示商品。','羊毛,开衫,预售','件','/demo/product-knit-v1.png','1,2',429.00),
  (1007,'柔软亲肤针织披肩，居家与出行两用。','披肩,针织','件','/demo/product-knit-v1.png','2',249.00),
  (1008,'城市通勤训练跑鞋，缓震透气。','跑鞋,通勤,运动','双','/demo/product-shoes-v1.png','2',529.00),
  (1101,'无火藤条香氛礼盒，客厅氛围演示商品。','香氛,藤条,礼盒','盒','/demo/product-fragrance-v1.png','2',299.00),
  (1102,'晚安助眠香薰蜡烛，助眠香调演示。','香薰,蜡烛,助眠','件','/demo/product-fragrance-v1.png','2',189.00),
  (1103,'恒温随行保温杯，出行随手一杯。','保温杯,随行','个','/demo/product-tumbler-v1.png','2',219.00),
  (1104,'晨间居家香氛套装，全款预售演示。','香氛,套装,居家','套','/demo/product-fragrance-v1.png','2',338.00),
  (1107,'客厅氛围香薰礼盒，定金预售演示。','香薰,礼盒,客厅','盒','/demo/product-fragrance-v1.png','2',369.00),
  (1108,'织物护理香氛喷雾，定金待审演示商品。','织物,护理,香氛喷雾','瓶','/demo/product-fragrance-v1.png','1,2',119.00)
ON DUPLICATE KEY UPDATE
  `brief`=VALUES(`brief`),
  `keyword`=VALUES(`keyword`),
  `unit_name`=VALUES(`unit_name`),
  `cover_url`=VALUES(`cover_url`),
  `delivery_way`=VALUES(`delivery_way`),
  `original_price`=VALUES(`original_price`);

-- 演示轮播：仅清理本补丁写入的 /demo/ 图，避免误删运营素材
DELETE FROM `qixi_crm_m_product_media`
WHERE `product_id` IN (1001,1002,1004,1007,1008,1101,1102,1103,1104,1107,1108)
  AND `media_type`='image'
  AND `url` LIKE '/demo/product-%';

INSERT INTO `qixi_crm_m_product_media` (`product_id`,`media_type`,`url`,`sort`) VALUES
  (1001,'image','/demo/product-knit-v1.png',0),
  (1001,'image','/demo/product-knit-v1.png',1),
  (1002,'image','/demo/product-bag-v1.png',0),
  (1002,'image','/demo/product-bag-v1.png',1),
  (1004,'image','/demo/product-knit-v1.png',0),
  (1004,'image','/demo/product-knit-v1.png',1),
  (1007,'image','/demo/product-knit-v1.png',0),
  (1008,'image','/demo/product-shoes-v1.png',0),
  (1008,'image','/demo/product-shoes-v1.png',1),
  (1101,'image','/demo/product-fragrance-v1.png',0),
  (1101,'image','/demo/product-fragrance-v1.png',1),
  (1102,'image','/demo/product-fragrance-v1.png',0),
  (1103,'image','/demo/product-tumbler-v1.png',0),
  (1104,'image','/demo/product-fragrance-v1.png',0),
  (1104,'image','/demo/product-fragrance-v1.png',1),
  (1107,'image','/demo/product-fragrance-v1.png',0),
  (1108,'image','/demo/product-fragrance-v1.png',0),
  (1108,'image','/demo/product-fragrance-v1.png',1);

INSERT INTO `qixi_crm_m_product_sku`
  (`id`,`product_id`,`spec_json`,`image`,`price`,`ot_price`,`stock`,`status`)
VALUES
  (61001,1001,JSON_OBJECT('默认','标准'),'/demo/product-knit-v1.png',299.00,399.00,60,1),
  (61002,1002,JSON_OBJECT('默认','标准'),'/demo/product-bag-v1.png',469.00,599.00,32,1),
  (61004,1004,JSON_OBJECT('默认','标准'),'/demo/product-knit-v1.png',329.00,429.00,36,1),
  (61008,1007,JSON_OBJECT('默认','标准'),'/demo/product-knit-v1.png',189.00,249.00,54,1),
  (61009,1008,JSON_OBJECT('默认','标准'),'/demo/product-shoes-v1.png',429.00,529.00,31,1),
  (61101,1101,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',239.00,299.00,72,1),
  (61102,1102,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',139.00,189.00,66,1),
  (61103,1103,JSON_OBJECT('默认','标准'),'/demo/product-tumbler-v1.png',159.00,219.00,80,1),
  (61104,1104,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',268.00,338.00,39,1),
  (61107,1107,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',299.00,369.00,26,1),
  (61108,1108,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',89.00,119.00,103,1)
ON DUPLICATE KEY UPDATE
  `product_id`=VALUES(`product_id`),
  `spec_json`=VALUES(`spec_json`),
  `image`=VALUES(`image`),
  `price`=VALUES(`price`),
  `ot_price`=VALUES(`ot_price`),
  `stock`=VALUES(`stock`),
  `status`=VALUES(`status`);
