SET NAMES utf8mb4;
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
  (13,10,'merchant.category','商户分类','lucide:tags','/merchant/categories','page',3),
  (14,10,'merchant.grouping','店铺分组','lucide:network','/merchant/grouping','page',4),
  (15,10,'merchant.type','店铺类型','lucide:badge-plus','/merchant/types','page',5),
  (16,10,'merchant.deposit','店铺保证金','lucide:shield-dollar','/merchant/deposits','page',6),
  (17,10,'merchant.applyments','店铺分账申请','lucide:split','/merchant/applyments','page',7),

  (20,0,'region','区域管理','lucide:map-pinned','/business-zones','directory',20),
  (21,20,'region.index','区域商圈','lucide:map-pinned','/region','page',1),
  (22,20,'region.agents','区域代理','lucide:users','/business-zones/agents','page',2),
  (23,20,'region.agent_review','代理审核','lucide:badge-check','/business-zones/agent-review','page',3),
  (24,20,'region.agent_settings','代理设置','lucide:settings-2','/business-zones/settings','page',4),

  (30,0,'service','客服管理','lucide:messages-square','/service','page',30),

  (40,0,'product','商品管理','lucide:puzzle','/product','directory',40),
  (41,40,'product.category','商品分类','lucide:folder-tree','/product/category','page',1),
  (42,40,'product.brand','品牌管理','lucide:award','/product/brand','page',2),
  (43,40,'product.audit','商品审核','lucide:shield-check','/product/audit','page',3),
  (44,40,'product.label','商品标签','lucide:tags','/product/label','page',4),
  (45,40,'product.guarantee','保障服务','lucide:shield','/product/guarantee','page',5),
  (46,40,'product.parameter','平台商品参数','lucide:list-tree','/product/specs','page',6),
  (47,40,'product.comment','评论管理','lucide:message-square','/product/comment','page',7),

  (50,0,'order','订单管理','lucide:receipt-text','/order','directory',50),
  (51,50,'order.list','订单列表','lucide:receipt-text','/order/list','page',1),
  (52,50,'order.refund','售后退款','lucide:wallet','/order/refund','page',2),
  (53,50,'order.cancellation','取消/退款订单','lucide:ban','/order/cancellation','page',3),

  (60,0,'marketing','营销管理','lucide:activity','/marketing','directory',60),
  (61,60,'marketing.coupon','平台优惠券','lucide:award','/marketing/coupon','page',1),
  (62,60,'marketing.seckill','秒杀活动','lucide:radio-tower','/marketing/seckill','page',2),
  (63,60,'marketing.combination','拼团活动','lucide:users','/marketing/combination','page',3),
  (64,60,'marketing.presell','预售活动','lucide:receipt-text','/marketing/presell','page',4),
  (65,60,'marketing.spread','分销管理','lucide:git-branch','/marketing/spread','page',5),
  (66,60,'marketing.broadcast','直播管理','lucide:radio-tower','/marketing/broadcast','page',6),
  (67,60,'marketing.assist','好友助力','lucide:hand-heart','/marketing/assist','page',7),
  (68,60,'marketing.points','积分商城','lucide:badge-plus','/marketing/points','page',8),
  (69,60,'marketing.recharge','用户充值','lucide:circle-dollar-sign','/marketing/recharge','page',9),
  (185,60,'marketing.coupon.send_records','优惠券发送记录','lucide:ticket-check','/marketing/coupon/send-records','page',10),
  (186,60,'marketing.coupon.receipt_records','优惠券领取记录','lucide:ticket','/marketing/coupon/receipt-records','page',11),

  (70,0,'user','用户管理','lucide:users','/user','directory',70),
  (71,70,'user.label','用户标签','lucide:badge-check','/user/label','page',1),
  (72,70,'user.group','用户分组','lucide:users','/user/group','page',2),
  (73,70,'user.svip','会员等级','lucide:award','/user/svip','page',3),
  (74,70,'user.feedback','用户反馈','lucide:message-circle-warning','/user/feedback','directory',4),
  (75,74,'user.feedback.list','反馈列表','lucide:message-circle','/user/feedback/list','page',1),
  (76,74,'user.feedback.category','反馈分类','lucide:tags','/user/feedback/categories','page',2),
  (77,70,'user.list','用户列表','lucide:contact-round','/user/list','page',5),
  (78,70,'user.assets_adjustment','用户资产调整','lucide:badge-dollar-sign','/user/assets-adjustment','page',6),
  (79,70,'user.member_adjustment','用户会员等级调整','lucide:award','/user/member-adjustment','page',7),
  -- 170–179 避开内容管理已占用的 80–84，保证全新初始化不会因主键碰撞覆盖用户页面。
  (170,70,'user.coupon_operation','用户优惠券操作','lucide:ticket-plus','/user/coupon-operation','page',8),
  (171,70,'user.referrer_adjustment','用户推荐关系调整','lucide:git-branch-plus','/user/referrer-adjustment','page',9),
  (172,70,'user.group_assignment','用户分组归属','lucide:users-round','/user/group-assignment','page',10),
  (173,70,'user.label_assignment','用户标签归属','lucide:tags','/user/label-assignment','page',11),
  (174,70,'user.status_adjustment','用户启停','lucide:user-cog','/user/status-adjustment','page',12),
  (175,70,'user.create','新增用户','lucide:user-plus','/user/create','page',13),
  (176,70,'user.profile_maintenance','用户资料维护','lucide:contact','/user/profile-maintenance','page',14),
  (177,70,'user.password_reset','用户密码重置','lucide:key-round','/user/password-reset','page',15),
  (178,70,'user.promoter_assignment','批量设置推广员','lucide:badge-percent','/user/promoter-assignment','page',16),
  (179,70,'user.notification','用户站内图文通知','lucide:send','/user/notification','page',17),
  (180,70,'user.svip.plan','会员类型','lucide:badge-check','/user/member/type','page',18),
  (181,70,'user.svip.record','会员记录','lucide:scroll-text','/user/member/record','page',19),
  (182,70,'user.search_record','搜索记录','lucide:search','/user/search-record','page',20),
  (183,70,'user.member.level','会员等级管理','lucide:medal','/user/member/levels','page',21),
  (184,70,'user.svip.interest','会员权益','lucide:heart-handshake','/user/member/equity','page',22),

  (80,0,'content','内容管理','lucide:images','/content','directory',80),
  (81,80,'content.notice','公告管理','lucide:messages-square','/content/notice','page',1),
  (82,80,'content.community','社区内容','lucide:images','/content/community','page',2),
  (83,80,'content.attachment','素材管理','lucide:images','/content/attachment','page',3),
  (84,80,'content.article','文章管理','lucide:pen-line','/cms/article','page',4),
  (85,80,'content.community.category','社区分类','lucide:tags','/community/category','page',5),
  (86,80,'content.community.topic','社区话题','lucide:hash','/community/topic','page',6),
  (87,80,'content.community.list','社区内容','lucide:images','/community/list','page',7),
  (88,80,'content.community.reply','社区评论','lucide:message-square','/community/reply','page',8),
  (48,40,'product.price_description','价格说明','lucide:badge-dollar-sign','/product/priceDescription','page',8),
  (49,40,'product.activity_label','活动标签','lucide:badge','/product/activityLabel','page',9),

  (90,0,'freight','物流配送','lucide:map-plus','/freight','directory',90),
  (91,90,'freight.express','快递公司','lucide:map-plus','/freight/express','page',1),

  (100,0,'accounts','财务管理','lucide:wallet','/accounts','directory',100),
  (101,100,'accounts.withdraw','提现审核','lucide:wallet','/accounts/withdraw','page',1),
  (102,100,'accounts.user_assets','用户资产流水','lucide:wallet-cards','/accounts/user-assets','page',2),
  (103,100,'accounts.merchant_settlement','店铺结算监管','lucide:landmark','/accounts/merchant-settlement','page',3),
  (188,100,'accounts.invoice','发票管理','lucide:receipt','/accounts/invoices','page',4),

  (110,0,'operations','运营装修','lucide:pen-line','/operations','directory',110),
  (111,110,'operations.diy','商城装修','lucide:pen-line','/operations/diy','page',1),

  (120,0,'setting','系统设置','lucide:settings','/setting','directory',120),
  (121,120,'setting.admin','管理员','lucide:user-round-cog','/setting/admin','page',1),
  (122,120,'setting.role','角色权限','lucide:shield-check','/setting/role','page',2),
  (123,120,'setting.menu','菜单管理','lucide:folder-tree','/setting/menu','page',3),
  (124,120,'setting.agreements','协议设置','lucide:receipt-text','/setting/agreements','page',4),
  (125,120,'setting.cloud_config','云服务配置','lucide:key-round','/setting/cloud-config','page',5),
  (126,120,'setting.sms','短信配置','lucide:messages-square','/setting/sms','page',6)
  ,(187,120,'setting.operation_log','操作日志','lucide:history','/setting/system-log','page',7)
  ,(189,120,'setting.login_log','登录日志','lucide:log-in','/setting/login-log','page',8)
  ,(190,120,'setting.shop','商城设置','lucide:store','/setting/shop','page',9)
  ,(191,120,'setting.pay','支付设置','lucide:credit-card','/setting/pay','page',10)
  ,(130,0,'app','应用管理','lucide:smartphone','/app','directory',130)
  ,(131,130,'app.wechat','公众号','lucide:message-circle','/app/wechat','page',1)
  ,(192,0,'statistic','数据统计','lucide:bar-chart-3','/statistic','directory',15)
  ,(193,192,'statistic.order','订单统计','lucide:receipt-text','/statistic/order','page',1)
  ,(194,192,'statistic.product','商品统计','lucide:package','/statistic/product','page',2)
  ,(195,192,'statistic.user','用户统计','lucide:users','/statistic/user','page',3)
  ,(196,120,'setting.storage','存储配置','lucide:hard-drive','/setting/storage','page',11)
  ,(197,0,'maintain','系统维护','lucide:wrench','/maintain','directory',197)
  ,(198,197,'maintain.cache','清除缓存','lucide:eraser','/maintain/cache','page',1)
  ,(199,197,'maintain.backup','数据备份','lucide:database-backup','/maintain/dataBackup','page',2)
  ,(200,197,'maintain.group_data','组合数据','lucide:layers','/group/list','page',3)
  ,(201,197,'maintain.hot_search','热门搜索','lucide:search','/group/config/67','page',4)
  ,(202,130,'app.routine','小程序','lucide:smartphone','/app/routine','page',2)
  ,(203,130,'app.wechat_reply','微信回复','lucide:message-square-reply','/admin/app/wechat/reply','page',3)
  ,(204,130,'app.wechat_menus','微信菜单','lucide:menu','/app/wechat/menus','page',4)
  ,(205,130,'app.wechat_template','模板消息','lucide:mail','/app/wechat/template','page',5)
  ,(206,130,'app.wechat_news','图文消息','lucide:newspaper','/app/wechat/newsCategory','page',6)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- 统一后台按钮节点必须使用 qixi_crm_a_* RBAC；运营写操作不再依赖旧系统菜单表。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  (34,52,'order.refund.approve','审核退款','','order/refund','button',1),
  (35,52,'order.refund.reject','拒绝退款','','order/refund','button',2),
  (20903,43,'product.audit.submit','审核商品','','product/audit','button',1),
  (20904,41,'product.category.manage','维护商品分类','','product/category','button',1),
  (20905,42,'product.brand.manage','维护商品品牌','','product/brand','button',1),
  (20910,84,'content.article.manage','维护文章','','cms/article','button',1),
  (20911,84,'content.article_category.manage','维护文章分类','','cms/article','button',2),
  (20912,81,'content.notice.manage','维护公告','','content/notice','button',1),
  (20913,124,'setting.agreement.manage','维护协议','','setting/agreements','button',1),
  (20914,111,'operations.diy.manage','维护商城装修','','operations/diy','button',1),
  (20915,62,'marketing.seckill.manage','维护秒杀活动','','marketing/seckill','button',1),
  (20916,63,'marketing.combination.manage','维护拼团活动','','marketing/combination','button',1),
  (20917,64,'marketing.presell.manage','维护预售活动','','marketing/presell','button',1),
  (20918,61,'marketing.coupon.manage','维护平台优惠券','','marketing/coupon','button',1),
  (20919,83,'content.attachment.manage','维护素材库','','content/attachment','button',1),
  (20920,82,'content.community.audit','审核社区内容','','content/community','button',1),
  (20921,82,'content.community.delete','删除社区内容','','content/community','button',2),
  (20922,66,'marketing.broadcast.audit','审核直播间','','marketing/broadcast','button',1),
  (20923,101,'accounts.withdraw.review','审核用户提现','','accounts/withdraw','button',1),
  (20924,102,'accounts.user_assets.read','查看用户资产流水','','accounts/user-assets','button',1),
  (20925,103,'accounts.merchant_settlement.read','查看店铺结算监管投影','','accounts/merchant-settlement','button',1),
  (20926,103,'accounts.merchant_settlement.review','审核并登记店铺结算凭证','','accounts/merchant-settlement','button',2),
  (20927,65,'marketing.spread.read','查看分销推广与佣金监管','','marketing/spread','button',1),
  (20928,67,'marketing.assist.manage','维护好友助力活动','','marketing/assist','button',1),
  (20929,21,'region.zone.manage','维护区域商圈','','region','button',1),
  (20930,22,'region.agent.manage','维护区域代理申请','','business-zones/agents','button',1),
  (20931,23,'region.agent.review','审核区域代理申请','','business-zones/agent-review','button',1),
  (20932,91,'freight.express.manage','维护快递公司与行政区划','','freight/express','button',1),
  (20933,11,'merchant.status.manage','启停店铺经营状态','','merchant/list','button',1),
  (20934,12,'merchant.intention.audit','审核商户入驻申请','','merchant/audit','button',1),
  (20935,12,'merchant.intention.assign_region','分配商户入驻审核区域','','merchant/audit','button',2),
  (20936,126,'setting.sms.manage','维护无密钥短信 stub 配置','','setting/sms','button',1),
  (20937,13,'merchant.category.manage','维护商户分类与佣金比例','','merchant/categories','button',1),
  (20938,14,'merchant.group.manage','维护店铺分组、关联店铺与装修模板','','merchant/grouping','button',1)
  ,(20939,15,'merchant.type.manage','维护店铺类型、保证金规则与授权菜单','','merchant/types','button',1)
  ,(20940,16,'merchant.deposit.review','维护保证金、扣减、退款审核及打款登记','','merchant/deposits','button',1)
  ,(20941,17,'merchant.profitsharing.review','审核店铺分账申请及维护审核备注','','merchant/applyments','button',1)
  ,(20942,44,'product.label.manage','维护平台商品标签','','product/label','button',1)
  ,(20943,45,'product.guarantee.manage','维护平台保障服务','','product/guarantee','button',1)
  ,(20944,46,'product.parameter.manage','维护平台商品参数模板','','product/specs','button',1)
  ,(20945,47,'product.comment.review','审核或隐藏商品评论','','product/comment','button',1)
  ,(20946,47,'product.comment.virtual.manage','新增或编辑虚拟商品评论','','product/comment','button',2)
  ,(20947,47,'product.comment.sort','调整虚拟商品评论排序','','product/comment','button',3)
  ,(20948,47,'product.comment.delete','删除虚拟商品评论','','product/comment','button',4)
  ,(20949,75,'user.feedback.read','查看用户反馈','','user/feedback/list','button',1)
  ,(20950,75,'user.feedback.manage','回复、关闭或删除用户反馈','','user/feedback/list','button',2)
  ,(20951,76,'user.feedback.category.manage','维护反馈分类','','user/feedback/categories','button',1)
  ,(20952,77,'user.list.read','查看脱敏用户列表','','user/list','button',1)
  ,(20953,77,'user.list.manage','人工调整用户余额、积分与会员等级','','user/list','button',2)
  ,(20954,78,'user.asset.adjust','提交用户余额或积分人工调整','','user/assets-adjustment','button',1)
  ,(20955,79,'user.member.adjust','提交用户会员等级人工调整','','user/member-adjustment','button',1)
  ,(20956,170,'user.coupon.manage','提交用户优惠券发放或撤销','','user/coupon-operation','button',1)
  ,(20957,171,'user.referrer.manage','提交用户推荐关系调整','','user/referrer-adjustment','button',1)
  ,(20958,172,'user.group.assign','提交用户单个或批量分组调整','','user/group-assignment','button',1)
  ,(20959,173,'user.label.assign','提交用户单个或批量标签调整','','user/label-assignment','button',1)
  ,(20960,174,'user.status.manage','提交用户启用或停用','','user/status-adjustment','button',1)
  ,(20961,175,'user.create.execute','创建本地 PC 用户','','user/create','button',1)
  ,(20962,176,'user.profile.manage','维护用户公开资料','','user/profile-maintenance','button',1)
  ,(20963,177,'user.password.reset','重置用户 PC 登录密码','','user/password-reset','button',1)
  ,(20964,178,'user.promoter.manage','批量启用或停用推广员资格','','user/promoter-assignment','button',1)
  ,(20965,179,'user.notification.send','发送用户站内图文通知','','user/notification','button',1)
  ,(20966,77,'user.list.export','导出脱敏用户信息','','user/list','button',3)
  ,(20967,22,'region.agent.password.reset','重置区域代理后台密码','','business-zones/agents','button',2)
  ,(20968,52,'order.refund.log','查看退款状态流转日志','','order/refund','button',3)
  ,(20969,52,'order.refund.export','导出受数据范围约束的退款监管清单','','order/refund','button',4)
  ,(20970,68,'marketing.points.manage','维护积分商品并查看积分订单','','marketing/points','button',1)
  ,(20971,69,'marketing.recharge.manage','维护充值计划并查看充值订单','','marketing/recharge','button',1)
  ,(20972,73,'user.svip.manage','查看并维护用户 SVIP 状态','','user/svip','button',1)
  ,(20973,180,'user.svip.plan.manage','维护可售会员类型','','user/member/type','button',1)
  ,(20974,181,'user.svip.record.read','查看会员购买记录与统计','','user/member/record','button',1)
  ,(20975,182,'user.search_record.read','查看用户搜索记录','','user/search-record','button',1)
  ,(20976,182,'user.search_record.clear','按用户清理搜索记录','','user/search-record','button',2)
  ,(20977,182,'user.search_record.export','导出用户搜索记录','','user/search-record','button',3)
  ,(20978,183,'user.member.level.manage','维护普通会员等级、规则和权益','','user/member/levels','button',1)
  ,(20979,184,'user.svip.interest.manage','维护付费会员权益','','user/member/equity','button',1)
  ,(20980,185,'marketing.coupon.send.read','查看平台人工发券与撤销审计','','marketing/coupon/send-records','button',1)
  ,(20981,186,'marketing.coupon.record.read','查看用户优惠券领取与使用状态','','marketing/coupon/receipt-records','button',1)
  ,(20982,187,'setting.operation_log.read','查看统一后台成功操作日志','','setting/system-log','button',1)
  ,(20983,188,'accounts.invoice.read','查看订单发票监管记录','','accounts/invoices','button',1)
  ,(20984,190,'setting.shop.manage','维护商城基础设置','','setting/shop','button',1)
  ,(20985,191,'setting.pay.manage','维护支付方式开关','','setting/pay','button',1)
  ,(20986,53,'order.cancellation.read','查看取消或退款订单','','order/cancellation','button',1)
  ,(20987,131,'app.wechat.manage','维护公众号基础开关','','app/wechat','button',1)
  ,(20988,196,'setting.storage.manage','维护存储配置开关','','setting/storage','button',1)
  ,(20989,198,'maintain.cache.manage','提交缓存清理','','maintain/cache','button',1)
  ,(20990,199,'maintain.backup.manage','维护备份记录 stub','','maintain/dataBackup','button',1)
  ,(20991,200,'maintain.group_data.manage','维护组合数据 stub','','group/list','button',1)
  ,(20992,201,'maintain.hot_search.manage','维护热门搜索 stub','','group/config/67','button',1)
  ,(20993,202,'app.routine.manage','维护小程序基础开关','','app/routine','button',1)
  ,(20994,203,'app.wechat_reply.manage','维护微信回复开关','','admin/app/wechat/reply','button',1)
  ,(20995,204,'app.wechat_menus.manage','维护微信菜单开关','','app/wechat/menus','button',1)
  ,(20996,205,'app.wechat_template.manage','维护模板消息开关','','app/wechat/template','button',1)
  ,(20997,206,'app.wechat_news.manage','维护图文消息开关','','app/wechat/newsCategory','button',1)
  ,(20998,77,'user.setup.manage','维护用户注册设置','','user/setup_user','button',4)
  ,(20999,101,'accounts.transfer_settings.manage','维护转账监管设置','','accounts/settings','button',2)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`title`=VALUES(`title`),`route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

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
-- 尚未完成服务端数据隔离的历史入口必须同时撤销旧授权；INSERT IGNORE 不能回收
-- 已存在的角色菜单，留下它们会造成菜单可见但请求被默认拒绝的假入口。
DELETE rm
FROM `qixi_crm_a_role_menu` AS rm
INNER JOIN `qixi_crm_a_role` AS r ON r.id = rm.role_id
INNER JOIN `qixi_crm_a_menu` AS m ON m.id = rm.menu_id
WHERE (r.code = 'merchant' AND m.code IN (
  'marketing','marketing.coupon','marketing.seckill','marketing.combination','marketing.presell','marketing.spread','marketing.broadcast','marketing.assist','marketing.points','marketing.recharge',
  'content','content.notice','content.community','content.community.category','content.community.topic','content.community.list','content.community.reply','content.attachment','content.article'
)) OR (r.code = 'region' AND m.code IN ('region','region.index','region.agents','region.agent_review'))
  OR (r.code = 'operations' AND m.code IN ('dashboard','user','user.label','user.group','user.svip.plan','user.svip.record','user.svip.interest'));

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
   OR (r.code = 'merchant' AND m.code IN (
      'dashboard','merchant','merchant.list','merchant.audit',
      'product','product.audit',
      'order','order.list','order.refund','order.cancellation'))
   OR (r.code = 'region' AND m.code IN (
      'dashboard','merchant','merchant.list','merchant.audit',
      'product','product.audit','order','order.list','order.refund','order.cancellation',
      'accounts','accounts.merchant_settlement','accounts.merchant_settlement.read','merchant.intention.audit'))
   OR (r.code = 'customer_service' AND m.code IN (
      'dashboard','service',
      'user','user.feedback','user.feedback.list','user.feedback.read','user.feedback.manage'))
   OR (r.code = 'operations' AND m.code IN (
      'marketing','marketing.coupon','marketing.seckill','marketing.combination',
      'marketing.presell','marketing.spread','marketing.broadcast','marketing.assist','marketing.points','marketing.recharge','marketing.spread.read',
      'user','user.svip.plan','user.svip.record','user.svip.interest','user.svip.plan.manage','user.svip.record.read','user.svip.interest.manage',
      'content','content.notice','content.community','content.community.category','content.community.topic','content.community.list','content.community.reply','content.attachment','content.article',
      'operations','operations.diy','setting','setting.agreements',
      'content.article.manage','content.article_category.manage','content.notice.manage','setting.agreement.manage',
      'operations.diy.manage','marketing.seckill.manage','marketing.combination.manage','marketing.presell.manage','marketing.coupon.manage','marketing.assist.manage','marketing.points.manage','marketing.recharge.manage',
      'content.attachment.manage','content.community.audit','content.community.delete','marketing.broadcast.audit'));
