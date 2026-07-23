export interface RegionRow {
  id: number;
  shortname: string;
  name: string;
  merger_name: string;
  level: number;
  lng: string;
  lat: string;
  ad_code?: string;
}

export interface RegionAreaItem {
  id: number;
  name: string;
  city?: Record<number, RegionAreaItem> | RegionAreaItem[];
}

export interface RegionSearchForm {
  name: string;
  level: number;
  province_id: number;
  city_id: number | string;
}

export interface RegionFormModel {
  id?: number;
  level: number;
  province_id?: number;
  city_id?: number;
  name: string;
  shortname: string;
  merger_name: string;
  pinyin: string;
  code?: string;
  ad_code: string;
  zip_code: string;
  first: string;
  lng: string;
  lat: string;
  sort: string | number;
}

/** PHP/Go 地区树中 city 为 id 映射对象，兼容旧数组结构 */
export function regionCityOptions(
  areaList: Record<number | string, RegionAreaItem>,
  provinceId: number | string,
): RegionAreaItem[] {
  const pid = Number(provinceId);
  if (!pid) return [];
  const province = areaList[pid] ?? areaList[String(pid)];
  const city = province?.city;
  if (!city) return [];
  return Array.isArray(city) ? city : Object.values(city);
}
