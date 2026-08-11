-- 删除「微信模板消息」菜单；补齐自动回复读写权限
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 软隐藏：公众号 → 微信模板消息（含子权限按钮）
UPDATE `qixi_crm_a_menu`
SET `status`=0, `title`='微信模板消息', `sort`=99
WHERE `id`=205 OR `code`='app.wechat_template';

UPDATE `qixi_crm_a_menu`
SET `status`=0
WHERE `id`=20996 OR `code`='app.wechat_template.manage';

-- 自动回复读权限
INSERT INTO `qixi_crm_a_menu`
  (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES
  (21121,203,'app.wechat_reply.read','查看自动回复','','admin/app/wechat/reply','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu`
SET `title`='维护自动回复', `sort`=2, `status`=1
WHERE `id`=20994 OR `code`='app.wechat_reply.manage';

INSERT INTO `qixi_crm_a_role_menu` (`role_id`, `menu_id`)
SELECT r.id, 21121
FROM `qixi_crm_a_role` r
WHERE r.code IN ('platform', 'operations')
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_a_role_menu` rm
    WHERE rm.role_id = r.id AND rm.menu_id = 21121
  );

-- 自动回复表
CREATE TABLE IF NOT EXISTS `qixi_crm_a_wechat_reply` (
  `wechat_reply_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '回复规则 ID',
  `reply_key` varchar(64) NOT NULL COMMENT '关键字；subscribe=关注欢迎 default=默认回复',
  `reply_type` varchar(32) NOT NULL DEFAULT 'text' COMMENT 'text/image/news',
  `content` text NOT NULL COMMENT '回复内容（text 为纯文本；其它为 JSON）',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '0停用 1启用',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序，越小越靠前',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`wechat_reply_id`),
  UNIQUE KEY `uk_reply_key` (`reply_key`),
  KEY `idx_status_sort` (`status`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公众号自动回复';
