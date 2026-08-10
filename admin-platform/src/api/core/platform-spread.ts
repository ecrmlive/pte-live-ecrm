import { requestClient } from '#/api/request';

export type CommissionStatus = 'available' | 'pending' | 'settled' | 'voided';

export interface DistributionPromoter {
  available_commission: number;
  avatar_url: string;
  commission_amount: number;
  direct_user_count: number;
  level_id: number;
  level_name: string;
  mobile: string;
  nickname: string;
  parent_nickname: string;
  parent_user_id: number;
  pending_commission: number;
  settled_commission: number;
  spread_order_amount: number;
  spread_order_count: number;
  spread_user_count: number;
  status: number;
  unwithdrawn_amount: number;
  updated_at: string;
  user_id: number;
  withdraw_count: number;
  withdrawn_amount: number;
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

export interface DistributionLevelTaskItem {
  info: string;
  name: string;
  num: number;
}

export interface DistributionLevelTaskRule {
  pay_money: DistributionLevelTaskItem;
  pay_num: DistributionLevelTaskItem;
  spread_money: DistributionLevelTaskItem;
  spread_pay_num: DistributionLevelTaskItem;
  spread_user: DistributionLevelTaskItem;
}

export interface DistributionLevel {
  extension_one: number;
  extension_two: number;
  icon_url: string;
  id: number;
  name: string;
  promoter_count: number;
  rank: number;
  status: number;
  task_rule: DistributionLevelTaskRule;
}

export type DistributionLevelSaveInput = {
  extension_one: number;
  extension_two: number;
  icon_url: string;
  name: string;
  rank: number;
  status?: number;
  task_rule: DistributionLevelTaskRule;
};

export interface DistributionChild {
  avatar_url: string;
  bound_at: string;
  is_promoter: number;
  level: number;
  mobile: string;
  nickname: string;
  pay_amount: number;
  pay_count: number;
  spread_user_count: number;
  user_id: number;
}

export interface DistributionSpreadOrder {
  buyer_id: number;
  buyer_name: string;
  commission: number;
  created_at: string;
  order_id: number;
  order_no: string;
  paid_at?: string;
  pay_amount: number;
  status: string;
}

export interface DistributionPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface DistributionSummary {
  active_promoter_count: number;
  available_commission: number;
  pending_commission: number;
  promoter_count: number;
  settled_commission: number;
  spread_order_amount: number;
  spread_order_count: number;
  spread_user_count: number;
  unwithdrawn_amount: number;
  withdrawn_amount: number;
}

export type PromoterListParams = {
  date_from?: string;
  date_to?: string;
  keyword?: string;
  keyword_type?: 'nickname' | 'phone' | 'uid';
  level_id?: number;
  limit: number;
  page: number;
  sort_field?: 'commission_amount' | 'unwithdrawn_amount' | 'withdrawn_amount';
  sort_order?: 'asc' | 'desc';
  status?: 0 | 1;
  user_id?: number;
};

export type PromoterSummaryParams = Omit<PromoterListParams, 'limit' | 'page' | 'sort_field' | 'sort_order'>;

export function getDistributionSummaryApi(params?: PromoterSummaryParams) {
  return requestClient.get<DistributionSummary>('/distribution/summary', { params });
}

export function listDistributionPromotersApi(params: PromoterListParams) {
  return requestClient.get<DistributionPage<DistributionPromoter>>('/distribution/promoters', {
    params,
  });
}

export function listDistributionCommissionsApi(params: {
  date_from?: string;
  date_to?: string;
  limit: number;
  page: number;
  status?: CommissionStatus;
  user_id?: number;
}) {
  return requestClient.get<DistributionPage<DistributionCommission>>(
    '/distribution/commissions',
    { params },
  );
}

export function listDistributionLevelsApi(params?: {
  keyword?: string;
  limit?: number;
  name?: string;
  page?: number;
}) {
  return requestClient.get<{
    limit?: number;
    list: DistributionLevel[];
    page?: number;
    total?: number;
  }>('/distribution/levels', { params });
}

export function getDistributionLevelApi(id: number) {
  return requestClient.get<DistributionLevel>(`/distribution/levels/${id}`);
}

export function createDistributionLevelApi(input: DistributionLevelSaveInput) {
  return requestClient.post<DistributionLevel>('/distribution/levels', input);
}

export function updateDistributionLevelConfigApi(
  id: number,
  input: DistributionLevelSaveInput,
) {
  return requestClient.put<DistributionLevel>(`/distribution/levels/${id}`, input);
}

export function deleteDistributionLevelApi(id: number) {
  return requestClient.delete(`/distribution/levels/${id}`);
}

export function listDistributionChildrenApi(
  userId: number,
  params: {
    date_from?: string;
    date_to?: string;
    keyword?: string;
    level?: 0 | 1 | 2;
    limit: number;
    page: number;
    sort_field?: 'spread_user_count';
    sort_order?: 'asc' | 'desc';
  },
) {
  return requestClient.get<DistributionPage<DistributionChild>>(
    `/distribution/promoters/${userId}/children`,
    { params },
  );
}

export function listDistributionSpreadOrdersApi(
  userId: number,
  params: { limit: number; page: number },
) {
  return requestClient.get<DistributionPage<DistributionSpreadOrder>>(
    `/distribution/promoters/${userId}/orders`,
    { params },
  );
}

export function clearDistributionParentApi(
  userId: number,
  input: { idempotency_key: string; reason: string },
) {
  return requestClient.post(`/distribution/promoters/${userId}/clear-parent`, input);
}

export function updateDistributionLevelApi(
  userId: number,
  input: { idempotency_key: string; level_id: number; reason: string },
) {
  return requestClient.post(`/distribution/promoters/${userId}/level`, input);
}

export interface WithdrawBank {
  created_at: string;
  id: number;
  name: string;
  sort: number;
  status: number;
  updated_at: string;
}

export type WithdrawBankSaveInput = {
  name: string;
  sort?: number;
  status?: number;
};

export function listWithdrawBanksApi(params?: {
  keyword?: string;
  limit?: number;
  page?: number;
  status?: 0 | 1;
}) {
  return requestClient.get<DistributionPage<WithdrawBank>>('/distribution/withdraw-banks', {
    params,
  });
}

export function createWithdrawBankApi(input: WithdrawBankSaveInput) {
  return requestClient.post<WithdrawBank>('/distribution/withdraw-banks', input);
}

export function updateWithdrawBankApi(id: number, input: WithdrawBankSaveInput) {
  return requestClient.put<WithdrawBank>(`/distribution/withdraw-banks/${id}`, input);
}

export function setWithdrawBankStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<WithdrawBank>(`/distribution/withdraw-banks/${id}/status`, {
    status,
  });
}

