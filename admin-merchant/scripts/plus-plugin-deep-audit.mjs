/**
 * Deep functional smoke for all 插件中心 card pages.
 * Checks tabs, toolbar actions, and API errors after navigation.
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
const TOKEN_KEY = 'qixiLiveShopToken';
const SECRET_KEY = 'jjj_shop_single_admin_2024';
const OUT = path.join(__dirname, '../.e2e-screenshots/plus-deep-audit.json');

function encryptToken(token) {
  return CryptoJS.AES.encrypt(JSON.stringify(token), SECRET_KEY).toString();
}

function readCaptchaFromRedis(codeKey) {
  const key = `${codeKey}_shop_code`;
  for (const cmd of [
    `docker exec pte_live_redis redis-cli -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
    `redis-cli -p 13379 -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
  ]) {
    try {
      const out = execSync(cmd, { encoding: 'utf8' }).trim();
      if (out && out !== '(nil)' && !out.includes('Warning')) return out.split('\n').pop();
    } catch {
      /* next */
    }
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
  const codeKey = baseJson.data?.codeData?.codeKey;
  const code = readCaptchaFromRedis(codeKey);
  const loginRes = await fetch(`${API_BASE}/shop/passport/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'adminm', password: 'm123456', code, codeKey }),
  });
  const loginJson = await loginRes.json();
  if (loginJson.code !== 1) throw new Error(`login failed: ${loginJson.msg}`);
  return loginJson.data?.token || loginJson.data?.accessToken;
}

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
  const plugins = [];
  for (const cat of json.data?.list || []) {
    for (const child of cat.children || []) {
      plugins.push({ name: child.name, path: child.path, category: cat.name });
    }
  }
  return plugins;
}

async function deepCheck(page, plugin) {
  const apiErrors = [];
  const onResponse = (res) => {
    const url = res.url();
    if (
      (url.includes('/shop/') || url.includes('/api/')) &&
      res.request().method() !== 'OPTIONS'
    ) {
      res
        .json()
        .then((j) => {
          if (j && typeof j.code === 'number' && j.code !== 1) {
            apiErrors.push(`${url.split('?')[0]} → ${j.msg}`);
          }
        })
        .catch(() => {});
    }
  };
  page.on('response', onResponse);

  const hashPath = plugin.path.startsWith('/') ? plugin.path : `/${plugin.path}`;
  await page.goto(`${BASE_URL}/#${hashPath}`, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(1500);

  const tabs = page.locator('.el-tabs__item');
  const tabCount = await tabs.count();
  const tabLabels = [];
  for (let i = 0; i < Math.min(tabCount, 8); i++) {
    tabLabels.push((await tabs.nth(i).innerText()).trim());
  }

  if (tabCount > 1) {
    for (let i = 1; i < Math.min(tabCount, 6); i++) {
      await tabs.nth(i).click();
      await page.waitForTimeout(800);
    }
  }

  const hasTable = (await page.locator('.el-table, .native-list-page table').count()) > 0;
  const hasForm = (await page.locator('.el-form').count()) > 0;
  const addBtn = page.locator(
    'button:has-text("添加"), button:has-text("新增"), button:has-text("创建")',
  );
  const hasAdd = (await addBtn.count()) > 0;
  const searchInput = page.locator(
    'input[placeholder*="搜索"], input[placeholder*="查询"], input[placeholder*="名称"]',
  );
  const hasSearch = (await searchInput.count()) > 0;

  const bodyText = await page.locator('body').innerText();
  const issues = [];
  if (/加载失败|请求失败|系统错误/i.test(bodyText)) issues.push('error text on page');
  if (apiErrors.length) issues.push(`api: ${apiErrors.slice(0, 2).join('; ')}`);
  if (!hasTable && !hasForm && bodyText.trim().length < 80) issues.push('no table/form and sparse content');

  page.off('response', onResponse);

  return {
    ...plugin,
    status: issues.length ? 'warn' : 'pass',
    issues,
    ui: { tabCount, tabLabels, hasTable, hasForm, hasAdd, hasSearch },
    apiErrors: apiErrors.slice(0, 5),
  };
}

async function main() {
  const token = await apiLogin();
  const plugins = await fetchPlugins(token);
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();
  await page.goto(BASE_URL);
  await page.evaluate(
    ({ key, value }) => {
      localStorage.setItem(key, value);
      sessionStorage.setItem(key, value);
    },
    { key: TOKEN_KEY, value: encryptToken(token) },
  );

  const results = [];
  for (const plugin of plugins) {
    const r = await deepCheck(page, plugin);
    results.push(r);
    const icon = r.status === 'pass' ? '✓' : '⚠';
    const extra = r.issues.length ? ` — ${r.issues.join('; ')}` : '';
    console.log(
      `${icon} ${plugin.name}: tabs=${r.ui.tabCount} table=${r.ui.hasTable} form=${r.ui.hasForm} add=${r.ui.hasAdd}${extra}`,
    );
  }

  await browser.close();
  const summary = {
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    warn: results.filter((r) => r.status === 'warn').length,
    warnings: results.filter((r) => r.status === 'warn'),
    results,
  };
  fs.mkdirSync(path.dirname(OUT), { recursive: true });
  fs.writeFileSync(OUT, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}, Warn: ${summary.warn}`);
  console.log(`Report: ${OUT}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
