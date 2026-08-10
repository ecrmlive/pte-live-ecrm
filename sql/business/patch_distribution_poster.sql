-- 分销海报（对齐 CRMEB /group/config/68 组合数据 spread_banner）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_poster` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名称',
  `pic_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '背景图(600*1000px)',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '1显示 0隐藏',
  `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_b_distribution_poster` (`id`,`name`,`pic_url`,`status`,`sort`,`created_at`) VALUES
  (4173,'618','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2260%22%20height%3D%22100%22%20viewBox%3D%220%200%2060%20100%22%3E%3Crect%20width%3D%2260%22%20height%3D%22100%22%20fill%3D%22%23F97316%22/%3E%3Ctext%20x%3D%2230%22%20y%3D%2254%22%20fill%3D%22white%22%20font-size%3D%2214%22%20text-anchor%3D%22middle%22%20font-family%3D%22sans-serif%22%3E618%3C/text%3E%3C/svg%3E',1,1,'2025-06-17 18:25:23'),
  (4472,'测试2','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2260%22%20height%3D%22100%22%20viewBox%3D%220%200%2060%20100%22%3E%3Crect%20width%3D%2260%22%20height%3D%22100%22%20fill%3D%22%230EA5E9%22/%3E%3Ctext%20x%3D%2230%22%20y%3D%2254%22%20fill%3D%22white%22%20font-size%3D%2212%22%20text-anchor%3D%22middle%22%20font-family%3D%22sans-serif%22%3E%E6%B5%8B%E8%AF%95%3C/text%3E%3C/svg%3E',1,2,'2026-04-08 17:11:25')
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `pic_url`=IF(`pic_url` IS NULL OR `pic_url`='',VALUES(`pic_url`),`pic_url`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`);
