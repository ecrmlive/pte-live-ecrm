import { http } from "@/utils/request";

export interface ReserveProduct {
  product_id: number;
  store_name: string;
  mer_name?: string;
  price: number;
  image?: string;
  ot_price?: number;
  show_reservation_days?: number;
  reservation_type?: number;
}

export interface SlotDay {
  attr_reservation_id: number;
  start_time: string;
  end_time: string;
  stock: number;
  remain: number;
  label: string;
}

/** 保留列表页使用的旧名称，避免同一领域出现两套 DTO。 */
export type ReservationProduct = ReserveProduct;
export type ReservationSlot = SlotDay;

export function fetchReserveProducts(page = 1, limit = 20) {
  return http.get<{ list: ReserveProduct[]; total: number }>(
    `/reservation/products?page=${page}&limit=${limit}`,
    false,
  );
}

export const fetchReservationProducts = fetchReserveProducts;

export function fetchReserveProduct(productId: number) {
  return http.get<ReserveProduct>(`/reservation/products/${productId}`, false);
}

export function fetchDaySlots(productId: number, date: string) {
  return http.get<{ list: SlotDay[]; date: string }>(
    `/reservation/products/${productId}/slots?date=${encodeURIComponent(date)}`,
    false,
  );
}

export const fetchReservationSlots = fetchDaySlots;

type ReservationOrderInput = {
  product_id: number;
  slot_id: number;
  date: string;
  address_id: number;
  mark?: string;
};

export function reservationCheck(body: ReservationOrderInput) {
  return http.post<{
    product_id: number;
    slot_id: number;
    date: string;
    pay_price: number;
    verify_hint: string;
  }>("/order/reservation/check", body, true);
}

export function reservationCreate(body: ReservationOrderInput) {
  return http.post<{ group_order_id: number }>("/order/reservation/create", body, true);
}

export const createReservationOrder = reservationCreate;
