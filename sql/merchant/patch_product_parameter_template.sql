-- 店铺商品参数模板（对齐 CRMEB eb_parameter_template + eb_parameter）
-- 平台后台「店铺商品参数」/product/merSpecs 只读列表与详情；复制写入平台参数模板。
-- 商户后台 /product/specs 新增/编辑（含平台分类、是否必选）。
-- 用法：导入 qixi_crm_merchant（make local-sync-sql 会自动执行）

SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_parameter_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `store_id` bigint unsigned NOT NULL COMMENT '店铺 ID（qixi_crm_m_store.id）',
  `template_name` varchar(64) NOT NULL COMMENT '参数模板名称',
  `cate_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '关联平台分类 ID',
  `is_required` tinyint NOT NULL DEFAULT 0 COMMENT '是否必选：1=店铺关联该分类时必须使用本模板',
  `params_json` json NOT NULL COMMENT '参数项数组 [{name,values[],required,sort}]',
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_store_sort` (`store_id`,`is_del`,`sort`,`id`),
  KEY `idx_store_cate` (`store_id`,`cate_id`,`is_del`),
  KEY `idx_name` (`template_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='店铺商品参数模板';

-- 幂等补列（已有库）
SET @db := DATABASE();
SET @sql := (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=@db AND TABLE_NAME='qixi_crm_m_product_parameter_template' AND COLUMN_NAME='cate_id'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_crm_m_product_parameter_template` ADD COLUMN `cate_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT ''关联平台分类 ID'' AFTER `template_name`'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=@db AND TABLE_NAME='qixi_crm_m_product_parameter_template' AND COLUMN_NAME='is_required'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_crm_m_product_parameter_template` ADD COLUMN `is_required` tinyint NOT NULL DEFAULT 0 COMMENT ''是否必选：1=店铺关联该分类时必须使用本模板'' AFTER `cate_id`'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

SET @sql := (
  SELECT IF(
    EXISTS(
      SELECT 1 FROM information_schema.STATISTICS
      WHERE TABLE_SCHEMA=@db AND TABLE_NAME='qixi_crm_m_product_parameter_template' AND INDEX_NAME='idx_store_cate'
    ),
    'SELECT 1',
    'ALTER TABLE `qixi_crm_m_product_parameter_template` ADD KEY `idx_store_cate` (`store_id`,`cate_id`,`is_del`)'
  )
);
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;
