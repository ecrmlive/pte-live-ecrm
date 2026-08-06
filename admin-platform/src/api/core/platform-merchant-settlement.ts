import { requestClient } from '#/api/request';

export type MerchantSettlementStatus =
  | 'approved'
  | 'bill_frozen'
  | 'bill_pending'
  | 'cancelled'
  | 'paid'
  | 'rejected'
  | 'withdraw_applied';

export interface MerchantSettlementRow {
  settlement_id: number;
  merchant_id: number;
  store_id: number;
  merchant_name: string;
  period_start: string;
  period_end: string;
  amount: number;
  status: MerchantSettlementStatus;
  updated_at: string;
}

export interface MerchantSettlementSummary {
  status: MerchantSettlementStatus;
  amount: number;
  count: number;
}

export interface MerchantSettlementPage {
  list: MerchantSettlementRow[];
  total: number;
  page: number;
  limit: number;
}

export interface MerchantSettlementCommandInput {
  idempotency_key: string;
  payout_reference?: string;
  review_note?: string;
}

export function listPlatformMerchantSettlementsApi(params: {
  limit: number;
  merchant_id?: number;
  page: number;
  status?: MerchantSettlementStatus;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<MerchantSettlementPage>('/finance/merchant-settlements', { params });
}

export function getPlatformMerchantSettlementSummaryApi() {
  return requestClient.get<{ list: MerchantSettlementSummary[] }>('/finance/merchant-settlements/summary');
}

export function approvePlatformMerchantSettlementApi(settlementId: number, data: MerchantSettlementCommandInput) {
  return requestClient.post(`/finance/merchant-settlements/${settlementId}/approve`, data);
}

export function rejectPlatformMerchantSettlementApi(settlementId: number, data: MerchantSettlementCommandInput) {
  return requestClient.post(`/finance/merchant-settlements/${settlementId}/reject`, data);
}

export function markPlatformMerchantSettlementPaidApi(settlementId: number, data: MerchantSettlementCommandInput) {
  return requestClient.post(`/finance/merchant-settlements/${settlementId}/mark-paid`, data);
}

export type TransferRecordStatus = Extract<MerchantSettlementStatus, 'approved' | 'paid' | 'rejected'>;

/** 转账记录：结算打款链路只读投影（approved / paid / rejected）。 */
export function listPlatformTransferRecordsApi(params: {
  limit: number;
  merchant_id?: number;
  page: number;
  status?: TransferRecordStatus;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<MerchantSettlementPage>('/finance/transfer-records', { params });
}
