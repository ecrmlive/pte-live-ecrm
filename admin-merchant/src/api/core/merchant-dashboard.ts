import { requestClient } from '#/api/request';

export interface MerchantDashboardSummary {
  available_balance: number;
  paid_order: number;
  pending_refund: number;
  product_total: number;
}

export function getMerchantDashboardSummaryApi() {
  return requestClient.get<MerchantDashboardSummary>('/dashboard/summary');
}
