import { requestClient } from '#/api/request';

export interface PlatformProduct {
  activity_labels?: string[];
  brand_name?: string;
  care_count?: number;
  cate_name?: string;
  create_time: string;
  ficti?: number;
  image: string;
  is_show: number;
  mer_cate_name?: string;
  mer_id: number;
  mer_name?: string;
  ot_price: number;
  price: number;
  product_id: number;
  product_type?: number;
  rank?: number;
  refusal?: string;
  sales: number;
  spec_type?: number;
  star?: number;
  status: number;
  stock: number;
  store_info: string;
  store_name: string;
  svip_price_type?: number;
  title: string;
}

export interface PlatformProductPage {
  limit: number;
  list: PlatformProduct[];
  page: number;
  total: number;
}

export interface PlatformProductStatusFilter {
  count: number;
  name: string;
  type: number;
}

export function listPlatformProductsApi(params: {
  brand_name?: string;
  cate_hot?: number;
  cate_id?: number;
  date_from?: string;
  date_to?: string;
  is_gift_bag?: 0 | 1;
  is_hot?: number;
  is_trader?: 0 | 1;
  is_used?: number;
  keyword?: string;
  limit: number;
  mer_category_id?: number;
  mer_id?: number;
  mer_type_id?: number;
  page: number;
  product_type?: number;
  star?: number;
  status?: number;
  store_name?: string;
  svip_price_type?: number;
  type?: number;
  us_status?: number;
}) {
  return requestClient.get<PlatformProductPage>('/products', { params });
}

export function getPlatformProductStatusFilterApi(params?: {
  brand_name?: string;
  cate_id?: number;
  is_gift_bag?: 0 | 1;
  keyword?: string;
  mer_category_id?: number;
  mer_id?: number;
  mer_type_id?: number;
  store_name?: string;
  svip_price_type?: number;
}) {
  return requestClient.get<{ list: PlatformProductStatusFilter[] }>(
    '/products/status-filter',
    { params },
  );
}

export function getPlatformProductApi(id: number) {
  return requestClient.get<PlatformProduct>(`/products/${id}`);
}

export interface PlatformProductEditSKU {
  bar_code?: string;
  code?: string;
  extension_one?: number;
  image?: string;
  ot_price: number;
  price: number;
  sku_id: number;
  spec: Record<string, string>;
  spec_text: string;
  status: number;
  stock: number;
  volume?: number;
  weight?: number;
}

export interface PlatformProductEditDetail {
  activity_labels?: string[];
  brand_name: string;
  care_count?: number;
  cate_hot: number;
  cate_id: number;
  cate_name: string;
  cate_path?: string;
  commission_text?: string;
  content: string;
  create_time: string;
  delivery_way: number[];
  ficti?: number;
  image: string;
  is_best: number;
  is_benefit: number;
  is_gift_bag?: number;
  is_hot: number;
  is_new: number;
  is_show: number;
  keyword: string;
  mer_cate_id?: number;
  mer_cate_name: string;
  mer_cate_options?: Array<{ id: number; name: string }>;
  mer_category_name: string;
  mer_form_id?: number | null;
  mer_id: number;
  mer_label_ids?: string[];
  mer_label_options?: Array<{ id: number; name: string }>;
  mer_labels: string[];
  mer_name: string;
  mer_params: Array<{ name: string; value: string }>;
  mer_recommend?: number;
  mer_type_name: string;
  once_min_count: number;
  ot_price: number;
  platform_params: Array<{ name: string; value: string }>;
  price: number;
  product_id: number;
  product_type: number;
  rank: number;
  refund_switch: number;
  refusal?: string;
  sales: number;
  skus: PlatformProductEditSKU[];
  slider_image: string[];
  spec_type: number;
  star: number;
  status: number;
  stock: number;
  store_id: number;
  store_info: string;
  store_name: string;
  svip_price?: number;
  svip_price_type: number;
  sys_label_names?: string[];
  sys_labels: string[];
  title: string;
  unit_name: string;
}

export interface PlatformProductOperateLog {
  action_label: string;
  created_at: string;
  id: number;
  index: number;
  operator_name: string;
  role_name: string;
  terminal: string;
}

export function getPlatformProductEditApi(id: number) {
  return requestClient.get<PlatformProductEditDetail>(`/products/${id}/edit`);
}

export interface PlatformProductStoreOption {
  merchant_id: number;
  merchant_name: string;
  store_id: number;
  store_name: string;
}

export interface PlatformProductStoreMeta {
  mer_cate_options: Array<{ id: number; name: string }>;
  mer_label_options: Array<{ id: number; name: string }>;
}

export type PlatformProductAdminSaveBody = {
  brand_name: string;
  cate_hot?: number;
  cate_id: number;
  content?: string;
  delivery_way: number[];
  image: string;
  is_best?: number;
  is_benefit?: number;
  is_hot?: number;
  is_new?: number;
  keyword: string;
  mer_cate_id: number;
  mer_label_ids: string[];
  once_min_count?: number;
  ot_price?: number;
  rank?: number;
  refund_switch?: number;
  skus: Array<{
    bar_code?: string;
    code?: string;
    extension_one?: number;
    image?: string;
    ot_price?: number;
    price: number;
    sku_id?: number;
    spec: Record<string, string>;
    status?: number;
    stock: number;
    volume?: number;
    weight?: number;
  }>;
  slider_image: string[];
  star?: number;
  store_info: string;
  sys_labels: string[];
  title: string;
  unit_name: string;
};

export function listPlatformProductStoresApi() {
  return requestClient.get<{ list: PlatformProductStoreOption[] }>('/product-stores');
}

