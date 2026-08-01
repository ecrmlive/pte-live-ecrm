#!/usr/bin/env node
/**
 * Full live room create E2E: open modal, verify fields, submit, verify list.
 * Run: node scripts/e2e-live-create-submit.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT = path.join(__dirname, '../.e2e-screenshots/live-create-test');
const SEED_LOG = path.join(__dirname, 'e2e-seed-log.json');

export function futureStartTime() {
  const future = new Date(Date.now() + 10 * 60 * 1000);
  const pad = (n) => String(n).padStart(2, '0');
  return `${future.getFullYear()}-${pad(future.getMonth() + 1)}-${pad(future.getDate())} ${pad(future.getHours())}:${pad(future.getMinutes())}:${pad(future.getSeconds())}`;
}

export async function fetchFirstAnchor(token, apiBase = API_BASE) {
  const res = await fetch(`${apiBase}/api/v1/shop/live/anchor/list`, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ list_rows: 5, page: 1, status: 1 }),
  });
  const json = await res.json();
  const row = json?.data?.list?.data?.[0];
  if (!row?.anchor_id) throw new Error('no anchor in seed data');
  return row;
}

export function createModal(page) {
  return page.locator('[role="dialog"]').filter({ hasText: '创建直播间' }).last();
}

export async function dismissPopovers(page) {
  for (let i = 0; i < 3; i += 1) {
    await page.keyboard.press('Escape');
    await page.waitForTimeout(150);
  }
}

async function setStartTime(page, modal, value) {
  const dt = modal.locator('input[placeholder="选择时间"]');
  await dt.click();
  await page.waitForTimeout(400);
  const shortcut = page
    .locator('.el-picker-panel:visible .el-picker-panel__shortcut')
    .filter({ hasText: '此刻' });
  if (await shortcut.isVisible().catch(() => false)) {
    await shortcut.click();
  } else {
    await dt.fill(value);
    await dt.press('Tab');
  }
  await dismissPopovers(page);
}

async function selectAnchor(page, modal) {
  const select = modal.locator('.anchor-select');
  await select.scrollIntoViewIfNeeded();
  await select.locator('.el-select__wrapper').click();
  await page.waitForTimeout(600);
  const option = page
    .locator('.native-modal-popper .el-select-dropdown__item, .el-select-dropdown .el-select-dropdown__item')
    .first();
  await option.waitFor({ state: 'attached', timeout: 8000 });
  await option.click();
  await page.waitForTimeout(300);
}

export async function patchLiveForm(page, payload) {
  return page.evaluate((data) => {
    const panel = document.querySelector('.live-room-add-panel');
    if (!panel) return false;
    let comp = panel.__vueParentComponent;
    while (comp) {
      const state = comp.setupState || comp.ctx;
      const form = state?.form;
      if (form && 'anchor_id' in form) {
        if (data.name) form.name = data.name;
        if (data.start_time) form.start_time = data.start_time;
        if (data.anchor_id) {
          form.anchor_id = data.anchor_id;
          if (typeof state.onAnchorChange === 'function') state.onAnchorChange(data.anchor_id);
        }
        if (data.cover_img) form.cover_img = data.cover_img;
        if (data.background_img) form.background_img = data.background_img;
        return true;
      }
      comp = comp.parent;
    }
    return false;
  }, payload);
}

export function appendSeedLog(entry) {
  let data = { entries: [] };
  try {
    data = JSON.parse(fs.readFileSync(SEED_LOG, 'utf8'));
  } catch {
    /* new file */
  }
  data.entries.push({ at: new Date().toISOString(), ...entry });
  fs.writeFileSync(SEED_LOG, `${JSON.stringify(data, null, 2)}\n`);
}

async function pickImage(page) {
  const picker = page.locator('.el-dialog').filter({ hasText: /选择图片|素材库|图片/ }).last();
  if (!(await picker.isVisible().catch(() => false))) return false;
  const img = picker.locator('img.media-library-panel__thumb, img').first();
  if (await img.isVisible().catch(() => false)) {
    await img.click();
    const ok = page
      .locator('.el-dialog button:has-text("确定"), .el-dialog button:has-text("确认")')
      .last();
    if (await ok.isVisible().catch(() => false)) await ok.click();
    return true;
  }
  return false;
}

async function verifyRecordRadio(page, modal) {
  await modal.locator('.el-radio:has-text("录播")').click();
  await page.waitForTimeout(400);
  const recordSelected = await modal
    .locator('.el-radio.is-checked:has-text("录播")')
    .count();
  if (!recordSelected) {
    throw new Error('录播 radio not selected after click (:value fix regression)');
  }
  await modal.locator('.el-radio:has-text("直播")').click();
  await page.waitForTimeout(200);
}

