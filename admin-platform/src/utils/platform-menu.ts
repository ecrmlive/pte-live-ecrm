import type { RouteRecordStringComponent } from '@vben/types';

export interface PlatformAccessMenuItem {
  access_id?: number;
  children?: PlatformAccessMenuItem[];
  component?: string;
  icon?: string;
  is_menu?: number | string;
  is_route?: number | string;
  is_show?: number | string;
  name: string;
  path: string;
  permission_code?: string;
  redirect_name?: string;
}

function routeNameFromPath(path: string) {
  return (
    String(path || '')
      .split('/')
      .filter(Boolean)
      .map((segment) => segment.replace(/[^a-zA-Z0-9]/g, ''))
      .join('') || 'Root'
  );
}

function resolveComponentPath(component?: string) {
  const raw = String(component || '').trim();
  if (!raw) {
    return undefined;
  }
  const normalized = raw.replace(/^\/+/, '');
  return normalized.startsWith('views/')
    ? `../${normalized}`
    : `../views/${normalized}`;
}

function convertNode(
  item: PlatformAccessMenuItem,
  parentMenuPath?: string,
): null | RouteRecordStringComponent {
  const hideInMenu = Number(item.is_menu) !== 1;
  const isRoute = Number(item.is_route) === 1;
  const menuPath = String(item.path || '').trim();
  const children = (item.children || [])
    .map((child) => convertNode(child, menuPath || parentMenuPath))
    .filter(Boolean) as RouteRecordStringComponent[];

  if (!isRoute && children.length === 0 && hideInMenu) {
    return null;
  }

  const node: RouteRecordStringComponent = {
    name: routeNameFromPath(item.path),
    path: item.path,
    meta: {
      hideInMenu,
      title: item.name,
      ...(hideInMenu && parentMenuPath ? { activePath: parentMenuPath } : {}),
    },
  };

  if (item.icon) {
    node.meta = { ...node.meta, icon: item.icon };
  }

  if (isRoute) {
    const component = resolveComponentPath(item.component);
    if (component) {
      node.component = component;
    }
  } else if (item.redirect_name) {
    node.redirect = item.redirect_name;
  } else if (children.length > 0 && children[0]?.path) {
    node.redirect = children[0].path;
  }

  if (children.length > 0) {
    node.children = children;
  }

  return node;
}

/**
 * 列表页（自带 component）下再挂新增/编辑子路由时，Vben 会删掉父级 component
 * 并 redirect 到第一个子页（如 /region/add）。将这类子路由提到同级注册。
 */
function detachNestedPageRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const result: RouteRecordStringComponent[] = [];
  for (const route of routes) {
    const detached: RouteRecordStringComponent[] = [];
    const kept: RouteRecordStringComponent[] = [];
    for (const child of route.children ?? []) {
      if (route.component && child.component) {
        detached.push(child);
      } else {
        const nested = detachNestedPageRoutes([child]);
        if (nested[0]) {
          kept.push(nested[0]);
        }
        detached.push(...nested.slice(1));
      }
    }
    const next: RouteRecordStringComponent = { ...route };
    if (kept.length > 0) {
      next.children = kept;
    } else {
      delete next.children;
    }
    result.push(next);
    result.push(...detached);
  }
  return result;
}

export function convertPlatformMenusToVben(
  menus: PlatformAccessMenuItem[],
): RouteRecordStringComponent[] {
  const roots = (menus || [])
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];
  return detachNestedPageRoutes(roots);
}

/**
 * 仅提取带 component 的页面路由（平铺）供 Vue Router 注册。
 * 侧栏菜单仍用完整 convertPlatformMenusToVben 树。
 *
 * 切勿把平铺结果再写回菜单树并重复注册 —— 会导致同名路由冲突、Tab 切换白屏。
 */
export function extractPlatformRouterRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const pages: RouteRecordStringComponent[] = [];

  function walk(route: RouteRecordStringComponent) {
    if (route.component) {
      pages.push({ ...route, children: undefined });
      return;
    }
    for (const child of route.children ?? []) {
      walk(child);
    }
  }

  for (const route of routes) {
    walk(route);
  }
  return pages;
}

export const QIXI_PLATFORM_MENU_KEY = 'pte_live_platform_menu';

export function clearPlatformMenuCache() {
  sessionStorage.removeItem(QIXI_PLATFORM_MENU_KEY);
}
