SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_role` (`code`,`name`,`status`,`role_type`) VALUES
  ('platform','平台管理',1,'platform'),('merchant','商户管理',1,'merchant'),('region','区域管理',1,'region'),
  ('customer_service','客服管理',1,'platform'),('operations','运营管理',1,'platform')
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `status`=VALUES(`status`),
  `role_type`=VALUES(`role_type`);

-- 平台、商户、区域、客服、运营共用同一套 Vben 应用，由角色决定可见菜单。
-- “统一后台”是系统身份，不是侧栏菜单；控制台是所有有权限账号的首个业务入口。
-- 顶层顺序对齐 CRMEB 平台侧栏（sort ASC）；顶层 icon 用离线 ant-design:*。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  -- 首页
  (1,0,'home','首页','ant-design:home-outlined','/','directory',1),
  (2,1,'dashboard','控制台','lucide:layout-dashboard','/dashboard','page',1),
  (216,1,'data.screen','数据大屏','lucide:monitor','/data-screen/index','page',2),
  (194,1,'statistic.product','商品统计','lucide:package','/statistic/product','page',3),
  (193,1,'statistic.order','订单统计','lucide:receipt-text','/statistic/order','page',4),
  (195,1,'statistic.user','用户统计','lucide:users','/statistic/user','page',5),

  -- 店铺：店铺管理 / 店铺设置 / 商户管理 / 区域代理
  (10,0,'store','店铺','ant-design:shop-outlined','/store','directory',2),
  (18,10,'store.manage','店铺管理','ant-design:shop-outlined','/mer/mer','directory',1),
  (11,18,'merchant.list','店铺列表','lucide:store','/merchant/list','page',1),
  (13,18,'merchant.category','店铺分类','lucide:folder-tree','/merchant/classify','page',2),
  (14,18,'merchant.grouping','店铺分组','lucide:git-branch','/merchant/grouping','page',3),
  (15,18,'merchant.type','店铺类型','lucide:award','/merchant/type','page',4),
  (12,18,'merchant.audit','店铺入驻申请','lucide:badge-check','/merchant/application','page',5),
  (17,18,'merchant.applyments','店铺分账申请','lucide:wallet','/merchant/applyments','page',6),
  (19,10,'store.settings','店铺设置','ant-design:setting-outlined','/mer/store','directory',2),
  (16,19,'merchant.deposit','店铺保证金','lucide:shield-check','/merchant/deposit_list','page',1),
  (213,19,'store.margin_config','保证金配置','lucide:key-round','/systemForm/Basics/margin','page',2),
  (214,19,'store.menu','店铺菜单','lucide:folder-tree','/merchant/system','page',3),
  (215,19,'store.description','说明提示','lucide:receipt-text','/merchant/type/description','page',4),
  (25,10,'merchant.mgmt','商户管理','ant-design:team-outlined','/merchant-mgmt','directory',3),
  (26,25,'merchant.mgmt.list','商户列表','lucide:store','/merchant/index','page',1),
  (27,25,'merchant.mgmt.review','店铺入驻申请','lucide:badge-check','/merchant/review','page',2),
  (28,25,'merchant.mgmt.admins','商户管理员','lucide:user-round-cog','/merchant/admin-list','page',3),
  (29,25,'merchant.mgmt.settings','商户设置','lucide:settings','/merchant/apply-setting','page',4),
  -- 区域代理侧栏仅 3 页：区域列表 / 代理人员 / 代理设置（代理审核页软隐藏，路由与按钮权限保留）
  (20,10,'region','区域代理','ant-design:cluster-outlined','/business-zones','directory',4),
  (21,20,'region.index','区域列表','lucide:map-pinned','/business-zones/index','page',1),
  (22,20,'region.agents','代理人员','lucide:users','/business-zones/agents','page',2),
  (23,20,'region.agent_review','代理审核','lucide:badge-check','/business-zones/agent-review','page',99),
  (24,20,'region.agent_settings','代理设置','lucide:settings-2','/business-zones/settings','page',3),

  -- 商品（对齐 CRMEB / 图片1：品牌管理、商品参数为二级目录）
  (40,0,'product','商品','ant-design:shopping-outlined','/product','directory',3),
  (43,40,'product.audit','商品管理','lucide:shield-check','/product/audit','page',1),
  (41,40,'product.category','商品分类','lucide:folder-tree','/product/category','page',2),
  (42,40,'product.brand','品牌管理','lucide:award','/product/brand','directory',3),
  (54,42,'product.brand.category','品牌分类','lucide:folder-tree','/product/band/brandClassify','page',1),
  (55,42,'product.brand.list','品牌列表','lucide:award','/product/band/brandList','page',2),
  (47,40,'product.comment','评论管理','lucide:message-square','/product/comment','page',4),
  (45,40,'product.guarantee','保障服务','lucide:shield','/product/guarantee','page',5),
  (44,40,'product.label','商品标签','lucide:tags','/product/label','page',6),
  (46,40,'product.parameter','商品参数','lucide:list-tree','/product/specsMain','directory',7),
  (56,46,'product.parameter.store','店铺商品参数','lucide:store','/product/merSpecs','page',1),
  (57,46,'product.parameter.platform','平台商品参数','lucide:list-tree','/product/specs','page',2),
  (48,40,'product.price_description','价格说明','lucide:badge-dollar-sign','/product/priceDescription','page',8),
  -- 活动标签：功能页保留，侧栏按图片1软隐藏（status 默认 1，见下方 UPDATE）
  (49,40,'product.activity_label','活动标签','lucide:badge','/product/activityLabel','page',99),

  -- 订单
  (50,0,'order','订单','ant-design:file-text-outlined','/order','directory',4),
  (51,50,'order.list','订单列表','lucide:receipt-text','/order/list','page',1),
  (52,50,'order.refund','退款订单','lucide:wallet','/order/refund','page',2),
  (53,50,'order.cancellation','核销记录','ant-design:audit-outlined','/order/cancellation','page',3),

  -- 分销（独立一级；嵌套结构见 patch_promoter_menus.sql）
  (220,0,'promoter','分销','ant-design:send-outlined','/promoter','directory',5),
  -- 65 保留为兼容页（侧栏由 patch 软隐藏）；正式导航为 522/1373 等
  (65,220,'marketing.spread','分销管理','ant-design:send-outlined','/marketing/spread','page',99),
  (522,220,'promoter.user','分销员列表','lucide:users','/promoter/user','page',1),
  (1373,220,'promoter.brokerage','分销等级','lucide:layers','/brokerage','directory',2),
  (1374,1373,'promoter.brokerage.level','分销员等级','lucide:badge-percent','/promoter/membership_level','page',1),
  (1375,1373,'promoter.brokerage.rule','等级规则','lucide:scale','/promoter/distribution','page',2),
  (677,220,'promoter.bank','提现银行','lucide:landmark','/group/config/76','page',3),
  (685,220,'promoter.privilege','分销特权','lucide:crown','/group/config/75','page',4),
  (686,220,'promoter.poster','分销海报','lucide:image','/group/config/68','page',5),
  (731,220,'promoter.gift','分销礼包','lucide:gift','/promoter/gift','page',6),
  (1296,220,'promoter.commission','佣金说明','lucide:file-text','/promoter/commission','page',7),
  (9169,220,'promoter.order','分销订单','lucide:receipt','/promoter/orderList','page',8),
  (9368,220,'promoter.explain','分销说明','lucide:book-open','/promoter/retail','page',9),
  (5122,220,'promoter.config','分销配置','lucide:settings-2','/systemForm/Basics/distribution_tabs','page',10),
  (21051,5122,'promoter.config.manage','保存分销配置','','systemForm/Basics/distribution_tabs','button',1),

  -- 营销（嵌套结构见 patch_marketing_menus.sql；下列为正式导航树）
  (60,0,'marketing','营销','ant-design:flag-outlined','/marketing','directory',6),
  (1657,60,'marketing.platform_coupon','平台优惠券','lucide:ticket','/marketing/platform_coupon','directory',1),
  (1658,1657,'marketing.platform_coupon.list','优惠券列表','lucide:list','/marketing/platform_coupon/list','page',1),
  (1659,1657,'marketing.platform_coupon.record','领取记录','lucide:history','/marketing/platform_coupon/couponRecord','page',2),
  (1662,1657,'marketing.platform_coupon.send','发送记录','lucide:send','/marketing/platform_coupon/couponSend','page',3),
  (1663,1657,'marketing.platform_coupon.help','使用说明','lucide:book-open','/marketing/platform_coupon/instructions','page',4),
  (720,60,'marketing.store_coupon','商户优惠券','lucide:tickets','/marketing/coupon','directory',2),
  (721,720,'marketing.store_coupon.list','优惠券列表','lucide:list','/marketing/coupon/list','page',1),
  (734,720,'marketing.store_coupon.user','领取记录','lucide:history','/marketing/coupon/user','page',2),
  (780,60,'marketing.seckill.dir','秒杀','lucide:zap','/marketing/seckill','directory',3),
  (779,780,'marketing.seckill.config','秒杀配置','lucide:settings-2','/marketing/seckill/seckillConfig','page',1),
  (794,780,'marketing.seckill.manage.page','秒杀管理','lucide:list','/marketing/seckill/list','page',2),
  (9287,780,'marketing.seckill.activity','秒杀活动','lucide:flame','/marketing/seckill/store_seckill/list','page',3),
  (782,60,'marketing.broadcast.dir','直播','lucide:radio','/marketing2','directory',4),
  (781,782,'marketing.broadcast.studio','直播间管理','lucide:video','/marketing/studio/list','page',1),
  (783,782,'marketing.broadcast.goods','直播商品管理','lucide:shopping-bag','/marketing/broadcast/list','page',2),
  (1022,60,'marketing.presell.dir','预售','lucide:calendar-clock','/marketing/presell','directory',5),
  (1023,1022,'marketing.presell.goods','预售商品','lucide:package','/marketing/presell/list','page',1),
  (1024,1022,'marketing.presell.agreement','预售协议','lucide:file-text','/marketing/presell/agreement','page',2),
  (1051,60,'marketing.assist.dir','助力','lucide:hand-heart','/assist','directory',6),
  (1095,1051,'marketing.assist.goods','活动商品','lucide:package','/marketing/assist/goods_list','page',1),
  (1096,1051,'marketing.assist.activity','助力活动','lucide:hand-heart','/marketing/assist/list','page',2),
  (1135,60,'marketing.combination.dir','拼团','lucide:users','/marketing/combination','directory',7),
  (1136,1135,'marketing.combination.set','拼团设置','lucide:settings-2','/marketing/combination/combination_set','page',1),
  (1137,1135,'marketing.combination.goods','拼团商品列表','lucide:package','/marketing/combination/combination_goods','page',2),
  (1138,1135,'marketing.combination.list','拼团活动列表','lucide:list','/marketing/combination/combination_list','page',3),
  (1289,60,'marketing.integral.dir','积分','lucide:coins','/marketing/integral','directory',8),
  (1290,1289,'marketing.integral.config','积分配置','lucide:settings-2','/marketing/integral/config','page',1),
  (1291,1289,'marketing.integral.log','积分日志','lucide:scroll-text','/marketing/integral/log','page',2),
  (21052,1291,'marketing.integral.log.read','查看积分日志','','marketing/integral/log','button',1),
  (9118,1289,'marketing.integral.classify','商品分类','lucide:folder-tree','/marketing/integral/classify','page',3),
  (21053,9118,'marketing.integral.classify.manage','管理积分商品分类','','marketing/integral/classify','button',1),
  (9119,1289,'marketing.integral.products','商品列表','lucide:list','/marketing/integral/proList','page',4),
  (9120,1289,'marketing.integral.orders','积分订单','lucide:receipt','/marketing/integral/orderList','page',5),
  (21054,9120,'marketing.integral.orders.read','查看积分订单','','marketing/integral/orderList','button',1),
  (21055,9120,'marketing.integral.orders.manage','管理积分订单','','marketing/integral/orderList','button',2),
  (9007,60,'marketing.atmosphere.nav','活动氛围图','lucide:sparkles','/marketing/atmosphere/list','page',9),
  (9008,60,'marketing.border.nav','活动边框图','lucide:frame','/marketing/border/list','page',10),
  (1470,60,'marketing.topic.nav','专场列表','lucide:layout-template','/group/topic/94','page',11),
  (1629,60,'marketing.discounts.nav','优惠套餐','lucide:package','/marketing/discounts/list','page',12),
  (5126,60,'marketing.balance.dir','余额充值','lucide:wallet','/banlace','directory',13),
  (667,5126,'marketing.balance.settings','余额设置','lucide:sliders-horizontal','/systemForm/Basics/balance','page',1),
  (21056,667,'marketing.balance.settings.read','查看余额设置','','systemForm/Basics/balance','button',1),
  (21057,667,'marketing.balance.settings.manage','保存余额设置','','systemForm/Basics/balance','button',2),
  (687,5126,'marketing.balance.config','余额充值设置','lucide:badge-dollar-sign','/group/config/69','page',2),
  (21058,687,'marketing.balance.config.read','查看余额充值设置','','group/config/69','button',1),
  (21059,687,'marketing.balance.config.manage','管理余额充值设置','','group/config/69','button',2),
  (9217,60,'marketing.application.nav','报名活动','lucide:clipboard-list','/marketing/application/list','page',14),
  -- 旧扁平营销页兼容占位（status 由 patch_marketing_menus 软隐藏）
  (61,60,'marketing.coupon','平台优惠券(旧)','lucide:award','/marketing/coupon','page',990),
  (62,60,'marketing.seckill','秒杀活动(旧)','lucide:radio-tower','/marketing/seckill','page',991),
  (63,60,'marketing.combination','拼团活动(旧)','lucide:users','/marketing/combination','page',992),
  (64,60,'marketing.presell','预售活动(旧)','lucide:receipt-text','/marketing/presell','page',993),
  (66,60,'marketing.broadcast','直播管理(旧)','lucide:radio-tower','/marketing/broadcast','page',994),
  (67,60,'marketing.assist','好友助力(旧)','lucide:hand-heart','/marketing/assist','page',995),
  (68,60,'marketing.points','积分商城(旧)','lucide:badge-plus','/marketing/points','page',996),
  (69,60,'marketing.recharge','用户充值(旧)','lucide:circle-dollar-sign','/marketing/recharge','page',997),
  (185,60,'marketing.coupon.send_records','优惠券发送记录(旧)','lucide:ticket-check','/marketing/coupon/send-records','page',998),
  (186,60,'marketing.coupon.receipt_records','优惠券领取记录(旧)','lucide:ticket','/marketing/coupon/receipt-records','page',999),
  (207,60,'marketing.discounts','优惠套餐(旧)','lucide:package','/marketing/discounts/list','page',1000),
  (208,60,'marketing.atmosphere','活动氛围(旧)','lucide:sparkles','/marketing/atmosphere/list','page',1001),
  (209,60,'marketing.border','活动边框(旧)','lucide:frame','/marketing/border/list','page',1002),
  (210,60,'marketing.topic','活动专题(旧)','lucide:layout-template','/group/topic/94','page',1003),
  (211,60,'marketing.application','活动报名(旧)','lucide:clipboard-list','/marketing/application/list','page',1004),

  -- 用户（嵌套见 patch_user_menus.sql）
  (70,0,'user','用户','ant-design:user-outlined','/user','directory',7),
  (77,70,'user.list','用户列表','lucide:contact-round','/user/list','page',1),
  (72,70,'user.group','用户分组','lucide:users','/user/group','page',2),
  (71,70,'user.label','用户标签','lucide:tags','/user/label','page',3),
  (74,70,'user.feedback','用户反馈','lucide:message-circle-warning','/user/feedback','directory',4),
  (75,74,'user.feedback.list','反馈列表','lucide:message-circle','/user/feedback/list','page',1),
  (76,74,'user.feedback.category','反馈分类','lucide:folder-tree','/user/feedback/categories','page',2),
  (182,70,'user.search_record','搜索记录','lucide:search','/user/search-record','page',5),
  (530,70,'user.level.dir','用户等级','lucide:medal','/user/member','directory',6),
  (183,530,'user.member.level','等级管理','lucide:medal','/user/member/list','page',1),
  (531,530,'user.level.description','等级说明','lucide:file-text','/user/member/description','page',2),
  (532,70,'user.setup','用户设置','lucide:settings-2','/user/setup_user','page',7),
  (73,70,'user.svip','付费会员','lucide:crown','/user/svip','directory',8),
  (180,73,'user.svip.plan','会员类型','lucide:badge-check','/user/member/type','page',1),
  (184,73,'user.svip.interest','会员权益','lucide:heart-handshake','/user/member/equity','page',2),
  (181,73,'user.svip.record','会员记录','lucide:scroll-text','/user/member/record','page',3),
  (533,73,'user.svip.agreement','会员协议','lucide:file-text','/user/member/vipAgreement','page',4),
  -- 运营辅助页兼容（patch 软隐藏）
  (78,70,'user.assets_adjustment','用户资产调整','lucide:badge-dollar-sign','/user/assets-adjustment','page',990),
  (79,70,'user.member_adjustment','用户会员等级调整','lucide:award','/user/member-adjustment','page',991),
  (170,70,'user.coupon_operation','用户优惠券操作','lucide:ticket-plus','/user/coupon-operation','page',992),
  (171,70,'user.referrer_adjustment','用户推荐关系调整','lucide:git-branch-plus','/user/referrer-adjustment','page',993),
  (172,70,'user.group_assignment','用户分组归属','lucide:users-round','/user/group-assignment','page',994),
  (173,70,'user.label_assignment','用户标签归属','lucide:tags','/user/label-assignment','page',995),
  (174,70,'user.status_adjustment','用户启停','lucide:user-cog','/user/status-adjustment','page',996),
  (175,70,'user.create','新增用户','lucide:user-plus','/user/create','page',997),
  (176,70,'user.profile_maintenance','用户资料维护','lucide:contact','/user/profile-maintenance','page',998),
  (177,70,'user.password_reset','用户密码重置','lucide:key-round','/user/password-reset','page',999),
  (178,70,'user.promoter_assignment','批量设置推广员','lucide:badge-percent','/user/promoter-assignment','page',1000),
  (179,70,'user.notification','用户站内图文通知','lucide:send','/user/notification','page',1001),

  -- 内容（嵌套见 patch_content_menus.sql）
  (80,0,'content','内容','ant-design:read-outlined','/content','directory',8),
  (534,80,'content.article.dir','文章','lucide:newspaper','/cms','directory',1),
  (84,534,'content.article','文章管理','lucide:pen-line','/cms/article','page',1),
  (535,534,'content.article.category','文章分类','lucide:folder-tree','/cms/articleCategory','page',2),
  (82,80,'content.community','社区','lucide:images','/community','directory',2),
  (85,82,'content.community.category','社区分类','lucide:tags','/community/category','page',1),
  (86,82,'content.community.topic','社区话题','lucide:hash','/community/topic','page',2),
  (87,82,'content.community.list','社区内容','lucide:images','/community/list','page',3),
  (88,82,'content.community.reply','社区评论','lucide:message-square','/community/reply','page',4),
  (81,80,'content.notice','公告管理(旧)','lucide:bell','/content/notice','page',990),
  (83,80,'content.attachment','素材管理(旧)','lucide:images','/content/attachment','page',991),

  -- 财务（嵌套见 patch_accounts_menus.sql）
  (100,0,'accounts','财务','ant-design:bar-chart-outlined','/accounts','directory',9),
  (536,100,'accounts.merchant.dir','店铺结算','lucide:landmark','/mer/accounts','directory',1),
  (538,536,'accounts.statement','平台账单','lucide:file-spreadsheet','/accounts/statement','page',1),
  (21090,538,'accounts.statement.read','查看平台账单','','accounts/statement','button',1),
  (21091,538,'accounts.statement.download','下载平台账单','','accounts/statement','button',2),
  (539,536,'accounts.transfer','转账记录','lucide:arrow-left-right','/accounts/transferRecord','page',2),
  (21092,539,'accounts.transfer.read','查看转账记录','','accounts/transferRecord','button',1),
  (21093,539,'accounts.transfer.manage','审核备注并登记转账','','accounts/transferRecord','button',2),
  (21094,539,'accounts.transfer.export','导出转账记录','','accounts/transferRecord','button',3),
  (540,536,'accounts.profitsharing','分账管理','lucide:split','/merchant/applyList','page',3),
  (21095,540,'accounts.profitsharing.read','查看分账管理','','merchant/applyList','button',1),
  (21096,540,'accounts.profitsharing.manage','重新发起分账','','merchant/applyList','button',2),
  (21097,540,'accounts.profitsharing.export','导出分账列表','','merchant/applyList','button',3),
  (537,100,'accounts.user.dir','用户结算','lucide:wallet','/accounts/record','directory',2),
  (101,537,'accounts.withdraw','提现管理','lucide:wallet','/accounts/extract','page',1),
  (21098,101,'accounts.withdraw.read','查看提现管理','','accounts/extract','button',1),
  (20923,101,'accounts.withdraw.review','审核用户提现','','accounts/extract','button',2),
  (21099,101,'accounts.withdraw.export','导出提现列表','','accounts/extract','button',3),
  (541,537,'accounts.recharge_record','充值记录','lucide:circle-dollar-sign','/accounts/bill','page',2),
  (21100,541,'accounts.recharge_record.read','查看充值记录','','accounts/bill','button',1),
  (21101,541,'accounts.recharge_record.refund','充值退款','','accounts/bill','button',2),
  (102,537,'accounts.user_assets','资金记录','lucide:wallet-cards','/accounts/capital','page',3),
  (20924,102,'accounts.user_assets.read','查看资金记录','','accounts/capital','button',1),
  (21102,102,'accounts.user_assets.export','导出资金记录','','accounts/capital','button',2),
  (542,537,'accounts.capital_flow','资金流水','lucide:list','/accounts/capitalFlow','page',4),
  (21103,542,'accounts.capital_flow.read','查看资金流水','','accounts/capitalFlow','button',1),
  (21104,542,'accounts.capital_flow.export','导出资金流水','','accounts/capitalFlow','button',2),
  (188,100,'accounts.invoice','发票管理','lucide:receipt','/accounts/accounts','directory',3),
  (543,188,'accounts.invoice.list','发票列表','lucide:receipt','/accounts/receipt','page',1),
  (544,188,'accounts.invoice.desc','发票说明','lucide:file-text','/accounts/invoiceDesc','page',2),
  (21105,544,'accounts.invoice.desc.read','查看发票说明','','accounts/invoiceDesc','button',1),
  (21106,544,'accounts.invoice.desc.manage','维护发票说明','','accounts/invoiceDesc','button',2),
  (103,100,'accounts.merchant_settlement','店铺结算监管(旧)','lucide:landmark','/accounts/merchant-settlement','page',990),

  -- 应用（嵌套见 patch_app_menus.sql）
  (130,0,'app','应用','ant-design:appstore-outlined','/app','directory',10),
  (131,130,'app.wechat','公众号','lucide:message-circle','/app/wechat','directory',1),
  (204,131,'app.wechat_menus','微信菜单','lucide:menu','/app/wechat/menus','page',1),
  (203,131,'app.wechat_reply','自动回复','lucide:message-square-reply','/admin/app/wechat/reply','page',2),
  (206,131,'app.wechat_news','图文管理','lucide:newspaper','/app/wechat/newsCategory','page',3),
  -- 205 微信模板消息已下线，保留 id 兼容历史 role_menu，status 由 patch 置 0
  (205,131,'app.wechat_template','微信模板消息','lucide:mail','/app/wechat/template','page',99),
  (202,130,'app.routine','小程序','lucide:smartphone','/app/routine','page',2),
  (217,130,'app.mobile','App','lucide:smartphone','/app/mobile','page',3),

  -- 装修（嵌套见 patch_operations_menus.sql）
  (110,0,'operations','装修','ant-design:format-painter-outlined','/operations','directory',11),
  (111,110,'operations.diy','页面装修','lucide:paintbrush','/setting/diy/list','page',1),
  (800,110,'operations.merchant_diy','店铺模板','lucide:store','/setting/merchant/diyList','page',2),
  (801,110,'operations.product_detail','商品详情','lucide:package','/setting/diy/product_detail','page',3),
  (802,110,'operations.page_config','页面配置','lucide:sliders-horizontal','/setting/system_visualization_data','page',4),
  (803,110,'operations.page_links','页面链接','lucide:link','/setting/page','directory',5),
  (807,803,'operations.page_links.platform_cat','平台页面分类','lucide:folder-tree','/setting/diy/plantform/category/list','page',1),
  (808,803,'operations.page_links.platform','平台页面链接','lucide:link','/setting/diy/links/list','page',2),
  (809,803,'operations.page_links.merchant_cat','商户页面分类','lucide:folder-tree','/setting/diy/merchant/category/list','page',3),
  (810,803,'operations.page_links.merchant','商户页面链接','lucide:link','/setting/diy/merLink/list','page',4),
  (804,110,'operations.material','素材管理','lucide:images','/config/picture','page',6),
  (212,110,'operations.system_form','系统表单','lucide:clipboard-pen','/systemForm/form_list','page',7),
  (805,110,'operations.fab','悬浮菜单','lucide:circle-dot','/setting/fab','page',8),
  (806,110,'operations.product_category','商品分类','lucide:folder-tree','/setting/product_category','page',9),

  -- 客服
  (30,0,'service','客服','ant-design:customer-service-outlined','/service','directory',12),
  (301,30,'service.auto_reply','客服自动回复','lucide:message-square-reply','/systemForm/customer_keyword','page',1),
  (302,30,'service.customer.list','客服列表','ant-design:customer-service-outlined','/service/customer/list','page',2),
  (303,30,'service.settings','客服设置','lucide:settings-2','/systemForm/Basics/service','page',3),

  -- 设置（嵌套见 patch_setting_menus.sql）
  (120,0,'setting','设置','ant-design:setting-outlined','/setting','directory',13),
  (1500,120,'setting.system','系统设置','lucide:settings','/sys','directory',1),
  (1511,1500,'setting.system.service','服务配置','lucide:settings-2','/service/settings','page',1),
  (1512,1500,'setting.system.sms','短信配置','lucide:mail','/setting/sms','page',2),
  (1502,1500,'setting.shop.dir','商城设置','lucide:store','/shop','directory',3),
  (1520,1502,'setting.shop.form','商城设置','lucide:store','/systemForm/Basics/shop_tabs','page',1),
  (1521,1502,'setting.shop.hot','热门搜索','lucide:search','/group/config/67','page',2),
  (1522,1502,'setting.shop.agreements','协议设置','lucide:file-text','/setting/agreements','page',3),
  (1503,1500,'setting.delivery.dir','配送配置','lucide:truck','/delivery_config','directory',4),
  (1530,1503,'setting.delivery.express','物流公司','lucide:truck','/freight/express','page',1),
  (1504,1500,'setting.rbac.dir','权限管理','lucide:shield-check','/setting/rbac','directory',5),
  (1540,1504,'setting.rbac.role','角色权限','lucide:shield-check','/setting/role','page',1),
  (1541,1504,'setting.rbac.admin','管理员管理','lucide:user-round-cog','/setting/admin','page',2),
  (1542,1504,'setting.rbac.menu','菜单管理','lucide:menu','/setting/menu','page',3),
  (1505,1500,'setting.notice.dir','消息管理','lucide:bell','/notice','directory',6),
  (1550,1505,'setting.notice.station','公告管理','lucide:megaphone','/station/notice','page',1),
  (1551,1505,'setting.notice.list','消息管理','lucide:bell','/setting/notification/index','page',2),
  -- 旧扁平设置兼容（patch 软隐藏）
  (121,120,'setting.admin','管理员(旧)','lucide:user-round-cog','/setting/admin','page',990),
  (122,120,'setting.role','角色权限(旧)','lucide:shield-check','/setting/role','page',991),
  (123,120,'setting.menu','菜单管理(旧)','lucide:folder-tree','/setting/menu','page',992),
  (124,120,'setting.agreements','协议设置(旧)','lucide:receipt-text','/setting/agreements','page',993),
  (125,120,'setting.cloud_config','云服务配置(旧)','lucide:key-round','/setting/cloud-config','page',994),
  (126,120,'setting.sms','短信配置(旧)','lucide:messages-square','/setting/sms','page',995),
  (187,120,'setting.operation_log','操作日志(旧)','lucide:history','/setting/system-log','page',996),
  (189,120,'setting.login_log','登录日志(旧)','lucide:log-in','/setting/login-log','page',997),
  (190,120,'setting.shop','商城设置(旧)','lucide:store','/setting/shop','page',998),
  (191,120,'setting.pay','支付设置(旧)','lucide:credit-card','/setting/pay','page',999),
  (196,120,'setting.storage','存储配置(旧)','lucide:hard-drive','/setting/storage','page',1000),
  (90,120,'freight','物流配送(旧)','lucide:map-plus','/freight','directory',1001),
  (91,90,'freight.express','快递公司(旧)','lucide:map-plus','/freight/express','page',1),

  -- 维护（嵌套见 patch_maintain_menus.sql）
  (197,0,'maintain','维护','ant-design:tool-outlined','/maintain','directory',14),
  (1700,197,'maintain.exploit','开发配置','lucide:code','/safe/exploit','directory',1),
  (200,1700,'maintain.group_data','组合数据','lucide:database','/group/list','page',1),
  (1701,197,'maintain.safe','安全维护','lucide:shield','/maintain','directory',2),
  (199,1701,'maintain.backup','数据备份','lucide:hard-drive','/maintain/dataBackup','page',1),
  (1706,1701,'maintain.auth','商业授权','lucide:badge-check','/setting/system/maintain/auth','page',2),
  (198,1701,'maintain.cache','缓存清除','lucide:eraser','/maintain/cache','page',3),
  (1702,197,'maintain.config_classify','配置分类','lucide:folder-tree','/config/classify','page',3),
  (1703,197,'maintain.config_setting','配置管理','lucide:sliders-horizontal','/config/setting','page',4),
  (1704,197,'maintain.system_log','操作日志','lucide:history','/setting/systemLog','page',5),
  (1705,197,'maintain.export','导出记录','lucide:download','/group/exportList','page',6),
  (201,197,'maintain.hot_search','热门搜索(旧)','lucide:search','/group/config/67','page',990)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),`code`=VALUES(`code`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- 分销/营销/用户…：旧扁平入口侧栏软隐藏（正式树见上方；patch_*_menus.sql 可重复纠偏）
UPDATE `qixi_crm_a_menu` SET `status`=0 WHERE `id`=65 OR `code`='marketing.spread';
UPDATE `qixi_crm_a_menu` SET `status`=0 WHERE `id` IN (61,62,63,64,66,67,68,69,185,186,207,208,209,210,211)
  OR `code` IN (
    'marketing.coupon','marketing.seckill','marketing.combination','marketing.presell',
    'marketing.broadcast','marketing.assist','marketing.points','marketing.recharge',
    'marketing.coupon.send_records','marketing.coupon.receipt_records',
    'marketing.discounts','marketing.atmosphere','marketing.border','marketing.topic','marketing.application'
  );
UPDATE `qixi_crm_a_menu` SET `status`=0 WHERE `id` IN (78,79,170,171,172,173,174,175,176,177,178,179,81,83,103,121,122,123,124,125,126,187,189,190,191,196,90,91,201);

-- 首页树纠偏（可重复执行）
UPDATE `qixi_crm_a_menu` SET `code`='home',`title`='首页',`icon`='ant-design:home-outlined',`route_path`='/',`kind`='directory',`sort`=1,`parent_id`=0 WHERE `id`=1;
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  (2,1,'dashboard','控制台','lucide:layout-dashboard','/dashboard','page',1),
  (216,1,'data.screen','数据大屏','lucide:monitor','/data-screen/index','page',2),
  (194,1,'statistic.product','商品统计','lucide:package','/statistic/product','page',3),
  (193,1,'statistic.order','订单统计','lucide:receipt-text','/statistic/order','page',4),
  (195,1,'statistic.user','用户统计','lucide:users','/statistic/user','page',5)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`code`=VALUES(`code`),`title`=VALUES(`title`),`icon`=VALUES(`icon`),`route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`);
DELETE FROM `qixi_crm_a_role_menu` WHERE `menu_id`=192;
DELETE FROM `qixi_crm_a_menu` WHERE `id`=192;

-- 统一后台按钮节点必须使用 qixi_crm_a_* RBAC；运营写操作不再依赖旧系统菜单表。
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`) VALUES
  (34,52,'order.refund.approve','审核退款','','order/refund','button',1),
  (35,52,'order.refund.reject','拒绝退款','','order/refund','button',2),
  (20903,43,'product.audit.submit','审核商品','','product/audit','button',1),
  (20904,41,'product.category.manage','维护商品分类','','product/category','button',1),
  (20905,55,'product.brand.manage','维护商品品牌','','product/band/brandList','button',1),
  (20906,54,'product.brand.category.manage','维护品牌分类','','product/band/brandClassify','button',1),
  (21073,84,'content.article.read','查看文章','','cms/article','button',1),
  (20910,84,'content.article.manage','维护文章','','cms/article','button',2),
  (21074,535,'content.article_category.read','查看文章分类','','cms/articleCategory','button',1),
  (20911,535,'content.article_category.manage','维护文章分类','','cms/articleCategory','button',2),
  (21075,85,'content.community_category.read','查看社区分类','','community/category','button',1),
  (21076,85,'content.community_category.manage','维护社区分类','','community/category','button',2),
  (21082,87,'content.community_list.read','查看社区内容','','community/list','button',1),
  (21083,87,'content.community_list.manage','维护社区内容','','community/list','button',2),
  (21084,88,'content.community_reply.read','查看社区评论','','community/reply','button',1),
  (21085,88,'content.community_reply.manage','维护社区评论','','community/reply','button',2),
  (21080,86,'content.community_topic.read','查看社区话题','','community/topic','button',1),
  (21081,86,'content.community_topic.manage','维护社区话题','','community/topic','button',2),
  (20912,81,'content.notice.manage','维护公告','','content/notice','button',1),
	(21910,1551,'setting.notice.config.manage','维护消息配置','','setting/notification/index','button',1),
  (20913,124,'setting.agreement.manage','维护协议','','setting/agreements','button',1),
  (20914,111,'operations.diy.manage','维护商城装修','','operations/diy','button',1),
  (20915,794,'marketing.seckill.manage','维护秒杀活动','','marketing/seckill/list','button',1),
  (20916,1138,'marketing.combination.manage','维护拼团活动','','marketing/combination/combination_list','button',1),
  (20917,1023,'marketing.presell.manage','维护预售活动','','marketing/presell/list','button',1),
  (20918,1658,'marketing.coupon.manage','维护平台优惠券','','marketing/platform_coupon/list','button',1),
  (21010,1629,'marketing.discounts.read','查看优惠套餐','','marketing/discounts/list','button',1),
  (21011,1629,'marketing.discounts.manage','上下架优惠套餐','','marketing/discounts/list','button',2),
  (21012,9007,'marketing.atmosphere.read','查看活动氛围','','marketing/atmosphere/list','button',1),
  (21013,9007,'marketing.atmosphere.manage','维护活动氛围','','marketing/atmosphere/list','button',2),
  (21014,9008,'marketing.border.read','查看活动边框','','marketing/border/list','button',1),
  (21015,9008,'marketing.border.manage','维护活动边框','','marketing/border/list','button',2),
  (21016,1470,'marketing.topic.read','查看专场','','group/topic/94','button',1),
  (21017,1470,'marketing.topic.manage','维护专场','','group/topic/94','button',2),
  (21018,9217,'marketing.application.read','查看活动报名','','marketing/application/list','button',1),
  (21019,9217,'marketing.application.manage','维护活动报名','','marketing/application/list','button',2),
  (21020,212,'operations.system_form.manage','维护系统表单','','systemForm/form_list','button',1),
  (20919,83,'content.attachment.manage','维护素材库','','content/attachment','button',1),
  (20920,87,'content.community.audit','审核社区内容','','community/list','button',3),
  (20921,87,'content.community.delete','删除社区内容','','community/list','button',4),
  (20922,781,'marketing.broadcast.audit','审核直播间','','marketing/studio/list','button',1),
  (20925,103,'accounts.merchant_settlement.read','查看店铺结算监管投影','','accounts/merchant-settlement','button',1),
  (20926,103,'accounts.merchant_settlement.review','审核并登记店铺结算凭证','','accounts/merchant-settlement','button',2),
  (20927,522,'marketing.spread.read','查看分销推广与佣金监管','','promoter/user','button',1),
  (20928,1096,'marketing.assist.manage','维护好友助力活动','','marketing/assist/list','button',1),
  (20929,21,'region.zone.manage','维护区域列表','','business-zones/index','button',1),
  (20930,22,'region.agent.manage','维护代理人员','','business-zones/agents','button',1),
  (20931,22,'region.agent.review','审核区域代理申请','','business-zones/agent-review','button',1),
  (20932,91,'freight.express.manage','维护快递公司与行政区划','','freight/express','button',1),
  (20933,11,'merchant.status.manage','启停店铺经营状态','','merchant/list','button',1),
  (20934,12,'merchant.intention.audit','审核店铺入驻申请','','merchant/audit','button',1),
  (20935,12,'merchant.intention.assign_region','分配店铺入驻申请区域','','merchant/audit','button',2),
  (21003,12,'merchant.intention.create','新增店铺入驻申请','','merchant/audit','button',0),
  (20936,1512,'setting.sms.manage','维护平台短信验证码配置','','setting/sms','button',1),
  (20937,13,'merchant.category.manage','维护商户分类与佣金比例','','merchant/categories','button',1),
  (20938,14,'merchant.group.manage','维护店铺分组、关联店铺与装修模板','','merchant/grouping','button',1)
  ,(20939,15,'merchant.type.manage','维护店铺类型、保证金规则与授权菜单','','merchant/types','button',1)
  ,(20940,16,'merchant.deposit.review','维护保证金、扣减、退款审核及打款登记','','merchant/deposits','button',1)
  ,(21001,213,'store.margin_config.manage','维护保证金补缴提醒配置','','systemForm/Basics/margin','button',1)
  ,(21002,29,'merchant.mgmt.settings.manage','维护商户入驻页面与表单设置','','merchant/apply-setting','button',1)
  ,(21004,24,'region.agent_settings.manage','维护区域代理默认提成与申请表单','','business-zones/settings','button',1)
  ,(20941,17,'merchant.profitsharing.review','审核店铺分账申请及维护审核备注','','merchant/applyments','button',1)
  ,(20942,44,'product.label.manage','维护平台商品标签','','product/label','button',1)
  ,(20943,45,'product.guarantee.manage','维护平台保障服务','','product/guarantee','button',1)
  ,(20944,57,'product.parameter.manage','维护平台商品参数模板','','product/specs','button',1)
  ,(20907,56,'product.parameter.store.manage','维护店铺商品参数模板','','product/merSpecs','button',1)
  ,(21021,48,'product.price_description.manage','维护平台价格说明','','product/priceDescription','button',1)
  ,(20945,47,'product.comment.review','审核或隐藏商品评论','','product/comment','button',1)
  ,(20946,47,'product.comment.virtual.manage','新增或编辑虚拟商品评论','','product/comment','button',2)
  ,(20947,47,'product.comment.sort','调整虚拟商品评论排序','','product/comment','button',3)
  ,(20948,47,'product.comment.delete','删除虚拟商品评论','','product/comment','button',4)
  ,(20949,75,'user.feedback.read','查看用户反馈','','user/feedback/list','button',1)
  ,(20950,75,'user.feedback.manage','回复、关闭或删除用户反馈','','user/feedback/list','button',2)
  ,(21064,76,'user.feedback.category.read','查看反馈分类','','user/feedback/categories','button',1)
  ,(20951,76,'user.feedback.category.manage','维护反馈分类','','user/feedback/categories','button',2)
  ,(20952,77,'user.list.read','查看脱敏用户列表','','user/list','button',1)
  ,(20953,77,'user.list.manage','人工调整用户余额、积分与会员等级','','user/list','button',2)
  ,(20954,78,'user.asset.adjust','提交用户余额或积分人工调整','','user/assets-adjustment','button',1)
  ,(20955,79,'user.member.adjust','提交用户会员等级人工调整','','user/member-adjustment','button',1)
  ,(20956,170,'user.coupon.manage','提交用户优惠券发放或撤销','','user/coupon-operation','button',1)
  ,(20957,171,'user.referrer.manage','提交用户推荐关系调整','','user/referrer-adjustment','button',1)
  ,(20958,172,'user.group.assign','提交用户单个或批量分组调整','','user/group-assignment','button',1)
  ,(21060,72,'user.group.read','查看用户分组','','user/group','button',1)
  ,(21061,72,'user.group.manage','维护用户分组','','user/group','button',2)
  ,(21062,71,'user.label.read','查看用户标签','','user/label','button',1)
  ,(21063,71,'user.label.manage','维护用户标签','','user/label','button',2)
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
  ,(20970,9119,'marketing.points.manage','维护积分商品并查看积分订单','','marketing/integral/proList','button',1)
  ,(20971,687,'marketing.recharge.manage','维护充值计划并查看充值订单','','group/config/69','button',1)
  ,(20972,73,'user.svip.manage','查看并维护用户 SVIP 状态','','user/svip','button',1)
  ,(20973,180,'user.svip.plan.manage','维护可售会员类型','','user/member/type','button',1)
  ,(20974,181,'user.svip.record.read','查看会员购买记录与统计','','user/member/record','button',1)
  ,(21070,533,'user.svip.agreement.read','查看会员协议','','user/member/vipAgreement','button',1)
  ,(21071,533,'user.svip.agreement.manage','维护会员协议','','user/member/vipAgreement','button',2)
  ,(20975,182,'user.search_record.read','查看用户搜索记录','','user/search-record','button',1)
  ,(20976,182,'user.search_record.clear','按用户清理搜索记录','','user/search-record','button',2)
  ,(20977,182,'user.search_record.export','导出用户搜索记录','','user/search-record','button',3)
  ,(20978,183,'user.member.level.manage','维护普通会员等级、规则和权益','','user/member/levels','button',1)
  ,(20979,184,'user.svip.interest.manage','维护付费会员权益','','user/member/equity','button',1)
  ,(20980,1662,'marketing.coupon.send.read','查看平台人工发券与撤销审计','','marketing/platform_coupon/couponSend','button',1)
  ,(20981,1659,'marketing.coupon.record.read','查看用户优惠券领取与使用状态','','marketing/platform_coupon/couponRecord','button',1)
  ,(20982,187,'setting.operation_log.read','查看统一后台成功操作日志','','setting/system-log','button',1)
  ,(20983,188,'accounts.invoice.read','查看订单发票监管记录','','accounts/invoices','button',1)
  ,(20984,190,'setting.shop.manage','维护商城基础设置','','setting/shop','button',1)
  ,(20985,191,'setting.pay.manage','维护支付方式开关','','setting/pay','button',1)
  ,(20986,53,'order.cancellation.read','查看核销记录','','order/cancellation','button',1)
  ,(20987,131,'app.wechat.manage','维护公众号基础开关','','app/wechat','button',1)
  ,(20988,196,'setting.storage.manage','维护存储配置开关','','setting/storage','button',1)
  ,(20989,198,'maintain.cache.manage','提交缓存清理','','maintain/cache','button',1)
  ,(20990,199,'maintain.backup.manage','维护备份记录','','maintain/dataBackup','button',1)
  ,(20991,200,'maintain.group_data.manage','维护组合数据','','group/list','button',1)
  ,(20992,201,'maintain.hot_search.manage','维护热门搜索','','group/config/67','button',1)
  ,(21570,1521,'setting.shop.hot.read','查看热门搜索','','group/config/67','button',1)
  ,(21571,1521,'setting.shop.hot.manage','维护热门搜索','','group/config/67','button',2)
  ,(20993,202,'app.routine.manage','维护小程序基础开关','','app/routine','button',1)
  ,(21122,217,'app.mobile.manage','维护 iOS、Android、HarmonyOS 应用配置','','app/mobile','button',1)
  ,(21123,217,'app.push.manage','维护 App 推送配置','','app/mobile','button',2)
  ,(21560,1511,'systemServeMerMealLst','服务套餐列表','','service/settings','button',1)
  ,(21561,1511,'systemServeMealDetail','服务套餐详情','','service/settings','button',2)
  ,(21562,1511,'systemServeMealCreate','新增服务套餐','','service/settings','button',3)
  ,(21563,1511,'systemServeMealUpdate','编辑服务套餐','','service/settings','button',4)
  ,(21564,1511,'systemServeMealDelete','删除服务套餐','','service/settings','button',5)
  ,(21565,1511,'systemServeMealStatus','修改服务套餐状态','','service/settings','button',6)
  ,(20994,203,'app.wechat_reply.manage','维护自动回复','','admin/app/wechat/reply','button',2)
  ,(21121,203,'app.wechat_reply.read','查看自动回复','','admin/app/wechat/reply','button',1)
  ,(20995,204,'app.wechat_menus.manage','维护微信菜单开关','','app/wechat/menus','button',1)
  ,(20996,205,'app.wechat_template.manage','维护模板消息开关','','app/wechat/template','button',1)
  ,(20997,206,'app.wechat_news.manage','维护图文管理','','app/wechat/newsCategory','button',2)
  ,(21120,206,'app.wechat_news.read','查看图文管理','','app/wechat/newsCategory','button',1)
  ,(21072,532,'user.setup.read','查看用户设置','','user/setup_user','button',1)
  ,(20998,532,'user.setup.manage','保存用户设置','','user/setup_user','button',2)
  ,(20999,101,'accounts.transfer_settings.manage','维护转账监管设置','','accounts/settings','button',2)
  ,(21000,12,'merchant.intention.delete','删除店铺入驻申请','','merchant/audit','button',3)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`title`=VALUES(`title`),`route_path`=VALUES(`route_path`),`kind`=VALUES(`kind`),`sort`=VALUES(`sort`),`status`=1;

