#!/usr/bin/env node
/**
 * Scan merchant-admin list pages for forbidden toolbar+search single-row layout.
 *
 * Run: node scripts/e2e-list-layout-audit.mjs
 *      node scripts/e2e-list-layout-audit.mjs --paths /product/product/index,/plus/agent/product/index
 * Env: E2E_BASE_URL, E2E_API_BASE (see dev-login.mjs)
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const VIEWPORT = { width: 1440, height: 900 };
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/list-layout-audit');
const GOLD_PATH = '/store/store/index';
const GOLD_NAME = 'gold-standard-store-list.png';
const REPORT_PATH = path.join(OUT_DIR, 'report.json');

function slug(text) {
  return String(text)
    .replace(/[^\w\u4e00-\u9fff-]+/g, '-')
    .replace(/^-|-$/g, '')
    .slice(0, 80);
}

function flattenMenus(items, out = []) {
  for (const item of items ?? []) {
    if (item.path && item.isMenu !== 0 && item.isShow !== 0) {
      out.push({
        name: item.name || item.meta?.title || item.path,
        path: item.path.startsWith('/') ? item.path : `/${item.path}`,
      });
    }
    if (item.children?.length) {
      flattenMenus(item.children, out);
    }
  }
  return out;
}

async function fetchMenuPaths(token) {
  const res = await fetch(`${API_BASE}/shop/auth/session`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'authori-zation': `Bearer ${token}`,
      AppID: '10000',
    },
    body: '{}',
  });
  const json = await res.json();
  if (json.code !== 1) {
    throw new Error(`session failed: ${json.msg}`);
  }
  const menus = json.data?.menus ?? [];
  const flat = flattenMenus(menus);
  const seen = new Set();
  return flat.filter((item) => {
    if (seen.has(item.path)) return false;
    seen.add(item.path);
    return true;
  });
}

function auditLayoutDom() {
  const grids = document.querySelectorAll(
    '.native-vxe-grid .vxe-grid, .list-panel .vxe-grid, .native-list-page .vxe-grid',
  );
  const violations = [];

  for (const grid of grids) {
    const form = grid.querySelector('.vxe-grid--form-wrapper');
    const toolbar = grid.querySelector('.vxe-grid--toolbar-wrapper');
    if (!form || !toolbar) continue;

    const formRect = form.getBoundingClientRect();
    const toolbarRect = toolbar.getBoundingClientRect();
    if (formRect.height < 8 || toolbarRect.height < 8) continue;

    const overlap =
      toolbarRect.top < formRect.bottom && toolbarRect.bottom > formRect.top;
    const singleRowMix = overlap && Math.abs(formRect.top - toolbarRect.top) < 24;

    if (singleRowMix) {
      violations.push({
        type: 'toolbar-search-single-row',
        formTop: Math.round(formRect.top),
        toolbarTop: Math.round(toolbarRect.top),
      });
    }

    const toolbarStyle = window.getComputedStyle(toolbar);
    if (toolbarStyle.position === 'absolute') {
      violations.push({ type: 'toolbar-absolute-overlay' });
    }
  }

  return violations;
}

async function auditPage(page, item) {
  const hashPath = item.path.startsWith('/') ? item.path : `/${item.path}`;
  const url = `${BASE_URL}/#${hashPath}`;
  const screenshotName = `${slug(item.name || hashPath)}.png`;
  const screenshotPath = path.join(OUT_DIR, screenshotName);

  await page.setViewportSize(VIEWPORT);

  try {
    await page.goto(url, { waitUntil: 'commit', timeout: 45000 });
    await page.waitForLoadState('domcontentloaded', { timeout: 15000 }).catch(() => {});
    await page.waitForLoadState('load', { timeout: 8000 }).catch(() => {});
    await page
      .waitForSelector('.vxe-grid, .native-form-page, .el-empty, .list-panel', {
        timeout: 15000,
      })
      .catch(() => {});
    await page.waitForTimeout(600);
  } catch (err) {
    return {
      ...item,
      hashPath,
      url,
      screenshot: null,
      layoutViolations: [],
      status: 'fail',
      issues: [`navigation error: ${err.message}`],
    };
  }

  const finalUrl = page.url();
  const issues = [];

  if (finalUrl.includes('/auth/login') || finalUrl.includes('/login')) {
    issues.push('redirected to login');
  }

  let layoutViolations = [];
  try {
    layoutViolations = await Promise.race([
      page.evaluate(auditLayoutDom),
      new Promise((_resolve, reject) =>
        setTimeout(() => reject(new Error('layout evaluate timeout')), 15000),
      ),
    ]);
  } catch {
    issues.push('layout evaluate failed');
  }

  for (const v of layoutViolations) {
    if (v.type === 'toolbar-search-single-row') {
      issues.push('forbidden single-row toolbar+search mix');
    } else if (v.type === 'toolbar-absolute-overlay') {
      issues.push('toolbar uses absolute overlay positioning');
    }
  }

  try {
    await page.screenshot({ fullPage: false, path: screenshotPath });
  } catch {
    /* ignore screenshot failure */
  }

  return {
    ...item,
    hashPath,
    url: finalUrl,
    screenshot: screenshotPath,
    layoutViolations,
    status: issues.length ? 'fail' : 'pass',
    issues,
  };
}

