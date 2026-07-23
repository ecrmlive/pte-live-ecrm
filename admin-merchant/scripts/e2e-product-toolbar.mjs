#!/usr/bin/env node
/**
 * Verify product list toolbar (single row) and add-product dialog popper cleanup.
 * Run: node scripts/e2e-product-toolbar.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT_DIR = path.join(__dirname, '../.e2e-screenshots/product-toolbar-fix');

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });

  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });

  await injectMerchantToken(page, login.token);
  await page.goto(`${BASE_URL}/#/product/product/index`, { waitUntil: 'networkidle' });
  await page.waitForTimeout(1200);

  await page.screenshot({ path: path.join(OUT_DIR, '01-product-list-toolbar.png'), fullPage: false });

  const toolbarMetrics = await page.evaluate(() => {
    const form = document.querySelector('.native-vxe-grid .vxe-grid--form-wrapper');
    const toolbar = document.querySelector('.native-vxe-grid .vxe-grid--toolbar-wrapper');
    if (!form || !toolbar) return { error: 'missing form or toolbar' };
    const formRect = form.getBoundingClientRect();
    const toolbarRect = toolbar.getBoundingClientRect();
    const overlap =
      toolbarRect.top < formRect.bottom && toolbarRect.bottom > formRect.top;
    const singleRowMix = overlap && Math.abs(formRect.top - toolbarRect.top) < 24;
    const toolbarAboveSearch = toolbarRect.bottom <= formRect.top + 8;
    const searchBelowToolbar = formRect.top >= toolbarRect.bottom - 4;
    const separatedRows = !singleRowMix && searchBelowToolbar;
    return {
      formTop: Math.round(formRect.top),
      toolbarTop: Math.round(toolbarRect.top),
      formBottom: Math.round(formRect.bottom),
      toolbarBottom: Math.round(toolbarRect.bottom),
      overlap,
      singleRowMix,
      separatedRows,
    };
  });

  console.log('Toolbar layout:', JSON.stringify(toolbarMetrics, null, 2));

  const typeSelectMetrics = await page.evaluate(() => {
    const form = document.querySelector('.native-vxe-grid .vxe-grid--form-wrapper');
    const typeField = form?.querySelector('.native-search-field--type');
    const select = typeField?.querySelector('.el-select-v2, .el-select');
    if (!select) return { error: 'missing product type select' };
    const rect = select.getBoundingClientRect();
    const style = window.getComputedStyle(select);
    return {
      width: Math.round(rect.width),
      minWidth: style.minWidth,
    };
  });
  console.log('Product type select:', JSON.stringify(typeSelectMetrics, null, 2));

  const typeSelect = page
    .locator('.native-search-field--type .el-select-v2, .native-search-field--type .el-select')
    .first();
  await typeSelect.click();
  await page.waitForTimeout(500);

  const dropdownMetrics = await page.evaluate(() => {
    const popper = document.querySelector(
      '.native-search-select-popper:not([style*="display: none"])',
    );
    if (!popper) return { error: 'dropdown popper not visible' };
    const rect = popper.getBoundingClientRect();
    const options = [...popper.querySelectorAll('.el-select-dropdown__item, .el-select-v2__option')];
    const labels = options.slice(0, 5).map((el) => el.textContent?.trim() ?? '');
    const truncated = labels.some((label) => label.length <= 2 && label.endsWith('...'));
    return {
      popperWidth: Math.round(rect.width),
      optionCount: options.length,
      sampleLabels: labels,
      truncated,
    };
  });
  console.log('Product type dropdown:', JSON.stringify(dropdownMetrics, null, 2));

  await page.screenshot({
    path: path.join(OUT_DIR, '01b-product-type-dropdown.png'),
    fullPage: false,
  });

  await page.keyboard.press('Escape');
  await page.waitForTimeout(300);

  const addBtn = page.locator('button').filter({ hasText: '添加商品' }).first();
  await addBtn.click();
  await page.waitForTimeout(800);

  const previewRadio = page.locator('.product-add-form .el-radio').filter({ hasText: '开启' }).first();
  if (await previewRadio.count()) {
    await previewRadio.click();
    await page.waitForTimeout(400);
  }

  const dateInput = page.locator('.product-add-form .el-date-editor input').first();
  if (await dateInput.count()) {
    await dateInput.click();
    await page.waitForTimeout(500);
  }

  await page.screenshot({ path: path.join(OUT_DIR, '02-add-dialog-picker-open.png'), fullPage: false });

  const pickerOpen = await page.evaluate(
    () => document.querySelectorAll('.el-picker__popper').length,
  );
  console.log('Picker poppers while open:', pickerOpen);

  const cancelBtn = page.locator('.product-add-form__footer button').filter({ hasText: '取消' }).first();
  await cancelBtn.click();
  await page.waitForTimeout(600);

  const pickerAfterClose = await page.evaluate(
    () => document.querySelectorAll('.el-picker__popper').length,
  );
  console.log('Picker poppers after close:', pickerAfterClose);

  await page.screenshot({ path: path.join(OUT_DIR, '03-after-close-no-picker.png'), fullPage: false });

  const passToolbar = toolbarMetrics.separatedRows === true;
  const passPicker = pickerAfterClose === 0;
  const passTypeWidth = (typeSelectMetrics.width ?? 0) >= 120;
  const passDropdown =
    !dropdownMetrics.error &&
    (dropdownMetrics.popperWidth ?? 0) >= 200 &&
    dropdownMetrics.truncated !== true &&
    (dropdownMetrics.optionCount ?? 0) > 0;

  console.log(passToolbar ? 'PASS toolbar/search separated rows' : 'FAIL toolbar/search not separated');
  console.log(passTypeWidth ? 'PASS product type select min width' : 'FAIL product type select squeezed');
  console.log(passDropdown ? 'PASS product type dropdown readable' : 'FAIL product type dropdown truncated/hidden');
  console.log(passPicker ? 'PASS no orphan picker' : 'FAIL orphan picker remains');

  await browser.close();
  process.exit(passToolbar && passPicker && passTypeWidth && passDropdown ? 0 : 1);
}

void main().catch((err) => {
  console.error(err);
  process.exit(1);
});