-- 公众号「微信模板消息」下线（软隐藏）
UPDATE `qixi_crm_a_menu` SET `status`=0, `sort`=99 WHERE `id` IN (205,20996) OR `code` IN ('app.wechat_template','app.wechat_template.manage');

-- 历史“快递公司”是物流承运商维护，不得成为一级菜单；正式入口为
-- 设置 → 系统设置 → 配送配置 → 物流公司（setting.delivery.express）。
UPDATE `qixi_crm_a_menu`
SET `status`=0, `sort`=1001
WHERE `id` IN (90,91) OR `code` IN ('freight','freight.express');

-- 商户管理「商户入驻审核」统一为「店铺入驻申请」（与店铺管理入口同功能，路由保持 /merchant/review）
UPDATE `qixi_crm_a_menu`
SET `title`='店铺入驻申请', `icon`='lucide:badge-check', `route_path`='/merchant/review', `kind`='page', `sort`=2, `parent_id`=25, `status`=1
WHERE `id`=27 OR `code`='merchant.mgmt.review';

-- 商品：活动标签页保留路由，侧栏按图片1软隐藏（status=0）
UPDATE `qixi_crm_a_menu`
SET `title`='活动标签', `icon`='lucide:badge', `route_path`='/product/activityLabel', `kind`='page', `sort`=99, `parent_id`=40, `status`=0
WHERE `id`=49 OR `code`='product.activity_label';

