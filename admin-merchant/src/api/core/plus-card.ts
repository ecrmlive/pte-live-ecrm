import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface CardCategoryItem {
  category_id: number;
  create_time?: string;
  name: string;
  sort?: number;
}

export interface CardListItem {
  card_id: number;
  card_price: number | string;
  card_sort: number;
  card_status: number;
  card_title: string;
  category?: { name: string };
  create_time: string;
  image?: { file_path: string };
  sell_num: number;
  sell_price: number | string;
  wait_num: number;
}

export interface CardProductOption {
  product_id: number;
  product_name: string;
}

export interface CardFormValues {
  card_content: string;
  card_id?: number;
  card_price: number | string;
  card_sort: number | string;
  card_status: number;
  card_title: string;
  category_id: number | string;
  image_id: number | string;
  product_attr: string;
  product_id: number | string;
  sell_price: number | string;
  stock_num?: number | string;
}

export interface CardCodeGenerateForm {
  card_id: number;
  code_count: number | string;
  code_len: number | string;
  end_time: string;
  prefix: string;
  start_num: number | string;
  start_time: string;
}

export interface CardCodeEditForm {
  code_id: number;
  code_pwd: string;
  code_status: number;
  end_time: string;
  is_delete: number;
  start_time: string;
}

export async function getCardListApi(params: { list_rows?: number; page?: number }) {
  return requestClient.post<{ list: PaginatedList<CardListItem> }>(
    '/shop/plus.card.card/index',
    params,
  );
}

export async function deleteCardApi(cardId: number) {
  return requestClient.post('/shop/plus.card.card/delete', { card_id: cardId });
}

export async function getCardCategoryListApi() {
  return requestClient.post<{ category: CardCategoryItem[] }>(
    '/shop/plus.card.category/index',
    {},
  );
}

export async function addCardCategoryApi(payload: { name: string; sort: number | string }) {
  return requestClient.post('/shop/plus.card.category/add', payload);
}

export async function editCardCategoryApi(payload: CardCategoryItem) {
  return requestClient.post('/shop/plus.card.category/edit', payload);
}

export async function deleteCardCategoryApi(categoryId: number) {
  return requestClient.post('/shop/plus.card.category/delete', { category_id: categoryId });
}

export async function getCardAddMetaApi() {
  return requestClient.get<{ category: CardCategoryItem[]; product: CardProductOption[] }>(
    '/shop/plus.card.card/add',
  );
}

export async function addCardApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.card.card/add', payload);
}

export async function getCardEditMetaApi(cardId: number) {
  return requestClient.get<{ category: CardCategoryItem[]; model: CardFormValues; product: CardProductOption[] }>(
    '/shop/plus.card.card/edit',
    { params: { card_id: cardId } },
  );
}

export async function editCardApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.card.card/edit', payload);
}

export async function getCardCodeGenerateMetaApi(cardId: number) {
  return requestClient.get<{ model: { card_title: string } }>('/shop/plus.card.card/code', {
    params: { card_id: cardId },
  });
}

export async function generateCardCodeApi(payload: CardCodeGenerateForm) {
  return requestClient.post('/shop/plus.card.card/code', payload);
}

export async function getCardCodeEditMetaApi(codeId: number) {
  return requestClient.get<{ model: CardCodeEditForm & { card?: { card_title: string }; code_no?: string } }>(
    '/shop/plus.card.code/edit',
    { params: { code_id: codeId } },
  );
}

export async function editCardCodeApi(payload: Record<string, unknown>) {
  return requestClient.post('/shop/plus.card.code/edit', payload);
}

export async function getCardSettingApi() {
  return requestClient.get<{ vars: { values: { image: string } } }>(
    '/shop/plus.card.setting/index',
  );
}

export async function saveCardSettingApi(payload: { image: string }) {
  return requestClient.post('/shop/plus.card.setting/index', payload);
}

export interface CardCodeListItem {
  card_title?: string;
  code_id: number;
  code_no: string;
  code_pwd?: string;
  code_status?: number;
  create_time?: string;
  end_time_str?: string;
  is_delete?: number;
  start_time_str?: string;
}

export interface CardCodeListQuery {
  card_id?: number | string;
  category_id?: number | string;
  code_no?: string;
  code_status?: number | string;
  is_delete?: number | string;
  list_rows?: number;
  page?: number;
}

export async function getCardCodeListApi(params: CardCodeListQuery) {
  return requestClient.post<{
    cardList: Array<{ card_id: number; card_title: string }>;
    categoryList: CardCategoryItem[];
    list: PaginatedList<CardCodeListItem>;
  }>('/shop/plus.card.code/index', params);
}

function saveCardExportBlob(blob: Blob, prefix: string) {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = `${prefix}-${Date.now()}.xlsx`;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(url);
}

export async function exportCardCodeApi(params: Omit<CardCodeListQuery, 'list_rows' | 'page'>) {
  const blob = await requestClient.download<Blob>('/shop/plus.card.code/export', {
    params,
  });
  saveCardExportBlob(blob, 'card-codes');
}

export interface CardOrderListItem {
  create_time?: string;
  delivery_status?: number;
  mobile?: string;
  name?: string;
  order_id: number;
  order_no?: string;
  product_name?: string;
  card?: { card_title?: string };
  code?: { code_no?: string };
  user?: { nickName?: string; user_id?: number };
}

export interface CardOrderListQuery {
  category_id?: number | string;
  code_no?: string;
  create_time?: string[];
  dataType?: string;
  delivery_status?: number | string;
  list_rows?: number;
  order_no?: string;
  page?: number;
}

export async function getCardOrderListApi(params: CardOrderListQuery) {
  return requestClient.get<{
    categoryList: CardCategoryItem[];
    list: PaginatedList<CardOrderListItem>;
    order_count: { wait: number };
  }>('/shop/plus.card.order/index', { params });
}

export interface CardOrderDetail {
  delivery_status?: number;
  delivery_time?: string;
  detail?: string;
  express?: { express_name?: string };
  express_no?: string;
  mobile?: string;
  name?: string;
  order_no?: string;
  product_attr?: string;
  product_image?: string;
  product_name?: string;
  product_price?: number | string;
  region?: { city?: string; province?: string; region?: string };
  user?: { nickName?: string; user_id?: number };
}

export async function getCardOrderDetailApi(orderId: number) {
  return requestClient.post<{
    detail: CardOrderDetail;
    expressList: Array<{ express_id: number; express_name: string }>;
    shopClerkList?: unknown[];
  }>('/shop/plus.card.order/detail', { order_id: orderId });
}

export async function deliverCardOrderApi(payload: {
  order_id: number;
  param: { express_id: number; express_no: string };
}) {
  return requestClient.post('/shop/plus.card.order/delivery', payload);
}
