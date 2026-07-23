#!/usr/bin/env node
/**
 * Phase 1 静态检查：PLUGIN_PATH_ORDER 每条 path 在 registry 或 alias 可解析。
 * Run: node scripts/verify-plus-registry.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import {
  PLUS_PLUGIN_ENTRY_ALIASES,
  resolvePlusPluginEntryPath,
} from './lib/plus-entry-path.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const registryPath = path.join(__dirname, '../src/views/native/registry.ts');
const plusIndexPath = path.join(
  __dirname,
  '../src/views/native/plus/plus/index.vue',
);

const registry = fs.readFileSync(registryPath, 'utf8');
const plusVue = fs.readFileSync(plusIndexPath, 'utf8');

const orderMatch = plusVue.match(
  /const PLUGIN_PATH_ORDER = \[([\s\S]*?)\];/,
);
if (!orderMatch) {
  console.error('PLUGIN_PATH_ORDER not found');
  process.exit(1);
}

const paths = [...orderMatch[1].matchAll(/'([^']+)'/g)].map((m) => m[1]);
const missing = [];

for (const cardPath of paths) {
  const entry = resolvePlusPluginEntryPath({ path: cardPath });
  const key = entry.endsWith('.vue') ? entry : `${entry}.vue`;
  const registryKey = `'${key.replace(/^\//, '/')}':`;
  const altKey1 = `'${entry}.vue':`;
  const altKey2 = `'${cardPath}.vue':`;
  const ok =
    registry.includes(altKey1) ||
    registry.includes(altKey2) ||
    registry.includes(`'${entry}.vue'`) ||
    registry.includes(`'${cardPath}.vue'`);
  if (!ok) {
    missing.push({ cardPath, entry, alias: PLUS_PLUGIN_ENTRY_ALIASES[cardPath] });
  }
}

if (missing.length) {
  console.error('Registry gaps:');
  for (const m of missing) {
    console.error(`  card ${m.cardPath} → ${m.entry}`);
  }
  process.exit(1);
}

console.log(`OK: ${paths.length} plugin cards resolve to registry entries.`);
