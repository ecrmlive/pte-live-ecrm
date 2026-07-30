import type { Router } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { useAccessStore, useTabbarStore, useUserStore } from '@vben/stores';
import { startProgress, stopProgress } from '@vben/utils';

import { ElMessage } from 'element-plus';

import { accessRoutes, coreRouteNames } from '#/router/routes';
import { mapShopSessionUser } from '#/api/core/shop-session';
import {
  hydrateMerchantSession,
} from '#/adapters/merchant-session-bridge';
import { useAuthStore } from '#/store';
import { STANDALONE_ROUTE_NAMES } from '#/utils/qixi-live-api';
import { MERCHANT_PLUS_HUB_PATH } from '#/utils/qixi-live-menu';
import { PLUS_PLUGIN_ENTRY_ALIASES, resolvePlusProductHubTabRedirect } from '#/utils/plus-navigation';
import { resolveNavigationAfterAccess } from '#/utils/post-access-redirect';
import {
  hydrateAccessTokenFromLegacy,
  setEncryptedToken,
  syncAuthFromOpener,
  syncLegacyUserSession,
  clearEncryptedToken,
  getLegacyUserInfo,
  resolveMerchantAccessToken,
} from '#/utils/qixi-live-token';
import { resolveShopAppId } from '#/utils/qixi-live-shop-app-id';
import {
  markShopSessionBootstrapped,
  refreshStoredJwtSession,
  restoreShopSessionBootstrappedFromStorage,
  shouldLoadShopBootstrap,
  shouldRefreshStoredJwt,
  isFreshLoginToken,
} from '#/utils/jwt-session';
import {
  loadShopBootstrapData,
  ShopBootstrapError,
} from '#/utils/shop-bootstrap';

import { resolveListModalRedirect, markListModalIntent } from '#/utils/list-modal-route';

import { generateAccess } from './access';

function applyMerchantToken(token: string) {
  const accessStore = useAccessStore();
  accessStore.setAccessToken(token);
  setEncryptedToken(token);
}

function clearMerchantToken() {
  const accessStore = useAccessStore();
  accessStore.setAccessToken(null);
  clearEncryptedToken();
}

function syncMerchantTokenFromStorage() {
  const accessStore = useAccessStore();
  return resolveMerchantAccessToken(
    () => accessStore.accessToken,
    (token) => applyMerchantToken(token),
  );
}

function restoreMerchantUserFromStorage() {
  restoreShopSessionBootstrappedFromStorage();
  const userStore = useUserStore();
  if (userStore.userInfo?.username) {
    return;
  }
  const legacy = getLegacyUserInfo();
  if (!legacy?.userName && !legacy?.shopName) {
    return;
  }
  userStore.setUserInfo(
    mapShopSessionUser({
      store_app_id: legacy.store_app_id,
      homePath: '/home',
      logoUrl: legacy.logoUrl,
      roles: ['merchant_user'],
      shop_name: legacy.shopName,
      user_name: legacy.userName,
      version: legacy.version,
    }),
  );
  hydrateMerchantSession();
}

function hydrateMerchantTokenFromStorage() {
  syncAuthFromOpener();
  hydrateAccessTokenFromLegacy((token) => {
    applyMerchantToken(token);
  });
  syncMerchantTokenFromStorage();
}

async function runMerchantStoredJwtRefresh(
  userInfo: Record<string, unknown> | null | undefined,
  force = false,
) {
  const accessStore = useAccessStore();
  const token =
    syncMerchantTokenFromStorage() || accessStore.accessToken?.trim() || null;
  if (!token) {
    return 'skipped' as const;
  }
  if (!accessStore.accessToken) {
    applyMerchantToken(token);
  }
  if (!force && !shouldRefreshStoredJwt(token)) {
    return 'skipped' as const;
  }
  return refreshStoredJwtSession({
    appId: resolveShopAppId(userInfo),
    force,
    loginPlatform: 'merchant_admin',
    onClear: clearMerchantToken,
    onToken: applyMerchantToken,
    token,
  });
}

/**
 * 通用守卫配置
 * @param router
 */
