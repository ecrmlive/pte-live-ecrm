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

export interface DiyHome {
  id: number;
  name?: string;
  title?: string;
  banners: DiyBanner[];
  menus: DiyMenu[];
}

export function fetchDiyHome() {
  return http.get<DiyHome>("/diy/home", false);
}
