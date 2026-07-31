USE `qixi_crm_merchant`;
SET NAMES utf8mb4;
-- 不初始化店铺账号、密码或任何真实手机号。仅提供无个人信息的店铺读写模型夹具。

INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`status`,`version`) VALUES
  (1001,1,'七禧臻选羊绒针织衫',101,'on_sale',1),
  (1002,1,'七禧头层牛皮通勤包',101,'on_sale',1),
  (1003,1,'七禧轻量跑步鞋',106,'on_sale',1),
  (1004,1,'七禧智能保温杯',103,'on_sale',1),
  (1005,1,'七禧真丝方巾',104,'on_sale',1),
  (1006,1,'七禧香氛礼盒',102,'on_sale',1)
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`),`category_id`=VALUES(`category_id`),`status`=VALUES(`status`),`version`=VALUES(`version`);

INSERT INTO `qixi_crm_m_product_sku` (`product_id`,`spec_json`,`price`,`stock`,`status`)
SELECT p.id, JSON_OBJECT('默认','标准'), CASE p.id
  WHEN 1001 THEN 299.00 WHEN 1002 THEN 469.00 WHEN 1003 THEN 369.00
  WHEN 1004 THEN 159.00 WHEN 1005 THEN 129.00 ELSE 239.00 END,
  50, 1
FROM `qixi_crm_m_product` AS p
WHERE p.id BETWEEN 1001 AND 1006
  AND NOT EXISTS (SELECT 1 FROM `qixi_crm_m_product_sku` AS s WHERE s.product_id = p.id);

INSERT INTO `qixi_crm_m_finance_ledger` (`store_id`,`entry_type`,`amount`,`reference_type`,`reference_id`,`idempotency_key`) VALUES
  (1,'test_seed',1000.00,'fixture','opening-balance','fixture-store-1-opening-balance')
ON DUPLICATE KEY UPDATE `amount`=VALUES(`amount`);

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
