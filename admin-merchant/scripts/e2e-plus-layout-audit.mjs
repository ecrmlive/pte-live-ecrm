/**
 * Screenshot all plugin index pages and flag layout issues:
 * - >1 visible vxe-table
 * - >1 visible search form toolbar
 *
 * Run: node scripts/e2e-plus-layout-audit.mjs
 *      node scripts/e2e-plus-layout-audit.mjs --interact   # delegate to e2e-plus-hub-tabs.mjs
 * Env: E2E_BASE_URL, E2E_API_BASE (see dev-login.mjs)
 */
import { spawnSync } from 'node:child_process';
import { chromium } from 'playwright';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const VIEWPORT = { width: 1440, height: 900 };
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/plus-layout-audit');
const REPORT_PATH = path.join(OUT_DIR, 'report.json');

/** Align with src/utils/plus-navigation.ts PLUS_PLUGIN_ENTRY_ALIASES */
const PLUGIN_ENTRY_ALIASES = {
  '/plus/sign': '/plus/sign/index',
  '/plus/live/index': '/plus/live/wx/index',
};

async function fetchPlugins(token) {
  const res = await fetch(`${API_BASE}/shop/plus.plus/index`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'authori-zation': `Bearer ${token}`,
      AppID: '10000',
    },
    body: '{}',
  });
  const json = await res.json();
  if (json.code !== 1) throw new Error(`plus center failed: ${json.msg}`);
  const plugins = [];
  for (const cat of json.data?.list || []) {
    for (const child of cat.children || []) {
      plugins.push({
        name: child.name,
        path: child.path,
        redirect_name: child.redirect_name || '',
        category: cat.name,
      });
    }
  }
  return plugins;
}

function resolveHashPath(plugin) {
  const rawPath = plugin.path.startsWith('/') ? plugin.path : `/${plugin.path}`;
  const redirect = String(plugin.redirect_name || '').trim();
  return PLUGIN_ENTRY_ALIASES[rawPath] || redirect || rawPath;
}

function slug(name) {
  return name.replace(/[^\w\u4e00-\u9fff-]+/g, '-').replace(/^-|-$/g, '');
}

async function countVisible(page, selector) {
  return page.locator(selector).evaluateAll((nodes) => {
    return nodes.filter((node) => {
      const el = node;
      const style = window.getComputedStyle(el);
      if (style.display === 'none' || style.visibility === 'hidden') return false;
      const rect = el.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    }).length;
  });
}

async function auditPage(page, plugin) {
  const hashPath = resolveHashPath(plugin);
  const url = `${BASE_URL}/#${hashPath}`;
  const screenshotName = `${slug(plugin.name) || slug(plugin.path)}.png`;
  const screenshotPath = path.join(OUT_DIR, screenshotName);

  await page.setViewportSize(VIEWPORT);
  await page.goto(url, { waitUntil: 'domcontentloaded' });
  await page
    .waitForSelector('.vxe-grid, .native-form-page form, .el-empty, .list-panel', {
      timeout: 8000,
    })
    .catch(() => {});
  await page.waitForTimeout(800);

  const finalUrl = page.url();
  const issues = [];

  if (finalUrl.includes('/auth/login') || finalUrl.includes('/login')) {
    issues.push('redirected to login');
  }

  const tableCount = await countVisible(page, '.native-vxe-grid');
  const formCount = await countVisible(
    page,
    '.native-vxe-grid .vxe-grid--form-wrapper form',
  );

  if (tableCount > 1) {
    issues.push(`${tableCount} visible tables (expected 1)`);
  }
  if (formCount > 1) {
    issues.push(`${formCount} visible search forms (expected 1)`);
  }

  await page.screenshot({ fullPage: false, path: screenshotPath });

  return {
    ...plugin,
    hashPath,
    url: finalUrl,
    screenshot: screenshotPath,
    tableCount,
    formCount,
    status: issues.length ? 'fail' : 'pass',
    issues,
  };
}

async function main() {
  if (process.argv.includes('--interact')) {
    const script = path.join(__dirname, 'e2e-plus-hub-tabs.mjs');
    const extra = process.argv.includes('--no-dialog') ? ['--no-dialog'] : [];
    const child = spawnSync(process.execPath, [script, ...extra], { stdio: 'inherit' });
    process.exit(child.status ?? 1);
    return;
  }

  const login = await merchantDevLogin();
  const plugins = await fetchPlugins(login.token);
  console.log(`Auditing ${plugins.length} plugin pages at ${VIEWPORT.width}px…`);

  fs.mkdirSync(OUT_DIR, { recursive: true });

  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();
  await page.goto(BASE_URL);
  await injectMerchantToken(page, login.token);

  const results = [];
  for (const plugin of plugins) {
    const result = await auditPage(page, plugin);
    results.push(result);
    const icon = result.status === 'pass' ? '✓' : '✗';
    const detail =
      result.issues.length > 0
        ? ` — ${result.issues.join('; ')}`
        : ` (tables=${result.tableCount}, forms=${result.formCount})`;
    console.log(`${icon} ${plugin.name}${detail}`);
  }

  await browser.close();

  const summary = {
    viewport: VIEWPORT,
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    fail: results.filter((r) => r.status === 'fail').length,
    layoutFailures: results.filter(
      (r) => r.tableCount > 1 || r.formCount > 1 || r.issues.some((i) => i.includes('tables') || i.includes('forms')),
    ),
    failures: results.filter((r) => r.status === 'fail'),
    results,
  };

  fs.writeFileSync(REPORT_PATH, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}/${summary.total}, Fail: ${summary.fail}`);
  console.log(`Screenshots: ${OUT_DIR}`);
  console.log(`Report: ${REPORT_PATH}`);

  process.exit(summary.fail > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
