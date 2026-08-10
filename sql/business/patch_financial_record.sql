-- 平台「平台账单」：商户/平台财务流水（对齐 CRMEB eb_financial_record，前缀 qixi_crm_b_）
-- 幂等：CREATE IF NOT EXISTS。utf8mb4。
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_financial_record` (
  `financial_record_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `financial_record_sn` varchar(32) NOT NULL DEFAULT '' COMMENT '流水号',
  `order_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '订单 ID',
  `order_sn` varchar(64) NOT NULL DEFAULT '' COMMENT '订单编号',
  `user_info` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名（演示用，非敏感）',
  `user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `financial_type` varchar(64) NOT NULL DEFAULT '' COMMENT '流水类型',
  `financial_pm` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '0 支出 / 1 获得',
  `number` decimal(16,2) NOT NULL DEFAULT 0.00 COMMENT '金额',
  `type` tinyint NOT NULL DEFAULT 2 COMMENT '0 商户 / 1 公共 / 2 平台',
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '商户 ID',
  `pay_type` int NOT NULL DEFAULT 0 COMMENT '支付类型（7=线下）',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`financial_record_id`),
  KEY `idx_mer_id` (`mer_id`),
  KEY `idx_financial_type` (`financial_type`),
  KEY `idx_create_time` (`create_time`),
  KEY `idx_pay_type` (`pay_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='商户/平台财务流水';
