-- 一号通 → 服务配置：套餐 CRUD、状态与平台按钮权限。
-- 字段与 CRMEB eb_serve_meal 对齐，统一后台使用 qixi_crm_a_ 前缀。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_serve_meal` (
  `meal_id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '套餐 ID',
  `name` varchar(30) NOT NULL COMMENT '套餐名称',
  `type` tinyint NOT NULL DEFAULT 1 COMMENT '套餐类型：1商品采集，2电子面单',
  `price` decimal(8,2) NOT NULL DEFAULT 0.00 COMMENT '价格',
  `num` int NOT NULL DEFAULT 0 COMMENT '购买数量（次数）',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示：1显示，0隐藏',
  `is_del` tinyint NOT NULL DEFAULT 0 COMMENT '逻辑删除',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '添加时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`meal_id`),
  KEY `idx_visible_sort` (`is_del`,`status`,`sort`,`meal_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='一号通服务套餐';

INSERT INTO `qixi_crm_a_menu`
  (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`)
VALUES
  (21560,1511,'systemServeMerMealLst','服务套餐列表','','service/settings','button',1,1),
  (21561,1511,'systemServeMealDetail','服务套餐详情','','service/settings','button',2,1),
  (21562,1511,'systemServeMealCreate','新增服务套餐','','service/settings','button',3,1),
  (21563,1511,'systemServeMealUpdate','编辑服务套餐','','service/settings','button',4,1),
  (21564,1511,'systemServeMealDelete','删除服务套餐','','service/settings','button',5,1),
  (21565,1511,'systemServeMealStatus','修改服务套餐状态','','service/settings','button',6,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform'
  AND m.code IN (
    'systemServeMerMealLst','systemServeMealDetail','systemServeMealCreate',
    'systemServeMealUpdate','systemServeMealDelete','systemServeMealStatus'
  );
