#!/usr/bin/env node
/**
 * Dev-only merchant login helper for automation (E2E, browser MCP, Playwright).
 *
 * Flow: POST /shop/index/base → read captcha from Redis → POST /shop/passport/login
 * Token is AES-encrypted with the same key as merchant-admin (qixi-live-token.ts).
 *
 * Usage:
 *   node scripts/dev-login.mjs              # print raw JWT + encrypted storage value
 *   node scripts/dev-login.mjs --json       # machine-readable output
 *   node scripts/dev-login.mjs --token-only # raw JWT only (for curl Authorization)
 *   node scripts/dev-login.mjs --inject     # Playwright: inject token + open plugin center
 *
 * Env:
 *   E2E_BASE_URL       default http://localhost:11525
 *   E2E_API_BASE       default http://127.0.0.1:11503
 *   E2E_REDIS_PASSWORD default M8n!xP4qR2vK7sD5cH1t
 *   E2E_USERNAME       default adminm
 *   E2E_PASSWORD       default m123456
 */
import { execSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';

import CryptoJS from 'crypto-js';
import { chromium } from 'playwright';

const BASE_URL = process.env.E2E_BASE_URL || 'http://localhost:11525';
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const REDIS_PASSWORD = process.env.E2E_REDIS_PASSWORD || 'M8n!xP4qR2vK7sD5cH1t';
const USERNAME = process.env.E2E_USERNAME || 'adminm';
const PASSWORD = process.env.E2E_PASSWORD || 'm123456';

const TOKEN_KEY = 'qixiLiveShopToken';
const SECRET_KEY = 'jjj_shop_single_admin_2024';

function encryptToken(token) {
  return CryptoJS.AES.encrypt(JSON.stringify(token), SECRET_KEY).toString();
}

function readCaptchaFromRedis(codeKey) {
  const key = `${codeKey}_shop_code`;
  const attempts = [
    `docker exec pte_live_redis redis-cli -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
    `redis-cli -p 13379 -a "${REDIS_PASSWORD}" GET "${key}" 2>/dev/null`,
  ];
  for (const cmd of attempts) {
    try {
      const out = execSync(cmd, { encoding: 'utf8' }).trim();
      if (out && out !== '(nil)' && !out.includes('Warning')) {
        return out.split('\n').pop();
      }
    } catch {
      /* try next */
    }
  }
  return null;
}

async function fetchJson(url, body) {
  const res = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body ?? {}),
  });
  return res.json();
}

/** API login; captcha answer read from local Redis (dev-only). */
export async function merchantDevLogin(options = {}) {
  const apiBase = options.apiBase || API_BASE;
  const username = options.username || USERNAME;
  const password = options.password || PASSWORD;

  const baseJson = await fetchJson(`${apiBase}/shop/index/base`);
  if (baseJson.code !== 1) {
    throw new Error(`login base failed: ${baseJson.msg}`);
  }

  const codeKey = baseJson.data?.codeData?.codeKey;
  let code;

  if (codeKey) {
    code = readCaptchaFromRedis(codeKey);
    if (!code) {
      throw new Error(
        `captcha not in redis for key ${codeKey}_shop_code — is sql running?`,
      );
    }
  }

  const loginBody = { username, password };
  if (codeKey && code) {
    loginBody.code = code;
    loginBody.codeKey = codeKey;
  }

  const loginJson = await fetchJson(`${apiBase}/shop/passport/login`, loginBody);
  if (loginJson.code !== 1) {
    throw new Error(`login failed: ${loginJson.msg}`);
  }

  const token = loginJson.data?.token;
  if (!token) {
    throw new Error('login succeeded but no token in response');
  }

  return {
    token,
    encrypted: encryptToken(token),
    tokenKey: TOKEN_KEY,
    user: loginJson.data,
    captchaUsed: Boolean(codeKey),
  };
}

/** @param {import('playwright').Page} page */
export async function injectMerchantToken(page, token, user = null) {
  const enc = encryptToken(token);
  const legacyUser = user
    ? {
        AppID: user.app_id,
        app_id: user.app_id,
        logoUrl: user.logoUrl,
        shopName: user.shop_name,
        shopUserId: user.shop_user_id ?? user.shopUserId,
        userId: user.shop_user_id ?? user.userId,
        userName: user.user_name ?? user.userName,
        version: user.version,
      }
    : null;
  await page.goto(BASE_URL);
  await page.evaluate(
    ({ key, value, legacyUser: info }) => {
      localStorage.setItem(key, value);
      sessionStorage.setItem(key, value);
      if (info) {
        sessionStorage.setItem('userInfo', JSON.stringify(info));
      }
    },
    { key: TOKEN_KEY, value: enc, legacyUser },
  );
}

async function main() {
  const args = new Set(process.argv.slice(2));
  const result = await merchantDevLogin();

  if (args.has('--token-only')) {
    console.log(result.token);
    return;
  }

  if (args.has('--json')) {
    console.log(JSON.stringify(result, null, 2));
    return;
  }

  if (args.has('--inject')) {
    const browser = await chromium.launch({ headless: true, channel: 'chrome' });
    const page = await browser.newPage();
    await injectMerchantToken(page, result.token);
    await page.goto(`${BASE_URL}/#/home`, { waitUntil: 'load' });
    await page.waitForTimeout(1500);
    const url = page.url();
    const onLogin = url.includes('/auth/login');
    console.log(onLogin ? 'FAIL: still on login page' : `OK: ${url}`);
    await browser.close();
    process.exit(onLogin ? 1 : 0);
    return;
  }

  console.log('Merchant dev login OK');
  console.log(`  user:     ${result.user.user_name} (${result.user.shop_name})`);
  console.log(`  tokenKey: ${result.tokenKey}`);
  console.log(`  captcha:  ${result.captchaUsed ? 'redis lookup' : 'skipped (shop_code off)'}`);
  console.log('');
  console.log('Raw JWT (curl Authorization: Bearer …):');
  console.log(result.token);
  console.log('');
  console.log('Browser localStorage injection (encrypted, matches qixi-live-token.ts):');
  console.log(`  localStorage.setItem('${result.tokenKey}', '${result.encrypted}')`);
  console.log(`  sessionStorage.setItem('${result.tokenKey}', '${result.encrypted}')`);
  console.log('');
  console.log('Then navigate to http://localhost:11525/#/home and reload.');
}

const isMain =
  process.argv[1] && fileURLToPath(import.meta.url) === process.argv[1];

if (isMain) {
  void main().catch((err) => {
    console.error(err.message || err);
    process.exit(1);
  });
}
