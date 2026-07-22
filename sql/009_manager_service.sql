-- 阶段 6d：店员/客服账号（核销 + 代退）
USE `qixi_mergers`;

-- 子单核销客服 id（权威字段见 docs/schema）
SET @db := DATABASE();
SET @has_vsi := (
  SELECT COUNT(*) FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA = @db AND TABLE_NAME = 'qixi_store_order' AND COLUMN_NAME = 'verify_service_id'
);
SET @sql := IF(@has_vsi = 0,
  'ALTER TABLE `qixi_store_order` ADD COLUMN `verify_service_id` int(10) unsigned DEFAULT NULL COMMENT ''核销客服 id'' AFTER `verify_code`',
  'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

CREATE TABLE IF NOT EXISTS `qixi_store_service` (
  `service_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '客服id',
  `mer_id` int(11) NOT NULL DEFAULT 0 COMMENT '商户id',
  `uid` int(11) NOT NULL DEFAULT 0 COMMENT '关联用户uid',
  `avatar` varchar(250) NOT NULL DEFAULT '',
  `nickname` varchar(50) NOT NULL DEFAULT '',
  `account` varchar(32) DEFAULT NULL COMMENT '客服账号',
  `pwd` varchar(64) DEFAULT NULL COMMENT '客服密码',
  `is_open` tinyint(3) NOT NULL DEFAULT 1 COMMENT '开启登录',
  `status` tinyint(3) NOT NULL DEFAULT 1 COMMENT '0隐藏1显示',
  `notify` int(11) DEFAULT 0,
  `phone` varchar(18) DEFAULT '',
  `customer` tinyint(1) NOT NULL DEFAULT 1 COMMENT '客服/统计',
  `is_verify` tinyint(3) NOT NULL DEFAULT 1 COMMENT '核销权限',
  `is_goods` tinyint(3) DEFAULT 0,
  `sort` tinyint(3) NOT NULL DEFAULT 0,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(3) NOT NULL DEFAULT 0,
  PRIMARY KEY (`service_id`),
  UNIQUE KEY `uk_account` (`account`),
  KEY `idx_mer` (`mer_id`,`status`,`is_del`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='店员/客服';

-- 密码 admin123（与其它演示账号相同 bcrypt）
INSERT INTO `qixi_store_service` (
  `service_id`, `mer_id`, `uid`, `avatar`, `nickname`, `account`, `pwd`,
  `is_open`, `status`, `customer`, `is_verify`, `is_goods`, `phone`
)
SELECT 1, 1, 0, '', '店员小栖', 'staff1',
  '$2a$10$g9WCcDmxUSOewBGinelwoOeK94b3svdlJ8FGKb2Cv5xzKBjXMYAIG',
  1, 1, 1, 1, 0, '13900001111'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_store_service` WHERE `account` = 'staff1');

-- 确保演示订单有核销码（已支付子单）
UPDATE `qixi_store_order`
SET `verify_code` = CONCAT('V', LPAD(`order_id`, 6, '0'))
WHERE (`verify_code` IS NULL OR `verify_code` = '') AND `paid` = 1 AND `is_del` = 0;

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6d-manager-service', '阶段6d：店员账号 + 核销码补齐'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6d-manager-service');
