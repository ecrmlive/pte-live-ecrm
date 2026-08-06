import { requestClient } from '#/api/request';

export interface RechargePlan {
  id: number;
  name: string;
  amount: number;
  bonus_amount: number;
  status: number;
  sort: number;
  version: number;
}

export interface RechargeOrder {
  id: number;
  recharge_no: string;
  user_id: number;
  amount: number;
  bonus_amount: number;
  status: string;
  created_at: string;
  paid_at?: string;
}

export const listRechargePlans = (params?: {
  keyword?: string;
  status?: number;
}) => requestClient.get<{ list: RechargePlan[] }>('/recharge/plans', { params });

export const createRechargePlan = (data: Omit<RechargePlan, 'id'>) =>
  requestClient.post<RechargePlan>('/recharge/plans', data);

export const updateRechargePlan = (id: number, data: Omit<RechargePlan, 'id'>) =>
  requestClient.put(`/recharge/plans/${id}`, data);

export const listRechargeOrders = (params: {
  page: number;
  limit: number;
  status?: string;
  date_from?: string;
  date_to?: string;
  keyword?: string;
}) => requestClient.get<{ list: RechargeOrder[]; total: number }>('/recharge/orders', { params });
