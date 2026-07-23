import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface LiveComplaintItem {
  complaint_id: number;
  content?: string;
  cover_img?: string;
  create_time_text?: string;
  handle_remark?: string;
  handle_time_text?: string;
  live_id: number;
  live_name?: string;
  reason_type?: string;
  session_id?: string;
  status: number;
  status_text?: string;
  user_avatar?: string;
  user_id: number;
  user_nick_name?: string;
}

export interface LiveComplaintListQuery {
  list_rows?: number;
  live_id?: number;
  page?: number;
  search?: string;
  status?: number;
}

export async function getLiveComplaintListApi(params: LiveComplaintListQuery) {
  return requestClient.post<{ list: PaginatedList<LiveComplaintItem> }>(
    '/api/v1/shop/live/complaint/list',
    params,
  );
}

export async function getLiveComplaintDetailApi(payload: {
  complaint_id: number;
}) {
  return requestClient.post<LiveComplaintItem>(
    '/api/v1/shop/live/complaint/detail',
    payload,
  );
}

export async function handleLiveComplaintApi(payload: {
  complaint_id: number;
  handle_remark?: string;
  status: number;
}) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/complaint/handle',
    payload,
  );
}
