import type {
  RouteLocationNormalized,
  RouteLocationRaw,
} from 'vue-router';

/**
 * 动态路由注册完成后重新导航，使 Vue Router 用新路由表重匹配。
 * 若仅 return true，首次进入 /Home 等页会沿用注册前已解析的 FallbackNotFound（404）。
 */
export function resolveNavigationAfterAccess(
  to: RouteLocationNormalized,
  from: RouteLocationNormalized,
  options: {
    defaultHomePath: string;
    homePath?: string;
  },
): RouteLocationRaw {
  const redirectPath = (from.query.redirect ??
    (to.path === options.defaultHomePath
      ? options.homePath || options.defaultHomePath
      : to.fullPath)) as string;
  const decoded = decodeURIComponent(redirectPath);
  if (decoded !== to.path && decoded !== to.fullPath) {
    return { path: decoded, replace: true };
  }
  return {
    path: to.path,
    query: to.query,
    hash: to.hash,
    replace: true,
  };
}
