-- 阶段 7：客服快捷回复落库 + 商户维护菜单/按钮
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_service_reply` (
  `service_reply_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `mer_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '商户 id',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '1:文字 2:图片',
  `keyword` varchar(64) NOT NULL COMMENT '回复的关键字/标题',
  `content` varchar(512) NOT NULL COMMENT '回复内容',
  `status` tinyint(3) unsigned DEFAULT '1' COMMENT '是否开启',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`service_reply_id`) USING BTREE,
  KEY `mer_id` (`mer_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='客服快捷回复';

INSERT INTO `qixi_store_service_reply` (`service_reply_id`, `mer_id`, `type`, `keyword`, `content`, `status`)
SELECT 1, 1, 1, '问候', '您好，请问有什么可以帮您？', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service_reply` WHERE `service_reply_id` = 1);

INSERT INTO `qixi_store_service_reply` (`service_reply_id`, `mer_id`, `type`, `keyword`, `content`, `status`)
SELECT 2, 1, 1, '查单中', '已为您查询订单，请稍等片刻。', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service_reply` WHERE `service_reply_id` = 2);

INSERT INTO `qixi_store_service_reply` (`service_reply_id`, `mer_id`, `type`, `keyword`, `content`, `status`)
SELECT 3, 1, 1, '发货说明', '您的订单已发货，可在「我的订单」查看物流。', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service_reply` WHERE `service_reply_id` = 3);

INSERT INTO `qixi_store_service_reply` (`service_reply_id`, `mer_id`, `type`, `keyword`, `content`, `status`)
SELECT 4, 1, 1, '售后引导', '如需售后，请在订单详情中申请退款，我们会尽快处理。', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service_reply` WHERE `service_reply_id` = 4);

INSERT INTO `qixi_store_service_reply` (`service_reply_id`, `mer_id`, `type`, `keyword`, `content`, `status`)
SELECT 5, 1, 1, '结束语', '感谢您的耐心，祝您购物愉快！', 1
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service_reply` WHERE `service_reply_id` = 5);

-- 商户设置：快捷回复页 + 写按钮
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 163, 110, '/setting/replies', '', '快捷回复', 'MerSettingReplies', '', 4, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 163 AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 164, 163, 'reply/write', '', '维护快捷回复', 'MerReplyWriteBtn', '', 1, 0, 2, 2, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `menu_id` = 164 AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',163,164')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%163%';

-- meract / mersub：不赋快捷回复页/写（不对称）

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase7-service-reply', '阶段7：客服快捷回复落库 + 商户维护'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase7-service-reply');
