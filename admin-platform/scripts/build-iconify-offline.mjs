/**
 * 从 @iconify/json 提取平台管理端用到的图标，写入 public/iconify/（离线加载，不请求 api.iconify.design）
 *
 * LUCIDE_ICONS 须与 sql/init_platform_access.sql 中 lucide:* 一致（去前缀）。
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

const outDir = path.resolve(import.meta.dirname, '../public/iconify');
fs.mkdirSync(outDir, { recursive: true });

const EP_ICONS = ['fold', 'expand'];
/** @see sql/init_platform_access.sql */
const LUCIDE_ICONS = [
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

const missing = LUCIDE_ICONS.filter((name) => !lucide.icons[name]);
if (missing.length) {
  console.warn('[build:icons] missing in @iconify/json lucide:', missing.join(', '));
}

fs.writeFileSync(path.join(outDir, 'ep.json'), JSON.stringify(epSubset));
fs.writeFileSync(path.join(outDir, 'lucide.json'), JSON.stringify(lucideSubset));

console.log('[build:icons] ep:', Object.keys(epSubset.icons).join(', '));
console.log('[build:icons] lucide:', Object.keys(lucideSubset.icons).join(', '));
