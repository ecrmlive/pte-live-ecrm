import { requestClient } from '#/api/request';

import type { EnumField, PaginatedList } from '#/api/core/member';

export interface RechargePlanItem {
  give_money: number | string;
  money: number | string;
  plan_id: number;
  plan_name: string;
  real_money: number | string;
  sort?: number | string;
}

export interface RechargePlanFormPayload {
  give_money: number | string;
  money: number | string;
  plan_id?: number;
  plan_name: string;
  real_money: number | string;
  sort: number | string;
}

export interface RechargeRecordItem {
  create_time: string;
  order_no: string;
  pay_price: number | string;
  pay_status: EnumField;
  pay_type: EnumField;
  real_money: number | string;
  snapshot?: { plan_name?: string };
  user?: { avatarUrl?: string; nickName?: string; user_id: number };
}

export interface CashPayTypeOption {
  id: number | string;
  name: string;
}

export interface CashSettingValues {
  alipay_type?: string;
  cash_ratio?: string;
  is_open?: string;
  min_money?: string;
  pay_type?: Array<number | string>;
  scene_id?: string;
}

export interface MemberCashItem {
  alipay_account?: string;
  alipay_name?: string;
  apply_status: EnumField;
  apply_time?: string;
  audit_time?: string;
  avatarUrl?: string;
  bank_account?: string;
  bank_card?: string;
  bank_name?: string;
  cancelStatus?: number;
  create_time: string;
  id: number;
  money: number | string;
  nickName?: string;
  pay_type: EnumField;
  real_money: number | string;
  real_name?: string;
  reject_reason?: string;
  user_id: number;
}

export interface MemberCashListQuery {
  apply_status?: number | string;
  list_rows?: number;
  page?: number;
  pay_type?: number | string;
  search?: string;
  user_id?: number | string;
}

function saveExportBlob(blob: Blob, prefix: string) {
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `${prefix}-${Date.now()}.csv`;
  link.click();
  URL.revokeObjectURL(url);
}

export async function getRechargePlanListApi() {
  return requestClient.post<{ list: RechargePlanItem[] }>('/shop/user.plan/index', {});
}

export async function addRechargePlanApi(payload: RechargePlanFormPayload) {
  return requestClient.get('/shop/user.plan/add', { params: payload });
}

export async function editRechargePlanApi(payload: RechargePlanFormPayload & { plan_id: number }) {
  return requestClient.post('/shop/user.plan/edit', payload);
}

export async function deleteRechargePlanApi(planId: number) {
  return requestClient.get('/shop/user.plan/delete', { params: { plan_id: planId } });
}

export async function getRechargeRecordListApi(params: {
  list_rows?: number;
  page?: number;
  search?: string;
  value1?: '' | string[];
}) {
  return requestClient.post<{ list: PaginatedList<RechargeRecordItem> }>(
    '/shop/user.plan/log',
    params,
  );
}

export async function getMemberCashSettingApi() {
  return requestClient.get<{
    pay_type: CashPayTypeOption[];
    values: CashSettingValues;
  }>('/shop/user.cash/setting');
}

export async function saveMemberCashSettingApi(payload: CashSettingValues) {
  return requestClient.post('/shop/user.cash/setting', payload);
}

export async function getMemberCashListApi(params: MemberCashListQuery) {
  return requestClient.post<{
    alipay_type: number | string;
    list: PaginatedList<MemberCashItem>;
  }>('/shop/user.cash/index', params);
}

export async function auditMemberCashApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/user.cash/audit', payload);
}

export async function confirmMemberCashPayApi(id: number) {
  return requestClient.post('/shop/user.cash/money', { id });
}

export async function memberCashWxPayApi(id: number) {
  return requestClient.post('/shop/user.cash/wxpay', { id });
}

export async function cancelMemberCashWxPayApi(id: number) {
  return requestClient.post('/shop/user.cash/cancel', { id });
}

export async function exportMemberCashApi(params: MemberCashListQuery) {
  const blob = await requestClient.download<Blob>('/shop/user.cash/export', {
    params,
  });
  saveExportBlob(blob, 'member-cash');
}