function setupCommonGuard(router: Router) {
  // 记录已经加载的页面
  const loadedPaths = new Set<string>();

  router.beforeEach((to) => {
    to.meta.loaded = loadedPaths.has(to.path);

    const listModalRedirect = resolveListModalRedirect(to);
    if (listModalRedirect) {
      return listModalRedirect;
    }

    const plusHubTabRedirect = resolvePlusProductHubTabRedirect(to);
    if (plusHubTabRedirect) {
      return plusHubTabRedirect;
    }

    if (to.path === '/live/index' && to.query.openAdd === '1') {
      markListModalIntent({ action: 'live-add' });
    }
    if (to.path === '/product/product/index' && to.query.openAdd === '1') {
      markListModalIntent({ action: 'product-add' });
    }

    // 侧栏「插件中心」→ 卡片页，避免 Vben 链式 redirect 到积分商品 add 等子页
    if (to.path === '/plus') {
      return { path: MERCHANT_PLUS_HUB_PATH, replace: true };
    }

    const plusEntry = PLUS_PLUGIN_ENTRY_ALIASES[to.path];
    if (plusEntry && plusEntry !== to.path) {
      return { path: plusEntry, query: to.query, hash: to.hash, replace: true };
    }

    // 页面加载进度条
    if (!to.meta.loaded && preferences.transition.progress) {
      startProgress();
    }
    return true;
  });

  router.afterEach((to) => {
    // 记录页面是否加载,如果已经加载，后续的页面切换动画等效果不在重复执行

    loadedPaths.add(to.path);

    // 插件中心等：路由 meta.title 同步到顶栏 Tab（页面内 setTabTitle 可覆盖）
    const routeTitle = to.meta?.title;
    if (
      typeof routeTitle === 'string' &&
      routeTitle.trim() &&
      !to.meta?.newTabTitle &&
      (to.path.startsWith('/plus/') || to.path === MERCHANT_PLUS_HUB_PATH)
    ) {
      const tabbarStore = useTabbarStore();
      void tabbarStore.setTabTitle(to, routeTitle.trim());
    }

    // 关闭页面加载进度条
    if (preferences.transition.progress) {
      stopProgress();
    }
  });
}

/**
 * 权限访问守卫配置
 * @param router
 */
