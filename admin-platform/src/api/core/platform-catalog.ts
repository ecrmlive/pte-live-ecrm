import { requestClient } from '#/api/request';

export interface PlatformProduct {
  cate_name?: string;
  create_time: string;
  image: string;
  is_show: number;
  mer_id: number;
  mer_name?: string;
  ot_price: number;
  price: number;
  product_id: number;
  refusal?: string;
  sales: number;
  status: number;
  stock: number;
  store_info: string;
  store_name: string;
  title: string;
}

export interface PlatformProductPage {
  limit: number;
  list: PlatformProduct[];
  page: number;
  total: number;
}

export function listPlatformProductsApi(params: { keyword?: string; limit: number; mer_id?: number; page: number; status?: number }) {
  return requestClient.get<PlatformProductPage>('/products', { params });
}

export function getPlatformProductApi(id: number) {
  return requestClient.get<PlatformProduct>(`/products/${id}`);
}

export function auditPlatformProductApi(id: number, body: { refusal?: string; status: number }) {
  return requestClient.post(`/products/${id}/audit`, body);
}

export interface PlatformCategory {
  cate_name: string;
  children?: PlatformCategory[];
  is_show: number;
  level: number;
  pid: number;
  sort: number;
  store_category_id: number;
}

export interface PlatformBrand {
  brand_id: number;
  brand_name: string;
  category_id: number;
  is_show: number;
  sort: number;
}
export interface PlatformBrandCategory {
  brand_category_id: number;
  cate_name: string;
  children?: PlatformBrandCategory[];
  is_show: number;
  pid: number;
  sort: number;
}
export interface CatalogSaveInput {
  brand_name?: string;
  cate_name?: string;
  category_id?: number;
  is_show: number;
  pid?: number;
  sort: number;
}

export function listPlatformCategoriesApi() { return requestClient.get<{ list: PlatformCategory[] }>('/product-categories'); }
export function createPlatformCategoryApi(data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.post('/product-categories', data); }
export function updatePlatformCategoryApi(id: number, data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'sort'>>) { return requestClient.put(`/product-categories/${id}`, data); }
export function deletePlatformCategoryApi(id: number) { return requestClient.delete(`/product-categories/${id}`); }
export function listPlatformBrandCategoriesApi() { return requestClient.get<{ list: PlatformBrandCategory[] }>('/brand-categories'); }
export function createPlatformBrandCategoryApi(data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.post('/brand-categories', data); }
export function updatePlatformBrandCategoryApi(id: number, data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.put(`/brand-categories/${id}`, data); }
export function deletePlatformBrandCategoryApi(id: number) { return requestClient.delete(`/brand-categories/${id}`); }
export function listPlatformBrandsApi(params?: { category_id?: number }) { return requestClient.get<{ list: PlatformBrand[] }>('/brands', { params }); }
export function createPlatformBrandApi(data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>> & { category_id?: number }) { return requestClient.post('/brands', data); }
export function updatePlatformBrandApi(id: number, data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>> & { category_id?: number }) { return requestClient.put(`/brands/${id}`, data); }
export function deletePlatformBrandApi(id: number) { return requestClient.delete(`/brands/${id}`); }
