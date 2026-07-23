/**
 * P4 fill panels: remove list-panel--fill / list-panel / native-vxe-grid shell
 * 用法：node scripts/migrate-vben-list-shell-fill.mjs [dirs...]
 */
import fs from 'node:fs';
import path from 'node:path';

const roots = process.argv.slice(2).length
  ? process.argv.slice(2)
  : ['src/views/native/plus'];

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

function stripNativeVxeGridClass(classValue) {
  return classValue
    .split(/\s+/)
    .filter((c) => c && c !== 'native-vxe-grid')
    .join(' ')
    .trim();
}

function stripListPanelClasses(classValue) {
  return classValue
    .split(/\s+/)
    .filter((c) => c && c !== 'list-panel' && c !== 'list-panel--fill')
    .join(' ')
    .trim();
}

function migrate(content) {
  if (!content.includes('list-panel') && !content.includes('native-vxe-grid')) {
    return null;
  }

  let next = content;

  // class="list-panel foo" or class="foo list-panel--fill"
  next = next.replace(/class="([^"]*)"/g, (_, classes) => {
    if (!classes.includes('list-panel') && !classes.includes('native-vxe-grid')) {
      return `class="${classes}"`;
    }
    let cleaned = stripListPanelClasses(classes);
    cleaned = stripNativeVxeGridClass(cleaned);
    return cleaned ? `class="${cleaned}"` : '';
  });

  // bare class="list-panel--fill" on hub child components → remove attribute
  next = next.replace(/\s+class=""\s*/g, ' ');
  next = next.replace(/\s+class=''\s*/g, ' ');

  // orphan empty class on self-closing tags
  next = next.replace(/<([A-Z][A-Za-z0-9]*)\s+(\/>)/g, '<$1 $2');
  next = next.replace(/<([A-Z][A-Za-z0-9]*)\s+>/g, '<$1>');

  return next === content ? null : next;
}

let changed = 0;
for (const root of roots) {
  const abs = path.resolve(root);
  for (const file of walk(abs)) {
    const content = fs.readFileSync(file, 'utf8');
    const next = migrate(content);
    if (next) {
      fs.writeFileSync(file, next);
      changed += 1;
      console.log('migrated', path.relative(process.cwd(), file));
    }
  }
}
console.log('done', { changed });
