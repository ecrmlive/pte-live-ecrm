import { requestClient } from '#/api/request';

export interface MerchantStaff {
  account: string;
  create_time: string;
  is_goods: number;
  is_open: number;
  is_verify: number;
  nickname: string;
  phone: string;
	role_code: 'clerk' | 'delivery' | 'manager' | 'service';
  service_id: number;
  status: number;
}

export interface MerchantStaffSaveInput {
  account?: string;
  is_goods?: number;
  is_open?: number;
  is_verify?: number;
  nickname?: string;
  password?: string;
  phone?: string;
  status?: number;
}

interface PageResult<T> { limit: number; list: T[]; page: number; total: number }

export interface MerchantStaffListParams {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  limit?: number;
  page?: number;
  staff_scope?: 'delivery' | 'service';
  status?: 0 | 1;
}

export function listMerchantStaffApi(params: MerchantStaffListParams = {}) {
  return requestClient.get<PageResult<MerchantStaff>>('/setting/staff', { params });
}

export function createMerchantStaffApi(data: MerchantStaffSaveInput) {
  return requestClient.post<MerchantStaff>('/setting/staff', data);
}

export function updateMerchantStaffApi(id: number, data: MerchantStaffSaveInput) {
  return requestClient.put<MerchantStaff>(`/setting/staff/${id}`, data);
}

export function removeMerchantStaffApi(id: number) {
  return requestClient.delete<{ removed: boolean; service_id: number }>(`/setting/staff/${id}`);
}