function setupAccessGuard(router: Router) {
  router.beforeEach(async (to, from) => {
    const accessStore = useAccessStore();
    const userStore = useUserStore();
    const authStore = useAuthStore();

    // 基本路由，这些路由不需要进入权限拦截
    if (coreRouteNames.includes(to.name as string)) {
      if (to.path === LOGIN_PATH) {
        if (authStore.loginLoading) {
          return true;
        }
        hydrateMerchantTokenFromStorage();
        if (!accessStore.accessToken) {
          return true;
        }
        if (!isFreshLoginToken(accessStore.accessToken)) {
          const refreshOutcome = await runMerchantStoredJwtRefresh(
            userStore.userInfo as Record<string, unknown>,
            true,
          );
          if (refreshOutcome === 'invalid' || !accessStore.accessToken) {
            return true;
          }
        }
        const redirectTarget =
          (to.query?.redirect as string) ||
          userStore.userInfo?.homePath ||
          preferences.app.defaultHomePath;
        return {
          path: decodeURIComponent(redirectTarget),
          replace: true,
        };
      }
      return true;
    }

    hydrateMerchantTokenFromStorage();
    restoreMerchantUserFromStorage();

    if (!accessStore.accessToken) {
      const resolved = syncMerchantTokenFromStorage();
      if (!resolved) {
        // accessToken 检查
        if (to.meta.ignoreAccess) {
          return true;
        }
        if (to.fullPath !== LOGIN_PATH) {
          return {
            path: LOGIN_PATH,
            query:
              to.fullPath === preferences.app.defaultHomePath
                ? {}
                : { redirect: encodeURIComponent(to.fullPath) },
            replace: true,
          };
        }
        return to;
      }
    }

    if (accessStore.accessToken && shouldRefreshStoredJwt(accessStore.accessToken)) {
      const refreshOutcome = await runMerchantStoredJwtRefresh(
        userStore.userInfo as Record<string, unknown>,
        false,
      );
      if (refreshOutcome === 'invalid') {
        await authStore.logout(false);
        return { path: LOGIN_PATH, replace: true };
      }
    }

    // 直播中控等独立页：有 token 即可，不阻塞菜单生成逻辑
    if (STANDALONE_ROUTE_NAMES.has(String(to.name))) {
      if (!accessStore.accessToken) {
        return {
          path: LOGIN_PATH,
          query: { redirect: encodeURIComponent(to.fullPath) },
          replace: true,
        };
      }
      if (!userStore.userInfo) {
        try {
          await authStore.fetchUserInfo();
        } catch {
          // 沿用 session 中的 userInfo
        }
      }
      return true;
    }

    // accessToken 检查
    if (!accessStore.accessToken) {
      // 明确声明忽略权限访问权限，则可以访问
      if (to.meta.ignoreAccess) {
        return true;
      }

      // 没有访问权限，跳转登录页面
      if (to.fullPath !== LOGIN_PATH) {
        return {
          path: LOGIN_PATH,
          // 如不需要，直接删除 query
          query:
            to.fullPath === preferences.app.defaultHomePath
              ? {}
              : { redirect: encodeURIComponent(to.fullPath) },
          // 携带当前跳转的页面，登录后重新跳转该页面
          replace: true,
        };
      }
      return to;
    }

    // 是否已经生成过动态路由（HMR 后 root.children 可能为空，需重新生成）
    if (accessStore.isAccessChecked) {
      const rootRoute = router.getRoutes().find((item) => item.path === '/');
      const hasAccessRoutes = (rootRoute?.children?.length ?? 0) > 0;
      if (hasAccessRoutes) {
        return true;
      }
      accessStore.setIsAccessChecked(false);
    }

    try {
      let userInfo = userStore.userInfo;
      let accessCodes = accessStore.accessCodes;

      const needBootstrap =
        shouldLoadShopBootstrap() && !userInfo?.username;
      if (needBootstrap) {
        const bootstrap = await loadShopBootstrapData();
        userInfo = bootstrap.userInfo;
        accessCodes = bootstrap.accessCodes;
        userStore.setUserInfo(userInfo);
        accessStore.setAccessCodes(accessCodes);
        markShopSessionBootstrapped();
        syncLegacyUserSession({
	        store_app_id: (userInfo as any).store_app_id,
          logoUrl: (userInfo as any).logoUrl,
          shopName: (userInfo as any).shopName || userInfo.realName,
          userName: userInfo.username,
        });
        hydrateMerchantSession();
      }

      const userRoles = userInfo?.roles ?? [];

      const { accessibleMenus, accessibleRoutes } = await generateAccess({
        roles: userRoles,
        router,
        routes: accessRoutes,
      });

      accessStore.setAccessMenus(accessibleMenus);
      accessStore.setAccessRoutes(accessibleRoutes);
      accessStore.setIsAccessChecked(true);

      return resolveNavigationAfterAccess(to, from, {
        defaultHomePath: preferences.app.defaultHomePath,
        homePath: userInfo?.homePath,
      });
    } catch (error) {
      stopProgress();
      const bootstrapError =
        error instanceof ShopBootstrapError
          ? error
          : new ShopBootstrapError(
              error instanceof Error ? error.message : '加载失败，请稍后重试',
              'unknown',
            );
      accessStore.setIsAccessChecked(false);

      if (bootstrapError.kind === 'auth') {
        await authStore.logout(false);
        return { path: LOGIN_PATH, replace: true };
      }

      if (accessStore.accessToken) {
        try {
          const userRoles = userStore.userInfo?.roles ?? [];
          const { accessibleMenus, accessibleRoutes } = await generateAccess({
            roles: userRoles,
            router,
            routes: accessRoutes,
          });
          accessStore.setAccessMenus(accessibleMenus);
          accessStore.setAccessRoutes(accessibleRoutes);
          accessStore.setIsAccessChecked(true);
          ElMessage.warning(bootstrapError.message);
          return true;
        } catch {
          // fall through
        }
      }

      return {
        path: LOGIN_PATH,
        query: { redirect: encodeURIComponent(to.fullPath) },
        replace: true,
      };
    }
  });
}

/**
 * 项目守卫配置
 * @param router
 */
function createRouterGuard(router: Router) {
  /** 通用 */
  setupCommonGuard(router);
  /** 权限访问 */
  setupAccessGuard(router);
}

export { createRouterGuard };
