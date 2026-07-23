#!/usr/bin/env node
/**
 * Live Control Center E2E — tabs, settings, products empty state, chat input.
 * Run: node scripts/e2e-live-control.mjs
 * Env: E2E_BASE_URL, E2E_API_BASE, E2E_LIVE_ID, E2E_ROOM_ID
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/live-control-e2e');
const REPORT_PATH = path.join(OUT_DIR, 'report.json');

const TABS = [
  { key: 'audience', label: '在线观众' },
  { key: 'products', label: '商品管理' },
  { key: 'redpack', label: '发送红包' },
  { key: 'audit', label: '消息审核' },
  { key: 'bot', label: '弹幕机器人' },
  { key: 'heat', label: '火力值' },
];

const report = {
  startedAt: new Date().toISOString(),
  baseUrl: BASE_URL,
  liveId: null,
  roomId: null,
  steps: [],
  pass: 0,
  fail: 0,
};

function step(name, status, note = '', extra = {}) {
  report.steps.push({ name, status, note, ...extra });
  const icon = status === 'pass' ? '✓' : status === 'fail' ? '✗' : '○';
  console.log(`${icon} ${name}${note ? ` — ${note}` : ''}`);
  if (status === 'pass') report.pass += 1;
  if (status === 'fail') report.fail += 1;
}

async function shot(page, name) {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const file = path.join(OUT_DIR, name);
  await page.screenshot({ path: file, fullPage: false });
  return file;
}

async function resolveLiveRoom(token) {
  if (process.env.E2E_LIVE_ID) {
    return {
      live_id: process.env.E2E_LIVE_ID,
      room_id: process.env.E2E_ROOM_ID || process.env.E2E_LIVE_ID,
    };
  }
  const res = await fetch(`${API_BASE}/api/v1/shop/live/list`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ page: 1, list_rows: 50 }),
  });
  const json = await res.json();
  const list = json.data?.list?.data ?? json.data?.list ?? [];
  const rows = Array.isArray(list) ? list : [];
  const named = rows.find((r) => String(r.name || '').includes('录播精选好物场'));
  const pick = named || rows[0];
  if (!pick) throw new Error('no live room in list API');
  return {
    live_id: String(pick.live_id),
    room_id: String(pick.roomid || pick.room_id || pick.live_id),
  };
}

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const login = await merchantDevLogin({ apiBase: API_BASE });
  const room = await resolveLiveRoom(login.token);
  report.liveId = room.live_id;
  report.roomId = room.room_id;

  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  const pageErrors = [];
  page.on('pageerror', (err) => pageErrors.push(err.message));

  await injectMerchantToken(page, login.token, login.user);
  const controlUrl = `${BASE_URL}/#/live/control/center?live_id=${room.live_id}&room_id=${room.room_id}`;
  await page.goto(controlUrl, { waitUntil: 'domcontentloaded', timeout: 45000 });
  await page.waitForSelector('.live-control-panel', { timeout: 30000 }).catch(() => {});
  await page.waitForTimeout(3500);
  await page.waitForSelector('.control-chat__input .el-textarea__inner, .control-chat__input textarea', {
    timeout: 15000,
  }).catch(() => {});
  for (let i = 0; i < 50; i += 1) {
    const imOk = await page.evaluate(() =>
      Boolean(document.querySelector('.control-header__im.is-connected')),
    );
    if (imOk) break;
    await page.waitForTimeout(500);
  }

  const panelOk = await page.evaluate(() => Boolean(document.querySelector('.live-control-panel')));
  step('load control panel', panelOk ? 'pass' : 'fail', panelOk ? room.live_id : 'panel missing');
  await shot(page, '01-loaded.png');

  const headerStats = await page.evaluate(() => {
    const metrics = [...document.querySelectorAll('.header-metric__value')].map((el) =>
      el.textContent?.trim(),
    );
    return metrics;
  });
  step('header metrics render', headerStats.length >= 5 ? 'pass' : 'fail', headerStats.join(', '));

  const productsUi = await page.evaluate(() => ({
    headers: document.querySelectorAll('.control-products-grid .vxe-header--row').length,
    settingsLeak: [...document.querySelectorAll('.settings-row')].filter(
      (el) => el.getBoundingClientRect().height > 0,
    ).length,
    wireframeVisible: (() => {
      const el = document.querySelector('.control-products-grid .vxe-table--empty-placeholder');
      if (!el) return false;
      const style = getComputedStyle(el);
      return style.display !== 'none' && el.getBoundingClientRect().height > 40;
    })(),
    customEmpty: Boolean(document.querySelector('.products-empty')),
  }));
  step(
    'product single header',
    productsUi.headers === 1 ? 'pass' : 'fail',
    `headers=${productsUi.headers}`,
  );
  step(
    'no settings leak',
    productsUi.settingsLeak === 0 ? 'pass' : 'fail',
    `visible=${productsUi.settingsLeak}`,
  );
  step(
    'product empty state',
    productsUi.customEmpty ? 'pass' : 'fail',
    JSON.stringify(productsUi),
  );
  await shot(page, '02-products.png');

  for (const tab of TABS) {
    await page.locator('.panel-tab').filter({ hasText: tab.label }).click();
    await page.waitForTimeout(600);
    const active = await page.evaluate((label) => {
      const btn = [...document.querySelectorAll('.panel-tab')].find((el) =>
        el.textContent?.includes(label),
      );
      return btn?.classList.contains('active');
    }, tab.label);
    step(`tab ${tab.label}`, active ? 'pass' : 'fail');
    await shot(page, `tab-${tab.key}.png`);
  }

  await page.locator('.control-header__settings-btn').click();
  await page.waitForTimeout(600);
  const settingsOpen = await page.evaluate(
    () =>
      [...document.querySelectorAll('.settings-row')].filter((el) => el.getBoundingClientRect().height > 0)
        .length >= 7,
  );
  step('settings dialog open', settingsOpen ? 'pass' : 'fail');
  await shot(page, '03-settings-open.png');
  await page.keyboard.press('Escape');
  await page.waitForTimeout(800);
  const settingsClosed = await page.evaluate(
    () => !document.querySelector('.control-settings-dialog'),
  );
  step('settings dialog close', settingsClosed ? 'pass' : 'fail');

  const chatInput = page.locator('.control-chat__input textarea, .control-chat__input .el-textarea__inner').first();
  if (await chatInput.count()) {
    const testMsg = `E2E中控弹幕${Date.now() % 100000}`;
    await chatInput.fill(testMsg);
    const sendBtn = page.locator('.control-chat__send');
    const sendDisabled = await sendBtn.isDisabled();
    step('chat send enabled', !sendDisabled ? 'pass' : 'fail', sendDisabled ? 'disabled' : 'ok');
    if (!sendDisabled) {
      await sendBtn.click();
      await page.waitForTimeout(2500);
      const sent = await page.evaluate((msg) => {
        const nodes = document.querySelectorAll('.chat-msg__text, .chat-msg__system-body');
        return [...nodes].some((el) => {
          const t = el.textContent || '';
          return t.includes(msg) || t.includes('提交审核');
        });
      }, testMsg);
      step('chat danmaku send', sent ? 'pass' : 'fail', sent ? testMsg : 'not in list');
    } else {
      step('chat danmaku send', 'fail', 'send disabled');
    }
    await shot(page, '04-chat-input.png');
  } else {
    step('chat input', 'fail', 'textarea not found');
    step('chat send enabled', 'fail', 'skipped');
    step('chat danmaku send', 'fail', 'skipped');
  }

  if (pageErrors.length) {
    const benign = pageErrors.filter((e) =>
      e.includes("Cannot read properties of undefined (reading 'type')"),
    );
    const blocking = pageErrors.filter(
      (e) => !e.includes("Cannot read properties of undefined (reading 'type')"),
    );
    if (blocking.length) {
      step('no page errors', 'fail', [...new Set(blocking)].slice(0, 3).join(' | '));
    } else if (benign.length) {
      step('no page errors', 'pass', `ignored echarts LineView (${benign.length})`);
    } else {
      step('no page errors', 'fail', [...new Set(pageErrors)].slice(0, 3).join(' | '));
    }
  } else {
    step('no page errors', 'pass');
  }

  report.finishedAt = new Date().toISOString();
  report.pageErrors = [...new Set(pageErrors)];
  fs.writeFileSync(REPORT_PATH, JSON.stringify(report, null, 2));
  console.log(`\nReport: ${REPORT_PATH}`);
  console.log(`Pass ${report.pass} / Fail ${report.fail}`);

  await browser.close();
  if (report.fail > 0) process.exit(1);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
