import { requestClient } from '#/api/request';

export interface PlatformUserRechargeRow {
  avatar?: string;
  can_refund?: boolean;
  create_time?: string;
  give_price: number;
  nickname?: string;
  order_id: string;
  paid: number;
  paid_name?: string;
  pay_time?: string;
  price: number;
  real_name?: string;
  recharge_id: number;
  recharge_type: string;
  recharge_type_name?: string;
  refund_price?: number;
  uid: number;
}

export interface PlatformUserRechargePage {
  limit: number;
  list: PlatformUserRechargeRow[];
  page: number;
  total: number;
}

export interface PlatformUserRechargeTotal {
  total_pay_price: number;
  total_refund_price: number;
  total_routine_price: number;
  total_wx_price: number;
}

export type PlatformUserRechargeQuery = {
  date_from?: string;
  date_to?: string;
  limit?: number;
  order_id?: string;
  page?: number;
  paid?: number;
  recharge_type?: string;
  user_keyword?: string;
  user_type?: string;
};

export function listPlatformUserRechargesApi(params: PlatformUserRechargeQuery) {
  return requestClient.get<PlatformUserRechargePage>('/finance/user-recharges', {
    params,
  });
}

export function getPlatformUserRechargeTotalApi() {
  return requestClient.get<PlatformUserRechargeTotal>(
    '/finance/user-recharges/total',
  );
}

export function refundPlatformUserRechargeApi(
  id: number,
  data: { amount?: number; idempotency_key: string },
) {
  return requestClient.post<{ idempotent?: boolean; ok: boolean; recharge_id: number }>(
    `/finance/user-recharges/${id}/refund`,
    data,
  );
}
