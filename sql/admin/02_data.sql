USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_role` (`code`,`name`,`status`) VALUES
  ('platform','平台管理',1),('merchant','商户管理',1),('region','区域管理',1),
  ('customer_service','客服管理',1),('operations','运营管理',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`);

-- 平台、商户、区域、客服、运营共用同一套 Vben 应用，由角色决定可见菜单。
-- “统一后台”是系统身份，不是侧栏菜单；控制台是所有有权限账号的首个业务入口。
-- icon 使用 admin-platform 本地打包的 lucide 图标，不依赖公网 Iconify 服务。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  (2,0,'dashboard','控制台','lucide:layout-dashboard','/dashboard','page',1),

  (10,0,'merchant','商户管理','lucide:store','/merchant','directory',10),
  (11,10,'merchant.list','商户列表','lucide:store','/merchant/list','page',1),
  (12,10,'merchant.audit','入驻审核','lucide:badge-check','/merchant/audit','page',2),

  (20,0,'region','区域管理','lucide:map-pinned','/business-zones','directory',20),
  (21,20,'region.index','区域商圈','lucide:map-pinned','/region','page',1),
  (22,20,'region.agents','区域代理','lucide:users','/business-zones/agents','page',2),
  (23,20,'region.agent_review','代理审核','lucide:badge-check','/business-zones/agent-review','page',3),

  (30,0,'service','客服管理','lucide:messages-square','/service','page',30),

  (40,0,'product','商品管理','lucide:puzzle','/product','directory',40),
  (41,40,'product.category','商品分类','lucide:folder-tree','/product/category','page',1),
  (42,40,'product.brand','品牌管理','lucide:award','/product/brand','page',2),
  (43,40,'product.audit','商品审核','lucide:shield-check','/product/audit','page',3),

  (50,0,'order','订单管理','lucide:receipt-text','/order','directory',50),
  (51,50,'order.list','订单列表','lucide:receipt-text','/order/list','page',1),
  (52,50,'order.refund','售后退款','lucide:wallet','/order/refund','page',2),

  (60,0,'marketing','营销管理','lucide:activity','/marketing','directory',60),
  (61,60,'marketing.coupon','平台优惠券','lucide:award','/marketing/coupon','page',1),
  (62,60,'marketing.seckill','秒杀活动','lucide:radio-tower','/marketing/seckill','page',2),
  (63,60,'marketing.combination','拼团活动','lucide:users','/marketing/combination','page',3),
  (64,60,'marketing.presell','预售活动','lucide:receipt-text','/marketing/presell','page',4),
  (65,60,'marketing.spread','分销管理','lucide:git-branch','/marketing/spread','page',5),
  (66,60,'marketing.broadcast','直播管理','lucide:radio-tower','/marketing/broadcast','page',6),

  (70,0,'user','用户管理','lucide:users','/user','directory',70),
  (71,70,'user.label','用户标签','lucide:badge-check','/user/label','page',1),
  (72,70,'user.group','用户分组','lucide:users','/user/group','page',2),
  (73,70,'user.svip','会员等级','lucide:award','/user/svip','page',3),

  (80,0,'content','内容管理','lucide:images','/content','directory',80),
  (81,80,'content.notice','公告管理','lucide:messages-square','/content/notice','page',1),
  (82,80,'content.community','社区内容','lucide:images','/content/community','page',2),
  (83,80,'content.attachment','素材管理','lucide:images','/content/attachment','page',3),
  (84,80,'content.article','文章管理','lucide:pen-line','/cms/article','page',4),

  (90,0,'freight','物流配送','lucide:map-plus','/freight','directory',90),
  (91,90,'freight.express','快递公司','lucide:map-plus','/freight/express','page',1),

  (100,0,'accounts','财务管理','lucide:wallet','/accounts','directory',100),
  (101,100,'accounts.withdraw','提现审核','lucide:wallet','/accounts/withdraw','page',1),

  (110,0,'operations','运营装修','lucide:pen-line','/operations','directory',110),
  (111,110,'operations.diy','商城装修','lucide:pen-line','/operations/diy','page',1),

  (120,0,'setting','系统设置','lucide:settings','/setting','directory',120),
  (121,120,'setting.admin','管理员','lucide:user-round-cog','/setting/admin','page',1),
  (122,120,'setting.role','角色权限','lucide:shield-check','/setting/role','page',2),
  (123,120,'setting.menu','菜单管理','lucide:folder-tree','/setting/menu','page',3),
  (124,120,'setting.agreements','协议设置','lucide:receipt-text','/setting/agreements','page',4),
  (125,120,'setting.cloud_config','云服务配置','lucide:key-round','/setting/cloud-config','page',5),
  (126,120,'setting.sms','短信配置','lucide:messages-square','/setting/sms','page',6)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- v1 初始化曾用的临时扁平入口。它们会与新的目录结构重复，且没有 icon，
-- 仅清理这三个系统种子代码，不影响运营在后台维护的正式菜单。
DELETE rm
FROM `qixi_crm_a_role_menu` AS rm
INNER JOIN `qixi_crm_a_menu` AS m ON m.id = rm.menu_id
WHERE m.code IN ('merchant.review','region.manage','service.desk');
DELETE FROM `qixi_crm_a_menu`
WHERE `code` IN ('merchant.review','region.manage','service.desk');

-- 迁移旧种子：历史版本以“统一后台”作为根目录；现在改为真正的一级业务菜单。
DELETE rm
FROM `qixi_crm_a_role_menu` AS rm
INNER JOIN `qixi_crm_a_menu` AS m ON m.id = rm.menu_id
WHERE m.code = 'console';
DELETE FROM `qixi_crm_a_menu` WHERE `code` = 'console';

-- 平台拥有全部监管权限；其余四种账号在同一套后台中按职责显示模块。
-- 目录节点必须随子页面一起授权，避免 Vben 菜单树出现孤儿页面。
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
   OR (r.code = 'merchant' AND m.code IN (
      'dashboard','merchant','merchant.list','merchant.audit',
      'product','product.category','product.brand','product.audit',
      'order','order.list','order.refund','marketing','marketing.coupon','marketing.seckill',
      'marketing.combination','marketing.presell','marketing.spread','marketing.broadcast',
      'content','content.notice','content.community','content.attachment','content.article'))
   OR (r.code = 'region' AND m.code IN (
      'dashboard','merchant','merchant.list','merchant.audit',
      'region','region.index','region.agents','region.agent_review',
      'product','product.category','product.brand','product.audit','order','order.list','order.refund'))
   OR (r.code = 'customer_service' AND m.code IN (
      'dashboard','service'))
   OR (r.code = 'operations' AND m.code IN (
      'dashboard','marketing','marketing.coupon','marketing.seckill','marketing.combination',
      'marketing.presell','marketing.spread','marketing.broadcast','user','user.label','user.group','user.svip',
      'content','content.notice','content.community','content.attachment','content.article',
      'operations','operations.diy'));
