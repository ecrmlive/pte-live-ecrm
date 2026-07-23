#!/usr/bin/env node
/**
 * Align merchant list pages with platform-admin layout:
 * - height: 'auto' (not 100%)
 * - list-panel without flex-1
 * - native-vxe-grid without flex-1 min-h-0
 * - Page content-class without overflow-hidden min-h-0 flex chain
 * - Remove scoped vxe height: 100% !important blocks
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '../src/views/native');

const SKIP_PATTERNS = [
  /\/components\//,
  /\/dialog\//,
  /\/control\//,
  /ControlPreview\.vue$/,
  /-field\.vue$/,
  /-sku-field\.vue$/,
  /fullreduce-rules-field\.vue$/,
  /agent-apply-product-grid\.vue$/,
  /recommend-product-grid\.vue$/,
  /package-coupon-grid\.vue$/,
  /package-product-grid\.vue$/,
  /invitation-coupon-grid\.vue$/,
  /sign-coupon-grid\.vue$/,
  /register-coupon-grid\.vue$/,
  /coupon-product-grid\.vue$/,
  /coupon-send-user-grid\.vue$/,
  /newactivity-product-grid\.vue$/,
  /lottery-prize-grid\.vue$/,
  /delivery-rule-grid\.vue$/,
  /seckill-active-product-field\.vue$/,
  /control-moderation-grid\.vue$/,
  /control-products-grid\.vue$/,
  /RoomProductDialog\.vue$/,
  /PickMallProductsDialog\.vue$/,
  /order-detail-content\.vue$/,
  /refund-detail-content\.vue$/,
  /product-add-modal\.vue$/,
  /statistics\.vue$/,
  /fullreduce\/product\.vue$/,
];

function walk(dir, out = []) {
  for (const name of fs.readdirSync(dir)) {
    const full = path.join(dir, name);
    const stat = fs.statSync(full);
    if (stat.isDirectory()) walk(full, out);
    else if (name.endsWith('.vue')) out.push(full);
  }
  return out;
}

function shouldSkip(file) {
  const rel = path.relative(ROOT, file);
  return SKIP_PATTERNS.some((re) => re.test(rel));
}

function stripFlexClasses(cls) {
  return cls
    .split(/\s+/)
    .filter(
      (c) =>
        c &&
        c !== 'flex-1' &&
        c !== 'min-h-0' &&
        c !== 'flex' &&
        c !== 'overflow-hidden' &&
        c !== 'flex-col',
    )
    .join(' ')
    .replace(/\s+/g, ' ')
    .trim();
}

function patchListPanelClass(cls) {
  const parts = cls.split(/\s+/).filter(Boolean);
  const hasFill = parts.includes('list-panel--fill');
  const kept = parts.filter(
    (c) =>
      c !== 'flex' &&
      c !== 'flex-1' &&
      c !== 'min-h-0' &&
      c !== 'flex-col' &&
      c !== 'overflow-hidden' &&
      c !== 'list-panel--fill',
  );
  if (!kept.includes('list-panel') && cls.includes('list-panel')) {
    kept.unshift('list-panel');
  }
  if (hasFill) kept.push('list-panel--fill');
  return kept.join(' ');
}

function patchGridClass(cls) {
  let next = stripFlexClasses(cls);
  if (!next.includes('native-vxe-grid')) {
    next = `native-vxe-grid ${next}`.trim();
  }
  return next.replace(/\s+/g, ' ').trim();
}

function removeScopedVxeHeightBlock(src) {
  // Remove :deep(.vxe-grid) { ... height: 100% !important ... } blocks
  return src.replace(
    /:deep\(\.vxe-grid\)\s*\{[^}]*height:\s*100%\s*!important[^}]*\}\s*/g,
    '',
  );
}

function removeScopedGridHeightChains(src) {
  // Common pattern: wrapper with height 100% on .h-full and .vxe-grid
  return src.replace(
    /:deep\(>\s*\.h-full\.rounded-md\.bg-card\)\s*\{[^}]+\}\s*/g,
    '',
  );
}

function migrateFile(file) {
  let src = fs.readFileSync(file, 'utf8');
  if (!src.includes('useVbenVxeGrid')) return null;

  const original = src;

  // gridOptions height
  src = src.replace(/height:\s*['"]100%['"]/g, "height: 'auto'");

  // Page content-class on native-list-page lines
  src = src.replace(
    /content-class="([^"]*native-list-page[^"]*)"/g,
    (match, cls) => {
      const cleaned = stripFlexClasses(cls);
      return `content-class="${cleaned}"`;
    },
  );

  // Other list page content-class patterns
  src = src.replace(
    /content-class="([^"]*(?:-list-page|-list-page|store-list|order-list|member-list|live-room-list|clerk-list|refund-list|category-list|comment-list|opt-log|login-log|shop-admin|shop-role|diy-preview|h5domain-list|table-event)[^"]*)"/g,
    (match, cls) => {
      if (!cls.includes('flex') && !cls.includes('min-h-0') && !cls.includes('overflow-hidden')) {
        return match;
      }
      const cleaned = stripFlexClasses(cls);
      return `content-class="${cleaned}"`;
    },
  );

  // list-panel classes
  src = src.replace(/class="([^"]*list-panel[^"]*)"/g, (match, cls) => {
    return `class="${patchListPanelClass(cls)}"`;
  });

  // Grid / named grid components
  src = src.replace(
    /<(Grid|JoinGrid|OrderGrid|ListGrid|LogGrid)\s+class="([^"]*)"/g,
    (match, tag, cls) => `<${tag} class="${patchGridClass(cls)}"`,
  );

  // Hub panel root: flex min-h-0 flex-1 flex-col on *-panel components (template root)
  src = src.replace(
    /class="([^"]*-panel)\s+flex\s+min-h-0\s+flex-1\s+flex-col"/g,
    'class="$1 list-panel--fill"',
  );
  src = src.replace(
    /class="flex\s+min-h-0\s+flex-1\s+flex-col"/g,
    'class="list-panel--fill"',
  );

  // NativeListPage hub: add list-panel--fill to list-panel when missing (in native-list-page.vue handled globally)

  // Scoped style cleanup
  src = removeScopedVxeHeightBlock(src);
  src = removeScopedGridHeightChains(src);

  // Remove empty scoped blocks that only had grid height rules
  src = src.replace(
    /<style scoped lang="scss">\s*\.[\w-]+\s*\{\s*\}\s*<\/style>\s*/g,
    '',
  );

  if (src === original) return null;
  fs.writeFileSync(file, src);
  return path.relative(ROOT, file);
}

const files = walk(ROOT).filter((f) => !shouldSkip(f));
const changed = [];
for (const file of files) {
  const rel = migrateFile(file);
  if (rel) changed.push(rel);
}

console.log(`Aligned ${changed.length} files:`);
for (const rel of changed.sort()) console.log(`  ${rel}`);
