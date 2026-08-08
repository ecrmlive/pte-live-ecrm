-- 平台「维护」侧栏：开发配置 / 安全维护 目录 + 二级叶子
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='维护',
    `icon`='ant-design:tool-outlined',
    `route_path`='/maintain',
    `kind`='directory',
    `sort`=14,
    `parent_id`=0,
    `status`=1
WHERE `id`=197 OR `code`='maintain';

-- 旧扁平：热门搜索已迁到设置/商城设置，侧栏隐藏
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=201 OR `code`='maintain.hot_search';

-- 二级：目录 + 叶子
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1700,197,'maintain.exploit','开发配置','lucide:code','/safe/exploit','directory',1,1),
  (1701,197,'maintain.safe','安全维护','lucide:shield','/maintain','directory',2,1),
  (1702,197,'maintain.config_classify','配置分类','lucide:folder-tree','/config/classify','page',3,1),
  (1703,197,'maintain.config_setting','配置管理','lucide:sliders-horizontal','/config/setting','page',4,1),
  (1704,197,'maintain.system_log','操作日志','lucide:history','/setting/systemLog','page',5,1),
  (1705,197,'maintain.export','导出记录','lucide:download','/group/exportList','page',6,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=197,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：开发配置
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (200,1700,'maintain.group_data','组合数据','lucide:database','/group/list','page',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1700,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=1,
  `status`=1;

-- 三级：安全维护
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (199,1701,'maintain.backup','数据备份','lucide:hard-drive','/maintain/dataBackup','page',1,1),
  (1706,1701,'maintain.auth','商业授权','lucide:badge-check','/setting/system/maintain/auth','page',2,1),
  (198,1701,'maintain.cache','缓存清除','lucide:eraser','/maintain/cache','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=1701,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu` SET `parent_id`=198, `route_path`='maintain/cache' WHERE `id`=20989 OR `code`='maintain.cache.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=199, `route_path`='maintain/dataBackup' WHERE `id`=20990 OR `code`='maintain.backup.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=200, `route_path`='group/list' WHERE `id`=20991 OR `code`='maintain.group_data.manage';

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (197,198,199,200,1700,1701,1702,1703,1704,1705,1706);
