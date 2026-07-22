-- 阶段 7：商户商品上下架 / 改库存按钮权限
USE `qixi_mergers`;

-- 挂在「商品列表」下
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 132, 103, 'product/show', '', '上下架', 'MerProductShowBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 132 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 133, 103, 'product/stock', '', '改库存', 'MerProductStockBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 133 AND `is_mer` = 2);

-- 商户主账号模板含双按钮
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',132,133')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%132%';

-- 演示：本店角色「商品运营」— 可看列表/编辑，可改库存，不可上下架
INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 5, '商品运营', '101,102,103,104,133', 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 5);

INSERT INTO `qixi_merchant_admin` (`merchant_admin_id`, `mer_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `level`, `status`, `is_del`)
SELECT 4, 1, 'merprod', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '商品运营', '13900000012', '5', 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_admin` WHERE `account` = 'merprod');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-product-button-perms', '阶段7：商品上下架/库存按钮 + merprod'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-product-button-perms');
