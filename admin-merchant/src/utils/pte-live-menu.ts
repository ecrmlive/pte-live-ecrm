import type { RouteRecordStringComponent } from '@vben/types';

import { resolveLegacyComponentPath } from './pte-live-api';
import { NATIVE_PAGE_OVERRIDES } from '#/views/native/registry';

export interface ShopAccessMenuItem {
  children?: ShopAccessMenuItem[];
  component?: string;
  icon?: string;
  is_menu?: number | string;
  is_route?: number | string;
  is_show?: number | string;
  name: string;
  path: string;
  redirect_name?: string;
  upload_icon?: string;
}

const SKIP_MENU_PATHS = new Set([
  '/live/control/center',
  /** 卡密管理依赖 product_id，仅能从商品列表进入，不可作独立侧栏菜单 */
  '/product/virtual/index',
]);

/** 插件中心卡片页 path（点击侧栏「插件中心」应进入此路由） */
export const MERCHANT_PLUS_HUB_PATH = '/plus/plus/index';

/** 分销等 Tab 容器页：带 component，子路由为页内 Tab / 隐藏操作，不可 redirect 到子页 */
export const MERCHANT_AGENT_HUB_PATH = '/plus/agent/index';

/** 文章插件 Tab 容器（文章列表 / 分类管理） */
export const MERCHANT_ARTICLE_HUB_PATH = '/plus/article/index';

/** App 设置 Tab 容器（基础 / 分享 / 升级） */
export const MERCHANT_APPOPEN_HUB_PATH = '/appsetting/appopen/event';

/** 直播插件 Tab 容器（直播房间 / 主播 / 直播商品；legacy hub 为 wx/index） */
export const MERCHANT_PLUS_LIVE_HUB_PATH = '/plus/live/wx/index';

/** 直播插件 hub + 页内 Tab / CRUD 权限路由（插件卡片 path 为 /plus/live/index，非 route） */
const MERCHANT_PLUS_LIVE_ROUTE_PATHS = [
  MERCHANT_PLUS_LIVE_HUB_PATH,
  '/plus/live/index',
  '/plus/live/wx/add',
  '/plus/live/wx/edit',
  '/plus/live/wx/delete',
  '/plus/live/wx/syn',
  '/plus/live/wx/settop',
  '/plus/live/wx/setSyn',
  '/plus/live/LiveProduct/index',
  '/plus/live/LiveProduct/add',
  '/plus/live/LiveProduct/push',
  '/plus/live/LiveProduct/onsale',
  '/plus/live/LiveProduct/delete',
  '/plus/live/room/qrcode',
  '/plus/live/product/index',
  '/plus/live/product/add',
  '/plus/live/product/edit',
  '/plus/live/product/delete',
  '/plus/live/anchor/index',
  '/plus/live/anchor/add',
  '/plus/live/anchor/edit',
  '/plus/live/anchor/delete',
] as const;

/** App 设置页内 Tab / 升级 CRUD 权限路由，不可出现在侧栏 */
const MERCHANT_APPOPEN_INTERNAL_PATHS = new Set([
  '/appsetting/appopen/index',
  '/appsetting/appshare/index',
  '/appsetting/appupdate/index',
  '/appsetting/appupdate/add',
  '/appsetting/appupdate/edit',
  '/appsetting/appupdate/delete',
]);

/**
 * 插件 Tab 容器（自身渲染 tabs，子 path 仅作权限/导出等）。
 * 若仍挂 children，accessible 会把 redirect 改到首个子路由（如分销→导出）。
 */
const TAB_HUB_PATHS = new Set([
  MERCHANT_PLUS_HUB_PATH,
  MERCHANT_AGENT_HUB_PATH,
  MERCHANT_ARTICLE_HUB_PATH,
  MERCHANT_APPOPEN_HUB_PATH,
  MERCHANT_PLUS_LIVE_HUB_PATH,
]);

/** 页面列表/首页装修列表：不可 redirect 到 add 等权限子路由 */
const MERCHANT_PAGE_LIST_HUB_PATHS = new Set([
  '/page/page/index',
  '/page/page/list',
]);

/** 与 legacy constantRoutes 一致：首页不在 shop_access，登录后默认跳转 /home */
const STATIC_HOME_MENU: ShopAccessMenuItem = {
  component: 'views/native/home/index.vue',
  is_menu: 1,
  is_route: 1,
  is_show: 1,
  name: '我的数据',
  path: '/home',
};

