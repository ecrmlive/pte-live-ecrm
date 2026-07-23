#!/usr/bin/env node
/**
 * Plus hub tab interaction — switch sub-tabs, verify single panel, optional add-dialog smoke.
 *
 * Run: node scripts/e2e-plus-hub-tabs.mjs
 *      node scripts/e2e-plus-hub-tabs.mjs --no-dialog
 */
import { chromium } from 'playwright';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/plus-hub-tabs');
const REPORT_PATH = path.join(OUT_DIR, 'report.json');
const skipDialog = process.argv.includes('--no-dialog');

/** Eight marketing hubs fixed for mount-all tab panels */
const HUBS = [
  {
    name: 'assemble',
    path: '/plus/assemble/index',
    tabs: ['拼团商品', '拼团列表', '基础设置'],
    addButton: '添加商品',
  },
  {
    name: 'bargain',
    path: '/plus/bargain/index',
    tabs: ['砍价商品', '砍价列表', '砍价设置'],
    addButton: '添加商品',
  },
  {
    name: 'seckill',
    path: '/plus/seckill/index',
    tabs: ['秒杀活动', '秒杀商品', '秒杀配置', '基础设置'],
    addButton: '添加活动',
  },
  {
    name: 'advance',
    path: '/plus/advance/index',
    tabs: ['预售商品', '基础设置'],
    addButton: null,
  },
  {
    name: 'coupon',
    path: '/plus/coupon/index',
    tabs: ['优惠券列表', '领取记录', '发送优惠券', '基础设置'],
    addButton: '添加优惠券',
  },
  {
    name: 'points',
    path: '/plus/points/index',
    tabs: ['商品设置', '兑换设置', '兑换记录'],
    addButton: null,
  },
  {
    name: 'live',
    path: '/plus/live/index',
    tabs: ['直播房间', '主播管理', '直播商品'],
    addButton: '新增主播',
  },
  {
    name: 'sign',
    path: '/plus/sign/index',
    tabs: ['签到设置', '签到记录'],
    addButton: null,
  },
];

async function countVisible(page, selector) {
  return page.locator(selector).evaluateAll((nodes) =>
    nodes.filter((node) => {
      const style = window.getComputedStyle(node);
      if (style.display === 'none' || style.visibility === 'hidden') return false;
      const rect = node.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    }).length,
  );
}

async function countListHeaders(page) {
  return countVisible(page, '.native-list-page .vxe-grid, .list-panel .vxe-grid');
}

async function countSearchForms(page) {
  return countVisible(page, '.native-list-page .vxe-grid--form-wrapper form, .list-panel .vxe-grid--form-wrapper form');
}

async function popperLeakCount(page) {
  return page.evaluate(() =>
    document.querySelectorAll(
      'body > .el-popper, body > .el-picker__popper, body > .el-select__popper',
    ).length,
  );
}

/** Unresolved `<el-*>` custom elements — EP not registered or import missing */
const DEAD_EL_TAGS = [
  'el-input',
  'el-select',
  'el-switch',
  'el-button',
  'el-date-picker',
  'el-radio-group',
  'el-tabs',
  'el-checkbox',
  'el-input-number',
];

async function auditDeadElTags(page) {
  return page.evaluate((tags) => {
    const dead = [];
    for (const tag of tags) {
      const nodes = document.querySelectorAll(tag);
      if (nodes.length) dead.push(`${tag}×${nodes.length}`);
    }
    return dead;
  }, DEAD_EL_TAGS);
}

