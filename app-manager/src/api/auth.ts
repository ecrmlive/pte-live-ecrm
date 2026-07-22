import { http } from './http';

export interface ManagerUser {
  service_id: number;
  mer_id: number;
  mer_name?: string;
  account: string;
  nickname: string;
  is_verify: number;
  is_goods: number; // 1=可发货
}

export function login(account: string, password: string) {
  return http.post<{
    token: { access_token: string; refresh_token: string; expires_in: number };
    user: ManagerUser;
  }>('/auth/login', { account, password });
}

export function fetchMe() {
  return http.get<ManagerUser>('/auth/me');
}
