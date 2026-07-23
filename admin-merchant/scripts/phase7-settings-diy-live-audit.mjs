/**
 * Phase 7 E2E: 设置 / 装修 DIY / 直播 / App / 财务。
 * Run: node scripts/phase7-settings-diy-live-audit.mjs
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
const OUT = path.join(__dirname, '../.e2e-screenshots/phase7-audit-report.json');

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

  const deliveries = await apiPost(token, '/shop/setting.delivery/index', {
    page: 1,
    list_rows: 1,
  });
  const deliveryId = deliveries.data?.list?.data?.[0]?.delivery_id;
  if (deliveryId) {
    dynamic.push({
      id: 'S2a',
      name: '运费-编辑弹窗',
      hash: `/setting/delivery/edit?delivery_id=${deliveryId}`,
      expectBody: ['模板名称', '添加运费模版', '编辑运费模版'],
    });
  }

  const addresses = await apiPost(token, '/shop/setting.address/index', {
    page: 1,
    list_rows: 1,
  });
  const addressId = addresses.data?.list?.data?.[0]?.address_id;
  if (addressId) {
    dynamic.push({
      id: 'S2b',
      name: '地址-编辑弹窗',
      hash: `/setting/address/edit?address_id=${addressId}`,
      expectBody: ['收货人姓名', '添加', '编辑'],
    });
  }

  dynamic.push({
    id: 'L2a',
    name: '直播间-创建弹窗',
    hash: '/live/room/add',
    expectBody: ['直播间', '直播间名称'],
  });

  const lives = await apiPost(token, '/api/v1/shop/live/list', {
    page: 1,
    list_rows: 1,
  });
  const liveId =
    lives.data?.list?.data?.[0]?.live_id ??
    lives.data?.list?.data?.[0]?.id ??
    lives.list?.data?.[0]?.live_id;
  if (liveId) {
    dynamic.push({
      id: 'L2b',
      name: '直播间-编辑弹窗',
      hash: `/live/room/edit?live_id=${liveId}`,
      expectBody: ['直播间', '直播间名称'],
    });
  }

  return dynamic;
}

/** @type {Array<{ id: string; name: string; hash: string; expectTab?: string; expectBody: string[] }>} */
function buildStaticCases() {
  return [
    // S 设置
    {
      id: 'S1a',
      name: '设置-商城',
      hash: '/setting/store/index',
      expectTab: '平台设置',
      expectBody: ['平台设置', '会员设置'],
    },
    {
      id: 'S1b',
      name: '设置-交易',
      hash: '/setting/trade/index',
      expectTab: '订单流程设置',
      expectBody: ['订单流程', '运费设置'],
    },
    { id: 'S1c', name: '设置-运费模板', hash: '/setting/delivery/index', expectBody: ['模板名称'] },
    { id: 'S1d', name: '设置-物流公司', hash: '/setting/express/index', expectBody: ['物流公司'] },
    { id: 'S1e', name: '设置-消息', hash: '/setting/message/index', expectBody: ['消息'] },
    { id: 'S1f', name: '设置-短信', hash: '/setting/sms/index', expectBody: ['短信'] },
    { id: 'S1g', name: '设置-地址', hash: '/setting/address/index', expectBody: ['收货人姓名'] },
    { id: 'S1h', name: '设置-打印机', hash: '/setting/printer/index', expectBody: ['打印机'] },
    { id: 'S1i', name: '设置-打印', hash: '/setting/printing/index', expectBody: ['小票打印'] },
    { id: 'S1j', name: '设置-协议', hash: '/setting/protocol/index', expectBody: ['协议'] },
    { id: 'S1k', name: '设置-客服', hash: '/setting/mpservice/index', expectBody: ['客服'] },
    { id: 'S1l', name: '设置-清理缓存', hash: '/setting/clear/index', expectBody: ['清理缓存'] },
    { id: 'S1m', name: '设置-支付', hash: '/appsetting/app/pay', expectBody: ['支付'] },

    // P 装修
    { id: 'P1a', name: '装修-首页列表', hash: '/page/page/list', expectBody: ['页面名称'] },
    { id: 'P1b', name: '装修-页面列表', hash: '/page/page/index', expectBody: ['页面名称'] },
    { id: 'P1c', name: '装修-分类模板', hash: '/page/page/category', expectBody: ['分类'] },
    { id: 'P1d', name: '装修-主题', hash: '/page/theme/index', expectBody: ['主题'] },
    { id: 'P1e', name: '装修-底部导航', hash: '/page/tabbar/index', expectBody: ['底部导航'] },
    { id: 'P1f', name: '装修-个人中心', hash: '/page/center/index', expectBody: ['个人中心'] },
    { id: 'P1g', name: '装修-添加页面', hash: '/page/page/add', expectBody: ['组件库'] },

    // L 直播
    { id: 'L1a', name: '直播-列表', hash: '/live/index', expectBody: ['直播间'] },
    { id: 'L1b', name: '直播-主播', hash: '/live/anchor/index', expectBody: ['主播'] },
    { id: 'L1c', name: '直播-记录', hash: '/live/session/index', expectBody: ['关键词'] },
    { id: 'L1d', name: '直播-统计', hash: '/live/stat/index', expectBody: ['直播统计'] },
    { id: 'L1e', name: '直播-流量统计', hash: '/live/traffic/account', expectBody: ['流量'] },
    { id: 'L1f', name: '直播-流量明细', hash: '/live/traffic/session', expectBody: ['流量'] },
    { id: 'L1g', name: '直播-充值记录', hash: '/live/traffic/recharge/index', expectBody: ['流量'] },
    { id: 'L1h', name: '直播-域名', hash: '/live/h5domain/index', expectBody: ['域名'] },
    { id: 'L1i', name: '直播-投诉', hash: '/live/complaint/index', expectBody: ['投诉'] },
    { id: 'L1j', name: '直播-素材', hash: '/live/material/index', expectBody: ['素材'] },
    { id: 'L1k', name: '直播-话术', hash: '/live/danmaku-bot/script/index', expectBody: ['话术'] },
    { id: 'L1l', name: '直播-机器人', hash: '/live/danmaku-bot/robot/index', expectBody: ['机器人'] },
    { id: 'L1m', name: '直播-敏感词', hash: '/live/sensitive-word/index', expectBody: ['敏感词'] },

    // A 应用
    { id: 'A1a', name: '应用-小程序', hash: '/appsetting/appwx/index', expectBody: ['小程序'] },
    { id: 'A1b', name: '应用-公众号', hash: '/appsetting/appmp/index', expectBody: ['公众号'] },
    {
      id: 'A1c',
      name: '应用-App Hub',
      hash: '/appsetting/appopen/event',
      expectTab: '基础设置',
      expectBody: ['App', '基础设置'],
    },
    {
      id: 'A1d',
      name: '应用-分享 Tab',
      hash: '/appsetting/appshare/index',
      expectTab: '分享设置',
      expectBody: ['分享'],
    },
    {
      id: 'A1e',
      name: '应用-升级 Tab',
      hash: '/appsetting/appupdate/index',
      expectTab: '升级管理',
      expectBody: ['升级'],
    },

    // F 财务
    { id: 'F1a', name: '财务-概况', hash: '/finance/financeSituation', expectBody: ['订单概况'] },
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
