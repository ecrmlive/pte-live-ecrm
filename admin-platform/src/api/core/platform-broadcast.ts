import { requestClient } from '#/api/request';

export interface RoomGoods {
  id: number;
  image?: string;
  on_sale: number;
  price?: number;
  product_id: number;
  sort: number;
  store_name?: string;
}

export interface Room {
  anchor_name: string;
  broadcast_room_id: number;
  cover_img?: string;
  create_time: string;
  feeds_img?: string;
  goods?: RoomGoods[];
  is_show: number;
  live_status: number;
  mark?: string;
  mer_id: number;
  mer_name?: string;
  name: string;
  play_url?: string;
  refusal?: string;
  start_time?: string;
  status: number;
}

export interface RoomPage {
  limit: number;
  list: Room[];
  page: number;
  total: number;
}

export interface AuditRoomPayload {
  is_show?: 0 | 1;
  refusal?: string;
  status: -1 | 0 | 2;
}

export function listRoomsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<RoomPage>('/broadcast/rooms', { params });
}

export function getPlatformBroadcastApi(id: number) {
  return requestClient.get<Room>(`/broadcast/rooms/${id}`);
}

export function auditPlatformBroadcastApi(id: number, data: AuditRoomPayload) {
  return requestClient.put<Room>(`/broadcast/rooms/${id}/status`, data);
}
