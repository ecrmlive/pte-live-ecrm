-- 阶段 7：平台/商户优惠券「新建」按钮（启停见 031）
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 141, 116, 'coupon/create', '', '新建优惠券', 'MerCouponCreateBtn', '', 3, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 141 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 42, 18, 'coupon/create', '', '新建优惠券', 'MarketingCouponCreateBtn', '', 3, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 42 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',141')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%141%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',42')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%42%';

-- meract：可启停店铺券，不可新建（不对称）
-- auditor：可启停平台券，不可新建（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-coupon-create-btns', '阶段7：平台/商户新建优惠券按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-coupon-create-btns');
