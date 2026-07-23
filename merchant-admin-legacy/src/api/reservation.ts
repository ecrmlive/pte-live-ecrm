import { http } from '@/api/http';

export interface ReservationProduct {
  product_id: number;
  store_name: string;
  price: number;
  stock: number;
}

export interface ReservationSlot {
  attr_reservation_id?: number;
  start_time: string;
  end_time: string;
  stock: number;
}

export interface ReservationConfigBody {
  config: {
    product_id: number;
    show_reservation_days: number;
    reservation_type: number;
  };
  slots: ReservationSlot[];
}

export function fetchReservationProducts(params: Record<string, unknown>) {
  return http.get<{ list: ReservationProduct[]; total: number }>('/reservation/products', { params });
}

export function fetchReservationConfig(id: number) {
  return http.get<ReservationConfigBody>(`/reservation/products/${id}/config`);
}

export function saveReservationConfig(id: number, data: Record<string, unknown>) {
  return http.put<ReservationConfigBody>(`/reservation/products/${id}/config`, data);
}