-- 区域代理：侧栏只保留区域列表 / 代理人员 / 代理设置；代理审核菜单软隐藏（status=0），API 按钮权限仍挂在代理人员下
UPDATE `qixi_crm_a_menu`
SET `title`='区域列表', `icon`='lucide:map-pinned', `route_path`='/business-zones/index', `kind`='page', `sort`=1, `parent_id`=20, `status`=1
WHERE `id`=21 OR `code`='region.index';
UPDATE `qixi_crm_a_menu`
SET `title`='代理人员', `icon`='lucide:users', `route_path`='/business-zones/agents', `kind`='page', `sort`=2, `parent_id`=20, `status`=1
WHERE `id`=22 OR `code`='region.agents';
UPDATE `qixi_crm_a_menu`
SET `title`='代理审核', `icon`='lucide:badge-check', `route_path`='/business-zones/agent-review', `kind`='page', `sort`=99, `parent_id`=20, `status`=0
WHERE `id`=23 OR `code`='region.agent_review';
UPDATE `qixi_crm_a_menu`
SET `title`='代理设置', `icon`='lucide:settings-2', `route_path`='/business-zones/settings', `kind`='page', `sort`=3, `parent_id`=20, `status`=1
WHERE `id`=24 OR `code`='region.agent_settings';
UPDATE `qixi_crm_a_menu`
SET `parent_id`=22, `title`='审核区域代理申请', `route_path`='business-zones/agent-review', `kind`='button', `sort`=1, `status`=1
WHERE `id`=20931 OR `code`='region.agent.review';
UPDATE `qixi_crm_a_menu`
SET `title`='维护区域列表', `route_path`='business-zones/index'
WHERE `id`=20929 OR `code`='region.zone.manage';
UPDATE `qixi_crm_a_menu`
SET `title`='维护代理人员'
WHERE `id`=20930 OR `code`='region.agent.manage';

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
  'marketing','marketing.coupon','marketing.seckill','marketing.combination','marketing.presell','marketing.spread','marketing.broadcast','marketing.assist','marketing.points','marketing.recharge','marketing.discounts','marketing.atmosphere','marketing.border','marketing.topic','marketing.application',
  'content','content.notice','content.community','content.community.category','content.community.topic','content.community.list','content.community.reply','content.attachment','content.article'
)) OR (r.code = 'region' AND m.code IN ('region','region.index','region.agents','region.agent_review'))
  OR (r.code = 'operations' AND m.code IN ('dashboard','user','user.label','user.group','user.svip.plan','user.svip.record','user.svip.interest'));

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
JOIN `qixi_crm_a_menu` AS m
WHERE r.code = 'platform'
   OR (r.code = 'merchant' AND m.code IN (
      'home','dashboard','data.screen','statistic.product','statistic.order','statistic.user',
      'store','store.manage','merchant.list','merchant.audit',
      'product','product.audit',
      'order','order.list','order.refund','order.cancellation'))
   OR (r.code = 'region' AND m.code IN (
      'home','dashboard','data.screen','statistic.product','statistic.order','statistic.user',
      'store','store.manage','merchant.list','merchant.audit',
      'product','product.audit','order','order.list','order.refund','order.cancellation',
      'accounts','accounts.merchant_settlement','accounts.merchant_settlement.read','merchant.intention.audit','merchant.intention.delete'))
   OR (r.code = 'customer_service' AND m.code IN (
      'home','dashboard','service','service.auto_reply','service.customer.list','service.settings',
      'user','user.feedback','user.feedback.list','user.feedback.read','user.feedback.manage'))
   OR (r.code = 'operations' AND m.code IN (
      'promoter','promoter.user','promoter.brokerage','promoter.brokerage.level','promoter.brokerage.rule',
      'promoter.bank','promoter.privilege','promoter.poster','promoter.gift','promoter.commission',
      'promoter.order','promoter.explain','promoter.config','promoter.config.manage','marketing.spread.read',
      'marketing','marketing.platform_coupon','marketing.platform_coupon.list','marketing.platform_coupon.record',
      'marketing.platform_coupon.send','marketing.platform_coupon.help',
      'marketing.store_coupon','marketing.store_coupon.list','marketing.store_coupon.user',
      'marketing.seckill.dir','marketing.seckill.config','marketing.seckill.manage.page','marketing.seckill.activity',
      'marketing.broadcast.dir','marketing.broadcast.studio','marketing.broadcast.goods',
      'marketing.presell.dir','marketing.presell.goods','marketing.presell.agreement',
      'marketing.assist.dir','marketing.assist.goods','marketing.assist.activity',
      'marketing.combination.dir','marketing.combination.set','marketing.combination.goods','marketing.combination.list',
      'marketing.integral.dir','marketing.integral.config','marketing.integral.log','marketing.integral.log.read','marketing.integral.classify','marketing.integral.classify.manage',
      'marketing.integral.products','marketing.integral.orders','marketing.integral.orders.read','marketing.integral.orders.manage',
      'marketing.atmosphere.nav','marketing.atmosphere.read','marketing.atmosphere.manage',
      'marketing.border.nav','marketing.border.read','marketing.border.manage',
      'marketing.topic.nav','marketing.topic.read','marketing.topic.manage',
      'marketing.discounts.nav','marketing.discounts.read','marketing.discounts.manage',
      'marketing.balance.dir','marketing.balance.settings','marketing.balance.settings.read','marketing.balance.settings.manage','marketing.balance.config','marketing.balance.config.read','marketing.balance.config.manage',
      'marketing.application.nav','marketing.application.read','marketing.application.manage',
      'user','user.feedback','user.feedback.category','user.feedback.category.read','user.feedback.category.manage','user.group','user.group.read','user.group.manage','user.label','user.label.read','user.label.manage','user.svip','user.svip.plan','user.svip.record','user.svip.interest','user.svip.agreement','user.svip.plan.manage','user.svip.record.read','user.svip.interest.manage','user.svip.agreement.read','user.svip.agreement.manage',
      'content','content.notice','content.community','content.community.category','content.community.topic','content.community.list','content.community.reply','content.attachment','content.article','content.article.dir','content.article.category','content.article.read',
      'operations','operations.diy','operations.system_form','operations.system_form.manage','setting','setting.agreements',
      'content.article.manage','content.article_category.read','content.article_category.manage','content.community_category.read','content.community_category.manage','content.community_topic.read','content.community_topic.manage','content.community_list.read','content.community_list.manage','content.community_reply.read','content.community_reply.manage','content.notice.manage','setting.agreement.manage','setting.notice.dir','setting.notice.list','setting.notice.config.manage',
      'operations.diy.manage','marketing.seckill.manage','marketing.combination.manage','marketing.presell.manage','marketing.coupon.manage','marketing.assist.manage','marketing.points.manage','marketing.recharge.manage',
      'maintain','maintain.cache','maintain.cache.manage','maintain.backup','maintain.backup.manage','maintain.group_data','maintain.group_data.manage','maintain.hot_search','maintain.hot_search.manage',
      'setting.shop.hot','setting.shop.hot.read','setting.shop.hot.manage',
      'content.attachment.manage','content.community.audit','content.community.delete','marketing.broadcast.audit'));

