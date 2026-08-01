/**
 * 平台菜单 path → views/ecrm 下组件（无后缀）
 * 未注册叶子走 placeholder
 */
const PATH_COMPONENT: Record<string, string> = {
  '/dashboard': 'ecrm/dashboard/index',
  '/freight/express': 'ecrm/freight/express',
  '/cms/article': 'ecrm/cms/article',
  '/user/label': 'ecrm/user/label',
  '/user/group': 'ecrm/user/group',
  '/user/svip': 'ecrm/user/svip',
  '/content/notice': 'ecrm/content/notice',
  '/content/community': 'ecrm/content/community',
  '/content/attachment': 'ecrm/content/attachment',
  '/config/picture': 'ecrm/content/attachment',
  '/merchant/list': 'ecrm/merchant/list',
  '/merchant/audit': 'ecrm/merchant/audit',
  '/merchant/review': 'ecrm/merchant/audit',
  '/business-zones/index': 'ecrm/business-zones/index',
  '/business-zones/agents': 'ecrm/business-zones/agents',
  '/business-zones/agent-review': 'ecrm/business-zones/agent-review',
  '/region': 'ecrm/business-zones/index',
  '/service': 'ecrm/customer-service/index',
  '/product/category': 'ecrm/product/category',
  '/product/brand': 'ecrm/product/brand',
  '/product/audit': 'ecrm/product/audit',
  '/order/list': 'ecrm/order/list',
  '/order/refund': 'ecrm/order/refund',
  '/accounts/withdraw': 'ecrm/accounts/withdraw',
  '/marketing/coupon': 'ecrm/marketing/coupon',
  '/marketing/spread': 'ecrm/marketing/spread',
  '/marketing/seckill': 'ecrm/marketing/seckill',
  '/marketing/combination': 'ecrm/marketing/combination',
  '/marketing/presell': 'ecrm/marketing/presell',
  '/marketing/broadcast': 'ecrm/marketing/broadcast',
  '/marketing/coupon/list': 'ecrm/marketing/coupon',
  '/marketing/platform_coupon': 'ecrm/marketing/coupon',
  '/marketing/platform_coupon/list': 'ecrm/marketing/coupon',
  '/setting/admin': 'ecrm/setting/admin',
  '/setting/role': 'ecrm/setting/role',
  '/setting/menu': 'ecrm/setting/menu',
  '/setting/agreements': 'ecrm/setting/agreements',
  '/setting/cloud-config': 'ecrm/setting/cloud-config',
  '/setting/sms': 'ecrm/setting/sms',
  '/setting/diy/list': 'ecrm/diy/list',
  '/setting/micro/list': 'ecrm/diy/list',
  '/setting/diy/index': 'ecrm/diy/editor',
  '/content/diy': 'ecrm/diy/list',
  '/operations/diy': 'ecrm/diy/list',
  '/setting/diy/plantform/category/list': 'ecrm/diy/page-category',
  '/setting/diy/merchant/category/list': 'ecrm/diy/page-category',
  '/setting/diy/links/list': 'ecrm/diy/page-link',
  '/setting/diy/merLink/list': 'ecrm/diy/page-link',
};

export function resolveMergersComponent(path: string): string | undefined {
  const normalized = path.startsWith('/') ? path : `/${path}`;
  return PATH_COMPONENT[normalized];
}
