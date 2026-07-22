-- 阶段 7：秒杀/拼团启停 + 直播审房按钮
USE `qixi_mergers`;

-- 商户：秒杀启停、拼团上下架
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 136, 117, 'seckill/toggle', '', '秒杀启停', 'MerSeckillToggleBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 136 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 137, 118, 'combination/toggle', '', '拼团上下架', 'MerCombinationToggleBtn', '', 2, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 137 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',136,137')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%136%';

-- 演示：营销运营（秒杀列表+启停，无拼团启停）
INSERT INTO `qixi_system_role` (`role_id`, `role_name`, `rules`, `status`, `mer_id`, `is_agent`, `is_default`)
SELECT 6, '营销运营', '101,115,117,136', 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_role` WHERE `role_id` = 6);

INSERT INTO `qixi_merchant_admin` (`merchant_admin_id`, `mer_id`, `account`, `pwd`, `real_name`, `phone`, `roles`, `level`, `status`, `is_del`)
SELECT 5, 1, 'meract', '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG', '营销运营', '13900000013', '6', 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_merchant_admin` WHERE `account` = 'meract');

-- 平台：直播审房
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 40, 27, 'broadcast/audit', '', '直播审房', 'MarketingBroadcastAuditBtn', '', 2, 1, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 40 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',40')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%40%';

-- auditor：直播监管 + 审房（可审；与删帖不对称一致）
UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',17,27,40')
WHERE `role_id` = 4 AND `role_name` = '平台运营' AND `rules` NOT LIKE '%40%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-marketing-broadcast-btns', '阶段7：秒杀/拼团启停 + 直播审房按钮 + meract'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-marketing-broadcast-btns');