async function main() {
  const login = await merchantDevLogin();
  let paths = [];

  const pathsArg = process.argv.find((a) => a.startsWith('--paths='));
  if (pathsArg) {
    paths = pathsArg
      .slice('--paths='.length)
      .split(',')
      .map((p) => ({ name: p.trim(), path: p.trim() }));
  } else {
    paths = await fetchMenuPaths(login.token);
  }

  console.log(`Layout audit: ${paths.length} pages at ${VIEWPORT.width}px…`);
  fs.mkdirSync(OUT_DIR, { recursive: true });

  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();
  await page.goto(BASE_URL, { waitUntil: 'domcontentloaded', timeout: 30000 });
  await injectMerchantToken(page, login.token);

  const results = [];
  const skipGold = Boolean(pathsArg);

  if (!skipGold) {
    const goldResult = await auditPage(page, { name: '门店列表(金标准)', path: GOLD_PATH });
    try {
      fs.copyFileSync(
        path.join(OUT_DIR, `${slug(goldResult.name || GOLD_PATH)}.png`),
        path.join(OUT_DIR, GOLD_NAME),
      );
      console.log(`Gold screenshot: ${path.join(OUT_DIR, GOLD_NAME)}`);
    } catch {
      /* gold copy optional */
    }
    results.push(goldResult);
  }

  for (const item of paths) {
    if (!skipGold && (item.path === GOLD_PATH || item.path.endsWith(GOLD_PATH))) continue;
    const result = await auditPage(page, item);
    results.push(result);
    const icon = result.status === 'pass' ? '✓' : '✗';
    const detail = result.issues.length ? ` — ${result.issues.join('; ')}` : '';
    console.log(`${icon} ${item.name}${detail}`);
  }

  await browser.close();

  const layoutFails = results.filter((r) =>
    r.issues.some((i) => i.includes('single-row') || i.includes('absolute overlay')),
  );

  const summary = {
    viewport: VIEWPORT,
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    fail: results.filter((r) => r.status === 'fail').length,
    layoutFailures: layoutFails,
    failures: results.filter((r) => r.status === 'fail'),
    results,
  };

  fs.writeFileSync(REPORT_PATH, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}/${summary.total}, Fail: ${summary.fail}`);
  console.log(`Layout violations: ${layoutFails.length}`);
  console.log(`Screenshots: ${OUT_DIR}`);
  console.log(`Report: ${REPORT_PATH}`);

  process.exit(layoutFails.length > 0 ? 1 : 0);
}

void main().catch((err) => {
  console.error(err);
  process.exit(1);
});
