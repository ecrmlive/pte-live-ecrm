#!/usr/bin/env node
/**
 * Deep layout + interaction E2E for plus plugin list hubs (8 marketing hubs).
 *
 * Usage:
 *   node scripts/e2e-plus-assemble-layout.mjs
 *   node scripts/e2e-plus-assemble-layout.mjs --plugin assemble
 *   node scripts/e2e-plus-assemble-layout.mjs --plugin bargain,seckill,advance
 *
 * Env: E2E_BASE_URL, E2E_API_BASE (see dev-login.mjs)
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';
import {
  closeTopOverlays,
  countOrphanPoppers,
  resetMerchantPage,
} from './e2e-cleanup.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const VIEWPORT = { width: 1440, height: 900 };
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/plus-assemble-layout');

/** First tab with product list + add button */
const PLUGINS = {
  assemble: {
    name: '限时拼团',
    path: '/plus/assemble/index',
    listTab: '拼团商品',
    addLabel: '添加商品',
    addTab: '拼团商品',
    settingTab: '基础设置',
  },
  bargain: {
    name: '限时砍价',
    path: '/plus/bargain/index',
    listTab: '砍价商品',
    addLabel: '添加商品',
    addTab: '砍价商品',
    settingTab: '砍价设置',
  },
  seckill: {
    name: '秒杀',
    path: '/plus/seckill/index',
    listTab: '秒杀商品',
    addLabel: '添加活动',
    addTab: '秒杀活动',
    settingTab: '基础设置',
  },
  advance: {
    name: '预售',
    path: '/plus/advance/index',
    listTab: '预售商品',
    addLabel: '选择商品',
    addTab: '预售商品',
    settingTab: '基础设置',
    expectSearchForm: false,
  },
  coupon: {
    name: '优惠券',
    path: '/plus/coupon/index',
    listTab: '优惠券列表',
    addLabel: '添加优惠券',
    addTab: '优惠券列表',
    settingTab: '基础设置',
    expectSearchForm: false,
  },
  points: {
    name: '积分商城',
    path: '/plus/points/index',
    listTab: '商品设置',
    addLabel: '选择商品',
    addTab: '商品设置',
    settingTab: '兑换设置',
    expectSearchForm: false,
  },
  live: {
    name: '直播',
    path: '/plus/live/index',
    listTab: '直播商品',
    addLabel: '添加商品',
    addTab: '直播商品',
    addButtonSelector: 'button',
    expectSearchForm: false,
  },
  sign: {
    name: '签到',
    path: '/plus/sign/index',
    listTab: '签到记录',
    settingTab: '签到设置',
    skipAdd: true,
  },
};

function parseArgs() {
  const args = process.argv.slice(2);
  let plugins = Object.keys(PLUGINS);
  for (let i = 0; i < args.length; i += 1) {
    if (args[i] === '--plugin' && args[i + 1]) {
      plugins = args[i + 1].split(',').map((s) => s.trim()).filter(Boolean);
    }
  }
  return plugins.filter((key) => PLUGINS[key]);
}

async function countLayout(page) {
  return page.evaluate(() => {
    const isVisible = (el) => {
      const style = window.getComputedStyle(el);
      if (style.display === 'none' || style.visibility === 'hidden') return false;
      const rect = el.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    };

    const pageRoot = document.querySelector('.native-list-page, .list-panel') ?? document.body;
    const grids = [...pageRoot.querySelectorAll('.vxe-grid')].filter(isVisible);
    const forms = [...pageRoot.querySelectorAll('.vxe-grid--form-wrapper form')].filter(isVisible);
    const floatingPickers = [
      ...document.querySelectorAll('.el-picker__popper, .el-select__popper, .media-library-panel'),
    ].filter(isVisible);
    const inlineForms = [...pageRoot.querySelectorAll('.native-form-page, [class*="-product-form-page"]')].filter(
      (el) => isVisible(el) && !el.closest('[role="dialog"], .vben-modal, [data-dismissable-modal]'),
    );

    const formWrap = pageRoot.querySelector('.native-vxe-grid .vxe-grid--form-wrapper');
    const toolbarWrap = pageRoot.querySelector('.native-vxe-grid .vxe-grid--toolbar-wrapper');
    let toolbarSearchSeparated = null;
    if (formWrap && toolbarWrap) {
      const formRect = formWrap.getBoundingClientRect();
      const toolbarRect = toolbarWrap.getBoundingClientRect();
      const overlap = toolbarRect.top < formRect.bottom && toolbarRect.bottom > formRect.top;
      const singleRowMix = overlap && Math.abs(formRect.top - toolbarRect.top) < 24;
      const searchBelowToolbar = formRect.top >= toolbarRect.bottom - 4;
      toolbarSearchSeparated = !singleRowMix && searchBelowToolbar;
    }

    return {
      gridCount: grids.length,
      searchFormCount: forms.length,
      floatingPickerCount: floatingPickers.length,
      inlineFormCount: inlineForms.length,
      toolbarSearchSeparated,
      hasTabHeader: !!document.querySelector('.native-status-tabs .el-tabs__header'),
    };
  });
}

