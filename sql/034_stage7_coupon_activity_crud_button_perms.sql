-- 阶段 7：券删除 + 秒杀/拼团/预售/助力 创建·删除按钮
USE `qixi_mergers`;

-- 商户：券删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 142, 116, 'coupon/delete', '', '删除优惠券', 'MerCouponDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 142 AND `is_mer` = 2);

-- 商户：秒杀创建/删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 143, 117, 'seckill/create', '', '新建秒杀', 'MerSeckillCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 143 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 144, 117, 'seckill/delete', '', '删除秒杀', 'MerSeckillDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 144 AND `is_mer` = 2);

-- 商户：拼团创建/删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 145, 118, 'combination/create', '', '新建拼团', 'MerCombinationCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 145 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 146, 118, 'combination/delete', '', '删除拼团', 'MerCombinationDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 146 AND `is_mer` = 2);

-- 商户：预售创建/删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 147, 121, 'presell/create', '', '新建预售', 'MerPresellCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 147 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 148, 121, 'presell/delete', '', '删除预售', 'MerPresellDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 148 AND `is_mer` = 2);

-- 商户：助力创建/删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 149, 124, 'assist/create', '', '新建助力', 'MerAssistCreateBtn', '', 4, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 149 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 150, 124, 'assist/delete', '', '删除助力', 'MerAssistDeleteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 150 AND `is_mer` = 2);

-- 平台：券删除
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 43, 18, 'coupon/delete', '', '删除优惠券', 'MarketingCouponDeleteBtn', '', 1, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 43 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',142,143,144,145,146,147,148,149,150')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%142%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',43')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%43%';

-- 若已执行过旧版 034（仅到 146），补预售/助力
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',147,148,149,150')
WHERE `role_id` = 2 AND `rules` LIKE '%146%' AND `rules` NOT LIKE '%147%';

-- meract：仅秒杀/券启停，不赋创建删除（不对称）
-- auditor：仅券启停，不赋删除

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-coupon-activity-crud-btns', '阶段7：券删除 + 秒杀/拼团/预售/助力创建删除'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-coupon-activity-crud-btns');
