/**
 * 从 @iconify/json 提取商户管理端侧栏 lucide 图标，写入 public/iconify/
 *
 * LUCIDE_ICONS 须与 sql/init_merchant_access.sql（及 gen_merchant_access.py ICON_BY_*）一致。
 * 插件中心卡片图标见 views/native/plus/plus/index.vue PLUGIN_PATH_ICONS。
 * 新增菜单 icon 时更新 LUCIDE_ICONS 后运行：
 *   pnpm --filter @pte-live/admin-merchant run build:icons
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
const LUCIDE_ICONS = [
  'app-window',
  'award',
  'banknote',
  'bell',
  'bot',
  'chart-bar',
  'chart-line',
  'circle-user',
  'clipboard-list',
  'coins',
  'contact-round',
  'credit-card',
  'eraser',
  'file-check',
  'file-clock',
  'file-text',
  'files',
  'film',
  'flag',
  'folder-tree',
  'gauge',
  'gift',
  'globe',
  'headphones',
  'headset',
  'history',
  'house',
  'layers',
  'layout-grid',
  'list',
  'link',
  'list-ordered',
  'map-pin',
  'map-pinned',
  'message-circle',
  'message-square',
  'package',
  'package-search',
  'paintbrush',
  'palette',
  'panel-bottom',
  'printer',
  'puzzle',
  'radio',
  'receipt',
  'scan-line',
  'scroll-text',
  'settings',
  'shield',
  'shield-ban',
  'shopping-cart',
  'smartphone',
  'star',
  'store',
  'tags',
  'trending-up',
  'truck',
  'user-check',
  'user-cog', // plus/agent、auth/role
  'user-round',
  'users',
  'users-round',
  'video',
  'wallet',
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

fs.writeFileSync(path.join(outDir, 'ep.json'), JSON.stringify(epSubset));
fs.writeFileSync(path.join(outDir, 'lucide.json'), JSON.stringify(lucideSubset));

console.log('[build:icons] ep:', Object.keys(epSubset.icons).join(', '));
console.log('[build:icons] lucide:', Object.keys(lucideSubset.icons).join(', '));
