import { http } from "@/utils/request";

export interface ProductAssist {
  product_assist_id: number;
  product_id: number;
  assist_price: number;
  assist_count: number;
  store_name?: string;
  mer_name?: string;
  ot_price?: number;
}

export interface AssistSet {
  product_assist_set_id: number;
  product_assist_id: number;
  yet_assist_count: number;
  assist_count: number;
  status: number;
  nickname?: string;
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

export function startAssist(id: number) {
  return http.post<AssistSet>(`/assist/actives/${id}/start`, {});
}

export function helpAssist(setID: number) {
  return http.post<AssistSet>(`/assist/sets/${setID}/help`, {});
}

export function assistCreate(body: {
  product_assist_set_id: number;
  cart_num?: number;
  address_id: number;
}) {
  return http.post<{ group_order_id: number; pay_price: number }>(
    "/order/assist/create",
    body as Record<string, unknown>,
  );
}
