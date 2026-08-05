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

export type MerchantSettlementStatus =
  | 'approved'
  | 'bill_frozen'
  | 'bill_pending'
  | 'paid'
  | 'rejected'
  | 'withdraw_applied';

/** GET 列表/详情响应字段以 settlement_id 为主键；applied_at / payout_reference / version 当前未在 view 投影中返回。 */
export interface MerchantSettlement {
  amount: number;
  application_no?: string | null;
  applied_at?: string | null;
  mer_id?: number;
  merchant_name?: string;
  paid_at?: string | null;
  payout_reference?: string | null;
  period_end: string;
  period_start: string;
  review_note?: string;
  settlement_id: number;
  status: MerchantSettlementStatus;
  store_id?: number;
  updated_at: string;
  version?: number;
}

export interface MerchantSettlementPage {
  limit: number;
  list: MerchantSettlement[];
  page: number;
  total: number;
}

export function listMerchantSettlementsApi(params: {
  limit: number;
  page: number;
  status?: MerchantSettlementStatus;
}) {
  return requestClient.get<MerchantSettlementPage>('/settlements', { params });
}

export function getMerchantSettlementApi(id: number) {
  return requestClient.get<MerchantSettlement>(`/settlements/${id}`);
}

export function applyMerchantSettlementApi(id: number, body: { idempotency_key: string }) {
  return requestClient.post<MerchantSettlement>(`/settlements/${id}/apply`, body);
}
