import { requestClient } from '#/api/request';

import type { PaginatedList } from './product';

export interface LiveVodVideoItem {
  cover_url?: string;
  create_time_text?: string;
  duration?: number;
  file_id: string;
  media_name?: string;
  media_url?: string;
}

export interface LiveVodApplyUploadResult {
  cover_storage_path?: string;
  media_storage_path?: string;
  storage_bucket?: string;
  storage_region?: string;
  temp_certificate?: {
    expired_time?: number;
    secret_id?: string;
    secret_key?: string;
    token?: string;
  };
  vod_session_key?: string;
}

export interface LiveVodCommitUploadResult {
  file_id: string;
  media_url: string;
}

export async function getLiveVodListApi(params: {
  keyword?: string;
  list_rows?: number;
  page?: number;
}) {
  return requestClient.post<{ list: PaginatedList<LiveVodVideoItem> }>(
    '/api/v1/shop/live/vod/list',
    params,
  );
}

export async function deleteLiveVodApi(payload: { file_id: string }) {
  return requestClient.post<{ msg?: string }>(
    '/api/v1/shop/live/vod/delete',
    payload,
  );
}

export async function applyLiveVodUploadApi(payload: {
  cover_type?: string;
  media_name?: string;
  media_type?: string;
}) {
  return requestClient.post<LiveVodApplyUploadResult>(
    '/api/v1/shop/live/vod/apply-upload',
    payload,
  );
}

export async function commitLiveVodUploadApi(payload: { vod_session_key: string }) {
  return requestClient.post<LiveVodCommitUploadResult>(
    '/api/v1/shop/live/vod/commit-upload',
    payload,
  );
}
