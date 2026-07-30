import { requestClient } from '#/api/request';

export interface MerchantStaff {
  account: string;
  create_time: string;
  is_goods: number;
  is_open: number;
  is_verify: number;
  nickname: string;
  phone: string;
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

export function listMerchantStaffApi(params: { limit?: number; page?: number }) {
  return requestClient.get<PageResult<MerchantStaff>>('/setting/staff', { params });
}

export function createMerchantStaffApi(data: MerchantStaffSaveInput) {
  return requestClient.post<MerchantStaff>('/setting/staff', data);
}

export function updateMerchantStaffApi(id: number, data: MerchantStaffSaveInput) {
  return requestClient.put<MerchantStaff>(`/setting/staff/${id}`, data);
}