function normalizeBackendComponentPath(raw: string) {
  const normalized = raw.replace(/^(\.\/|\.\.\/)+/, '');
  const viewPath = normalized.startsWith('/') ? normalized : `/${normalized}`;
  let path = viewPath.replace(/^\/views/, '');
  // pageMap key 与 legacy path 一致，不含 native 段
  if (path.startsWith('/native/')) {
    path = path.slice('/native'.length);
  }
  // legacy 首页组件路径为 /home/home.vue
  if (path === '/home/index' || path === '/home/index.vue') {
    return '/home/home';
  }
  return path;
}

function nativePageRegistryKey(path: string): string {
  const legacy = resolveLegacyComponentPath(path);
  const base = legacy.startsWith('/') ? legacy : `/${legacy}`;
  return base.endsWith('.vue') ? base : `${base}.vue`;
}

function hasNativePageComponent(path: string): boolean {
  return Object.prototype.hasOwnProperty.call(
    NATIVE_PAGE_OVERRIDES,
    nativePageRegistryKey(path),
  );
}

/** 后端 component 优先；目录节点（无 component 有 children）不注册 /product.vue 等无效页 */
function resolveMenuComponentPath(item: ShopAccessMenuItem): string | undefined {
  const backendComponent = String(item.component || '').trim();
  if (backendComponent) {
    return normalizeBackendComponentPath(backendComponent);
  }

  const rawChildren = item.children ?? [];
  const hasChildren = rawChildren.length > 0;
  /** 仅挂 add/edit/delete 等按钮权限子路由时仍是列表页，须保留 component */
  const onlyPermissionChildren =
    hasChildren &&
    rawChildren.every((child) => Number(child.is_menu) !== 1);

  const isDirectoryOnly =
    hasChildren &&
    !onlyPermissionChildren &&
    !TAB_HUB_PATHS.has(item.path) &&
    !MERCHANT_PAGE_LIST_HUB_PATHS.has(item.path);
  if (isDirectoryOnly) {
    if (Number(item.is_route) === 1 && hasNativePageComponent(item.path)) {
      return resolveLegacyComponentPath(item.path);
    }
    return undefined;
  }

  if (Number(item.is_route) === 1) {
    return resolveLegacyComponentPath(item.path);
  }

  return undefined;
}

function routeNameFromPath(path: string) {
  return (
    String(path || '')
      .split('/')
      .filter(Boolean)
      .join('') || 'Root'
  );
}

function hideAppopenInternalMenu(
  route: RouteRecordStringComponent,
): RouteRecordStringComponent {
  if (!MERCHANT_APPOPEN_INTERNAL_PATHS.has(route.path)) {
    return route;
  }
  return {
    ...route,
    meta: {
      ...route.meta,
      hideInMenu: true,
    },
  };
}

function convertNode(
  item: ShopAccessMenuItem,
): null | RouteRecordStringComponent {
  if (SKIP_MENU_PATHS.has(item.path)) {
    return null;
  }

  const hideInMenu = Number(item.is_menu) !== 1;
  const isRoute = Number(item.is_route) === 1;
  const redirect =
    TAB_HUB_PATHS.has(item.path) || MERCHANT_PAGE_LIST_HUB_PATHS.has(item.path)
      ? ''
      : String(item.redirect_name || '').trim();
  const children = (item.children || [])
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];

  const visibleChildren = children.filter(
    (child) => child.meta?.hideInMenu !== true,
  );

  if (!isRoute && children.length === 0 && hideInMenu) {
    return null;
  }

  const node: RouteRecordStringComponent = {
    name: routeNameFromPath(item.path),
    path: item.path,
    meta: {
      hideInMenu,
      ...(hideInMenu ? { hideInTab: true } : {}),
      title: item.name,
    },
  };

  if (item.icon) {
    node.meta = { ...node.meta, icon: item.icon };
  }

  if (redirect) {
    node.redirect = redirect;
    // 侧栏仅展示一级入口，点击直接走 redirect（对齐 legacy LeftMenu）
    if (children.length > 0 && visibleChildren.length === 0) {
      node.meta = { ...node.meta, hideChildrenInMenu: true };
    }
  }

  if (isRoute) {
    const component = resolveMenuComponentPath(item);
    if (component) {
      node.component = component;
    }
  }

  if (item.path === MERCHANT_APPOPEN_HUB_PATH) {
    node.meta = {
      ...node.meta,
      hideChildrenInMenu: true,
      title: item.name || 'App设置',
    };
  }

  // 插件中心等：侧栏不展开隐藏子项，仅保留 redirect 跳转
  if (visibleChildren.length > 0) {
    node.children = visibleChildren;
  } else if (children.length > 0 && redirect) {
    node.children = children;
  } else if (children.length > 0) {
    node.children = children;
  }

  return node;
}

