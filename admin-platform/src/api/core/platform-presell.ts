import { requestClient } from '#/api/request';

export interface PlatformPresell {
  action_status?: number;
  attend_num?: number;
  delivery_day?: number;
  delivery_type?: number;
  down_price?: number;
  end_time: string;
  final_end_time?: string;
  final_price?: number;
  final_start_time?: string;
  image?: string;
  is_show?: number;
  is_trader?: number;
  mer_id: number;
  mer_name?: string;
  ot_price?: number;
  pay_count?: number;
  presell_status?: number;
  presell_status_text?: string;
  presell_type?: number;
  price: number;
  product_id?: number;
  product_presell_id: number;
  product_status?: number;
  product_status_name?: string;
  refusal?: string;
  seles?: number;
  star?: number;
  start_time: string;
  status: number;
  stock?: number;
  stock_count?: number;
  store_info?: string;
  store_name: string;
  success_num?: number;
  sys_labels?: string;
  trader_name?: string;
  us_status?: number;
}

export interface PlatformPresellPage {
  limit: number;
  list: PlatformPresell[];
  page: number;
  total: number;
}

export interface PlatformPresellListParams {
  page: number;
  limit: number;
  presell_type?: number;
  mer_id?: number;
  keyword?: string;
  is_trader?: number;
  star?: number;
  product_status?: number;
  type?: number;
  us_status?: number;
  sys_labels?: string;
}

export interface PlatformPresellTypeFilter {
  type: number;
  name: string;
  count: number;
}

export interface PlatformPresellInput {
  store_name?: string;
  store_info?: string;
  price?: number;
  down_price?: number;
  final_price?: number;
  stock?: number;
  stock_count?: number;
  pay_count?: number;
  delivery_type?: number;
  delivery_day?: number;
  start_time?: string;
  end_time?: string;
  final_start_time?: string;
  final_end_time?: string;
  is_show?: number;
  status?: number;
  product_status?: number;
  refusal?: string;
  star?: number;
  sys_labels?: string;
  presell_type?: number;
}

export function listPlatformPresellsApi(params: PlatformPresellListParams) {
  return requestClient.get<PlatformPresellPage>('/presell/actives', { params });
}

export function getPlatformPresellTypeFilterApi(
  params: Omit<PlatformPresellListParams, 'page' | 'limit' | 'presell_type'>,
) {
  return requestClient.get<{ list: PlatformPresellTypeFilter[] }>(
    '/presell/actives/filter',
    { params },
  );
}

export function getPlatformPresellApi(id: number) {
  return requestClient.get<PlatformPresell>(`/presell/actives/${id}`);
}

export function updatePlatformPresellApi(
  id: number,
  payload: PlatformPresellInput,
) {
  return requestClient.put<PlatformPresell>(`/presell/actives/${id}`, payload);
}

export function setPlatformPresellShowApi(id: number, isShow: number) {
  return requestClient.put(`/presell/actives/${id}/show`, { is_show: isShow });
}

export function setPlatformPresellStarApi(id: number, star: number) {
  return requestClient.put(`/presell/actives/${id}/star`, { star });
}

export function setPlatformPresellLabelsApi(id: number, sysLabels: string) {
  return requestClient.put(`/presell/actives/${id}/labels`, {
    sys_labels: sysLabels,
  });
}

export function auditPlatformPresellApi(
  id: number,
  status: number,
  refusal?: string,
) {
  return requestClient.post<PlatformPresell>('/presell/actives/audit', {
    id,
    status,
    refusal: refusal || '',
  });
}

export function forceOffPlatformPresellApi(ids: number[], reason: string) {
  return requestClient.post<{ ok: true }>('/presell/actives/force-off', {
    ids,
    reason,
  });
}

export function deletePlatformPresellApi(id: number) {
  return requestClient.delete<{ ok: true }>(`/presell/actives/${id}`);
}
