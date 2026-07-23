import type { Router } from 'vue-router';

import {
  MERCHANT_AGENT_HUB_PATH,
  MERCHANT_ARTICLE_HUB_PATH,
  MERCHANT_PLUS_LIVE_HUB_PATH,
} from '#/utils/qixi-live-menu';

/** Tab 容器：插件卡片 path 即 hub，勿跟 redirect 进导出等子路由 */
export const PLUS_TAB_HUB_ENTRY_PATHS = new Set<string>([
  MERCHANT_AGENT_HUB_PATH,
  MERCHANT_ARTICLE_HUB_PATH,
  MERCHANT_PLUS_LIVE_HUB_PATH,
]);

/** 插件中心卡片 path → 实际 hub 路由 */
export const PLUS_PLUGIN_ENTRY_ALIASES: Record<string, string> = {
  '/plus/sign': '/plus/sign/index',
  '/plus/article': '/plus/article/index',
  '/plus/live/index': '/plus/live/wx/index',
};

export function resolvePlusPluginEntryPath(item: {
  path?: string;
  redirect_name?: string;
}): string {
  const path = String(item.path || '').trim();
  const redirect = String(item.redirect_name || '').trim();
  if (path && PLUS_TAB_HUB_ENTRY_PATHS.has(path)) {
    return path;
  }
  if (path && PLUS_PLUGIN_ENTRY_ALIASES[path]) {
    return PLUS_PLUGIN_ENTRY_ALIASES[path];
  }
  return redirect || path;
}

/** 插件中心卡片 → 路由跳转（保留插件中心 tab） */
export async function openPlusPluginPage(router: Router, path: string) {
  const target = String(path || '').trim();
  if (!target) {
    return;
  }
  await router.push(target);
}

/** 商品型 hub Tab 直链 → 统一 hub + ?type=（避免多 route 同名组件 keep-alive Tab 错乱） */
export const PLUS_PRODUCT_HUB_TAB_CANONICAL: Record<
  string,
  { hub: string; type: string }
> = {
  '/plus/bargain/product/index': { hub: '/plus/bargain/index', type: 'product' },
  '/plus/bargain/setting/index': { hub: '/plus/bargain/index', type: 'setting' },
  '/plus/bargain/task/index': { hub: '/plus/bargain/index', type: 'task' },
  '/plus/seckill/active/index': { hub: '/plus/seckill/index', type: 'first' },
  '/plus/seckill/product/index': { hub: '/plus/seckill/index', type: 'second' },
  '/plus/seckill/setting/index': { hub: '/plus/seckill/index', type: 'fourth' },
  '/plus/seckill/time/index': { hub: '/plus/seckill/index', type: 'third' },
};

export function resolvePlusProductHubTabRedirect(
  to: { path: string; query: Record<string, unknown> },
) {
  const hit = PLUS_PRODUCT_HUB_TAB_CANONICAL[to.path];
  if (!hit) {
    return null;
  }
  if (to.path === hit.hub && String(to.query.type ?? '') === hit.type) {
    return null;
  }
  return {
    path: hit.hub,
    query: { ...to.query, type: hit.type },
    replace: true as const,
  };
}