function flattenHubSubtree(
  children: RouteRecordStringComponent[],
  lifted: RouteRecordStringComponent[],
) {
  for (const child of children) {
    const nested = child.children ? [...child.children] : [];
    const node = hideAppopenInternalMenu({ ...child });
    delete node.children;
    delete node.redirect;
    lifted.push(node);
    if (nested.length) {
      flattenHubSubtree(nested, lifted);
    }
  }
}

function sanitizePlusTabHub(
  route: RouteRecordStringComponent,
): RouteRecordStringComponent {
  const hub: RouteRecordStringComponent = { ...route };
  delete hub.children;
  delete hub.redirect;
  return hub;
}

/**
 * Tab 容器 hub 不可挂子路由，否则 accessible 会 redirect 到首个子页
 *（插件中心→添加积分商品；分销→导出）。将子树提升到 /plus 下，hub 本身仅保留 component。
 */
function fixTabHubRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const lifted: RouteRecordStringComponent[] = [];

  function walk(route: RouteRecordStringComponent): RouteRecordStringComponent {
    if (TAB_HUB_PATHS.has(route.path)) {
      if (route.children?.length) {
        flattenHubSubtree(route.children, lifted);
      }
      return sanitizePlusTabHub(route);
    }

    const next: RouteRecordStringComponent = { ...route };
    if (next.children?.length) {
      next.children = next.children.map(walk);
    }
    return next;
  }

  const walked = routes.map(walk);

  return walked.map((route) => {
    if (route.path !== '/appsetting') {
      return route;
    }

    const children = [...(route.children ?? [])];
    const seen = new Set(children.map((child) => child.path));
    for (const child of lifted) {
      if (!child.path.startsWith('/appsetting/')) {
        continue;
      }
      if (!seen.has(child.path)) {
        children.push(child);
        seen.add(child.path);
      }
    }

    for (let i = 0; i < children.length; i++) {
      if (TAB_HUB_PATHS.has(children[i]!.path)) {
        children[i] = sanitizePlusTabHub(children[i]!);
      }
    }

    const hubIndex = children.findIndex(
      (child) => child.path === MERCHANT_APPOPEN_HUB_PATH,
    );
    if (hubIndex >= 0) {
      children[hubIndex] = {
        ...sanitizePlusTabHub(children[hubIndex]!),
        meta: {
          ...children[hubIndex]!.meta,
          hideChildrenInMenu: true,
          title: 'App设置',
        },
      };
    }

    return {
      ...route,
      children,
    };
  });
}

/** DIY 编辑器权限子路由：不参与列表页 permission 子路由提升 */
const MERCHANT_PAGE_DIY_ROUTE_PATHS = new Set([
  '/page/page/add',
  '/page/page/edit',
  '/page/page/addPage',
  '/page/page/editPage',
  '/page/center/add',
  '/page/center/edit',
]);

/** 页面列表权限子路由（v-auth 用，不可作 index 的 children） */
const MERCHANT_PAGE_LIST_PERMISSION_PATHS = new Set([
  '/page/page/delete',
  '/page/page/set',
  '/page/page/deletePage',
  '/page/page/setPage',
]);

function isListPageWithPermissionChildren(
  route: RouteRecordStringComponent,
): boolean {
  if (!route.component) {
    return false;
  }
  const kids = route.children ?? [];
  return (
    kids.length > 0 &&
    kids.every((child) => child.meta?.hideInMenu === true)
  );
}

function collectRoutePaths(
  routes: RouteRecordStringComponent[],
  paths = new Set<string>(),
): Set<string> {
  for (const route of routes) {
    paths.add(route.path);
    if (route.children?.length) {
      collectRoutePaths(route.children, paths);
    }
  }
  return paths;
}

function resolvePermissionRouteAttachPath(
  prefix: string,
  routePaths: Set<string>,
): string | null {
  if (routePaths.has(prefix)) {
    return prefix;
  }
  const parts = prefix.split('/').filter(Boolean);
  while (parts.length > 0) {
    parts.pop();
    const candidate = parts.length ? `/${parts.join('/')}` : '/';
    if (routePaths.has(candidate)) {
      return candidate;
    }
  }
  return null;
}

