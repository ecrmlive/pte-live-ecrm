import { requestClient } from '#/api/request';

export interface PlatformCombination {
  buying_count_num: number;
  create_time: string;
  end_time: string;
  mer_id: number;
  mer_name?: string;
  price: number;
  product_group_id: number;
  product_id: number;
  start_time: string;
  status: number;
  store_name?: string;
}

export interface PlatformCombinationPage {
  limit: number;
  list: PlatformCombination[];
  page: number;
  total: number;
}

export interface PlatformCombinationInput {
  price?: number;
  buying_count_num?: number;
  time?: number;
  start_time?: string;
  end_time?: string;
  is_show?: number;
  status?: number;
}

export function listPlatformCombinationsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformCombinationPage>('/combination/groups', { params });
}

export function getPlatformCombinationApi(id: number) {
  return requestClient.get<PlatformCombination & { time: number; is_show: number }>(
    '/combination/groups/' + id,
  );
}

export function updatePlatformCombinationApi(
  id: number,
  payload: PlatformCombinationInput,
) {
  return requestClient.put(`/combination/groups/${id}`, payload);
}

export function deletePlatformCombinationApi(id: number) {
  return requestClient.delete<{ ok: true }>('/combination/groups/' + id);
}
