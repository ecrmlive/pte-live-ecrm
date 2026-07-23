#!/usr/bin/env node
/**
 * Merchant-admin CRUD flow E2E — product / order / live / assemble / member / coupon.
 *
 * Run: node scripts/e2e-crud-flows.mjs
 * Page DIY routes: node scripts/e2e-page-diy-routes.mjs
 * Env: E2E_BASE_URL, E2E_API_BASE (see dev-login.mjs)
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';
import { closeTopOverlays, countOrphanPoppers, resetMerchantPage } from './e2e-cleanup.mjs';
import { dismissPopovers, runLiveCreateSubmit } from './e2e-live-create-submit.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/crud-flows');
const REPORT_PATH = path.join(OUT_DIR, 'report.json');
const VIEWPORT = { width: 1440, height: 900 };

const report = {
  startedAt: new Date().toISOString(),
  baseUrl: BASE_URL,
  steps: [],
  pass: 0,
  fail: 0,
};

function step(module, name, ok, detail = '') {
  report.steps.push({ module, name, ok, detail });
  const icon = ok ? '✓' : '✗';
  console.log(`${icon} [${module}] ${name}${detail ? ` — ${detail}` : ''}`);
  if (ok) report.pass += 1;
  else report.fail += 1;
}

async function shot(page, name) {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  await page.screenshot({ path: path.join(OUT_DIR, name), fullPage: false });
}

async function waitStable(page, ms = 1000) {
  await page.waitForLoadState('networkidle', { timeout: 15000 }).catch(() => {});
  await page.waitForTimeout(ms);
}

async function fillProductMinimal(page) {
  const form = page.locator('.product-add-form');
  await form.waitFor({ state: 'visible', timeout: 12000 });

  const productName = form.getByLabel('商品名称', { exact: false });
  if (await productName.count()) {
    await productName.fill(`E2E商品${Date.now() % 100000}`);
  } else {
    await form.locator('.vben-form-item:has-text("商品名称") input, label:has-text("商品名称") + * input').first().fill(`E2E商品${Date.now() % 100000}`);
  }

  const productNo = form.locator('.vben-form-item:has-text("商品编码") input, label:has-text("商品编码") + * input').first();
  if (await productNo.count()) {
    await productNo.fill(`E2E${Date.now() % 100000}`);
  }

  const categorySelect = form.locator('.vben-form-item:has-text("所属分类") .el-select, label:has-text("所属分类") + * .el-select').first();
  if (await categorySelect.count()) {
    await categorySelect.locator('.el-select__wrapper, .el-input__wrapper').click();
    await page.waitForTimeout(400);
    const option = page.locator('.el-select-dropdown:visible .el-select-dropdown__item').first();
    if (await option.count()) {
      await option.click();
      await page.waitForTimeout(300);
    }
    await dismissPopovers(page);
  }

  const previewRadio = form.locator('.el-radio').filter({ hasText: '开启' }).first();
  if (await previewRadio.count()) {
    await previewRadio.click();
    await page.waitForTimeout(300);
    const dateInput = form.locator('.el-date-editor input').first();
    if (await dateInput.count()) {
      await dateInput.click();
      await page.waitForTimeout(400);
    }
  }
}

async function flowProduct(page) {
  const mod = '商品';
  await page.goto(`${BASE_URL}/#/product/product/index`, { waitUntil: 'domcontentloaded' });
  await waitStable(page);
  await shot(page, '01-product-list.png');

  await page.locator('button').filter({ hasText: '添加商品' }).first().click();
  await page.waitForTimeout(900);
  const addOpen = await page.locator('.product-add-form').isVisible().catch(() => false);
  step(mod, 'add modal open', addOpen);

  if (addOpen) {
    await fillProductMinimal(page);
    step(mod, 'fill minimal fields', true);

    const cancelBtn = page.locator('.product-add-form__footer button').filter({ hasText: '取消' }).first();
    await cancelBtn.click();
    await page.waitForTimeout(700);
    await dismissPopovers(page);

    const poppers = await countOrphanPoppers(page);
    step(mod, 'cancel clean (0 poppers)', poppers === 0, `poppers=${poppers}`);
    await shot(page, '02-product-add-cancel.png');
  }

  const editBtn = page.locator('.vxe-body--column button').filter({ hasText: '编辑' }).first();
  if (await editBtn.isVisible().catch(() => false)) {
    await editBtn.click();
    await waitStable(page, 1500);
    const onEdit =
      page.url().includes('/product/product/edit') ||
      (await page.locator('.product-edit-page, .product-add-form').count()) > 0;
    step(mod, 'edit first row', onEdit, page.url());
    await shot(page, '03-product-edit.png');

    if (onEdit) {
      await page.locator('.product-add-form__footer button').filter({ hasText: '取消' }).first().click().catch(() => {});
      await page.goto(`${BASE_URL}/#/product/product/index`, { waitUntil: 'domcontentloaded' });
      await waitStable(page);
    }
  } else {
    step(mod, 'edit first row', false, 'no edit button');
  }

  const deleteBtn = page.locator('.vxe-body--column button').filter({ hasText: '删除' }).first();
  if (await deleteBtn.isVisible().catch(() => false)) {
    await deleteBtn.click();
    await page.waitForTimeout(400);
    const msgBox = page.locator('.el-message-box');
    const confirmVisible = await msgBox.isVisible().catch(() => false);
    step(mod, 'delete confirm open', confirmVisible);

    const cancelConfirm = msgBox.locator('button').filter({ hasText: '取消' });
    if (await cancelConfirm.count()) {
      await cancelConfirm.click();
    } else {
      await page.keyboard.press('Escape');
    }
    await page.waitForTimeout(500);
    step(mod, 'delete cancel', !(await msgBox.isVisible().catch(() => false)));
    await shot(page, '04-product-delete-cancel.png');
  } else {
    step(mod, 'delete cancel', false, 'no delete button');
  }
}

async function flowOrder(page) {
  const mod = '订单';
  await page.goto(`${BASE_URL}/#/order/order/index`, { waitUntil: 'domcontentloaded' });
  await waitStable(page);
  await shot(page, '05-order-list.png');

  const detailBtn = page.locator('.vxe-body--column button, button').filter({ hasText: '详情' }).first();
  if (!(await detailBtn.isVisible().catch(() => false))) {
    step(mod, 'open detail', false, 'no detail button');
    step(mod, 'back to list', false, 'skipped');
    return;
  }

  await detailBtn.click();
  await waitStable(page, 1500);
  const onDetail = page.url().includes('/order/order/detail');
  step(mod, 'open detail', onDetail, page.url());
  await shot(page, '06-order-detail.png');

  const backBtn = page.locator('button').filter({ hasText: '返回列表' }).first();
  if (await backBtn.isVisible().catch(() => false)) {
    await backBtn.click();
    await waitStable(page);
  } else {
    await page.goto(`${BASE_URL}/#/order/order/index`, { waitUntil: 'domcontentloaded' });
    await waitStable(page);
  }

  const backOk = page.url().includes('/order/order/index');
  step(mod, 'back to list', backOk, page.url());
  await shot(page, '07-order-back.png');
}

async function flowLive(page, context, token, user) {
  const mod = '直播';
  let liveId;
  let roomId;
  try {
    const result = await runLiveCreateSubmit(page, { token, baseUrl: BASE_URL });
    liveId = result.liveId;
    roomId = result.roomId;
    step(mod, 'create room', true, `live_id=${result.liveId}`);
    await shot(page, '08-live-after-create.png');
  } catch (err) {
    step(mod, 'create room', false, err.message);
    return;
  }

  let controlPage = null;
  let openedVia = 'navigate';
  const controlBtn = page.locator('.table-buttons button, .vxe-body--column button').filter({ hasText: '中控台' }).first();
  await controlBtn.scrollIntoViewIfNeeded().catch(() => {});

  if ((await controlBtn.count()) > 0) {
    try {
      [controlPage] = await Promise.all([
        context.waitForEvent('page', { timeout: 8000 }),
        controlBtn.click(),
      ]);
      openedVia = 'button';
    } catch {
      controlPage = null;
    }
  }

  if (!controlPage) {
    controlPage = await context.newPage();
    await injectMerchantToken(controlPage, token, user);
    await controlPage.goto(
      `${BASE_URL}/#/live/control/center?live_id=${liveId}&room_id=${roomId}`,
      { waitUntil: 'domcontentloaded', timeout: 45000 },
    );
  }

  await controlPage.waitForSelector('.live-control-panel', { timeout: 35000 }).catch(() => {});
  await controlPage.waitForTimeout(2500);

  const panelOk = await controlPage.evaluate(() => Boolean(document.querySelector('.live-control-panel')));
  step(mod, 'open control center', panelOk, openedVia);
  await controlPage.screenshot({ path: path.join(OUT_DIR, '09-live-control.png'), fullPage: false }).catch(() => {});
  await controlPage.close().catch(() => {});
}

async function flowAssemble(page) {
  const mod = '拼团';
  await page.goto(`${BASE_URL}/#/plus/assemble/index`, { waitUntil: 'domcontentloaded' });
  await waitStable(page);

  for (const tab of ['拼团列表', '基础设置', '拼团商品']) {
    const tabEl = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: tab });
    if (await tabEl.count()) {
      await tabEl.click();
      await page.waitForTimeout(900);
    }
  }
  await page.waitForSelector('.assemble-product-panel .vxe-grid', { timeout: 12000 }).catch(() => {});
  await page.waitForTimeout(800);
  step(mod, 'switch tabs', true);
  await shot(page, '10-assemble-tabs.png');

  const addBtn = page.locator('.list-panel button').filter({ hasText: '添加商品' }).first();
  await addBtn.scrollIntoViewIfNeeded().catch(() => {});
  if ((await addBtn.count()) === 0) {
    step(mod, 'add dialog open', false, 'button missing');
    step(mod, 'close clean (0 poppers)', false, 'skipped');
    return;
  }

  let dialogOpen = false;
  try {
    await Promise.all([
      page.waitForFunction(
        () =>
          Boolean(
            document.querySelector('[role="dialog"]') ||
              document.querySelector('[data-dismissable-modal]') ||
              document.querySelector('.vben-modal'),
          ),
        { timeout: 10000 },
      ),
      addBtn.click(),
    ]);
    dialogOpen = true;
  } catch {
    dialogOpen =
      (await page.locator('[role="dialog"], [data-dismissable-modal], .vben-modal').count()) > 0;
  }
  step(mod, 'add dialog open', dialogOpen, dialogOpen ? 'product picker' : 'no dialog');
  await shot(page, '11-assemble-add-dialog.png');

  if (dialogOpen) {
    await closeTopOverlays(page);
  }
  const poppers = await countOrphanPoppers(page).catch(() => -1);
  step(mod, 'close clean (0 poppers)', poppers === 0, `poppers=${poppers}`);
  await shot(page, '12-assemble-closed.png');
}

async function flowMember(page) {
  const mod = '会员';
  await page.goto(`${BASE_URL}/#/user/user/index`, { waitUntil: 'domcontentloaded' });
  await waitStable(page);
  await shot(page, '13-member-list.png');

  const addBtn = page.locator('button').filter({ hasText: '添加会员' }).first();
  if (!(await addBtn.isVisible().catch(() => false))) {
    step(mod, 'add dialog open', false, 'button missing');
    step(mod, 'add dialog close', false, 'skipped');
    return;
  }

  await addBtn.click();
  await page.waitForTimeout(900);
  const dialogOpen =
    (await page.locator('[role="dialog"], .vben-modal').filter({ hasText: '添加会员' }).count()) > 0 ||
    (await page.getByText('添加会员').count()) > 0;
  step(mod, 'add dialog open', dialogOpen);
  await shot(page, '14-member-add-dialog.png');

  await closeTopOverlays(page);
  const stillOpen = await page.locator('[role="dialog"], .vben-modal').filter({ hasText: '添加会员' }).isVisible().catch(() => false);
  step(mod, 'add dialog close', !stillOpen);
  await shot(page, '15-member-closed.png');
}

async function flowCoupon(page) {
  const mod = '优惠券';
  await page.goto(`${BASE_URL}/#/plus/coupon/index`, { waitUntil: 'domcontentloaded' });
  await waitStable(page);
  await shot(page, '16-coupon-list.png');

  const listTab = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: '优惠券列表' });
  if (await listTab.count()) {
    await listTab.click();
    await page.waitForSelector('.list-panel .vxe-grid', { timeout: 8000 }).catch(() => {});
    await page.waitForTimeout(600);
  }

  const addBtn = page.locator('.list-panel button').filter({ hasText: '添加优惠券' }).first();
  if (!(await addBtn.isVisible().catch(() => false))) {
    step(mod, 'add activity dialog', false, 'button missing');
    step(mod, 'dialog close', false, 'skipped');
    return;
  }

  let dialogOpen = false;
  try {
    await Promise.all([
      page.waitForSelector('[role="dialog"], .vben-modal, .el-dialog', {
        state: 'visible',
        timeout: 12000,
      }),
      addBtn.click(),
    ]);
    dialogOpen = true;
  } catch {
    await page.waitForTimeout(1200);
    dialogOpen =
      (await page.locator('[role="dialog"], .vben-modal, .el-dialog').filter({ hasText: /添加优惠券|优惠券/ }).count()) > 0;
  }
  step(mod, 'add activity dialog', dialogOpen);
  await shot(page, '17-coupon-add-dialog.png');

  await closeTopOverlays(page);
  const poppers = await countOrphanPoppers(page);
  step(mod, 'dialog close', poppers === 0, `poppers=${poppers}`);
  await shot(page, '18-coupon-closed.png');
}

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });

  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const context = await browser.newContext({ viewport: VIEWPORT });
  const page = await context.newPage();

  await injectMerchantToken(page, login.token, login.user);
  if (process.env.E2E_CHAIN === '1') {
    await resetMerchantPage(page, { reload: true });
  } else {
    await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' });
    await waitStable(page, 800);
  }

  const flows = [
    () => flowProduct(page),
    () => flowOrder(page),
    () => flowLive(page, context, login.token, login.user),
    () => flowAssemble(page),
    () => flowMember(page),
    () => flowCoupon(page),
  ];

  for (const run of flows) {
    try {
      await closeTopOverlays(page);
      await run();
    } catch (err) {
      const mod = err.module || 'unknown';
      step(mod, 'flow error', false, err.message);
      await shot(page, `error-${Date.now()}.png`).catch(() => {});
    }
    await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'domcontentloaded' }).catch(() => {});
    await resetMerchantPage(page);
  }

  report.finishedAt = new Date().toISOString();
  fs.writeFileSync(REPORT_PATH, JSON.stringify(report, null, 2));

  console.log(`\nPass ${report.pass} / Fail ${report.fail}`);
  console.log(`Report: ${REPORT_PATH}`);
  console.log(`Screenshots: ${OUT_DIR}`);

  await browser.close();
  process.exit(report.fail > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