async function auditHub(page, hub) {
  const url = `${BASE_URL}/#${hub.path}`;
  const result = {
    hub: hub.name,
    path: hub.path,
    tabs: [],
    dialog: null,
    status: 'pass',
    issues: [],
  };

  await page.goto(url, { waitUntil: 'networkidle', timeout: 45000 });
  await page
    .waitForSelector('.native-status-tabs .el-tabs__item, .el-tabs__item', { timeout: 15000 })
    .catch(() => {});
  await page.waitForTimeout(2000);

  if (page.url().includes('/auth/login')) {
    result.status = 'fail';
    result.issues.push('redirected to login');
    return result;
  }

  const notFound = (await page.locator('text=未找到页面').count()) > 0;
  if (notFound) {
    result.status = 'skip';
    result.issues.push('route not registered (404)');
    return result;
  }

  for (const tabLabel of hub.tabs) {
    const tabItem = page.getByRole('tab', { name: tabLabel, exact: true });
    await tabItem.first().click({ timeout: 15000 });
    await page
      .waitForFunction(
        (label) => {
          const tab = [...document.querySelectorAll('.native-status-tabs .el-tabs__item, .el-tabs__item')].find(
            (el) => el.textContent?.trim() === label,
          );
          return tab?.classList.contains('is-active') ?? false;
        },
        tabLabel,
        { timeout: 8000 },
      )
      .catch(() => {});
    await page.waitForTimeout(350);

    const tables = await countVisible(page, '.vxe-grid');
    const forms = await countSearchForms(page);
    const active = await tabItem.first().evaluate((el) => el.classList.contains('is-active'));

    const deadTags = await auditDeadElTags(page);
    const tabResult = { label: tabLabel, active, tables, forms, deadTags };
    result.tabs.push(tabResult);

    if (!active) result.issues.push(`tab "${tabLabel}" not active`);
    if (tables > 1) result.issues.push(`tab "${tabLabel}": ${tables} visible tables`);
    if (forms > 1) result.issues.push(`tab "${tabLabel}": ${forms} visible search forms`);
    if (deadTags.length) {
      result.issues.push(`tab "${tabLabel}": dead EP tags ${deadTags.join(', ')}`);
    }
  }

  if (!skipDialog && hub.addButton) {
    const firstTab = hub.tabs[0];
    await page.getByRole('tab', { name: firstTab }).first().click();
    await page.waitForTimeout(500);

    const poppersBefore = await popperLeakCount(page);
    const addBtn = page.locator('button').filter({ hasText: hub.addButton }).first();
    if (await addBtn.count()) {
      await addBtn.click();
      await page.waitForTimeout(700);
      const dialogs = await page.locator('.el-dialog, [role="dialog"]').count();
      result.dialog = { opened: dialogs > 0, poppersBefore, poppersAfter: await popperLeakCount(page) };
      if (dialogs === 0) result.issues.push(`add dialog not opened for "${hub.addButton}" (soft)`);
      await page.keyboard.press('Escape');
      await page.waitForTimeout(400);
      const poppersAfterClose = await popperLeakCount(page);
      if (poppersAfterClose > poppersBefore + 1) {
        result.issues.push(`popper leak after dialog close (${poppersAfterClose})`);
      }
    } else {
      result.dialog = { opened: false, skipped: 'button not found' };
    }
  }

  const hardIssues = result.issues.filter((i) => !i.includes('(soft)'));
  if (hardIssues.length) result.status = 'fail';
  else if (result.status !== 'skip') result.status = 'pass';
  const slug = hub.name;
  await page.screenshot({ path: path.join(OUT_DIR, `${slug}.png`), fullPage: false });
  return result;
}

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await page.goto(BASE_URL);
  await injectMerchantToken(page, login.token);

  const results = [];
  for (const hub of HUBS) {
    const result = await auditHub(page, hub);
    results.push(result);
    const icon = result.status === 'pass' ? '✓' : result.status === 'skip' ? '○' : '✗';
    const detail = result.issues.length ? ` — ${result.issues.join('; ')}` : '';
    console.log(`${icon} ${hub.name}${detail}`);
  }

  await browser.close();

  const summary = {
    hubs: HUBS.length,
    pass: results.filter((r) => r.status === 'pass').length,
    skip: results.filter((r) => r.status === 'skip').length,
    fail: results.filter((r) => r.status === 'fail').length,
    skipDialog,
    results,
  };
  fs.writeFileSync(REPORT_PATH, JSON.stringify(summary, null, 2));
  console.log(`\nPass: ${summary.pass}/${summary.hubs}, Skip: ${summary.skip}, Fail: ${summary.fail}`);
  console.log(`Report: ${REPORT_PATH}`);
  process.exit(summary.fail > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
