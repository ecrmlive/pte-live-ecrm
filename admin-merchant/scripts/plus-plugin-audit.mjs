/**
 * Audit every plugin card page under 插件中心.
 * Run: node scripts/plus-plugin-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/plus-audit-report.json');

import { resolvePlusPluginEntryPath } from './lib/plus-entry-path.mjs';

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

async function auditPage(page, plugin) {
  const rawPath = plugin.path.startsWith('/') ? plugin.path : `/${plugin.path}`;
  const hashPath = resolvePlusPluginEntryPath(plugin);
  const url = `${BASE_URL}/#${hashPath}`;
  const errors = [];
  const consoleErrors = [];

  const onConsole = (msg) => {
    if (msg.type() === 'error') {
      const t = msg.text();
      if (!t.includes('favicon')) consoleErrors.push(t.slice(0, 300));
    }
  };
  const onPageError = (err) => consoleErrors.push(err.message.slice(0, 300));
  page.on('console', onConsole);
  page.on('pageerror', onPageError);

  await page.goto(url, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(2000);

  const finalUrl = page.url();
  if (finalUrl.includes('/auth/login') || finalUrl.includes('/login')) {
    errors.push('redirected to login');
  }

  const bodyText = await page.locator('body').innerText().catch(() => '');
  if (/404|页面不存在|Not Found/i.test(bodyText) && bodyText.length < 500) {
    errors.push('404 or not found text');
  }

  const legacyFrame = page.locator('iframe[src*="legacy"], .legacy-bridge, [data-legacy-bridge]');
  const legacyCount = await legacyFrame.count();
  const isLegacy = legacyCount > 0;

  const blank = bodyText.trim().length < 30;
  if (blank) errors.push('blank page');

  const elError = await page.locator('.vben-exception').count();
  if (elError > 0) errors.push('visible error UI');

  page.off('console', onConsole);
  page.off('pageerror', onPageError);

  return {
    ...plugin,
    url: finalUrl,
    status: errors.length ? 'fail' : 'pass',
    errors,
    consoleErrors: consoleErrors.slice(0, 3),
    legacyBridge: isLegacy,
  };
}

async function main() {
  const token = await apiLogin();
  const plugins = await fetchPlugins(token);
  console.log(`Auditing ${plugins.length} plugins…`);

  const browser = await chromium.launch({
    headless: true,
    channel: 'chrome',
  });
  const page = await browser.newPage();
  await page.goto(BASE_URL);
  await page.evaluate(
    ({ key, enc }) => {
      localStorage.setItem(key, enc);
      sessionStorage.setItem(key, enc);
    },
    { key: TOKEN_KEY, enc: encryptToken(token) },
  );

  const results = [];
  for (const plugin of plugins) {
    const result = await auditPage(page, plugin);
    results.push(result);
    const icon = result.status === 'pass' ? '✓' : '✗';
    const legacy = result.legacyBridge ? ' [legacy]' : '';
    console.log(
      `${icon} ${plugin.name} (${plugin.path})${legacy}${result.errors.length ? ' — ' + result.errors.join('; ') : ''}`,
    );
  }

  await browser.close();

  const summary = {
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    fail: results.filter((r) => r.status === 'fail').length,
    legacyBridge: results.filter((r) => r.legacyBridge).length,
    failures: results.filter((r) => r.status === 'fail'),
    legacyPlugins: results.filter((r) => r.legacyBridge).map((r) => ({ name: r.name, path: r.path })),
    results,
  };

  fs.mkdirSync(path.dirname(OUT), { recursive: true });
  fs.writeFileSync(OUT, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}/${summary.total}, Fail: ${summary.fail}, Legacy: ${summary.legacyBridge}`);
  console.log(`Report: ${OUT}`);
  process.exit(summary.fail > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