/**
 * Create a live room via UI and verify it appears in the list.
 * @param {import('playwright').Page} page — already logged in
 * @param {{ token: string, roomName?: string, baseUrl?: string }} options
 */
export async function runLiveCreateSubmit(page, options) {
  const { token, roomName = `E2E循环${Date.now() % 100000}`, baseUrl = BASE_URL } = options;
  const anchor = await fetchFirstAnchor(token);
  const coverPath = 'https://cos.qxkejiwl.top/pte-live/image/bg/bg_3.png';

  await page.goto(`${baseUrl}/#/live/index`, { waitUntil: 'networkidle' });
  await page.waitForTimeout(2000);

  const anchorWait = page
    .waitForResponse((r) => r.url().includes('/live/anchor/list') && r.status() === 200, {
      timeout: 12000,
    })
    .catch(() => null);
  await page.locator('button:has-text("创建直播间")').first().click();
  await anchorWait;
  await page.waitForTimeout(800);

  const modal = createModal(page);
  if (!(await modal.isVisible().catch(() => false))) {
    throw new Error('create modal not visible');
  }

  await verifyRecordRadio(page, modal);

  const patched = await patchLiveForm(page, {
    anchor_id: anchor.anchor_id,
    background_img: coverPath,
    cover_img: coverPath,
    name: roomName,
    start_time: futureStartTime(),
  });
  if (!patched) {
    await modal.locator('input[placeholder*="直播间名称"]').fill(roomName);
    await selectAnchor(page, modal);
    await setStartTime(page, modal, futureStartTime());
    await patchLiveForm(page, { background_img: coverPath, cover_img: coverPath });
  }

  const [resp] = await Promise.all([
    page
      .waitForResponse(
        (r) => r.url().includes('/live/create') && r.request().method() === 'POST',
        { timeout: 20000 },
      )
      .catch(() => null),
    modal.locator('button:has-text("提交")').click(),
  ]);

  if (!resp) {
    const toast = await page.locator('.el-message').allTextContents();
    throw new Error(`no create API response${toast.length ? `: ${toast.join('; ')}` : ''}`);
  }
  const body = await resp.json();
  if (body.code !== 1) {
    throw new Error(body.msg || 'live create failed');
  }

  const liveId = body.data?.live_id;
  const roomId = body.data?.roomid || body.data?.room_id || liveId;
  await page.waitForTimeout(1500);

  const searchInput = page.locator('input[placeholder*="直播间名称"]');
  if (await searchInput.count()) {
    await searchInput.first().fill(roomName);
    const searchWait = page.waitForResponse(
      (r) => r.url().includes('/live/list') && r.request().method() === 'POST',
      { timeout: 12000 },
    );
    await page.locator('button:has-text("查询")').first().click();
    await searchWait.catch(() => null);
    await page.waitForTimeout(800);
  }

  let inList =
    (await page.locator('.vxe-body--column', { hasText: roomName }).count()) > 0 ||
    (await page.getByText(roomName, { exact: true }).count()) > 0;
  if (!inList && liveId) {
    inList = (await page.locator(`text=${liveId}`).count()) > 0;
  }
  if (!inList) {
    throw new Error('room not in list after create');
  }

  appendSeedLog({
    module: '直播',
    action: 'create',
    detail: `${roomName} (live_id=${liveId})`,
    method: 'ui',
  });

  return { roomName, liveId, roomId, anchor };
}

async function main() {
  fs.mkdirSync(OUT, { recursive: true });
  const roomName = `E2E循环${Date.now() % 100000}`;

  const { token } = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });

  await injectMerchantToken(page, token);

  try {
    const result = await runLiveCreateSubmit(page, { token, roomName });
    console.log('Fields check skipped in standalone — use CRUD suite for full audit');
    console.log('Anchor seed:', result.anchor.nick_name);
    console.log('live_id:', result.liveId);
    console.log('PASS: live create flow complete');
  } catch (err) {
    console.error('FAIL:', err.message);
    await page.screenshot({ path: path.join(OUT, 'fail-submit.png'), fullPage: true }).catch(() => {});
    process.exit(1);
  }

  await browser.close();
}

const isMain =
  process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1];

if (isMain) {
  main().catch((err) => {
    console.error(err);
    process.exit(1);
  });
}
