-- 阶段 7：平台/商户优惠券启停按钮
USE `qixi_mergers`;

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 138, 116, 'coupon/toggle', '', '优惠券启停', 'MerCouponToggleBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 138 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 41, 18, 'coupon/toggle', '', '优惠券启停', 'MarketingCouponToggleBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 41 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',138')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%138%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',41')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%41%';

-- meract：店铺券列表 + 启停（有秒杀启停，无拼团启停；此处补券）
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',116,138')
WHERE `role_id` = 6 AND `role_name` = '营销运营' AND `rules` NOT LIKE '%138%';

-- auditor：平台券列表 + 启停
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',18,41')
WHERE `role_id` = 4 AND `role_name` = '平台运营' AND `rules` NOT LIKE '%41%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-coupon-button-perms', '阶段7：平台/商户优惠券启停按钮'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-coupon-button-perms');
