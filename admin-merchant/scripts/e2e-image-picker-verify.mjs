#!/usr/bin/env node
/**
 * Screenshot verify: unified ImagePickerTrigger on key merchant-admin pages.
 * Usage: node scripts/e2e-image-picker-verify.mjs
 */
import { mkdirSync } from 'node:fs';
import { join } from 'node:path';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT = join(process.cwd(), '.e2e-screenshots/image-picker');

const PAGES = [
  { name: 'card-add', path: '#/plus/card/card/add', waitFor: '.image-picker-trigger__btn' },
  { name: 'card-setting', path: '#/plus/card/setting/index', waitFor: '.image-picker-trigger__btn' },
  { name: 'coupon-setting', path: '#/plus/coupon/index?type=setting', waitFor: '.image-picker-trigger__btn' },
  { name: 'agent-background', path: '#/plus/agent/setting/index', tab: '页面背景图', waitFor: '.image-picker-trigger__btn' },
  { name: 'agent-grade-add', path: '#/plus/agent/grade/index', click: '添加等级', waitFor: '.image-picker-trigger__btn' },
  { name: 'agent-poster-add', path: '#/plus/agent/index?type=poster', click: '添加', waitFor: '.image-picker-trigger__btn' },
  { name: 'table-form-add', path: '#/plus/table/table/add', waitFor: '添加表单' },
];

async function main() {
  mkdirSync(OUT, { recursive: true });
  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { height: 900, width: 1440 } });
  await injectMerchantToken(page, login.token, login.user);

  const results = [];

  for (const item of PAGES) {
    try {
      await page.goto(`${BASE_URL}/${item.path}`, { waitUntil: 'networkidle', timeout: 30000 });
      if (item.tab) {
        await page.getByRole('tab', { name: item.tab, exact: true }).click();
        await page.waitForTimeout(800);
      }
      if (item.click) {
        const btn = page.getByRole('button', { name: item.click }).first();
        await btn.click({ timeout: 8000 });
        await page.waitForTimeout(1000);
      }
      if (item.waitFor?.startsWith('.')) {
        await page.locator(item.waitFor).first().waitFor({ timeout: 12000 });
      } else if (item.waitFor) {
        await page.getByText(item.waitFor, { exact: false }).first().waitFor({ timeout: 12000 });
      }
      const file = join(OUT, `${item.name}.png`);
      await page.screenshot({ fullPage: true, path: file });
      const btnCount = await page.locator('.image-picker-trigger__btn').count();
      results.push({ btnCount, name: item.name, ok: true, path: file });
    } catch (err) {
      results.push({ error: err.message, name: item.name, ok: false });
    }
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
