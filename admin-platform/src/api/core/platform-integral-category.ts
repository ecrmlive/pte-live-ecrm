import { requestClient } from '#/api/request';

export interface PlatformIntegralCategoryRow {
  cate_name: string;
  create_time: string;
  has_product?: number;
  is_show: number;
  pic?: string;
  pid?: number;
  sort: number;
  store_category_id: number;
  type?: number;
}

export interface PlatformIntegralCategorySaveInput {
  cate_name: string;
  is_show: number;
  pic?: string;
  pid?: number;
  sort: number;
}

export function listPlatformIntegralCategoriesApi() {
  return requestClient.get<{ list: PlatformIntegralCategoryRow[]; total: number }>(
    '/points/categories',
  );
}

export function selectPlatformIntegralCategoriesApi() {
  return requestClient.get<{ list: PlatformIntegralCategoryRow[] }>(
    '/points/categories/select',
  );
}

export function createPlatformIntegralCategoryApi(
  payload: PlatformIntegralCategorySaveInput,
) {
  return requestClient.post<{ ok: boolean }>('/points/categories', payload);
}

export function updatePlatformIntegralCategoryApi(
  id: number,
  payload: PlatformIntegralCategorySaveInput,
) {
  return requestClient.put<{ ok: boolean }>(`/points/categories/${id}`, payload);
}

export function updatePlatformIntegralCategoryStatusApi(
  id: number,
  isShow: number,
) {
  return requestClient.put<{ ok: boolean; is_show: number }>(
    `/points/categories/${id}/status`,
    { is_show: isShow },
  );
}

export function deletePlatformIntegralCategoryApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/points/categories/${id}`);
}
