import { requestClient } from '#/api/request';

export interface PlatformAssistActive {
  action_status: number;
  assist_count: number;
  assist_price: number;
  assist_user_count: number;
  end_time: string;
  is_show: number;
  mer_id: number;
  mer_name?: string;
  product_assist_id: number;
  product_id: number;
  product_status: number;
  stock: number;
  store_name: string;
  start_time: string;
  status: number;
}

export interface PlatformAssistPage {
  limit: number;
  list: PlatformAssistActive[];
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
