-- 本地 DEMO：店铺商品参数模板（utf8mb4 中文可读假数据）
-- 依赖：qixi_crm_m_store 至少存在 id=1、2；平台分类 id 见 qixi_crm_a_platform_category
-- 幂等按 store_id+template_name

SET NAMES utf8mb4;
USE `qixi_crm_merchant`;

CREATE TABLE IF NOT EXISTS `qixi_crm_m_product_parameter_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `store_id` bigint unsigned NOT NULL,
  `template_name` varchar(64) NOT NULL,
  `cate_id` bigint unsigned NOT NULL DEFAULT 0,
  `is_required` tinyint NOT NULL DEFAULT 0,
  `params_json` json NOT NULL,
  `sort` int NOT NULL DEFAULT 0,
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_store_sort` (`store_id`,`is_del`,`sort`,`id`),
  KEY `idx_name` (`template_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `qixi_crm_m_product_parameter_template`
  (`store_id`,`template_name`,`cate_id`,`is_required`,`params_json`,`sort`,`is_del`,`created_at`)
SELECT v.store_id, v.template_name, v.cate_id, v.is_required, v.params_json, v.sort, 0, v.created_at
FROM (
  SELECT 1 AS store_id, '测试' AS template_name, 7613 AS cate_id, 0 AS is_required,
    JSON_ARRAY(
      JSON_OBJECT('name','材质','values',JSON_ARRAY('棉','涤纶'),'required',1,'sort',10),
      JSON_OBJECT('name','厚度','values',JSON_ARRAY('薄','适中','厚'),'required',0,'sort',5)
    ) AS params_json,
    150 AS sort, '2026-05-26 17:01:15' AS created_at
  UNION ALL
  SELECT 1, '上衣参数', 7613, 1,
    JSON_ARRAY(
      JSON_OBJECT('name','版型','values',JSON_ARRAY('修身','宽松'),'required',1,'sort',20),
      JSON_OBJECT('name','袖长','values',JSON_ARRAY('短袖','长袖'),'required',0,'sort',10)
    ),
    99, '2025-12-10 17:18:53'
  UNION ALL
  SELECT 1, '尺码对应表', 7613, 0,
    JSON_ARRAY(
      JSON_OBJECT('name','尺码','values',JSON_ARRAY('S','M','L','XL'),'required',1,'sort',30),
      JSON_OBJECT('name','胸围(cm)','values',JSON_ARRAY('84','88','92','96'),'required',0,'sort',20)
    ),
    0, '2025-09-17 11:53:53'
  UNION ALL
  SELECT 1, '鞋子参数', 7603, 0,
    JSON_ARRAY(
      JSON_OBJECT('name','鞋码','values',JSON_ARRAY('36','37','38','39','40'),'required',1,'sort',40),
      JSON_OBJECT('name','鞋面','values',JSON_ARRAY('真皮','网布'),'required',0,'sort',10)
    ),
    0, '2025-09-17 11:48:36'
  UNION ALL
  SELECT 2, '居家日用参数', 7612, 0,
    JSON_ARRAY(
      JSON_OBJECT('name','容量','values',JSON_ARRAY('小','中','大'),'required',1,'sort',10)
    ),
    80, '2026-03-01 10:00:00'
) v
WHERE EXISTS (SELECT 1 FROM `qixi_crm_m_store` s WHERE s.id = v.store_id)
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_m_product_parameter_template` t
    WHERE t.store_id = v.store_id AND t.template_name = v.template_name AND t.is_del = 0
  );

-- 已有 DEMO 行对齐有效平台分类（避免复制到平台时校验失败）
UPDATE `qixi_crm_m_product_parameter_template` SET `cate_id`=7613
WHERE `is_del`=0 AND `store_id`=1 AND `template_name` IN ('测试','上衣参数','尺码对应表');
UPDATE `qixi_crm_m_product_parameter_template` SET `cate_id`=7603
WHERE `is_del`=0 AND `store_id`=1 AND `template_name`='鞋子参数';
UPDATE `qixi_crm_m_product_parameter_template` SET `cate_id`=7612
WHERE `is_del`=0 AND `store_id`=2 AND `template_name`='居家日用参数';
