#!/usr/bin/env node
/**
 * Screenshot verify: card hub tabs (卡券管理 / 分类 / 提货码 / 提货订单 / 卡券设置).
 * Run: node scripts/e2e-card-hub-tabs.mjs
 */
import { mkdirSync } from 'node:fs';
import { join } from 'node:path';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT = join(process.cwd(), '.e2e-screenshots/card-hub-tabs');

const TABS = [
  { file: '01-card-list', label: '卡券管理', checks: { noTitle: '卡券列表', hasAdd: '添加卡券' } },
  { file: '02-category', label: '分类管理', checks: { hasAdd: '添加分类' } },
  { file: '03-code', label: '提货码管理', checks: { noTitle: '提货码管理', hasExport: '导出' } },
  { file: '04-order', label: '提货订单', checks: { noTitle: '提货订单' } },
  {
    file: '05-setting',
    label: '卡券设置',
    checks: { imagePicker: true },
  },
];

async function main() {
  mkdirSync(OUT, { recursive: true });
  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { height: 900, width: 1440 } });
  await injectMerchantToken(page, login.token, login.user);

  await page.goto(`${BASE_URL}/#/plus/card/event`, { waitUntil: 'networkidle', timeout: 30000 });
  await page.waitForTimeout(1200);

  const results = [];

  for (const tab of TABS) {
    const tabItem = page.getByRole('tab', { name: tab.label, exact: true });
    await tabItem.click();
    await page.waitForTimeout(900);

    const metrics = await page.evaluate(({ checks, label }) => {
      const titleEl = document.querySelector(
        '.native-vxe-grid .vxe-grid--toolbar-wrapper .text-\\[1rem\\], .native-vxe-grid .font-bold.text-\\[1rem\\]',
      );
      const duplicateGridTitle = checks.noTitle
        ? Boolean(titleEl?.textContent?.includes(checks.noTitle))
        : false;
      const toolbar = document.querySelector('.native-vxe-grid .vxe-grid--toolbar-wrapper');
      const toolbarRect = toolbar?.getBoundingClientRect();
      const refreshBtn = toolbar?.querySelector('.vxe-tools--operate .vxe-button--icon, [title="刷新"]');
      const refreshRect = refreshBtn?.getBoundingClientRect();
      const addBtn = [...document.querySelectorAll('button')].find((btn) =>
        checks.hasAdd ? btn.textContent?.includes(checks.hasAdd) : false,
      );
      const addRect = addBtn?.getBoundingClientRect();

      return {
        duplicateGridTitle,
        hasImagePicker: checks.imagePicker
          ? document.querySelectorAll('.image-picker-trigger__btn').length
          : undefined,
        refreshRightOfAdd:
          !addRect || !refreshRect ? null : refreshRect.left > addRect.right + 8,
        toolbarWidth: toolbarRect ? Math.round(toolbarRect.width) : null,
        tab: label,
      };
    }, tab);

    const file = join(OUT, `${tab.file}.png`);
    await page.screenshot({ fullPage: false, path: file });

    results.push({
      ...tab,
      file,
      metrics,
      ok:
        !metrics.duplicateGridTitle &&
        (tab.checks.imagePicker ? (metrics.hasImagePicker ?? 0) > 0 : true),
    });
  }

  await browser.close();
  console.log(JSON.stringify(results, null, 2));
  const failed = results.filter((r) => !r.ok);
  process.exit(failed.length ? 1 : 0);
}

void main().catch((err) => {
  console.error(err.message || err);
  process.exit(1);
});