-- 系统设置不包含一号通/呼叫系统登录入口；服务配置与短信配置直接挂在系统设置。
DELETE FROM `qixi_crm_a_role_menu` WHERE `menu_id` IN (1501,1510);
DELETE FROM `qixi_crm_a_menu` WHERE `id` IN (1501,1510) OR `code` IN ('setting.serve','setting.serve.login');
UPDATE `qixi_crm_a_menu`
SET `parent_id`=1512, `route_path`='setting/sms', `title`='维护平台短信验证码配置', `status`=1
WHERE `id`=20936 OR `code`='setting.sms.manage';
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform' AND m.id IN (1500,1511,1512,20936);

-- 公开参数允许初始化；App Key 由后台加密保存，严禁写入 SQL。
INSERT INTO `qixi_crm_a_cloud_config`
  (`provider`,`config_key`,`ciphertext`,`key_version`,`updated_by`)
VALUES
  ('tencent_sms','enabled','true','bootstrap-public-v1',0),
  ('tencent_sms','sdk_app_id','1401165606','bootstrap-public-v1',0),
  ('tencent_sms','sign_id','711884','bootstrap-public-v1',0),
  ('tencent_sms','sign_content','杭州乐成体育','bootstrap-public-v1',0),
  ('tencent_sms','template_id','2701987','bootstrap-public-v1',0)
