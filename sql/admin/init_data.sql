SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_role` (`code`,`name`,`status`) VALUES
  ('platform','平台管理',1),('merchant','商户管理',1),('region','区域管理',1),
  ('customer_service','客服管理',1),('operations','运营管理',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`);

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
  (1501,1500,'setting.serve','一号通','lucide:cloud','/serve','directory',1),
  (1510,1501,'setting.serve.login','登陆入口','lucide:log-in','/setting/sms/sms_config/index','page',1),
  (1511,1501,'setting.serve.config','服务配置','lucide:settings-2','/service/settings','page',2),
  (1512,1501,'setting.serve.sms','短信设置','lucide:mail','/sms','page',3),
  (1502,1500,'setting.shop.dir','商城设置','lucide:store','/shop','directory',2),
  (1520,1502,'setting.shop.form','商城设置','lucide:store','/systemForm/Basics/shop_tabs','page',1),
  (1521,1502,'setting.shop.hot','热门搜索','lucide:search','/group/config/67','page',2),
  (1522,1502,'setting.shop.agreements','协议规则','lucide:file-text','/setting/agreements','page',3),
  (1503,1500,'setting.delivery.dir','配送配置','lucide:truck','/delivery_config','directory',3),
  (1530,1503,'setting.delivery.express','物流公司','lucide:truck','/freight/express','page',1),
  (1504,1500,'setting.rbac.dir','权限管理','lucide:shield-check','/setting/rbac','directory',4),
  (1540,1504,'setting.rbac.role','角色权限','lucide:shield-check','/setting/role','page',1),
  (1541,1504,'setting.rbac.admin','管理员管理','lucide:user-round-cog','/setting/admin','page',2),
  (1542,1504,'setting.rbac.menu','菜单管理','lucide:menu','/setting/menu','page',3),
  (1505,1500,'setting.notice.dir','消息管理','lucide:bell','/notice','directory',5),
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
  (20936,126,'setting.sms.manage','维护无密钥短信 stub 配置','','setting/sms','button',1),
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
  ,(20993,202,'app.routine.manage','维护小程序基础开关','','app/routine','button',1)
  ,(21122,217,'app.mobile.manage','维护 iOS、Android、HarmonyOS 应用配置','','app/mobile','button',1)
  ,(21123,217,'app.push.manage','维护 App 推送配置','','app/mobile','button',2)
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
      'content.article.manage','content.article_category.read','content.article_category.manage','content.community_category.read','content.community_category.manage','content.community_topic.read','content.community_topic.manage','content.community_list.read','content.community_list.manage','content.community_reply.read','content.community_reply.manage','content.notice.manage','setting.agreement.manage',
      'operations.diy.manage','marketing.seckill.manage','marketing.combination.manage','marketing.presell.manage','marketing.coupon.manage','marketing.assist.manage','marketing.points.manage','marketing.recharge.manage',
      'maintain','maintain.cache','maintain.cache.manage','maintain.backup','maintain.backup.manage','maintain.group_data','maintain.group_data.manage','maintain.hot_search','maintain.hot_search.manage',
      'content.attachment.manage','content.community.audit','content.community.delete','marketing.broadcast.audit'));

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
