import type { UserInfo } from '@vben/types';

import { fetchPlatformSessionApi } from '#/api/core/platform-session';
import type { PlatformAccessMenuItem } from '#/utils/platform-menu';
import { formatUserFacingApiError } from '#/utils/api-error';

/** 进后台时 session 接口超时（毫秒） */
export const PLATFORM_STARTUP_TIMEOUT_MS = 12_000;

export type PlatformBootstrapErrorKind = 'auth' | 'network' | 'timeout' | 'unknown';

export class PlatformBootstrapError extends Error {
  readonly kind: PlatformBootstrapErrorKind;

  constructor(message: string, kind: PlatformBootstrapErrorKind) {
    super(message);
    this.name = 'PlatformBootstrapError';
    this.kind = kind;
  }
}

type ApiErrorLike = {
  code?: number;
  data?: { code?: number; msg?: string };
  message?: string;
  response?: { data?: { code?: number; msg?: string } };
};

export function withTimeout<T>(
  promise: Promise<T>,
  ms: number,
  _label = '请求',
): Promise<T> {
  return new Promise((resolve, reject) => {
    const timer = window.setTimeout(() => {
      reject(
        new PlatformBootstrapError(
          '内部错误，请稍后重试',
          'timeout',
        ),
      );
    }, ms);

    promise
      .then((value) => {
        window.clearTimeout(timer);
        resolve(value);
      })
      .catch((error) => {
        window.clearTimeout(timer);
        reject(error);
      });
  });
}

function classifyBootstrapError(error: unknown): PlatformBootstrapError {
  if (error instanceof PlatformBootstrapError) {
    return error;
  }

  const err = error as ApiErrorLike;
  const code = err?.response?.data?.code ?? err?.data?.code ?? err?.code;
  if (code === -1) {
    return new PlatformBootstrapError('登录已失效，请重新登录', 'auth');
  }

  const message = String(err?.message || '');
  if (message.includes('timeout') || message.includes('超时')) {
    return new PlatformBootstrapError(
      formatUserFacingApiError(error, '内部错误，请稍后重试'),
      'timeout',
    );
  }

  return new PlatformBootstrapError(
    formatUserFacingApiError(error, '内部错误，请稍后重试'),
    'network',
  );
}

export interface PlatformBootstrapData {
  accessCodes: string[];
  menus: PlatformAccessMenuItem[];
  userInfo: UserInfo;
}

/** api-platform 一次拉取菜单 + 权限码 + 用户信息（Redis 缓存） */
export async function loadPlatformBootstrapData(): Promise<PlatformBootstrapData> {
  try {
    const session = await withTimeout(
      fetchPlatformSessionApi(),
      PLATFORM_STARTUP_TIMEOUT_MS,
      '加载平台权限',
    );
    return {
      accessCodes: session.accessCodes,
      menus: session.menus ?? [],
      userInfo: session.userInfo,
    };
  } catch (error) {
    throw classifyBootstrapError(error);
  }
}

export async function loadPlatformMenusWithTimeout() {
  const { convertPlatformMenusToVben, QIXI_PLATFORM_MENU_KEY } = await import(
    '#/utils/platform-menu'
  );
  const cached = sessionStorage.getItem(QIXI_PLATFORM_MENU_KEY);
  if (cached) {
    try {
      return convertPlatformMenusToVben(JSON.parse(cached));
    } catch {
      // ignore invalid cache
    }
  }
  try {
    const session = await withTimeout(
      fetchPlatformSessionApi(),
      PLATFORM_STARTUP_TIMEOUT_MS,
      '获取菜单',
    );
    return convertPlatformMenusToVben(session.menus);
  } catch (error) {
    throw classifyBootstrapError(error);
  }
}
