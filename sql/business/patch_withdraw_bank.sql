-- 分销提现银行（对齐 CRMEB /group/config/76 组合数据）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_withdraw_bank` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '1显示 0隐藏',
  `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_b_withdraw_bank` (`id`,`name`,`status`,`sort`) VALUES
  (1,'中国银行',1,1),
  (2,'招商银行',1,2),
  (3,'工商银行',1,3)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`);
