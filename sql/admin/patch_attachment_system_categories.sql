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
  (5108,0,'','其他图片','other_image',20,0,1,NOW()),
  (5111,0,'','店铺视频','store_video',19,0,1,NOW()),
  (5112,0,'','商品视频','product_video',18,0,1,NOW()),
  (5113,0,'','其他视频','other_video',17,0,1,NOW())
ON DUPLICATE KEY UPDATE
  `attachment_category_name`=VALUES(`attachment_category_name`),
  `attachment_category_enname`=VALUES(`attachment_category_enname`),
  `sort`=VALUES(`sort`),
  `mer_id`=VALUES(`mer_id`),
  `is_system`=VALUES(`is_system`);

-- 素材行级 is_system：侧栏「系统素材」只查 is_system=1，与「挂在系统分类下的运营图」区分
SET @asset_col := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA='qixi_crm_admin'
    AND TABLE_NAME='qixi_crm_a_attachment_asset'
    AND COLUMN_NAME='is_system'
);
SET @asset_ddl := IF(
  @asset_col=0,
  'ALTER TABLE `qixi_crm_a_attachment_asset` ADD COLUMN `is_system` tinyint NOT NULL DEFAULT 0 COMMENT ''1=系统预置素材'' AFTER `attachment_type`, ADD KEY `idx_owner_system` (`user_type`,`is_system`,`attachment_type`)',
  'SELECT 1'
);
PREPARE asset_stmt FROM @asset_ddl;
EXECUTE asset_stmt;
DEALLOCATE PREPARE asset_stmt;

-- 历史运营/演示风景图：可留在背景图片等系统分类，但不得标为系统预置素材
UPDATE `qixi_crm_a_attachment_asset`
   SET `is_system`=0
 WHERE `attachment_id` IN (5313,5314,5315)
    OR `attachment_name` IN ('bg_1.png','bg_4.png','bg_6.png');

DELETE FROM `qixi_crm_a_attachment_category`
 WHERE `mer_id`=0 AND `attachment_category_enname` IN ('demo-image','demo-video');
