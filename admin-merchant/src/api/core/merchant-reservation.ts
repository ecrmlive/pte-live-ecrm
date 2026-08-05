import { requestClient } from '#/api/request';

export interface ReservationProduct {
  image: string;
  mer_id: number;
  mer_name: string;
  ot_price: number;
  price: number;
  product_id: number;
  reservation_type: number;
  show_reservation_days: number;
  stock: number;
  store_name: string;
}

export interface ReservationSlot {
  attr_reservation_id: number;
  end_time: string;
  product_id: number;
  start_time: string;
  stock: number;
  unique: string;
  use_num: number;
}

export interface ReservationConfig {
  is_cancel_reservation: number;
  product_id: number;
  product_reservation_id: number;
  reservation_type: number;
  show_reservation_days: number;
  time_period: string;
}

export interface ReservationPage<T> {
  limit: number;
  list: T[];
  page: number;
  total: number;
}

export interface ReservationConfigSaveInput {
  reservation_type: number;
  show_reservation_days: number;
  slots: Array<{
    end_time: string;
    start_time: string;
    stock: number;
  }>;
}

export function listReservationProductsApi(params: { limit: number; page: number }) {
  return requestClient.get<ReservationPage<ReservationProduct>>('/reservation/products', { params });
}

export function getReservationConfigApi(productId: number) {
  return requestClient.get<{ config: ReservationConfig; slots: ReservationSlot[] }>(
    `/reservation/products/${productId}/config`,
  );
}

export function saveReservationConfigApi(productId: number, data: ReservationConfigSaveInput) {
  return requestClient.put<{ config: ReservationConfig; slots: ReservationSlot[] }>(
    `/reservation/products/${productId}/config`,
    data,
  );
}
