import useMerchantSessionStore from '#/store/merchant-session';

/** 将 Vben 登录写入的 session 同步到 merchantSession（中控 / live-api） */
export function hydrateMerchantSession() {
  useMerchantSessionStore().hydrate();
}

/** 登出时清理 merchantSession */
export function resetMerchantSession() {
  useMerchantSessionStore().afterLogout();
}

/** @deprecated 使用 resetMerchantSession */
export const resetLegacyUserStore = resetMerchantSession;
