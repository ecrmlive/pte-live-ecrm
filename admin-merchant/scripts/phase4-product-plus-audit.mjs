/**
 * Phase 4 E2E: 商品型营销 M-19～M-22 — hub Tab · add/edit · statistics。
 * Run: node scripts/phase4-product-plus-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase4-audit-report.json');

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

async function apiPost(token, path, body = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
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

/** @type {Array<{ id: string; name: string; hash: string; expectTab?: string; expectBody: string[] }>} */
function buildCases(stats) {
  const cases = [
    // M-19 秒杀 hub
    { id: 'M-19a', name: '秒杀-活动', hash: '/plus/seckill/index', expectTab: '秒杀活动', expectBody: ['秒杀', '活动'] },
    { id: 'M-19b', name: '秒杀-商品', hash: '/plus/seckill/product/index', expectTab: '秒杀商品', expectBody: ['商品'] },
    { id: 'M-19c', name: '秒杀-时段', hash: '/plus/seckill/time/index', expectTab: '秒杀配置', expectBody: ['时段', '秒杀'] },
    { id: 'M-19d', name: '秒杀-设置', hash: '/plus/seckill/setting/index', expectTab: '基础设置', expectBody: ['设置'] },
    { id: 'M-19e', name: '秒杀-新增活动', hash: '/plus/seckill/active/add', expectBody: ['基础设置', '商品设置'] },
    // M-20 拼团
    { id: 'M-20a', name: '拼团-商品', hash: '/plus/assemble/index', expectTab: '拼团商品', expectBody: ['拼团'] },
    { id: 'M-20b', name: '拼团-列表', hash: '/plus/assemble/record/index', expectTab: '拼团列表', expectBody: ['拼团'] },
    { id: 'M-20c', name: '拼团-设置', hash: '/plus/assemble/setting/index', expectTab: '基础设置', expectBody: ['设置'] },
    { id: 'M-20d', name: '拼团-新增商品', hash: '/plus/assemble/product/add', expectBody: ['基础设置', '规格库存', '拼团详情'] },
    // M-21 砍价
    { id: 'M-21a', name: '砍价-商品', hash: '/plus/bargain/index', expectTab: '砍价商品', expectBody: ['砍价'] },
    { id: 'M-21b', name: '砍价-列表', hash: '/plus/bargain/task/index', expectTab: '砍价列表', expectBody: ['砍价'] },
    { id: 'M-21c', name: '砍价-设置', hash: '/plus/bargain/setting/index', expectTab: '砍价设置', expectBody: ['设置'] },
    { id: 'M-21d', name: '砍价-新增商品', hash: '/plus/bargain/product/add', expectBody: ['砍价', '选择商品'] },
    // M-22 预售
    { id: 'M-22a', name: '预售-商品', hash: '/plus/advance/index', expectTab: '预售商品', expectBody: ['预售'] },
    { id: 'M-22b', name: '预售-设置', hash: '/plus/advance/setting/index', expectTab: '基础设置', expectBody: ['设置'] },
    // add 权限路由 → hub + 商品选择弹窗（与拼团 M-20d 同模式，hub Tab 文案可命中）
    { id: 'M-22c', name: '预售-新增商品', hash: '/plus/advance/product/add', expectBody: ['预售', '选择商品'] },
  ];

  if (stats.seckill) {
    cases.push({
      id: 'M-19f',
      name: '秒杀-统计',
      hash: stats.seckill,
      expectBody: ['返回', '参与'],
    });
  }
  if (stats.assemble) {
    cases.push({
      id: 'M-20e',
      name: '拼团-统计',
      hash: stats.assemble,
      expectBody: ['返回', '参与'],
    });
  }
  if (stats.bargain) {
    cases.push({
      id: 'M-21e',
      name: '砍价-统计',
      hash: stats.bargain,
      expectBody: ['返回', '参与'],
    });
  }

  return cases;
}

async function resolveStatisticsHashes(token) {
  const stats = { seckill: null, assemble: null, bargain: null };

  const seckill = await apiPost(token, '/shop/plus.seckill.product/index', {
    page: 1,
    list_rows: 1,
  });
  const sRow = seckill.data?.list?.data?.[0];
  if (sRow) {
    const q = new URLSearchParams({
      product_id: String(sRow.product_id ?? sRow.productId ?? ''),
      seckill_product_id: String(sRow.seckill_product_id ?? sRow.seckillProductId ?? ''),
      seckill_activity_id: String(sRow.seckill_activity_id ?? sRow.seckillActivityId ?? ''),
    });
    stats.seckill = `/plus/seckill/product/statistics?${q}`;
  }

  const assemble = await apiPost(token, '/shop/plus.assemble.product/index', {
    page: 1,
    list_rows: 1,
  });
  const aRow = assemble.data?.list?.data?.[0];
  if (aRow) {
    const q = new URLSearchParams({
      product_id: String(aRow.product_id ?? ''),
      assemble_product_id: String(aRow.assemble_product_id ?? aRow.assembleProductId ?? ''),
    });
    stats.assemble = `/plus/assemble/product/statistics?${q}`;
  }

  const bargain = await apiPost(token, '/shop/plus.bargain.product/index', {
    page: 1,
    list_rows: 1,
  });
  const bRow = bargain.data?.list?.data?.[0];
  if (bRow) {
    const q = new URLSearchParams({
      product_id: String(bRow.product_id ?? ''),
      bargain_product_id: String(bRow.bargain_product_id ?? bRow.bargainProductId ?? ''),
    });
    stats.bargain = `/plus/bargain/product/statistics?${q}`;
  }

  return stats;
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
  const stats = await resolveStatisticsHashes(token);
  const cases = buildCases(stats);

  console.log('Statistics URLs:', stats);

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
    statsUrls: stats,
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
