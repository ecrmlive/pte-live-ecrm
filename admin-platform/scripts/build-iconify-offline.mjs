/**
 * 从 @iconify/json 提取平台管理端用到的图标，写入 public/iconify/（离线加载，不请求 api.iconify.design）
 *
 * LUCIDE_ICONS 须与 sql/init_platform_access.sql 中 lucide:* 一致（去前缀）。
 * ANT_DESIGN_ICONS 须与店铺菜单种子 / CRMEB_STORE_MENU_ICON_MAP 中 ant-design:* 一致。
 * 新增菜单 icon 时：改 SQL 种子 → 本列表 → 运行：
 *   pnpm --filter @pte-live/admin-platform run build:icons
 */
import fs from 'node:fs';
import path from 'node:path';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const jsonRoot = path.dirname(require.resolve('@iconify/json/json/ep.json'));
const ep = require(path.join(jsonRoot, 'ep.json'));
const lucide = require(path.join(jsonRoot, 'lucide.json'));
const antDesign = require(path.join(jsonRoot, 'ant-design.json'));

const outDir = path.resolve(import.meta.dirname, '../public/iconify');
fs.mkdirSync(outDir, { recursive: true });

const EP_ICONS = ['fold', 'expand'];
/** 店铺菜单 Iconify（qixi_crm_m_menu 全树 ant-design:* / CRMEB_STORE_MENU_ICON_MAP） */
const ANT_DESIGN_ICONS = [
  'account-book-outlined',
  'appstore-outlined',
  'audit-outlined',
  'bar-chart-outlined',
  'barcode-outlined',
  'bell-outlined',
  'calendar-outlined',
  'car-outlined',
  'cloud-outlined',
  'cluster-outlined',
  'comment-outlined',
  'credit-card-outlined',
  'crown-outlined',
  'customer-service-outlined',
  'dashboard-outlined',
  'field-time-outlined',
  'file-done-outlined',
  'file-protect-outlined',
  'file-text-outlined',
  'flag-outlined',
  'font-size-outlined',
  'form-outlined',
  'format-painter-outlined',
  'gift-outlined',
  'gold-outlined',
  'history-outlined',
  'home-outlined',
  'idcard-outlined',
  'key-outlined',
  'like-outlined',
  'line-chart-outlined',
  'message-outlined',
  'notification-outlined',
  'percentage-outlined',
  'picture-outlined',
  'printer-outlined',
  'profile-outlined',
  'read-outlined',
  'rollback-outlined',
  'safety-outlined',
  'schedule-outlined',
  'search-outlined',
  'send-outlined',
  'setting-outlined',
  'share-alt-outlined',
  'shop-outlined',
  'shopping-outlined',
  'star-outlined',
  'swap-outlined',
  'tags-outlined',
  'team-outlined',
  'thunderbolt-outlined',
  'tool-outlined',
  'transaction-outlined',
  'trophy-outlined',
  'user-outlined',
  'video-camera-outlined',
  'wallet-outlined',
];
/** @see sql/init_platform_access.sql */
const LUCIDE_ICONS = [
  'activity',
  'align-left',
  'arrow-left-right',
  'award',
  'badge',
  'badge-check',
  'badge-dollar-sign',
  'badge-percent',
  'badge-plus',
  'ban',
  'bar-chart-3',
  'bell',
  'book-open',
  'boxes',
  'calendar-clock',
  'calendar-check-2',
  'chart-area',
  'chart-no-axes-column',
  'chevron-right',
  'circle-dollar-sign',
  'circle-dot',
  'circle-user-round',
  'clipboard-list',
  'clipboard-pen',
  'cloud',
  'code',
  'coins',
  'contact',
  'contact-round',
  'credit-card',
  'crown',
  'database',
  'database-backup',
  'download',
  'eraser',
  'file',
  'file-spreadsheet',
  'file-check-2',
  'file-text',
  'flame',
  'flag',
  'folder-tree',
  'frame',
  'gift',
  'git-branch',
  'git-branch-plus',
  'grip',
  'grip-vertical',
  'hand-heart',
  'hand-coins',
  'handshake',
  'hard-drive',
  'hash',
  'headphones',
  'headset',
  'heart',
  'heart-handshake',
  'history',
  'house',
  'image',
  'images',
  'key-round',
  'landmark',
  'layers',
  'layout-dashboard',
  'layout-grid',
  'layout-template',
  'life-buoy',
  'link',
  'list',
  'list-tree',
  'log-in',
  'mail',
  'map',
  'map-pin',
  'map-pinned',
  'map-plus',
  'medal',
  'megaphone',
  'menu',
  'message-circle',
  'message-circle-more',
  'message-circle-warning',
  'message-square',
  'message-square-reply',
  'message-square-text',
  'messages-square',
  'monitor',
  'mouse-pointer',
  'newspaper',
  'package',
  'package-check',
  'panel-bottom',
  'party-popper',
  'paintbrush',
  'palette',
  'pen-line',
  'percent',
  'puzzle',
  'radio',
  'radio-tower',
  'receipt',
  'receipt-text',
  'rotate-ccw',
  'scale',
  'scan-line',
  'scroll-text',
  'search',
  'send',
  'settings',
  'settings-2',
  'share-2',
  'shield',
  'shield-check',
  'shopping-bag',
  'shopping-cart',
  'sliders-horizontal',
  'smartphone',
  'sparkles',
  'split',
  'store',
  'tags',
  'ticket',
  'ticket-check',
  'ticket-plus',
  'tickets',
  'timer',
  'truck',
  'user-cog',
  'user-plus',
  'user-round',
  'user-round-cog',
  'users',
  'users-round',
  'video',
  'wallet',
  'wallet-cards',
  'wrench',
  'x',
  'zap',
];

function subset(collection, names) {
  return {
    prefix: collection.prefix,
    icons: Object.fromEntries(
      names.filter((name) => collection.icons[name]).map((name) => [name, collection.icons[name]]),
    ),
    width: collection.width,
    height: collection.height,
  };
}

const epSubset = subset(ep, EP_ICONS);
const lucideSubset = subset(lucide, LUCIDE_ICONS);
const antDesignSubset = subset(antDesign, ANT_DESIGN_ICONS);

const missingLucide = LUCIDE_ICONS.filter((name) => !lucide.icons[name]);
if (missingLucide.length) {
  console.warn('[build:icons] missing in @iconify/json lucide:', missingLucide.join(', '));
}
const missingAnt = ANT_DESIGN_ICONS.filter((name) => !antDesign.icons[name]);
if (missingAnt.length) {
  console.warn('[build:icons] missing in @iconify/json ant-design:', missingAnt.join(', '));
}

fs.writeFileSync(path.join(outDir, 'ep.json'), JSON.stringify(epSubset));
fs.writeFileSync(path.join(outDir, 'lucide.json'), JSON.stringify(lucideSubset));
fs.writeFileSync(path.join(outDir, 'ant-design.json'), JSON.stringify(antDesignSubset));

console.log('[build:icons] ep:', Object.keys(epSubset.icons).join(', '));
console.log('[build:icons] lucide:', Object.keys(lucideSubset.icons).join(', '));
console.log('[build:icons] ant-design:', Object.keys(antDesignSubset.icons).join(', '));
