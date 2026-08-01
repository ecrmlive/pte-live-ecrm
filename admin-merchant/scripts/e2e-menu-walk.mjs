#!/usr/bin/env node
/**
 * Merchant admin menu walk — login, navigate each menu, screenshot, basic interactions.
 *
 * Usage:
 *   node scripts/e2e-menu-walk.mjs
 *   node scripts/e2e-menu-walk.mjs --filter live
 *   node scripts/e2e-menu-walk.mjs --only /home,/live/index
 *
 * Env: E2E_BASE_URL (default http://localhost:11525), E2E_API_BASE (default http://127.0.0.1:11503)
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/menu-walk');
const TOKEN_KEY = 'qixiLiveShopToken';

/** 与 plus-navigation.ts 对齐 */
const PLUS_PLUGIN_ENTRY_ALIASES = {
  '/plus/sign': '/plus/sign/index',
  '/plus/live/index': '/plus/live/wx/index',
};

/** 顶级菜单顺序（init_merchant_access sort） */
const TOP_LEVEL_ORDER = [
  '/home',
  '/plus',
  '/product',
  '/live',
  '/order',
  '/user',
  '/statistics',
  '/store',
  '/page',
  '/auth',
  '/appsetting',
  '/setting',
  '/cash',
];

/** 按钮权限路由，非侧栏可访问页（对齐 pte-live-menu 权限子路由） */
const PERMISSION_ONLY_RE =
  /\/(add|edit|delete|state|audit|receipt|detail|send|end|partake|pushs|receive|qrcode|export|syn|setSyn|set-top)$/i;

/** 仅作分组 redirect 的顶级 path，侧栏点击会跳子页 */
const TOP_REDIRECT_HUBS = new Set([
  '/plus',
  '/product',
  '/order',
  '/user',
  '/statistics',
  '/store',
  '/page',
  '/auth',
  '/appsetting',
  '/setting',
  '/cash',
  '/live',
]);

function isWalkableMenuPath(path, parentPath = '') {
  if (!path || path === '/') return false;
  if (PERMISSION_ONLY_RE.test(path)) return false;
  if (TOP_REDIRECT_HUBS.has(path)) return false;
  // /plus 下除插件中心外均由 fetchPlusPlugins 注入
  if (path.startsWith('/plus/') && path !== '/plus/plus/index') {
    return false;
  }
  return true;
}

function dedupePages(pages) {
  const seen = new Set();
  return pages.filter((p) => {
    if (seen.has(p.path)) return false;
    seen.add(p.path);
    return true;
  });
}

