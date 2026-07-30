import { requestClient } from '#/api/request';

export interface MerchantBalance {
  mer_id: number;
  mer_money: number;
}

export interface MerchantWithdraw {
  create_time?: string;
  extract_money: number;
  financial_account: string;
  financial_id: number;
  financial_sn: string;
  financial_status: number;
  financial_type: number;
  image?: string;
  mark?: string;
  mer_money: number;
  refusal?: string;
  status: number;
  status_time?: string;
}

export interface MerchantWithdrawPage {
  limit: number;
  list: MerchantWithdraw[];
  page: number;
  total: number;
}

export function getMerchantBalanceApi() {
  return requestClient.get<MerchantBalance>('/finance/balance');
}

export function listMerchantWithdrawsApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantWithdrawPage>('/finance/withdraws', { params });
}

export function getMerchantWithdrawApi(id: number) {
  return requestClient.get<MerchantWithdraw>(`/finance/withdraws/${id}`);
}

export function applyMerchantWithdrawApi(body: {
  extract_money: number;
  financial_account: string;
  financial_type: number;
  mark: string;
}) {
  return requestClient.post<MerchantWithdraw>('/finance/withdraw', body);
}
