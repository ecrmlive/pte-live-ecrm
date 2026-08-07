-- 平台素材库：补齐系统分类字段与固定分类种子（可重复执行）
-- 系统分类用途：客户端（H5/小程序/App）与装修页常用的图标、图片、背景
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @col_exists := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA='qixi_crm_admin'
    AND TABLE_NAME='qixi_crm_a_attachment_category'
    AND COLUMN_NAME='is_system'
);
SET @ddl := IF(
  @col_exists=0,
  'ALTER TABLE `qixi_crm_a_attachment_category` ADD COLUMN `is_system` tinyint NOT NULL DEFAULT 0 AFTER `mer_id`',
  'SELECT 1'
);
PREPARE stmt FROM @ddl;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

INSERT INTO `qixi_crm_a_attachment_category`
  (`attachment_category_id`,`pid`,`path`,`attachment_category_name`,`attachment_category_enname`,`sort`,`mer_id`,`is_system`,`create_time`)
VALUES
  (5101,0,'','店铺封面','store_cover',90,0,1,NOW()),
  (5102,0,'','支付图标','pay_icon',80,0,1,NOW()),
  (5103,0,'','物流图标','logistics_icon',70,0,1,NOW()),
  (5104,0,'','客服图标','service_icon',60,0,1,NOW()),
  (5105,0,'','商品图片','product_image',50,0,1,NOW()),
  (5106,0,'','背景图片','background_image',40,0,1,NOW()),
  (5107,0,'','列表图标','list_icon',30,0,1,NOW()),
  (5108,0,'','其他图片','other_image',20,0,1,NOW())
ON DUPLICATE KEY UPDATE
  `attachment_category_name`=VALUES(`attachment_category_name`),
  `attachment_category_enname`=VALUES(`attachment_category_enname`),
  `sort`=VALUES(`sort`),
  `mer_id`=VALUES(`mer_id`),
  `is_system`=VALUES(`is_system`);

UPDATE `qixi_crm_a_attachment_asset`
   SET `attachment_category_id`=5101
 WHERE `attachment_id`=5311;
UPDATE `qixi_crm_a_attachment_asset`
   SET `attachment_category_id`=5108
 WHERE `attachment_id`=5312;
DELETE FROM `qixi_crm_a_attachment_category`
 WHERE `mer_id`=0 AND `attachment_category_enname` IN ('demo-image','demo-video');
