-- 平台「分销」侧栏对齐 CRMEB 嵌套结构（与「商品」同级缩进规范）
-- 二级：叶子与「分销等级」目录同级；三级挂在分销等级下
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_admin' < sql/admin/patch_promoter_menus.sql
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 一级「分销」
UPDATE `qixi_crm_a_menu`
SET `title`='分销',
    `icon`='ant-design:send-outlined',
    `route_path`='/promoter',
    `kind`='directory',
    `sort`=5,
    `parent_id`=0,
    `status`=1
WHERE `id`=220 OR `code`='promoter';

-- 旧扁平「分销管理」侧栏隐藏（功能由下方叶子承接；按钮权限仍挂原节点）
UPDATE `qixi_crm_a_menu`
SET `title`='分销管理',
    `icon`='ant-design:send-outlined',
    `route_path`='/marketing/spread',
    `kind`='page',
    `sort`=99,
    `parent_id`=220,
    `status`=0
WHERE `id`=65 OR `code`='marketing.spread';

-- 二级叶子 + 「分销等级」目录（sort 对齐侧栏观感：列表 → 等级 → 配置项 → 订单/说明）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (522,220,'promoter.user','分销员列表','lucide:users','/promoter/user','page',1,1),
  (1373,220,'promoter.brokerage','分销等级','lucide:layers','/brokerage','directory',2,1),
  (677,220,'promoter.bank','提现银行','lucide:landmark','/group/config/76','page',3,1),
  (685,220,'promoter.privilege','分销特权','lucide:crown','/group/config/75','page',4,1),
  (686,220,'promoter.poster','分销海报','lucide:image','/group/config/68','page',5,1),
  (731,220,'promoter.gift','分销礼包','lucide:gift','/promoter/gift','page',6,1),
  (1296,220,'promoter.commission','佣金说明','lucide:file-text','/promoter/commission','page',7,1),
  (9169,220,'promoter.order','分销订单','lucide:receipt','/promoter/orderList','page',8,1),
  (9368,220,'promoter.explain','分销说明','lucide:book-open','/promoter/retail','page',9,1),
  (5122,220,'promoter.config','分销配置','lucide:settings-2','/systemForm/Basics/distribution_tabs','page',10,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：分销等级子页
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1374,1373,'promoter.brokerage.level','分销员等级','lucide:badge-percent','/promoter/membership_level','page',1,1),
  (1375,1373,'promoter.brokerage.rule','等级规则','lucide:scale','/promoter/distribution','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 原查看权限挂到分销员列表
UPDATE `qixi_crm_a_menu`
SET `parent_id`=522,
    `title`='查看分销推广与佣金监管',
    `route_path`='promoter/user',
    `kind`='button',
    `sort`=1,
    `status`=1
WHERE `id`=20927 OR `code`='marketing.spread.read';

-- 平台角色补授权
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (220,522,1373,1374,1375,677,685,686,731,1296,9169,9368,5122,20927);
