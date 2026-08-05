import { requestClient } from '#/api/request';

export interface MerchantCombinationGroup {
  action_status: number;
  buying_count_num: number;
  buying_num?: number;
  create_time: string;
  end_time: string;
  image?: string;
  is_show: number;
  mer_id?: number;
  once_pay_count?: number;
  ot_price?: number;
  pay_count?: number;
  price: number;
  product_group_id: number;
  product_id: number;
  product_status: number;
  start_time: string;
  status: number;
  store_name?: string;
  success_num?: number;
  time: number;
}

export interface MerchantCombinationPage {
  limit: number;
  list: MerchantCombinationGroup[];
  page: number;
  total: number;
}

export interface MerchantCombinationSaveInput {
  buying_count_num: number;
  end_time: string;
  is_show?: number;
  price: number;
  product_id: number;
  start_time: string;
  status?: number;
  time: number;
}

export function listMerchantCombinationGroupsApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantCombinationPage>('/combination/groups', { params });
}

export function createMerchantCombinationGroupApi(body: MerchantCombinationSaveInput) {
  return requestClient.post<MerchantCombinationGroup>('/combination/groups', body);
}

export function updateMerchantCombinationGroupApi(id: number, body: MerchantCombinationSaveInput) {
  return requestClient.put<MerchantCombinationGroup>(`/combination/groups/${id}`, body);
}

export function setMerchantCombinationShowApi(id: number, isShow: number) {
  return requestClient.put<MerchantCombinationGroup>(`/combination/groups/${id}/show`, { is_show: isShow });
}

export function deleteMerchantCombinationGroupApi(id: number) {
  return requestClient.delete(`/combination/groups/${id}`);
}
