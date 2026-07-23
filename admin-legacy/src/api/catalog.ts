import { http } from './http';

export interface CategoryNode {
  store_category_id: number;
  pid: number;
  cate_name: string;
  sort: number;
  is_show: number;
  level: number;
  children?: CategoryNode[];
}

export interface Brand {
  brand_id: number;
  brand_name: string;
  sort: number;
  is_show: number;
  pic: string;
}

export interface Product {
  product_id: number;
  mer_id: number;
  mer_name?: string;
  store_name: string;
  store_info: string;
  status: number;
  is_show: number;
  cate_id: number;
  cate_name?: string;
  price: number;
  stock: number;
  refusal: string;
  type: number;
  product_type: number;
  delivery_way: string;
  spec_type: number;
  create_time: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchCategoryTree() {
  return http.get<{ list: CategoryNode[] }>('/product-categories');
}

export function createCategory(body: { pid?: number; cate_name: string; sort?: number; is_show?: number }) {
  return http.post('/product-categories', body);
}

export function updateCategory(id: number, body: { cate_name: string; sort?: number; is_show?: number }) {
  return http.put(`/product-categories/${id}`, body);
}

export function deleteCategory(id: number) {
  return http.delete(`/product-categories/${id}`);
}

export function fetchBrands() {
  return http.get<{ list: Brand[] }>('/brands');
}

export function createBrand(body: { brand_name: string; sort?: number; is_show?: number }) {
  return http.post('/brands', body);
}

export function updateBrand(id: number, body: { brand_name: string; sort?: number; is_show?: number }) {
  return http.put(`/brands/${id}`, body);
}

export function deleteBrand(id: number) {
  return http.delete(`/brands/${id}`);
}

export function fetchProducts(params: Record<string, unknown>) {
  return http.get<PageResult<Product>>('/products', { params });
}

export function auditProduct(id: number, body: { status: number; refusal?: string }) {
  return http.post(`/products/${id}/audit`, body);
}
