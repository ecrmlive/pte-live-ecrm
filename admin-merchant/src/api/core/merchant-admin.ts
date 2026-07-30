import { requestClient } from '#/api/request';

export interface MerchantRoleOption {
  mer_id: number;
  role_id: number;
  role_name: string;
  rules: string;
  status: number;
}

export interface MerchantMenuNode {
  children?: MerchantMenuNode[];
  is_menu: number;
  menu_id: number;
  menu_name: string;
  path: string;
}

export interface MerchantRoleSaveInput {
  menu_ids: number[];
  role_name: string;
  status: number;
}

export interface MerchantAdmin {
  account: string;
  level: number;
  merchant_admin_id: number;
  phone: string;
  real_name: string;
  roles: string;
  status: number;
}

export interface MerchantAdminSaveInput {
  account?: string;
  password?: string;
  phone?: string;
  real_name?: string;
  roles?: string;
  status?: number;
}

interface PageResult<T> { limit: number; list: T[]; page: number; total: number }

export function listMerchantAdminsApi(params: { limit?: number; page?: number }) {
  return requestClient.get<PageResult<MerchantAdmin>>('/setting/admins', { params });
}

export function createMerchantAdminApi(data: MerchantAdminSaveInput) {
  return requestClient.post<MerchantAdmin>('/setting/admins', data);
}

export function updateMerchantAdminApi(id: number, data: MerchantAdminSaveInput) {
  return requestClient.put<MerchantAdmin>(`/setting/admins/${id}`, data);
}

export function listMerchantRoleOptionsApi() {
  return requestClient.get<PageResult<MerchantRoleOption>>('/setting/roles', { params: { page: 1, limit: 100 } });
}

export function createMerchantRoleApi(data: MerchantRoleSaveInput) {
  return requestClient.post<MerchantRoleOption>('/setting/roles', data);
}

export function updateMerchantRoleApi(id: number, data: MerchantRoleSaveInput) {
  return requestClient.put<MerchantRoleOption>(`/setting/roles/${id}`, data);
}

export function getMerchantMenuTreeApi() {
  return requestClient.get<{ list: MerchantMenuNode[] }>('/setting/menus/tree');
}
