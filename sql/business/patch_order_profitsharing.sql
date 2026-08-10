-- 订单分账（对齐 CRMEB eb_store_order_profitsharing）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_store_order_profitsharing` (
  `profitsharing_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `profitsharing_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '分账单号',
  `order_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '订单 ID',
  `order_sn` varchar(64) NOT NULL DEFAULT '' COMMENT '订单编号快照',
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '商户 ID',
  `mer_name` varchar(128) NOT NULL DEFAULT '' COMMENT '店铺名称快照',
  `transaction_id` varchar(60) NOT NULL DEFAULT '' COMMENT '支付渠道交易号（演示）',
  `profitsharing_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '订单/分账金额',
  `profitsharing_refund` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '已退款分账金额',
  `profitsharing_mer_price` decimal(10,2) NOT NULL DEFAULT 0.00 COMMENT '平台手续费',
  `type` varchar(32) NOT NULL DEFAULT 'order' COMMENT 'order=订单支付 presell=尾款支付',
  `status` tinyint NOT NULL DEFAULT 0 COMMENT '2分账中 1已分账 0待分账 -1已退款 -2分账失败',
  `error_msg` varchar(255) NOT NULL DEFAULT '' COMMENT '失败原因',
  `profitsharing_time` datetime DEFAULT NULL COMMENT '分账时间',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `is_combine` tinyint NOT NULL DEFAULT 1 COMMENT '1平台收付通 2服务商',
  PRIMARY KEY (`profitsharing_id`),
  KEY `idx_order_id` (`order_id`),
  KEY `idx_mer_id` (`mer_id`),
  KEY `idx_status` (`status`),
  KEY `idx_type` (`type`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_profitsharing_time` (`profitsharing_time`),
  KEY `idx_order_sn` (`order_sn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单分账';
