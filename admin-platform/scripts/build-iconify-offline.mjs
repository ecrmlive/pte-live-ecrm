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
  'badge-check',
  'chart-area',
  'circle-dot',
  'folder-tree',
  'git-branch',
  'house',
  'images',
  'key-round',
  'layout-grid',
  'log-in',
  'map-pinned',
  'map-plus',
  'messages-square',
  'pen-line',
  'puzzle',
  'radio-tower',
  'receipt-text',
  'settings',
  'shield-check',
  'store',
  'user-round-cog',
  'users',
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

const missing = LUCIDE_ICONS.filter((name) => !lucide.icons[name]);
if (missing.length) {
  console.warn('[build:icons] missing in @iconify/json lucide:', missing.join(', '));
}

fs.writeFileSync(path.join(outDir, 'ep.json'), JSON.stringify(epSubset));
fs.writeFileSync(path.join(outDir, 'lucide.json'), JSON.stringify(lucideSubset));

console.log('[build:icons] ep:', Object.keys(epSubset.icons).join(', '));
console.log('[build:icons] lucide:', Object.keys(lucideSubset.icons).join(', '));
