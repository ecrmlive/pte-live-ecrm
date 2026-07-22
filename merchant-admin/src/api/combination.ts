import { http } from '@/api/http';

export interface ProductGroup {
  product_group_id: number;
  product_id: number;
  price: number;
  buying_count_num: number;
  time: number;
  is_show: number;
  store_name?: string;
  ot_price?: number;
}

export function fetchGroups(params: Record<string, unknown>) {
  return http.get<{ list: ProductGroup[]; total: number }>('/combination/groups', { params });
}

export function createGroup(data: Record<string, unknown>) {
  return http.post<ProductGroup>('/combination/groups', data);
}

export function updateGroup(id: number, data: Record<string, unknown>) {
  return http.put<ProductGroup>(`/combination/groups/${id}`, data);
}

export function setGroupShow(id: number, is_show: number) {
  return http.put<ProductGroup>(`/combination/groups/${id}/show`, { is_show });
}

export function deleteGroup(id: number) {
  return http.delete<{ ok: boolean }>(`/combination/groups/${id}`);
}
