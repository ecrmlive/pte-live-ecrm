import { http } from "@/utils/request";

export interface PresellActive {
  product_presell_id: number;
  product_id: number;
  store_name: string;
  mer_name?: string;
  price: number;
  down_price?: number;
  final_price?: number;
  presell_type?: number;
  ot_price?: number;
  stock: number;
  in_window?: boolean;
  start_time?: string;
  end_time?: string;
}

export function fetchPresells(page = 1, limit = 20) {
  return http.get<{ list: PresellActive[]; total: number }>(
    `/presell/actives?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchPresell(id: number) {
  return http.get<PresellActive>(`/presell/actives/${id}`, false);
}

export function presellCreate(body: {
  product_presell_id: number;
  cart_num: number;
  address_id: number;
}) {
  return http.post<{ group_order_id: number; pay_price: number }>(
    "/order/presell/create",
    body as unknown as Record<string, unknown>,
  );
}

export interface PresellFinal {
  presell_order_id: number;
  presell_order_sn: string;
  order_id: number;
  store_name?: string;
  pay_price: number;
  paid: number;
  status: number;
  final_end_time?: string;
}

export function fetchPresellFinals() {
  return http.get<{ list: PresellFinal[]; total: number }>("/presell/finals");
}

export function payPresellFinal(id: number, type: "balance" = "balance") {
  return http.post<PresellFinal>(`/presell/pay/${id}`, { type } as Record<string, unknown>);
}
