import { requestClient } from '#/api/request';

export interface ProductCacheListItem {
  id: string;
  name: string;
  enabled: boolean;
  remark: string;
}

export function fetchPriceDescriptions() {
  return requestClient.get<{ list: ProductCacheListItem[] }>('/product/price-descriptions');
}

export function savePriceDescriptions(list: ProductCacheListItem[]) {
  return requestClient.put<{ list: ProductCacheListItem[] }>('/product/price-descriptions', { list });
}

export function fetchActivityLabels() {
  return requestClient.get<{ list: ProductCacheListItem[] }>('/product/activity-labels');
}

export function saveActivityLabels(list: ProductCacheListItem[]) {
  return requestClient.put<{ list: ProductCacheListItem[] }>('/product/activity-labels', { list });
}
