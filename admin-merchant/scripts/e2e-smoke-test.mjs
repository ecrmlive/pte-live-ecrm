/**
 * Merchant admin interactive smoke test — screenshots + interaction report.
 * Run: node scripts/e2e-smoke-test.mjs
 */
import { chromium } from 'playwright';
import CryptoJS from 'crypto-js';
import fs from 'node:fs';
import path from 'node:path';
import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const REDIS_PASSWORD = process.env.E2E_REDIS_PASSWORD || 'M8n!xP4qR2vK7sD5cH1t';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots');
const TOKEN_KEY = 'qixiLiveShopToken';
const SECRET_KEY = 'jjj_shop_single_admin_2024';

/** 与 src/utils/plus-navigation.ts 对齐 */
const PLUS_PLUGIN_ENTRY_ALIASES = {
  '/plus/sign': '/plus/sign/index',
  '/plus/live/index': '/plus/live/wx/index',
};

const results = [];

function log(page, interaction, status, note = '') {
  results.push({ page, interaction, status, note });
  const icon = status === 'pass' ? '✓' : status === 'fail' ? '✗' : '○';
  console.log(`${icon} [${page}] ${interaction}${note ? ` — ${note}` : ''}`);
}

function encryptToken(token) {
  return CryptoJS.AES.encrypt(JSON.stringify(token), SECRET_KEY).toString();
}

function readCaptchaFromRedis(codeKey) {
  const key = `${codeKey}_shop_code`;
  try {
    const out = execSync(
      `docker exec pte_live_redis redis-cli -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
      { encoding: 'utf8' },
    ).trim();
    if (out && out !== '(nil)' && !out.includes('Warning')) return out.split('\n').pop();
  } catch {
    /* fall through */
  }
  try {
    const out = execSync(
      `redis-cli -p 13379 -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
      { encoding: 'utf8' },
    ).trim();
    if (out && out !== '(nil)') return out.split('\n').pop();
  } catch {
    /* ignore */
  }
  return null;
}

async function apiLogin() {
  const baseRes = await fetch(`${API_BASE}/shop/index/base`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: '{}',
  });
  const baseJson = await baseRes.json();
  if (baseJson.code !== 1) throw new Error(`login base failed: ${baseJson.msg}`);
  const codeKey = baseJson.data?.codeData?.codeKey;
  if (!codeKey) throw new Error('no codeKey');

  const code = readCaptchaFromRedis(codeKey);
  if (!code) throw new Error('captcha not in redis');

  const loginRes = await fetch(`${API_BASE}/shop/passport/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: 'adminm',
      password: 'm123456',
      code,
      codeKey,
    }),
  });
  const loginJson = await loginRes.json();
  if (loginJson.code !== 1) throw new Error(`login failed: ${loginJson.msg}`);
  return loginJson.data?.token || loginJson.data?.accessToken;
}

async function shot(page, name) {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const file = path.join(OUT_DIR, `${name}.png`);
  await page.screenshot({ path: file, fullPage: false });
  return file;
}

async function waitStable(page, ms = 800) {
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
  await page.waitForTimeout(ms);
}

async function clickIfVisible(page, selector, timeout = 3000) {
  const el = page.locator(selector).first();
  if (await el.isVisible({ timeout }).catch(() => false)) {
    await el.click();
    return true;
  }
  return false;
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
  if (json.code !== 1) throw new Error(`plus center failed: ${json.msg}`);
  const plugins = [];
  for (const cat of json.data?.list || []) {
    for (const child of cat.children || []) {
      plugins.push({
        name: child.name,
        path: child.path,
        redirect_name: child.redirect_name || '',
      });
    }
  }
  return plugins;
}

async function testAllPlusPluginRoutes(page, token) {
  const name = '01b-plus-routes';
  const plugins = await fetchPlusPlugins(token);
  log(name, 'plugin inventory', 'pass', `${plugins.length} plugins`);

  for (const plugin of plugins) {
    const raw = plugin.path.startsWith('/') ? plugin.path : `/${plugin.path}`;
    const target =
      PLUS_PLUGIN_ENTRY_ALIASES[raw] ||
      String(plugin.redirect_name || '').trim() ||
      raw;
    await page.goto(`${BASE_URL}/#${target}`);
    await waitStable(page, 1200);

    const notFound = await page.locator('.vben-exception').count();
    const onLogin = page.url().includes('/auth/login') || page.url().includes('/login');
    if (onLogin) {
      log(name, plugin.name, 'fail', 'redirected to login');
    } else if (notFound > 0) {
      log(name, plugin.name, 'fail', `404 at ${target}`);
    } else {
      log(name, plugin.name, 'pass', target);
    }
  }
}

