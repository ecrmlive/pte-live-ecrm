import { http } from "@/utils/request";

export interface ProductAssist {
  product_assist_id: number;
  product_id: number;
  assist_price: number;
  assist_count: number;
  assist_user_count: number;
  stock: number;
  store_name?: string;
  mer_name?: string;
  ot_price?: number;
  image?: string;
}

export interface AssistSet {
  product_assist_set_id: number;
  product_assist_id: number;
  product_id: number;
  uid: number;
  status: number;
  assist_count: number;
  yet_assist_count: number;
  assist_price?: number;
  store_name?: string;
  nickname?: string;
  helpers?: { uid: number; nickname: string }[];
}

export function fetchAssists(page = 1, limit = 20) {
  return http.get<{ list: ProductAssist[]; total: number }>(
    `/assist/actives?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchAssist(id: number) {
  return http.get<ProductAssist>(`/assist/actives/${id}`, false);
}

export function fetchSets(assistID: number) {
  return http.get<{ list: AssistSet[] }>(`/assist/actives/${assistID}/sets`, false);
}

export function fetchSet(id: number) {
  return http.get<AssistSet>(`/assist/sets/${id}`, false);
}

export function fetchAssistSet(id: number) {
  return fetchSet(id);
}

export function fetchMyAssistSet(assistID: number) {
  return http.get<AssistSet | null>(`/assist/actives/${assistID}/mine`, true);
}

export function startAssist(id: number) {
  return http.post<AssistSet>(`/assist/actives/${id}/start`, {}, true);
}

export function helpAssist(setID: number) {
  return http.post<AssistSet>(`/assist/sets/${setID}/help`, {}, true);
}

export function assistCheck(data: {
  product_assist_set_id: number;
  cart_num?: number;
  address_id?: number;
}) {
  return http.post<{
    pay_price: number;
    price: number;
    store_name: string;
    product_assist_set_id: number;
  }>("/order/assist/check", data, true);
}

export function assistCreate(data: {
  product_assist_set_id: number;
  cart_num?: number;
  address_id: number;
}) {
  return http.post<{ group_order_id: number; pay_price: number }>(
    "/order/assist/create",
    data,
    true,
  );
}
