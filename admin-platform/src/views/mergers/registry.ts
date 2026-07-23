/**
 * 平台菜单 path → views/mergers 下组件（无后缀）
 * 未注册叶子走 placeholder
 */
const PATH_COMPONENT: Record<string, string> = {
  '/dashboard': 'mergers/dashboard/index',
  '/freight/express': 'mergers/freight/express',
  '/cms/article': 'mergers/cms/article',
  '/user/label': 'mergers/user/label',
  '/content/notice': 'mergers/content/notice',
  '/merchant/list': 'mergers/merchant/list',
  '/merchant/audit': 'mergers/merchant/audit',
  '/product/category': 'mergers/product/category',
  '/product/brand': 'mergers/product/brand',
  '/product/audit': 'mergers/product/audit',
  '/order/list': 'mergers/order/list',
  '/order/refund': 'mergers/order/refund',
  '/accounts/withdraw': 'mergers/accounts/withdraw',
  '/marketing/coupon': 'mergers/marketing/coupon',
  '/setting/admin': 'mergers/setting/admin',
  '/setting/role': 'mergers/setting/role',
  '/setting/menu': 'mergers/setting/menu',
  '/setting/agreements': 'mergers/setting/agreements',
  '/setting/diy/list': 'mergers/diy/list',
  '/setting/diy/index': 'mergers/diy/editor',
};

export function resolveMergersComponent(path: string): string | undefined {
  const normalized = path.startsWith('/') ? path : `/${path}`;
  return PATH_COMPONENT[normalized];
}
