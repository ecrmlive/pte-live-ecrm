import { http } from "@/utils/request";

export interface ReserveProduct {
  product_id: number;
  mer_id: number;
  mer_name?: string;
  store_name: string;
  image?: string;
  price: number;
  ot_price?: number;
  show_reservation_days?: number;
}

export interface SlotDay {
  attr_reservation_id: number;
  start_time: string;
  end_time: string;
  remain: number;
  label: string;
}

export function fetchReserveProducts(page = 1, limit = 20) {
  return http.get<{ list: ReserveProduct[]; total: number }>(
    `/reservation/products?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchDaySlots(productId: number, date: string) {
  return http.get<{ list: SlotDay[]; date: string }>(
    `/reservation/products/${productId}/slots?date=${encodeURIComponent(date)}`,
    false,
  );
}

export function reservationCreate(body: {
  product_id: number;
  slot_id: number;
  date: string;
  address_id: number;
}) {
  return http.post<{ group_order_id: number }>(
    "/order/reservation/create",
    body as unknown as Record<string, unknown>,
  );
}
