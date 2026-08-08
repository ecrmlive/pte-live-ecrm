/**
 * 平台侧栏可用 Lucide 图标（离线 bundle，见 scripts/build-iconify-offline.mjs）。
 * 店铺菜单顶层为 ant-design:*，须同步打入 public/iconify/ant-design.json。
 * 新增菜单 icon 时同步改 SQL 种子与本列表，并运行 pnpm --filter @pte-live/admin-platform run build:icons
 */
export const PLATFORM_LUCIDE_ICON_NAMES = [
  'activity',
  'award',
  'badge',
  'badge-check',
  'badge-dollar-sign',
  'badge-percent',
  'badge-plus',
  'ban',
  'bar-chart-3',
  'bell',
  'boxes',
  'chart-area',
  'circle-dollar-sign',
  'circle-dot',
  'clipboard-list',
  'clipboard-pen',
  'contact',
  'contact-round',
  'credit-card',
  'database-backup',
  'eraser',
  'file-text',
  'folder-tree',
  'frame',
  'git-branch',
  'git-branch-plus',
  'hand-heart',
  'handshake',
  'hard-drive',
  'hash',
  'headphones',
  'headset',
  'heart-handshake',
  'history',
  'house',
  'images',
  'key-round',
  'landmark',
  'layers',
  'layout-dashboard',
  'layout-grid',
  'layout-template',
  'life-buoy',
  'list-tree',
  'log-in',
  'mail',
  'map-pinned',
  'map-plus',
  'medal',
  'megaphone',
  'menu',
  'message-circle',
  'message-circle-warning',
  'message-square',
  'message-square-reply',
  'messages-square',
  'monitor',
  'newspaper',
  'package',
  'paintbrush',
  'palette',
  'pen-line',
  'puzzle',
  'radio-tower',
  'receipt',
  'receipt-text',
  'scroll-text',
  'search',
  'send',
  'settings',
  'settings-2',
  'share-2',
  'shield',
  'shield-check',
  'shopping-bag',
  'smartphone',
  'sparkles',
  'store',
  'tags',
  'ticket',
  'ticket-check',
  'ticket-plus',
  'truck',
  'user-cog',
  'user-plus',
  'user-round',
  'user-round-cog',
  'users',
  'users-round',
  'wallet',
  'wallet-cards',
  'wrench',
] as const;

export const PLATFORM_LUCIDE_ICONS = PLATFORM_LUCIDE_ICON_NAMES.map(
  (name) => `lucide:${name}`,
);

/** 平台侧栏顶层 / 关键二级 ant-design 图标（离线 ant-design.json） */
export const PLATFORM_ANT_DESIGN_ICON_NAMES = [
  'appstore-outlined',
  'bar-chart-outlined',
  'cluster-outlined',
  'customer-service-outlined',
  'dashboard-outlined',
  'file-text-outlined',
  'flag-outlined',
  'format-painter-outlined',
  'home-outlined',
  'read-outlined',
  'send-outlined',
  'setting-outlined',
  'share-alt-outlined',
  'shop-outlined',
  'shopping-outlined',
  'team-outlined',
  'tool-outlined',
  'user-outlined',
] as const;

/** 菜单图标选择器：lucide + 平台顶层 ant-design */
export const PLATFORM_MENU_PICKER_ICONS = [
  ...PLATFORM_ANT_DESIGN_ICON_NAMES.map((name) => `ant-design:${name}`),
  ...PLATFORM_LUCIDE_ICONS,
] as const;

/**
 * CRMEB 商户菜单旧 icon（无 iview/remix 名或历史 lucide）→ Vben Iconify。
 * 店铺菜单列表「菜单图标」列与侧栏预览共用；已是 `prefix:name` 且未命中表则原样返回。
 */
