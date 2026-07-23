/**
 * Phase 3 E2E: Tab 型营销插件 M-09～M-18 — hub 可达 + 子 path / query 激活正确 Tab。
 * Run: node scripts/phase3-tab-plus-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase3-audit-report.json');

/** @type {Array<{ id: string; name: string; path: string; query?: string; expectTab?: string; expectBody: string[] }>} */
const CASES = [
  // M-09 签到
  { id: 'M-09a', name: '签到-设置', path: '/plus/sign/index', expectTab: '签到设置', expectBody: ['连续', '签到'] },
  { id: 'M-09b', name: '签到-记录', path: '/plus/sign/lists', expectTab: '签到记录', expectBody: ['昵称', '签到'] },
  // M-10 积分
  { id: 'M-10a', name: '积分-商品', path: '/plus/points/index', expectTab: '商品设置', expectBody: ['积分', '商品'] },
  { id: 'M-10b', name: '积分-兑换设置', path: '/plus/points/product/settings', expectTab: '兑换设置', expectBody: ['兑换'] },
  { id: 'M-10c', name: '积分-兑换记录', path: '/plus/points/product/record', expectTab: '兑换记录', expectBody: ['兑换'] },
  // M-11 优惠券
  { id: 'M-11a', name: '优惠券-列表', path: '/plus/coupon/index', expectTab: '优惠券列表', expectBody: ['优惠券'] },
  { id: 'M-11b', name: '优惠券-领取', path: '/plus/coupon/coupon/receive', expectTab: '领取记录', expectBody: ['领取'] },
  { id: 'M-11c', name: '优惠券-发送', path: '/plus/coupon/coupon/SendCoupon', expectTab: '发送优惠券', expectBody: ['发送', '会员'] },
  { id: 'M-11d', name: '优惠券-设置', path: '/plus/coupon/setting/index', expectTab: '基础设置', expectBody: ['设置', '图片'] },
  // M-12 满减
  { id: 'M-12a', name: '满减-活动', path: '/plus/fullreduce/index', expectBody: ['活动名称', '满减'] },
  { id: 'M-12b', name: '满减-商品', path: '/plus/fullreduce/product', expectBody: ['商品', '满减'] },
  // M-13 新人专区
  { id: 'M-13a', name: '新人-基本', path: '/plus/newactivity/index', expectTab: '基本信息', expectBody: ['活动状态', '运费'] },
  { id: 'M-13b', name: '新人-商品', path: '/plus/newactivity/index', query: 'type=product', expectTab: '活动商品', expectBody: ['商品'] },
  { id: 'M-13c', name: '新人-规则', path: '/plus/newactivity/index', query: 'type=rule', expectTab: '活动规则', expectBody: ['规则'] },
  // M-14 买送
  { id: 'M-14', name: '买送活动', path: '/plus/buyactivity/index', expectBody: ['活动', '买'] },
  // M-15 礼包
  { id: 'M-15', name: '礼包购', path: '/plus/package/index', expectBody: ['活动名称', '礼包'] },
  // M-16 邀请有礼
  { id: 'M-16a', name: '邀请-列表', path: '/plus/invitation/active/index', expectBody: ['活动', '邀请'] },
  { id: 'M-16b', name: '邀请-礼品', path: '/plus/invitation/active/receive', expectTab: '邀请人', expectBody: ['邀请'] },
  // M-17 大转盘
  { id: 'M-17a', name: '转盘-设置', path: '/plus/lottery/index', expectTab: '基础设置', expectBody: ['抽奖', '奖品'] },
  { id: 'M-17b', name: '转盘-记录', path: '/plus/lottery/record', expectTab: '抽奖记录', expectBody: ['记录', '抽奖'] },
  // M-18 文章
  { id: 'M-18a', name: '文章-列表', path: '/plus/article/index', expectTab: '文章列表', expectBody: ['文章'] },
  { id: 'M-18b', name: '文章-分类', path: '/plus/article/category', expectTab: '分类管理', expectBody: ['分类'] },
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

async function auditCase(page, c) {
  const hash = c.query ? `${c.path}?${c.query}` : c.path;
  const url = `${BASE_URL}/#${hash}`;
  const errors = [];
  await page.goto(url, { waitUntil: 'domcontentloaded' });
  await page.waitForTimeout(2800);

  const finalUrl = page.url();
  if (finalUrl.includes('/login')) errors.push('redirected to login');

  const bodyText = await page.locator('body').innerText().catch(() => '');
  if (/404|页面不存在|Not Found/i.test(bodyText) && bodyText.length < 800) {
    errors.push('404');
  }
  if (bodyText.trim().length < 30) errors.push('blank');

  const bodyOk = c.expectBody.some((kw) => bodyText.includes(kw));
  if (!bodyOk) errors.push(`body missing: ${c.expectBody.join('|')}`);

  if (c.expectTab) {
    const activeTabEl = page.locator('.el-tabs__item.is-active').first();
    const tabText = (await activeTabEl.innerText().catch(() => '')).trim();
    if (!tabText.includes(c.expectTab)) {
      errors.push(`tab expected「${c.expectTab}」got「${tabText}」`);
    }
  }

  return {
    ...c,
    url: finalUrl,
    status: errors.length ? 'fail' : 'pass',
    errors,
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
  for (const c of CASES) {
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
