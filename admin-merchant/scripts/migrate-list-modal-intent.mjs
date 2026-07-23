#!/usr/bin/env node
/**
 * 批量移除 syncRouteAction，改为 consumeListEditIntent / consumeListRouteIntent
 * Run: node scripts/migrate-list-modal-intent.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const root = path.join(path.dirname(fileURLToPath(import.meta.url)), '..');
const views = path.join(root, 'src/views/native');

/** @type {Array<{ file: string; listKey: string; editBody: string; extraImport?: string; onMountedExtra?: string; useRouteIntent?: boolean }>} */
const MIGRATIONS = [
  {
    file: 'auth/user/index.vue',
    listKey: 'auth-user',
    editBody: `curModel.value = { shop_user_id: id };\n    openEdit.value = true;`,
  },
  {
    file: 'auth/role/index.vue',
    listKey: 'auth-role',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'user/user/index.vue',
    listKey: 'user-member',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'user/tag/index.vue',
    listKey: 'user-tag',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'user/equity/index.vue',
    listKey: 'user-equity',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'user/grade/grade-list-panel.vue',
    listKey: 'user-grade',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'setting/address/index.vue',
    listKey: 'setting-address',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'setting/delivery/index.vue',
    listKey: 'setting-delivery',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'setting/express/index.vue',
    listKey: 'setting-express',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'setting/printer/index.vue',
    listKey: 'setting-printer',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/fullreduce/index.vue',
    listKey: 'plus-fullreduce',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/buyactivity/index.vue',
    listKey: 'plus-buyactivity',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/coupon/coupon-list-panel.vue',
    listKey: 'plus-coupon',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/seckill/seckill-active-panel.vue',
    listKey: 'plus-seckill-active',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/seckill/seckill-time-panel.vue',
    listKey: 'plus-seckill-time',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/article/article-article-panel.vue',
    listKey: 'plus-article',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/surface/surface-template-panel.vue',
    listKey: 'plus-surface-template',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/surface/surface-setting-panel.vue',
    listKey: 'plus-surface-setting',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/table/table-table-panel.vue',
    listKey: 'plus-table-table',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/agent/agent-poster-panel.vue',
    listKey: 'plus-agent-poster',
    editBody: `editPosterId.value = id;\n    openForm.value = true;`,
  },
  {
    file: 'plus/card/card-code-panel.vue',
    listKey: 'plus-card-code',
    editBody: `editingId.value = id;\n    openForm.value = true;`,
  },
];

function stripSyncRouteAction(content) {
  let next = content;
  // remove function syncRouteAction ... } block (greedy until next function or onMounted/watch at col 0)
  next = next.replace(
    /\n(?:async )?function syncRouteAction\(\) \{[\s\S]*?\n\}\n/g,
    '\n',
  );
  next = next.replace(
    /\nwatch\(\(\) => \[route\.path[\s\S]*?\);\n/g,
    '\n',
  );
  next = next.replace(/\nwatch\(\(\) => route\.path, syncRouteAction\);\n/g, '\n');
  next = next.replace(/\n  syncRouteAction\(\);\n/g, '\n');
  next = next.replace(/\n  void syncRouteAction\(\);\n/g, '\n');
  return next;
}

function ensureImport(content, importLine) {
  if (content.includes(importLine.trim())) return content;
  const marker = "from '#/utils/list-modal-route';";
  if (content.includes(marker)) {
    return content.replace(
      /import \{([^}]+)\} from '#\/utils\/list-modal-route';/,
      (m, inner) => {
        const names = inner.split(',').map((s) => s.trim()).filter(Boolean);
        const add = importLine.match(/\{([^}]+)\}/)?.[1]?.split(',').map((s) => s.trim()) ?? [];
        for (const n of add) {
          if (!names.includes(n)) names.push(n);
        }
        return `import { ${names.join(', ')} } from '#/utils/list-modal-route';`;
      },
    );
  }
  const lastImport = content.lastIndexOf("\nimport ");
  const end = content.indexOf('\n', lastImport + 1);
  return `${content.slice(0, end + 1)}\nimport ${importLine} from '#/utils/list-modal-route';${content.slice(end + 1)}`;
}

