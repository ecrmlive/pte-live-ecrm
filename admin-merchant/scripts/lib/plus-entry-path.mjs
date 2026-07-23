/** 与 src/utils/plus-navigation.ts 保持同步（E2E / audit 脚本用） */

export const PLUS_TAB_HUB_ENTRY_PATHS = new Set([
  '/plus/agent/index',
  '/plus/article/index',
  '/plus/live/wx/index',
]);

export const PLUS_PLUGIN_ENTRY_ALIASES = {
  '/plus/sign': '/plus/sign/index',
  '/plus/article': '/plus/article/index',
  '/plus/live/index': '/plus/live/wx/index',
};

/** 插件中心卡片 → 实际打开的路由 */
export function resolvePlusPluginEntryPath(item) {
  const path = String(item?.path || '').trim();
  const redirect = String(item?.redirect_name || '').trim();
  if (path && PLUS_TAB_HUB_ENTRY_PATHS.has(path)) {
    return path;
  }
  if (path && PLUS_PLUGIN_ENTRY_ALIASES[path]) {
    return PLUS_PLUGIN_ENTRY_ALIASES[path];
  }
  return redirect || path;
}