function slug(p) {
  return p.replace(/^\//, '').replace(/[/?=&]/g, '_') || 'root';
}

const report = [];

function logEntry(entry) {
  report.push(entry);
  const icon = entry.status === 'pass' ? '✓' : entry.status === 'fail' ? '✗' : '○';
  console.log(`${icon} ${entry.path} — ${entry.status}${entry.errors?.length ? ` (${entry.errors.join('; ')})` : ''}`);
}

async function fetchMenus(token) {
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
  if (json.code !== 1) throw new Error(`session failed: ${json.msg}`);
  return json.data?.menus ?? [];
}

function flattenMenus(nodes, parentPath = '') {
  const items = [];
  const sorted = [...nodes].sort((a, b) => {
    const ai = TOP_LEVEL_ORDER.indexOf(a.path);
    const bi = TOP_LEVEL_ORDER.indexOf(b.path);
    if (parentPath === '' && ai >= 0 && bi >= 0) return ai - bi;
    return (a.access_id ?? 0) - (b.access_id ?? 0);
  });

  for (const node of sorted) {
    const isMenu = node.is_menu === 1 || node.is_menu === true;
    const isRoute = node.is_route === 1 || node.is_route === true;

    if (isMenu && isRoute && node.path && isWalkableMenuPath(node.path, parentPath)) {
      items.push({
        name: node.name,
        path: node.path.startsWith('/') ? node.path : `/${node.path}`,
      });
    }

    if (node.children?.length) {
      items.push(...flattenMenus(node.children, node.path.startsWith('/') ? node.path : `/${node.path}`));
    }
  }
  return items;
}

async function fetchPlusPlugins(token) {
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
  if (json.code !== 1) return [];
  const plugins = [];
  for (const cat of json.data?.list ?? []) {
    for (const child of cat.children ?? []) {
      const raw = child.path?.startsWith('/') ? child.path : `/${child.path ?? ''}`;
      plugins.push({
        name: child.name,
        path:
          PLUS_PLUGIN_ENTRY_ALIASES[raw] ||
          String(child.redirect_name || '').trim() ||
          raw,
      });
    }
  }
  return plugins;
}

function isLoginPageUrl(url) {
  return /#\/auth\/login(?:\?|$)/.test(url) || /#\/login(?:\?|$)/.test(url);
}

async function waitStable(page, ms = 900) {
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
  await page.waitForTimeout(ms);
}

async function clickIfVisible(page, selector, timeout = 2500) {
  const el = page.locator(selector).first();
  if (await el.isVisible({ timeout }).catch(() => false)) {
    await el.click({ timeout: 5000 }).catch(() => {});
    return true;
  }
  return false;
}

async function testPageInteractions(page, item) {
  const notes = [];
  const p = item.path;

  if (p.includes('/live/sensitive-word')) {
    return notes;
  }

  if (p === '/live/index') {
    const opened = await clickIfVisible(page, 'button:has-text("创建直播间")');
    if (opened) {
      await waitStable(page, 1200);
      const nameOk = await page.locator('input[placeholder*="直播间名称"]').isVisible().catch(() => false);
      notes.push(nameOk ? 'create-modal-ok' : 'create-modal-missing-fields');
      await clickIfVisible(page, '[role=dialog] button:has-text("取消"), button:has-text("取消")');
      await page.keyboard.press('Escape').catch(() => {});
    } else {
      notes.push('create-button-missing');
    }
    return notes;
  }

  if (p.includes('/product/product/index')) {
    await clickIfVisible(page, '.el-tabs__item:has-text("出售中")');
    await clickIfVisible(page, 'button:has-text("添加商品"), button:has-text("新增")');
    await waitStable(page, 800);
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
    return notes;
  }

  if (p.includes('/order/')) {
    await clickIfVisible(page, '.el-tabs__item:has-text("全部")');
    return notes;
  }

  if (p.includes('/setting/store')) {
    const tabs = page.locator('.el-tabs__item');
    const n = Math.min(await tabs.count(), 3);
    for (let i = 0; i < n; i++) {
      await tabs.nth(i).click().catch(() => {});
      await page.waitForTimeout(350);
    }
    return notes;
  }

  // 设置/统计/流量等表单页勿点「添加」，易触发 Vue 卸载级联错误
  if (
    p.includes('/setting/') ||
    p.includes('/statistics/') ||
    p.includes('/traffic/') ||
    p.endsWith('/index') && p.includes('/plus/')
  ) {
    return notes;
  }

  const addBtn = page.locator('button:has-text("添加"), button:has-text("新增"), button:has-text("创建")').first();
  if (await addBtn.isVisible().catch(() => false)) {
    await addBtn.click({ force: true }).catch(() => {});
    await waitStable(page, 800);
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消"), button:has-text("关闭")');
    await page.keyboard.press('Escape').catch(() => {});
    notes.push('add-dialog-tried');
  }

  return notes;
}

async function recoverSession(page, token) {
  await injectMerchantToken(page, token);
  await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'networkidle', timeout: 30000 });
  await waitStable(page, 800);
}

function isRecoverableError(errors) {
  return errors.some(
    (e) =>
      e.includes('parentNode') ||
      e.includes('vnode') ||
      e.includes('typeName'),
  );
}

