import { http } from "@/utils/request";

export interface ReservationProduct {
  product_id: number;
  store_name: string;
  mer_name?: string;
  price: number;
  image?: string;
}

export interface ReservationSlot {
  attr_reservation_id: number;
  start_time: string;
  end_time: string;
  stock: number;
  remain?: number;
}

export function fetchReservationProducts(page = 1, limit = 20) {
  return http.get<{ list: ReservationProduct[]; total: number }>(
    `/reservation/products?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchReservationSlots(productId: number, date: string) {
  return http.get<{ list: ReservationSlot[]; date: string }>(
    `/reservation/products/${productId}/slots?date=${encodeURIComponent(date)}`,
    false,
  );
}

export function createReservationOrder(body: {
  product_id: number;
  slot_id: number;
  date: string;
  address_id: number;
  mark?: string;
}) {
  return http.post<{ group_order_id: number }>("/order/reservation/create", body, true);
}
