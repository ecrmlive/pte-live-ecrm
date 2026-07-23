#!/usr/bin/env node
/**
 * Audit merchant-admin image picker patterns against the standard.
 * Usage: node scripts/audit-image-picker.mjs
 */

import { readFileSync, readdirSync, statSync } from 'node:fs';
import { join, relative } from 'node:path';
import { fileURLToPath } from 'node:url';

const ROOT = fileURLToPath(new URL('..', import.meta.url));
const SRC = join(ROOT, 'src');

/** @type {Array<{ file: string; line: number; rule: string; snippet: string }>} */
const violations = [];

const SKIP_DIRS = new Set(['legacy', 'node_modules', 'dist']);

const RULES = [
  {
    id: 'raw-select-button',
    re: /<el-button[^>]*>[^<]*(?:选择图片|上传图片)/,
    message: 'Raw el-button for 选择图片/上传图片 — use ImageField or ImagePickerTrigger',
  },
  {
    id: 'upload-label',
    re: /(?:button-text|buttonText|label:\s*['"])上传图片['"]/,
    message: 'Button label 上传图片 — use default 选择图片',
  },
  {
    id: 'picker-v-if',
    re: /ImagePickerDialog[^>]*v-if|MediaPickerDialog[^>]*v-if|v-if[^>]*ImagePickerDialog|v-if[^>]*MediaPickerDialog/,
    message: 'v-if on picker dialog — mount always, control with v-model:open',
  },
  {
    id: 'round-picker-button',
    re: /ImagePickerTrigger[^>]*round|image-picker-trigger[^>]*round|ImageField[^>]*round/,
    message: 'round on image picker trigger',
  },
  {
    id: 'legacy-upload-v-if',
    re: /<Upload\s+v-if=/,
    message: 'Legacy Upload v-if pattern in src/',
  },
];

function walk(dir) {
  for (const name of readdirSync(dir)) {
    if (SKIP_DIRS.has(name)) continue;
    const path = join(dir, name);
    const st = statSync(path);
    if (st.isDirectory()) {
      walk(path);
      continue;
    }
    if (!/\.(vue|ts|tsx|js|jsx)$/.test(name)) continue;
    scanFile(path);
  }
}

function scanFile(absPath) {
  const rel = relative(ROOT, absPath);
  const text = readFileSync(absPath, 'utf8');
  const lines = text.split('\n');

  for (const rule of RULES) {
    lines.forEach((line, index) => {
      if (rule.re.test(line)) {
        violations.push({
          file: rel,
          line: index + 1,
          rule: rule.id,
          snippet: line.trim().slice(0, 120),
        });
      }
    });
  }
}

walk(SRC);

if (!violations.length) {
  console.log('audit-image-picker: OK (no violations in src/, excluding legacy/)');
  process.exit(0);
}

console.log(`audit-image-picker: ${violations.length} violation(s)\n`);
for (const v of violations) {
  console.log(`${v.file}:${v.line} [${v.rule}]`);
  console.log(`  ${v.snippet}\n`);
}
process.exit(1);
