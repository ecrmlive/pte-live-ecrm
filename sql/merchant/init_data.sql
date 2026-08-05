SET NAMES utf8mb4;
USE `qixi_crm_merchant`;
INSERT INTO `qixi_crm_m_merchant` (`id`,`name`,`status`,`region_id`) VALUES
  (1,'七禧演示茶铺',1,10),(2,'七禧居家优选店',1,20),(3,'七禧数码生活店',1,10)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`),`region_id`=VALUES(`region_id`);
INSERT INTO `qixi_crm_m_store` (`id`,`merchant_id`,`app_id`,`name`,`status`) VALUES
  (1,1,'qixi.store.demo.1','七禧演示茶铺',1),(2,2,'qixi.store.demo.2','七禧居家优选店',1),(3,3,'qixi.store.demo.3','七禧数码生活店',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`app_id`=VALUES(`app_id`),`name`=VALUES(`name`),`status`=VALUES(`status`);

-- 平台商品监管中文夹具：区域 10 对应商户 1、3，区域 20 对应商户 2；不含真实个人或商户资料。
INSERT INTO `qixi_crm_m_product` (`id`,`store_id`,`title`,`category_id`,`status`,`version`) VALUES
  (5001,1,'七禧春日茉莉花茶礼盒',101,'on_sale',1),
  (5002,2,'居家香氛扩香礼盒',102,'pending_review',1),
  (5003,3,'智能温控随行杯',103,'draft',1)
ON DUPLICATE KEY UPDATE `store_id`=VALUES(`store_id`),`title`=VALUES(`title`),`category_id`=VALUES(`category_id`),`status`=VALUES(`status`),`version`=VALUES(`version`);
INSERT INTO `qixi_crm_m_product_detail` (`product_id`,`brief`,`keyword`,`unit_name`,`cover_url`,`original_price`) VALUES
  (5001,'虚构中文商品：清香茉莉花茶礼盒，供平台商品范围验收。','茉莉花茶,礼盒','盒','/demo/product-tea-v1.png',129.00),
  (5002,'虚构中文商品：居家香氛礼盒，处于待审核状态。','香氛,家居','盒','/demo/product-fragrance-v1.png',199.00),
  (5003,'虚构中文商品：智能温控随行杯草稿。','温控杯,数码','个','/demo/product-tumbler-v1.png',259.00)
ON DUPLICATE KEY UPDATE `brief`=VALUES(`brief`),`keyword`=VALUES(`keyword`),`unit_name`=VALUES(`unit_name`),`cover_url`=VALUES(`cover_url`),`original_price`=VALUES(`original_price`);
INSERT INTO `qixi_crm_m_product_sku` (`id`,`product_id`,`spec_json`,`price`,`stock`,`status`) VALUES
  (5101,5001,JSON_OBJECT('规格','250克'),99.00,60,1),(5102,5002,JSON_OBJECT('规格','礼盒装'),159.00,36,1),(5103,5003,JSON_OBJECT('颜色','深空灰'),199.00,48,1)
ON DUPLICATE KEY UPDATE `product_id`=VALUES(`product_id`),`spec_json`=VALUES(`spec_json`),`price`=VALUES(`price`),`stock`=VALUES(`stock`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`) VALUES
  (1,0,'dashboard','经营概览','/dashboard','views/ecrm/dashboard/index.vue','lucide:layout-dashboard',1,1,1,1),
  (10,0,'product','商品','/product','', 'lucide:package',1,0,10,1),
  (11,10,'product.list','商品管理','/product/list','views/ecrm/product/list.vue','',1,1,1,1),
  (12,10,'product.specs','规格模板','/product/specs','views/ecrm/product/specs.vue','',1,1,2,1),
  (13,10,'product.category','商品分类','/product/classify','views/ecrm/product/category.vue','',1,1,3,1),
  (20,0,'order','订单','/order','', 'lucide:receipt-text',1,0,20,1),
  (21,20,'order.list','订单管理','/order/list','views/ecrm/order/list.vue','',1,1,1,1),
  (22,20,'order.refund','售后管理','/order/refund','views/ecrm/order/refund.vue','',1,1,2,1),
  (23,20,'order.verify','核销记录','/order/verify','views/ecrm/order/verify.vue','',1,1,3,1),
  (30,0,'marketing','营销','/marketing','', 'lucide:megaphone',1,0,30,1),
  (31,30,'marketing.coupon','优惠券','/marketing/coupon','views/ecrm/marketing/coupon.vue','',1,1,1,1),
  (32,30,'marketing.seckill','秒杀','/marketing/seckill/list','views/ecrm/marketing/seckill.vue','',1,1,2,1),
  (33,30,'marketing.assist','好友助力','/marketing/assist/list','views/ecrm/marketing/assist.vue','',1,1,3,1),
  (34,30,'marketing.combination','拼团','/marketing/combination/combination_goods','views/ecrm/marketing/combination.vue','',1,1,4,1),
  (40,0,'finance','财务','/finance','', 'lucide:wallet-cards',1,0,40,1),
  (41,40,'finance.balance','资金账户','/finance/balance','views/ecrm/finance/balance.vue','',1,1,1,1),
  (42,40,'finance.withdraw','提现管理','/finance/withdraw','views/ecrm/finance/withdraw.vue','',1,1,2,1),
  (43,40,'finance.settlement','结算管理','/finance/settlement','views/ecrm/finance/settlement.vue','',1,1,3,1),
  (50,0,'setting','设置','/setting','', 'lucide:settings',1,0,50,1),
  (51,50,'setting.shop','店铺设置','/setting/shop','views/ecrm/setting/shop.vue','',1,1,1,1),
  (52,50,'setting.staff','员工管理','/setting/staff','views/ecrm/setting/staff.vue','',1,1,2,1),
  (53,50,'setting.admins','账号管理','/setting/admins','views/ecrm/setting/admins.vue','',1,1,3,1),
  (54,50,'setting.role','角色权限','/setting/role','views/ecrm/setting/role.vue','',1,1,4,1),
  (55,50,'setting.payment','支付方式','/setting/payment','views/ecrm/setting/payment.vue','',1,1,5,1),
  (56,50,'setting.im_sdk_app','IM SDK AppId','/setting/im-sdk-app','views/ecrm/setting/im-sdk-app.vue','',1,1,6,1),
  (57,50,'setting.integral','积分抵扣','/setting/integral-policy','views/ecrm/setting/integral.vue','',1,1,7,1),
  (14,10,'product.reservation','预约设置','/product/reservation','views/ecrm/marketing/reservation.vue','',1,1,4,1),
  (24,20,'order.reservation','预约订单','/order/reservation','views/ecrm/marketing/reservation.vue','',1,1,4,1),
  (60,0,'config','配送设置','/config','', 'lucide:truck',1,0,60,1),
  (61,60,'config.freight.templates','运费模板','/config/freight/shippingTemplates','views/ecrm/shipping/templates.vue','',1,1,1,1),
  (62,60,'config.freight.express','物流公司','/config/freight/express','views/ecrm/shipping/express.vue','',1,1,2,1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`path`=VALUES(`path`),`component`=VALUES(`component`),`icon`=VALUES(`icon`),`is_menu`=VALUES(`is_menu`),`is_route`=VALUES(`is_route`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

-- 新增店铺后台页面菜单：须与 admin-merchant registry.ts 中 PATH_COMPONENT 对齐。
INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`) VALUES
  (70,0,'statistic','数据统计','/statistic','','lucide:bar-chart-3',1,0,15,1),
  (71,70,'statistic.order','订单统计','/statistic/order','views/ecrm/statistic/order.vue','',1,1,1,1),
  (72,70,'statistic.product','商品统计','/statistic/product','views/ecrm/statistic/product.vue','',1,1,2,1),
  (15,10,'product.unit','商品单位','/product/unit','views/ecrm/product/unit.vue','',1,1,5,1),
  (16,10,'product.cdkey','卡密库','/product/cdkey','views/ecrm/product/cdkey.vue','',1,1,6,1),
  (17,10,'product.reviews','商品评论','/product/reviews','views/ecrm/product/reviews.vue','',1,1,7,1),
  (25,20,'order.customer','代客下单','/order/customer','views/ecrm/order/customer.vue','',1,1,5,1),
  (26,20,'order.invoice','发票管理','/order/invoice','views/ecrm/order/invoice.vue','',1,1,6,1),
  (44,40,'finance.statement','对账单','/accounts/statement','views/ecrm/finance/statement.vue','',1,1,4,1),
  (45,40,'finance.capitalFlow','资金流水','/accounts/capitalFlow','views/ecrm/finance/capital-flow.vue','',1,1,5,1),
  (46,40,'finance.transfer','转账管理','/accounts/transManagement','views/ecrm/finance/transfer.vue','',1,1,6,1),
  (47,40,'finance.profitsharing','分账申请','/systemForm/applyments','views/ecrm/finance/profitsharing.vue','',1,1,7,1),
  (35,30,'marketing.discounts','优惠套餐','/marketing/discounts/list','views/ecrm/marketing/discounts.vue','',1,1,5,1),
  (36,30,'marketing.topic','专题活动','/group/topic/95','views/ecrm/marketing/topic.vue','',1,1,6,1),
  (80,0,'user','用户','/user','','lucide:users',1,0,35,1),
  (81,80,'user.list','用户列表','/user/list','views/ecrm/user/list.vue','',1,1,1,1),
  (82,80,'user.label','用户标签','/user/label','views/ecrm/user/label.vue','',1,1,2,1),
  (83,80,'user.autoLabel','自动标签','/user/maticlabel','views/ecrm/user/auto-label.vue','',1,1,3,1),
  (84,80,'user.searchRecord','搜索记录','/user/searchRecord','views/ecrm/user/search-record.vue','',1,1,4,1),
  (58,50,'setting.serviceStaff','客服员工','/config/service_staff','views/ecrm/setting/service-staff.vue','',1,1,8,1),
  (59,50,'setting.printer','小票打印','/setting/printer/list','views/ecrm/setting/printer.vue','',1,1,9,1),
  (63,50,'setting.operationLog','操作日志','/setting/systemLog','views/ecrm/setting/operation-log.vue','',1,1,10,1),
  (64,50,'setting.autoReply','自动回复','/systemForm/customer_keyword','views/ecrm/setting/auto-reply.vue','',1,1,11,1),
  (65,50,'setting.openAuth','开放授权','/systemForm/openAuth/list','views/ecrm/setting/open-auth.vue','',1,1,12,1),
  (90,0,'delivery','同城配送','/delivery','','lucide:bike',1,0,65,1),
  (91,90,'delivery.personnel','配送员管理','/delivery/personnel_manage','views/ecrm/delivery/personnel.vue','',1,1,1,1),
  (92,90,'delivery.storePoint','自提点管理','/delivery/store_manage','views/ecrm/delivery/store-point.vue','',1,1,2,1),
  (100,0,'diy','装修','/devise','','lucide:palette',1,0,70,1),
  (101,100,'diy.systemForm','系统表单','/systemForm/form_list','views/ecrm/diy/system-form.vue','',1,1,3,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`path`=VALUES(`path`),`component`=VALUES(`component`),`icon`=VALUES(`icon`),`is_menu`=VALUES(`is_menu`),`is_route`=VALUES(`is_route`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);

-- 新店铺后台的写操作只使用 qixi_crm_m_ RBAC；禁止回退至旧 qixi_m_* 菜单表。
INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`) VALUES
  (1101,11,'product.create','新增商品','','','',2,0,1,1),
  (1102,11,'product.update','编辑商品','','','',2,0,2,1),
  (1103,11,'product.delete','移入回收站','','','',2,0,3,1),
  (1104,11,'product.restore','恢复商品','','','',2,0,4,1),
  (1105,11,'product.show','上下架商品','','','',2,0,5,1),
  (1106,11,'product.stock','调整库存','','','',2,0,6,1),
  (1201,21,'order.deliver','订单发货','','','',2,0,1,1),
  (1202,21,'order.proxy','代客下单','','','',2,0,2,1),
  (1203,23,'order.verify.action','确认核销','','','',2,0,1,1),
  (1301,22,'refund.approve','同意售后','','','',2,0,1,1),
  (1302,22,'refund.reject','驳回售后','','','',2,0,2,1),
  (1303,22,'refund.log','查看售后操作记录','','','',2,0,3,1),
  (1304,22,'refund.export','导出店铺退款清单','','','',2,0,4,1),
  (1305,22,'refund.express','查看退货物流快照','','','',2,0,5,1),
  (1306,22,'refund.remark','添加售后处理备注','','','',2,0,6,1),
  (1307,22,'refund.delete','从本店售后列表隐藏退款单','','','',2,0,7,1),
  (1401,52,'staff.create','新增员工','','','',2,0,1,1),
  (1402,52,'staff.update','编辑员工','','','',2,0,2,1),
  (1403,52,'staff.delete','移除员工','','','',2,0,3,1),
  (1410,11,'product.label.create','新增商品标签','','','',2,0,10,1),
  (1411,11,'product.label.update','编辑商品标签','','','',2,0,11,1),
  (1412,11,'product.label.delete','删除商品标签','','','',2,0,12,1),
  (1420,13,'product.category.create','新增商品分类','','','',2,0,1,1),
  (1421,13,'product.category.update','编辑商品分类','','','',2,0,2,1),
  (1422,13,'product.category.delete','删除商品分类','','','',2,0,3,1),
  (1501,43,'finance.settlement.apply','提交结算申请','','','',2,0,1,1),
  (1601,26,'invoice.audit','发票审核','','','',2,0,1,1),
  (1701,35,'marketing.discounts.manage','维护优惠套餐','','','',2,0,1,1)
ON DUPLICATE KEY UPDATE `parent_id`=VALUES(`parent_id`),`name`=VALUES(`name`),`is_menu`=VALUES(`is_menu`),`is_route`=VALUES(`is_route`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);
INSERT IGNORE INTO `qixi_crm_m_role_menu` (`role_code`,`menu_id`)
SELECT roles.role_code, menus.id FROM (SELECT 'owner' AS role_code UNION ALL SELECT 'manager') AS roles CROSS JOIN `qixi_crm_m_menu` AS menus;
