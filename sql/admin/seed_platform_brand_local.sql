-- 本地开发：平台品牌分类 + 品牌样例数据（utf8mb4 中文可读假数据）
-- 用法：在 qixi_crm_admin 执行；可重复执行（按名称幂等）

SET NAMES utf8mb4;

-- 补齐创建时间列（若已有则跳过）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_brand_category' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_brand_category` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_platform_brand' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_a_platform_brand` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `status`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 一级分类
INSERT INTO `qixi_crm_a_platform_brand_category` (`parent_id`, `name`, `sort`, `status`, `created_at`)
SELECT 0, v.name, v.sort, 1, v.created_at
FROM (
  SELECT '服饰' AS name, 10 AS sort, '2026-01-23 15:25:27' AS created_at UNION ALL
  SELECT '家电', 20, '2026-01-20 10:00:00' UNION ALL
  SELECT '美妆个护', 30, '2026-01-18 11:30:00' UNION ALL
  SELECT '食品饮料', 40, '2026-01-15 09:12:00' UNION ALL
  SELECT '数码3C', 50, '2026-01-12 14:08:00'
) v
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_platform_brand_category` c
  WHERE c.parent_id = 0 AND c.name = v.name
);

-- 服饰下级
INSERT INTO `qixi_crm_a_platform_brand_category` (`parent_id`, `name`, `sort`, `status`, `created_at`)
SELECT p.id, v.name, v.sort, 1, v.created_at
FROM `qixi_crm_a_platform_brand_category` p
JOIN (
  SELECT '男装' AS name, 1 AS sort, '2026-01-23 16:00:00' AS created_at UNION ALL
  SELECT '女装', 2, '2026-01-23 16:01:00' UNION ALL
  SELECT '羽绒服', 3, '2026-01-23 16:02:00'
) v
WHERE p.parent_id = 0 AND p.name = '服饰'
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_a_platform_brand_category` c
    WHERE c.parent_id = p.id AND c.name = v.name
  );

-- 家电下级
INSERT INTO `qixi_crm_a_platform_brand_category` (`parent_id`, `name`, `sort`, `status`, `created_at`)
SELECT p.id, v.name, v.sort, 1, v.created_at
FROM `qixi_crm_a_platform_brand_category` p
JOIN (
  SELECT '大家电' AS name, 1 AS sort, '2026-01-20 10:10:00' AS created_at UNION ALL
  SELECT '小家电', 2, '2026-01-20 10:11:00'
) v
WHERE p.parent_id = 0 AND p.name = '家电'
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_a_platform_brand_category` c
    WHERE c.parent_id = p.id AND c.name = v.name
  );

-- 品牌（挂到对应分类；对齐设计图常见品牌名）
INSERT INTO `qixi_crm_a_platform_brand` (`name`, `category_id`, `logo_url`, `sort`, `status`, `created_at`)
SELECT v.name, COALESCE(c.id, 0), '', v.sort, 1, v.created_at
FROM (
  SELECT '南极人' AS name, '服饰' AS cat, 90 AS sort, '2026-01-23 15:30:00' AS created_at UNION ALL
  SELECT '三彩', '服饰', 80, '2026-01-22 11:00:00' UNION ALL
  SELECT '鸭鸭', '羽绒服', 70, '2026-01-21 09:20:00' UNION ALL
  SELECT '波司登', '羽绒服', 60, '2026-01-20 16:40:00' UNION ALL
  SELECT '雪中飞', '羽绒服', 50, '2025-12-18 10:00:00' UNION ALL
  SELECT '海尔', '大家电', 40, '2025-11-02 14:22:00' UNION ALL
  SELECT '小米', '数码3C', 30, '2025-10-15 08:45:00' UNION ALL
  SELECT '无印良品', '美妆个护', 20, '2025-09-08 13:10:00' UNION ALL
  SELECT '心心相印', '美妆个护', 10, '2025-07-29 17:05:00'
) v
LEFT JOIN `qixi_crm_a_platform_brand_category` c ON c.name = v.cat
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_platform_brand` b WHERE b.name = v.name
);
