import { requestClient } from '#/api/request';

export interface PlatformWithdraw {
  create_time?: string;
  extract_money: number;
  financial_account: string;
  financial_id: number;
  financial_sn: string;
  financial_status: number;
  financial_type: number;
  image?: string;
  mark?: string;
  mer_id: number;
  user_id: number;
  mer_money: number;
	paid_at?: string;
	payout_reference?: string;
  refusal?: string;
  status: number;
  status_time?: string;
  withdrawal_status?: 'applied' | 'approved' | 'paid' | 'paying' | 'rejected' | 'reviewing';
}

export interface PlatformWithdrawPage {
  limit: number;
  list: PlatformWithdraw[];
  page: number;
  total: number;
}

export function listPlatformWithdrawsApi(params: {
  limit: number;
  page: number;
  status?: number;
}) {
  return requestClient.get<PlatformWithdrawPage>('/finance/withdraws', { params });
}

export function getPlatformWithdrawApi(id: number) {
  return requestClient.get<PlatformWithdraw>(`/finance/withdraws/${id}`);
}

export function approvePlatformWithdrawApi(id: number) {
  return requestClient.post(`/finance/withdraws/${id}/approve`);
}

export function rejectPlatformWithdrawApi(id: number, refusal: string) {
  return requestClient.post(`/finance/withdraws/${id}/reject`, { refusal });
}

export function markPlatformWithdrawPaidApi(id: number, data: { idempotency_key: string; payout_reference: string }) {
  return requestClient.post(`/finance/withdraws/${id}/mark-paid`, data);
}
