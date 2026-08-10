-- 分销特权（对齐 CRMEB /group/config/75 组合数据 promoter_config）
-- 用法：make local-sync-sql 或 scripts/local-dev-sync.sh sql
SET NAMES utf8mb4;
USE `qixi_crm_business`;

CREATE TABLE IF NOT EXISTS `qixi_crm_b_distribution_privilege` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
  `img_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '图片(90*90px)',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '1显示 0隐藏',
  `sort` int NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status_sort` (`status`,`sort`,`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_b_distribution_privilege` (`id`,`title`,`img_url`,`status`,`sort`,`created_at`) VALUES
  (245,'零成本','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2290%22%20height%3D%2290%22%20viewBox%3D%220%200%2090%2090%22%3E%3Ccircle%20cx%3D%2245%22%20cy%3D%2245%22%20r%3D%2228%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22/%3E%3Cellipse%20cx%3D%2245%22%20cy%3D%2245%22%20rx%3D%2240%22%20ry%3D%2214%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%223%22%20transform%3D%22rotate(-24%2045%2045)%22/%3E%3C/svg%3E',1,1,'2020-06-26 15:12:12'),
  (246,'高佣金','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2290%22%20height%3D%2290%22%20viewBox%3D%220%200%2090%2090%22%3E%3Ccircle%20cx%3D%2245%22%20cy%3D%2234%22%20r%3D%2214%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22/%3E%3Cpath%20d%3D%22M20%2072c3-14%2012-22%2025-22s22%208%2025%2022%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22%20stroke-linecap%3D%22round%22/%3E%3C/svg%3E',1,2,'2020-06-26 15:12:12'),
  (247,'持续收入','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2290%22%20height%3D%2290%22%20viewBox%3D%220%200%2090%2090%22%3E%3Crect%20x%3D%2222%22%20y%3D%2228%22%20width%3D%2246%22%20height%3D%2236%22%20rx%3D%226%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22/%3E%3Cpath%20d%3D%22M34%2040h22M34%2050h14%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%223%22%20stroke-linecap%3D%22round%22/%3E%3C/svg%3E',1,3,'2020-06-26 15:12:12'),
  (248,'佣金抵现','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2290%22%20height%3D%2290%22%20viewBox%3D%220%200%2090%2090%22%3E%3Ccircle%20cx%3D%2238%22%20cy%3D%2242%22%20r%3D%2218%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22/%3E%3Cpath%20d%3D%22M52%2058l14%2014M58%2052h12v12%22%20fill%3D%22none%22%20stroke%3D%22%233B82F6%22%20stroke-width%3D%224%22%20stroke-linecap%3D%22round%22%20stroke-linejoin%3D%22round%22/%3E%3C/svg%3E',1,4,'2020-06-26 15:12:12')
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `img_url`=IF(`img_url` IS NULL OR `img_url`='',VALUES(`img_url`),`img_url`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`);
