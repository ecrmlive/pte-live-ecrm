import { http } from '@/api/http';

export interface Financial {
  financial_id: number;
  financial_sn: string;
  mer_money: number;
  extract_money: number;
  financial_type: number;
  financial_account: string;
  financial_status: number;
  status: number;
  refusal?: string;
  mer_id: number;
  mark?: string;
  create_time?: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchWithdraws(params: Record<string, unknown>) {
  return http.get<PageResult<Financial>>('/finance/withdraws', { params });
}

export function fetchWithdraw(id: number) {
  return http.get<Financial>(`/finance/withdraws/${id}`);
}

export function approveWithdraw(id: number) {
  return http.post(`/finance/withdraws/${id}/approve`, {});
}

export function rejectWithdraw(id: number, refusal: string) {
  return http.post(`/finance/withdraws/${id}/reject`, { refusal });
}
