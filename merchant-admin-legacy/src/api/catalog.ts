import { http } from './http';

export interface Product {
  product_id: number;
  mer_id: number;
  store_name: string;
  store_info: string;
  keyword: string;
  cate_id: number;
  cate_name?: string;
  brand_id: number;
  status: number;
  is_show: number;
  unit_name: string;
  price: number;
  ot_price: number;
  stock: number;
  sales: number;
  spec_type: number;
  type: number;
  delivery_way: string;
  image: string;
  slider_image: string;
  refusal: string;
  svip_price_type?: number;
  svip_price?: number;
  mer_svip_status?: number;
}

export interface CategoryNode {
  store_category_id: number;
  pid: number;
  cate_name: string;
  children?: CategoryNode[];
}

export interface Brand {
  brand_id: number;
  brand_name: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export type ProductSave = {
  store_name: string;
  store_info?: string;
  keyword?: string;
  cate_id: number;
  brand_id?: number;
  unit_name?: string;
  price: number;
  ot_price?: number;
  cost?: number;
  stock: number;
  image?: string;
  slider_image?: string;
  delivery_way?: string;
  type?: number;
  spec_type?: number;
  is_show?: number;
  svip_price_type?: number;
  svip_price?: number;
  mer_svip_status?: number;
};

export function fetchProducts(params: Record<string, unknown>) {
  return http.get<PageResult<Product>>('/products', { params });
}

export function fetchProduct(id: number) {
  return http.get<Product>(`/products/${id}`);
}

export function createProduct(body: ProductSave) {
  return http.post<Product>('/products', body);
}

export function updateProduct(id: number, body: ProductSave) {
  return http.put<Product>(`/products/${id}`, body);
}

export function deleteProduct(id: number) {
  return http.delete(`/products/${id}`);
}

export function setProductShow(id: number, is_show: boolean) {
  return http.put(`/products/${id}/show`, { is_show: is_show ? 1 : 0 });
}

export function setProductStock(id: number, stock: number) {
  return http.put(`/products/${id}/stock`, { stock });
}

export function fetchCategories() {
  return http.get<{ list: CategoryNode[] }>('/product-categories');
}

export function fetchBrands() {
  return http.get<{ list: Brand[] }>('/brands');
}
