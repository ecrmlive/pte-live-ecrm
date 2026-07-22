import { http } from './http';

export interface Balance {
  mer_id: number;
  mer_money: number;
}

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

export function fetchBalance() {
  return http.get<Balance>('/finance/balance');
}

export function fetchWithdraws(params: Record<string, unknown>) {
  return http.get<PageResult<Financial>>('/finance/withdraws', { params });
}

export function createWithdraw(body: {
  extract_money: number;
  financial_type: number;
  financial_account: string;
  mark?: string;
}) {
  return http.post<Financial>('/finance/withdraw', body);
}

export function fetchWithdraw(id: number) {
  return http.get<Financial>(`/finance/withdraws/${id}`);
}
