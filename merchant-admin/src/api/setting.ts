import { http } from '@/api/http';

export interface ShopProfile {
  mer_id: number;
  mer_name: string;
  real_name: string;
  mer_phone: string;
  mer_address: string;
  mer_info: string;
}

export interface Staff {
  service_id: number;
  mer_id: number;
  account: string;
  nickname: string;
  phone: string;
  status: number;
  is_open: number;
  is_verify: number;
  is_goods: number;
}

export interface MerchantAdmin {
  merchant_admin_id: number;
  mer_id: number;
  account: string;
  real_name: string;
  phone?: string;
  roles: string;
  status: number;
  level: number;
}

export interface SystemRole {
  role_id: number;
  role_name: string;
  rules: string;
  status: number;
  mer_id: number;
}

export interface MenuNode {
  menu_id: number;
  pid: number;
  path: string;
  menu_name: string;
  route: string;
  sort: number;
  is_menu?: number;
  children?: MenuNode[];
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchShop() {
  return http.get<ShopProfile>('/setting/shop');
}

export function updateShop(data: Record<string, unknown>) {
  return http.put<ShopProfile>('/setting/shop', data);
}

export function fetchStaff(params: Record<string, unknown>) {
  return http.get<PageResult<Staff>>('/setting/staff', { params });
}

export function createStaff(data: Record<string, unknown>) {
  return http.post<Staff>('/setting/staff', data);
}

export function updateStaff(id: number, data: Record<string, unknown>) {
  return http.put<Staff>(`/setting/staff/${id}`, data);
}

export function fetchAdmins(params: Record<string, unknown>) {
  return http.get<PageResult<MerchantAdmin>>('/setting/admins', { params });
}

export function createAdmin(data: Record<string, unknown>) {
  return http.post<MerchantAdmin>('/setting/admins', data);
}

export function updateAdmin(id: number, data: Record<string, unknown>) {
  return http.put<MerchantAdmin>(`/setting/admins/${id}`, data);
}

export function fetchRoles(params: Record<string, unknown> = {}) {
  return http.get<PageResult<SystemRole>>('/setting/roles', { params });
}

export function createRole(data: Record<string, unknown>) {
  return http.post<SystemRole>('/setting/roles', data);
}

export function updateRole(id: number, data: Record<string, unknown>) {
  return http.put<SystemRole>(`/setting/roles/${id}`, data);
}

export function fetchMenuTree() {
  return http.get<{ list: MenuNode[] }>('/setting/menus/tree');
}
