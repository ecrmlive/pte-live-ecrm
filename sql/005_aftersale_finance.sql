-- 阶段 4：售后仅退款 + 商户提现（审核拒绝退余额）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_store_refund_order` (
  `refund_order_id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '退款单id',
  `refund_order_sn` varchar(32) NOT NULL COMMENT '退款单号',
  `order_id` int(10) unsigned NOT NULL COMMENT '子单id',
  `uid` int(10) unsigned NOT NULL DEFAULT '0',
  `mer_id` int(10) unsigned NOT NULL,
  `delivery_type` varchar(32) DEFAULT NULL,
  `delivery_id` varchar(32) DEFAULT NULL,
  `delivery_mark` varchar(200) DEFAULT NULL,
  `phone` varchar(18) DEFAULT NULL,
  `mark` varchar(200) NOT NULL DEFAULT '',
  `mer_mark` varchar(255) NOT NULL DEFAULT '',
  `admin_mark` varchar(255) NOT NULL DEFAULT '',
  `pics` text,
  `refund_type` tinyint(1) NOT NULL COMMENT '1仅退款 2退货退款',
  `refund_message` varchar(128) NOT NULL DEFAULT '' COMMENT '退款原因',
  `refund_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `platform_refund_price` decimal(8,2) DEFAULT '0.00',
  `refund_postage` decimal(8,2) DEFAULT '0.00',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0',
  `fail_message` varchar(200) DEFAULT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0待审 -1拒绝 1待退货 2待收货 3已退款 4平台介入 -2取消',
  `status_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `is_del` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `is_system_del` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`refund_order_id`),
  UNIQUE KEY `refund_order_sn` (`refund_order_sn`),
  KEY `idx_order` (`order_id`),
  KEY `idx_mer_status` (`mer_id`,`status`),
  KEY `idx_uid` (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款单';

CREATE TABLE IF NOT EXISTS `qixi_store_refund_product` (
  `refund_product_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `refund_order_id` int(10) unsigned NOT NULL,
  `order_product_id` int(10) unsigned NOT NULL,
  `refund_price` decimal(8,2) NOT NULL DEFAULT '0.00',
  `refund_num` int(10) unsigned NOT NULL DEFAULT '0',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`refund_product_id`),
  KEY `refund_order_id` (`refund_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款商品行';

CREATE TABLE IF NOT EXISTS `qixi_store_refund_status` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `refund_order_id` int(10) unsigned NOT NULL,
  `change_type` varchar(32) NOT NULL DEFAULT '',
  `change_message` varchar(256) NOT NULL DEFAULT '',
  `change_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `refund_order_id` (`refund_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='退款状态日志';

CREATE TABLE IF NOT EXISTS `qixi_financial` (
  `financial_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `financial_sn` varchar(32) NOT NULL,
  `mer_money` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '申请时余额快照',
  `extract_money` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '提现金额',
  `financial_type` int(10) DEFAULT '1' COMMENT '1银行卡 2微信 3支付宝',
  `financial_account` varchar(500) NOT NULL DEFAULT '',
  `financial_status` int(10) DEFAULT '0' COMMENT '0未打款 1已打款',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '0待审 1通过 -1拒绝',
  `refusal` varchar(128) DEFAULT NULL,
  `mer_id` int(10) unsigned NOT NULL,
  `image` varchar(1000) DEFAULT NULL,
  `admin_id` int(11) DEFAULT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `status_time` timestamp NULL DEFAULT NULL,
  `update_time` timestamp NULL DEFAULT NULL,
  `is_del` int(10) DEFAULT '0',
  `mark` varchar(255) DEFAULT NULL,
  `admin_mark` varchar(255) DEFAULT NULL,
  `mer_admin_id` int(11) DEFAULT NULL,
  `type` int(10) DEFAULT '0' COMMENT '0余额提现 1保证金',
  PRIMARY KEY (`financial_id`),
  UNIQUE KEY `financial_sn` (`financial_sn`),
  KEY `idx_mer_status` (`mer_id`,`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户提现申请';

CREATE TABLE IF NOT EXISTS `qixi_financial_record` (
  `financial_record_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `financial_record_sn` varchar(32) NOT NULL,
  `order_id` int(10) unsigned NOT NULL DEFAULT '0',
  `order_sn` varchar(32) NOT NULL DEFAULT '',
  `user_info` varchar(32) NOT NULL DEFAULT '',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0',
  `financial_type` varchar(32) NOT NULL DEFAULT '' COMMENT '流水类型如 order_pay/refund/extract',
  `financial_pm` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0支出 1收入',
  `number` decimal(8,2) NOT NULL DEFAULT '0.00',
  `type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0商户',
  `mer_id` int(10) unsigned NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `pay_type` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`financial_record_id`),
  KEY `idx_mer` (`mer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户财务流水';

-- 演示商户可提现余额
UPDATE `qixi_merchant` SET `mer_money` = 5000.00 WHERE `mer_id` IN (1, 2) AND (`mer_money` IS NULL OR `mer_money` < 100);

-- 平台菜单：退款监管 + 提现审核
-- 注意：menu_id=14 已被 002 品牌占用，此处使用 30/31/32
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 30, 8, '/order/refund', '', '退款监管', 'OrderRefund', '', 68, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/order/refund' AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 31, 0, '/accounts', 'AccountBookOutlined', '财务', 'Accounts', '', 60, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/accounts' AND `is_mer` = 1);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 32, 31, '/accounts/withdraw', '', '提现审核', 'AccountsWithdraw', '', 59, 1, 1, 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/accounts/withdraw' AND `is_mer` = 1);

-- 商户菜单：售后 + 提现
INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 113, 105, '/order/refund', '', '售后处理', 'MerOrderRefund', '', 77, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/order/refund' AND `is_mer` = 2);

INSERT INTO `qixi_system_menu` (`menu_id`, `pid`, `path`, `icon`, `menu_name`, `route`, `params`, `sort`, `is_show`, `is_mer`, `is_menu`, `is_agent`)
SELECT 114, 108, '/finance/withdraw', '', '提现申请', 'MerFinanceWithdraw', '', 68, 1, 2, 1, 2
WHERE NOT EXISTS (SELECT 1 FROM `qixi_system_menu` WHERE `path` = '/finance/withdraw' AND `is_mer` = 2);

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',30,31,32')
WHERE `role_id` = 1 AND `rules` NOT LIKE '%30%';

UPDATE `qixi_system_role`
SET `rules` = CONCAT(`rules`, ',113,114')
WHERE `role_id` = 2 AND `rules` NOT LIKE '%113%';

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase4-aftersale-finance', '阶段4：仅退款状态机 + 商户提现审核拒绝退余额'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase4-aftersale-finance');
