import { http } from '@/api/http';

export interface PlatformAdmin {
  admin_id: number;
  account: string;
  real_name: string;
  phone: string;
  roles: string;
  status: number;
  level: number;
}

export interface SystemRole {
  role_id: number;
  role_name: string;
  rules: string;
  status: number;
}

export interface SystemMenu {
  menu_id: number;
  pid: number;
  path: string;
  menu_name: string;
  route: string;
  sort: number;
  is_show: number;
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

export function fetchAdmins(params: Record<string, unknown>) {
  return http.get<PageResult<PlatformAdmin>>('/setting/admins', { params });
}

export function createAdmin(data: Record<string, unknown>) {
  return http.post<PlatformAdmin>('/setting/admins', data);
}

export function updateAdmin(id: number, data: Record<string, unknown>) {
  return http.put<PlatformAdmin>(`/setting/admins/${id}`, data);
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

export function fetchMenus() {
  return http.get<{ list: SystemMenu[] }>('/setting/menus');
}

export function fetchMenuTree() {
  return http.get<{ list: MenuNode[] }>('/setting/menus/tree');
}

export function updateMenu(id: number, data: Record<string, unknown>) {
  return http.put<SystemMenu>(`/setting/menus/${id}`, data);
}
