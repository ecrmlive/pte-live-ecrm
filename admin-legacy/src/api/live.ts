import { http } from '@/api/http';

export interface LiveRoom {
  broadcast_room_id: number;
  mer_id: number;
  mer_name?: string;
  name: string;
  anchor_name?: string;
  live_status: number;
  status: number;
  is_show: number;
  refusal?: string;
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

export function updateLiveRoomStatus(
  id: number,
  data: { status?: number; refusal?: string; is_show?: number },
) {
  return http.put<LiveRoom>(`/broadcast/rooms/${id}/status`, data);
}
