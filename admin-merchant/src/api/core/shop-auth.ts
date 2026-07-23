import { requestClient } from '#/api/request';

export interface PaginatedList<T> {
  data: T[];
  total: number;
}

export interface ShopRoleItem {
  create_time: string;
  role_id: number;
  role_name_h1: string;
  sort: number;
}

export interface ShopRoleListResult {
  list: ShopRoleItem[];
}

export async function getShopRoleListApi() {
  return requestClient.post<ShopRoleListResult>('/shop/auth.role/index', {});
}

export async function deleteShopRoleApi(roleId: number) {
  return requestClient.post('/shop/auth.role/delete', { role_id: roleId });
}

export interface ShopAccessMenuItem {
  access_id: number;
  children?: ShopAccessMenuItem[];
  name: string;
}

export interface ShopRoleFormModel {
  access_id?: number[];
  parent_id?: number | string;
  role_id?: number;
  role_name: string;
  sort: number;
}

export interface ShopRoleAddInfoResult {
  menu: ShopAccessMenuItem[];
  roleList: ShopRoleItem[];
}

export interface ShopRoleEditInfoResult extends ShopRoleAddInfoResult {
  model: ShopRoleFormModel;
  select_menu: number[];
}

export async function getShopRoleAddInfoApi() {
  return requestClient.get<ShopRoleAddInfoResult>('/shop/auth.role/add');
}

export async function addShopRoleApi(form: ShopRoleFormModel) {
  return requestClient.post('/shop/auth.role/add', {
    params: JSON.stringify(form),
  });
}

export async function getShopRoleEditInfoApi(roleId: number) {
  return requestClient.get<ShopRoleEditInfoResult>('/shop/auth.role/edit', {
    params: { role_id: roleId },
  });
}

export async function editShopRoleApi(roleId: number, form: ShopRoleFormModel) {
  return requestClient.post('/shop/auth.role/edit', {
    role_id: roleId,
    params: JSON.stringify(form),
  });
}

export interface ShopAdminUserItem {
  create_time: string;
  is_super: number;
  shop_user_id: number;
  userRole?: Array<{ role: { role_name: string } }>;
  user_name: string;
}

export interface ShopAdminUserRoleOption {
  role_id: number;
  role_name_h1: string;
}

export interface ShopAdminUserListResult {
  list: PaginatedList<ShopAdminUserItem>;
  roleList: ShopAdminUserRoleOption[];
}

export async function getShopAdminUserListApi(params: {
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<ShopAdminUserListResult>(
    '/shop/auth.user/index',
    params,
  );
}

export async function deleteShopAdminUserApi(shopUserId: number) {
  return requestClient.post('/shop/auth.user/delete', {
    shop_user_id: shopUserId,
  });
}

export interface ShopAdminUserForm {
  access_id?: number[];
  confirm_password?: string;
  password?: string;
  real_name: string;
  role_id?: number[];
  shop_user_id?: number;
  user_name: string;
}

export interface ShopAdminUserEditInfo {
  info: ShopAdminUserForm;
  roleList: ShopAdminUserRoleOption[];
  role_arr: number[];
}

export async function addShopAdminUserApi(form: ShopAdminUserForm) {
  return requestClient.post<{ msg?: string }>('/shop/auth.user/add', form);
}

export async function getShopAdminUserEditInfoApi(shopUserId: number) {
  return requestClient.get<ShopAdminUserEditInfo>('/shop/auth.user/edit', {
    params: { shop_user_id: shopUserId },
  });
}

export async function editShopAdminUserApi(form: ShopAdminUserForm) {
  return requestClient.post<{ msg?: string }>('/shop/auth.user/edit', form);
}

export interface ShopLoginLogItem {
  app_id?: number;
  create_time: string;
  ip: string;
  login_log_id: number;
  result: string;
  username: string;
}

export interface ShopLoginLogListResult {
  list: PaginatedList<ShopLoginLogItem>;
}

export async function getShopLoginLogListApi(params: {
  list_rows?: number;
  page?: number;
  username?: string;
}) {
  return requestClient.post<ShopLoginLogListResult>(
    '/shop/auth.loginlog/index',
    params,
  );
}

export interface ShopOptLogItem {
  browser: string;
  create_time: string;
  ip: string;
  opt_log_id: number;
  real_name: string;
  request_type?: string;
  shop_user_id?: number;
  title: string;
  url: string;
  user_name: string;
  agent?: string;
  content?: string;
}

export interface ShopOptLogListResult {
  list: PaginatedList<ShopOptLogItem>;
}

export async function getShopOptLogListApi(params: {
  list_rows?: number;
  page?: number;
  username?: string;
}) {
  return requestClient.post<ShopOptLogListResult>(
    '/shop/auth.optlog/index',
    params,
  );
}
