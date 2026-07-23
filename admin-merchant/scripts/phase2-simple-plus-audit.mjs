/**
 * Phase 2 E2E: 简单营销插件 M-01～M-07 页面可达 + 表单可见。
 * Run: node scripts/phase2-simple-plus-audit.mjs
 * 需已登录或脚本自动登录（同 plus-plugin-audit）。
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase2-audit-report.json');

const MODULES = [
  { id: 'M-01', name: '引导收藏', path: '/plus/collection/index', expect: ['引导收藏'] },
  { id: 'M-02', name: '公众号关注', path: '/plus/officia/index', expect: ['公众号', '关注'] },
  { id: 'M-03', name: '满额包邮', path: '/plus/fullfree/index', expect: ['满额', '包邮'] },
  { id: 'M-04', name: '首页推送', path: '/plus/homepush/index', expect: ['首页推送', '活动名称'] },
  { id: 'M-05', name: '商品推荐', path: '/plus/recommend/index', expect: ['推荐', '商品'] },
  { id: 'M-06', name: '注册有礼', path: '/plus/register/index', expect: ['注册'] },
  { id: 'M-07', name: '任务中心', path: '/plus/task/index', expect: ['任务', '成长'] },
];

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

async function auditModule(page, mod) {
  const url = `${BASE_URL}/#${mod.path}`;
  const errors = [];
  await page.goto(url, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(2500);
  const finalUrl = page.url();
  if (finalUrl.includes('/login')) errors.push('redirected to login');
  const bodyText = await page.locator('body').innerText().catch(() => '');
  if (/404|页面不存在|Not Found/i.test(bodyText) && bodyText.length < 800) {
    errors.push('404');
  }
  if (bodyText.trim().length < 40) errors.push('blank');
  const matched = mod.expect.some((kw) => bodyText.includes(kw));
  if (!matched) errors.push(`missing keywords: ${mod.expect.join('|')}`);
  const apiErrors = [];
  return {
    ...mod,
    url: finalUrl,
    status: errors.length ? 'fail' : 'pass',
    errors,
    apiErrors,
  };
}

async function main() {
  const token = await apiLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
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
  for (const mod of MODULES) {
    const result = await auditModule(page, mod);
    results.push(result);
    console.log(
      `${result.status === 'pass' ? '✓' : '✗'} ${mod.id} ${mod.name} (${mod.path})${result.errors.length ? ' — ' + result.errors.join('; ') : ''}`,
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
