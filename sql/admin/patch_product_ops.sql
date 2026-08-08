-- 平台商品运营字段表（显示/星级/排序/推荐/详情 HTML）。
-- 修复：编辑商品 PUT /products/:id 因缺表 qixi_crm_a_product_ops 返回 500。
CREATE TABLE IF NOT EXISTS `qixi_crm_a_product_ops` (
  `product_id` bigint unsigned NOT NULL,
  `is_used` tinyint NOT NULL DEFAULT 1 COMMENT '平台是否显示',
  `star` tinyint NOT NULL DEFAULT 0 COMMENT '推荐级别 0-5',
  `rank_sort` int NOT NULL DEFAULT 0 COMMENT '平台排序',
  `is_hot` tinyint NOT NULL DEFAULT 0,
  `is_best` tinyint NOT NULL DEFAULT 0,
  `is_benefit` tinyint NOT NULL DEFAULT 0,
  `is_new` tinyint NOT NULL DEFAULT 0,
  `cate_hot` tinyint NOT NULL DEFAULT 0,
  `sys_labels` varchar(500) NOT NULL DEFAULT '',
  `content_html` mediumtext COMMENT '平台侧商品详情/营销详情 HTML',
  `refund_switch` tinyint NOT NULL DEFAULT 1 COMMENT '支持退款',
  `once_min_count` int NOT NULL DEFAULT 1 COMMENT '最少购买件数',
  `ficti` int NOT NULL DEFAULT 0 COMMENT '虚拟已售数量',
  `updated_by` bigint unsigned NOT NULL DEFAULT 0,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`), KEY `idx_used_star` (`is_used`,`star`,`rank_sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_ops' AND COLUMN_NAME='content_html')=0,
    'ALTER TABLE `qixi_crm_a_product_ops` ADD COLUMN `content_html` mediumtext NULL AFTER `sys_labels`, ADD COLUMN `refund_switch` tinyint NOT NULL DEFAULT 1 AFTER `content_html`, ADD COLUMN `once_min_count` int NOT NULL DEFAULT 1 AFTER `refund_switch`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_ops' AND COLUMN_NAME='ficti')=0,
    'ALTER TABLE `qixi_crm_a_product_ops` ADD COLUMN `ficti` int NOT NULL DEFAULT 0 COMMENT ''虚拟已售数量'' AFTER `once_min_count`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
