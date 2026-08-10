import { requestClient } from '#/api/request';

export interface PlatformAssistActive {
  action_status: number;
  all?: number;
  assist_count: number;
  assist_price: number;
  assist_status?: number;
  assist_status_text?: string;
  assist_user_count: number;
  create_time?: string;
  end_time: string;
  image?: string;
  is_show: number;
  is_trader?: number;
  mer_id: number;
  mer_name?: string;
  pay?: number;
  pay_count?: number;
  product_assist_id: number;
  product_id: number;
  product_status: number;
  product_status_name?: string;
  refusal?: string;
  start_time: string;
  status: number;
  stock: number;
  stock_count?: number;
  store_info?: string;
  store_name: string;
  success?: number;
  trader_name?: string;
}

export interface PlatformAssistPage {
  limit: number;
  list: PlatformAssistActive[];
  page: number;
  total: number;
}

/** 用户发起的助力实例（CRMEB StoreProductAssistSet） */
export interface PlatformAssistSet {
  assist_count: number;
  assist_price: number;
  assist_user_count: number;
  create_time: string;
  end_time?: string;
  helpers?: PlatformAssistHelper[];
  image?: string;
  is_trader?: number;
  mer_id: number;
  mer_name?: string;
  nickname?: string;
  product_assist_id: number;
  product_assist_set_id: number;
  product_id: number;
  start_time?: string;
  status: number;
  store_name?: string;
  trader_name?: string;
  uid: number;
  yet_assist_count: number;
}

export interface PlatformAssistHelper {
  avatar_img?: string;
  create_time?: string;
  nickname: string;
  product_assist_user_id: number;
  uid: number;
}

export interface PlatformAssistSetPage {
  limit: number;
  list: PlatformAssistSet[];
  page: number;
  total: number;
}

export function listPlatformAssistApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  status?: number;
  is_show?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformAssistPage>('/assist/actives', { params });
}

export function getPlatformAssistApi(id: number) {
  return requestClient.get<PlatformAssistActive>(`/assist/actives/${id}`);
}

export function updatePlatformAssistApi(
  id: number,
  payload: { is_show: number },
) {
  return requestClient.put(`/assist/actives/${id}`, payload);
}

export function listPlatformAssistSetsApi(params: {
  limit: number;
  page: number;
  mer_id?: number;
  keyword?: string;
  user_name?: string;
  is_trader?: number;
  status?: number;
  date_from?: string;
  date_to?: string;
}) {
  return requestClient.get<PlatformAssistSetPage>('/assist/sets', { params });
}

export function getPlatformAssistSetApi(id: number) {
  return requestClient.get<PlatformAssistSet>(`/assist/sets/${id}`);
}
