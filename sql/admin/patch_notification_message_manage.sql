-- 消息管理：通知会员 / 通知店铺的默认行为与固定文本配置。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_notification_config` (
  `notification_id` bigint unsigned NOT NULL,
  `audience` enum('member','store') NOT NULL,
  `notice_type` varchar(100) NOT NULL,
  `scene` varchar(255) NOT NULL,
  `wechat_enabled` tinyint NOT NULL DEFAULT 0,
  `mini_program_enabled` tinyint NOT NULL DEFAULT 0,
  `sms_enabled` tinyint NOT NULL DEFAULT 0,
  `wechat_text` varchar(500) NOT NULL DEFAULT '',
  `mini_program_text` varchar(500) NOT NULL DEFAULT '',
  `sms_text` varchar(500) NOT NULL DEFAULT '',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`notification_id`),
  KEY `idx_notification_audience` (`audience`,`notification_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='默认通知行为与固定文本';

INSERT INTO `qixi_crm_a_menu`
  (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES
  (21910,1551,'setting.notice.config.manage','维护消息配置','','setting/notification/index','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` r
JOIN `qixi_crm_a_menu` m ON m.code IN ('setting.notice.dir','setting.notice.list','setting.notice.config.manage')
WHERE r.code='platform';