ON DUPLICATE KEY UPDATE
  `ciphertext`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`ciphertext`),`ciphertext`),
  `key_version`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`key_version`),`key_version`),
  `updated_by`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`updated_by`),`updated_by`);

-- 平台素材库系统预设分类（不可增删改；侧栏「全部素材」为虚拟入口不落库）
-- 用途：客户端（H5/小程序/App）与装修页常用的图标、图片、背景、视频
INSERT INTO `qixi_crm_a_attachment_category`
  (`attachment_category_id`,`pid`,`path`,`attachment_category_name`,`attachment_category_enname`,`sort`,`mer_id`,`is_system`,`create_time`)
VALUES
  (5101,0,'','店铺封面','store_cover',90,0,1,NOW()),
  (5102,0,'','支付图标','pay_icon',80,0,1,NOW()),
  (5103,0,'','物流图标','logistics_icon',70,0,1,NOW()),
  (5104,0,'','客服图标','service_icon',60,0,1,NOW()),
  (5105,0,'','商品图片','product_image',50,0,1,NOW()),
  (5106,0,'','背景图片','background_image',40,0,1,NOW()),
  (5107,0,'','列表图标','list_icon',30,0,1,NOW()),
  (5108,0,'','其他图片','other_image',20,0,1,NOW()),
  (5111,0,'','店铺视频','store_video',19,0,1,NOW()),
  (5112,0,'','商品视频','product_video',18,0,1,NOW()),
  (5113,0,'','其他视频','other_video',17,0,1,NOW())
ON DUPLICATE KEY UPDATE
  `attachment_category_name`=VALUES(`attachment_category_name`),
  `attachment_category_enname`=VALUES(`attachment_category_enname`),
  `sort`=VALUES(`sort`),
  `mer_id`=VALUES(`mer_id`),
  `is_system`=VALUES(`is_system`);

-- 以下均为平台运行所需的默认配置，不属于商品、店铺或用户演示数据。
-- 使用 INSERT IGNORE / 保守更新：仅补齐缺失项，不覆盖运营在后台保存的配置。

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('mall_shop_config',0,'{"site_name":"七禧商城","site_url":"","enabled":true,"remark":"","auto_parse_clipboard":true,"arrival_notice_enabled":true,"product_comment_enabled":true,"auto_positive_review_enabled":true,"default_copy_times":8,"order_auto_cancel_minutes":15,"order_auto_receive_days":7,"after_sale_days":1,"merchant_refund_auto_days":1,"refund_reasons":["商品质量问题","不想要了","未收到货"],"platform_rights_enabled":true,"platform_rights_days":1,"merge_payment_enabled":true,"merchant_apply_enabled":true,"merchant_qualification_required":true,"merchant_margin_badge_enabled":false,"merchant_margin_badge_image":"","merchant_category_limit":5,"mall_show_stores":true,"mall_recommend_enabled":true,"mall_recommend_distance_enabled":true,"mall_recommend_sort":"star","live_stream_auto_approve":false,"live_product_auto_approve":false,"hot_ranking_enabled":true,"hot_ranking_category_level":2,"hot_ranking_refresh_hours":24,"mall_search_mode":"fuzzy","product_ranking_period":"month","product_ranking_metric":"sales_amount","shop_ranking_period":"month","shop_ranking_metric":"product_count","dashboard_display_name":"数据大屏"}'),
  ('sys_merchant_type',0,'<p>请根据经营模式选择适用的店铺类型，并遵守平台入驻规则。</p>'),
  ('sys_merchant_category',0,'<p>请根据主营商品或服务选择店铺分类，便于平台提供运营支持。</p>'),
  ('sys_brokerage',0,'<ol><li><p>一级分销佣金按平台规则结算。</p></li><li><p>二级分销佣金按平台规则结算。</p></li></ol>'),
  ('sys_extension_agree',0,'<p>分销佣金以平台已发布的规则和订单结算结果为准。</p>'),
  ('promoter_explain',0,'<p>推广服务说明以平台公示规则为准。</p>'),
  ('sys_coupon_agree',0,'<p>优惠券请在有效期内使用，具体适用范围以优惠券说明为准。</p>'),
  ('sys_product_presell_agree',0,'<p>预售商品以页面公示的发货时间和尾款支付规则为准。</p>'),
  ('sms_config',0,'{"enabled":true,"provider":"tencent_sms"}'),
  ('merchant_apply_setting',0,'{"background_image":"","form_fields":[]}'),
  ('distribution_config',0,'{"extension_status":true,"extension_self":true,"extension_limit":false,"extension_limit_day":15,"promoter_type":0,"promoter_low_money":0,"extension_pop":0,"extension_one_rate":0.15,"extension_two_rate":0.05,"user_extract_min":10,"lock_brokerage_timer":7,"sys_extension_type":0,"withdraw_type":["0","1","2"],"extract_switch":1,"transfer_scene_id":0,"max_bag_number":10}'),
  ('group_buying_config',0,'{"ficti_status":1,"group_buying_rate":30}'),
  ('integral_config',0,'{"integral_status":1,"integral_money":0.1,"integral_order_rate":1,"integral_freeze":0,"integral_clear_time":24,"integral_user_give":50,"integral_community_give":10,"integral_community_give_limit":10,"rule":""}')
ON DUPLICATE KEY UPDATE
  `result`=IF(`key`='mall_shop_config' AND (JSON_VALID(`result`)=0 OR `result`=''),VALUES(`result`),`result`),
  `expire_time`=IF(`result`='' OR `result` LIKE '%本地验收%',VALUES(`expire_time`),`expire_time`);

INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_user_agree',0,'<h2>用户协议</h2><p>使用平台服务前，请仔细阅读并理解本协议内容。</p><h3>一、账户与使用</h3><p>请妥善保管账户信息，并在遵守平台规则的前提下使用服务。</p><h3>二、服务说明</h3><p>平台将持续优化服务内容，并依法保障用户的合法权益。</p>'),
  ('sys_userr_privacy',0,'<h2>隐私协议</h2><p>平台将按照适用法律法规处理提供服务所必需的个人信息。</p><h3>一、信息使用</h3><p>仅在提供服务、保障安全和履行法定义务所必需的范围内使用信息。</p>'),
  ('the_cancellation_prompt',0,'<h2>注销提示</h2><p>提交注销前，请确认账户内不存在未完成订单、售后或其他待处理事项。</p><p>注销后，部分服务和权益可能无法恢复。</p>'),
  ('platform_rule',0,'<h2>平台规则</h2><p>请遵守法律法规及平台公示规则，共同维护安全、有序的交易环境。</p>'),
  ('sys_intention_agree',0,'<h2>店铺入驻申请协议</h2><p>申请店铺入驻前，请确认提交的信息真实、准确、完整，并同意遵守平台入驻规则。</p>'),
  ('circle_entry_agree',0,'<h2>代理入驻申请协议</h2><p>申请代理入驻前，请确认具备相应资质，并同意遵守平台代理合作规则。</p>'),
  ('business_entry_agree',0,'<h2>商户入驻申请协议</h2><p>申请商户入驻前，请确认提交资料真实有效，并同意遵守平台商户经营规范。</p>'),
  ('the_cancellation_msg',0,'<h2>注销声明</h2><p>账户注销申请完成后，平台将按照适用规则处理账户资料和相关服务记录。</p>'),
  ('sys_about_us',0,'<h2>关于我们</h2><p>请在此维护平台介绍、服务理念及实际主体联系信息。</p>'),
  ('sys_certificate',0,'<h2>资质证照</h2><p>请在此维护平台依法应公示的主体资质与相关证照信息。</p>')
ON DUPLICATE KEY UPDATE
  `result`=IF(`result`='' OR `result` LIKE '%本地验收%' OR `result` LIKE '%虚构中文%',VALUES(`result`),`result`),
  `expire_time`=VALUES(`expire_time`);

INSERT IGNORE INTO `qixi_crm_a_config` (`config_key`,`config_value`,`updated_by`) VALUES
  ('customer_service.settings',JSON_OBJECT('auto_reply_enabled',TRUE,'auto_reply_text','您好，客服将在工作时间内回复您。','queue_mode','round_robin','max_sessions_per_agent',12),NULL);

-- 通知会员、通知店铺的默认发送行为与固定文本。渠道开关可在后台逐项调整。
INSERT IGNORE INTO `qixi_crm_a_notification_config`
  (`notification_id`,`audience`,`notice_type`,`scene`,`wechat_enabled`,`mini_program_enabled`,`sms_enabled`,`wechat_text`,`mini_program_text`,`sms_text`)
