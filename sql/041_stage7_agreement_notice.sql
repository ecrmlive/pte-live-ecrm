-- 阶段 7：协议规则（qixi_cache）+ C 端公告闭环配套
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_cache` (
  `key` varchar(32) NOT NULL COMMENT '协议/缓存键',
  `expire_time` int(11) NOT NULL DEFAULT '0' COMMENT '0=永久',
  `result` longtext NOT NULL COMMENT '缓存数据/协议正文',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '缓存时间',
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='协议与缓存';

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_user_agree', 0, '【用户协议】欢迎使用栖息多商户商城。使用本服务即表示您同意遵守平台交易、支付与售后规则。本页为演示文案，非法律意见。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_user_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_userr_privacy', 0, '【隐私政策】我们仅在提供服务所必需的范围内处理账号、订单与设备信息，并采取合理安全措施保护您的数据。演示环境不采集真实敏感信息。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_userr_privacy');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_svip', 0, '【付费会员协议】开通 SVIP 后可享受会员价等权益；权益以店铺配置为准。演示账号可在后台调整文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_svip');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_product_presell_agree', 0, '【预售协议】预售商品按活动规则支付定金/全款；定金预售需在尾款期内支付尾款，逾期规则以活动说明为准。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_product_presell_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'business_entry_agree', 0, '【商户入驻协议】申请入驻需提交真实店铺信息，审核通过后可开展经营；平台有权对违规店铺采取限制措施。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'business_entry_agree');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'promoter_explain', 0, '【分销说明】推广员可分享商品获得佣金；佣金结算以订单完成且无售后为准。演示环境仅展示说明文案。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'promoter_explain');

INSERT INTO `qixi_cache` (`key`, `expire_time`, `result`)
SELECT 'sys_about_us', 0, '【关于我们】栖息多商户商城 · qixi-live-mergers。功能基线对齐 CRMEB Merchant v4.0，技术栈为 Go + Vben + uni-app。'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_cache` WHERE `key` = 'sys_about_us');

-- 平台：协议规则页 + 保存按钮
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 52, 10, '/setting/agreements', '', '协议规则', 'SettingAgreements', '', 6, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 52 AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 53, 52, 'agreement/update', '', '保存协议', 'SettingAgreementUpdateBtn', '', 1, 0, 1, 2, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 53 AND `is_mer` = 1);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',52,53')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%52%';

-- auditor：不赋协议写（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-agreement-notice', '阶段7：协议规则 qixi_cache + 菜单；公告 C 端闭环'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-agreement-notice');