function cleanupImports(content) {
  let next = content;
  if (!next.includes('route.') && !next.includes('route,')) {
    next = next.replace(/\nimport \{([^}]*?)useRoute,?\s*/g, '\nimport {');
    next = next.replace(/,\s*useRoute/g, '');
    next = next.replace(/useRoute,\s*/g, '');
    next = next.replace(/\nconst route = useRoute\(\);\n/g, '\n');
  }
  if (!next.includes('router.')) {
    next = next.replace(/\nimport \{([^}]*?)useRouter,?\s*/g, '\nimport {');
    next = next.replace(/,\s*useRouter/g, '');
    next = next.replace(/useRouter,\s*/g, '');
    next = next.replace(/\nconst router = useRouter\(\);\n/g, '\n');
  }
  if (!next.includes('watch(')) {
    next = next.replace(/,\s*watch/g, '');
    next = next.replace(/watch,\s*/g, '');
    next = next.replace(/\{ watch \}/g, '{}');
    next = next.replace(/import \{\s*\} from 'vue';\n/g, '');
  }
  next = next.replace(/import \{\s*,/g, 'import {');
  next = next.replace(/,\s*\} from 'vue'/g, " } from 'vue'");
  return next;
}

function applyMigration({ file, listKey, editBody, onMountedExtra }) {
  const full = path.join(views, file);
  if (!fs.existsSync(full)) {
    console.warn('skip missing', file);
    return;
  }
  let content = fs.readFileSync(full, 'utf8');
  if (!content.includes('syncRouteAction')) {
    console.log('skip (no syncRouteAction)', file);
    return;
  }

  content = stripSyncRouteAction(content);
  content = ensureImport(content, '{ consumeListEditIntent }');

  const applyFn = `function applyListModalIntent() {
  consumeListEditIntent('${listKey}', (id) => {
    ${editBody}
  });
}`;

  if (!content.includes('function applyListModalIntent')) {
    const mountMatch = content.match(/\nonMounted\(\(\) => \{/);
    if (mountMatch) {
      content = content.replace(/\nonMounted\(\(\) => \{/, `\n${applyFn}\n\nonMounted(() => {\n  applyListModalIntent();`);
    } else {
      content = content.replace(
        /<\/script>/,
        `\n${applyFn}\n\nonMounted(() => {\n  applyListModalIntent();\n});\n</script>`,
      );
    }
  }

  if (onMountedExtra) {
    content = content.replace(
      'applyListModalIntent();',
      `applyListModalIntent();\n  ${onMountedExtra}`,
    );
  }

  content = cleanupImports(content);
  fs.writeFileSync(full, content);
  console.log('migrated', file);
}

for (const m of MIGRATIONS) {
  applyMigration(m);
}

// h5domain: add only — strip sync, no consume
for (const file of ['setting/h5domain/index.vue']) {
  const full = path.join(views, file);
  let content = fs.readFileSync(full, 'utf8');
  if (!content.includes('syncRouteAction')) continue;
  content = stripSyncRouteAction(content);
  content = cleanupImports(content);
  fs.writeFileSync(full, content);
  console.log('migrated (add-only)', file);
}

// agent grade/user: add only
for (const file of ['plus/agent/agent-grade-panel.vue', 'plus/agent/agent-user-panel.vue']) {
  const full = path.join(views, file);
  let content = fs.readFileSync(full, 'utf8');
  if (!content.includes('syncRouteAction')) continue;
  content = stripSyncRouteAction(content);
  content = cleanupImports(content);
  fs.writeFileSync(full, content);
  console.log('migrated (add-only hub)', file);
}

console.log('done');
