-- 阶段 6：SVIP 会员价 + 与店铺券互斥（FUNCTIONAL-TRUTH §8.2）
USE `qixi_mergers`;

-- 商品 SVIP 字段
SET @sql := (
  SELECT IF(
    EXISTS (
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_product' AND COLUMN_NAME = 'svip_price_type'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_store_product`
       ADD COLUMN `svip_price_type` tinyint(3) unsigned NOT NULL DEFAULT 0 COMMENT ''0不参加 1默认比例 2自定义'' AFTER `ot_price`,
       ADD COLUMN `svip_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT ''会员价'' AFTER `svip_price_type`,
       ADD COLUMN `mer_svip_status` tinyint(3) unsigned NOT NULL DEFAULT 1 COMMENT ''商户会员状态'' AFTER `svip_price`'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 商户：SVIP 价与店铺券是否可叠加（1=可叠加）
SET @sql := (
  SELECT IF(
    EXISTS (
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_merchant' AND COLUMN_NAME = 'svip_coupon_merge'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_merchant`
       ADD COLUMN `svip_coupon_merge` tinyint(1) NOT NULL DEFAULT 0 COMMENT ''1=SVIP价可叠店铺券'' AFTER `mer_integral_rate`'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 子单 SVIP 优惠金额
SET @sql := (
  SELECT IF(
    EXISTS (
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'svip_discount'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_store_order`
       ADD COLUMN `svip_discount` decimal(8,2) NOT NULL DEFAULT 0.00 COMMENT ''svip优惠金额'' AFTER `platform_coupon_price`'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 演示：商品1 自定义会员价 19.90；SKU 同步
UPDATE `qixi_store_product`
SET `svip_price_type` = 2, `svip_price` = 19.90, `mer_svip_status` = 1
WHERE `product_id` = 1;

UPDATE `qixi_store_product_attr_value`
SET `svip_price` = 19.90
WHERE `product_id` = 1 AND `value_id` = 1;

-- 演示用户开通永久 SVIP
UPDATE `qixi_user`
SET `is_svip` = 3, `svip_endtime` = NULL
WHERE `uid` = 1;

-- 商户1 默认不可叠店铺券（验互斥）
UPDATE `qixi_merchant` SET `svip_coupon_merge` = 0 WHERE `mer_id` = 1;

-- 平台菜单：付费会员
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 25, 0, '/user/svip', 'CrownOutlined', '付费会员', 'UserSvip', '', 52, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 25 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',25')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%25%';

-- 商户菜单：SVIP 与店铺券叠加
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 120, 110, '/setting/svip', '', '会员价叠加', 'MerSettingSvip', '', 7, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 120 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',120')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%120%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6-svip', '阶段6：SVIP会员价与店铺券互斥'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6-svip');
