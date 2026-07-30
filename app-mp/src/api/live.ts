import { http } from "@/utils/request";

export interface LiveGoods {
  product_id: number;
  store_name?: string;
  image?: string;
  price?: number;
  on_sale?: number;
}

export interface LiveRoom {
  broadcast_room_id: number;
  mer_id: number;
  name: string;
  cover_img?: string;
  feeds_img?: string;
  anchor_name?: string;
  live_status: number;
  status: number;
  mer_name?: string;
  start_time?: string;
  end_time?: string;
  goods?: LiveGoods[];
}

export function fetchLiveRooms(page = 1, limit = 20) {
  return http.get<{ list: LiveRoom[]; total: number }>(
    `/live/rooms?page=${page}&limit=${limit}`,
    false,
  );
}

export function fetchLiveRoom(id: number) {
  return http.get<LiveRoom>(`/live/rooms/${id}`, false);
}
