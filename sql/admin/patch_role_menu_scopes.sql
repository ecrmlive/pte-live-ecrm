-- 商户、区域角色使用独立的可授权菜单树；所有菜单统一按 sort 升序展示。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_menu`
  (`id`, `parent_id`, `code`, `title`, `icon`, `route_path`, `kind`, `sort`, `status`, `menu_scope`)
VALUES
  -- 商户端
  (22001, 0,     'merchant.portal.home',               '首页',       'lucide:house',               '/dashboard',                        'directory', 1, 1, 'merchant'),
  (22002, 22001, 'merchant.portal.dashboard',          '控制台',     'lucide:layout-dashboard',    '/dashboard',                        'page',      1, 1, 'merchant'),
  (22010, 0,     'merchant.portal.product',            '商品',       'lucide:shopping-bag',         '/product',                          'directory', 2, 1, 'merchant'),
  (22011, 22010, 'merchant.portal.product.list',       '商品列表',   'lucide:list',                 '/product/list',                     'page',      1, 1, 'merchant'),
  (22012, 22010, 'merchant.portal.product.category',   '商品分类',   'lucide:folder-tree',          '/product/classify',                 'page',      2, 1, 'merchant'),
  (22013, 22010, 'merchant.portal.product.review',     '商品评价',   'lucide:message-square-text',  '/product/reviews',                  'page',      3, 1, 'merchant'),
  (22020, 0,     'merchant.portal.order',              '订单',       'lucide:clipboard-list',       '/order',                            'directory', 3, 1, 'merchant'),
  (22021, 22020, 'merchant.portal.order.list',         '订单列表',   'lucide:list',                 '/order/list',                       'page',      1, 1, 'merchant'),
  (22022, 22020, 'merchant.portal.order.refund',       '售后订单',   'lucide:rotate-ccw',           '/order/refund',                     'page',      2, 1, 'merchant'),
  (22023, 22020, 'merchant.portal.order.cancellation', '取消订单',   'lucide:ban',                  '/order/cancellation',               'page',      3, 1, 'merchant'),
  (22030, 0,     'merchant.portal.marketing',          '营销',       'lucide:flag',                 '/marketing',                        'directory', 4, 1, 'merchant'),
  (22031, 22030, 'merchant.portal.marketing.coupon',   '优惠券',     'lucide:ticket',               '/marketing/coupon',                 'page',      1, 1, 'merchant'),
  (22032, 22030, 'merchant.portal.marketing.seckill',  '秒杀活动',   'lucide:timer',                '/marketing/seckill/list',           'page',      2, 1, 'merchant'),
  (22033, 22030, 'merchant.portal.marketing.groupbuy', '拼团活动',   'lucide:users-round',          '/marketing/combination',            'page',      3, 1, 'merchant'),
  (22034, 22030, 'merchant.portal.marketing.presell',  '预售活动',   'lucide:calendar-clock',       '/marketing/presell/list',           'page',      4, 1, 'merchant'),
  (22040, 0,     'merchant.portal.finance',            '财务',       'lucide:chart-no-axes-column', '/accounts',                         'directory', 5, 1, 'merchant'),
  (22041, 22040, 'merchant.portal.finance.statement',  '账单明细',   'lucide:receipt-text',         '/accounts/statement',               'page',      1, 1, 'merchant'),
  (22042, 22040, 'merchant.portal.finance.balance',    '账户余额',   'lucide:wallet-cards',         '/accounts/balance',                 'page',      2, 1, 'merchant'),
  (22050, 0,     'merchant.portal.setting',            '设置',       'lucide:settings',             '/config',                           'directory', 6, 1, 'merchant'),
  (22051, 22050, 'merchant.portal.setting.store',      '店铺设置',   'lucide:store',                '/systemForm/modifyStoreInfo',       'page',      1, 1, 'merchant'),
  (22052, 22050, 'merchant.portal.setting.logistics',  '物流公司',   'lucide:truck',                '/config/freight/express',           'page',      2, 1, 'merchant'),
  (22053, 22050, 'merchant.portal.setting.freight',    '运费模板',   'lucide:map',                  '/config/freight/shippingTemplates', 'page',      3, 1, 'merchant'),
  (22054, 22050, 'merchant.portal.setting.staff',      '店员管理',   'lucide:headphones',           '/config/service',                   'page',      4, 1, 'merchant'),
  (22055, 22050, 'merchant.portal.setting.reply',      '自动回复',   'lucide:message-square-reply', '/systemForm/customer_keyword',      'page',      5, 1, 'merchant'),
  -- 区域端
  (22101, 0,     'region.portal.home',                 '首页',       'lucide:house',               '/dashboard',                        'directory', 1, 1, 'region'),
  (22102, 22101, 'region.portal.dashboard',            '控制台',     'lucide:layout-dashboard',    '/dashboard',                        'page',      1, 1, 'region'),
  (22110, 0,     'region.portal.management',           '区域管理',   'lucide:map-pinned',          '/business-zones',                   'directory', 2, 1, 'region'),
  (22111, 22110, 'region.portal.management.list',      '区域列表',   'lucide:list',                 '/business-zones/index',             'page',      1, 1, 'region'),
  (22112, 22110, 'region.portal.management.agents',    '代理人员',   'lucide:users-round',          '/business-zones/agents',            'page',      2, 1, 'region'),
  (22113, 22110, 'region.portal.management.setting',   '区域设置',   'lucide:settings',             '/business-zones/settings',          'page',      3, 1, 'region'),
  (22120, 0,     'region.portal.merchant',             '店铺监管',   'lucide:store',                '/merchant',                         'directory', 3, 1, 'region'),
  (22121, 22120, 'region.portal.merchant.list',        '店铺列表',   'lucide:list',                 '/merchant/list',                    'page',      1, 1, 'region'),
  (22122, 22120, 'region.portal.merchant.apply',       '入驻申请',   'lucide:file-check-2',         '/merchant/application',             'page',      2, 1, 'region'),
  (22130, 0,     'region.portal.order',                '订单监管',   'lucide:clipboard-list',       '/order',                            'directory', 4, 1, 'region'),
  (22131, 22130, 'region.portal.order.list',           '订单列表',   'lucide:list',                 '/order/list',                       'page',      1, 1, 'region'),
  (22132, 22130, 'region.portal.order.refund',         '售后订单',   'lucide:rotate-ccw',           '/order/refund',                     'page',      2, 1, 'region'),
  (22140, 0,     'region.portal.finance',              '财务监管',   'lucide:chart-no-axes-column', '/accounts',                         'directory', 5, 1, 'region'),
  (22141, 22140, 'region.portal.finance.settlement',   '商户结算',   'lucide:hand-coins',           '/accounts/merchant-settlement',     'page',      1, 1, 'region')
ON DUPLICATE KEY UPDATE
  `parent_id` = VALUES(`parent_id`),
  `title` = VALUES(`title`),
  `icon` = VALUES(`icon`),
  `route_path` = VALUES(`route_path`),
  `kind` = VALUES(`kind`),
  `sort` = VALUES(`sort`),
  `status` = VALUES(`status`),
  `menu_scope` = VALUES(`menu_scope`);

-- 默认商户、区域身份拥有本范围全部菜单；保留已有平台权限，避免影响旧账号的接口权限。
INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`, `menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` r
JOIN `qixi_crm_a_menu` m
  ON (r.code = 'merchant' AND m.menu_scope = 'merchant')
  OR (r.code = 'region' AND m.menu_scope = 'region');

-- 身份菜单树严格隔离。历史初始化曾把平台、商户、区域树同时授权给同一身份，
-- 会造成侧栏出现多个“首页/控制台”。角色类型只能拥有同类型菜单。
DELETE rm
FROM `qixi_crm_a_role_menu` AS rm
INNER JOIN `qixi_crm_a_role` AS r ON r.id = rm.role_id
INNER JOIN `qixi_crm_a_menu` AS m ON m.id = rm.menu_id
WHERE m.menu_scope <> r.role_type;