async function waitSettingPanel(page) {
  await page
    .waitForFunction(
      () => {
        const isVisible = (el) => {
          const style = window.getComputedStyle(el);
          if (style.display === 'none' || style.visibility === 'hidden') return false;
          const rect = el.getBoundingClientRect();
          return rect.width > 0 && rect.height > 0;
        };
        const pageRoot = document.querySelector('.native-list-page, .list-panel') ?? document.body;
        const grids = [...pageRoot.querySelectorAll('.vxe-grid')].filter(isVisible);
        return grids.length === 0;
      },
      { timeout: 8000 },
    )
    .catch(() => {});
  await page.waitForTimeout(400);
}

async function auditPlugin(page, key) {
  const plugin = PLUGINS[key];
  const slug = key;
  const checks = [];
  const push = (name, ok, detail = '') => {
    checks.push({ name, ok, detail });
    console.log(`  ${ok ? '✓' : '✗'} ${name}${detail ? ` — ${detail}` : ''}`);
  };

  await page.setViewportSize(VIEWPORT);
  await page.goto(`${BASE_URL}/#${plugin.path}`, { waitUntil: 'domcontentloaded' });
  await page
    .waitForSelector('.native-status-tabs .el-tabs__header, .vxe-grid, .native-form-page', {
      timeout: 15000,
    })
    .catch(() => {});
  await page.waitForTimeout(1500);

  if (page.url().includes('/login')) {
    push('auth', false, 'redirected to login');
    return { plugin: key, checks, pass: 0, fail: checks.length };
  }

  if (plugin.listTab) {
    const listTab = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: plugin.listTab });
    if (await listTab.count()) {
      await listTab.click();
      await page.waitForTimeout(1200);
    }
  }

  const indexLayout = await countLayout(page);
  push('tab-header', indexLayout.hasTabHeader, `hasTabHeader=${indexLayout.hasTabHeader}`);
  push('single-grid', indexLayout.gridCount === 1, `grids=${indexLayout.gridCount}`);
  const expectSearchForm = plugin.expectSearchForm !== false;
  push(
    expectSearchForm ? 'single-search-form' : 'no-search-form',
    expectSearchForm ? indexLayout.searchFormCount === 1 : indexLayout.searchFormCount === 0,
    `forms=${indexLayout.searchFormCount}`,
  );
  push('no-floating-pickers', indexLayout.floatingPickerCount === 0, `pickers=${indexLayout.floatingPickerCount}`);
  push('list-only-index', indexLayout.inlineFormCount === 0, `inlineForms=${indexLayout.inlineFormCount}`);
  if (indexLayout.toolbarSearchSeparated != null && plugin.expectSearchForm !== false) {
    push(
      'toolbar-search-separated',
      indexLayout.toolbarSearchSeparated === true,
      String(indexLayout.toolbarSearchSeparated),
    );
  }

  await page.screenshot({
    path: path.join(OUT_DIR, `${slug}-01-index-toolbar.png`),
    fullPage: false,
  });

  const addTabLabel = plugin.addTab || plugin.listTab;
  if (addTabLabel && addTabLabel !== plugin.listTab) {
    const addTab = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: addTabLabel });
    if (await addTab.count()) {
      await addTab.click();
      await page.waitForSelector('.list-panel .native-vxe-grid, .list-panel .vxe-grid', {
        timeout: 8000,
      }).catch(() => {});
      await page.waitForTimeout(600);
    }
  }

  const addBtn = page
    .locator(plugin.addButtonSelector || '.list-panel .native-vxe-grid button, .list-panel button')
    .filter({ hasText: plugin.addLabel })
    .first();
  await page.waitForTimeout(300);
  if (plugin.skipAdd) {
    push('add-button', true, 'skipped');
  } else if (await addBtn.count()) {
    await addBtn.click();
    await page.waitForTimeout(1200);

    const afterAdd = await countLayout(page);
    push('add-opens-dialog-not-inline', afterAdd.inlineFormCount === 0, `inlineForms=${afterAdd.inlineFormCount}`);

    await page.screenshot({
      path: path.join(OUT_DIR, `${slug}-02-add-dialog.png`),
      fullPage: false,
    });

    const cancel = page
      .locator('[role="dialog"] button, .vben-modal button, .el-dialog button')
      .filter({ hasText: /^取消$/ })
      .first();
    if (await cancel.isVisible().catch(() => false)) {
      await cancel.click();
      await page.waitForTimeout(800);
    } else {
      await closeTopOverlays(page);
    }

    const orphanPoppers = await countOrphanPoppers(page);
    push('no-orphan-poppers', orphanPoppers === 0, `poppers=${orphanPoppers}`);
  } else {
    push('add-button', false, 'not found');
  }

  if (plugin.settingTab) {
    await closeTopOverlays(page);
    const settingTab = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: plugin.settingTab });
    if (await settingTab.count()) {
      await settingTab.click();
      await waitSettingPanel(page);
      const settingLayout = await countLayout(page);
      push('setting-no-list-grid', settingLayout.gridCount === 0, `grids=${settingLayout.gridCount}`);
      push('setting-no-floating-pickers', settingLayout.floatingPickerCount === 0, `pickers=${settingLayout.floatingPickerCount}`);

      await page.screenshot({
        path: path.join(OUT_DIR, `${slug}-03-setting-tab.png`),
        fullPage: false,
      });

      const productTab = page.locator('.native-status-tabs .el-tabs__item').filter({ hasText: plugin.listTab });
      await productTab.click();
      await page.waitForSelector('.vxe-grid', { timeout: 8000 }).catch(() => {});
      await page.waitForTimeout(2000);
      let backLayout = await countLayout(page);
      if (backLayout.gridCount === 0) {
        await productTab.click();
        await page.waitForSelector('.vxe-grid', { timeout: 8000 }).catch(() => {});
        await page.waitForTimeout(1500);
        backLayout = await countLayout(page);
      }
      push('back-to-product-single-grid', backLayout.gridCount === 1, `grids=${backLayout.gridCount}`);
    }
  }

  const pass = checks.filter((c) => c.ok).length;
  const fail = checks.filter((c) => !c.ok).length;
  return { plugin: key, name: plugin.name, checks, pass, fail };
}

