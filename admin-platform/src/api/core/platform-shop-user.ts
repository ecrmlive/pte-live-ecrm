import { apiLiveAdminPost } from '#/utils/api-live-admin';

export interface ShopUserRoleOption {
  role_id: number;
  role_name: string;
  role_name_h1?: string;
}

export interface ShopUserRoleRef {
  role_id: number;
  role: { role_id: number; role_name: string };
}

export interface ShopUserRow {
  create_time?: string;
  is_status?: number;
  is_super?: number;
  real_name?: string;
  shop_user_id: number;
  userRole?: ShopUserRoleRef[];
  user_name: string;
}

export interface ShopUserForm {
  app_id: number;
  confirm_password?: string;
  password?: string;
  real_name: string;
  role_id?: number[];
  shop_user_id?: number;
  user_name: string;
}

const PlatformShopUserApi = {
  index(
    data: { app_id: number; list_rows?: number; page?: number },
    errorback?: boolean,
  ) {
    return apiLiveAdminPost<{
      list: { data: ShopUserRow[]; total: number };
      roleList: ShopUserRoleOption[];
    }>('/admin/platform/shop-user/index', data, errorback);
  },
  addInfo(appId: number, errorback?: boolean) {
    return apiLiveAdminPost<{ roleList: ShopUserRoleOption[] }>(
      '/admin/platform/shop-user/add-info',
      { app_id: appId },
      errorback,
    );
  },
  add(data: ShopUserForm, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/shop-user/add', data, errorback);
  },
  editInfo(appId: number, shopUserId: number, errorback?: boolean) {
    return apiLiveAdminPost<{
      info: ShopUserForm;
      roleList: ShopUserRoleOption[];
      role_arr: number[];
    }>(
      '/admin/platform/shop-user/edit-info',
      { app_id: appId, shop_user_id: shopUserId },
      errorback,
    );
  },
  edit(data: ShopUserForm, errorback?: boolean) {
    return apiLiveAdminPost('/admin/platform/shop-user/edit', data, errorback);
  },
  delete(appId: number, shopUserId: number, errorback?: boolean) {
    return apiLiveAdminPost(
      '/admin/platform/shop-user/delete',
      { app_id: appId, shop_user_id: shopUserId },
      errorback,
    );
  },
  setStatus(
    appId: number,
    shopUserId: number,
    isStatus: 0 | 1,
    errorback?: boolean,
  ) {
    return apiLiveAdminPost(
      '/admin/platform/shop-user/set-status',
      { app_id: appId, shop_user_id: shopUserId, is_status: isStatus },
      errorback,
    );
  },
};

export default PlatformShopUserApi;

export function shopUserRoleNames(row: ShopUserRow): string[] {
  const names = (row.userRole ?? [])
    .map((item) => item.role?.role_name)
    .filter((name): name is string => Boolean(name));
  return [...new Set(names)];
}
