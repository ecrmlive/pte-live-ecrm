#!/usr/bin/env node
/**
 * Seed minimal demo rows for merchant-admin modules that are empty.
 * Appends results to scripts/e2e-seed-log.json — safe to re-run (skips non-empty).
 */
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

import { merchantDevLogin } from './dev-login.mjs';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const API_BASE = process.env.E2E_API_BASE || 'http://127.0.0.1:11503';
const SEED_LOG = path.join(__dirname, 'e2e-seed-log.json');

function appendSeedLog(entry) {
  let data = { entries: [] };
  try {
    data = JSON.parse(fs.readFileSync(SEED_LOG, 'utf8'));
  } catch {
    /* new */
  }
  data.entries.push({ at: new Date().toISOString(), ...entry });
  fs.writeFileSync(SEED_LOG, `${JSON.stringify(data, null, 2)}\n`);
}

async function shopPost(token, path, body = {}) {
  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: {
      AppID: '10000',
      'Content-Type': 'application/json',
      'authori-zation': `Bearer ${token}`,
    },
    body: JSON.stringify(body),
  });
  return res.json();
}

async function listTotal(token, path) {
  const json = await shopPost(token, path, { list_rows: 1, page: 1 });
  if (json.code !== 1) return { ok: false, total: -1, msg: json.msg };
  return { ok: true, total: json.data?.list?.total ?? 0 };
}

async function seedFullreduce(token) {
  const { ok, total, msg } = await listTotal(token, '/shop/plus.fullreduce/index');
  if (!ok) return { skipped: true, reason: msg };
  if (total > 0) return { skipped: true, reason: `already ${total}` };
  const json = await shopPost(token, '/shop/plus.fullreduce/add', {
    active_name: 'E2E满减活动',
    full_type: 1,
    full_value: 100,
    reduce_type: 1,
    reduce_value: 10,
  });
  if (json.code !== 1) return { ok: false, msg: json.msg };
  appendSeedLog({ module: '满减', action: 'create', detail: 'E2E满减活动', method: 'api' });
  return { ok: true };
}

async function livePost(token, apiPath, body = {}) {
  const res = await fetch(`${API_BASE}${apiPath}`, {
    method: 'POST',
    headers: {
      AppID: '10000',
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/x-www-form-urlencoded',
      'authori-zation': `Bearer ${token}`,
    },
    body: new URLSearchParams(
      Object.fromEntries(
        Object.entries(body).map(([k, v]) => [k, Array.isArray(v) ? v.join(',') : String(v)]),
      ),
    ),
  });
  return res.json();
}

async function seedLiveControl(token) {
  const roomsJson = await livePost(token, '/api/v1/shop/live/list', {
    list_rows: 5,
    page: 1,
  });
  const room = roomsJson?.data?.list?.data?.[0];
  if (!room?.live_id) return { skipped: true, reason: 'no live room' };

  const liveId = room.live_id;
  const productsJson = await livePost(token, '/api/v1/shop/live/product/list', {
    live_id: liveId,
    list_rows: 5,
    page: 1,
  });
  const productTotal = productsJson?.data?.list?.length ?? 0;
  if (productTotal > 0) {
    return { skipped: true, reason: `room ${liveId} already has ${productTotal} products` };
  }

  const mallJson = await shopPost(token, '/shop/product.product/index', {
    list_rows: 5,
    page: 1,
    type: 'sell',
  });
  const mallRow = mallJson?.data?.list?.data?.[0];
  if (!mallRow?.product_id) return { skipped: true, reason: 'no mall product to link' };

  const addJson = await livePost(token, '/api/v1/shop/live/product/add', {
    live_id: liveId,
    product_ids: String(mallRow.product_id),
  });
  if (addJson.code !== 1 && addJson.code !== 200) {
    return { ok: false, msg: addJson.msg || addJson.message || 'live product add failed' };
  }
  appendSeedLog({
    module: '直播中控',
    action: 'link-product',
    detail: `live_id=${liveId} product_id=${mallRow.product_id}`,
    method: 'api',
  });
  return { ok: true, liveId, productId: mallRow.product_id };
}

async function main() {
  const { token } = await merchantDevLogin();
  const tasks = [
    ['fullreduce', () => seedFullreduce(token)],
    ['live-control', () => seedLiveControl(token)],
  ];

  for (const [name, fn] of tasks) {
    const result = await fn();
    if (result.ok) console.log(`✓ seeded ${name}`);
    else if (result.skipped) console.log(`○ skip ${name}: ${result.reason}`);
    else console.log(`✗ ${name}: ${result.msg || 'failed'}`);
  }
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
