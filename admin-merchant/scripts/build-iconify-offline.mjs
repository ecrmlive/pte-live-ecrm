/**
 * 从 @iconify/json 提取商户管理端侧栏图标，写入 public/iconify/
 *
 * LUCIDE_ICONS 须与 sql/init_merchant_access.sql（及 gen_merchant_access.py ICON_BY_*）一致。
 * ANT_DESIGN_ICONS 须与店铺菜单 qixi_crm_m_menu 顶层 ant-design:* 一致。
 * 插件中心卡片图标见 views/native/plus/plus/index.vue PLUGIN_PATH_ICONS。
 * 新增菜单 icon 时更新列表后运行：
 *   pnpm --filter @pte-live/admin-merchant run build:icons
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
/** 店铺菜单 Iconify（qixi_crm_m_menu 全树 ant-design:*，与 platform 对齐） */
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
  'rollback-outlined',
  'safety-outlined',
  'schedule-outlined',
  'search-outlined',
  'setting-outlined',
  'shop-outlined',
  'shopping-outlined',
  'star-outlined',
  'swap-outlined',
  'tags-outlined',
  'team-outlined',
  'thunderbolt-outlined',
  'transaction-outlined',
  'trophy-outlined',
  'user-outlined',
  'video-camera-outlined',
  'wallet-outlined',
];
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
const antDesignSubset = subset(antDesign, ANT_DESIGN_ICONS);

fs.writeFileSync(path.join(outDir, 'ep.json'), JSON.stringify(epSubset));
fs.writeFileSync(path.join(outDir, 'lucide.json'), JSON.stringify(lucideSubset));
fs.writeFileSync(path.join(outDir, 'ant-design.json'), JSON.stringify(antDesignSubset));

console.log('[build:icons] ep:', Object.keys(epSubset.icons).join(', '));
console.log('[build:icons] lucide:', Object.keys(lucideSubset.icons).join(', '));
console.log('[build:icons] ant-design:', Object.keys(antDesignSubset.icons).join(', '));
