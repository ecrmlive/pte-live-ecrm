import { requestClient } from '#/api/request';

export interface StoreUserLabel {
  label_id: number;
  name: string;
  sort: number;
  status: number;
}

export interface StoreAutoLabelRule {
  name: string;
  rule_id: number;
  rule_type: string;
  status: number;
}

export interface StoreUserLabelListParams {
  keyword?: string;
  status?: 0 | 1;
}

export function listStoreUserLabelsApi(params: StoreUserLabelListParams = {}) {
  return requestClient.get<{ list: StoreUserLabel[] }>('/store-user-labels', { params });
}

export function saveStoreUserLabelsApi(body: { list: StoreUserLabel[] }) {
  return requestClient.post<{ list: StoreUserLabel[] }>('/store-user-labels', body);
}

export interface StoreAutoLabelRuleListParams {
  keyword?: string;
  status?: 0 | 1;
}

export function listStoreAutoLabelRulesApi(params: StoreAutoLabelRuleListParams = {}) {
  return requestClient.get<{ list: StoreAutoLabelRule[] }>('/store-auto-label-rules', { params });
}

export function saveStoreAutoLabelRulesApi(body: { list: StoreAutoLabelRule[] }) {
  return requestClient.post<{ list: StoreAutoLabelRule[] }>('/store-auto-label-rules', body);
}
