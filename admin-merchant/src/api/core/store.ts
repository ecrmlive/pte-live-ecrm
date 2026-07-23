import { requestClient } from '#/api/request';

import type { PaginatedList } from './member';

export interface EnumField {
  text: string;
  value: number;
}

export interface StoreListItem {
  create_time: string;
  detail_address?: string;
  is_check: EnumField;
  linkman: string;
  logo?: { file_path: string };
  phone: string;
  shop_hours: string;
  status: EnumField;
  store_id: number;
  store_name: string;
}

export interface StoreListQuery {
  is_check?: number;
  list_rows?: number;
  page?: number;
  status?: number;
  store_name?: string;
  value1?: string[];
}

export interface StoreOption {
  store_id: number | string;
  store_name: string;
}

export interface StoreFormPayload {
  address: string;
  city_id: number | string;
  coordinate?: string;
  is_check: number;
  linkman: string;
  logo_image_id: number | string;
  phone: string;
  province_id: number | string;
  region_id: number | string;
  shop_hours: string;
  sort: number;
  status: number;
  store_id?: number;
  store_name: string;
  summary?: string;
}

export interface StoreAddMeta {
  key?: string;
  regionData: import('./data').RegionTree;
}

export interface StoreEditMeta extends StoreAddMeta {
  model: StoreFormPayload & {
    coordinate?: string;
    logo?: { file_path: string };
  };
}

export interface ClerkListItem {
  clerk_id: number;
  create_time: string;
  mobile: string;
  real_name: string;
  status: EnumField;
  store?: { store_name: string };
  store_id: number;
  user?: { avatarUrl?: string; nickName: string; user_id: number };
  user_id: number;
}

export interface ClerkListQuery {
  list_rows?: number;
  page?: number;
  search?: string;
  store_id?: number | string;
}

export interface ClerkFormPayload {
  clerk_id?: number;
  mobile: string;
  real_name: string;
  status: number;
  store_id: number | string;
  user_id?: number | string;
}

export interface StoreOrderListItem {
  cancel_time?: string;
  clerk?: { real_name: string };
  create_time: string;
  id: number;
  order?: { order_no: string };
  order_id: number;
  order_type: EnumField;
  store?: { store_name: string };
  verify_num: number;
  verify_status: number;
}

export interface StoreOrderListQuery {
  create_time?: string[];
  list_rows?: number;
  order_no?: string;
  page?: number;
  search?: string;
  store_id?: number | string;
}

export async function getStoreListApi(params: StoreListQuery) {
  return requestClient.post<{ list: PaginatedList<StoreListItem> }>(
    '/shop/store.store/index',
    params,
  );
}

export async function getStoreAddMetaApi() {
  return requestClient.get<StoreAddMeta>('/shop/store.store/add');
}

export async function addStoreApi(payload: StoreFormPayload) {
  return requestClient.post('/shop/store.store/add', payload);
}

export async function getStoreEditMetaApi(storeId: number) {
  return requestClient.get<StoreEditMeta>('/shop/store.store/edit', {
    params: { store_id: storeId },
  });
}

export async function editStoreApi(payload: StoreFormPayload) {
  return requestClient.post('/shop/store.store/edit', payload);
}

export async function deleteStoreApi(storeId: number) {
  return requestClient.post('/shop/store.store/delete', { store_id: storeId });
}

export async function getClerkListApi(params: ClerkListQuery) {
  return requestClient.post<{ list: PaginatedList<ClerkListItem>; store_list: StoreOption[] }>(
    '/shop/store.clerk/index',
    params,
  );
}

export async function getClerkAddMetaApi() {
  return requestClient.get<{ store_list: StoreOption[] }>('/shop/store.clerk/add');
}

export async function addClerkApi(payload: ClerkFormPayload) {
  return requestClient.post('/shop/store.clerk/add', payload);
}

export async function getClerkEditMetaApi(clerkId: number) {
  return requestClient.get<{
    detail: ClerkListItem;
    store_list: StoreOption[];
  }>('/shop/store.clerk/edit', { params: { clerk_id: clerkId } });
}

export async function editClerkApi(payload: ClerkFormPayload & { clerk_id: number }) {
  return requestClient.post('/shop/store.clerk/edit', payload);
}

export async function deleteClerkApi(clerkId: number) {
  return requestClient.post('/shop/store.clerk/delete', { clerk_id: clerkId });
}

export async function getStoreOrderListApi(params: StoreOrderListQuery) {
  return requestClient.post<{
    list: PaginatedList<StoreOrderListItem>;
    store_list: StoreOption[];
  }>('/shop/store.order/index', params);
}

export async function cancelStoreOrderApi(id: number) {
  return requestClient.post('/shop/store.order/cancel', { id });
}