export const CRMEB_STORE_MENU_ICON_MAP: Record<string, string> = {
  // CRMEB bare / Element 风格
  house: 'ant-design:home-outlined',
  's-home': 'ant-design:home-outlined',
  goods: 'ant-design:shopping-outlined',
  's-goods': 'ant-design:shopping-outlined',
  tickets: 'ant-design:file-text-outlined',
  's-order': 'ant-design:file-text-outlined',
  bell: 'ant-design:bell-outlined',
  's-flag': 'ant-design:bell-outlined',
  'pie-chart': 'ant-design:field-time-outlined',
  's-data': 'ant-design:field-time-outlined',
  user: 'ant-design:user-outlined',
  'user-solid': 'ant-design:user-outlined',
  headset: 'ant-design:customer-service-outlined',
  brush: 'ant-design:format-painter-outlined',
  's-open': 'ant-design:format-painter-outlined',
  setting: 'ant-design:setting-outlined',
  's-tools': 'ant-design:setting-outlined',
  'notebook-2': 'ant-design:notification-outlined',
  notebook: 'ant-design:notification-outlined',
  // 历史种子 lucide → 与 CRMEB 顶层观感对齐的 ant-design
  'lucide:home': 'ant-design:home-outlined',
  'lucide:house': 'ant-design:home-outlined',
  'lucide:package': 'ant-design:shopping-outlined',
  'lucide:shopping-bag': 'ant-design:shopping-outlined',
  'lucide:tickets': 'ant-design:file-text-outlined',
  'lucide:file-text': 'ant-design:file-text-outlined',
  'lucide:bell': 'ant-design:bell-outlined',
  'lucide:pie-chart': 'ant-design:field-time-outlined',
  'lucide:users': 'ant-design:user-outlined',
  'lucide:user-round': 'ant-design:user-outlined',
  'lucide:headset': 'ant-design:customer-service-outlined',
  'lucide:headphones': 'ant-design:customer-service-outlined',
  'lucide:paintbrush': 'ant-design:format-painter-outlined',
  'lucide:shield': 'ant-design:safety-outlined',
  'lucide:shield-check': 'ant-design:safety-outlined',
  'lucide:settings': 'ant-design:setting-outlined',
  'lucide:notebook': 'ant-design:notification-outlined',
  'lucide:megaphone': 'ant-design:notification-outlined',
};

/**
 * 店铺菜单名称 → Iconify 兜底（DB icon 为空时列表仍可渲染）。
 * 仅覆盖目录/页面常用名；纯按钮权限叶子不在表内则仍空。
 */
export const STORE_MENU_NAME_ICON_MAP: Record<string, string> = {
  首页: 'ant-design:home-outlined',
  控制台: 'ant-design:dashboard-outlined',
  商品统计: 'ant-design:bar-chart-outlined',
  订单统计: 'ant-design:line-chart-outlined',
  权限: 'ant-design:safety-outlined',
  权限管理: 'ant-design:safety-outlined',
  商品: 'ant-design:shopping-outlined',
  商品列表: 'ant-design:shopping-outlined',
  商品分类: 'ant-design:appstore-outlined',
  商品规格: 'ant-design:cluster-outlined',
  商品评价: 'ant-design:star-outlined',
  订单: 'ant-design:file-text-outlined',
  订单管理: 'ant-design:file-text-outlined',
  退款订单: 'ant-design:rollback-outlined',
  营销: 'ant-design:bell-outlined',
  财务: 'ant-design:field-time-outlined',
  资金流水: 'ant-design:transaction-outlined',
  账单管理: 'ant-design:account-book-outlined',
  申请分账商户: 'ant-design:account-book-outlined',
  用户: 'ant-design:user-outlined',
  用户管理: 'ant-design:user-outlined',
  员工: 'ant-design:customer-service-outlined',
  装修: 'ant-design:format-painter-outlined',
  设置: 'ant-design:setting-outlined',
  公告: 'ant-design:notification-outlined',
  服务统计: 'ant-design:bar-chart-outlined',
  配送统计: 'ant-design:bar-chart-outlined',
};

/** 侧栏/列表展示用：CRMEB 映射或补全 lucide: 前缀，空值返回 undefined */
export function normalizePlatformMenuIcon(icon?: string) {
  const raw = String(icon || '').trim();
  if (!raw) {
    return undefined;
  }
  const mapped = CRMEB_STORE_MENU_ICON_MAP[raw] || CRMEB_STORE_MENU_ICON_MAP[raw.toLowerCase()];
  if (mapped) {
    return mapped;
  }
  if (raw.includes(':')) {
    return raw;
  }
  return `lucide:${raw.replace(/^lucide-/, '')}`;
}

/**
 * 店铺菜单列表图标：优先 DB icon，空则按菜单名兜底。
 * 纯按钮（is_menu=2 且无 path）且名称不在兜底表时保持空。
 */
export function resolveStoreMenuIcon(opts: {
  icon?: string;
  name?: string;
  isMenu?: number;
  path?: string;
}) {
  const fromDb = normalizePlatformMenuIcon(opts.icon);
  if (fromDb) {
    return fromDb;
  }
  const name = String(opts.name || '').trim();
  const byName = name ? STORE_MENU_NAME_ICON_MAP[name] : undefined;
  if (byName) {
    return byName;
  }
  const isButton = Number(opts.isMenu) === 2;
  const hasPath = Boolean(String(opts.path || '').trim());
  if (isButton && !hasPath) {
    return undefined;
  }
  return undefined;
}