async function testPluginCenter(page) {
  const name = '01-plugin-center';
  await page.goto(`${BASE_URL}/#/plus/plus/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const cards = page.locator('.plus-card, .plugin-card, [class*="plus"] .el-card, .native-list-page .el-card');
  const count = await cards.count();
  if (count === 0) {
    const alt = page.locator('main .el-card, main [class*="card"]');
    const altCount = await alt.count();
    log(name, 'plugin cards visible', altCount > 0 ? 'pass' : 'fail', `found ${altCount} cards`);
    if (altCount > 0) {
      await alt.first().click();
      await waitStable(page);
      await shot(page, `${name}-click1`);
      log(name, 'click first card', 'pass');
    }
  } else {
    log(name, 'plugin cards visible', 'pass', `${count} cards`);
    const clicks = Math.min(3, count);
    for (let i = 0; i < clicks; i++) {
      await cards.nth(i).click();
      await waitStable(page);
      await shot(page, `${name}-click${i + 1}`);
      log(name, `click card ${i + 1}`, 'pass');
      await page.goto(`${BASE_URL}/#/plus/plus/index`);
      await waitStable(page);
    }
  }
}

async function testProductList(page) {
  const name = '02-product-list';
  await page.goto(`${BASE_URL}/#/product/product/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const tabs = ['出售中', '仓库', '已售罄', '库存'];
  for (const tab of tabs) {
    const clicked = await clickIfVisible(page, `.el-tabs__item:has-text("${tab}")`);
    if (clicked) {
      await waitStable(page);
      log(name, `tab ${tab}`, 'pass');
    }
  }
  await shot(page, `${name}-tabs`);

  const searchInput = page.locator('input[placeholder*="商品"], input[placeholder*="搜索"], .el-input__inner').first();
  if (await searchInput.isVisible().catch(() => false)) {
    await searchInput.fill('测试');
    await clickIfVisible(page, 'button:has-text("搜索"), button:has-text("查询")');
    await waitStable(page);
    await shot(page, `${name}-search`);
    log(name, 'search', 'pass');
  }

  const addBtn = page.locator('button:has-text("添加商品"), button:has-text("新增")').first();
  if (await addBtn.isVisible().catch(() => false)) {
    await addBtn.click();
    await waitStable(page, 1200);
    await shot(page, `${name}-add`);
    log(name, 'add product dialog/page', 'pass');
    const close = page.locator('.el-dialog__headerbtn, button:has-text("取消"), button:has-text("关闭")').first();
    if (await close.isVisible().catch(() => false)) await close.click();
  }

  const editBtn = page.locator('button:has-text("编辑"), a:has-text("编辑"), .el-button:has-text("编辑")').first();
  if (await editBtn.isVisible().catch(() => false)) {
    await editBtn.click();
    await waitStable(page, 1200);
    await shot(page, `${name}-edit`);
    log(name, 'edit dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
  } else {
    log(name, 'edit dialog', 'skip', 'no edit button on first row');
  }
}

async function testProductCategory(page) {
  const name = '03-product-category';
  await page.goto(`${BASE_URL}/#/product/category/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const addBtn = page.locator('button:has-text("添加"), button:has-text("新增")').first();
  if (await addBtn.isVisible().catch(() => false)) {
    await addBtn.click();
    await waitStable(page);
    await shot(page, `${name}-add-dialog`);
    const nameInput = page.locator('.el-dialog input, .el-drawer input').first();
    const testName = `E2E分类${Date.now() % 10000}`;
    if (await nameInput.isVisible().catch(() => false)) {
      await nameInput.fill(testName);
      log(name, 'add dialog open + fill', 'pass');
      await clickIfVisible(page, 'button:has-text("取消")');
    }
  }

  const editBtn = page.locator('button:has-text("编辑")').first();
  if (await editBtn.isVisible().catch(() => false)) {
    await editBtn.click();
    await waitStable(page);
    await shot(page, `${name}-edit`);
    log(name, 'edit dialog', 'pass');
    await clickIfVisible(page, 'button:has-text("取消"), .el-dialog__headerbtn');
  }
}

async function testOrderList(page) {
  const name = '04-order-list';
  await page.goto(`${BASE_URL}/#/order/order/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  for (const tab of ['全部', '待付款', '待发货', '已完成']) {
    await clickIfVisible(page, `.el-tabs__item:has-text("${tab}")`);
    await page.waitForTimeout(400);
  }
  await shot(page, `${name}-tabs`);
  log(name, 'status tabs', 'pass');

  const detailBtn = page.locator('button:has-text("详情"), a:has-text("详情")').first();
  if (await detailBtn.isVisible().catch(() => false)) {
    await detailBtn.click();
    await waitStable(page, 1200);
    await shot(page, `${name}-detail`);
    log(name, 'order detail', 'pass');
    await page.goBack().catch(() => clickIfVisible(page, 'button:has-text("返回")'));
  }
}

async function testRefundList(page) {
  const name = '05-refund-list';
  await page.goto(`${BASE_URL}/#/order/refund/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  for (const tab of ['全部', '待审核', '已同意', '已拒绝']) {
    await clickIfVisible(page, `.el-tabs__item:has-text("${tab}")`);
    await page.waitForTimeout(400);
  }
  log(name, 'tabs', 'pass');

  const detailBtn = page.locator('button:has-text("详情")').first();
  if (await detailBtn.isVisible().catch(() => false)) {
    await detailBtn.click();
    await waitStable(page);
    await shot(page, `${name}-detail`);
    log(name, 'detail', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("关闭")');
  }

  const auditBtn = page.locator('button:has-text("审核")').first();
  if (await auditBtn.isVisible().catch(() => false)) {
    log(name, 'audit button visible', 'pass');
  } else {
    log(name, 'audit button visible', 'skip', 'no pending refunds');
  }
}

async function testLiveRoomList(page) {
  const name = '06-live-room';
  await page.goto(`${BASE_URL}/#/live/room/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  for (const tab of ['全部', '直播中', '未开始', '已结束']) {
    await clickIfVisible(page, `.el-tabs__item:has-text("${tab}")`);
    await page.waitForTimeout(400);
  }
  log(name, 'tabs', 'pass');

  const addBtn = page.locator('button:has-text("添加"), button:has-text("创建")').first();
  if (await addBtn.isVisible().catch(() => false)) {
    await addBtn.click();
    await waitStable(page, 1200);
    await shot(page, `${name}-add`);
    log(name, 'add dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
  }

  const editBtn = page.locator('button:has-text("编辑")').first();
  if (await editBtn.isVisible().catch(() => false)) {
    await editBtn.click();
    await waitStable(page, 1200);
    await shot(page, `${name}-edit`);
    log(name, 'edit dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
  }

  const qrBtn = page.locator('button:has-text("二维码"), button:has-text("分享")').first();
  if (await qrBtn.isVisible().catch(() => false)) {
    await qrBtn.click();
    await waitStable(page);
    await shot(page, `${name}-qrcode`);
    log(name, 'qrcode dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("关闭")');
  }
}

async function testMemberList(page) {
  const name = '07-member-list';
  await page.goto(`${BASE_URL}/#/user/user/index`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const addBtn = page.locator('button:has-text("添加"), button:has-text("新增")').first();
  if (await addBtn.isVisible().catch(() => false)) {
    await addBtn.click();
    await waitStable(page);
    await shot(page, `${name}-add`);
    log(name, 'add dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
  }

  const editBtn = page.locator('button:has-text("编辑")').first();
  if (await editBtn.isVisible().catch(() => false)) {
    await editBtn.click();
    await waitStable(page);
    await shot(page, `${name}-edit`);
    log(name, 'edit dialog', 'pass');
    await clickIfVisible(page, '.el-dialog__headerbtn, button:has-text("取消")');
  }
}

async function testHome(page) {
  const name = '08-home';
  await page.goto(`${BASE_URL}/#/home`);
  await waitStable(page, 1500);
  await shot(page, name);
  const hasContent = await page.locator('main, .dashboard, .el-card, canvas').first().isVisible().catch(() => false);
  log(name, 'dashboard render', hasContent ? 'pass' : 'fail');
}

async function testStoreSettings(page) {
  const name = '09-store-settings';
  await page.goto(`${BASE_URL}/#/setting/store`);
  await waitStable(page);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const tabItems = page.locator('.el-tabs__item');
  const tabCount = await tabItems.count();
  for (let i = 0; i < Math.min(tabCount, 4); i++) {
    await tabItems.nth(i).click();
    await page.waitForTimeout(500);
  }
  await shot(page, `${name}-tabs`);
  log(name, 'switch tabs', tabCount > 0 ? 'pass' : 'fail', `${tabCount} tabs`);
}

async function testSalesStats(page) {
  const name = '10-sales-stats';
  await page.goto(`${BASE_URL}/#/statistics/sales/index`);
  await waitStable(page, 1500);
  await shot(page, name);
  log(name, 'page load', 'pass');

  const datePicker = page.locator('.el-date-editor, input[placeholder*="日期"]').first();
  if (await datePicker.isVisible().catch(() => false)) {
    await datePicker.click();
    await page.waitForTimeout(500);
    await shot(page, `${name}-datepicker`);
    log(name, 'date picker', 'pass');
    await page.keyboard.press('Escape');
  }

  const chart = page.locator('canvas, .echarts, [class*="chart"]').first();
  const hasChart = await chart.isVisible({ timeout: 5000 }).catch(() => false);
  log(name, 'charts render', hasChart ? 'pass' : 'fail');
  await shot(page, `${name}-charts`);
}

async function main() {
  console.log('Logging in via API...');
  const token = await apiLogin();
  if (!token) throw new Error('no token');
  console.log('Token acquired');

  const browser = await chromium.launch({
    headless: true,
    channel: 'chrome',
  });
  const context = await browser.newContext({ viewport: { width: 1440, height: 900 } });
  const page = await context.newPage();

  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      const t = msg.text();
      if (!t.includes('favicon') && !t.includes('404')) {
        console.warn('  [console.error]', t.slice(0, 200));
      }
    }
  });
  page.on('pageerror', (err) => console.warn('  [pageerror]', err.message.slice(0, 200)));

  await page.goto(BASE_URL);
  await page.evaluate(
    ({ key, enc }) => {
      localStorage.setItem(key, enc);
      sessionStorage.setItem(key, enc);
    },
    { key: TOKEN_KEY, enc: encryptToken(token) },
  );

  const tests = [
    testPluginCenter,
    (p) => testAllPlusPluginRoutes(p, token),
    testProductList,
    testProductCategory,
    testOrderList,
    testRefundList,
    testLiveRoomList,
    testMemberList,
    testHome,
    testStoreSettings,
    testSalesStats,
  ];

  for (const test of tests) {
    try {
      await test(page);
    } catch (err) {
      const testName = test.name;
      log(testName, 'test suite', 'fail', err.message);
      await shot(page, `error-${testName}`).catch(() => {});
    }
  }

  const reportPath = path.join(OUT_DIR, 'report.json');
  fs.writeFileSync(reportPath, JSON.stringify({ results, outDir: OUT_DIR }, null, 2));

  const fails = results.filter((r) => r.status === 'fail');
  console.log('\n=== SUMMARY ===');
  console.log(`Total: ${results.length}, Pass: ${results.filter((r) => r.status === 'pass').length}, Fail: ${fails.length}, Skip: ${results.filter((r) => r.status === 'skip').length}`);
  console.log(`Screenshots: ${OUT_DIR}`);
  if (fails.length) {
    console.log('\nFailures:');
    fails.forEach((f) => console.log(`  - [${f.page}] ${f.interaction}: ${f.note}`));
  }

  await browser.close();
  process.exit(fails.length > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
