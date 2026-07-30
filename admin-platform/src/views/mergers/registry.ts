/**
 * 平台菜单 path → views/mergers 下组件（无后缀）
 * 未注册叶子走 placeholder
 */
const PATH_COMPONENT: Record<string, string> = {
  '/dashboard': 'mergers/dashboard/index',
  '/freight/express': 'mergers/freight/express',
  '/cms/article': 'mergers/cms/article',
  '/user/label': 'mergers/user/label',
  '/user/group': 'mergers/user/group',
  '/user/svip': 'mergers/user/svip',
  '/content/notice': 'mergers/content/notice',
  '/content/community': 'mergers/content/community',
  '/content/attachment': 'mergers/content/attachment',
  '/config/picture': 'mergers/content/attachment',
  '/merchant/list': 'mergers/merchant/list',
  '/merchant/audit': 'mergers/merchant/audit',
  '/merchant/review': 'mergers/merchant/audit',
  '/business-zones/index': 'mergers/business-zones/index',
  '/business-zones/agents': 'mergers/business-zones/agents',
  '/business-zones/agent-review': 'mergers/business-zones/agent-review',
  '/region': 'mergers/business-zones/index',
  '/service': 'mergers/customer-service/index',
  '/product/category': 'mergers/product/category',
  '/product/brand': 'mergers/product/brand',
  '/product/audit': 'mergers/product/audit',
  '/order/list': 'mergers/order/list',
  '/order/refund': 'mergers/order/refund',
  '/accounts/withdraw': 'mergers/accounts/withdraw',
  '/marketing/coupon': 'mergers/marketing/coupon',
  '/marketing/spread': 'mergers/marketing/spread',
  '/marketing/seckill': 'mergers/marketing/seckill',
  '/marketing/combination': 'mergers/marketing/combination',
  '/marketing/presell': 'mergers/marketing/presell',
  '/marketing/broadcast': 'mergers/marketing/broadcast',
  '/marketing/coupon/list': 'mergers/marketing/coupon',
  '/marketing/platform_coupon': 'mergers/marketing/coupon',
  '/marketing/platform_coupon/list': 'mergers/marketing/coupon',
  '/setting/admin': 'mergers/setting/admin',
  '/setting/role': 'mergers/setting/role',
  '/setting/menu': 'mergers/setting/menu',
  '/setting/agreements': 'mergers/setting/agreements',
  '/setting/cloud-config': 'mergers/setting/cloud-config',
  '/setting/sms': 'mergers/setting/sms',
  '/setting/diy/list': 'mergers/diy/list',
  '/setting/micro/list': 'mergers/diy/list',
  '/setting/diy/index': 'mergers/diy/editor',
  '/content/diy': 'mergers/diy/list',
  '/operations/diy': 'mergers/diy/list',
  '/setting/diy/plantform/category/list': 'mergers/diy/page-category',
  '/setting/diy/merchant/category/list': 'mergers/diy/page-category',
  '/setting/diy/links/list': 'mergers/diy/page-link',
  '/setting/diy/merLink/list': 'mergers/diy/page-link',
};

export function resolveMergersComponent(path: string): string | undefined {
  const normalized = path.startsWith('/') ? path : `/${path}`;
  return PATH_COMPONENT[normalized];
}
