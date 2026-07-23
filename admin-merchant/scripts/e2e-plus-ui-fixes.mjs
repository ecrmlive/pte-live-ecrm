/**
 * Screenshot pages touched by plus UI layout fixes (1440px).
 * Run: node scripts/e2e-plus-ui-fixes.mjs [--after]
 */
import { chromium } from 'playwright';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const VIEWPORT = { width: 1440, height: 900 };
const phase = process.argv.includes('--after') ? 'after' : 'before';
const OUT_DIR = path.join(__dirname, `../.e2e-screenshots/plus-ui-fixes-${phase}`);

const TARGETS = [
  { file: '01-fullreduce', path: '/plus/fullreduce/index' },
  { file: '02-table-form-list', path: '/plus/table/event?type=table' },
  { file: '03-recommend', path: '/plus/recommend/index' },
  { file: '04-agent-order', path: '/plus/agent/index?type=order' },
  { file: '05-agent-cash', path: '/plus/agent/index?type=cash' },
  { file: '06-agent-setting-basic', path: '/plus/agent/index?type=setting' },
  { file: '07-table-form-add', path: '/plus/table/table/add' },
];

async function main() {
  fs.mkdirSync(OUT_DIR, { recursive: true });
  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage({ viewport: VIEWPORT });
  await page.goto(`${BASE_URL}/auth/login`);
  await injectMerchantToken(page, login);

  for (const target of TARGETS) {
    await page.goto(`${BASE_URL}${target.path}`, { waitUntil: 'networkidle' });
    await page.waitForTimeout(1200);
    if (target.file.startsWith('06-agent-setting')) {
      const navItem = page.locator('.agent-setting-nav__item').first();
      if (await navItem.count()) {
        await navItem.waitFor({ timeout: 15000 });
      } else {
        await page.locator('.agent-hub-tabs').waitFor({ timeout: 15000 });
      }
    }
    await page.screenshot({
      fullPage: true,
      path: path.join(OUT_DIR, `${target.file}.png`),
    });
    console.log(`saved ${target.file}.png`);
  }

  await browser.close();
  console.log(`Done (${phase}) -> ${OUT_DIR}`);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
