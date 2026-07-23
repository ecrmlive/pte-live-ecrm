import {
  refreshLiveJwtToken,
  type LoginPlatform,
  type RefreshTokenResult,
} from '#/utils/live-token-refresh';

/** 与后端 Unauthorized 一致：HTTP 401 或 JSON code -1 */
export function isApiUnauthorized(httpStatus?: number, bodyCode?: number) {
  return httpStatus === 401 || bodyCode === -1;
}

let jwtIssuedFromLogin = false;
let storedJwtRefreshDone = false;

/** passport/login 成功后调用，本页会话内跳过 startup refresh */
export function markJwtIssuedFromLogin() {
  jwtIssuedFromLogin = true;
}

export function resetJwtSessionState() {
  jwtIssuedFromLogin = false;
  storedJwtRefreshDone = false;
}

export function shouldRefreshStoredJwt() {
  return !jwtIssuedFromLogin && !storedJwtRefreshDone;
}

export type StoredJwtRefreshOutcome = 'invalid' | 'network' | 'ok' | 'skipped';

export interface StoredJwtRefreshHandlers {
  appId?: number;
  loginPlatform: LoginPlatform;
  onClear: () => void;
  onToken: (token: string) => void;
  token: string;
  /** 登录页带历史 token 时强制校验，不受「刚登录」影响 */
  force?: boolean;
}

/**
 * 有本地 JWT 时：refresh 校验并续期；无效则清 token（等同 401）。
 * 刚完成 passport/login 且非 force 时跳过。
 */
export async function refreshStoredJwtSession(
  handlers: StoredJwtRefreshHandlers,
): Promise<StoredJwtRefreshOutcome> {
  const token = handlers.token?.trim();
  if (!token) {
    return 'skipped';
  }
  if (!handlers.force && !shouldRefreshStoredJwt()) {
    return 'skipped';
  }

  if (!handlers.force) {
    storedJwtRefreshDone = true;
  }

  const result: RefreshTokenResult = await refreshLiveJwtToken({
    loginPlatform: handlers.loginPlatform,
    token,
    appId: handlers.appId,
  });

  if (result.clearToken) {
    handlers.onClear();
    return 'invalid';
  }
  if (result.ok && result.token) {
    handlers.onToken(result.token);
    return 'ok';
  }
  return 'network';
}
