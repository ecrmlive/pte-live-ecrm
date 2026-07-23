/**
 * Phase 6 E2E: 核心业务 B1～B6 — 商品 / 订单 / 会员 / 门店 / 统计 / 权限。
 * Run: node scripts/phase6-core-business-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase6-audit-report.json');

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

async function apiPost(token, apiPath, body = {}) {
  const res = await fetch(`${API_BASE}${apiPath}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'authori-zation': `Bearer ${token}`,
      AppID: '10001',
    },
    body: JSON.stringify(body),
  });
  return res.json();
}

async function resolveDynamicCases(token) {
  const dynamic = [];

  const orders = await apiPost(token, '/shop/order.order/index', {
    page: 1,
    list_rows: 1,
    dataType: 'all',
  });
  const orderId = orders.data?.list?.data?.[0]?.order_id;
  if (orderId) {
    dynamic.push({
      id: 'B2c',
      name: '订单-详情弹窗',
      hash: `/order/order/detail?order_id=${orderId}`,
      expectBody: ['订单', '详情'],
    });
  }

  const refunds = await apiPost(token, '/shop/order.refund/index', {
    page: 1,
    list_rows: 1,
    status: 0,
  });
  const refundId = refunds.data?.list?.data?.[0]?.order_refund_id;
  if (refundId) {
    dynamic.push({
      id: 'B2d',
      name: '售后-详情弹窗',
      hash: `/order/refund/detail?order_refund_id=${refundId}`,
      expectBody: ['退款', '订单'],
    });
  }

  return dynamic;
}

/** @type {Array<{ id: string; name: string; hash: string; expectTab?: string; expectBody: string[] }>} */
function buildStaticCases() {
  return [
    // B1 商品
    { id: 'B1a', name: '商品-列表', hash: '/product/product/index', expectTab: '出售中', expectBody: ['商品'] },
    { id: 'B1b', name: '商品-分类', hash: '/product/category/index', expectBody: ['分类名称'] },
    { id: 'B1c', name: '商品-评价', hash: '/product/comment/index', expectBody: ['评价'] },
    { id: 'B1d', name: '商品-新增', hash: '/product/product/add', expectBody: ['添加商品', '出售中'] },
    // B2 订单
    { id: 'B2a', name: '订单-列表', hash: '/order/order/index', expectBody: ['订单'] },
    { id: 'B2b', name: '售后-列表', hash: '/order/refund/index', expectBody: ['进行中'] },
    // B3 会员
    { id: 'B3a', name: '会员-列表', hash: '/user/user/index', expectBody: ['关键词'] },
    { id: 'B3b', name: '会员-等级', hash: '/user/grade/index', expectTab: '等级列表', expectBody: ['等级'] },
    { id: 'B3c', name: '会员-标签', hash: '/user/tag/index', expectBody: ['标签名称'] },
    { id: 'B3d', name: '会员-余额', hash: '/user/balance/index', expectBody: ['余额'] },
    { id: 'B3e', name: '会员-积分', hash: '/user/points/index', expectBody: ['积分'] },
    { id: 'B3f', name: '会员-权益', hash: '/user/equity/index', expectBody: ['名称'] },
    // B4 门店
    { id: 'B4a', name: '门店-列表', hash: '/store/store/index', expectBody: ['门店名称'] },
    { id: 'B4b', name: '门店-店员', hash: '/store/clerk/index', expectBody: ['店员'] },
    { id: 'B4c', name: '门店-核销', hash: '/store/order/index', expectBody: ['门店'] },
    // B5 统计
    { id: 'B5a', name: '统计-销售', hash: '/statistics/sales/index', expectBody: ['成交额'] },
    { id: 'B5b', name: '统计-会员', hash: '/statistics/user/index', expectBody: ['会员'] },
    // B6 权限
    { id: 'B6a', name: '权限-角色', hash: '/auth/role/index', expectBody: ['角色'] },
    { id: 'B6b', name: '权限-管理员', hash: '/auth/user/index', expectBody: ['用户名'] },
  ];
}

async function auditCase(page, c) {
  const url = `${BASE_URL}/#${c.hash}`;
  const errors = [];
  await page.goto(url, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(3000);

  const finalUrl = page.url();
  if (finalUrl.includes('/login')) errors.push('redirected to login');

  const bodyText = await page.locator('body').innerText().catch(() => '');
  if (/404|页面不存在|Not Found/i.test(bodyText) && bodyText.length < 800) {
    errors.push('404');
  }
  if (bodyText.trim().length < 25) errors.push('blank');

  const bodyOk = c.expectBody.some((kw) => bodyText.includes(kw));
  if (!bodyOk) errors.push(`body missing: ${c.expectBody.join('|')}`);

  if (c.expectTab) {
    const tabText = (await page.locator('.el-tabs__item.is-active').first().innerText().catch(() => '')).trim();
    if (!tabText.includes(c.expectTab)) {
      errors.push(`tab expected「${c.expectTab}」got「${tabText}」`);
    }
  }

  return { ...c, url: finalUrl, status: errors.length ? 'fail' : 'pass', errors };
}

async function main() {
  const token = await apiLogin();
  const dynamic = await resolveDynamicCases(token);
  const cases = [...buildStaticCases(), ...dynamic];
  console.log('Dynamic cases:', dynamic.map((c) => c.id));

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
  for (const c of cases) {
    const result = await auditCase(page, c);
    results.push(result);
    console.log(
      `${result.status === 'pass' ? '✓' : '✗'} ${c.id} ${c.name}${result.errors.length ? ' — ' + result.errors.join('; ') : ''}`,
    );
  }
  await browser.close();

  const summary = {
    total: results.length,
    pass: results.filter((r) => r.status === 'pass').length,
    fail: results.filter((r) => r.status === 'fail').length,
    dynamicIds: dynamic.map((c) => c.id),
    failures: results.filter((r) => r.status === 'fail'),
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
