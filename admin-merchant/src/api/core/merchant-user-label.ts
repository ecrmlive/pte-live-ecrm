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

export function listStoreUserLabelsApi() {
  return requestClient.get<{ list: StoreUserLabel[] }>('/store-user-labels');
}

export function saveStoreUserLabelsApi(body: { list: StoreUserLabel[] }) {
  return requestClient.post<{ list: StoreUserLabel[] }>('/store-user-labels', body);
}

export function listStoreAutoLabelRulesApi() {
  return requestClient.get<{ list: StoreAutoLabelRule[] }>('/store-auto-label-rules');
}

export function saveStoreAutoLabelRulesApi(body: { list: StoreAutoLabelRule[] }) {
  return requestClient.post<{ list: StoreAutoLabelRule[] }>('/store-auto-label-rules', body);
}
