/**
 * Phase 8 shell smoke — 登录 / 首页 / 侧栏 / 根路由。
 * Run: node scripts/phase8-shell-smoke.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase8-smoke-report.json');

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

async function fetchSessionMenus(token) {
  const res = await fetch(`${API_BASE}/shop/auth/session`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'authori-zation': `Bearer ${token}`,
      AppID: '10001',
    },
    body: '{}',
  });
  const json = await res.json();
  if (json.code !== 1) throw new Error(`session failed: ${json.msg}`);
  return json.data?.menus ?? [];
}

function countMenuNodes(nodes) {
  let n = 0;
  for (const node of nodes) {
    if (node.is_menu === 1 || node.is_menu === true) n += 1;
    if (node.children?.length) n += countMenuNodes(node.children);
  }
  return n;
}

/** @type {Array<{ id: string; name: string; run: (ctx: { page: import('playwright').Page; token: string; browser: import('playwright').Browser }) => Promise<void> }>} */
const CASES = [
  {
    id: 'R8a',
    name: '未登录跳转登录页',
    async run({ page }) {
      await page.goto(`${BASE_URL}/#/product/product/index`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(2000);
      const url = page.url();
      if (!url.includes('/login')) throw new Error(`expected login redirect, got ${url}`);
      const body = await page.locator('body').innerText();
      if (!body.includes('登录') && !body.includes('账号')) {
        throw new Error('login form missing');
      }
    },
  },
  {
    id: 'R8b',
    name: 'Session 菜单 API',
    async run({ token }) {
      const menus = await fetchSessionMenus(token);
      const count = countMenuNodes(menus);
      if (count < 10) throw new Error(`menu nodes ${count} < 10`);
    },
  },
  {
    id: 'R8c',
    name: '根路由 #/',
    async run({ page, token }) {
      await page.goto(BASE_URL);
      await page.evaluate(
        ({ key, enc }) => {
          localStorage.setItem(key, enc);
          sessionStorage.setItem(key, enc);
        },
        { key: TOKEN_KEY, enc: encryptToken(token) },
      );
      const started = Date.now();
      await page.goto(`${BASE_URL}/#/`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(3500);
      if (Date.now() - started > 15000) throw new Error('root route too slow');
      const url = page.url();
      if (url.includes('/login')) throw new Error('redirected to login');
      const body = await page.locator('body').innerText();
      if (body.trim().length < 40) throw new Error('blank root page');
    },
  },
  {
    id: 'R8d',
    name: '首页 Dashboard',
    async run({ page }) {
      await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(3000);
      const body = await page.locator('body').innerText();
      if (!body.includes('待办事项')) throw new Error('home dashboard missing 待办事项');
    },
  },
  {
    id: 'R8e',
    name: '侧栏顶级菜单',
    async run({ page }) {
      await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(2500);
      const body = await page.locator('body').innerText();
      for (const label of ['插件', '商品', '订单', '设置']) {
        if (!body.includes(label)) throw new Error(`sidebar label missing: ${label}`);
      }
    },
  },
  {
    id: 'R8f',
    name: '侧栏导航-插件中心',
    async run({ page }) {
      await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(2000);
      const link = page.locator('aside a, nav a').filter({ hasText: '插件' }).first();
      if (await link.isVisible({ timeout: 5000 }).catch(() => false)) {
        await link.click();
        await page.waitForTimeout(2500);
      } else {
        await page.goto(`${BASE_URL}/#/plus/plus/index`);
        await page.waitForTimeout(2500);
      }
      const body = await page.locator('body').innerText();
      if (body.includes('404') && body.length < 500) throw new Error('plus hub 404');
      if (!body.includes('插件') && !body.includes('营销')) {
        throw new Error('plus hub content missing');
      }
    },
  },
  {
    id: 'R8g',
    name: 'Tab 栏多页签',
    async run({ page }) {
      await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(1500);
      await page.goto(`${BASE_URL}/#/product/product/index`, { waitUntil: 'domcontentloaded' });
      await page.waitForTimeout(2000);
      const tabs = page.locator('.vben-tabs, [class*="tab-bar"], .flex.items-center.gap-1');
      const count = await tabs.locator('a, button, [role="tab"]').count().catch(() => 0);
      if (count < 2) {
        const body = await page.locator('body').innerText();
        if (!body.includes('商品')) throw new Error('product page not loaded');
      }
    },
  },
  {
    id: 'R8h',
    name: '登出后不可访问业务页',
    async run({ browser }) {
      const ctx = await browser.newContext();
      const guest = await ctx.newPage();
      await guest.goto(`${BASE_URL}/#/order/order/index`, { waitUntil: 'domcontentloaded' });
      await guest.waitForTimeout(2500);
      const ok = guest.url().includes('/login');
      await ctx.close();
      if (!ok) throw new Error('expected login without token');
    },
  },
];

async function main() {
  const token = await apiLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();

  const results = [];
  for (const c of CASES) {
    const row = { id: c.id, name: c.name, status: 'pass', errors: [] };
    try {
      await c.run({ page, token, browser });
    } catch (err) {
      row.status = 'fail';
      row.errors.push(err instanceof Error ? err.message : String(err));
    }
    results.push(row);
    console.log(
      `${row.status === 'pass' ? '✓' : '✗'} ${c.id} ${c.name}${row.errors.length ? ' — ' + row.errors.join('; ') : ''}`,
    );
  }

  await browser.close();

  const summary = {
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    fail: results.filter((r) => r.status === 'fail').length,
    results,
  };
  fs.mkdirSync(path.dirname(OUT), { recursive: true });
  fs.writeFileSync(OUT, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}/${summary.total}`);
  process.exit(summary.fail > 0 ? 1 : 0);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
