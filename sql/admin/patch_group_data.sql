SET NAMES utf8mb4;

-- 组合数据：数据组与组内数据分表，避免与普通平台配置条目混用。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_data_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL COMMENT '数据组名称',
  `group_key` varchar(64) NOT NULL COMMENT '数据组 key',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '数据组说明',
  `fields` json NOT NULL COMMENT '组内数据字段定义',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_key` (`group_key`),
  KEY `idx_visible_sort` (`is_del`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台组合数据组';

CREATE TABLE IF NOT EXISTS `qixi_crm_a_data_group_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `group_id` bigint unsigned NOT NULL COMMENT '数据组 ID',
  `data` json NOT NULL COMMENT '组合数据内容',
  `sort` int NOT NULL DEFAULT 0,
  `status` tinyint NOT NULL DEFAULT 1,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_group_visible_sort` (`group_id`,`is_del`,`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台组合数据项';

-- 将早期错误落入普通配置条目的组合数据迁入新的数据组结构；保留用户维护内容。
INSERT IGNORE INTO `qixi_crm_a_data_group`
  (`id`,`name`,`group_key`,`description`,`fields`,`sort`,`is_del`,`created_at`,`updated_at`)
SELECT `id`,`name`,IF(`code`='',CONCAT('legacy_group_',`id`),`code`),`remark`,JSON_ARRAY(),`sort`,`is_del`,`created_at`,`updated_at`
FROM `qixi_crm_a_config_item`
WHERE `item_type`='group_data';
