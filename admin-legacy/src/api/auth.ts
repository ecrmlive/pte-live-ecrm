import { http } from './http';

export interface TokenPair {
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface PlatformUser {
  admin_id: number;
  account: string;
  real_name: string;
  phone: string;
  roles: string;
  level: number;
}

export interface MenuNode {
  menu_id: number;
  pid: number;
  path: string;
  icon: string;
  menu_name: string;
  route: string;
  sort: number;
  children?: MenuNode[];
}

export function login(account: string, password: string) {
  return http.post<{ token: TokenPair; user: PlatformUser }>('/auth/login', {
    account,
    password,
  });
}

export function fetchMe() {
  return http.get<PlatformUser>('/auth/me');
}

export function fetchMenus() {
  return http.get<{ menus: MenuNode[] }>('/auth/menus');
}

export function fetchPermissions() {
  return http.get<{ permissions: string[] }>('/auth/permissions');
}

export function changePassword(old_password: string, new_password: string) {
  return http.put('/auth/password', { old_password, new_password });
}
