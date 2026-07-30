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

export interface Buying {
  group_buying_id: number;
  product_group_id: number;
  yet_buying_num: number;
  buying_count_num: number;
  remain: number;
  status: number;
  end_time: number;
}

export function fetchGroups(page = 1, limit = 20) {
  return http.get<{ list: ProductGroup[]; total: number }>(
    `/combination/groups?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchGroup(id: number) {
  return http.get<ProductGroup>(`/combination/groups/${id}`, false);
}

export function fetchBuyings(productGroupID: number) {
  return http.get<{ list: Buying[] }>(`/combination/groups/${productGroupID}/buyings`, false);
}

export function groupCheck(data: {
  product_group_id: number;
  group_buying_id?: number;
  cart_num?: number;
  address_id?: number;
}) {
  return http.post<{
    pay_price: number;
    price: number;
    store_name: string;
    product_group_id: number;
  }>("/order/group/check", data, true);
}

export function groupCreate(data: {
  product_group_id: number;
  group_buying_id?: number;
  cart_num?: number;
  address_id: number;
}) {
  return http.post<{ group_order_id: number; pay_price: number }>("/order/group/create", data, true);
}
