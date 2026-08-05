import { requestClient } from '#/api/request';

export type CommissionStatus = 'available' | 'pending' | 'settled' | 'voided';

export interface DistributionPromoter {
  available_commission: number;
  direct_user_count: number;
  pending_commission: number;
  settled_commission: number;
  status: number;
  updated_at: string;
  user_id: number;
}

export interface DistributionCommission {
  amount: number;
  available_at?: string;
  commission_id: number;
  created_at: string;
  order_id: number;
  status: CommissionStatus;
  user_id: number;
}

export interface DistributionPage<T> { limit: number; list: T[]; page: number; total: number; }

export interface DistributionSummary {
  active_promoter_count: number;
  available_commission: number;
  pending_commission: number;
  promoter_count: number;
  settled_commission: number;
}

export function getDistributionSummaryApi() {
  return requestClient.get<DistributionSummary>('/distribution/summary');
}

export function listDistributionPromotersApi(params: { limit: number; page: number; status?: 0 | 1; user_id?: number }) {
  return requestClient.get<DistributionPage<DistributionPromoter>>('/distribution/promoters', { params });
}

export function listDistributionCommissionsApi(params: { limit: number; page: number; status?: CommissionStatus; user_id?: number }) {
  return requestClient.get<DistributionPage<DistributionCommission>>('/distribution/commissions', { params });
}
