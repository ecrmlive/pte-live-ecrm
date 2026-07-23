import { http } from '@/api/http';

export interface ProductGroup {
  product_group_id: number;
  product_id: number;
  mer_id: number;
  price: number;
  buying_count_num: number;
  is_show: number;
  status: number;
  store_name?: string;
  mer_name?: string;
  ot_price?: number;
  start_time?: string;
  end_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchGroups(params: Record<string, unknown>) {
  return http.get<PageResult<ProductGroup>>('/combination/groups', { params });
}

export function updateGroup(id: number, data: Record<string, unknown>) {
  return http.put<ProductGroup>(`/combination/groups/${id}`, data);
}

export function deleteGroup(id: number) {
  return http.delete<{ ok: boolean }>(`/combination/groups/${id}`);
}
