import { updatePreferences } from '@vben/preferences';

/** 产品名（登录页 / 未登录壳层），与 i18n authentication.pageTitle 一致 */
export const MERCHANT_PRODUCT_NAME = '栖息多商户·商户';

/**
 * 管理后台的侧栏、页签和浏览器标题统一使用产品品牌。
 * 商城名称属于业务数据，不能覆盖后台产品名；否则历史租户名称会把
 * 商户店铺名不能覆盖后台产品名，避免把本项目后台误显示为历史项目名称。
 */
export function resolveShopDisplayName(_name?: null | string) {
  return MERCHANT_PRODUCT_NAME;
}

/** 登录页 / 退出后：恢复产品名，可选更新 Logo（不用 sys_config.shop_name 覆盖标题） */
export function applyLoginPageBranding(logoUrl?: null | string) {
  const patch: Parameters<typeof updatePreferences>[0] = {
    app: { name: MERCHANT_PRODUCT_NAME },
  };
  const logo = String(logoUrl ?? '').trim();
  if (logo) {
    patch.logo = { source: logo };
  }
  updatePreferences(patch);
  if (typeof document !== 'undefined') {
    document.title = MERCHANT_PRODUCT_NAME;
  }
}

/**
 * 登录成功后：保持统一产品品牌；租户名称仅用于业务内容，不能覆盖后台名称。
 */
export function applyTenantAppBranding(
  _name?: null | string,
  logoUrl?: null | string,
) {
  const displayName = MERCHANT_PRODUCT_NAME;

  const patch: Parameters<typeof updatePreferences>[0] = {
    app: { name: displayName },
  };
  const logo = String(logoUrl ?? '').trim();
  if (logo) {
    patch.logo = { source: logo };
  }
  updatePreferences(patch);
  if (typeof document !== 'undefined') {
    document.title = displayName;
  }
  return displayName;
}

/** @deprecated 使用 applyTenantAppBranding / applyLoginPageBranding */
export const DEFAULT_MERCHANT_PLATFORM_NAME = MERCHANT_PRODUCT_NAME;
export const applyMerchantAppBranding = applyTenantAppBranding;
