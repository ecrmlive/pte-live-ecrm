/**
 * P1：移除 list-panel / native-vxe-grid / native-list-page 壳，改为标准 Page + Grid
 * 用法：node scripts/migrate-vben-list-shell.mjs [glob dirs...]
 */
import fs from 'node:fs';
import path from 'node:path';

const roots = process.argv.slice(2).length
  ? process.argv.slice(2)
  : [
      'src/views/native/auth',
      'src/views/native/user',
      'src/views/native/setting/address',
      'src/views/native/setting/delivery',
      'src/views/native/setting/express',
      'src/views/native/setting/printer',
      'src/views/native/setting/h5domain',
    ];

function walk(dir, out = []) {
  if (!fs.existsSync(dir)) return out;
  for (const name of fs.readdirSync(dir)) {
    const full = path.join(dir, name);
    const stat = fs.statSync(full);
    if (stat.isDirectory()) walk(full, out);
    else if (name.endsWith('.vue')) out.push(full);
  }
  return out;
}

function migrate(content) {
  if (!content.includes('list-panel') && !content.includes('native-vxe-grid')) {
    return null;
  }
  // Hub fill 内嵌 panel 暂不在此脚本处理
  if (content.includes('list-panel--fill')) {
    return null;
  }

  let next = content;

  next = next.replace(
    /<Page(\s+)content-class="native-list-page[^"]*"(\s*)>/g,
    '<Page$1>',
  );
  next = next.replace(
    /<Page(\s+)content-class='native-list-page[^']*'(\s*)>/g,
    '<Page$1>',
  );
  next = next.replace(/<Page\s+content-class="native-list-page[^"]*"\s*\n\s*>/g, '<Page>');

  next = next.replace(/\n\s*<div class="list-panel">\n/g, '\n');
  next = next.replace(/\n\s*<\/div>\n(\s*<(?:Shop|Return|Member|Store|.*Form|.*Modal|.*Dialog))/g, '\n$1');

  next = next.replace(/<Grid class="native-vxe-grid([^"]*)">/g, (_, rest) => {
    const extra = rest.trim();
    return extra ? `<Grid class="${extra}">` : '<Grid>';
  });
  next = next.replace(/<Grid class='native-vxe-grid([^']*)'>/g, (_, rest) => {
    const extra = rest.trim();
    return extra ? `<Grid class="${extra}">` : '<Grid>';
  });

  // 去掉 Grid 后多余的 list-panel 闭合
  next = next.replace(/(\s*)<\/Grid>\s*\n\s*<\/div>\s*\n(\s*<)/g, '$1</Grid>\n$2');

  return next === content ? null : next;
}

let changed = 0;
let skippedFill = 0;
for (const root of roots) {
  const abs = path.resolve(root);
  for (const file of walk(abs)) {
    const content = fs.readFileSync(file, 'utf8');
    if (content.includes('list-panel--fill')) {
      skippedFill += 1;
      continue;
    }
    const next = migrate(content);
    if (next) {
      fs.writeFileSync(file, next);
      console.log('migrated', path.relative(process.cwd(), file));
      changed += 1;
    }
  }
}
console.log('done', { changed, skippedFill });
