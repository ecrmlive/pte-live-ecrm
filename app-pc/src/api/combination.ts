import { http } from "@/utils/request";

export interface ProductGroup {
  product_group_id: number;
  product_id: number;
  price: number;
  buying_count_num: number;
  store_name?: string;
  mer_name?: string;
  ot_price?: number;
  image?: string;
}

export function fetchGroups(page = 1, limit = 20) {
  return http.get<{ list: ProductGroup[]; total: number }>(
    `/combination/groups?page=${page}&limit=${limit}`,
    false,
  );
}