function mergeRouteChildren(
  route: RouteRecordStringComponent,
  attachMap: Map<string, RouteRecordStringComponent[]>,
): RouteRecordStringComponent {
  const bucket = attachMap.get(route.path);
  let next: RouteRecordStringComponent = { ...route };
  if (bucket?.length) {
    const children = [...(route.children ?? [])];
    const seen = new Set(children.map((child) => child.path));
    for (const child of bucket) {
      if (!seen.has(child.path)) {
        children.push(child);
        seen.add(child.path);
      }
    }
    next = { ...next, children };
  }
  if (next.children?.length) {
    next = {
      ...next,
      children: next.children.map((child) => mergeRouteChildren(child, attachMap)),
    };
  }
  return next;
}

/**
 * 列表页（带 component）下仅挂按钮权限子路由时，Vben accessible 会 redirect 到首个子页
 *（如 /live/index → /live/room/add）。将权限路由提升到父级，列表页不再挂 children。
 */
function fixListPagePermissionRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const lifted = new Map<string, RouteRecordStringComponent[]>();

  function walk(route: RouteRecordStringComponent): RouteRecordStringComponent {
    if (isListPageWithPermissionChildren(route)) {
      const parentPrefix =
        route.path.split('/').slice(0, -1).join('/') || '/';
      const bucket = lifted.get(parentPrefix) ?? [];
      flattenHubSubtree(route.children!, bucket);
      lifted.set(parentPrefix, bucket);
      const hub = { ...route };
      delete hub.children;
      delete hub.redirect;
      return hub;
    }

    const next: RouteRecordStringComponent = { ...route };
    if (next.children?.length) {
      next.children = next.children.map(walk);
    }
    return next;
  }

  const walked = routes.map(walk);
  if (lifted.size === 0) {
    return walked;
  }

  const routePaths = collectRoutePaths(walked);
  const attachMap = new Map<string, RouteRecordStringComponent[]>();

  for (const [prefix, bucket] of lifted.entries()) {
    const attachable = bucket.filter(
      (child) =>
        !MERCHANT_PAGE_DIY_ROUTE_PATHS.has(child.path) &&
        !MERCHANT_PAGE_LIST_PERMISSION_PATHS.has(child.path),
    );
    if (attachable.length === 0) {
      continue;
    }
    const attachPath = resolvePermissionRouteAttachPath(prefix, routePaths);
    if (!attachPath) {
      continue;
    }
    const existing = attachMap.get(attachPath) ?? [];
    attachMap.set(attachPath, [...existing, ...attachable]);
  }

  return walked.map((route) => mergeRouteChildren(route, attachMap));
}

/**
 * DIY 编辑器为列表页权限子路由，须提升到根级以便 router 注册；保持在 BasicLayout 内以显示 Tab。
 */
function promotePageDiyEditorRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const promoted: RouteRecordStringComponent[] = [];
  const seen = new Set<string>();

  function strip(route: RouteRecordStringComponent): RouteRecordStringComponent {
    const next: RouteRecordStringComponent = { ...route };
    if (!next.children?.length) {
      return next;
    }
    const keptChildren: RouteRecordStringComponent[] = [];
    for (const child of next.children) {
      const processed = strip(child);
      if (MERCHANT_PAGE_DIY_ROUTE_PATHS.has(processed.path)) {
        if (!seen.has(processed.path)) {
          seen.add(processed.path);
          promoted.push({
            ...processed,
            children: undefined,
            meta: {
              ...processed.meta,
              hideInMenu: true,
            },
          });
        }
        continue;
      }
      keptChildren.push(processed);
    }
    if (keptChildren.length > 0) {
      next.children = keptChildren;
    } else {
      delete next.children;
    }
    return next;
  }

  const stripped = routes.map(strip);
  return promoted.length > 0 ? [...stripped, ...promoted] : stripped;
}

/**
 * 页面列表 hub 不可保留 children/redirect（DIY add 等已提升为同级路由）。
 */
function sanitizePageListHubRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const promoted: RouteRecordStringComponent[] = [];
  const seen = new Set<string>();

  function walk(route: RouteRecordStringComponent): RouteRecordStringComponent {
    const next: RouteRecordStringComponent = { ...route };
    if (next.children?.length) {
      const kept: RouteRecordStringComponent[] = [];
      for (const child of next.children) {
        const processed = walk(child);
        if (MERCHANT_PAGE_LIST_PERMISSION_PATHS.has(processed.path)) {
          if (!seen.has(processed.path)) {
            seen.add(processed.path);
            promoted.push({
              ...processed,
              children: undefined,
              meta: { ...processed.meta, hideInMenu: true },
            });
          }
          continue;
        }
        kept.push(processed);
      }
      if (kept.length > 0) {
        next.children = kept;
      } else {
        delete next.children;
      }
    }
    if (MERCHANT_PAGE_LIST_HUB_PATHS.has(next.path)) {
      delete next.children;
      delete next.redirect;
    }
    return next;
  }

  const walked = routes.map(walk);
  return promoted.length > 0 ? [...walked, ...promoted] : walked;
}

