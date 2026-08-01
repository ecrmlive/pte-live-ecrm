import type { Router } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';
import { startProgress, stopProgress } from '@vben/utils';

import { ElMessage } from 'element-plus';

import { accessRoutes, coreRouteNames } from '#/router/routes';
import { useAuthStore } from '#/store';
import { usePlatformUserStore } from '#/store/platform-user';
import {
  QIXI_PLATFORM_APP_ID,
} from '#/utils/pte-live-api';
import {
  clearEncryptedToken,
  hydrateAccessTokenFromLegacy,
  setEncryptedToken,
} from '#/utils/pte-live-token';
import {
  refreshStoredJwtSession,
  shouldRefreshStoredJwt,
} from '#/utils/jwt-session';
import {
  loadPlatformBootstrapData,
  PlatformBootstrapError,
} from '#/utils/platform-bootstrap';

import { resolveNavigationAfterAccess } from '#/utils/post-access-redirect';

import { generateAccess } from './access';
import { applySubMenuDefaultIcons } from './menu-default-icons';

const SERVICE_ERROR_PATH = '/auth/service-error';

function applyPlatformToken(token: string) {
  const accessStore = useAccessStore();
  accessStore.setAccessToken(token);
  setEncryptedToken(token);
  usePlatformUserStore().setToken(token);
}

function clearPlatformToken() {
  const accessStore = useAccessStore();
  accessStore.setAccessToken(null);
  clearEncryptedToken();
}

async function runPlatformStoredJwtRefresh(force = false) {
  const accessStore = useAccessStore();
  const token = accessStore.accessToken;
  if (!token) {
    return 'skipped' as const;
  }
  if (!force && !shouldRefreshStoredJwt()) {
    return 'skipped' as const;
  }
  return refreshStoredJwtSession({
    appId: QIXI_PLATFORM_APP_ID,
    force,
    loginPlatform: 'platform_admin',
    onClear: clearPlatformToken,
    onToken: applyPlatformToken,
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

    // 页面加载进度条
    if (!to.meta.loaded && preferences.transition.progress) {
      startProgress();
    }
    return true;
  });

  router.afterEach((to) => {
    // 记录页面是否加载,如果已经加载，后续的页面切换动画等效果不在重复执行

    loadedPaths.add(to.path);

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
        if (!accessStore.accessToken) {
          hydrateAccessTokenFromLegacy((token) => {
            applyPlatformToken(token);
          });
        }
        if (!accessStore.accessToken) {
          return true;
        }
        const refreshOutcome = await runPlatformStoredJwtRefresh(true);
        if (refreshOutcome === 'invalid' || !accessStore.accessToken) {
          return true;
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
      if (to.name === 'ServiceError') {
        return true;
      }
      return true;
    }

    if (!accessStore.accessToken) {
      hydrateAccessTokenFromLegacy((token) => {
        applyPlatformToken(token);
      });
    }

    if (accessStore.accessToken && shouldRefreshStoredJwt()) {
      const refreshOutcome = await runPlatformStoredJwtRefresh(false);
      if (refreshOutcome === 'invalid') {
        await authStore.logout(false);
        return { path: LOGIN_PATH, replace: true };
      }
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

    // 生成路由表
    try {
      // 始终刷新 session，避免 Pinia 持久化的旧 accessCodes 缺新权限导致 v-access 隐藏按钮
      const bootstrap = await loadPlatformBootstrapData();
      const userInfo = bootstrap.userInfo;
      userStore.setUserInfo(userInfo);
      accessStore.setAccessCodes(bootstrap.accessCodes);

      const userRoles = userInfo?.roles ?? [];

      const { accessibleMenus, accessibleRoutes } = await generateAccess({
        platformMenus: bootstrap.menus,
        roles: userRoles,
        router,
        routes: accessRoutes,
      });

      accessStore.setAccessMenus(applySubMenuDefaultIcons(accessibleMenus));
      accessStore.setAccessRoutes(accessibleRoutes);
      accessStore.setIsAccessChecked(true);

      return resolveNavigationAfterAccess(to, from, {
        defaultHomePath: preferences.app.defaultHomePath,
        homePath: userInfo?.homePath,
      });
    } catch (error) {
      stopProgress();
      const bootstrapError =
        error instanceof PlatformBootstrapError
          ? error
          : new PlatformBootstrapError(
              error instanceof Error ? error.message : '加载失败，请稍后重试',
              'unknown',
            );
      accessStore.setIsAccessChecked(false);

      if (bootstrapError.kind === 'auth') {
        if (to.path !== LOGIN_PATH) {
          ElMessage.error(bootstrapError.message);
        }
        await authStore.logout(false);
        return { path: LOGIN_PATH, replace: true };
      }

      return {
        path: SERVICE_ERROR_PATH,
        query: { from: encodeURIComponent(to.fullPath) },
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
