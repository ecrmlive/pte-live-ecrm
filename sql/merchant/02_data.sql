USE `qixi_crm_merchant`;
INSERT INTO `qixi_crm_m_merchant` (`id`,`name`,`status`) VALUES (1,'七禧演示商户',1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`status`=VALUES(`status`);
INSERT INTO `qixi_crm_m_store` (`id`,`merchant_id`,`app_id`,`name`,`status`) VALUES (1,1,'qixi.store.demo.1','七禧演示店铺',1)
ON DUPLICATE KEY UPDATE `merchant_id`=VALUES(`merchant_id`),`app_id`=VALUES(`app_id`),`name`=VALUES(`name`),`status`=VALUES(`status`);

INSERT INTO `qixi_crm_m_menu` (`id`,`parent_id`,`code`,`name`,`path`,`component`,`icon`,`is_menu`,`is_route`,`sort`,`status`) VALUES
  (1,0,'dashboard','经营概览','/dashboard','views/mergers/dashboard/index.vue','lucide:layout-dashboard',1,1,1,1),
  (10,0,'product','商品','/product','', 'lucide:package',1,0,10,1),
  (11,10,'product.list','商品管理','/product/list','views/mergers/product/list.vue','',1,1,1,1),
  (12,10,'product.specs','规格模板','/product/specs','views/mergers/product/specs.vue','',1,1,2,1),
  (20,0,'order','订单','/order','', 'lucide:receipt-text',1,0,20,1),
  (21,20,'order.list','订单管理','/order/list','views/mergers/order/list.vue','',1,1,1,1),
  (22,20,'order.refund','售后管理','/order/refund','views/mergers/order/refund.vue','',1,1,2,1),
  (30,0,'marketing','营销','/marketing','', 'lucide:megaphone',1,0,30,1),
  (31,30,'marketing.coupon','优惠券','/marketing/coupon','views/mergers/marketing/coupon.vue','',1,1,1,1),
  (32,30,'marketing.seckill','秒杀','/marketing/seckill/list','views/mergers/marketing/seckill.vue','',1,1,2,1),
  (40,0,'finance','财务','/finance','', 'lucide:wallet-cards',1,0,40,1),
  (41,40,'finance.balance','资金账户','/finance/balance','views/mergers/finance/balance.vue','',1,1,1,1),
  (42,40,'finance.withdraw','提现管理','/finance/withdraw','views/mergers/finance/withdraw.vue','',1,1,2,1),
  (50,0,'setting','设置','/setting','', 'lucide:settings',1,0,50,1),
  (51,50,'setting.shop','店铺设置','/setting/shop','views/mergers/setting/shop.vue','',1,1,1,1),
  (52,50,'setting.staff','员工管理','/setting/staff','views/mergers/setting/staff.vue','',1,1,2,1),
  (53,50,'setting.admins','账号管理','/setting/admins','views/mergers/setting/admins.vue','',1,1,3,1),
  (54,50,'setting.role','角色权限','/setting/role','views/mergers/setting/role.vue','',1,1,4,1),
  (55,50,'setting.payment','支付方式','/setting/payment','views/mergers/setting/payment.vue','',1,1,5,1),
  (56,50,'setting.im_sdk_app','IM SDK AppId','/setting/im-sdk-app','views/mergers/setting/im-sdk-app.vue','',1,1,6,1)
ON DUPLICATE KEY UPDATE `name`=VALUES(`name`),`path`=VALUES(`path`),`component`=VALUES(`component`),`icon`=VALUES(`icon`),`is_menu`=VALUES(`is_menu`),`is_route`=VALUES(`is_route`),`sort`=VALUES(`sort`),`status`=VALUES(`status`);
INSERT IGNORE INTO `qixi_crm_m_role_menu` (`role_code`,`menu_id`)
SELECT roles.role_code, menus.id FROM (SELECT 'owner' AS role_code UNION ALL SELECT 'manager') AS roles CROSS JOIN `qixi_crm_m_menu` AS menus;
