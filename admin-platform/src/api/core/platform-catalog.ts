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

export interface PlatformBrand { brand_id: number; brand_name: string; is_show: number; sort: number; }
export interface CatalogSaveInput { cate_name?: string; brand_name?: string; is_show: number; pid?: number; sort: number; }

export function listPlatformCategoriesApi() { return requestClient.get<{ list: PlatformCategory[] }>('/product-categories'); }
export function createPlatformCategoryApi(data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'pid' | 'sort'>>) { return requestClient.post('/product-categories', data); }
export function updatePlatformCategoryApi(id: number, data: Required<Pick<CatalogSaveInput, 'cate_name' | 'is_show' | 'sort'>>) { return requestClient.put(`/product-categories/${id}`, data); }
export function deletePlatformCategoryApi(id: number) { return requestClient.delete(`/product-categories/${id}`); }
export function listPlatformBrandsApi() { return requestClient.get<{ list: PlatformBrand[] }>('/brands'); }
export function createPlatformBrandApi(data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>>) { return requestClient.post('/brands', data); }
export function updatePlatformBrandApi(id: number, data: Required<Pick<CatalogSaveInput, 'brand_name' | 'is_show' | 'sort'>>) { return requestClient.put(`/brands/${id}`, data); }
export function deletePlatformBrandApi(id: number) { return requestClient.delete(`/brands/${id}`); }
