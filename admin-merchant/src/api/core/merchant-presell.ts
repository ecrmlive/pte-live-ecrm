import { requestClient } from '#/api/request';

export interface MerchantPresellActive {
  down_price: number;
  end_time: string;
  final_end_time: string;
  final_price: number;
  final_start_time: string;
  image?: string;
  is_show: number;
  price: number;
  presell_type: number;
  product_id: number;
  product_presell_id: number;
  seles: number;
  start_time: string;
  stock: number;
  store_info: string;
  store_name: string;
}

export interface MerchantPresellPage {
  limit: number;
  list: MerchantPresellActive[];
  page: number;
  total: number;
}

export interface MerchantPresellSaveInput {
  down_price: number;
  end_time: string;
  final_end_time: string;
  final_price: number;
  final_start_time: string;
  is_show?: number;
  presell_type: number;
  price: number;
  product_id: number;
  stock: number;
  store_info: string;
  store_name: string;
  start_time: string;
}

export interface MerchantPresellListParams {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  is_show?: 0 | 1;
  limit: number;
  page: number;
}

export function listMerchantPresellActivesApi(params: MerchantPresellListParams) {
  return requestClient.get<MerchantPresellPage>('/presell/actives', { params });
}

export function createMerchantPresellActiveApi(body: MerchantPresellSaveInput) {
  return requestClient.post<MerchantPresellActive>('/presell/actives', body);
}

export function updateMerchantPresellActiveApi(id: number, body: MerchantPresellSaveInput) {
  return requestClient.put<MerchantPresellActive>(`/presell/actives/${id}`, body);
}

export function setMerchantPresellShowApi(id: number, isShow: number) {
  return requestClient.put<MerchantPresellActive>(`/presell/actives/${id}/show`, { is_show: isShow });
}

export function deleteMerchantPresellActiveApi(id: number) {
  return requestClient.delete(`/presell/actives/${id}`);
}
