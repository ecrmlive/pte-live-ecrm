-- 平台「转账记录」：商户余额提现申请投影（对齐 CRMEB eb_financial type=0）
-- 幂等：可重复执行
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_financial` (
  `financial_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `financial_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '单号',
  `mer_money` decimal(12,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '申请时店铺余额快照',
  `extract_money` decimal(12,2) unsigned NOT NULL DEFAULT 0.00 COMMENT '提现金额',
  `financial_type` int unsigned NOT NULL DEFAULT 0 COMMENT '收款类型 1银行卡 2微信 3支付宝',
  `financial_account` varchar(1000) NOT NULL DEFAULT '' COMMENT '收款账户 JSON',
  `financial_status` int unsigned NOT NULL DEFAULT 0 COMMENT '到账状态 0未到账 1已到账',
  `status` int NOT NULL DEFAULT 0 COMMENT '审核 0待审核 1通过 -1未通过',
  `refusal` varchar(255) NOT NULL DEFAULT '' COMMENT '拒绝理由',
  `mer_id` bigint unsigned NOT NULL COMMENT '商户 ID',
  `image` varchar(1000) NOT NULL DEFAULT '' COMMENT '转账凭证',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '平台审核/打款管理员',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '申请时间',
  `status_time` datetime DEFAULT NULL COMMENT '审核时间',
  `update_time` datetime DEFAULT NULL COMMENT '凭证更新时间',
  `is_del` tinyint unsigned NOT NULL DEFAULT 0,
  `mark` varchar(255) NOT NULL DEFAULT '' COMMENT '商户备注',
  `admin_mark` varchar(255) NOT NULL DEFAULT '' COMMENT '平台备注',
  `mer_admin_id` bigint unsigned DEFAULT NULL COMMENT '商户管理员',
  `type` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '0余额提现 1保证金',
  PRIMARY KEY (`financial_id`),
  KEY `idx_mer_type_status` (`mer_id`,`type`,`status`,`financial_status`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_admin_id` (`admin_id`),
  UNIQUE KEY `uk_financial_sn` (`financial_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户转账申请（平台监管）';

-- 店铺监管投影补余额 / 冻结金额（统计卡用）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='mer_money')=0,
    'ALTER TABLE `qixi_crm_a_merchant_view` ADD COLUMN `mer_money` decimal(12,2) NOT NULL DEFAULT 0 AFTER `mark`, ADD COLUMN `freeze_money` decimal(12,2) NOT NULL DEFAULT 0 AFTER `mer_money`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
