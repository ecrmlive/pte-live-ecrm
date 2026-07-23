#!/usr/bin/env node
/** Pass 2: wrap standalone grids in list-panel; fix content-class; clean scoped grid flex */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '../src/views/native');

const SKIP = [
  /\/components\//, /\/dialog\//, /\/control\//, /-field\.vue$/, /-grid\.vue$/,
  /statistics\.vue$/, /fullreduce\/product\.vue$/, /order-detail/,
];

function walk(dir, out = []) {
  for (const name of fs.readdirSync(dir)) {
    const full = path.join(dir, name);
    if (fs.statSync(full).isDirectory()) walk(full, out);
    else if (name.endsWith('.vue')) out.push(full);
  }
  return out;
}

function shouldSkip(file) {
  return SKIP.some((re) => re.test(path.relative(ROOT, file)));
}

function ensureNativeListPageClass(cls) {
  const parts = cls.split(/\s+/).filter(Boolean);
  const filtered = parts.filter(
    (c) => c !== 'flex' && c !== 'flex-1' && c !== 'min-h-0' && c !== 'flex-col' && c !== 'overflow-hidden',
  );
  if (!filtered.includes('native-list-page')) filtered.unshift('native-list-page');
  return filtered.join(' ');
}

function wrapGridInListPanel(src) {
  if (src.includes('list-panel')) return src;
  if (!src.includes('<Page') || !src.includes('useVbenVxeGrid')) return src;

  // Wrap first Grid/JoinGrid block inside Page
  return src.replace(
    /(<Page[^>]*>\s*)(<(Grid|JoinGrid|OrderGrid|ListGrid|LogGrid)\s)/,
    '$1<div class="list-panel">\n    $2',
  ).replace(
    /(<\/(Grid|JoinGrid|OrderGrid|ListGrid|LogGrid)>)(\s*(?:<[^/][\s\S]*?<\/[^>]+>\s*)*<\/Page>)/,
    (m, closeGrid, _tag, rest) => {
      if (!m.includes('list-panel')) {
        return `${closeGrid}\n    </div>${rest}`;
      }
      return m;
    },
  );
}

function fixPanelRoot(src) {
  if (!src.includes('-panel.vue') && !src.endsWith('Panel.vue')) return src;
  if (!src.includes('useVbenVxeGrid')) return src;
  // template single root div without list-panel--fill
  return src.replace(
    /<template>\s*\n\s*<div class="([^"]+)">\s*\n\s*<(Grid|JoinGrid|ListGrid|LogGrid|ElTabs)/,
    (match, cls, next) => {
      if (cls.includes('list-panel')) return match;
      const nextCls = cls.includes('list-panel--fill') ? cls : `${cls} list-panel--fill`.trim();
      return `<template>\n  <div class="${nextCls}">\n    <${next}`;
    },
  );
}

function cleanScopedGridFlex(src) {
  // Remove blocks like .xxx-grid { flex:1; min-height:0; ... :deep(.vxe-grid) height 100% }
  return src
    .replace(
      /\.[\w-]+(?:-grid|Grid)\s*\{\s*display:\s*flex;\s*flex:\s*1;\s*flex-direction:\s*column;\s*min-height:\s*0;\s*\}/g,
      '',
    )
    .replace(
      /\.[\w-]+(?:-grid|Grid)\s*\{\s*flex:\s*1;\s*min-height:\s*0;\s*\}/g,
      '',
    );
}

function fixNamedGrids(src) {
  return src.replace(
    /<(ListGrid|LogGrid|JoinGrid|OrderGrid)\s+([^>]*?)class="([^"]*)"/g,
    (match, tag, mid, cls) => {
      const cleaned = cls
        .split(/\s+/)
        .filter((c) => c && c !== 'flex-1' && c !== 'min-h-0')
        .join(' ');
      const withNative = cleaned.includes('native-vxe-grid')
        ? cleaned
        : `native-vxe-grid ${cleaned}`.trim();
      return `<${tag} ${mid}class="${withNative}"`;
    },
  );
}

function migrateFile(file) {
  let src = fs.readFileSync(file, 'utf8');
  if (!src.includes('useVbenVxeGrid')) return null;
  const original = src;

  src = src.replace(/content-class="([^"]*)"/g, (match, cls) => {
    if (
      !src.includes('useVbenVxeGrid') ||
      cls.includes('native-form') ||
      cls.includes('diy-editor') ||
      cls.includes('hub-panel') && !cls.includes('list')
    ) {
      // still fix overflow on list-like pages
      if (
        (cls.includes('overflow-hidden') || cls.includes('flex-col')) &&
        (file.includes('/index.vue') || file.includes('-panel.vue')) &&
        !cls.includes('native-form')
      ) {
        const isList =
          src.includes('pagerConfig') ||
          src.includes('withNativeListSearchLayout') ||
          src.includes('#toolbar-actions');
        if (isList) return `content-class="${ensureNativeListPageClass(cls)}"`;
      }
      return match;
    }
    return `content-class="${ensureNativeListPageClass(cls)}"`;
  });

  src = fixNamedGrids(src);
  src = wrapGridInListPanel(src);
  src = fixPanelRoot(src);
  src = cleanScopedGridFlex(src);

  // ListGrid flex-1 in class attr (agent-grade)
  src = src.replace(
    /class="([^"]*native-vxe-grid[^"]*)"/g,
    (m, cls) =>
      `class="${cls
        .split(/\s+/)
        .filter((c) => c && c !== 'flex-1' && c !== 'min-h-0')
        .join(' ')}"`,
  );

  if (src === original) return null;
  fs.writeFileSync(file, src);
  return path.relative(ROOT, file);
}

const changed = [];
for (const file of walk(ROOT).filter((f) => !shouldSkip(f))) {
  const rel = migrateFile(file);
  if (rel) changed.push(rel);
}
console.log(`Pass 2: ${changed.length} files`);
changed.sort().forEach((r) => console.log(`  ${r}`));
