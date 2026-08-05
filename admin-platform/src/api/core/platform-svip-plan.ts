import { requestClient } from '#/api/request';

export interface SvipPlan {
  id: number;
  name: string;
  price: number;
  plan_type: 'lifetime' | 'period' | 'trial';
  duration_days?: number | null;
  benefits: string;
  status: number;
  sort: number;
}

export interface SvipPlanInput {
  name: string;
  price: number;
  plan_type: SvipPlan['plan_type'];
  duration_days: number;
  benefits: string[];
  status: number;
  sort: number;
}

export interface SvipOrder {
  id: number;
  order_no: string;
  user_id: number;
  plan_id: number;
  plan_name: string;
  plan_type: string;
  duration_days?: number | null;
  amount: number;
  status: 'closed' | 'paid' | 'pending';
  created_at: string;
  paid_at?: string | null;
}

export interface SvipOrderSummary {
  total: number;
  pending: number;
  paid: number;
  closed: number;
  paid_amount: number;
}

export const listSvipPlans = () => requestClient.get<{ list: SvipPlan[] }>('/svip/plans');
export const createSvipPlan = (data: SvipPlanInput) => requestClient.post<SvipPlan>('/svip/plans', data);
export const updateSvipPlan = (id: number, data: SvipPlanInput) => requestClient.put(`/svip/plans/${id}`, data);
export const listSvipOrders = (params: { page: number; limit: number; status?: string }) => requestClient.get<{ list: SvipOrder[]; total: number }>('/svip/orders', { params });
export const getSvipOrderSummary = () => requestClient.get<SvipOrderSummary>('/svip/orders/summary');
