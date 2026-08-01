import { http } from "@/utils/request";

export interface CategoryItem {
  id: number;
  name: string;
  pid: number;
  pic?: string;
}

export interface ProductItem {
  id: number;
  /** 新消费视图字段；旧夹具/接口仍可仅提供 store_name。 */
  title?: string;
  mer_id: number;
  mer_name: string;
  store_name: string;
  shop_name?: string;
  category_id?: number;
  image: string;
  price: string;
  ot_price?: string;
  sales: number;
  stock: number;
}

export type ProductSort = "default" | "sales" | "price";

export interface ProductPage {
  list: ProductItem[];
  total: number;
  page: number;
  limit: number;
}

/** 平台首页装修中的“商品组 → 自动获取”投影。 */
export interface HomeDisplayType {
  category_id: number;
  initial_limit: number;
  sort: ProductSort;
}

export interface HomePayload {
  banners: { id: number; title: string; image: string; url?: string }[];
  display_types: HomeDisplayType[];
  hot: ProductItem[];
}

export interface StoreDirectoryItem {
  store_id: number;
  mer_id: number;
  name: string;
  product_count: number;
  sales_count: number;
  cover_url?: string;
}

export interface ProductDetail extends ProductItem {
  unit_name: string;
  store_info: string;
  slider_image: string[];
  spec_type: number;
  delivery_way: string;
  merchant_app_id?: string;
}

export function fetchCategories() {
  return http.get<CategoryItem[]>("/catalog/categories", false);
}

export function fetchProducts(params?: {
  cate_id?: number;
  keyword?: string;
  min_price?: number;
  max_price?: number;
  page?: number;
  limit?: number;
  mer_id?: number;
  sort?: ProductSort;
  order?: "asc" | "desc";
}) {
  const q = new URLSearchParams();
  if (params?.cate_id) q.set("cate_id", String(params.cate_id));
  if (params?.keyword) q.set("keyword", params.keyword);
  if (Number.isFinite(params?.min_price) && (params?.min_price || 0) >= 0) q.set("min_price", String(params?.min_price));
  if (Number.isFinite(params?.max_price) && (params?.max_price || 0) >= 0) q.set("max_price", String(params?.max_price));
  if (params?.page) q.set("page", String(params.page));
  if (params?.limit) q.set("limit", String(params.limit));
  if (params?.mer_id) q.set("mer_id", String(params.mer_id));
  if (params?.sort && params.sort !== "default") q.set("sort", params.sort);
  if (params?.order) q.set("order", params.order);
  const qs = q.toString();
  return http.get<ProductPage>(
    `/catalog/products${qs ? `?${qs}` : ""}`,
    false,
  );
}

export function fetchStores() {
  return http.get<{ list: StoreDirectoryItem[]; total: number }>("/catalog/stores", false);
}

export function fetchProductDetail(id: number) {
  return http.get<ProductDetail>(`/catalog/products/${id}`, false);
}

export function fetchHome() {
  return http.get<HomePayload>("/catalog/home", false);
}

export function fetchStoreHome(merId: number) {
  return http.get<{ mer_id: number; mer_name: string; store_id: number; merchant_app_id: string; products: ProductItem[]; total: number }>(
    `/catalog/stores/${merId}`,
    false,
  );
}
