#!/usr/bin/env node
/**
 * Phase 8 — 全站回归：Phase 1 静态 + Phase 1～7 E2E + Shell smoke。
 *
 * Usage:
 *   node scripts/phase8-full-regression.mjs
 *   node scripts/phase8-full-regression.mjs --smoke-only
 *   PHASE8_SKIP=phase2,phase3 node scripts/phase8-full-regression.mjs
 *
 * Env: E2E_BASE_URL, E2E_API_BASE
 */
import { spawnSync } from 'node:child_process';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.join(__dirname, '..');
const NODE = process.execPath;
const REPORT_DIR = path.join(__dirname, '../.e2e-screenshots');
const OUT = path.join(REPORT_DIR, 'phase8-full-regression-report.json');

const smokeOnly = process.argv.includes('--smoke-only');
const skipSet = new Set(
  (process.env.PHASE8_SKIP || '')
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean),
);

/** @type {{ id: string; name: string; script: string; report?: string; staticOk?: string }[]} */
const SUITES = [
  {
    id: 'p1-static',
    name: 'Phase 1 静态 registry',
    script: 'verify-plus-registry.mjs',
    staticOk: 'OK:',
  },
  {
    id: 'p1-e2e',
    name: 'Phase 1 插件中心 E2E',
    script: 'plus-plugin-audit.mjs',
    report: 'plus-audit-report.json',
  },
  {
    id: 'phase2',
    name: 'Phase 2 简单营销',
    script: 'phase2-simple-plus-audit.mjs',
    report: 'phase2-audit-report.json',
  },
  {
    id: 'phase3',
    name: 'Phase 3 Tab 型插件',
    script: 'phase3-tab-plus-audit.mjs',
    report: 'phase3-audit-report.json',
  },
  {
    id: 'phase4',
    name: 'Phase 4 商品型插件',
    script: 'phase4-product-plus-audit.mjs',
    report: 'phase4-audit-report.json',
  },
  {
    id: 'phase5',
    name: 'Phase 5 重型 Plus',
    script: 'phase5-heavy-plus-audit.mjs',
    report: 'phase5-audit-report.json',
  },
  {
    id: 'phase6',
    name: 'Phase 6 核心业务',
    script: 'phase6-core-business-audit.mjs',
    report: 'phase6-audit-report.json',
  },
  {
    id: 'phase7',
    name: 'Phase 7 设置/装修/直播',
    script: 'phase7-settings-diy-live-audit.mjs',
    report: 'phase7-audit-report.json',
  },
  {
    id: 'phase8-smoke',
    name: 'Phase 8 Shell smoke',
    script: 'phase8-shell-smoke.mjs',
    report: 'phase8-smoke-report.json',
  },
];

function formatDuration(ms) {
  if (ms < 1000) return `${ms}ms`;
  const sec = ms / 1000;
  if (sec < 60) return `${sec.toFixed(1)}s`;
  const min = Math.floor(sec / 60);
  const rem = (sec % 60).toFixed(0);
  return `${min}m ${rem}s`;
}

function readReportCounts(reportFile) {
  if (!reportFile) return null;
  const full = path.join(REPORT_DIR, reportFile);
  if (!fs.existsSync(full)) return null;
  try {
    const data = JSON.parse(fs.readFileSync(full, 'utf8'));
    if (typeof data.pass === 'number' && typeof data.total === 'number') {
      return { pass: data.pass, total: data.total, fail: data.fail ?? data.total - data.pass };
    }
  } catch {
    /* ignore */
  }
  return null;
}

function runSuite(suite) {
  const scriptPath = path.join(__dirname, suite.script);
  const started = Date.now();

  console.log(`\n${'='.repeat(72)}`);
  console.log(`>>> ${suite.name}`);
  console.log(`    node scripts/${suite.script}`);
  console.log('='.repeat(72));

  const result = spawnSync(NODE, [scriptPath], {
    cwd: ROOT,
    stdio: 'inherit',
    env: { ...process.env, E2E_CHAIN: '1' },
  });

  const durationMs = Date.now() - started;
  const exitCode = result.status ?? (result.signal ? 128 : 1);
  const counts = readReportCounts(suite.report);

  return {
    id: suite.id,
    name: suite.name,
    script: suite.script,
    ok: exitCode === 0,
    exitCode,
    durationMs,
    cases: counts,
  };
}

function printSummary(results) {
  const nameWidth = Math.max(8, ...results.map((r) => r.name.length));
  console.log('\n');
  console.log('='.repeat(72));
  console.log('PHASE 8 — FULL REGRESSION SUMMARY');
  console.log('='.repeat(72));
  console.log(
    `${'Suite'.padEnd(nameWidth)}  ${'Status'.padEnd(6)}  ${'Cases'.padEnd(12)}  ${'Time'.padEnd(10)}`,
  );
  console.log('-'.repeat(72));

  let casePass = 0;
  let caseTotal = 0;

  for (const r of results) {
    const status = r.ok ? 'PASS' : 'FAIL';
    const cases =
      r.cases != null ? `${r.cases.pass}/${r.cases.total}` : r.id === 'p1-static' ? '27 static' : '—';
    if (r.cases) {
      casePass += r.cases.pass;
      caseTotal += r.cases.total;
    } else if (r.ok && r.id === 'p1-static') {
      casePass += 27;
      caseTotal += 27;
    }
    console.log(
      `${r.name.padEnd(nameWidth)}  ${status.padEnd(6)}  ${cases.padEnd(12)}  ${formatDuration(r.durationMs).padEnd(10)}`,
    );
    if (!r.ok) console.log(`  ↳ exit code: ${r.exitCode}`);
  }

  console.log('-'.repeat(72));
  const suitePass = results.filter((r) => r.ok).length;
  const suiteFail = results.length - suitePass;
  const totalMs = results.reduce((sum, r) => sum + r.durationMs, 0);
  console.log(`Suites: ${suitePass}/${results.length} passed`);
  if (caseTotal > 0) console.log(`E2E cases (aggregated): ${casePass}/${caseTotal}`);
  console.log(`Total time: ${formatDuration(totalMs)}`);
  console.log('='.repeat(72));
}

function main() {
  const baseUrl = process.env.E2E_BASE_URL || 'http://localhost:11525';
  console.log('Merchant Admin — Phase 8 full regression');
  console.log(`Base URL: ${baseUrl}`);
  console.log(`Started:  ${new Date().toISOString()}`);
  if (skipSet.size) console.log(`Skipping: ${[...skipSet].join(', ')}`);

  const activeSuites = smokeOnly
    ? SUITES.filter((s) => s.id === 'phase8-smoke')
    : SUITES.filter((s) => !skipSet.has(s.id));

  const results = [];
  for (const suite of activeSuites) {
    results.push(runSuite(suite));
  }

  printSummary(results);

  const aggregate = {
    startedAt: new Date().toISOString(),
    smokeOnly,
    skip: [...skipSet],
    suitePass: results.filter((r) => r.ok).length,
    suiteTotal: results.length,
    suites: results,
    casePass: results.reduce((s, r) => s + (r.cases?.pass ?? (r.ok && r.id === 'p1-static' ? 27 : 0)), 0),
    caseTotal: results.reduce(
      (s, r) => s + (r.cases?.total ?? (r.ok && r.id === 'p1-static' ? 27 : 0)),
      0,
    ),
  };
  fs.mkdirSync(REPORT_DIR, { recursive: true });
  fs.writeFileSync(OUT, JSON.stringify(aggregate, null, 2));
  console.log(`\nReport: ${OUT}`);

  process.exit(results.some((r) => !r.ok) ? 1 : 0);
}

main();
