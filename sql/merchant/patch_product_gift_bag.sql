-- 分销礼包标记（对齐 CRMEB Product.is_gift_bag）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_m_product' AND COLUMN_NAME='is_gift_bag')=0,
    'ALTER TABLE `qixi_crm_m_product` ADD COLUMN `is_gift_bag` tinyint NOT NULL DEFAULT 0 COMMENT ''1分销礼包'' AFTER `status`, ADD KEY `idx_gift_bag_status` (`is_gift_bag`,`status`)',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 本地演示：3 条待审核分销礼包（虚构商品，无真实手机号/收款资料）
INSERT INTO `qixi_crm_m_product`
  (`id`,`store_id`,`title`,`category_id`,`store_category_id`,`brand_name`,`status`,`is_gift_bag`,`version`)
VALUES
  (1901,1,'分销新人礼包·尝鲜装',101,0,'七禧分销','pending_review',1,1),
  (1902,1,'分销进阶礼包·权益包',101,0,'七禧分销','pending_review',1,1),
  (1903,1,'分销尊享礼包·年度装',101,0,'七禧分销','pending_review',1,1)
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `status`=VALUES(`status`),
  `is_gift_bag`=VALUES(`is_gift_bag`),
  `brand_name`=VALUES(`brand_name`);

INSERT INTO `qixi_crm_m_product_detail`
  (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`original_price`)
VALUES
  (1901,'分销新人入门礼包','分销,礼包,新人','套','',199.00),
  (1902,'分销进阶权益礼包','分销,礼包,进阶','套','',399.00),
  (1903,'分销尊享年度礼包','分销,礼包,尊享','套','',999.00)
ON DUPLICATE KEY UPDATE
  `brief`=VALUES(`brief`),
  `keyword`=VALUES(`keyword`),
  `unit_name`=VALUES(`unit_name`),
  `original_price`=VALUES(`original_price`);

INSERT INTO `qixi_crm_m_product_sku`
  (`id`,`product_id`,`spec_json`,`price`,`stock`,`status`)
VALUES
  (61901,1901,JSON_OBJECT('默认','标准'),99.00,100,1),
  (61902,1902,JSON_OBJECT('默认','标准'),199.00,80,1),
  (61903,1903,JSON_OBJECT('默认','标准'),499.00,50,1)
ON DUPLICATE KEY UPDATE
  `product_id`=VALUES(`product_id`),
  `spec_json`=VALUES(`spec_json`),
  `price`=VALUES(`price`),
  `stock`=VALUES(`stock`),
  `status`=VALUES(`status`);
