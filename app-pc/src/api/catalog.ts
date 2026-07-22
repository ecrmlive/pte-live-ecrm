import { http } from "@/utils/request";

export interface CategoryItem {
  id: number;
  name: string;
  pid: number;
  pic?: string;
}

export interface ProductItem {
  id: number;
  mer_id: number;
  mer_name: string;
  store_name: string;
  image: string;
  price: string;
  ot_price?: string;
  sales: number;
  stock: number;
}

export interface ProductDetail extends ProductItem {
  unit_name: string;
  store_info: string;
  slider_image: string[];
  spec_type: number;
  delivery_way: string;
}

export function fetchCategories() {
  return http.get<CategoryItem[]>("/catalog/categories", false);
}

export function fetchProducts(params?: {
  cate_id?: number;
  keyword?: string;
  page?: number;
  mer_id?: number;
}) {
  const q = new URLSearchParams();
  if (params?.cate_id) q.set("cate_id", String(params.cate_id));
  if (params?.keyword) q.set("keyword", params.keyword);
  if (params?.page) q.set("page", String(params.page));
  if (params?.mer_id) q.set("mer_id", String(params.mer_id));
  const qs = q.toString();
  return http.get<{ list: ProductItem[]; total: number }>(
    `/catalog/products${qs ? `?${qs}` : ""}`,
    false,
  );
}

export function fetchProductDetail(id: number) {
  return http.get<ProductDetail>(`/catalog/products/${id}`, false);
}

export function fetchHome() {
  return http.get<{ banners: { id: number; title: string; image: string }[]; hot: ProductItem[] }>(
    "/catalog/home",
    false,
  );
}

export function fetchStoreHome(merId: number) {
  return http.get<{ mer_id: number; mer_name: string; products: ProductItem[]; total: number }>(
    `/catalog/stores/${merId}`,
    false,
  );
}
