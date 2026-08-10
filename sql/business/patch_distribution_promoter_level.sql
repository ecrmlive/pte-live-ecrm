-- 分销员等级 + 推广员 level_id（对齐 CRMEB user_brokerage / promoter list）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_level` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '',
  `rank` int NOT NULL DEFAULT 0,
  `icon_url` varchar(1024) NOT NULL DEFAULT '',
  `task_rule` json DEFAULT NULL,
  `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT 0.00,
  `extension_two` decimal(8,2) unsigned NOT NULL DEFAULT 0.00,
  `status` tinyint NOT NULL DEFAULT 1,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rank` (`rank`),
  KEY `idx_status_rank` (`status`,`rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_distribution_promoter' AND COLUMN_NAME='level_id')=0,
    'ALTER TABLE `qixi_crm_b_distribution_promoter` ADD COLUMN `level_id` bigint unsigned NOT NULL DEFAULT 0 AFTER `status`, ADD KEY `idx_level_id` (`level_id`)',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_distribution_level' AND COLUMN_NAME='task_rule')=0,
    'ALTER TABLE `qixi_crm_b_distribution_level` ADD COLUMN `task_rule` json DEFAULT NULL AFTER `icon_url`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_distribution_level' AND COLUMN_NAME='extension_one')=0,
    'ALTER TABLE `qixi_crm_b_distribution_level` ADD COLUMN `extension_one` decimal(8,2) unsigned NOT NULL DEFAULT 0.00 AFTER `task_rule`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_distribution_level' AND COLUMN_NAME='extension_two')=0,
    'ALTER TABLE `qixi_crm_b_distribution_level` ADD COLUMN `extension_two` decimal(8,2) unsigned NOT NULL DEFAULT 0.00 AFTER `extension_one`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

INSERT INTO `qixi_crm_b_distribution_level` (`id`,`name`,`rank`,`icon_url`,`task_rule`,`extension_one`,`extension_two`,`status`) VALUES
  (1,'普通分销员',1,'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Crect%20x%3D%2214%22%20y%3D%2210%22%20width%3D%2236%22%20height%3D%2244%22%20rx%3D%224%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3Cpath%20d%3D%22M22%2024h20M22%2034h20M22%2044h12%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22%20stroke-linecap%3D%22round%22/%3E%3C/svg%3E',
   CAST('{"spread_user":{"name":"邀请新人","num":3,"info":""},"pay_money":{"name":"消费门槛","num":0,"info":""},"pay_num":{"name":"","num":0,"info":""},"spread_money":{"name":"","num":0,"info":""},"spread_pay_num":{"name":"","num":0,"info":""}}' AS JSON),
   0.00,0.00,1),
  (2,'金牌分销员',2,'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Crect%20x%3D%2212%22%20y%3D%2212%22%20width%3D%2216%22%20height%3D%2216%22%20rx%3D%223%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3Crect%20x%3D%2236%22%20y%3D%2212%22%20width%3D%2216%22%20height%3D%2216%22%20rx%3D%223%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3Crect%20x%3D%2212%22%20y%3D%2236%22%20width%3D%2216%22%20height%3D%2216%22%20rx%3D%223%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3Crect%20x%3D%2236%22%20y%3D%2236%22%20width%3D%2216%22%20height%3D%2216%22%20rx%3D%223%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3C/svg%3E',
   CAST('{"spread_user":{"name":"邀请新人","num":5,"info":""},"pay_money":{"name":"消费门槛","num":5000,"info":""},"pay_num":{"name":"消费单量","num":3,"info":""},"spread_money":{"name":"下级消费金额","num":3000,"info":""},"spread_pay_num":{"name":"下级消费单量","num":5,"info":""}}' AS JSON),
   1.50,1.50,1),
  (3,'钻石分销员',3,'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2224%22%20r%3D%2210%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22/%3E%3Cpath%20d%3D%22M14%2052c2-10%2010-16%2018-16s16%206%2018%2016%22%20fill%3D%22none%22%20stroke%3D%22%23F59E0B%22%20stroke-width%3D%223%22%20stroke-linecap%3D%22round%22/%3E%3C/svg%3E',
   CAST('{"spread_user":{"name":"邀请新人","num":10,"info":""},"pay_money":{"name":"消费门槛","num":20000,"info":""},"pay_num":{"name":"消费单量","num":5,"info":""},"spread_money":{"name":"下级消费金额","num":10000,"info":""},"spread_pay_num":{"name":"下级消费单量","num":10,"info":""}}' AS JSON),
   3.00,2.00,1)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `rank`=VALUES(`rank`),
  `icon_url`=IF(`icon_url` IS NULL OR `icon_url`='',VALUES(`icon_url`),`icon_url`),
  `task_rule`=VALUES(`task_rule`),
  `extension_one`=VALUES(`extension_one`),
  `extension_two`=VALUES(`extension_two`),
  `status`=VALUES(`status`);

-- 演示推广员挂等级（幂等）
UPDATE `qixi_crm_b_distribution_promoter` SET `level_id`=2 WHERE `user_id`=9001 AND `level_id`=0;
UPDATE `qixi_crm_b_distribution_promoter` SET `level_id`=1 WHERE `user_id`=9101 AND `level_id`=0;
