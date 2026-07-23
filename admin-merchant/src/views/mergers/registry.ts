const PATH_COMPONENT: Record<string, string> = {
  '/dashboard': 'mergers/dashboard/index',
  '/config/freight/shippingTemplates': 'mergers/shipping/templates',
  '/product/label': 'mergers/product/label',
  '/product/specs': 'mergers/product/specs',
  '/config/guarantee': 'mergers/product/guarantee',
  '/order/invoice': 'mergers/order/invoice',
  '/order/list': 'mergers/order/list',
  '/product/list': 'mergers/product/list',
  '/finance/balance': 'mergers/finance/balance',
  '/finance/withdraw': 'mergers/finance/withdraw',
  '/marketing/coupon': 'mergers/marketing/coupon',
  '/setting/shop': 'mergers/setting/shop',
  '/devise/diy/list': 'mergers/diy/list',
  '/devise/diy/index': 'mergers/diy/editor',
};

export function resolveMergersComponent(path: string): string | undefined {
  const normalized = path.startsWith('/') ? path : `/${path}`;
  return PATH_COMPONENT[normalized];
}
