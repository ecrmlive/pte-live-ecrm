-- 用户余额充值记录（对齐 CRMEB eb_user_recharge）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_recharge` (
  `recharge_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '充值用户 UID',
  `order_id` varchar(64) NOT NULL DEFAULT '' COMMENT '订单号',
  `price` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '充值金额',
  `give_price` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '赠送金额',
  `recharge_type` varchar(32) NOT NULL DEFAULT '' COMMENT '充值类型 routine/weixin/h5/alipay',
  `paid` tinyint unsigned NOT NULL DEFAULT 0 COMMENT '是否支付 0未支付 1已支付',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `refund_price` decimal(12,2) NOT NULL DEFAULT 0.00 COMMENT '已退款金额',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`recharge_id`),
  UNIQUE KEY `uk_order_id` (`order_id`),
  KEY `idx_uid` (`uid`),
  KEY `idx_paid` (`paid`),
  KEY `idx_recharge_type` (`recharge_type`),
  KEY `idx_pay_time` (`pay_time`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户余额充值记录';