VALUES
  (1,'member','订单发货通知','发货给用户的提醒',1,1,0,'您的订单已发货，请留意物流信息。','您的订单已发货，请留意物流信息。','您的订单已发货，请留意物流信息。'),
  (3,'member','订单已签收通知','确认收货给用户的提醒',1,1,0,'您的订单已签收，感谢您的支持。','您的订单已签收，感谢您的支持。','您的订单已签收，感谢您的支持。'),
  (4,'member','订单支付成功通知','用户支付成功提醒',1,1,0,'您的订单支付成功，我们将尽快为您处理。','您的订单支付成功，我们将尽快为您处理。','您的订单支付成功，我们将尽快为您处理。'),
  (5,'member','改价提醒','订单价格调整提醒',0,0,0,'您的订单价格已调整，请查看订单详情。','您的订单价格已调整，请查看订单详情。','您的订单价格已调整，请查看订单详情。'),
  (6,'member','提醒付款通知','订单关闭前提醒付款通知',0,0,0,'您的订单待付款，请及时完成支付。','您的订单待付款，请及时完成支付。','您的订单待付款，请及时完成支付。'),
  (7,'member','退货退款申请结果通知','商家拒绝退款给用户提醒',1,0,0,'您的退货退款申请暂未通过，请查看订单详情。','您的退货退款申请暂未通过，请查看订单详情。','您的退货退款申请暂未通过，请查看订单详情。'),
  (8,'member','退货退款申请结果通知','商家同意退款给用户提醒',1,0,0,'您的退货退款申请已通过，请按订单指引操作。','您的退货退款申请已通过，请按订单指引操作。','您的退货退款申请已通过，请按订单指引操作。'),
  (9,'member','退款成功通知','退款成功给用户通知',1,1,0,'您的退款已成功处理，请留意到账情况。','您的退款已成功处理，请留意到账情况。','您的退款已成功处理，请留意到账情况。'),
  (14,'member','预售尾款支付通知','预售尾款支付通知',0,0,0,'您的预售订单待支付尾款，请及时完成支付。','您的预售订单待支付尾款，请及时完成支付。','您的预售订单待支付尾款，请及时完成支付。'),
  (17,'member','入驻申请未通过提醒','入驻申请未通过提醒',0,0,0,'您的入驻申请暂未通过，请根据提示完善资料。','您的入驻申请暂未通过，请根据提示完善资料。','您的入驻申请暂未通过，请根据提示完善资料。'),
  (19,'member','到货提醒通知','到货提醒用户通知',0,1,0,'您关注的商品已到货。','您关注的商品已到货。','您关注的商品已到货。'),
  (20,'member','积分即将到期提醒','积分即将到期提醒',0,0,0,'您的部分积分即将到期，请及时使用。','您的部分积分即将到期，请及时使用。','您的部分积分即将到期，请及时使用。'),
  (21,'member','账户资金变动提醒','账户资金变动提醒',0,1,0,'您的账户资金发生变动，请查看详情。','您的账户资金发生变动，请查看详情。','您的账户资金发生变动，请查看详情。'),
  (23,'member','拼团成功','拼团成功',0,0,0,'您的拼团已成功，请查看订单详情。','您的拼团已成功，请查看订单详情。','您的拼团已成功，请查看订单详情。'),
  (25,'member','用户提现结果通知','用户提现审核通知',0,1,0,'您的提现申请已有处理结果，请查看详情。','您的提现申请已有处理结果，请查看详情。','您的提现申请已有处理结果，请查看详情。'),
  (26,'member','保证金退回申请通过通知','保证金退回申请通过给商户通知',0,0,0,'您的保证金退回申请已通过。','您的保证金退回申请已通过。','您的保证金退回申请已通过。'),
  (27,'member','保证金退回申请未通过通知','保证金退回申请未通过给商户的通知',0,0,0,'您的保证金退回申请暂未通过，请查看详情。','您的保证金退回申请暂未通过，请查看详情。','您的保证金退回申请暂未通过，请查看详情。'),
  (28,'member','短信验证码','短信验证码',0,0,1,'','', '您的验证码为：{code}，请勿泄露给他人。'),
  (29,'member','商户申请分账通过','开启分账商户申请分账通过短信提醒',0,0,0,'商户分账申请已通过。','商户分账申请已通过。','商户分账申请已通过。'),
  (30,'member','商户申请分账未通过','开启分账商户申请分账未通过短信提醒',0,0,0,'商户分账申请暂未通过。','商户分账申请暂未通过。','商户分账申请暂未通过。'),
  (31,'member','商户申请分账待验证','开启分账商户申请分账需操作验证的短信提醒',0,0,0,'商户分账申请需要完成验证。','商户分账申请需要完成验证。','商户分账申请需要完成验证。'),
  (32,'member','付费会员开通','付费会员付费成功后给用户通知',0,0,0,'您的付费会员已开通。','您的付费会员已开通。','您的付费会员已开通。'),
  (33,'member','订单配送通知','订单配送通知',1,1,0,'您的订单正在配送，请留意收货。','您的订单正在配送，请留意收货。','您的订单正在配送，请留意收货。'),
  (34,'member','助力成功通知','助力成功给用户的通知',0,0,0,'您的助力已成功。','您的助力已成功。','您的助力已成功。'),
  (2,'store','用户下单成功通知','创建订单给商户通知',1,0,0,'有新的用户订单，请及时处理。','有新的用户订单，请及时处理。','有新的用户订单，请及时处理。'),
  (10,'store','订单支付成功通知','订单支付成功通知给管理员',1,0,0,'订单已支付成功，请及时安排发货。','订单已支付成功，请及时安排发货。','订单已支付成功，请及时安排发货。'),
  (11,'store','退款通知','申请退款客服通知',1,0,0,'有新的退款申请，请及时处理。','有新的退款申请，请及时处理。','有新的退款申请，请及时处理。'),
  (12,'store','订单已签收通知','确认收货给管理员通知',1,0,0,'订单已确认收货。','订单已确认收货。','订单已确认收货。'),
  (13,'store','退货信息提醒','退货信息提醒',0,0,0,'用户已提交退货信息，请及时查看。','用户已提交退货信息，请及时查看。','用户已提交退货信息，请及时查看。'),
  (15,'store','直播审核通过主播通知','直播审核通过主播通知',0,0,0,'您的直播审核已通过。','您的直播审核已通过。','您的直播审核已通过。'),
  (16,'store','入驻申请通过提醒','入驻申请通过提醒',0,0,1,'您的入驻申请已通过，请登录后台完成后续配置。','您的入驻申请已通过，请登录后台完成后续配置。','您的入驻申请已通过，请登录后台完成后续配置。'),
  (18,'store','直播未通过通知','直播未通过通知',0,0,0,'您的直播审核未通过，请根据提示调整后重新提交。','您的直播审核未通过，请根据提示调整后重新提交。','您的直播审核未通过，请根据提示调整后重新提交。'),
  (24,'store','客服消息通知','用户咨询消息通知',0,0,0,'您有新的用户咨询，请及时回复。','您有新的用户咨询，请及时回复。','您有新的用户咨询，请及时回复。');

