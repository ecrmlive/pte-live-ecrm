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
  anchor_wechat?: string;
  broadcast_room_id: number;
  cover_img?: string;
  create_time: string;
  end_time?: string;
  feeds_img?: string;
  goods?: RoomGoods[];
  is_show: number;
  is_trader?: number;
  live_status: number;
  mark?: string;
  mer_id: number;
  mer_name?: string;
  name: string;
  play_url?: string;
  refusal?: string;
  sort: number;
  star: number;
  start_time?: string;
  status: number;
  trader_name?: string;
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

export interface RecommendRoomPayload {
  sort?: number;
  star?: number;
}

export function listRoomsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status_tag?: number;
  status?: number;
  show_type?: number;
  live_status?: number;
  star?: number;
  is_trader?: number;
}) {
  return requestClient.get<RoomPage>('/broadcast/rooms', { params });
}

export function getPlatformBroadcastApi(id: number) {
  return requestClient.get<Room>(`/broadcast/rooms/${id}`);
}

export function auditPlatformBroadcastApi(id: number, data: AuditRoomPayload) {
  return requestClient.put<Room>(`/broadcast/rooms/${id}/status`, data);
}

export function setPlatformBroadcastShowApi(id: number, isShow: 0 | 1) {
  return requestClient.put<Room>(`/broadcast/rooms/${id}/show`, {
    is_show: isShow,
  });
}

export function setPlatformBroadcastRecommendApi(
  id: number,
  data: RecommendRoomPayload,
) {
  return requestClient.put<Room>(`/broadcast/rooms/${id}/recommend`, data);
}

export function deletePlatformBroadcastApi(id: number) {
  return requestClient.delete<{ ok: boolean }>(`/broadcast/rooms/${id}`);
}
