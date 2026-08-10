import { requestClient } from '#/api/request';

export interface SvipPlan {
  id: number;
  name: string;
  cost_price?: number;
  price: number;
  plan_type: 'lifetime' | 'period' | 'trial';
  duration_days?: number | null;
  benefits: string;
  status: number;
  sort: number;
}

export interface SvipPlanInput {
  name: string;
  cost_price?: number;
  price: number;
  plan_type: SvipPlan['plan_type'];
  duration_days: number;
  benefits: string[];
  status: number;
  sort: number;
}

export type SvipPayType =
  | 'weixin'
  | 'wechat'
  | 'alipay'
  | 'routine'
  | 'sys'
  | 'free'
  | '';

export interface SvipOrder {
  id: number;
  order_no: string;
  user_id: number;
  nickname: string;
  phone: string;
  plan_id: number;
  plan_name: string;
  plan_type: string;
  duration_days?: number | null;
  amount: number;
  pay_type: SvipPayType | string;
  status: 'closed' | 'paid' | 'pending';
  created_at: string;
  paid_at?: string | null;
  end_time?: string | null;
}

export interface SvipOrderSummary {
  paid_member_count: number;
  paid_amount: number;
  expired_member_count: number;
}

export const listSvipPlans = (params?: { page?: number; limit?: number }) =>
  requestClient.get<{ list: SvipPlan[]; total: number }>('/svip/plans', {
    params,
  });
export const createSvipPlan = (data: SvipPlanInput) =>
  requestClient.post<SvipPlan>('/svip/plans', data);
export const updateSvipPlan = (id: number, data: SvipPlanInput) =>
  requestClient.put(`/svip/plans/${id}`, data);
export const updateSvipPlanStatus = (id: number, status: number) =>
  requestClient.put(`/svip/plans/${id}/status`, { status });
export const deleteSvipPlan = (id: number) =>
  requestClient.delete(`/svip/plans/${id}`);
export const listSvipOrders = (params: {
  page: number;
  limit: number;
  pay_type?: string;
  title?: string;
  nickname?: string;
  status?: string;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) =>
  requestClient.get<{ list: SvipOrder[]; total: number }>('/svip/orders', {
    params,
  });
export const getSvipOrderSummary = (params?: {
  pay_type?: string;
  title?: string;
  nickname?: string;
  status?: string;
  keyword?: string;
  date_from?: string;
  date_to?: string;
}) =>
  requestClient.get<SvipOrderSummary>('/svip/orders/summary', { params });
