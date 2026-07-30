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
  image?: string;
  delivery_day?: number;
  final_start_time?: string;
  final_end_time?: string;
}

export interface PresellFinal {
  presell_order_id: number;
  presell_order_sn: string;
  order_id: number;
  product_presell_id: number;
  store_name?: string;
  pay_price: number;
  paid: number;
  status: number;
  final_start_time?: string;
  final_end_time?: string;
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

export function presellCheck(body: {
  product_presell_id: number;
  cart_num?: number;
  address_id?: number;
}) {
  return http.post<{
    pay_price: number;
    price: number;
    down_price?: number;
    final_price?: number;
    presell_type?: number;
    cart_num: number;
    store_name: string;
  }>("/order/presell/check", body);
}

export function presellCreate(body: {
  product_presell_id: number;
  cart_num?: number;
  address_id: number;
}) {
  return http.post<{ group_order_id: number; pay_price: number }>("/order/presell/create", body);
}

export function fetchPresellFinals(unpaid = true) {
  return http.get<{ list: PresellFinal[]; total: number }>(
    `/presell/finals?unpaid=${unpaid ? 1 : 0}`,
  );
}

export function fetchPresellFinal(id: number) {
  return http.get<PresellFinal>(`/presell/finals/${id}`);
}

export function payPresellFinal(id: number, type: "mock" | "balance" = "mock") {
  return http.post<PresellFinal>(`/presell/pay/${id}`, { type });
}
