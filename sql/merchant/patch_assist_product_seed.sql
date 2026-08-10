-- 助力「活动商品」编辑 / 编辑标签依赖 GET /products/:id/edit（读 qixi_crm_m_product）。
-- 业务库 qixi_crm_b_assist 关联 1005/1107/1207，商户源表缺失或不完整会导致「商品不存在」。
-- 幂等补齐助力演示关联商品（含详情/封面/轮播/SKU）。
USE `qixi_crm_merchant`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_m_product`
  (`id`,`store_id`,`title`,`category_id`,`store_category_id`,`brand_name`,`status`,`version`)
VALUES
  (1005,1,'真丝印花方巾礼盒',104,0,'云锦织造','on_sale',1),
  (1107,2,'客厅氛围香薰礼盒',102,0,'澄日生活','on_sale',1),
  (1207,3,'桌面恒温杯垫礼盒',103,0,'CRM Live精选','on_sale',1)
ON DUPLICATE KEY UPDATE
  `store_id`=VALUES(`store_id`),
  `title`=VALUES(`title`),
  `category_id`=VALUES(`category_id`),
  `brand_name`=VALUES(`brand_name`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_product_detail`
  (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`delivery_way`,`original_price`)
VALUES
  (1005,'真丝印花方巾礼盒，好友助力演示商品。','真丝,方巾,礼盒','盒','/demo/product-scarf-v1.png','2',169.00),
  (1107,'客厅氛围香薰礼盒，好友助力演示商品。','香薰,礼盒,客厅','盒','/demo/product-fragrance-v1.png','2',369.00),
  (1207,'桌面恒温杯垫礼盒，好友助力演示商品。','恒温,杯垫,礼盒','盒','/demo/product-tumbler-v1.png','2',279.00)
ON DUPLICATE KEY UPDATE
  `brief`=VALUES(`brief`),
  `keyword`=VALUES(`keyword`),
  `unit_name`=VALUES(`unit_name`),
  `cover_url`=VALUES(`cover_url`),
  `delivery_way`=VALUES(`delivery_way`),
  `original_price`=VALUES(`original_price`);

-- 演示轮播：仅清理本补丁写入的 /demo/ 图，避免误删运营素材
DELETE FROM `qixi_crm_m_product_media`
WHERE `product_id` IN (1005,1107,1207)
  AND `media_type`='image'
  AND `url` LIKE '/demo/product-%';

INSERT INTO `qixi_crm_m_product_media` (`product_id`,`media_type`,`url`,`sort`) VALUES
  (1005,'image','/demo/product-scarf-v1.png',0),
  (1005,'image','/demo/product-scarf-v1.png',1),
  (1107,'image','/demo/product-fragrance-v1.png',0),
  (1107,'image','/demo/product-fragrance-v1.png',1),
  (1207,'image','/demo/product-tumbler-v1.png',0),
  (1207,'image','/demo/product-tumbler-v1.png',1);

INSERT INTO `qixi_crm_m_product_sku`
  (`id`,`product_id`,`spec_json`,`image`,`price`,`ot_price`,`stock`,`status`)
VALUES
  (61005,1005,JSON_OBJECT('默认','标准'),'/demo/product-scarf-v1.png',129.00,169.00,90,1),
  (61107,1107,JSON_OBJECT('默认','标准'),'/demo/product-fragrance-v1.png',299.00,369.00,26,1),
  (61207,1207,JSON_OBJECT('默认','标准'),'/demo/product-tumbler-v1.png',219.00,279.00,38,1)
ON DUPLICATE KEY UPDATE
  `product_id`=VALUES(`product_id`),
  `spec_json`=VALUES(`spec_json`),
  `image`=VALUES(`image`),
  `price`=VALUES(`price`),
  `ot_price`=VALUES(`ot_price`),
  `stock`=VALUES(`stock`),
  `status`=VALUES(`status`);
