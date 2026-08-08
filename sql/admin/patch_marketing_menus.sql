-- 平台「营销」侧栏对齐 CRMEB 嵌套结构（与「商品」「分销」同级缩进规范）
-- 二级目录/叶子同列；三级挂在各活动目录下；「余额充值」两子页同为 page 同级
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4 qixi_crm_admin' < sql/admin/patch_marketing_menus.sql
-- 执行后需重新登录平台后台以刷新菜单缓存。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 一级「营销」
UPDATE `qixi_crm_a_menu`
SET `title`='营销',
    `icon`='ant-design:flag-outlined',
    `route_path`='/marketing',
    `kind`='directory',
    `sort`=6,
    `parent_id`=0,
    `status`=1
WHERE `id`=60 OR `code`='marketing';

-- 旧扁平叶子：侧栏隐藏（由下方 CRMEB 树承接；保留 code/按钮挂载点）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=990 WHERE `id`=61 OR `code`='marketing.coupon';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=991 WHERE `id`=62 OR `code`='marketing.seckill';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=992 WHERE `id`=63 OR `code`='marketing.combination';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=993 WHERE `id`=64 OR `code`='marketing.presell';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=994 WHERE `id`=66 OR `code`='marketing.broadcast';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=995 WHERE `id`=67 OR `code`='marketing.assist';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=996 WHERE `id`=68 OR `code`='marketing.points';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=997 WHERE `id`=69 OR `code`='marketing.recharge';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=998 WHERE `id`=185 OR `code`='marketing.coupon.send_records';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=999 WHERE `id`=186 OR `code`='marketing.coupon.receipt_records';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1000 WHERE `id`=207 OR `code`='marketing.discounts';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1001 WHERE `id`=208 OR `code`='marketing.atmosphere';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1002 WHERE `id`=209 OR `code`='marketing.border';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1003 WHERE `id`=210 OR `code`='marketing.topic';
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=1004 WHERE `id`=211 OR `code`='marketing.application';

