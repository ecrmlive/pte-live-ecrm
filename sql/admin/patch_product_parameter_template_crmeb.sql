-- 平台商品参数模板：关联分类 + 参数项（名称/候选值/必填/排序），对齐 CRMEB parameter_template
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_parameter_template' AND COLUMN_NAME='cate_ids_json')=0,
    'ALTER TABLE `qixi_crm_a_product_parameter_template` ADD COLUMN `cate_ids_json` json NULL COMMENT ''关联平台分类 ID 列表'' AFTER `name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_product_parameter_template' AND COLUMN_NAME='params_json')=0,
    'ALTER TABLE `qixi_crm_a_product_parameter_template` ADD COLUMN `params_json` json NULL COMMENT ''参数项 [{name,values,required,sort}]'' AFTER `cate_ids_json`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

-- 旧 values_json（字符串数组）迁移为单条参数项；可重复执行
UPDATE `qixi_crm_a_product_parameter_template`
SET `params_json` = JSON_ARRAY(
  JSON_OBJECT(
    'name', `name`,
    'values', `values_json`,
    'required', 0,
    'sort', 0
  )
)
WHERE (`params_json` IS NULL OR JSON_TYPE(`params_json`) = 'NULL' OR JSON_LENGTH(`params_json`) = 0)
  AND `values_json` IS NOT NULL
  AND JSON_TYPE(`values_json`) = 'ARRAY'
  AND JSON_LENGTH(`values_json`) > 0;

UPDATE `qixi_crm_a_product_parameter_template`
SET `cate_ids_json` = JSON_ARRAY()
WHERE `cate_ids_json` IS NULL OR JSON_TYPE(`cate_ids_json`) = 'NULL';

UPDATE `qixi_crm_a_product_parameter_template`
SET `params_json` = JSON_ARRAY()
WHERE `params_json` IS NULL OR JSON_TYPE(`params_json`) = 'NULL';

-- 演示模板（对齐截图分类：电脑/连衣裙、手机/香水/时尚潮牌/CRMEB上门…）
INSERT INTO `qixi_crm_a_product_parameter_template`
  (`id`,`name`,`cate_ids_json`,`params_json`,`values_json`,`sort`,`status`)
VALUES
  (
    7521,
    '测试',
    JSON_ARRAY(7635, 7628),
    JSON_ARRAY(
      JSON_OBJECT('name','颜色','values',JSON_ARRAY('中国红','竹青色','云白色'),'required',0,'sort',100),
      JSON_OBJECT('name','尺码','values',JSON_ARRAY('S','M','L'),'required',0,'sort',90)
    ),
    JSON_ARRAY('中国红','竹青色','云白色'),
    0,
    1
  ),
  (
    7522,
    '通用参数',
    JSON_ARRAY(7636, 7632, 7629, 7621, 7635, 7628, 7631),
    JSON_ARRAY(
      JSON_OBJECT('name','品牌','values',JSON_ARRAY('七禧','演示品牌'),'required',0,'sort',100),
      JSON_OBJECT('name','产地','values',JSON_ARRAY('中国','进口'),'required',0,'sort',90)
    ),
    JSON_ARRAY('七禧','演示品牌'),
    0,
    1
  )
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `cate_ids_json`=VALUES(`cate_ids_json`),
  `params_json`=VALUES(`params_json`),
  `values_json`=VALUES(`values_json`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`);