function buildPlusLiveStaticRoutes(): RouteRecordStringComponent[] {
  return MERCHANT_PLUS_LIVE_ROUTE_PATHS.map((path) => ({
    name: routeNameFromPath(path),
    path,
    component: resolveLegacyComponentPath(path),
    meta: {
      hideInMenu: true,
      title: path === MERCHANT_PLUS_LIVE_HUB_PATH ? '直播' : '',
    },
  }));
}

function fixPlusHubRoutes(
  routes: RouteRecordStringComponent[],
): RouteRecordStringComponent[] {
  const lifted: RouteRecordStringComponent[] = [];

  function walk(route: RouteRecordStringComponent): RouteRecordStringComponent {
    if (TAB_HUB_PATHS.has(route.path) && route.path !== MERCHANT_APPOPEN_HUB_PATH) {
      if (route.children?.length) {
        flattenHubSubtree(route.children, lifted);
      }
      return sanitizePlusTabHub(route);
    }

    const next: RouteRecordStringComponent = { ...route };
    if (next.children?.length) {
      next.children = next.children.map(walk);
    }
    return next;
  }

  const walked = routes.map(walk);

  return walked.map((route) => {
    if (route.path !== '/plus') {
      return route;
    }

    const children = [...(route.children ?? [])];
    const seen = new Set(children.map((child) => child.path));
    for (const child of lifted) {
      if (!seen.has(child.path)) {
        children.push(child);
        seen.add(child.path);
      }
    }

    for (const child of buildPlusLiveStaticRoutes()) {
      if (!seen.has(child.path)) {
        children.push(child);
        seen.add(child.path);
      }
    }

    for (let i = 0; i < children.length; i++) {
      const childPath = children[i]!.path;
      if (
        TAB_HUB_PATHS.has(childPath) &&
        childPath !== MERCHANT_APPOPEN_HUB_PATH
      ) {
        children[i] = sanitizePlusTabHub(children[i]!);
      }
      if (childPath === MERCHANT_PLUS_LIVE_HUB_PATH) {
        children[i] = {
          ...sanitizePlusTabHub(children[i]!),
          meta: {
            ...children[i]!.meta,
            title: '直播',
          },
        };
      }
    }

    let hubIndex = children.findIndex(
      (child) => child.path === MERCHANT_PLUS_HUB_PATH,
    );
    const plusHubRoute: RouteRecordStringComponent = {
      name: routeNameFromPath(MERCHANT_PLUS_HUB_PATH),
      path: MERCHANT_PLUS_HUB_PATH,
      component: resolveLegacyComponentPath(MERCHANT_PLUS_HUB_PATH),
      meta: { hideInMenu: true, title: '插件中心' },
    };
    if (hubIndex < 0) {
      children.unshift(plusHubRoute);
    } else {
      children[hubIndex] = sanitizePlusTabHub({
        ...children[hubIndex]!,
        ...plusHubRoute,
        meta: {
          ...children[hubIndex]!.meta,
          ...plusHubRoute.meta,
        },
      });
    }

    return {
      ...route,
      redirect: MERCHANT_PLUS_HUB_PATH,
      meta: {
        ...route.meta,
        hideChildrenInMenu: true,
        title: route.meta?.title || '插件中心',
      },
      children,
    };
  });
}

function withStaticHomeMenu(menus: ShopAccessMenuItem[]) {
  const list = [...(menus || [])].map((item) =>
    item.path === '/home'
      ? {
          ...item,
          name: '我的数据',
          is_route: 1,
          component: item.component || 'views/native/home/index.vue',
        }
      : item,
  );
  if (!list.some((item) => item.path === '/home')) {
    list.unshift(STATIC_HOME_MENU);
  }
  return list;
}

export function convertShopMenusToVben(
  menus: ShopAccessMenuItem[],
): RouteRecordStringComponent[] {
  const roots = withStaticHomeMenu(menus)
    .map(convertNode)
    .filter(Boolean) as RouteRecordStringComponent[];

  // Tab 容器 hub 不可挂子路由，否则 Vben accessible 会 redirect 到错误子页
  return fixTabHubRoutes(
    sanitizePageListHubRoutes(
      fixListPagePermissionRoutes(
        promotePageDiyEditorRoutes(fixPlusHubRoutes(roots)),
      ),
    ),
  );
}
