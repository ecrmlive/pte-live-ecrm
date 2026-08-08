-- 平台保障服务：店铺/商品使用数计数列（对齐 CRMEB eb_guarantee.mer_count / product_cout）
-- 商户保障模板与商品关联尚未落地时，计数保持 0，供后续定时/事件回填。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_guarantee' AND COLUMN_NAME='mer_count')=0,
    'ALTER TABLE `qixi_crm_a_product_guarantee` ADD COLUMN `mer_count` int NOT NULL DEFAULT 0 COMMENT ''使用的店铺数'' AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_guarantee' AND COLUMN_NAME='product_count')=0,
    'ALTER TABLE `qixi_crm_a_product_guarantee` ADD COLUMN `product_count` int NOT NULL DEFAULT 0 COMMENT ''使用商品数'' AFTER `mer_count`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 本地演示数据（utf8mb4 中文；图标留空，验收时从素材库选择）
INSERT INTO `qixi_crm_a_product_guarantee` (`id`,`name`,`content`,`icon_url`,`sort`,`status`,`mer_count`,`product_count`) VALUES
  (7511,'正品保障','平台演示：商品来源与售后承诺以订单规则为准。','',100,1,0,0),
  (7512,'极速退款','极速退款服务，符合规则的退款优先进入售后状态机处理。','',90,1,0,0),
  (7513,'15天价保','商品自签收之日起15天内，如出现官方降价，可申请价格保护差价退还。','',80,1,12,86),
  (7514,'7天无理由退货','商品自签收之日起7天内，在不影响二次销售的前提下可申请无理由退货。','',70,1,28,156),
  (7515,'退货宝','退货流程协助与运费补贴说明（演示数据，以实际售后规则为准）。','',60,1,9,42)
ON DUPLICATE KEY UPDATE
  `content`=VALUES(`content`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `mer_count`=VALUES(`mer_count`),
  `product_count`=VALUES(`product_count`);
