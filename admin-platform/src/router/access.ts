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

  const routerRoutes = extractPlatformRouterRoutes(menuTree);

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
