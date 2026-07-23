import { http } from '@/api/http';

export interface BannerItem {
  id: number;
  title: string;
  image?: string;
  url?: string;
}

export interface MenuItem {
  id: number;
  name: string;
  icon?: string;
  url?: string;
}

export interface DiyPage {
  id: number;
  name: string;
  title: string;
  template_name: string;
  value: string;
  status: number;
  is_default: number;
  parsed?: { banners: BannerItem[]; menus: MenuItem[] };
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchDiyPages(params: Record<string, unknown>) {
  return http.get<PageResult<DiyPage>>('/diy/pages', { params });
}

export function createDiyPage(data: Record<string, unknown>) {
  return http.post<DiyPage>('/diy/pages', data);
}

export function updateDiyPage(id: number, data: Record<string, unknown>) {
  return http.put<DiyPage>(`/diy/pages/${id}`, data);
}

export function activateDiyPage(id: number) {
  return http.post<DiyPage>(`/diy/pages/${id}/active`, {});
}

export function deleteDiyPage(id: number) {
  return http.delete<{ ok: boolean }>(`/diy/pages/${id}`);
}
