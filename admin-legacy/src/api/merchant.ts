import { http } from './http';

export interface Merchant {
  mer_id: number;
  category_id: number;
  mer_name: string;
  real_name: string;
  mer_phone: string;
  mer_address: string;
  mark: string;
  status: number;
  mer_state: number;
  is_audit: number;
  create_time: string;
}

export interface Intention {
  mer_intention_id: number;
  phone: string;
  mer_name: string;
  name: string;
  create_time: string;
  status: number;
  fail_msg: string;
  mark: string;
  mer_id: number;
  merchant_category_id: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchMerchants(params: Record<string, unknown>) {
  return http.get<PageResult<Merchant>>('/merchants', { params });
}

export function setMerchantStatus(id: number, enabled: boolean) {
  return http.put(`/merchants/${id}/status`, { enabled });
}

export function fetchIntentions(params: Record<string, unknown>) {
  return http.get<PageResult<Intention>>('/merchant-intentions', { params });
}

export function auditIntention(
  id: number,
  body: { status: number; fail_msg?: string; mark?: string; account?: string; password?: string },
) {
  return http.post<{ mer_id?: number; account?: string; intention: Intention }>(
    `/merchant-intentions/${id}/audit`,
    body,
  );
}

export function fetchMerchantCategories() {
  return http.get<{ list: { merchant_category_id: number; category_name: string; commission_rate: number }[] }>(
    '/merchant-categories',
  );
}
