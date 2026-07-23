import fs from 'node:fs';
import path from 'node:path';

function walk(dir) {
  const out = [];
  for (const name of fs.readdirSync(dir)) {
    const full = path.join(dir, name);
    const stat = fs.statSync(full);
    if (stat.isDirectory()) out.push(...walk(full));
    else if (name.endsWith('.vue')) out.push(full);
  }
  return out;
}

const root = new URL('../src/views/native', import.meta.url).pathname;
let fixed = 0;
for (const file of walk(root)) {
  const content = fs.readFileSync(file, 'utf8');
  const next = content.replace(/\nimport \{\} from 'vue-router';\n/g, '\n');
  if (next !== content) {
    fs.writeFileSync(file, next);
    console.log('fixed', path.relative(process.cwd(), file));
    fixed += 1;
  }
}
console.log('total fixed', fixed);
