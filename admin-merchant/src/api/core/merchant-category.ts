import { requestClient } from '#/api/request';

export interface MerchantCategory { category_id: number; parent_id: number; name: string; sort: number; status: 0 | 1; children?: MerchantCategory[]; }
export interface MerchantCategoryInput { parent_id: number; name: string; sort: number; status: 0 | 1; }
export interface MerchantCategoryListParams {
  keyword?: string;
  status?: 0 | 1;
}

export function listMerchantCategoriesApi(params: MerchantCategoryListParams = {}) {
  return requestClient.get<{ list: MerchantCategory[] }>('/store-categories', { params });
}
export function createMerchantCategoryApi(body: MerchantCategoryInput) { return requestClient.post<MerchantCategory>('/store-categories', body); }
export function updateMerchantCategoryApi(id: number, body: MerchantCategoryInput) { return requestClient.put<MerchantCategory>(`/store-categories/${id}`, body); }
export function deleteMerchantCategoryApi(id: number) { return requestClient.delete(`/store-categories/${id}`); }