export function getPlatformProductStoreOptionsApi(storeId: number) {
  return requestClient.get<PlatformProductStoreMeta>(`/product-stores/${storeId}/options`);
}

export function createPlatformProductAdminApi(
  body: PlatformProductAdminSaveBody & { store_id: number },
) {
  return requestClient.post<{ ok: boolean; product_id: number }>('/products', body);
}

export function updatePlatformProductAdminApi(
  id: number,
  body: PlatformProductAdminSaveBody,
) {
  return requestClient.put(`/products/${id}`, body);
}

export function listPlatformProductOperateLogsApi(
  id: number,
  params: {
    date_from?: string;
    date_to?: string;
    limit: number;
    page: number;
    terminal?: string;
  },
) {
  return requestClient.get<{
    limit: number;
    list: PlatformProductOperateLog[];
    page: number;
    total: number;
  }>(`/products/${id}/operate-logs`, { params });
}

export function auditPlatformProductApi(id: number, body: { refusal?: string; status: number }) {
  return requestClient.post(`/products/${id}/audit`, body);
}

export function setPlatformProductShowApi(id: number, status: 0 | 1) {
  return requestClient.post(`/products/${id}/show`, { status });
}

export function forceOffPlatformProductApi(id: number, reason: string) {
  return requestClient.post(`/products/${id}/force-off`, { reason });
}

export function batchForceOffPlatformProductsApi(ids: number[], reason: string) {
  return requestClient.post('/products/batch/force-off', { ids, reason });
}

export function batchShowPlatformProductsApi(ids: number[], status: 0 | 1) {
  return requestClient.post('/products/batch/show', { ids, status });
}

/** TODO backend: 批量标签 */
export function batchLabelsPlatformProductsApi(ids: number[], sysLabels: string) {
  return requestClient.post('/products/batch/labels', { ids, sys_labels: sysLabels });
}

/** TODO backend: 批量推荐 */
export function batchRecommendPlatformProductsApi(
  ids: number[],
  body: Record<string, number | string>,
) {
  return requestClient.post('/products/batch/recommend', { ids, ...body });
}

/** TODO backend: 批量复制到店铺 */
export function batchCopyPlatformProductsApi(productIds: number[], storeId: number) {
  return requestClient.post('/products/batch/copy', {
    product_ids: productIds,
    store_id: storeId,
  });
}

export function updatePlatformProductOpsApi(
  id: number,
  body: {
    cate_hot?: number;
    is_benefit?: number;
    is_best?: number;
    is_hot?: number;
    is_new?: number;
    rank?: number;
    star?: number;
  },
) {
  return requestClient.put(`/products/${id}/ops`, body);
}

/** 修改虚拟已售数量：type 1 增加 / 2 减少 */
export function setPlatformProductFictiApi(
  id: number,
  body: { ficti: number; type: 1 | 2 },
) {
  return requestClient.post<{ ficti: number; ok: boolean; sales: number }>(
    `/products/${id}/ficti`,
    body,
  );
}

export interface PlatformCategory {
  cate_name: string;
  children?: PlatformCategory[];
  create_time?: string;
  is_hot: number;
  is_show: number;
  level: number;
  pic?: string;
  pid: number;
  sort: number;
  store_category_id: number;
}

export interface PlatformBrand {
  brand_id: number;
  brand_name: string;
  category_id: number;
  create_time?: string;
  is_show: number;
  sort: number;
}
export interface PlatformBrandCategory {
  brand_category_id: number;
  cate_name: string;
  children?: PlatformBrandCategory[];
  create_time?: string;
  is_show: number;
  pid: number;
  sort: number;
}
export interface CatalogSaveInput {
  brand_name?: string;
  cate_name?: string;
  category_id?: number;
  is_hot?: number;
  is_show: number;
  pic?: string;
  pid?: number;
  sort: number;
}

export function listPlatformCategoriesApi() { return requestClient.get<{ list: PlatformCategory[] }>('/product-categories'); }
export function createPlatformCategoryApi(data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>> & Pick<CatalogSaveInput, 'pic' | 'is_hot'>) { return requestClient.post('/product-categories', data); }
export function updatePlatformCategoryApi(id: number, data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'sort'>> & Pick<CatalogSaveInput, 'pic' | 'is_hot' | 'pid'>) { return requestClient.put(`/product-categories/${id}`, data); }
export function updatePlatformCategoryStatusApi(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/product-categories/${id}/status`, {
    is_show: enabled ? 1 : 0,
  });
}
export function updatePlatformCategoryRecommendApi(id: number, enabled: boolean) {
  return requestClient.put<{ ok: boolean }>(`/product-categories/${id}/recommend`, {
    is_hot: enabled ? 1 : 0,
  });
}
export function deletePlatformCategoryApi(id: number) { return requestClient.delete(`/product-categories/${id}`); }
export function listPlatformBrandCategoriesApi() { return requestClient.get<{ list: PlatformBrandCategory[] }>('/brand-categories'); }
export function createPlatformBrandCategoryApi(data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.post('/brand-categories', data); }
export function updatePlatformBrandCategoryApi(id: number, data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.put(`/brand-categories/${id}`, data); }
export function deletePlatformBrandCategoryApi(id: number) { return requestClient.delete(`/brand-categories/${id}`); }
export function listPlatformBrandsApi(params?: { category_id?: number }) { return requestClient.get<{ list: PlatformBrand[] }>('/brands', { params }); }
export function createPlatformBrandApi(data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>> & { category_id?: number }) { return requestClient.post('/brands', data); }
export function updatePlatformBrandApi(id: number, data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>> & { category_id?: number }) { return requestClient.put(`/brands/${id}`, data); }
export function deletePlatformBrandApi(id: number) { return requestClient.delete(`/brands/${id}`); }
