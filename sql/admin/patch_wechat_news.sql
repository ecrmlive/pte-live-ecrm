-- 公众号图文管理（对齐 CRMEB wechat_news）
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_wechat_news` (
  `wechat_news_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '图文组 ID',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '状态 1启用 0停用',
  `items` json NOT NULL COMMENT '图文条目 [{title,author,synopsis,image,content}]',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`wechat_news_id`),
  KEY `idx_status_ctime` (`status`, `create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公众号图文消息';

-- 读权限（列表/详情）；manage 保留写操作
INSERT INTO `qixi_crm_a_menu`
  (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES
  (21120,206,'app.wechat_news.read','查看图文管理','','app/wechat/newsCategory','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu`
SET `title`='维护图文管理', `sort`=2, `status`=1
WHERE `id`=20997 OR `code`='app.wechat_news.manage';

-- 平台超管角色挂上读权限
INSERT INTO `qixi_crm_a_role_menu` (`role_id`, `menu_id`)
SELECT r.id, 21120
FROM `qixi_crm_a_role` r
WHERE r.code IN ('platform', 'operations')
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_a_role_menu` rm
    WHERE rm.role_id = r.id AND rm.menu_id = 21120
  );