async function walkPage(page, item, index, token) {
  const target = item.path;
  const shotName = `${String(index + 1).padStart(3, '0')}-${slug(target)}.png`;
  const screenshotPath = path.join(OUT_DIR, shotName);
  const errors = [];

  page.removeAllListeners('pageerror');
  page.removeAllListeners('console');
  page.on('pageerror', (e) => errors.push(e.message.slice(0, 120)));
  page.on('console', (m) => {
    if (m.type() === 'error') {
      const t = m.text();
      if (
        !t.includes('favicon') &&
        !t.includes('route component is invalid') &&
        !t.includes('ERR_NAME_NOT_RESOLVED')
      ) {
        errors.push(t.slice(0, 120));
      }
    }
  });

  try {
    await page.goto(`${BASE_URL}/#${target}`, { waitUntil: 'networkidle', timeout: 30000 });
    await waitStable(page);

    const onLogin = isLoginPageUrl(page.url());
    const notFound = (await page.locator('.vben-exception').count()) > 0;

    if (onLogin) {
      await recoverSession(page, token);
      await page.goto(`${BASE_URL}/#${target}`, { waitUntil: 'networkidle', timeout: 30000 });
      await waitStable(page);
      if (isLoginPageUrl(page.url())) {
        logEntry({ path: target, name: item.name, status: 'fail', errors: ['redirected to login'], screenshotPath });
        return false;
      }
    }
    if (notFound) {
      await page.screenshot({ path: screenshotPath, fullPage: false });
      logEntry({ path: target, name: item.name, status: 'fail', errors: ['404'], screenshotPath });
      return false;
    }

    const interactionNotes = await testPageInteractions(page, item);
    await page.screenshot({ path: screenshotPath, fullPage: false });

    logEntry({
      path: target,
      name: item.name,
      status: errors.length ? 'fail' : 'pass',
      errors: errors.length ? errors : undefined,
      interactions: interactionNotes,
      screenshotPath,
    });
    return isRecoverableError(errors);
  } catch (err) {
    await page.screenshot({ path: screenshotPath, fullPage: false }).catch(() => {});
    logEntry({
      path: target,
      name: item.name,
      status: 'fail',
      errors: [err.message || String(err)],
      screenshotPath,
    });
    return true;
  }
}

function parseArgs() {
  const args = process.argv.slice(2);
  let filter = '';
  let only = null;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--filter') filter = args[++i] ?? '';
    if (args[i] === '--only') only = (args[++i] ?? '').split(',').map((s) => s.trim()).filter(Boolean);
  }
  return { filter, only };
}

async function main() {
  const { filter, only } = parseArgs();
  fs.mkdirSync(OUT_DIR, { recursive: true });

  console.log('Logging in…');
  const { token } = await merchantDevLogin();
  console.log('Login OK');

  let pages = flattenMenus(await fetchMenus(token));

  const plugins = await fetchPlusPlugins(token);
  const plusCenterIdx = pages.findIndex((p) => p.path === '/plus/plus/index');
  if (plusCenterIdx < 0 && plugins.length) {
    const homeIdx = pages.findIndex((p) => p.path === '/home');
    pages.splice(homeIdx >= 0 ? homeIdx + 1 : 0, 0, {
      name: '插件中心',
      path: '/plus/plus/index',
    });
  }
  const hubIdx = pages.findIndex((p) => p.path === '/plus/plus/index');
  if (hubIdx >= 0 && plugins.length) {
    pages.splice(
      hubIdx + 1,
      0,
      ...plugins.map((p) => ({ name: `插件:${p.name}`, path: p.path })),
    );
  }

  if (filter) {
    pages = pages.filter((p) => p.path.includes(filter) || p.name.includes(filter));
  }
  if (only?.length) {
    pages = pages.filter((p) => only.includes(p.path));
  }

  pages = dedupePages(pages);

  console.log(`Walking ${pages.length} pages → ${OUT_DIR}\n`);

  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await injectMerchantToken(page, token);

  await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'networkidle' });
  await waitStable(page, 1200);
  if (isLoginPageUrl(page.url())) {
    throw new Error('dev-login inject failed — still on login page');
  }
  console.log('Session OK at #/home\n');

  for (let i = 0; i < pages.length; i++) {
    const needsRecover = await walkPage(page, pages[i], i, token);
    if (needsRecover) {
      await recoverSession(page, token);
    }
  }

  const summary = {
    generatedAt: new Date().toISOString(),
    baseUrl: BASE_URL,
    total: report.length,
    pass: report.filter((r) => r.status === 'pass').length,
    fail: report.filter((r) => r.status === 'fail').length,
    results: report,
    outDir: OUT_DIR,
  };

  const reportPath = path.join(OUT_DIR, 'report.json');
  fs.writeFileSync(reportPath, JSON.stringify(summary, null, 2));

  console.log('\n=== SUMMARY ===');
  console.log(`Pass: ${summary.pass} / ${summary.total}, Fail: ${summary.fail}`);
  console.log(`Report: ${reportPath}`);

  await browser.close();
  process.exit(summary.fail > 0 ? 1 : 0);
}

const isMain = process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1];
if (isMain) {
  main().catch((err) => {
    console.error(err);
    process.exit(1);
  });
}

export { flattenMenus, fetchMenus, walkPage };
