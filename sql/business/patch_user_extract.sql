-- 用户佣金提现（对齐 CRMEB eb_user_extract）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_extract` (
  `extract_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `extract_sn` varchar(64) NOT NULL DEFAULT '' COMMENT '提现单号',
  `real_name` varchar(64) NOT NULL DEFAULT '' COMMENT '户名',
  `extract_type` tinyint NOT NULL DEFAULT 0 COMMENT '0银行卡 1微信 2支付宝 3微信零钱 4余额',
  `bank_code` varchar(64) NOT NULL DEFAULT '' COMMENT '银行卡号（演示脱敏）',
  `bank_address` varchar(256) NOT NULL DEFAULT '' COMMENT '开户地址',
  `bank_name` varchar(128) NOT NULL DEFAULT '' COMMENT '银行名称',
  `alipay_code` varchar(64) NOT NULL DEFAULT '' COMMENT '支付宝账号（演示）',
  `wechat` varchar(64) NOT NULL DEFAULT '' COMMENT '微信号（演示）',
  `extract_pic` varchar(1024) NOT NULL DEFAULT '' COMMENT '收款码',
  `extract_price` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '提现金额',
  `balance` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '提现后佣金余额',
  `mark` varchar(512) NOT NULL DEFAULT '' COMMENT '管理员备注',
  `admin_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '审核管理员',
  `fail_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '拒绝原因',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '-1已拒绝 0审核中 1已通过',
  `status_time` datetime DEFAULT NULL COMMENT '审核时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  PRIMARY KEY (`extract_id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_status` (`status`),
  KEY `idx_extract_type` (`extract_type`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_extract_sn` (`extract_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户佣金提现';
