import { requestClient } from '#/api/request';

export interface RechargePlan {
  id: number;
  name: string;
  amount: number;
  bonus_amount: number;
  status: number;
  sort: number;
  version: number;
  created_at?: string;
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

export type RechargePlanSaveInput = {
  name?: string;
  amount: number;
  bonus_amount: number;
  status: number;
  sort: number;
  version?: number;
};

export const listRechargePlans = (params?: {
  page?: number;
  limit?: number;
  status?: number;
}) =>
  requestClient.get<{ list: RechargePlan[]; total: number }>(
    '/recharge/plans',
    { params },
  );

export const createRechargePlan = (data: RechargePlanSaveInput) =>
  requestClient.post<RechargePlan>('/recharge/plans', data);

export const updateRechargePlan = (id: number, data: RechargePlanSaveInput) =>
  requestClient.put(`/recharge/plans/${id}`, data);

export const updateRechargePlanStatus = (id: number, status: number) =>
  requestClient.put(`/recharge/plans/${id}/status`, { status });

export const deleteRechargePlan = (id: number) =>
  requestClient.delete(`/recharge/plans/${id}`);

export const listRechargeOrders = (params: {
  page: number;
  limit: number;
  status?: string;
  date_from?: string;
  date_to?: string;
  keyword?: string;
}) =>
  requestClient.get<{ list: RechargeOrder[]; total: number }>(
    '/recharge/orders',
    { params },
  );
