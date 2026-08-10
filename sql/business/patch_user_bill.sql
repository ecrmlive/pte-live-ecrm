-- 平台「积分日志」：用户积分/资产流水（对齐 CRMEB eb_user_bill，前缀 qixi_crm_b_）
-- 幂等：CREATE IF NOT EXISTS。utf8mb4。
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_user_bill` (
  `bill_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户 ID',
  `link_id` varchar(64) NOT NULL DEFAULT '' COMMENT '关联业务 ID',
  `pm` tinyint NOT NULL DEFAULT 1 COMMENT '1 收入 / 0 支出',
  `title` varchar(128) NOT NULL DEFAULT '' COMMENT '积分标题',
  `category` varchar(32) NOT NULL DEFAULT '' COMMENT '分类：integral / mer_integral 等',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '类型：sign_integral / deduction / lock 等',
  `number` decimal(16,2) NOT NULL DEFAULT 0.00 COMMENT '变动数量（绝对值）',
  `balance` decimal(16,2) NOT NULL DEFAULT 0.00 COMMENT '变动后余额',
  `mark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
  `mer_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '商户 ID（商户积分场景）',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '1 有效；积分冻结 lock 时 0=冻结中',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`bill_id`),
  KEY `idx_uid_category_time` (`uid`,`category`,`create_time`),
  KEY `idx_category_type` (`category`,`type`),
  KEY `idx_type_status` (`type`,`status`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户积分/资金流水';
