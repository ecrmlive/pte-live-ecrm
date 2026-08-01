import { clearShopSessionCache } from '#/api/core/shop-session';
import { QIXI_SHOP_MENU_KEY } from '#/utils/pte-live-api';
import {
  refreshLiveJwtToken,
  type LoginPlatform,
  type RefreshTokenResult,
} from '#/utils/live-token-refresh';

/** 与后端 Unauthorized 一致：HTTP 401 或 JSON code -1 */
export function isApiUnauthorized(httpStatus?: number, bodyCode?: number) {
  return httpStatus === 401 || bodyCode === -1;
}

const JWT_FRESH_SKIP_REFRESH_SEC = 300;

let jwtIssuedFromLogin = false;
let storedJwtRefreshDone = false;
let jwtRefreshRejected = false;
let shopSessionBootstrapped = false;

/** 解析 JWT iat（不校验签名），用于判断是否为刚签发的 token */
export function parseTokenIat(token: string): number | null {
  try {
    const part = token.split('.')[1];
    if (!part) return null;
    const normalized = part.replace(/-/g, '+').replace(/_/g, '/');
    const payload = JSON.parse(atob(normalized)) as { iat?: number };
    return typeof payload.iat === 'number' ? payload.iat : null;
  } catch {
    return null;
  }
}

/** 刚登录/刚 refresh 的 token 在有效期内不必再 refresh，避免 F5 把有效 token 拉黑 */
export function isFreshLoginToken(token: string | null | undefined) {
  if (!token?.trim()) return false;
  const iat = parseTokenIat(token);
  if (!iat) return false;
  return Date.now() / 1000 - iat < JWT_FRESH_SKIP_REFRESH_SEC;
}

/** passport/login 成功后调用，本页会话内跳过 startup refresh */
export function markJwtIssuedFromLogin() {
  jwtIssuedFromLogin = true;
}

/** authLogin 已拉过 /shop/auth/session，路由守卫勿重复 bootstrap */
export function markShopSessionBootstrapped() {
  shopSessionBootstrapped = true;
}

export function resetJwtSessionState() {
  jwtIssuedFromLogin = false;
  storedJwtRefreshDone = false;
  jwtRefreshRejected = false;
  shopSessionBootstrapped = false;
  clearShopSessionCache();
}

export function markJwtRefreshRejected() {
  jwtRefreshRejected = true;
  storedJwtRefreshDone = true;
}

export function shouldRefreshStoredJwt(token?: string | null) {
  if (jwtRefreshRejected) return false;
  if (jwtIssuedFromLogin) return false;
  if (storedJwtRefreshDone) return false;
  if (isFreshLoginToken(token)) return false;
  return true;
}

export function shouldLoadShopBootstrap() {
  return !shopSessionBootstrapped;
}

/** 刷新后 sessionStorage 已有 userInfo + 非空菜单时，跳过重复 /shop/auth/session */
export function restoreShopSessionBootstrappedFromStorage() {
  if (shopSessionBootstrapped) {
    return;
  }
  try {
    const user = sessionStorage.getItem('userInfo');
    if (!user) {
      return;
    }
    const menusRaw = sessionStorage.getItem(QIXI_SHOP_MENU_KEY);
    if (!menusRaw) {
      return;
    }
    const menus = JSON.parse(menusRaw);
    if (Array.isArray(menus) && menus.length > 0) {
      shopSessionBootstrapped = true;
    }
  } catch {
    // ignore
  }
}

export type StoredJwtRefreshOutcome = 'invalid' | 'network' | 'ok' | 'skipped';

export interface StoredJwtRefreshHandlers {
  appId?: number | string;
  loginPlatform: LoginPlatform;
  onClear: () => void;
  onToken: (token: string) => void;
  token: string;
  /** 登录页带历史 token 时强制校验，不受「刚登录」影响 */
  force?: boolean;
}

let refreshInflight: Promise<StoredJwtRefreshOutcome> | null = null;

/**
 * 有本地 JWT 时：refresh 校验并续期；无效则清 token（等同 401）。
 * 同页并发去重，避免双 refresh 把 token 拉黑。
 */
export async function refreshStoredJwtSession(
  handlers: StoredJwtRefreshHandlers,
): Promise<StoredJwtRefreshOutcome> {
  const token = handlers.token?.trim();
  if (!token) {
    return 'skipped';
  }
  if (jwtRefreshRejected) {
    handlers.onClear();
    return 'invalid';
  }
  if (!handlers.force && !shouldRefreshStoredJwt(token)) {
    return 'skipped';
  }

  if (refreshInflight) {
    return refreshInflight;
  }

  refreshInflight = (async () => {
    try {
      const result: RefreshTokenResult = await refreshLiveJwtToken({
        loginPlatform: handlers.loginPlatform,
        token,
        appId: handlers.appId,
      });

      if (result.clearToken) {
        handlers.onClear();
        markJwtRefreshRejected();
        return 'invalid' as const;
      }
      if (result.ok && result.token) {
        handlers.onToken(result.token);
        storedJwtRefreshDone = true;
        jwtRefreshRejected = false;
        return 'ok' as const;
      }
      storedJwtRefreshDone = true;
      return 'network' as const;
    } finally {
      refreshInflight = null;
    }
  })();

  return refreshInflight;
}
