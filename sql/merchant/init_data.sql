SET NAMES utf8mb4;
USE `qixi_crm_merchant`;
INSERT INTO `qixi_crm_m_merchant` (`id`,`name`,`status`,`region_id`) VALUES
  (1,'七禧演示茶铺',1,10),(2,'七禧居家优选店',1,20),(3,'七禧数码生活店',1,10)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`),`region_id`=VALUES(`region_id`);
INSERT INTO `qixi_crm_m_store` (`id`,`merchant_id`,`app_id`,`name`,`status`) VALUES
  (1,1,'qixi.store.demo.1','七禧演示茶铺',1),(2,2,'qixi.store.demo.2','七禧居家优选店',1),(3,3,'qixi.store.demo.3','七禧数码生活店',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`app_id`=VALUES(`app_id`),`name`=VALUES(`name`),`status`=VALUES(`status`);

-- 店铺后台 owner：平台「店铺列表 → 登录」依赖 qixi_crm_m_account；密码演示值 123456m（仅本地夹具）。
INSERT INTO `qixi_crm_m_account` (`id`,`store_id`,`username`,`password_hash`,`role_code`,`display_name`,`phone`,`status`,`auth_version`) VALUES
  (1,1,'demo_store_1','$2a$10$7e1OmptO8l5P3lJ7ziIfOeC0GXY0MGUNNY/QS6LQKgXLNq2Z6TFJe','owner','演示店长','13900000000',1,1),
  (2,2,'demo_store_2','$2a$10$7e1OmptO8l5P3lJ7ziIfOeC0GXY0MGUNNY/QS6LQKgXLNq2Z6TFJe','owner','居家店长','13900000001',1,1),
  (3,3,'demo_store_3','$2a$10$7e1OmptO8l5P3lJ7ziIfOeC0GXY0MGUNNY/QS6LQKgXLNq2Z6TFJe','owner','数码店长','13900000002',1,1)
ON DUPLICATE KEY UPDATE
  `store_id`=VALUES(`store_id`),`username`=VALUES(`username`),`password_hash`=VALUES(`password_hash`),
  `role_code`=VALUES(`role_code`),`display_name`=VALUES(`display_name`),`phone`=VALUES(`phone`),
  `status`=VALUES(`status`);

-- 平台商品监管中文夹具：区域 10 对应商户 1、3，区域 20 对应商户 2；不含真实个人或商户资料。
INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`status`,`version`) VALUES
  (5001,1,'七禧春日茉莉花茶礼盒',101,'on_sale',1),
  (5002,2,'居家香氛扩香礼盒',102,'pending_review',1),
  (5003,3,'智能温控随行杯',103,'draft',1)
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`title`=VALUES(`title`),`category_id`=VALUES(`category_id`),`status`=VALUES(`status`),`version`=VALUES(`version`);
INSERT INTO `qixi_crm_m_product_detail` (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`original_price`) VALUES
  (5001,'虚构中文商品：清香茉莉花茶礼盒，供平台商品范围验收。','茉莉花茶,礼盒','盒','/demo/product-tea-v1.png',129.00),
  (5002,'虚构中文商品：居家香氛礼盒，处于待审核状态。','香氛,家居','盒','/demo/product-fragrance-v1.png',199.00),
  (5003,'虚构中文商品：智能温控随行杯草稿。','温控杯,数码','个','/demo/product-tumbler-v1.png',259.00)
ON DUPLICATE KEY UPDATE `brief`=VALUES(`brief`),`keyword`=VALUES(`keyword`),`unit_name`=VALUES(`unit_name`),`cover_url`=VALUES(`cover_url`),`original_price`=VALUES(`original_price`);
INSERT INTO `qixi_crm_m_product_sku` (`id`,`product_id`,`spec_json`,`price`,`stock`,`status`) VALUES
  (5101,5001,JSON_OBJECT('规格','250克'),99.00,60,1),(5102,5002,JSON_OBJECT('规格','礼盒装'),159.00,36,1),(5103,5003,JSON_OBJECT('颜色','深空灰'),199.00,48,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`spec_json`=VALUES(`spec_json`),`price`=VALUES(`price`),`stock`=VALUES(`stock`),`status`=VALUES(`status`);

-- 店铺菜单（qixi_crm_m_menu / role_menu）由 init_menu_crmeb_full.sql 全量导入：
-- 源：CRMEB 线上 GET /sys/merchant/menu/lst（https://mer.crmeb.net/admin/merchant/system），共 711 条。
-- 由 scripts/release/db-reset.sh 与 scripts/qixi-crm.sh db-init 在 merchant/init_data 之后执行。
