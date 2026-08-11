import type {
  ComponentRecordType,
  GenerateMenuAndRoutesOptions,
} from '@vben/types';

import { generateAccessible } from '@vben/access';
import { preferences } from '@vben/preferences';
import { generateMenus } from '@vben/utils';

import { fetchPlatformSessionApi } from '#/api/core/platform-session';
import { BasicLayout, IFrameView } from '#/layouts';
import {
  convertPlatformMenusToVben,
  extractPlatformRouterRoutes,
  QIXI_PLATFORM_MENU_KEY,
  type PlatformAccessMenuItem,
} from '#/utils/platform-menu';
import {
  PlatformBootstrapError,
  PLATFORM_STARTUP_TIMEOUT_MS,
  withTimeout,
} from '#/utils/platform-bootstrap';

import { STANDALONE_ROUTE_PATHS } from './routes/standalone';

const forbiddenComponent = () => import('#/views/_core/fallback/forbidden.vue');

async function loadPlatformMenuTree(
  preloaded?: PlatformAccessMenuItem[],
): Promise<PlatformAccessMenuItem[]> {
  if (preloaded?.length) {
    return preloaded;
  }

  const session = await withTimeout(
    fetchPlatformSessionApi(),
    PLATFORM_STARTUP_TIMEOUT_MS,
    '获取菜单',
  );
  sessionStorage.setItem(
    QIXI_PLATFORM_MENU_KEY,
    JSON.stringify(session.menus || []),
  );
  return session.menus || [];
}

type GenerateAccessOptions = GenerateMenuAndRoutesOptions & {
  platformMenus?: PlatformAccessMenuItem[];
};

async function generateAccess(options: GenerateAccessOptions) {
  const pageMap: ComponentRecordType = Object.fromEntries(
    Object.entries(import.meta.glob('../views/**/*.vue')).filter(
      ([path]) => !path.includes('/native/'),
    ),
  );

  const layoutMap: ComponentRecordType = {
    BasicLayout,
    IFrameView,
  };

  let menuTree;
  try {
    const rawMenus = await loadPlatformMenuTree(options.platformMenus);
    menuTree = convertPlatformMenusToVben(rawMenus);
  } catch (error) {
    if (error instanceof PlatformBootstrapError) {
      throw error;
    }
    throw error;
  }

  // DIY / 数据大屏等已在 standalone 注册为无 BasicLayout 页；
  // 菜单树仍保留入口，但动态路由不再挂到侧栏布局下。
  const routerRoutes = extractPlatformRouterRoutes(menuTree).filter(
    (route) => !STANDALONE_ROUTE_PATHS.has(String(route.path || '')),
  );

  // 图文新增/编辑：侧栏不展示，从列表跳入（activePath 高亮「图文管理」）
  routerRoutes.push(
    {
      name: 'AppWechatNewsCreate',
      path: '/app/wechat/newsCategory/create',
      component: '../views/ecrm/app/wechat-news-save.vue',
      meta: {
        activePath: '/app/wechat/newsCategory',
        hideInMenu: true,
        title: '新增图文',
      },
    },
    {
      name: 'AppWechatNewsEdit',
      path: '/app/wechat/newsCategory/edit/:id',
      component: '../views/ecrm/app/wechat-news-save.vue',
      meta: {
        activePath: '/app/wechat/newsCategory',
        hideInMenu: true,
        title: '编辑图文',
      },
    },
  );

  const result = await generateAccessible(preferences.app.accessMode, {
    ...options,
    fetchMenuListAsync: async () => routerRoutes,
    forbiddenComponent,
    layoutMap,
    pageMap,
  });

  return {
    ...result,
    accessibleMenus: generateMenus(menuTree, options.router),
  };
}

export { generateAccess };
