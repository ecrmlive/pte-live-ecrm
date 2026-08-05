import { requestClient } from '#/api/request';

export type MerchantDiscountStatus = 'active' | 'closed' | 'draft' | 'pending' | 'rejected';

export interface MerchantDiscount {
  activity_id: number;
  ends_at: string;
  free_shipping: boolean;
  name: string;
  package_price: number;
  product_ids: number[];
  remark: string;
  starts_at: string;
  status: MerchantDiscountStatus;
  store_id: number;
}

export interface MerchantDiscountPage {
  limit: number;
  list: MerchantDiscount[];
  page: number;
  total: number;
}

export interface MerchantDiscountInput {
  ends_at?: string;
  free_shipping?: boolean;
  name: string;
  package_price: number;
  product_ids: number[];
  remark?: string;
  starts_at?: string;
  status?: MerchantDiscountStatus;
}

export function listMerchantDiscountsApi(params: {
  keyword?: string;
  limit: number;
  page: number;
  status?: MerchantDiscountStatus;
}) {
  return requestClient.get<MerchantDiscountPage>('/marketing/discounts', { params });
}

export function getMerchantDiscountApi(id: number) {
  return requestClient.get<MerchantDiscount>(`/marketing/discounts/${id}`);
}

export function createMerchantDiscountApi(body: MerchantDiscountInput) {
  return requestClient.post<MerchantDiscount>('/marketing/discounts', body);
}

export function updateMerchantDiscountApi(id: number, body: MerchantDiscountInput) {
  return requestClient.put<MerchantDiscount>(`/marketing/discounts/${id}`, body);
}

export function setMerchantDiscountStatusApi(id: number, status: MerchantDiscountStatus) {
  return requestClient.put<MerchantDiscount>(`/marketing/discounts/${id}/status`, { status });
}

export function deleteMerchantDiscountApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/marketing/discounts/${id}`);
}