-- 物流公司：完整 CRMEB 快递目录（419 项），来源于 CRMEB_MER_v4.0/extend/express.xlsx。
-- 仅恢复标准名称和软删除状态，不覆盖后台维护的排序与显示状态。
INSERT INTO `qixi_crm_a_express` (`name`,`code`,`sort`,`is_show`,`is_del`) VALUES
  ('A2U速递','a2u',0,1,0),
  ('AAE快递','aae',0,1,0),
  ('爱彼西快递','abc',0,1,0),
  ('德方物流','ahdf',0,1,0),
  ('航空快递','airgtc',0,1,0),
  ('阿里物流','ALP',0,1,0),
  ('安得物流','ande',0,1,0),
  ('安捷快递','anjie',0,1,0),
  ('安能物流','anneng',0,1,0),
  ('安信达快递','anxinda',0,1,0),
  ('安迅物流','anxun',0,1,0),
  ('AOL快递','aol',0,1,0),
  ('AOL澳通速递','aolau',0,1,0),
  ('Aramex','aramex',0,1,0),
  ('方舟速递','arke',0,1,0),
  ('澳邮中国快运','auexpress',0,1,0),
  ('卡行天下','B2B',0,1,0),
  ('百千诚国际物流','baiqian',0,1,0),
  ('百腾物流','baitengwuliu',0,1,0),
  ('八梁物流','BALIANGWL',0,1,0),
  ('巴伦支快递','balunzhi',0,1,0),
  ('邦送物流','bangsongwuliu',0,1,0),
  ('宝通达物流','baotongda',0,1,0),
  ('BCWELT','bcwelt',0,1,0),
  ('奔腾物流','benteng',0,1,0),
  ('滨发物流','BFWL',0,1,0),
  ('布谷鸟快递','bgn',0,1,0),
  ('挂号信','bgpyghx',0,1,0),
  ('BHT','bht',0,1,0),
  ('华慧快递','BHTEXP',0,1,0),
  ('笨鸟海淘','birdex',0,1,0),
  ('速方国际物流','bphchina',0,1,0),
  ('百事亨通','bsht',0,1,0),
  ('百世快运','bsky',0,1,0),
  ('博源恒通','byht',0,1,0),
  ('河南次晨达','ccd',0,1,0),
  ('CCES快递','cces',0,1,0),
  ('长通物流','changtong',0,1,0),
  ('程光快递','chengguang',0,1,0),
  ('城际速递','chengji',0,1,0),
  ('城市100快递','chengshi100',0,1,0),
  ('同舟行物流','chinatzx',0,1,0),
  ('传志快递','chuanzhi',0,1,0),
  ('出口易','chukouyi',0,1,0),
  ('CityLink快递','citylink',0,1,0),
  ('CE易欧通国际速递','cloudexpress',0,1,0),
  ('GLS快递','CNGLS',0,1,0),
  ('中环快递','cnpex',0,1,0),
  ('东方快递','coe',0,1,0),
  ('城市之星','cszx',0,1,0),
  ('云南中诚','czwlyn',0,1,0),
  ('大达物流','dada',0,1,0),
  ('大金物流','dajin',0,1,0),
  ('大顺物流','dashun',0,1,0),
  ('达速物流','dasu',0,1,0),
  ('大田物流','datian',0,1,0),
  ('大洋物流快递','dayang',0,1,0),
  ('大众佐川急便','dazhong',0,1,0),
  ('德邦物流','debang',0,1,0),
  ('德创物流','dechuangwuliu',0,1,0),
  ('德中快递','decnlh',0,1,0),
  ('德坤供应链','dekuncn',0,1,0),
  ('达方物流','dfpost',0,1,0),
  ('DHL快递','dhl',0,1,0),
  ('店通快递','diantong',0,1,0),
  ('递达快递','dida',0,1,0),
  ('叮咚澳洲转运','dindon',0,1,0),
  ('递四方速递','disifang',0,1,0),
  ('东瀚物流','donghanwl',0,1,0),
  ('东红物流','donghong',0,1,0),
  ('东骏快捷物流','dongjun',0,1,0),
  ('DPEX快递','dpex',0,1,0),
  ('D速快递','dsu',0,1,0),
  ('易满客','ecmscn',0,1,0),
  ('益递物流','edlogistics',0,1,0),
  ('百福东方快递','ees',0,1,0),
  ('易联通达物流','el56',0,1,0),
  ('EMS','ems',0,1,0),
  ('俄顺达','eshunda',0,1,0),
  ('欧亚专线','euasia',0,1,0),
  ('EWE全球快递','ewe',0,1,0),
  ('安鲜达','exfresh',0,1,0),
  ('E邮宝','eyoubao',0,1,0),
  ('伍圆速递','F5XM',0,1,0),
  ('颿达国际快递','fandaguoji',0,1,0),
  ('方方达物流','fangfangda',0,1,0),
  ('凡宇速递','fanyu',0,1,0),
  ('泛远国际物流','farlogistis',0,1,0),
  ('FedEx英国','fedexuk',0,1,0),
  ('飞邦物流','feibang',0,1,0),
  ('飞豹快递','feibao',0,1,0),
  ('原飞航快递','feihang',0,1,0),
  ('飞狐快递','feihu',0,1,0),
  ('飞快达物流','feikuaida',0,1,0),
  ('飞特物流','feite',0,1,0),
  ('飞洋快递','feiyang',0,1,0),
  ('飞远物流','feiyuan',0,1,0),
  ('丰达快递','fengda',0,1,0),
  ('风行天下','fengxingtianxia',0,1,0),
  ('飞康达物流','fkd',0,1,0),
  ('飞力士物流','flysman',0,1,0),
  ('FOX国际速递','fox',0,1,0),
  ('港快速递','gangkuai',0,1,0),
  ('GATI快递','gaticn',0,1,0),
  ('广东ems快递','gdems',0,1,0),
  ('国际包裹','gjbg',0,1,0),
  ('英脉物流','gml',0,1,0),
  ('国内小包','gnxb',0,1,0),
  ('共速达物流','gongsuda',0,1,0),
  ('GSM','gsm',0,1,0),
  ('万通快递','gswtkd',0,1,0),
  ('GTS快递','gts',0,1,0),
  ('高铁速递','gtsd',0,1,0),
  ('冠达快递','guada',0,1,0),
  ('广东邮政','guangdongyouzhengwuliu',0,1,0),
  ('广通速递','guangtong',0,1,0),
  ('国通快递','guotong',0,1,0),
  ('文捷航空速递','GZWENJIE',0,1,0),
  ('山东海红快递','haihong',0,1,0),
  ('海盟速递','haimeng',0,1,0),
  ('海外环球','haiwaihuanqiu',0,1,0),
  ('航宇快递','hangyu',0,1,0),
  ('韩润物流','hanrun',0,1,0),
  ('好来运快递','haolaiyun',0,1,0),
  ('昊盛物流','haosheng',0,1,0),
  ('好又快物流','haoyoukuai',0,1,0),
  ('河北建华物流','hebeijianhua',0,1,0),
  ('恒诚物流','HENGCHENGWL',0,1,0),
  ('恒丰物流','HENGFENGWL',0,1,0),
  ('恒路物流','henglu',0,1,0),
  ('恒宇运通','hengyu',0,1,0),
  ('和丰同城','hfwuxi',0,1,0),
  ('黑狗物流','higo',0,1,0),
  ('海派通','hipito',0,1,0),
  ('猴急送','hjs',0,1,0),
  ('香港邮政','hkpost',0,1,0),
  ('飞鹰物流','hnfy',0,1,0),
  ('宏捷国际物流','hongjie',0,1,0),
  ('鸿讯物流','hongxun',0,1,0),
  ('环球通达','hqtd',0,1,0),
  ('汇通天下物流','httx56',0,1,0),
  ('华通务达物流','htwd',0,1,0),
  ('华诚物流','huacheng',0,1,0),
  ('华达快运','huada',0,1,0),
  ('华翰物流','huahan',0,1,0),
  ('华航快递','huahang',0,1,0),
  ('黄马甲快递','huangmajia',0,1,0),
  ('环球速运','huanqiu',0,1,0),
  ('华企快运','huaqi',0,1,0),
  ('华夏龙物流','huaxialong',0,1,0),
  ('天地华宇物流','huayu',0,1,0),
  ('辉联物流','huilian',0,1,0),
  ('汇强快递','huiqiang',0,1,0),
  ('百世快递','huitong',0,1,0),
  ('汇文配送','huiwen',0,1,0),
  ('伙伴物流','huoban',0,1,0),
  ('户通物流','hutongwuliu',0,1,0),
  ('百成大达物流','idada',0,1,0),
  ('中国邮政','intmail',0,1,0),
  ('京东快递','jd',0,1,0),
  ('景光物流','jgwl',0,1,0),
  ('佳惠尔快递','jiahuier',0,1,0),
  ('佳吉物流','jiaji',0,1,0),
  ('佳家通','jiajiatong56',0,1,0),
  ('佳怡物流','jiayi',0,1,0),
  ('佳宇物流','JIAYU',0,1,0),
  ('加运美物流','jiayunmei',0,1,0),
  ('捷特快递','jiete',0,1,0),
  ('锦程国际物流','jinchengwuliu',0,1,0),
  ('金大物流','jindawuliu',0,1,0),
  ('京世物流','jingshi',0,1,0),
  ('京广速递快件','jinguangsudikuaijian',0,1,0),
  ('晋越快递','jinyue',0,1,0),
  ('九曳供应链','jiuye',0,1,0),
  ('久易快递','jiuyi',0,1,0),
  ('急先达物流','jixianda',0,1,0),
  ('嘉里大通','jldt',0,1,0),
  ('金马甲','jmjss',0,1,0),
  ('日本邮政','jppost',0,1,0),
  ('吉日优派','jrypex',0,1,0),
  ('骏川物流','JUNCHUANWL',0,1,0),
  ('骏丰国际速递','junfengguoji',0,1,0),
  ('吉祥邮','jxy',0,1,0),
  ('康力物流','klwl',0,1,0),
  ('直邮易','kuachangwuliu',0,1,0),
  ('快捷速递','kuaijie',0,1,0),
  ('快淘快递','kuaitao',0,1,0),
  ('快优达速递','kuaiyouda',0,1,0),
  ('宽容物流','kuanrong',0,1,0),
  ('跨越快递','kuayue',0,1,0),
  ('蓝镖快递','lanbiao',0,1,0),
  ('蓝弧快递','lanhu',0,1,0),
  ('宝凯物流','lbbk',0,1,0),
  ('联邦物流','LBWL',0,1,0),
  ('乐递供应链','ledii',0,1,0),
  ('乐捷递','lejiedi',0,1,0),
  ('云豹国际货运','leopard',0,1,0),
  ('联昊通快递','lianhaotong',0,1,0),
  ('成都立即送快递','lijisong',0,1,0),
  ('利民物流','LIMINWL',0,1,0),
  ('一号线','lineone',0,1,0),
  ('龙邦快运','longbang',0,1,0),
  ('隆浪快递','longlangkuaidi',0,1,0),
  ('龙胜物流','LONGSHENWL',0,1,0),
  ('恒通快递','lqht',0,1,0),
  ('乐天速递','ltexp',0,1,0),
  ('论道国际物流','lundao',0,1,0),
  ('鲁通快运','lutong',0,1,0),
  ('麦力快递','mailikuaidi',0,1,0),
  ('木春货运','mchy',0,1,0),
  ('美国快递','meiguo',0,1,0),
  ('美龙快递','meilong',0,1,0),
  ('美快国际物流','meiquick',0,1,0),
  ('美西快递','meixi',0,1,0),
  ('门对门','menduimen',0,1,0),
  ('蒙速快递','mengsu',0,1,0),
  ('民邦速递','minbang',0,1,0),
  ('明亮物流','mingliang',0,1,0),
  ('民航快递','minhang',0,1,0),
  ('闽盛物流','minsheng',0,1,0),
  ('南北快递','nanbei',0,1,0),
  ('中国南方航空股份有限公司','NANHANG',0,1,0),
  ('红马速递','nedahm',0,1,0),
  ('港中能达快递','nengda',0,1,0),
  ('新蛋奥硕物流','neweggozzo',0,1,0),
  ('华赫物流','nmhuahe',0,1,0),
  ('腾达速递','nntengda',0,1,0),
  ('偌亚奥国际','nuoyaao',0,1,0),
  ('OCS国际快递','ocs',0,1,0),
  ('一号仓','onehcang',0,1,0),
  ('onTrac','ontrac',0,1,0),
  ('中欧快运','otobv',0,1,0),
  ('澳大利亚PCA快递','pca',0,1,0),
  ('配思货运','PEISI',0,1,0),
  ('陪行物流','peixing',0,1,0),
  ('彪记快递','PEWKEE',0,1,0),
  ('皇家物流','pfcexpress',0,1,0),
  ('凤凰快递','PHOENIXEXP',0,1,0),
  ('平安达快递','pinganda',0,1,0),
  ('平安达腾飞','pingandatengfei',0,1,0),
  ('小包','pingyou',0,1,0),
  ('品骏快递','pjbest',0,1,0),
  ('贝邮宝','ppbyb',0,1,0),
  ('急顺通','pzhjst',0,1,0),
  ('秦邦快运','qbexpress',0,1,0),
  ('启辰国际物流','qichen',0,1,0),
  ('秦远物流','qinyuan',0,1,0),
  ('千顺快递','qskdyxgs',0,1,0),
  ('全晨快递','quanchen',0,1,0),
  ('全峰快递','quanfeng',0,1,0),
  ('全际通快递','quanjitong',0,1,0),
  ('全日通快递','quanritong',0,1,0),
  ('全速快运','quansu',0,1,0),
  ('全速通国际快递','quansutong',0,1,0),
  ('全信通快递','quanxintong',0,1,0),
  ('全一快递','quanyi',0,1,0),
  ('全之鑫物流','qzx56',0,1,0),
  ('日日顺物流','ririshun',0,1,0),
  ('日昱物流','riyu',0,1,0),
  ('荣庆物流','rongqing',0,1,0),
  ('RPX保时达','rpx',0,1,0),
  ('捷网俄全通','ruexp',0,1,0),
  ('如风达快递','rufeng',0,1,0),
  ('凡客如风达','rufengda',0,1,0),
  ('瑞达国际速递','ruidaex',0,1,0),
  ('瑞丰速递','ruifeng',0,1,0),
  ('全时速运','runhengfeng',0,1,0),
  ('日益通速递','rytsd',0,1,0),
  ('赛澳递','saiaodi',0,1,0),
  ('三态速递','santai',0,1,0),
  ('丰程物流','sccod',0,1,0),
  ('泰国138','sd138',0,1,0),
  ('优配速运','sdyoupei',0,1,0),
  ('速递中国','sendtochina',0,1,0),
  ('七天连锁','sevendays',0,1,0),
  ('十方通物流','sfift',0,1,0),
  ('圣安物流','shengan',0,1,0),
  ('晟邦物流','shengbang',0,1,0),
  ('盛丰物流','shengfeng',0,1,0),
  ('盛辉物流','shenghui',0,1,0),
  ('申通快递','shentong',0,1,0),
  ('昊昕物流','SHHX',0,1,0),
  ('世运快递','shiyun',0,1,0),
  ('上海林道货运','shlindao',0,1,0),
  ('顺发物流','SHUNFAWL',0,1,0),
  ('顺丰速运','shunfeng',0,1,0),
  ('顺捷丰达','shunjiefengda',0,1,0),
  ('四海快递','sihaiet',0,1,0),
  ('思迈快递','simai',0,1,0),
  ('信联通','sinatone',0,1,0),
  ('新加坡邮政','singpost',0,1,0),
  ('宋军物流','SJWL',0,1,0),
  ('荷兰','Sky',0,1,0),
  ('春风物流','spring',0,1,0),
  ('星晨急便','STARS',0,1,0),
  ('顺通快递','stkd',0,1,0),
  ('速必达物流','subida',0,1,0),
  ('速呈宅配','suchengzhaipei',0,1,0),
  ('穗佳物流','suijia',0,1,0),
  ('郑州速捷','sujievip',0,1,0),
  ('上大物流','SUNDAPOST',0,1,0),
  ('苏宁快递','suning',0,1,0),
  ('新杰物流','sunjex',0,1,0),
  ('新速航','sunspeedy',0,1,0),
  ('速尔物流','sure',0,1,0),
  ('速腾快递','suteng',0,1,0),
  ('速通物流','sutong',0,1,0),
  ('苏粤货运','SUYUE',0,1,0),
  ('盛旺货运','SWHY',0,1,0),
  ('山西红马甲','sxhongmajia',0,1,0),
  ('沈阳佳惠尔','syjiahuier',0,1,0),
  ('华宇物流','tiandihuayu',0,1,0),
  ('天河物流','TIANHEWL',0,1,0),
  ('天天快递','tiantian',0,1,0),
  ('天纵物流','tianzong',0,1,0),
  ('万家通','timedg',0,1,0),
  ('天联快运','tlky',0,1,0),
  ('TNT快递','tnt',0,1,0),
  ('通成物流','tongcheng',0,1,0),
  ('通达兴物流','tongdaxing',0,1,0),
  ('通和天下物流','tonghe',0,1,0),
  ('中运全速','topspeedex',0,1,0),
  ('汤氏物流','TSWL',0,1,0),
  ('合众速递','ucs',0,1,0),
  ('UEQ快递','ueq',0,1,0),
  ('UEX','uex',0,1,0),
  ('UPS快递','ups',0,1,0),
  ('USPS快递','usps',0,1,0),
  ('美通快递','valueway',0,1,0),
  ('鹰运国际速递','vipexpress',0,1,0),
  ('万博快递','wanbo',0,1,0),
  ('万家物流','wanjia',0,1,0),
  ('万象物流','wanxiang',0,1,0),
  ('微特派快递','weitepai',0,1,0),
  ('渥途国际速运','wotu',0,1,0),
  ('威时沛运','wtdchina',0,1,0),
  ('五环速递','wuhuan',0,1,0),
  ('微转运','wzhaunyun',0,1,0),
  ('西安胜峰','xaetc',0,1,0),
  ('鑫飞鸿物流快递','XFHONG',0,1,0),
  ('西安城联速递','xianchenglian',0,1,0),
  ('先锋快递','xianfeng',0,1,0),
  ('北青小红帽','xiaohongmao',0,1,0),
  ('喜来快递','xilaikd',0,1,0),
  ('新邦物流','xinbang',0,1,0),
  ('新蛋物流','xindan',0,1,0),
  ('信丰物流','xinfeng',0,1,0),
  ('星程宅配','xingchengzhaipei',0,1,0),
  ('鑫天顺物流','XINTIAN',0,1,0),
  ('信天捷快递','xintianjie',0,1,0),
  ('西邮寄','xipost',0,1,0),
  ('希优特快递','xiyoute',0,1,0),
  ('祥龙运通','xlyt',0,1,0),
  ('鑫世锐达','xsrd',0,1,0),
  ('鑫通宝物流','xtb',0,1,0),
  ('源安达快递','yad',0,1,0),
  ('亚风速递','yafeng',0,1,0),
  ('亚马逊物流','yamaxunwuliu',0,1,0),
  ('燕文物流','yanwen',0,1,0),
  ('邮联物流','YBWL',0,1,0),
  ('远成快运','ycgky',0,1,0),
  ('一邦快递','yibang',0,1,0),
  ('易达通快递','yidatong',0,1,0),
  ('亿领速运','yilingsuyun',0,1,0),
  ('英超物流','yingchao',0,1,0),
  ('顺捷丰达','yinjie',0,1,0),
  ('音速速运','yinsu',0,1,0),
  ('一柒物流','yiqiguojiwuliu',0,1,0),
  ('亿顺航','yishunhang',0,1,0),
  ('易通达','yitongda',0,1,0),
  ('永昌物流','yongchang',0,1,0),
  ('永旺达快递','yongwangda',0,1,0),
  ('邮必佳','youbijia',0,1,0),
  ('UC优速快递','youshuwuliu',0,1,0),
  ('优速快递','yousu',0,1,0),
  ('优速通达','yousutongda',0,1,0),
  ('挂号信','youzhengguonei',0,1,0),
  ('壹品速递','ypsd',0,1,0),
  ('一统飞鸿快递','ytfh',0,1,0),
  ('远成物流','yuancheng',0,1,0),
  ('圆通速递','yuantong',0,1,0),
  ('圆圆物流','YUANYUANWL',0,1,0),
  ('元智捷诚快递','yuanzhi',0,1,0),
  ('越丰快递','yuefeng',0,1,0),
  ('御风速运','yufeng',0,1,0),
  ('煜嘉物流','yujiawuliu',0,1,0),
  ('誉美捷快递','yumeijie',0,1,0),
  ('韵达快运','yunda',0,1,0),
  ('韵达美国件','yundaexus',0,1,0),
  ('运通快递','yuntong',0,1,0),
  ('云物流','yunwuliu',0,1,0),
  ('宇鑫物流','yuxin',0,1,0),
  ('源伟丰快递','ywfex',0,1,0),
  ('宇鑫物流','yxwl',0,1,0),
  ('远洋国际','yyexpress',0,1,0),
  ('一运全成物流','yyqc56',0,1,0),
  ('增益快递','zengyi',0,1,0),
  ('增益速递','zengyisudi',0,1,0),
  ('振刚物流','ZGWL',0,1,0),
  ('宅急送','zhaijisong',0,1,0),
  ('众辉达物流','zhdwl',0,1,0),
  ('至诚通达快递','zhichengtongda',0,1,0),
  ('芝麻开门','zhima',0,1,0),
  ('中睿速递','zhongruisudi',0,1,0),
  ('中速快件','zhongsukuaidi',0,1,0),
  ('中天快运','zhongtian',0,1,0),
  ('中铁快运','zhongtie',0,1,0),
  ('中通快递','zhongtong',0,1,0),
  ('中外运速递','zhongwaiyun',0,1,0),
  ('中信达快递','zhongxinda',0,1,0),
  ('中邮物流','zhongyou',0,1,0),
  ('纵行物流','zongxing',0,1,0),
  ('准实快运','zsky123',0,1,0),
  ('中铁物流','ZTKY',0,1,0),
  ('智通物流','ztong',0,1,0),
  ('中天万运快递','ztwy',0,1,0),
  ('佐川急便','zuochuan',0,1,0),
  ('中外速运','zwsy',0,1,0),
  ('郑州建华快递','zzjh',0,1,0)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`), `is_del`=0;

INSERT IGNORE INTO `qixi_crm_a_serve_meal`
  (`meal_id`,`name`,`type`,`price`,`num`,`sort`,`status`,`is_del`,`create_time`) VALUES
  (151101,'商品采集基础套餐',1,0.00,0,30,1,0,'2026-05-15 01:17:40'),
  (151102,'电子面单标准套餐',2,1.00,1,20,1,0,'2025-12-19 21:03:11'),
  (151103,'商品采集扩展套餐',1,0.00,8,10,1,0,'2025-12-06 12:49:50');

INSERT IGNORE INTO `qixi_crm_a_config_item`
  (`id`,`item_type`,`name`,`code`,`remark`,`payload`,`status`,`sort`,`is_del`) VALUES
  (81101,'hot_search','海鲜','hot-seafood','默认热门搜索',JSON_OBJECT('weight',10),1,10,0),
  (81102,'hot_search','预制菜','hot-ready-meal','默认热门搜索',JSON_OBJECT('weight',20),1,20,0),
  (81103,'hot_search','国潮','hot-chinese-style','默认热门搜索',JSON_OBJECT('weight',30),1,30,0),
  (81104,'hot_search','箱包','hot-bags','默认热门搜索',JSON_OBJECT('weight',40),1,40,0),
  (81105,'hot_search','家居','hot-home','默认热门搜索',JSON_OBJECT('weight',50),1,50,0),
  (81106,'hot_search','口红','hot-lipstick','默认热门搜索',JSON_OBJECT('weight',60),1,60,0),
  (81107,'hot_search','运动鞋','hot-sneakers','默认热门搜索',JSON_OBJECT('weight',70),1,70,0),
  (81108,'hot_search','手机','hot-mobile','默认热门搜索',JSON_OBJECT('weight',80),1,80,0),
  (81202,'system_form','通用售后申请表','aftersale-form','默认系统表单',JSON_OBJECT('fields',JSON_ARRAY('reason','images')),1,10,0);

-- 配置类默认数据：不包含商品、店铺、用户模拟数据，也不写入密钥值。
INSERT IGNORE INTO `qixi_crm_a_config_classification`
  (`id`,`parent_id`,`name`,`classify_key`,`description`,`icon`,`status`,`sort`,`is_del`) VALUES
  (81401,0,'公众号配置','wechat','公众号基础配置','lucide:message-circle',1,100,0),
  (81402,0,'短信配置','message','平台短信验证码配置','lucide:message-square',1,90,0),
  (81403,0,'商户基础配置','mer_base','商户基础配置','lucide:store',1,80,0),
  (81404,0,'对象存储配置','storage','对象存储与媒体资源配置','lucide:cloud',1,70,0),
  (81405,0,'小程序配置','smallapp','小程序基础配置','lucide:smartphone',1,60,0),
  (81406,0,'余额/充值设置','balance','余额与充值规则配置','lucide:wallet',1,50,0),
  (81407,0,'商城基础设置','mall','商城默认行为配置','lucide:shopping-bag',1,40,0),
  (81408,0,'客服设置','customer_service','客服接待与自动回复配置','lucide:headphones',1,30,0),
  (81409,0,'通知设置','notification','公众号、小程序与短信通知配置','lucide:bell',1,20,0);

-- 默认配置项只写入非敏感、可公开维护的业务默认值；云服务密钥等仅允许保存在本地 init_key.sql。
INSERT IGNORE INTO `qixi_crm_a_config_classification_item`
  (`id`,`classification_id`,`name`,`config_key`,`field_type`,`backend_type`,`content`,`description`,`status`,`sort`,`is_del`) VALUES
  (81501,81407,'网站名称','site_name','input',0,'商城','平台网站名称',1,100,0),
  (81502,81407,'网站开启','site_open','switch',0,'1','关闭后仅允许平台后台访问',1,90,0),
  (81503,81407,'自动解析复制口令','parse_copy_command','switch',0,'1','开启后小程序和 App 自动读取剪贴板口令',1,80,0),
  (81504,81407,'默认赠送复制次数','default_copy_count','number',1,'8','默认给商户赠送的商品采集次数',1,70,0),
  (81505,81402,'启用短信验证码','sms_verification_enabled','switch',0,'1','平台短信验证码开关',1,100,0),
  (81506,81405,'小程序名称','mini_program_name','input',0,'商城小程序','移动端小程序显示名称',1,100,0),
  (81507,81408,'默认客服类型','default_customer_service_type','radio',0,'system','默认由平台系统客服在线接待',1,100,0),
  (81508,81409,'订单支付成功通知','order_payment_notice_enabled','switch',0,'1','订单支付成功后通知会员',1,100,0),
  (81509,81403,'商户入驻开关','merchant_admission_enabled','switch',0,'1','允许商户提交入驻申请',1,100,0),
  (81510,81406,'余额功能开关','balance_enabled','switch',0,'0','平台余额与充值功能开关',1,100,0);

-- 组合数据属于系统默认配置，不使用商品、店铺、用户等演示数据。
INSERT IGNORE INTO `qixi_crm_a_data_group`
  (`id`,`name`,`group_key`,`description`,`fields`,`sort`,`is_del`) VALUES
  (81306,'开屏广告','open_screen_advertising','移动端开屏广告配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),200,0),
  (81307,'积分范围配置','points_mall_scope','积分商城范围配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','积分范围','field','scope','type','input')),190,0),
  (81308,'积分金刚区','points_mall_district','积分商城金刚区配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),180,0),
  (81309,'积分商城轮播图','points_mall_banner','积分商城轮播图配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),170,0),
  (81310,'付费会员购买类型','svip_pay','付费会员购买套餐配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','价格','field','price','type','number')),160,0),
  (81311,'首页顶部广告','pc_top_banner','PC 首页顶部广告图配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),150,0),
  (81312,'社区热门搜索','community_hot_keyword','社区热门搜索配置',JSON_ARRAY(JSON_OBJECT('name','关键词','field','keyword','type','input')),140,0),
  (81313,'平台专题管理','sys_activity','平台专题配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','封面图','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),130,0),
  (81314,'商户专题管理','mer_activity','商户专题配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','封面图','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),120,0),
  (81315,'服务套餐','service_meal','商户可购买的服务套餐配置',JSON_ARRAY(JSON_OBJECT('name','套餐名称','field','name','type','input'),JSON_OBJECT('name','价格','field','price','type','number')),115,0),
  (81316,'PC 底部版权信息','pc_copyright','PC 页脚版权信息配置',JSON_ARRAY(JSON_OBJECT('name','版权内容','field','content','type','textarea')),110,0),
  (81317,'商户 PC 轮播图','mer_pc_banner','商户 PC 轮播图配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),105,0),
  (81318,'PC 首页分类广场','pc_home_rec','PC 首页分类广场配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),95,0),
  (81319,'PC 首页轮播图','pc_home_banner','PC 首页轮播图配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),85,0),
  (81301,'热门搜索','hot_keyword','移动端搜索页热门搜索设置',JSON_ARRAY(JSON_OBJECT('name','关键词','field','keyword','type','input')),100,0),
  (81302,'首页轮播图','home_banner','移动端首页轮播图配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),90,0),
  (81303,'首页菜单','home_menu','移动端首页菜单配置',JSON_ARRAY(JSON_OBJECT('name','名称','field','name','type','input'),JSON_OBJECT('name','图片','field','pic','type','image'),JSON_OBJECT('name','链接','field','url','type','input')),80,0),
  (81304,'提现银行列表','bank_list','用户提现可选银行配置',JSON_ARRAY(JSON_OBJECT('name','银行名称','field','name','type','input')),70,0),
  (81305,'连续签到奖励','sign_day_config','连续签到积分奖励配置',JSON_ARRAY(JSON_OBJECT('name','签到天数','field','sign_day','type','number'),JSON_OBJECT('name','赠送积分','field','sign_integral','type','number')),60,0);
INSERT IGNORE INTO `qixi_crm_a_data_group_item`
  (`id`,`group_id`,`data`,`sort`,`status`,`is_del`) VALUES
  (81401,81301,JSON_OBJECT('keyword','热销'),100,1,0),
  (81402,81301,JSON_OBJECT('keyword','新品'),90,1,0),
  (81403,81304,JSON_OBJECT('name','中国工商银行'),100,1,0),
  (81404,81304,JSON_OBJECT('name','中国建设银行'),90,1,0),
  (81405,81305,JSON_OBJECT('sign_day',1,'sign_integral',1),100,1,0),
  (81406,81305,JSON_OBJECT('sign_day',7,'sign_integral',10),90,1,0);

INSERT INTO `qixi_crm_a_platform_category` (`id`,`parent_id`,`name`,`sort`,`status`) VALUES
  (7601,0,'维修服务',100,1),(7602,0,'家居百货',90,1),(7603,0,'服饰',80,1),(7604,0,'美妆',70,1),(7605,0,'数码',60,1),
  (7611,7601,'清洗保养',100,1),(7612,7602,'日用品',100,1),(7613,7603,'女装',100,1),(7614,7604,'香水彩妆',100,1),(7615,7605,'数码设备',100,1),
  (7621,7611,'上门维修',100,1),(7622,7611,'到店维修',90,1),(7623,7611,'家电维修',80,1),
  (7624,7612,'生活日用',100,1),(7625,7612,'床上用品',90,1),(7626,7612,'母婴用品',80,1),
  (7627,7613,'家居服',100,1),(7628,7613,'连衣裙',90,1),(7629,7613,'时尚潮牌',80,1),
  (7630,7614,'洗发护发',100,1),(7631,7614,'面部护肤',90,1),(7632,7614,'香水',80,1),
  (7633,7615,'数码配件',100,1),(7634,7615,'智能设备',90,1),(7635,7615,'电脑',80,1),(7636,7615,'手机',70,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

-- 入驻、商品元数据与价格规则属于平台可维护的默认配置，不依赖任何演示店铺或商品。
INSERT IGNORE INTO `qixi_crm_a_merchant_category` (`id`,`name`,`commission_rate`,`status`) VALUES
  (701,'综合零售',8.50,1),(702,'生活服务',6.00,1);
INSERT IGNORE INTO `qixi_crm_a_merchant_type` (`id`,`name`,`type_info`,`is_margin`,`margin`,`description`,`remark`,`status`) VALUES
  (711,'标准店铺','适用于常规经营主体',0,0.00,'请按平台规则完成入驻与经营。','',1),
  (712,'保证金店铺','适用于需要缴纳保证金的经营主体',1,500.00,'保证金规则以平台公示内容为准。','',1);
INSERT IGNORE INTO `qixi_crm_a_merchant_type_menu` (`merchant_type_id`,`menu_code`) VALUES
  (711,'merchant.dashboard'),(712,'merchant.dashboard'),(712,'merchant.catalog');

INSERT IGNORE INTO `qixi_crm_a_product_label` (`id`,`name`,`description`,`color`,`sort`,`status`) VALUES
  (7501,'七天无理由','符合规则的商品支持七天无理由退货。','#16a34a',100,1),
  (7502,'新品','平台新品标识。','#2563eb',90,1);
INSERT IGNORE INTO `qixi_crm_a_product_guarantee` (`id`,`name`,`content`,`icon_url`,`sort`,`status`,`mer_count`,`product_count`) VALUES
  (7511,'正品保障','商品来源与售后承诺以平台规则和订单约定为准。','',100,1,0,0),
  (7512,'极速退款','符合规则的退款将按售后流程处理。','',90,1,0,0);
INSERT IGNORE INTO `qixi_crm_a_product_parameter_template`
  (`id`,`name`,`cate_ids_json`,`params_json`,`values_json`,`sort`,`status`) VALUES
  (7521,'服饰规格',JSON_ARRAY(7628),JSON_ARRAY(JSON_OBJECT('name','颜色','values',JSON_ARRAY('中国红','竹青色','云白色'),'required',0,'sort',100),JSON_OBJECT('name','尺码','values',JSON_ARRAY('S','M','L'),'required',0,'sort',90)),JSON_ARRAY('中国红','竹青色','云白色'),0,1),
  (7522,'通用参数',JSON_ARRAY(),JSON_ARRAY(JSON_OBJECT('name','品牌','values',JSON_ARRAY(),'required',0,'sort',100),JSON_OBJECT('name','产地','values',JSON_ARRAY(),'required',0,'sort',90)),JSON_ARRAY(),0,1);
INSERT IGNORE INTO `qixi_crm_a_product_price_rule`
  (`id`,`name`,`cate_ids_json`,`is_default`,`content`,`sort`,`status`) VALUES
  (7531,'家电价格说明',JSON_ARRAY(7623),'0','<p>商品价格、规格与运费以商品详情页及订单结算页为准。</p>',1,1),
  (7532,'服饰价格说明',JSON_ARRAY(7613),'0','<p>服饰商品请以商品规格、尺码与订单结算页信息为准。</p>',0,1),
  (7533,'通用价格说明',JSON_ARRAY(),'1','<p>未指定分类时默认适用于全部商品。</p>',0,1);
