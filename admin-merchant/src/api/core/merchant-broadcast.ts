import { requestClient } from '#/api/request';

export interface MerchantBroadcastGoods {
  image?: string;
  on_sale: number;
  price?: number;
  product_id: number;
  sort: number;
  store_name?: string;
}

export interface MerchantBroadcastRoom {
  anchor_name: string;
  broadcast_room_id: number;
  cover_img?: string;
  create_time: string;
  end_time?: string;
  feeds_img?: string;
  goods?: MerchantBroadcastGoods[];
  is_show: number;
  live_status: number;
  mark?: string;
  name: string;
  phone?: string;
  play_url?: string;
  push_url?: string;
  refusal?: string;
  sort: number;
  star: number;
  start_time?: string;
  status: number;
}

export interface MerchantBroadcastRoomInput {
  anchor_name: string;
  cover_img: string;
  end_time: string;
  feeds_img: string;
  is_show: number;
  mark: string;
  name: string;
  phone: string;
  play_url: string;
  product_ids?: number[];
  push_url: string;
  sort: number;
  star: number;
  start_time: string;
}

export interface MerchantBroadcastPage {
  limit: number;
  list: MerchantBroadcastRoom[];
  page: number;
  total: number;
}

export function listMerchantBroadcastRoomsApi(params: { limit: number; page: number }) {
  return requestClient.get<MerchantBroadcastPage>('/broadcast/rooms', { params });
}

export function getMerchantBroadcastRoomApi(id: number) {
  return requestClient.get<MerchantBroadcastRoom>(`/broadcast/rooms/${id}`);
}

export function createMerchantBroadcastRoomApi(data: MerchantBroadcastRoomInput) {
  return requestClient.post<MerchantBroadcastRoom>('/broadcast/rooms', data);
}

export function updateMerchantBroadcastRoomApi(id: number, data: MerchantBroadcastRoomInput) {
  return requestClient.put<MerchantBroadcastRoom>(`/broadcast/rooms/${id}`, data);
}

export function deleteMerchantBroadcastRoomApi(id: number) {
  return requestClient.delete(`/broadcast/rooms/${id}`);
}

export function setMerchantBroadcastLiveApi(id: number, liveStatus: number) {
  return requestClient.put<MerchantBroadcastRoom>(`/broadcast/rooms/${id}/live`, { live_status: liveStatus });
}

export function setMerchantBroadcastGoodsApi(id: number, productIDs: number[]) {
  return requestClient.put<MerchantBroadcastRoom>(`/broadcast/rooms/${id}/goods`, { product_ids: productIDs });
}
