import { requestClient } from '#/api/request';

export interface MerchantOrderStats {
  completed: number;
  paid_order: number;
  pending_refund: number;
  pending_ship: number;
  shipped: number;
}

export interface MerchantProductStats {
  draft: number;
  off_sale: number;
  on_sale: number;
  pending_review: number;
  total: number;
}

export function getMerchantOrderStatsApi() {
  return requestClient.get<MerchantOrderStats>('/dashboard/order-stats');
}

export function getMerchantProductStatsApi() {
  return requestClient.get<MerchantProductStats>('/dashboard/product-stats');
}
