/**
 * Phase 5 E2E: 重型 Plus M-23～M-27 — 分销 / 卡券 / 表单 / 面单 / 直播插件。
 * Run: node scripts/phase5-heavy-plus-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase5-audit-report.json');

/** @type {Array<{ id: string; name: string; hash: string; expectTab?: string; expectBody: string[]; allowRedirect?: boolean }>} */
const CASES = [
  // M-23 分销
  { id: 'M-23a', name: '分销-入驻', hash: '/plus/agent/index?type=apply', expectTab: '入驻申请', expectBody: ['申请', '分销'] },
  { id: 'M-23b', name: '分销-用户', hash: '/plus/agent/user/index', expectTab: '分销商用户', expectBody: ['分销'], allowRedirect: true },
  { id: 'M-23c', name: '分销-等级', hash: '/plus/agent/grade/index', expectTab: '分销商等级', expectBody: ['等级'], allowRedirect: true },
  { id: 'M-23d', name: '分销-商品', hash: '/plus/agent/product/index', expectTab: '分销商品', expectBody: ['商品'], allowRedirect: true },
  { id: 'M-23e', name: '分销-订单', hash: '/plus/agent/order/index', expectTab: '分销订单', expectBody: ['订单'], allowRedirect: true },
  { id: 'M-23f', name: '分销-提现', hash: '/plus/agent/cash/index', expectTab: '提现申请', expectBody: ['提现'], allowRedirect: true },
  { id: 'M-23g', name: '分销-设置', hash: '/plus/agent/setting/index', expectTab: '分销设置', expectBody: ['设置'], allowRedirect: true },
  { id: 'M-23h', name: '分销-海报', hash: '/plus/agent/poster/index', expectTab: '分销海报', expectBody: ['海报'], allowRedirect: true },
  // M-24 卡券
  { id: 'M-24a', name: '卡券-管理', hash: '/plus/card/event', expectTab: '卡券管理', expectBody: ['卡券'] },
  { id: 'M-24b', name: '卡券-分类', hash: '/plus/card/category/index', expectTab: '分类管理', expectBody: ['分类'] },
  { id: 'M-24c', name: '卡券-提货码', hash: '/plus/card/code/index', expectTab: '提货码管理', expectBody: ['提货'] },
  { id: 'M-24d', name: '卡券-订单', hash: '/plus/card/order/index', expectTab: '提货订单', expectBody: ['订单'] },
  { id: 'M-24e', name: '卡券-设置', hash: '/plus/card/setting/index', expectTab: '卡券设置', expectBody: ['设置'] },
  // M-25 万能表单
  { id: 'M-25a', name: '表单-管理', hash: '/plus/table/event', expectTab: '表单管理', expectBody: ['表单'] },
  { id: 'M-25b', name: '表单-记录', hash: '/plus/table/record/index', expectTab: '表单记录', expectBody: ['记录'] },
  { id: 'M-25c', name: '表单-新增', hash: '/plus/table/table/add', expectBody: ['表单', '字段'], allowRedirect: true },
  // M-26 面单
  { id: 'M-26a', name: '面单-模板', hash: '/plus/surface/index', expectTab: '面单模板', expectBody: ['模板', '面单'] },
  { id: 'M-26b', name: '面单-配置', hash: '/plus/surface/setting/index', expectTab: '面单配置', expectBody: ['配置', '面单'] },
  { id: 'M-26c', name: '面单-新增模板', hash: '/plus/surface/template/add', expectBody: ['模板', '面单'], allowRedirect: true },
  // M-27 直播插件
  { id: 'M-27a', name: '直播-房间', hash: '/plus/live/wx/index', expectTab: '直播房间', expectBody: ['直播'] },
  { id: 'M-27b', name: '直播-主播', hash: '/plus/live/anchor/index', expectTab: '主播管理', expectBody: ['主播'] },
  { id: 'M-27c', name: '直播-商品', hash: '/plus/live/product/index', expectTab: '直播商品', expectBody: ['商品'] },
  { id: 'M-27d', name: '直播-卡片别名', hash: '/plus/live/index', expectTab: '直播房间', expectBody: ['直播'], allowRedirect: true },
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
