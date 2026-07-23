#!/usr/bin/env node
/**
 * Screenshot high-traffic merchant-admin pages and flag common layout issues.
 * Run: node scripts/ui-audit-pages.mjs
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { chromium } from 'playwright';

import { injectMerchantToken, merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const BASE = process.env.E2E_BASE_URL || 'http://localhost:11525';
const OUT = path.join(__dirname, '../.e2e-screenshots/ui-audit-pass');
fs.mkdirSync(OUT, { recursive: true });

const PAGES = [
  { name: 'home', path: '/home' },
  { name: 'page-diy-list', path: '/page/page/index' },
  { name: 'page-center', path: '/page/center/index' },
  { name: 'setting-store', path: '/setting/store/index' },
  { name: 'setting-trade', path: '/setting/trade/index' },
  { name: 'setting-message', path: '/setting/message/index' },
  { name: 'setting-print', path: '/setting/print/index' },
  { name: 'setting-pay', path: '/setting/pay/index' },
  { name: 'plus-bargain', path: '/plus/bargain/index' },
  { name: 'plus-assemble', path: '/plus/assemble/index' },
  { name: 'plus-seckill', path: '/plus/seckill/index' },
  { name: 'plus-coupon', path: '/plus/coupon/index' },
  { name: 'plus-invitation', path: '/plus/invitation/index' },
  { name: 'live-session', path: '/live/session/index' },
  { name: 'order-list', path: '/order/order/index' },
  { name: 'order-refund', path: '/order/refund/index' },
  { name: 'page-tabbar', path: '/page/tabbar/index' },
  { name: 'page-theme', path: '/page/theme/index' },
];

async function auditPage(page, { name, path: hashPath }) {
  const url = `${BASE}/#${hashPath}`;
  await page.goto(url, { waitUntil: 'load', timeout: 60000 });
  await page.waitForTimeout(1500);
  const metrics = await page.evaluate(() => {
    const dupTitles = [...document.querySelectorAll('.table-title, .page-title, h2.title')].map(
      (el) => el.textContent?.trim(),
    ).filter(Boolean);
    const titleCounts = {};
    dupTitles.forEach((t) => {
      titleCounts[t] = (titleCounts[t] || 0) + 1;
    });
    const dupTitleTexts = Object.entries(titleCounts)
      .filter(([, c]) => c > 1)
      .map(([t, c]) => `${t}(${c}x)`);
    const grids = document.querySelectorAll('.native-vxe-grid .vxe-grid, .list-panel .vxe-grid');
    let searchToolbarMerged = false;
    for (const grid of grids) {
      const form = grid.querySelector('.vxe-grid--form-wrapper');
      const toolbar = grid.querySelector('.vxe-grid--toolbar-wrapper');
      if (form && toolbar) {
        const fr = form.getBoundingClientRect();
        const tr = toolbar.getBoundingClientRect();
        if (Math.abs(fr.top - tr.top) < 8 && fr.height > 0 && tr.height > 0) {
          searchToolbarMerged = true;
        }
      }
    }
    const visibleTables = [...document.querySelectorAll('.vxe-table')].filter((el) => {
      const r = el.getBoundingClientRect();
      const s = getComputedStyle(el);
      return r.width > 0 && r.height > 0 && s.display !== 'none' && s.visibility !== 'hidden';
    }).length;
    const pager = document.querySelector('.vxe-pager, .el-pagination');
    const pagerBottom = pager ? window.innerHeight - pager.getBoundingClientRect().bottom : null;
    return {
      dupTitleTexts,
      searchToolbarMerged,
      visibleTables,
      pagerBottom,
      url: location.href,
    };
  });
  await page.screenshot({ path: path.join(OUT, `${name}.png`), fullPage: false });
  const pageIssues = [];
  if (metrics.dupTitleTexts.length) {
    pageIssues.push(`duplicate titles: ${metrics.dupTitleTexts.join(', ')}`);
  }
  if (metrics.searchToolbarMerged) pageIssues.push('search+toolbar merged row');
  if (metrics.visibleTables > 1) pageIssues.push(`${metrics.visibleTables} visible tables`);
  if (metrics.pagerBottom !== null && metrics.pagerBottom > 120) {
    pageIssues.push(`pager ${Math.round(metrics.pagerBottom)}px above viewport bottom`);
  }
  if (metrics.url.includes('/auth/login')) pageIssues.push('redirected to login');
  console.log(
    `${pageIssues.length ? 'FAIL' : 'OK  '} ${name}${pageIssues.length ? ` — ${pageIssues.join('; ')}` : ''}`,
  );
  return { name, path: hashPath, url: metrics.url, issues: pageIssues, metrics };
}

async function main() {
  const login = await merchantDevLogin();
  const browser = await chromium.launch({ headless: true, channel: 'chrome' });
  const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
  await injectMerchantToken(page, login.token, login.user);
  await page.goto(`${BASE}/#/home`, { waitUntil: 'load', timeout: 60000 });
  await page.waitForTimeout(2000);
  const onLogin = page.url().includes('/auth/login');
  if (onLogin) {
    throw new Error('token injection failed — still on login page');
  }
  const issues = [];
  for (const p of PAGES) {
    try {
      issues.push(await auditPage(page, p));
    } catch (e) {
      console.log(`ERR  ${p.name}: ${e.message}`);
      issues.push({ name: p.name, path: p.path, issues: [e.message] });
    }
  }
  fs.writeFileSync(path.join(OUT, 'report.json'), JSON.stringify({ scanned: PAGES.length, issues }, null, 2));
  await browser.close();
  console.log(`\nReport: ${path.join(OUT, 'report.json')}`);
}

void main().catch((err) => {
  console.error(err.message || err);
  process.exit(1);
});