-- 二级：目录与叶子（sort 对齐截图顺序）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1657,60,'marketing.platform_coupon','平台优惠券','lucide:ticket','/marketing/platform_coupon','directory',1,1),
  (720,60,'marketing.store_coupon','商户优惠券','lucide:tickets','/marketing/coupon','directory',2,1),
  (780,60,'marketing.seckill.dir','秒杀','lucide:zap','/marketing/seckill','directory',3,1),
  (782,60,'marketing.broadcast.dir','直播','lucide:radio','/marketing2','directory',4,1),
  (1022,60,'marketing.presell.dir','预售','lucide:calendar-clock','/marketing/presell','directory',5,1),
  (1051,60,'marketing.assist.dir','助力','lucide:hand-heart','/assist','directory',6,1),
  (1135,60,'marketing.combination.dir','拼团','lucide:users','/marketing/combination','directory',7,1),
  (1289,60,'marketing.integral.dir','积分','lucide:coins','/marketing/integral','directory',8,1),
  (9007,60,'marketing.atmosphere.nav','活动氛围图','lucide:sparkles','/marketing/atmosphere/list','page',9,1),
  (9008,60,'marketing.border.nav','活动边框图','lucide:frame','/marketing/border/list','page',10,1),
  (1470,60,'marketing.topic.nav','专场列表','lucide:layout-template','/group/topic/94','page',11,1),
  (1629,60,'marketing.discounts.nav','优惠套餐','lucide:package','/marketing/discounts/list','page',12,1),
  (5126,60,'marketing.balance.dir','余额充值','lucide:wallet','/banlace','directory',13,1),
  (9217,60,'marketing.application.nav','报名活动','lucide:clipboard-list','/marketing/application/list','page',14,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：平台优惠券
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1658,1657,'marketing.platform_coupon.list','优惠券列表','lucide:list','/marketing/platform_coupon/list','page',1,1),
  (1659,1657,'marketing.platform_coupon.record','领取记录','lucide:history','/marketing/platform_coupon/couponRecord','page',2,1),
  (1662,1657,'marketing.platform_coupon.send','发送记录','lucide:send','/marketing/platform_coupon/couponSend','page',3,1),
  (1663,1657,'marketing.platform_coupon.help','使用说明','lucide:book-open','/marketing/platform_coupon/instructions','page',4,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：商户优惠券
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (721,720,'marketing.store_coupon.list','优惠券列表','lucide:list','/marketing/coupon/list','page',1,1),
  (734,720,'marketing.store_coupon.user','领取记录','lucide:history','/marketing/coupon/user','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：秒杀
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (779,780,'marketing.seckill.config','秒杀配置','lucide:settings-2','/marketing/seckill/seckillConfig','page',1,1),
  (794,780,'marketing.seckill.manage.page','秒杀管理','lucide:list','/marketing/seckill/list','page',2,1),
  (9287,780,'marketing.seckill.activity','秒杀活动','lucide:flame','/marketing/seckill/store_seckill/list','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：直播
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (781,782,'marketing.broadcast.studio','直播间管理','lucide:video','/marketing/studio/list','page',1,1),
  (783,782,'marketing.broadcast.goods','直播商品管理','lucide:shopping-bag','/marketing/broadcast/list','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：预售
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1023,1022,'marketing.presell.goods','预售商品','lucide:package','/marketing/presell/list','page',1,1),
  (1024,1022,'marketing.presell.agreement','预售协议','lucide:file-text','/marketing/presell/agreement','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：助力
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1095,1051,'marketing.assist.goods','活动商品','lucide:package','/marketing/assist/goods_list','page',1,1),
  (1096,1051,'marketing.assist.activity','助力活动','lucide:hand-heart','/marketing/assist/list','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：拼团
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1136,1135,'marketing.combination.set','拼团设置','lucide:settings-2','/marketing/combination/combination_set','page',1,1),
  (1137,1135,'marketing.combination.goods','拼团商品列表','lucide:package','/marketing/combination/combination_goods','page',2,1),
  (1138,1135,'marketing.combination.list','拼团活动列表','lucide:list','/marketing/combination/combination_list','page',3,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：积分
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1290,1289,'marketing.integral.config','积分配置','lucide:settings-2','/marketing/integral/config','page',1,1),
  (1291,1289,'marketing.integral.log','积分日志','lucide:scroll-text','/marketing/integral/log','page',2,1),
  (9118,1289,'marketing.integral.classify','商品分类','lucide:folder-tree','/marketing/integral/classify','page',3,1),
  (9119,1289,'marketing.integral.products','商品列表','lucide:list','/marketing/integral/proList','page',4,1),
  (9120,1289,'marketing.integral.orders','积分订单','lucide:receipt','/marketing/integral/orderList','page',5,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`=VALUES(`kind`),
  `sort`=VALUES(`sort`),
  `status`=1;

-- 三级：余额充值（两子页必须同为 page、同 parent，避免侧栏同级缩进错位）
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (667,5126,'marketing.balance.settings','余额设置','lucide:sliders-horizontal','/systemForm/Basics/balance','page',1,1),
  (687,5126,'marketing.balance.config','余额充值配置','lucide:badge-dollar-sign','/group/config/69','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=5126,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

-- 按钮权限挂到新页面节点（保留旧 code，便于既有 RBAC）
UPDATE `qixi_crm_a_menu` SET `parent_id`=1658, `route_path`='marketing/platform_coupon/list' WHERE `id`=20918 OR `code`='marketing.coupon.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1662, `route_path`='marketing/platform_coupon/couponSend' WHERE `id`=20980 OR `code`='marketing.coupon.send.read';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1659, `route_path`='marketing/platform_coupon/couponRecord' WHERE `id`=20981 OR `code`='marketing.coupon.record.read';
UPDATE `qixi_crm_a_menu` SET `parent_id`=794, `route_path`='marketing/seckill/list' WHERE `id`=20915 OR `code`='marketing.seckill.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1138, `route_path`='marketing/combination/combination_list' WHERE `id`=20916 OR `code`='marketing.combination.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1023, `route_path`='marketing/presell/list' WHERE `id`=20917 OR `code`='marketing.presell.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=781, `route_path`='marketing/studio/list' WHERE `id`=20922 OR `code`='marketing.broadcast.audit';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1096, `route_path`='marketing/assist/list' WHERE `id`=20928 OR `code`='marketing.assist.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=9119, `route_path`='marketing/integral/proList' WHERE `id`=20970 OR `code`='marketing.points.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=667, `route_path`='systemForm/Basics/balance' WHERE `id`=20971 OR `code`='marketing.recharge.manage';
UPDATE `qixi_crm_a_menu` SET `parent_id`=1629, `route_path`='marketing/discounts/list' WHERE `id` IN (21010,21011) OR `code` IN ('marketing.discounts.read','marketing.discounts.manage');
UPDATE `qixi_crm_a_menu` SET `parent_id`=9007, `route_path`='marketing/atmosphere/list' WHERE `id` IN (21012,21013) OR `code` IN ('marketing.atmosphere.read','marketing.atmosphere.manage');
UPDATE `qixi_crm_a_menu` SET `parent_id`=9008, `route_path`='marketing/border/list' WHERE `id` IN (21014,21015) OR `code` IN ('marketing.border.read','marketing.border.manage');
UPDATE `qixi_crm_a_menu` SET `parent_id`=1470, `route_path`='group/topic/94' WHERE `id` IN (21016,21017) OR `code` IN ('marketing.topic.read','marketing.topic.manage');
UPDATE `qixi_crm_a_menu` SET `parent_id`=9217, `route_path`='marketing/application/list' WHERE `id` IN (21018,21019) OR `code` IN ('marketing.application.read','marketing.application.manage');

-- 平台角色补授权（导航节点；按钮沿用原 role_menu）
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
  AND m.id IN (
    1657,1658,1659,1662,1663,
    720,721,734,
    780,779,794,9287,
    782,781,783,
    1022,1023,1024,
    1051,1095,1096,
    1135,1136,1137,1138,
    1289,1290,1291,9118,9119,9120,
    9007,9008,1470,1629,
    5126,667,687,
    9217
  );
