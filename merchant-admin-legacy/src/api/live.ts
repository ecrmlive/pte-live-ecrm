import { http } from '@/api/http';

export interface LiveRoom {
  broadcast_room_id: number;
  mer_id: number;
  name: string;
  cover_img?: string;
  feeds_img?: string;
  anchor_name?: string;
  phone?: string;
  live_status: number;
  status: number;
  is_show: number;
  start_time?: string;
  end_time?: string;
  goods?: Array<{ product_id: number; store_name?: string; price?: number }>;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  limit: number;
}

export function fetchLiveRooms(params: Record<string, unknown>) {
  return http.get<PageResult<LiveRoom>>('/broadcast/rooms', { params });
}

export function fetchLiveRoom(id: number) {
  return http.get<LiveRoom>(`/broadcast/rooms/${id}`);
}

export function createLiveRoom(data: Record<string, unknown>) {
  return http.post<LiveRoom>('/broadcast/rooms', data);
}

export function updateLiveRoom(id: number, data: Record<string, unknown>) {
  return http.put<LiveRoom>(`/broadcast/rooms/${id}`, data);
}

export function putLiveRoomGoods(id: number, product_ids: number[]) {
  return http.put<LiveRoom>(`/broadcast/rooms/${id}/goods`, { product_ids });
}

export function deleteLiveRoom(id: number) {
  return http.delete<{ ok: boolean }>(`/broadcast/rooms/${id}`);
}

export function setLiveStatus(id: number, live_status: number) {
  return http.put<LiveRoom>(`/broadcast/rooms/${id}/live`, { live_status });
}
