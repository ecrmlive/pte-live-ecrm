import { requestClient } from '#/api/request';

export interface MerchantCategoryNode {
  cate_name: string;
  children?: MerchantCategoryNode[];
  is_show: number;
  store_category_id: number;
}

export interface MerchantProduct {
  cate_id: number;
  cate_name?: string;
  create_time: string;
  image: string;
  is_show: number;
  keyword: string;
  ot_price: number;
  price: number;
  product_id: number;
  refusal?: string;
  status: number;
  stock: number;
  store_info: string;
  store_name: string;
  unit_name: string;
}

export interface MerchantProductSaveInput {
  cate_id: number;
  cost: number;
  image: string;
  is_show?: number;
  keyword: string;
  ot_price: number;
  price: number;
  slider_image: string;
  spec_type: number;
  stock: number;
  store_info: string;
  store_name: string;
  type: number;
  unit_name: string;
}

export interface MerchantProductPage {
  limit: number;
  list: MerchantProduct[];
  page: number;
  total: number;
}

export function listMerchantProductsApi(params: { keyword?: string; limit: number; page: number; status?: number }) {
  return requestClient.get<MerchantProductPage>('/products', { params });
}

export function getMerchantProductCategoriesApi() {
  return requestClient.get<{ list: MerchantCategoryNode[] }>('/product-categories');
}

export function createMerchantProductApi(body: MerchantProductSaveInput) {
  return requestClient.post<MerchantProduct>('/products', body);
}

export function updateMerchantProductApi(id: number, body: MerchantProductSaveInput) {
  return requestClient.put<MerchantProduct>(`/products/${id}`, body);
}

export function deleteMerchantProductApi(id: number) {
  return requestClient.delete(`/products/${id}`);
}

export function setMerchantProductShowApi(id: number, isShow: number) {
  return requestClient.put(`/products/${id}/show`, { is_show: isShow });
}

export function setMerchantProductStockApi(id: number, stock: number) {
  return requestClient.put(`/products/${id}/stock`, { stock });
}
