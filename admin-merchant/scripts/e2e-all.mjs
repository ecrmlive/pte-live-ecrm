#!/usr/bin/env node
/**
 * Master merchant-admin E2E runner — executes sub-suites in fixed order.
 *
 * Usage:
 *   node scripts/e2e-all.mjs
 *   pnpm e2e:all
 *
 * Env: E2E_BASE_URL, E2E_API_BASE (see dev-login.mjs)
 */
import { spawnSync } from 'node:child_process';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '..');
const NODE = process.execPath;

/** @type {{ name: string; script: string; args?: string[] }[]} */
const SUITES = [
  { name: 'dev-login (inject)', script: 'dev-login.mjs', args: ['--inject'] },
  { name: 'menu-walk', script: 'e2e-menu-walk.mjs' },
  { name: 'plus-layout-audit', script: 'e2e-plus-layout-audit.mjs' },
  { name: 'plus-hub-tabs', script: 'e2e-plus-hub-tabs.mjs' },
  { name: 'plus-assemble-layout', script: 'e2e-plus-assemble-layout.mjs' },
  { name: 'crud-flows', script: 'e2e-crud-flows.mjs' },
  { name: 'live-control-deep', script: 'e2e-live-control-deep.mjs' },
  { name: 'product-toolbar', script: 'e2e-product-toolbar.mjs' },
];

// Phase 2 (not implemented yet):
// { name: 'settings / store', script: 'e2e-settings-store.mjs' },

function formatDuration(ms) {
  if (ms < 1000) return `${ms}ms`;
  const sec = ms / 1000;
  if (sec < 60) return `${sec.toFixed(1)}s`;
  const min = Math.floor(sec / 60);
  const rem = (sec % 60).toFixed(0);
  return `${min}m ${rem}s`;
}

function runSuite(suite) {
  const scriptPath = path.join(__dirname, suite.script);
  const args = [scriptPath, ...(suite.args ?? [])];
  const label = `${suite.script}${suite.args?.length ? ` ${suite.args.join(' ')}` : ''}`;
  const started = Date.now();

  console.log(`\n${'='.repeat(72)}`);
  console.log(`>>> ${suite.name}`);
  console.log(`    node ${label}`);
  console.log('='.repeat(72));

  const result = spawnSync(NODE, args, {
    cwd: ROOT,
    stdio: 'inherit',
    env: {
      ...process.env,
      E2E_CHAIN: '1',
    },
  });

  const durationMs = Date.now() - started;
  const exitCode = result.status ?? (result.signal ? 128 : 1);
  const ok = exitCode === 0;

  return {
    name: suite.name,
    script: label,
    ok,
    exitCode,
    signal: result.signal,
    durationMs,
  };
}

function printSummary(results) {
  const nameWidth = Math.max(4, ...results.map((r) => r.name.length));
  const statusWidth = 6;
  const timeWidth = 10;

  console.log('\n');
  console.log('='.repeat(72));
  console.log('E2E ALL — SUMMARY');
  console.log('='.repeat(72));
  console.log(
    `${'Suite'.padEnd(nameWidth)}  ${'Status'.padEnd(statusWidth)}  ${'Time'.padEnd(timeWidth)}  Script`,
  );
  console.log('-'.repeat(72));

  for (const r of results) {
    const status = r.ok ? 'PASS' : 'FAIL';
    console.log(
      `${r.name.padEnd(nameWidth)}  ${status.padEnd(statusWidth)}  ${formatDuration(r.durationMs).padEnd(timeWidth)}  ${r.script}`,
    );
    if (!r.ok && r.signal) {
      console.log(`  ↳ signal: ${r.signal}`);
    }
    if (!r.ok && r.exitCode !== 0) {
      console.log(`  ↳ exit code: ${r.exitCode}`);
    }
  }

  console.log('-'.repeat(72));
  const pass = results.filter((r) => r.ok).length;
  const fail = results.length - pass;
  const totalMs = results.reduce((sum, r) => sum + r.durationMs, 0);
  console.log(`Total: ${pass} passed, ${fail} failed (${results.length} suites, ${formatDuration(totalMs)})`);
  console.log('='.repeat(72));
}

function main() {
  const baseUrl = process.env.E2E_BASE_URL || 'http://localhost:11525';
  console.log('Merchant Admin — E2E all suites');
  console.log(`Base URL: ${baseUrl}`);
  console.log(`Started:  ${new Date().toISOString()}`);

  const results = [];
  for (const suite of SUITES) {
    const row = runSuite(suite);
    results.push(row);
  }

  printSummary(results);

  const anyFail = results.some((r) => !r.ok);
  process.exit(anyFail ? 1 : 0);
}

main();
