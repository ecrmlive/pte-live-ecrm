import { requestClient } from '#/api/request';

export interface MerchantAssistActive {
  action_status: number;
  assist_count: number;
  assist_price: number;
  assist_user_count: number;
  create_time?: string;
  end_time: string;
  image?: string;
  is_show: number;
  ot_price?: number;
  pay_count?: number;
  product_assist_id: number;
  product_id: number;
  product_status: number;
  start_time: string;
  status: number;
  stock: number;
  store_info?: string;
  store_name: string;
}

export interface MerchantAssistPage {
  limit: number;
  list: MerchantAssistActive[];
  page: number;
  total: number;
}

export interface MerchantAssistSaveInput {
  assist_count: number;
  assist_price: number;
  assist_user_count: number;
  end_time: string;
  is_show?: number;
  product_id: number;
  start_time: string;
  status?: number;
  stock: number;
  store_info?: string;
  store_name?: string;
}

export function listMerchantAssistActivesApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantAssistPage>('/assist/actives', { params });
}

export function createMerchantAssistActiveApi(body: MerchantAssistSaveInput) {
  return requestClient.post<MerchantAssistActive>('/assist/actives', body);
}

export function updateMerchantAssistActiveApi(id: number, body: MerchantAssistSaveInput) {
  return requestClient.put<MerchantAssistActive>(`/assist/actives/${id}`, body);
}

export function setMerchantAssistShowApi(id: number, isShow: number) {
  return requestClient.put<MerchantAssistActive>(`/assist/actives/${id}/show`, { is_show: isShow });
}

export function deleteMerchantAssistActiveApi(id: number) {
  return requestClient.delete(`/assist/actives/${id}`);
}
