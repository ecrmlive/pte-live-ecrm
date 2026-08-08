-- 区域代理侧栏对齐为 3 项：区域列表 / 代理人员 / 代理设置
-- 可重复执行；代理审核页软隐藏（status=0），路由与审核按钮权限保留，不删 API。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

UPDATE `qixi_crm_a_menu`
SET `title`='区域列表',
    `icon`='lucide:map-pinned',
    `route_path`='/business-zones/index',
    `kind`='page',
    `sort`=1,
    `parent_id`=20,
    `status`=1
WHERE `id`=21 OR `code`='region.index';

UPDATE `qixi_crm_a_menu`
SET `title`='代理人员',
    `icon`='lucide:users',
    `route_path`='/business-zones/agents',
    `kind`='page',
    `sort`=2,
    `parent_id`=20,
    `status`=1
WHERE `id`=22 OR `code`='region.agents';

UPDATE `qixi_crm_a_menu`
SET `title`='代理审核',
    `icon`='lucide:badge-check',
    `route_path`='/business-zones/agent-review',
    `kind`='page',
    `sort`=99,
    `parent_id`=20,
    `status`=0
WHERE `id`=23 OR `code`='region.agent_review';

UPDATE `qixi_crm_a_menu`
SET `title`='代理设置',
    `icon`='lucide:settings-2',
    `route_path`='/business-zones/settings',
    `kind`='page',
    `sort`=3,
    `parent_id`=20,
    `status`=1
WHERE `id`=24 OR `code`='region.agent_settings';

UPDATE `qixi_crm_a_menu`
SET `parent_id`=22,
    `title`='审核区域代理申请',
    `route_path`='business-zones/agent-review',
    `kind`='button',
    `sort`=1,
    `status`=1
WHERE `id`=20931 OR `code`='region.agent.review';

UPDATE `qixi_crm_a_menu`
SET `title`='维护区域列表',
    `route_path`='business-zones/index'
WHERE `id`=20929 OR `code`='region.zone.manage';

UPDATE `qixi_crm_a_menu`
SET `title`='维护代理人员'
WHERE `id`=20930 OR `code`='region.agent.manage';

UPDATE `qixi_crm_a_menu`
SET `parent_id`=10,
    `title`='区域代理',
    `icon`='ant-design:cluster-outlined',
    `sort`=4,
    `kind`='directory',
    `code`='region',
    `status`=1
WHERE `id`=20 OR `code`='region';

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21004,24,'region.agent_settings.manage','维护区域代理默认提成与申请表单','','business-zones/settings','button',1,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`);

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, 21004
FROM `qixi_crm_a_role` AS r
WHERE r.code = 'platform';
