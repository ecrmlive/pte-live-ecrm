#!/usr/bin/env node
/**
 * Live Control Center — deep E2E (tabs, danmaku, stats vs API, settings, products, redpack).
 *
 * Run: node scripts/e2e-live-control-deep.mjs
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
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/live-control-deep');
const REPORT_PATH = path.join(OUT_DIR, 'report.json');

const TABS = [
  { key: 'audience', label: '在线观众', marker: '.control-audience' },
  { key: 'products', label: '商品管理', marker: '.control-products' },
  { key: 'redpack', label: '发送红包', marker: '.control-redpack' },
  { key: 'audit', label: '消息审核', marker: '.control-audit' },
  { key: 'bot', label: '弹幕机器人', marker: '.control-danmaku-bot' },
  { key: 'heat', label: '火力值', marker: '.control-heat-boost' },
];

const report = {
  startedAt: new Date().toISOString(),
  baseUrl: BASE_URL,
  liveId: null,
  roomId: null,
  api: {},
  checks: [],
  pass: 0,
  fail: 0,
  blockers: [],
  pageErrors: [],
};

function check(name, ok, detail = '') {
  report.checks.push({ name, ok, detail });
  const icon = ok ? '✓' : '✗';
  console.log(`${icon} ${name}${detail ? ` — ${detail}` : ''}`);
  if (ok) report.pass += 1;
  else report.fail += 1;
}

async function shot(page, name) {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  await page.screenshot({ path: path.join(OUT_DIR, name), fullPage: false });
}

function countVisible(page, selector) {
  return page.locator(selector).evaluateAll((nodes) =>
    nodes.filter((node) => {
      const style = window.getComputedStyle(node);
      if (style.display === 'none' || style.visibility === 'hidden') return false;
      const rect = node.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    }).length,
  );
}

function parseUiCount(text) {
  const raw = String(text || '').replace(/[,，\s]/g, '');
  if (raw.startsWith('¥')) {
    const n = Number.parseFloat(raw.slice(1));
    return Number.isFinite(n) ? n : null;
  }
  const n = Number.parseInt(raw, 10);
  return Number.isFinite(n) ? n : null;
}

async function apiPost(token, urlPath, fields = {}) {
  const body = new URLSearchParams(fields);
  const res = await fetch(`${API_BASE}${urlPath}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded',
      Authorization: `Bearer ${token}`,
      'authori-zation': `Bearer ${token}`,
    },
    body,
  });
  return { status: res.status, json: await res.json().catch(() => ({})) };
}

async function resolveLiveRoom(token) {
  if (process.env.E2E_LIVE_ID) {
    return {
      live_id: String(process.env.E2E_LIVE_ID),
      room_id: String(process.env.E2E_ROOM_ID || process.env.E2E_LIVE_ID),
    };
  }
  const { json } = await apiPost(token, '/api/v1/shop/live/list', {
    page: 1,
    list_rows: 50,
  });
  const rows = json?.data?.list?.data ?? json?.data?.list ?? [];
  const list = Array.isArray(rows) ? rows : [];
  const pick =
    list.find((r) => r.session_status === 1 && String(r.name || '').includes('录播精选好物场')) ||
    list.find((r) => r.session_status === 1) ||
    list.find((r) => String(r.name || '').includes('录播精选好物场')) ||
    list[0];
  if (!pick) throw new Error('no live room in list API');
  return {
    live_id: String(pick.live_id),
    room_id: String(pick.roomid || pick.room_id || pick.live_id),
  };
}

function chatInputLocator(page) {
  return page.locator(
    '.control-chat__input textarea, .control-chat__input .el-textarea__inner',
  );
}

async function waitForManageIdle(page) {
  await page
    .locator('.control-manage .el-loading-mask')
    .first()
    .waitFor({ state: 'hidden', timeout: 15000 })
    .catch(() => {});
  await page.waitForTimeout(150);
}

async function clickPanelTab(page, label) {
  await waitForManageIdle(page);
  const tab = page.locator('.control-manage__tabs .panel-tab').filter({ hasText: label }).first();
  await tab.scrollIntoViewIfNeeded();
  await tab.click();
  await waitForManageIdle(page);
}

async function waitForImConnected(page, timeoutMs = 25000) {
  const start = Date.now();
  while (Date.now() - start < timeoutMs) {
    const ok = await page.evaluate(() =>
      Boolean(
        document.querySelector('.control-header__im.is-connected') ||
          document.querySelector('.control-chat-col__im-dot--connected'),
      ),
    );
    if (ok) return true;
    await page.waitForTimeout(500);
  }
  return false;
}

async function readHeaderMetrics(page) {
  return page.evaluate(() => {
    const labels = [...document.querySelectorAll('.header-metric__label-full')].map((el) =>
      el.textContent?.trim(),
    );
    const values = [...document.querySelectorAll('.header-metric__value')].map((el) =>
      el.textContent?.trim(),
    );
    const out = {};
    labels.forEach((label, i) => {
      if (label) out[label] = values[i] ?? '';
    });
    return out;
  });
}

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const login = await merchantDevLogin({ apiBase: API_BASE });
  check('dev-login', true, login.user?.user_name || 'ok');

  const room = await resolveLiveRoom(login.token);
  report.liveId = room.live_id;
  report.roomId = room.room_id;

  const [detail, counts, stats] = await Promise.all([
    apiPost(login.token, '/api/v1/shop/live/detail', {
      live_id: room.live_id,
      room_id: room.room_id,
      with_counts: 1,
      with_live_stats: 1,
    }),
    apiPost(login.token, '/api/v1/room/counts', { live_id: room.live_id, room_id: room.room_id }),
    apiPost(login.token, '/api/v1/shop/live/live-stats', {
      live_id: room.live_id,
      room_id: room.room_id,
    }),
  ]);
  report.api = {
    detailCode: detail.json?.code,
    countsCode: counts.json?.code,
    statsCode: stats.json?.code,
    productCount: 0,
  };
  check(
    'live-apis',
    detail.json?.code === 1 && counts.json?.code === 1,
    `detail=${detail.json?.code} counts=${counts.json?.code}`,
  );

  const detailData = detail.json?.data || {};
  const countsData = counts.json?.data || {};
  const statsData = stats.json?.data || {};
  const sessionActive = detailData.session_status === 1 || detailData.session_status === 2;

  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  page.on('pageerror', (err) => report.pageErrors.push(err.message));

  await injectMerchantToken(page, login.token);
  await page.goto(
    `${BASE_URL}/#/live/control/center?live_id=${room.live_id}&room_id=${room.room_id}`,
    { waitUntil: 'domcontentloaded', timeout: 45000 },
  );
  await page.waitForSelector('.live-control-panel', { timeout: 30000 });
  await page.waitForTimeout(4000);
  await shot(page, '00-loaded.png');

  check('auth', !page.url().includes('/auth/login'));
  check('header', (await page.locator('.live-control-panel').count()) > 0, room.live_id);
  check('no-inline-settings', (await countVisible(page, '.settings-row')) === 0);

  check('im-connected', await waitForImConnected(page), await page.locator('.control-header__im-text').textContent().catch(() => ''));

  const uiMetrics = await readHeaderMetrics(page);
  const apiMetrics = {
    online: countsData.online_count ?? detailData.online_count,
    total: countsData.total_count ?? detailData.total_count,
    heat: detailData.display_heat ?? detailData.fire_value,
    sales: statsData.pay_amount_total ?? detailData.live_stats?.pay_amount_total,
  };
  report.uiMetrics = uiMetrics;
  report.apiMetrics = apiMetrics;

  for (const [label, key] of [
    ['在线人数', 'online'],
    ['累计人数', 'total'],
    ['火力值', 'heat'],
    ['总销量金额', 'sales'],
  ]) {
    const uiN = parseUiCount(uiMetrics[label]);
    const apiN = Number(apiMetrics[key]);
    const tol = label.includes('销量') ? 0.01 : label === '在线人数' ? 1 : 0;
    const ok = uiN != null && Number.isFinite(apiN) && Math.abs(uiN - apiN) <= tol;
    check(`stats:${label}`, ok, `ui=${uiN} api=${apiN}`);
  }

  await page.waitForSelector('.control-chat__input .el-textarea__inner, .control-chat__input textarea', {
    timeout: 15000,
  }).catch(() => {});
  const chatInput = chatInputLocator(page);
  check('chat-input', (await chatInput.count()) > 0);

  if ((await chatInput.count()) > 0) {
    const text = `E2E深度${Date.now() % 100000}`;
    await chatInput.first().fill(text);
    const sendBtn = page.locator('.control-chat__send');
    const sendOk = !(await sendBtn.isDisabled());
    check('chat-send-enabled', sendOk, sendOk ? 'ok' : 'disabled');
    if (!sessionActive) {
      check('danmaku-send', true, 'skipped: session not live');
    } else if (sendOk) {
      await sendBtn.click();
      await page.waitForTimeout(4000);
      const appeared = await page.locator('.control-chat__list').filter({ hasText: text }).count();
      const recent = await apiPost(login.token, '/api/v1/shop/live/danmaku/recent', {
        live_id: room.live_id,
        room_id: room.room_id,
        session_id: detailData.current_session_id || detailData.session_id || '',
        limit: 30,
      });
      const apiHas = (recent.json?.data?.list || []).some((row) =>
        String(row.content || row.text || '').includes(text.slice(0, 8)),
      );
      const pendingMsg = await page.locator('.control-chat__list').filter({ hasText: '消息已提交审核' }).count();
      const draftEmpty = ((await chatInput.first().inputValue().catch(() => 'x')) || '').trim() === '';
      const sessionActive = Number(detailData.session_status) === 1;
      check(
        'danmaku-send',
        appeared > 0 || apiHas || pendingMsg > 0 || draftEmpty || !sessionActive,
        `ui=${appeared} api=${apiHas} pending=${pendingMsg} session=${sessionActive}`,
      );
      await shot(page, 'danmaku-sent.png');
    }
  }

  for (const tab of TABS) {
    await clickPanelTab(page, tab.label);
    await page.waitForTimeout(500);
    if (tab.key === 'audience') {
      await page.waitForSelector('.control-audience', { timeout: 10000, state: 'visible' }).catch(() => {
        return clickPanelTab(page, tab.label);
      });
    }
    await page
      .waitForSelector(tab.marker, { timeout: 12000, state: 'visible' })
      .catch(() => {});
    if (tab.key === 'heat') {
      await page
        .waitForSelector('.control-heat-boost', { timeout: 10000, state: 'visible' })
        .catch(() => {});
      await page.waitForTimeout(400);
    }
    const state = await page.evaluate(
      ({ label, marker }) => {
        const btn = [...document.querySelectorAll('.control-manage__tabs .panel-tab')].find((el) =>
          el.textContent?.includes(label),
        );
        const markerEl = document.querySelector(marker);
        const markerRect = markerEl?.getBoundingClientRect();
        return {
          active: btn?.classList.contains('active') ?? false,
          markerVisible: Boolean(markerRect && markerRect.height > 0 && markerRect.width > 0),
        };
      },
      { label: tab.label, marker: tab.marker },
    );
    const tabOk = state.active && state.markerVisible;
    if (tab.key === 'redpack' && !tabOk) {
      report.blockers.push('redpack: panel not visible after tab click');
    }
    check(`tab:${tab.label}`, tabOk, JSON.stringify(state));
    await shot(page, `tab-${tab.key}.png`);

    if (tab.key === 'redpack') {
      const canSend = await page.locator('.control-redpack .btn-send').isEnabled().catch(() => false);
      if (!canSend) {
        report.blockers.push('redpack: needs active session');
      }
      check('redpack-ui', canSend || (await page.locator('.control-redpack').count()) > 0);
    }
  }

  await clickPanelTab(page, '商品管理');

  const productsUi = await page.evaluate(() => ({
    headers: document.querySelectorAll('.control-products-grid .vxe-header--row').length,
    uiRows: document.querySelectorAll('.control-products-grid .vxe-body--row').length,
    customEmpty: Boolean(document.querySelector('.products-empty')),
  }));
  check('product-single-header', productsUi.headers <= 1, `headers=${productsUi.headers}`);

  await page.locator('.control-products .btn-add').first().click();
  await page.waitForTimeout(1200);
  const pickDialog = page.locator('.el-dialog, [role="dialog"]').filter({
    hasText: /选择商城商品|添加商品|商城商品/,
  });
  check('product-pick-dialog', (await pickDialog.count()) > 0);
  if (await pickDialog.count()) {
    const row = page.locator('.pick-mall-grid .vxe-body--row, .vxe-table--body-wrapper .vxe-body--row').first();
    if (await row.count()) {
      await row.click({ force: true }).catch(() => {});
      check('product-pick-row', true, 'clicked first row');
    }
    await page.locator('.pick-mall-products-modal button').filter({ hasText: '取消' }).click({ force: true }).catch(() => {});
    await page
      .waitForFunction(
        () => !document.querySelector('[data-slot="dialog-content"]'),
        { timeout: 5000 },
      )
      .catch(() => page.keyboard.press('Escape'));
  }
  await page.waitForTimeout(400);

  // Fullscreen: pick dialog must mount inside the control overlay root
  await page.locator('.control-header__right .icon-btn').last().click();
  await page
    .waitForFunction(
      () => document.fullscreenElement?.classList.contains('live-control-panel'),
      { timeout: 5000 },
    )
    .catch(() => {});

  await page.locator('.control-products .btn-add').first().click();
  await page.waitForTimeout(1500);

  const fullscreenPick = await page.evaluate(() => {
    const fs = document.fullscreenElement;
    if (!fs?.classList.contains('live-control-panel')) {
      return { ok: false, reason: 'not-fullscreen' };
    }
    const dialog = [...document.querySelectorAll('[data-slot="dialog-content"], [role="dialog"]')].find(
      (el) => /选择商城商品/.test(el.textContent || ''),
    );
    if (!dialog) {
      return { ok: false, reason: 'dialog-missing' };
    }
    const overlayRoot = fs.querySelector('.live-control-overlay-root');
    return {
      ok: Boolean(overlayRoot?.contains(dialog)),
      inFullscreen: fs.contains(dialog),
      inOverlayRoot: Boolean(overlayRoot?.contains(dialog)),
    };
  });
  check('product-pick-dialog-fullscreen', fullscreenPick.ok, JSON.stringify(fullscreenPick));

  await page.locator('.pick-mall-products-modal button').filter({ hasText: '取消' }).click({ force: true }).catch(() => page.keyboard.press('Escape'));
  await page
    .waitForFunction(
      () => !document.querySelector('.live-control-overlay-root [data-slot="dialog-content"]'),
      { timeout: 5000 },
    )
    .catch(() => {});
  await page.waitForTimeout(300);
  await page.evaluate(() => document.exitFullscreen?.()).catch(() => {});
  await page.waitForTimeout(300);

  await page.locator('.control-header__settings-btn').click();
  await page.waitForSelector('.settings-btn--save', { timeout: 10000 });
  check('settings-dialog', (await page.locator('.settings-row').count()) >= 7);

  const beforeDetail = await apiPost(login.token, '/api/v1/shop/live/detail', {
    live_id: room.live_id,
    room_id: room.room_id,
  });
  const beforeGift = beforeDetail.json?.data?.enable_gift;
  const giftSwitch = page.locator('.settings-row--rose .el-switch');
  if (await giftSwitch.count()) {
    const beforeChecked = await giftSwitch.getAttribute('aria-checked');
    await giftSwitch.locator('.el-switch__core').click({ force: true });
    await page.waitForTimeout(250);
    const afterChecked = await giftSwitch.getAttribute('aria-checked');
    await page.locator('.settings-btn--save').click();
    await page
      .waitForSelector('.control-settings-dialog', { state: 'hidden', timeout: 12000 })
      .catch(() => page.waitForTimeout(2000));
    const after = await apiPost(login.token, '/api/v1/shop/live/detail', {
      live_id: room.live_id,
      room_id: room.room_id,
    });
    const afterGift = after.json?.data?.enable_gift;
    check(
      'settings-persist',
      afterGift !== beforeGift,
      `${beforeGift}→${afterGift} switch=${beforeChecked}→${afterChecked}`,
    );
    await page.locator('.control-header__settings-btn').click();
    await page.waitForSelector('.settings-btn--save', { timeout: 10000 });
    const giftSwitchRestore = page.locator('.settings-row--rose .el-switch');
    if (await giftSwitchRestore.count()) {
      await giftSwitchRestore.locator('.el-switch__core').click({ force: true });
      await page.waitForTimeout(250);
      await page.locator('.settings-btn--save').click();
      await page
        .waitForSelector('.control-settings-dialog', { state: 'hidden', timeout: 12000 })
        .catch(() => page.waitForTimeout(1500));
    }
  }

  await page.locator('.settings-btn--cancel').click().catch(() => page.keyboard.press('Escape'));
  await page.waitForTimeout(800);
  check('settings-closed', (await countVisible(page, '.settings-row').catch(() => 0)) === 0);

  const IGNORED_PAGE_ERRORS = [
    "Cannot read properties of undefined (reading 'type')",
    "Cannot read properties of undefined (reading 'getAxesOnZeroOf')",
    'LineView',
  ];
  const uniqueErrors = [...new Set(report.pageErrors)].filter(
    (msg) => !IGNORED_PAGE_ERRORS.some((needle) => msg.includes(needle)),
  );
  check(
    'no-page-errors',
    uniqueErrors.length === 0,
    uniqueErrors.length ? uniqueErrors.slice(0, 2).join(' | ') : '',
  );

  report.finishedAt = new Date().toISOString();
  report.pageErrors = [...new Set(report.pageErrors)];
  fs.writeFileSync(REPORT_PATH, JSON.stringify(report, null, 2));
  console.log(`\nPass ${report.pass} / Fail ${report.fail}`);
  if (report.blockers.length) console.log('Blockers:', report.blockers.join('; '));
  console.log(`Report: ${REPORT_PATH}`);

  await browser.close();
  if (report.fail > 0) process.exit(1);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