async function main() {
  const pluginKeys = parseArgs();
  if (!pluginKeys.length) {
    console.error('No valid --plugin keys. Use: assemble, bargain, seckill, advance, coupon, points, live, sign');
    process.exit(1);
  }

  fs.mkdirSync(OUT_DIR, { recursive: true });

  const { token } = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage();
  await page.goto(BASE_URL);
  await injectMerchantToken(page, token);

  if (process.env.E2E_CHAIN === '1') {
    await resetMerchantPage(page, { reload: true });
  }

  const results = [];
  for (const key of pluginKeys) {
    console.log(`\n[${PLUGINS[key].name}] ${PLUGINS[key].path}`);
    await resetMerchantPage(page, { reload: true });
    results.push(await auditPlugin(page, key));
  }

  await browser.close();

  const summary = {
    viewport: VIEWPORT,
    results,
    pass: results.reduce((n, r) => n + r.pass, 0),
    fail: results.reduce((n, r) => n + r.fail, 0),
  };
  fs.writeFileSync(path.join(OUT_DIR, 'report.json'), JSON.stringify(summary, null, 2));
  console.log(`\nPass ${summary.pass}, Fail ${summary.fail} — ${OUT_DIR}`);
  process.exit(summary.fail > 0 ? 1 : 0);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
