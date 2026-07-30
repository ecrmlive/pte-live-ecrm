import { http } from "@/utils/request";

export interface DiyBanner {
  id: number;
  title: string;
  image?: string;
  url?: string;
}

export interface DiyMenu {
  id: number;
  name: string;
  icon?: string;
  url?: string;
}

export interface DiyItem {
  type: string;
  name?: string;
  params?: Record<string, unknown>;
  style?: Record<string, unknown>;
  data?: Array<Record<string, unknown>>;
  images?: Array<Record<string, unknown>>;
  [key: string]: unknown;
}

export interface DiyHome {
  id: number;
  name?: string;
  title?: string;
  page?: Record<string, unknown>;
  items?: DiyItem[];
  banners: DiyBanner[];
  menus: DiyMenu[];
}

export function fetchDiyHome(merId?: number) {
  const q = merId ? `?mer_id=${merId}` : "";
  return http.get<DiyHome>(`/diy/home${q}`, false);
}

export function fetchDiyPage(id: number) {
  return http.get<DiyHome>(`/diy/pages/${id}`, false);
}