export function deleteWithdrawBankApi(id: number) {
  return requestClient.delete(`/distribution/withdraw-banks/${id}`);
}

export interface DistributionPrivilege {
  created_at: string;
  id: number;
  img_url: string;
  sort: number;
  status: number;
  title: string;
  updated_at: string;
}

export type DistributionPrivilegeSaveInput = {
  img_url: string;
  sort?: number;
  status?: number;
  title: string;
};

export function listDistributionPrivilegesApi(params?: {
  keyword?: string;
  limit?: number;
  page?: number;
  status?: 0 | 1;
}) {
  return requestClient.get<DistributionPage<DistributionPrivilege>>(
    '/distribution/privileges',
    { params },
  );
}

export function createDistributionPrivilegeApi(
  input: DistributionPrivilegeSaveInput,
) {
  return requestClient.post<DistributionPrivilege>(
    '/distribution/privileges',
    input,
  );
}

export function updateDistributionPrivilegeApi(
  id: number,
  input: DistributionPrivilegeSaveInput,
) {
  return requestClient.put<DistributionPrivilege>(
    `/distribution/privileges/${id}`,
    input,
  );
}

export function setDistributionPrivilegeStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<DistributionPrivilege>(
    `/distribution/privileges/${id}/status`,
    { status },
  );
}

export function deleteDistributionPrivilegeApi(id: number) {
  return requestClient.delete(`/distribution/privileges/${id}`);
}

export interface DistributionPoster {
  created_at: string;
  id: number;
  name: string;
  pic_url: string;
  sort: number;
  status: number;
  updated_at: string;
}

export type DistributionPosterSaveInput = {
  name: string;
  pic_url: string;
  sort?: number;
  status?: number;
};

export function listDistributionPostersApi(params?: {
  keyword?: string;
  limit?: number;
  page?: number;
  status?: 0 | 1;
}) {
  return requestClient.get<DistributionPage<DistributionPoster>>(
    '/distribution/posters',
    { params },
  );
}

export function createDistributionPosterApi(input: DistributionPosterSaveInput) {
  return requestClient.post<DistributionPoster>(
    '/distribution/posters',
    input,
  );
}

export function updateDistributionPosterApi(
  id: number,
  input: DistributionPosterSaveInput,
) {
  return requestClient.put<DistributionPoster>(
    `/distribution/posters/${id}`,
    input,
  );
}

export function setDistributionPosterStatusApi(id: number, status: 0 | 1) {
  return requestClient.put<DistributionPoster>(
    `/distribution/posters/${id}/status`,
    { status },
  );
}

export function deleteDistributionPosterApi(id: number) {
  return requestClient.delete(`/distribution/posters/${id}`);
}
